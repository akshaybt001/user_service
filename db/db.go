package db

import (
	"github.com/akshaybt001/user_service/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(connect string) (*gorm.DB,error){

	db,err :=gorm.Open(postgres.Open(connect),&gorm.Config{})
	 if err !=nil{
		return nil,err
	 }

	 db.AutoMigrate(&entities.User{})
	 db.AutoMigrate(&entities.Admin{})
	 db.AutoMigrate(&entities.SupAdmin{})

	 return db,nil
}