package picker

import (
	"github.com/phandox/mealplanner/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const testDataSevenDays = `"name","kind","portions"
"l1","lunch","2"
"l2","lunch","4"
"l3","lunch","6"
"l4","lunch","2"
"l5","lunch","6"
"l6","lunch","2"
"l7","lunch","4"
`

func TestPlanLunchesMultiDays(t *testing.T) {
	tests := []struct {
		name    string
		DB      data.MealsDB
		days    int
		wantErr bool
	}{
		{
			"7 day lunch plan",
			data.NewMealsDB(strings.NewReader(testDataSevenDays)),
			7,
			false,
		},
		{
			"5 day lunch plan",
			data.NewMealsDB(strings.NewReader(testDataSevenDays)),
			5,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewPicker(&test.DB, DefaultPeople)
			got, err := p.PlanLunches(test.days)
			if test.wantErr {
				require.Error(t, err)
			}
			require.Len(t, got, test.days)
		})
	}
}

func TestPlanRandomLunches(t *testing.T) {
	tests := []struct {
		name string
		db   data.MealsDB
	}{
		{
			"different values repeated call",
			data.NewMealsDB(strings.NewReader(testDataSevenDays)),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewPicker(&test.db, DefaultPeople)
			first, err := p.PlanLunches(5)
			require.NoError(t, err)
			second, err := p.PlanLunches(5)
			require.NoError(t, err)
			require.NotEqual(t, first, second)
		})
	}
}

func TestPicker_PlanLunchesLogic(t *testing.T) {
	tests := []struct {
		name   string
		people int
		db     data.MealsDB
	}{
		{
			"Minimum of 2 days per planned meal",
			DefaultPeople,
			data.NewMealsDB(strings.NewReader(testDataSevenDays)),
		},
		{
			"planned meal have to be aligned in week",
			DefaultPeople,
			data.NewMealsDB(strings.NewReader(testDataSevenDays)),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewPicker(&test.db, test.people)
			got, err := p.PlanLunches(7)
			require.NoError(t, err)
			assert.Truef(t, checkPortions(got, 1, 7), "Portions check: got %v", got)
			assert.Truef(t, checkMealBoundary(got, test.people), "Boundary check: got %v", got)
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
