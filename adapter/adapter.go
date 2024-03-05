package adapter

import (
	"github.com/akshaybt001/user_service/entities"
	"gorm.io/gorm"
)

type UserAdapter struct {
	DB *gorm.DB
}

func NewUserAdapter(db *gorm.DB) *UserAdapter {
	return &UserAdapter{
		DB: db,
	}
}

func (u *UserAdapter) UserSignUp(req entities.User) (entities.User, error) {
	var res entities.User

	query := "INSERT INTO users (name,email,password) VALUES ($1,$2,$3) RETURNING id ,name,email"

	return res, u.DB.Raw(query, req.Name, req.Email, req.Password).Scan(&res).Error
}


func (u *UserAdapter) UserLogin(email string) (entities.User, error) {
	var res entities.User

	query:="SELECT * FROM users WHERE email = $1"

	return res, u.DB.Raw(query,email).Scan(&res).Error
}
