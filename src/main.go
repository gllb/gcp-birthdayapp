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

type birthday struct {
	name        string
	dateOfBirth time.Time `json:dateOfBirth`
}

const shortForm = "2006-01-02"

// func newBirthday(date date, name string) (*birthday, error) {
// 	return &birthday{
// 		dateOfBirth: t,
// 		name: name,
// 	}, nil
// }

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

// func (birthday *Birthday) MarshalJSON() ([]byte, error) {
// 	log.Println("MarshalJSON")
// 	type Alias Birthday
// 	return json.Marshal(&struct {
// 		*Alias
// 		dateOfBirth string `json:"dateOfBirth"`
// 	}{
// 		Alias: (*Alias)(birthday),
// 		dateOfBirth: birthday.dateOfBirth.Format(shortForm),
// 	})
// }

func formatBirthday(birthday birthday) string {
	return birthday.dateOfBirth.Format(shortForm)
}


func main() {
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

func initDatabaseConnection() {
	var err error

	host := os.Getenv("DBHOST")
	if host == "" {
		panic("environment variable DBHOST is empty")
	}

	port := os.Getenv("DBPORT")
	if port == "" {
		port = "5432"
	}

	dbname := os.Getenv("DBNAME")
	if dbname == "" {
		panic("environment variable DBNAME is empty")
	}

	user := os.Getenv("DBUSER")
	if user == "" {
		panic("environment variable DBUSER is empty")
	}

	password := os.Getenv("DBPASSWORD")
	if password == "" {
		panic("environment variable DBPASSWORD is empty")
	}

	sslmode := os.Getenv("DBSSLMODE")
	if sslmode == "" {
		sslmode = "require"
	}
	pgsqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err = sql.Open("postgres", pgsqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to the database!")
}

func initDatabaseTable() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var count int
	const count_table_query = `SELECT COUNT(*)
FROM pg_catalog.pg_tables
WHERE tablename = $1;`

	rows, err := db.Query(count_table_query, tablename)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	rows.Next()
	if err := rows.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		log.Printf("Create %s table", tablename)
		const create_table_query = `CREATE TABLE birthday (
name varchar(30),
birthday date
);`
		if _, err := db.Query(create_table_query); err != nil {
			log.Fatal(err)
		}
	// 	columns, _ := rows.Columns()
	// 	log.Println(strings.Join(columns[:], ", "))
	// 	rows.Next()
	// 	if err := rows.Scan(&count); err != nil {
	// 		log.Fatal(err)
	// 	}
	} else {
		log.Printf("Found existing %s database", tablename)
	}
}

func birthdaySave(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("ERROR:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(data["dateOfBirth"])
	name := strings.Split(r.URL.Path, "/")[2]
	birthday, err := newBirthday(data["dateOfBirth"], name)
	if err != nil {
		log.Println("ERROR:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Insert new birthday: ", birthday.dateOfBirth, birthday.name)
	const insert_birthday_query = `INSERT INTO birthday
VALUES ($1, $2)`
	if _, err := db.Query(insert_birthday_query, birthday.name, birthday.dateOfBirth); err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func birthdayGet(w http.ResponseWriter, r *http.Request) {
	const get_birthday_query = `SELECT birthday
FROM birthday
WHERE name = $1`
	name := strings.Split(r.URL.Path, "/")[2]
	var dateOfBirth time.Time
	log.Println("Read birthday of", name)
	rows, err := db.Query(get_birthday_query, name)
	if err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	rows.Next()
	if err := rows.Scan(&dateOfBirth); err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var birthday birthday
	birthday.dateOfBirth = dateOfBirth
	birthday.name = name
	if err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	age := time.Now().Year() - birthday.dateOfBirth.Year()
	birthday_current_year := birthday.dateOfBirth.AddDate(age, 0, 0)
	remainingDays := int(birthday_current_year.Sub(time.Now()).Hours()) / 24

	var message string
	// TODO: Handle the case of a birthday already passed this year (remainingDays < 0)
	if remainingDays < 1 && birthday.dateOfBirth.Day() == time.Now().Day() {
		message = fmt.Sprintf("Hello, %s! Happy birthday!", birthday.name)
	} else {
		message = fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", birthday.name, remainingDays)
	}
	data := map[string]string{"message": message}
	jsondata, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsondata)
}
