package gsheets

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

func GetClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the client: %w", err)
	}
	return config.Client(ctx, tok), nil
}
