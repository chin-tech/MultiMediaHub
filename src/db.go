package main

import (
	"database/sql"
	"fmt"
	"log"
	// "net/http"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"time"
)

func GetDb(database_type string) *sql.DB {
	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("[-] DATABASE ERROR: Please enter a port between [5000 - 65535] | \t[-] Port provided: [%s]\n", os.Getenv("DB_PORT"))
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	con := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, db_name)

	db, err := sql.Open(database_type, con)
	if err != nil {
		log.Fatalf("[-] Invalid database arguments: %v\n", err)
	}

	waitTime := 60
	for ; waitTime > 0; waitTime-- {
		err = db.Ping()
		if err != nil {
			fmt.Printf("[-] Connection could not be established: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	if waitTime == 0 {
		log.Fatalf("[-] Database connection could not be established after 15 seconds: %v \n", err)
	}
	fmt.Println("[+] Database successfully connected")
	return db
}

func GetWord(db *sql.DB, word string) (string, error) {
	rows, err := db.Query(`SELECT "definition" FROM "definitions" WHERE word=$1`, word)
	fmt.Printf("[DEBUG] - Searching for word: %s\n", word)
	if err != nil {
		fmt.Printf("[-] Error in query: %v\n", err)
		return "", err
	}

	defer rows.Close()
	for rows.Next() {
		var definition string
		err = rows.Scan(&definition)
		if err != nil {
			fmt.Printf("[-] Unknown error: %v\n", err)
			return "", err
		}
		return definition, nil
	}

	return "", nil

}

func InsertWord(db *sql.DB, word, definition string) error {
	insStmt := `INSERT into "definitions"("word", "definition") values($1, $2)`
	_, err := db.Exec(insStmt, word, definition)
	if err != nil {
		return fmt.Errorf("[-] Failed to insert: %v\n", err)
	}

	return nil

}

func UpdateWord(db *sql.DB, word, definition string) error {
	insStmt := `UPDATE "definitions" set "definition"=$2 where "word"=$1`
	_, err := db.Exec(insStmt, word, definition)
	if err != nil {
		return fmt.Errorf("[-] Failed to insert: %v\n", err)
	}

	return nil

}
