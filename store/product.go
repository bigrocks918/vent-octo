package store

import (
	"database/sql"
	"fmt"
	"octo-api/model"
)

func GetProductsFromDB(db *sql.DB) ([]model.Product, error) {
	rows, err := db.Query("SELECT id, name, capacity, price, currency FROM products")
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Capacity, &p.Price, &p.Currency); err != nil {
			// log.Fatal(err)
			fmt.Println(err.Error())
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetProductFromDB(db *sql.DB, productId string) (*model.Product, error) {
	var p model.Product
	err := db.QueryRow("SELECT id, name, capacity, price, currency FROM products WHERE id = $1", productId).Scan(&p.ID, &p.Name, &p.Capacity, &p.Price, &p.Currency)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		return nil, err
	}
	return &p, nil
}

func InsertProductIntoDB(db *sql.DB, productInfo model.Product) error {
	tx, err := db.Begin()
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		return err
	}

	// Insert the product
	productStmt := "INSERT INTO products (id, name, capacity, price, currency) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.Exec(productStmt, productInfo.ID, productInfo.Name, productInfo.Capacity, productInfo.Price, productInfo.Currency)
	if err != nil {
		tx.Rollback()
		// log.Fatal(err)
		fmt.Println(err.Error())
		return err
	}

	return tx.Commit()
}
