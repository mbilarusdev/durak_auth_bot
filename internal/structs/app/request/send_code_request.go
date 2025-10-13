package app_request

type SendCodeRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required" example:"+79268566814"`
}
