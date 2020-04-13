package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type Meta struct {
	Type string
	URL  string
}

type Location struct {
	Type    string
	Polygon [][2]int
	Id      string
}

type Data struct {
	Type       string
	LocationId string
	Value      string
	Id         string
}

type PiFFStruct struct {
	Meta     Meta
	Location []Location
	Data     []Data
	Children []int
	Parent   int
}

// You will be using this Trainer type later in the program
type Picture struct {
	// Id in db
	Id primitive.ObjectID `bson:"_id" json:"Id"`
	// Piff
	PiFF     PiFFStruct `json:"PiFF"`
	Url      string     `json:"Url"`      //The URL on our fileserver
	Filename string     `json:"Filename"` //The original name of the file
	// Flags
	Annotated  bool `json:"Annotated"`
	Corrected  bool `json:"Corrected"`
	SentToReco bool `json:"SentToReco"`
	Unreadable bool `json:"Unreadable"`
	//
	Annotator string `json:"Annotator"`
}

type Modification struct {
	Id    primitive.ObjectID `json:"Id"`
	Flag  string             `json:"Flag"`
	Value bool               `json:"Value"`
}

type Annotation struct {
	Id    primitive.ObjectID `json:"Id"`
	Value string             `json:"Value"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Connect() *mongo.Client {
	URI := ""

	if os.Getenv("MICRO_ENVIRONMENT") == "production" {
		URI = "mongodb://pinky.local:27017/"
		log.Println("Started in production environment.")
	} else if os.Getenv("MICRO_ENVIRONMENT") == "dev" {
		URI = "mongodb://pinky.local:27017/"
		log.Println("Started in dev environment.")
	} else {
		URI = "mongodb://localhost:27017/"
		log.Println("Started in local environment.")
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(URI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	checkError(err)

	log.Printf("Establishing connection to mongodb on %v\n", URI)

	// Check the connection
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	checkError(err)

	log.Printf("Connection successful!\n")

	if os.Getenv("MICRO_ENVIRONMENT") == "production" {
		Database = client.Database("taliesin").Collection("prod")
	} else if os.Getenv("MICRO_ENVIRONMENT") == "dev" {
		Database = client.Database("taliesin").Collection("dev")
	} else {
		Database = client.Database("taliesin").Collection("local")
	}

	return client
}

func Disconnect(client *mongo.Client) {
	//Disconnection
	err := client.Disconnect(context.TODO())
	checkError(err)
	log.Printf("Connection to MongoDB closed.\n")
}

/**
From a json flow, insert multiple entries in the database
byte : Flot JSON
*/
func InsertMany(b []byte, collection *mongo.Collection) ([]interface{}, error) {
	var pics []interface{}
	err := json.Unmarshal(b, &pics)
	if err != nil {
		return nil, err
	}
	insertManyResult, err := collection.InsertMany(context.TODO(), pics)
	if err != nil {
		return nil, err
	}

	log.Printf("Inserted multiple documents: %v\n", insertManyResult.InsertedIDs)
	return insertManyResult.InsertedIDs, nil
}

func FindOne(id primitive.ObjectID, collection *mongo.Collection) (Picture, error) {
	filter := bson.D{{"_id", id}}
	var result Picture

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return Picture{}, err
	} else {
		log.Printf("Found a single document: %+v\n", result)
	}

	return result, nil
}

func FindManyUnused(amount int, collection *mongo.Collection) ([]Picture, error) {
	// Pass these options to the Find method
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"$and",
			bson.A{
				bson.D{{"Annotated", false}},
				bson.D{{"Unreadable", false}},
			}}}}},
		bson.D{{"$sample", bson.D{{"size", amount}}}},
	}

	opts := options.Aggregate()

	// Here's an array in which you can store the decoded documents
	var results []Picture

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Picture
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	} else {
		log.Printf("Found multiple documents : %v\n", results)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	return results, nil
}

func FindManyWithSuggestion(amount int, collection *mongo.Collection) ([]Picture, error) {
	// Pass these options to the Find method
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"$and",
			bson.A{
				bson.D{{"Annotated", true}},
				bson.D{{"Unreadable", false}},
				bson.D{{"Annotator", "$taliesin_recognizer"}},
			}}}}},
		bson.D{{"$sample", bson.D{{"size", amount}}}},
	}

	opts := options.Aggregate()

	// Here's an array in which you can store the decoded documents
	var results []Picture

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem Picture
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	} else {
		log.Printf("Found multiple documents : %v\n", results)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	return results, nil
}

func FindManyForSuggestion(amount int, collection *mongo.Collection) ([]Picture, error) {
	// Pass these options to the Find method
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"$and",
			bson.A{
				bson.D{{"Annotated", false}},
				bson.D{{"Unreadable", false}},
				bson.D{{"SentToReco", false}},
			}}}}},
		bson.D{{"$sample", bson.D{{"size", amount}}}},
	}

	opts := options.Aggregate()

	// Here's an array in which you can store the decoded documents
	var results []Picture

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Picture
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		filter := bson.D{{"_id", elem.Id}}
		update := bson.D{
			{"$set", bson.D{
				{"SentToReco", true},
			}},
		}
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	} else {
		log.Printf("Found multiple documents : %v\n", results)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	return results, nil
}

func FindAll(collection *mongo.Collection) ([]Picture, error) {
	// Pass these options to the Find method
	findOptions := options.Find()
	//findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents
	var results []Picture

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Picture
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		var result bson.D
		if err := cur.Decode(&result); err != nil {
			log.Printf("%v\n", err)
		}
		log.Printf("%v\n", result)

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	} else {
		log.Printf("Found multiple documents : %+v\n", results)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	return results, nil
}

/**
Modify the différents flags
byte : Flot JSON a list of Modification objects
*/
func UpdateFlags(b []byte, collection *mongo.Collection) error {
	var modifications []Modification
	var filter, update bson.D
	err := json.Unmarshal(b, &modifications)
	if err != nil {
		return err
	}

	for _, modif := range modifications {
		filter = bson.D{{"_id", primitive.ObjectID(modif.Id)}}
		update = bson.D{
			{"$set", bson.D{
				{modif.Flag, modif.Value},
			}},
		}
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}

		log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}
	return nil
}

/**
Annote multiple documents.
Set the annotated flag to true
byte : Flot JSON a list of Annotation objects
*/
func UpdateValue(b []byte, collection *mongo.Collection, annotator string) error {
	var annotations []Annotation
	var filter, update bson.D
	err := json.Unmarshal(b, &annotations)
	if err != nil {
		return err
	}

	log.Printf("Value : %v\n", annotations)

	for _, annot := range annotations {
		filter = bson.D{{"_id", annot.Id}}
		update = bson.D{{"$set", bson.D{
			{"PiFF.Data.0.Value", annot.Value},
			{"Annotated", true},
			{"Annotator", annotator},
		}}}
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}

	return nil
}

/**
  Flush the database
*/
func DeleteAll(collection *mongo.Collection) error {
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		return err
	}
	log.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil
}

func CountSnippets(collection *mongo.Collection) (int64, error) {
	filter := bson.D{{}}
	opts := options.Count()
	res, err := collection.CountDocuments(context.TODO(), filter, opts)
	return res, err
}

func CountFlag(collection *mongo.Collection, flag string) (int64, error) {
	filter := bson.D{{flag, true}}
	opts := options.Count()
	res, err := collection.CountDocuments(context.TODO(), filter, opts)
	return res, err
}

func CountFlagIgnoringReco(collection *mongo.Collection, flag string) (int64, error) {
	pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"$and",
		bson.A{
			bson.D{{flag, true}},
			bson.D{{"Annotator", bson.D{{"$not", bson.D{{"$regex", "\\$taliesin_recognizer"}}}}}}},
	}}}}}

	opts := options.Aggregate()

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		return -1, err
	}

	res := int64(0)

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		res++
	}

	return res, err
}
