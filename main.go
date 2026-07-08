package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Source  string `json:"source"`
	Port    int    `json:"port"`
	UrlPath string `json:"urlpath"`
}

var cfg Config

func main() {
	configFile, err := os.ReadFile("config.json")
	check(err)

	err = json.Unmarshal(configFile, &cfg)
	check(err)

	fmt.Println("Listening to :8080")
	http.HandleFunc(cfg.UrlPath, download)
	port := ":" + strconv.Itoa(cfg.Port)
	http.ListenAndServe(port, nil)
}

func download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"pack.zip\"")

	archive := zip.NewWriter(w)
	defer archive.Close()

	filepath.WalkDir(cfg.Source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			srcFile, err := os.Open(path)
			check(err)

			relativePath, err := filepath.Rel(cfg.Source, path)
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
