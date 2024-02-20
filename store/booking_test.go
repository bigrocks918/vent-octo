package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// func TestCreateBooking(t *testing.T) {
// 	db, mock := NewMock()
// 	defer db.Close()

// 	booking := model.Booking{
// 		ID:             "booking_id",
// 		Status:         "NEW",
// 		AvailabilityId: "availability_id",
// 		Units:          2,
// 		Price:          200.0,
// 		Currency:       "USD",
// 	}

// 	mock.ExpectBegin()

// 	mock.ExpectQuery("SELECT vacancies FROM availabilities WHERE id = \\$1").
// 		WithArgs(booking.AvailabilityId).
// 		WillReturnRows(sqlmock.NewRows([]string{"vacancies"}).AddRow(10))

// 	mock.ExpectExec("INSERT INTO bookings \\(id, status, availability_id, units, price, currency\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\)").
// 		WithArgs(booking.ID, booking.Status, booking.AvailabilityId, booking.Units, booking.Price, booking.Currency).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	mock.ExpectExec("UPDATE availabilities SET vacancies = vacancies - \\$1 WHERE id = \\$2").
// 		WithArgs(booking.Units, booking.AvailabilityId).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	mock.ExpectCommit()

// 	err := CreateBooking(db, booking)
// 	if err != nil {
// 		t.Fatalf("error was not expected while creating booking: %s", err)
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Fatalf("there were unmet expectations: %s", err)
// 	}
// }

// func TestGetAllBookings(t *testing.T) {
// 	db, mock := NewMock()
// 	defer db.Close()

// 	mock.ExpectQuery("SELECT id, status, availability_id, price, currency FROM bookings").
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "availability_id", "price", "currency"}).
// 			AddRow("booking_id", "CONFIRMED", "availability_id", 100.0, "USD"))

// 	mock.ExpectQuery("SELECT id, booking_id, price, currency FROM booking_units WHERE booking_id = \\$1").
// 		WithArgs("booking_id").
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "booking_id", "price", "currency"}).
// 			AddRow("unit_id", "booking_id", 100.0, "USD"))

// 	_, err := GetAllBookings(db)
// 	if err != nil {
// 		t.Fatalf("error was not expected while fetching all bookings: %s", err)
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Fatalf("there were unmet expectations: %s", err)
// 	}
// }

func TestGetBookingByID(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	bookingID := "booking_id"
	mock.ExpectQuery("SELECT id, status, availability_id, price, currency FROM bookings WHERE id = \\$1").
		WithArgs(bookingID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "availability_id", "price", "currency"}).
			AddRow(bookingID, "CONFIRMED", "availability_id", 100.0, "USD"))

	mock.ExpectQuery("SELECT id, booking_id, price, currency FROM booking_units WHERE booking_id = \\$1").
		WithArgs(bookingID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "booking_id", "price", "currency"}).
			AddRow("unit_id", bookingID, 100.0, "USD"))

	_, err := GetBookingByID(db, bookingID)
	if err != nil {
		t.Fatalf("error was not expected while fetching booking by ID: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unmet expectations: %s", err)
	}
}
