package player

import "qttf/internal/models"

type UseCase interface {
	GetPlayers() ([]models.Player, error)
	GetById(id int) (models.Player, error)
}
