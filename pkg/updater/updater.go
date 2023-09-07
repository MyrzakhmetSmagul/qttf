package updater

import (
	"context"
	"log"
	"os"
	"os/signal"
	"qttf/config"
	"qttf/internal/city"
	"qttf/internal/player"
	"qttf/internal/rating"
	"qttf/pkg/sheet"
	"syscall"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
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
		log.Println("Updater start update info from google sheets")
		if err := u.update(); err != nil {
			log.Printf("update was failed: %v", err)
		}

		// Ожидание до следующего запуска
		<-ticker.C
	}
}

func (u *Updater) update() error {
	fields := "sheets(data(rowData(values(formattedValue,hyperlink))))"
	resp, err := u.sheetsService.Spreadsheets.Get(u.sheetsConfig.SpreadsheetId).Ranges(u.sheetsConfig.SheetName).Fields(googleapi.Field(fields)).Do()
	if err != nil {
		return err
	}

	rating, err := sheet.GetRatingFromSheet(resp)

	return nil
}
