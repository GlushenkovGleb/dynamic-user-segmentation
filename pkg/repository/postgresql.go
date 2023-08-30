package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5434
	user     = "admin"
	password = "admin"
	dbname   = "user_segmentation"
)

const initialSQL = `
	CREATE TABLE IF NOT EXISTS users (
	  id VARCHAR(30) PRIMARY KEY,
	  order_id SERIAL NOT NULL UNIQUE,
	  created_at TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP 
	);

	CREATE TABLE IF NOT EXISTS test_members (
	  id SERIAL PRIMARY KEY,
	  user_id VARCHAR NOT NULL ,
	  slug_id INTEGER NOT NULL ,
	  created_at TIMESTAMP NOT NULL DEFAULT  CURRENT_TIMESTAMP,
	  expires_at TIMESTAMP,
	  UNIQUE (user_id, slug_id)
	);

	CREATE TABLE IF NOT EXISTS slugs (
	  id SERIAL PRIMARY KEY,
	  name VARCHAR(30) NOT NULL UNIQUE ,
	  percentage NUMERIC (5, 2),
	  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  CHECK ((percentage > 0 AND percentage <= 100) OR percentage is NULL)
	);


	CREATE TYPE test_status AS ENUM ('deleted', 'added');

	CREATE TABLE IF NOT EXISTS test_members_history (
	  id SERIAL PRIMARY KEY,
	  user_id VARCHAR(30) NOT NULL,
	  slug_name VARCHAR(30) NOT NULL,
	  status test_status,
	  happened_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	ALTER TABLE "test_members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

	ALTER TABLE "test_members" ADD FOREIGN KEY ("slug_id") REFERENCES "slugs" ("id");

	CREATE INDEX idx_tests_user_id ON test_members (user_id);
	CREATE INDEX idx_tests_slug_id ON test_members (slug_id);
	CREATE INDEX idx_timestamp ON test_members (created_at, expires_at);
	CREATE INDEX idx_history_user_id ON test_members_history (user_id);

`

func NewPostgresDB() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Println("failed to connect to database")
		return nil, err
	}
	_, err = db.Exec(initialSQL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Successfully connected!")
	return db, nil
}
