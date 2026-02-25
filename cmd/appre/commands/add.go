package commands

import (
	"strings"

	"github.com/Robben-Media/apple-reminders-cli/internal/service"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var (
		listName string
		dueIn    string
		priority string
		notes    string
		url      string
	)

	cmd := &cobra.Command{
		Use:   "add <title>",
		Short: "Add a reminder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			remindersService, _, err := ensureServices()
			if err != nil {
				return err
			}

			due, err := parseDateTimeFlag(dueIn)
			if err != nil {
				return err
			}
			parsedPriority, err := parsePriorityFlag(priority)
			if err != nil {
				return err
			}

			input := service.AddReminderInput{
				Title:    strings.TrimSpace(args[0]),
				List:     strings.TrimSpace(listName),
				Due:      due,
				Priority: parsedPriority,
			}
			if cmd.Flags().Changed("notes") {
				input.Notes = &notes
			}
			if cmd.Flags().Changed("url") {
				input.URL = &url
			}

			item, err := remindersService.Add(input)
			if err != nil {
				return err
			}
			return writeJSON(item)
		},
	}

	cmd.Flags().StringVar(&listName, "list", "", "List name")
	cmd.Flags().StringVar(&dueIn, "due", "", "Due date/time (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&priority, "priority", "", "Priority: none|high|medium|low or 0|1|5|9")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes")
	cmd.Flags().StringVar(&url, "url", "", "URL")
	return cmd
}
