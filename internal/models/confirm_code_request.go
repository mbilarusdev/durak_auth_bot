package models

type ConfirmCodeRequest struct {
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}
