package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"qttf/internal/city"
	"qttf/internal/models"

	_ "github.com/lib/pq"
)

type cityRepository struct {
	db *sql.DB
}

func NewCityRepository(db *sql.DB) city.Repository {
	return &cityRepository{db: db}
}

func (c *cityRepository) Create(city *models.City) error {
	_, err := c.db.Exec(createCity, city.Name, city.Hyperlink)
	if err != nil {
		return fmt.Errorf("cityRepository.Create was failed cause: %w", err)
	}

	err = c.db.QueryRow(getId, city.Hyperlink).Scan(&city.Id)
	if err != nil {
		return fmt.Errorf("cityRepository.Create was failed cause: %w", err)
	}

	return nil
}
func (c *cityRepository) GetCities() ([]models.City, error) {
	rows, err := c.db.Query(getCities)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.City{}, nil
		}

		return nil, fmt.Errorf("cityRepository.GetCities was failed cause: %w", err)
	}
	cities := make([]models.City, 0)
	for rows.Next() {
		city := models.City{}
		err = rows.Scan(&city.Id, &city.Name, &city.Hyperlink)
		if err != nil {
			return nil, fmt.Errorf("cityRepository.GetCities was failed.\nrow scan was failed: %w", err)
		}
		cities = append(cities, city)
	}
	return cities, nil
}

func (c *cityRepository) GetById(id int) (models.City, error) {
	city := models.City{Id: id}
	err := c.db.QueryRow(getById, id).Scan(&city.Name, &city.Hyperlink)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return city, fmt.Errorf("cityRepository.GetById was failed cause: %w", err)
	}

	return city, nil
}

func (c *cityRepository) GetIdByLink(link string) (int, error) {
	var id int
	err := c.db.QueryRow(getById, id).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cityRepository.GetIdByLink was failed cause: %w", err)
	}

	return id, nil
}
