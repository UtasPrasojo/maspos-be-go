package dto

import "mime/multipart"

type ProductRequest struct {
	CategoryID string                `form:"category_id" binding:"required"`
	Name       string                `form:"name" binding:"required"`
	Price      float64               `form:"price" binding:"required"`
	Picture    *multipart.FileHeader `form:"picture" binding:"required"`
}

type ProductResponse struct {
	ID         string  `json:"id"`
	CategoryID string  `json:"category_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Picture    string  `json:"picture"`
}