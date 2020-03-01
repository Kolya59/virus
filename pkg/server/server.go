package server

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/kolya59/virus/proto"
)

var (
	crt = "client.crt"
	key = "client.key"
)

type server struct {
}

func (s *server) register() {

}

func (s *server) startServer(host, port string) {
	// Create the channel to listen on
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load TLS keys")
	}

	// Create the gRPC server with the credentials
	srv := grpc.NewServer(grpc.Creds(creds))

	// Register the handler object
	pb.RegisterWorkerSaverServer(srv, s)

	// Serve and Listen
	if err := srv.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}

func (s *server) Check(context.Context, *pb.HealthCheckReq) (*pb.HealthCheckRes, error) {
	return &pb.HealthCheckRes{Status: pb.HealthCheckRes_SERVING}, nil
}

func (s *server) SaveMachine(context.Context, *pb.SaveMachineReq) (*pb.SaveMachineRes, error) {
	panic("implement me")
}
