package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"roomate/sqlc"
)

var (
	Queries *sqlc.Queries
	DB      *sql.DB
)

func Init() error {
	connection := "postgres://postgres.umxljfuepfivmaspbiym:r@hS37z*6Pg8Ts4@aws-0-eu-central-1.pooler.supabase.com:5432/postgres"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	Queries = sqlc.New(db)
	DB = db
	return nil
}
