package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	SERVER = "0.0.0.0:9000"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Welcome Home!")
	case "POST":
		fmt.Fprintf(w, "Thanks for the POST!")
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("[+] Container built successfully")
		return
	}
	http.HandleFunc("/", HomePage)
	log.Println("[+] Hope the pipeline works!!")
	log.Printf("Listening: %s\n", SERVER)
	http.ListenAndServe(SERVER, nil)
}
