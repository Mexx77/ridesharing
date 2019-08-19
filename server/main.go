package main

import (
    "./logging"
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "net/http"
    "time"
)

const RECORDS_URI = "/records"

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
    db, err := sql.Open("sqlite3", "./sqlite.db")
    checkErr(err)
    return &Server{
        database: db,
    }
}

func (server *Server) startHttpServer() {
    port := ":8080"
    srv := &http.Server{
        Addr: port,
        ErrorLog: logging.Error,
    }

    // routes
    http.HandleFunc(RECORDS_URI, server.recordsHandler)

    logging.Info.Print("Starting insecure http server on port ", port)
    if err := srv.ListenAndServe(); err != nil {
        panic(err)
    }

}

func (server *Server) recordsHandler(w http.ResponseWriter, r *http.Request) {{
    rows, err := server.database.Query("SELECT * FROM userinfo")
    checkErr(err)
    var uid int
    var username string
    var department string
    var created time.Time

    for rows.Next() {
        err = rows.Scan(&uid, &username, &department, &created)
        checkErr(err)
        fmt.Fprintf(w, "username is %s",username)
    }
}}

func test() {
    db, err := sql.Open("sqlite3", "./sqlite.db")
    checkErr(err)

    // insert
    stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
    checkErr(err)

    res, err := stmt.Exec("astaxie", "mydepartment", "2012-12-09")
    checkErr(err)

    id, err := res.LastInsertId()
    checkErr(err)

    fmt.Println(id)
    // update
    stmt, err = db.Prepare("update userinfo set username=? where uid=?")
    checkErr(err)

    res, err = stmt.Exec("astaxieupdate", id)
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)

    fmt.Println(affect)

    // query
    rows, err := db.Query("SELECT * FROM userinfo")
    checkErr(err)
    var uid int
    var username string
    var department string
    var created time.Time

    for rows.Next() {
        err = rows.Scan(&uid, &username, &department, &created)
        checkErr(err)
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

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}