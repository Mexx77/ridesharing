package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mexx77/ridesharing/logging"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"regexp"
	"time"
)

const tokenExpiryTimeMinutes = 60 * 24

type user struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName     string             `json:"firstName,omitempty" bson:"firstName"`
	LastName      string             `json:"lastName,omitempty" bson:"lastName"`
	Username      string             `json:"username,omitempty" bson:",omitempty"`
	UsernamePhone string             `json:"usernamePhone,omitempty" bson:",omitempty"`
	Password      string             `json:"password,omitempty"`
	Phone         string             `json:"phone,omitempty"`
	Token         string             `json:"token" bson:",omitempty"`
	Expires       time.Time          `json:"expires" bson:",omitempty"`
	IsAdmin       bool               `json:"isAdmin" bson:"isAdmin"`
}

type Claims struct {
	UserId primitive.ObjectID `json:"userId"`
	Phone  string             `json:"phone"`
	jwt.StandardClaims
}

type Response struct {
	Message string `json:"message"`
}

func (s *server) authenticateHandler() http.HandlerFunc {
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

		var payload user
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

		filter := bson.D{
			{"$or", bson.A{
				bson.D{
					{"username", payload.UsernamePhone},
				},
				bson.D{
					{"phone", payload.UsernamePhone},
				},
			}},
		}
		res := s.users.FindOne(context.TODO(), filter, options.FindOne())
		var user user
		err = res.Decode(&user)
		if err == mongo.ErrNoDocuments {
			logging.Warning.Printf("Username/phone not found. Auth failed. %s\n", body)
			rsp := Response{Message: "Kein Benutzer mit diesem/r Benutzername/Handy-Nr. gefunden"}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusUnauthorized)
			return
		} else if err != nil {
			logging.Error.Printf("Failed to lookup user at mongo %s\n", body)
			rsp := Response{Message: "Sorry, konnte deine Daten nicht überprüfen"}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusInternalServerError)
			return
		}
		if !checkPasswordHash(payload.Password, user.Password) {
			logging.Warning.Printf("Password incorrect. Auth failed. %s\n", body)
			rsp := Response{Message: "Benutzername/Handy-Nr. und Passwort stimmen nicht überein"}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(tokenExpiryTimeMinutes * time.Minute)
		claims := &Claims{
			Phone:  user.Phone,
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		user.Token, err = token.SignedString(s.config.JwtSecret)
		if err != nil {
			logging.Error.Printf("unable to create tokenString from token %v", token)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user.Expires = expirationTime
		user.Password = "" // Don't need to send it back

		userJson, _ := json.Marshal(user)
		fmt.Fprint(w, string(userJson))
	}
}

func (s *server) refreshTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errorMsg := fmt.Sprintf("Invalid request method %s. POST is allowed only", r.Method)
			logging.Error.Println(errorMsg)
			http.Error(w, errorMsg, http.StatusMethodNotAllowed)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		oldTokenString := buf.String()

		claims, err := s.tokenIsValid(oldTokenString)
		if err != nil {
			logging.Warning.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(tokenExpiryTimeMinutes * time.Minute)
		claims.ExpiresAt = expirationTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		newTokenString, err := token.SignedString(s.config.JwtSecret)
		if err != nil {
			logging.Error.Printf("unable to create newTokenString from token %v", token)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user := user{
			Token:   newTokenString,
			Expires: expirationTime,
		}

		userJson, _ := json.Marshal(user)
		fmt.Fprint(w, string(userJson))
	}
}

func (s *server) validateTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errorMsg := fmt.Sprintf("Invalid request method %s. POST is allowed only", r.Method)
			logging.Error.Print(errorMsg)
			http.Error(w, errorMsg, http.StatusMethodNotAllowed)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		tknStr := buf.String()

		if _, err := s.tokenIsValid(tknStr); err != nil {
			logging.Warning.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) registerHandler() http.HandlerFunc {
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

		var payload user
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
		logging.Info.Println("Going to register new user:")
		logging.Info.Println(body)

		if len(payload.FirstName) < 2 {
			rsp := Response{Message: "Vorname muss mind. 2 Zeichen lang sein"}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusBadRequest)
			return
		}
		if len(payload.LastName) < 2 {
			rsp := Response{Message: "Nachname muss mind. 2 Zeichen lang sein"}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusBadRequest)
			return
		}
		if payload.Username != "" {
			if len(payload.Username) < 3 {
				rsp := Response{Message: "Benutzername muss mind. 3 Zeichen lang sein"}
				rspJson, _ := json.Marshal(rsp)
				http.Error(w, string(rspJson), http.StatusBadRequest)
				return
			}

			filter := bson.D{{"username", payload.Username}}
			var user user
			err := s.users.FindOne(context.TODO(), filter, options.FindOne()).Decode(&user)
			if err == nil {
				rsp := Response{Message: fmt.Sprintf("Der Benutzername %s existiert bereits", payload.Username)}
				rspJson, _ := json.Marshal(rsp)
				http.Error(w, string(rspJson), http.StatusBadRequest)
				return
			} else if err != mongo.ErrNoDocuments {
				logging.Error.Printf("Fehler beim Überprüfen eines neuen Benutzernamens (%s): %s\n", payload.Username, err.Error())
				rsp := Response{Message: "Sorry, Fehler beim Prüfen des Benutzernamens"}
				rspJson, _ := json.Marshal(rsp)
				http.Error(w, string(rspJson), http.StatusBadRequest)
				return
			}
		}

		if len(payload.Password) < 8 {
			rsp := Response{Message: "Passwort muss mind. 8 Zeichen lang sein"}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusBadRequest)
			return
		}
		err = validPassword(payload.Password)
		if err != nil {
			rsp := Response{Message: err.Error()}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusBadRequest)
			return
		}
		payload.Password, err = hashPassword(payload.Password)
		if err != nil {
			logging.Error.Printf("Hashing password failed: %s\n", err.Error())
			rsp := Response{Message: "Sorry, Fehler beim sicheren Speichern deines Passworts"}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusInternalServerError)
			return
		}

		phoneMatched, _ := regexp.Match(`^01[567][0-9]{8,11}$`, []byte(payload.Phone))
		if !phoneMatched {
			rsp := Response{Message: "Handy-Nr. muss im Format 01 {5,6,7} 12345678 [901] sein"}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusBadRequest)
			return
		}
		filter := bson.D{{"phone", payload.Phone}}
		var user user
		err = s.users.FindOne(context.TODO(), filter, options.FindOne()).Decode(&user)
		if err == nil {
			rsp := Response{Message: fmt.Sprintf("Die Handy-Nr. %s wird bereits von %s verwendet", payload.Phone, user.Username)}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusBadRequest)
			return
		} else if err != mongo.ErrNoDocuments {
			logging.Error.Printf("Fehler beim Überprüfen der Handy-Nr. (%s): %s\n", payload.Phone, err.Error())
			rsp := Response{Message: "Sorry, Fehler beim Überprüfen der Handy-Nr."}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusInternalServerError)
			return
		}

		_, err = s.users.InsertOne(context.TODO(), payload)
		if err != nil {
			logging.Error.Printf("Unable writing user %v to mongo: %v", payload, err)
			rsp := Response{Message: "Sorry, konnte dich nicht registieren"}
			rspJson, _ := json.Marshal(rsp)
			http.Error(w, string(rspJson), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) loggedInOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tknStr := r.Header.Get("Authorization")
		_, err := s.tokenIsValid(tknStr)
		if err != nil {
			logging.Warning.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdmin(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

func (s *server) getUserId(r *http.Request) (primitive.ObjectID, error) {
	tknStr := r.Header.Get("Authorization")
	claims, err := s.tokenIsValid(tknStr)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return claims.UserId, nil
}

func (s *server) isAdmin(r *http.Request) bool {
	tknStr := r.Header.Get("Authorization")
	claims, err := s.tokenIsValid(tknStr)
	if err != nil {
		logging.Warning.Printf("token invalid: %s", err)
		return false
	}
	var user user
	err = s.users.FindOne(context.TODO(), bson.D{{"phone", claims.Phone}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return false
	} else if err != nil {
		logging.Warning.Println("db error when checking if admin", err)
		return false
	}
	if !user.IsAdmin {
		return false
	}
	return true
}

func (s *server) tokenIsValid(tknStr string) (claims *Claims, err error) {
	if tknStr == "" {
		return nil, errors.New("no token provided")
	}

	cla := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, cla, func(token *jwt.Token) (interface{}, error) {
		return s.config.JwtSecret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("token has invalid signature")
		}
		return nil, errors.New("token is invalid " + err.Error())
	}
	if !tkn.Valid {
		return nil, errors.New("token is invalid")
	}
	return cla, nil
}
