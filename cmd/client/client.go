package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func handler(h http.Handler) {
	// Get request
	// Get id from google functions
	// Publish request to queue
	// Wait response from queue
	// Handle gotten request
}

func listenRequests(port string, callback func(w http.ResponseWriter, r *http.Request)) {
	mux := chi.NewRouter()
	// TODO Check other paths
	mux.HandleFunc("/", callback)
	s := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
	if err := s.ListenAndServe(); err != nil {

	}

}

func main() {
	// Server definition
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: mux,
	}

	// Graceful shutdown
	done := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM)
		signal.Notify(sigint, syscall.SIGINT)
		<-sigint

		select {
		case <-done:
			return
		default:
			close(done)
		}

	}()

	// Collect metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, profilerPort), nil); err != http.ErrServerClosed {
			log.Error().Err(err).Msgf("Metric server ListenAndServe")
			select {
			case <-done:
				return
			default:
				close(done)
			}
		}
	}()

	// Listen requests
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error().Err(err).Msgf("Server ListenAndServe")
			select {
			case <-done:
				return
			default:
				close(done)
			}
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Error().Err(err).Msgf("Could not shutdown server")
	}
}
