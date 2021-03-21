package picker

import (
	"context"
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/phandox/mealplanner/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockDataManager struct {
	stored []data.Meal
}

func (d *MockDataManager) GetDb() *sql.DB {
	panic("not supported on MockDataManager")
}

func (d *MockDataManager) Close() error {
	panic("implement me")
}

func (d *MockDataManager) AddMeal(ctx context.Context, name string, portions string, kind string) error {
	panic("implement me")
}

func (d *MockDataManager) GetMeals(ctx context.Context, kind string) ([]data.Meal, error) {
	var meals []data.Meal
	for _, v := range d.stored {
		if v.Kind == kind {
			var m data.Meal
			m = v
			meals = append(meals, m)
		}
	}
	return meals, nil
}

func (d *MockDataManager) LoadMeals(r *csv.Reader) error {
	_, err := r.Read() // read header and throw away
	if err != nil {
		log.Fatal(err)
	}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var m data.Meal
		p, err := strconv.Atoi(rec[2])
		if err != nil {
			log.Fatal(err)
		}
		m = data.Meal{
			Name:     rec[0],
			Kind:     rec[1],
			Portions: p,
		}
		d.stored = append(d.stored, m)
	}
	return nil
}

const testDataSevenDays = `"name","kind","portions"
"l1","lunch","4"
"l2","lunch","4"
"l3","lunch","6"
"l4","lunch","6"
"l5","lunch","6"
"l6","lunch","8"
"l7","lunch","4"
"l8","lunch","4"
"l9","lunch","4"
"l10","lunch","2"
`
const testDataNotEnoughPortions = `"name","kind","portions"
"l1","lunch","2"
"l2","lunch","2"
"l3","lunch","2"
"l4","lunch","2"
"l5","lunch","2"
"l6","lunch","2"
"l7","lunch","2"
`
const testDataBoundaryOverflow = `"name","kind","portions"
"l1","lunch","3"
"l2","lunch","3"
"l3","lunch","3"
"l4","lunch","3"
"l5","lunch","3"
"l6","lunch","3"
"l7","lunch","3"
`
const testDataNotEnoughUniqueFood = `"name","kind","portions"
"l1","lunch","4"
"l2","lunch","4"
`

func TestPlanLunchesMultiDays(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		days    int
		seed    int64
		wantErr bool
	}{
		{
			"7 day lunch plan",
			testDataSevenDays,
			7,
			2,
			false,
		},
		{
			"5 day lunch plan",
			testDataSevenDays,
			5,
			2,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var db MockDataManager
			err := db.LoadMeals(csv.NewReader(strings.NewReader(test.data)))
			if err != nil {
				t.Fatal(err)
			}
			p := NewPicker(&db, DefaultPeople, 2)
			got, err := p.PlanLunches(test.days)
			if test.wantErr {
				require.Error(t, err)
			}
			require.NoError(t, err)
			require.Len(t, got, test.days)
		})
	}
}

func TestPlanRandomLunches(t *testing.T) {
	tests := []struct {
		name string
		data string
		seed int64
	}{
		{
			"different values repeated call",
			testDataSevenDays,
			2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var db MockDataManager
			err := db.LoadMeals(csv.NewReader(strings.NewReader(test.data)))
			if err != nil {
				t.Fatal(err)
			}
			p := NewPicker(&db, DefaultPeople, test.seed)
			first, err := p.PlanLunches(7)
			require.NoError(t, err)
			second, err := p.PlanLunches(7)
			require.NoError(t, err)
			require.NotEqual(t, first, second)
		})
	}
}

func TestPicker_PlanLunchesLogic(t *testing.T) {
	tests := []struct {
		name   string
		people int
		data   string
		seed   int64
	}{
		{
			"Minimum of 2 days per planned meal",
			DefaultPeople,
			testDataSevenDays,
			2,
		},
		{
			"planned meal have to be aligned in week",
			DefaultPeople,
			testDataSevenDays,
			2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var db MockDataManager
			err := db.LoadMeals(csv.NewReader(strings.NewReader(test.data)))
			if err != nil {
				t.Fatal(err)
			}
			p := NewPicker(&db, test.people, test.seed)
			got, err := p.PlanLunches(7)
			require.NoError(t, err)
			assert.Truef(t, checkPortions(got, 2, 7), "Portions check: got %v", got)
			assert.Truef(t, checkMealBoundary(got, test.people), "Boundary check: got %v", got)
		})
	}
}

func TestPicker_PlanLunchesFailures(t *testing.T) {
	tests := []struct {
		name   string
		people int
		data   string
		seed   int64
	}{
		{
			"No food with minimum portions available",
			DefaultPeople,
			testDataNotEnoughPortions,
			2,
		},
		{
			"can not satisfy boundary",
			DefaultPeople,
			testDataBoundaryOverflow,
			2,
		},
		{
			"not enough unique food",
			DefaultPeople,
			testDataNotEnoughUniqueFood,
			2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var db MockDataManager
			err := db.LoadMeals(csv.NewReader(strings.NewReader(test.data)))
			if err != nil {
				t.Fatal(err)
			}
			p := NewPicker(&db, test.people, test.seed)
			got, err := p.PlanLunches(7)
			assert.Error(t, err)
			assert.Nil(t, got)
		})
	}
}

// TODO when the food is picked, it can't be picked again

func checkMealBoundary(m []data.Meal, people int) bool {
	occurrences := make(map[string]int, len(m))
	for _, v := range m {
		occurrences[v.Name]++
	}
	for _, v := range m {
		if v.Portions != occurrences[v.Name]*people {
			return false
		}
	}
	return true
}

func checkPortions(m []data.Meal, minDays int, days int) bool {
	occurrences := make(map[string]int, days)
	for _, v := range m {
		occurrences[v.Name]++
	}
	for _, v := range occurrences {
		if v < minDays {
			return false
		}
	}
	return true
}
