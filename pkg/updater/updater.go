package updater

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"qttf/config"
	"qttf/internal/city"
	"qttf/internal/models"
	"qttf/internal/player"
	"qttf/internal/rating"
	"qttf/pkg/sheet"
	"syscall"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

var (
	cities  = make(map[string]models.City)
	ratings = make(map[string]models.Rating)
)

type Updater struct {
	cityRepo      city.Repository
	playerRepo    player.Repository
	ratingRepo    rating.Repository
	sheetsService *sheets.Service
	sheetsConfig  *config.SpreadsheetsConfig
}

func NewUpdater(cityRepo city.Repository, playerRepo player.Repository, ratingRepo rating.Repository, sheetsService *sheets.Service, sheetsConfig *config.SpreadsheetsConfig) *Updater {
	return &Updater{cityRepo: cityRepo, playerRepo: playerRepo, ratingRepo: ratingRepo, sheetsService: sheetsService, sheetsConfig: sheetsConfig}
}

func (u *Updater) Run(hours time.Duration) {
	interval := hours * time.Hour
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		_, shutdown := context.WithTimeout(context.Background(), time.Second)
		defer shutdown()

		log.Println("Updater stoped")
	}()

	for {
		// Запуск парсера
		log.Println("Updater start fetch info from google sheets")
		if err := u.fetch(); err != nil {
			log.Printf("fetch was failed: %v", err)
		}

		// Ожидание до следующего запуска
		<-ticker.C
	}
}

func (u *Updater) fetch() error {
	fields := "sheets(data(rowData(values(formattedValue,hyperlink))))"
	resp, err := u.sheetsService.Spreadsheets.Get(u.sheetsConfig.SpreadsheetId).Ranges(u.sheetsConfig.SheetName).Fields(googleapi.Field(fields)).Do()
	if err != nil {
		return err
	}

	rating, err := sheet.GetRatingFromSheet(resp)
	for i := 0; i < len(rating); i++ {
		if v, ok := cities[rating[i].Player.City.Hyperlink]; ok {
			rating[i].Player.City.Id = v.Id
		} else {
			err = u.cityRepo.Create(&rating[i].Player.City)
			if err != nil {
				return fmt.Errorf("an error occurred while retrieving data from the google table:%s %w\n%v",
					"\nerror occurred while creating a city record:", err, rating[i].Player.City)
			}

			cities[rating[i].Player.City.Hyperlink] = rating[i].Player.City
		}

		if _, ok := ratings[rating[i].Player.Hyperlink]; !ok {
			err = u.playerRepo.Create(&rating[i].Player)
			if err != nil {
				return fmt.Errorf("an error occurred while retrieving data from the google table:%s %w\n%v",
					"\nerror occurred while creating a player:", err, rating[i].Player)
			}

			err = u.ratingRepo.Create(&rating[i])
			if err != nil {
				return fmt.Errorf("an error occurred while retrieving data from the google table:%s %w",
					"\nerror occurred while creating a rating record:", err)
			}
			ratings[rating[i].Player.Hyperlink] = rating[i]
		}
	}

	return nil
}
