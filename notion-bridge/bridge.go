package main

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// ServerSource is implemented by the manager: it returns the currently usable upstreams.
type ServerSource interface {
	ActiveUpstreams() []ActiveServer
	StatusText() string
}

type ActiveServer struct {
	ID, Name string
	Client   Upstream
}

type route struct{ serverID, tool string }

type Bridge struct {
	src      ServerSource
	token    func() string
	log      func(string)
	mu       sync.Mutex
	routes   map[string]route
	CallTimeout time.Duration
}

func NewBridge(src ServerSource, token func() string, log func(string)) *Bridge {
	return &Bridge{src: src, token: token, log: log, routes: map[string]route{}, CallTimeout: 10 * time.Minute}
}

func exposedName(serverID, tool string) string {
	n := serverID + "_" + tool
	if len(n) <= 64 { return n }
	h := sha256.Sum256([]byte(n))
	return n[:55] + "_" + hex.EncodeToString(h[:])[:8]
}

func (b *Bridge) authorized(r *http.Request) bool {
	want := b.token()
	if want == "" { return false }
	got := ""
	if a := r.Header.Get("Authorization"); len(a) > 7 && strings.EqualFold(a[:7], "bearer ") {
		got = strings.TrimSpace(a[7:])
	} else if strings.HasPrefix(r.URL.Path, "/k/") {
		rest := strings.TrimPrefix(r.URL.Path, "/k/")
		if i := strings.IndexByte(rest, '/'); i > 0 { got = rest[:i] }
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}

func (b *Bridge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/health" {
		w.Write([]byte("ok"))
		return
	}
	isMCP := path == "/mcp" || path == "/mcp/" || (strings.HasPrefix(path, "/k/") && strings.HasSuffix(strings.TrimSuffix(path, "/"), "/mcp"))
	if !isMCP { http.NotFound(w, r); return }
	if !b.authorized(r) {
		w.Header().Set("WWW-Authenticate", `Bearer realm="PlazCode Notion Bridge"`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodPost:
	case http.MethodDelete:
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.Header().Set("Allow", "POST, DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 32<<20))
	if err != nil { http.Error(w, "bad request", 400); return }
	trimmed := strings.TrimSpace(string(body))
	var msgs []rpcMessage
	batch := strings.HasPrefix(trimmed, "[")
	if batch {
		err = json.Unmarshal(body, &msgs)
	} else {
		var m rpcMessage
		err = json.Unmarshal(body, &m)
		msgs = []rpcMessage{m}
	}
	if err != nil {
		writeJSON(w, map[string]any{"jsonrpc": "2.0", "id": nil, "error": map[string]any{"code": -32700, "message": "Parse error"}})
		return
	}
	var out []any
	for _, m := range msgs {
		hasID := len(m.ID) > 0 && string(m.ID) != "null"
		if m.Method == "initialize" {
			w.Header().Set("Mcp-Session-Id", randomToken()[:32])
		}
		res, rerr := b.handle(r.Context(), m)
		if !hasID { continue }
		if rerr != nil {
			out = append(out, map[string]any{"jsonrpc": "2.0", "id": m.ID, "error": rerr})
		} else {
			out = append(out, map[string]any{"jsonrpc": "2.0", "id": m.ID, "result": res})
		}
	}
	if len(out) == 0 { w.WriteHeader(http.StatusAccepted); return }
	if batch { writeJSON(w, out) } else { writeJSON(w, out[0]) }
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (b *Bridge) handle(ctx context.Context, m rpcMessage) (any, *RPCError) {
	switch m.Method {
	case "initialize":
		var p struct{ ProtocolVersion string `json:"protocolVersion"` }
		json.Unmarshal(m.Params, &p)
		v := LatestProtocol
		for _, s := range SupportedProtocols { if s == p.ProtocolVersion { v = s } }
		return map[string]any{
			"protocolVersion": v,
			"capabilities":    map[string]any{"tools": map[string]any{"listChanged": false}},
			"serverInfo":      map[string]any{"name": "PlazCode Notion Bridge", "version": AppVersion},
			"instructions":    "Bridge zu lokalen MCP-Servern dieses PCs. Tool-Namen beginnen mit der Server-ID (z. B. pc_, roblox_). bridge_status zeigt, welche Server aktiv sind.",
		}, nil
	case "ping":
		return map[string]any{}, nil
	case "tools/list":
		return map[string]any{"tools": b.listTools(ctx)}, nil
	case "tools/call":
		return b.callTool(ctx, m.Params)
	case "resources/list":
		return map[string]any{"resources": []any{}}, nil
	case "resources/templates/list":
		return map[string]any{"resourceTemplates": []any{}}, nil
	case "prompts/list":
		return map[string]any{"prompts": []any{}}, nil
	}
	if strings.HasPrefix(m.Method, "notifications/") { return nil, nil }
	return nil, &RPCError{Code: -32601, Message: "Methode nicht gefunden: " + m.Method}
}

func statusTool() map[string]any {
	return map[string]any{
		"name":        "bridge_status",
		"description": "Zeigt, welche MCP-Server hinter der PlazCode Notion Bridge laufen und erreichbar sind.",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{}},
		"annotations": map[string]any{"readOnlyHint": true},
	}
}

func (b *Bridge) listTools(ctx context.Context) []any {
	servers := b.src.ActiveUpstreams()
	type res struct { idx int; tools []map[string]any; err error }
	results := make([]res, len(servers))
	var wg sync.WaitGroup
	for i, s := range servers {
		wg.Add(1)
		go func(i int, s ActiveServer) {
			defer wg.Done()
			c, cancel := context.WithTimeout(ctx, 45*time.Second)
			defer cancel()
			tools, err := listAllTools(c, s.Client)
			results[i] = res{i, tools, err}
		}(i, s)
	}
	wg.Wait()
	all := []any{statusTool()}
	routes := map[string]route{}
	for i, r := range results {
		s := servers[i]
		if r.err != nil { b.log(fmt.Sprintf("tools/list %s: %v", s.Name, r.err)); continue }
		sort.SliceStable(r.tools, func(a, c int) bool { return fmt.Sprint(r.tools[a]["name"]) < fmt.Sprint(r.tools[c]["name"]) })
		for _, t := range r.tools {
			orig, _ := t["name"].(string)
			if orig == "" { continue }
			name := exposedName(s.ID, orig)
			routes[name] = route{s.ID, orig}
			t["name"] = name
			if d, ok := t["description"].(string); ok {
				t["description"] = "[" + s.Name + "] " + d
			} else {
				t["description"] = "[" + s.Name + "]"
			}
			if title, ok := t["title"].(string); ok && title != "" { t["title"] = s.Name + ": " + title }
			all = append(all, t)
		}
	}
	b.mu.Lock()
	for k, v := range routes { b.routes[k] = v }
	b.mu.Unlock()
	return all
}

func listAllTools(ctx context.Context, c Upstream) ([]map[string]any, error) {
	var out []map[string]any
	cursor := ""
	for page := 0; page < 50; page++ {
		var params any
		if cursor != "" { params = map[string]any{"cursor": cursor} }
		raw, err := c.Request(ctx, "tools/list", params)
		if err != nil { return nil, err }
		var r struct {
			Tools      []map[string]any `json:"tools"`
			NextCursor string           `json:"nextCursor"`
		}
		if err := json.Unmarshal(raw, &r); err != nil { return nil, err }
		out = append(out, r.Tools...)
		if r.NextCursor == "" { break }
		cursor = r.NextCursor
	}
	return out, nil
}

func textResult(text string, isErr bool) map[string]any {
	return map[string]any{"content": []any{map[string]any{"type": "text", "text": text}}, "isError": isErr}
}

func (b *Bridge) resolve(name string) (route, bool) {
	b.mu.Lock()
	r, ok := b.routes[name]
	b.mu.Unlock()
	if ok { return r, true }
	if i := strings.IndexByte(name, '_'); i > 0 {
		return route{name[:i], name[i+1:]}, true
	}
	return route{}, false
}

func (b *Bridge) callTool(ctx context.Context, raw json.RawMessage) (any, *RPCError) {
	var p map[string]any
	if err := json.Unmarshal(raw, &p); err != nil { return nil, &RPCError{Code: -32602, Message: "Ungültige Parameter"} }
	name, _ := p["name"].(string)
	if name == "bridge_status" { return textResult(b.src.StatusText(), false), nil }
	rt, ok := b.resolve(name)
	if !ok { return nil, &RPCError{Code: -32602, Message: "Unbekanntes Tool: " + name} }
	var target *ActiveServer
	for _, s := range b.src.ActiveUpstreams() {
		if s.ID == rt.serverID { s := s; target = &s; break }
	}
	if target == nil {
		return textResult(fmt.Sprintf("Der MCP-Server %q läuft gerade nicht. Starte ihn in PlazCode Notion und versuche es dann erneut.", rt.serverID), true), nil
	}
	params := map[string]any{"name": rt.tool}
	if a, ok := p["arguments"]; ok { params["arguments"] = a }
	if meta, ok := p["_meta"]; ok { params["_meta"] = meta }
	c, cancel := context.WithTimeout(ctx, b.CallTimeout)
	defer cancel()
	start := time.Now()
	res, err := target.Client.Request(c, "tools/call", params)
	b.log(fmt.Sprintf("Tool %s (%s) – %.1fs", rt.tool, target.Name, time.Since(start).Seconds()))
	if err != nil {
		if re, ok := err.(*RPCError); ok { return nil, re }
		return textResult(fmt.Sprintf("Verbindung zu %s fehlgeschlagen: %v. Prüfe vor einer Wiederholung den aktuellen Zustand, da der Befehl eventuell bereits ausgeführt wurde.", target.Name, err), true), nil
	}
	return json.RawMessage(res), nil
}
