package seed

import (
	"log"

	"pasarmalam/models"

	"gorm.io/gorm"
)

type seedProduct struct {
	Name        string
	Description string
	Price       float64
	Stock       int
	Category    string
	ImageURL    string
}

var products = []seedProduct{
	{"Sate Ayam Madura", "Sate ayam khas Madura dengan bumbu kacang pedas, 10 tusuk.", 25000, 50, "Makanan", ""},
	{"Bakso Urat", "Bakso urat sapi kenyal dengan kuah kaldu panas, mangkok.", 18000, 40, "Makanan", ""},
	{"Bakso Bakar", "Bakso bakar saus padang pedas manis, 3 tusuk.", 15000, 30, "Makanan", ""},
	{"Soto Ayam Kampung", "Soto bening ayam kampung dengan tauge, seledri, dan jeruk nipis.", 22000, 25, "Makanan", ""},
	{"Nasi Goreng Spesial", "Nasi goreng dengan telur, ayam suwir, dan kerupuk.", 20000, 35, "Makanan", ""},
	{"Mie Ayam Bakso", "Mie ayam dengan bakso sapi dan pangsit goreng.", 18000, 40, "Makanan", ""},
	{"Es Teh Manis", "Es teh manis segar, gelas besar.", 5000, 100, "Minuman", ""},
	{"Es Jeruk Peras", "Es jeruk peras asli tanpa gula buatan.", 10000, 60, "Minuman", ""},
	{"Bajigur", "Minuman jahe hangat khas Sunda dengan santan dan gula aren.", 8000, 30, "Minuman", ""},
	{"Bandrek", "Bandrek pedas manis dengan jahe, gula merah, dan serai.", 8000, 30, "Minuman", ""},
	{"Pisang Goreng", "Pisang goreng crispy dengan keju dan coklat.", 12000, 40, "Snack", ""},
	{"Tahu Bulat", "Tahu bulat crispy dengan bumbu pedas, 5 buah.", 5000, 80, "Snack", ""},
}

func Run(db *gorm.DB) {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		log.Printf("[seed] products sudah ada (%d) — skip", count)
		return
	}
	for _, p := range products {
		db.Create(&models.Product{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Category:    p.Category,
			ImageURL:    p.ImageURL,
			IsActive:    true,
		})
	}
	log.Printf("[seed] inserted %d products", len(products))
}
