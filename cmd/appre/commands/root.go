package commands

import (
	"fmt"
	"sync"

	ekrem "github.com/BRO3886/go-eventkit/reminders"
	"github.com/Robben-Media/apple-reminders-cli/internal/service"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:           "appre",
		Short:         "Apple Reminders CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("no command provided")
		},
	}

	initOnce     sync.Once
	initServices struct {
		reminders *service.RemindersService
		lists     *service.ListsService
		err       error
	}
)

func init() {
	rootCmd.AddCommand(newListsCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newShowCmd())
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newEditCmd())
	rootCmd.AddCommand(newCompleteCmd())
	rootCmd.AddCommand(newUncompleteCmd())
	rootCmd.AddCommand(newDeleteCmd())
	rootCmd.AddCommand(newWatchCmd())
	rootCmd.AddCommand(newRPCCmd())
	rootCmd.AddCommand(newVersionCmd())
}

func Execute() error {
	return rootCmd.Execute()
}

func ensureServices() (*service.RemindersService, *service.ListsService, error) {
	initOnce.Do(func() {
		client, err := ekrem.New()
		if err != nil {
			initServices.err = err
			return
		}
		initServices.reminders = service.NewRemindersService(client)
		initServices.lists = service.NewListsService(client)
	})
	return initServices.reminders, initServices.lists, initServices.err
}
