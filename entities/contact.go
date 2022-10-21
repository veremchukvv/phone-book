package entities

type ContactList []Contact

type Contact struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}
