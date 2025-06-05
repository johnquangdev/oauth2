package usecase

import (
	imlp "github.com/johnquangdev/oauth2/usecase/imlp"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
)

type imlpUseCase struct{}

func (u imlpUseCase) Auth() interfaces.Auth {
	return imlp.NewOAuthUsecase()
}

func NewUseCase() (interfaces.UseCase, error) {
	return &imlpUseCase{}, nil
}
