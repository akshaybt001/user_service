package adapter

import "github.com/akshaybt001/user_service/entities"

type UserInterface interface {
	UserSignUp(req entities.User) (entities.User, error)
	UserLogin(email string) (entities.User, error)
	AdminLogin(email string) (entities.Admin, error)
	SupAdminLogin(email string) (entities.SupAdmin, error)
	AddAdmin(req entities.Admin)(entities.Admin,error)
	GetAllAdmins() ([]entities.Admin,error)
	GetAllUsers() ([]entities.User,error)
}
