package service

import (
	"fmt"
	"strings"

	ekrem "github.com/BRO3886/go-eventkit/reminders"
	"github.com/Robben-Media/apple-reminders-cli/internal/model"
)

type listsClient interface {
	Lists() ([]ekrem.List, error)
	CreateList(input ekrem.CreateListInput) (*ekrem.List, error)
	DeleteList(id string) error
}

// ListsService wraps go-eventkit list APIs with internal app models.
type ListsService struct {
	client listsClient
}

func NewListsService(client listsClient) *ListsService {
	return &ListsService{client: client}
}

func (s *ListsService) List() ([]model.List, error) {
	lists, err := s.client.Lists()
	if err != nil {
		return nil, err
	}

	out := make([]model.List, 0, len(lists))
	for _, item := range lists {
		out = append(out, toModelList(item))
	}
	return out, nil
}

func (s *ListsService) Create(name, source string) (*model.List, error) {
	title := strings.TrimSpace(name)
	if title == "" {
		return nil, fmt.Errorf("list name is required")
	}

	source = strings.TrimSpace(source)
	if source == "" {
		lists, err := s.client.Lists()
		if err != nil {
			return nil, err
		}
		source = firstSource(lists)
		if source == "" {
			return nil, fmt.Errorf("list source is required; pass --source")
		}
	}

	item, err := s.client.CreateList(ekrem.CreateListInput{
		Title:  title,
		Source: source,
	})
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("empty create response")
	}

	modelList := toModelList(*item)
	return &modelList, nil
}

func (s *ListsService) Delete(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return s.client.DeleteList(id)
}

func firstSource(lists []ekrem.List) string {
	for _, item := range lists {
		if src := strings.TrimSpace(item.Source); src != "" {
			return src
		}
	}
	return ""
}

func toModelList(item ekrem.List) model.List {
	return model.List{
		ID:       item.ID,
		Title:    item.Title,
		Color:    item.Color,
		Source:   item.Source,
		Count:    item.Count,
		ReadOnly: item.ReadOnly,
	}
}
