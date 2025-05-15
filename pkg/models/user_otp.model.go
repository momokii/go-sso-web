package models

type OtpPurpose string
type OtpChannel string

const (
	// OtpPurpose constants
	OTPPurpose_PhoneVerification OtpPurpose = "phone_verification"
	OTPPurpose_PasswordReset     OtpPurpose = "password_reset"
	OTPPurpose_TwoFactorAuth     OtpPurpose = "login_2fa"

	// OtpChannel constants
	OTPChannel_Email    OtpChannel = "email"
	OTPChannel_SMS      OtpChannel = "sms"
	OTPChannel_Whatsapp OtpChannel = "whatsapp"
)

type UserOtps struct {
	Id        int        `json:"id"`
	UserId    int        `json:"user_id"`
	Purpose   OtpPurpose `json:"purpose"`
	Channel   OtpChannel `json:"channel"`
	CodeHash  string     `json:"code_hash"`
	CreatedAt string     `json:"created_at"`
	ExpiresAt string     `json:"expires_at"`
	Used      bool       `json:"used"`
	UsedAt    string     `json:"used_at"`
}

type UserOtpInput struct {
	UserId   int        `json:"user_id" validate:"required"`
	Purpose  OtpPurpose `json:"purpose" validate:"required"`
	Channel  OtpChannel `json:"channel" validate:"required"`
	CodeHash string     `json:"code_hash" validate:"required"`
}

type UserOtpUpdate struct {
	Id     int `json:"id" validate:"required"`
	UserId int `json:"user_id" validate:"required"`
}

type UserOtpVerify struct {
	Otp string `json:"otp" validate:"required"`
}
