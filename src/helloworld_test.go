package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"fmt"
	// "log"
)

func TestGetBirthday(t *testing.T) {
	initDatabaseConnection()
	defer db.Close()
	initDatabaseTable()

	name := "guillaume"
	birthday := map[string]string{"dateOfBirth": "1994-12-20" }
	jsondata, err := json.Marshal(birthday)
	if err != nil {
		panic(err)
	}

	t.Run("Save birthday date", func(t *testing.T) {
		url := fmt.Sprintf("/hello/%s", name)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsondata))
		if err != nil {
			panic(err)
		}
		response := httptest.NewRecorder()
		birthdaySave(response, request)

		got := response.Code
		want := 204
		if got != want {
			t.Errorf("got HTTP status %d, want %d", got, want)
		}
	})
	t.Run("returns birthday date", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/hello/guillaume", nil)
		if err != nil {
			panic(err)
		}
		response := httptest.NewRecorder()

		birthdayGet(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("got %d, want %d", response.Code, http.StatusOK)
		}


		remainingDays := 47
		got := response.Body.String()
		want := fmt.Sprintf(`{"message":"Hello, %s! Your birthday is in %d day(s)"}`, name, remainingDays)
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
