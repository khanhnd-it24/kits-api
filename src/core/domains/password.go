package domains

import (
	"context"
	"time"
)

type Password struct {
	Id        int64
	UserId    int64
	Name      string
	Password  string
	Note      *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type PasswordRepo interface {
	Save(ctx context.Context, password *Password) (*Password, error)
	FindByUserId(ctx context.Context, userId int64) ([]*Password, error)
}
