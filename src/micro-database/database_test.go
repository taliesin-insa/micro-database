package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"testing"
)

var EmptyPiFF = PiFFStruct{
	Meta: Meta{
		Type: "line",
		URL:  "",
	},
	Location: []Location{
		{Type: "line",
			Polygon: [][2]int{
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 0},
			},
			Id: "loc_0",
		},
	},
	Data: []Data{
		{
			Type:       "line",
			LocationId: "loc_0",
			Value:      "",
			Id:         "0",
		},
	},
	Children: nil,
	Parent:   0,
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	os.Setenv("MICRO_ENVIRONMENT", "test")
	Database = Client.Database("taliesin").Collection("test")
	log.Println("Started Tests")
}

func shutdown() {
	Disconnect(Client)
	log.Println("Ended Tests")
}

func TestInsert(t *testing.T) {
	p0 := PiFFStruct{
		Meta:     Meta{},
		Location: nil,
		Data:     nil,
		Children: nil,
		Parent:   0,
	}
	doc0 := Picture{primitive.NewObjectID(), p0, "/temp/none0", "", false, false, false, false, ""}

	p1 := PiFFStruct{
		Meta:     Meta{},
		Location: nil,
		Data:     nil,
		Children: nil,
		Parent:   0,
	}
	doc1 := Picture{primitive.NewObjectID(), p1, "/temp/none1", "", false, false, false, false, ""}

	tab := [2]Picture{doc0, doc1}

	b, _ := json.Marshal(tab)
	_, err := InsertMany(b, Database)
	assert.Nil(t, err)
}

func TestFindFail(t *testing.T) {
	_, err := FindOne(primitive.NewObjectID(), Database)
	assert.NotNil(t, err)
}

func TestFind(t *testing.T) {
	p0 := PiFFStruct{
		Meta:     Meta{},
		Location: nil,
		Data:     nil,
		Children: nil,
		Parent:   0,
	}

	doc0 := Picture{primitive.NewObjectID(), p0, "/temp/none0", "", false, false, false, false, ""}

	tab := [1]Picture{doc0}
	b, _ := json.Marshal(tab)
	res, _ := InsertMany(b, Database)

	id := res[0].(primitive.ObjectID)

	pic, err := FindOne(id, Database)
	assert.Nil(t, err)

	doc0.Id = id
	assert.Equal(t, doc0, pic)

}
