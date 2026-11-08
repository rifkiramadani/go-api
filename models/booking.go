package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	BookingCode string `json:"booking_code"`
	Phone       string `json:"phone"`

	UserID int  `json:"userId"`
	User   User `gorm:"foreignKey:UserId" json:"user"`

	EventID int   `json:"eventId"`
	Event   Event `gorm:"foreignKey:EventId" json:"event"`
}
