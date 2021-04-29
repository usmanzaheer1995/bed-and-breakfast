package dbrepo

import (
	"errors"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/models"
	"time"
)

func (t *testDBRepo) AllUsers() bool {
	return true
}

func (t *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2 then fail
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (t *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

func (t *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {
	if roomId == 1000 {
		return false, errors.New("my error")
	}
	return false, nil
}

func (t *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var room []models.Room

	return room, nil
}

func (t *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}

func (t *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User

	return u, nil
}

func (t *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (t *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}