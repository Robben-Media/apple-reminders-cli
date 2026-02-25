package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	ekrem "github.com/BRO3886/go-eventkit/reminders"
	"github.com/Robben-Media/apple-reminders-cli/internal/model"
)

type remindersClient interface {
	Reminders(opts ...ekrem.ListOption) ([]ekrem.Reminder, error)
	Reminder(id string) (*ekrem.Reminder, error)
	CreateReminder(input ekrem.CreateReminderInput) (*ekrem.Reminder, error)
	UpdateReminder(id string, input ekrem.UpdateReminderInput) (*ekrem.Reminder, error)
	DeleteReminder(id string) error
	CompleteReminder(id string) (*ekrem.Reminder, error)
	UncompleteReminder(id string) (*ekrem.Reminder, error)
	WatchChanges(ctx context.Context) (<-chan struct{}, error)
}

// RemindersService wraps go-eventkit with internal app models.
type RemindersService struct {
	client remindersClient
}

func NewRemindersService(client remindersClient) *RemindersService {
	return &RemindersService{client: client}
}

// ReminderQuery is the list filter model used by CLI and RPC layers.
type ReminderQuery struct {
	List      string
	ListID    string
	Completed *bool
	Search    string
	DueBefore *time.Time
	DueAfter  *time.Time
}

// AddReminderInput is the create reminder input model used by CLI and RPC layers.
type AddReminderInput struct {
	Title    string
	List     string
	Due      *time.Time
	Priority *model.Priority
	Notes    *string
	URL      *string
}

// EditReminderInput is the update reminder input model used by CLI and RPC layers.
type EditReminderInput struct {
	Title    *string
	Due      *time.Time
	ClearDue bool
	Priority *model.Priority
	Notes    *string
	URL      *string
	List     *string
}

func (s *RemindersService) List(query ReminderQuery) ([]model.Reminder, error) {
	opts := make([]ekrem.ListOption, 0, 6)
	if query.List != "" {
		opts = append(opts, ekrem.WithList(query.List))
	}
	if query.ListID != "" {
		opts = append(opts, ekrem.WithListID(query.ListID))
	}
	if query.Completed != nil {
		opts = append(opts, ekrem.WithCompleted(*query.Completed))
	}
	if query.Search != "" {
		opts = append(opts, ekrem.WithSearch(query.Search))
	}
	if query.DueBefore != nil {
		opts = append(opts, ekrem.WithDueBefore(*query.DueBefore))
	}
	if query.DueAfter != nil {
		opts = append(opts, ekrem.WithDueAfter(*query.DueAfter))
	}

	items, err := s.client.Reminders(opts...)
	if err != nil {
		return nil, err
	}

	out := make([]model.Reminder, 0, len(items))
	for _, item := range items {
		out = append(out, toModelReminder(item))
	}
	return out, nil
}

func (s *RemindersService) Show(id string) (*model.Reminder, error) {
	item, err := s.client.Reminder(strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("reminder not found")
	}
	modelReminder := toModelReminder(*item)
	return &modelReminder, nil
}

func (s *RemindersService) Add(input AddReminderInput) (*model.Reminder, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	req := ekrem.CreateReminderInput{
		Title:    title,
		ListName: strings.TrimSpace(input.List),
		DueDate:  input.Due,
	}
	if input.Priority != nil {
		req.Priority = toEKPriority(*input.Priority)
	}
	if input.Notes != nil {
		req.Notes = *input.Notes
	}
	if input.URL != nil {
		req.URL = *input.URL
	}

	item, err := s.client.CreateReminder(req)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("empty create response")
	}
	modelReminder := toModelReminder(*item)
	return &modelReminder, nil
}

func (s *RemindersService) Edit(id string, input EditReminderInput) (*model.Reminder, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	req := ekrem.UpdateReminderInput{
		Title:        input.Title,
		Notes:        input.Notes,
		ListName:     input.List,
		DueDate:      input.Due,
		ClearDueDate: input.ClearDue,
		URL:          input.URL,
	}
	if input.Priority != nil {
		p := toEKPriority(*input.Priority)
		req.Priority = &p
	}

	item, err := s.client.UpdateReminder(id, req)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("empty update response")
	}
	modelReminder := toModelReminder(*item)
	return &modelReminder, nil
}

func (s *RemindersService) Complete(id string) (*model.Reminder, error) {
	item, err := s.client.CompleteReminder(strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("empty complete response")
	}
	modelReminder := toModelReminder(*item)
	return &modelReminder, nil
}

func (s *RemindersService) Uncomplete(id string) (*model.Reminder, error) {
	item, err := s.client.UncompleteReminder(strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("empty uncomplete response")
	}
	modelReminder := toModelReminder(*item)
	return &modelReminder, nil
}

func (s *RemindersService) Delete(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return s.client.DeleteReminder(id)
}

func (s *RemindersService) Watch(ctx context.Context) (<-chan struct{}, error) {
	return s.client.WatchChanges(ctx)
}

func ParsePriority(value string) (model.Priority, error) {
	s := strings.ToLower(strings.TrimSpace(value))
	switch s {
	case "", "none", "0":
		return model.PriorityNone, nil
	case "high", "h", "1":
		return model.PriorityHigh, nil
	case "medium", "med", "m", "5":
		return model.PriorityMedium, nil
	case "low", "l", "9":
		return model.PriorityLow, nil
	default:
		return 0, fmt.Errorf("invalid priority: %s", value)
	}
}

func ParsePriorityInt(value int) (model.Priority, error) {
	switch value {
	case int(model.PriorityNone):
		return model.PriorityNone, nil
	case int(model.PriorityHigh):
		return model.PriorityHigh, nil
	case int(model.PriorityMedium):
		return model.PriorityMedium, nil
	case int(model.PriorityLow):
		return model.PriorityLow, nil
	default:
		return 0, fmt.Errorf("invalid priority: %d", value)
	}
}

func toEKPriority(p model.Priority) ekrem.Priority {
	switch p {
	case model.PriorityHigh:
		return ekrem.PriorityHigh
	case model.PriorityMedium:
		return ekrem.PriorityMedium
	case model.PriorityLow:
		return ekrem.PriorityLow
	default:
		return ekrem.PriorityNone
	}
}

func toModelReminder(item ekrem.Reminder) model.Reminder {
	return model.Reminder{
		ID:             item.ID,
		Title:          item.Title,
		Notes:          item.Notes,
		List:           item.List,
		ListID:         item.ListID,
		DueDate:        item.DueDate,
		CompletionDate: item.CompletionDate,
		CreatedAt:      item.CreatedAt,
		ModifiedAt:     item.ModifiedAt,
		Priority:       model.Priority(item.Priority),
		Completed:      item.Completed,
		URL:            item.URL,
	}
}
