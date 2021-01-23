package picker

import (
	"github.com/phandox/mealplanner/internal/data"
	"math/rand"
)

type Picker struct {
}

func NewPicker() *Picker {
	return &Picker{}
}

func (p Picker) PlanLunches(db *data.MealsDB, days int) ([]data.Meal, error) {
	lunches := db.Meals("lunch")
	var plan []data.Meal
	for i := 0; i < days; i++ {
		plan = append(plan, lunches[rand.Intn(len(lunches))])
	}
	return plan, nil
}
