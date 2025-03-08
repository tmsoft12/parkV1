package database

import (
	"log"
	"os"
	"park/models/camera"
	modelscar "park/models/modelsCar"
	modelsuser "park/models/modelsUser"
	modeloperator "park/models/operatorModel"
	"park/models/tarif"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}
	err = database.AutoMigrate(
		&modelscar.Car_Model{},
		&modelsuser.User{},
		&camera.Cameras{},
		&modeloperator.Operator{},
		&tarif.Tarif{},
		&camera.CamFix{},
		&modelsuser.MacUser{},
	)
	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}
	DB = database
	log.Println("Successfully connected to PostgreSQL")
}
