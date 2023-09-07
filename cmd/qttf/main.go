package main

import (
	"context"
	"log"
	"path"
	"qttf/config"
	cityRepository "qttf/internal/city/repository"
	playerRepository "qttf/internal/player/repository"
	ratingRepository "qttf/internal/rating/repository"
	"qttf/internal/server"
	"qttf/pkg/db/postgres"
	"qttf/pkg/sheet"
	"qttf/pkg/updater"

	_ "github.com/lib/pq"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	log.Println("Starting api server")

	cnf, err := config.ParseConfig(path.Clean("./config.json"))
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	psqlDB, err := postgres.NewPsqlDB(cnf)
	if err != nil {
		log.Fatalf("NewPsqlDB: %v", err)
	}

	defer psqlDB.Close()
	log.Printf("Postgres connected, Status: %v", psqlDB.Stats())

	cityRepo := cityRepository.NewCityRepository(psqlDB)
	playerRepo := playerRepository.NewPlayerRepository(psqlDB)
	ratingRepo := ratingRepository.NewRatingRepository(psqlDB)

	ctx := context.Background()
	client, err := sheet.GetClient(ctx, cnf.GoogleCOnfig, cnf.TokenPath)
	if err != nil {
		log.Fatalf("sheets.GetClient: %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("sheets.NewService: %v", err)
	}

	updater := updater.NewUpdater(cityRepo, playerRepo, ratingRepo, srv, &cnf.Spreadsheet)

	go func() {
		updater.Run(24)
	}()
	s := server.NewServer(cnf, psqlDB)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
