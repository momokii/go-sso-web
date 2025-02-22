package models

type ConnRoomCreditReserved struct {
	Id                   int    `json:"id"`
	RoomCode             string `json:"room_code"`
	UserCreditReservedId int    `json:"user_credit_reserved_id"`
	CreatedAt            string `json:"created_at"`
}
