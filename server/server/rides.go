package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mexx77/ridesharing/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type ride struct {
	Driver       string  `json:"driver"`
	CarName      string	 `json:"carName" bson:"carName"`
	Destination  string  `json:"destination"`
	Start        string  `json:"start"`
	End          string  `json:"end"`
	Confirmed    bool    `json:"confirmed"`
	BigCarNeeded bool    `json:"bigCarNeeded" bson:"bigCarNeeded"`
	Name         string  `json:"name"`
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
		}else {
			end = end + "T23:59:59"
		}

		rides := make([]ride, 0)

		filter := bson.D{
			{"start",  bson.D{{"$gt", start}}},
			{"end",    bson.D{{"$lt", end}}},
		}

		cur, err := s.rides.Find(context.TODO(), filter, options.Find())
		if err != nil {
			logging.Error.Print(err)
		}
		for cur.Next(context.TODO()) {
			var ride ride
			err := cur.Decode(&ride)
			if err != nil {
				logging.Error.Print(err)
			}
			ride.Name = ride.Driver + " ↦ " + ride.Destination
			rides = append(rides, ride)
		}
		if err := cur.Err(); err != nil {
			logging.Error.Print(err)
		}
		cur.Close(context.TODO())

		rideJson, _ := json.Marshal(rides)
		fmt.Fprint(w, string(rideJson))
	}
}

func (s *server) rideHandler() http.HandlerFunc {
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

		_, err = s.rides.InsertOne(context.TODO(), payload)
		if err != nil {
			logging.Error.Printf("Unable writing ride %v to mongo: %v", payload, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logging.Debug.Println("We got this record: ", body)
		w.WriteHeader(http.StatusNoContent)
	}
}