package contacts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

var tokenFile = "token.json"

const CARD_TEXT_FIELD = "Joulukorttiteksti"

// Define Contact struct
type Contact struct {
	Name          string
	StreetAddress string
	City          string
	PostalCode    string
}

type ContactGroup struct {
	Name string
	Id   string
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config, authCode string) *http.Client {
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	saveToken(tokenFile, tok)
	return config.Client(context.Background(), tok)
}

func GetClientFromFile(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		log.Fatalf("Unable to retrieve token from file: %v", err)
	}
	return config.Client(context.Background(), tok)
}

func TokenFileExists() bool {
	_, error := os.Stat(tokenFile)
	return !os.IsNotExist(error)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func CreateService(client *http.Client) (*people.Service, error) {
	ctx := context.Background()
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func parseContact(contact *people.Person) *Contact {
	// Create a new contact
	c := &Contact{}

	// Check for card text
	for _, u := range contact.UserDefined {
		if u.Key == CARD_TEXT_FIELD {
			c.Name = u.Value
		}
	}
	// Check if contact has a name
	if c.Name == "" && len(contact.Names) > 0 {
		c.Name = contact.Names[0].DisplayName
	}

	// Check if contact has an email
	if len(contact.Addresses) > 0 {
		c.StreetAddress = contact.Addresses[0].StreetAddress
	}

	// Check if contact has an email
	if len(contact.Addresses) > 0 {
		c.City = contact.Addresses[0].City
		c.PostalCode = contact.Addresses[0].PostalCode
	}

	return c
}

// Get contacts from Google People API
func GetContactGroups(peopleService *people.Service) ([]*ContactGroup, error) {
	// Create a new call to people api
	call := peopleService.ContactGroups.List()
	call = call.PageSize(100)

	// Get response from people api
	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	var groups []*ContactGroup
	// Loop through response and add contacts to slice
	for _, c := range res.ContactGroups {
		group := &ContactGroup{
			Name: c.Name,
			Id:   c.ResourceName,
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func GetContactGroupMembers(peopleService *people.Service, groupId string) ([]*Contact, error) {
	call := peopleService.ContactGroups.Get(groupId)
	call = call.MaxMembers(100)

	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	batchCall := peopleService.People.GetBatchGet()
	batchCall = batchCall.ResourceNames(res.MemberResourceNames...)
	batchCall = batchCall.PersonFields("names,addresses,userDefined")

	batchRes, err := batchCall.Do()
	if err != nil {
		return nil, err
	}

	var contacts []*Contact
	// Loop through response and add contacts to slice
	for _, r := range batchRes.Responses {
		person := r.Person
		contact := parseContact(person)
		contacts = append(contacts, contact)
	}

	// Return slice of contacts
	return contacts, nil
}
