package picker

import (
	"github.com/phandox/mealplanner/internal/data"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const testDataSevenDays = `"name","kind"
"l1","lunch"
"l2","lunch"
"l3","lunch"
"l4","lunch"
"l5","lunch"
"l6","lunch"
"l7","lunch"
`

func TestPlanLunches(t *testing.T) {
	tests := []struct {
		name    string
		p       *Picker
		DB      data.MealsDB
		days    int
		wantErr bool
	}{
		{
			"7 day lunch plan",
			NewPicker(),
			data.NewMealsDB(strings.NewReader(testDataSevenDays)),
			7,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.p.PlanLunches(&test.DB, test.days)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Len(t, got, test.days)
		})
	}
}
