package main

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	_"github.com/mattn/go-sqlite3"
	"strings"
)

/*type weightLifting struct{
	exercise string
	weight int
	unit string
	sets int
}*/

func main() {
	database, _ :=
		sql.Open("sqlite3", "./gym.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS weights (id INTEGER PRIMARY KEY, exercise TEXT NOT NULL, weight INT NOT NULL, units TEXT NOT NULL, sets INT NOT NULL DEFAULT 3, date TEXT NOT NULL)")
	statement.Exec()

	exercises := []string{"bench", "curls", "dips", "farmers-walks", "good-mornings", "overhead-press", "pullups", "rows", "squats", "tricep-dips", "tricep-ext"}
	unitOptions := []string{"bodyweight", "kg", "kgs", "lb", "lbs"}
	var weight, sets int
	var unit, date, exercise, formatting string
	date = time.Now().Format("02-01-2006")

	fmt.Println("How many different exercises have you done?")
	exercisesDone := 1
	fmt.Scanf("%d", &exercisesDone)
	if exercisesDone < 1{
		fmt.Println("Must have done at least 1 exercise.")
		os.Exit(1)
	}
	var check bool
	for i := 0; i < exercisesDone; i++{
		fmt.Println("What was exercise", i+1, "? Please choose from the following:\n", exercises) 
		fmt.Scanf("%s", &exercise)
		exercise = strings.ToLower(exercise)
		check, exercise = validValue(exercise, exercises)
		if !check {
			fmt.Println("Invalid exercise.")
			os.Exit(2)
		}

		fmt.Println("Please input weight for", exercise, "in the format:  5 kg x 3 || 40 lbs x 1")
		fmt.Scanf("%d %s %s %d", &weight, &unit, &formatting, &sets)
		unit = strings.ToLower(unit)
		check, unit = validValue(unit, unitOptions)
		if !check {
			fmt.Println("Invalid units, please use one of:", unitOptions)
			os.Exit(3)
		}

		if !((weight > 0 && weight < 1000) && (sets > 0 && sets < 1000)){
			fmt.Println("Weight must be a positive integer.")
			os.Exit(4)
		}
		addToDatabase(weight, sets, exercise, unit, date, database)
	}

	
	var id int
	rows, _ :=
		database.Query("SELECT * FROM weights")
	for rows.Next() {
		rows.Scan(&id, &exercise, &weight, &unit, &sets, &date)
		fmt.Println(id, ":", exercise, weight, unit, "x", sets, "on", date)
	}
}

func addToDatabase(weight, sets int, exercise, unit, date string, database *sql.DB){
	statement, _ :=
			database.Prepare("INSERT INTO weights (exercise, weight, units, sets, date) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(exercise, weight, unit, sets, date)
} 


func validValue(value string, list []string) (bool, string) {
	for _, item := range list{
		if strings.Contains(item, value){
			return true, item
		}
	}
	return false, ""
}