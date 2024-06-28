package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	API_PATH = "/apis/v2/songs"
)

type songs struct {
	Name        string
	Artist_Name string
	Id          string
}

type server struct {
	dbHost string
	dbName string
	dbPass string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "bat2001"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "server"
	}

	s := server{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}

	r := mux.NewRouter()
	r.HandleFunc(apiPath, s.getSongs).Methods(http.MethodGet)
	http.ListenAndServe(":4000", r)
}

func (s server) getSongs(w http.ResponseWriter, r *http.Request) {

	db := s.openConnection()

	//read all the songs

	rows, err := db.Query("select * from songs")

	if err != nil {
		log.Fatalf("querying the songs table %s\n", err.Error())
	}

	song := []Song{}
	for rows.Next() {
		var id, name, artist_name string
		err := rows.Scan(&id, &name, &artist_name)
		if err != nil {
			log.Fatalf("while scanning the row %s\n", err.Error())
		}
		aSong := Song{
			Id:          id,
			Name:        name,
			Artist_Name: artist_name,
		}
		song = append(song, aSong)
	}
	json.NewEncoder(w).Encode(song)

	//close connection

	s.closeConnection(db)
}

func (s server) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", "root", s.dbPass, s.dbHost, s.dbName))
	if err != nil {
		log.Fatalf("opening the connection to the database: %s\n", err.Error())
	}
	return db
}

func (s server) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("closing connection %s\n", err.Error())
	}
}
