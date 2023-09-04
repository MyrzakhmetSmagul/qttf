package main

import (
	"database/sql"
	"fmt"
	"log"
	"qttf/internal/models"

	_ "github.com/lib/pq"
)

func main() {
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
