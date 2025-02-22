package user

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/momokii/go-sso-web/pkg/models"
)

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
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"

	if _, err := tx.Exec(query, user.Username, user.Password); err != nil {
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

// TODO: update this function to comply with ROOM CHAT AI SIMULATION business logic
func (r *UserRepo) ResetUserDailyToken(tx *sql.Tx) error {
	time_now_last_day := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	log.Println("time_now: ", time_now_last_day)
	query := `
		UPDATE users 
		SET credit_token = $1, last_first_llm_used = NULL 
		WHERE last_first_llm_used IS NOT NULL 
		AND credit_token < $2 
		AND last_first_llm_used < $3
	`

	if _, err := tx.Exec(query, USER_MAX_DAILY_CREDIT_TOKEN, USER_MAX_DAILY_CREDIT_TOKEN, time_now_last_day); err != nil {
		return err
	}

	return nil
}
