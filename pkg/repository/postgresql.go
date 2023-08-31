package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
)

// TODO: вынести в конфиг
const (
	host     = "localhost"
	port     = 5434
	user     = "admin"
	password = "admin"
	dbname   = "user_segmentation"
)

const initDBQuery = `
	CREATE TABLE IF NOT EXISTS users (
	  id VARCHAR(30) PRIMARY KEY,
	  order_id SERIAL NOT NULL UNIQUE,
	  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS test_members (
	  id SERIAL PRIMARY KEY,
	  user_id VARCHAR NOT NULL ,
	  segment_id INTEGER NOT NULL ,
	  created_at TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  expires_at TIMESTAMP,
	  CONSTRAINT UC_MEMBER UNIQUE (user_id, segment_id)
	);

	CREATE TABLE IF NOT EXISTS segments (
	  id SERIAL PRIMARY KEY,
	  slug VARCHAR(30) NOT NULL UNIQUE ,
	  percentage NUMERIC (4, 2),
	  is_auto BOOLEAN NOT NULL,
	  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  CHECK ((percentage > 0 AND percentage <= 100) OR is_auto is false)
	);


	CREATE TYPE test_status AS ENUM ('DELETED', 'ADDED');

	CREATE TABLE IF NOT EXISTS test_members_history (
	  id SERIAL PRIMARY KEY,
	  user_id VARCHAR(30) NOT NULL,
	  slug VARCHAR(30) NOT NULL,
	  status test_status,
	  happened_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	ALTER TABLE "test_members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

	ALTER TABLE "test_members" ADD FOREIGN KEY ("segment_id") REFERENCES "segments" ("id");

	CREATE INDEX idx_tests_user_id ON test_members (user_id);
	CREATE INDEX idx_tests_segment_id ON test_members (segment_id);
	CREATE INDEX idx_timestamp ON test_members (created_at, expires_at);
	CREATE INDEX idx_history_user_id ON test_members_history (user_id);

`

func NewPostgresDB(log *slog.Logger) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Error("Failed to connect to database")
		return nil, err
	}
	//_, err = db.Exec(initDBQuery)
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, err
	//}
	connectMessage := fmt.Sprintf("successful connected to database: %s", "user_segmentation")
	log.Info(connectMessage)
	return db, nil
}
