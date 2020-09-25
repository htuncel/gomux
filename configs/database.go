package configs

import (
	"github.com/jinzhu/gorm"
	// sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"main/models"
)

var (
	// DB singleton instance
	DB *gorm.DB
)

func init() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&models.Book{})

	DB = database
}
