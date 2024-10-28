package db

import (
	"example.com/go-crud-api/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Db *gorm.DB
var err error

func InitPostgresDB() {
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	var (
		//host     = "host.docker.internal"
		//host     = "localhost"
		//port     = "5432"
		//dbUser   = "postgres"
		//dbName   = "demo"
		//password = "Varanasi@123"

		host     = os.Getenv("DB_HOST")
		dbUser   = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbName   = os.Getenv("DB_NAME")
		port     = os.Getenv("DB_PORT")
	)
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)

	Db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	Db.AutoMigrate(models.User{}, models.LimitUpdateJson{}, models.AsyncProcessStatus{})
}
