package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pasarmalam/middleware"
	"pasarmalam/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	DB *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{DB: db}
}

func (h *CartHandler) getOrCreateCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	// Try to fetch existing cart first
	if err := h.DB.Where("user_id = ?", userID).First(&cart).Error; err == nil {
		return &cart, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Cart belum ada — buat baru. Handle race condition: jika ada concurrent
	// request yang berhasil INSERT duluan, kita tangkap duplicate-key error
	// lalu re-fetch.
	now := time.Now()
	newCart := models.Cart{
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := h.DB.Create(&newCart).Error
	if err != nil {
		// Deteksi MySQL duplicate-key (1062) → kemungkinan race winner sudah INSERT
		if isDuplicateKeyError(err) {
			if err := h.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return nil, err
	}
	return &newCart, nil
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Error 1062") || strings.Contains(msg, "Duplicate entry")
}

type addToCartReq struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type updateCartReq struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

func (h *CartHandler) List(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	cart, err := h.getOrCreateCart(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	}
	// Reload dengan relasi
	if err := h.DB.Preload("Items.Product").First(cart, cart.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    cart,
	})
}

func (h *CartHandler) Add(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	var req addToCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "product_id dan quantity wajib diisi"})
		return
	}

	var product models.Product
	if err := h.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}

	cart, err := h.getOrCreateCart(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	}

	var item models.CartItem
	err = h.DB.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		item = models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Subtotal:  product.Price * float64(req.Quantity),
		}
		if err := h.DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal tambah item: " + err.Error()})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	} else {
		item.Quantity += req.Quantity
		item.Subtotal = product.Price * float64(item.Quantity)
		if err := h.DB.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal update item: " + err.Error()})
			return
		}
	}

	// Reload cart
	h.DB.Preload("Items.Product").First(cart, cart.ID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item ditambahkan ke keranjang",
		"data":    cart,
	})
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	itemID, _ := strconv.Atoi(c.Param("id"))
	var req updateCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "quantity wajib diisi"})
		return
	}

	cart, err := h.getOrCreateCart(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	}

	var item models.CartItem
	if err := h.DB.Where("id = ? AND cart_id = ?", itemID, cart.ID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Item tidak ditemukan"})
		return
	}

	var product models.Product
	h.DB.First(&product, item.ProductID)

	item.Quantity = req.Quantity
	item.Subtotal = product.Price * float64(req.Quantity)
	if err := h.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal update: " + err.Error()})
		return
	}

	h.DB.Preload("Items.Product").First(cart, cart.ID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item diperbarui",
		"data":    cart,
	})
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	itemID, _ := strconv.Atoi(c.Param("id"))

	cart, err := h.getOrCreateCart(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	}

	if err := h.DB.Where("id = ? AND cart_id = ?", itemID, cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal hapus: " + err.Error()})
		return
	}

	h.DB.Preload("Items.Product").First(cart, cart.ID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item dihapus",
		"data":    cart,
	})
}

func (h *CartHandler) Clear(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	cart, err := h.getOrCreateCart(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + err.Error()})
		return
	}
	if err := h.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal clear: " + err.Error()})
		return
	}
	h.DB.Preload("Items.Product").First(cart, cart.ID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Keranjang dikosongkan",
		"data":    cart,
	})
}
