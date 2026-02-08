package dto

type CategoryRequest struct {
    Name string `json:"name" example:"Makanan" binding:"required"`
}

type CategoryResponse struct {
    ID   string `json:"id" example:"uuid-string-123"`
    Name string `json:"name" example:"Makanan"`
}