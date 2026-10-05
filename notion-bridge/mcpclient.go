package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const LatestProtocol = "2025-06-18"

var SupportedProtocols = []string{"2025-06-18", "2025-03-26", "2024-11-05"}

type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *RPCError) Error() string { return fmt.Sprintf("MCP-Fehler %d: %s", e.Code, e.Message) }

type rpcMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// Upstream is an MCP client connection to one server behind the bridge.
type Upstream interface {
	Request(ctx context.Context, method string, params any) (json.RawMessage, error)
	Close()
}

var errSessionExpired = errors.New("MCP-Sitzung abgelaufen")

func initializeUpstream(ctx context.Context, send func(ctx context.Context, method string, params any, notify bool) (json.RawMessage, error)) (string, error) {
	res, err := send(ctx, "initialize", map[string]any{
		"protocolVersion": LatestProtocol,
		"capabilities":    map[string]any{},
		"clientInfo":      map[string]any{"name": "plazcode-notion-bridge", "version": AppVersion},
	}, false)
	if err != nil { return "", err }
	var r struct{ ProtocolVersion string `json:"protocolVersion"` }
	_ = json.Unmarshal(res, &r)
	if _, err := send(ctx, "notifications/initialized", nil, true); err != nil { return "", err }
	if r.ProtocolVersion == "" { r.ProtocolVersion = LatestProtocol }
	return r.ProtocolVersion, nil
}

// ---------- Streamable HTTP ----------

type HTTPUpstream struct {
	URL, Token string
	hc         *http.Client
	mu         sync.Mutex
	smu        sync.Mutex
	session    string
	proto      string
	ready      bool
	nextID     atomic.Int64
}

func NewHTTPUpstream(url, token string) *HTTPUpstream {
	return &HTTPUpstream{URL: url, Token: token, hc: &http.Client{Timeout: 0}}
}

func (h *HTTPUpstream) send(ctx context.Context, method string, params any, notify bool) (json.RawMessage, error) {
	msg := map[string]any{"jsonrpc": "2.0", "method": method}
	if params != nil { msg["params"] = params }
	var id int64
	if !notify { id = h.nextID.Add(1); msg["id"] = id }
	body, _ := json.Marshal(msg)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.URL, bytes.NewReader(body))
	if err != nil { return nil, err }
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if h.Token != "" { req.Header.Set("Authorization", "Bearer "+h.Token) }
	h.smu.Lock()
	session, proto := h.session, h.proto
	h.smu.Unlock()
	if session != "" { req.Header.Set("Mcp-Session-Id", session) }
	if proto != "" { req.Header.Set("MCP-Protocol-Version", proto) }
	resp, err := h.hc.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if sid := resp.Header.Get("Mcp-Session-Id"); sid != "" && method == "initialize" {
		h.smu.Lock(); h.session = sid; h.smu.Unlock()
	}
	if resp.StatusCode == http.StatusNotFound && session != "" && method != "initialize" {
		return nil, errSessionExpired
	}
	if notify {
		io.Copy(io.Discard, resp.Body)
		if resp.StatusCode >= 300 { return nil, fmt.Errorf("HTTP %d", resp.StatusCode) }
		return nil, nil
	}
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	want := fmt.Sprint(id)
	ct := resp.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "text/event-stream") {
		return readSSEResponse(resp.Body, want)
	}
	var m rpcMessage
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil { return nil, fmt.Errorf("ungültige Antwort: %w", err) }
	if m.Error != nil { return nil, m.Error }
	return m.Result, nil
}

func readSSEResponse(r io.Reader, wantID string) (json.RawMessage, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 64*1024*1024)
	var data strings.Builder
	flush := func() (json.RawMessage, error, bool) {
		if data.Len() == 0 { return nil, nil, false }
		raw := data.String()
		data.Reset()
		var m rpcMessage
		if json.Unmarshal([]byte(raw), &m) != nil { return nil, nil, false }
		if m.Method != "" || strings.TrimSpace(string(m.ID)) != wantID { return nil, nil, false }
		if m.Error != nil { return nil, m.Error, true }
		return m.Result, nil, true
	}
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			if res, err, ok := flush(); ok { return res, err }
			continue
		}
		if strings.HasPrefix(line, "data:") {
			if data.Len() > 0 { data.WriteByte('\n') }
			data.WriteString(strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
		}
	}
	if res, err, ok := flush(); ok { return res, err }
	if err := sc.Err(); err != nil { return nil, err }
	return nil, errors.New("Stream endete ohne Antwort")
}

func (h *HTTPUpstream) ensure(ctx context.Context) error {
	if h.ready { return nil }
	h.smu.Lock(); h.session, h.proto = "", ""; h.smu.Unlock()
	ictx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	p, err := initializeUpstream(ictx, h.send)
	if err != nil { return fmt.Errorf("Initialisierung fehlgeschlagen: %w", err) }
	h.smu.Lock(); h.proto = p; h.smu.Unlock()
	h.ready = true
	return nil
}

func (h *HTTPUpstream) Request(ctx context.Context, method string, params any) (json.RawMessage, error) {
	h.mu.Lock()
	err := h.ensure(ctx)
	h.mu.Unlock()
	if err != nil { return nil, err }
	res, err := h.send(ctx, method, params, false)
	if errors.Is(err, errSessionExpired) {
		// The server never processed the request, so one retry with a new session is safe.
		h.mu.Lock()
		h.ready = false
		err = h.ensure(ctx)
		h.mu.Unlock()
		if err != nil { return nil, err }
		res, err = h.send(ctx, method, params, false)
	}
	return res, err
}

func (h *HTTPUpstream) Close() {}

// ---------- stdio ----------

type StdioUpstream struct {
	w       io.WriteCloser
	wmu     sync.Mutex
	pmu     sync.Mutex
	pending map[string]chan rpcMessage
	nextID  atomic.Int64
	closed  chan struct{}
	once    sync.Once
	initMu  sync.Mutex
	ready   bool
	onClose func()
}

func NewStdioUpstream(stdin io.WriteCloser, stdout io.Reader, onClose func()) *StdioUpstream {
	s := &StdioUpstream{w: stdin, pending: map[string]chan rpcMessage{}, closed: make(chan struct{}), onClose: onClose}
	go s.readLoop(stdout)
	return s
}

func (s *StdioUpstream) write(v any) error {
	b, _ := json.Marshal(v)
	s.wmu.Lock()
	defer s.wmu.Unlock()
	_, err := s.w.Write(append(b, '\n'))
	return err
}

func (s *StdioUpstream) readLoop(r io.Reader) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 64*1024*1024)
	for sc.Scan() {
		var m rpcMessage
		if json.Unmarshal(sc.Bytes(), &m) != nil { continue }
		hasID := len(m.ID) > 0 && string(m.ID) != "null"
		switch {
		case m.Method != "" && hasID:
			if m.Method == "ping" {
				s.write(map[string]any{"jsonrpc": "2.0", "id": m.ID, "result": map[string]any{}})
			} else {
				s.write(map[string]any{"jsonrpc": "2.0", "id": m.ID, "error": map[string]any{"code": -32601, "message": "Von der Bridge nicht unterstützt"}})
			}
		case m.Method == "" && hasID:
			key := strings.TrimSpace(string(m.ID))
			s.pmu.Lock()
			ch := s.pending[key]
			delete(s.pending, key)
			s.pmu.Unlock()
			if ch != nil { ch <- m }
		}
	}
	s.Close()
}

func (s *StdioUpstream) send(ctx context.Context, method string, params any, notify bool) (json.RawMessage, error) {
	msg := map[string]any{"jsonrpc": "2.0", "method": method}
	if params != nil { msg["params"] = params }
	if notify { return nil, s.write(msg) }
	id := s.nextID.Add(1)
	msg["id"] = id
	ch := make(chan rpcMessage, 1)
	key := fmt.Sprint(id)
	s.pmu.Lock()
	s.pending[key] = ch
	s.pmu.Unlock()
	defer func() { s.pmu.Lock(); delete(s.pending, key); s.pmu.Unlock() }()
	if err := s.write(msg); err != nil { return nil, err }
	select {
	case m := <-ch:
		if m.Error != nil { return nil, m.Error }
		return m.Result, nil
	case <-s.closed:
		return nil, errors.New("Serverprozess wurde beendet")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *StdioUpstream) Request(ctx context.Context, method string, params any) (json.RawMessage, error) {
	s.initMu.Lock()
	if !s.ready {
		ictx, cancel := context.WithTimeout(ctx, 60*time.Second)
		_, err := initializeUpstream(ictx, s.send)
		cancel()
		if err != nil { s.initMu.Unlock(); return nil, fmt.Errorf("Initialisierung fehlgeschlagen: %w", err) }
		s.ready = true
	}
	s.initMu.Unlock()
	return s.send(ctx, method, params, false)
}

func (s *StdioUpstream) Close() {
	s.once.Do(func() {
		close(s.closed)
		s.w.Close()
		if s.onClose != nil { s.onClose() }
	})
}
