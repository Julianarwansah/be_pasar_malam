package routes

import (
	"pasarmalam/handlers"
	"pasarmalam/middleware"
	"pasarmalam/seed"
	"pasarmalam/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, jwtSvc *services.JWTService, fbSvc *services.FirebaseAuthService) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// Auto-migrate
	if err := db.AutoMigrate(
		&databaseUser{}, &databaseProduct{}, &databaseCart{}, &databaseCartItem{},
		&databaseOrder{}, &databaseOrderItem{},
	); err != nil {
		panic(err)
	}
	seed.Run(db)

	authH := handlers.NewAuthHandler(db, jwtSvc, fbSvc, nil)
	productH := handlers.NewProductHandler(db)
	cartH := handlers.NewCartHandler(db)
	orderH := handlers.NewOrderHandler(db)

	v1 := r.Group("/v1")
	{
		v1.GET("/health", handlers.HealthCheck)

		auth := v1.Group("/auth")
		{
			auth.POST("/verify-token", authH.VerifyToken)
			auth.POST("/dev-verify-email", authH.DevVerifyEmail) // dev only

			authRequired := auth.Group("")
			authRequired.Use(middleware.Auth(jwtSvc))
			{
				authRequired.GET("/me", authH.Me)
				authRequired.PUT("/fcm-token", authH.UpdateFCMToken)
			}
		}

		products := v1.Group("/products")
		{
			products.GET("", productH.List)
			products.GET("/:id", productH.Get)
		}

		cart := v1.Group("/cart")
		cart.Use(middleware.Auth(jwtSvc))
		{
			cart.GET("", cartH.List)
			cart.POST("", cartH.Add)
			cart.PUT("/:id", cartH.UpdateItem)
			cart.DELETE("/:id", cartH.DeleteItem)
			cart.DELETE("", cartH.Clear)
		}

		orders := v1.Group("/orders")
		orders.Use(middleware.Auth(jwtSvc))
		{
			orders.GET("", orderH.MyOrders)
			orders.POST("/checkout", orderH.Checkout)
			orders.GET("/:id", orderH.Detail)
		}
	}

	return r
}
