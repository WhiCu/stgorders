package cli

import (
	"github.com/WhiCu/stgorders/cmd/app"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.NewAPP(cfg, log.WithGroup("app"))
	},
}

func init() {
	rootCmd.AddCommand(appCmd)
}
