package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KHarshit1203/simple-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config := util.Config{
		Database: util.DBConfig{
			Source: "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable",
		},
		App: util.APPConfig{
			Address: "localhost",
			Port:    "8080",
		},
	}

	testConnPool, err := pgxpool.New(context.Background(), config.Database.Source)
	if err != nil {
		log.Fatal("cannot connnect to database", err)
	}
	defer testConnPool.Close()

	testStore = NewSQLStore(testConnPool)
	os.Exit(m.Run())
}
