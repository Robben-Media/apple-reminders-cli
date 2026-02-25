package commands

import (
	"strings"

	"github.com/spf13/cobra"
)

func newUncompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uncomplete <id>",
		Short: "Mark a reminder incomplete",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			remindersService, _, err := ensureServices()
			if err != nil {
				return err
			}
			item, err := remindersService.Uncomplete(strings.TrimSpace(args[0]))
			if err != nil {
				return err
			}
			return writeJSON(item)
		},
	}
}
