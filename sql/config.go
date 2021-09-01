package sql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	rrs "github.com/valentinvstoyanov/rating-review-system"
	"github.com/valentinvstoyanov/rating-review-system/env"
	"log"
	"os"
)

var db *gorm.DB

const (
	dbHostVar     = "DB_HOST"
	dbUserVar     = "DB_USER"
	dbPassVar     = "DB_PASS"
	dbNameVar     = "DB_NAME"
	dbPortVar     = "DB_PORT"
	dbDriverVar   = "DB_DRIVER"
	dbFileNameVar = "DB_FILE_NAME"
)

func CreateDatabaseConnection() {
	driver := env.GetEnvVar(dbDriverVar)
	var err error

	if env.IsProd() {
		host := env.GetEnvVar(dbHostVar)
		user := env.GetEnvVar(dbUserVar)
		pass := os.Getenv(dbPassVar)
		name := env.GetEnvVar(dbNameVar)
		port := env.GetEnvVar(dbPortVar)
		connection := user + ":" + pass + "@(" + host + ":" + port + ")/" + name + "?charset=utf8&parseTime=True&loc=Local"

		db, err = gorm.Open(driver, connection)
		log.Printf("%s db setup: %s\n", driver, connection)
	} else {
		fileName := env.GetEnvVar(dbFileNameVar)
		db, err = gorm.Open(driver, fileName)
		log.Printf("%s db setup: %s\n", driver, fileName)
	}

	if err != nil {
		panic(err)
	}
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
