package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	APIPath = "/apis/v2/songs"
)

type Song struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ArtistName string `json:"artist_name"`
}

type Server struct {
	DBHost string
	DBName string
	DBPass string
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
		apiPath = APIPath
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "server"
	}

	s := Server{
		DBHost: dbHost,
		DBPass: dbPass,
		DBName: dbName,
	}

	r := mux.NewRouter()
	r.HandleFunc(apiPath, s.getSongs).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":4000", r))
}

func (s Server) getSongs(w http.ResponseWriter, r *http.Request) {
	db := s.openConnection()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, artist_name FROM songs")
	if err != nil {
		log.Fatalf("Error querying the songs table: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	songs := []Song{}
	for rows.Next() {
		var song Song
		err := rows.Scan(&song.ID, &song.Name, &song.ArtistName)
		if err != nil {
			log.Fatalf("Error scanning song row: %s\n", err.Error())
			continue
		}
		songs = append(songs, song)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (s Server) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(%s)/%s", s.DBPass, s.DBHost, s.DBName))
	if err != nil {
		log.Fatalf("Error opening the database connection: %s\n", err.Error())
	}
	return db
}
