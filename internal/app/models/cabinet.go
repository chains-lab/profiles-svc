package models

type Cabinet struct {
	Profile   Profile   `json:"profile"`
	Job       JobResume `json:"job"`
	Biography Biography `json:"biography"`
}
