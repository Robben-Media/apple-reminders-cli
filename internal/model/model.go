package model

import "time"

// Priority mirrors Apple's EventKit reminder priority values.
type Priority int

const (
	PriorityNone   Priority = 0
	PriorityHigh   Priority = 1
	PriorityMedium Priority = 5
	PriorityLow    Priority = 9
)

// Reminder is the internal reminder model used by CLI and RPC layers.
type Reminder struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Notes          string     `json:"notes,omitempty"`
	List           string     `json:"list"`
	ListID         string     `json:"listID"`
	DueDate        *time.Time `json:"dueDate,omitempty"`
	CompletionDate *time.Time `json:"completionDate,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	ModifiedAt     *time.Time `json:"modifiedAt,omitempty"`
	Priority       Priority   `json:"priority"`
	Completed      bool       `json:"completed"`
	URL            string     `json:"url,omitempty"`
}

// List is the internal reminder list model used by CLI and RPC layers.
type List struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Color    string `json:"color,omitempty"`
	Source   string `json:"source,omitempty"`
	Count    int    `json:"count"`
	ReadOnly bool   `json:"readOnly"`
}
