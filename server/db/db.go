package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "")
	if err != nil {
		fmt.Print(err)
	}
	_, err = db.Exec(`DROP TABLE IF EXISTS public.users;`)

	//make email unique

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.users (
    "id" bigserial PRIMARY KEY,
    "username" varchar(255) NOT NULL,
    "email" varchar(255) NOT NULL UNIQUE,
    "password" varchar(255) NOT NULL
);`)

	// drop table users;

	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
