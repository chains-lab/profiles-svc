package models

type UserData struct {
	Profile Profile   `json:"profile"`
	Job     Job       `json:"job"`
	Bio     Biography `json:"biography"`
}
