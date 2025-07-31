package handler

import (
	"bytes"
	"fmt"
	"net/http"

	"golang-woktopup/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type InvoiceHandler struct {
	DB *gorm.DB
}

func NewInvoiceHandler(db *gorm.DB) *InvoiceHandler {
	return &InvoiceHandler{DB: db}
}

// GetInvoiceDetail handles GET /invoice/:id
func (h *InvoiceHandler) GetInvoiceDetail(c *gin.Context) {
	id := c.Param("id")
	var invoice model.Invoice
	err := h.DB.Preload("Order.User").
		Preload("Order.Product.Game").
		Preload("Order.Voucher").
		Where("order_id = ?", id).
		First(&invoice).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Invoice not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// DownloadInvoicePDF handles GET /invoice/:id/download
func (h *InvoiceHandler) DownloadInvoicePDF(c *gin.Context) {
	id := c.Param("id")
	var invoice model.Invoice
	err := h.DB.Preload("Order.User").
		Preload("Order.Product.Game").
		Preload("Order.Voucher").
		Where("order_id = ?", id).
		First(&invoice).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Invoice not found", "error": err.Error()})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice #%s", invoice.InvoiceCode))

	pdf.SetFont("Arial", "", 12)
	pdf.Ln(12)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", invoice.GeneratedAt.Format("2006-01-02 15:04:05")))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", invoice.Order.User.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Game: %s", invoice.Order.Product.Game.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Product: %s", invoice.Order.Product.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Payment Method: %s", invoice.Order.PaymentMethod))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", invoice.Order.Status))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Amount: Rp %.0f", invoice.Order.TotalPrice))

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate PDF", "error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=invoice-%s.pdf", invoice.InvoiceCode))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}
