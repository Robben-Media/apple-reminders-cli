package commands

import (
	"strings"

	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show a reminder by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			remindersService, _, err := ensureServices()
			if err != nil {
				return err
			}
			item, err := remindersService.Show(strings.TrimSpace(args[0]))
			if err != nil {
				return err
			}
			return writeJSON(item)
		},
	}
}
