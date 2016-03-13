package server

import (
	"log"
	"net/http"
	"sync"
)

const (
	// BuildPath is the path to the builds. The directory is fully
	// managed, so just choose one that is solely for builds; it
	// may get deleted.
	BuildPath = "builds"

	// BuildExpiry is how long builds live before being deleted.
	// The build server used to update its own dependencies by
	// running go get -u before builds, but we found out this was
	// a bad idea, so now I just manually update dependencies by
	// deleting the package folder and running go get (without -u).
	// Since the build server has to be taken offline to perform
	// these updates anyway, we don't have a need to expire the
	// builds anymore, as long the builds folder is empty when
	// the build server is restarted. So this used to be 24 hours,
	// but now is 0. This is especially useful since Go 1.5 and 1.6
	// got longer build times.
	BuildExpiry = 0

	// MainCaddyPackage is the canonical package name of Caddy's main.
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
