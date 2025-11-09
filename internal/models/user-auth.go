package models

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Town     string `json:"town"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
