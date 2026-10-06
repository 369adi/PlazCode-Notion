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
- **1.0.35** - USAGE_UI_JS/CHATS_JS Escape-Fix, Notion-KI-Tab per closest(role=tab) klicken, auf Usage-Tab warten.
- **1.0.36**: Prompts/Protokoll gekuerzt (Token sparen)
- **1.0.37**: Fix Pruefe-Konten-Haenger bei 58 %
- **1.0.38**: Refresh-Button (Konto loeschen nach Bestaetigung), Kein-Chat-Erkennung
- **1.0.39**: Usage-Check Fix (CDP-Timeout), Refresh-Button bei Kein-Chat/ohne KI
- **1.0.40**: Refresh oeffnet danach Gmail-Login mit derselben E-Mail
- **1.0.40**: Refresh oeffnet danach Gmail-Login mit derselben E-Mail
- **1.0.40**: Refresh oeffnet danach Gmail-Login mit derselben E-Mail
- **1.0.41**: Live-Gedankengang von Main im AdiCode-Chat
- **1.0.42**: Auto-Einrichtung nach Refresh (Onboarding, Trial starten und kuendigen)
- **1.0.43**: Onboarding waehlt For work (Onboarding, Trial starten und kuendigen)
- **1.0.44**: Uploads in Bilder-Ordner, keine Extra-Fenster (Onboarding, Trial starten und kuendigen)
- **1.0.45**: Uploads in Screenshots-Ordner, automatisch loeschen (Onboarding, Trial starten und kuendigen)
- **1.0.46**: Google-Login automatisch (Onboarding, Trial starten und kuendigen)
- **1.0.47**: Notion-Login klickt Google (Onboarding, Trial starten und kuendigen)
- **1.0.48**: Popup-Blocker aus fuer Google-Login (Onboarding, Trial starten und kuendigen)
- **1.0.49**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.50**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.51**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.52**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.53**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.54**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.55**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.56**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.57**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.58**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.59**: Auto-Login im Check (Onboarding, Trial starten und kuendigen)
- **1.0.60**: Google-400-Fix (Login ueber notion.so/login), DE-Onboarding, Klick-Entprellung (ein Willkommen-Chat), Chat-Zaehlung
- **1.0.61**: Co-Work schneller (GATE 45 s, STEAL 25 s, wait weckt bei neuen freien Tasks, Nudger 4 s/90 s, parallele Worker-Prompts, Tempo-Hinweise in Prompts)
- **1.0.62**: Nudger entfernt, Willkommen-Chat-Guard in insert_prompt, Token-Modus low/high (Datei profiles/token-mode)
- **1.0.63**: Aufwecken-Button (cowork_wake), Popup-Timeout 8 s, kurze project_memory_update-Antwort
- **1.0.64**: Token-Modus-Dropdown sichtbar (eigene Klasse cwtok statt cwbrow), .bak-Dateien aus Patch entfernt
- **1.0.65**: Check ohne Cache (MCP/Modell), Chat-Modell-Knopf (CHAT_MODEL_JS), Token-Modus -> Check, Trial-Fix + Auto-Probe-Abo (TRIAL_START_JS), MCP-Erkennung nur im Dialog
- **1.0.66**: skip-tabs (aussetzen), timeline.log
- **1.0.67**: Refresh->check_one, Skip-to-content-Fix, Google-Chooser sichtbar, Trial vor MCP, Willkommen-Chat entsperren, Trial kuendigen
- **1.0.68**: MCP-URL mit Schluessel (/k/<token>/mcp) fuer fremde PCs, Header ODER Pfad-Token, Custom-MCP-Menuepunkt, Connect im Dialog, Cookie-OK + Chat ohne Modell-Knopf ok
- **1.0.69**: Bearer-Fix fremde PCs (clean_tok: doppeltes Bearer/Quotes/Whitespace, x-api-key, ?token=, Root-URL + /k/<token> ohne /mcp), OPTIONS/CORS, Diagnose abgelehnter Anfragen (UI-Zeile Abgelehnt, logs\mcp-auth.log)
- **1.0.70**: Fenster-Waechter fuer versteckte Worker (alle Chrome_WidgetWin-Fenster der Tab-Prozesse, 0,9 s; Zeigen nur Hauptfenster), ngrok-stderr wird gesammelt (ERROR:-Mehrzeiler)
- **1.0.71**: ngrok-Schnellstart (Aktion ngrok_quick_start, quick_start/detect_domain in notion_mcp.rs: install -> add-authtoken -> alle ngrok beenden -> Domain aus ngrok-Log ohne --url -> Tunnel -> POST /mcp muss 401 ohne Ngrok-Error-Code geben), Knopf nQuickStart
- **1.0.72**: has_authtoken findet Store-ngrok (Packages\ngrok.ngrok_*\LocalCache\Local\ngrok\ngrok.yml), Schnellstart ohne Token-Vorabsperre
- **1.0.73**: find_ngrok bevorzugt base_dir\ngrok.exe, install_ngrok killt alle ngrok + setzt ngrok_path, quick_start nutzt eigenes ngrok + Reinstall bei zu alt, authtoken_mask im State/UI, Install-Knopf mit Feedback
- **1.0.74**: ensure_real_ngrok_config (Store-ngrok.yml -> %LOCALAPPDATA%\ngrok)
- **1.0.75**: open_settings (Menue oben/unten, trusted CDP-Klicks, Strg+, Fallback), FILL_JS Name+Bearer-Token, wait_ready/poll_js
- **1.0.76**: set_english (LANG_*_JS, Aktion cowork_set_english, in check_one+setup_mcp), ALREADY_JS/mcp_already, SETTINGS_ITEM_JS nur im Menue, CHATS_JS nur Seitenleiste
- **1.0.77**: ensure_fullscreen (FS_BTN_JS/FS_ITEM_JS) in prepare_chat, welcome_url nur mit /chat, open-Skript: Connected/Manage/All connections, Add-Knopf auch Configure (click+Pointer-Events), mehr Custom-MCP-Bezeichnungen
- **1.0.78**: load_watchdog (STUCK_LOAD_JS, Reload nach 20 s, 3. Mal Navigate) + COOKIE_JS alle 5 s, gestartet in launch(); READY_JS braucht Text; english_once/mem lang_en; Check-Schritte gemerkt (mcp_ok2, models_ok, chat_model, usage_ui_ok); mem_forget behaelt lang_en
- **1.0.79**: Chat-Verlauf alle 4 s neu laden, Titel aus document.title (ntitle) bzw. 7 Woerter, STATE erkennt DE-Labels, cowork_main_context in send() bei neuem Chat, hide_all_tabs nach gruenem Check/Kickoff (set_main_hidden)
- **1.0.34** - Memory-Schluessel mcp_ok -> mcp_ok2 (alte falsche Merkung ungueltig).
- **1.0.34** - Memory-Schluessel mcp_ok -> mcp_ok2 (alte falsche Merkung ungueltig).
- **1.0.33** - USAGE_UI_JS liest Settings>Notion KI>Usage (x% used, Resets on) -> ui_pct/usage; CHATS_JS zaehlt Sidebar-Chats (chats_multi gelb); inspect_account behaelt ui_pct.
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
## STAND 1.0.60 (in Arbeit, noch NICHT released) - Login-/Onboarding-Bugs
- Ursache Google-400: Google lehnt `accounts.google.com/AccountChooser?...&continue=https://www.notion.so/login` IMMER mit 400 ab (continue auf Nicht-Google-Domain). Alle Google-URL-Varianten getestet (AccountChooser, ServiceLogin, v3/signin/identifier) -> 400. google_close schloss dann das einzige Fenster = 'Tab schliesst sich'.
- Fix im Code (build-src ...\agent\src\notion_profiles.rs, Backup .bak159): login_url() = https://www.notion.so/login (Notion-Button 'Mit Google fortfahren' -> Popup funktioniert, im Test an Tab 4 erfolgreich); cowork_login nutzt login_url; google_close schliesst nie das letzte Fenster (navigiert dann zu notion.so/login); google-error-Regex enger.
- Onboarding DE: Regexe ergaenzt (Fortfahren, 'Wer ist sonst noch in deinem Team', 'Waehle dein Abo', 'Vorerst ueberspringen'); Klick-Entprellung (window.__clk, 20 s je Schritt/Regex) gegen Mehrfach-Klicks -> vermutete Ursache der mehreren 'Willkommen bei Notion'-Chats (Beweis nur mit frischem Konto moeglich).
- CHATS_JS zaehlte gleichnamige Chats als 1 -> jetzt Dedupe nach Position (Test: 3 Chats -> total 3 / welcome 3).
- Testwerkzeuge: work\bugtest (harness2.js, logout.js, urltest.js, runjs.js, patch160*.js). Dateien haben gemischte Zeilenenden (CRLF+LF) -> beim Patchen nach '\n' splitten.
- OFFEN: Test mit frischem Konto (Refresh Tab 4 adiabi26444 = Konto loeschen -> braucht OK des Nutzers), Build/Release 1.0.60 nur nach OK.

- 1.0.60 RELEASED (Commit e52f5fb, fast-release). Test: Tab4-Konto per REFRESH_JS geloescht, Login ueber notion.so/login + Onboarding mit harness -> genau 1 Willkommen-Chat (vorher 3). CHATS_JS oeffnet jetzt eingeklappte Seitenleiste.

## Arbeitsregeln vom Nutzer (Mia/Adrian)
- Nicht unnoetig warten: waehrend etwas laeuft/baut/wartet IMMER andere Aufgaben erledigen (parallele Tool-Calls, Hintergrundjobs per Start-Process).
- Co-Work-Geschwindigkeit (meine + andere Co-Work-Chats) deutlich erhoehen, aber Leistung/Qualitaet behalten. (Aufgabe offen seit 1.0.60)

## HARTE REGELN (Nutzer, 1.0.62)
- NIEMALS ein anderer Chat als 'Willkommen bei Notion' (auch nicht von AdiCode automatisch). Mehrere Willkommen-Chats sind erlaubt: jeder ist ein eigenes Kontextfeld. Andere Chats machen Konten kaputt. Kontextlimit wird erfasst (Opus 5.5 sehr hoch).
- Automatische Erinnerungen (Nudger) sind ENTFERNT (1.0.62) - nicht wieder einbauen. AdiCode schreibt nie von selbst in Chats. insert_prompt hat einen Guard (nur Willkommen-Chat).
- Token-Modus (Co-Work-Seite, Datei profiles\token-mode, Aktion cowork_token_mode): low = Main Opus 5.5 + Worker Sonnet 5.5; high = alle Agents Opus 5.5 (setzt personal_agent_model_policy je Konto).
- Mehrfach-Willkommen-Chats entstehen durch wiederholte Onboarding-Klicks (Debounce in 1.0.60). 'token_tb%' wurde in 133 geladenen Notion-JS-Dateien NICHT gefunden (Suche token_tb/tbToken/welcomeChat).

## STAND 1.0.68 (released, Commit 7cef7ba, installiert)
### Einrichtung auf einem fremden PC (z. B. Freund) - so klappt es
1. `PlazCode-Notion.exe` vom neuesten Release starten (aktualisiert sich alle 15 s selbst).
2. Seite **Notion AI** > Tunnel: ngrok installieren, **eigenen** ngrok-Authtoken eintragen, **eigene feste ngrok-Domain** eintragen (dashboard.ngrok.com > Domains, kostenlos 1 Domain). Jeder PC braucht sein EIGENES ngrok-Konto + Domain (eine Domain kann nur an einem Tunnel haengen).
3. Tunnel muss "online" zeigen. Token wird pro PC zufaellig erzeugt (`%LOCALAPPDATA%\PlazCodeNotion\plazcode-notion.json`).
4. In Notion: Einstellungen > Verbindungen > Add connection > **Custom MCP connection** > als URL die **"URL mit Schluessel"** eintragen (`https://<domain>/k/<token>/mcp`) > Verbinden. Kein Bearer-Feld noetig. Braucht Business-Plan/Trial.
### Ursache des Fehlers bei anderen (gefixt in 1.0.68)
- Auto-Einrichtung trug nur `https://<domain>/mcp` ein und erwartete danach ein Bearer-/Passwort-Feld. Neue Notion-Dialoge ("Individueller MCP-Server") haben nur das URL-Feld -> Notion verband ohne Token -> 401 -> "klappt nicht". Bei Adrian lief es nur, weil seine Verbindung schon von frueher bestand.
- Fix: `cowork_connection()` liefert jetzt die URL mit Schluessel; `authorized()` akzeptiert Header-Token ODER Pfad-Token (ein falscher Header blockiert den Pfad-Token nicht mehr). Getestet: key-url lokal 200, key-url+falscher Header 200, Bearer 200, ohne Token 401, key-url oeffentlich ueber ngrok 200.
- setup_mcp: "Custom MCP connection" ist ein `[role=menuitem]` (nicht button) -> wird jetzt gefunden (Retry 10x); Verbinden wird nur im obersten Dialog geklickt (vorher erwischte es die Verbinden-Knoepfe der Kartenliste); Verify prueft AdiCode im Dialog.
- CHAT_MODEL_JS: Cookie-OK-Dialog ("Notion wird neu geladen ... Cookie-Einstellungen") wird bestaetigt (der Reject-Klick traf vorher den Cookie-Banner); ohne Modell-Knopf (Policy erlaubt nur 1 Modell) gilt Chat als ok; Wartezeit 40x3 s.
### Offen / naechste Schritte
- Tab 4 (adiabi26444): nach 1.0.68 `node api.js act cowork_check` -> MCP-Auto-Einrichtung mit neuer URL testen. Trial-Kuendigung meldet "nocancel" (Abrechnung-Tab: Knopf-Text pruefen). Usage-Tab nicht gefunden (neue Konten).
- Tab 1 (trinix, "Adrian Hotz's Space", EN) + Tab 2 (qoleqabi, "Jens Steffen", EN) sind eingeloggt, Einstellungen oeffnen wieder -> Refresh erneut testen (`node api.js act cowork_refresh_confirm 1` bzw. 2).
- Danach: cowork_start_all, `node sendmain.js 9301 task66.txt`, Zeiten aus `profiles\timeline.log`, Ergebnis `node test.js` in work\cowork-test66.
- Neue Testwerkzeuge in work\bugtest: shot.js (CDP-Screenshot, bei versteckten Fenstern leer), tclick.js (trusted Klick x y), settabs.txt (__TAB__=Regex), addconn.txt, custmcp.txt, popup.txt, overlay.txt, authtest.js (Auth-Test ohne Token-Ausgabe), chkinline.js (Syntax der inline-JS in setup_mcp).

## STAND 1.0.69 (released, Commit 81c0553, installiert)
- Problem: Kollege konnte mit eigener ngrok-URL + Bearer-Token keinen MCP-Server in Notion einrichten.
- Fix (notion_mcp.rs): clean_tok() entfernt Leerzeichen/Quotes/<>/Zero-Width und mehrfache Praefixe 'Authorization:', 'Bearer ', 'Token ' (z. B. 'Bearer abc' ins Bearer-Feld kopiert). token_candidates(): Authorization, x-api-key, api-key, x-auth-token, x-adicode-token, /k/<token>, ?token=/key=/api_key=/access_token=. is_mcp_path(): Root '/', /mcp, /k/<t>, /k/<t>/mcp. OPTIONS -> 204 vor Auth, CORS-Header auf allen MCP-Antworten. GET / ohne Token = Info-Text.
- Diagnose: auth_check() liefert Grund (ohne Token-Inhalt, nur Laenge/Quelle); note_auth_fail() -> Ring-Puffer 15 (state.auth_failures, auth_fail_count), UI-Zeile 'Abgelehnt' auf Notion AI (nAuthRow/nAuthFails), Datei %LOCALAPPDATA%\PlazCodeNotion\logs\mcp-auth.log. 401-Body enthaelt den Grund.
- Test: Unit-Test notion_auth_variants ok (15/15). Oeffentlich ueber ngrok: Bearer, Bearer Bearer, ohne Bearer, x-api-key, Root-URL, key-url, key-url ohne /mcp = 200; falsch/ohne = 401; OPTIONS = 204.
- Naechster Schritt: Kollege auf 1.0.69 (Auto-Update), erneut verbinden; klappt es nicht, Zeile 'Abgelehnt' bzw. logs\mcp-auth.log schicken lassen. Steht dort nichts, erreicht Notion seinen PC gar nicht (Tunnel/Domain/Firewall). Release-Skript: work\bugtest\rel69.ps1. Backups: bugtest\notion_mcp.rs.bak169, desktop.html.bak169.

## STAND 1.0.70 (released)
- Kollege: seine Domain liefert ERR_NGROK_3200 (Endpoint offline) -> Notion: 'URL doesn't point to a supported MCP endpoint'. Ursache = sein ngrok-Tunnel startet nicht, NICHT der Token. Verdacht: Adrians ngrok-Authtoken in seiner ngrok.yml (Domain gehoert nicht zum Konto / ERR_NGROK_108). Loesung: eigener Authtoken (dashboard.ngrok.com > Your Authtoken) auf Notion AI > Tunnel eintragen, eigene Domain, Tunnel muss 'online' zeigen. Pruefen: curl -i https://<domain>/mcp -> 401 = ok, Header Ngrok-Error-Code = Tunnel-Problem.
- 1.0.70: ngrok-stderr wird gesammelt (vorher nur 'ERROR:'), neue Fehlertexte 8012/121/105/107. Fenster-Waechter in notion_profiles.rs windows (WANT/PIDCACHE/APPLY, guard 0,9 s) versteckt neue Fenster versteckter Worker; Zeigen nur Hauptfenster. Backups bugtest\*.bak170.

## STAND 1.0.79
- Chat-Tab (notion_chat.rs/notion_chat_ui.html): UI laedt Verlauf alle 4 s neu (vorher nur beim Oeffnen der Seite -> neue Chats fehlten links). Titel: STATE liefert ntitle (document.title ohne ' | Notion'); remember() nimmt ntitle, sonst die ersten 7 Woerter der ersten Nachricht. STATE erkennt DE-Labels (Nachricht bearbeiten, Antwort kopieren, Vorschau).
- Co-Work + Main-Chat: Ursache 'Co-Work klappt nicht im Main-Chat' = 'Neuer Chat' im Chat-Tab hatte den Main-Prompt nicht. send(): Chat ohne Nachrichten + STAGE 2..5 -> np::cowork_main_context() (Main-Prompt) wird vor 'Meine Aufgabe:' gesetzt.
- hide_all_tabs: nach gruenem Check (check_action) und im Kickoff alle Notion-Fenster inkl. Main verstecken, notion_chat::set_main_hidden(true).
- Nicht getestet: kompletter Vier-Konten-E2E (bewusst nicht ausgefuehrt).

## STAND 1.0.78
- load_watchdog (einmal gestartet in launch()): alle 5 s pro Tab STUCK_LOAD_JS (Notion-Host, <25 Zeichen Text, <2 beschriftete Bedienelemente) -> nach 20 s Page.reload(ignoreCache), beim 3. Mal Navigate NOTION_URL; Feld load_state. Gleiche Schleife: COOKIE_JS lehnt Cookie-Banner ab (Reject all/Alle ablehnen/Nur notwendige ...), Feld cookie.
- READY_JS verlangt >25 Zeichen Text (Skelett gilt nicht als fertig).
- Checks einmal pro Konto (accounts.json): lang_en (english_once), mcp_ok2 (setup_mcp uebersprungen, 'Verbunden (gemerkt)'), models_ok==tid, chat_model==tname, usage_ui_ok. Trial-Schritt unveraendert. mem_forget (Vergessen + Refresh) loescht alles ausser lang_en.

## STAND 1.0.77
- Notion hat mehrere Layouts: (a) Chat-App app.notion.com/chat (Tabs Home/Chat/Meetings/Inbox), (b) klassisch mit Chat als Seitenleiste/Floating rechts (Ansicht-Menue Sidebar/Floating/Full screen). prepare_chat ruft nach FIND_WELCOME ensure_fullscreen (FS_BTN_JS: Knopf oben rechts mit view/sidebar/..-Label, nicht close/new/share/pin; FS_ITEM_JS: Full screen/Vollbild). welcome_url wird nur genutzt, wenn sie /chat enthaelt.
- Verbindungen-Seite: Tabs 'Browse Connections | Connected N | Manage' (EN) bzw. 'Verbindungen durchsuchen | Verwalten' (DE). Existing-Check jetzt Connected/Verbunden, sonst Manage > All connections. Add-Knopf: Add connection, sonst Configure/Konfigurieren, sonst Add; abwechselnd click und Pointer-Events.
- Ungetestet live: Layout (b) Full screen (Tab 3 Mia hat es, darf nicht angefasst werden), Configure-Variante.

## STAND 1.0.76
- set_english(port): Fast-Path html lang=en; sonst open_settings > Sprach-Knopf (LANG_BTN_JS, ggf. Tab Meine Einstellungen/Preferences) > Option English (US) > Bestaetigen (Aktualisieren/Update) > Reload. Laeuft in check_one (vor Trial) und am Anfang von setup_mcp. Aktion cowork_set_english <id>, Feld lang. Live auf Tab 4 manuell getestet (de -> en).
- mcp_already/ALREADY_JS: 'Dieser MCP-Server wurde bereits ... hinzugefuegt' / 'already added' => Escape, mcp=Bereits verbunden, Ok (vorher Fehler -> 'MCP not configured' beim Kollegen).
- SETTINGS_ITEM_JS sucht zuerst im offenen Menue/Dialog (traf vorher den Einstellungen-Knopf im Chat-Eingabefeld -> open_settings schlug fehl).
- CHATS_JS nur Seitenleiste (left<420, nicht in Dialogen); check_chats schliesst offene Dialoge (DIALOG_OPEN_JS+Escape), wait_ready, 1x Retry bei 0.
- Testskripte bugtest: runjs.js <port> <CONST_NAME>, evalfile.js, tclick.js.

## STAND 1.0.75
- setup_mcp/check_usage_ui nutzen open_settings(port): Konto-Menue oben ODER unten links, trusted CDP-Klicks, Escape, Fallback Strg+,.
- wait_ready/poll_js statt fester Sleeps (Seite darf laden).
- FILL_JS: Name AdiCode + Bearer-Token (auch deutscher Dialog 'Individueller MCP-Server'), Auth ggf. auf Bearer, wartet auf aktiven Verbinden-Knopf.
- Ungetestet live: Kollegen-Layout mit Menue unten.

## STAND 1.0.74 (released)
- ngrok-Schnellstart (Aktion ngrok_quick_start, Knopf nQuickStart): Authtoken einfuegen -> Start. Installiert eigenes ngrok (base_dir), add-authtoken, killt alle ngrok, ermittelt Gratis-Domain (ngrok ohne --url, url= aus Log), startet Tunnel, testet POST /mcp = 401 ohne Ngrok-Error-Code.
- 'ngrok ist zu alt' beim Kollegen: find_ngrok nahm altes ngrok aus PATH, Neu-installieren landete in base_dir und wurde ignoriert. Jetzt: eigenes ngrok zuerst, install killt alte Prozesse + setzt ngrok_path, Schnellstart reinstalliert bei 'zu alt'.
- Authtoken-Feld leert sich nach Speichern (Absicht) - jetzt Anzeige 'gespeichert abcd...wxyz' (authtoken_mask). Store-ngrok: Config liegt unter Packages\ngrok.ngrok_*\LocalCache\Local\ngrok, wird erkannt und nach %LOCALAPPDATA%\ngrok kopiert.
- Offen: Schnellstart live auf Adrians PC testen (trennt Tunnel kurz), Kollege 1.0.74 Start druecken.
## SCHNELL-ARBEITSREGELN FUER AI-CHATS (immer zuerst lesen und befolgen)
Ziel: maximal schnell antworten, ohne Qualitaet zu verlieren. Wenige, grosse, gezielte Schritte statt vieler kleiner.

**1. Start**
- Nur `project_memory_open "PlazCode-Notion"` laden, dann sofort loslegen. Keine Tool-Listen, Docs oder Status vorab abrufen, wenn das Handout reicht.
- Pfade, Ablaeufe und Skripte stehen hier: nicht neu erkunden, direkt benutzen.

**2. Tool-Calls buendeln**
- Ein `pc_PowerShell`-Call macht mehrere Dinge: lesen + suchen + Backup + Patch + Check in einem Befehl.
- Unabhaengige Calls parallel im selben Schritt (z. B. project_memory_update + HANDOUT-Commit).
- Code lesen: `Select-String -Pattern 'a|b|c'` mit Zeilennummern, dann gezielt `(Get-Content f)[x..y]`. Nie ganze grosse Dateien ausgeben.
- Ausgaben klein halten: `Select -First N`, `-Tail N`, `Select-String 'error|test result|TEST'`.

**3. Aendern**
- Patches als Node-Skript `work\bugtest\patchNNN.js` per `files_write_file` (exakte Anker, Abbruch bei fehlendem/doppeltem Anker). Vorher Backups `bugtest\<datei>.bakNNN` (nie in agent/src, sonst landen sie im Patch).
- JS in desktop.html pruefen: `node work\bugtest\chk169.js` (muss `bad 0` melden).
- Bei Unsicherheit ueber Ursache: erst EIN gezielter Messbefehl (curl -i auf die URL, Header `Ngrok-Error-Code`, Logs), dann fixen. Nicht raten.

**4. Testen + Release in EINEM Hintergrundjob**
- `tr74.ps1` kopieren -> `trNN.ps1` (cargo test, nur bei Erfolg relNN.ps1). `relNN.ps1` aus dem letzten relXX.ps1 ableiten (Versionen + Notizen ersetzen). Start mit `Start-Process powershell -WindowStyle Hidden -ArgumentList '-NoProfile','-ExecutionPolicy','Bypass','-File',...` OHNE -RedirectStandardOutput (das blockiert den Call) - Logs per `*>` im Skript.
- Dauer: Test ~40 s, Release ~80 s. NICHT mit `sleep`/`Start-Sleep` darauf warten. In der Zwischenzeit HANDOUT/Memory schreiben oder den naechsten Fix vorbereiten, dann EINMAL Log pruefen.
- MCP-Timeout 55 s: nie laengere Befehle synchron.
- Nach einem Release startet AdiCode neu (Tunnel ~30-60 s weg, 'ERR_NGROK_3004'/'Unknown tool'). Nicht pollen; erst die Antwort an den Nutzer schreiben, spaeter verifizieren.

**5. Abschluss**
- HANDOUT (Abschnitt STAND x.y.z) + `project_memory_update` parallel in einem Schritt, Commit mit `[skip ci]`.
- Antwort kurz: Ergebnis zuerst, Ursache in 1-2 Saetzen, was der Nutzer/Kollege tun soll, was noch nicht getestet ist. Keine langen Erklaerungen.

**6. Nicht tun**
- Keine Mehrfach-Abfragen desselben Status, keine Warte-Schleifen, keine erneuten Tool-Beschreibungen fuer bekannte Tools.
- Keine Rueckfragen, wenn die naechste Aktion klar ist - einfach machen. Nur fragen, wenn die Antwort die Arbeit wirklich aendert.
