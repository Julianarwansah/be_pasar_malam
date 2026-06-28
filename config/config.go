package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	JWTExpiryHours       int
	FirebaseCredsPath    string
	FirebaseAPIKey       string
}

func Load() *Config {
	_ = godotenv.Load()

	c := &Config{
		Port:              getenv("PORT", "8082"),
		DBHost:            getenv("DB_HOST", "localhost"),
		DBPort:            getenv("DB_PORT", "3306"),
		DBUser:            getenv("DB_USER", "useremoney"),
		DBPassword:        getenv("DB_PASSWORD", "Password#123"),
		DBName:            getenv("DB_NAME", "pasarmalam"),
		JWTSecret:         getenv("JWT_SECRET", "pasarmalam-super-secret-jwt-key"),
		JWTExpiryHours:    168,
		FirebaseCredsPath: getenv("FIREBASE_CREDENTIALS_PATH", "firebase_service_account.json"),
		FirebaseAPIKey:    getenv("FIREBASE_API_KEY", ""),
	}
	log.Printf("[config] Loaded — port=%s db=%s@%s/%s", c.Port, c.DBUser, c.DBHost, c.DBName)
	return c
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
