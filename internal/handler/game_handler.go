package handler

import (
	"golang-woktopup/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GameHandler struct {
	DB *gorm.DB
}

func NewGameHandler(db *gorm.DB) *GameHandler {
	return &GameHandler{DB: db}
}

// GET /api/games
func (h *GameHandler) GetGames(c *gin.Context) {
	var games []model.Game
	if err := h.DB.Where("status = ?", "available").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch games"})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (h *GameHandler) GetGameDetail(c *gin.Context) {
	id := c.Param("id")
	var game model.Game

	err := h.DB.Preload("Products", "status = ?", "available").First(&game, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, game)
}
