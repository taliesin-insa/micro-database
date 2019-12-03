package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// You will be using this Trainer type later in the program
type Picture struct {
	Id int
	Type string
	Value  string
	Zone string
	Children string
	Parent string
	Url string
	Annotated bool
	Corrected bool
	SentToReco bool
	SentToUser bool
	Unreadable bool
}

type Modification struct {
	Id int
	Flag string
	Value bool
}

type Annotation struct {
	Id int
	Value string
}

/**
	From a json flow, insert one entry in the database -> Useless
	TODO : vérification des champs avant l'insertion
 */
func InsertOne(b []byte, collection *mongo.Collection) {
	var pic Picture
	err := json.Unmarshal(b, &pic)
	if err != nil {
		log.Fatal(err)
	}
	insertResult, err := collection.InsertOne(context.TODO(), pic)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

/**
From a json flow, insert multiple entries in the database
TODO : vérification des champs avant l'insertion
*/
func InsertMany(b []byte, collection *mongo.Collection) {
	var pics []interface{}
	err := json.Unmarshal(b, &pics)
	insertManyResult, err := collection.InsertMany(context.TODO(), pics)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

func FindOne(key string, value int, collection *mongo.Collection) Picture{
	filter := bson.D{{key, value}}
	var result Picture

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Found a single document: %+v\n", result)
	}
	return result
}

func FindMany(key string, value int, collection *mongo.Collection) []Picture {

	// Pass these options to the Find method
	findOptions := options.Find()
	//findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents
	var results []Picture

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{key, value}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

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
 */
func UpdateFlags(b []byte, collection *mongo.Collection) {
	var modifications []Modification
	var filter, update bson.D
	err := json.Unmarshal(b,&modifications)
	if err != nil {
		log.Fatal(err)
	}


	for _, modif := range modifications {
		filter = bson.D{{"Id",modif.Id}}
		update = bson.D{
			{"$set", bson.D{
				{modif.Flag, modif.Value},
			}},
		}
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}
}

/**
	Annote multiple documents. Take a JSON comporting the id of annotated documents and the value of annotation text.
	Set the annotated flag to true
 */
func UpdateValue(b []byte, collection *mongo.Collection) {
	var annotations []Annotation
	var filter, update bson.D
	err := json.Unmarshal(b,&annotations)
	if err != nil {
		log.Fatal(err)
	}

	for _, annot := range annotations {
		filter = bson.D{{"Id",annot.Id}}
		update = bson.D{
			{"$set", bson.D{
				{"Value", annot.Value},
				{"Annotated", true},
			}},
		}
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}
}

func DeleteAll(collection *mongo.Collection) {
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
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

func Disconnect(client *mongo.Client) {
	//Disconnection
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func Connect() *mongo.Client{
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

	return client
}

func main() {
	client := Connect()

	//access the DB
	collection := client.Database("test").Collection("Docs")

	//create some trainers
	doc0 := Picture{0, "","","","","","",false,false,false,false, false}
	doc1 := Picture{1, "","","","","","",false,false,false,false, false}
	doc2 := Picture{2, "","","","","","",false,false,false,false, false}
	doc3 := Picture{3, "","","","","","",false,false,false,false, false}
	doc4 := `[{"Id": 4, "Type": "", "Value": "", "Children": "", "Parent": "", "Url" : "", "Annotated":false}]`	//JSON with holes

	//Test
	fmt.Println("TEST")
	b, _ := json.Marshal(doc0)
	InsertOne(b, collection)

	//insert multiple entries
	docs := []interface{}{doc1, doc2, doc3}
	bs, _ := json.Marshal(docs)
	InsertMany(bs, collection)

	InsertMany([]byte(doc4), collection)

	FindOne("id", 0, collection) // WHY does name need to be in lowercase for the first one ???

	FindOne("Id",1, collection)

	FindMany("Id",4, collection)

	DeleteOne("Id", 3, collection)

	//FindMany("Value", "", collection)


	modification := `[{"Id": 2,"Flag": "Annotated","Value": true}]`

	modif := []byte(modification)

	UpdateFlags(modif, collection)

	FindOne("Id", 2, collection)

	modification = `[{"Id": 1,"Value": "This text is annotated"}]`
	modif = []byte(modification)

	UpdateValue(modif, collection)

	FindOne("Id", 1, collection)

	DeleteAll(collection)

	Disconnect(client)
}
