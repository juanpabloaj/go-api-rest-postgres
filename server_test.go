package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "Alice", "alice@email.com").
		AddRow(2, "Bob", "bob@email.com")

	mock.ExpectQuery("^SELECT	\\* FROM users$").WillReturnRows(rows)

	server := NewServer(db)

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	server.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("http status code should be %v", http.StatusOK)
	}

	users := []User{}
	err = json.NewDecoder(rr.Body).Decode(&users)
	if err != nil {
		t.Fatal("expected users to be JSON decoded")
	}

	if len(users) != 2 {
		t.Errorf("expected users length to be 2, %v", users)
	}
	if name := users[1].Name; name != "Bob" {
		t.Errorf("got name %v, expected Bob", name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations <%s>", err)
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	newUser := User{Name: "Alice", Email: "alice@email.com"}
	jsonBody, err := json.Marshal(newUser)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(newUser.Name, newUser.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	req, err := http.NewRequest(
		"POST", "/users", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server := NewServer(db)

	server.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("got %v expected %v", rr.Code, http.StatusOK)
	}

	var returnedUser User
	err = json.NewDecoder(rr.Body).Decode(&returnedUser)
	if err != nil {
		t.Fatal(err)
	}

	if returnedUser.Name != newUser.Name {
		t.Errorf(
			"got %v expected %v, %v", returnedUser.Name, newUser.Name, returnedUser)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations <%s>", err)
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "Alice", "alice@email.com")

	mock.ExpectQuery(
		"SELECT id, name, email FROM users WHERE id = \\$1").
		WithArgs("1").
		WillReturnRows(rows)

	server := NewServer(db)

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	server.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("http status code should be %v", http.StatusOK)
	}

	var u User
	err = json.NewDecoder(rr.Body).Decode(&u)
	if err != nil {
		t.Error(err)
	}

	if u.Name != "Alice" {
		t.Errorf("got %v expected %v, %v", u.Name, "Alice", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations <%s>", err)
	}
}

func TestUpdateUser(t *testing.T) {
	t.Skip("to develop")
}

func TestDeleteUser(t *testing.T) {
	t.Skip("to develop")
}
