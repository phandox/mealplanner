package data

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

type MockManager struct {
	Calls  int
	stored map[string]meal
}

func (m *MockManager) Close() error {
	return nil
}

func (m *MockManager) AddMeal(ctx context.Context, name string, portions int, kind string) error {
	food := meal{name: name, portions: portions, kind: kind}
	if _, ok := m.stored[food.name]; ok {
		return errors.New("meal already stored")
	}
	m.stored[food.name] = food
	m.Calls++
	return nil
}

type meal struct {
	name     string
	kind     string
	portions int
}

func TestAddMealToEmptyDB(t *testing.T) {
	tests := []struct {
		name string
		m    meal
		err  error
	}{
		{
			"add single meal",
			meal{
				name:     "meal1",
				kind:     "lunch",
				portions: 2,
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := MockManager{
				Calls:  0,
				stored: map[string]meal{},
			}
			ctx := context.TODO()
			err := m.AddMeal(ctx, test.m.name, test.m.portions, test.m.kind)
			if test.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAddMealMultipleTimes(t *testing.T) {
	tests := []struct {
		name     string
		meals    []meal
		expCalls int
	}{
		{
			"add 2 different foods",
			[]meal{
				{
					"meal1",
					"lunch",
					2,
				},
				{
					"meal2",
					"lunch",
					2,
				},
			},
			2,
		},
		{
			"add 2 same food",
			[]meal{
				{
					"meal1",
					"lunch",
					2,
				},
				{
					"meal1",
					"lunch",
					2,
				},
			},
			1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := MockManager{
				Calls:  0,
				stored: map[string]meal{},
			}
			ctx := context.TODO()
			for _, f := range test.meals {
				_ = m.AddMeal(ctx, f.name, f.portions, f.kind)
			}
			assert.Equal(t, test.expCalls, m.Calls)
		})
	}
}
