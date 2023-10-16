package services

import (
	"context"
	"kits/api/src/common/configs"
	hashprovider "kits/api/src/common/crypto/hash"
	"kits/api/src/core/domains"
)

type AuthService struct {
	userRepo     domains.UserRepo
	hashProvider hashprovider.HashProvider
	cf           *configs.Config
}

func NewAuthService(userRepo domains.UserRepo, hash hashprovider.HashProvider, cf *configs.Config) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		hashProvider: hash,
		cf:           cf,
	}
}

func (a *AuthService) Login(ctx context.Context, authUser *domains.AuthUser) (*domains.Token, *domains.User, error) {
	return nil, nil, nil
}
