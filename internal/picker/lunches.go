package picker

import (
	"math/rand"

	"github.com/phandox/mealplanner/internal/data"
)

const DefaultPeople = 2

type Picker struct {
	db     *data.MealsDB
	people int
}

func NewPicker(db *data.MealsDB, people int) *Picker {
	return &Picker{db: db, people: people}
}

func (p Picker) Plan(kind string, days int) ([]data.Meal, error) {
	f := p.db.Meals(kind)
	return p.pick(f, days)
}

func pickFood(food []data.Meal, minPortions int, maxPortions int) data.Meal {
	var daysOk, portionsOk bool
	var m data.Meal
	for daysOk != true || portionsOk != true {
		daysOk = false
		portionsOk = false
		m = food[rand.Intn(len(food))]
		if m.Portions <= maxPortions {
			portionsOk = true
		}
		if m.Portions >= minPortions {
			daysOk = true
		}
	}
	return m
}

func (p Picker) PlanLunches(days int) ([]data.Meal, error) {
	food := p.db.Meals("lunch")
	week := make([]data.Meal, days)
	portions := days * p.people
	i := 0
	for fill := 0; fill < portions; {
		m := pickFood(food, 2*p.people, portions-fill)
		fill = fill + m.Portions
		for j := 0; j < m.Portions/p.people && i < cap(week); j++ {
			week[i] = m
			i++
		}
	}
	return week, nil
}

func (p Picker) pick(m []data.Meal, days int) ([]data.Meal, error) {
	var plan []data.Meal
	for i := 0; i < days; i++ {
		plan = append(plan, m[rand.Intn(len(m))])
	}
	return plan, nil

}
