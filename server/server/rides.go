package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ride struct {
	Driver      string `json:"driver"`
	CarName     string `json:"carName"`
	CarId       int    `json:"carId"`
	CarColor    string `json:"carColor"`
	Destination string `json:"destination"`
	Start       string `json:"start"`
	End 		string `json:"end"`
	Name 		string `json:"name"`
}

func (s *server) ridesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := s.database.Query("" +
			"SELECT driver, carName, car, carColor, destination, start, end FROM rides " +
			"JOIN cars on rides.car = cars.id")
		if err != nil {
			panic(err)
		}

		var rides []ride
		for rows.Next() {
			ride := ride{}
			err = rows.Scan(&ride.Driver, &ride.CarName, &ride.CarId, &ride.CarColor, &ride.Destination, &ride.Start, &ride.End)
			if err != nil {
				panic(err)
			}
			ride.Name = ride.Driver + " ↦ " + ride.Destination
			rides = append(rides, ride)
		}
		rideJson, _ := json.Marshal(rides)
		fmt.Fprint(w, string(rideJson))
	}
}

func test() {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		panic(err)
	}

	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec("astaxie", "mydepartment", "2012-12-09")
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println(id)
	// update
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	if err != nil {
		panic(err)
	}

	res, err = stmt.Exec("astaxieupdate", id)
	if err != nil {
		panic(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	if err != nil {
		panic(err)
	}
	var uid int
	var username string
	var department string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&uid, &username, &department, &created)
		if err != nil {
			panic(err)
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