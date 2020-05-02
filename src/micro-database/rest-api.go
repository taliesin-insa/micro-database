package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	lib_auth "github.com/taliesin-insa/lib-auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	_, err, authStatusCode := lib_auth.AuthenticateUser(r)

	// check if there was an error during the authentication or if the user wasn't authenticated
	if err != nil {
		log.Printf("[ERROR] Check authentication: %v", err.Error())
		w.WriteHeader(authStatusCode)
		w.Write([]byte("[MICRO-DATABASE] Couldn't verify identity"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[MICRO-DATABASE] Could not read request"))
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
		w.Write([]byte("[MICRO-DATABASE] Could not marshal answer data"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func selectById(w http.ResponseWriter, r *http.Request) {
	_, err, authStatusCode := lib_auth.AuthenticateUser(r)

	// check if there was an error during the authentication or if the user wasn't authenticated
	if err != nil {
		log.Printf("[ERROR] Check authentication: %v", err.Error())
		w.WriteHeader(authStatusCode)
		w.Write([]byte("[MICRO-DATABASE] Couldn't verify identity"))
		return
	}

	entryId, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-DATABASE] Could not decode ID"))
		return
	}

	entry, err := FindOne(entryId, Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}
	body, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-DATABASE] Could not marshal answer data"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func newPageWithSuggestions(w http.ResponseWriter, r *http.Request) {
	_, err, authStatusCode := lib_auth.AuthenticateUser(r)

	// check if there was an error during the authentication or if the user wasn't authenticated
	if err != nil {
		log.Printf("[ERROR] Check authentication: %v", err.Error())
		w.WriteHeader(authStatusCode)
		w.Write([]byte("[MICRO-DATABASE] Couldn't verify identity"))
		return
	}

	entryAmnt := mux.Vars(r)["amount"]
	amount, err := strconv.Atoi(entryAmnt)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[MICRO-DATABASE] Could not read specified amount"))
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
	body, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-DATABASE] Could not marshal answer data"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func newBatchForReco(w http.ResponseWriter, r *http.Request) {
	password := r.Header.Get("Authorization")
	expectedPassword := os.Getenv("CLUSTER_INTERNAL_PASSWORD")

	if password != expectedPassword {
		log.Printf("[ERROR] : Wrong password, expected %v but got %v", expectedPassword, password)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("[MICRO-DATABASE] Recognizer didn't have correct password"))
		return
	}

	entryAmnt := mux.Vars(r)["amount"]
	amount, err := strconv.Atoi(entryAmnt)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[MICRO-DATABASE] Could not read specified amount"))
		return
	}

	entry, err := FindManyForSuggestion(amount, Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[MICRO-DATABASE] %v", err.Error())))
		return
	}

	body, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-DATABASE] Could not marshal answer data"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

	body, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-DATABASE] Could not marshal answer data"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func updateFlags(w http.ResponseWriter, r *http.Request) {
	_, err, authStatusCode := lib_auth.AuthenticateUser(r)

	// check if there was an error during the authentication or if the user wasn't authenticated
	if err != nil {
		log.Printf("[ERROR] Check authentication: %v", err.Error())
		w.WriteHeader(authStatusCode)
		w.Write([]byte("[MICRO-DATABASE] Couldn't verify identity"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[MICRO-DATABASE] Could not read request"))
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
	password := r.Header.Get("Authorization")
	expectedPassword := os.Getenv("CLUSTER_INTERNAL_PASSWORD")

	if password != expectedPassword {
		_, err, authStatusCode := lib_auth.AuthenticateUser(r)

		// check if there was an error during the authentication or if the user wasn't authenticated
		if err != nil {
			log.Printf("[ERROR] Check authentication: %v", err.Error())
			w.WriteHeader(authStatusCode)
			w.Write([]byte("[MICRO-DATABASE] Couldn't verify identity"))
			return
		}
	}

	log.Println("Update value : ")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[MICRO-DATABASE] Could not read request"))
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
	_, err, authStatusCode := lib_auth.AuthenticateUser(r)

	// check if there was an error during the authentication or if the user wasn't authenticated
	if err != nil {
		log.Printf("[ERROR] Check authentication: %v", err.Error())
		w.WriteHeader(authStatusCode)
		w.Write([]byte("[MICRO-DATABASE] Couldn't verify identity"))
		return
	}

	annotator := mux.Vars(r)["annotator"]
	log.Println("Update value by " + annotator + " : ")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[MICRO-DATABASE] Could not read request"))
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
	_, err, authStatusCode := lib_auth.AuthenticateUser(r)

	// check if there was an error during the authentication or if the user wasn't authenticated
	if err != nil {
		log.Printf("[ERROR] Check authentication: %v", err.Error())
		w.WriteHeader(authStatusCode)
		w.Write([]byte("[MICRO-DATABASE] Couldn't verify identity"))
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	res := new(Status)
	err = Client.Ping(ctx, readpref.Primary())
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
		w.Write([]byte("[MICRO-DATABASE] Error during MongoDB counting"))
		return
	}
	res.Total = total

	annotated, err := CountAnnotatedIgnoringRecoOrUnreadable(Database)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-DATABASE] Error during MongoDB counting"))
		return
	}
	res.Annotated = annotated

	unreadable, err := CountFlag(Database, "Unreadable")
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-DATABASE] Error during MongoDB counting"))
		return
	}
	res.Unreadable = unreadable

	body, err := json.Marshal(res)
	if err != nil {
		log.Printf("[ERROR] : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-DATABASE] Could not marshal answer data"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func deleteAll(w http.ResponseWriter, r *http.Request) {
	user, err, authStatusCode := lib_auth.AuthenticateUser(r)

	// check if there was an error during the authentication or if the user wasn't authenticated
	if err != nil {
		log.Printf("[ERROR] Check authentication: %v", err.Error())
		w.WriteHeader(authStatusCode)
		w.Write([]byte("[MICRO-DATABASE] Couldn't verify identity"))
		return
	}

	// check if the authenticated user has sufficient permissions to
	if user.Role != lib_auth.RoleAdmin {
		log.Printf("[WRONG_ROLE] Insufficient permission: want %v, was %v", lib_auth.RoleAdmin, user.Role)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("[MICRO-DATABASE] Insufficient permissions to delete"))
		return
	}

	err = DeleteAll(Database)
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
