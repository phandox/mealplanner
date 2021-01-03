package data

import (
	"encoding/csv"
	"io"
	"log"
)

type MealsDB struct {
	storage string
	data    []Meal
}

func NewMealsDB(r io.Reader) MealsDB {
	d := csv.NewReader(r)
	_, err := d.Read() // skip header
	if err != nil {
		log.Fatal(err)
	}
	recs, err := d.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var m []Meal
	for _, r := range recs {
		m = append(m, Meal{
			Name: r[0],
			Kind: r[1],
		})
	}
	return MealsDB{data: m}
}

func (db MealsDB) GetMeal(kind string) *Meal {
	for _, m := range db.data {
		if m.Kind == kind {
			return &m
		}
	}
	return nil
}

type Meal struct {
	Name string
	Kind string
}
