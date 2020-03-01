package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	roundrobin "github.com/kolya59/virus/pkg/round-robin"
	pb "github.com/kolya59/virus/proto"
)

const (
	crt       = "server.crt"
	key       = "server.key"
	host      = "127.0.0.1"
	port      = "8080"
	storePath = "./data"
)

type server struct {
	rr    roundrobin.RoundRobin
	certs map[string][]byte
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

	// Create the gRPC worker_server with the credentials
	srv := grpc.NewServer(grpc.Creds(creds))

	// Register the handler object
	pb.RegisterServerDispatcherServer(srv, s)

	// Serve and Listen
	if err := srv.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}

func (s *server) Check(context.Context, *pb.HealthCheckReq) (*pb.HealthCheckRes, error) {
	return &pb.HealthCheckRes{Status: pb.HealthCheckRes_SERVING}, nil
}

func (s *server) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	var ip string
	if md, ok := peer.FromContext(ctx); ok {
		ip = md.Addr.String()
	} else {
		log.Error().Msg("Failed to parse client ip")
		return nil, fmt.Errorf("failed to resolve ip")
	}

	next, err := url.Parse(ip)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse next url")
		return nil, fmt.Errorf("failed to parse next url")
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", storePath, ip), req.GetCertificate(), 0644); err != nil {
		log.Error().Msg("Failed to save client cert")
		return nil, fmt.Errorf("failed to save client cert")
	}

	s.certs[next.String()] = req.Certificate
	s.rr.Add(*next)

	return &pb.RegisterRes{Status: pb.RegisterRes_REGISTERED}, nil
}

func (s *server) GetNextServer(context.Context, *pb.GetNextServerReq) (*pb.GetNextServerRes, error) {
	next, err := s.rr.Next()
	if err != nil {
		log.Error().Msg("Failed to get next client")
		return nil, fmt.Errorf("failed to get next client")
	}

	return pb.GetNextServerRes{
		Ip:                   next.String(),
		Certificate:          nil,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}, nil
}

func main() {
	rr := roundrobin.NewRoundRobin(healthCheck)
}

func healthCheck(target url.URL) bool {
	_, err := http.Get(target.String())
	log.Printf("Server: %v Available: %v", target.String(), err != nil)
	return err != nil
}
