package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/kolya59/virus/models"
	"github.com/kolya59/virus/pkg/machine"
	"github.com/kolya59/virus/pkg/pubsub"
)

var (
	publishTimeout = 5 * time.Second
)

type service struct {
	dataClient     *pubsub.Client
	commandsClient *pubsub.Client
	upgrader       websocket.Upgrader
}

// Save machine handler
func (s service) SaveMachine(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Save machine")

	// Get body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read msg")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get machine from body
	var m machine.Machine
	if err := json.Unmarshal(data, &m); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal msg")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Bind external IP
	m.ExternalIP = r.RemoteAddr

	// Marshal machine
	data, err = json.Marshal(m)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal machine")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Publish raw bytes to pubsub
	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()
	if err := s.dataClient.Publish(ctx, data); err != nil {
		log.Error().Err(err).Msg("Failed to publish msg")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Subscribe to changes handler
func (s service) Subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade conn")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer c.Close()

	// Subscribe to commands topic
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := s.commandsClient.Consume(ctx, func(ctx context.Context, data []byte) (bool, error) {
		if err := c.WriteJSON(data); err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				ctx.Done()
			}
			log.Error().Err(err).Msg("Failed to write json")
			return false, err
		}

		var res models.WSAck
		if err := c.ReadJSON(&res); err != nil {
			log.Error().Err(err).Msg("Failed to read json")
			return false, err
		}

		// TODO: Think about ACK
		if res.Err != nil {
			log.Error().Err(err).Msg("Failed to do requests")
			return true, err
		}

		return true, nil
	}); err != nil {
		log.Error().Err(err).Msg("Failed to consume")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Publish commands handler
func (s service) PublishCommand(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read msg")
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	var msg models.WSCommand
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal msg")
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	if err := s.commandsClient.Publish(context.Background(), data); err != nil {
		log.Error().Err(err).Msg("Failed to publish msg")
		w.WriteHeader(http.StatusInternalServerError)
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
	dataTopicName := os.Getenv("DATA_TOPIC")
	if dataTopicName == "" {
		dataTopicName = "machines"
	}

	// Get sub name from ENV
	dataSubName := os.Getenv("DATA_SUB")
	if dataSubName == "" {
		dataSubName = "machines-sub"
	}

	// Get topic name from ENV
	commandsTopicName := os.Getenv("COM_TOPIC")
	if commandsTopicName == "" {
		commandsTopicName = "machines-command"
	}

	// Get sub name from ENV
	commandsSubName := os.Getenv("COM_SUB")
	if commandsSubName == "" {
		commandsSubName = "machines-sub-command"
	}

	// Get port from ENV
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize pub sub dataClient
	dataClient, err := pubsub.NewClient(projectID, dataTopicName, dataSubName, 5*time.Second)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize pubsub dataClient")
	}

	// Initialize pub sub commandsClient
	commandsClient, err := pubsub.NewClient(projectID, commandsTopicName, commandsSubName, 2*time.Minute)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize pubsub commandsClient")
	}

	s := service{
		dataClient:     dataClient,
		commandsClient: commandsClient,
		upgrader:       websocket.Upgrader{},
	}

	// Initialize server
	r := chi.NewRouter()
	r.Post("/machine", s.SaveMachine)
	r.HandleFunc("/subscribe", s.Subscribe)
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
