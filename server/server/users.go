package server

import (
	"../logging"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

		logging.Debug.Print(body)

		stmt, err := s.database.Prepare("SELECT isAdmin FROM users WHERE username = ? AND password = ?")
		if err != nil {
			logging.Error.Print(err)
		}
		defer stmt.Close()
		var isAdmin bool
		err = stmt.QueryRow(payload.Username, payload.Password).Scan(&isAdmin)
		if err == sql.ErrNoRows {
			logging.Warning.Printf("Login request with wrong credentials %v", body)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			logging.Error.Print(err)
			return
		}

		expirationTime := time.Now().Add(60 * time.Minute)
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