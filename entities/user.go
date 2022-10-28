package entities

type User struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}
