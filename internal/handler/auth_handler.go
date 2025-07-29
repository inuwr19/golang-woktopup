package handler

import (
	"golang-woktopup/internal/model"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Name                 string `json:"name"`
		Email                string `json:"email"`
		Phone                string `json:"phone_number"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	if req.Password != req.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password mismatch"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := model.User{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.Phone,
		Role:        "user",
		Password:    string(hashedPassword),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		log.Println("Register error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registered successfully"})
}

// POST /auth/spa/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var user model.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// GET /user
func (h *AuthHandler) CurrentUser(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthenticated"})
		return
	}

	var user model.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
