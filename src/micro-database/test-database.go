package main

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestAll(t *testing.T) {

	//create some trainers
	p0 := PiFFStruct{
		Meta:     Meta{},
		Location: nil,
		Data:     nil,
		Children: nil,
		Parent:   0,
	}
	doc0 := Picture{*new(primitive.ObjectID), p0, "", false, false, false, false}

	//doc1 := Picture{1, "","","","","","",false,false,false,false, false}
	//doc2 := Picture{2, "","","","","","",false,false,false,false, false}
	//doc3 := Picture{3, "","","","","","",false,false,false,false, false}
	//doc4 := `[{"Id": 4, "Type": "", "Value": "", "Children": "", "Parent": "", "Url" : "", "Annotated":false}]`	//JSON with holes

	//Test
	fmt.Println("TEST")
	b, _ := json.Marshal(doc0)
	InsertOne(b, Database)

	/*
		//insert multiple entries
		docs := []interface{}{doc1, doc2, doc3}
		bs, _ := json.Marshal(docs)
		InsertMany(bs, collection)

		InsertMany([]byte(doc4), collection)
	*/

	FindOne("0", Database) // WHY does name need to be in lowercase for the first one ???

	FindOne("1", Database)

	//FindMany("Id",4, collection)

	DeleteOne("Id", 3, Database)

	//FindMany("Value", "", collection)

	modification := `[{"Id": 2,"Flag": "Annotated","Value": true}]`

	modif := []byte(modification)

	UpdateFlags(modif, Database)

	FindOne("2", Database)

	modification = `[{"Id": 1,"Value": "This text is annotated"}]`
	modif = []byte(modification)

	UpdateValue(modif, Database)

	FindOne("1", Database)

	DeleteAll(Database)

	Disconnect(Client)
}
