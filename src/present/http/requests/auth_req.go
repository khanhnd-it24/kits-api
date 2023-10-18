package requests

type AuthLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthRefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
