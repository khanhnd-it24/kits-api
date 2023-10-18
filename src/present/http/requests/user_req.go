package requests

type UserCreateReq struct {
	FullName string `json:"full_name" binding:"required,max=100"`
	Username string `json:"username" binding:"required,max=30,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=30,strong_password"`
}
