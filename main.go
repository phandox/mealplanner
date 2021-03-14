package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/phandox/mealplanner/internal/handlers"
	"github.com/phandox/mealplanner/internal/picker"

	"github.com/phandox/mealplanner/internal/data"
)

const mealsPath = "internal/data/meals.csv"
const defaultDb = "internal/data/meals.sqlite"

func main() {
	rand.Seed(time.Now().UnixNano())
	fd, err := os.Open(mealsPath)
	if err != nil {
		log.Fatal("can't open data source")
	}
	defer fd.Close()
	dbpool, err := sql.Open("sqlite3", defaultDb)
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
	db := data.NewManager(dbpool)
	if err = db.LoadMeals(csv.NewReader(fd)); err != nil {
		log.Fatal(err)
	}

	p := picker.NewPicker(db, picker.DefaultPeople)
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
