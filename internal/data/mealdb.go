package data

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
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
		n, err := strconv.Atoi(r[2])
		if err != nil {
			panic(err)
		}
		m = append(m, Meal{
			Name:     r[0],
			Kind:     r[1],
			Portions: n,
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

func (db MealsDB) Meals(kind string) []Meal {
	var r []Meal
	for _, m := range db.data {
		if strings.ToLower(kind) == strings.ToLower(m.Kind) {
			r = append(r, m)
		}
	}
	return r
}

type Meal struct {
	Name     string
	Kind     string
	Portions int
}

func (m Meal) IsKind(kind string) bool {
	return strings.ToLower(kind) == m.Kind
}
