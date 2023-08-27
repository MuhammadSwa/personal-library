package models

import "time"

type Book struct {
	ISBN            string
	Title           string
	Author          string
	Category        string
	Publisher       string
	PuplicationDate string
	ImgPath         string
	PagesNumber     int
	PersonalRating  float32
	PersonalNotes   string
	ReadStatus      bool
	ReadDate        time.Time
}
