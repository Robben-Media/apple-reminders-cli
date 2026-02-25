package commands

import (
	"github.com/Robben-Media/apple-reminders-cli/internal/rpc"
	"github.com/spf13/cobra"
)

func newRPCCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rpc",
		Short: "Run JSON-RPC 2.0 server on stdin/stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			remindersService, listsService, err := ensureServices()
			if err != nil {
				return err
			}

			server := rpc.NewServer(remindersService, listsService, nil, nil)
			return server.Serve(cmd.Context())
		},
	}
}
