package repository

import (
	"database/sql"
	"fmt"
	"qttf/internal/models"
	"qttf/internal/player"

	_ "github.com/lib/pq"
)

type playerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) player.Repository {
	return &playerRepository{
		db: db,
	}
}

func (p *playerRepository) Create(player *models.Player) error {
	err := p.db.QueryRow(create, player.Name, player.Surname, player.Hyperlink, player.Id).Scan(&player.Id)
	if err != nil {
		return fmt.Errorf("playerRepository.Create was failed cause: %w", err)
	}
	return nil
}

func (p *playerRepository) GetPlayers() ([]models.Player, error) {
	rows, err := p.db.Query(getPlayers)
	if err != nil {
		return nil, fmt.Errorf("playerRepository.GetPlayers was failed cause: %w", err)
	}

	players := make([]models.Player, 0)
	for rows.Next() {
		player := models.Player{}
		err = rows.Scan(&player.Id, &player.Name, &player.Surname, &player.Hyperlink, &player.City.Id)
		if err != nil {
			return nil, fmt.Errorf("playerRepository.GetPlayers, row scanning was failed: %w", err)
		}

		players = append(players, player)
	}

	return players, nil
}

func (p *playerRepository) GetById(id int) (models.Player, error) {
	player := models.Player{Id: id}
	err := p.db.QueryRow(getPlayerById, id).Scan(&player.Name, &player.Surname, &player.Hyperlink, &player.City.Id)
	if err != nil {
		return player, fmt.Errorf("playerRepository.GetById was failed: %w", err)
	}

	return player, nil
}
