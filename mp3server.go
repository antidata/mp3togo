package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

type MusicFile struct {
	Name      string
	Mode      os.FileMode
	Directory bool
}

type ByFileName []MusicFile

func (a ByFileName) Len() int           { return len(a) }
func (a ByFileName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFileName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type Configuration struct {
	Mp3Root string
	Prefix  string
}

var configuration = Configuration{}

func init() {
	configuration = Configuration{"./mp3", "/mp3/"}
}

func servePath(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Serving Path " + path)
		http.ServeFile(w, r, path)
	}
}

func failIfError(e error) {
	if e != nil {
		panic(e)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	path := filepath.Join(configuration.Mp3Root, r.URL.Path[len(configuration.Prefix):])
	fmt.Println("Path: " + path)
	stat, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	switch stat.IsDir() {
	case true:
		serveDirectory(w, r, path)
	case false:
		http.ServeFile(w, r, path)
	}
}

func serveDirectory(w http.ResponseWriter, r *http.Request, path string) {
	file, err := os.Open(path)

	failIfError(err)

	defer file.Close()

	directoryFiles, err := file.Readdir(-1)

	failIfError(err)

	fileList := make([]MusicFile, len(directoryFiles), len(directoryFiles))

	for i := range directoryFiles {
		fileList[i].Name = directoryFiles[i].Name()
		fileList[i].Directory = directoryFiles[i].IsDir()
		fileList[i].Mode = directoryFiles[i].Mode()
	}

	jsonEncoder := json.NewEncoder(w)

	// sorting by filename
	sort.Sort(ByFileName(fileList))

	if jsonEncoder.Encode(&fileList) != nil {
		panic(err)
	}

	defer func() {
		err, ok := recover().(error)
		if ok {
			fmt.Println("Error serving: " + path + " message: " + err.Error())
			http.Error(w, "Error loading "+path, http.StatusInternalServerError)
		}
	}()
}

