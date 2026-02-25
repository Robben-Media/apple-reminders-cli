package commands

import (
	"strings"

	"github.com/Robben-Media/apple-reminders-cli/internal/service"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var (
		listName    string
		listID      string
		completed   string
		search      string
		dueBeforeIn string
		dueAfterIn  string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List reminders",
		RunE: func(cmd *cobra.Command, args []string) error {
			remindersService, _, err := ensureServices()
			if err != nil {
				return err
			}

			completedFilter, err := parseCompletedFilter(completed)
			if err != nil {
				return err
			}
			dueBefore, err := parseDateTimeFlag(dueBeforeIn)
			if err != nil {
				return err
			}
			dueAfter, err := parseDateTimeFlag(dueAfterIn)
			if err != nil {
				return err
			}

			items, err := remindersService.List(service.ReminderQuery{
				List:      strings.TrimSpace(listName),
				ListID:    strings.TrimSpace(listID),
				Completed: completedFilter,
				Search:    strings.TrimSpace(search),
				DueBefore: dueBefore,
				DueAfter:  dueAfter,
			})
			if err != nil {
				return err
			}

			return writeJSON(items)
		},
	}

	cmd.Flags().StringVar(&listName, "list", "", "List name filter")
	cmd.Flags().StringVar(&listID, "list-id", "", "List ID filter")
	cmd.Flags().StringVar(&completed, "completed", "", "Completion filter: true/false")
	cmd.Flags().StringVar(&search, "search", "", "Search query")
	cmd.Flags().StringVar(&dueBeforeIn, "due-before", "", "Due before (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&dueAfterIn, "due-after", "", "Due after (RFC3339 or YYYY-MM-DD)")
	return cmd
}
