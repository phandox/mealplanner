package handlers

import (
	"github.com/phandox/mealplanner/internal/data"
	"html/template"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

const TMPLDIR = "../templates/"

func TestMainPageTemplateRender(t *testing.T) {
	tests := []struct {
		name    string
		th      MainPageTable
		tpath   string
		expcode int
	}{
		{
			"bad path",
			MainPageTable{
				Days:      nil,
				MealTypes: nil,
			},
			filepath.Join(TMPLDIR, "bad_path.gohtml"),
			http.StatusInternalServerError,
		},
		{
			"empty table header",
			MainPageTable{
				Days:      nil,
				MealTypes: nil,
				Fm: template.FuncMap{
					"lower": strings.ToLower,
					"chooseMeal": func(food map[string][]*data.Meal, k string) (*data.Meal, error) {
						return nil, nil
					},
				},
			},
			filepath.Join(TMPLDIR, "mainpage.gohtml"),
			http.StatusOK,
		},
		{
			"not full template",
			MainPageTable{
				Days:      []string{"Monday", "Tuesday"},
				MealTypes: nil,
				Fm: template.FuncMap{
					"lower": strings.ToLower,
					"chooseMeal": func(food map[string][]*data.Meal, k string) (*data.Meal, error) {
						return nil, nil
					},
				},
			},
			filepath.Join(TMPLDIR, "mainpage.gohtml"),
			http.StatusOK,
		},
		{
			"filled template",
			MainPageTable{
				Days:      []string{"Monday", "Tuesday"},
				MealTypes: []string{"Breakfast", "Snack"},
				Fm: template.FuncMap{
					"lower": strings.ToLower,
					"chooseMeal": func(food map[string][]*data.Meal, k string) (*data.Meal, error) {
						return nil, nil
					},
				},
			},
			filepath.Join(TMPLDIR, "mainpage.gohtml"),
			http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(MainPage(test.th, test.tpath))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; rr.Code != test.expcode {
				t.Errorf("handler returned wrong status code: got %v want %v. Error: %q", status, test.expcode, rr.Body)
			}
		})
	}
}
