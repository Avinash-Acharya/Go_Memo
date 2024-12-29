package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Note struct {
	Id     int    `json:"id"`
	Head string `json:"head"`
	Body   string `json:"body"`
	CreatedAt string `json:"createdAt"`
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create table if not exists
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS notes (id SERIAL PRIMARY KEY, head TEXT, body TEXT, createdAt TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// create router
	router := mux.NewRouter()
	router.HandleFunc("/api/go/users", getNotes(db)).Methods("GET")
	router.HandleFunc("/api/go/users", createNote(db)).Methods("POST")
	router.HandleFunc("/api/go/users/{id}", getNote(db)).Methods("GET")
	router.HandleFunc("/api/go/users/{id}", updateNote(db)).Methods("PUT")
	router.HandleFunc("/api/go/users/{id}", deleteNote(db)).Methods("DELETE")

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	// start server
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
	fmt.Println("Server is running at 127.0.0.1:8080")
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getNotes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM notes")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		notes := []Note{} // array of notes
		for rows.Next() {
			var n Note
			if err := rows.Scan(&n.Id, &n.Head, &n.Body, &n.CreatedAt); err != nil {
				log.Fatal(err)
			}
			notes = append(notes, n)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(notes)
	}
}

func getNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		var n Note
		err := db.QueryRow("SELECT * FROM notes WHERE id = $1", id).Scan(&n.Id, &n.Head, &n.Body, &n.CreatedAt)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(n)
	}
}

func createNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var n Note
		json.NewDecoder(r.Body).Decode(&n)

		err := db.QueryRow("INSERT INTO notes (head, body, createdAt) VALUES ($1, $2, $3) RETURNING id", n.Head, n.Body, n.CreatedAt).Scan(&n.Id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(n)
	}
}

func updateNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var n Note
		json.NewDecoder(r.Body).Decode(&n)

		vars := mux.Vars(r)
		id := vars["id"]

		// Execute the update query
		_, err := db.Exec("UPDATE notes SET head = $1, body = $2, createdAt = $3 WHERE id = $4", n.Head, n.Body, n.CreatedAt, id)
		if err != nil {
			log.Fatal(err)
		}

		// Retrieve the updated user data from the database
		var updatedNote Note
		err = db.QueryRow("SELECT id, head, body, createdAt FROM notes WHERE id = $1", id).Scan(&updatedNote.Id, &updatedNote.Head, &updatedNote.Body, &updatedNote.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}

		// Send the updated user data in the response
		json.NewEncoder(w).Encode(updatedNote)
	}
}

func deleteNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var n Note
		err := db.QueryRow("SELECT * FROM notes WHERE id = $1", id).Scan(&n.Id, &n.Head, &n.Body, &n.CreatedAt)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM notes WHERE id = $1", id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("Note deleted")
		}
	}
}