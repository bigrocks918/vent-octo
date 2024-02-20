package store

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAvailabilitiesFromDB(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	// Mock rows data
	rows := sqlmock.NewRows([]string{"id", "local_date", "status", "product_name", "vacancies", "available", "availability_price", "availability_currency"}).
		AddRow("id1", time.Now(), "AVAILABLE", "Product 1", 10, true, 100.0, "USD")

	// Expectations
	mock.ExpectQuery("^SELECT (.+) FROM availabilities a INNER JOIN products p").WillReturnRows(rows)

	_, err := GetAvailabilitiesFromDB(db, time.Now(), time.Now())
	if err != nil {
		t.Errorf("Error was not expected, got %v", err)
	}

	// Assert all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetAvailabilityByIdFromDB(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	query := "SELECT id, local_date, status, product_id, vacancies, available, price, currency FROM availabilities WHERE id = \\$1"
	mock.ExpectQuery(query).WithArgs("test_id").WillReturnRows(sqlmock.NewRows([]string{"id", "local_date", "status", "product_id", "vacancies", "available", "price", "currency"}).
		AddRow("test_id", time.Now(), "AVAILABLE", "product_id", 5, true, 100.0, "USD"))

	_, err := GetAvailabilityByIdFromDB(db, "test_id")
	if err != nil {
		t.Fatalf("error was not expected while fetching data: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddAvailabilityIntoDB(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	// Begin transaction
	mock.ExpectBegin()

	// Expect the product select query
	mock.ExpectQuery("SELECT id, name, capacity, price, currency FROM products WHERE id = \\$1").
		WithArgs("product_id").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "capacity", "price", "currency"}).
			AddRow("product_id", "Product Name", 100, 50.0, "USD"))

	// Expect the insert into availabilities
	mock.ExpectExec("INSERT INTO availabilities").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Commit transaction
	mock.ExpectCommit()

	err := AddAvailabilityIntoDB(db, "product_id", time.Now(), time.Now(), 100.0, "USD")
	if err != nil {
		t.Errorf("error was not expected while inserting data: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
