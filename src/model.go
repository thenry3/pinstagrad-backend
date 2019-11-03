package main

import (
	"time"
)

// Season types
type Season int

const (
	Spring Season = 0
	Summer Season = 1
	Fall   Season = 2
	Winter Season = 3
)

// TimeDay types
type TimeDay int

const (
	Night TimeDay = 1
	Day   TimeDay = 0
)

type User struct {
	Name     string   `json: "name"`
	Liked    []*Photo `json: "liked"`
	Pictures int      `json: "pictures"`
}

type Photo struct {
	Id           int       `json: "id"`
	Tags         []string  `json: "tags"`
	Uploadtime   time.Time `json: "uploadtime"`
	Photographer string    `json: "photographer"`
}

type Tags struct {
	Location string  `json: "location"`
	Season   Season  `json: "season"`
	Time     TimeDay `json: "time"`
}
