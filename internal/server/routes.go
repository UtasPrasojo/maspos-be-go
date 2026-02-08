package server

import (
	"net/http"

	_ "maspos-be-go/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// ===== BASIC ROUTES =====
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	r.Static("/uploads", "./uploads")

	// ===== SWAGGER ROUTE =====
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := r.Group("/auth")
	{

		auth.POST("/register", s.RegisterHandler)
		auth.POST("/login", s.LoginHandler)
	
	}
	cat := r.Group("/categories")
    {
        cat.POST("", s.CreateCategoryHandler)      // Create
        cat.GET("", s.GetAllCategoriesHandler)     // Get All
        cat.GET("/:id", s.GetCategoryByIDHandler)  // Get By ID
        cat.PATCH("/:id", s.UpdateCategoryHandler) // Update
        cat.DELETE("/:id", s.DeleteCategoryHandler) // Delete
    }
	prod := r.Group("/products")
    {
        prod.POST("", s.CreateProductHandler)       // Create
        prod.GET("", s.GetAllProductsHandler)      // Read All
        prod.GET("/:id", s.GetProductByIDHandler)  // Read One
        prod.PATCH("/:id", s.UpdateProductHandler) // Update
        prod.DELETE("/:id", s.DeleteProductHandler)// Delete
    }
	return r
}



func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
