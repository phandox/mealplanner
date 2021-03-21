package picker

import (
	"context"
	"errors"
	"math/rand"

	"github.com/phandox/mealplanner/internal/data"
)

const DefaultPeople = 2

type Picker struct {
	db     data.Manager
	people int
	seed   int64
}

func NewPicker(db data.Manager, people int, seed int64) *Picker {
	rand.Seed(seed)
	return &Picker{db: db, people: people, seed: seed}
}

func (p Picker) PlanRandom(kind string, days int) ([]data.Meal, error) {
	f, err := p.db.GetMeals(context.TODO(), kind)
	if err != nil {
		return nil, err
	}
	return p.pick(f, days)
}

func allPicked(p *[]bool) bool {
	for _, v := range *p {
		if v == false {
			return false
		}
	}
	return true
}

func pickFood(food []data.Meal, minPortions int, maxPortions int, picked *[]bool) (data.Meal, error) {
	var daysOk, portionsOk bool
	var m data.Meal
	for !allPicked(picked) {
		daysOk = false
		portionsOk = false
		idx := rand.Intn(len(food))
		if (*picked)[idx] {
			continue
		}
		m = food[idx]
		if m.Portions <= maxPortions {
			portionsOk = true
		}
		if m.Portions >= minPortions {
			daysOk = true
		}
		(*picked)[idx] = true
		if portionsOk && daysOk {
			return m, nil
		}
	}
	return data.Meal{}, errors.New("no food available")
}

func (p Picker) PlanLunches(days int) ([]data.Meal, error) {
	food, err := p.db.GetMeals(context.TODO(), "lunch")
	if err != nil {
		return nil, err
	}
	week := make([]data.Meal, days)
	portions := days * p.people
	picked := make([]bool, len(food))
	i := 0
	for fill := 0; fill < portions; {
		m, err := pickFood(food, 2*p.people, portions-fill, &picked)
		if err != nil {
			return nil, err
		}
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
