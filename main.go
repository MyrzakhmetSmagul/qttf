package main

import (
	"context"
	"database/sql"
	"encoding/json"
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
	"google.golang.org/api/option"
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

type DB struct {
	db *sql.DB
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

func main() {
	// go func() {
	// 	http.HandleFunc("/callback", handleGoogleCallback)
	// 	fmt.Println(http.ListenAndServe(":8080", nil))
	// }()

	ctx := context.Background()
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := "1aWTIE7G7R_p-CSUTpU7UswTOLSFHkPYB_rGG5qzbTYg"
	sheetName := "KZ"

	resp, err := srv.Spreadsheets.Get(spreadsheetId).Ranges(sheetName).Fields("sheets(data(rowData(values(formattedValue,hyperlink))))").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	saveResp("resp.json", resp)
	rating, err := getInfoFromSheet(resp)
	if err != nil {
		log.Fatal(err)
	}

	db := getDB()
	err = cityActualization(db.db)
	if err != nil {
		log.Fatal("cityActualization: ", err)
	}
	ratingActualization(db.db, rating)
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
	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatal(err)
	}
	tok <- token
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

	token := <-tok
	return token
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

func getDB() DB {
	configFile, _ := os.Open(path.Clean("./sql_config.json"))

	var config DatabaseConfig
	err := json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User,
		config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database!")
	return DB{
		db: db,
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
