package commands

import (
	"fmt"
	"strings"

	"github.com/Robben-Media/apple-reminders-cli/internal/service"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	var (
		title    string
		listName string
		dueIn    string
		clearDue bool
		priority string
		notes    string
		url      string
	)

	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit a reminder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			remindersService, _, err := ensureServices()
			if err != nil {
				return err
			}
			if clearDue && cmd.Flags().Changed("due") {
				return fmt.Errorf("--due and --clear-due are mutually exclusive")
			}

			input := service.EditReminderInput{}
			if cmd.Flags().Changed("title") {
				input.Title = &title
			}
			if cmd.Flags().Changed("list") {
				input.List = &listName
			}
			if cmd.Flags().Changed("due") {
				parsed, err := parseDateTimeFlag(dueIn)
				if err != nil {
					return err
				}
				if parsed == nil {
					return fmt.Errorf("--due must not be empty")
				}
				input.Due = parsed
			}
			if clearDue {
				input.ClearDue = true
			}
			if cmd.Flags().Changed("priority") {
				parsedPriority, err := parsePriorityFlag(priority)
				if err != nil {
					return err
				}
				if parsedPriority == nil {
					return fmt.Errorf("--priority must not be empty")
				}
				input.Priority = parsedPriority
			}
			if cmd.Flags().Changed("notes") {
				input.Notes = &notes
			}
			if cmd.Flags().Changed("url") {
				input.URL = &url
			}

			if input.Title == nil && input.List == nil && input.Due == nil && !input.ClearDue && input.Priority == nil && input.Notes == nil && input.URL == nil {
				return fmt.Errorf("no update fields supplied")
			}

			item, err := remindersService.Edit(strings.TrimSpace(args[0]), input)
			if err != nil {
				return err
			}
			return writeJSON(item)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New title")
	cmd.Flags().StringVar(&listName, "list", "", "Move reminder to list")
	cmd.Flags().StringVar(&dueIn, "due", "", "Due date/time (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().BoolVar(&clearDue, "clear-due", false, "Clear due date")
	cmd.Flags().StringVar(&priority, "priority", "", "Priority: none|high|medium|low or 0|1|5|9")
	cmd.Flags().StringVar(&notes, "notes", "", "New notes")
	cmd.Flags().StringVar(&url, "url", "", "New URL")
	return cmd
}
