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
		log.Fatal().Msg(err.Error())
	}
	router := mux.NewRouter()
	router.HandleFunc("/events", getEvents).Methods("GET")
	log.Fatal().Msg(http.ListenAndServe(":80", router).Error())
}

func createLogger() {
	today := time.Now().Format(time.DateOnly)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filename := wd + "/logs/backend" + today + ".log"
	file, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	logger := ecszerolog.New(file)
	log.Logger = logger
}

func getEvents(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	events, err := readEvents()
	if err != nil {
		log.Error().Msg(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	error := json.NewEncoder(response).Encode(events)
	if error != nil {
		log.Fatal().Msg(error.Error())
	} else {
		log.Info().Msg("getEvents() was successfull")
	}
}
