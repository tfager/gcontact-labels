package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gcontact-labels/contacts"
	"gcontact-labels/web"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

func FetchContacts(client *http.Client) {
	ctx := context.Background()
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create people Client %v", err)
	}
	contacts, err := contacts.GetContacts(srv)
	if err != nil {
		log.Fatalf("Unable to retrieve contacts: %v", err)
	}

	if len(contacts) > 0 {
		for _, c := range contacts {
			fmt.Printf("%s\n", c.Name)
		}
	} else {
		fmt.Print("No connections found.")
	}
}

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

	// Get command from command line arguments with a parser library
	var groupId string
	flag.StringVar(&groupId, "group", "", "Group ID to get contacts from")
	flag.Parse()
	fmt.Printf("Command = %s\n", flag.Args()[0])

	if flag.Args()[0] == "group" {
		contacts, err := contacts.GetContactGroupMembers(peopleService, groupId)
		if err != nil {
			log.Fatalf("Unable to retrieve contacts: %v", err)
		}
		fmt.Printf("Contacts: %v\n", contacts)
	} else if flag.Args()[0] == "allgroups" {
		contacts.GetContactGroups(peopleService)
	} else {
		fmt.Println("Invalid command")
	}

}
