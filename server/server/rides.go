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
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Driver       string             `json:"driver"`
	CarName      string             `json:"carName,omitempty" bson:"carName,omitempty"`
	CarColor     string             `json:"carColor,omitempty" bson:"carColor,omitempty"`
	Destination  string             `json:"destination"`
	Start        string             `json:"start"`
	End          string             `json:"end"`
	StartTime    string             `json:"startTime" bson:"startTime,omitempty"`
	EndTime      string             `json:"endTime" bson:"endTime,omitempty"`
	BigCarNeeded bool               `json:"bigCarNeeded" bson:"bigCarNeeded"`
	Date         string             `json:"date" bson:",omitempty"`
	Name         string             `json:"name" bson:"-"`
	Details      string             `json:"details" bson:"-"`
	UserId       primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
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
		payload.UserId, err = s.getUserId(r)
		if err != nil {
			logging.Error.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
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

		isAdmin := s.isAdmin(r)
		if payload.CarName != "" {
			if !isAdmin {
				payload.CarName = ""
				payload.CarColor = ""
			}
		}

		if !isAdmin {
			// Not admin, need to check if user owns ride
			userObjId, err := s.getUserId(r)
			if err != nil {
				logging.Error.Println(err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if userObjId != payload.UserId {
				logging.Warning.Printf("user %s is not authorized to update record %s", userObjId, body)
				http.Error(w, "unauthorized to update this record", http.StatusUnauthorized)
				return
			}
		}

		result, err := s.rides.ReplaceOne(context.TODO(), bson.D{{"_id", payload.Id}}, payload)
		if err != nil {
			logging.Error.Printf("Unable to update ride %v at mongo: %v", body, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logging.Debug.Printf("We updated %v: %+v", result.ModifiedCount, payload)
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

		userObjId, err := s.getUserId(r)
		if err != nil {
			logging.Error.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var ride ride
		err = s.rides.FindOne(context.TODO(), bson.D{{"_id", objId}}).Decode(&ride)
		if err != nil {
			logging.Error.Printf("Error finding ride for deletion %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		isAdmin := s.isAdmin(r)
		if userObjId == ride.UserId || isAdmin {
			_, err = s.rides.DeleteOne(context.TODO(), bson.D{{"_id", objId}})
			if err != nil {
				logging.Error.Println("unable to delete ride: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		} else {
			logging.Error.Printf("delete unauthorized: userId %s, isAdmin %t, ride %+v", userObjId, isAdmin, ride)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}

func (s *server) unconfirmedRidesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noUnconfirmed, err := s.rides.CountDocuments(context.TODO(), bson.D{{"carName", primitive.Null{}}})
		if err != nil {
			logging.Error.Printf("Error counting unconfirmed rides %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, noUnconfirmed)
	}
}

func treatRide(ride ride) ride {
	ride.Name = fmt.Sprintf("%s 🚌 %s", ride.Destination, ride.Driver)
	if ride.CarName != "" {
		ride.Details = fmt.Sprintf(`<table>
			<tr><td>Ziel</td><td>%s</td></tr>
			<tr><td>Fahrer</td><td>%s</td></tr>
			<tr><td>Auto</td><td>%s</td></tr>
			<tr><td>Startzeit</td><td>%s</td></tr>
			<tr><td>Rückkehr</td><td>%s</td></tr></table>`,
			ride.Destination,
			ride.Driver,
			ride.CarName,
			ride.StartTime,
			ride.EndTime,
		)
	} else {
		bigCarTxt := "normale Größe"
		if ride.BigCarNeeded {
			bigCarTxt = "groß"
		}
		ride.Details = fmt.Sprintf(`<table>
			<tr><td><b>Achtung</b></td><td><b>Unbestätigte Fahrt</b></td></tr>
			<tr><td>Ziel</td><td>%s</td></tr>
			<tr><td>Fahrer</td><td>%s</td></tr>
			<tr><td>Auto</td><td>%s</td></tr>
			<tr><td>Startzeit</td><td>%s</td></tr>
			<tr><td>Rückkehr</td><td>%s</td></tr></table>`,
			ride.Destination,
			ride.Driver,
			bigCarTxt,
			ride.StartTime,
			ride.EndTime,

		)
	}
	return ride
}
