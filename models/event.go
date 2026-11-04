package models

import "time"

type Event struct {
	Id          int
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Location    string `binding:"required"`
	DateTime    time.Time
	UserId      int
}

var events []Event = []Event{}

// Fungsi untuk menyimpan event
func (e Event) Save() {
	events = append(events, e)
}

// Fungsi Menampilkan Semua Event
func GetAllEvents() []Event {
	return events
}
