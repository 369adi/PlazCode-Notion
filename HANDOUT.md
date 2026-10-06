# AdiCode (ehemals PlazCode Notion) – Handout für neue Chats

> Diese Datei ist das Gedächtnis des Projekts. Neue Chats zuerst diese Datei lesen.
> Nach jeder Änderung hier unten im **Changelog** und, falls nötig, in den anderen Abschnitten nachtragen.

## 1. Was ist das?
**PlazCode Notion** ist ein inoffizieller Fork der PlazCode-Desktop-App (stoveez/PlazCode, GPL-3.0, Basis-Version 1.19.36).
Der Fork baut einen MCP-Server in die App ein, damit **Notion AI** über einen ngrok-Tunnel Roblox Studio, den PC, den Browser und weitere MCP-Server steuern kann. Eine Browser-Erweiterung ist dafür nicht nötig.

- Repo: https://github.com/369adi/PlazCode-Notion (öffentlich)
- Releases: Tag `notion-desktop-v<VERSION>` mit einer einzigen Datei `PlazCode-Notion.exe`
- Besitzer: GitHub-User **369adi** (Adrian). Kommunikation auf Deutsch, per „du“.
- Keine Verbindung zu den PlazCode-Autoren oder zu Notion Labs. `LICENSE` und `BRANDING-NOTICE.txt` müssen erhalten bleiben.

## 2. Aufbau des Repos
| Pfad | Zweck |
|---|---|
| `PlazCode-source-1.19.36.zip` | Original-Quellcode (Build-Basis, wird nicht verändert) |
| `PlazCode-1.19.36.zip` | Original-Paket; darin wird nur `PlazCode.exe` ersetzt |
| `notion-desktop/plazcode-notion.patch` | **Alle** Code-Änderungen ggü. dem Original (git diff, inkl. neuer Dateien) |
| `notion-desktop/VERSION` | Release-Version, z. B. `1.0.14` (jede Erhöhung = neues Release + Auto-Update bei allen Nutzern) |
| `notion-desktop/README.md` | Release-Notes (wird als Release-Text verwendet; Abschnitte `## Neu in X.Y.Z` erscheinen in der App unter *Updates*) |
| `notion-desktop/launcher/Launcher.cs` | C#-Launcher = die verteilte `PlazCode-Notion.exe` (enthält app.zip, entpackt, startet, Auto-Update) |
| `notion-desktop/extension/` | Chrome-Erweiterung „PlazCode Notion Web Agents“ (nicht in der exe, siehe 7.) |
| `notion-desktop/branding/` | Icon usw. |
| `.github/workflows/build-notion-desktop.yml` | Build auf `windows-latest`: Patch anwenden → `cargo build --release --locked` → Paket + Launcher → Release |
| `HANDOUT.md` | diese Datei |

## 3. Build- und Release-Ablauf
1. Workflow entpackt `PlazCode-source-1.19.36.zip` nach `build-src/PlazCode` und führt vom Repo-Root aus `git apply --directory=build-src/PlazCode notion-desktop/plazcode-notion.patch` aus.
2. `cargo build --release --locked` in `build-src/PlazCode/agent`.
3. `PlazCode-1.19.36.zip` wird entpackt, `PlazCode.exe` ersetzt, alles als `app.zip` in den Launcher eingebettet (csc, `Launcher.cs` + `Version.cs` aus VERSION).
4. Release `notion-desktop-v<VERSION>` mit `PlazCode-Notion.exe` (als *latest*).

**Lokale Arbeitskopie auf Adrians PC** (über das MCP „Lube Runner“ erreichbar):
- Repo: `C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion`
- Gepatchter Quellcode: `...\PlazCode-Notion\build-src\PlazCode` (eigenes git-Repo, nicht im Haupt-Repo eingecheckt)
  - Commit `original` = `c29ca985ab37e5d91801fe8aaf712c9a50feec19` = unveränderter Quellcode
  - `agent/target/` steht in `.git/info/exclude` (nie mit-committen!)
- Patch neu erzeugen (in `build-src/PlazCode`): `git add -A agent/src` und dann `git diff --cached --binary c29ca985… > ..\..\notion-desktop\plazcode-notion.patch` (per `cmd /c`, damit die Bytes stimmen).
- Prüfen: `git worktree add -f --detach %TEMP%\pcverify c29ca985…`, dort `git apply --check <patch>`, danach Worktree entfernen.
- Tests: `cargo test --release --locked notion` in `agent/` (dauert ca. 2–4 min; im Hintergrund per `Start-Process` starten, Log in `work\test.log`).
- JS-Syntax von `desktop.html` prüfen: die `<script>`-Blöcke herauskopieren und `node --check` ausführen.
- Release: VERSION erhöhen, Notizen in `notion-desktop/README.md` + Changelog hier, committen, pushen. Ergebnis über die GitHub-API unter `actions/runs` prüfen (Build dauert ca. 5–6 min).

**Push:** Den Token im Windows-User-Env `GITHUB_PERSONAL_ACCESS_TOKEN` (Scope `repo`) per `git -c http.extraHeader="Authorization: Basic <base64 x-access-token:TOKEN>" push origin HEAD` verwenden. Den Token **nie** ausgeben oder in die git-Config schreiben.
⚠️ Der Token hat **keinen `workflow`-Scope**: Änderungen an `.github/workflows/*` werden beim Push abgelehnt. Solche Änderungen muss Adrian selbst machen oder den Scope ergänzen.

## 4. Architektur (agent/src)
| Datei | Inhalt |
|---|---|
| `notion_mcp.rs` | MCP-Endpunkt `127.0.0.1:8787/mcp` (Bearer oder `/k/<token>/mcp`), Tool-Registrierung und Routing, ngrok-Supervisor (feste Domain, `--url`/`--domain`), Einstellungen `%LOCALAPPDATA%\PlazCodeNotion\plazcode-notion.json` (token, domain, port, ngrok_*, web_token, safe_shell), API `/api/notion/state` + `/api/notion/action` |
| `notion_cowork.rs` | Co-Work-Board für mehrere Notion-AI-Chats: Tools `cowork_start/join/board/claim/update/add_tasks/message/lock/wait/results/close`, Rollen, Reviews, Datei-Locks (`lock_owner`) |
| `notion_chat.rs` + `notion_chat_ui.html` | **Chat-Tab** (ab 1.0.22): steuert den Main-Tab (Notion AI) per CDP. API `/api/notion/chat/{state,history,send,control,upload}`; Uploads in `%LOCALAPPDATA%\PlazCodeNotion\chat-uploads`, Verlauf `chat-history.json`. UI wird über den Platzhalter `__ADICODE_CHAT_UI__` in desktop.html eingesetzt. DOM-Selektoren: Composer `[role=textbox][contenteditable]`, `input[type=file]`, Buttons `data-testid=agent-chat-{send,stop,dictation,stop-recording,cancel-recording}-button`, Zeilen `[data-agent-service-find-row]` |
| `notion_profiles.rs` | **Co-Work-Tabs** = getrennte Browser-Container (siehe 6.), API `/api/notion/cowork-tabs` |
| `notion_project_memory.rs` | Dauerhafte, getrennte Projekt-Memory (project_memory_*), JSON-Store + automatische HANDOUT.md-Synchronisierung |
| `notion_updates.rs` | Seite *Updates* (GitHub-Releases, Notizen pro Version, Cache 10 min), API `/api/notion/updates`, ngrok-Installation (`bin.equinox.io` + `tar`), Authtoken (`ngrok config add-authtoken`), Trigger-Datei `update-now` |
| `notion_web_agents.rs` | WebSocket `/extension/ws` für die Chrome-Erweiterung, Tools `web_sites/web_agents/web_chat`, Direkt-Modus `/direct-agent/:op` (Freigabe-Panel, 300 s, 50 Aufrufe, sperrt Notion-Schreibzugriffe) |
| `notion_roblox_plus.rs` + `roblox_helpers.luau` | `roblox_studio` (23 Aktionen), `roblox_workflow`, `roblox_project_memory`, automatische Snapshots vor Schreibzugriffen, BW-Helper werden bei `execute_luau` automatisch geladen |
| `notion_safety.rs` | `safe_read_file/safe_write_file/safe_search_files` (SHA-256, Schutz von Secret-Dateien, Co-Work-Locks), Sicherheitsmodus für die Shell (`guard_tool`) |
| `notion_skills.rs` | `plazcode_skill` (eingebaute Coding-Skills) |
| `mcp_addons.rs` | Katalog der MCP-Add-ons (pc = Windows-MCP, browser = Playwright, files, git, github, memory, fetch, context7, thinking, blender …) |
| `updater.rs` | Original-Auto-Update deaktiviert (sonst würde das Original den Fork überschreiben) |
| `desktop.html` | Komplette Oberfläche (eine Datei, per `include_str!` eingebettet) |
| `main.rs` | Module, Routen (`/api/notion/*`), Start |

**Launcher (`Launcher.cs`):** Mutex `Local\PlazCodeNotionLauncher`; entpackt `app.zip` nach `%LOCALAPPDATA%\PlazCodeNotion\app` (Marker `.version`); prüft **alle 15 s** `github.com/369adi/PlazCode-Notion/releases/latest` (HEAD-Redirect, kein API-Limit); ist eine neuere Version da: lädt die exe, ersetzt sich selbst (`.old`), startet mit `--updated`, beendet die alte App (Kill-Tree) und installiert die neue. Ist die Datei `%LOCALAPPDATA%\PlazCodeNotion\update-now` da, prüft er sofort. Log: `%LOCALAPPDATA%\PlazCodeNotion\logs\launcher.log`.

## 5. Oberfläche (Seitenleiste)
Home · **Chat** · **Notion AI** · **Co-Work** · Tools · MCP Servers · Terminal · Settings · **Updates**
- Ausgeblendet (Code bleibt drin, Weiterleitung auf Notion AI): Model Builder, UI Builder, Toolkit, Templates, die Original-Updates-Seite (`page-updates`). Die neue Updates-Seite heißt `page-fupdates`.
- Notion AI: Verbindung (Server-URL **editierbar**, Token, URL mit Schlüssel), Web-Agenten & Sicherheit, Einrichtung in Notion, Tunnel (ngrok installieren, Authtoken, Domain / URL, ngrok.exe-Pfad), Tool-Liste.
- MCP Servers: oben die Karte „Notion-Verbindung (ngrok)“ mit editierbarer Server-URL.
- Eingabefelder sind gegen das Status-Polling (alle 2,5 s) geschützt (`data-dirty`).
- Offener Vorschlag (noch nicht umgesetzt, wartet auf Adrians OK): In den Settings alles ausblenden, was nur die alte Erweiterung nutzt (Engram, Execution, Stop mode, Permissions, Reasoning, Automation, Instructions). Appearance und Desktop bleiben.

## 6. Co-Work-Tabs (Container mit eigenem Gmail-Login)
- Jeder Tab = eigener Browser-Profilordner `%LOCALAPPDATA%\PlazCodeNotion\profiles\agent-N` → eigene Cookies. Aus- oder Einloggen in einem Tab ändert nichts an den anderen.
- Konfiguration in `profiles\tabs.json`: `id, label, email (Gmail), role, url (leer = https://www.notion.so/ai), autostart`.
- Ablauf für Nutzer: *+ Tab hinzufügen* → Gmail eintragen → *Speichern* → *Gmail-Login* (öffnet Google AccountChooser mit dem Konto, danach Notion-Login) → später *▶ Alle starten* (alle Tabs mit Auto) / *■ Alle stoppen*.
- Gestartet wird Chrome → Edge → Brave mit `--user-data-dir=… --app=<url>`. Ob ein Tab läuft, wird per sysinfo an der Kommandozeile erkannt; Stoppen beendet diese Prozesse.
- Aktionen (`/api/notion/action`): `cowork_add_tab, cowork_save_tabs, cowork_start_all, cowork_stop_all, cowork_stop, cowork_login, cowork_delete, open_notion_window`.
- **Alle starten** führt alle vier Tabs einmalig auf die Chat-Startseite (`/ai`, „Willkommen in Notion“) und prüft sichtbare Composer. Erst danach werden die drei Worker versteckt. Nach Main-`cowork_start` bekommen sie feste Rollenprompts; Coder/Reviewer/Tester müssen per Rollen-Handshake auf dem neuesten gemeinsamen Board erscheinen, sonst zeigt AdiCode alle Fenster zur Kontrolle.

## 7. Chrome-Erweiterung (Web-Agenten)
- Ordner `notion-desktop/extension` (lokal auch `%LOCALAPPDATA%\PlazCodeNotion\extension`). Installation: `chrome://extensions` → Entwicklermodus → „Entpackte Erweiterung laden“.
- Optionen: `ws://127.0.0.1:8787/extension/ws` + Erweiterungs-Token (Notion AI → „Web-Agenten & Sicherheit“).
- Nicht in der exe enthalten, weil der Workflow geändert werden müsste (fehlender `workflow`-Scope).
- Herkunft: portiert aus dem privaten Repo `369adi/Secretscript`. Dort wurde `config.json` mit Tokens entfernt; Adrian sollte diese Tokens rotieren.

## 8. Fallstricke (bisher gelernt)
- **PowerShell:** Bei `@($a, 'x'+$nl+'y')` bindet das Komma stärker als `+`. Verkettungen immer in Klammern setzen, sonst fehlen Teile (das hat in 1.0.11 das Layout zerschossen).
- **PowerShell:** Typografische Anführungszeichen (`„ “ ’`) im Befehl beenden Strings. In Skripttexten vermeiden oder `[char]` verwenden.
- **PowerShell-Befehle über ~30 KB** scheitern (WinError 206) → große Dateien mit `files_write_file` schreiben.
- Bei `desktop.html` nach Änderungen die `div`-Bilanz pro `<section class="page">` und `node --check` prüfen. Alle `section.page` müssen direkte Kinder von `#content` sein.
- Startet die App neu (Auto-Update), ist der Tunnel kurz weg (ngrok-Fehler ERR_NGROK_3004 oder „Unknown tool“). Kurz warten, dann `adicode_status` aufrufen.
- Hängt ein Add-on: `adicode_restart_server` (z. B. `pc`). Dabei werden Kindprozesse beendet, also auch laufende `cargo`-Jobs.
- Lange Befehle (> ca. 60–90 s) laufen ins MCP-Timeout → im Hintergrund starten und abfragen.

- **Workflow-Prüfung:** Der Build bricht ab, wenn in desktop.html der Text `plazcode-notion-theme` fehlt. Steht deshalb als Kommentar im AdiCode-Theme – nicht entfernen. Workflow selbst nicht ändern (Token ohne workflow-Scope).
- **Name vs. Technik:** Sichtbar heißt alles AdiCode. Technisch bleiben `PlazCode-Notion.exe` (Release-Asset), Tag-Präfix `notion-desktop-v`, Repo, Datenordner `%LOCALAPPDATA%\PlazCodeNotion` und Prozess `PlazCode.exe` – sonst finden installierte Launcher keine Updates mehr.

## 10. Offene Aufgaben (Stand nach 1.0.25)
**Release B = 1.0.17 „Co-Work 2.0“ ist umgesetzt.**

Als Nächstes Nutzerentscheidung einholen: Sollen in Settings **Engram, Execution, Stop mode, Permissions, Reasoning, Automation und Instructions** entfernt werden? Außerdem die leere Karte oben entfernen. Appearance und Desktop bleiben.
## 11. Changelog
- **1.0.32** - Fortschrittsbalken (STAGE/step in notion_profiles.rs, CProg in desktop.html, state.stage), Trial-Check (TRIAL_JS: getSubscriptionData isSubscribed && !hasPaidNonzero = Trial; Chip in CChk).
- **1.0.31** - Account-Check: ensure_models (MODELS_JS via saveTransactionsFanout, personal_agent_model_policy; Main Opus 5.5 albuquerque-quinn, Worker Sonnet 5.5 achira-donut), Gedaechtnis accounts.json (mem_get/mem_set: mcp_ok, models_ok, welcome_url), close_settings vor prepare_chat, Worker nach Kickoff versteckt, state.hidden. UI: CChk-Chips, Zum Chat (cowork_goto_chat), neu pruefen (cowork_forget), Icon-Labels, Browser-Picker CBrIcon/CBrRender.
- **1.0.30** - Co-Work: Session-Leiste (cwsession) mit Start rechts, Konten als Zeilen (cwrow, CRow neu), Usage-Dashboard-Karte versteckt (cwUsage bleibt fuer JS). Schwarze Fenster gefixt: notion_profiles.rs windows::visit merkt versteckte HWNDs (HIDDEN) und zeigt nur diese bzw. Chrome_WidgetWin_1 mit Titel ohne Owner.
- **1.0.29** - Co-Work-UI kompakt (Start/Stopp im page-head `cwtop`, Werkzeugleiste `cwtool`, Karten-Buttons icon-only via `.lbl`, IDs unveraendert). RELEASE-WEG NEU: `work\fast-release.ps1` baut in `work\fastbuild` (robocopy von build-src + Branding, CARGO_TARGET_DIR=fastbuild\target, Schnellprofil LTO aus/16 CGU/inkrementell/opt 1 -> ~5 s) und laedt `PlazCode-Notion.exe` per GitHub-API als Release `notion-desktop-vX` hoch. Commits IMMER mit `[skip ci]`, sonst baut Actions zusaetzlich. Offen: Usage-% aus Einstellungen > Notion KI > Usage und automatische Modellwahl (Main Opus 5.5, Worker Sonnet 5.5 ueber "Allowed models for Notion Agent") - braucht Live-Mitschnitt mit cdp.js.
- **1.0.28** - Ein Check: `check_one`/`full_check` (parallel je Tab inkl. Main: `wait_cdp` -> `logged_in` -> `setup_mcp` -> `inspect_account` -> `prepare_chat`), Feld `check` im Account-State, Aktion `cowork_check` (= Alias `cowork_check_accounts`, `check_action` startet fehlende Tabs, Sperre `CHECKING`). `kickoff` nutzt `full_check` und versteckt NICHTS mehr (nur `cowork_hide_all`). `prepare_chat` oeffnet nur den Sidebar-Chat "Welcome to/Willkommen bei Notion" (`FIND_WELCOME`, `CHAT_READY`, merkt `welcome_url`), nie New Chat. `USAGE_JS` neu: v1 `getAIUsageEligibility` + V2 `premiumCredits.perSource` -> basic_*/premium_*/plan_type. `launch` loescht `usage_error`. Versionsanzeige: `updater.rs` Feld `notion_version` aus `%LOCALAPPDATA%\PlazCodeNotion\app\.version`, desktop.html ueberschreibt sidebarVersion/titleVersion. Fakten: Custom Agents und Workers kosten Notion-Credits (nicht das normale KI-Kontingent) -> Orchestrierung selbst bauen (1.0.29: kurze Task-Pakete, Ergebnis-Zusammenfassungen, Startprompt nur einmal, weniger Nudges). Offen: Script-Mode-Sperre stammt nicht aus agent/src (vermutlich PlazCode-Extension AGENTSCRIPT, core/config.js).
- **1.0.27** - Browserwahl: `installed_browsers()`/`chosen_browser()` (Datei `profiles\browser`), `browser()` liefert nur noch den gewaehlten Browser, Aktion `cowork_set_browser` (gesperrt solange Tabs laufen), state liefert `browsers`. Neue Aktionen `cowork_hide_all`, `cowork_show_tab`, `cowork_login_view`. Login-Ansicht: `login_view` startet/navigiert den Tab zur Google-Kontoauswahl und versteckt das Fenster; Route `/api/notion/cowork-view` (`view_handler`, ops shot/click/type/key/scroll/back per CDP, `cdp_ws` + `any_target`). Co-Work-UI neu (Karten, SVG-Icons via `data-ic`/`CI()`, Modal `#cwLogin`). Offen: Google kann Logins in ferngesteuerten Fenstern blockieren -> dann "Echtes Fenster zeigen". Idee in Pruefung: Notion Custom Agents fuer Orchestrierung (kosten Notion-Credits).
- **1.0.26** - Usage-Check neu: `USAGE_JS` in notion_profiles.rs fragt `/api/v3/getSpaces` + `getAIUsageEligibility(V2)` im Tab ab (Felder used/limit/exhausted/usage_source/usage_error), `usage_monitor` alle 30 s fuer laufende Tabs, Live-Dashboard `#cwUsage` auf Co-Work. Worker-Chats bedienbar: Aktionen `cowork_takeover`/`cowork_release` (PAUSED-Set, Nudger + Watchdog ueberspringen den Tab, Fenster wird gezeigt, Generierung gestoppt), Nudger schickt nichts, wenn `document.hasFocus()` im Worker. MCP-Connector-Karte auf Notion AI (`nMcpTab`, `nMcpSetup`, `nMcpAuto`), Aktion `cowork_auto_mcp` (Datei `profiles\auto-mcp`), `auto_mcp_watcher` alle 20 s (eingeloggt + kein MCP -> `setup_mcp`, max. alle 5 min pro Tab). Monitore starten ueber `ensure_monitors()` (tabs_handler/resume_background). Offen: API-Feldnamen der Usage-Antwort im echten Konto pruefen (Quelle steht im Dashboard).
- **1.0.25** – cowork_start-Ausgabe: Main bittet den Nutzer nicht mehr, Tabs zu öffnen (Worker-Tabs treten automatisch bei; Hinweis erst nach 2 Min ohne Beitritt).
- **1.0.24** – E2E-Fixes Co-Work: natürlicher Worker-/Main-Prompt (alter wurde von Notion AI als Prompt-Injection abgelehnt; Hinweis Verbindung heißt „asf“/„AdiCode“), Nudge-Texte als „automatische Erinnerung von AdiCode (von mir eingerichtet)“, Kontingent-Erkennung (`QUOTA_JS` in insert_prompt) + Kickoff läuft mit den übrigen Workern weiter, `safe_write_file` Schema mit `agent`, Gate ohne Worker 90 s. E2E-Ergebnis (Projekt p3 Todo-App): Coder+Tester arbeiten parallel, Cross-Reviews und Fix-Runden funktionieren, Rollen-Stealing funktioniert (Coder machte UI/UX- und Tester-Tasks). Reviewer-Konto (mia.delgado261) hat kein KI-Kontingent mehr.
- **1.0.23** – Co-Work 3.0: Rollen = Schwerpunkt + Bewertungs-Brille (`role_ok` mit `specialist_idle`/`STEAL_SECS=90`), Prompt-Wellen (`open_wave`, Assess-Tasks mit `lens()`, Gate `GATE_SECS=150`, neues Tool `cowork_prompt`, `cowork_update tasks=[...]`), Claim-Reihenfolge eigene Bewertung → eigene Rolle → Reviews → ohne Rolle → Rest. Nudger in `notion_profiles.rs` (alle 10 s, `nudge_for(agent)` aus notion_cowork, nur wenn Chat idle, Wiederholung nur bei neuer Signatur/nach 4 min, max. 3×). Main-Instruktion beim Kickoff. Skill `cowork` neu. Chat-Tab: pending-Erkennung per bekannter User-Nachrichten (Fix Dauer-Loading). Releases bei Actions-Störung lokal: `work\rel.ps1` + `work\pub.ps1`.
- **1.0.22** – Neuer Tab **Chat** (ChatGPT-artig) über den Main-Tab: Text, Foto/Video/Datei-Anhänge (Upload-Endpunkt, 512 MB), Notion-Diktat + Sprachmodus (VAD, Vorlesen per speechSynthesis), Stopp, Freigabe-Knöpfe, lokaler Verlauf, Co-Work-Start/Status. Fix: echte 0x08-Zeichen statt `\b` in den Regexen von `setup_mcp` (Erkennung bestehender Verbindung griff nie). Der 1.0.21-Run wurde mangels Runner abgebrochen – Workflow ist ok, einfach neu pushen bzw. Rerun.
- **1.0.21** – Weitere reale Vier-Account-E2E-Fixes: `setup_mcp` folgt nun der aktuellen Notion-UI über Workspace-Menü → Settings → Connections, erkennt bereits installiertes AdiCode, öffnet sonst Discover → Custom MCP, bearbeitet den zweistufigen Dialog und nutzt für Connect vertrauenswürdige CDP-Mausereignisse. Chat-Boot klickt explizit „Welcome to Notion“/„Willkommen bei Notion“; ohne diesen Verlauf wird ein neuer Chat geöffnet, statt ein bestehendes Projektgespräch zu übernehmen. Worker senden über den sichtbaren Submit-Knopf (Ctrl+Enter-Fallback), damit unterschiedliche Enter-Einstellungen nicht blockieren.
- **1.0.20** – Hotfix aus realem Vier-Account-E2E: Worker-Prompts scheiterten in 1.0.19, weil der Composer fälschlich unterhalb von 45 % der Fensterhöhe liegen musste. `insert_prompt` navigiert den vorbereiteten „Willkommen in Notion“-Chat nicht mehr weg, sucht sichtbare Composer robust nach Platzhalter/Position und wartet bis zu 10 s. Der neue Selektor wurde per CDP gegen alle vier echten Account-Layouts geprüft.
- **1.0.19** – Sicherer Co-Work-Boot: Jeder Account wird beim Start einmalig auf `/ai` bzw. Chat/„Willkommen in Notion“ geführt und per CDP auf ein sichtbares Eingabefeld geprüft. Worker werden erst versteckt, wenn alle vier Tabs bereit sind. Nach Main-`cowork_start` senden die Worker feste Rollenprompts; ein 90-s-Handshake prüft, ob Coder, Reviewer und Tester dem neuesten Board wirklich beigetreten sind. Live-Status und Tab-Zeilen zeigen den Fortschritt; bei Chat-, Prompt- oder Rollenfehlern werden alle Fenster sichtbar. Account-Checks erhalten MCP-/Chat-Status statt ihn zu überschreiben.
- **1.0.18** – Co-Work-Browser überleben Launcher-Updates (KillTree überspringt Browser unter `profiles\agent-*`); vor Update-Neustart und normalem App-Beenden werden versteckte Worker-Fenster sichtbar gemacht. „Alle zeigen“ nutzt `ShowWindowAsync(SW_RESTORE)` + Vordergrund. Main-First: Nutzer schreibt nur in der sichtbaren Main-Notion-App; Worker bleiben versteckt und werden erst nach einem neuen `cowork_start` gestartet. Eingebautes `project_memory_*` mit getrennten Projekten, Suche, durablem JSON-Store und 20-s-Watcher für HANDOUT.md unter `PlazCode-Shared\work`; Skill/Instructions erzwingen Load vor Projektarbeit und Save nach Änderungen/Themenwechsel. Codebase-Memory bleibt zusätzlich für Codegraph/Indexierung zuständig.
- **1.0.17** – Co-Work 2.0: max. vier Tabs; Rollen automatisch Main/Coder/Reviewer/Tester; Start nur mit vier E-Mails und Aufgabe; CDP-Auto-Kickoff plus Kopier-Fallback; nur Main sichtbar. Chromium-Hintergrunddrosselung und Native Window Occlusion deaktiviert, zusätzlicher CDP-Keeper für versteckte/unfokussierte Tabs. Account-Prüfung zeigt Usage und erkannte verfügbare Modelle. `codebase-memory-mcp@0.11.0` als standardmäßig aktiviertes Add-on. Co-Work-Antworten auf neue Meldungen gekürzt. Nach Nutzerfeedback: Auto-Kickoff navigiert gezielt zu `/ai` und nutzt nur den sichtbaren unteren AI-Composer; automatische MCP-Einrichtung pro E-Mail-Tab über Settings → Connections, mit Button/Fallback.
- **1.0.16** – Hotfix fuer korrekte, aber zeitweise als „nicht unterstuetzter MCP-Endpunkt“ gemeldete URLs: Der ngrok-Supervisor prueft die oeffentliche `/health`-Route alle 15 Sekunden und startet den Tunnel nach drei Fehlern automatisch neu. Release B wurde dadurch auf 1.0.17 verschoben.
- **1.0.15** – Umbenannt in **AdiCode** (sichtbare Texte, Fenstertitel, Tray, Erweiterung, Launcher-Meldungen), neues Logo (`notion-desktop/branding/make-icon.ps1` erzeugt `plazcode.ico/png`, Dateinamen bleiben wegen Workflow). Tools `plazcode_*` → `adicode_*`. Neues dunkles Theme (`<style id="adicode-theme">` in desktop.html, ersetzt das helle Notion-Theme). Stabilität: `call_tool_guarded` in notion_mcp.rs (eigene Task, Panic-Schutz, 55-s-Timeout, Ausgabe > 120 000 Zeichen gekürzt); `AddonManager` mit eigenem Lock pro Add-on (`mcp_addons::call_shared`), Cache für Tool-Listen, Auto-Reset bei Timeout/Absturz; Co-Work speichert gebündelt im Hintergrund-Thread, max. 300 Nachrichten.
- **1.0.14** – Neue Seite **Co-Work**: mehrere Notion-Tabs als getrennte Container, je ein Gmail-Konto, Speichern, Gmail-Login, Alle starten/stoppen, Status „läuft“. Dieses Handout neu geschrieben.
- **1.0.13** – Server-URL direkt eintragbar (Notion AI + Karte auf MCP Servers).
- **1.0.12** – Fix: Notion-AI-Inhalte erschienen auf jeder Seite; Felder für Authtoken und Domain fehlten.
- **1.0.11** – Launcher prüft alle 15 s statt alle 90 s.
- **1.0.10** – Seite *Updates* (Versionen, Notizen, Jetzt prüfen/aktualisieren, Launcher-Log); ngrok-Ersteinrichtung in der App; Fix: Domain-Feld wurde überschrieben.
- **1.0.9** – Oberfläche aufgeräumt (Model/UI Builder, Toolkit, Templates, Updates raus); erste getrennte Notion-Fenster.
- **1.0.8** – Aus Secretscript übernommen: Web-Agenten + Erweiterung, Direkt-Modus, roblox_studio/workflow/project_memory, Snapshots, BW-Helper, safe_* Dateitools, Shell-Sicherheitsmodus.
- **1.0.7** – Co-Work mit festen Experten-Rollen + Live-Status.
- **1.0.6** – Co-Work (mehrere Notion-AI-Chats an einem Projekt).
- **1.0.5** – mehr Tools und Skills fürs Programmieren.
- **1.0.1** – Build-Fix (`git apply --directory`). 1.0.0 war fehlerhaft.

## 12. Arbeitsstand
- Releases bei GitHub-Actions-Störung: `C:\Users\liket\PlazCode-Shared\work\rel.ps1` (baut exakt wie der Workflow nach `work\relbuild\r\PlazCode-Notion.exe`) und `work\pub.ps1` (Release `notion-desktop-v<VERSION>` anlegen + Asset hochladen). Beide per `Start-Process … -RedirectStandardOutput` starten (MCP-Timeout ~55 s).
- Lokale API zum Testen: `http://127.0.0.1:3000/api/notion/...` mit `Authorization: Bearer <%LOCALAPPDATA%\PlazCode\bridge-key>`.
- Update-Stolperstein: „PlazCode.exe wird von einem anderen Prozess verwendet“ → `taskkill /F /IM PlazCode.exe /T` und neu starten. Idee: Launcher killt vor dem Entpacken alle PlazCode.exe aus dem app-Ordner und versucht es erneut.
- Als Nächstes: echter E2E-Test von Co-Work 3.0 mit vier Konten (Bewertungen, Arbeitsverteilung, Weckrufe, finale Zusammenfassung) und die Werte (Gate 150 s, Steal 90 s, Nudge 10 s/4 min) nach dem Test feinjustieren.