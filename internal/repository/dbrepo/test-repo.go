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

// GetUserByID returns a user by id
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

// UpdateUser updates a user in database
func (m *testDBRepo) UpdateUser(u models.User) error {

	return nil
}

// Authenticate authenticates a user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "me@hehe.com" {
		return 1, "", nil
	}
	return 0, "", errors.New("some error")
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// GetReservationByID takes reservation by id
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservation models.Reservation

	return reservation, nil
}

// UpdateReservation updates a reservation in database
func (m *testDBRepo) UpdateReservation(u models.Reservation) error {

	return nil
}

// DeleteReservation deletes a reservation by id from database
func (m *testDBRepo) DeleteReservation(id int) error {

	return nil
}

// UpdateProcessedForReservation updates processed of reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}

// AllRooms gets all rooms in database
func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	rooms := []models.Room{
		{
			ID:       1,
			RoomName: "Room1",
		},
		{
			ID:       2,
			RoomName: "Room2",
		},
	}
	return rooms, nil
}

// GetRestrictionsForRoomByDate returns restrictions for room by date range
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	restrictions := []models.RoomRestriction{
		{
			ID:            1,
			ReservationID: 1,
		},
		{
			ID:            2,
			ReservationID: 0,
		},
		{
			ID:            3,
			ReservationID: 1,
		},
	}
	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil
}

// DeleteBlockByID deletes a room restriction
func (m *testDBRepo) DeleteBlockByID(id int) error {

	return nil
}
