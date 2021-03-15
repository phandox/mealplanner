package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/phandox/mealplanner/internal/data"
	"github.com/phandox/mealplanner/internal/picker"
)

func populateMealPlan(p *picker.Picker) map[string][]data.Meal {
	var r = make(map[string][]data.Meal)
	r["breakfast"], _ = p.PlanRandom("breakfast", 7)
	r["lunch"], _ = p.PlanLunches(7)
	r["dinner"], _ = p.PlanRandom("dinner", 7)
	r["snack"], _ = p.PlanRandom("snack", 7)
	return r
}

type MainPageTable struct {
	Days      []string
	MealTypes []string
	Food      map[string][]data.Meal
	Fm        template.FuncMap
}

func (t MainPageTable) FetchMeals(kind string) []data.Meal {
	return t.Food[kind]
}

func MainPage(tmpl string, p *picker.Picker) func(w http.ResponseWriter, r *http.Request) {
	fm := template.FuncMap{
		"lower": strings.ToLower,
	}
	tableHeader := MainPageTable{
		Days:      []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		MealTypes: []string{"Breakfast", "Snack", "Lunch", "Snack", "Dinner"},
		Fm:        fm,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		name := filepath.Base(tmpl)
		tp, err := filepath.Abs(tmpl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t, err := template.New(name).Funcs(tableHeader.Fm).ParseFiles(tp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tableHeader.Food = populateMealPlan(p)
		err = t.Execute(w, tableHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}
