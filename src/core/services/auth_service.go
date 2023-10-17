package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"kits/api/src/common/configs"
	"kits/api/src/common/crypto/aes"
	hashprovider "kits/api/src/common/crypto/hash"
	"kits/api/src/common/fault"
	"kits/api/src/core/domains"
	"time"
)

type AuthService struct {
	userRepo           domains.UserRepo
	refreshTokenRepo   domains.RefreshTokenRepo
	accessTokenCache   domains.AccessTokenCache
	hashProvider       hashprovider.HashProvider
	cf                 *configs.Config
	refreshAesProvider *aes.GcmProvider
	refreshExpire      time.Duration
	accessAesProvider  *aes.GcmProvider
	accessExpire       time.Duration
}

func NewAuthService(userRepo domains.UserRepo, refreshTokenRepo domains.RefreshTokenRepo, accessTokenCache domains.AccessTokenCache, hash hashprovider.HashProvider, cf *configs.Config) *AuthService {
	refreshKey := cf.Aes["refresh_token"].Key
	refreshExpire := cf.Aes["refresh_token"].Expire
	refreshAesProvider, _ := aes.NewAesGcmProvider(refreshKey)

	accessKey := cf.Aes["access_token"].Key
	accessExpire := cf.Aes["access_token"].Expire
	accessAesProvider, _ := aes.NewAesGcmProvider(accessKey)

	return &AuthService{
		userRepo:           userRepo,
		refreshTokenRepo:   refreshTokenRepo,
		accessTokenCache:   accessTokenCache,
		hashProvider:       hash,
		cf:                 cf,
		refreshAesProvider: refreshAesProvider,
		refreshExpire:      refreshExpire,
		accessAesProvider:  accessAesProvider,
		accessExpire:       accessExpire,
	}
}

func (a *AuthService) Login(ctx context.Context, authUser *domains.AuthUser) (*domains.Token, *domains.User, error) {
	caller := "AuthService.Login"
	user, err := a.userRepo.FindByUsername(ctx, authUser.Username)

	if fault.IsTag(err, fault.TagNotFound) {
		return nil, nil, fault.Wrapf(err, "[%v] username %s not found", caller, authUser.Username).
			SetTag(fault.TagUnauthenticated).
			SetKey(fault.KeyAuthInvalidIdentify)
	}

	if err != nil {
		return nil, nil, fault.Wrapf(err, "[%v] failed to find user", caller)
	}

	err = a.hashProvider.ComparePassword(authUser.Password, user.Password)
	if err != nil {
		return nil, nil, fault.Wrapf(err, "[%v] password not same", caller).
			SetTag(fault.TagUnauthenticated).
			SetKey(fault.KeyAuthInvalidIdentify)
	}

	token, err := a.generateToken(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return token, user, nil
}

func (a *AuthService) generateToken(ctx context.Context, user *domains.User) (*domains.Token, error) {
	caller := "AuthService.generateToken"
	refreshTokenString, refreshExpiredAt, err := a.generateAesToken(user, a.refreshAesProvider, a.refreshExpire)
	refreshToken := &domains.RefreshToken{
		UserId:    user.Id,
		Token:     refreshTokenString,
		ExpiredAt: refreshExpiredAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	refreshToken, err = a.refreshTokenRepo.Save(ctx, refreshToken)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to save refresh token", caller)
	}
	accessTokenString, _, err := a.generateAesToken(user, a.accessAesProvider, a.accessExpire)
	accessToken := &domains.AccessToken{
		Token:      accessTokenString,
		ExpireTime: a.accessExpire,
	}

	err = a.accessTokenCache.Save(ctx, user.Id, accessToken)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to save access token", caller)
	}

	return &domains.Token{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    int64(accessToken.ExpireTime.Seconds()),
	}, nil
}

func (a *AuthService) generateAesToken(user *domains.User, provider *aes.GcmProvider, expireTime time.Duration) (string, time.Time, error) {
	caller := "AuthService.generateRefreshToken"

	expiredAt := time.Now().Add(expireTime)
	token := &domains.TokenAes{
		UserId:    user.Id,
		ExpiredAt: expiredAt.Format(time.RFC3339),
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		wErr := fault.Wrapf(err, "[%v] failed to marshal token %s", caller, token).SetTag(fault.TagInternal)
		return "", expiredAt, wErr
	}

	cipherTokenBytes, err := provider.Seal(tokenBytes)
	if err != nil {
		wErr := fault.Wrapf(err, "[%v] failed to seal token %s", caller, token).SetTag(fault.TagInternal)
		return "", expiredAt, wErr
	}

	return base64.StdEncoding.EncodeToString(cipherTokenBytes), expiredAt, nil
}
