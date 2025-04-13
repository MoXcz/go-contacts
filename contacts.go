package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Contact struct {
	FirstName string
	LastName  string
	Email     string
	Phone     int
	Id        uuid.UUID
}

func newContacts() []Contact {
	return []Contact{
		newContactData("Pedro", "Sanchez", "pedro@gm.com", 113),
		newContactData("Juan", "Mama", "juan@gm.com", 112),
	}
}

func newContactData(firstName, lastName, email string, phone int) Contact {
	return Contact{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Id:        uuid.New(),
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
		fmt.Printf("Error rendering template: %v", err)
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
	// TODO: handle errors of form values
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

func getContactByID(contact_id uuid.UUID) (Contact, error) {
	for _, contact := range contacts {
		if contact_id == contact.Id {
			return contact, nil
		}
	}
	return Contact{}, errors.New("No matches found")
}

func getContact(w http.ResponseWriter, r *http.Request) {
	contact_id, err := uuid.Parse(r.PathValue("contact_id"))
	if err != nil {
		log.Printf("Error parsing contact id: %v", err)
		return
	}

	contact, err := getContactByID(contact_id)
	if err != nil {
		log.Printf("Error finding contact: %v", err)
		return
	}
	renderTemplate(w, "view", contact)
}
