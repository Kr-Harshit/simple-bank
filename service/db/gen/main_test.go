package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KHarshit1203/simple-bank/service/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testConnPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connnect to database", err)
	}
	defer testConnPool.Close()

	testStore = NewSQLStore(testConnPool)
	os.Exit(m.Run())
}
