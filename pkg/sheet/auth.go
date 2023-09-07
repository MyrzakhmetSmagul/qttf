package sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"golang.org/x/oauth2"
)

// Retrieves a token from a local file.
func tokenFromFile(tokenPath string) (*oauth2.Token, error) {
	f, err := os.Open(path.Clean(tokenPath))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func GetClient(ctx context.Context, config *oauth2.Config, tokenPath string) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the client: %w", err)
	}
	return config.Client(ctx, tok), nil
}
