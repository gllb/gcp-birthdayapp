package main

import (
	"fmt"
	"errors"
	"os"
	"strings"
	"time"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB
const tablename = "birthday"

const count_table_query = `SELECT COUNT(*)
FROM pg_catalog.pg_tables
WHERE tablename = $1;`
const create_table_query = `CREATE TABLE birthday (
name varchar(30),
birthday date
);`
const get_birthday_query = `SELECT birthday
FROM birthday
WHERE name = $1`
const insert_birthday_query = `INSERT INTO birthday
VALUES ($1, $2)`


type birthday struct {
	name        string
	dateOfBirth time.Time `json:dateOfBirth`
}

const shortForm = "2006-01-02"

func newBirthday(date string, name string) (*birthday, error) {
	t, _ := time.Parse(shortForm, date)
	if t.After(time.Now()) {
		return nil, errors.New("Birthday is in the future")
	}
	return &birthday{
		dateOfBirth: t,
		name: name,
	}, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	initDatabaseConnection()
	defer db.Close()
	initDatabaseTable()

	r := mux.NewRouter()
	r.HandleFunc("/birthday/{person}", birthdaySave).
		Methods("PUT").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/birthday/{person}", birthdayGet).
		Methods("GET")
	err := http.ListenAndServe(":8080", r)
	log.Fatal(err)
}

/*
   initDatabaseConnection read environment variable and connect to the configured database
*/
func initDatabaseConnection() {
	var err error

	host := os.Getenv("DBHOST")
	if host == "" {
		log.Fatal("environment variable DBHOST is empty")
	}

	port := os.Getenv("DBPORT")
	if port == "" {
		port = "5432"
	}

	dbname := os.Getenv("DBNAME")
	if dbname == "" {
		log.Fatal("environment variable DBNAME is empty")
	}

	user := os.Getenv("DBUSER")
	if user == "" {
		log.Fatal("environment variable DBUSER is empty")
	}

	password := os.Getenv("DBPASSWORD")
	if password == "" {
		log.Fatal("environment variable DBPASSWORD is empty")
	}

	sslmode := os.Getenv("DBSSLMODE")
	if sslmode == "" {
		sslmode = "require"
	}
	pgsqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err = sql.Open("postgres", pgsqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to the database!")
}
/*
   initDatabaseTable create the necessary table if it does not exist
 */
func initDatabaseTable() {
	var count int

	rows, err := db.Query(count_table_query, tablename)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	rows.Next()
	if err := rows.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		log.Printf("Create %s table", tablename)

		if _, err := db.Query(create_table_query); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("Found existing %s database", tablename)
	}
}
/*
   birthdaySave is an http.Handler:
   Read username from url and birthday from json body
   and Write a birthday related to a username in the database
 */
func birthdaySave(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR", err)
		return
	}
	log.Println(data["dateOfBirth"])
	name := strings.Split(r.URL.Path, "/")[2]
	birthday, err := newBirthday(data["dateOfBirth"], name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR", err)
		return
	}
	log.Println("Insert new birthday: ", birthday.dateOfBirth, birthday.name)

	if _, err := db.Query(insert_birthday_query, birthday.name, birthday.dateOfBirth); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/*
   birthdayGet is an http.Handler:
   Fetch a birthday in database from a username in the url request
   and write a birthday information message in the response
*/
func birthdayGet(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.Path, "/")[2]
	var dateOfBirth time.Time
	log.Println("Read birthday of", name)
	rows, err := db.Query(get_birthday_query, name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR", err)
		return
	}
	rows.Next()
	if err := rows.Scan(&dateOfBirth); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR", err)
		return
	}
	var birthday birthday
	birthday.dateOfBirth = dateOfBirth
	birthday.name = name
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR", err)
		return
	}
	year, month, day := time.Now().Date()
	today, err := time.Parse(shortForm, fmt.Sprintf("%d-%02d-%02d", year, month, day))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR", err)
		return
	}
	age := time.Now().Year() - birthday.dateOfBirth.Year()
	birthday_current_year := birthday.dateOfBirth.AddDate(age, 0, 0)

	var message string

	if today.Equal(birthday_current_year) {
		// Birthday is today
		message = fmt.Sprintf("Hello, %s! Happy birthday!", birthday.name)
	} else {
		var remainingDays int
		if today.After(birthday_current_year) {
			// Birthday is next year
			birthday_next_year := birthday_current_year.AddDate(1, 0, 0)
			remainingDays = int(birthday_next_year.Sub(time.Now()).Hours()) / 24		
		} else {
			// Birthday is this year
			remainingDays = int(birthday_current_year.Sub(time.Now()).Hours()) / 24
		}
		message = fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", birthday.name, remainingDays)
	}

	data := map[string]string{"message": message}
	jsondata, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsondata)
}

func computeRemainingDays(date string) (int, error) {
	birthday_time, err := time.Parse(shortForm, date)
	if err != nil {
		return -1, err
	}
	year, month, day := time.Now().Date()
	today, err := time.Parse(shortForm, fmt.Sprintf("%d-%02d-%02d", year, month, day))
	if err != nil {
		return -1, err
	}
	age := time.Now().Year() - birthday_time.Year()
	birthday_current_year := birthday_time.AddDate(age, 0, 0)

	var remainingDays int
	if today.After(birthday_current_year) {
		// Birthday is next year
		birthday_next_year := birthday_current_year.AddDate(1, 0, 0)
		remainingDays = int(birthday_next_year.Sub(time.Now()).Hours()) / 24		
	} else {
		// Birthday is this year
		remainingDays = int(birthday_current_year.Sub(time.Now()).Hours()) / 24
	}

	return remainingDays, nil
}
