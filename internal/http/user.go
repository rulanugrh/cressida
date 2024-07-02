package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/middleware"
	"github.com/rulanugrh/cressida/internal/service"
)

type UserHandler interface {
	// Endpoint for register user
	Register(w http.ResponseWriter, r *http.Request)
	// Endpoint for login user
	Login(w http.ResponseWriter, r *http.Request)
}

type user struct {
	service service.UserServive
	middleware middleware.InterfaceJWT
}

func NewUserHandler(service service.UserServive) UserHandler {
	return &user{
		service: service,
		middleware: middleware.NewJSONWebToken(),
	}
}

// @Summary register new account
// @ID register
// @Tags users
// @Accept json
// @Produce json
// @Param request body web.Register true "request body for register new user"
// @Router /api/user/register [post]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func(u *user) Register(w http.ResponseWriter, r *http.Request) {
	// Decode request bdoy
	var request web.Register
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(500)
		w.Write(web.Marshalling(web.InternalServerError("error while decode request body")))
		return
	}

	// parsing request body into service
	data, err := u.service.Register(request)
	if err != nil {
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest(err.Error())))
		return
	}

	w.WriteHeader(201)
	w.Write(web.Marshalling(web.Created("success create account", data)))
}

// @Summary login account
// @ID login
// @Tags users
// @Accept json
// @Produce json
// @Param request body web.Login true "request body for login user"
// @Router /api/user/login [post]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func(u *user) Login(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var request web.Login
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(500)
		w.Write(web.Marshalling(web.InternalServerError("error while decode request body")))
		return
	}

	// parsing request body into layer service
	data, err := u.service.Login(request)
	if err != nil {
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest(err.Error())))
		return
	}

	// create access token jwt
	accessToken, err := u.middleware.GenerateAccessToken(*data)
	if err != nil {
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest("cannot create access token")))
		return
	}

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("succes login user", accessToken)))
}

// @Summary get detail account
// @ID getme
// @Tags users
// @Accept json
// @Produce json
// @Router /api/user/getme [get]
// @Security ApiKeyAuth
// @Success 302 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 401 {object} web.Response
func (u *user) GetMe(w http.ResponseWriter, r *http.Request) {
	// read header authorization
	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(401)
		w.Write(web.Marshalling(web.Unauthorized("token is required")))
		return
	}

	// read token for get email user
	email, err := u.middleware.CheckEmail(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(401)
		w.Write(web.Marshalling(web.Unauthorized("sorry you must login first")))
		return
	}

	// parsing email into function get me
	data, err := u.service.GetMe(*email)
	if err != nil {
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest(err.Error())))
		return
	}

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Found("success found account", data)))
}
