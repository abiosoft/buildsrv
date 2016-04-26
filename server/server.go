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
	// See https://golang.org/doc/install/source#environment
	// Commented builds are problematic.
	allowed = combos{
		{os: "darwin", arch: "386"},
		{os: "darwin", arch: "amd64"},
		{os: "darwin", arch: "arm"},
		//{os: "darwin", arch: "arm64"},
		//{os: "dragonfly", arch: "amd64"},
		{os: "freebsd", arch: "386"},
		{os: "freebsd", arch: "amd64"},
		{os: "freebsd", arch: "arm"},
		{os: "linux", arch: "386"},
		{os: "linux", arch: "amd64"},
		{os: "linux", arch: "arm"},
		{os: "linux", arch: "arm64"},
		{os: "linux", arch: "ppc64"},
		{os: "linux", arch: "ppc64le"},
		{os: "linux", arch: "mips64"},
		//{os: "linux", arch: "mips64le"},
		{os: "netbsd", arch: "386"},
		{os: "netbsd", arch: "amd64"},
		{os: "netbsd", arch: "arm"},
		{os: "openbsd", arch: "386"},
		{os: "openbsd", arch: "amd64"},
		{os: "openbsd", arch: "arm"},
		//{os: "plan9", arch: "386"},
		//{os: "plan9", arch: "amd64"},
		{os: "solaris", arch: "amd64"},
		{os: "windows", arch: "386"},
		{os: "windows", arch: "amd64"},
	}

	allowedARM = list{"5", "6", "7"}
	defaultARM = 7

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

type list []string

func (l list) contains(target string) bool {
	for _, s := range l {
		if s == target {
			return true
		}
	}
	return false
}

type combos []platform

func (c combos) valid(os, arch string) bool {
	for _, pl := range c {
		if pl.os == os && pl.arch == arch {
			return true
		}
	}
	return false
}

type platform struct {
	os, arch string
}
