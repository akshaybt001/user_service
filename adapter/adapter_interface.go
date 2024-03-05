package adapter

import "github.com/akshaybt001/user_service/entities"

type UserInterface interface {
	UserSignUp(req entities.User) (entities.User,error)
	UserLogin(email string)(entities.User,error)
}