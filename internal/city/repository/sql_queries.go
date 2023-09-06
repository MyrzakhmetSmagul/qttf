package repository

const (
	createCity = `INSERT INTO city (city_name, city_link) VALUES($1, $2) ON CONFLICT (city_link) DO NOTHING RETURNING *;`

	getCities = `SELECT city_id, city_name, city_link FROM city`

	getById = `SELECT city_name, city_link FROM city WHERE city_id = $1`

	getIdByLink = `SELECT city_id FROM city WHERE city_link = $1;`
)
