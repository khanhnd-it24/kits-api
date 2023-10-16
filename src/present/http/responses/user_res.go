package responses

import (
	"kits/api/src/core/domains"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLoginRes struct {
	Token *Token `json:"token"`
	User  *User  `json:"user"`
}

func UserFromDomain(d *domains.User) *User {
	return &User{
		Id:        d.Id,
		Username:  d.Username,
		FullName:  d.FullName,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
