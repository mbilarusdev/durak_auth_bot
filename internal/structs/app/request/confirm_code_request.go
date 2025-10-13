package app_request

type ConfirmCodeRequest struct {
	Username    string `json:"username"     validate:"required,min=3,max=20" example:"Vasiliy"`
	PhoneNumber string `json:"phone_number" validate:"required"              example:"+79268566814"`
	Code        string `json:"code"         validate:"required,min=6,max=6"  example:"852469"`
}
