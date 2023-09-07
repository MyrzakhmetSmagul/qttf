package usecase

import (
	"context"
	"qttf/config"
	"qttf/internal/auth"
	"qttf/pkg/jsonsaver"

	"golang.org/x/oauth2"
)

type authUC struct {
	cnf *config.Config
}

// GetGoogleToken implements auth.UseCase.
func (a *authUC) GetGoogleToken() string {
	return a.cnf.GoogleCOnfig.AuthCodeURL("stateToken", oauth2.AccessTypeOffline)
}

// SaveGoogleToken implements auth.UseCase.
func (a *authUC) SaveGoogleToken(code string) error {
	tok, err := a.cnf.GoogleCOnfig.Exchange(context.TODO(), code)
	if err != nil {
		return err
	}

	return jsonsaver.SaveJson(a.cnf.TokenPath, tok)
}

func NewAuthUseCase(cnf *config.Config) auth.UseCase {
	return &authUC{cnf: cnf}
}
