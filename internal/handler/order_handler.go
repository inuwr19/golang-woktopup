package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"golang-woktopup/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB *gorm.DB
}

func generateTransactionID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("TRX-%d", rand.Intn(1_000_000_000))
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input model.Order

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	input.TransactionID = generateTransactionID() // Pindah ke sini
	input.Status = "pending"
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	if err := h.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order", "error": err.Error()})
		return
	}

	var user model.User
	if err := h.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	id := c.Param("id")

	var order model.Order
	err := h.DB.
		Preload("User").
		Preload("Product").
		Preload("Product.Game").
		Preload("Voucher").
		First(&order, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}
