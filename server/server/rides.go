package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mexx77/ridesharing/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strings"
)

type ride struct {
	Id 			 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Driver       string `json:"driver"`
	CarName      string `json:"carName" bson:"carName"`
	CarColor     string `json:"carColor" bson:"carColor,omitempty"`
	Destination  string `json:"destination"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Confirmed    bool   `json:"confirmed"`
	BigCarNeeded bool   `json:"bigCarNeeded" bson:"bigCarNeeded"`
	IsBig        bool   `json:"isBig" bson:"isBig,omitempty"`
	Name         string `json:"name"`
	Details      string `json:"details"`
}

func (s *server) ridesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := r.URL.Query().Get("start")
		if len(start) == 0 {
			start = "1970-01-01"
		} else {
			start = start + "T00:00:00"
		}
		end := r.URL.Query().Get("end")
		if len(end) == 0 {
			end = "9999-12-31"
		} else {
			end = end + "T23:59:59"
		}

		ridesCarsPipeline := bson.A{
			bson.D{
				{"$match", bson.D{
					{"start", bson.D{{"$gt", start}}},
					{"end", bson.D{{"$lt", end}}},
				}},
			},
			bson.D{
				{"$lookup", bson.D{
					{"from", "cars"},
					{"localField", "carName"},
					{"foreignField", "carName"},
					{"as", "fromCars"},
				}},
			},
			bson.D{
				{"$replaceRoot", bson.D{{
					"newRoot", bson.D{
						{"$mergeObjects", bson.A{
							bson.D{{
								"$arrayElemAt", bson.A{"$fromCars", 0},
							}},
							"$$ROOT",
						}},
					},
				}}},
			},
			bson.D{{"$project", bson.D{{"fromCars", 0}}}},
		}

		rides := make([]ride, 0)

		cur, err := s.rides.Aggregate(context.TODO(), ridesCarsPipeline, options.Aggregate())
		if err != nil {
			logging.Error.Print(err)
		}
		for cur.Next(context.TODO()) {
			var ride ride
			err := cur.Decode(&ride)
			if err != nil {
				logging.Error.Println(err)
			}
			ride = treatRide(ride)
			rides = append(rides, ride)
		}
		if err := cur.Err(); err != nil {
			logging.Error.Println(err)
		}
		cur.Close(context.TODO())

		rideJson, _ := json.Marshal(rides)
		fmt.Fprint(w, string(rideJson))
	}
}

func (s *server) rideAddHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errorMsg := fmt.Sprintf("Invalid request method %s. POST is allowed only", r.Method)
			logging.Error.Print(errorMsg)
			http.Error(w, errorMsg, http.StatusMethodNotAllowed)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		var payload ride
		err := json.Unmarshal(buf.Bytes(), &payload)
		if err != nil {
			errorMsg := "Cannot decode payload: " + err.Error()
			logging.Error.Print(errorMsg)
			if body == "" {
				logging.Info.Print("[empty Body]")
			} else {
				logging.Info.Print(body)
			}

			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		payload.Id = primitive.NilObjectID

		var result *mongo.InsertOneResult
		result, err = s.rides.InsertOne(context.TODO(), payload)
		if err != nil {
			logging.Error.Printf("Unable writing ride %v to mongo: %v", payload, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logging.Debug.Println("We got this record: ", body)
		payload.Id = result.InsertedID.(primitive.ObjectID)
		payload = treatRide(payload)
		rideJson, _ := json.Marshal(payload)
		fmt.Fprint(w, string(rideJson))
	}
}

func (s *server) rideDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			errorMsg := "No id provided for ride deletion"
			logging.Error.Println(errorMsg)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			logging.Error.Println("unable to read ObjectID from string: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		deletePipeline := bson.D{
			{"_id", objId},
		}
		_, err = s.rides.DeleteOne(context.TODO(),deletePipeline)
		if err != nil {
			logging.Error.Println("unable to delete ride: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func treatRide(ride ride) ride {
	startSlice := strings.Split(strings.Split(ride.Start,"T")[1],":")
	timeStartStr := startSlice[0] + ":" + startSlice[1]

	ride.Name = ride.Driver + " ↦ " + ride.Destination
	if ride.Confirmed {
		ride.Details = fmt.Sprintf(
			"%s fährt mit dem %s um %s nach %s",
			ride.Driver,
			ride.CarName,
			timeStartStr,
			ride.Destination,
		)
	} else {
		bigCarTxt := ""
		if ride.BigCarNeeded {
			bigCarTxt = "mit einem großen Auto "
		}
		ride.Details = fmt.Sprintf(
			"%s möchte %sum %s nach %s fahren",
			ride.Driver,
			bigCarTxt,
			timeStartStr,
			ride.Destination,
		)
	}
	return ride
}