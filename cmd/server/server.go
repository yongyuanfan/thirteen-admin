package server

import (
	appRuntime "thirteen-admin/internal/runtime"

	"github.com/spf13/cobra"
)

var (
	configYAML string
	NewCommand = &cobra.Command{
		Use:          "server",
		Short:        "Start HTTP server",
		Example:      "thirteen-admin server -c config/settings.yml",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

// init 注册服务启动命令参数。
func init() {
	NewCommand.PersistentFlags().StringVarP(&configYAML, "config", "c", "config/settings.yml", "Start http with provided configuration file")
}

// run 启动 HTTP 与 gRPC 服务并处理优雅退出。
func run() error {
	return appRuntime.RunServer(configYAML)
}
