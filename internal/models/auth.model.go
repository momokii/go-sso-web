package models

type AuthLogin struct {
	Username string `json:"username" validate:"required,min=5,max=25,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=50,containsany=1234567890,containsany=QWERTYUIOPASDFGHJKLZXCVBNM"`
}
