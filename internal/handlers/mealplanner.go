package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
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

type MainPageTable struct {
	Days      []string
	MealTypes []string
}

func MainPage(th MainPageTable, tmpl string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := filepath.Base(tmpl)
		tp, err := filepath.Abs(tmpl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t, err := template.New(name).ParseFiles(tp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, th)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}
