package handler

import (
	"fmt"
	"golang-woktopup/config"
	"golang-woktopup/internal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"

	//"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	DB *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{DB: db}
}

func (h *PaymentHandler) CreateMidtransTransaction(c *gin.Context) {
	type RequestBody struct {
		OrderID int `json:"order_id"`
	}

	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	var order model.Order
	if err := h.DB.First(&order, body.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	var user model.User
	if err := h.DB.First(&user, order.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if order.TransactionID == "" {
		order.TransactionID = fmt.Sprintf("ORDER-%d-%d", order.ID, time.Now().Unix())
		h.DB.Save(&order)
	}

	// Snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.TransactionID,
			GrossAmt: int64(order.TotalPrice),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
	}

	snapResp, err := config.MidtransClient.CreateTransaction(snapReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create midtrans transaction", "error": err.Error()})
		return
	}

	// Simpan status pending pada order dan payments
	order.Status = "pending"
	if err := h.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order", "error": err.Error()})
		return
	}

	payment := model.Payment{
		OrderID:       order.ID,
		PaymentMethod: "midtrans",
		TransactionID: order.TransactionID,
		Amount:        order.TotalPrice,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := h.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save payment", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    snapResp.Token,
		"redirect": snapResp.RedirectURL,
	})
}

func (h *PaymentHandler) CreateInvoiceIfNotExists(orderID uint) {
	var count int64
	h.DB.Model(&model.Invoice{}).Where("order_id = ?", orderID).Count(&count)
	if count == 0 {
		invoice := model.Invoice{
			OrderID:     orderID,
			InvoiceCode: fmt.Sprintf("INV-%d-%d", orderID, time.Now().Unix()),
			GeneratedAt: time.Now(),
		}
		h.DB.Create(&invoice)
	}
}

func (h *PaymentHandler) HandleNotification(c *gin.Context) {
	type NotificationBody struct {
		Result         map[string]interface{} `json:"result"`
		OverrideStatus string                 `json:"override_status"` // optional: "pending", "failed"
	}

	var body NotificationBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	internalTrxID, ok := body.Result["order_id"].(string)
	if !ok || internalTrxID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid or missing order_id in result"})
		return
	}

	var order model.Order
	if err := h.DB.Where("transaction_id = ?", internalTrxID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found with transaction_id", "order_id": internalTrxID})
		return
	}

	status := body.OverrideStatus
	if status == "" {
		if s, ok := body.Result["transaction_status"]; ok {
			status = fmt.Sprintf("%v", s)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing transaction_status in result"})
			return
		}
	}

	midtransID := ""
	if txID, ok := body.Result["transaction_id"].(string); ok {
		midtransID = txID
	}

	// Update order
	order.Status = status
	order.MidtransID = midtransID
	h.DB.Save(&order)

	// Update atau buat payment
	var payment model.Payment
	if err := h.DB.Where("order_id = ?", order.ID).First(&payment).Error; err == nil {
		payment.Status = status
		payment.MidtransID = midtransID
		payment.UpdatedAt = time.Now()
		h.DB.Save(&payment)
	} else {
		paymentMethod := "-"
		if pt, ok := body.Result["payment_type"].(string); ok {
			paymentMethod = pt
		}

		payment = model.Payment{
			OrderID:       order.ID,
			PaymentMethod: paymentMethod,
			TransactionID: order.TransactionID,
			MidtransID:    midtransID,
			Amount:        order.TotalPrice,
			Status:        status,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		h.DB.Create(&payment)
	}

	// ✅ BUAT INVOICE jika status settlement
	if status == "settlement" {
		h.CreateInvoiceIfNotExists(order.ID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification handled successfully"})
}
