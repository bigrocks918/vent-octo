package store

import (
	"database/sql"
	"errors"
	"fmt"
	"octo-api/model"
)

// CreateBooking inserts a new booking into the database and updates availability, with a check for sufficient vacancies.
func CreateBooking(db *sql.DB, booking model.Booking) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Check if there are enough vacancies for the booking
	var vacancies int
	checkStmt := "SELECT vacancies FROM availabilities WHERE id = $1"
	err = tx.QueryRow(checkStmt, booking.AvailabilityId).Scan(&vacancies)
	if err != nil {
		tx.Rollback()
		return err
	}
	if vacancies < booking.Units {
		tx.Rollback()
		return errors.New("insufficient vacancies for the requested booking")
	}

	// Insert the booking
	bookingStmt := "INSERT INTO bookings (id, status, availability_id, units, price, currency) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = tx.Exec(bookingStmt, booking.ID, booking.Status, booking.AvailabilityId, booking.Units, booking.Price, booking.Currency)
	if err != nil {
		tx.Rollback()
		return err
	}

	// check if the vacancies will be 0
	emptyFlg := (vacancies == booking.Units)

	// Update the availability
	var updateStmt string
	var result sql.Result
	if emptyFlg {
		// No need to update availability.status and availability.available
		updateStmt = "UPDATE availabilities SET vacancies = vacancies - $1 WHERE id = $2"
		result, err = tx.Exec(updateStmt, booking.Units, booking.AvailabilityId)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// Update availability.status to SOLD_OUT and availability.available to false
		updateStmt = "UPDATE availabilities SET vacancies = vacancies - $1, status = $2, available = $3 WHERE id = $4"
		result, err = tx.Exec(updateStmt, booking.Units, "SOLD_OUT", false, booking.AvailabilityId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	affected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if affected == 0 {
		tx.Rollback()
		return errors.New("failed to update availability, possibly due to insufficient vacancies")
	}

	return tx.Commit()
}

// ConfirmBooking updates the booking's status to CONFIRMED and generates tickets.
func ConfirmBooking(db *sql.DB, bookingID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Generate tickets and update booking status
	// This is a simplified approach. Adjust according to your schema and requirements.
	updateStmt := "UPDATE bookings SET status = 'CONFIRMED' WHERE id = $1 RETURNING availability_id, units"
	var availabilityID string
	var units int
	err = tx.QueryRow(updateStmt, bookingID).Scan(&availabilityID, &units)
	if err != nil {
		tx.Rollback()
		return err
	}

	ava, err := GetAvailabilityByIdFromDB(db, availabilityID)
	if err != nil {
		return err
	}

	prod, err := GetProductFromDB(db, ava.ProductId)
	if err != nil {
		return err
	}

	// Generate tickets for each unit. This could be more complex in a real scenario.
	for i := 0; i < units; i++ {
		ticketID := fmt.Sprintf("TICKET-%d-%s", i, bookingID)
		insertTicketStmt := "INSERT INTO booking_units (id, booking_id, price, currency) VALUES ($1, $2, $3, $4)"
		_, err := tx.Exec(insertTicketStmt, ticketID, bookingID, prod.Price, prod.Currency)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetAllBookings get all lists of booking information
func GetAllBookings(db *sql.DB) ([]model.BookingPayload_Rs, error) {
	var query string
	var rows *sql.Rows
	var err error
	fmt.Println("------------------ Insides")

	query = "SELECT id, status, availability_id, price, currency FROM bookings"
	rows, err = db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var bookings []model.BookingPayload_Rs
	for rows.Next() {
		var curBooking model.BookingPayload_Rs
		if err := rows.Scan(
			&curBooking.ID,
			&curBooking.Status,
			&curBooking.AvailabilityId,
			&curBooking.Price,
			&curBooking.Currency,
		); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		// Retrieve booking units
		unitsQuery := "SELECT id, booking_id, price, currency FROM booking_units WHERE booking_id = $1"
		unitsRows, err := db.Query(unitsQuery, curBooking.ID)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		defer unitsRows.Close()

		curBooking.Units = []model.BookingUnitPayload_Rs{}
		for unitsRows.Next() {
			var unit model.BookingUnitPayload_Rs
			if err := unitsRows.Scan(&unit.ID, &unit.BookingId, &unit.Price, &unit.Currency); err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			curBooking.Units = append(curBooking.Units, unit)
		}
		bookings = append(bookings, curBooking)
	}

	return bookings, nil
}

// GetBookingByID retrieves a booking and its units by ID.
func GetBookingByID(db *sql.DB, bookingID string) (*model.BookingPayload_Rs, error) {
	booking := &model.BookingPayload_Rs{}

	// Retrieve the booking
	bookingQuery := "SELECT id, status, availability_id, price, currency FROM bookings WHERE id = $1"
	err := db.QueryRow(bookingQuery, bookingID).Scan(&booking.ID, &booking.Status, &booking.AvailabilityId, &booking.Price, &booking.Currency)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Retrieve booking units
	unitsQuery := "SELECT id, booking_id, price, currency FROM booking_units WHERE booking_id = $1"
	rows, err := db.Query(unitsQuery, bookingID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	booking.Units = []model.BookingUnitPayload_Rs{}
	for rows.Next() {
		var unit model.BookingUnitPayload_Rs
		if err := rows.Scan(&unit.ID, &unit.BookingId, &unit.Price, &unit.Currency); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		booking.Units = append(booking.Units, unit)
	}

	return booking, nil
}
