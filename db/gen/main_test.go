package db

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v4"
)

const (
	DATABASE_URL = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var (
	testQueries *Queries
	testDb      *pgx.Conn
)

func TestMain(m *testing.M) {
	var err error
	testDb, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatal("cannot connnect to database", err)
	}
	defer testDb.Close(context.Background())

	testQueries = New(testDb)
	m.Run()
}
