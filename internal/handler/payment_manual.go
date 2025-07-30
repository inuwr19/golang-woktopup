package handler

import (
	"fmt"
	"golang-woktopup/internal/model"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
)

type ManualPaymentHandler struct {
	DB *gorm.DB
}

type ManualPaymentNotification struct {
	Result         map[string]interface{} `json:"result"`
	OverrideStatus string                 `json:"override_status"` // optional: "pending", "failed"
}

func NewManualPaymentHandler(db *gorm.DB) *ManualPaymentHandler {
	return &ManualPaymentHandler{DB: db}
}

func (h *ManualPaymentHandler) CheckPaymentStatus(c *gin.Context) {
	var req struct {
		TransactionID string `json:"transaction_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.TransactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Transaction ID is required"})
		return
	}

	// Inisialisasi coreapi client
	core := coreapi.Client{}
	core.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	// Check transaksi dari Midtrans
	transactionStatusResp, err := core.CheckTransaction(req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to check transaction", "error": err.Error()})
		return
	}

	// Temukan order berdasarkan transaction_id
	var order model.Order
	if err := h.DB.Where("transaction_id = ?", req.TransactionID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	status := string(transactionStatusResp.TransactionStatus)

	if status == "settlement" || status == "capture" {
		order.Status = "paid"
		h.DB.Save(&order)

		// Cek dan simpan payment jika belum ada
		var existing model.Payment
		if err := h.DB.Where("order_id = ?", order.ID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			payment := model.Payment{
				OrderID:       order.ID,
				PaymentMethod: order.PaymentMethod,
				TransactionID: order.TransactionID,
				Amount:        order.TotalPrice,
				Status:        "success",
			}
			h.DB.Create(&payment)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Payment checked",
		"transaction_id": req.TransactionID,
		"status":         status,
		"order_id":       order.ID,
	})
}

func (h *ManualPaymentHandler) HandleMidtransNotification(c *gin.Context) {
	var body ManualPaymentNotification

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	transactionID := fmt.Sprintf("%v", body.Result["order_id"])
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order_id in result"})
		return
	}

	var order model.Order
	if err := h.DB.Where("transaction_id = ?", transactionID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	status := body.OverrideStatus
	if status == "" {
		status = fmt.Sprintf("%v", body.Result["transaction_status"])
	}

	order.Status = status
	h.DB.Save(&order)

	// Buat atau update payment
	payment := model.Payment{
		OrderID:       order.ID,
		PaymentMethod: order.PaymentMethod,
		TransactionID: transactionID,
		Amount:        order.TotalPrice,
		Status:        status,
	}
	h.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated manually"})
}
