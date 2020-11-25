package main

import (
	"fmt"
	"github.com/phandox/mealplanner/internal/handlers"
	"net/http"
)

func main() {
	db := handlers.MealDB{}
	db.Add("delicious meal 1")
	db.Add("delicious meal 2")
	m := http.NewServeMux()
	m.HandleFunc("/meals", handlers.GetMeals(&db))
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
