package commands

import (
	"strings"

	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a reminder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			remindersService, _, err := ensureServices()
			if err != nil {
				return err
			}
			id := strings.TrimSpace(args[0])
			if err := remindersService.Delete(id); err != nil {
				return err
			}
			return writeJSON(map[string]any{"deleted": true, "id": id})
		},
	}
}
