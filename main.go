package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

var (
	Version, Build string
)

func main() {
	version := flag.Bool("version", false, "build version")
	host := flag.String("host", "0.0.0.0:9981", "listen address")
	flag.Parse()
	if *version {
		fmt.Printf("Version: %s Build: %s\nGo Version: %s\nGo OS/ARCH: %s %s\n", Version, Build, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return
	}
	s, err := newServer()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	log.Printf("Listening on :%v ...", *host)
	log.Fatalf("Error listening on %v: %v", *host, http.ListenAndServe(*host, s))
}
