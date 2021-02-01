package data

import (
	"encoding/csv"
	"log"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

const testData = `"name","kind","portions"
"lunch 1","lunch","4"
"breakfast 1","breakfast","4"
"dinner 1","dinner","4"
"snack 1","snack","4"
`
const testDataMulti = `"name","kind","portions"
"lunch 1","lunch","4"
"lunch 2","lunch","4"
"dinner 1","dinner","4"
"snack 1","snack","4"
`

func loadMeals(t *testing.T, s string) []Meal {
	t.Helper()
	d := csv.NewReader(strings.NewReader(s))
	_, err := d.Read() // skip header
	if err != nil {
		log.Fatal(err)
	}
	rec, err := d.ReadAll()
	if err != nil {
		t.Fatal("can't read test data")
	}
	var m []Meal
	for _, r := range rec {
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
	return m
}

func TestNewMealsDB(t *testing.T) {
	tests := []struct {
		name string
		err  error
		data string
	}{
		{
			"load CSV data",
			nil,
			testData,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expDB := MealsDB{
				data: loadMeals(t, test.data),
			}
			db := NewMealsDB(strings.NewReader(test.data))
			switch test.err {
			default:
				if !reflect.DeepEqual(db.data, expDB.data) {
					t.Errorf("got %v != want %v", db.data, expDB.data)
				}
			}
		})
	}
}

func TestGetMeal(t *testing.T) {
	tests := []struct {
		name string
		kind string
		want Meal
	}{
		{
			"get lunch meal",
			"lunch",
			Meal{
				Kind: "lunch",
			},
		},
		{
			"get breakfast meal",
			"breakfast",
			Meal{
				Kind: "breakfast",
			},
		},
		{
			"get snack meal",
			"snack",
			Meal{
				Kind: "snack",
			},
		},
		{
			"get dinner meal",
			"dinner",
			Meal{
				Kind: "dinner",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := MealsDB{
				data: loadMeals(t, testData),
			}
			got := db.GetMeal(test.kind)

			if got == nil {
				t.Fatal("unexpected failure: no meal returned")
			}

			if got.Kind != test.want.Kind {
				t.Errorf("got %v != want %v", got.Kind, test.want.Kind)
			}
		})
	}
}

func TestMeals(t *testing.T) {
	tests := []struct {
		name string
		kind string
		want int
		data string
	}{
		{
			"single meal in DB",
			"breakfast",
			1,
			testData,
		},
		{
			"two meals of type in DB",
			"lunch",
			2,
			testDataMulti,
		},
		{
			"mismatch cases of kind",
			"LUnch",
			2,
			testDataMulti,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := MealsDB{
				storage: "",
				data:    loadMeals(t, test.data),
			}
			got := db.Meals(test.kind)
			if len(got) != test.want {
				t.Errorf("got %v != want %v elements of %s kind", len(got), test.want, test.kind)
			}
		})
	}
}
