package cmd

import (
	"dongtzu/pkg/repository/arangodb"
	"dongtzu/pkg/service/scheduler"
	"os"
	"os/signal"
	"syscall"

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
	arangodb.Init()

	scheduler.Init()
	scheduler.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	return nil
}
