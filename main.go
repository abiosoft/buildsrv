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
	"time"

	"github.com/caddyserver/buildsrv/server"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// Get GOPATH and path of caddy project
	cmd := exec.Command("go", "env", "GOPATH")
	result, err := cmd.Output()
	if err != nil {
		log.Fatal("Cannot locate GOPATH:", err)
	}
	server.CaddyPath = strings.TrimSpace(string(result)) + "/src/" + server.MainCaddyPackage
}

func main() {
	go func() {
		// Delete existing builds on quit
		interrupt := make(chan os.Signal)
		signal.Notify(interrupt, os.Interrupt, os.Kill)
		<-interrupt

		err := os.RemoveAll(server.BuildPath)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	http.HandleFunc("/download/build", server.BuildHandler)
	http.Handle("/download/builds/", http.StripPrefix("/download/builds/", http.FileServer(http.Dir(server.BuildPath))))
	http.HandleFunc("/online", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	fmt.Println("Example URL:")
	fmt.Println("http://localhost:5050/download/build?os=darwin&arch=amd64&features=markdown,git,templates")
	http.ListenAndServe(":5050", nil)
}
