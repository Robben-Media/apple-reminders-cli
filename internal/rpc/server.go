package rpc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/Robben-Media/apple-reminders-cli/internal/service"
)

// Server implements a JSON-RPC 2.0 line-delimited stdin/stdout server.
type Server struct {
	in      io.Reader
	out     io.Writer
	handler *Handler

	writeMu sync.Mutex

	watchMu     sync.Mutex
	watchCancel context.CancelFunc
	watching    bool
	watchID     uint64
	baseCtx     context.Context
	reminders   *service.RemindersService
}

func NewServer(reminders *service.RemindersService, lists *service.ListsService, in io.Reader, out io.Writer) *Server {
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}

	s := &Server{
		in:        in,
		out:       out,
		reminders: reminders,
	}
	s.handler = &Handler{
		Reminders:        reminders,
		Lists:            lists,
		SubscribeWatch:   s.subscribeWatch,
		UnsubscribeWatch: s.unsubscribeWatch,
	}
	return s
}

func (s *Server) Serve(ctx context.Context) error {
	s.baseCtx = ctx
	defer func() { _ = s.unsubscribeWatch() }()

	scanner := bufio.NewScanner(s.in)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		response := s.handleLine(ctx, []byte(line))
		if response != nil {
			if err := s.writeJSON(response); err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleLine(ctx context.Context, line []byte) *Response {
	decoder := json.NewDecoder(bytes.NewReader(line))
	decoder.DisallowUnknownFields()

	var req Request
	if err := decoder.Decode(&req); err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      json.RawMessage("null"),
			Error:   newError(ErrCodeParse, "parse error"),
		}
	}
	if req.JSONRPC != "2.0" || strings.TrimSpace(req.Method) == "" {
		id := json.RawMessage("null")
		if len(req.ID) > 0 {
			id = req.ID
		}
		return &Response{
			JSONRPC: "2.0",
			ID:      id,
			Error:   newError(ErrCodeInvalidRequest, "invalid request"),
		}
	}

	result, rpcErr := s.handler.Handle(ctx, req)
	if len(req.ID) == 0 {
		return nil
	}

	response := &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
	}
	if rpcErr != nil {
		response.Error = rpcErr
		return response
	}
	response.Result = result
	return response
}

func (s *Server) subscribeWatch(_ context.Context) error {
	s.watchMu.Lock()
	if s.watching {
		s.watchMu.Unlock()
		return nil
	}
	watchCtx, cancel := context.WithCancel(s.baseCtx)
	changes, err := s.reminders.Watch(watchCtx)
	if err != nil {
		cancel()
		s.watchMu.Unlock()
		return err
	}
	s.watchID++
	currentWatchID := s.watchID
	s.watching = true
	s.watchCancel = cancel
	s.watchMu.Unlock()

	go func(expectedWatchID uint64) {
		for range changes {
			notification := Notification{
				JSONRPC: "2.0",
				Method:  "watch.event",
				Params:  map[string]any{},
			}
			if err := s.writeJSON(notification); err != nil {
				break
			}
		}

		s.watchMu.Lock()
		if s.watchID == expectedWatchID {
			s.watchCancel = nil
			s.watching = false
		}
		s.watchMu.Unlock()
	}(currentWatchID)

	return nil
}

func (s *Server) unsubscribeWatch() error {
	s.watchMu.Lock()
	cancel := s.watchCancel
	s.watchCancel = nil
	s.watching = false
	s.watchMu.Unlock()
	if cancel != nil {
		cancel()
	}
	return nil
}

func (s *Server) writeJSON(v any) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	encoder := json.NewEncoder(s.out)
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("rpc write failed: %w", err)
	}
	return nil
}
