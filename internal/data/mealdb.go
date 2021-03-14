package data

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"log"

	"github.com/mattn/go-sqlite3"
)

const initDbStmt = `
CREATE TABLE IF NOT EXISTS meal(
	id INTEGER PRIMARY KEY, 
	name TEXT NOT NULL UNIQUE,
	portions INTEGER,
	kind TEXT NOT NULL
);
`
const addMealStmt = `
INSERT INTO meal(name, portions, kind) VALUES(
 :name,
 :portions,
 :kind
)
`

const selectMealStmt = `SELECT name,kind,portions FROM meal WHERE kind = :kind`

type Manager interface {
	Close() error
	AddMeal(ctx context.Context, name string, portions string, kind string) error
	GetMeals(ctx context.Context, kind string) ([]Meal, error)
	LoadMeals(reader *csv.Reader) error
}

type Meal struct {
	name     string
	kind     string
	portions int
}

func (m Meal) Name() string {
	return m.name
}

func (m Meal) Kind() string {
	return m.kind
}

func (m Meal) Portions() int {
	return m.portions
}

type manager struct {
	db *sql.DB
}

func NewManager(dbpool *sql.DB) Manager {
	_, err := dbpool.Exec(initDbStmt)
	if err != nil {
		log.Fatal("Failed to prepare tables:", err)
	}
	return &manager{db: dbpool}
}

func (m *manager) Close() error {
	return m.db.Close()
}

func (m *manager) AddMeal(ctx context.Context, name string, portions string, kind string) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(
		ctx,
		addMealStmt,
		sql.Named("name", name),
		sql.Named("portions", portions),
		sql.Named("kind", kind),
	)
	// don't report unique constrain error, continue with operation
	// todo: add notice here to log?
	if err != nil && errors.Is(err, sqlite3.ErrConstraintUnique) {
		return nil
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (m *manager) GetMeals(ctx context.Context, kind string) ([]Meal, error) {
	var meals []Meal
	rows, err := m.db.QueryContext(ctx, selectMealStmt, sql.Named("kind", kind))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m Meal
		if err = rows.Scan(&m.name, &m.kind, &m.portions); err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return meals, nil
}

func (m *manager) LoadMeals(r *csv.Reader) error {
	_, err := r.Read()
	if err != nil {
		return err
	}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = m.AddMeal(context.TODO(), rec[0], rec[2], rec[1])
		if err != nil {
			return err
		}
	}
	return nil
}
