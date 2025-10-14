package grpcendpoint

import (
	"context"

	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	app_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/model"
	"github.com/mbilarusdev/durak_auth_bot/internal/utils"
	"github.com/mbilarusdev/durak_proto/proto/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CheckAuthEndpoint struct {
	authpb.UnimplementedAuthEndpointServer
	tokenService service.TokenManager
}

func NewGrpcCheckAuthEndpoint(tokenService service.TokenManager) *CheckAuthEndpoint {
	server := new(CheckAuthEndpoint)
	server.tokenService = tokenService
	return server
}

func (endpoint *CheckAuthEndpoint) CheckAuth(
	ctx context.Context,
	req *authpb.CheckAuthRequest,
) (*authpb.CheckAuthResponse, error) {
	token := req.GetJwt()
	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Token is empty")
	}

	actualToken, err := endpoint.tokenService.FindActualByToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	if actualToken == nil || actualToken.Status != app_model.TokenAvailable {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	resp := &authpb.CheckAuthResponse{
		Token: &authpb.Token{
			Id:       actualToken.ID,
			PlayerId: actualToken.PlayerID,
			Jwt:      actualToken.Jwt,
			Status:   utils.ConvertTokenStatus(actualToken.Status),
		},
	}
	return resp, nil
}
