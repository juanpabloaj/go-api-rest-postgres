package main

import (
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
