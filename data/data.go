package data

import (
    "time"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type Run struct {
    Id int
    Miles float64
    Date time.Time
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

func LastRun() (Run) {
    db := connect()
    defer db.Close()

    rows, err := db.Query("select id, miles, date from runs order by datetime(date) desc limit 1")
    if err != nil {
        panic(err)
    }
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
