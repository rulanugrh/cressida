package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/entity/web"
)

type InterfaceJWT interface {
	GenerateRefreshToken(request domain.User) (*string, error)
	GenerateAccessToken(request domain.User) (*string, error)
	CheckUserID(token string) (*uint, error)
}

type JWT struct {
	config *config.App
}

type jwtClaim struct {
	ID     uint   `json:"id"`
	FName  string `json:"f_name"`
	LName  string `json:"l_name"`
	Email  string `json:"email"`
	RoleID uint   `json:"role_id"`
	jwt.RegisteredClaims
}

type refreshClaim struct {
	ID    uint   `json:"id"`
	FName string `json:"f_name"`
	LName string `json:"l_name"`
	jwt.RegisteredClaims
}

func NewJSONWebToken() InterfaceJWT {
	return &JWT{
		config: config.GetConfig(),
	}
}


func(j *JWT) GenerateRefreshToken(request domain.User) (*string, error) {
	claim := &refreshClaim{
		ID: request.ID,
		FName: request.FName,
		LName: request.LName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
		},
	}

	refreshtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := refreshtoken.SignedString([]byte(j.config.Server.Secret))
	if err != nil {
		return nil, web.BadRequest("cannot sign secret token")
	}

	return &token, nil
}

func(j *JWT) GenerateAccessToken(request domain.User) (*string, error) {
	claim := &jwtClaim{
		ID: request.ID,
		FName: request.FName,
		LName: request.LName,
		Email: request.Email,
		RoleID: request.RoleID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := accessToken.SignedString([]byte(j.config.Server.Secret))
	if err != nil {
		return nil, web.BadRequest("Sorry cannot generate for access token")
	}

	return &token, nil
}

func(j *JWT) CheckUserID(token string) (*uint, error) {
	tkn, err := jwt.ParseWithClaims(token, &jwtClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.config.Server.Secret), nil
	})

	if err != nil {
		return nil, web.BadRequest("error while parsing claim token")
	}

	claim, valid := tkn.Claims.(*jwtClaim)
	if !valid {
		return nil, web.Unauthorized("Sorry you're not loggin into app")
	}

	return &claim.ID, nil
}
