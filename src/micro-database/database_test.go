package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	lib_auth "github.com/taliesin-insa/lib-auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
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

// for mockedAuthServer
type VerifyRequest struct {
	Token string
}

func TestMain(m *testing.M) {
	// fake server to replace the authentication call (in lib_auth)
	mockedAuthServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/auth/verifyToken" {
				reqBody, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Printf("[TEST_ERROR] Create authentication mocked server (read body): %v", err.Error())
					panic(m)
				}

				var reqData VerifyRequest
				err = json.Unmarshal(reqBody, &reqData)
				if err != nil {
					log.Printf("[TEST_ERROR] Create authentication mocked server (unmarsal body): %v", err.Error())
					panic(m)
				}

				var result []byte
				if reqData.Token == "admin_token" {
					result, err = json.Marshal(lib_auth.UserData{Username: "morpheus", Role: lib_auth.RoleAdmin})
				} else {
					result, err = json.Marshal(lib_auth.UserData{Username: "morpheus", Role: lib_auth.RoleAnnotator})
				}

				if err != nil {
					log.Printf("[TEST_ERROR] Create authentication mocked server (marshal result): %v", err.Error())
					panic(m)
				}

				w.WriteHeader(http.StatusOK)
				w.Write(result)
			}
		}))

	// replace the redirect to authentication microservice
	previousAuthUrl := os.Getenv("AUTH_API_URL")
	os.Setenv("AUTH_API_URL", mockedAuthServer.URL)

	errSetup := setupDB()
	if errSetup != nil {
		log.Println("Could not drop database on test start")
	}

	code := m.Run()

	errShutdown := shutdown()
	if errShutdown != nil {
		log.Println("Could not drop database on test end")
	}

	os.Setenv("AUTH_API_URL", previousAuthUrl)

	os.Exit(code)
}

func setupDB() error {
	database := Client.Database("taliesin_test")
	err := database.Drop(context.TODO())
	if err != nil {
		return err
	}
	log.Println("Started Tests")
	return nil
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

//Note that I use the global Database when I test requests as a whole and I make a local variable when I can to limit concurrent accesses to this variable

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
	coll := Client.Database("taliesin_test").Collection("test_find_unused")
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

	picsTest1, err := FindManyUnused(1, coll)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(picsTest1))

	picsTest2, err := FindManyUnused(2, coll)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(picsTest2))

}

func TestFindManyForSuggestion(t *testing.T) {
	coll := Client.Database("taliesin_test").Collection("test_find_suggestion")
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

	picsTest1, err := FindManyForSuggestion(1, coll)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(picsTest1))

	picTest1 := picsTest1[0]

	picsTest2, err := FindManyForSuggestion(2, coll)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(picsTest2))

	picTest2 := picsTest2[0]
	assert.NotEqual(t, picTest1, picTest2)

	picsTest3, err := FindManyForSuggestion(1, coll)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(picsTest3))
}

func TestUpdateFlags(t *testing.T) {
	Database = Client.Database("taliesin_test").Collection("test_update_value")
	data := Data{}
	p0 := PiFFStruct{
		Meta:     Meta{},
		Location: nil,
		Data:     []Data{data},
		Children: nil,
		Parent:   0,
	}
	fakeid, _ := primitive.ObjectIDFromHex("face")
	doc0 := Picture{fakeid, p0, "/temp/none0", "", false, false, false, false, ""}
	p1 := PiFFStruct{
		Meta:     Meta{},
		Location: nil,
		Data:     []Data{data},
		Children: nil,
		Parent:   0,
	}
	doc1 := Picture{fakeid, p1, "/temp/none1", "", false, false, false, false, ""}
	tab := [2]Picture{doc0, doc1}
	b, _ := json.Marshal(tab)
	res, _ := InsertMany(b, Database)
	doc0.Id = res[0].(primitive.ObjectID)
	doc1.Id = res[1].(primitive.ObjectID)

	mod0 := Modification{
		Id:    doc0.Id,
		Flag:  "Unreadable",
		Value: true,
	}
	mod1 := Modification{
		Id:    doc1.Id,
		Flag:  "Corrected",
		Value: true,
	}
	tab0 := [2]Modification{mod0, mod1}
	body0, _ := json.Marshal(tab0)
	request, _ := http.NewRequest("PUT", "/db/update/flags", bytes.NewBuffer(body0))

	recorder0 := httptest.NewRecorder()
	updateFlags(recorder0, request)
	assert.Equal(t, http.StatusNoContent, recorder0.Code)

	pic, _ := FindOne(doc0.Id, Database)
	assert.False(t, pic.Annotated)
	assert.True(t, pic.Unreadable)
	assert.False(t, pic.Corrected)
	assert.False(t, pic.SentToReco)

	pic, _ = FindOne(doc1.Id, Database)
	assert.False(t, pic.Annotated)
	assert.False(t, pic.Unreadable)
	assert.True(t, pic.Corrected)
	assert.False(t, pic.SentToReco)

}

func TestUpdateValue(t *testing.T) {
	Database = Client.Database("taliesin_test").Collection("test_update_value")
	data := Data{}
	p0 := PiFFStruct{
		Meta:     Meta{},
		Location: nil,
		Data:     []Data{data},
		Children: nil,
		Parent:   0,
	}
	fakeid, _ := primitive.ObjectIDFromHex("face")
	doc0 := Picture{fakeid, p0, "/temp/none0", "", false, false, false, false, ""}
	p1 := PiFFStruct{
		Meta:     Meta{},
		Location: nil,
		Data:     []Data{data},
		Children: nil,
		Parent:   0,
	}
	doc1 := Picture{fakeid, p1, "/temp/none1", "", false, false, false, false, ""}
	tab := [2]Picture{doc0, doc1}
	b, _ := json.Marshal(tab)
	res, _ := InsertMany(b, Database)
	doc0.Id = res[0].(primitive.ObjectID)
	doc1.Id = res[1].(primitive.ObjectID)

	annot0 := Annotation{
		Id:    doc0.Id,
		Value: "Test without annotator",
	}
	tab0 := [1]Annotation{annot0}
	body0, _ := json.Marshal(tab0)
	request, _ := http.NewRequest("PUT", "/db/update/value", bytes.NewBuffer(body0))

	recorder0 := httptest.NewRecorder()
	updateValue(recorder0, request)
	assert.Equal(t, http.StatusNoContent, recorder0.Code)

	annot1 := Annotation{
		Id:    doc1.Id,
		Value: "Test with annotator",
	}
	tab1 := [1]Annotation{annot1}
	body1, _ := json.Marshal(tab1)
	request, _ = http.NewRequest("PUT", "/db/update/value/test", bytes.NewBuffer(body1))
	request = mux.SetURLVars(request, map[string]string{"annotator": "test"})

	recorder1 := httptest.NewRecorder()
	updateValueWithAnnotator(recorder1, request)
	assert.Equal(t, http.StatusNoContent, recorder1.Code)

	pic, _ := FindOne(doc0.Id, Database)
	assert.Equal(t, pic.PiFF.Data[0].Value, "Test without annotator")
	assert.True(t, pic.Annotated)

	pic, _ = FindOne(doc1.Id, Database)
	assert.Equal(t, pic.PiFF.Data[0].Value, "Test with annotator")
	assert.True(t, pic.Annotated)
	assert.Equal(t, "test", pic.Annotator)

}

func TestDeleteAll(t *testing.T) {
	coll := Client.Database("taliesin_test").Collection("test_delete_all")
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

	pics, _ := FindAll(coll)
	assert.NotEqual(t, 0, len(pics))

	err := DeleteAll(coll)
	assert.Nil(t, err)
	pics, _ = FindAll(coll)
	assert.Equal(t, 0, len(pics))

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
