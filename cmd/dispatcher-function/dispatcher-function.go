package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/kolya59/virus/proto"
)

const (
	dispatcherServerHost = "127.0.0.1"
	dispatcherServerPort = "8081"
)

var (
	getNextServerTimeout = 5 * time.Second
)

type server struct {
	dispatcherServerClient pb.ServerDispatcherClient
}

func (s server) Check(context.Context, *pb.HealthCheckReq) (*pb.HealthCheckRes, error) {
	return &pb.HealthCheckRes{Status: pb.HealthCheckRes_SERVING}, nil
}

func (s server) GetTarget(ctx context.Context, req *pb.GetTargetReq) (*pb.GetTargetRes, error) {
	ctx, cancel := context.WithTimeout(ctx, getNextServerTimeout)
	defer cancel()
	res, err := s.dispatcherServerClient.GetNextServer(ctx, &pb.GetNextServerReq{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get next server")
		return nil, err
	}
	return &pb.GetTargetRes{
		Ip:          res.Ip,
		Certificate: res.Certificate,
	}, nil
}

func main() {
	// Get port from ENV
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create the channel to listen on
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	/*// TODO Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load TLS keys")
	}*/

	// Create the gRPC worker-server with the credentials
	srv := grpc.NewServer()
	// TODO srv := grpc.NewServer(grpc.Creds(creds))

	// TODO Set up a connection to the worker-server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", dispatcherServerHost, dispatcherServerPort), grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to dial dispatcher server")
	}
	defer conn.Close()

	// Initialize the client
	c := pb.NewServerDispatcherClient(conn)
	s := server{c}

	// Register the handler object
	pb.RegisterFunctionDispatcherServer(srv, s)

	// Serve and Listen
	if err := srv.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
