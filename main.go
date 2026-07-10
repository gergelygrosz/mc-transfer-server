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
	if err != nil {
		fmt.Println("failed to open config.json")
		panic(err)
	}

	err = json.Unmarshal(configFile, &cfg)
	if err != nil {
		fmt.Println("failed to parse config.json")
		panic(err)
	}

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

	err := filepath.WalkDir(cfg.Source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			srcFile, err := os.Open(path)
			if err != nil {
				fmt.Println("failed to open source file")
				return err
			}

			relativePath, err := filepath.Rel(cfg.Source, path)
			if err != nil {
				fmt.Println("failed to relativise file path in respect to source root")
				return err
			}

			destFile, err := archive.CreateHeader(&zip.FileHeader{
				Name:   filepath.ToSlash(relativePath),
				Method: zip.Store,
			})
			if err != nil {
				fmt.Println("failed to create file inside of the archive")
				return err
			}

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				fmt.Println("failed to copy file contents into the archive")
				return err
			}

			srcFile.Close()
		}

		return nil
	})
	if err != nil {
		fmt.Println("archive creation failed as described above")
		panic(err)
	}
}
