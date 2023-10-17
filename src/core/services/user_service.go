package services

import (
	"context"
	"fmt"
	hashprovider "kits/api/src/common/crypto/hash"
	"kits/api/src/common/fault"
	"kits/api/src/core/domains"
)

type UserService struct {
	userRepo     domains.UserRepo
	hashProvider hashprovider.HashProvider
}

func NewUserService(userRepo domains.UserRepo, hash hashprovider.HashProvider) *UserService {
	return &UserService{
		userRepo:     userRepo,
		hashProvider: hash,
	}
}

func (u *UserService) Create(ctx context.Context, userCreate *domains.UserCreate) (int64, error) {
	caller := "UserService.Create"
	existedUser, err := u.userRepo.FindByUsername(ctx, userCreate.Username)

	if existedUser != nil && existedUser.IsAvailable() {
		wErr := fmt.Errorf("[%v]: conflict name %v", caller, userCreate.Username)
		return -1, fault.Wrap(wErr).SetTag(fault.TagAlreadyExists).SetKey(fault.KeyUserAlreadyExist)
	}

	if err != nil && !fault.IsTag(err, fault.TagNotFound) {
		return -1, fault.Wrapf(err, "[%v] failed to find user", caller)
	}

	hashPw, err := u.hashProvider.Hash(userCreate.Password)
	if err != nil {
		return -1, fault.Wrapf(err, "[%v]: failed to hash password", caller)
	}

	userCreate.Password = hashPw

	user, err := u.userRepo.Save(ctx, userCreate)
	if err != nil {
		return -1, fault.Wrapf(err, "[%v] failed to create user", caller)
	}

	return user.Id, nil
}
