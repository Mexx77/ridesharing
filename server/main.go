package main

import (
    "./logging"
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "net/http"
    "os"
    "time"
)

const RIDES_URI = "/rides"

func main() {

    logging.Init()
    logging.Info.Print("This is ridesharing!")

    server := NewServer()
    server.startHttpServer()
}

type Server struct {
    database *sql.DB
}

func NewServer() *Server {
    const databaseFile = "./sqlite.db"
    if _, err := os.Stat(databaseFile); os.IsNotExist(err) {
        panic(err)
    }
    db, err := sql.Open("sqlite3", databaseFile)
    if err != nil {
        panic(err)
    }
    return &Server{
        database: db,
    }
}

func (s *Server) startHttpServer() {
    port := ":8080"
    srv := &http.Server{
        Addr: port,
        ErrorLog: logging.Error,
    }

    // routes
    http.HandleFunc(RIDES_URI, s.ridesHandler)

    logging.Info.Print("Starting insecure http s on port ", port)
    if err := srv.ListenAndServe(); err != nil {
        panic(err)
    }

}

type ride struct {
    Driver      string `json:"driver"`
    CarName     string `json:"carName"`
    CarId       int    `json:"carId"`
    Destination string `json:"destination"`
    Start       string `json:"start"`
}

func (s *Server) ridesHandler(w http.ResponseWriter, r *http.Request) {{
    rows, err := s.database.Query("" +
        "SELECT driver, carName, car, destination, start FROM rides " +
        "JOIN cars on rides.car = cars.id")
    if err != nil {
        panic(err)
    }

    ride := &ride{}
    for rows.Next() {
        err = rows.Scan(&ride.Driver, &ride.CarName, &ride.CarId, &ride.Destination, &ride.Start)
        if err != nil {
            panic(err)
        }
        rideJson, _ := json.Marshal(ride)
        fmt.Fprint(w, string(rideJson))
    }
}}

func test() {
    db, err := sql.Open("sqlite3", "./sqlite.db")
    if err != nil {
        panic(err)}

    // insert
    stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
    if err != nil {
        panic(err)}

    res, err := stmt.Exec("astaxie", "mydepartment", "2012-12-09")
    if err != nil {
        panic(err)}

    id, err := res.LastInsertId()
    if err != nil {
        panic(err)}

    fmt.Println(id)
    // update
    stmt, err = db.Prepare("update userinfo set username=? where uid=?")
    if err != nil {
        panic(err)}

    res, err = stmt.Exec("astaxieupdate", id)
    if err != nil {
        panic(err)}

    affect, err := res.RowsAffected()
    if err != nil {
        panic(err)}

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