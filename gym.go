package main

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	_"github.com/mattn/go-sqlite3"
)

func main() {
	database, _ :=
		sql.Open("sqlite3", "./gym.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS weights (id INTEGER PRIMARY KEY, weight INT NOT NULL, units TEXT NOT NULL, sets INT NOT NULL DEFAULT 3, date TEXT NOT NULL)")
	statement.Exec()
		

	var weight, sets int
	var unit, date, formatting string
	date = time.Now().Format("02-01-2006")

	fmt.Println("How many different exercises have you done?")
	exercisesDone := 1
	fmt.Scanf("%d", &exercisesDone)
	if exercisesDone < 1{
		fmt.Println("Must have done at least 1 exercise.")
		os.Exit(1)
	}

	for i := 0; i < exercisesDone; i++{
		fmt.Println("Please input exercise", i+1, ", in the format:  5 kg x 3 || 40 lbs x 1",)
		fmt.Scanf("%d %s %s %d", &weight, &unit, &formatting, &sets)
		if !((weight > 0 && weight < 1000) && (sets > 0 && sets < 1000) && (unit == "kg" || unit == "lbs")){
			fmt.Println("Weight must be an integer, units must be \"kg\" or \"lbs\".")
			os.Exit(2)
		}
		addToDatabase(weight, sets, unit, date, database)
	}

	
	var id int
	rows, _ :=
		database.Query("SELECT * FROM weights")
	for rows.Next() {
		rows.Scan(&id, &weight, &unit, &sets, &date)
		fmt.Println(id, ":", weight, unit, "x", sets, "on", date)
	}
	for rows.Next() {
		rows.Scan(&id, &weight, &unit, &sets, &date)
		fmt.Println(id, ":", weight, unit, "x", sets, "on", date)
	}
}

func addToDatabase(weight, sets int, unit, date string, database *sql.DB){
	statement, _ :=
			database.Prepare("INSERT INTO weights (weight, units, sets, date) VALUES (?, ?, ?, ?)")
	statement.Exec(weight, unit, sets, date)
} 