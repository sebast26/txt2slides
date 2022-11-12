package google

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
)

type OAuthClient struct {
	HttpClient *http.Client
}

// NewOAuthClient creates new OAuthClient that can be used to create Google services.
func NewOAuthClient(credentialsFilePath, tokenFilePath string, forceWeb bool) (*OAuthClient, error) {
	b, err := os.ReadFile(credentialsFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "reading client secret file failed")
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/drive", "https://www.googleapis.com/auth/presentations")
	if err != nil {
		return nil, errors.Wrap(err, "parsing client secret file to config failed")
	}

	tok, err := tokenFromFile(tokenFilePath)
	if err != nil && forceWeb {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		err = saveToken(tokenFilePath, tok)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, errors.Wrapf(err, "count not get token from %s", tokenFilePath)
	}

	return &OAuthClient{HttpClient: config.Client(context.Background(), tok)}, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, errors.Wrap(err, "failed to read authorization code")
	}

	ctx := context.Background()
	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve token from web")
	}
	return tok, nil
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return errors.Wrap(err, "failed to cache OAuth token")
	}
	return json.NewEncoder(f).Encode(token)
}
