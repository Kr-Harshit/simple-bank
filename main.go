package main

import (
	"context"
	"log"

	"github.com/KHarshit1203/simple-bank/api"
	db "github.com/KHarshit1203/simple-bank/db/gen"
	"github.com/jackc/pgx/v4"
)

const (
	DATABASE_URL  = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	SEVER_ADDRESS = "0.0.0.0:8080"
)

func main() {
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatal("cannot connnect to database: ", err)
	}
	defer conn.Close(context.Background())

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(SEVER_ADDRESS)
	if err != nil {
		log.Fatal("cannot start the server: ", err)
	}
}
