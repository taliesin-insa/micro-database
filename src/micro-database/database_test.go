package main

import (
	"context"
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
	errShutdown := shutdown()
	if errShutdown != nil {
		log.Println("Could not drop database")
	}
	os.Exit(code)
}

func setup() {
	os.Setenv("MICRO_ENVIRONMENT", "test")
	log.Println("Started Tests")
}

func shutdown() error {
	database := Client.Database("taliesin_test")
	err := database.Drop(context.TODO())
	if err != nil {
		return err
	}
	Disconnect(Client)
	log.Println("Ended Tests")
	return nil
}

func TestInsert(t *testing.T) {
	coll := Client.Database("taliesin_test").Collection("test_insert")
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
	_, err := InsertMany(b, coll)
	assert.Nil(t, err)
}

func TestFindFail(t *testing.T) {
	coll := Client.Database("taliesin_test").Collection("test_find_fail")
	_, err := FindOne(primitive.NewObjectID(), coll)
	assert.NotNil(t, err)
}

func TestFind(t *testing.T) {
	coll := Client.Database("taliesin_test").Collection("test_find")
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
	res, _ := InsertMany(b, coll)

	id := res[0].(primitive.ObjectID)

	pic, err := FindOne(id, coll)
	assert.Nil(t, err)

	doc0.Id = id
	assert.Equal(t, doc0, pic)

}
