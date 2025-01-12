package models

import (
	"time"
)

// User is the user model
type User struct {
	Id          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the room model
type Room struct {
	Id        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the restriction model
type Restriction struct {
	Id              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the reservation model
type Reservation struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	CheckIn   time.Time
	CheckOut  time.Time
	RoomId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
	Processed int
}

// RoomRestriction is room_restriction model
type RoomRestriction struct {
	Id            int
	CheckIn       time.Time
	CheckOut      time.Time
	RoomId        int
	ReservationId int
	RestrictionId int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Reservation   Reservation
	Restriction   Restriction
}

// MailData holds an email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
