package user_otp

import (
	"database/sql"

	"github.com/momokii/go-sso-web/pkg/models"
)

type UserOTPRepo struct{}

func NewUserOTPRepo() *UserOTPRepo {
	return &UserOTPRepo{}
}

func (r *UserOTPRepo) GetNewest(tx *sql.Tx, id int, purpose models.OtpPurpose) (*models.UserOtps, error) {
	query := "SELECT id, user_id, purpose, channel, code_hash, created_at, expires_at, used, COALESCE(used_at::text, '') as used_at FROM user_otps WHERE user_id = $1 AND used = false AND purpose = $2 ORDER BY created_at DESC LIMIT 1"

	row := tx.QueryRow(query, id, purpose)

	otp := &models.UserOtps{}
	if err := row.Scan(&otp.Id, &otp.UserId, &otp.Purpose, &otp.Channel, &otp.CodeHash, &otp.CreatedAt, &otp.ExpiresAt, &otp.Used, &otp.UsedAt); err != nil {
		return nil, err
	}

	return otp, nil
}

func (r *UserOTPRepo) Create(tx *sql.Tx, otp_input *models.UserOtps) error {

	query := "INSERT INTO user_otps (user_id, purpose, channel, code_hash, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	if _, err := tx.Exec(query, otp_input.UserId, otp_input.Purpose, otp_input.Channel, otp_input.CodeHash, otp_input.ExpiresAt); err != nil {
		return err
	}

	return nil
}

func (r *UserOTPRepo) Update(tx *sql.Tx, otp_update *models.UserOtps) error {

	query := "UPDATE user_otps SET used = $1, used_at = $2 WHERE id = $3 AND user_id = $4"

	if _, err := tx.Exec(query, otp_update.Used, otp_update.UsedAt, otp_update.Id, otp_update.UserId); err != nil {
		return err
	}

	return nil
}

func (r *UserOTPRepo) DeletesByUserId(tx *sql.Tx, otp_data *models.UserOtps) error {

	query := "DELETE FROM user_otps WHERE user_id = $1 and purpose = $2"

	if _, err := tx.Exec(query, otp_data.UserId, otp_data.Purpose); err != nil {
		return err
	}

	return nil
}
