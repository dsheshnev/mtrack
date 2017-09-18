package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var createStmt, deleteStmt, editStmt, getStmt, getByIDStmt *sql.Stmt

// wd, _ := os.Getwd()

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
	wd, _ := os.Getwd()
	db, err = sql.Open("sqlite3", fmt.Sprintf("%s/records.db", wd))
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
	editStmt, err = db.Prepare("UPDATE records SET date=?, amount=?, comment=? WHERE id=?")
	logError(err)
	getStmt, err = db.Prepare("SELECT * FROM records")
	logError(err)
	getByIDStmt, err = db.Prepare("SELECT * FROM records WHERE id=?")
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

func editRecord(rw http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		log.Fatal(err)
	}
	id, _ := strconv.Atoi(request.Form.Get("id"))
	record := getByID(id)
	date := request.Form.Get("date")
	amount, _ := strconv.ParseFloat(request.Form.Get("amount"), 64)
	comment := strings.Join(request.Form["comment"], "\n")
	if date != "" {
		date, err := time.Parse(fmt.Sprint(time.Now()), date)
		if err != nil {
			logError(err)
		} else {
			record.Date = date
		}
	}
	if amount != 0 {
		record.Amount = amount
	}
	if comment != "" {
		record.Comment = comment
	}

	res, err := editStmt.Exec(time.Now(), amount, comment, id)
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
	rows.Close()
	return
}

func getByID(id int) (record Record) {
	rows, err := getByIDStmt.Query(id)
	defer rows.Close()
	logError(err)
	for rows.Next() {
		err = rows.Scan(&record.ID, &record.Date, &record.Amount, &record.Comment)
		logError(err)
	}
	return
}

func getRecords(rw http.ResponseWriter, request *http.Request) {
	// b, err := json.Marshal(getAll())
	// logError(err)
	// rw.Write(b)
	wd, _ := os.Getwd()
	templateFiles := []string{wd + "/templates/index.html", wd + "/templates/record.html"}
	pageTmpl, err := template.ParseFiles(templateFiles...)
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
	http.HandleFunc("/edit/", editRecord)

	http.ListenAndServe(":9090", nil)
}
