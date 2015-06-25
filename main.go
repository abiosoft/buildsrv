package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/caddyserver/buildsrv/features"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	go func() {
		// Delete existing builds on quit
		interrupt := make(chan os.Signal)
		signal.Notify(interrupt, os.Interrupt, os.Kill)
		<-interrupt

		err := os.RemoveAll(buildPath)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	http.HandleFunc("/download/build", handleBuild)
	http.Handle("/download/builds/", http.StripPrefix("/download/builds/", http.FileServer(http.Dir(buildPath))))
	http.HandleFunc("/online", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	fmt.Println("Example URL:")
	fmt.Println("http://localhost:5050/download/build?os=darwin&arch=amd64&features=markdown,git,templates")
	http.ListenAndServe(":5050", nil)
}

func handleError(w http.ResponseWriter, r *http.Request, err error, status int) {
	if status >= 500 {
		log.Printf("[%d %s] %v", status, r.URL.String(), err)
		http.Error(w, http.StatusText(status), status)
	} else {
		http.Error(w, err.Error(), status)
	}
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

// buildHash creates a string that uniquely identifies a kind of build
func buildHash(goOS, goArch, goARM, orderedFeatures string) string {
	return fmt.Sprintf("%s:%s:%s:%s", goOS, goArch, goARM, orderedFeatures)
}

// codeGen is the function that mutates a copy of the project so that
// a custom build can be performed.
// TODO - @abiosoft
func codeGen(repo string, packages []string) error {
	// TODO
	return nil
}

// list is just any list of strings that can determine if
// a string belongs to it.
type list []string

// contains determines if target is in l.
func (l list) contains(target string) bool {
	for _, str := range l {
		if str == target {
			return true
		}
	}
	return false
}

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

var (
	allowedOS   = list{"linux", "darwin", "windows", "freebsd", "openbsd"}
	allowedArch = list{"386", "amd64", "arm"}
	allowedARM  = list{"5", "6", "7"}

	builds      = make(map[string]Build)
	buildsMutex sync.Mutex // protects the builds map
)

const (
	buildPath   = "builds"
	caddyPath   = "/Users/matt/Dev/src/github.com/mholt/caddy" // TODO: Get from env
	buildExpiry = 24 * time.Hour
)
