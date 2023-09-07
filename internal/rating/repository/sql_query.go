package repository

const (
	create    = `INSERT INTO rating(player_id, rating, last_update) VALUES($1, $2, $3) ON CONFLICT (player_id) DO NOTHING RETURNING rating_id;`
	getRating = `SELECT r.rating_id, p.player_id, p.player_name, p.player_surname, p.profile_link, c.city_id, c.city_name, c.city_link, r.rating, r.last_update
	FROM rating r 
	INNER JOIN player p ON p.player_id = r.player_id
	INNER JOIN city c ON c.city_id = p.city_id
	ORDER BY r.rating DESC`
)
