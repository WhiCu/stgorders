package cli

import (
	"github.com/WhiCu/stgorders/cmd/migrate"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrate.MigrateDB(cfg.Storage.DSN(), cfg.Migrate.Dir, log.WithGroup("migrate"))
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
