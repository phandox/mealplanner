package main

import (
	"fmt"
	"github.com/phandox/mealplanner/internal/data"
	"github.com/phandox/mealplanner/internal/handlers"
	"log"
	"net/http"
	"os"
)

const mealsPath = "internal/data/meals.csv"

func main() {
	fd, err := os.Open(mealsPath)
	if err != nil {
		log.Fatal("can't open data source")
	}
	defer fd.Close()
	db := data.NewMealsDB(fd)
	tableHeader := handlers.MainPageTable{
		Days:      []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		MealTypes: []string{"Breakfast", "Snack", "Lunch", "Snack", "Dinner"},
		Food:      populateMealPlan(db),
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

func populateMealPlan(db data.MealsDB) []*data.Meal {
	var meals []*data.Meal
	meals = append(meals, db.GetMeal("breakfast"))
	meals = append(meals, db.GetMeal("lunch"))
	meals = append(meals, db.GetMeal("dinner"))
	meals = append(meals, db.GetMeal("snack"))
	return meals
}
