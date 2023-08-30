package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"qttf/internal/models"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var (
	config *oauth2.Config
	// TODO: randomize it
	oauthStateString = "pseudo-random"
	tok              = make(chan *oauth2.Token)
	cities           = make(map[string]models.City)
)

func init() {
	credential, err := os.ReadFile(path.Clean("./credential.json"))
	if err != nil {
		log.Fatalf("Unable to read of credential json: %v", err)
	}
	config, err = google.ConfigFromJSON(credential, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
}

func main() {
}

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
			cities[rating[i].Player.City.Name] = rating[i].Player.City
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

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	_, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatal(err)
	}
}

func cityActualization(db *sql.DB) error {
	querySelect := `SELECT city_id FROM city WHERE city_link=$1`

	for _, city := range cities {
		_, err := db.Exec(city.ToInsertScript())
		if err != nil {
			return fmt.Errorf("cityInsertion: %w\nScript: %s", err, city.ToInsertScript())
		}

		var id int
		err = db.QueryRow(querySelect, city.Hyperlink).Scan(&id)
		if err != nil {
			return fmt.Errorf("citySelection: %w", err)
		}
		city.Id = id
		cities[city.Name] = city
	}

	return nil
}

func ratingActualization(db *sql.DB, rating []models.Rating) {
	for _, rat := range rating {
		err := rat.Player.PlayerActualization(db)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(rat.ToInsertScript())
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("done!")
}
