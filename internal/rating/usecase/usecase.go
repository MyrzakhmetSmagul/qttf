package usecase

import (
	"qttf/internal/models"
	"qttf/internal/rating"
)

type ratingUC struct {
	ratingRepo rating.Repository
}

// GetRating implements rating.UseCase.
func (r *ratingUC) GetRating() ([]models.Rating, error) {
	return r.ratingRepo.GetRating()
}

func NewRatingUseCase(ratingRepo rating.Repository) rating.UseCase {
	return &ratingUC{ratingRepo: ratingRepo}
}
