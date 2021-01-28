package handlers

import (
	"github.com/phandox/mealplanner/internal/data"
	"html/template"
	"net/http"
	"path/filepath"
)

type MainPageTable struct {
	Days      []string
	MealTypes []string
	Food      map[string][]data.Meal
	Fm        template.FuncMap
}

func (t MainPageTable) FetchMeals(kind string) []data.Meal {
	return t.Food[kind]
}

func MainPage(th MainPageTable, tmpl string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := filepath.Base(tmpl)
		tp, err := filepath.Abs(tmpl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t, err := template.New(name).Funcs(th.Fm).ParseFiles(tp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, th)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}
