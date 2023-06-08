package main

import (
	"context"
	"log"

	"github.com/KHarshit1203/simple-bank/api"
	db "github.com/KHarshit1203/simple-bank/db/gen"
	"github.com/KHarshit1203/simple-bank/util"
	"github.com/jackc/pgx/v4"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configurations: ", err)
	}

	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connnect to database: ", err)
	}
	defer conn.Close(context.Background())

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start the server: ", err)
	}
}
