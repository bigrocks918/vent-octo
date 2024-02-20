package store

import (
	"octo-api/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetProductsFromDB(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, capacity, price, currency FROM products").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "capacity", "price", "currency"}).
			AddRow("product_id", "Product 1", 100, 1000.0, "USD").
			AddRow("product_id2", "Product 2", 200, 2000.0, "EUR"))

	products, err := GetProductsFromDB(db)
	if err != nil {
		t.Fatalf("error was not expected while fetching products: %s", err)
	}

	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unmet expectations: %s", err)
	}
}

func TestGetProductFromDB(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	query := "SELECT id, name, capacity, price, currency FROM products WHERE id = \\$1"
	mock.ExpectQuery(query).WithArgs("product_id").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "capacity", "price", "currency"}).
		AddRow("product_id", "Product Name", 100, 50.0, "USD"))

	_, err := GetProductFromDB(db, "product_id")
	if err != nil {
		t.Errorf("error was not expected while fetching product by ID: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

func TestInsertProductIntoDB(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO products").WithArgs(sqlmock.AnyArg(), "Product Name", 100, 50.0, "USD").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := InsertProductIntoDB(db, model.Product{ID: "product_id", Name: "Product Name", Capacity: 100, Price: 50.0, Currency: "USD"})
	if err != nil {
		t.Errorf("error was not expected while inserting product: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
