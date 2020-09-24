package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

// use viper package to read .env file
// return the value of the key
func viperEnvVariable(key string) string {

  // SetConfigFile explicitly defines the path, name and extension of the config file.
  // Viper will use this and not check any of the config paths.
  // .env - It will search for the .env file in the current directory
  viper.SetConfigFile("./.env")

  // Find and read the config file
  err := viper.ReadInConfig()

  if err != nil {
    log.Fatalf("Error while reading config file %s", err)
  }

  // viper.Get() returns an empty interface{}
  // to get the underlying type of the key,
  // we have to do the type assertion, we know the underlying value is string
  // if we type assert to other type it will throw an error
  value, ok := viper.Get(key).(string)

  // If the type is a string then ok will be true
  // ok will make sure the program not break
  if !ok {
    log.Fatalf("Invalid type assertion")
  }

  return value
}


func homeLink(w http.ResponseWriter, r *http.Request) {
	viperenv := viperEnvVariable("TEST_VAR")
	fmt.Fprintf(w, "Welcome home!" + viperenv)
}

func initEvents() {
	events = allEvents{
		{
			ID:          "1",
			Title:       "Introduction to Golang",
			Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
		},
		{
			ID:          "2",
			Title:       "Introduction to Golang",
			Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
		},
		{
			ID:          "3",
			Title:       "Introduction to Golang",
			Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
		},
	}
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event id, title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)

	// Add the newly created event to the array of events
	events = append(events, newEvent)

	// Return the 201 created status code
	w.WriteHeader(http.StatusCreated)
	// Return the newly created event
	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	eventID := mux.Vars(r)["id"]

	// Get the details from an existing event
	// Use the blank identifier to avoid creating a value that will not be used
	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	eventID := mux.Vars(r)["id"]
	var updatedEvent event
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events[i] = singleEvent
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	eventID := mux.Vars(r)["id"]

	// Get the details from an existing event
	// Use the blank identifier to avoid creating a value that will not be used
	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func main() {
	initEvents()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}