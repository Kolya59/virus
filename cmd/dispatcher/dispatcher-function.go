package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/kolya59/virus/pkg/pubsub"
	pb "github.com/kolya59/virus/proto"
)

type server struct {
	client *pubsub.Client
}

func (s server) SaveMachine(ctx context.Context, req *pb.SaveMachineReq) (*pb.SaveMachineRes, error) {
	// Unmarshal Machine from gRPC msg to raw bytes
	data, err := proto.Marshal(req.Machine)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal msg")
		return nil, err
	}

	// Publish raw bytes to pubsub
	if err := s.client.Publish(ctx, data); err != nil {
		log.Error().Err(err).Msg("Failed to publish msg")
		return nil, err
	}

	return &pb.SaveMachineRes{Status: pb.SaveMachineRes_ACCEPTED}, nil
}

func (s server) Check(context.Context, *pb.HealthCheckReq) (*pb.HealthCheckRes, error) {
	return &pb.HealthCheckRes{Status: pb.HealthCheckRes_SERVING}, nil
}

func main() {
	// Get project ID from ENV
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		projectID = "trrp-virus"
	}

	// Get topic name from ENV
	topicName := os.Getenv("TOPIC")
	if topicName == "" {
		topicName = "machines"
	}

	// Get sub name from ENV
	subName := os.Getenv("SUB")
	if subName == "" {
		subName = "machines-sub"
	}

	log.Info().Msgf("ProjectID: %v Topic: %v Sub: %v", projectID, topicName, subName)

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

	// Create the gRPC worker with the credentials
	srv := grpc.NewServer()
	// TODO srv := grpc.NewServer(grpc.Creds(creds))

	psClient, err := pubsub.NewClient(projectID, topicName, subName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize pubsub client")
	}

	s := server{client: psClient}

	// Register the handler object
	pb.RegisterFunctionDispatcherServer(srv, s)

	log.Info().Msg("Start to serve")
	// Serve and Listen
	if err := srv.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
