package models

import (
	"kits/api/src/core/domains"
	"time"
)

type User struct {
	Id        int64      `gorm:"column:id"`
	Username  string     `gorm:"column:username"`
	Password  string     `gorm:"column:password"`
	FullName  string     `gorm:"column:full_name"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (u *User) ToDomain() *domains.User {
	return &domains.User{
		Id:        u.Id,
		Username:  u.Username,
		Password:  u.Password,
		FullName:  u.FullName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.CreatedAt,
		DeletedAt: u.DeletedAt,
	}
}

func UserCreateFromDomain(u *domains.UserCreate) *User {
	return &User{
		Username:  u.Username,
		Password:  u.Password,
		FullName:  u.FullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) TableName() string {
	return "users"
}
