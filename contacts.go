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
	Phone     int
	Id        int
}

func newContacts() []Contact {
	return []Contact{
		newContactData("Pedro", "Sanchez", "pedro@gm.com", 113),
		newContactData("Juan", "Mama", "juan@gm.com", 112),
	}
}

var id = 0

func newContactData(firstName, lastName, email string, phone int) Contact {
	id++
	return Contact{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Id:        id,
	}
}

var contacts = newContacts()

func getContacts(w http.ResponseWriter, r *http.Request) {
	c := contacts
	search := r.URL.Query().Get("q")
	if search != "" {
		c = searchContacts(search, c)
	}

	data := PageData{
		Contacts:   c,
		SearchTerm: search,
	}

	err := renderTemplate(w, "contacts", data)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}

func newContact(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, "new", nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}

func createNewContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	phone, err := strconv.Atoi(r.FormValue("phone"))
	if err != nil {
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	contact := newContactData(firstName, lastName, email, phone)
	contacts = append(contacts, contact)

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func searchContacts(q string, contacts []Contact) []Contact {
	filteredContacts := []Contact{}
	for _, contact := range contacts {
		if strings.Contains(contact.FirstName, q) ||
			strings.Contains(contact.LastName, q) ||
			strings.Contains(contact.Email, q) ||
			strings.Contains(strconv.Itoa(contact.Phone), q) {
			filteredContacts = append(filteredContacts, contact)
		}
	}
	return filteredContacts
}
