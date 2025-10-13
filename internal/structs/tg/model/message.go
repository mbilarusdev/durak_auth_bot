package tg_model

type Message struct {
	MessageID int64      `json:"message_id"`
	From      *User      `json:"from"`
	Date      int        `json:"date"`
	Chat      *Chat      `json:"chat"`
	Text      string     `json:"text"`
	Contact   *TgContact `json:"contact"`
}
