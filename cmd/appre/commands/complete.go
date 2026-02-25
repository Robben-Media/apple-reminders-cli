package commands

import (
	"strings"

	"github.com/spf13/cobra"
)

func newCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "complete <id>",
		Short: "Mark a reminder complete",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			remindersService, _, err := ensureServices()
			if err != nil {
				return err
			}
			item, err := remindersService.Complete(strings.TrimSpace(args[0]))
			if err != nil {
				return err
			}
			return writeJSON(item)
		},
	}
}
