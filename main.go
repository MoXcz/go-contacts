package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type config struct {
	templates *template.Template
}

type PageData struct {
	Contacts   []Contact
	SearchTerm string
}

func main() {
	listenAddr := ":3000"
	mux := http.NewServeMux()
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Printf("Error: Could not find templates %v", err)
		return
	}
	config := config{
		templates: t,
	}

	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/contacts", config.contacts)

	srv := http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	fmt.Println("Listening on port", listenAddr)
	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusPermanentRedirect)
}
