package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gcontact-labels/contacts"
	"gcontact-labels/web"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

func FetchContacts(config *oauth2.Config) {
	ctx := context.Background()
	client := contacts.GetClient(config)
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

	web.StartWebServer(config)
}
