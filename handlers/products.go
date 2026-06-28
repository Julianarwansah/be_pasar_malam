package handlers

import (
	"net/http"
	"strconv"

	"pasarmalam/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	category := c.Query("category")

	q := h.DB.Model(&models.Product{}).Where("is_active = ?", true)
	if category != "" {
		q = q.Where("category = ?", category)
	}

	var total int64
	q.Count(&total)

	var products []models.Product
	if err := q.Order("id ASC").Limit(limit).Offset((page - 1) * limit).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    products,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p models.Product
	if err := h.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    p,
	})
}
