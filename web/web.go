package web

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func serveTemplate(w http.ResponseWriter, r *http.Request, data any) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Print(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Printf("Data = %v\n", data)
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

type ReqAuthPageData struct {
	Url string
}

func StartWebServer(config *oauth2.Config) {
	// Get the URL to OAuth2 consent page
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// TODO: handle content types
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := ReqAuthPageData{
			Url: url,
		}
		fmt.Printf("Data.url = %v\n", data.Url)
		serveTemplate(w, r, data)
	})

	// Create a handler to handle the callback from the OAuth2 consent page
	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		// Get the authorization code from the request
		code := r.URL.Query().Get("code")
		if code == "" {
			log.Fatal("Code not found in OAuth callback")
		}
		fmt.Fprintf(w, "Code received, you may now close this browser window")
	})
	fmt.Printf("Starting server on http://localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}
