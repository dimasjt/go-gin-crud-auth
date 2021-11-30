package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProviderDB struct {
	Storage *gorm.DB
}

func ProvideDB() ProviderDB {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Article{})
	db.AutoMigrate(&User{})

	return ProviderDB{Storage: db}
}
