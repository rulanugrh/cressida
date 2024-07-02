package repository

import (
	"fmt"

	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/helper"
)

type UserRepository interface {
	Create(request web.Register) (*domain.User, error)
	CheckEmail(email string) error
	Login(request web.Login) (*domain.User, error)
	GetMe(email string) (*domain.User, error)
}

type user struct {
	conn *config.SDatabase
	log helper.ILog
}

func NewUserRepository(conn *config.SDatabase) UserRepository {
	return &user{conn: conn, log: helper.NewLogger()}

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
		u.log.Error(fmt.Sprintf("[REPOSITORY] - [CreateUser] %s", err.Error()))
		return nil, err
	}

	u.log.Info("[REPOSITORY] - [CreateUser] success create user")
	return &response, nil

}

func(u *user) CheckEmail(email string) error {
	find := u.conn.DB.Exec("SELECT * FROM users WHERE email = ?",email)
	if find.RowsAffected > 0 {
		u.log.Error(fmt.Sprintf("[REPOSITORY] - [CheckEmail] %s", find.Error))
		return web.BadRequest("Sorry email has been taken")
	}

	u.log.Info(fmt.Sprintf("[REPOSITORY] - [CheckEmail] %s bisa digunakan", email))
	return nil
}

func(u *user) Login(request web.Login) (*domain.User, error) {
	var response domain.User
	find := u.conn.DB.Exec("SELECT * FROM users WHERE email = ?",request.Email).Preload("Role").Find(&response)
	if find.RowsAffected < 1 {
		u.log.Error(fmt.Sprintf("[REPOSITORY] - [Login] %s", find.Error))
		return nil, find.Error
	}

	u.log.Info("[REPOSITORY] - [Login] berhasil login")
	return &response, nil
}


func(u *user) GetMe(email string) (*domain.User, error) {
	var response domain.User
	find := u.conn.DB.Exec("SELECT * FROM users WHERE email = ?",email).Preload("Role").Find(&response)
	if find.RowsAffected < 1 {
		u.log.Error(fmt.Sprintf("[REPOSITORY] - [GetMe] %s", find.Error))
		return nil, find.Error
	}

	u.log.Info(fmt.Sprintf("[REPOSITORY] - [GetMe] %s berhasil mendapatkan detail user", email))
	return &response, nil
}
