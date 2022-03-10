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
	reps int
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
		database.Prepare("CREATE TABLE IF NOT EXISTS weights (id INTEGER PRIMARY KEY, exercise INT NOT NULL, weight INT NOT NULL, units INT NOT NULL, reps INT NOT NULL DEFAULT 5, date TEXT NOT NULL)")
	statement.Exec()

	wl := weightLifting {Tricepext, 0, Resistancebands, 0}
	wl.prediction(database)

	fmt.Println("How many different exercises have you done?")
	exercisesDone := 1
	fmt.Scanf("%d", &exercisesDone)
	
	for i := 0; i < exercisesDone; i++{
		fmt.Println("What was exercise", i+1, "? Please choose from the following:") 
		for i := Exercises(0); i < ELimit; i++ {
			fmt.Println(">", i)
		} 
		wl.exerciseDone(database)
	}
	var records string
	fmt.Println("Would you like to see your previous workouts? y/n/best")
	fmt.Scanf("%s", &records)
	switch strings.ToLower(records) {
	case "y", "yes":
		wl.printAll(database)
	case "best", "bests", "b":
		wl.printBest(database)
	}
	
}



func (wl weightLifting) addToDatabase(database *sql.DB){
	statement, _ :=
			database.Prepare("INSERT INTO weights (exercise, weight, units, reps, date) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(wl.exercise, wl.weight, wl.unit, wl.reps, time.Now().Format("02-01-2006"))
	//statement.Exec(wl.exercise, wl.weight, wl.unit, wl.reps, time.Now().AddDate(0, 0, -7).Format("02-01-2006")) //to add historical entry for testing
} 



func copyPrevious(history string, copyAll bool, database *sql.DB) {
	if copyAll == true {
		statement, _ :=
			database.Prepare("INSERT INTO weights (date, exercise, weight, units, reps) SELECT ?, exercise, weight, units, reps FROM weights WHERE date = ?")
		statement.Exec(time.Now().Format("02-01-2006"), history)
	}
}



func (wl weightLifting) exerciseDone(database *sql.DB) {
	var check bool
	var exerciseInput, unitInput string
	fmt.Scanf("%s", &exerciseInput)
	check, wl.exercise = parseExercise(exerciseInput)
	if !check {
		fmt.Fprintf(os.Stderr, "Invalid exercise: %s.", exerciseInput)
		os.Exit(2)
	}

	fmt.Println("Please input weight for", wl.exercise, "in the format:  5 kg x 3 || 40 lbs x 1")
	fmt.Scanf("%d %s x %d", &wl.weight, &unitInput, &wl.reps)
	check, wl.unit = parseUnits(unitInput)
	if !check {
		fmt.Fprintf(os.Stderr, "Invalid units: %s", unitInput)
		os.Exit(3)
	}

	if !((wl.weight > 0 && wl.weight < 1000) && (wl.reps > 0 && wl.reps < 1000)){
		fmt.Fprintf(os.Stderr, "Weights and reps must be positive integers.")
		os.Exit(4)
	}
	var sets int
	fmt.Println("How many times did you do this exercise?")
	fmt.Scanf("%d", &sets)
	if sets < 1 {
		fmt.Fprintf(os.Stderr,"Must have done the exercise at least once.")
		os.Exit(5)
	}
	for i := 0; i < sets; i++ {
		wl.addToDatabase(database)
	}
}


func (e Exercises) String() string {
	switch e {
	case Bench:
		return "bench"
	case Curls:
		return "curls"
	case Dips:
		return "dips"
	case Farmerswalks:
		return "farmers-walks"
	case Goodmornings:
		return "good-mornings"
	case Overheadpress:
		return "overhead-press"
	case Pullups:
		return "pullups"
	case Rows:
		return "rows"
	case Squats:
		return "squats"
	case Tricepdips:
		return "tricep-dips"
	case Tricepext:
		return "tricep-extensions"
	}
	return "error"
}


func parseExercise(value string) (bool, Exercises) {
	value = strings.ToLower(value)
	for i := Exercises(0); i < ELimit; i++ {
		workout := Exercises(i).String()
		if strings.Contains(workout, value){
			return true, Exercises(i)
		}
	} 
	return false, ELimit
}


func parseUnits(value string) (bool, UnitOptions) {
	value = strings.ToLower(value)
	for i := UnitOptions(0); i < ULimit; i++ {
		units := UnitOptions(i).String()
		if strings.Contains(units, value){
			return true, UnitOptions(i)
		}
	} 
	return false, ULimit
}


func (wl weightLifting) prediction(database *sql.DB) {
	history := time.Now().AddDate(0, 0, -7).Format("02-01-2006")
	var id int
	rows, _ :=
			database.Query("SELECT id, exercise, weight, units, reps FROM weights WHERE date = ?", history)	
	for rows.Next() {
		rows.Scan(&id, &wl.exercise, &wl.weight, &wl.unit, &wl.reps)
		if id != 0 {
			fmt.Println( "On", history, "you did:", wl.exercise, wl.weight, wl.unit, "x", wl.reps)
		}
	}
	if id != 0 {
		var reuse string
		fmt.Println("Would you like to reuse this data? (y/n)")
		fmt.Scanln(&reuse)
		if strings.ToLower(reuse) == "y" {
			copyPrevious(history, true, database)
		}
	}
}


func (wl weightLifting) printAll(database *sql.DB) {
	var id int
	var date string
	rows, _ :=
		database.Query("SELECT * FROM weights")
	for rows.Next() {
		rows.Scan(&id, &wl.exercise, &wl.weight, &wl.unit, &wl.reps, &date)
		fmt.Println(id, ":", wl.exercise, wl.weight, wl.unit, "x", wl.reps, "on", date)
	}
}

func (wl weightLifting) printBest(database *sql.DB){
	var id int
	var date string
	for i := Exercises(0); i < ELimit; i++ {
		workout := Exercises(i)
		for j := UnitOptions(0); j < ULimit; j++ {
			units := UnitOptions(j)
			rows, _ :=
				database.Query("SELECT id, exercise, weight, units, reps, date FROM weights WHERE exercise = ? AND units = ? ORDER BY weight, reps DESC LIMIT 1", workout, units)// ORDER BY weight ASC LIMIT 1", workout, units)
			for rows.Next() {
				rows.Scan(&id, &wl.exercise, &wl.weight, &wl.unit, &wl.reps, &date)
				if id != 0 {
					fmt.Println("Your best", wl.exercise, "is", wl.weight, wl.unit, "x", wl.reps, "on", date)
				}
			}
		}
	}
}


func (u UnitOptions) String() string {
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


