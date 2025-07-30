package handler

import (
	"golang-woktopup/internal/model"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VoucherHandler struct {
	DB *gorm.DB
}

func NewVoucherHandler(db *gorm.DB) *VoucherHandler {
	return &VoucherHandler{DB: db}
}

func (h *VoucherHandler) ApplyVoucher(c *gin.Context) {
	var request struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.Code) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid voucher code"})
		return
	}

	var voucher model.Voucher
	if err := h.DB.Where("code = ? AND status = ?", request.Code, "active").First(&voucher).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Voucher not found or inactive"})
		return
	}

	if time.Now().After(voucher.ExpiryDate) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Voucher has expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       voucher.ID,
		"discount": voucher.Discount,
	})
}
