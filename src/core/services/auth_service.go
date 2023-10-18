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

	token, err := a.generateToken(ctx, user.Id)
	if err != nil {
		return nil, nil, err
	}

	return token, user, nil
}

func (a *AuthService) generateToken(ctx context.Context, userId int64, refreshTokenIds ...int64) (*domains.Token, error) {
	caller := "AuthService.generateToken"
	refreshTokenString, refreshExpiredAt, err := a.generateAesToken(userId, a.refreshAesProvider, a.refreshExpire)
	refreshToken := &domains.RefreshToken{
		UserId:    userId,
		Token:     refreshTokenString,
		ExpiredAt: refreshExpiredAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if len(refreshTokenIds) > 0 {
		refreshToken.Id = refreshTokenIds[0]
	}

	refreshToken, err = a.refreshTokenRepo.Save(ctx, refreshToken)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to save refresh token", caller)
	}
	accessTokenString, _, err := a.generateAesToken(userId, a.accessAesProvider, a.accessExpire)
	accessToken := &domains.AccessToken{
		Token:      accessTokenString,
		ExpireTime: a.accessExpire,
	}

	err = a.accessTokenCache.Save(ctx, userId, accessToken)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to save access token", caller)
	}

	return &domains.Token{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    int64(accessToken.ExpireTime.Seconds()),
	}, nil
}

func (a *AuthService) generateAesToken(userId int64, provider *aes.GcmProvider, expireTime time.Duration) (string, time.Time, error) {
	caller := "AuthService.generateRefreshToken"

	expiredAt := time.Now().Add(expireTime)
	token := &domains.TokenAes{
		UserId:    userId,
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

func (a *AuthService) RefreshToken(ctx context.Context, authRefreshToken *domains.AuthRefreshToken) (*domains.Token, error) {
	caller := "AuthService.RefreshToken"
	refreshToken, err := a.refreshTokenRepo.FindByToken(ctx, authRefreshToken.Token)

	if fault.IsTag(err, fault.TagNotFound) {
		return nil, fault.Wrapf(err, "[%v] token %s invalid", caller, authRefreshToken.Token).
			SetTag(fault.TagUnauthenticated).
			SetKey(fault.KeyAuthInvalidToken)
	}

	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to find refresh token", caller)
	}

	err = a.checkRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	token, err := a.generateToken(ctx, refreshToken.UserId, refreshToken.Id)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *AuthService) checkRefreshToken(refreshToken *domains.RefreshToken) error {
	caller := "AuthService.checkRefreshToken"
	refreshTokenDecode, err := base64.StdEncoding.DecodeString(refreshToken.Token)
	if err != nil {
		return fault.Wrapf(err, "[%v] failed to decode refresh token", caller).SetTag(fault.TagInternal)
	}
	byteRefreshTokenAes, err := a.refreshAesProvider.Open(refreshTokenDecode)
	if err != nil {
		return fault.Wrapf(err, "[%v] failed to open refresh token", caller).SetTag(fault.TagInternal)
	}
	var refreshTokenAes domains.TokenAes
	if err := json.Unmarshal(byteRefreshTokenAes, &refreshTokenAes); err != nil {
		return fault.Wrapf(err, "[%v] failed to unmarshal refresh token", caller).SetTag(fault.TagInternal)
	}

	if refreshTokenAes.UserId != refreshToken.UserId {
		return fault.Wrapf(err, "[%v] user id not compare", caller).
			SetTag(fault.TagUnauthenticated).SetKey(fault.KeyAuthInvalidToken)
	}

	if refreshToken.ExpiredAt.Before(time.Now()) {
		return fault.Wrapf(err, "[%v] token expire", caller).
			SetTag(fault.TagUnauthenticated).SetKey(fault.KeyAuthTokenExpire)
	}

	return nil
}
