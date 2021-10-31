package helloworld

import (
	"fmt"
	"os"

	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	initDatabaseConnection()
}

func initDatabaseConnection() {
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
		sslmode = "enable"
	}
	pgsqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", pgsqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connect to the database!")
}

func birthdaySave(w http.ResponseWriter, r *http.Request) {

}

func birthdayGet(w http.ResponseWriter, r *http.Request) {

}
