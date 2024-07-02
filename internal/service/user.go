package service

import (
	"fmt"

	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/helper"
	"github.com/rulanugrh/cressida/internal/middleware"
	"github.com/rulanugrh/cressida/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserServive interface {
	Register(request web.Register) (*web.ResponseRegister, error)
	Login(request web.Login) (*web.ResponseLogin, error)
	GetMe(email string) (*web.ResponseGetUser, error)
}

type user struct {
	repository repository.UserRepository
	validate   middleware.IValidation
	log        helper.ILog
}

func NewUserService(repository repository.UserRepository) UserServive {
	return &user{
		repository: repository,
		validate:   middleware.NewValidation(),
		log:        helper.NewLogger(),
	}
}

func (u *user) Register(request web.Register) (*web.ResponseRegister, error) {
	// validation struct for request
	err := u.validate.Validate(request)
	if err != nil {
		u.log.Error(fmt.Sprintf("[SERVICE] - [Register] Error while validate request: %s", err.Error()))
		return nil, u.validate.ValidationMessage(err)
	}

	// check email is has been taken
	err = u.repository.CheckEmail(request.Email)
	if err != nil {
		u.log.Debug(fmt.Sprintf("[SERVICE] - [Register] email: %s trying to taken again", request.Email))
		return nil, web.BadRequest("Email has been taken")
	}

	// hashing password before insert into db
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		u.log.Debug(fmt.Sprintf("[SERVICE] - [Register] email: %s error while hashing password", request.Email))
		return nil, web.InternalServerError("Error while hashing password")
	}

	// parsing value from request and hashing password
	req := web.Register{
		FName:    request.FName,
		LName:    request.LName,
		Email:    request.Email,
		Password: string(hashedPassword),
		Address:  request.Address,
		Phone:    request.Phone,
		RoleID:   2,
	}

	// save data into db
	result, errCreate := u.repository.Create(req)
	if errCreate != nil {
		u.log.Error(fmt.Sprintf("[SERVICE] - [Register] Error while input into db: %s", errCreate.Error()))
		return nil, web.BadRequest("cannot create user")
	}

	// parsing into response for handler
	response := web.ResponseRegister{
		FName: result.FName,
		LName: result.LName,
	}

	u.log.Info(fmt.Sprintf("[SERVICE] - [Register] %s success create user", result.Email))
	return &response, nil
}

func (u *user) Login(request web.Login) (*web.ResponseLogin, error) {
	// validation request user
	err := u.validate.Validate(request)
	if err != nil {
		u.log.Error(fmt.Sprintf("[SERVICE] - [Login] Error while validate request: %s", err.Error()))
		return nil, u.validate.ValidationMessage(err)
	}

	// checking data in database
	result, err := u.repository.Login(request)
	if err != nil {
		u.log.Error(fmt.Sprintf("[SERVICE] - [Login] Error while get into db: %s", err.Error()))
		return nil, web.BadRequest("sorry you account not found")
	}

	// checking hash password is valid
	errCompare := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(request.Password))
	if errCompare != nil {
		u.log.Warn(fmt.Sprintf("[SERVICE] - [Login] email: %s, request password but not matched", request.Email))
		return nil, web.Unauthorized("sorry your password is not matched")
	}

	// parsing value into response
	response := web.ResponseLogin{
		ID:     result.ID,
		FName:  result.FName,
		LName:  result.LName,
		Email:  result.Email,
		RoleID: result.RoleID,
	}

	u.log.Info(fmt.Sprintf("[SERVICE] - [Login] %s success login", result.Email))
	return &response, nil
}

func (u *user) GetMe(email string) (*web.ResponseGetUser, error) {
	// checking data in database
	result, err := u.repository.GetMe(email)
	if err != nil {
		u.log.Error(fmt.Sprintf("[SERVICE] - [Login] Error while get into db: %s", err.Error()))
		return nil, web.BadRequest("sorry your account not found")
	}

	// parsing value into response
	response := web.ResponseGetUser{
		FName:   result.FName,
		LName:   result.LName,
		Email:   result.Email,
		Address: result.Address,
		Phone:   result.Phone,
		Role:    result.Role.Name,
	}

	u.log.Info(fmt.Sprintf("[SERVICE] - [GetMe] %s success found", result.Email))
	return &response, nil
}
