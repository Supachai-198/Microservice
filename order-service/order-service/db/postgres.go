package db

import (
	"fmt"
	"os"

	"example.com/order-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbDNS := os.Getenv("DATABASE_DNS")

	db, err := gorm.Open(postgres.Open(dbDNS), &gorm.Config{})

	if err != nil {
		fmt.Println("ติดต่อฐานข้อมูลไม่ได้")
		panic(err)
	}

	fmt.Println("ติดต่อฐานข้อมูลสำเร็จ")

	//Auto Migrations
	db.AutoMigrate(&models.Order{})

	DB = db
}
