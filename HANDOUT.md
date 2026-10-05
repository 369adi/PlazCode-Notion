# PlazCode Notion – Handout (Stand 05.10.2026)

## Ziel
Original-PlazCode-Desktop-App (stoveez/PlazCode, GPL-3.0) als **inoffiziellen Fork „PlazCode Notion“** umbauen:
Notion-Design + eingebauter MCP-Server, damit Notion AI Roblox Studio, PC und Browser **ohne Browser-Erweiterung** steuert.

## Aufbau des Repos
| Datei | Zweck |
|---|---|
| `notion-desktop/plazcode-notion.patch` | Alle Änderungen ggü. PlazCode 1.19.36 |
| `notion-desktop/VERSION` | Release-Version (aktuell 1.0.1) |
| `PlazCode-source-1.19.36.zip` | Original-Quellcode (Build-Basis) |
| `PlazCode-1.19.36.zip` | Original-Paket, darin wird nur `PlazCode.exe` ersetzt |
| `.github/workflows/build-notion-desktop.yml` | Baut auf windows-latest mit Cargo, veröffentlicht Release `notion-desktop-v<VERSION>` |

## Inhalt des Patches (agent/src)
- `notion_mcp.rs` (neu): MCP-Endpunkt `127.0.0.1:8787/mcp` (Bearer-Token oder `/k/<token>/mcp`), Tools `plazcode_status`, `roblox_*`, Add-on-Tools (`pc_*`, `browser_*` …); ngrok-Supervisor mit fester Domain; API `/api/notion/state` + `/api/notion/action`.
  Einstellungen: `%LOCALAPPDATA%\PlazCodeNotion\plazcode-notion.json` (übernimmt beim ersten Start Token/Domain der alten Go-Bridge `config.json`).
- `main.rs`: Modul registriert, Routen, Start, ngrok beim Beenden stoppen.
- `mcp_addons.rs`: Katalog + „PC (Windows-MCP)“ (`uvx --python 3.14 --from windows-mcp==0.8.7 windows-mcp`) und „Browser (Playwright)“ (`npx -y @playwright/mcp@latest`, Brave/Chrome/Edge automatisch).
- `updater.rs`: Auto-Update deaktiviert (sonst überschreibt das Original den Fork).
- `gui.rs`: Fenstertitel „PlazCode Notion“, GitHub-Link auf den Fork.
- `desktop.html`: Notion-Design-CSS, Seite „Notion AI“, Statuskarte auf Home, Branding.

## Status
- Release 1.0.0: **fehlerhaft** – `git apply` lief in einem Repo-Unterordner und hat alle Dateien still übersprungen → exe = Original-Design ohne Notion-Tab. Nicht verwenden.
- Fix in 1.0.1: `git apply --directory=build-src/PlazCode` vom Repo-Root + Prüfschritt (bricht ab, wenn der Patch fehlt). Build lief beim Erstellen dieses Handouts – Ergebnis unter *Actions* prüfen.
  Jetzt wird der neue Rust-Code erstmals wirklich kompiliert → evtl. Compiler-Fehler im Actions-Log beheben, Patch neu erzeugen, VERSION erhöhen, pushen.

## Nutzung
1. Original-PlazCode **und** alte `PlazCode-Notion-Bridge.exe` beenden (blockieren Port 3000/8787 und die ngrok-Domain). Desktop-Skript „PlazCode Notion starten.cmd“ erledigt das.
2. `PlazCode-Notion-1.0.1-windows.zip` aus Releases entpacken, `PlazCode-Notion\PlazCode.exe` starten.
3. Seite **Notion AI**: Server-URL `https://<ngrok-domain>/mcp` + Bearer-Token in Notion als benutzerdefinierten MCP-Server eintragen.
4. Unter **MCP Servers** „PC (Windows-MCP)“ und „Browser (Playwright)“ aktivieren; Roblox-Tools erscheinen, sobald Studio verbunden ist.

## Offene Punkte / Ideen
- Build-Ergebnis 1.0.1 prüfen, ggf. Compiler-Fehler fixen.
- Altes Release `notion-desktop-v1.0.0` löschen.
- Optional: Dark-Mode im Notion-Stil, Steuer-Tools (Server neu starten) per MCP.

## Hinweise
Inoffizieller Fork, keine Verbindung zu den PlazCode-Autoren oder Notion Labs. `LICENSE` und `BRANDING-NOTICE.txt` müssen erhalten bleiben. Tokens nie veröffentlichen.
## Neu in 1.0.13
- Server-URL direkt eintragbar: Auf der Seite **MCP Servers** gibt es oben die Karte *Notion-Verbindung (ngrok)*, und auf der Seite **Notion AI** ist die Server-URL jetzt ein Eingabefeld. ngrok-Domain oder ganze URL einfügen, Enter oder Speichern.

## Neu in 1.0.12
- Fix: Die Inhalte der Seite Notion AI (Tool-Liste) erschienen auf jeder Seite, und die Felder für Authtoken und Domain / URL fehlten. Beides ist repariert.

## Neu in 1.0.11
- Updates kommen schneller: Der Launcher prüft jetzt alle 15 Sekunden statt alle 90 Sekunden auf neue Versionen.

## Neu in 1.0.10: Updates-Seite + einfache Ersteinrichtung
- Neue Seite **Updates**: installierte und neueste Version, „Was ist neu?“ für jede Version, Launcher-Protokoll, Knöpfe „Jetzt prüfen“ und „Jetzt aktualisieren“. Updates installiert der Launcher weiterhin automatisch bei jedem Nutzer.
- Ngrok-Ersteinrichtung direkt in der App: **ngrok installieren** (lädt ngrok automatisch herunter), **Authtoken** eintragen, **Domain / URL** eintragen (ganze URL geht auch, Enter speichert).
- Fix: Das Domain-Feld wurde alle 2,5 s vom Status überschrieben, wenn man woanders hingeklickt hatte – Eingaben bleiben jetzt stehen, bis man speichert.
## Neu in 1.0.9: Aufgeräumte Oberfläche + getrennte Notion-Fenster
- Model Builder, UI Builder, Toolkit, Templates und Updates sind aus der Oberfläche entfernt (Code bleibt intern erhalten).
- Seite Notion AI → Karte **„Notion-Fenster (getrennte Logins)“**: „Neues Notion-Fenster“ startet Chrome/Edge/Brave mit eigenem Profil (`%LOCALAPPDATA%\PlazCodeNotion\profiles\agent-N`) als App-Fenster mit Notion AI. Jedes Profil hat eigene Cookies – Ein-/Ausloggen wirkt nur dort. Profile können benannt, geöffnet und gelöscht werden (`notion_profiles.rs`).
## Neu in 1.0.8: Funktionen aus Secretscript
- **Web-Agenten** (`notion_web_agents.rs` + Chrome-Erweiterung in `notion-desktop/extension`): Tools `web_sites`, `web_agents`, `web_chat`. Notion AI kann freigegebene KI-Chat-Tabs (ChatGPT, Claude, Gemini, DeepSeek, Qwen, Kimi, GLM, Arena) nutzen. WebSocket `ws://127.0.0.1:8787/extension/ws`, eigener Erweiterungs-Token (Seite Notion AI → „Web-Agenten & Sicherheit“).
- **Direkt-Modus**: Ein KI-Chat arbeitet als Roblox-Agent (`/direct-agent/*`, 300 s Laufzeit, max. 50 Aufrufe); jede Änderung braucht eine Freigabe im Panel auf der Seite. Solange er läuft, sind Notion-Schreibzugriffe auf Studio gesperrt.
- **Roblox Plus** (`notion_roblox_plus.rs`): `roblox_studio` mit 23 High-Level-Aktionen, `roblox_workflow` (Playbooks), `roblox_project_memory` (`%LOCALAPPDATA%\PlazCodeNotion\roblox-project-memory.json`), automatische Snapshots vor Schreibzugriffen + Wiederherstellen, BW-Helper-Bibliothek (`roblox_helpers.luau`) wird bei `execute_luau` automatisch geladen.
- **Sicheres Datei-Schreiben** (`notion_safety.rs`): `safe_read_file`/`safe_write_file` (SHA-256-Prüfung gegen Überschreiben veralteter Stände)/`safe_search_files`, Secret-Dateien (.env, Schlüssel, Tokens …) werden blockiert, Co-Work-Locks werden beachtet.
- **Sicherheitsmodus für die Shell** (Schalter, Standard aus): PowerShell nur mit erlaubten Einzelbefehlen, Registry und unsicherer Browser-Code gesperrt.
## Neu in 1.0.7: Co-Work mit festen Experten-Rollen + Live-Status

Der Main Chat ist der Lead: Er zerlegt das Projekt, vergibt Tasks nur an die Rollen, die wirklich gebraucht werden, und fuehrt am Ende alle Ergebnisse zu EINER optimierten Antwort zusammen (cowork_results).

Rollen: Prompt-Schreiber, Programmierer, Code-Bewerter, UI/UX-Kritiker, Tester, Roblox-Spezialist, Generalist (eigene Rollen moeglich).

So geht's:
1. Main Chat: "Starte ein Co-Work-Projekt: <Ziel>". Die KI nennt Projekt-ID und welche Rollen-Tabs du oeffnen sollst.
2. Pro weiterem Tab: "Tritt Co-Work-Projekt p1 bei als Code-Bewerter" (bzw. Programmierer, UI/UX-Kritiker ...).
3. Fertige Code-/UI-Tasks gehen automatisch in die Pruefung (Review). Fordert der Pruefer Aenderungen an (Status changes), entsteht automatisch ein Fix-Task fuer die Original-Rolle (max. 2 Runden, danach entscheidet der Lead).
4. In der Desktop-App zeigt der Button "Co-Work" unten rechts live alle Agents, Rollen und Task-Status.

Robustheit: atomare Task-Uebernahme, Datei-Sperren, niemand prueft seine eigene Arbeit, unbesetzte Rollen werden von anderen uebernommen, haengende Agents geben ihre Tasks nach einstellbarer Zeit (stale_minutes, Standard 20) frei, Lead-Uebernahme wenn der Main Chat ausfaellt, Ergebnis Pflicht bei done/changes/blocked.
## Neu in 1.0.6: Co-Work (mehrere Notion-AI-Chats an einem Projekt)

Mehrere Notion-AI-Tabs im Browser arbeiten gleichzeitig am selben Projekt. Sie koordinieren sich ueber ein gemeinsames Board in PlazCode Notion: Tasks mit Abhaengigkeiten, atomare Uebernahme (kein Task wird doppelt bearbeitet), Datei-Sperren und Nachrichten zwischen den Chats.

So geht's:
1. Tab 1 (Lead): "Starte ein Co-Work-Projekt: <Ziel>". Die KI zerlegt das Ziel in Tasks und nennt eine Projekt-ID (z. B. p1).
2. Weitere Tabs (2-5 empfohlen): "Tritt Co-Work-Projekt p1 bei". Jeder Tab holt sich automatisch freie Tasks, arbeitet sie ab und meldet das Ergebnis.
3. Der Lead arbeitet mit, prueft Ergebnisse, ergaenzt Fix-Tasks und macht am Ende die Endkontrolle.

Tools: cowork_start, cowork_join, cowork_board, cowork_claim, cowork_update, cowork_add_tasks, cowork_message, cowork_lock, cowork_wait, cowork_close. Skill: plazcode_skill name=cowork.
Das Board wird in %LOCALAPPDATA%\PlazCodeNotion\cowork.json gespeichert und uebersteht Neustarts. Inaktive Agents (20 Min.) geben ihre Tasks automatisch frei.
## Neu in 1.0.5 – mehr Tools & Skills fürs Programmieren
- **Local files** (`files`): `@modelcontextprotocol/server-filesystem` auf einen Shared-Ordner (Standard `%USERPROFILE%\PlazCode-Shared`, änderbar per Umgebungsvariable `PLAZCODE_SHARED_DIR`). Pfade sind relativ zum Shared-Ordner (Arbeitsverzeichnis des Servers).
- **GitHub** (`github`): `@modelcontextprotocol/server-github`. Vorher Benutzer-Umgebungsvariable `GITHUB_PERSONAL_ACCESS_TOKEN` setzen und App neu starten. Token nie ins Repo.
- Bereits im Katalog: Context7, Fetch, Git, Memory, Sequential Thinking – unter **MCP Servers** aktivieren.
- Neues Modul `notion_skills.rs`:
  - Tool `plazcode_skill` mit eingebauten Skills `roblox`, `roblox-studio`, `app`, `debug`, `git`.
  - Tool `plazcode_restart_server` startet ein hängendes Add-on (oder `all`) neu.
  - Ausführlichere MCP-`instructions` für Notion AI.
- Patch-Hinweis: `desktop.html`, `gui.rs`, `updater.rs` unverändert aus 1.0.4 übernommen; `main.rs`, `mcp_addons.rs`, `notion_mcp.rs`, `notion_skills.rs` neu erzeugt und per `git apply` gegen die Original-Quelle geprüft.
