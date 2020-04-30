package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	"github.com/kolya59/virus/pkg/pubsub"
)

var (
	publishTimeout = 5 * time.Second
)

type service struct {
	client *pubsub.Client
}

// Save machine handler
func (s service) SaveMachine(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Save machine")

	// Get body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal msg")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Publish raw bytes to pubsub
	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()
	if err := s.client.Publish(ctx, data); err != nil {
		log.Error().Err(err).Msg("Failed to publish msg")
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	w.WriteHeader(http.StatusOK)
}

// Check server handler
func (s service) Check(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Health check")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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

	// Initialize pub sub client
	psClient, err := pubsub.NewClient(projectID, topicName, subName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize pubsub client")
	}

	s := service{client: psClient}

	// Initialize server
	r := chi.NewRouter()
	r.Post("/machine", s.SaveMachine)
	r.Get("/health", s.Check)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r,
	}

	log.Info().Msg("Start to serve")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to listen and serve")
	}
}
