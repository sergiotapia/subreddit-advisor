package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	api_endpoint = "http://www.reddit.com/r/%s/top.json?t=month&limit=100"
	longForm     = "Jan 2, 2006 at 3:04pm (MST)"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Please enter 1 argument: Subreddits (eg: askreddit, askscience, technology)")
		os.Exit(1)
	}

	subreddit := os.Args[1]
	top_by_subreddit_endpoint := fmt.Sprintf(api_endpoint, subreddit)

	var top_submissions struct {
		Kind interface{}
		Data struct {
			ModHash  interface{}
			Children []struct {
				Kind interface{}
				Data Submission
			}
			After  string
			Before interface{}
		}
	}

	// Hit the API service and marshal into top_submissions.
	response, err := http.Get(top_by_subreddit_endpoint)
	if err != nil {
		panic(err)
	} else {
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		} else {
			if err := json.Unmarshal(content, &top_submissions); err != nil {
				panic(err)
			}
		}
	}

	// Make a slice of only Submissions.
	submissions := make([]Submission, len(top_submissions.Data.Children))
	for i, v := range top_submissions.Data.Children {
		submissions[i] = v.Data
	}

	day_of_week := best_day_of_week(submissions)
	time_of_day := best_time_of_day(submissions)
	fmt.Printf("/r/%s: It's best to submit your story on a %s from %s.", subreddit, day_of_week, time_of_day)
}

func best_day_of_week(submissions []Submission) string {
	// Create a map to save Weekday and it's associated submission count.
	submission_count := map[time.Weekday]int{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
		6: 0,
	}

	// Tally the top submission count per day.
	for _, submission := range submissions {
		time_of_submission := time.Unix(int64(submission.CreatedUtc), 0)
		submission_count[time_of_submission.Weekday()] += 1
	}

	// Find which Weekday has the highest amount of top submissions.
	best_day := time.Now().Weekday()
	highest_value := 0
	for key, weekday := range submission_count {
		if weekday > highest_value {
			highest_value = weekday
			best_day = key
		}
	}

	// Return the Weekday as a string.
	switch best_day {
	case 0:
		return "Sunday"
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	case 6:
		return "Saturday"
	default:
		return "Monday"
	}
}

func best_time_of_day(submissions []Submission) string {
	// Create a map to save Hour intervals and it's associated submission count.
	submission_count := map[string]int{
		"0-3":   0,
		"4-7":   0,
		"8-11":  0,
		"12-15": 0,
		"16-19": 0,
		"20-23": 0,
	}

	// Tally the top submission count per day.
	for _, submission := range submissions {
		time_of_submission := time.Unix(int64(submission.CreatedUtc), 0)
		hour := time_of_submission.Hour()

		if hour >= 0 && hour <= 3 {
			submission_count["0-3"] += 1
		} else if hour >= 4 && hour <= 7 {
			submission_count["4-7"] += 1
		} else if hour >= 8 && hour <= 11 {
			submission_count["8-11"] += 1
		} else if hour >= 12 && hour <= 15 {
			submission_count["12-15"] += 1
		} else if hour >= 16 && hour <= 19 {
			submission_count["16-19"] += 1
		} else {
			submission_count["20-23"] += 1
		}
	}

	// Find which Hour has the highest amount of top submissions.
	best_hours_range := ""
	highest_hour_range := 0
	for key, hour_range := range submission_count {
		if hour_range > highest_hour_range {
			highest_hour_range = hour_range
			best_hours_range = key
		}
	}
	return best_hours_range
}

type Submission struct {
	Id         string  `json:"id"`
	Subreddit  string  `json:"subreddit"`
	Score      int     `json:"score"`
	Ups        int     `json:"ups"`
	Downs      int     `json:"downs"`
	CreatedUtc float64 `json:"created_utc"`
}
