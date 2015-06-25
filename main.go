package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"
)

const (
	// Path to the builds. The directory is fully managed, so just
	// choose one that is solely for builds; it may get deleted.
	buildPath = "builds"

	// How long builds live before being deleted
	buildExpiry = 24 * time.Hour

	// Canonical package name of Caddy's main
	mainCaddyPackage = "github.com/mholt/caddy"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// Get GOPATH and path of caddy project
	cmd := exec.Command("go", "env", "GOPATH")
	result, err := cmd.Output()
	if err != nil {
		log.Fatal("Cannot locate GOPATH:", err)
	}
	caddyPath = strings.TrimSpace(string(result)) + "/src/" + mainCaddyPackage
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

	http.HandleFunc("/download/build", buildHandler)
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

var (
	allowedOS   = list{"linux", "darwin", "windows", "freebsd", "openbsd"}
	allowedArch = list{"386", "amd64", "arm"}
	allowedARM  = list{"5", "6", "7"}

	builds      = make(map[string]*Build)
	buildsMutex sync.Mutex // protects the builds map

	// Path to the caddy project repository
	caddyPath string
)
