package main

import (
	"log"

	"pasarmalam/config"
	"pasarmalam/database"
	"pasarmalam/routes"
	"pasarmalam/services"
)

func main() {
	cfg := config.Load()

	db := database.InitMySQL(cfg)
	fbApp := database.InitFirebase(cfg)

	jwtSvc := services.NewJWTService(cfg.JWTSecret, cfg.JWTExpiryHours)
	fbSvc := services.NewFirebaseAuthService(fbApp)

	r := routes.Setup(db, jwtSvc, fbSvc)

	log.Printf("Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
