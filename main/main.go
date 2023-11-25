package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gcontact-labels/contacts"
	"gcontact-labels/web"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
)

func main() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, people.ContactsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	config.RedirectURL = "http://localhost:8080/oauth2callback"
	var client *http.Client
	if !contacts.TokenFileExists() {
		var code string
		// Get the URL to OAuth2 consent page
		// TODO: state-token should be random string
		authCodeUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		web.StartWebServer(authCodeUrl, &code)
		client = contacts.GetClient(config, code)
	} else {
		client = contacts.GetClientFromFile(config)
	}

	//FetchContacts(client)
	peopleService, err := contacts.CreateService(client)
	if err != nil {
		log.Fatalf("Unable to create people service: %v", err)
	}

	// Get command from command line arguments with a parser library
	var groupId string
	flag.StringVar(&groupId, "group", "", "Group ID to get contacts from")
	flag.Parse()

	if flag.Args()[0] == "contacts" {
		if groupId == "" {
			log.Fatalf("Group ID is required")
		}
		contacts, err := contacts.GetContactGroupMembers(peopleService, groupId)
		if err != nil {
			log.Fatalf("Unable to retrieve contacts: %v", err)
		}
		for _, contact := range contacts {
			fmt.Printf("Contact: %v, %v, %v, %v\n", contact.Name, contact.StreetAddress, contact.City, contact.PostalCode)
		}
	} else if flag.Args()[0] == "groups" {
		groups, err := contacts.GetContactGroups(peopleService)
		if err != nil {
			log.Fatalf("Unable to retrieve contacts: %v", err)
		}
		for _, group := range groups {
			fmt.Printf("%v, %v\n", group.Name, group.Id)
		}
	} else {
		fmt.Println("Invalid command. Usage: contacts groups|contacts -group <group id>")
	}

}
