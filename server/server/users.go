package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mexx77/ridesharing/logging"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Expires  time.Time `json:"expires"`
}

// TODO: load secret from elsewhere
var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
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
			{"username",  payload.Username},
			{"password",  payload.Password},
		}
		res := s.users.FindOne(context.TODO(), filter, options.FindOne())
		var user user
		err = res.Decode(&user)
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			logging.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		expirationTime := time.Now().Add(1 * time.Minute)
		claims := &Claims{
			Username: payload.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		payload.Token, err = token.SignedString(jwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			logging.Error.Printf("unable to create tokenString from token %v", token)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload.Expires = expirationTime

		userJson, _ := json.Marshal(payload)
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

		if !tokenIsValid(tknStr) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) loggedInOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tknStr := r.Header.Get("Authorization")
		if !tokenIsValid(tknStr) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

func (s *server) test() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "visible to logged-on users only")
	}
}

func tokenIsValid(tknStr string) bool {
	if tknStr == "" {
		logging.Warning.Print("No token provided")
		return false
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logging.Warning.Print("Token has invalid signature")
			return false
		}
		logging.Warning.Printf("Unable to validate token: %v", err)
		return false
	}
	if !tkn.Valid {
		logging.Warning.Print("Token is not valid")
		return false
	}
	return true
}