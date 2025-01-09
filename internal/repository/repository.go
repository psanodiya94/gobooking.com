package repository

import (
	"github.com/psanodiya94/gobooking.com/internal/models"
	"time"
)

type DBRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityForDatesByRoomId(roomId int, checkIn, checkOut time.Time) (bool, error)
	SearchAvailabilityForAllRooms(checkIn, checkOut time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
}
