package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Contact struct {
	FirstName string
	LastName  string
	Email     string
	Phone     int
	Id        int
	Errors    map[string]string
}

func newContacts() []Contact {
	return []Contact{
		newContactData("Pedro", "Sanchez", "pedro@gm.com", 113),
		newContactData("Juan", "Mama", "juan@gm.com", 112),
	}
}

func newContactData(firstName, lastName, email string, phone int) Contact {
	c := Contact{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Id:        counter,
		Errors:    map[string]string{},
	}
	counter += 1
	return c
}

var contacts = newContacts()
var counter = 0

func getContacts(w http.ResponseWriter, r *http.Request) {
	c := contacts
	search := r.URL.Query().Get("q")
	if search != "" {
		c = searchContacts(search, c)
	}

	page := r.URL.Query().Get("page")
	pageNum, err := strconv.Atoi(page)
	if page != "" && err != nil {
		log.Println("Invalid page number")
		return
	}

	if pageNum <= 0 {
		c = c[:10]
	} else if (pageNum+1)*10 > len(c) {
		c = c[pageNum*10:]
	} else {
		c = c[pageNum*10 : (pageNum+1)*10]
	}

	data := PageData{
		Contacts:   c,
		SearchTerm: search,
		Page:       pageNum,
	}

	err = renderTemplate(w, "contacts", data)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}

func createNewContactGet(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, "new", Contact{Errors: map[string]string{}})
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

	contact, err := parseContactForm(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	contacts = append(contacts, contact)

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func getContact(w http.ResponseWriter, r *http.Request) {
	contact_id, err := strconv.Atoi(r.PathValue("contact_id"))
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
	contact_id, err := strconv.Atoi(r.PathValue("contact_id"))
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

func getEmailValidation(w http.ResponseWriter, r *http.Request) {
	msg := ""
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
	}

	email := r.FormValue("email")
	if !isEmailValid(email) {
		msg = "Invalid email"
	}

	w.Write([]byte(msg))
}

func editContactPost(w http.ResponseWriter, r *http.Request) {
	contact_id, err := strconv.Atoi(r.PathValue("contact_id"))
	if err != nil {
		http.Error(w, "Unable to parse contact_id", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	contact, err := parseContactForm(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	for i, c := range contacts {
		if c.Id == contact_id {
			contacts[i].FirstName = contact.FirstName
			contacts[i].LastName = contact.LastName
			contacts[i].Email = contact.Email
			contacts[i].Phone = contact.Phone
		}
	}
	http.Redirect(w, r, "/contacts/"+fmt.Sprintf("%d", contact_id), http.StatusSeeOther)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	contact_id, err := strconv.Atoi(r.PathValue("contact_id"))
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

func parseContactForm(w http.ResponseWriter, r *http.Request) (Contact, error) {
	errors := map[string]string{}
	firstName := r.FormValue("firstName")
	if firstName == "" {
		errors["firstName"] = "Invalid first name"
	}
	lastName := r.FormValue("lastName")
	if lastName == "" {
		errors["lastName"] = "Invalid last name"
	}
	email := r.FormValue("email")
	if !isEmailValid(email) {
		errors["mail"] = "Invalid email"
	}
	phone := r.FormValue("phone")
	phoneNum, err := strconv.Atoi(phone)
	if err != nil {
		errors["phone"] = "Invalid phone number"
		log.Println(err)
	}

	if len(errors) != 0 {
		err := renderTemplate(w, "new", Contact{
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phoneNum,
			Email:     email,
			Errors:    errors,
		})
		if err != nil {
			fmt.Printf("Error rendering template: %v", err)
			return Contact{}, err
		}
		return Contact{}, fmt.Errorf("Error parsing values in contact form")
	}

	return newContactData(firstName, lastName, email, phoneNum), nil
}

func isEmailValid(email string) bool {
	// invalidate repeated emails
	for _, contact := range contacts {
		if contact.Email == email {
			return false
		}
	}

	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
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

func getContactByID(contact_id int) (Contact, error) {
	for _, contact := range contacts {
		if contact_id == contact.Id {
			return contact, nil
		}
	}
	return Contact{}, errors.New("No matches found")
}
