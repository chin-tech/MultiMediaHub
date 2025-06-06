package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	DEFAULT_PORT = "80"
)

var dbGlobal *sql.DB

func DynamicHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		HomePage(w, r)
		return
	}
	if dbGlobal == nil {
		http.Error(w, "[-] Database not yet available!", http.StatusInternalServerError)
		return
	}
	word := path[1:]
	// Insert extra container for handling this because we want to emulate microservices
	switch r.Method {
	case "GET":
		definition, err := GetWord(dbGlobal, word)
		if err != nil {
			fmt.Fprintf(w, "[-] Error! %v", err)
			return
		}
		if definition == "" {
			fmt.Fprintf(w, "This word doesn't yet have a definition friend, contribute!")
		}
		fmt.Fprintf(w, "%s", definition)

	case "POST":
		postBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "[-] Failed to read request body", http.StatusBadRequest)
			return
		}
		user_definition := strings.TrimSpace(string(postBody))
		err = InsertWord(dbGlobal, word, user_definition)
		if err != nil {
			fmt.Fprintf(w, "[-] Error! %v", err)
			return
		}
		fmt.Fprintf(w, "[+] Success! Added %s : %s", word, user_definition)

	case "PUT":
		postBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "[-] Failed to read request body", http.StatusBadRequest)
			return
		}
		user_definition := strings.TrimSpace(string(postBody))
		err = UpdateWord(dbGlobal, word, user_definition)
		if err != nil {
			fmt.Fprintf(w, "[-] Error! %v", err)
			return
		}
		fmt.Fprintf(w, "[+] Success! Updated %s : %s", word, user_definition)

	default:
		http.Error(w, "Leeroy Jenkins!", http.StatusMethodNotAllowed)

	}

}

func HomePage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Congratulations you're here! Visit the /<word> to get a word")
	case "POST":
		fmt.Fprintf(w, "Not here! post to a /<word>")
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("[+] Container built successfully")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}
	dbGlobal = GetDb("postgres")
	defer dbGlobal.Close()
	http.HandleFunc("/", DynamicHandler)
	log.Println("[+] Hope the pipeline works!!")
	server := fmt.Sprintf(":%s", port)
	log.Printf("Listening: %s\n", server)
	http.ListenAndServe(server, nil)
}
