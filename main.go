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

// const (
// 	ARCHIVE_NAME string = "pack.zip"
// 	SOURCE_PATH  string = "PrismLauncher/instances/Zéta/minecraft/blueprints/zeta"
// )

func download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"név.zip\"")

	archive := zip.NewWriter(w)
	defer archive.Close()

	filepath.WalkDir("folder", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			file, e := os.Open(path)
			check(e)

			secondaryWriter, e := archive.CreateHeader(&zip.FileHeader{
				Name:   path,
				Method: zip.Store,
			})
			check(e)

			_, e = io.Copy(secondaryWriter, file)
			check(e)
			file.Close()
		}

		return nil
	})
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
