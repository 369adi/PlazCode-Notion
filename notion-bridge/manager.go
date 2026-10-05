package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type ServerRuntime struct {
	Cfg      ServerConfig
	State    string
	Err      string
	client   Upstream
	cmd      *exec.Cmd
	stopping bool
	gen      int
}

type Manager struct {
	mu          sync.Mutex
	cfg         *Config
	cfgPath     string
	home        string
	logDir      string
	logFn       func(string)
	logMu       sync.Mutex
	servers     []*ServerRuntime
	httpSrv     *http.Server
	bridgeState string
	ngrokCmd    *exec.Cmd
	ngrokState  string
	ngrokErr    string
	ngrokStop   bool
	bridge      *Bridge
}

func NewManager(cfgPath, home, logDir string, logFn func(string)) (*Manager, error) {
	m := &Manager{cfgPath: cfgPath, home: home, logDir: logDir, logFn: logFn, bridgeState: "gestoppt", ngrokState: "gestoppt"}
	os.MkdirAll(logDir, 0o700)
	if err := m.ReloadConfig(); err != nil { return nil, err }
	m.bridge = NewBridge(m, func() string { m.mu.Lock(); defer m.mu.Unlock(); return m.cfg.BridgeToken }, m.Log)
	return m, nil
}

func (m *Manager) Log(s string) {
	line := time.Now().Format("15:04:05") + "  " + s
	m.logMu.Lock()
	if f, err := os.OpenFile(filepath.Join(m.logDir, "bridge.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600); err == nil {
		f.WriteString(time.Now().Format("2006-01-02 ") + line + "\n")
		f.Close()
	}
	m.logMu.Unlock()
	if m.logFn != nil { m.logFn(line) }
}

func (m *Manager) Config() Config { m.mu.Lock(); defer m.mu.Unlock(); return *m.cfg }

func (m *Manager) ReloadConfig() error {
	cfg, err := LoadConfig(m.cfgPath, m.home)
	if err != nil { return err }
	m.mu.Lock()
	defer m.mu.Unlock()
	old := map[string]*ServerRuntime{}
	for _, s := range m.servers { old[s.Cfg.ID] = s }
	var list []*ServerRuntime
	for _, sc := range cfg.Servers {
		if rt, ok := old[sc.ID]; ok {
			rt.Cfg = sc
			list = append(list, rt)
			delete(old, sc.ID)
		} else {
			list = append(list, &ServerRuntime{Cfg: sc, State: "gestoppt"})
		}
	}
	for _, rt := range old { go m.stopRuntime(rt) }
	m.cfg, m.servers = cfg, list
	return nil
}

func (m *Manager) save() error { return SaveConfig(m.cfgPath, m.cfg) }

func (m *Manager) find(id string) *ServerRuntime {
	for _, s := range m.servers { if s.Cfg.ID == id { return s } }
	return nil
}

// ---------- ServerSource ----------

func (m *Manager) ActiveUpstreams() []ActiveServer {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []ActiveServer
	for _, s := range m.servers {
		if s.State == "läuft" && s.client != nil { out = append(out, ActiveServer{s.Cfg.ID, s.Cfg.Name, s.client}) }
	}
	return out
}

func (m *Manager) StatusText() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var b strings.Builder
	fmt.Fprintf(&b, "Bridge: %s   |   ngrok: %s", m.bridgeState, m.ngrokState)
	if m.ngrokErr != "" && m.ngrokState != "online" { fmt.Fprintf(&b, " (%s)", m.ngrokErr) }
	for _, s := range m.servers {
		st := s.State
		if !s.Cfg.Enabled { st = "deaktiviert" }
		fmt.Fprintf(&b, "\n%s [%s_*]: %s", s.Cfg.Name, s.Cfg.ID, st)
		if s.Err != "" && s.State != "läuft" { fmt.Fprintf(&b, " – %s", s.Err) }
	}
	return b.String()
}

// ---------- Bridge HTTP server ----------

func (m *Manager) StartBridge() error {
	m.mu.Lock()
	if m.httpSrv != nil { m.mu.Unlock(); return nil }
	addr := fmt.Sprintf("127.0.0.1:%d", m.cfg.BridgePort)
	m.mu.Unlock()
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		m.setBridge("Fehler")
		return fmt.Errorf("Bridge-Port %s belegt oder gesperrt: %w", addr, err)
	}
	srv := &http.Server{Handler: m.bridge, ReadHeaderTimeout: 15 * time.Second}
	m.mu.Lock()
	m.httpSrv = srv
	m.bridgeState = "läuft auf " + addr
	m.mu.Unlock()
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) { m.Log("Bridge-Fehler: " + err.Error()); m.setBridge("Fehler") }
	}()
	m.Log("Bridge läuft auf http://" + addr + "/mcp")
	return nil
}

func (m *Manager) setBridge(s string) { m.mu.Lock(); m.bridgeState = s; m.mu.Unlock() }

func (m *Manager) StopBridge() {
	m.mu.Lock()
	srv := m.httpSrv
	m.httpSrv = nil
	m.bridgeState = "gestoppt"
	m.mu.Unlock()
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		srv.Shutdown(ctx)
		cancel()
	}
}

// ---------- processes ----------

func (m *Manager) logFile(name string) (*os.File, error) {
	return os.OpenFile(filepath.Join(m.logDir, name+".log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
}

func findExe(name string) (string, error) {
	if p, err := exec.LookPath(name); err == nil { return p, nil }
	if runtime.GOOS != "windows" { return "", fmt.Errorf("%s nicht gefunden", name) }
	up, la, ad, pf := os.Getenv("USERPROFILE"), os.Getenv("LOCALAPPDATA"), os.Getenv("APPDATA"), os.Getenv("ProgramFiles")
	dirs := []string{
		filepath.Join(up, ".local", "bin"), filepath.Join(la, "Microsoft", "WinGet", "Links"),
		filepath.Join(pf, "nodejs"), filepath.Join(ad, "npm"), filepath.Join(pf, "Git", "cmd"),
		filepath.Join(la, "ngrok"), filepath.Join(up, "ngrok"), filepath.Join(up, "scoop", "shims"),
		filepath.Join(la, "Programs", "ngrok"), `C:\ProgramData\chocolatey\bin`,
	}
	for _, d := range dirs {
		for _, ext := range []string{".exe", ".cmd", ".bat"} {
			p := filepath.Join(d, name+ext)
			if st, err := os.Stat(p); err == nil && !st.IsDir() { return p, nil }
		}
	}
	return "", fmt.Errorf("%s wurde nicht gefunden. Bitte installieren bzw. zum PATH hinzufügen", name)
}

func findBrowser() string {
	pf, pf86, la := os.Getenv("ProgramFiles"), os.Getenv("ProgramFiles(x86)"), os.Getenv("LOCALAPPDATA")
	for _, p := range []string{
		filepath.Join(pf, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
		filepath.Join(pf86, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
		filepath.Join(la, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
		filepath.Join(pf86, "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(pf, "Google", "Chrome", "Application", "chrome.exe"),
	} {
		if st, err := os.Stat(p); err == nil && !st.IsDir() { return p }
	}
	return ""
}

func waitPort(ctx context.Context, port int, exited <-chan struct{}, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		c, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 500*time.Millisecond)
		if err == nil { c.Close(); return nil }
		select {
		case <-exited: return errors.New("Prozess wurde beim Start beendet – siehe Log")
		case <-ctx.Done(): return ctx.Err()
		case <-time.After(700 * time.Millisecond):
		}
	}
	return fmt.Errorf("Port %d wurde nach %s nicht erreichbar", port, timeout)
}

func portBusy(port int) bool {
	c, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 300*time.Millisecond)
	if err == nil { c.Close(); return true }
	return false
}

func (m *Manager) StartServer(id string) error {
	m.mu.Lock()
	rt := m.find(id)
	if rt == nil { m.mu.Unlock(); return fmt.Errorf("Server %s unbekannt", id) }
	if rt.State == "läuft" || rt.State == "startet" { m.mu.Unlock(); return nil }
	rt.State, rt.Err, rt.stopping = "startet", "", false
	rt.gen++
	gen := rt.gen
	sc := rt.Cfg
	cfg := *m.cfg
	m.mu.Unlock()
	m.Log("Starte " + sc.Name + " …")
	client, cmd, err := m.launch(sc, cfg, rt, gen)
	m.mu.Lock()
	defer m.mu.Unlock()
	if rt.gen != gen { if client != nil { client.Close() }; return nil }
	if err != nil {
		rt.State, rt.Err = "Fehler", err.Error()
		if cmd != nil && cmd.Process != nil { killTree(cmd) }
		m.Log(sc.Name + ": " + err.Error())
		return err
	}
	rt.client, rt.cmd, rt.State = client, cmd, "läuft"
	m.Log(sc.Name + " läuft")
	return nil
}

func (m *Manager) watch(rt *ServerRuntime, gen int, cmd *exec.Cmd, exited chan struct{}, name string) {
	err := cmd.Wait()
	close(exited)
	m.mu.Lock()
	defer m.mu.Unlock()
	if rt.gen != gen || rt.stopping { return }
	if rt.client != nil { rt.client.Close() }
	rt.client, rt.cmd = nil, nil
	if rt.State == "läuft" {
		rt.State = "Fehler"
		rt.Err = fmt.Sprintf("Prozess beendet (%v)", err)
		m.Log(name + ": Prozess unerwartet beendet – Neustart in 5 s")
		go func() {
			time.Sleep(5 * time.Second)
			m.mu.Lock()
			again := rt.gen == gen && !rt.stopping && rt.Cfg.Enabled
			m.mu.Unlock()
			if again { m.StartServer(rt.Cfg.ID) }
		}()
	}
}

func (m *Manager) startProcess(cmd *exec.Cmd, logName string, rt *ServerRuntime, gen int, name string, keepStdout bool) (chan struct{}, error) {
	lf, err := m.logFile(logName)
	if err != nil { return nil, err }
	fmt.Fprintf(lf, "\n===== %s Start %s =====\n", name, time.Now().Format(time.RFC3339))
	cmd.Stderr = lf
	if !keepStdout { cmd.Stdout = lf }
	prepareCmd(cmd)
	if err := cmd.Start(); err != nil { lf.Close(); return nil, fmt.Errorf("Start fehlgeschlagen: %w", err) }
	afterStart(cmd)
	exited := make(chan struct{})
	go func() { m.watch(rt, gen, cmd, exited, name); lf.Close() }()
	return exited, nil
}

func (m *Manager) launch(sc ServerConfig, cfg Config, rt *ServerRuntime, gen int) (Upstream, *exec.Cmd, error) {
	ctx := context.Background()
	switch sc.Kind {
	case "http":
		c := NewHTTPUpstream(sc.URL, sc.BearerToken)
		ictx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if _, err := c.Request(ictx, "ping", nil); err != nil {
			var re *RPCError
			if !errors.As(err, &re) { return nil, nil, err }
		}
		return c, nil, nil
	case "pc":
		if portBusy(cfg.PCPort) { return nil, nil, fmt.Errorf("Port %d ist schon belegt (alter PC-MCP?). Nutze „Alte Prozesse beenden“", cfg.PCPort) }
		uvx, err := findExe("uvx")
		if err != nil { return nil, nil, err }
		tok := randomToken()
		cmd := exec.Command(uvx, "--python", "3.14", "--from", "windows-mcp==0.8.7", "windows-mcp", "serve",
			"--transport", "streamable-http", "--host", "127.0.0.1", "--port", fmt.Sprint(cfg.PCPort), "--auth-key", tok)
		exited, err := m.startProcess(cmd, "pc-mcp", rt, gen, sc.Name, false)
		if err != nil { return nil, nil, err }
		if err := waitPort(ctx, cfg.PCPort, exited, 180*time.Second); err != nil { return nil, cmd, err }
		c := NewHTTPUpstream(fmt.Sprintf("http://127.0.0.1:%d/mcp", cfg.PCPort), tok)
		return c, cmd, probe(c)
	case "roblox":
		entry := filepath.Join(cfg.RobloxRepoDir, "server", "dist", "index.js")
		if _, err := os.Stat(entry); err != nil { return nil, nil, errors.New("Roblox MCP ist nicht installiert – zuerst „Roblox MCP installieren“ klicken") }
		if portBusy(cfg.RobloxHTTPPort) || portBusy(3667) { return nil, nil, errors.New("Port 3667/3668 ist belegt (alter Roblox-MCP?). Nutze „Alte Prozesse beenden“") }
		node, err := findExe("node")
		if err != nil { return nil, nil, err }
		tok := randomToken()
		cmd := exec.Command(node, "server/dist/index.js", "--transport", "http")
		cmd.Dir = cfg.RobloxRepoDir
		pub := "0"
		if cfg.RobloxAllowPublish { pub = "1" }
		cmd.Env = append(os.Environ(), "ROBLOX_MCP_TOKEN="+cfg.RobloxPluginToken, "ROBLOX_MCP_HTTP_TOKEN="+tok,
			"ROBLOX_MCP_TRANSPORT=http", "ROBLOX_MCP_HTTP_HOST=127.0.0.1", fmt.Sprintf("ROBLOX_MCP_HTTP_PORT=%d", cfg.RobloxHTTPPort),
			"ROBLOX_MCP_ALLOW_PUBLISH="+pub)
		exited, err := m.startProcess(cmd, "roblox-mcp", rt, gen, sc.Name, false)
		if err != nil { return nil, nil, err }
		if err := waitPort(ctx, cfg.RobloxHTTPPort, exited, 90*time.Second); err != nil { return nil, cmd, err }
		c := NewHTTPUpstream(fmt.Sprintf("http://127.0.0.1:%d/mcp", cfg.RobloxHTTPPort), tok)
		return c, cmd, probe(c)
	case "browser":
		npx, err := findExe("npx")
		if err != nil { return nil, nil, errors.New("npx nicht gefunden – Node.js installieren") }
		bp := cfg.BrowserPath
		if bp == "" { bp = findBrowser() }
		if bp == "" { return nil, nil, errors.New("Brave/Edge/Chrome nicht gefunden – browserPath in config.json setzen") }
		prof := cfg.BrowserProfileDir
		if prof == "" { prof = filepath.Join(filepath.Dir(m.cfgPath), "browser-profile") }
		os.MkdirAll(prof, 0o700)
		cmd := exec.Command(npx, "-y", "@playwright/mcp@latest", "--executable-path", bp, "--user-data-dir", prof)
		cmd.Env = os.Environ()
		stdin, err := cmd.StdinPipe()
		if err != nil { return nil, nil, err }
		stdout, err := cmd.StdoutPipe()
		if err != nil { return nil, nil, err }
		if _, err := m.startProcess(cmd, "browser-mcp", rt, gen, sc.Name, true); err != nil { return nil, nil, err }
		c := NewStdioUpstream(stdin, stdout, nil)
		return c, cmd, probe(c)
	case "stdio":
		exe, err := findExe(sc.Command)
		if err != nil { exe = sc.Command }
		cmd := exec.Command(exe, sc.Args...)
		cmd.Dir = sc.WorkDir
		cmd.Env = os.Environ()
		for k, v := range sc.Env { cmd.Env = append(cmd.Env, k+"="+v) }
		stdin, err := cmd.StdinPipe()
		if err != nil { return nil, nil, err }
		stdout, err := cmd.StdoutPipe()
		if err != nil { return nil, nil, err }
		if _, err := m.startProcess(cmd, "server-"+sc.ID, rt, gen, sc.Name, true); err != nil { return nil, nil, err }
		c := NewStdioUpstream(stdin, stdout, nil)
		return c, cmd, probe(c)
	}
	return nil, nil, fmt.Errorf("unbekannter Servertyp %s", sc.Kind)
}

// probe initializes the session and checks that tools can be listed.
func probe(c Upstream) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if _, err := c.Request(ctx, "tools/list", nil); err != nil { return fmt.Errorf("MCP antwortet nicht: %w", err) }
	return nil
}

func (m *Manager) stopRuntime(rt *ServerRuntime) {
	m.mu.Lock()
	rt.stopping = true
	rt.gen++
	cmd, client := rt.cmd, rt.client
	rt.cmd, rt.client = nil, nil
	rt.State, rt.Err = "gestoppt", ""
	m.mu.Unlock()
	if client != nil { client.Close() }
	if cmd != nil && cmd.Process != nil { killTree(cmd) }
}

func (m *Manager) StopServer(id string) {
	m.mu.Lock()
	rt := m.find(id)
	m.mu.Unlock()
	if rt != nil { m.stopRuntime(rt); m.Log(rt.Cfg.Name + " gestoppt") }
}

func (m *Manager) SetEnabled(id string, on bool) error {
	m.mu.Lock()
	for i := range m.cfg.Servers { if m.cfg.Servers[i].ID == id { m.cfg.Servers[i].Enabled = on } }
	if rt := m.find(id); rt != nil { rt.Cfg.Enabled = on }
	err := m.save()
	m.mu.Unlock()
	return err
}

// ---------- ngrok ----------

func (m *Manager) StartNgrok() error {
	m.mu.Lock()
	if m.ngrokCmd != nil { m.mu.Unlock(); return nil }
	domain, port := m.cfg.NgrokDomain, m.cfg.BridgePort
	m.ngrokStop = false
	m.ngrokState, m.ngrokErr = "startet", ""
	m.mu.Unlock()
	if domain == "" { m.setNgrok("nicht konfiguriert", "ngrokDomain fehlt in config.json"); return errors.New("ngrokDomain fehlt") }
	exe, err := findExe("ngrok")
	if err != nil { m.setNgrok("Fehler", err.Error()); return err }
	for _, flag := range []string{"--url", "--domain"} {
		err := m.runNgrok(exe, flag, domain, port)
		if err == nil { return nil }
		if !strings.Contains(err.Error(), "unknown flag") { m.setNgrok("Fehler", err.Error()); return err }
	}
	return errors.New("ngrok konnte nicht gestartet werden")
}

func (m *Manager) setNgrok(state, e string) { m.mu.Lock(); m.ngrokState, m.ngrokErr = state, e; m.mu.Unlock() }

func (m *Manager) runNgrok(exe, flag, domain string, port int) error {
	cmd := exec.Command(exe, "http", fmt.Sprintf("127.0.0.1:%d", port), flag+"=https://"+domain, "--log", "stdout", "--log-format", "logfmt")
	if flag == "--domain" { cmd.Args[3] = flag + "=" + domain }
	out, err := cmd.StdoutPipe()
	if err != nil { return err }
	lf, err := m.logFile("ngrok")
	if err != nil { return err }
	cmd.Stderr = lf
	prepareCmd(cmd)
	if err := cmd.Start(); err != nil { return fmt.Errorf("ngrok-Start fehlgeschlagen: %w", err) }
	afterStart(cmd)
	online := make(chan struct{})
	lastErr := make(chan string, 1)
	go func() {
		var once sync.Once
		sc := bufio.NewScanner(out)
		var lastE string
		for sc.Scan() {
			line := sc.Text()
			io.WriteString(lf, line+"\n")
			if strings.Contains(line, "started tunnel") || (strings.Contains(line, "url=https://") && strings.Contains(line, domain)) {
				once.Do(func() { close(online) })
			}
			if strings.Contains(line, "lvl=eror") || strings.Contains(line, "lvl=crit") || strings.Contains(line, "ERR_NGROK") || strings.Contains(line, "unknown flag") {
				lastE = ngrokErrText(line)
			}
		}
		lastErr <- lastE
	}()
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	select {
	case <-online:
	case e := <-exited:
		msg := <-lastErr
		lf.Close()
		if msg == "" { msg = fmt.Sprint(e) }
		return errors.New(msg)
	case <-time.After(40 * time.Second):
		killTree(cmd)
		return errors.New("ngrok hat innerhalb von 40 s keinen Tunnel aufgebaut")
	}
	m.mu.Lock()
	m.ngrokCmd = cmd
	m.ngrokState, m.ngrokErr = "online", ""
	m.mu.Unlock()
	m.Log("ngrok online: https://" + domain + "/mcp")
	go func() {
		e := <-exited
		msg := <-lastErr
		lf.Close()
		m.mu.Lock()
		stop := m.ngrokStop
		m.ngrokCmd = nil
		if !stop {
			m.ngrokState, m.ngrokErr = "Fehler", strings.TrimSpace(msg+" "+fmt.Sprint(e))
		}
		m.mu.Unlock()
		if !stop {
			m.Log("ngrok wurde beendet – Neustart in 10 s")
			time.Sleep(10 * time.Second)
			m.mu.Lock(); stop = m.ngrokStop; m.mu.Unlock()
			if !stop { m.StartNgrok() }
		}
	}()
	return nil
}

func ngrokErrText(line string) string {
	switch {
	case strings.Contains(line, "ERR_NGROK_4018"), strings.Contains(line, "authtoken"):
		return "ngrok-Authtoken fehlt: einmal „ngrok config add-authtoken <TOKEN>“ ausführen"
	case strings.Contains(line, "ERR_NGROK_108"):
		return "Es läuft bereits eine andere ngrok-Sitzung. Nutze „Alte Prozesse beenden“"
	case strings.Contains(line, "ERR_NGROK_334"):
		return "Die Domain ist schon an einen anderen ngrok-Tunnel gebunden. Nutze „Alte Prozesse beenden“"
	case strings.Contains(line, "ERR_NGROK_313"), strings.Contains(line, "ERR_NGROK_320"):
		return "Diese Domain gehört nicht zu deinem ngrok-Konto"
	}
	if i := strings.Index(line, "err="); i >= 0 { return strings.Trim(line[i+4:], `"`) }
	return line
}

func (m *Manager) StopNgrok() {
	m.mu.Lock()
	m.ngrokStop = true
	cmd := m.ngrokCmd
	m.ngrokCmd = nil
	m.ngrokState, m.ngrokErr = "gestoppt", ""
	m.mu.Unlock()
	if cmd != nil { killTree(cmd) }
}

// ---------- combined ----------

func (m *Manager) StartAll() error {
	var errs []string
	if err := m.StartBridge(); err != nil { errs = append(errs, err.Error()) }
	m.mu.Lock()
	var ids []string
	for _, s := range m.servers { if s.Cfg.Enabled { ids = append(ids, s.Cfg.ID) } }
	m.mu.Unlock()
	var wg sync.WaitGroup
	var emu sync.Mutex
	wg.Add(1)
	go func() { defer wg.Done(); if err := m.StartNgrok(); err != nil { emu.Lock(); errs = append(errs, "ngrok: "+err.Error()); emu.Unlock() } }()
	for _, id := range ids {
		wg.Add(1)
		go func(id string) { defer wg.Done(); if err := m.StartServer(id); err != nil { emu.Lock(); errs = append(errs, id+": "+err.Error()); emu.Unlock() } }(id)
	}
	wg.Wait()
	if len(errs) > 0 { return errors.New(strings.Join(errs, "\n")) }
	return nil
}

func (m *Manager) StopAll() {
	m.StopNgrok()
	m.mu.Lock()
	list := append([]*ServerRuntime(nil), m.servers...)
	m.mu.Unlock()
	var wg sync.WaitGroup
	for _, rt := range list { wg.Add(1); go func(rt *ServerRuntime) { defer wg.Done(); m.stopRuntime(rt) }(rt) }
	wg.Wait()
	m.StopBridge()
	m.Log("Alles gestoppt")
}

func (m *Manager) RotateToken() string {
	m.mu.Lock()
	m.cfg.BridgeToken = randomToken()
	t := m.cfg.BridgeToken
	m.save()
	m.mu.Unlock()
	m.Log("Neuer Bridge-Token erzeugt – in Notion aktualisieren")
	return t
}

func (m *Manager) runLogged(dir, name string, args ...string) error {
	exe, err := findExe(name)
	if err != nil { return err }
	cmd := exec.Command(exe, args...)
	cmd.Dir = dir
	lf, err := m.logFile("roblox-setup")
	if err != nil { return err }
	defer lf.Close()
	cmd.Stdout, cmd.Stderr = lf, lf
	prepareCmd(cmd)
	m.Log(name + " " + strings.Join(args, " "))
	if err := cmd.Run(); err != nil { return fmt.Errorf("%s %s fehlgeschlagen: %w (siehe roblox-setup.log)", name, strings.Join(args, " "), err) }
	return nil
}

func (m *Manager) SetupRoblox() error {
	cfg := m.Config()
	dir := cfg.RobloxRepoDir
	if _, err := os.Stat(filepath.Join(dir, ".git")); err != nil {
		os.MkdirAll(filepath.Dir(dir), 0o755)
		if err := m.runLogged(filepath.Dir(dir), "git", "clone", "https://github.com/Meganugger/roblox-studio-mcp.git", dir); err != nil { return err }
	} else if err := m.runLogged(dir, "git", "pull", "--ff-only"); err != nil {
		m.Log("git pull übersprungen: " + err.Error())
	}
	if err := m.runLogged(dir, "npm", "install"); err != nil { return err }
	if err := m.runLogged(dir, "npm", "run", "build"); err != nil { return err }
	if err := m.runLogged(dir, "node", "server/dist/index.js", "--install-plugin"); err != nil { return err }
	m.Log("Roblox MCP installiert. Plugin-Token einmal in Roblox Studio → Plugins → MCP eintragen.")
	return nil
}
