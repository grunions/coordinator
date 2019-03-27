package main

import (
	"fmt"
	"net/http"
	"os"
)

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World, %s", r.RemoteAddr)
}

func main() {

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, &handler{})
}
