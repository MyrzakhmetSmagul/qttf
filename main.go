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

	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Rating struct {
	Id         int
	Position   int
	User       User
	Rating     int
	UpdateTime string
}

type User struct {
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
	ctx := context.Background()
	credential, err := os.ReadFile(path.Clean("./credential.json"))
	fmt.Println(string(credential))
	if err != nil {
		log.Fatalf("Unable to read of credential json: %v", err)
	}
	config, err := google.ConfigFromJSON(credential, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

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

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	}
	fmt.Println(resp.Values[0]...)
	// xlFile, err := xlsx.OpenFile("excelFileName")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// userFile, _ := os.Create(path.Clean("./sql/user.sql"))
	// userScript := ""
	// ratingFile, _ := os.Create(path.Clean("./sql/rating.sql"))
	// ratingScript := ""

	// sheet := xlFile.Sheet["KZ"]
	// defer userFile.Close()
	// defer ratingFile.Close()
	// for i, row := range sheet.Rows {
	// 	if i == 0 || row.Cells[0].String() == "" {
	// 		continue
	// 	}
	// 	record := Rating{}
	// 	temp, _ := strconv.Atoi(row.Cells[0].String())
	// 	record.Position = temp
	// 	fullName := strings.Split(row.Cells[1].String(), " ")
	// 	record.User.Surname = fullName[0]
	// 	if len(fullName) > 1 {
	// 		record.User.Name = fullName[1]
	// 	}
	// 	city := row.Cells[2].String()
	// 	record.User.City.Name = city
	// 	cities[city]++

	// 	userScript += record.User.ToScript()
	// 	temp, _ = strconv.Atoi(row.Cells[3].String())
	// 	record.Rating = temp
	// 	excelTime, _ := row.Cells[4].GetTime(false) // Get the time value from the cell
	// 	goTime := excelTime.UTC()                   // Convert to UTC time
	// 	record.UpdateTime = goTime.Format("02.01.2006")
	// 	ratingScript += record.ToScript()
	// }
	// userFile.Write([]byte(userScript))
	// ratingFile.Write([]byte(ratingScript))

	// cityScript, _ := os.ReadFile(path.Clean("./sql/city.sql"))

	// db := getDB()
	// // err := db.db.QueryRow("INSERT INTO city (city_name) VALUES ('Алматы') ON CONFLICT (city_name) DO NOTHING;SELECT city_id FROM city WHERE city_name='Алматы'").Scan(&i)
	// _, err = db.db.Exec(string(cityScript))
	// if err != nil {
	// 	log.Fatal("city: ", err)
	// }
	// _, err = db.db.Exec(userScript)
	// if err != nil {
	// 	log.Fatal("user: ", err)
	// }
	// _, err = db.db.Exec(ratingScript)
	// if err != nil {
	// 	log.Fatal("rating: ", err)
	// }
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
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
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
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
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

func (u User) ToScript() string {
	return fmt.Sprintf("INSERT INTO \"user\"(user_name, user_surname, city_id) VALUES('%s', '%s', (SELECT city_id FROM city WHERE city_name='%s'));\n", u.Name, u.Surname, u.City.Name)
}

func (r Rating) ToScript() string {
	return fmt.Sprintf("INSERT INTO rating(user_id, rating, position, last_update) VALUES((SELECT user_id FROM \"user\" WHERE user_name='%s' AND user_surname='%s'), %d, %d, '%s');\n", r.User.Name, r.User.Surname, r.Rating, r.Position, r.UpdateTime)
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

// func playerUpdate
