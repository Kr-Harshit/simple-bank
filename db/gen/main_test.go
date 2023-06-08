package db

import (
	"context"
	"log"
	"testing"

	"github.com/KHarshit1203/simple-bank/util"
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
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDb, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connnect to database", err)
	}
	defer testDb.Close(context.Background())

	testQueries = New(testDb)
	m.Run()
}
