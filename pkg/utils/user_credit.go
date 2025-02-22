package utils

import (
	"database/sql"
	"errors"

	"github.com/momokii/go-sso-web/pkg/models"
	"github.com/momokii/go-sso-web/pkg/repository/user"
)

func UpdateUserCredit(tx *sql.Tx, user_repo user.UserRepo, user *models.User, feature_cost int) error {
	var is_first_day_used bool

	if user.CreditToken < feature_cost || user.CreditToken-feature_cost <= 0 {
		return errors.New("user credit is not enough")
	}

	if user.LastFirstLLMUsed == "" {
		is_first_day_used = true
	} else {
		is_first_day_used = false
	}

	updated_credit := user.CreditToken - feature_cost

	if err := user_repo.UpdateCreditToken(tx, user, updated_credit, is_first_day_used); err != nil {
		return err
	}

	return nil
}
