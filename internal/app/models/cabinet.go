package models

type Cabinet struct {
	Profile   Profile   `json:"profile"`
	Job       Job       `json:"job"`
	Biography Biography `json:"biography"`
}
