package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"pasarmalam/config"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Logger level Error: hanya log error fatal. RecordNotFound dan
	// duplicate-key (1062) yang kita tangani eksplisit di handler
	// tidak akan muncul sebagai noise.
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
		},
	)

	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
		if err == nil {
			break
		}
		log.Printf("[mysql] retry %d/30: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("[mysql] gagal konek: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("[mysql] connected")
	return db
}

func InitFirebase(cfg *config.Config) *firebase.App {
	credFile := cfg.FirebaseCredsPath
	if _, err := os.Stat(credFile); err != nil {
		log.Printf("[firebase] credentials file tidak ditemukan di %s — token verify dev mode", credFile)
		return nil
	}
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "dompet-digital-13b1b",
	}, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Printf("[firebase] gagal inisialisasi: %v (token verify dev mode)", err)
		return nil
	}
	log.Println("[firebase] initialized")
	return app
}
