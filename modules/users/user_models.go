package users

type User struct {
	ID   string `json:"id"`
	Body Body   `json:"body"`
}

type Body struct {
	Name     string `json:"name" validate:"required"`
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserLogin struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
}
