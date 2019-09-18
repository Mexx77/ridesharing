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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

const tokenExpiryTimeMinutes = 60

type user struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token"`
	Expires  time.Time `json:"expires"`
}

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
			logging.Warning.Printf("Authentication failed for %s\n", body)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			logging.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		expirationTime := time.Now().Add(tokenExpiryTimeMinutes * time.Minute)
		claims := &Claims{
			Username: payload.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		payload.Token, err = token.SignedString(s.config.JwtSecret)
		if err != nil {
			logging.Error.Printf("unable to create tokenString from token %v", token)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload.Expires = expirationTime
		payload.Password = "" // Don't need to send it back

		userJson, _ := json.Marshal(payload)
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
			Username: claims.Username,
			Token:    newTokenString,
			Expires:  expirationTime,
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

func (s *server) loggedInOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tknStr := r.Header.Get("Authorization")
		if _, err := s.tokenIsValid(tknStr); err != nil {
			logging.Warning.Println(err)
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