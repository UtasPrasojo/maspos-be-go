package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"maspos-be-go/internal/database/repository"
	"maspos-be-go/internal/server/dto"
	"maspos-be-go/internal/utils"
)

// Register user
// @Summary Register new user
// @Description Create new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "Register payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
func (s *Server) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	repo := repository.NewUserRepository(s.db.DB())
	ctx := c.Request.Context()

	// cek email
	exists, err := repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "email already registered",
		})
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to hash password",
		})
		return
	}

	// insert user
	if err := repo.Create(ctx, req.Name, req.Email, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully",
		"data": dto.RegisterResponse{
			Name:  req.Name,
			Email: req.Email,
		},
	})
}
// Login user
// @Summary Login user
// @Description Authenticate user and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (s *Server) LoginHandler(c *gin.Context) {
    var req dto.LoginRequest

    // Validasi input berdasarkan tag 'binding' di DTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": err.Error(),
        })
        return
    }

    // Panggil repository untuk mencari user berdasarkan email
    repo := repository.NewUserRepository(s.db.DB())
    ctx := c.Request.Context()
    
    user, err := repo.GetByEmail(ctx, req.Email) // Pastikan method GetByEmail ada di repository
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "invalid email or password",
        })
        return
    }

    // Verifikasi password (membandingkan input plain text dengan hash di DB)
    if err := utils.CheckPassword(req.Password, user.Password); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "invalid email or password",
        })
        return
    }

   token, err := utils.GenerateToken(user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to generate token",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Login successful",
        "data": dto.LoginResponse{
            Token: token, // Sekarang berisi token JWT asli
            Email: user.Email,
        },
    })
}
