package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	roundrobin "github.com/kolya59/virus/pkg/round-robin"
)

const (
	host = "127.0.0.1"
	port = "8080"
)

func main() {
	rr := roundrobin.NewRoundRobin(healthCheck)

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	r.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			next, err := rr.Next()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if _, err := w.Write([]byte(next.String())); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case "POST":
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Failed to read body: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			next, err := url.Parse(string(data))
			if err != nil {
				log.Printf("Failed to parse next url: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			rr.Add(*next)

			w.WriteHeader(http.StatusOK)
		default:
			log.Fatalf("Unexpected method %v", r.Method)
		}
	}).Methods("GET", "POST")

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(fmt.Errorf("failed to listen and serve: %w", err))
	}
}

func healthCheck(target url.URL) bool {
	_, err := http.Get(target.String())
	log.Printf("Server: %v Available: %v", target.String(), err != nil)
	return err != nil
}
