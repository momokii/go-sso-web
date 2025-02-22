package user_credit_reserved

import (
	"database/sql"

	"github.com/momokii/go-sso-web/pkg/models"
)

type UserCreditReserved struct{}

func NewUserCreditReserved() *UserCreditReserved {
	return &UserCreditReserved{}
}

// status is for room chat status (active or not)
func (r *UserCreditReserved) FindById(tx *sql.Tx, id int, status bool) (*models.UserCreditReservedResp, error) {
	var userCreditReserved models.UserCreditReservedResp

	query := `
		SELECT 
			ucr.id, 
			ucr.user_id, 
			ucr.credit, 
			ucr.feature_type, 
			ucr.status, 
			rct.is_still_continue AS is_have_room_active
		FROM 
			user_credit_reserved ucr 
		LEFT JOIN 
			room_credit_reserved_conn rcrc ON ucr.id = rcrc.user_credit_reserved_id
		LEFT JOIN 
			room_chat_train rct ON rcrc.room_code = rct.room_code
		WHERE 
			ucr.id = $1 
			AND rct.status = $2
	`

	if err := tx.QueryRow(query, id, status).Scan(&userCreditReserved.Id, &userCreditReserved.UserId, &userCreditReserved.Credit, &userCreditReserved.FeatureType, &userCreditReserved.Status, &userCreditReserved.IsHaveRoomActive); err != nil && err != sql.ErrNoRows {
		return &userCreditReserved, err
	}

	return &userCreditReserved, nil
}

func (r *UserCreditReserved) Create(tx *sql.Tx, userCreditReserved *models.UserCreditReserved) (int, error) {
	var created_id int

	query := "INSERT INTO user_credit_reserved (user_id, credit, feature_type, status) VALUES ($1, $2, $3, $4) RETURNING id"

	if err := tx.QueryRow(query, userCreditReserved.UserId, userCreditReserved.Credit, userCreditReserved.FeatureType, userCreditReserved.Status).Scan(&created_id); err != nil {
		return 0, err
	}

	return created_id, nil
}

func (r *UserCreditReserved) Update(tx *sql.Tx, userCreditReserved *models.UserCreditReserved) error {
	query := "UPDATE user_credit_reserved SET credit = $1, feature_type = $2, status = $3 WHERE id = $4"

	if _, err := tx.Exec(query, userCreditReserved.Credit, userCreditReserved.FeatureType, userCreditReserved.Status, userCreditReserved.Id); err != nil {
		return err
	}

	return nil
}
