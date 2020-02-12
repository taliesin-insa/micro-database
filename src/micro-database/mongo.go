package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

// mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[database][?options]]
const URI_SUR_CLUSTER = "mongodb://pinky.local:27017/" //TODO : ne marche pas encore à cause de la config
const URI_TESTS_LOCAUX = "mongodb://localhost:27017/"

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
	PiFF PiFFStruct `json:"PiFF"`
	// Url fileserver
	Url string `json:"Url"`
	// Flags
	Annotated  bool `json:"Annotated"`
	Corrected  bool `json:"Corrected"`
	SentToReco bool `json:"SentToReco"`
	Unreadable bool `json:"Unreadable"`
}

type Modification struct {
	Id    primitive.ObjectID `json:"Id"`
	Flag  string             `json:"flag"`
	Value bool               `json:"value"`
}

type Annotation struct {
	Id    primitive.ObjectID `json:"Id"`
	Value string             `json:"value"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Connect() *mongo.Client {
	// Set client options
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	URI := URI_TESTS_LOCAUX
	clientOptions := options.Client().ApplyURI(URI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	checkError(err)

	fmt.Println("Establishing connection to mongodb on " + URI)

	// Check the connection
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	checkError(err)

	fmt.Println("Connection successful!")

	return client
}

func Disconnect(client *mongo.Client) {
	//Disconnection
	err := client.Disconnect(context.TODO())
	checkError(err)
	fmt.Println("Connection to MongoDB closed.")
}

/**
From a json flow, insert one entry in the database -> Useless
*/
func InsertOne(b []byte, collection *mongo.Collection) error {
	var pic Picture
	var piff PiFFStruct

	err := json.Unmarshal(b, &piff)
	if err != nil {
		return err
	}

	pic.PiFF = piff

	insertResult, err := collection.InsertOne(context.TODO(), pic)
	if err != nil {
		return err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil
}

/**
From a json flow, insert multiple entries in the database
byte : Flot JSON
*/
func InsertMany(b []byte, collection *mongo.Collection) error {
	var pics []interface{}
	err := json.Unmarshal(b, &pics)
	if err != nil {
		return err
	}
	insertManyResult, err := collection.InsertMany(context.TODO(), pics)
	if err != nil {
		return err
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	return nil
}

func FindOne(id string, collection *mongo.Collection) (Picture, error) {
	filter := bson.D{{"id", id}}
	var result Picture

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return Picture{}, err
	} else {
		fmt.Printf("Found a single document: %+v\n", result)
	}

	return result, nil
}

func FindManyUnused(amount int, collection *mongo.Collection) ([]Picture, error) {
	// Pass these options to the Find method
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "Annotated", Value: false}}, bson.D{{Key: "Unreadable", Value: false}}}}}}},
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: amount}}}},
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
		fmt.Printf("Found multiple documents : %+v\n", results)
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
			log.Fatal(err)
		}
		fmt.Println(result)

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	} else {
		fmt.Printf("Found multiple documents : %+v\n", results)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	return results, nil
}

func FindList(key string, value int, collection *mongo.Collection) []Picture {

	// Pass these options to the Find method
	findOptions := options.Find()
	//findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents
	var results []Picture

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{key, value}}, findOptions)
	checkError(err)

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Picture
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Found multiple documents : %+v\n", results)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	return results
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

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}
	return nil
}

/**
Annote multiple documents.
Set the annotated flag to true
byte : Flot JSON a list of Annotation objects
*/
func UpdateValue(b []byte, collection *mongo.Collection) error {
	var annotations []Annotation
	var filter, update bson.D
	err := json.Unmarshal(b, &annotations)
	if err != nil {
		return err
	}

	for _, annot := range annotations {
		filter = bson.D{{"_id", annot.Id}}
		update = bson.D{
			{"$set", bson.D{
				{"Value", annot.Value},
				{"Annotated", true},
			}},
		}
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}

	return nil
}

/**
  Flush the database
*/
func DeleteAll(collection *mongo.Collection) {
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	checkError(err)
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

func DeleteOne(key string, value int, collection *mongo.Collection) {
	filter := bson.D{{key, value}}

	del, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Deleted %+v document\n", del.DeletedCount)
	}
}
