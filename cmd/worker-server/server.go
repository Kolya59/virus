package main

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"

	"github.com/kolya59/virus/pkg/pubsub"
	pb "github.com/kolya59/virus/proto"
)

// TODO Take out
const (
	projectID = "trrp-virus"
	topicName = "machines"
)

var saveTimeout = 5 * time.Second

type server struct {
	firestore *firestore.Client
	pubsub    *pubsub.Client
}

// handleMsg handles messages from pubsub and pass it to firestore
func (s *server) handleMsg(ctx context.Context, data []byte) error {
	// Unmarshal data from pubsub to Machine
	var msg pb.Machine
	if err := proto.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal msg")
		return err
	}

	// Save Machine to firestore
	if err := s.saveMachine(ctx, msg); err != nil {
		log.Error().Err(err).Msg("Failed to save machine")
		return err
	}

	return nil
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
	srv.pubsub, err = pubsub.NewClient(projectID, topicName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize pubsub client")
	}

	// Start to listen events
	if err := srv.pubsub.Consume(ctx, srv.handleMsg); err != nil {
		log.Fatal().Err(err).Msg("Failed to handle msgs")
	}
}
