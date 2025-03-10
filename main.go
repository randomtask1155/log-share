package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = *sessions.NewCookieStore([]byte("scripts-log-share-key"))
var username string
var password string
var directory string
var disableSSL = true
var duration string
var sleepDuration time.Duration
var cert string
var key string
var listenPort string

func exitProgram() {
	log.Println("server will shutdown in ", duration)
	time.Sleep(sleepDuration)
	log.Println("server duration has expired")
	os.Exit(0)
}

func main() {

	flag.StringVar(&username, "u", "admin", "define username defaults to admin")
	flag.StringVar(&password, "p", "", "define password or skip and program will generate one")
	flag.StringVar(&directory, "d", ".", "define a folder to share defaults to current working directory")
	flag.StringVar(&duration, "t", "24h", "How long to run before exiting: Valid time units are “ns”, “us” (or “µs”), “ms”, “s”, “m”, “h”.")

	flag.StringVar(&listenPort, "port", "0", "set a static listening port")
	flag.StringVar(&cert, "cert", "", "server cert")
	flag.StringVar(&key, "key", "", "server key")

	flag.Parse()
	if password == "" {
		generatePassword()
	}

	fi, err := os.Stat(directory)
	if err != nil {
		log.Fatal(err)
	}
	if !fi.IsDir() {
		log.Fatal("-d %s is not a directory", directory)
	}

	if cert != "" && key != "" {
		disableSSL = false
	}

	sleepDuration, err = time.ParseDuration(duration)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(basicAuthMiddleware)
	r.HandleFunc("/", serveRoot)
	r.HandleFunc("/path", servePath)
	//r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(directory))))

	listener, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		panic(err)
	}
	log.Println("staring server: ", getURLString(fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)))

	go exitProgram()
	if disableSSL {
		log.Fatal(http.Serve(listener, r))
	} else {
		log.Fatal(http.ServeTLS(listener, r, cert, key))
	}

}
