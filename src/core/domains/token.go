package domains

import (
	"context"
	"time"
)

type AccessToken struct {
	Token      string
	ExpireTime time.Duration
}

type TokenAes struct {
	UserId    int64
	ExpiredAt string
}

type AccessTokenCache interface {
	Save(ctx context.Context, userId int64, token *AccessToken) error
	FindByUserId(ctx context.Context, userId int64) (*AccessToken, error)
}

type RefreshToken struct {
	Id        int64
	UserId    int64
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshTokenRepo interface {
	Save(ctx context.Context, token *RefreshToken) (*RefreshToken, error)
	FindByUserIdAndDelete(ctx context.Context, userId int64) error
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
