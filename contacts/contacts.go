package contacts

import (
        "context"
        "fmt"
        "log"
        "net/http"

        "golang.org/x/oauth2"
        "google.golang.org/api/people/v1"
)
// Define Contact struct
type Contact struct {
    Name string
    StreetAddress string
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
        tok := GetTokenFromWeb(config)
        return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code: %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web: %v", err)
        }
        return tok
}

// Get contacts from Google People API
func GetContacts(peopleService *people.Service) ([]*Contact, error) {
	// Create a new call to people api
	call := peopleService.People.Connections.List("people/me")
	call = call.PersonFields("names,addresses")
	call = call.PageSize(20)

	// Get response from people api
	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	// Create a slice of contacts
	var contacts []*Contact

	// Loop through response and add contacts to slice
	for _, c := range res.Connections {
		// Create a new contact
		contact := &Contact{}

		// Check if contact has a name
		if len(c.Names) > 0 {
			contact.Name = c.Names[0].DisplayName
		}

		// Check if contact has an email
		if len(c.Addresses) > 0 {
			contact.StreetAddress = c.Addresses[0].StreetAddress
		}

		// Add contact to slice
		contacts = append(contacts, contact)
	}

	// Return slice of contacts
	return contacts, nil
}

