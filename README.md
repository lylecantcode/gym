# Gym (Tracker)

#### Video Demo: https://youtu.be/0Em6o0eiW0I 
#### Description:

**Part 1**
An exercise tracker to be used in CLI, I created this using GO (golang).
This program saves files using SQL, creating a new database "weights", if once doesn't already exist.
On run through, the program first checks the users previous entries into the database, to see if there is any data from the last training 'period', for a cyclical program. If found, then it will show the data and ask if the user want to reuse it.

**Part 2**
Regardless of the choice to reuse the data, the user is then asked how many new exercises should be added. 
If the day was exactly the same as last week, if the user decided to reuse the data, they can type in 0 here, and skip towards the part 4. Otherwise the user is asked which exercise they have done.

**Part 3**
The exercises are printed out to show the user their options, the users input is then compared to the exercise enumerated list to check for (partial) matches (therefore names can be abbreviated for convenience). Once the exercise is chosen, the user puts in their weights, units and reps in a fixed format of "WEIGHT UNIT x REPS". They are then asked how many times this has been repeated.

**Part 4**
This data, if no errors have occured (and the program exited), will have now been sent to the SQL database.
At this point, the user is asked if they'd like to see previous workouts, they have 3 options:
They can see all their previous workouts, the best weight (and rep, but with weight taking precedent), for each submitted exercise (with different units counting as different exercises, as these will often have different "bests").

**Latest Update**
Complete up to and including v8.

## Project Plan:

~v1) Take unputs of weights, units and date into SQL, display it after an input~  
~v2) Number of sets and multiple exercises per day. //Take input of number of exercises, then weights, sets, etc~  
~v3) Types of exercise done, e.g. Bench Press~  
~v4) Reorganised to use structs~  
~v5) Predictions based on day of week (for 7 day rotating schedule)~  
~v6) Add reps (replace sets with separate values, so can have different weights for each set)~
~v7) Add enums in for exercise and unit options~  
~v8) Personal bests per exercise, as a comparison to new workout~  

## Potential extras:

> Overwrite date option  
> Delete entries  
