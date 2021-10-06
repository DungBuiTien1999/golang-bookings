package dbrepo

import (
	"errors"
	"time"

	"github.com/DungBuiTien1999/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("some errors")
	}

	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some errors")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exist for roomID otherwise false
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID == 3 {
		return false, errors.New("some errors")
	}
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of availability room, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	if start.Format("2006-01-02") == "2050-01-01" && end.Format("2006-01-02") == "2050-01-03" {
		rooms = append(rooms, models.Room{
			ID:       1,
			RoomName: "General's quarter",
		})
	}
	if start.Format("2006-01-02") == "2000-01-01" && end.Format("2006-01-02") == "2000-01-03" {
		return rooms, errors.New("some errors")
	}
	return rooms, nil
}

// GetRoomByID return room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some errors")
	}
	return room, nil
}
