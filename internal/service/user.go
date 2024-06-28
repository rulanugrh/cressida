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
	InsertToken(email string, token string) error
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
		u.log.Error(err)
		return nil, u.validate.ValidationMessage(err)
	}

	// check email is has been taken
	err = u.repository.CheckEmail(request.Email)
	if err != nil {
		u.log.Debug("someone request duplicate email")
		return nil, web.BadRequest("Email has been taken")
	}

	// hashing password before insert into db
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		u.log.Debug("error hashing password")
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
		u.log.Error(errCreate)
		return nil, web.BadRequest("cannot create user")
	}

	// parsing into response for handler
	response := web.ResponseRegister{
		FName: result.FName,
		LName: result.LName,
	}

	u.log.Info("Success Create Account", result.FName)
	return &response, nil
}

func (u *user) Login(request web.Login) (*web.ResponseLogin, error) {
	// validation request user
	err := u.validate.Validate(request)
	if err != nil {
		u.log.Error(err)
		return nil, u.validate.ValidationMessage(err)
	}

	// checking data in database
	result, err := u.repository.Login(request)
	if err != nil {
		u.log.Debug(fmt.Sprintf("%s not found", request.Email))
		return nil, web.BadRequest("sorry you account not found")
	}

	// checking hash password is valid
	errCompare := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(request.Password))
	if errCompare != nil {
		u.log.Warn(fmt.Sprintf("%s request password but not matched", request.Email))
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

	u.log.Info("Success Login", request.Email)
	return &response, nil
}

func (u *user) GetMe(email string) (*web.ResponseGetUser, error) {
	// checking data in database
	result, err := u.repository.GetMe(email)
	if err != nil {
		u.log.Error(err)
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

	u.log.Info("Account Found", email)
	return &response, nil
}

func (u *user) InsertToken(email string, token string) error {
	// checking and insert token into db
	err := u.repository.InsertToken(email, token)
	if err != nil {
		u.log.Error(err)
		return web.BadRequest("error while input token")
	}

	u.log.Info("Success Insert Token", email)
	return nil
}