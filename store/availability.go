package store

import (
	"database/sql"
	"fmt"
	"octo-api/model"
	"time"

	"github.com/google/uuid"
)

// GetAvailabilitiesFromDB queries all availabilities from the database.
func GetAvailabilitiesFromDB(db *sql.DB, startDate, endDate time.Time) ([]model.AvailabilityShow, error) {
	var query string
	var rows *sql.Rows
	var err error

	if startDate.Equal(endDate) {
		// Single date query
		query = "SELECT a.id, a.local_date, a.status, p.name AS product_name, a.vacancies, a.available, a.price AS availability_price, a.currency AS availability_currency FROM availabilities a INNER JOIN products p ON a.product_id = p.id WHERE a.local_date = $1"
		rows, err = db.Query(query, startDate)
	} else {
		// Date range query
		query = "SELECT a.id, a.local_date, a.status, p.name AS product_name, a.vacancies, a.available, a.price AS availability_price, a.currency AS availability_currency FROM availabilities a INNER JOIN products p ON a.product_id = p.id WHERE a.local_date BETWEEN $1 AND $2"
		rows, err = db.Query(query, startDate, endDate)
	}

	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var availabilities []model.AvailabilityShow
	for rows.Next() {
		var cur model.AvailabilityShow
		if err := rows.Scan(
			&cur.ID,
			&cur.LocalDate,
			&cur.Status,
			&cur.ProductName,
			&cur.Vacancies,
			&cur.Available,
			&cur.Price,
			&cur.Currency,
		); err != nil {
			// log.Fatal(err)
			fmt.Println(err.Error())
			return nil, err
		}
		availabilities = append(availabilities, cur)
	}
	return availabilities, nil
}

func GetAvailabilityByIdFromDB(db *sql.DB, id string) (*model.Availability, error) {
	var a model.Availability
	err := db.QueryRow(
		"SELECT id, local_date, status, product_id, vacancies, available, price, currency FROM availabilities WHERE id = $1",
		id,
	).Scan(
		&a.ID,
		&a.LocalDate,
		&a.Status,
		&a.ProductId,
		&a.Vacancies,
		&a.Available,
		&a.Price,
		&a.Currency,
	)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		return nil, err
	}
	return &a, nil
}

func AddAvailabilityIntoDB(db *sql.DB, productID string, startDate, endDate time.Time, price float64, currency string) error {

	tx, err := db.Begin()
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		return err
	}

	indDate := startDate

	curProduct, err := GetProductFromDB(db, productID)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		return err
	}

	for !indDate.After(endDate) {

		insertAvaStmt := "INSERT INTO availabilities (id, local_date, status, product_id, vacancies, available, price, currency) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		_, err = tx.Exec(
			insertAvaStmt,
			uuid.NewString(),
			indDate,
			"AVAILABLE",
			productID,
			curProduct.Capacity,
			true,
			price,
			currency,
		)
		if err != nil {
			tx.Rollback()
			// log.Fatal(err)
			fmt.Println(err.Error())
			continue
		}

		indDate = indDate.AddDate(0, 0, 1)
	}

	return tx.Commit()
}
