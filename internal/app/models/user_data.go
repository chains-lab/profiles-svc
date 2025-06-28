package models

type UserData struct {
	Profile   Profile   `json:"profile"`
	Job       Job       `json:"job"`
	Bio       Bio       `json:"biography"`
	Residence Residence `json:"residence"`
}
