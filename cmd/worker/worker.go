package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"

	"github.com/kolya59/virus/pkg/pubsub"
)

var saveTimeout = 5 * time.Second

type server struct {
	firestore *firestore.Client
	pubsub    *pubsub.Client
}

// handleMsg handles messages from pubsub and pass it to firestore
func (s *server) handleMsg(ctx context.Context, data []byte) (bool, error) {
	// Unmarshal data from pubsub to Machine
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
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
func (s *server) saveMachine(ctx context.Context, data map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, saveTimeout)
	defer cancel()
	if _, _, err := s.firestore.Collection("machines").Add(ctx, data); err != nil {
		return err
	}

	return nil
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
