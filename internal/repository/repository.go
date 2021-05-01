package repository

import (
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedReservation(id, processed int) error
	AllRooms() ([]models.Room, error)
	GetRestrictionsForRoomByDay(roomID int, startDate, endDate time.Time) ([]models.RoomRestriction, error)
	InsertBlockForRoom(id int, startDate time.Time) error
	DeleteBlockForRoom(id int) error
}
