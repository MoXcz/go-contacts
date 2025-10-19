package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Contacts   []Contact
	SearchTerm string
}

func main() {
	listenAddr := ":3000"
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", index)
	mux.HandleFunc("GET /contacts", getContacts)
	mux.HandleFunc("GET /contacts/new", createNewContactGet)
	mux.HandleFunc("POST /contacts/new", createNewContactPost)
	mux.HandleFunc("GET /contacts/{contact_id}", getContact)
	mux.HandleFunc("GET /contacts/{contact_id}/edit", editContactGet)
	mux.HandleFunc("POST /contacts/{contact_id}/edit", editContactPost)
	mux.HandleFunc("GET /contacts/{contact_id}/email", getEmailValidation)
	mux.HandleFunc("DELETE /contacts/{contact_id}", deleteContact)

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

func renderTemplate(w http.ResponseWriter, name string, data any) error {
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		fmt.Sprintf("templates/%s.html", name),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("Error: Content template not found. %w", err)
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		return fmt.Errorf("Error: Content template not executed. %w", err)
	}
	return nil
}
