package initializer

import (
	"github.com/akshaybt001/user_service/adapter"
	"github.com/akshaybt001/user_service/service"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB) *service.UserService{
	adapter:=adapter.NewUserAdapter(db)
	service:=service.NewUserService(adapter)
	return service
}