package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
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
	Error     string
}

// map[string][]string
// "phoneError -> "Invalid phone num"

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

func createNewContactGet(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, "new", nil)
	if err != nil {
		fmt.Printf("Error rendering template: %v", err)
		return
	}
}

func createNewContactPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	// TODO: phone := r.FormValue("phone") to preserve state on errors
	phoneNum, err := strconv.Atoi(r.FormValue("phone"))
	if err != nil {
		// use text phone here?
		err := renderTemplate(w, "new", Contact{
			FirstName: firstName,
			LastName:  lastName,
			Phone:     0,
			Email:     email,
			Error:     "Invalid phone number",
		})
		// this is really bad <3
		if err != nil {
			fmt.Printf("Error rendering template: %v", err)
			return
		}
		return
	}

	contact := newContactData(firstName, lastName, email, phoneNum)
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

func editContactGet(w http.ResponseWriter, r *http.Request) {
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
	renderTemplate(w, "edit", contact)
}

func editContactPost(w http.ResponseWriter, r *http.Request) {
	contact_id, err := uuid.Parse(r.PathValue("contact_id"))
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	phoneNum, err := strconv.Atoi(r.FormValue("phone"))
	if err != nil {
		err := renderTemplate(w, "edit", Contact{
			FirstName: firstName,
			LastName:  lastName,
			Phone:     0,
			Email:     email,
			Error:     "Invalid phone number",
			Id:        contact_id,
		})
		if err != nil {
			fmt.Printf("Error rendering template: %v", err)
			return
		}
		return
	}

	for i, contact := range contacts {
		if contact.Id == contact_id {
			contacts[i] = Contact{
				FirstName: firstName,
				LastName:  lastName,
				Phone:     phoneNum,
				Email:     email,
				Id:        contact.Id,
			}
		}
	}
	http.Redirect(w, r, "/contacts/"+contact_id.String(), http.StatusSeeOther)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	contact_id, err := uuid.Parse(r.PathValue("contact_id"))
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	for i, contact := range contacts {
		if contact.Id == contact_id {
			contacts = slices.Delete(contacts, i, i+1)
		}
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}
