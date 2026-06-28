package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"pasarmalam/config"
	"pasarmalam/middleware"
	"pasarmalam/models"
	"pasarmalam/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB    *gorm.DB
	JWT   *services.JWTService
	FB    *services.FirebaseAuthService
	Cfg   *config.Config
}

type verifyTokenReq struct {
	FirebaseToken string `json:"firebase_token" binding:"required"`
}

type fcmTokenReq struct {
	FCMToken string `json:"fcm_token" binding:"required"`
}

func NewAuthHandler(db *gorm.DB, jwt *services.JWTService, fb *services.FirebaseAuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, JWT: jwt, FB: fb, Cfg: cfg}
}

func (h *AuthHandler) VerifyToken(c *gin.Context) {
	var req verifyTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "firebase_token wajib diisi", "error_code": "INVALID_REQUEST"})
		return
	}

	claims, err := h.FB.VerifyIDToken(req.FirebaseToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Firebase token tidak valid: " + err.Error(), "error_code": "INVALID_FIREBASE_TOKEN"})
		return
	}

	// Upsert user
	var user models.User
	tx := h.DB.Where("firebase_uid = ?", claims.UID).First(&user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		user = models.User{
			FirebaseUID:   claims.UID,
			Email:         claims.Email,
			Name:          claims.Name,
			EmailVerified: claims.Verified,
			Role:          "customer",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := h.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat user: " + err.Error()})
			return
		}
	} else if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB error: " + tx.Error.Error()})
		return
	} else {
		// Update data dari firebase kalau ada
		updates := map[string]interface{}{"updated_at": time.Now()}
		if claims.Email != "" {
			updates["email"] = claims.Email
		}
		if claims.Name != "" {
			updates["name"] = claims.Name
		}
		updates["email_verified"] = claims.Verified
		h.DB.Model(&user).Updates(updates)
	}

	access, expSec, err := h.JWT.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data": gin.H{
			"access_token": access,
			"token_type":   "Bearer",
			"expires_in":   expSec,
			"user":         user,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	var user models.User
	if err := h.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    user,
	})
}

func (h *AuthHandler) UpdateFCMToken(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	var req fcmTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "fcm_token wajib diisi"})
		return
	}
	if err := h.DB.Model(&models.User{}).Where("id = ?", uid).Update("fcm_token", req.FCMToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal update fcm_token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "FCM token diperbarui"})
}

type devVerifyReq struct {
	Email string `json:"email" binding:"required"`
}

// DevVerifyEmail — endpoint development untuk menandai user sebagai
// email_verified di Firebase & database. HANYA aktif kalau ENABLE_DEV_AUTH=true.
// Berguna untuk testing tanpa harus klik link verifikasi di email.
func (h *AuthHandler) DevVerifyEmail(c *gin.Context) {
	if os.Getenv("ENABLE_DEV_AUTH") != "true" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Dev auth endpoint tidak aktif"})
		return
	}
	var req devVerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "email wajib diisi"})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User tidak ditemukan di DB lokal"})
		return
	}

	// Update di Firebase
	if h.FB != nil {
		if err := h.FB.SetEmailVerified(c.Request.Context(), user.FirebaseUID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal set verified di Firebase: " + err.Error()})
			return
		}
	}

	// Update di DB
	if err := h.DB.Model(&user).Update("email_verified", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal update DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email " + req.Email + " ditandai verified",
	})
}
