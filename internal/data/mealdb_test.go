package data

import (
	"context"
	"database/sql"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testData = `"name","kind","portions"
"lunch 1","lunch","4"
"breakfast 1","breakfast","4"
"dinner 1","dinner","4"
"snack 1","snack","4"
`
const testDuplicates = `"name","kind","portions"
"lunch 1","lunch","4"
"lunch 1","lunch","4"
`
const testDataMulti = `"name","kind","portions"
"lunch 1","lunch",4
"lunch 2","lunch",4
"dinner 1","dinner",4
"snack 1","snack",4
`

const testDb = "testMeals.sqlite"

var db *sql.DB

func setup(dbname string) (*sql.DB, string) {
	dir, err := ioutil.TempDir("", "mealplanner-integration")
	if err != nil {
		log.Fatal(err)
	}
	dbpath := filepath.Join(dir, dbname)
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		defer os.RemoveAll(dir)
		log.Fatal(err)
	}
	return db, dir
}

func teardown(db *sql.DB, tmpDir string) {
	defer os.RemoveAll(tmpDir)
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func cleanDb(t *testing.T, db *sql.DB) {
	t.Helper()
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		t.Fatal(err)
	}
	//goland:noinspection SqlWithoutWhere
	_, execErr := tx.Exec("DELETE FROM meal;")
	if execErr != nil {
		_ = tx.Rollback()
		t.Fatal(execErr)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func countRows(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	rows, err := db.Query("SELECT * FROM meal;", table)
	if err != nil {
		t.Fatal(err)
	}
	var count int
	for rows.Next() {
		count++
	}
	return count
}

func TestIntegration_LoadMealsValidCSV(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		expRows int
	}{
		{
			"only unique names",
			testData,
			4,
		},
		{
			"skip duplicate names",
			testDuplicates,
			1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := NewManager(db)
			err := m.LoadMeals(csv.NewReader(strings.NewReader(test.data)))
			defer cleanDb(t, db)
			assert.NoError(t, err)
			rows := countRows(t, m.GetDb(), "meal")
			assert.Equal(t, test.expRows, rows)
		})
	}
}

func TestMain(m *testing.M) {
	var dir string
	db, dir = setup(testDb)
	retCode := m.Run()
	teardown(db, dir)
	os.Exit(retCode)
}
