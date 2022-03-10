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
	exercise Exercises
	weight int
	unit UnitOptions
	sets int
}

type Exercises int

const (
	Bench Exercises = iota
	Curls
	Dips
	Farmerswalks
	Goodmornings
	Overheadpress
	Pullups
	Rows
	Squats
	Tricepdips
	Tricepext
	ELimit
)

type UnitOptions int

const (
	Kg UnitOptions = iota
	Bodyweight
	Lbs
	Resistancebands
	ULimit
)



func main() {
	database, _ :=
		sql.Open("sqlite3", "./gym.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS weights (id INTEGER PRIMARY KEY, exercise TEXT NOT NULL, weight INT NOT NULL, units TEXT NOT NULL, sets INT NOT NULL DEFAULT 3, date TEXT NOT NULL)")
	statement.Exec()

	wl := weightLifting {Tricepext, 0, Resistancebands, 0}
	wl.prediction(database)

	date := time.Now().Format("02-01-2006")

	fmt.Println("How many different exercises have you done?")
	exercisesDone := 1
	fmt.Scanf("%d", &exercisesDone)
	
	for i := 0; i < exercisesDone; i++{
		fmt.Println("What was exercise", i+1, "? Please choose from the following:") 
		for i := Exercises(0); i < ELimit; i++ {
			fmt.Println(exerciseString(i))
		} 
		wl.exerciseDone(database)
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



func (wl weightLifting) addToDatabase(database *sql.DB){
	statement, _ :=
			database.Prepare("INSERT INTO weights (exercise, weight, units, sets, date) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(exerciseString(wl.exercise), wl.weight, unitString(wl.unit), wl.sets, time.Now().Format("02-01-2006"))
} 



func copyPrevious(history string, copyAll bool, database *sql.DB) {
	if copyAll == true {
		today := time.Now().Format("02-01-2006")
		statement, _ :=
			database.Prepare("INSERT INTO weights (date, exercise, weight, units, sets) SELECT CASE WHEN date = ? THEN date = ? ELSE date END, exercise, weight, units, sets FROM weights WHERE date = ?")
		statement.Exec(history, today, history)
		statement, _ =
			database.Prepare("UPDATE weights SET date = ? WHERE date = 0;")
		statement.Exec(today)
	}
}



func (wl weightLifting) exerciseDone(database *sql.DB) {
	var check bool
	var exerciseInput, unitInput string
	fmt.Scanf("%s", &exerciseInput)
	check, wl.exercise = parseExercise(strings.ToLower(exerciseInput))
	if !check {
		fmt.Println("Invalid exercise.")
		os.Exit(2)
	}

	fmt.Println("Please input weight for", exerciseString(wl.exercise), "in the format:  5 kg x 3 || 40 lbs x 1")
	fmt.Scanf("%d %s x %d", &wl.weight, &unitInput, &wl.sets)
	check, wl.unit = parseUnits(strings.ToLower(unitInput))
	if !check {
		fmt.Println("Invalid units, please use one of:")
		for i := UnitOptions(0); i < ULimit; i++ {
			fmt.Println(unitString(UnitOptions(i)))
		}
		os.Exit(3)
	}

	if !((wl.weight > 0 && wl.weight < 1000) && (wl.sets > 0 && wl.sets < 1000)){
		fmt.Println("Weights and Sets must be positive integers.")
		os.Exit(4)
	}
	wl.addToDatabase(database)
}


func exerciseString(e Exercises) string {
	switch e {
	case Bench:
		return "bench"
	case Curls:
		return "curls"
	case Dips:
		return "dips"
	case Farmerswalks:
		return "farmers walks"
	case Goodmornings:
		return "good mornings"
	case Overheadpress:
		return "overhead press"
	case Pullups:
		return "pullups"
	case Rows:
		return "rows"
	case Squats:
		return "squats"
	case Tricepdips:
		return "tricep dips"
	case Tricepext:
		return "tricep extensions"
	}
	return "error"
}


func parseExercise(value string) (bool, Exercises) {
	for i := Exercises(0); i < ELimit; i++ {
		workout := exerciseString(Exercises(i))
		if strings.Contains(workout, value){
			return true, Exercises(i)
		}
	} 
	return false, ELimit
}


func parseUnits(value string) (bool, UnitOptions) {
	for i := UnitOptions(0); i < ULimit; i++ {
		units := unitString(UnitOptions(i))
		if strings.Contains(units, value){
			return true, UnitOptions(i)
		}
	} 
	return false, ULimit
}


func unitString(u UnitOptions) string {
	switch u {
	case Kg:
		return "kgs"
	case Bodyweight:
		return "bodyweight"
	case Lbs:
		return "lbs"
	case Resistancebands:
		return "resistance bands"
	}
	return "error"
}


func (wl weightLifting) prediction(database *sql.DB) {
	history := time.Now().AddDate(0, 0, -7).Format("02-01-2006")
	var id int
	var date string

	rows, _ :=
			database.Query("SELECT * FROM weights WHERE date = ?", history)
	fmt.Println("Last week you did:")		
	for rows.Next() {
		rows.Scan(&id, &wl.exercise, &wl.weight, &wl.unit, &wl.sets, &date)
		fmt.Println(id, ":", wl.exercise, wl.weight, wl.unit, "x", wl.sets, "on", date)
	}
	var reuse string
	fmt.Println("Would you like to reuse this data? (y/n)")
	fmt.Scanln(&reuse)
	if strings.ToLower(reuse) == "y" {
		copyPrevious(history, true, database)
	}
}