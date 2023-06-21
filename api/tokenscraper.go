// Package that contains the models of the JSON objects used in the requests and responses and the methods to create
// new objects that describes the API actions, like the Upload or the AtTokenScraper
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"log"

	"github.com/garfieldkwong/gphotosuploader/auth"
	"golang.org/x/net/html"
)

const (
	GooglePhotoUrl = "https://photos.google.com/"
)

// AtTokenScraper used to scape tokens to upload images
type AtTokenScraper struct {
	credentials auth.CookieCredentials
}

// Create a new scraper for the at token. This token is user-dependent, so you need to create a new token scraper
// for each Credentials object.
func NewAtTokenScraper(credentials auth.CookieCredentials) *AtTokenScraper {
	return &AtTokenScraper{
		credentials: credentials,
	}
}

// Use this method to get a new at token. The method makes an http request to Google and uses the user credentials
func (ts *AtTokenScraper) ScrapeNewAtToken() (string, error) {
	page, err := ts.getHomePage()
	if err != nil {
		return "", err
	}

	script, err := findScript(page)
	if err != nil {
		return "", err
	}

	return findTokenInScript(script)
}

func (ts *AtTokenScraper) getHomePage() (*http.Response, error) {
	req, err := http.NewRequest("GET", GooglePhotoUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create the request to get the Google Photos homepage (%v)", err)
	}

	// Make the request
	if res, err := ts.credentials.Client.Do(req); err != nil {
		return nil, fmt.Errorf("can't complete the request to get the Google Photos homepage (%v)", err)
	} else {
		return res, nil
	}
}

func findScript(page *http.Response) (string, error) {
	t := html.NewTokenizer(page.Body)
	for {
		tt := t.Next()
                curToken := t.Token()                

		switch {
		case tt == html.ErrorToken: // End of html document
			return "", errors.New("can't find the script tag with the token in the response")

		case tt == html.StartTagToken && curToken.Data == "script": // We need the first script tag
                        found := false
                        for _, attr := range curToken.Attr {
                                if (attr.Key == "data-id" && attr.Val == "_gd") {
                                	found = true
                                        break
				}
			}
                        if (found) {
				t.Next()

				// Get the script string
				return t.Token().Data, nil
			}
		}
	}
}

func findTokenInScript(script string) (string, error) {
	// The script assigns an object to the global window object. We are going to parse the script as a JSON
	// so we need to get rid of the assignment code
	// Parse the json
	equalsIndex := strings.Index(script, "=")
	start := equalsIndex + 1
	end := len(script) - 1
	script = script[start:end]
        log.Println("xxxx", script);
	object := ApiTokenContainer{}
	if err := json.NewDecoder(strings.NewReader(script)).Decode(&object); err != nil {
		return "", fmt.Errorf("can't parse the JSON object that contains the at token (%v)", err)
	}

	return object.Token, nil
}
