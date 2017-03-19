package main

import (
	"flag"
	"net/http"
)

var dir = flag.String("d", "./client/public", "client location")

func main() {
	flag.Parse()

	http.ListenAndServe(":8080", http.FileServer(http.Dir(*dir)))
}
