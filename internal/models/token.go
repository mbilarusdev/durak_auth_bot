package models

const (
	TokenAvailable string = "available"
	TokenExpired   string = "expired"
	TokenBlocked   string = "blocked"
)

type Token struct {
	ID       uint64
	PlayerID uint64
	Jwt      string
	Status   string
}
