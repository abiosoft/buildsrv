package server

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caddyserver/buildsrv/features"
)

// BuildHandler is the endpoint which creates and/or responds with builds.
func BuildHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Expose-Headers", "Location")

	goOS := r.URL.Query().Get("os")
	goArch := r.URL.Query().Get("arch")
	goARM := r.URL.Query().Get("arm")
	featureList := strings.Split(r.URL.Query().Get("features"), ",")
	if len(featureList) == 1 && featureList[0] == "" {
		featureList = []string{}
	}

	err := checkInput(goOS, goArch, goARM, featureList)
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}

	// Keep build hashes consistent with varying input
	if goArch != "arm" {
		goARM = ""
	}

	// Put features in order to keep hashes consistent and for use in the codegen function
	orderedFeatures := sortFeatures(featureList)

	// Create 'hash' to identify this build
	hash := buildHash(goOS, goArch, goARM, orderedFeatures.String())

	// Get the path from which to download the file
	buildsMutex.Lock()
	b, ok := builds[hash]
	buildsMutex.Unlock()

	if ok {
		// build exists; wait for it to complete if not done yet
		<-b.DoneChan
	} else {
		// no build yet; reserve it so we don't duplicate the build job
		ts := time.Now().Format("060201150405") // YearMonthDayHourMinSec
		var downloadPath string
		for {
			// find a suitable random number not already in use
			random := strconv.Itoa(rand.Intn(100) + 899)
			downloadPath = BuildPath + "/" + ts + random
			_, err := os.Stat(downloadPath)
			if os.IsNotExist(err) {
				break
			}
		}

		// Determine the remaining build information and reserve the build job

		downloadFileCompression := CompressZip
		if goOS == "linux" || goOS == "freebsd" || goOS == "openbsd" {
			downloadFileCompression = CompressTarGz
		}

		buildFilename := "caddy"
		if goOS == "windows" {
			buildFilename += ".exe"
		}

		downloadFilename := "caddy_" + goOS + "_" + goArch + "_custom"
		if downloadFileCompression == CompressZip {
			downloadFilename += ".zip"
		} else {
			downloadFilename += ".tar.gz"
		}

		b = &Build{
			DoneChan:                make(chan struct{}),
			OutputFile:              downloadPath + "/" + buildFilename,
			DownloadFile:            downloadPath + "/" + downloadFilename,
			DownloadFilename:        downloadFilename,
			DownloadFileCompression: downloadFileCompression,
			GoOS:     goOS,
			GoArch:   goArch,
			GoARM:    goARM,
			Features: orderedFeatures,
			Hash:     hash,
		}

		// Save the build, indicating currently in progress
		buildsMutex.Lock()
		builds[hash] = b
		buildsMutex.Unlock()

		// Perform build (blocking)
		err = b.Build()
		if err != nil {
			handleError(w, r, err, http.StatusInternalServerError)
			deleteBuildJob(hash) // delete the build; it didn't succeed
			return
		}
	}

	// Update our copy of the build information
	buildsMutex.Lock()
	b, ok = builds[hash]
	buildsMutex.Unlock()
	if !ok {
		handleError(w, r, errors.New("Build doesn't exist"), http.StatusInternalServerError)
		return
	}

	// Open download file
	f, err := os.Open(b.DownloadFile)
	if err != nil {
		handleError(w, r, err, http.StatusInternalServerError)
		deleteBuildJob(hash)
		return
	}
	defer f.Close()

	w.Header().Set("Location", "/download/"+b.DownloadFile)
	w.Header().Set("Expires", b.Expires.Format(http.TimeFormat))
	w.Header().Set("Content-Disposition", "attachment; filename=\""+b.DownloadFilename+"\"")

	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if r.Method == "GET" {
		io.Copy(w, f)
	}
}

// deleteBuild deletes a build from the map.
// It is safe for concurrent use. It does NOT
// delete the build from the file system.
func deleteBuildJob(hash string) {
	buildsMutex.Lock()
	delete(builds, hash)
	buildsMutex.Unlock()
}

// checkInput checks the arguments for valid values and returns an error
// if any one of them is invalid.
func checkInput(goOS, goArch, goARM string, featureList []string) error {
	// Check for required fields
	if goOS == "" {
		return errors.New("missing os parameter")
	}
	if goArch == "" {
		return errors.New("missing arch parameter")
	}

	// Check for valid input
	if !allowedOS.contains(goOS) {
		return errors.New("os not supported")
	}
	if !allowedArch.contains(goArch) {
		return errors.New("arch not supported")
	}
	if goARM != "" && !allowedARM.contains(goARM) {
		return errors.New("arm version not supported")
	}

	// Check features
	for _, feature := range featureList {
		if !features.Registry.Contains(feature) {
			return errors.New("unknown feature '" + feature + "'")
		}
	}

	return nil
}

// sortFeatures sorts features to the order in which they are registered.
func sortFeatures(featureList []string) features.Middlewares {
	var orderedFeatures features.Middlewares // TODO - could this be a []string instead? Would make things a little simpler, not needing that String() method
loop:
	for _, m := range features.Registry {
		for _, feature := range featureList {
			if feature == m.Directive {
				orderedFeatures = append(orderedFeatures, m)
				continue loop
			}
		}
	}
	return orderedFeatures
}

const (
	CompressZip = iota
	CompressTarGz
)
