package main

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"

	"github.com/kolya59/virus/pkg/pubsub"
	pb "github.com/kolya59/virus/proto"
)

var saveTimeout = 5 * time.Second

type server struct {
	firestore *firestore.Client
	pubsub    *pubsub.Client
}

// handleMsg handles messages from pubsub and pass it to firestore
func (s *server) handleMsg(ctx context.Context, data []byte) (bool, error) {
	// Unmarshal data from pubsub to Machine
	var msg pb.Machine
	if err := proto.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal msg")
		return true, err
	}

	// Save Machine to firestore
	if err := s.saveMachine(ctx, msg); err != nil {
		log.Error().Err(err).Msg("Failed to save machine")
		return false, err
	}

	return true, nil
}

// saveMachine saves machine's to firestore
func (s *server) saveMachine(ctx context.Context, data pb.Machine) error {
	converted := convertData(data)

	ctx, cancel := context.WithTimeout(ctx, saveTimeout)
	defer cancel()
	if _, _, err := s.firestore.Collection("users").Add(ctx, converted); err != nil {
		return err
	}

	return nil
}

// convertData converts pb.Machine to map
func convertData(raw pb.Machine) map[string]interface{} {
	res := make(map[string]interface{})
	// TODO Check
	res["machine"] = raw
	return res
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

	srv := server{}

	// Get a Firestore firestore.
	ctx := context.Background()
	var err error
	srv.firestore, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create firestore")
	}
	// Close firestore when done.
	defer srv.firestore.Close()

	// Initialize pubsub client
	srv.pubsub, err = pubsub.NewClient(projectID, topicName, subName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize pubsub client")
	}

	log.Info().Msg("Start to consume")
	// Start to listen events
	if err := srv.pubsub.Consume(ctx, srv.handleMsg); err != nil {
		log.Fatal().Err(err).Msg("Failed to handle msgs")
	}
}
