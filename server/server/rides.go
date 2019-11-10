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
	"net/http"
)

type ride struct {
	Id 			 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Driver       string `json:"driver"`
	CarName      string `json:"carName,omitempty" bson:"carName,omitempty"`
	CarColor     string `json:"carColor,omitempty" bson:"carColor,omitempty"`
	Destination  string `json:"destination"`
	Start        string `json:"start"`
	End          string `json:"end"`
	StartTime    string `json:"startTime" bson:"startTime,omitempty"`
	EndTime      string `json:"endTime" bson:"endTime,omitempty"`
	BigCarNeeded bool   `json:"bigCarNeeded" bson:"bigCarNeeded"`
	Date		 string `json:"date" bson:",omitempty"`
	Name         string `json:"name" bson:"-"`
	Details      string `json:"details" bson:"-"`
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

		filter := bson.D{
			{"start", bson.D{{"$gt", start}}},
			{"end", bson.D{{"$lt", end}}},
		}
		cur, err := s.rides.Find(context.TODO(), filter)
		if err != nil {
			logging.Error.Print(err)
		}

		rides := make([]ride, 0)
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

		if payload.CarName != "" {
			logging.Info.Println("wants to set carName. Checking if admin...")
			if !s.isAdmin(r) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		var result *mongo.InsertOneResult
		result, err = s.rides.InsertOne(context.TODO(), payload)
		if err != nil {
			logging.Error.Printf("Unable writing ride %v to mongo: %v", payload, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logging.Debug.Println("We add: ", body)
		payload.Id = result.InsertedID.(primitive.ObjectID)
		payload = treatRide(payload)
		rideJson, _ := json.Marshal(payload)
		fmt.Fprint(w, string(rideJson))
	}
}

func (s *server) rideUpdateHandler() http.HandlerFunc {
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

		result, err := s.rides.ReplaceOne(context.TODO(), bson.D{{"_id", payload.Id}}, payload)
		if err != nil {
			logging.Error.Printf("Unable updating ride %v at mongo: %v", body, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logging.Debug.Printf("We updated %v: %s", result.ModifiedCount, body)
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
	ride.Name = ride.Driver + " ↦ " + ride.Destination
	if ride.CarName != "" {
		ride.Details = fmt.Sprintf(
			"Fahrer: %s</br>Ziel: %s</br>Auto: %s</br>Startzeit: %s</br>Rückkehr: %s",
			ride.Driver,
			ride.Destination,
			ride.CarName,
			ride.StartTime,
			ride.EndTime,
		)
	} else {
		bigCarTxt := "normale Größe"
		if ride.BigCarNeeded {
			bigCarTxt = "groß"
		}
		ride.Details = fmt.Sprintf(
			"<b>Achtung: Unbestätigte Fahrt</b></br>Fahrer: %s</br>Ziel: %s</br>Auto: %s</br>Startzeit: %s</br>Rückkehr: %s",
			ride.Driver,
			ride.Destination,
			bigCarTxt,
			ride.StartTime,
			ride.EndTime,

		)
	}
	return ride
}