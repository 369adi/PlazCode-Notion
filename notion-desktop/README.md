# PlazCode Notion Desktop (inoffizieller Fork)

**PlazCode Notion** ist ein inoffizieller Fork der PlazCode-Desktop-App (GPL-3.0) von stoveez/PlazCode.
Er steht in keiner Verbindung zu den PlazCode-Autoren oder zu Notion Labs. Alle Original-Lizenz-,
Copyright- und Branding-Hinweise (LICENSE, BRANDING-NOTICE.txt, THIRD-PARTY-NOTICES) bleiben erhalten.

## Was ist anders?
- Notion-Design (hell, ruhig, Notion-typische Farben und Typografie), Name „PlazCode Notion“.
- Neue Seite **Notion AI**: eingebauter MCP-Server (`127.0.0.1:8787/mcp`) + ngrok-Tunnel mit fester Domain.
  Notion AI verbindet sich direkt – **ohne Browser-Erweiterung**.
- Notion sieht `roblox_*` (Roblox Studio MCP), `pc_*` (Windows-MCP), `browser_*` (Playwright) und alle weiteren aktivierten MCP-Server.
- Neue Katalog-Einträge unter **MCP Servers**: „PC (Windows-MCP)“ und „Browser (Playwright)“.
- Automatische Updates sind deaktiviert, damit der Fork nicht durch die Original-Version ersetzt wird.

## Neu in 1.0.18
- Updates beenden Co-Work-Browser nicht mehr: Der Launcher schliesst beim App-Neustart weiterhin AdiCode, ngrok und Add-ons, bewahrt aber alle Browserprozesse der vier Profilordner. Vor Update-Neustarts und normalem Beenden werden versteckte Worker-Fenster wieder sichtbar gemacht.
- „Alle zeigen“ stellt versteckte Worker-Fenster mit Windows `SW_RESTORE` wirklich wieder her.
- Main-First-Ablauf: **Alle starten** öffnet nur Main sichtbar. Aufgaben werden direkt im Main-Notion-Chat geschrieben; Coder, Reviewer und Tester bleiben versteckt und treten bei, sobald Main ein neues Co-Work-Projekt startet.
- Neues eingebautes `project_memory_*`: getrennte dauerhafte Kontexte pro Projekt, Suche/Öffnen/Aktualisieren und automatische Synchronisierung von HANDOUT.md-Dateien im Shared-Work-Ordner. Damit kann ein langer Main-Chat zwischen Aurelune, PlazCode usw. wechseln, ohne die Projekte zu vermischen.
- Neuer Skill `project-memory`: vor Weiterarbeit Projekt-Memory laden; nach wichtigen Änderungen HANDOUT.md und Memory aktualisieren. Roher Chatverlauf wird nicht als dauerhaft oder unbegrenzt angenommen.
## Neu in 1.0.17
- **Co-Work 2.0:** genau vier Tabs mit festen Rollen nach E-Mail-Reihenfolge (Main, Coder, Reviewer, Tester). Start erst mit vier Konten und eingetragener Aufgabe.
- **Auto-Kickoff per CDP:** Jeder Tab wird zuerst gezielt auf Notion AI (`/ai`) navigiert; nur der sichtbare untere AI-Composer wird verwendet. Pro Tab gibt es weiterhin „Startprompt kopieren“ als Fallback.
- **MCP pro Konto:** AdiCode versucht die gespeicherte MCP-URL und den Bearer-Token automatisch in jedem E-Mail-Tab unter Notion Connections einzurichten; eigener Button und manueller Fallback bleiben verfügbar.
- **Zuverlaessige Hintergrundarbeit:** Chromium-Drosselung und Native Window Occlusion sind deaktiviert; ein CDP-Keeper haelt alle Tabs auch versteckt und unfokussiert aktiv. Nur Main bleibt sichtbar, „Alle zeigen“ holt die anderen zurueck.
- **Account-Check:** Usage-Fortschritt und erkannte verfuegbare Notion-AI-Modelle pro Konto.
- **Codebase Memory:** `codebase-memory-mcp` 0.11.0 als standardmaessig aktiviertes Add-on.
- Co-Work-Antworten enthalten nur noch neue Meldungen statt immer das ganze Board.
## Neu in 1.0.16
- Tunnel-Selbstheilung: AdiCode prueft die oeffentliche ngrok-URL regelmaessig ueber `/health`. Wenn ngrok noch laeuft, die reservierte Domain aber nicht mehr erreichbar ist, wird der Tunnel automatisch neu verbunden. Damit meldet Notion nicht mehr dauerhaft, dass eine korrekte MCP-URL kein unterstuetzter Endpoint sei.

## Neu in 1.0.15
- Neuer Name **AdiCode** und neues Logo (A-Monogramm mit Farbverlauf), auch als Icon der exe. Die Tools heissen jetzt `adicode_status`, `adicode_skill` und `adicode_restart_server`.
- Komplett neues **dunkles Design**: ruhige dunkle Flaechen, Akzent-Verlauf Violett-Cyan, neue Karten, Buttons, Eingabefelder, Badges und Tabellen.
- Stabilerer MCP-Server: Jeder Tool-Aufruf laeuft abgesichert (ein Absturz reisst den Server nicht mehr mit) und antwortet nach spaetestens 55 s. Sehr grosse Antworten werden gekuerzt.
- Schneller: Jedes Add-on (pc, browser, git, ...) hat jetzt eine eigene Warteschlange. Ein langer Befehl auf einem Add-on blockiert die anderen nicht mehr. Haengende oder abgestuerzte Add-ons werden automatisch neu gestartet.
- Co-Work: Das Board wird gebuendelt im Hintergrund gespeichert, der Nachrichtenverlauf ist auf 300 begrenzt.

## Neu in 1.0.14
- Neue Seite **Co-Work**: mehrere Notion-AI-Tabs gleichzeitig, jeder in einem eigenen Container (eigenes Browser-Profil, eigene Cookies). Pro Tab ein Gmail-Konto, Name, Rolle, Startseite und Auto-Start. **Speichern** merkt sich alles, **Gmail-Login** meldet den Tab einmal an, **Alle starten** oeffnet alle Tabs mit ihren eingeloggten Sessions, **Alle stoppen** schliesst sie.
- HANDOUT.md komplett neu: Projektueberblick, Build-Ablauf, Architektur, Fallstricke und Changelog fuer neue Chats.

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
## Neu in 1.0.8
- Web-Agenten: Notion AI nutzt freigegebene KI-Chat-Tabs über die Chrome-Erweiterung (Ordner `notion-desktop/extension` im Repo) (Installation: `chrome://extensions` → Entwicklermodus → „Entpackte Erweiterung laden“; Token von der Seite Notion AI → „Web-Agenten & Sicherheit“).
- Direkt-Modus: ein KI-Chat als Roblox-Agent, jede Änderung mit Freigabe.
- Roblox: `roblox_studio` (23 Aktionen), `roblox_workflow`, Projekt-Gedächtnis, automatische Snapshots + Wiederherstellen, BW-Helper werden automatisch geladen.
- Sicheres Datei-Schreiben mit SHA-256-Prüfung und Schutz von Secret-Dateien.
- Sicherheitsmodus für die Shell (Schalter auf der Seite Notion AI, Standard aus).
## Start
1. Original-PlazCode und die alte `PlazCode-Notion-Bridge.exe` beenden (Port 3000 / 8787 und ngrok-Domain werden sonst blockiert).
2. `PlazCode-Notion.exe` herunterladen und starten (entpackt sich beim ersten Start nach `%LOCALAPPDATA%\PlazCodeNotion\app`).
3. Seite **Notion AI** öffnen → Server-URL und Token in Notion als benutzerdefinierten MCP-Server eintragen.
   Token und Domain der alten Bridge werden beim ersten Start automatisch übernommen.

Einstellungen: `%LOCALAPPDATA%\PlazCodeNotion\plazcode-notion.json`

## Build
GitHub Actions (`.github/workflows/build-notion-desktop.yml`) entpackt `PlazCode-source-1.19.36.zip`,
wendet `notion-desktop/plazcode-notion.patch` an, baut `agent` mit Cargo (MSVC) und ersetzt `PlazCode.exe`
im Original-Paket `PlazCode-1.19.36.zip` und packt alles in eine einzelne `PlazCode-Notion.exe`.
