package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pasarmalam/middleware"
	"pasarmalam/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

type checkoutReq struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	Notes           string `json:"notes"`
	PaymentMethod   string `json:"payment_method" binding:"required"`
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	var req checkoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "shipping_address dan payment_method wajib diisi"})
		return
	}

	// Ambil cart user
	var cart models.Cart
	if err := h.DB.Preload("Items.Product").Where("user_id = ?", uid).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Keranjang kosong"})
		return
	}
	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Keranjang kosong"})
		return
	}

	// Hitung total & buat order items
	var total float64
	orderItems := make([]models.OrderItem, 0, len(cart.Items))
	for _, it := range cart.Items {
		if it.Product == nil {
			continue
		}
		sub := it.Product.Price * float64(it.Quantity)
		total += sub
		orderItems = append(orderItems, models.OrderItem{
			ProductID:   it.ProductID,
			ProductName: it.Product.Name,
			Price:       it.Product.Price,
			Quantity:    it.Quantity,
			Subtotal:    sub,
		})
	}

	order := models.Order{
		UserID:          uid,
		TotalAmount:     total,
		Status:          "pending",
		ShippingAddress: req.ShippingAddress,
		Notes:           req.Notes,
		PaymentMethod:   req.PaymentMethod,
		Items:           orderItems,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Simulasi payment per method
	switch req.PaymentMethod {
	case "virtual_account":
		order.VANumber = generateVA()
	case "gopay":
		order.GopayDeeplink = fmt.Sprintf("pasarmalam://pay?order_id=%d&amount=%.0f", 0, total)
	}

	// Transaksi: insert order + clear cart items
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal checkout: " + err.Error()})
		return
	}

	// Set deeplink dgn ID asli
	if order.GopayDeeplink != "" {
		order.GopayDeeplink = fmt.Sprintf("pasarmalam://pay?order_id=%d&amount=%.0f", order.ID, order.TotalAmount)
		h.DB.Model(&order).Update("gopay_deeplink", order.GopayDeeplink)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order berhasil dibuat",
		"data":    order,
	})
}

func (h *OrderHandler) MyOrders(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var total int64
	h.DB.Model(&models.Order{}).Where("user_id = ?", uid).Count(&total)

	var orders []models.Order
	if err := h.DB.Preload("Items").Where("user_id = ?", uid).
		Order("id DESC").Limit(limit).Offset((page - 1) * limit).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    orders,
		"meta":    gin.H{"page": page, "limit": limit, "total": total},
	})
}

func (h *OrderHandler) Detail(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	orderID, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	err := h.DB.Preload("Items").Where("id = ? AND user_id = ?", orderID, uid).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Order tidak ditemukan"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    order,
	})
}

func generateVA() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "8808" + hex.EncodeToString(b)
}
