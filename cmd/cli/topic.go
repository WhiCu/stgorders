package cli

import (
	"github.com/WhiCu/stgorders/cmd/kafka"
	"github.com/spf13/cobra"
)

var topicCmd = &cobra.Command{
	Use:   "topic",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return kafka.CreateTopic(cfg.Kafka.Topic, cfg.Kafka.Brokers, log.WithGroup("topic"))
	},
}

func init() {
	rootCmd.AddCommand(topicCmd)
}
