package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	host = flag.String("h", "0.0.0.0", "Listen hostname")
	port = flag.String("p", "5000", "Listen port")
	dir  = flag.String("d", ".", "Path to files")
)

func init() {
	if os.Getenv("PORT") != "" {
		*port = os.Getenv("PORT")
	}
}

func main() {
	flag.Parse()

	bind := *host + ":" + *port

	log.Println("Serving files from", *dir, "on", bind)
	log.Fatal(http.ListenAndServe(bind, http.FileServer(http.Dir(*dir))))
}
