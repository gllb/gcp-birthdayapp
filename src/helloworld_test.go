package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"fmt"
	"time"
	// "log"
)

func TestSaveBirthdayInFuture(t *testing.T) {
	initDatabaseConnection()
	defer db.Close()
	initDatabaseTable()

	name := "toto"
	birthday := map[string]string{"dateOfBirth": "2050-10-10"}
	jsondata, err := json.Marshal(birthday)
	if err != nil {
		panic(err)
	}
	t.Run("Save birthday in the future", func(t *testing.T) {
		url := fmt.Sprintf("/hello/%s", name)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsondata))
		if err != nil {
			panic(err)
		}
		response := httptest.NewRecorder()
		birthdaySave(response, request)

		got := response.Code
		want := http.StatusBadRequest
		if got != want {
			t.Errorf("got HTTP status %d, want %d", got, want)
		}

	})
}

func TestGetBirthday(t *testing.T) {
	initDatabaseConnection()
	defer db.Close()
	initDatabaseTable()

	name := "guillaume"
	birthday := map[string]string{"dateOfBirth": "1994-12-20"}
	jsondata, err := json.Marshal(birthday)
	if err != nil {
		panic(err)
	}

	year, month, day := time.Now().Date()
	today, err := time.Parse(shortForm, fmt.Sprintf("%d-%02d-%02d", year, month, day))
	if err != nil {
		panic(err)
	}
	birthday_time, err := time.Parse(shortForm, birthday["dateOfBirth"])
	if err != nil {
		panic(err)
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

	t.Run("Save birthday date", func(t *testing.T) {
		url := fmt.Sprintf("/hello/%s", name)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsondata))
		if err != nil {
			panic(err)
		}
		response := httptest.NewRecorder()
		birthdaySave(response, request)

		got := response.Code
		want := http.StatusNoContent
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


		got := response.Body.String()
		want := fmt.Sprintf(`{"message":"Hello, %s! Your birthday is in %d day(s)"}`, name, remainingDays)
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
