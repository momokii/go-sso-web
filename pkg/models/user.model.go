package models

type User struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	CreditToken      int    `json:"credit_token"`
	LastFirstLLMUsed string `json:"last_first_llm_used"`
	MultiFAEnabled   bool   `json:"multifa_enabled"`
	PhoneNumber      string `json:"phone_number"`
}

type UserSession struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	CreditToken      int    `json:"credit_token"`
	LastFirstLLMUsed string `json:"last_first_llm_used"`
	MultiFAEnabled   bool   `json:"multifa_enabled"`
	PhoneNumber      string `json:"phone_number"`
}

type UserChangeUsernameInput struct {
	Id       int    `json:"id"`
	Username string `json:"username" validate:"required,min=5,max=25,alphanum"`
}

type UserChangePhoneInput struct {
	Id          int    `json:"id"`
	PhoneNumber string `json:"phone_number" validate:"required,min=9,max=15,number"`
}

type UserChangePasswordInput struct {
	Id          int    `json:"id"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=50,containsany=1234567890,containsany=QWERTYUIOPASDFGHJKLZXCVBNM"`
}

type UserChangePhoneSendOTPInput struct {
	PhoneNumber string `json:"phone_number" validate:"required,min=9,max=15,number"`
}

type UserEditPhoneVerifyOTPInput struct {
	OTPCode     string `json:"otp_code" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,min=9,max=15,number"`
}
