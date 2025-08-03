package seed

import (
	"time"

	"golang-woktopup/internal/model"

	"gorm.io/gorm"
)

func SeedGamesAndProducts(db *gorm.DB) {
	var count int64
	db.Model(&model.Game{}).Count(&count)
	db.Model(&model.Voucher{}).Count(&count)
	if count > 0 {
		return // jangan duplikat data
	}

	// ⛳ Games
	games := []model.Game{
		{
			Name:        "Mobile Legends",
			Description: "Mobile Legends: Bang Bang is multiplayer online battle arena.",
			Image:       "https://mrwallpaper.com/images/hd/mobile-legends-logo-on-blue-background-2m7kdws7prybsvmt.jpg",
			Status:      "available",
		},
		{
			Name:        "Genshin Impact",
			Description: "Genshin Impact is an open-world, action role-playing game.",
			Image:       "https://i.ytimg.com/vi/fIuhg0bjvQI/maxresdefault.jpg",
			Status:      "available",
		},
		{
			Name:        "PUBG Mobile",
			Description: "PUBG MOBILE is the original battle royale game on mobile.",
			Image:       "https://www.charlieintel.com/cdn-image/wp-content/uploads/2021/01/53f7d969762523414600a6549e36202a.jpg?width=1200&quality=60&format=auto",
			Status:      "available",
		},
	}
	db.Create(&games)

	// ⛳ Products
	products := []model.Product{
		{GameID: 1, Name: "100 Diamonds", Price: 15000, Description: "Product description.", Status: "available"},
		{GameID: 1, Name: "500+100 Diamonds", Price: 45000, Description: "Product description.", Status: "available"},
		{GameID: 2, Name: "60 Genesis Crystals", Price: 35000, Description: "Product description.", Status: "available"},
		{GameID: 2, Name: "300+30 Genesis Crystals", Price: 65000, Description: "Product description.", Status: "available"},
		{GameID: 3, Name: "75 UC", Price: 75000, Description: "Product description.", Status: "available"},
		{GameID: 3, Name: "260+25 UC", Price: 125000, Description: "Product description.", Status: "available"},
	}

	for i := range products {
		products[i].CreatedAt = time.Now()
		products[i].UpdatedAt = time.Now()
	}
	db.Create(&products)

	vouchers := []model.Voucher{
		{
			Code:       "DISCOUNT10",
			Discount:   10000,
			ExpiryDate: time.Date(2025, 8, 28, 0, 0, 0, 0, time.UTC),
			Status:     "active",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Code:       "PROMO20",
			Discount:   20000,
			ExpiryDate: time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
			Status:     "active",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Code:       "HEMAT50",
			Discount:   50000,
			ExpiryDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			Status:     "inactive",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}
	db.Create(&vouchers)

}
