package cmd

import (
	"github.com/DYSN-Project/auth/config"
	"github.com/DYSN-Project/auth/internal/packages/db"
	"github.com/DYSN-Project/auth/internal/packages/log"
	"github.com/pressly/goose"
	"github.com/spf13/cobra"
)

const path = "migrations"

var mgrCmd = &cobra.Command{
	Use: "migration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		logger := log.NewLogger()
		database := db.StartDB(cfg, logger)
		defer db.CloseDB(database, logger)
		if err := goose.Run(args[0], database.DB(), path, args[1:]...); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mgrCmd)
}
