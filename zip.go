package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

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
