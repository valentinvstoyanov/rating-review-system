package sql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	rrs "github.com/valentinvstoyanov/rating-review-system"
)

var db *gorm.DB

func CreateDatabaseConnection() {
	d, err := gorm.Open("sqlite3", "rating_reviews_system.db")
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}

func AutoMigrateAll() {
	db.AutoMigrate(&rrs.User{})
	db.AutoMigrate(&rrs.Entity{})
	db.AutoMigrate(&rrs.Review{})
	db.AutoMigrate(&rrs.RatingAlert{})
}
