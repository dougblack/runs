package data

import (
    "time"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func AddRun(miles float64, date time.Time) {
    db, err := sql.Open("sqlite3", "./data/runs.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        panic(err)
    }
    insert, err := tx.Prepare("insert into runs(miles, date) values(?, ?)")
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
