package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/phandox/mealplanner/internal/data"
	"github.com/phandox/mealplanner/internal/handlers"
	"github.com/phandox/mealplanner/internal/picker"
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
	m := http.NewServeMux()
	m.HandleFunc("/", handlers.MainPage("internal/templates/mainpage.gohtml", p))
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
