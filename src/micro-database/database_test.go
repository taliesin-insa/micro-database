package main

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestFindAll(t *testing.T) {
	coll := Client.Database("taliesin_test").Collection("test_findall")
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
	InsertMany(b, coll)

	pics, err := FindAll(coll)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(pics))
	assert.Equal(t, p0, pics[0].PiFF)
	assert.Equal(t, "/temp/none1", pics[1].Url)

}

func TestFindManyUnused(t *testing.T) {

}

func TestFindManyForSuggestion(t *testing.T) {

}

func TestUpdateFlags(t *testing.T) {

}

func TestUpdateValue(t *testing.T) {

}

func TestDeleteAll(t *testing.T) {

}

func TestStatusZero(t *testing.T) {
	Database = Client.Database("taliesin_test").Collection("test_status_zero")
	request, err := http.NewRequest("GET", "/db/status", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()
	status(recorder, request)

	statusCode := recorder.Code
	assert.Equal(t, http.StatusOK, statusCode)

	var status Status

	err = json.Unmarshal(recorder.Body.Bytes(), &status)
	assert.Nil(t, err)

	assert.True(t, status.DbUp)
	assert.Equal(t, int64(0), status.Total)
	assert.Equal(t, int64(0), status.Annotated)
	assert.Equal(t, int64(0), status.Unreadable)

}

func TestStatusTotal(t *testing.T) {
	Database = Client.Database("taliesin_test").Collection("test_status_total")

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
	InsertMany(b, Database)

	request, err := http.NewRequest("GET", "/db/status", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()
	status(recorder, request)

	statusCode := recorder.Code
	assert.Equal(t, http.StatusOK, statusCode)

	var status Status

	err = json.Unmarshal(recorder.Body.Bytes(), &status)
	assert.Nil(t, err)

	assert.True(t, status.DbUp)
	assert.Equal(t, int64(2), status.Total)
	assert.Equal(t, int64(0), status.Annotated)
	assert.Equal(t, int64(0), status.Unreadable)

}

func TestStatusAnnotated(t *testing.T) {
	Database = Client.Database("taliesin_test").Collection("test_status_annotated")

}

func TestStatusUnreadable(t *testing.T) {
	Database = Client.Database("taliesin_test").Collection("test_status_unreadable")

}
