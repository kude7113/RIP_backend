package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"RIP/internal/app/ds"
	"RIP/internal/app/dsn"
)

func main() {

	_ = godotenv.Load()
	env := dsn.FromEnv()
	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})
	if err != nil {
		println("failed to connect database" + err.Error())
	}
	// Migrate the schema
	if err = db.AutoMigrate(
		&ds.Fines{},
		&ds.Resolutions{},
		&ds.Fine_Resolutions{},
		&ds.Users{},
	); err != nil {
		println("cant migrate db")
	}
}
