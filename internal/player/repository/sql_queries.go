package repository

const (
	create = `INSERT INTO player(player_name, player_surname, profile_link, city_id) VALUES($1, $2, $3, $4) ON CONFLICT (profile_link) DO NOTHING RETURNING player_id;`

	getPlayers = `SELECT p.player_id, p.player_name, p.player_surname, p.profile_link, c.city_id, c.city_name, c.city_link
	FROM player p 
	INNER JOIN city c ON c.city_id = p.city_id;`

	getPlayerById = `SELECT player_name, player_surname, profile_link, city_id FROM player WHERE player_id = $1;`
)
