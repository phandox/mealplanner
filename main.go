package main

import (
	"fmt"
	"github.com/phandox/mealplanner/internal/handlers"
	"net/http"
)

func main() {
	tableHeader := handlers.MainPageTable{
		Days:      []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		MealTypes: []string{"Breakfast", "Snack", "Lunch", "Snack", "Dinner"},
	}
	db := handlers.MealDB{}
	db.Add("delicious meal 1")
	db.Add("delicious meal 2")
	m := http.NewServeMux()
	m.HandleFunc("/meals", handlers.GetMeals(&db))
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
