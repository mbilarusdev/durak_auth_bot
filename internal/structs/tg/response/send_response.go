package tg_response

type SendResponse struct {
	Ok               bool           `json:"ok"`
	Result           map[string]any `json:"result"`
	ErrorDescription string         `json:"description"`
}
