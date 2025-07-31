package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"golang-woktopup/internal/model"

	"github.com/gin-contrib/sessions"
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

func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var orders []model.Order
	err := h.DB.
		Preload("Product.Game").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetInvoiceByOrder(c *gin.Context) {
	orderIDParam := c.Param("order_id")
	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}

	var invoice model.Invoice
	err = h.DB.Preload("Order.Product.Game").
		Preload("Order.User").
		Where("order_id = ?", orderID).
		First(&invoice).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Invoice not found for this order"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

func (h *OrderHandler) GetLatestTransactions(c *gin.Context) {
	var orders []model.Order
	err := h.DB.
		Preload("User").
		Preload("Product").
		Preload("Product.Game").
		Where("status = ?", "settlement").
		Order("created_at DESC").
		Limit(10).
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch transactions"})
		return
	}

	var result []map[string]interface{}
	for _, o := range orders {
		result = append(result, map[string]interface{}{
			"id":        o.ID,
			"game":      o.Product.Game.Name,
			"amount":    o.Product.Name,
			"user":      o.User.Name,
			"createdAt": o.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}
