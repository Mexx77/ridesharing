package server

import (
	"../logging"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ride struct {
	Driver       string         `json:"driver"`
	CarName      sql.NullString `json:"carName"`
	CarId        sql.NullInt64  `json:"carId"`
	CarColor     sql.NullString `json:"carColor"`
	Destination  string         `json:"destination"`
	Start        string         `json:"start"`
	End          string         `json:"end"`
	Confirmed    bool           `json:"confirmed"`
	BigCarNeeded bool           `json:"bigCarNeeded"`
	Name         string         `json:"name"`
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

		stmt, err := s.database.Prepare("" +
			"SELECT driver, carName, car, carColor, destination, start, end, confirmed, bigCarNeeded FROM rides " +
			"LEFT JOIN cars on rides.car = cars.id where start > ? AND start < ?")
		if err != nil {
			logging.Error.Print(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query(start,end)
		if err != nil {
			logging.Error.Print(err)
		}

		rides := make([]ride, 0)
		for rows.Next() {
			ride := ride{}
			err = rows.Scan(
				&ride.Driver,
				&ride.CarName,
				&ride.CarId,
				&ride.CarColor,
				&ride.Destination,
				&ride.Start,
				&ride.End,
				&ride.Confirmed,
				&ride.BigCarNeeded)
			if err != nil {
				logging.Error.Print(err)
			}
			ride.Name = ride.Driver + " ↦ " + ride.Destination
			rides = append(rides, ride)
		}
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

		stmt, err := s.database.Prepare("INSERT INTO rides(driver, destination, start, end, bigCarNeeded)" +
			" values(?,?,?,?,?)")
		if err != nil {
			logging.Error.Print(err)
		}

		_, err = stmt.Exec(payload.Driver, payload.Destination, payload.Start, payload.End, payload.BigCarNeeded)
		if err != nil {
			logging.Error.Print(err)
		}

		logging.Debug.Println("We got this record:")
		logging.Debug.Print(body)
		w.WriteHeader(http.StatusNoContent)
	}
}

func test() {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		logging.Error.Print(err)
	}

	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	if err != nil {
		logging.Error.Print(err)
	}

	res, err := stmt.Exec("astaxie", "mydepartment", "2012-12-09")
	if err != nil {
		logging.Error.Print(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		logging.Error.Print(err)
	}

	fmt.Println(id)
	// update
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	if err != nil {
		logging.Error.Print(err)
	}

	res, err = stmt.Exec("astaxieupdate", id)
	if err != nil {
		logging.Error.Print(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		logging.Error.Print(err)
	}

	fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	if err != nil {
		logging.Error.Print(err)
	}
	var uid int
	var username string
	var department string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&uid, &username, &department, &created)
		if err != nil {
			logging.Error.Print(err)
		}
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	rows.Close() //good habit to close

	// delete
	//stmt, err = db.Prepare("delete from userinfo where uid=?")
	//checkErr(err)
	//
	//res, err = stmt.Exec(id)
	//checkErr(err)
	//
	//affect, err = res.RowsAffected()
	//checkErr(err)
	//
	//fmt.Println(affect)

	db.Close()
}
