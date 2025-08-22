package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/WhiCu/stgorders/cmd/app"
	"github.com/WhiCu/stgorders/cmd/kafka"
	"github.com/WhiCu/stgorders/cmd/migrate"
	"github.com/WhiCu/stgorders/internal/config"
	"github.com/WhiCu/stgorders/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg *config.Config
var log *slog.Logger

var rootCmd = &cobra.Command{
	Use:     "stgorders",
	Short:   "",
	Long:    ``,
	Version: "0.0.1",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		cfg, err = config.LoadWithDefault("config/config.yaml")
		if err != nil {
			return err
		}

		log = logger.GetLogger(&cfg.Logger)
		log.Info("logger created",
			slog.String("level", cfg.Logger.Level),
			slog.String("path", cfg.Logger.Path),
			slog.Int("size", cfg.Logger.Size),
			slog.Bool("compress", cfg.Logger.Compress),
		)
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = migrate.MigrateDB(cfg.Storage.DSN(), cfg.Migrate.Dir, log.WithGroup("migrate"))
		if err != nil {
			return err
		}
		if viper.GetBool("topic") {
			kafka.CreateTopic(cfg.Kafka.Topic, cfg.Kafka.Brokers, log.WithGroup("topic"))
			if err != nil {
				log.Error("could not create kafka topic", slog.String("ERR", err.Error()))
				// return err
			}
		}
		return app.NewAPP(cfg, log.WithGroup("app"))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().BoolP("topic", "t", false, "create default kafka topic")
	err := viper.BindPFlag("topic", rootCmd.Flags().Lookup("topic"))
	if err != nil {
		panic(err)
	}
}
