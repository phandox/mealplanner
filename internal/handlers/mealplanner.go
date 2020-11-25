package handlers

import (
	"encoding/json"
	"net/http"
)

type Meal struct {
	Name string `json:"name"`
}

type MealDB struct {
	m []Meal
}

func (mdb *MealDB) Add(name string) {
	mdb.m = append(mdb.m, Meal{Name: name})
}

func GetMeals(db *MealDB) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := json.Marshal(db.m)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write(body)
		}
	}
}
