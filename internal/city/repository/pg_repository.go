package repository

import (
	"database/sql"
	"fmt"
	"qttf/internal/city"
	"qttf/internal/models"
)

type cityRepository struct {
	db *sql.DB
}

func NewCityRepository(db *sql.DB) city.Repository {
	return &cityRepository{db: db}
}

func (c *cityRepository) Create(city *models.City) error {
	err := c.db.QueryRow(createCity, city.Name, city.Hyperlink).Scan(&city.Id)
	if err != nil {
		return fmt.Errorf("cityRepository.Create was failed cause: %w", err)
	}

	return nil
}

func (c *cityRepository) GetById(id int) (models.City, error) {
	city := models.City{Id: id}
	err := c.db.QueryRow(getById, id).Scan(&city.Name, &city.Hyperlink)
	if err != nil {
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
