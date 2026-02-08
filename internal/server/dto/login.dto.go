package dto

type LoginRequest struct {
    Email    string `json:"email" example:"admin@example.com" binding:"required,email"`
    Password string `json:"password" example:"password" binding:"required"`
}

type LoginResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
    Email string `json:"email" example:"admin@example.com"`
}