package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/caddyserver/buildsrv/features"
	"github.com/caddyserver/caddydev/caddybuild"
)

// Build represents a custom build job.
type Build struct {
	DoneChan                chan struct{}
	OutputFile              string
	DownloadFilename        string
	DownloadFileCompression int
	DownloadFile            string
	GoOS                    string
	GoArch                  string
	GoARM                   string
	Features                features.Middlewares
	Hash                    string
	Expires                 time.Time
	finished                bool
}

// Build performs a build job. This function is blocking. If the build
// job succeeds, it will automatically delete itself when it expires.
// If it fails, resources are not automatically cleaned up.
func (b *Build) Build() error {
	// Prepare the build
	builder, err := caddybuild.PrepareBuild(b.Features, false) // TODO: PullLatest (go get -u) DISABLED for stability; updates are manual for now
	defer builder.Teardown()                                   // always perform cleanup
	if err != nil {
		return err
	}

	// Perform the build
	if b.GoArch == "arm" {
		var armInt int
		if b.GoARM != "" {
			armInt, err = strconv.Atoi(b.GoARM)
			if err != nil {
				return err
			}
		} else {
			armInt = 7
		}
		err = builder.BuildStaticARM(b.GoOS, armInt, b.OutputFile)
	} else if b.GoOS == "darwin" { // At time of writing, building with CGO_ENABLED=0 for darwin can break stuff: https://www.reddit.com/r/golang/comments/46bd5h/ama_we_are_the_go_contributors_ask_us_anything/d03rmc9
		err = builder.Build(b.GoOS, b.GoArch, b.OutputFile)
	} else {
		err = builder.BuildStatic(b.GoOS, b.GoArch, b.OutputFile)
	}
	if err != nil {
		return err
	}

	// File list to include with build, then compress the build
	fileList := []string{
		filepath.Join(CaddyPath, "/dist/README.txt"),
		filepath.Join(CaddyPath, "/dist/LICENSES.txt"),
		filepath.Join(CaddyPath, "/dist/CHANGES.txt"),
		b.OutputFile,
	}
	if b.DownloadFileCompression == CompressZip {
		err = Zip(b.DownloadFile, fileList)
	} else if b.DownloadFileCompression == CompressTarGz {
		err = TarGz(b.DownloadFile, fileList)
	} else {
		return fmt.Errorf("unknown compress type %v", b.DownloadFileCompression)
	}
	if err != nil {
		return err
	}

	// Delete uncompressed binary
	err = os.Remove(b.OutputFile)
	if err != nil {
		return err
	}

	// Finalize the build and have it clean itself
	// up after its expiration
	b.finish()

	return nil
}

// finish finishes a job. Call this after the job is
// done with its build process and the result is ready
// for use. When this method is called, its lifetime
// begins and the build will be deleted after the
// expiration time.
func (b *Build) finish() {
	if b.finished {
		return
	}

	// Notify anyone waiting for the job to finish that it's done
	close(b.DoneChan)

	// Save the build in the master list
	buildsMutex.Lock()
	builds[b.Hash] = b
	buildsMutex.Unlock()

	// Make this idempotent
	b.finished = true

	if BuildExpiry > 0 {
		// Build lifetime starts now
		b.Expires = time.Now().Add(BuildExpiry)

		// Delete build after expiration time
		go func() {
			time.Sleep(BuildExpiry)

			// Delete the job
			deleteBuildJob(b.Hash)

			// Delete file and its folder
			err := os.RemoveAll(filepath.Dir(b.DownloadFile))
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

// buildHash creates a string that uniquely identifies a kind of build
func buildHash(goOS, goArch, goARM, orderedFeatures string) string {
	return fmt.Sprintf("%s:%s:%s:%s", goOS, goArch, goARM, orderedFeatures)
}
