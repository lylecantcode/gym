package main

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	_"github.com/mattn/go-sqlite3"
	"strings"
)

type weightLifting struct{
	exercise string
	weight int
	unit string
	sets int
}

func main() {
	database, _ :=
		sql.Open("sqlite3", "./gym.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS weights (id INTEGER PRIMARY KEY, exercise TEXT NOT NULL, weight INT NOT NULL, units TEXT NOT NULL, sets INT NOT NULL DEFAULT 3, date TEXT NOT NULL)")
	statement.Exec()

	exercises := []string{"bench", "curls", "dips", "farmers-walks", "good-mornings", "overhead-press", "pullups", "rows", "squats", "tricep-dips", "tricep-ext"}
	unitOptions := []string{"bodyweight", "kg", "kgs", "lb", "lbs"}
	var date, formatting string
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
		wl := weightLifting {"exercise", 0, "units", 0}
		fmt.Println("What was exercise", i+1, "? Please choose from the following:\n", exercises) 
		fmt.Scanf("%s", &wl.exercise)
		wl.exercise = strings.ToLower(wl.exercise)
		check, wl.exercise = validValue(wl.exercise, exercises)
		if !check {
			fmt.Println("Invalid exercise.")
			os.Exit(2)
		}

		fmt.Println("Please input weight for", wl.exercise, "in the format:  5 kg x 3 || 40 lbs x 1")
		fmt.Scanf("%d %s %s %d", &wl.weight, &wl.unit, &formatting, &wl.sets)
		wl.unit = strings.ToLower(wl.unit)
		check, wl.unit = validValue(wl.unit, unitOptions)
		if !check {
			fmt.Println("Invalid units, please use one of:", unitOptions)
			os.Exit(3)
		}

		if !((wl.weight > 0 && wl.weight < 1000) && (wl.sets > 0 && wl.sets < 1000)){
			fmt.Println("Weight must be a positive integer.")
			os.Exit(4)
		}
		wl.addToDatabase(date, database)
	}

	
	var id, weight, sets int
	var exercise, unit string
	rows, _ :=
		database.Query("SELECT * FROM weights")
	for rows.Next() {
		rows.Scan(&id, &exercise, &weight, &unit, &sets, &date)
		fmt.Println(id, ":", exercise, weight, unit, "x", sets, "on", date)
	}
}


func (wl weightLifting) addToDatabase(date string, database *sql.DB){
	statement, _ :=
			database.Prepare("INSERT INTO weights (exercise, weight, units, sets, date) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(wl.exercise, wl.weight, wl.unit, wl.sets, date)
} 


func validValue(value string, list []string) (bool, string) {
	for _, item := range list{
		if strings.Contains(item, value){
			return true, item
		}
	}
	return false, ""
}