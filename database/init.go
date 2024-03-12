package database

import (
	"os"

	"github.com/Hunter-club/cloudman/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	if os.Getenv("MODE") == "prod" {
		dsn := "root:123456@tcp(127.0.0.1:3306)/cloud?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		Db = db
		return
	}
	db, err := gorm.Open(sqlite.Open("/Users/csh0101/lab/go-playground/cloudman/db/gorm.sqlite"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Db = db
	err = Db.AutoMigrate(&models.Host{}, &models.HostTransfer{}, &models.HostOrderAllocate{}, &models.OrderSub{}, &models.Transfer{})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return Db
}
