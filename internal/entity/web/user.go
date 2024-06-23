package web

type Register struct {
	FName    string `json:"f_name" form:"f_name" validate:"required"`
	LName    string `json:"l_name" form:"l_name" validate:"required"`
	Email    string `json:"email" form:"email" gorm:"type:unique" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
	Address  string `json:"address" form:"address" validate:"required"`
	Phone    string `json:"phone" form:"phone" validate:"required"`
	RoleID   uint   `json:"role_id" form:"role_id"`
}

type Login struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

type ResponseRegister struct {
	FName string `json:"f_name" form:"f_name"`
	LName string `json:"l_name" form:"l_name"`
}

type ResponseLogin struct {
	ID     uint   `json:"id" form:"id"`
	FName  string `json:"f_name" form:"f_name"`
	LName  string `json:"l_name" form:"l_name"`
	Email  string `json:"email" form:"email"`
	RoleID uint   `json:"role_id" form:"role_id"`
}

type ResponseGetUser struct {
	FName    string `json:"f_name" form:"f_name"`
	LName    string `json:"l_name" form:"l_name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Address  string `json:"address" form:"address"`
	Phone    string `json:"phone" form:"phone"`
	Role     string `json:"role" form:"role"`
}
