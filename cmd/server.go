package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/KHarshit1203/simple-bank/api"
	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	"github.com/KHarshit1203/simple-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

func init() {
	cmd.PersistentFlags().StringP("config", "c", "", "Configuration file (required)")
}

var cmd = &cobra.Command{
	Use:   "simplebank",
	Short: "simplebank is a fiber based http server",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("error parsing config flag, %v", err)
		}

		config, err := util.LoadConfig(cfgFile)
		if err != nil {
			log.Fatalf("error loading configurations, %v", err)
		}

		runServer(config)
	},
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing simplebank cmd, %v", err)
	}
}

// runs simple bank server
func runServer(config util.Config) {
	connPool, err := pgxpool.New(context.Background(), config.Database.Source)
	if err != nil {
		log.Fatal("cannot connnect to database: ", err)
	}
	defer connPool.Close()

	store := db.NewSQLStore(connPool)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("unable to create server: %v", err)
	}

	address := fmt.Sprintf("%s:%s", config.App.Address, config.App.Port)
	err = server.Start(address)
	if err != nil {
		log.Fatal("cannot start the server: ", err)
	}
}
