package main

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
	tm     map[string]*template.Template
}

func main() {
	listenAddr := ":3000"
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	tm, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		logger: logger,
		tm:     tm,
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", index)
	mux.HandleFunc("GET /contacts", app.getContacts)
	mux.HandleFunc("GET /contacts/new", app.createNewContactGet)
	mux.HandleFunc("POST /contacts/new", app.createNewContactPost)
	mux.HandleFunc("GET /contacts/{contact_id}", app.getContact)
	mux.HandleFunc("GET /contacts/{contact_id}/edit", app.editContactGet)
	mux.HandleFunc("POST /contacts/{contact_id}/edit", app.editContactPost)
	mux.HandleFunc("GET /contacts/{contact_id}/email", getEmailValidation)
	mux.HandleFunc("DELETE /contacts/{contact_id}", deleteContact)

	srv := http.Server{
		Addr:    listenAddr,
		Handler: logRequest(mux),
	}

	fmt.Println("Listening on port", listenAddr)
	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusPermanentRedirect)
}

