package entity_models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	Request    string                 `json:"request" bson:"request"`
	Response   map[string]interface{} `json:"response,omitempty" bson:"response,omitempty"`
	LastUpdate primitive.DateTime     `json:"lastUpdate,omitempty" bson:"lastUpdate,omitmpty"`
}

type APIData struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Description  string             `json:"description" bson:"description"`
	ExternalURL  string             `json:"externalURL" bson:"externalURL"`
	DailyCount   int                `json:"dailyCount" bson:"dailyCount"`
	WeeklyCount  int                `json:"weeklyCount" bson:"weeklyCount"`
	MonthlyCount int                `json:"monthlyCount" bson:"monthlyCount"`
	YearlyCount  int                `json:"yearlyCount" bson:"yearlyCount"`
	TotalCount   int                `json:"totalCount" bson:"totalCount"`
	Categories   []string           `json:"categories" bson:"categories"`
	Base         string             `json:"base" bson:"base"`
	Requests     []*Request         `json:"requests" bson:"requests"`
}

type APIDataOverview struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Categories  []string           `json:"categories" bson:"categories"`
}

// Given a bson.M primitive from the apiData collected in MongoDB, make a new APIData struct
func NewAPIDataStruct(bsonData bson.M) (APIData, error) {
	fmt.Println(bsonData)
	var data APIData
	jsonData, err := bson.Marshal(bsonData)
	if err != nil {
		return APIData{}, err
	}
	err = bson.Unmarshal(jsonData, &data)
	fmt.Println(data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (a APIData) GetURL(index int) (string, *Request, error) {
	if index >= len(a.Requests) {
		return "", &Request{}, fmt.Errorf("Request index out of range (Min:0, Max: %v, Had: %v)", len(a.Requests)-1, index)
	}
	currentRequest := a.Requests[index]
	url := a.Base + currentRequest.Request
	return url, currentRequest, nil
}

func (r *Request) UpdateResponse(body []byte) error {
	if time.Since(r.LastUpdate.Time()).Hours() > 24 {
		var res any
		err := json.Unmarshal(body, &res)
		if err != nil {
			return err
		}
		if reflect.TypeOf(res).Kind() != reflect.Map {
			r.Response = bson.M{"data": res}
		} else {
			r.Response = res.(map[string]interface{})
		}
		r.LastUpdate = primitive.NewDateTimeFromTime(time.Now())
	}
	return nil
}
