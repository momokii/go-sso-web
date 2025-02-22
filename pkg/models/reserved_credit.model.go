package models

type UserCreditReserved struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Credit      int    `json:"credit"`
	FeatureType string `json:"feature_type"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type UserCreditReservedResp struct {
	UserCreditReserved
	IsHaveRoomActive bool `json:"is_have_room_active"`
}
