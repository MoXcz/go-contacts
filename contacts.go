package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Contact struct {
	FirstName string
	LastName  string
	Email     string
	Num       int
	Id        int
}

func newContacts() []Contact {
	return []Contact{
		newContact("Pedro", "Sanchez", "pedro@gm.com", 113),
		newContact("Juan", "Mama", "juan@gm.com", 112),
	}
}

var id = 0

func newContact(firstName, lastName, email string, num int) Contact {
	id++
	return Contact{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Num:       num,
		Id:        id,
	}
}

func (cfg *config) contacts(w http.ResponseWriter, r *http.Request) {
	contacts := newContacts()
	search := r.URL.Query().Get("q")
	if search != "" {
		contacts = searchContacts(search, contacts)
	}

	data := PageData{
		Contacts:   contacts,
		SearchTerm: search,
	}

	err := cfg.templates.ExecuteTemplate(w, "index", data)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}

func searchContacts(q string, contacts []Contact) []Contact {
	filteredContacts := []Contact{}
	for _, contact := range contacts {
		if strings.Contains(contact.FirstName, q) ||
			strings.Contains(contact.LastName, q) ||
			strings.Contains(contact.Email, q) ||
			strings.Contains(strconv.Itoa(contact.Num), q) {
			filteredContacts = append(filteredContacts, contact)
		}
	}
	return filteredContacts
}
