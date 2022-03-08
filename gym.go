package main

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	_"github.com/mattn/go-sqlite3"
)

func main() {
	var weight int
	var unit, date string
	date = time.Now().Format("02-01-2006")
	fmt.Println("What weight did you lift? e.g. 5 kg, 10 lbs")
	fmt.Scanf("%d %s", &weight, &unit)
	if !((weight > 0 && weight < 1000) && (unit == "kg" || unit == "lbs")){
		fmt.Println("Weight must be an integer, units must be \"kg\" or \"lbs\".")
		os.Exit(1) 
	}
	database, _ :=
		sql.Open("sqlite3", "./gym.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS weights (id INTEGER PRIMARY KEY, weight INT NOT NULL, units TEXT NOT NULL, date TEXT NOT NULL)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO weights (weight, units, date) VALUES (?,?,?)")
	statement.Exec(weight, unit, date)
	rows, _ :=
		database.Query("SELECT * FROM weights")
	
	var id int 
	//var units, date string
	for rows.Next() {
		rows.Scan(&id, &weight, &unit, &date)
		fmt.Println(id, ":", weight, unit, "on", date)
	}

}

