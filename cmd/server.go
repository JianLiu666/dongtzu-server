package cmd

import (
	"dongtzu/pkg/repository/arangodb"
	"dongtzu/pkg/repository/lineapi"
	"dongtzu/pkg/repository/zoomapi"
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
	// TODO: 之後在改成每個服務啟動時, 用 DI 把 repository 注入進去
	// NOTICE: 在啟動服務之前，一定要先把需要的 repository 初始化完畢
	arangodb.Init()
	lineapi.Init()
	zoomapi.Init()

	scheduler.Init()
	scheduler.Start()

	server.Init()
	server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	return nil
}
