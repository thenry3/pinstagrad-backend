package database

import (
	"time"
)

// Season types
type Season int

const (
	// Spring corresponding seasonal value
	Spring Season = 0
	// Summer corresponding seasonal value
	Summer Season = 1
	// Fall corresponding seasonal value
	Fall Season = 2
	// Winter corresponding seasonal value
	Winter Season = 3
)

// TimeDay types
type TimeDay int

const (
	// Night corresponding time value
	Night TimeDay = 1
	// Day corresponding time value
	Day TimeDay = 0
)

// User structure of user
type User struct {
	Name        string   `json:"name"`
	Liked       []string `json:"liked"`
	Pictures    int      `json:"pictures"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	PhoneNumber string   `json:"phonenumber"`
	ProfilePic  string   `json:"profilepic"`
	UID         string   `json:"uid"`
}

// Photo structure of photo
type Photo struct {
	UserID       int       `json:"id"`
	Pointer      string    `json:"pointer"`
	Location     string    `json:"location"`
	Time         TimeDay   `json:"time"`
	Season       Season    `json:"season"`
	Uploadtime   time.Time `json:"uploadtime"`
	Photographer string    `json:"photographer"`
	UUID         string    `json:"uuid"`
}

// Tags structure of tags
type Tags struct {
	Location string  `json:"location"`
	Season   Season  `json:"season"`
	Time     TimeDay `json:"time"`
}
