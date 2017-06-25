package data

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Run struct {
	Id    int
	Miles float64
	Date  time.Time
}

func connect() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/runs.db")
	if err != nil {
		panic(err)
	}
	return db
}

func begin(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	return tx
}

func AddRun(miles float64, date time.Time) {
	db := connect()
	defer db.Close()
	tx := begin(db)

	insert, err := tx.Prepare("INSERT INTO runs(miles, date) VALUES(?, ?)")
	if err != nil {
		panic(err)
	}
	defer insert.Close()
	_, err = insert.Exec(miles, date)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}

func LastRun() Run {
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, miles, date FROM runs ORDER BY datetime(date) DESC LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var miles float64
		var date time.Time
		err = rows.Scan(&id, &miles, &date)
		if err != nil {
			panic(err)
		}
		return Run{id, miles, date}
	}
	return Run{1, 100.0, time.Now().UTC()}
}

func LastMonth() (runs []Run) {
	db := connect()
	defer db.Close()

	now := time.Now().UTC()
	first := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	rows, err := db.Query("SELECT id, miles, date FROM runs WHERE date >= ?", first)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var miles float64
		var date time.Time
		err = rows.Scan(&id, &miles, &date)
		if err != nil {
			panic(err)
		}
		runs = append(runs, Run{id, miles, date})
	}
	return runs
}
