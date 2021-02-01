package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phandox/mealplanner/internal/picker"

	"github.com/phandox/mealplanner/internal/data"
)

const TMPLDIR = "../templates/"

const testDataSevenDays = `"name","kind","portions"
"l1","lunch","2"
"l2","breakfast","4"
"l3","lunch","6"
"l4","dinner","2"
"l5","snack","6"
"l6","snack","2"
"l7","lunch","4"
`

func TestMainPageTemplateRender(t *testing.T) {
	tests := []struct {
		name    string
		tpath   string
		expcode int
	}{
		{
			"bad path",
			filepath.Join(TMPLDIR, "bad_path.gohtml"),
			http.StatusInternalServerError,
		},
		{
			"empty table header",
			filepath.Join(TMPLDIR, "mainpage.gohtml"),
			http.StatusOK,
		},
		{
			"not full template",
			filepath.Join(TMPLDIR, "mainpage.gohtml"),
			http.StatusOK,
		},
		{
			"filled template",
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
			db := data.NewMealsDB(strings.NewReader(testDataSevenDays))
			p := picker.NewPicker(&db, 2)
			handler := http.HandlerFunc(MainPage(test.tpath, p))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; rr.Code != test.expcode {
				t.Errorf("handler returned wrong status code: got %v want %v. Error: %q", status, test.expcode, rr.Body)
			}
		})
	}
}
