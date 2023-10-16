package models

import (
	"kits/api/src/core/domains"
	"time"
)

type RefreshToken struct {
	Id        int64     `gorm:"column:id"`
	UserId    int64     `gorm:"column:user_id"`
	Token     string    `gorm:"column:token"`
	ExpiredAt time.Time `gorm:"column:expired_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (r *RefreshToken) ToDomain() *domains.RefreshToken {
	return &domains.RefreshToken{
		Id:        r.Id,
		UserId:    r.UserId,
		Token:     r.Token,
		ExpiredAt: r.ExpiredAt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.CreatedAt,
	}
}

func RefreshTokenFromDomain(r *domains.RefreshToken) *RefreshToken {
	return &RefreshToken{
		Id:        r.Id,
		UserId:    r.UserId,
		Token:     r.Token,
		ExpiredAt: r.ExpiredAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (r *RefreshToken) TableName() string {
	return "refresh_tokens"
}
