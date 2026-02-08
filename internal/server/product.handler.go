package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"maspos-be-go/internal/database/repository"
	"maspos-be-go/internal/server/dto"
)

// @Summary Create new product
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Param category_id formData string true "Category ID"
// @Param name formData string true "Product Name"
// @Param price formData number true "Price"
// @Param picture formData file true "Product Picture"
// @Success 201 {object} dto.ProductResponse
// @Router /products [post]
func (s *Server) CreateProductHandler(c *gin.Context) {
	var req dto.ProductRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Tangani Upload File
	file := req.Picture
	// Buat nama file unik
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), filepath.Ext(file.Filename))
	dst := filepath.Join("uploads", filename)
	
	// Simpan file ke folder 'uploads' (pastikan folder ini sudah dibuat)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// 2. Simpan ke Database
	repo := repository.NewProductRepository(s.db.DB())
	id, err := repo.Create(c.Request.Context(), req.CategoryID, req.Name, req.Price, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ProductResponse{
		ID:         id,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Price:      req.Price,
		Picture:    dst,
	})
}
// @Summary Get all products
// @Tags Product
// @Produce json
// @Success 200 {array} repository.Product
// @Router /products [get]
func (s *Server) GetAllProductsHandler(c *gin.Context) {
	repo := repository.NewProductRepository(s.db.DB())
	products, err := repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// @Summary Get product by ID
// @Tags Product
// @Param id path string true "Product ID"
// @Success 200 {object} repository.Product
// @Router /products/{id} [get]
func (s *Server) GetProductByIDHandler(c *gin.Context) {
	id := c.Param("id")
	repo := repository.NewProductRepository(s.db.DB())
	product, err := repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// @Summary Update product
// @Tags Product
// @Accept multipart/form-data
// @Param id path string true "Product ID"
// @Param category_id formData string true "Category ID"
// @Param name formData string true "Product Name"
// @Param price formData number true "Price"
// @Param picture formData file false "Product Picture"
// @Success 200 {object} map[string]string
// @Router /products/{id} [patch]
func (s *Server) UpdateProductHandler(c *gin.Context) {
	id := c.Param("id")
	var req dto.ProductRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo := repository.NewProductRepository(s.db.DB())
	
	// Ambil data lama untuk mendapatkan path gambar lama jika tidak ada upload baru
	oldProduct, err := repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	picturePath := oldProduct.Picture
	file, _ := c.FormFile("picture")
	if file != nil {
		filename := fmt.Sprintf("%d%s", time.Now().Unix(), filepath.Ext(file.Filename))
		picturePath = filepath.Join("uploads", filename)
		c.SaveUploadedFile(file, picturePath)
	}

	err = repo.Update(c.Request.Context(), id, req.CategoryID, req.Name, req.Price, picturePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// @Summary Delete product
// @Tags Product
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Router /products/{id} [delete]
func (s *Server) DeleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	repo := repository.NewProductRepository(s.db.DB())
	if err := repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}