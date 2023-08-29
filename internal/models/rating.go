package models

import "fmt"

type Rating struct {
	Id         int
	Position   int
	Player     Player
	Rating     int
	UpdateTime string
}

func (r Rating) ToInsertScript() string {
	return fmt.Sprintf("INSERT INTO rating(player_id, rating, last_update) VALUES( %d, %d, '%s') ON CONFLICT DO NOTHING;\n", r.Player.Id, r.Rating, r.UpdateTime)
}
