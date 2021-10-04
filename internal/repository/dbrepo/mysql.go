package dbrepo

import (
	"context"
	"time"

	"github.com/DungBuiTien1999/bookings/internal/models"
)

func (m *mysqlDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database
func (m *mysqlDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// request last longer 3 second so discard write record into db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservations 
	(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
	values (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	var newId int
	qry := `select LAST_INSERT_ID()`
	err = m.DB.QueryRowContext(ctx, qry).Scan(&newId)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into database
func (m *mysqlDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	// request last longer 3 second so discard write record into db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions 
	(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at) 
	values (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exist for roomID otherwise false
func (m *mysqlDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	// request last longer 3 second so discard write record into db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id) from room_restrictions where room_id = ? and ? < end_date and ? > start_date`

	var numRows int

	err := m.DB.QueryRowContext(ctx, query, roomID, start, end).Scan(&numRows)
	if err != nil {
		return false, nil
	}

	return numRows == 0, nil
}

// SearchAvailabilityForAllRooms returns a slice of availability room, if any, for given date range
func (m *mysqlDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	// request last longer 3 second so discard write record into db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select
				r.id, r.room_name
			from
				rooms as r
			where
				r.id not in 
				(select rr.room_id from room_restrictions as rr where ? < rr.end_date and ? > rr.start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID return room by id
func (m *mysqlDBRepo) GetRoomByID(id int) (models.Room, error) {
	// request last longer 3 second so discard write record into db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select id, room_name, created_at, updated_at from rooms where id = ?
	`

	var room models.Room
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil
}
