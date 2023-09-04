package city

import "context"

type UseCase interface {
	GetCities(ctx context.Context)
}
