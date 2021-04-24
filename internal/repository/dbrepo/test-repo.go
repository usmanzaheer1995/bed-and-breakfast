package dbrepo

import (
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/models"
	"time"
)

func (t *testDBRepo) AllUsers() bool {
	return true
}

func (t *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

func (t *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

func (t *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {
	return false, nil
}

func (t *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var room []models.Room

	return room, nil
}

func (t *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	return room, nil
}
