package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"qttf/internal/models"
	"qttf/internal/rating"

	_ "github.com/lib/pq"
)

type ratingRepository struct {
	db *sql.DB
}

// GetRating implements rating.Repository.
func (r *ratingRepository) GetRating() ([]models.Rating, error) {
	rows, err := r.db.Query(getRating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Rating{}, nil
		}

		return nil, fmt.Errorf("ratingRepository.GetRating was failed: %w", err)
	}

	rating := make([]models.Rating, 0)
	for position := 1; rows.Next(); position++ {
		r := models.Rating{Position: position}
		err = rows.Scan(&r.Id, &r.Player.Id, &r.Player.Name,
			&r.Player.Surname, &r.Player.Hyperlink,
			&r.Player.City.Id, &r.Player.City.Name,
			&r.Player.City.Hyperlink, &r.Rating, &r.UpdateTime)

		if err != nil {
			return nil, fmt.Errorf("ratingRepository.GetRating was failed: %w", err)
		}

		rating = append(rating, r)
	}

	return rating, nil
}

func NewRatingRepository(db *sql.DB) rating.Repository {
	return &ratingRepository{db: db}
}
