package utils

import (
	app_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/model"
	"github.com/mbilarusdev/durak_proto/proto/authpb"
)

func ConvertTokenStatus(status app_model.TokenStatus) authpb.TokenStatus {
	switch status {
	case app_model.TokenAvailable:
		return authpb.TokenStatus_available
	case app_model.TokenExpired:
		return authpb.TokenStatus_expired
	case app_model.TokenBlocked:
		return authpb.TokenStatus_blocked
	default:
		return authpb.TokenStatus_unknown
	}
}
