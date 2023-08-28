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
	"strconv"
	"strings"

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

type Rating struct {
	Id         int
	Position   int
	Player     Player
	Rating     int
	UpdateTime string
}

type Player struct {
	Id      int
	Surname string
	Name    string
	City    City
}

type City struct {
	Id   int
	Name string
}

type DB struct {
	db *sql.DB
}

type Config struct {
	Database DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

var cities = make(map[string]int)

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

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, sheetName).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	playerFile, _ := os.Create(path.Clean("./sql/player.sql"))
	playerScript := ""
	ratingFile, _ := os.Create(path.Clean("./sql/rating.sql"))
	ratingScript := ""

	defer playerFile.Close()
	defer ratingFile.Close()
	saveResp("resp.json", resp)
	for i, row := range resp.Values {
		if i == 0 || row[0] == "" {
			continue
		}
		record := Rating{}
		num, _ := strconv.Atoi(row[0].(string))
		record.Position = num
		fullName := strings.Split(row[1].(string), " ")

		record.Player.Surname = fullName[0]
		if len(fullName) > 1 {
			record.Player.Name = fullName[1]
		}
		city := row[2].(string)
		record.Player.City.Name = city
		cities[city]++

		playerScript += record.Player.ToScript()
		num, _ = strconv.Atoi(row[3].(string))
		record.Rating = num
		record.UpdateTime = row[4].(string)
		ratingScript += record.ToScript()
	}
	playerFile.Write([]byte(playerScript))
	ratingFile.Write([]byte(ratingScript))

	// cityScript, _ := os.ReadFile(path.Clean("./sql/city.sql"))

	// db := getDB()
	// err := db.db.QueryRow("INSERT INTO city (city_name) VALUES ('Алматы') ON CONFLICT (city_name) DO NOTHING;SELECT city_id FROM city WHERE city_name='Алматы'").Scan(&i)
	// _, err = db.db.Exec(string(cityScript))
	// if err != nil {
	// 	log.Fatal("city: ", err)
	// }
	// _, err = db.db.Exec(PlayerScript)
	// if err != nil {
	// 	log.Fatal("Player: ", err)
	// }
	// _, err = db.db.Exec(ratingScript)
	// if err != nil {
	// 	log.Fatal("rating: ", err)
	// }
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
	configFile, _ := os.Open(path.Clean("./sql/sql_config.json"))

	var config Config
	err := json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port, config.Database.User,
		config.Database.Password, config.Database.DBName, config.Database.SSLMode,
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

func (p Player) ToScript() string {
	return fmt.Sprintf("INSERT INTO player (player_name, player_surname, city_id) VALUES('%s', '%s', (SELECT city_id FROM city WHERE city_name='%s'));\n", p.Name, p.Surname, p.City.Name)
}

func (r Rating) ToScript() string {
	return fmt.Sprintf("INSERT INTO rating(player_id, rating, last_update) VALUES((SELECT player_id FROM player WHERE player_name='%s' AND player_surname='%s'), %d, '%s');\n", r.Player.Name, r.Player.Surname, r.Rating, r.UpdateTime)
}

func cityActualization(db DB) error {
	query := `INSERT INTO city (city_name) VALUES ($1) ON CONFLICT (city_name) DO NOTHING;
	SELECT city_id FROM city WHERE city_name=$1`

	for city, _ := range cities {
		var id int
		err := db.db.QueryRow(query, city).Scan(&id)
		if err != nil {
			return fmt.Errorf("cityActualization%w", err)
		}

		cities[city] = id
	}

	return nil
}
