package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("try to listend 8080 port...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s", r.RemoteAddr)
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal("Hostname() error")
		os.Exit(1)
	}
	content := "Hello Kubernetes Beginners! Server in " + hostName
	fmt.Fprintln(w, content)
}
