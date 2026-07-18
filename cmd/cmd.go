package cmd

import (
	"thirteen-admin/cmd/server"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "thirteen-admin",
	Short: "Vision AI command line tool",
}

// init 注册可用子命令。
func init() {
	rootCommand.AddCommand(server.NewCommand)
}

// Execute 执行根命令入口。
func Execute() error {
	return rootCommand.Execute()
}
