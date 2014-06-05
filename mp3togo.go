package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", servePath("html/index.html"))
	http.HandleFunc("/js/mp3player.js", servePath("js/mp3player.js"))
	http.HandleFunc(configuration.Prefix, serveFile)
	http.ListenAndServe(":8080", nil)
}
