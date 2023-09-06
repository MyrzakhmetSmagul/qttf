package repository

const (
	create = `INSERT INTO player(player_name, player_surname, profile_link, city_id) VALUES($1, $2, $3, $4) ON CONFLICT (profile_link) DO NOTHING RETURNING player_id;`

	getPlayers = `SELECT player_id, player_name, player_surname, profile_link, city_id FROM player;`

	getPlayerById = `SELECT player_name, player_surname, profile_link, city_id FROM player WHERE player_id = $1;`
)
