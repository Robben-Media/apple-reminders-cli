package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Robben-Media/apple-reminders-cli/internal/model"
	"github.com/Robben-Media/apple-reminders-cli/internal/service"
)

// Handler dispatches JSON-RPC methods to internal services.
type Handler struct {
	Reminders *service.RemindersService
	Lists     *service.ListsService

	SubscribeWatch   func(context.Context) error
	UnsubscribeWatch func() error
}

func (h *Handler) Handle(ctx context.Context, req Request) (any, *RPCError) {
	switch req.Method {
	case "lists.list":
		items, err := h.Lists.List()
		if err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return items, nil
	case "lists.create":
		var params struct {
			Name   string `json:"name"`
			Source string `json:"source"`
		}
		if rpcErr := decodeParams(req.Params, &params); rpcErr != nil {
			return nil, rpcErr
		}
		item, err := h.Lists.Create(params.Name, params.Source)
		if err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return item, nil
	case "lists.delete":
		var params struct {
			ID string `json:"id"`
		}
		if rpcErr := decodeParams(req.Params, &params); rpcErr != nil {
			return nil, rpcErr
		}
		if strings.TrimSpace(params.ID) == "" {
			return nil, newError(ErrCodeInvalidParams, "id is required")
		}
		if err := h.Lists.Delete(params.ID); err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return map[string]any{"deleted": true, "id": params.ID}, nil
	case "reminders.list":
		var params struct {
			List           string `json:"list"`
			ListID         string `json:"list_id"`
			ListIDAlt      string `json:"listID"`
			Completed      *bool  `json:"completed"`
			Search         string `json:"search"`
			DueBefore      string `json:"due_before"`
			DueBeforeCamel string `json:"dueBefore"`
			DueAfter       string `json:"due_after"`
			DueAfterCamel  string `json:"dueAfter"`
		}
		if rpcErr := decodeParams(req.Params, &params); rpcErr != nil {
			return nil, rpcErr
		}

		query := service.ReminderQuery{
			List:      strings.TrimSpace(params.List),
			Completed: params.Completed,
			Search:    strings.TrimSpace(params.Search),
		}
		if params.ListID != "" {
			query.ListID = strings.TrimSpace(params.ListID)
		} else {
			query.ListID = strings.TrimSpace(params.ListIDAlt)
		}

		if dueBefore := firstNonEmpty(params.DueBefore, params.DueBeforeCamel); dueBefore != "" {
			parsed, err := parseDateTime(dueBefore)
			if err != nil {
				return nil, newError(ErrCodeInvalidParams, "invalid due_before: "+err.Error())
			}
			query.DueBefore = parsed
		}
		if dueAfter := firstNonEmpty(params.DueAfter, params.DueAfterCamel); dueAfter != "" {
			parsed, err := parseDateTime(dueAfter)
			if err != nil {
				return nil, newError(ErrCodeInvalidParams, "invalid due_after: "+err.Error())
			}
			query.DueAfter = parsed
		}

		items, err := h.Reminders.List(query)
		if err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return items, nil
	case "reminders.show":
		var params struct {
			ID string `json:"id"`
		}
		if rpcErr := decodeParams(req.Params, &params); rpcErr != nil {
			return nil, rpcErr
		}
		if strings.TrimSpace(params.ID) == "" {
			return nil, newError(ErrCodeInvalidParams, "id is required")
		}
		item, err := h.Reminders.Show(params.ID)
		if err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return item, nil
	case "reminders.add":
		var params struct {
			Title    string          `json:"title"`
			List     string          `json:"list"`
			Due      string          `json:"due"`
			Priority json.RawMessage `json:"priority"`
			Notes    string          `json:"notes"`
			URL      string          `json:"url"`
		}
		if rpcErr := decodeParams(req.Params, &params); rpcErr != nil {
			return nil, rpcErr
		}

		input := service.AddReminderInput{
			Title: strings.TrimSpace(params.Title),
			List:  strings.TrimSpace(params.List),
		}
		if strings.TrimSpace(params.Notes) != "" {
			input.Notes = &params.Notes
		}
		if strings.TrimSpace(params.URL) != "" {
			input.URL = &params.URL
		}
		if strings.TrimSpace(params.Due) != "" {
			parsed, err := parseDateTime(params.Due)
			if err != nil {
				return nil, newError(ErrCodeInvalidParams, "invalid due: "+err.Error())
			}
			input.Due = parsed
		}
		if len(params.Priority) > 0 {
			parsedPriority, rpcErr := parsePriority(params.Priority)
			if rpcErr != nil {
				return nil, rpcErr
			}
			input.Priority = parsedPriority
		}

		item, err := h.Reminders.Add(input)
		if err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return item, nil
	case "reminders.edit":
		var params struct {
			ID       string           `json:"id"`
			Title    *string          `json:"title"`
			List     *string          `json:"list"`
			Due      *string          `json:"due"`
			Priority *json.RawMessage `json:"priority"`
			Notes    *string          `json:"notes"`
			URL      *string          `json:"url"`
		}
		if rpcErr := decodeParams(req.Params, &params); rpcErr != nil {
			return nil, rpcErr
		}
		if strings.TrimSpace(params.ID) == "" {
			return nil, newError(ErrCodeInvalidParams, "id is required")
		}

		input := service.EditReminderInput{
			Title: params.Title,
			List:  params.List,
			Notes: params.Notes,
			URL:   params.URL,
		}
		if params.Due != nil {
			due := strings.TrimSpace(*params.Due)
			if due == "" {
				input.ClearDue = true
			} else {
				parsed, err := parseDateTime(due)
				if err != nil {
					return nil, newError(ErrCodeInvalidParams, "invalid due: "+err.Error())
				}
				input.Due = parsed
			}
		}
		if params.Priority != nil && len(*params.Priority) > 0 {
			parsedPriority, rpcErr := parsePriority(*params.Priority)
			if rpcErr != nil {
				return nil, rpcErr
			}
			input.Priority = parsedPriority
		}

		item, err := h.Reminders.Edit(params.ID, input)
		if err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return item, nil
	case "reminders.complete":
		var params struct {
			ID string `json:"id"`
		}
		if rpcErr := decodeParams(req.Params, &params); rpcErr != nil {
			return nil, rpcErr
		}
		if strings.TrimSpace(params.ID) == "" {
			return nil, newError(ErrCodeInvalidParams, "id is required")
		}
		item, err := h.Reminders.Complete(params.ID)
		if err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return item, nil
	case "reminders.uncomplete":
		var params struct {
			ID string `json:"id"`
		}
		if rpcErr := decodeParams(req.Params, &params); rpcErr != nil {
			return nil, rpcErr
		}
		if strings.TrimSpace(params.ID) == "" {
			return nil, newError(ErrCodeInvalidParams, "id is required")
		}
		item, err := h.Reminders.Uncomplete(params.ID)
		if err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return item, nil
	case "reminders.delete":
		var params struct {
			ID string `json:"id"`
		}
		if rpcErr := decodeParams(req.Params, &params); rpcErr != nil {
			return nil, rpcErr
		}
		if strings.TrimSpace(params.ID) == "" {
			return nil, newError(ErrCodeInvalidParams, "id is required")
		}
		if err := h.Reminders.Delete(params.ID); err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return map[string]any{"deleted": true, "id": params.ID}, nil
	case "watch.subscribe":
		if h.SubscribeWatch == nil {
			return nil, newError(ErrCodeInternal, "watch not configured")
		}
		if err := h.SubscribeWatch(ctx); err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return map[string]any{"subscribed": true}, nil
	case "watch.unsubscribe":
		if h.UnsubscribeWatch == nil {
			return nil, newError(ErrCodeInternal, "watch not configured")
		}
		if err := h.UnsubscribeWatch(); err != nil {
			return nil, newError(ErrCodeInternal, err.Error())
		}
		return map[string]any{"subscribed": false}, nil
	default:
		return nil, newError(ErrCodeMethodNotFound, "method not found")
	}
}

func decodeParams(raw json.RawMessage, dst any) *RPCError {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return newError(ErrCodeInvalidParams, "invalid params")
	}
	return nil
}

func parseDateTime(input string) (*time.Time, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("value is empty")
	}

	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return &t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", input, time.Local); err == nil {
		return &t, nil
	}
	return nil, fmt.Errorf("expected RFC3339 or YYYY-MM-DD")
}

func parsePriority(raw json.RawMessage) (*model.Priority, *RPCError) {
	var asInt int
	if err := json.Unmarshal(raw, &asInt); err == nil {
		priority, parseErr := service.ParsePriorityInt(asInt)
		if parseErr != nil {
			return nil, newError(ErrCodeInvalidParams, parseErr.Error())
		}
		return &priority, nil
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		priority, parseErr := service.ParsePriority(asString)
		if parseErr != nil {
			return nil, newError(ErrCodeInvalidParams, parseErr.Error())
		}
		return &priority, nil
	}

	return nil, newError(ErrCodeInvalidParams, "priority must be string or number")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
