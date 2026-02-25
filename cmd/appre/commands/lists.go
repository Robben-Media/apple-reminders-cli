package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newListsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lists",
		Short: "List, create, and delete reminder lists",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, listsService, err := ensureServices()
			if err != nil {
				return err
			}
			lists, err := listsService.List()
			if err != nil {
				return err
			}
			return writeJSON(lists)
		},
	}

	createCmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a reminder list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, listsService, err := ensureServices()
			if err != nil {
				return err
			}
			source, _ := cmd.Flags().GetString("source")
			item, err := listsService.Create(strings.TrimSpace(args[0]), source)
			if err != nil {
				return err
			}
			return writeJSON(item)
		},
	}
	createCmd.Flags().String("source", "", "Account/source name (e.g. iCloud)")

	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a reminder list by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, listsService, err := ensureServices()
			if err != nil {
				return err
			}
			id := strings.TrimSpace(args[0])
			if id == "" {
				return fmt.Errorf("id is required")
			}
			if err := listsService.Delete(id); err != nil {
				return err
			}
			return writeJSON(map[string]any{"deleted": true, "id": id})
		},
	}

	cmd.AddCommand(createCmd)
	cmd.AddCommand(deleteCmd)
	return cmd
}
