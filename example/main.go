package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// You will be using this Trainer type later in the program
type Trainer struct {
	Name string
	Age  int
	City string
}

func insertOne(str []byte, collection *mongo.Collection) {
	var trainer Trainer
	err := json.Unmarshal(str, &trainer)

	insertResult, err := collection.InsertOne(context.TODO(), trainer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//access the DB
	collection := client.Database("example").Collection("trainers")

	//create some trainers
	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	//Test
	fmt.Println(reflect.TypeOf(collection))
	fmt.Println("TEST")
	gregre:= Trainer{"Gregre", 21, "Rennes"}
	b, err := json.Marshal(gregre)
	insertOne(b, collection)

	//insert one
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	//insert multiple entries
	trainers := []interface{}{misty, brock}

	//actually insert
	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)


	//prepares the update of one document (using bson)
	filter := bson.D{{"name", "Ash"}}

	update := bson.D{
		{"$set", bson.D{
			{"age", 12},
		}},
	}

	//actually update
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	//Find one document
	// create a value into which the result can be decoded
	var result Trainer
	//using filter again
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)

	//Find multiple documents
	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents
	var results []*Trainer

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	//Delete
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	//Disconnection
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
