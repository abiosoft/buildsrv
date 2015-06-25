package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/caddyserver/buildsrv/features"
	"github.com/mholt/custombuild"
)

// Build represents a custom build job.
type Build struct {
	sync.Mutex
	DoneChan         chan struct{}
	OutputFile       string
	DownloadFilename string
	DownloadFile     string
	GoOS             string
	GoArch           string
	GoARM            string
	Features         features.Middlewares
	Hash             string
	Expires          time.Time
}

// buildHash creates a string that uniquely identifies a kind of build
func buildHash(goOS, goArch, goARM, orderedFeatures string) string {
	return fmt.Sprintf("%s:%s:%s:%s", goOS, goArch, goARM, orderedFeatures)
}

// build performs a build job. It's blocking, so run it in a goroutine.
// It writes errors to the standard log.
func build(job Build) error {
	builder, err := custombuild.New(caddyPath, codeGen, []string{ /*TODO*/ })
	defer builder.Teardown() // always perform cleanup
	if err != nil {
		return err
	}

	if job.GoArch == "arm" {
		armInt, err := strconv.Atoi(job.GoARM)
		if err != nil {
			return err
		}
		err = builder.BuildARM(job.GoOS, armInt, job.OutputFile)
	} else {
		err = builder.Build(job.GoOS, job.GoArch, job.OutputFile)
	}
	if err != nil {
		return err
	}

	// Create archive
	out, err := os.Create(job.DownloadFile)
	if err != nil {
		return err
	}

	w := zip.NewWriter(out)
	for _, fpath := range []string{
		filepath.Join(caddyPath, "/dist/README.txt"),
		filepath.Join(caddyPath, "/dist/LICENSES.txt"),
		filepath.Join(caddyPath, "/dist/CHANGES.txt"),
		job.OutputFile,
	} {
		fin, err := os.Open(fpath)
		if err != nil {
			return err
		}

		fout, err := w.Create(filepath.Base(fpath))
		if err != nil {
			w.Close()
			fin.Close()
			return err
		}

		_, err = io.Copy(fout, fin)
		if err != nil {
			w.Close()
			fin.Close()
			return err
		}

		fin.Close()
	}

	// Finish and close zip file
	err = w.Close()
	if err != nil {
		return err
	}

	// Delete uncompressed binary
	err = os.Remove(job.OutputFile)
	if err != nil {
		return err
	}

	// Build is ready
	close(job.DoneChan)
	job.Expires = time.Now().Add(buildExpiry)
	buildsMutex.Lock()
	builds[job.Hash] = job
	buildsMutex.Unlock()

	// Build expires after some time; run in goroutine
	// so this function can return and cleanup right away
	go func() {
		time.Sleep(buildExpiry)

		// Delete the job
		buildsMutex.Lock()
		delete(builds, job.Hash)
		buildsMutex.Unlock()

		// Delete file and its folder
		err := os.RemoveAll(filepath.Dir(job.DownloadFile))
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}
