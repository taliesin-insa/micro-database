package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Database *mongo.Collection
var Client = Connect()

type Status struct {
	DbUp       bool  `json:"isDBUp"`
	Total      int64 `json:"total"`
	Annotated  int64 `json:"annotated"`
	Unreadable int64 `json:"unreadable"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	log.Printf("Homelink Joined")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[MICRO-DATABASE] Homelink Joined"))
}

func createEntry(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	ids, err := InsertMany(reqBody, Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	body, err := json.Marshal(ids)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func selectById(w http.ResponseWriter, r *http.Request) {
	entryId, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	entry, err := FindOne(entryId, Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(entry)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func newPageWithSuggestions(w http.ResponseWriter, r *http.Request) {
	entryAmnt := mux.Vars(r)["amount"]
	amount, err := strconv.Atoi(entryAmnt)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
	}

	entry, err := FindManyWithSuggestion(amount, Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	if len(entry) < amount {
		unsused, err := FindManyUnused(amount-len(entry), Database)
		if err != nil {
			log.Printf("[ERROR] : %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
			return
		}
		for _, pic := range unsused {
			entry = append(entry, pic)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(entry)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func newBatchForReco(w http.ResponseWriter, r *http.Request) {
	entryAmnt := mux.Vars(r)["amount"]
	amount, err := strconv.Atoi(entryAmnt)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
	}

	entry, err := FindManyForSuggestion(amount, Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
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
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
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
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	err = UpdateFlags(reqBody, Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateValue(w http.ResponseWriter, r *http.Request) {
	log.Println("Update value : ")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	err = UpdateValue(reqBody, Database, "unspecified")
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateValueWithAnnotator(w http.ResponseWriter, r *http.Request) {
	annotator := mux.Vars(r)["annotator"]
	log.Println("Update value by " + annotator + " : ")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	err = UpdateValue(reqBody, Database, annotator)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func status(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	res := new(Status)
	err := Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.Write([]byte("{ 'isDBUp': false }"))
		return
	} else {
		res.DbUp = true
	}

	total, err := CountSnippets(Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}
	res.Total = total

	annotated, err := CountFlag(Database, "Annotated")
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}
	res.Annotated = annotated

	unreadable, err := CountFlag(Database, "Unreadable")
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}
	res.Unreadable = unreadable

	body, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func deleteAll(w http.ResponseWriter, r *http.Request) {
	err := DeleteAll(Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
	}
	w.WriteHeader(http.StatusOK)
}

// Actual API
func main() {

	// Define the routing
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/db/", homeLink).Methods("GET")

	router.HandleFunc("/db/select/{id}", selectById).Methods("GET")
	router.HandleFunc("/db/retrieve/all", getAll).Methods("GET")
	router.HandleFunc("/db/retrieve/snippets/{amount}", newPageWithSuggestions).Methods("GET")
	router.HandleFunc("/db/retrieve/recognizer/{amount}", newBatchForReco).Methods("GET")
	router.HandleFunc("/db/status", status).Methods("GET")

	router.HandleFunc("/db/insert", createEntry).Methods("POST")

	router.HandleFunc("/db/update/flags", updateFlags).Methods("PUT")
	router.HandleFunc("/db/update/value", updateValue).Methods("PUT")
	router.HandleFunc("/db/update/value/{annotator}", updateValueWithAnnotator).Methods("PUT")

	router.HandleFunc("/db/delete/all", deleteAll).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
