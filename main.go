package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Listening to :8080")
	http.HandleFunc("/download", download)
	http.ListenAndServe(":8080", nil)
}

var sourcePath string = os.Getenv("mc-transfer-source")

func download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"pack.zip\"")

	archive := zip.NewWriter(w)
	defer archive.Close()

	filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			srcFile, err := os.Open(path)
			check(err)

			relativePath, err := filepath.Rel(sourcePath, path)
			check(err)

			destFile, err := archive.CreateHeader(&zip.FileHeader{
				Name:   filepath.ToSlash(relativePath),
				Method: zip.Store,
			})
			check(err)

			_, err = io.Copy(destFile, srcFile)
			check(err)
			srcFile.Close()
		}

		return nil
	})
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
