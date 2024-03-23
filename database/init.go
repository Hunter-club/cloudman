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
	//如果环境变量Mode = prod
	if os.Getenv("MODE") == "prod" {
		// 阿里云账号密码
		dsn := "root:Cc^T#YuUwBA%kb0y@tcp(rm-uf6effv0820ppw455no.mysql.rds.aliyuncs.com:3306)/cloud?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		Db = db
		err = Db.AutoMigrate(&models.Host{}, &models.HostTransfer{}, &models.HostOrderAllocate{}, &models.OrderSub{}, &models.Transfer{})
		if err != nil {
			panic(err)
		}
		return
	} else {

		db, err := gorm.Open(sqlite.Open("./gorm.sqlite"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		Db = db
		err = Db.AutoMigrate(&models.Host{}, &models.HostTransfer{}, &models.HostOrderAllocate{}, &models.OrderSub{}, &models.Transfer{})
		if err != nil {
			panic(err)
		}
	}
}

func GetDB() *gorm.DB {
	return Db
}
