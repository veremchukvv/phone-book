package entities

type FriendsList []Friend

type Friend struct {
	UserID         int    `json:"user_id"`
	RelationUserID int    `json:"relation_user_id"`
	PhoneNumber    string `json:"phone_number"`
}
