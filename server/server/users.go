package server

import (
	"../logging"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

		w.WriteHeader(http.StatusUnauthorized)

		userJson, _ := json.Marshal(payload)
		fmt.Fprint(w, string(userJson))
	}
}