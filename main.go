package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.elastic.co/ecszerolog"
)

type EventType string

const (
	Education EventType = "education"
	Skill     EventType = "skill"
	Position  EventType = "position"
)

type Event struct {
	Date  time.Time `json:"date" bson:"date"`
	Title string    `json:"title" bson:"title"`
	Event string    `json:"event" bson:"event"`
	Type  EventType `json:"type" bson:"type"`
}

func main() {
	createLogger()
	err := godotenv.Load()
	if err != nil {
		log.Error().Msg(err.Error())
	}
	router := mux.NewRouter()
	router.HandleFunc("/events", getEvents).Methods("GET")
	log.Fatal().Msg(http.ListenAndServe(":80", router).Error())
}

// createLogger creates an logger in ECS format to standard output
func createLogger() {
	logger := ecszerolog.New(os.Stdout)
	log.Logger = logger
}

// getEvents fetches events from database and returns them as json response
func getEvents(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	events, err := readEvents()
	if err != nil {
		log.Error().Msg(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(response).Encode(events)
	if err != nil {
		log.Error().Msg(err.Error())
	} else {
		log.Info().Msg("getEvents() was successfull")
	}
}
