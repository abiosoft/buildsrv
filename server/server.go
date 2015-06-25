package server

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	// Path to the builds. The directory is fully managed, so just
	// choose one that is solely for builds; it may get deleted.
	BuildPath = "builds"

	// How long builds live before being deleted
	BuildExpiry = 24 * time.Hour

	// Canonical package name of Caddy's main
	MainCaddyPackage = "github.com/mholt/caddy"
)

var (
	allowedOS   = list{"linux", "darwin", "windows", "freebsd", "openbsd"}
	allowedArch = list{"386", "amd64", "arm"}
	allowedARM  = list{"5", "6", "7"}

	builds      = make(map[string]*Build)
	buildsMutex sync.Mutex // protects the builds map

	// Path to the caddy project repository
	CaddyPath string
)

func handleError(w http.ResponseWriter, r *http.Request, err error, status int) {
	if status >= 500 {
		log.Printf("[%d %s] %v", status, r.URL.String(), err)
		http.Error(w, http.StatusText(status), status)
	} else {
		http.Error(w, err.Error(), status)
	}
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
