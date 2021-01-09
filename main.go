package main

import (
	"errors"
	"fmt"
	"github.com/phandox/mealplanner/internal/data"
	"github.com/phandox/mealplanner/internal/handlers"
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
	fm := template.FuncMap{
		"chooseMeal": ChooseMeal,
		"lower":      strings.ToLower,
	}
	tableHeader := handlers.MainPageTable{
		Days:      []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		MealTypes: []string{"Breakfast", "Snack", "Lunch", "Snack", "Dinner"},
		Food:      populateMealPlan(db),
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

func populateMealPlan(db data.MealsDB) map[string][]data.Meal {
	var r = make(map[string][]data.Meal)
	r["breakfast"] = db.Meals("breakfast")
	r["lunch"] = db.Meals("lunch")
	r["dinner"] = db.Meals("dinner")
	r["snack"] = db.Meals("snack")
	return r
}

// business logic function
func ChooseMeal(food map[string][]data.Meal, k string) (string, error) {
	r, ok := food[k]
	if !ok {
		return "", errors.New("meal kind not found")
	}
	i := rand.Intn(len(r))
	return r[i].Name, nil
}
