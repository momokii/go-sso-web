package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
)

func GenerateRandomOTPCode(length int) string {
	var random_code string

	const numbers = "0123456789"

	for i := 0; i < length; i++ {
		random_code += string(numbers[rand.Intn(len(numbers))])
	}

	return random_code
}

// GenerateOTPHash generates a random OTP code and its hash.
//
// Parameters:
//   - length: The desired length of the OTP code. Must be between 1 and 12 digits inclusive.
//
// Returns:
//   - string: The generated OTP code.
//   - string: The hash of the OTP code in hexadecimal format.
//   - error: An error if length is invalid (less than 1 or greater than 12).
//
// The function generates a random OTP code of the specified length and then creates
// a HMAC-SHA256 hash of the code. Both the original code and its hash are returned.
// The hash can be used for secure verification without exposing the actual OTP.
func GenerateOTPHash(length int, secret_key string) (string, string, error) {

	if length < 0 {
		return "", "", errors.New("length must be greater than 0")
	}

	if length > 12 {
		return "", "", errors.New("length must be less than 12")
	}

	otp_code := GenerateRandomOTPCode(length)

	otp_hash := hmac.New(sha256.New, []byte(secret_key))
	otp_hash.Write([]byte(otp_code))

	return otp_code, hex.EncodeToString(otp_hash.Sum(nil)), nil
}

func SendOTPMessage(targets, otp_code string) error {
	if len(targets) <= 0 {
		return errors.New("targets cannot be empty")
	}

	message_template := "Your verification code is: %s\n\nThis code will expire in 2 minutes.\n\nIf you did not request this code, please ignore this message.\n\nRegards,\nGo Klan SSO Team"
	message := fmt.Sprintf(message_template, otp_code)

	target_number := []string{targets}

	wa_response, err := SendWACustomMessage(target_number, message)
	if err != nil {
		return err
	}

	if wa_response.Error {
		return errors.New(wa_response.Message)
	}

	return nil
}

func VerifyOTPHash(otp_code, otp_hash, secret_key string) (bool, error) {

	stored_bytes, err := hex.DecodeString(otp_hash)
	if err != nil {
		return false, err
	}

	// re calculate the new hash for inputted code
	mac := hmac.New(sha256.New, []byte(secret_key))
	mac.Write([]byte(otp_code))
	computed_hash := mac.Sum(nil)

	// compare and return false if not equal
	if !hmac.Equal(stored_bytes, computed_hash) {
		return false, errors.New("invalid OTP code") // invalid OTP code
	}

	return true, nil // valid OTP code
}
