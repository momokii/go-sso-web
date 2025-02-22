package conn_room_credit_reserved

import (
	"database/sql"

	"github.com/momokii/go-sso-web/pkg/models"
)

type ConnRoomCreditReserved struct{}

func NewConnRoomCreditReserved() *ConnRoomCreditReserved {
	return &ConnRoomCreditReserved{}
}

func (r *ConnRoomCreditReserved) Create(tx *sql.Tx, connRoomCreditReserved *models.ConnRoomCreditReserved) error {
	query := "INSERT INTO room_credit_reserved_conn (room_code, user_credit_reserved_id) VALUES ($1, $2)"

	if _, err := tx.Exec(query, connRoomCreditReserved.RoomCode, connRoomCreditReserved.UserCreditReservedId); err != nil {
		return err
	}

	return nil
}
