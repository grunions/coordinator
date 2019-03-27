package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// database driver
	_ "github.com/lib/pq"
)

type handler struct {
	db *sql.DB
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World, %s\n", r.RemoteAddr)

	var name string
	rows, err := h.db.Query("SELECT name FROM count")
	if err != nil {
		fmt.Fprintf(w, "Error listing DB entries: %s\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Entry: %s\n", name)
	}
}

func main() {

	port := os.Getenv("PORT")
	dbrn := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbrn)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":"+port, &handler{db: db})
}
