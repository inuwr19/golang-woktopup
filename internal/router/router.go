package router

import (
	"time"

	"golang-woktopup/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Middleware CORS untuk izinkan akses dari frontend React (localhost:5173)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-XSRF-TOKEN"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware session dengan cookie store
	store := cookie.NewStore([]byte("secret")) // Ganti secret key dengan yang lebih aman di production
	r.Use(sessions.Sessions("woktopup_session", store))

	// Kelompokkan semua endpoint ke dalam prefix /api
	api := r.Group("/api")
	{
		// Endpoint CSRF dummy (untuk React + Axios)
		api.GET("/csrf-cookie", func(c *gin.Context) {
			c.Status(204)
		})

		// Handler auth (register, login, logout, get user)
		auth := handler.NewAuthHandler(db)
		api.POST("/auth/register", auth.Register)
		api.POST("/auth/login", auth.Login)
		api.POST("/auth/logout", auth.Logout)
		api.GET("/user", auth.CurrentUser)

		// Handler game (get available games)
		game := handler.NewGameHandler(db)
		api.GET("/games", game.GetGames)
		api.GET("/games/:id", game.GetGameDetail)

		// Handler product (get products by game)
		order := handler.NewOrderHandler(db)
		api.GET("/transactions/latest", order.GetLatestTransactions)
		api.POST("/orders", order.CreateOrder)
		api.GET("/orders/:id", order.GetOrderDetail)
		api.GET("/invoice-by-order/:order_id", order.GetInvoiceByOrder)

		invoiceHandler := handler.NewInvoiceHandler(db)
		invoices := api.Group("/invoices")
		{
			invoices.GET("/:id", invoiceHandler.GetInvoiceDetail)
			invoices.GET("/:id/download", invoiceHandler.DownloadInvoicePDF)
		}

		// Handler order (create order)
		voucher := handler.NewVoucherHandler(db)
		api.POST("/voucher/apply", voucher.ApplyVoucher)

		// Handler payment (create Midtrans transaction)
		payment := handler.NewPaymentHandler(db)
		api.POST("/payment/create", payment.CreateMidtransTransaction)
		api.POST("/payments/handle-notification", payment.HandleNotification)

		// Handler manual payment (check payment status)
		manualPayment := handler.NewManualPaymentHandler(db)
		api.POST("/payments/manual-check", manualPayment.CheckPaymentStatus)

	}

	return r
}
