package helloworld

import (
	"testing"
	"time"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"

	"fmt"
)
type Birthday struct {
	dateOfBirth time.Time
}

const shortForm = "2006-01-02"

func newBirthday(date string) Birthday {
	t, _ := time.Parse(shortForm, date)
	return Birthday{
		dateOfBirth: t,
	}
}

func (birthday *Birthday) MarshalJSON() ([]byte, error) {
	type Alias Birthday
	return json.Marshal(&struct {
		*Alias
		dateOfBirth string `json:"stamp"`
	}{
		Alias: (*Alias)(birthday),
		dateOfBirth: birthday.dateOfBirth.Format(shortForm),
	})
}

func formatBirthday(birthday Birthday) string {
	return birthday.dateOfBirth.Format(shortForm)
}

func TestGetBirthday(t *testing.T) {
	fmt.Println("hello")
	birthday := newBirthday("1994-20-12")
	jsondata, err := json.Marshal(birthday)
	if err != nil {
		panic(err)
	}

	t.Run("Save birthday date", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPut, "http://127.0.0.1/hello/guillaume", bytes.NewBuffer(jsondata))
		if err != nil {
			panic(err)
		}

		response := httptest.NewRecorder()

		birthdaySave(response, request)

		got := response.Code
		want := 200
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("returns birthday date", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/hello/guillaume", nil)
		if err != nil {
			panic(err)
		}
		response := httptest.NewRecorder()

		birthdayGet(response, request)

		// got := response.Code
		// want := 200
		// if got != want {
		// 	t.Errorf("got %q, want %q", got, want)
		// }

		got := response.Body.String()
		want := string(jsondata[:])
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
