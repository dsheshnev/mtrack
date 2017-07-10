package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var createStmt, deleteStmt, editStmt, getStmt *sql.Stmt

// Record - type to store
type Record struct {
	ID      int       `json:"id"`
	Date    time.Time `json:"date"`
	Amount  float64   `json:"amount"`
	Comment string    `json:"comment"`
}

func incorrectRequest(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusBadRequest)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	io.WriteString(rw, `{"error": ""}`)
}

func logError(err error) {
	if err != nil {
		//fmt.Println(err)
		log.Println(err)
	}
}

func connectDB() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "./records.db")
	logError(err)
	ct, err := db.Prepare(`CREATE TABLE records(id integer PRIMARY KEY AUTOINCREMENT, date datetime, amount float, comment text);`)
	if err == nil {
		_, err = ct.Exec()
		logError(err)
	}
	return
}

func setStatements() {
	db, err := connectDB()
	logError(err)
	createStmt, err = db.Prepare("INSERT INTO records(date, amount, comment) values(?,?,?)")
	logError(err)
	deleteStmt, err = db.Prepare("DELETE FROM records WHERE id=?")
	logError(err)
	editStmt, err = db.Prepare("UPDATE records SET date=?, amount=?, comment=?")
	logError(err)
	getStmt, err = db.Prepare("SELECT * FROM records")
	logError(err)
}

func createRecord(rw http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		log.Fatal(err)
	}
	amount, _ := strconv.ParseFloat(request.Form["amount"][0], 64)
	comment := strings.Join(request.Form["comment"], "\n")
	res, err := createStmt.Exec(time.Now(), amount, comment)
	logError(err)
	fmt.Println(res)
	incorrectRequest(rw, err)
}

func deleteByID(id int) (err error) {
	_, err = deleteStmt.Exec(id)
	logError(err)
	return
}

func getAll() (records []Record) {
	rows, err := getStmt.Query()
	logError(err)
	var newRecord Record
	for rows.Next() {
		err = rows.Scan(&newRecord.ID, &newRecord.Date, &newRecord.Amount, &newRecord.Comment)
		logError(err)
		records = append(records, newRecord)
	}
	rows.Close() //good habit to close
	return
}

func getRecords(rw http.ResponseWriter, request *http.Request) {
	// b, err := json.Marshal(getAll())
	// logError(err)
	// rw.Write(b)
	pageTmpl, err := template.ParseFiles("./templates/index.html", "./templates/record.html")
	logError(err)
	err = pageTmpl.Execute(rw, getAll())
	logError(err)
}

func deleteRecord(rw http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		log.Fatal(err)
	}
	id, _ := strconv.ParseInt(request.Form["id"][0], 10, 64)
	err := deleteByID(int(id))
	if err != nil {
		incorrectRequest(rw, err)
	}
	var b []byte
	rw.Write(b)
}

func main() {
	setStatements()

	http.HandleFunc("/create/", createRecord)
	http.HandleFunc("/get/", getRecords)
	http.HandleFunc("/delete/", deleteRecord)
	//http.HandleFunc("/edit/", deleteRecord)

	http.ListenAndServe(":9090", nil)
}
