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

	query := "SELECT * FROM users WHERE email = $1"

	return res, u.DB.Raw(query, email).Scan(&res).Error
}

func (u *UserAdapter) AdminLogin(email string) (entities.Admin, error) {
	var res entities.Admin

	query := "SELECT * FROM admins Where email = $1"

	return res, u.DB.Raw(query, email).Scan(&res).Error
}

func (u *UserAdapter) SupAdminLogin(email string) (entities.SupAdmin, error) {
	var res entities.SupAdmin

	query := "SELECT * FROM sup_admins WHERE email = $1"

	return res, u.DB.Raw(query, email).Scan(&res).Error
}

func (u *UserAdapter) AddAdmin(req entities.Admin) (entities.Admin, error) {
	var res entities.Admin

	query := "INSERT INTO admins (name,email,password) VALUES($1,$2,$3) RETURNING id,name,email"

	return res, u.DB.Raw(query, req.Name, req.Email, req.Password).Scan(&res).Error
}

func (sup *UserAdapter) GetAllAdmins() ([]entities.Admin, error) {
	var res []entities.Admin

	query := "SELECT * FROM admins"

	if err := sup.DB.Raw(query).Scan(&res).Error; err != nil {
		return []entities.Admin{}, err
	}
	return res, nil
}

func (admin *UserAdapter) GetAllUsers() ([]entities.User, error) {
	var res []entities.User

	query := "SELECT * FROM users"
	if err := admin.DB.Raw(query).Scan(&res).Error; err != nil {
		return []entities.User{}, err
	}
	return res, nil
}
