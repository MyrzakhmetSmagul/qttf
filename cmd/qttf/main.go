package main

import (
	"log"
	"path"
	"qttf/config"
	"qttf/internal/server"
	"qttf/pkg/db/postgres"

	_ "github.com/lib/pq"
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

	s := server.NewServer(cnf, psqlDB)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
