package server

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zip"
)

// Zip creates a .zip file in the location zipPath containing
// the contents of files listed in filePaths.
func Zip(zipPath string, filePaths []string) error {
	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()

	w := zip.NewWriter(out)
	for _, fpath := range filePaths {
		infile, err := os.Open(fpath)
		if err != nil {
			return err
		}

		outfile, err := w.Create(filepath.Base(fpath))
		if err != nil {
			w.Close()
			infile.Close()
			return err
		}

		_, err = io.Copy(outfile, infile)
		if err != nil {
			w.Close()
			infile.Close()
			return err
		}

		infile.Close()
	}

	return w.Close()
}

// TarGz creates a .tar.gz file at targzPath containing
// the contents of files listed in filePaths.
func TarGz(targzPath string, filePaths []string) error {
	out, err := os.Create(targzPath)
	if err != nil {
		return err
	}
	defer out.Close()

	gzWriter := gzip.NewWriter(out)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	for _, fpath := range filePaths {
		infile, err := os.Open(fpath)
		if err != nil {
			return err
		}

		infileInfo, err := infile.Stat()
		if err != nil {
			infile.Close()
			return err
		}

		fileHeader, err := tar.FileInfoHeader(infileInfo, fpath)
		if err != nil {
			infile.Close()
			return err
		}

		err = tarWriter.WriteHeader(fileHeader)
		if err != nil {
			infile.Close()
			return err
		}

		_, err = io.Copy(tarWriter, infile)
		if err != nil {
			infile.Close()
			return err
		}

		infile.Close()
	}

	return nil
}
