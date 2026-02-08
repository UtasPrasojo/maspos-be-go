package dto

type RegisterRequest struct {
	Name     string `json:"name" example:"Jonatan" binding:"required"`
	Email    string `json:"email" example:"admin@example.com" binding:"required,email"`
	Password string `json:"password" example:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
