package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	server := NewServer(db)

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting ...")
	log.Fatal(http.ListenAndServe(":8080", jsonContentTypeMiddleware(server.Router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type Server struct {
	DB     *sql.DB
	Router *mux.Router
}

func NewServer(db *sql.DB) *Server {
	server := new(Server)

	router := mux.NewRouter()
	router.HandleFunc("/users", server.GetUsers).Methods("GET")
	router.HandleFunc("/users", server.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", server.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", server.UpdateUser).Methods("PUT")

	server.DB = db
	server.Router = router

	return server
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Print(err)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.DB.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		u.Name, u.Email)

	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var u User
	err := s.DB.QueryRow(
		"SELECT id, name, email FROM users WHERE id = $1", userID).
		Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err = s.DB.Exec(stmt, u.Name, u.Email, userID)
	if err != nil {
		log.Print(err)
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Print(err)
	}
}
