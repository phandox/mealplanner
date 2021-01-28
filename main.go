package main

import (
	"fmt"
	"github.com/phandox/mealplanner/internal/data"
	"github.com/phandox/mealplanner/internal/handlers"
	"github.com/phandox/mealplanner/internal/picker"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const mealsPath = "internal/data/meals.csv"

func main() {
	rand.Seed(time.Now().UnixNano())
	fd, err := os.Open(mealsPath)
	if err != nil {
		log.Fatal("can't open data source")
	}
	defer fd.Close()
	db := data.NewMealsDB(fd)
	p := picker.NewPicker(&db, picker.DefaultPeople)
	fm := template.FuncMap{
		"lower": strings.ToLower,
	}
	tableHeader := handlers.MainPageTable{
		Days:      []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		MealTypes: []string{"Breakfast", "Snack", "Lunch", "Snack", "Dinner"},
		Food:      populateMealPlan(p),
		Fm:        fm,
	}
	m := http.NewServeMux()
	m.HandleFunc("/", handlers.MainPage(tableHeader, "internal/templates/mainpage.gohtml"))
	s := http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           m,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}
	fmt.Print(s.ListenAndServe())
}

func populateMealPlan(p *picker.Picker) map[string][]data.Meal {
	var r = make(map[string][]data.Meal)
	r["breakfast"], _ = p.Plan("breakfast", 7)
	r["lunch"], _ = p.PlanLunches(7)
	r["dinner"], _ = p.Plan("dinner", 7)
	r["snack"], _ = p.Plan("snack", 7)
	return r
}
