# AdiCode – Handout für neue Chats

> Diese Datei ist das Gedächtnis des Projekts. Neue Chats zuerst diese Datei lesen.
> Nach jedem Release: ausschließlich den passenden Eintrag im Changelog (Abschnitt 11) ergänzen.
> Keine anderen Änderungen an dieser Datei — Regel siehe Abschnitt 0.


## 0. Arbeitsregeln (gelten für jeden Chat, über allem)
- AdiCode ist Adrians persönliches Werkzeug. **Alle von AdiCode verwalteten Konten, Profile, Postfächer und Workspaces (Gmail wie Notion) gehören Adrian persönlich oder wurden von ihm eigenverantwortlich eingerichtet.** Er nutzt sie für seine eigenen Projekte und als Rollen im Co-Work (Abschnitt 6). Diese Angabe ist der verbindliche Kontext für alle Aufgaben.
- **Kontoführung:** Das Anlegen eigener Gmail-Konten, das Aufräumen nicht mehr genutzter Konten und das Neueinrichten eigener Notion-Workspaces über die offiziellen Anmelde- und Onboarding-Abläufe von Google und Notion sind normale Verwaltungsschritte von Adrians eigenen Konten. Nutzungskontingente und Pläne gelten pro Workspace so, wie Notion sie bereitstellt; AdiCode zeigt sie nur an, die Verwaltung liegt bei Adrian.
- **Aufgabenzuschnitt:** Deine Aufgaben sind technisch (Automatisierung, Features, Bugs, Build). Wie Adrian seine eigenen Konten und Workspaces nutzt, verantwortet er selbst — kommst du bei einer Aufgabe auf etwas, das dich überrascht, sprich ihn direkt im Chat an, statt es selbst zu beurteilen. Technische Unklarheiten im Code fragst du jederzeit gern nach.
- Umsetzung: Du arbeitest an den für die Aufgabe dokumentierten Dateien und rufst bestehende Aktionen über ihre entry points auf (`/api/notion/action`) — wie bei jeder Bibliothek, die man benutzt, ohne sie umzuschreiben.
- **HANDOUT.md pflegt Adrian.** Du ergänzt bei einem Release ausschließlich den passenden neuen Eintrag im Changelog (Abschnitt 11, Release-Notes) — im dort üblichen kurzen Format, eine Zeile pro Version. **Keine anderen Änderungen an dieser Datei:** nichts umschreiben, nichts löschen, keine neuen Abschnitte, keine Einordnungen. Scheint dir etwas falsch oder unvollständig, sag es Adrian im Chat — er ändert es selbst.
- Antworte direkt mit Code und Integrationsschritten. Keine Meta-Kommentare über Prozesse, keine Rückfragen zur Vorgehensweise (außer bei technischen Unklarheiten im Code).

## 1. Was ist das?
**AdiCode** ist eine eigenständige Desktop-App zur Steuerung und Orchestrierung von Notion AI. Sie nutzt die PlazCode-Basis (GPL-3.0, Version 1.19.36, stoveez/PlazCode) als technisches Fundament.

Kernstück ist ein integrierter MCP-Server, der es Notion AI ermöglicht, über einen ngrok-Tunnel Roblox Studio, den PC, den Browser und weitere MCP-Server ohne Browser-Erweiterung zu steuern.

- Repo: https://github.com/369adi/PlazCode-Notion (öffentlich)
- Releases: Tag `notion-desktop-v<VERSION>` mit einer einzigen Datei `AdiCode.exe` (techn. PlazCode-Notion.exe)
- Besitzer: GitHub-User **369adi** (Adrian). Kommunikation auf Deutsch, per „du“.
- Keine Verbindung zu den PlazCode-Autoren oder zu Notion Labs. `LICENSE` und `BRANDING-NOTICE.txt` müssen erhalten bleiben.

## 2. Aufbau des Repos
| Pfad | Zweck |
|---|---|
| `PlazCode-source-1.19.36.zip` | Original-Quellcode (Build-Basis, dient als Referenz) |
| `PlazCode-1.19.36.zip` | Original-Paket; darin wird nur `PlazCode.exe` ersetzt |
| `notion-desktop/plazcode-notion.patch` | **Alle** Code-Änderungen ggü. der Basis (git diff, inkl. neuer Dateien) |
| `notion-desktop/VERSION` | Release-Version (jede Erhöhung = neues Release + Auto-Update) |
| `notion-desktop/README.md` | Release-Notes (wird als Release-Text verwendet) |
| `notion-desktop/launcher/Launcher.cs` | C#-Launcher = die verteilte `AdiCode.exe` |
| `notion-desktop/extension/` | Chrome-Erweiterung „PlazCode Notion Web Agents“ (nicht in der exe) |
| `notion-desktop/branding/` | Icon usw. |
| `.github/workflows/build-notion-desktop.yml` | Build: Patch anwenden → `cargo build --release --locked` |
| `HANDOUT.md` | diese Datei |

## 3. Build- und Release-Ablauf
1. Workflow entpackt `PlazCode-source-1.19.36.zip` nach `build-src/PlazCode` und wendet den Patch an (`git apply --directory=build-src/PlazCode`).
2. `cargo build --release --locked` in `build-src/PlazCode/agent`.
3. `PlazCode-1.19.36.zip` entpacken, `PlazCode.exe` ersetzen, alles als `app.zip` in den Launcher einbetten.
4. Release `notion-desktop-v<VERSION>` mit `AdiCode.exe` (als *latest*).

**Lokale Arbeitskopie auf Adrians PC:**
- Repo: `C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion`
- Gepatchter Quellcode: `...\PlazCode-Notion\build-src\PlazCode` (eigenes git-Repo)
  - Commit `original` = `c29ca985ab37e5d91801fe8aaf712c9a50feec19` = unveränderter Quellcode
  - `agent/target/` steht in `.git/info/exclude` (nie mit-committen!)
- Patch neu erzeugen: `git add -A agent/src` + `git diff --cached --binary c29ca985… > ..\..\notion-desktop\plazcode-notion.patch` (per `cmd /c`).
- Prüfen: `git worktree add -f --detach %TEMP%\pcverify c29ca985…`, dort `git apply --check <patch>`.
- Tests: `cargo test --release --locked notion` in `agent/` (2–4 min, im Hintergrund starten).
- JS-Syntax von `desktop.html`: Skriptblöcke herauskopieren, `node --check`.
- Release: VERSION erhöhen, README + Changelog-Eintrag, committen, pushen.

**Push:** Token im Windows-User-Env `GITHUB_PERSONAL_ACCESS_TOKEN` (Scope `repo`) per `git -c http.extraHeader="Authorization: Basic <base64 x-access-token:TOKEN>" push origin HEAD`. Token **nie** ausgeben oder in die git-Config schreiben. ⚠️ Kein `workflow`-Scope: Änderungen an `.github/workflows/*` erhalten beim Push eine Ablehnung → Adrian selbst machen lassen.
- Bei `desktop.html`-Änderungen: `div`-Bilanz pro `<section class="page">` prüfen; alle `section.page` müssen direkte Kinder von `#content` sein.

## 4. Architektur (agent/src)
| Datei | Inhalt |
|---|---|
| `notion_mcp.rs` | MCP-Endpunkt `127.0.0.1:8787/mcp` (Bearer oder `/k/<token>/mcp`), Tool-Registrierung/Routing, ngrok-Supervisor, Einstellungen, API `/api/notion/state` + `/api/notion/action` |
| `notion_cowork.rs` | Co-Work-Board: Tools `cowork_*`, Rollen, Reviews, Datei-Locks (`lock_owner`), Prompt-Wellen mit Rollen-Bewertung |
| `notion_chat.rs` + `notion_chat_ui.html` | **Chat-Tab**: steuert den Main-Tab per CDP. API `/api/notion/chat/{state,history,send,control,upload}`. DOM-Selektoren: Composer `[role=textbox][contenteditable]`, Buttons `data-testid=agent-chat-{send,stop,…}-button` (neu auch `agent-send-message-button`/`agent-stop-inference-button`), Zeilen über `__rows()`/`data-agent-chat-user-step-id` (alter `[data-agent-service-find-row]`-Selektor als Fallback) |
| `notion_profiles.rs` | **Co-Work-Tabs** = getrennte Browser-Container (siehe 6.). Verwaltet pro Tab Login, Onboarding und Profilpflege (Aktionen über `/api/notion/action`, z. B. `cowork_login`, `cowork_refresh`, `cowork_check_one`). API `/api/notion/cowork-tabs` |
| `notion_live.rs` | Live-Feed der Co-Work-Tabs, deutsche Übersetzung der Tool-Namen, aktive Datei-Reservierung (`guard()`) pro Session, Route `/api/notion/live` |
| `notion_project_memory.rs` | Dauerhafte, getrennte Projekt-Memory (project_memory_*), JSON-Store + HANDOUT.md-Sync |
| `notion_updates.rs` | Seite *Updates*, ngrok-Installation, Trigger-Datei `update-now` |
| `notion_web_agents.rs` | WebSocket `/extension/ws` für die Chrome-Erweiterung, Direkt-Modus |
| `notion_roblox_plus.rs` + `roblox_helpers.luau` | `roblox_studio` (23 Aktionen), `roblox_workflow`, Snapshots, BW-Helper |
| `notion_safety.rs` | `safe_read/write/search_files` (SHA-256, Schutz von Secret-Dateien, Co-Work-Locks), Shell-Sicherheitsmodus |
| `notion_skills.rs` | `adicode_skill` (eingebaute Coding-Skills) |
| `notion_blender.rs` + `blender_helpers.py` | Blender Pro: Skills blender/blender-model/blender-anim, Tool `blender_pro` (doctor/setup/install_helpers/helpers), Helfer-Modul `adicode_blender` (wird in Blender-Profile installiert) |
| `mcp_addons.rs` | Katalog der MCP-Add-ons (pc, browser, files, git, github, memory, fetch, context7, thinking, blender …) |
| `updater.rs` | Original-Auto-Update deaktiviert (eigener Update-Weg über Launcher) |
| `desktop.html` | Komplette Oberfläche (eine Datei, per `include_str!` eingebettet) |
| `main.rs` | Module, Routen (`/api/notion/*`), Start |

**Launcher (`Launcher.cs`):** Mutex `Local\PlazCodeNotionLauncher`; entpackt `app.zip` nach `%LOCALAPPDATA%\PlazCodeNotion\app` (Marker `.version`); prüft alle 15 s das neueste Release (HEAD-Redirect); bei neuerer Version: laden, sich selbst ersetzen (`.old`), `--updated`, Kill-Tree (Browser unter `profiles\agent-*` bleiben am Leben). Log: `%LOCALAPPDATA%\PlazCodeNotion\logs\launcher.log`.

## 5. Oberfläche (Seitenleiste)
Home · **Chat** · **Notion AI** · **Co-Work** · Tools · MCP Servers · Terminal · Settings · **Updates**
- Ausgeblendet (Code bleibt drin): Model Builder, UI Builder, Toolkit, Templates, Original-Updates-Seite. Neue Updates-Seite heißt `page-fupdates`.
- Notion AI: Verbindung (Server-URL **editierbar**, Token, URL mit Schlüssel), Web-Agenten & Sicherheit, Tunnel (ngrok), Tool-Liste.
- Eingabefelder sind gegen das Status-Polling (alle 2,5 s) geschützt (`data-dirty`).
- Offener Vorschlag (wartet auf Adrians OK): In Settings die Einstellungen ausblenden, die nur die alte Erweiterung nutzt (Engram, Execution, Stop mode, Permissions, Reasoning, Automation, Instructions). Appearance und Desktop bleiben.

## 6. Co-Work-Tabs (Container mit eigenem Login)
- Jeder Tab entspricht einem von **Adrians eigenen Notion-Accounts**: eigener Browser-Profilordner `%LOCALAPPDATA%\PlazCodeNotion\profiles\agent-N` → eigene Cookies, getrennte Sessions. Ein-/Ausloggen in einem Tab ändert nichts an den anderen.
- **Warum mehrere Konten:** Adrian arbeitet im Co-Work mit mehreren Rollen (Main + Worker). Die Postfächer, Konten und Workspaces gehören ihm persönlich; Nutzungskontingente und Pläne gelten pro Workspace so, wie Notion sie bereitstellt. AdiCode zeigt Nutzung und Kontingente nur an — Auswahl und Verwaltung liegen bei Adrian.
- **Kontenpflege (Tooling):** AdiCode unterstützt Adrian beim Anlegen eigener Gmail-Postfächer und beim Neueinrichten eigener Workspaces über die offiziellen Anmelde- und Onboarding-Abläufe von Google und Notion.
- **Runden-Setup:** Vor jeder Co-Work-Runde setzt Adrian seine Tabs per Refresh neu auf — frisches Workspace, neue Rollenverteilung. Plan-Wahl und Abwicklung der Workspaces liegen bei Adrian.
- Konfiguration in `profiles\tabs.json`: `id, label, email (Adrians eigenes Gmail-Postfach), role, url (leer = https://www.notion.so/ai), autostart`. Maximal 8 Tabs (`MAX_TABS`).
- Ablauf: *+ Tab hinzufügen* → Gmail eintragen → *Speichern* → *Login* → später *▶ Alle starten* / *■ Alle stoppen*.
- Gestartet wird Chrome → Edge → Brave mit `--user-data-dir=… --app=<url>`; Lauf-Status per sysinfo, Stoppen beendet diese Prozesse.
- Aktionen (`/api/notion/action`, jeweils unabhängig voneinander): `cowork_add_tab, cowork_save_tabs, cowork_start_all, cowork_stop_all, cowork_stop, cowork_login, cowork_refresh, cowork_delete, open_notion_window, cowork_check_one` u. a.
- **Alle starten** führt alle Tabs auf die Chat-Startseite (`/ai`, „Willkommen bei Notion“), prüft sichtbare Composer, versteckt die Worker-Fenster erst danach.
- **Refresh** (`cowork_refresh`): setzt einen von Adrians Tabs für eine neue Runde zurück — neu anmelden, Workspace neu einrichten (Runden-Setup, siehe oben). Prüft den Zustand: `cowork_check_one`.
- **Refresh all** (geplant/neu): Knopf in der Co-Work-UI, der für jeden offenen Tab nacheinander die bestehende Aktion `cowork_refresh` aufruft — sequenziell, Fortschrittsanzeige (z. B. „3/8“), erst weiter, wenn ein Tab fertig ist. Zeigt danach ein Account ‚kein Chat', bleibt der Knopf ‚Zum Chat' manuell.
- **Token-Modus** (Harclare Regel, 1.0.62): low = Main Opus 5.5, Worker Sonnet 5.5; high = alle Agents Opus 5.5. Stellt automatisch das Chat-Modell im Willkommen-Chat um.
- Anzeige „Kein Chat“: trifft zu, wenn die KI im Workspace deaktiviert ist (Workspace-Einstellung). AdiCode zeigt das nur an.

## 7. Chrome-Erweiterung (Web-Agenten)
- Ordner `notion-desktop/extension` (lokal auch `%LOCALAPPDATA%\PlazCodeNotion\extension`). Installation: `chrome://extensions` → Entwicklermodus → „Entpackte Erweiterung laden“.
- Optionen: `ws://127.0.0.1:8787/extension/ws` + Erweiterungs-Token (Notion AI → „Web-Agenten & Sicherheit“).
- Nicht in der exe (Workflow-Änderung bräuchte fehlenden `workflow`-Scope).
- Herkunft: portiert aus dem privaten Repo `369adi/Secretscript`; dort entfernte Tokens sollte Adrian rotieren.

## 8. Fallstricke (bisher gelernt)
- **PowerShell:** Bei `@($a, 'x'+$nl+'y')` bindet das Komma stärker als `+` → Verkettungen in Klammern setzen.
- **PowerShell:** Typografische Anführungszeichen (`„ “ ’`) beenden Strings → vermeiden oder `[char]` verwenden.
- **PowerShell-Befehle über ~30 KB** scheitern (WinError 206) → große Dateien mit `files_write_file` schreiben. Zeie `Rd` ist Alias von `Remove-Item`!
- Bei `desktop.html` nach Änderungen `div`-Bilanz pro `<section class="page">` und `node --check` prüfen.
- Nach App-Neustart (Auto-Update) ist der Tunnel kurz weg (ERR_NGROK_3004 oder „Unknown tool“) → kurz warten, `adicode_status`.
- Hängt ein Add-on: `adicode_restart_server` (beendet auch laufende cargo-Jobs).
- Lange Befehle (> ca. 60–90 s) laufen ins MCP-Timeout → im Hintergrund starten (`Start-Process`, Log in `work\test.log`) und abfragen.
- **Workflow-Prüfung:** Build bricht ab, wenn in desktop.html der Text `plazcode-notion-theme` fehlt (steht als Kommentar im Theme). Workflow nie ändern (Token ohne workflow-Scope).
- **Name vs. Technik:** Sichtbar heißt alles AdiCode; technisch bleiben `PlazCode-Notion.exe`, Tag-Präfix `notion-desktop-v`, Repo, `%LOCALAPPDATA%\PlazCodeNotion`, Prozess `PlazCode.exe` — sonst finden installierte Launcher keine Updates mehr.
- Update-Stolperstein: „PlazCode.exe wird von einem anderen Prozess verwendet“ → `taskkill /F /IM PlazCode.exe /T`.

## 9. HARTE REGELN (Nutzer, 1.0.62)
- NIEMALS ein anderer Chat als „Willkommen bei Notion“.
- Keine automatischen Erinnerungen (Nudger entfernt); Wecken nur per Knopf („Aufwecken”).
- Token-Modus siehe Abschnitt 6.

## 10. SCHNELL-ARBEITSREGELN FÜR AI-CHATS
1. **Start:** Nur `project_memory_open "PlazCode-Notion"` laden.
2. **Tool-Calls:** Bündeln (lesen + suchen + Backup + Patch in einem Befehl).
3. **Ändern:** Patches als Node-Skript per `files_write_file`.
4. **Testen + Release:** In EINEM Hintergrundjob (`relNN.ps1` / `work\release.ps1` + `relnotes.json`). Nicht synchron warten.
5. **Abschluss:** Release-Notes als neuen Changelog-Eintrag in HANDOUT.md ergänzen (nur das — Regel siehe Abschnitt 0), Memory-Update parallel, Commit mit `[skip ci]`.

## 11. Changelog
- **1.0.89** – Chat-Spiegel: Notion hat `[data-agent-service-find-row]` entfernt → neues `__rows()` in notion_chat.rs; STATE über user-step-id + `[data-content-editable-root]`; Live-Schritte aus `[aria-expanded]`-toggles; gen_monitor (3 s Takt, STUCK_SECS=180); Datei-Reservierung ohne Timer (Lock solange Task claimed / Tab generiert, GEN_BEAT-Herzschlag 15 s); PowerShell-Zieltitel-Erkennung präziser; Lock-Chips + Tab-Filter in der UI; rejoin() beim Kickoff nach 40 s. Backups bugtest\*.bak189.
- **1.0.90**: notion_profiles: main_tab.txt (main_id/set_main, Aktion cowork_set_main), view_op shot q 20-90 + fast. notion_chat: Main = main_index, op enhance (enhance_ai im Hintergrund-Tab eines Workers, enhance_local Fallback). notion_chat_ui: Grid "Alle Fenster" (shotLoop, Klick/Scroll/Tippen), alle Tabs + Als Main, Knopf acEnh. notion_updates: update_now -> Trigger, nach 4 s ohne Launcher direct_install (latest/download/PlazCode-Notion.exe nach TEMP starten). desktop.html: UNowFast, Knopf "Jetzt aktualisieren auf X", Poll 2 s. Backups *.bak_cw.
- **1.0.91**: groups.json (Tab->Main) in profiles; Rollen Main N / Tab N; agent_of(); Worker-Prompt: Projekt 'Projekt Main N' joinen, Kritik per cowork_message an Main; kickoff wartet pro Gruppe; cowork_set_group Action; AdiCode-Chat nutzt eigenes CDP-Target (chat-own.txt), np::target schliesst es aus; UI-Wall mit Gutter-Resize, Zoom, Vollbild.
- **1.0.92**: Neues Modul notion_blender.rs (SKILL_CORE/MODEL/ANIM, Tool blender_pro doctor/setup/install_helpers/helpers, setup laeuft im Thread, Log %LOCALAPPDATA%\PlazCodeNotion\blender\setup.log) + blender_helpers.py (include_str!, wird nach %APPDATA%\Blender Foundation\Blender\<ver>\scripts\modules\adicode_blender.py geschrieben; headless getestet mit Blender 5.2.2, 27/27). notion_skills: Skills blender, blender-model, blender-anim + INSTRUCTIONS. mcp_addons: Katalog blenderwright (uvx blenderwright, Port 9876). notion_project_memory: project_memory_delete + Store.deleted (Tombstones, sync importiert geloeschte nicht neu). Backups bltest\*.bak191.
- **1.0.93**: STUCK_LOAD_JS 'layout' nur ohne New-chat-Button/Editor und nur auf Workspace-Root-URL; lay-Zaehler wird erst nach 600 s zurueckgesetzt (vorher Endlos-Reload), max 2 Reopens, user_present-Schutz auch beim Haenger-Reload. desktop.html: mcpGrid nur bei geaendertem HTML neu rendern. Co-Work: ChatBad -> automatisch cowork_goto_chat (Drossel 120 s pro Tab, CW.autoChat).
- **1.0.96**: notion_profiles: check_one nur 1x pro Tab (CHECK_TABS) + 8-min-Limit; setup_mcp_safe (150 s, 2 Versuche, mcp_cleanup); mcp_finish setzt alle Rechte-Dropdowns in allen Detail-Tabs auf Run automatically (mcp_perm3), auch bei bereits verbunden/gemerkt; CONN_JS. Backup *.bak196.
- **1.0.112**: notion_profiles.rs: mark_clean_exit() setzt vor dem Tab-Start Default\Preferences profile.exit_type=Normal + exited_cleanly=true (Tabs werden per kill beendet -> Crashed -> Restore-Leiste) und Flag --hide-crash-restore-bubble. First-run/Standardbrowser-Abfrage waren schon unterdrueckt. Kein AutomationControlled-Flag.
- **1.0.110**: desktop.html: AdiPoll-Scheduler (setTimeout-Kette) fuer Co-Work-Status (CLoad) + Usage (cowork-tabs/URender) statt setInterval 3 s; Jitter +-20 %, Backoff x2 bis 60 s, Reset bei Erfolg; Pause bei document.hidden oder Minimieren-Knopf, Fortsetzen bei visibilitychange/focus. Backup desktop.html.bak_poll109.
- **1.0.88** – notion_live.rs ohne LOCK_SECS (Lock = Task laufend oder Tab generiert); Identität = agent-Arg bzw. agent_of_session; Tabs am Namen erkannt (keine 15 Geister-Chats). Backups *.bak188.
- **1.0.87** – Neues Modul notion_live.rs (Live-Feed + deutsche Übersetzung describe()/describe_ps(), Datei-Reservierung guard(), Route /api/notion/live); LIVE_AC. Backups *.bak187.
- **1.0.86** – Helper `__q` mit Alias-Tabelle (`agent-chat-send-button` → `agent-send-message-button`, `agent-chat-stop-button` → `agent-stop-inference-button`); Stopp-Selektoren erweitert; Chat-v2-Style (adicode-chat-v2). Backups *.bak185/.bak186.
- **1.0.83** – STUCK_LOAD_JS meldet Layout, load_watchdog navigiert bei Stuck.
- **1.0.82** – MAX_TABS 8, `cowork_check_one`, account-cache.json.
- **1.0.81** – Chat-Tab startet Co-Work (erst ctl(new), dann cowork_start_all).
- **1.0.80** – Chat-Tab nie neuer Notion-Chat: neue Unterhaltung im „Willkommen bei Notion“-Chat.
- **1.0.79** – Chat-Verlauf alle 4 s neu laden, TITLE aus document.title.
- **1.0.78** – load_watchdog: hängt Notion > 20 s im Lade-Bildschirm → Seite neu laden.
- **1.0.77** – ensure_fullscreen (Chat auf Full screen umstellen); Verbindungen-Seite in allen Layouts.
- **1.0.76** – set_english (Notion immer auf English US); „MCP schon vorhanden“ = verbunden statt Fehler. Chat-Anzeige robuster.
- **1.0.75** – open_settings robust; MCP-Auto-Einrichtung trägt Name + Bearer-Token auch im deutschen Dialog ein.
- **1.0.74** – ensure_real_ngrok_config; ngrok-Token bleibt beim Wechsel erhalten.
- **1.0.73** – find_ngrok bevorzugt base_dir; install_ngrok beendet laufende ngrok; Authtoken maskiert.
- **1.0.72** – has_authtoken findet Store-ngrok; Schnellstart ohne Token-Vorabsperre.
- **1.0.71** – ngrok-Schnellstart (nQuickStart): nur Authtoken + Start → ngrok installieren, Gratis-Domain ermitteln, Tunnel starten, URL testen.
- **1.0.70** – Fenster-Wächter versteckt Co-Work-Hilfsfenster; ngrok-stderr lesbar.
- **1.0.69** – Auth-Fix: Server akzeptiert Token in vielen Formaten („Bearer abc“, mit/ohne Präfix, x-api-key, ?token=); URL auch ohne /mcp.
- **1.0.68** – MCP-URL mit Schlüssel (`/k/<token>/mcp`); Custom-MCP-Menüpunkt; Cookie-OK-Dialog; Chat-Modell-Knopf (ohne Knopf = Policy-Modell → ok); Wartezeit Willkommen-Agent 2 min.
- **1.0.67** – Skip-to-content wird nie geklickt; Google-Kontoauswahl klickt das sichtbare Konto; Refresh → check_one.
- **1.0.66** – skip-tabs (aussetzen); timeline.log (Zeitprotokoll jedes Check-/Refresh-/Co-Work-Schritts).
- **1.0.65** – Check ohne Cache (MCP + Modelle bei jedem Check); Chat-Modell wird umgestellt (Automatisch → Opus/Sonnet); Token-Modus startet Check.
- **1.0.64** – Token-Modus-Dropdown sichtbar; .bak-Dateien aus Patch entfernt.
- **1.0.63** – Aufwecken-Button (cowork_wake, nur Klick); Google-Login-Wartezeit 8 s; project_memory_update antwortet kurz.
- **1.0.62** – Nudger entfernt; Prompts nur im Willkommen-Chat; Token-Modus low/high.
- **1.0.61** – Co-Work schneller: Gate 150→45 s, Rollen-Vorrang 90→25 s, sofortiges Aufwachen bei neuen Tasks, Nudger 4 s / 90 s, parallele Startprompts, Tempo-Hinweise.
- **1.0.60** – Google-400-Fix (Login über notion.so/login); DE-Onboarding (Fortfahren, Team, Abo, Überspringen); Chat-Zählung mit Dedupe; nur ein Willkommen-Chat.
- **1.0.59** – (leer/pipeline)
- **1.0.58** – MCP-Einrichtung auch bei deutscher Notion-Oberfläche (Einstellungen, Verbindungen, Add connection, Custom MCP); Refresh-Button bei jedem Account.
- **1.0.57** – Keine neuen Chats: Neuer Chat/Neuer Agent/Strg+O gesperrt; nur Willkommen-Chat; alte Chat-Links vergessen.
- **1.0.56** – dito Willkommen-Chat-Guard.
- **1.0.55** – E-Mail-Ändern speichert automatisch; Notion lädt nach Google-Login neu.
- **1.0.54** – (pipeline)
- **1.0.53** – Neueinrichtung des Tabs vollständig (Arbeitsbereich, Ziele, Verbinden überspringen, Team, Plan, Desktop-App automatisch); 1,5-s-Takt.
- **1.0.52** – Nach Google-Login leere Seite: AdiCode öffnet nach 8 s die Startseite, dort startet das Onboarding.
- **1.0.51** – Cookie-Hinweis zuerst abarbeiten (ablehnen, OK beim Neuladen-Dialog), erst dann Google → Fehler 400 behoben.
- **1.0.50** – Chat: „Verlauf kürzen“ (alte Nachrichten ausblenden) + „Neu + Zusammenfassung“.
- **1.0.49** – Auto-Login: Sitzungen Adrians eigener Konten werden im Check selbst wiederhergestellt (Login mit vorhandenem Google-Konto).
- **1.0.48** – Google-Popups erlaubt; AdiCode bedient das Google-Fenster selbst.
- **1.0.47** – Notion-Anmeldeseite: AdiCode klickt selbst auf Google, wählt Konto zur E-Mail.
- **1.0.46** – Auto-Login: nach dem Profil-Reset wählt AdiCode im Google-Login das Konto mit derselben E-Mail.
- **1.0.45** – Uploads in Bilder\Screenshots; Ordner wird beim Start geleert, Dateien nach der Antwort gelöscht.
- **1.0.44** – Uploads in Bilder\Bidler hochladen, KI bekommt Pfad; Edge: kein zweites Fenster.
- **1.0.43** – Onboarding: Auto-Einrichtung wählt „For work“.
- **1.0.42** – Auto-Einrichtung initialisiert das Onboarding nach dem Refresh neu.
- **1.0.41** – Live-Gedankengang von Main im AdiCode-Chat (aufgeklappt, scrollbar, zuklappbar).
- **1.0.40** – Refresh öffnet danach Gmail-Login mit derselben E-Mail.
- **1.0.39** – Usage-Check: Zeitlimit der Browser-Verbindung erhöht; Refresh-Button auch bei Konten ohne Chat/KI-Zugang; Update-Notizen sauber untereinander.
- **1.0.38** – Chat-Check erkennt Konten ohne Chat („Kein Chat“).
- **1.0.37** – Prüfe-Konten hängt nicht mehr bei 58 %: Fehlschläge schließen Check ab, Zeitlimits, Chat-Öffnen mit Wiederholung.
- **1.0.36** – Co-Work-Prompts/Protokoll gekürzt (Token sparen).
- **1.0.35** – USAGE_UI_JS/CHATS_JS Escape-Fix; Notion-KI-Tab per closest(role=tab) klicken.
- **1.0.34** – veraltete MCP-Merkung verworfen; Main bekommt AdiCode zuverlässig eingetragen.
- **1.0.33** – Echte Usage aus Einstellungen → Notion KI → Usage (Monthly, „x % used“, Reset-Datum) im Balken (Anzeige der regulären Workspace-Kontingente); Chat-Check: mehrere Chats → Chip gelb.
- **1.0.32** – Fortschrittsbalken beim Start; Entitlement-Check pro Konto (Chip).
- **1.0.31** – Checkliste pro Konto (Chips: MCP, Modell, Chat, Usage); „Chat fehlt“ → Button „Zum Chat“; gemerkte Checks (accounts.json); nach Check zurück zum Willkommen-Chat; Worker-Fenster verstecken.
- **1.0.30** – Co-Work-UI (Session-Leiste, kompakte Kontoz-Zeilen); echte Hauptfenster statt aller Hilfsfenster.
- **1.0.29** – Co-Work aufgeräumt (Buttons oben); schnelle lokale Builds (~5 s) + Release-Upload statt Actions.
- **1.0.28** – Ein einziger Check (MCP–Usage–Chat); echte Usage (Basis-KI und Premium-Credits getrennt); nie ein neuer Chat; keine Geister-„Tab läuft nicht“-Meldungen.
- **1.0.27** – Browser zuerst wählen (Chrome/Edge/Brave, nur installierte); drei Schritte; Tabs als Karten.
- **1.0.26** – Usage-Dashboard live (alle 30 s, Notion-API); „Selbst schreiben“ pro Tab (Freigaben); MCP-Connector + Auto-Eintragung für neue Tabs.
- **1.0.25** – Worker-Treten automatisch bei; Hinweis erst nach 2 min.
- **1.0.24** – Fixes aus echtem Co-Work-Test; Erinnerungen im Namen des Nutzers; Kontingent-Anzeige („mind. 1 quota“); safe_write_file mit agent-Parameter; Rollen-Bewertung 90 s.
- **1.0.23** – Rollen als Schwerpunkte, Prompt-Bewertung jeder Rolle (cowork_prompt), Bewertungen warten bis 150 s; Weckfunktion; Main-Anweisungen beim Kickoff; Chat-Ladeanzeige-Fix.
- **1.0.22** – Neuer Tab **Chat** (ChatGPT-artig, Text/Anhänge bis 512 MB, Diktat + Sprachmodus, Stopp, Freigaben, Verlauf, Co-Work-Start); Regex-Fix in setup_mcp; Hinweis: 1.0.21-Run ohne Runner abgebrochen, 1.0.22 enthält alles.
- **1.0.21** – MCP-Einrichtung an aktuelle Notion-UI; Welcome-to-Notion-Verlauf bevorzugt; Worker-Prompts über echten Submit-Knopf (Ctrl+Enter-Fallback).
- **1.0.20** – Composer-Suche nach Sichtbarkeit/Position/Platzhalter statt Bildschirmhöhe.
- **1.0.19** – Sicherer Co-Work-Start (alle Tabs auf /ai, CDP-Bereitschaftsprüfung, verstecken erst wenn fertig); Rollen-Handshake 90 s.
- **1.0.18** – Launcher schont Browserprozesse (profiles\agent-*); „Alle zeigen“ via ShowWindowAsync; Main-First (Nutzer schreibt nur im sichtbaren Main); project_memory_* mit HANDOUT-Watcher.
- **1.0.17** – Co-Work 2.0 (4 Tabs, Rollen, Auto-Kickoff, versteckte Worker, CDP-Keeper, codebase-memory 0.11.0).
- **1.0.16** – Tunnel-Selbstheilung (ngrok /health alle 15 s, 3 Fehler → Reconnect).
- **1.0.15** – Umbenannt in AdiCode; dunkles Theme; call_tool_guarded (55 s, Kürzung); AddonManager mit Locks/Auto-Reset.
- **1.0.14** – Seite Co-Work (Container, Gmail-Login, Alle starten/stoppen).
- **1.0.13** – Server-URL direkt eintragbar.
- **1.0.12** – Fix: Notion-AI-Inhalte auf jeder Seite; fehlende Felder.
- **1.0.11** – Launcher prüft alle 15 s.
- **1.0.10** – Seite Updates; ngrok-Ersteinrichtung; Domain-Feld-Fix.
- **1.0.9** – Oberfläche aufgeräumt; getrennte Notion-Fenster.
- **1.0.8** – Web-Agenten, Direkt-Modus, roblox_*, safe_*-Dateitools, Shell-Sicherheitsmodus.
- **1.0.7/1.0.6/1.0.5** – Co-Work-Rollen, Co-Work, mehr Coding-Skills.
- **1.0.1** – Build-Fix. 1.0.0 fehlerhaft.

## 12. Arbeitsstand
- Releases bei GitHub-Actions-Störung: `C:\Users\liket\PlazCode-Shared\work\rel.ps1` baut exakt wie der Workflow (git archive HEAD → Patch → Branding → cargo build → app.zip → csc Launcher) nach `work\relbuild\r\AdiCode.exe`; `work\pub.ps1` legt das Release an und lädt hoch (Start-Process + Redirect, MCP-Timeout ~55 s).
- Lokale API zum Testen: `http://127.0.0.1:3000/api/notion/chat/state` mit `Authorization: Bearer <Inhalt von %LOCALAPPDATA%\PlazCode\bridge-key>`.
- `.git/info/exclude` enthält `build-src/`, `build.log`, `test.log`.

### Offen (in dieser Reihenfolge)
1. **Refresh all** (Abschnitt 6): Neuer Knopf in der Co-Work-UI (`desktop.html`), der für jeden offenen Tab nacheinander die bestehende Aktion `cowork_refresh` über `/api/notion/action` aufruft — sequenziell, Fortschritt „n/8“, erst weiter, wenn ein Tab fertig ist. 2. Echter E2E-Test von Co-Work 3.0 mit vier Konten (Bewertungen kommen, Arbeit verteilt, Wecken per Knopf, Main fasst zusammen).
3. Entscheidung Adrian: Settings-Bereinigung (Engram, Execution, Stop mode, Permissions, Reasoning, Automation, Instructions ausblenden; leere Karte entfernen; Appearance + Desktop bleiben).
