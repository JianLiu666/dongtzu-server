package cmd

import (
	"dongtzu/pkg/repository/arangodb"

	"github.com/spf13/cobra"
)

var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "enjoy for it.",
	Long:  `No more description.`,
	RunE:  RunMigrationCmd,
}

func init() {
	rootCmd.AddCommand(migrationCmd)
}

func RunMigrationCmd(cmd *cobra.Command, args []string) error {
	// TODO: 之後在改成每個服務啟動時, 用 DI 把 repository 注入進去
	// NOTICE: 在啟動服務之前，一定要先把需要的 repository 初始化完畢

	// 1. Initial repositories
	arangodb.Init()
	arangodb.Migration()

	return nil
}
