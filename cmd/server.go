package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/KHarshit1203/simple-bank/api"
	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	"github.com/KHarshit1203/simple-bank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	// init DB connection
	connPool, err := pgxpool.New(context.Background(), config.Database.Source)
	if err != nil {
		log.Fatalf("error connnecting to database; %v", err)
	}
	defer connPool.Close()
	log.Print("successfully established database connection!!!")

	// DB migartion
	if err := migrateDB(config.Database.MigrateSource, config.Database.Source); err != nil {
		log.Fatalf("error migrating database; %v", err)
	}
	log.Print("sucessfully migrated database!!!")

	// init Store
	store := db.NewSQLStore(connPool)

	// init server
	if err := runHttpServer(config, store); err != nil {
		log.Fatalf("error starting HTTP server; %v", err)
	}
}

// runHttpServer starts HTTP server at
func runHttpServer(config util.Config, store db.Store) error {
	server, err := api.NewServer(config, store)
	if err != nil {
		return fmt.Errorf("error creating new api server: %v", err)
	}

	address := fmt.Sprintf("%s:%s", config.App.Address, config.App.Port)
	return server.Start(address)
}

// migateDB migrates the database to latest migration version
func migrateDB(migrateURL, dbSource string) error {
	m, err := migrate.New(migrateURL, dbSource)
	if err != nil {
		return fmt.Errorf("error creating new migartion instance, %v", err)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("error migrating to latest version. %v", err)
		} else {
			log.Printf("no migration change!!!")
		}
	}
	return nil
}
