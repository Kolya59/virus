package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/kolya59/virus/proto"
)

const (
	// TODO Fill it
	host           = ""
	port           = ""
	dispatcherHost = ""
	dispatcherPort = ""
)

var (
	crt           = "server.crt"
	key           = "server.key"
	dispatcherCrt = "dispatcher.crt"
	interval      = 5 * time.Minute
	timeout       = 60 * time.Minute
)

type server struct{}

func (s *server) register(done chan interface{}) {
	// Create the client TLS credentials
	creds, err := credentials.NewServerTLSFromFile(dispatcherCrt, "")
	if err != nil {
		return
	}

	// Set up a connection to the worker_server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", dispatcherHost, dispatcherPort), grpc.WithTransportCredentials(creds))
	if err != nil {
		return
	}
	defer conn.Close()

	// Initialize the client
	c := pb.NewServerDispatcherClient(conn)

	// Read cert
	cert, err := ioutil.ReadFile(crt)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read crt")
	}

	// Register
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Register(ctx, &pb.RegisterReq{Certificate: cert})
	if err != nil {
		return
	}

	if r.Status != pb.RegisterRes_REGISTERED {
		close(done)
	}
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
		log.Fatal().Err(err).Msg("Failed to load TLS keys")
	}

	// Create the gRPC worker_server with the credentials
	srv := grpc.NewServer(grpc.Creds(creds))

	// Register the handler object
	pb.RegisterServerWorkerServer(srv, s)

	// Serve and Listen
	if err := srv.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}

func (s *server) Check(context.Context, *pb.HealthCheckReq) (*pb.HealthCheckRes, error) {
	return &pb.HealthCheckRes{Status: pb.HealthCheckRes_SERVING}, nil
}

func (s *server) SaveMachine(context.Context, *pb.SaveMachineReq) (*pb.SaveMachineRes, error) {
	// TODO Implement it
	panic("implement me")
}

func main() {
	srv := server{}

	done := make(chan interface{})
	ticker := time.NewTicker(interval)
	timer := time.NewTimer(timeout)

loop:
	for {
		select {
		case <-ticker.C:
			srv.register(done)
		case <-timer.C:
			log.Fatal().Msgf("Register time is out")
		case <-done:
			break loop
		}
	}

	srv.startServer(host, port)
}
