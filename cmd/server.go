package cmd

import (
	"dongtzu/pkg/repository/arangodb"
	"dongtzu/pkg/service/scheduler"
	"dongtzu/pkg/service/server"
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

	server.Init()
	server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	return nil
}
