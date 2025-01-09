package dbrepo

import (
	"errors"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"time"
)

// InsertReservation insert a reservation into database
func (psql *testdbPostgresRepository) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomId == 2 {
		return 0, errors.New("can't insert reservation for room id 2")
	}
	return 1, nil
}

// InsertRoomRestriction insert a room restriction into database
func (psql *testdbPostgresRepository) InsertRoomRestriction(res models.RoomRestriction) error {
	// if the room id is 1000, then fail; otherwise, pass
	if res.RoomId == 1000 {
		return errors.New("can't insert room restriction for room id 1000")
	}
	return nil
}

// SearchAvailabilityForDatesByRoomId query database with dates if available for booking room
func (psql *testdbPostgresRepository) SearchAvailabilityForDatesByRoomId(_ int, checkIn, _ time.Time) (bool, error) {
	// set up a test time
	layout := "2006-01-02"
	testDate, _ := time.Parse(layout, "2049-12-31")
	testDateToFail, _ := time.Parse(layout, "2060-01-01")

	if checkIn == testDateToFail {
		return false, errors.New("invalid date")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	if checkIn.After(testDate) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any for given date range
func (psql *testdbPostgresRepository) SearchAvailabilityForAllRooms(checkIn, _ time.Time) ([]models.Room, error) {
	var rooms []models.Room
	// set up a test time
	layout := "2006-01-02"
	testDate, _ := time.Parse(layout, "2049-12-31")
	testDateToFail, _ := time.Parse(layout, "2060-01-01")

	if checkIn == testDateToFail {
		return rooms, errors.New("invalid date")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	if checkIn.After(testDate) {
		return rooms, nil
	}
	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		Id: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

// GetRoomById get a room by id
func (psql *testdbPostgresRepository) GetRoomById(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("can't find room with id greater than 2")
	}
	return room, nil
}
