package commands

import (
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Stream JSON lines on reminders changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			remindersService, _, err := ensureServices()
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()

			changes, err := remindersService.Watch(ctx)
			if err != nil {
				return err
			}

			for range changes {
				if err := writeJSON(map[string]any{"event": "change"}); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
