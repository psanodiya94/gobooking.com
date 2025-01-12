package dbrepo

import (
	"context"
	"errors"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// InsertReservation insert a reservation into database
func (psql *dbPostgresRepo) InsertReservation(res models.Reservation) (int, error) {
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
func (psql *dbPostgresRepo) InsertRoomRestriction(res models.RoomRestriction) error {
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
func (psql *dbPostgresRepo) SearchAvailabilityForDatesByRoomId(roomId int, checkIn, checkOut time.Time) (bool, error) {
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
func (psql *dbPostgresRepo) SearchAvailabilityForAllRooms(checkIn, checkOut time.Time) ([]models.Room, error) {
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
func (psql *dbPostgresRepo) GetRoomById(id int) (models.Room, error) {
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

// GetUserById gets a user by id
func (psql *dbPostgresRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select
    			id, first_name, last_name, email, password, access_level, created_at, updated_at
			from
			    users
            where
                id = $1;`
	// indent on

	var user models.User

	row := psql.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser modifies user in database
func (psql *dbPostgresRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			update
			    users
			set
			    first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
			where
			    id = $6;`
	// indent on

	_, err := psql.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
		user.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (psql *dbPostgresRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select
    			id, password
			from
			    users
            where
                email = $1;`
	// indent on

	var id int
	var hash string

	row := psql.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hash)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hash, nil
}

// AllReservations returns a slice of all reservations
func (psql *dbPostgresRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select
    			r.id, r.first_name, r.last_name, r.email, r.phone, r.check_in, r.check_out,
                r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
			from
			    reservations r
            left join
                rooms rm
            on
                (r.room_id = rm.id)
            order by
                r.check_in asc;`
	// indent on

	var reservations []models.Reservation

	rows, err := psql.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.Id,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.CheckIn,
			&reservation.CheckOut,
			&reservation.RoomId,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Processed,
			&reservation.Room.Id,
			&reservation.Room.RoomName,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations
func (psql *dbPostgresRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select
    			r.id, r.first_name, r.last_name, r.email, r.phone, r.check_in, r.check_out,
                r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
			from
			    reservations r
            left join
                rooms rm
            on
                (r.room_id = rm.id)
            where
                processed = 0
            order by
                r.check_in asc;`
	// indent on

	var reservations []models.Reservation

	rows, err := psql.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.Id,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.CheckIn,
			&reservation.CheckOut,
			&reservation.RoomId,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Processed,
			&reservation.Room.Id,
			&reservation.Room.RoomName,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// GetReservationById returns one reservation by id
func (psql *dbPostgresRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select
    			r.id, r.first_name, r.last_name, r.email, r.phone, r.check_in, r.check_out,
                r.room_id, r.created_at, r.updated_at, r.processed,  rm.id, rm.room_name
			from
			    reservations r
            left join
                rooms rm
            on
                (r.room_id = rm.id)
            where
                r.id = $1;`
	// indent on

	var reservation models.Reservation

	row := psql.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&reservation.Id,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.CheckIn,
		&reservation.CheckOut,
		&reservation.RoomId,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.Id,
		&reservation.Room.RoomName,
	)
	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

// UpdateReservation updates reservation in database
func (psql *dbPostgresRepo) UpdateReservation(reservation models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			update
			    reservations
			set
			    first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
			where
			    id = $6;`
	// indent on

	_, err := psql.DB.ExecContext(ctx, query,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		time.Now(),
		reservation.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation deletes reservation from database
func (psql *dbPostgresRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			delete from
			    reservations
			where
			    id = $1;`
	// indent on

	_, err := psql.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProcessedForReservation updates processed field for reservation
func (psql *dbPostgresRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			update
			    reservations
			set
			    processed = $1
			where
			    id = $2;`
	// indent on

	_, err := psql.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil
}

// AllRooms returns all rooms
func (psql *dbPostgresRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select
    			id, room_name, created_at, updated_at
			from
			    rooms
            order by
                room_name;`
	// indent on

	var rooms []models.Room

	rows, err := psql.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.Id,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdatedAt,
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

func (psql *dbPostgresRepo) GetRestrictionsForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			select
    			id, coalesce(reservation_id, 0), restriction_id, room_id, check_in, check_out
			from
			    room_restrictions
            where
                room_id = $1 and $2 < check_out and $3 >= check_in;`
	// indent on

	var restrictions []models.RoomRestriction

	rows, err := psql.DB.QueryContext(ctx, query, roomId, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var restriction models.RoomRestriction
		err := rows.Scan(
			&restriction.Id,
			&restriction.ReservationId,
			&restriction.RestrictionId,
			&restriction.RoomId,
			&restriction.CheckIn,
			&restriction.CheckOut,
		)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, restriction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil
}

// InsertBlockForRoom insert block in calendar for date
func (psql *dbPostgresRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			insert into
			    room_restrictions (check_in, check_out, room_id, restriction_id, created_at, updated_at)
			values
			    ($1, $2, $3, $4, $5, $6);`
	// indent on

	_, err := psql.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

// DeleteBlockById delete block by id
func (psql *dbPostgresRepo) DeleteBlockById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// indent off
	query := `
			delete from
			    room_restrictions
			where
			    id = $1;`
	// indent on

	_, err := psql.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
