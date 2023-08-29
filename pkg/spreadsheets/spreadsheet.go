package spreadsheets

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"qttf/internal/models"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
)

func getInfoFromSheet(resp *sheets.Spreadsheet) ([]models.Rating, error) {
	if len(resp.Sheets) > 0 && len(resp.Sheets[0].Data) > 0 && len(resp.Sheets[0].Data[0].RowData) > 0 {
		rating := make([]models.Rating, len(resp.Sheets[0].Data[0].RowData[1:]))
		for i, row := range resp.Sheets[0].Data[0].RowData[1:] {
			num, err := strconv.Atoi(row.Values[0].FormattedValue)
			if err != nil {
				return nil, fmt.Errorf("getInfoFromSheet: %w", err)
			}

			rating[i].Position = num
			fullName := strings.Split(row.Values[1].FormattedValue, " ")
			rating[i].Player.Surname = fullName[0]
			if len(fullName) > 1 {
				rating[i].Player.Name = fullName[1]
			}

			rating[i].Player.Hyperlink = row.Values[1].Hyperlink
			rating[i].Player.City.Name = row.Values[2].FormattedValue
			rating[i].Player.City.Hyperlink = row.Values[2].Hyperlink
			// cities[rating[i].Player.City.Name] = rating[i].Player.City
			num, err = strconv.Atoi(row.Values[3].FormattedValue)
			if err != nil {
				return nil, fmt.Errorf("getInfoFromSheet: %w", err)
			}

			rating[i].Rating = num
			t, err := time.Parse("2.1.2006", row.Values[4].FormattedValue)
			if err != nil {
				return nil, fmt.Errorf("getInfoFromSheet: %w", err)
			}

			rating[i].UpdateTime = t.Format("2006/01/02")
		}
		return rating, nil
	}
	return nil, fmt.Errorf("getInfoFromSheet: data from resp not valid")
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveResp(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	return nil
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
func saveResp(path string, token interface{}) {
	fmt.Printf("Saving resp json to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
