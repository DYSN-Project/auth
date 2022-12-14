package cmd

import (
	"dysn/auth/config"
	"dysn/auth/pkg/db"
	"dysn/auth/pkg/log"
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

		dbSql, err := database.DB()
		if err != nil {
			panic(err)
		}
		if err := goose.Run(args[0], dbSql, path, args[1:]...); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mgrCmd)
}
