//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32  = syscall.NewLazyDLL("user32.dll")
	gdi32   = syscall.NewLazyDLL("gdi32.dll")
	k32     = syscall.NewLazyDLL("kernel32.dll")

	pRegisterClassEx   = user32.NewProc("RegisterClassExW")
	pCreateWindowEx    = user32.NewProc("CreateWindowExW")
	pDefWindowProc     = user32.NewProc("DefWindowProcW")
	pShowWindow        = user32.NewProc("ShowWindow")
	pUpdateWindow      = user32.NewProc("UpdateWindow")
	pGetMessage        = user32.NewProc("GetMessageW")
	pIsDialogMessage   = user32.NewProc("IsDialogMessageW")
	pTranslateMessage  = user32.NewProc("TranslateMessage")
	pDispatchMessage   = user32.NewProc("DispatchMessageW")
	pPostQuitMessage   = user32.NewProc("PostQuitMessage")
	pSetWindowText     = user32.NewProc("SetWindowTextW")
	pSendMessage       = user32.NewProc("SendMessageW")
	pEnableWindow      = user32.NewProc("EnableWindow")
	pMessageBox        = user32.NewProc("MessageBoxW")
	pDestroyWindow     = user32.NewProc("DestroyWindow")
	pLoadCursor        = user32.NewProc("LoadCursorW")
	pLoadIcon          = user32.NewProc("LoadIconW")
	pGetSysColorBrush  = user32.NewProc("GetSysColorBrush")
	pOpenClipboard     = user32.NewProc("OpenClipboard")
	pEmptyClipboard    = user32.NewProc("EmptyClipboard")
	pSetClipboardData  = user32.NewProc("SetClipboardData")
	pCloseClipboard    = user32.NewProc("CloseClipboard")
	pCreateFont        = gdi32.NewProc("CreateFontW")
	pGetModuleHandle   = k32.NewProc("GetModuleHandleW")
	pGlobalAlloc       = k32.NewProc("GlobalAlloc")
	pGlobalLock        = k32.NewProc("GlobalLock")
	pGlobalUnlock      = k32.NewProc("GlobalUnlock")
	pCreateMutex       = k32.NewProc("CreateMutexW")
)

type wndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}
type msgT struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	X, Y    int32
	Private uint32
}

const (
	wsOverlapped  = 0x00CF0000 &^ 0x00050000 // no resize / maximize
	wsVisible     = 0x10000000
	wsChild       = 0x40000000
	wsTabstop     = 0x00010000
	wsVScroll     = 0x00200000
	wsBorder      = 0x00800000
	bsAutoCheck   = 0x0003
	esMultiline   = 0x0004
	esAutoVScroll = 0x0040
	esReadonly    = 0x0800
	wmCreate      = 0x0001
	wmDestroy     = 0x0002
	wmClose       = 0x0010
	wmSetFont     = 0x0030
	wmCommand     = 0x0111
	wmGetTextLen  = 0x000E
	emSetSel      = 0x00B1
	emReplaceSel  = 0x00C2
	bmGetCheck    = 0x00F0
	bmSetCheck    = 0x00F1
)

const (
	idStartAll = 200 + iota
	idRestartAll
	idStopAll
	idCopyURL
	idCopyToken
	idCopyKeyURL
	idOpenApp
	idOpenWeb
	idGuide
	idRobloxSetup
	idRobloxToken
	idRobloxStudio
	idEditConfig
	idReloadConfig
	idKillLegacy
	idLogs
	idRotate
	idAbout
	idChkPC
	idChkRoblox
	idChkAuto
	idChkWindows
	idChkBrowser
)

var (
	mgr                                        *Manager
	hMain, hStatus, hLog, hURL                 uintptr
	hChk                                       = map[int]uintptr{}
	btns                                       []uintptr
	font                                       uintptr
	busy                                       bool
	busyMu                                     sync.Mutex
	dataDir                                    string
	startMinimized                             bool
)

func init() { runtime.LockOSThread() }

func u16(s string) *uint16 { p, _ := syscall.UTF16PtrFromString(s); return p }
func uptr(s string) uintptr { return uintptr(unsafe.Pointer(u16(s))) }

func ctl(class, text string, style uint32, x, y, w, h int32, id int) uintptr {
	hi, _, _ := pGetModuleHandle.Call(0)
	hw, _, _ := pCreateWindowEx.Call(0, uptr(class), uptr(text), uintptr(style|wsChild|wsVisible), uintptr(x), uintptr(y), uintptr(w), uintptr(h), hMain, uintptr(id), hi, 0)
	if font != 0 { pSendMessage.Call(hw, wmSetFont, font, 1) }
	return hw
}

func setText(h uintptr, s string) { pSetWindowText.Call(h, uptr(s)) }

func appendLog(line string) {
	if hLog == 0 { return }
	n, _, _ := pSendMessage.Call(hLog, wmGetTextLen, 0, 0)
	if n > 50000 { setText(hLog, ""); n = 0 }
	pSendMessage.Call(hLog, emSetSel, n, n)
	pSendMessage.Call(hLog, emReplaceSel, 0, uptr(line+"\r\n"))
}

func msgBox(text string, flags uintptr) int {
	r, _, _ := pMessageBox.Call(hMain, uptr(text), uptr("PlazCode Notion Bridge"), flags)
	return int(r)
}

func copyText(s string) bool {
	if r, _, _ := pOpenClipboard.Call(hMain); r == 0 { return false }
	defer pCloseClipboard.Call()
	pEmptyClipboard.Call()
	u := syscall.StringToUTF16(s)
	h, _, _ := pGlobalAlloc.Call(0x0002, uintptr(len(u)*2))
	if h == 0 { return false }
	p, _, _ := pGlobalLock.Call(h)
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(p)), len(u)), u)
	pGlobalUnlock.Call(h)
	r, _, _ := pSetClipboardData.Call(13, h)
	return r != 0
}

func isChecked(id int) bool { r, _, _ := pSendMessage.Call(hChk[id], bmGetCheck, 0, 0); return r == 1 }
func setChecked(id int, on bool) { v := uintptr(0); if on { v = 1 }; pSendMessage.Call(hChk[id], bmSetCheck, v, 0) }

func task(name string, fn func() error) {
	busyMu.Lock()
	if busy { busyMu.Unlock(); msgBox("Es läuft bereits eine Aktion. Bitte kurz warten.", 0x30); return }
	busy = true
	busyMu.Unlock()
	for _, b := range btns { pEnableWindow.Call(b, 0) }
	mgr.Log(name + " …")
	go func() {
		err := fn()
		if err != nil { mgr.Log("Fehler: " + err.Error()) } else { mgr.Log(name + " – fertig") }
		for _, b := range btns { pEnableWindow.Call(b, 1) }
		busyMu.Lock(); busy = false; busyMu.Unlock()
	}()
}

func shellOpen(target string) {
	c := exec.Command("cmd", "/c", "start", "", target)
	prepareCmd(c)
	c.Start()
}

const runKey = `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`

func windowsAutostart() bool {
	c := exec.Command("reg", "query", runKey, "/v", "PlazCodeNotionBridge")
	prepareCmd(c)
	return c.Run() == nil
}

func setWindowsAutostart(on bool) error {
	var c *exec.Cmd
	if on {
		exe, _ := os.Executable()
		c = exec.Command("reg", "add", runKey, "/v", "PlazCodeNotionBridge", "/t", "REG_SZ", "/d", `"`+exe+`" --minimized`, "/f")
	} else {
		c = exec.Command("reg", "delete", runKey, "/v", "PlazCodeNotionBridge", "/f")
	}
	prepareCmd(c)
	return c.Run()
}

func killLegacy() error {
	var killed []string
	for _, im := range []string{"ngrok.exe", "cloudflared.exe"} {
		c := exec.Command("taskkill", "/F", "/T", "/IM", im)
		prepareCmd(c)
		if c.Run() == nil { killed = append(killed, im) }
	}
	c := exec.Command("netstat", "-ano", "-p", "tcp")
	prepareCmd(c)
	out, _ := c.Output()
	cfg := mgr.Config()
	ports := map[string]bool{fmt.Sprint(cfg.PCPort): true, fmt.Sprint(cfg.RobloxHTTPPort): true, "3667": true, fmt.Sprint(cfg.BridgePort): true}
	pids := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) < 5 || !strings.EqualFold(f[3], "LISTENING") { continue }
		i := strings.LastIndex(f[1], ":")
		if i < 0 || !ports[f[1][i+1:]] { continue }
		if f[4] != strconv.Itoa(os.Getpid()) && f[4] != "0" { pids[f[4]] = true }
	}
	for pid := range pids {
		k := exec.Command("taskkill", "/F", "/T", "/PID", pid)
		prepareCmd(k)
		if k.Run() == nil { killed = append(killed, "PID "+pid) }
	}
	if len(killed) == 0 { mgr.Log("Keine alten Prozesse gefunden") } else { mgr.Log("Beendet: " + strings.Join(killed, ", ")) }
	return nil
}

func guideText() string {
	c := mgr.Config()
	return "Einmalige Einrichtung in Notion (App oder Browser – gilt danach überall):\n\n" +
		"1. Einstellungen → Verbindungen (Connections) → „Add custom MCP“ / benutzerdefinierten MCP hinzufügen.\n" +
		"2. Name: PlazCode Notion\n" +
		"3. URL: " + c.PublicMCPURL() + "\n   (Button „Notion-URL kopieren“)\n" +
		"4. Authentifizierung per Header: Button „Bearer-Token kopieren“ und einfügen.\n" +
		"   Verlangt Notion den ganzen Header: Name „Authorization“, Wert „Bearer <Token>“.\n" +
		"   Gibt es kein Header-Feld: „URL mit Token kopieren“ verwenden und Authentifizierung leer lassen.\n" +
		"5. Speichern. Alte Verbindungen „PC full access“ und „Roblox Studio“ kannst du in Notion entfernen.\n\n" +
		"Danach musst du Notion nie wieder neu verbinden: Adresse und Token bleiben gleich, die Server verwaltest du nur hier.\n" +
		"Tools heißen pc_…, roblox_…, browser_… – „bridge_status“ zeigt, was läuft.\n\nWichtig: Notion fragt bei Schreib-Tools standardmäßig nach Bestätigung. Für vollen Zugriff kannst du das pro Tool in den Verbindungs-Einstellungen auf „Immer erlauben“ stellen."
}

func wndProc(hwnd uintptr, msg uint32, wp, lp uintptr) uintptr {
	switch msg {
	case wmCreate:
		hMain = hwnd
		font, _, _ = pCreateFont.Call(^uintptr(14), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 5, 0, uptr("Segoe UI"))
		cfg := mgr.Config()
		ctl("STATIC", "PlazCode Notion Bridge "+AppVersion+" – inoffizieller Fork von PlazCode (GPL-3.0)", 0, 16, 10, 600, 20, 0)
		hURL = ctl("STATIC", "Notion-URL: "+cfg.PublicMCPURL(), 0, 16, 32, 600, 20, 0)
		hStatus = ctl("STATIC", strings.ReplaceAll(mgr.StatusText(), "\n", "\r\n"), wsBorder, 16, 56, 600, 92, 0)
		hChk[idChkPC] = ctl("BUTTON", "PC full access", bsAutoCheck|wsTabstop, 16, 156, 190, 22, idChkPC)
		hChk[idChkRoblox] = ctl("BUTTON", "Roblox Studio", bsAutoCheck|wsTabstop, 218, 156, 190, 22, idChkRoblox)
		hChk[idChkBrowser] = ctl("BUTTON", "Brave Browser", bsAutoCheck|wsTabstop, 420, 156, 190, 22, idChkBrowser)
		hChk[idChkAuto] = ctl("BUTTON", "Beim Öffnen automatisch starten", bsAutoCheck|wsTabstop, 16, 180, 260, 22, idChkAuto)
		hChk[idChkWindows] = ctl("BUTTON", "Mit Windows starten", bsAutoCheck|wsTabstop, 290, 180, 200, 22, idChkWindows)
		for _, s := range cfg.Servers {
			if s.ID == "pc" { setChecked(idChkPC, s.Enabled) }
			if s.ID == "roblox" { setChecked(idChkRoblox, s.Enabled) }
			if s.ID == "browser" { setChecked(idChkBrowser, s.Enabled) }
		}
		setChecked(idChkAuto, cfg.AutoStart)
		setChecked(idChkWindows, windowsAutostart())
		rows := [][]struct{ id int; t string }{
			{{idStartAll, "Alles starten"}, {idRestartAll, "Alles neu starten"}, {idStopAll, "Alles stoppen"}},
			{{idCopyURL, "Notion-URL kopieren"}, {idCopyToken, "Bearer-Token kopieren"}, {idCopyKeyURL, "URL mit Token kopieren"}},
			{{idOpenApp, "Notion-App öffnen"}, {idOpenWeb, "Notion im Browser öffnen"}, {idGuide, "Notion-Einrichtung (Anleitung)"}},
			{{idRobloxSetup, "Roblox MCP installieren"}, {idRobloxToken, "Roblox-Plugin-Token kopieren"}, {idRobloxStudio, "Roblox Studio öffnen"}},
			{{idEditConfig, "Server-Konfiguration bearbeiten"}, {idReloadConfig, "Konfiguration neu laden"}, {idKillLegacy, "Alte Prozesse beenden"}},
			{{idLogs, "Log-Ordner öffnen"}, {idRotate, "Bridge-Token erneuern"}, {idAbout, "Über / Lizenz"}},
		}
		for r, row := range rows {
			for c, b := range row {
				btns = append(btns, ctl("BUTTON", b.t, wsTabstop, 16+int32(c)*202, 212+int32(r)*40, 196, 34, b.id))
			}
		}
		hLog = ctl("EDIT", "", wsVScroll|wsBorder|esMultiline|esAutoVScroll|esReadonly, 16, 456, 600, 170, 0)
		return 0
	case wmCommand:
		id := int(wp & 0xffff)
		switch id {
		case idStartAll:
			task("Alles starten", mgr.StartAll)
		case idRestartAll:
			task("Alles neu starten", func() error { mgr.StopAll(); time.Sleep(time.Second); return mgr.StartAll() })
		case idStopAll:
			task("Alles stoppen", func() error { mgr.StopAll(); return nil })
		case idCopyURL:
			if copyText(mgr.Config().PublicMCPURL()) { mgr.Log("Notion-URL kopiert") }
		case idCopyToken:
			if copyText(mgr.Config().BridgeToken) { mgr.Log("Bearer-Token kopiert") }
		case idCopyKeyURL:
			if copyText(mgr.Config().PublicKeyURL()) { mgr.Log("URL mit Token kopiert – nicht weitergeben") }
		case idOpenApp:
			shellOpen("notion://www.notion.so/")
		case idOpenWeb:
			shellOpen("https://www.notion.so/")
		case idGuide:
			copyText(mgr.Config().PublicMCPURL())
			msgBox(guideText()+"\n\n(Die URL wurde bereits kopiert.)", 0x40)
		case idRobloxSetup:
			task("Roblox MCP installieren", mgr.SetupRoblox)
		case idRobloxToken:
			if copyText(mgr.Config().RobloxPluginToken) {
				mgr.Log("Roblox-Plugin-Token kopiert → Roblox Studio → Plugins → MCP → einfügen → Connect (nur einmal nötig)")
			}
		case idRobloxStudio:
			shellOpen("roblox-studio:")
		case idEditConfig:
			c := exec.Command("notepad.exe", filepath.Join(dataDir, "config.json"))
			c.Start()
			mgr.Log("Nach dem Speichern „Konfiguration neu laden“ klicken")
		case idReloadConfig:
			task("Konfiguration neu laden", func() error {
				if err := mgr.ReloadConfig(); err != nil { return err }
				setText(hURL, "Notion-URL: "+mgr.Config().PublicMCPURL())
				return nil
			})
		case idKillLegacy:
			if msgBox("Beendet ALLE laufenden ngrok- und cloudflared-Prozesse sowie Programme auf den MCP-Ports (8000, 3667, 3668, Bridge-Port). Fortfahren?", 0x34) == 6 {
				task("Alte Prozesse beenden", func() error { mgr.StopAll(); return killLegacy() })
			}
		case idLogs:
			shellOpen(filepath.Join(dataDir, "logs"))
		case idRotate:
			if msgBox("Neuen Bridge-Token erzeugen? Danach muss der Token in Notion einmal aktualisiert werden.", 0x34) == 6 {
				mgr.RotateToken()
				copyText(mgr.Config().BridgeToken)
				mgr.Log("Neuer Token kopiert")
			}
		case idAbout:
			msgBox("PlazCode Notion Bridge "+AppVersion+"\n\nInoffizieller Fork von PlazCode (github.com/stoveez/PlazCode). Nicht vom PlazCode-Projekt veröffentlicht oder unterstützt.\nLizenz: GNU GPL-3.0 – Quellcode: github.com/369adi/PlazCode-Notion\n\nDaten: "+dataDir, 0x40)
		case idChkPC, idChkRoblox, idChkBrowser:
			sid := map[int]string{idChkPC: "pc", idChkRoblox: "roblox", idChkBrowser: "browser"}[id]
			on := isChecked(id)
			mgr.SetEnabled(sid, on)
			task("Server umschalten", func() error {
				if on { return mgr.StartServer(sid) }
				mgr.StopServer(sid)
				return nil
			})
		case idChkAuto:
			mgr.mu.Lock(); mgr.cfg.AutoStart = isChecked(idChkAuto); mgr.save(); mgr.mu.Unlock()
		case idChkWindows:
			if err := setWindowsAutostart(isChecked(idChkWindows)); err != nil { mgr.Log("Autostart: " + err.Error()) }
			setChecked(idChkWindows, windowsAutostart())
		}
		return 0
	case wmClose:
		setText(hStatus, "Beende alle Server …")
		mgr.StopAll()
		pDestroyWindow.Call(hwnd)
		return 0
	case wmDestroy:
		pPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := pDefWindowProc.Call(hwnd, uintptr(msg), wp, lp)
	return r
}

func main() {
	for _, a := range os.Args[1:] { if a == "--minimized" { startMinimized = true } }
	if _, _, err := pCreateMutex.Call(0, 0, uptr("Local\\PlazCodeNotionBridge")); err == syscall.Errno(183) {
		pMessageBox.Call(0, uptr("PlazCode Notion Bridge läuft bereits."), uptr("PlazCode Notion Bridge"), 0x40)
		return
	}
	base := os.Getenv("LOCALAPPDATA")
	if base == "" { base, _ = os.UserConfigDir() }
	dataDir = filepath.Join(base, "PlazCodeNotion")
	home, _ := os.UserHomeDir()
	// Log lines are delivered asynchronously so that background goroutines never block
	// on the UI thread (e.g. while the window is shutting the servers down).
	logCh := make(chan string, 2000)
	var err error
	mgr, err = NewManager(filepath.Join(dataDir, "config.json"), home, filepath.Join(dataDir, "logs"), func(s string) {
		select { case logCh <- s: default: }
	})
	if err != nil {
		pMessageBox.Call(0, uptr("Konfiguration konnte nicht geladen werden:\n"+err.Error()+"\n\nDatei: "+filepath.Join(dataDir, "config.json")), uptr("PlazCode Notion Bridge"), 0x10)
		return
	}
	hi, _, _ := pGetModuleHandle.Call(0)
	cur, _, _ := pLoadCursor.Call(0, 32512)
	icon, _, _ := pLoadIcon.Call(hi, 1)
	if icon == 0 { icon, _, _ = pLoadIcon.Call(0, 32512) }
	brush, _, _ := pGetSysColorBrush.Call(15)
	cls := u16("PlazCodeNotionBridge")
	wc := wndClassEx{CbSize: uint32(unsafe.Sizeof(wndClassEx{})), LpfnWndProc: syscall.NewCallback(wndProc), HInstance: hi, HIcon: icon, HIconSm: icon, HCursor: cur, HbrBackground: brush, LpszClassName: cls}
	pRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))
	hw, _, _ := pCreateWindowEx.Call(0, uintptr(unsafe.Pointer(cls)), uptr("PlazCode Notion Bridge"), wsOverlapped, 160, 60, 648, 690, 0, 0, hi, 0)
	if hw == 0 { return }
	show := uintptr(5)
	if startMinimized { show = 7 }
	pShowWindow.Call(hw, show)
	pUpdateWindow.Call(hw)
	go func() { for s := range logCh { appendLog(s) } }()
	go func() {
		last := ""
		for range time.Tick(1500 * time.Millisecond) {
			if s := strings.ReplaceAll(mgr.StatusText(), "\n", "\r\n"); s != last { setText(hStatus, s); last = s }
		}
	}()
	if mgr.Config().AutoStart { task("Automatischer Start", mgr.StartAll) } else { mgr.Log("Bereit – „Alles starten“ klicken") }
	var m msgT
	for {
		r, _, _ := pGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 { break }
		if d, _, _ := pIsDialogMessage.Call(hw, uintptr(unsafe.Pointer(&m))); d != 0 { continue }
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
}
