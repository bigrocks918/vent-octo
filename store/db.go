package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ConnectToDB() *sql.DB {
	var (
		// host     = "localhost"
		host = os.Getenv("DB_HOST")
		port = 5432
		// user     = "postgres"
		user = os.Getenv("DB_USER")
		// password = "postgres"
		password = os.Getenv("DB_PASSWORD")
		// dbname   = "ventrata_octo"
		dbname = os.Getenv("DB_NAME")
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
	}
	err = db.Ping()
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
	}
	fmt.Println("Successfully connected!")
	return db
}
