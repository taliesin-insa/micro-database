package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Client = Connect()
var Database = Client.Database("example").Collection("trainers")

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createEntry(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = InsertMany(reqBody, Database)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func selectById(w http.ResponseWriter, r *http.Request) {
	entryId := mux.Vars(r)["id"]

	entry, err := FindOne(entryId, Database)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(entry)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func newPage(w http.ResponseWriter, r *http.Request) {
	entryAmnt := mux.Vars(r)["amount"]
	amount, err := strconv.Atoi(entryAmnt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	entry, err := FindManyUnused(amount, Database)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(entry)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	entry, err := FindAll(Database)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(entry)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func updateFlags(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = UpdateFlags(reqBody, Database)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateValue(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = UpdateValue(reqBody, Database)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func status(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := Client.Ping(ctx, readpref.Primary())
	if err != nil {
		w.Write([]byte("{ 'isDBUp': false }"))
	} else {
		w.Write([]byte("{ 'isDBUp': true }"))
	}
}

func deleteAll(w http.ResponseWriter, r *http.Request) {
	DeleteAll(Database)
	w.WriteHeader(http.StatusNoContent)
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
	router.HandleFunc("/db/retrieve/all", getAll).Methods("GET")
	router.HandleFunc("/db/retrieve/snippets/{amount}", newPage).Methods("GET")
	router.HandleFunc("/db/status", status).Methods("GET")

	router.HandleFunc("/db/update/flags", updateFlags).Methods("PUT")
	router.HandleFunc("/db/update/value", updateValue).Methods("PUT")

	router.HandleFunc("/db/delete/all", deleteAll).Methods("PUT")

	//router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	//router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	//router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
