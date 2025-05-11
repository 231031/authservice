package authservice

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/231031/authservice/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	errLog  *log.Logger
	pb.UnimplementedAuthServiceServer
}

func ListenGRPC(s Service, port int) error {
	address := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	l := log.New(nil, "error from service server : ", log.Ldate|log.LUTC|log.Lshortfile)
	srv := grpc.NewServer()
	grpcService := &grpcServer{
		service: s,
		errLog:  l,
	}

	pb.RegisterAuthServiceServer(srv, grpcService)
	reflection.Register(srv)

	err = srv.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *grpcServer) PostAuth(ctx context.Context, in *pb.PostAuthRequest) (*pb.PostAuthResponse, error) {
	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	return nil, ErrUnauthorized
	// }

	return &pb.PostAuthResponse{}, nil
}
