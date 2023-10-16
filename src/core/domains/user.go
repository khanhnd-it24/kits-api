package domains

import (
	"context"
	"time"
)

type User struct {
	Id        int64
	Username  string
	Password  string
	FullName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (u *User) IsAvailable() bool {
	return u.DeletedAt == nil
}

type UserCreate struct {
	Username string
	Password string
	FullName string
}

type AuthUser struct {
	Username string
	Password string
}

type UserRepo interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindById(ctx context.Context, id string) (*User, error)
	FindByIdAndDelete(ctx context.Context, id string) error
	Save(ctx context.Context, user *UserCreate) (*User, error)
}
