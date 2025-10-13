package app_model

type TokenStatus string

const (
	TokenAvailable TokenStatus = "available"
	TokenExpired   TokenStatus = "expired"
	TokenBlocked   TokenStatus = "blocked"
)

type Token struct {
	ID       uint64      `json:"id"        example:"1"`
	PlayerID uint64      `json:"player_id" example:"1"`
	Jwt      string      `json:"jwt"       example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	Status   TokenStatus `json:"status"    example:"available"`
}
