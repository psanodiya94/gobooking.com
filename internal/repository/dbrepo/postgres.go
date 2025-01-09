package dbrepo

import (
	"context"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"time"
)

// InsertReservation insert a reservation into database
func (psql *dbPostgresRepository) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	stmt := `insert into 
    				reservations (
                        first_name, last_name, email, phone, 
                        check_in, check_out, room_id, 
                        created_at, updated_at
            		) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	// indent on

	var id int

	err := psql.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.CheckIn,
		res.CheckOut,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// InsertRoomRestriction insert a room restriction into database
func (psql *dbPostgresRepository) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	stmt := `insert into
    				room_restrictions (
                    	check_in, check_out, room_id, reservation_id,
                        restriction_id, created_at, updated_at
            		)
			values ($1, $2, $3, $4, $5, $6, $7)`
	// indent on

	_, err := psql.DB.ExecContext(ctx, stmt,
		res.CheckIn,
		res.CheckOut,
		res.RoomId,
		res.ReservationId,
		res.RestrictionId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityForDatesByRoomId query database with dates if available for booking room
func (psql *dbPostgresRepository) SearchAvailabilityForDatesByRoomId(roomId int, checkIn, checkOut time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select 
    			count(id) 
			from 
			    room_restrictions 
            where 
                room_id = $1 
            and 
                $2 < check_out and $3 > check_in;`
	// indent on

	var count int

	err := psql.DB.QueryRowContext(
		ctx, query, roomId, checkIn, checkOut,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any for given date range
func (psql *dbPostgresRepository) SearchAvailabilityForAllRooms(checkIn, checkOut time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select
    			r.id, r.room_name
			from
			    rooms r
            where
                r.id not in (
                    select
                        rr.room_id
                    from
                        room_restrictions rr
                    where
                        $1 < rr.check_out and $2 > rr.check_in
                );`
	// indent on

	rows, err := psql.DB.QueryContext(ctx, query, checkIn, checkOut)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.Id,
			&room.RoomName,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

// GetRoomById get a room by id
func (psql *dbPostgresRepository) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select
    			id, room_name, created_at, updated_at
			from
			    rooms
            where
                id = $1;`
	// indent on

	var room models.Room

	row := psql.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.Id,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}
