package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	// database driver
	_ "github.com/lib/pq"
)

type app struct {
	db *sql.DB
}

func (a *app) Healthcheck(w http.ResponseWriter, r *http.Request) {
	var success bool
	err := a.db.QueryRow("SELECT 1 = 1;").Scan(&success)

	if err != nil {
		fmt.Fprint(w, "Database error\n")
		return
	}

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Wrong database response\n")
		return
	}

	fmt.Fprint(w, "Healthy\n")
	return
}

func (a app) Debug(w http.ResponseWriter, r *http.Request) {
	gameid := chi.URLParam(r, "gameid")
	modid := chi.URLParam(r, "modid")
	taskid := chi.URLParam(r, "taskid")

	fmt.Fprintf(w, "gameid:%s modid:%s taskid:%s", gameid, modid, taskid)
}

func (a app) List(w http.ResponseWriter, r *http.Request) {

	var name string
	rows, err := a.db.Query("SELECT name FROM count")
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

	app := &app{
		db: db,
	}

	r := chi.NewMux()

	r.HandleFunc("/healthcheck", app.Healthcheck)

	r.Route("/games", func(r chi.Router) {
		r.Get("/", app.List)
		r.Route("/{gameid}", func(r chi.Router) {
			r.Get("/", app.Debug)
			r.Put("/", app.Debug)
			r.Route("/mods", func(r chi.Router) {
				r.Route("/{modid}", func(r chi.Router) {
					r.Get("/", app.Debug)
					r.Put("/", app.Debug)
				})
			})
		})
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", app.Debug)
		r.Delete("/{taskid}", app.Debug)
	})

	http.ListenAndServe(":"+port, r)
}
