package models

type Player struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	ChatID      int64  `json:"chat_id"`
	CreatedAt   int64  `json:"created_at"`
}
