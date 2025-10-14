package server

import (
	"log"
	"net"

	grpcendpoint "github.com/mbilarusdev/durak_auth_bot/internal/grpc_endpoint"
	"github.com/mbilarusdev/durak_proto/proto/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	CheckAuthEndpont *grpcendpoint.CheckAuthEndpoint
}

func NewGrpcServer(checkAuthEndpoint *grpcendpoint.CheckAuthEndpoint) *GrpcServer {
	server := new(GrpcServer)
	server.CheckAuthEndpont = checkAuthEndpoint
	return server
}

func (grpcServer *GrpcServer) ListenAndServe() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	authpb.RegisterAuthEndpointServer(srv, grpcServer.CheckAuthEndpont)

	reflection.Register(srv)

	log.Printf("Starting gRPC server at %s", lis.Addr().String())

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
