package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Robben-Media/apple-reminders-cli/internal/model"
	"github.com/Robben-Media/apple-reminders-cli/internal/service"
)

func writeJSON(v any) error {
	encoder := json.NewEncoder(os.Stdout)
	return encoder.Encode(v)
}

func WriteErrorJSON(msg string) {
	_ = json.NewEncoder(os.Stderr).Encode(map[string]string{"error": msg})
}

func parseDateTimeFlag(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return &parsed, nil
	}
	if parsed, err := time.ParseInLocation("2006-01-02", value, time.Local); err == nil {
		return &parsed, nil
	}
	return nil, fmt.Errorf("invalid date/time %q (expected RFC3339 or YYYY-MM-DD)", value)
}

func parseCompletedFilter(value string) (*bool, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return nil, nil
	}
	switch value {
	case "true", "1", "yes", "y":
		v := true
		return &v, nil
	case "false", "0", "no", "n":
		v := false
		return &v, nil
	default:
		return nil, fmt.Errorf("invalid --completed value: %s", value)
	}
}

func parsePriorityFlag(value string) (*model.Priority, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	priority, err := service.ParsePriority(value)
	if err != nil {
		return nil, err
	}
	return &priority, nil
}
