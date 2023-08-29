package models

import (
	"database/sql"
	"fmt"
)

type Player struct {
	Id        int
	Surname   string
	Name      string
	Hyperlink string
	City      City
}

func (p Player) ToInsertScript() string {
	return fmt.Sprintf("INSERT INTO player (player_name, player_surname, profile_link, city_id) VALUES('%s', '%s', '%s', (SELECT city_id FROM city WHERE city_name='%s')) ON CONFLICT DO NOTHING;\n", p.Name, p.Surname, p.Hyperlink, p.City.Name)
}

func (p *Player) PlayerActualization(db *sql.DB) error {
	query := `SELECT player_id FROM player WHERE profile_link = $1`

	_, err := db.Exec(p.ToInsertScript())
	if err != nil {
		return fmt.Errorf("Player.playerActualization insert: %w", err)
	}

	err = db.QueryRow(query, p.Hyperlink).Scan(&p.Id)
	if err != nil {
		return fmt.Errorf("Player.playerActualization scan: %w", err)
	}

	return nil
}
