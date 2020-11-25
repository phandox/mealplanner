package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMeals(t *testing.T) {
	tests := []struct {
		name   string
		status int
		db     *MealDB
	}{
		{
			"no meals stored",
			http.StatusOK,
			&MealDB{},
		},
		{
			"meals stored",
			http.StatusOK,
			&MealDB{m: []Meal{{Name: "delicious meal 1"}, {Name: "delicious meal 2"}}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/meals", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetMeals(test.db))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.status)
			}

			body, err := json.Marshal(test.db.m)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Compare(rr.Body.Bytes(), body) != 0 {
				t.Errorf(
					"handler returned unexpected body: got %v want %v",
					rr.Body.String(),
					body,
				)
			}
		})
	}
}
