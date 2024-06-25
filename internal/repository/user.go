package repository

import (
	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/entity/web"
)

type UserRepository interface {
	Create(request web.Register) (*domain.User, error)
	CheckEmail(email string) error
	Login(request web.Login) (*domain.User, error)
	GetMe(email string) (*domain.User, error)
}

type user struct {
	conn *config.SDatabase
}

func NewUserRepository(conn *config.SDatabase) UserRepository {
	return &user{conn: conn}

}

func(u *user) Create(request web.Register) (*domain.User, error) {
	var response domain.User

	create := u.conn.DB.Exec("INSERT INTO users(f_name, l_name, email, address, password, phone, role_id) VALUES (?,?,?,?,?,?,?)",
		request.FName,
		request.LName,
		request.Email,
		request.Address,
		request.Password,
		request.Phone,
		request.RoleID,
	).Scan(&response)

	err := create.Model(&response.Role).Association("Users").Append(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

func(u *user) CheckEmail(email string) error {
	find := u.conn.DB.Exec("SELECT * FROM users WHERE email = ?",email)
	if find.RowsAffected > 0 {
		return web.BadRequest("Sorry email has been taken")
	}

	return nil
}

func(u *user) Login(request web.Login) (*domain.User, error) {
	var response domain.User
	find := u.conn.DB.Exec("SELECT * FROM users WHERE email = ?",request.Email).Preload("Role").Find(&response)
	if find.RowsAffected < 1 {
		return nil, find.Error
	}

	return &response, nil
}


func(u *user) GetMe(email string) (*domain.User, error) {
	var response domain.User
	find := u.conn.DB.Exec("SELECT * FROM users WHERE email = ?",email).Preload("Role").Find(&response)
	if find.RowsAffected < 1 {
		return nil, find.Error
	}

	return &response, nil
}