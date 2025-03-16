package user

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/momokii/go-sso-web/pkg/models"
)

// define generate placeholder here and not in utils for avoid circular import
func generatePlaceholders(n, startIdx int) string {
	placeholders := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "$" + strconv.Itoa(startIdx+i)
	}
	return placeholders
}

const (
	USER_MAX_DAILY_CREDIT_TOKEN = 15
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) FindByID(tx *sql.Tx, id int) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, username, password, credit_token, 
		COALESCE(last_first_llm_used::text, '') as last_first_llm_used 
		FROM users WHERE id = $1
	`

	if err := tx.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Password, &user.CreditToken, &user.LastFirstLLMUsed); err != nil && err != sql.ErrNoRows {
		return &user, err
	}

	return &user, nil
}

func (r *UserRepo) FindByUsername(tx *sql.Tx, username string) (*models.User, error) {
	var user models.User

	query := "SELECT id, username, password FROM users WHERE username = $1"

	if err := tx.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password); err != nil && err != sql.ErrNoRows {
		return &user, err
	}

	return &user, nil
}

func (r *UserRepo) Create(tx *sql.Tx, user *models.User) error {
	query := "INSERT INTO users (username, password, credit_token) VALUES ($1, $2, $3)"

	if _, err := tx.Exec(query, user.Username, user.Password, user.CreditToken); err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Update(tx *sql.Tx, user *models.User) error {
	query := "UPDATE users SET username = $1 WHERE id = $2"

	if _, err := tx.Exec(query, user.Username, user.Id); err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) UpdatePassword(tx *sql.Tx, user *models.User) error {
	query := "UPDATE users SET password = $1 WHERE id = $2"

	if _, err := tx.Exec(query, user.Password, user.Id); err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) FindUserCreditTokenForUpdate(tx *sql.Tx, user_id int) (*models.User, error) {
	var user models.User

	query := "SELECT id, username, credit_token, last_first_llm_used FROM users WHERE id = $1 FOR UPDATE"

	if err := tx.QueryRow(query, user_id).Scan(&user.Id, &user.Username, &user.CreditToken, &user.LastFirstLLMUsed); err != nil && err != sql.ErrNoRows {
		return &user, err
	}

	return &user, nil
}

func (r *UserRepo) UpdateCreditToken(tx *sql.Tx, user *models.User, newTotalCredit int, is_first_day_used bool) error {
	// Lock the row for update
	query := "SELECT id FROM users WHERE id = $1 FOR UPDATE"
	if _, err := tx.Exec(query, user.Id); err != nil {
		return err
	}

	// Update the row
	idx := 1
	param := make([]interface{}, 0)
	query = "UPDATE users SET credit_token = $" + strconv.Itoa(idx)
	param = append(param, newTotalCredit)
	idx++

	if is_first_day_used {
		time_now := time.Now().Format(time.RFC3339)
		query += ", last_first_llm_used = $" + strconv.Itoa(idx)
		param = append(param, time_now)
		idx++
	}

	query += " WHERE id = $" + strconv.Itoa(idx)
	param = append(param, user.Id)

	if _, err := tx.Exec(query, param...); err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) ResetUserDailyToken(tx *sql.Tx) error {
	// flow reset token explained:
	// 1. the flow for now is reset all user credit token to 5 if they have pending room (chat ai) and the room is still active
	// and reset all user credit token to 15 if they have used the first credit token and the last used is more than 24 hours
	// 2. the reset will be done in the same time

	var users_reset []int

	query_get_user := `
		SELECT 
			ucr.id, 
			ucr.user_id
		FROM 
			user_credit_reserved ucr 
		LEFT JOIN 
			room_credit_reserved_conn rcrc ON ucr.id = rcrc.user_credit_reserved_id
		LEFT JOIN 
			room_chat_train rct ON rcrc.room_code = rct.room_code
		WHERE
			rct.is_still_continue IS TRUE
			AND ucr.status = 'pending'
	`
	rows, err := tx.Query(query_get_user)
	if err != nil {
		return err
	}

	for rows.Next() {
		var user models.UserCreditReservedResp
		if err := rows.Scan(&user.Id, &user.UserId); err != nil {
			return err
		}
		users_reset = append(users_reset, user.UserId)
	}

	timeNowLastDay := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)

	in_query := "0" // for empty array
	if len(users_reset) > 0 {
		in_query = generatePlaceholders(len(users_reset), 4)
	}

	query := `
		UPDATE users
		SET 
			credit_token = CASE
				WHEN id IN (` + in_query + `) THEN 5
				ELSE $1
			END,
			last_first_llm_used = NULL
		WHERE last_first_llm_used IS NOT NULL 
			AND credit_token < $2 
			AND last_first_llm_used < $3
	`

	params := []interface{}{USER_MAX_DAILY_CREDIT_TOKEN, USER_MAX_DAILY_CREDIT_TOKEN, timeNowLastDay}
	if len(users_reset) > 0 {
		for _, id := range users_reset {
			params = append(params, id)
		}
	}

	if _, err := tx.Exec(query, params...); err != nil {
		return err
	}

	return nil
}
