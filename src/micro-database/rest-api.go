package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var Client = Connect()
var Collection = Client.Database("example").Collection("trainers")

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createEntry(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	InsertMany(reqBody, Collection)

	w.WriteHeader(http.StatusCreated)
}

func updateFlags(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	UpdateFlags(reqBody, Collection)

	w.WriteHeader(http.StatusAccepted)
}

func updateValue(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	UpdateValue(reqBody, Collection)

	w.WriteHeader(http.StatusAccepted)
}

func selectById(w http.ResponseWriter, r *http.Request) {
	entryId := mux.Vars(r)["id"]

	entry := FindOne("Id", entryId, Collection)
	json.NewEncoder(w).Encode(entry)

	w.WriteHeader(http.StatusFound)
}

func newPage(w http.ResponseWriter, r *http.Request) {
	entryAmnt := mux.Vars(r)["amount"]
	amount, err := strconv.Atoi(entryAmnt)
	if err != nil {
		fmt.Fprintf(w, "Error while converting argument to number")
	}

	entry := FindManyUnused(amount, Collection)
	json.NewEncoder(w).Encode(entry)

	w.WriteHeader(http.StatusFound)
}

func deleteAll(w http.ResponseWriter, r *http.Request) {
	DeleteAll(Collection)
	w.WriteHeader(http.StatusAccepted)
}

//func getOneEvent(w http.ResponseWriter, r *http.Request) {
//	eventID := mux.Vars(r)["id"]
//
//	for _, singleEvent := range events {
//		if singleEvent.ID == eventID {
//			json.NewEncoder(w).Encode(singleEvent)
//		}
//	}
//}
//
//func getAllEvents(w http.ResponseWriter, r *http.Request) {
//	json.NewEncoder(w).Encode(events)
//}
//
//func updateEvent(w http.ResponseWriter, r *http.Request) {
//	eventID := mux.Vars(r)["id"]
//	var updatedEvent event
//
//	reqBody, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
//	}
//	json.Unmarshal(reqBody, &updatedEvent)
//
//	for i, singleEvent := range events {
//		if singleEvent.ID == eventID {
//			singleEvent.Title = updatedEvent.Title
//			singleEvent.Description = updatedEvent.Description
//			events = append(events[:i], singleEvent)
//			json.NewEncoder(w).Encode(singleEvent)
//		}
//	}
//}
//
//func deleteEvent(w http.ResponseWriter, r *http.Request) {
//	eventID := mux.Vars(r)["id"]
//
//	for i, singleEvent := range events {
//		if singleEvent.ID == eventID {
//			events = append(events[:i], events[i+1:]...)
//			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
//		}
//	}
//}

// Actual API
func main() {

	// Define the routing
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/db/", homeLink)
	router.HandleFunc("/db/insert", createEntry).Methods("POST")
	router.HandleFunc("/db/select/{id}", selectById).Methods("GET")
	router.HandleFunc("/db/retrieve/snippets/{amount}", newPage).Methods("GET")

	router.HandleFunc("/db/update/flags", updateFlags).Methods("PUT")
	router.HandleFunc("/db/update/value/user", updateValue).Methods("PUT")

	router.HandleFunc("/db/delete/all", deleteAll).Methods("PUT")

	//router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	//router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	//router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
