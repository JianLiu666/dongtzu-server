package cmd

import (
	"github.com/spf13/cobra"
)

var subscriberCmd = &cobra.Command{
	Use:   "server",
	Short: "enjoy for it.",
	Long:  `No more description.`,
	RunE:  RunServerCmd,
}

func init() {
	rootCmd.AddCommand(subscriberCmd)
}

func RunServerCmd(cmd *cobra.Command, args []string) error {

	return nil
}
