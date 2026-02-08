package server

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "maspos-be-go/internal/database/repository"
    "maspos-be-go/internal/server/dto"
)

// @Summary Create new category
// @Tags Category
// @Accept json
// @Produce json
// @Param body body dto.CategoryRequest true "Category payload"
// @Success 201 {object} dto.CategoryResponse
// @Router /categories [post]
func (s *Server) CreateCategoryHandler(c *gin.Context) {
    var req dto.CategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    repo := repository.NewCategoryRepository(s.db.DB())
    id, err := repo.Create(c.Request.Context(), req.Name)
    
    if err != nil {
        // PERHATIKAN DI SINI: Kita kirim err.Error() asli dari database
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(), 
        })
        return
    }
    
    c.JSON(http.StatusCreated, dto.CategoryResponse{ID: id, Name: req.Name})
}

// @Summary Get all categories
// @Tags Category
// @Produce json
// @Success 200 {array} repository.Category
// @Router /categories [get]
func (s *Server) GetAllCategoriesHandler(c *gin.Context) {
    repo := repository.NewCategoryRepository(s.db.DB())
    res, err := repo.GetAll(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, res)
}

// @Summary Get category by ID
// @Tags Category
// @Param id path string true "Category ID"
// @Success 200 {object} repository.Category
// @Router /categories/{id} [get]
func (s *Server) GetCategoryByIDHandler(c *gin.Context) {
    id := c.Param("id")
    repo := repository.NewCategoryRepository(s.db.DB())
    res, err := repo.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        return
    }
    c.JSON(http.StatusOK, res)
}

// @Summary Update category
// @Tags Category
// @Param id path string true "Category ID"
// @Param body body dto.CategoryRequest true "Category Name"
// @Success 200 {object} map[string]string
// @Router /categories/{id} [patch]
func (s *Server) UpdateCategoryHandler(c *gin.Context) {
    id := c.Param("id")
    var req dto.CategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    repo := repository.NewCategoryRepository(s.db.DB())
    if err := repo.Update(c.Request.Context(), id, req.Name); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Category updated"})
}

// @Summary Delete category
// @Tags Category
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Router /categories/{id} [delete]
func (s *Server) DeleteCategoryHandler(c *gin.Context) {
    id := c.Param("id")
    repo := repository.NewCategoryRepository(s.db.DB())
    if err := repo.Delete(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}