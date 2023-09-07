package auth

type UseCase interface {
	SaveGoogleToken(code string) error
	GetGoogleToken() string
}
