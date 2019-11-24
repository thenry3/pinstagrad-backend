package main

import (
	"time"
)

// Season types
type Season int

// Seasons
const (
	Spring Season = 0
	Summer Season = 1
	Fall   Season = 2
	Winter Season = 3
)

// TimeDay types
type TimeDay int

// Times of Day
const (
	Night TimeDay = 1
	Day   TimeDay = 0
)

// User structure
type User struct {
	Name     string   `json: "name"`
	Liked    []*Photo `json: "liked"`
	Pictures int      `json: "pictures"`
}

// Photo structure model
type Photo struct {
	ID           int       `json: "id"`
	Tags         []string  `json: "tags"`
	Uploadtime   time.Time `json: "uploadtime"`
	Photographer string    `json: "photographer"`
}

// Tags model
type Tags struct {
	Location string  `json: "location"`
	Season   Season  `json: "season"`
	Time     TimeDay `json: "time"`
}
