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

## Neu in 1.0.33
- **Echte Usage**: Der Check liest den Verbrauch direkt aus Einstellungen -> Notion KI -> Usage (Monthly, "x% used", Reset-Datum) und zeigt ihn im Balken.
- **Chat-Check**: zeigt, ob ein Konto nur den Chat "Willkommen bei Notion" hat; bei mehreren Chats wird der Chip gelb.
## Neu in 1.0.32
- **Fortschrittsbalken beim Start**: Alle starten zeigt einen Balken (Konten pruefen, Main-Chat, Projekt, Worker beitreten) und am Ende gruen "Aktiv".
- **Trial-Check**: Der Check erkennt pro Konto, ob ein Trial laeuft (Abo aktiv, noch nichts bezahlt) und zeigt es als Chip.
## Neu in 1.0.31
- **Modell-Check pro Rolle**: Der Account-Check stellt unter Notion KI -> Model controls -> Allowed models for Notion Agent automatisch genau ein Modell ein - Main = Opus 5.5 (Orchestrator), Coder/Reviewer/Tester = Sonnet 5.5 (deutlich guenstiger). Alle anderen Anbieter werden abgeschaltet.
- **Checkliste pro Konto**: gruene/rote Chips fuer MCP configured, Modell, Chat und Usage. Fehlt der Chat, erscheint rot "Chat fehlt" mit Button "Zum Chat".
- **Gemerkte Checks**: MCP und Modelle werden pro Gmail-Konto gespeichert (accounts.json) und nicht erneut geprueft; "neu pruefen" setzt das zurueck.
- **Nach dem Check zurueck zum Chat**: offene Einstellungen werden per X geschlossen, jedes Konto landet im Chat "Willkommen bei Notion"; Worker-Fenster werden nach erfolgreichem Start automatisch versteckt.
- **UI**: kleine Beschriftungen unter allen Icons, Worker zeigen/verstecken erkennt den echten Zustand, kompakter Browser-Waehler mit Browser-Icons.
## Neu in 1.0.30
- **Co-Work neu gedacht**: oben eine Session-Leiste mit Status (x von 4 bereit), Browser, Check, Worker zeigen, Stoppen und dem grossen Alle-starten-Button rechts daneben. Die Konten sind kompakte Zeilen (Rolle + Status, E-Mail, Usage-Balken, Icon-Buttons) - das separate Usage-Dashboard ist darin aufgegangen.
- **Keine schwarzen Fenster mehr**: Zeigen/Verstecken fasst nur noch echte Notion-Fenster an statt aller unsichtbaren Browser-Hilfsfenster.
## Neu in 1.0.29
- **Co-Work aufgeraeumt**: "Alle starten" / "Stoppen" stehen jetzt ganz oben neben dem Titel. Browser, Check, Worker zeigen, Tab hinzufuegen und Speichern sind in einer kompakten Leiste. Die Konto-Karten haben nur noch Icon-Buttons (Name als Tooltip), "Entfernen" sitzt rechts abgesetzt.
- **Schnelle Releases**: Updates werden lokal inkrementell gebaut (ca. 5 s) und direkt als GitHub-Release hochgeladen - kein 5-Minuten-Build in GitHub Actions mehr.
## Neu in 1.0.28
- **Ein einziger Check** (Button "Check (MCP - Usage - Chat)"): startet fehlende Tabs, wartet bis sie erreichbar und eingeloggt sind, prueft ob der AdiCode-MCP installiert ist (auch unter dem Namen "asf" oder "PlazCode") und richtet ihn sonst ein, liest das echte Usage-Limit und oeffnet am Ende den Chat "Willkommen bei Notion". Jeder Co-Work-Start macht diesen Check automatisch - auch fuer den Main-Tab.
- **Echtes Usage**: Basis-KI (z. B. 143/50, im Plan "unlimited") und Premium-Credits (z. B. 0/300) werden getrennt angezeigt. Vorher stand faelschlich immer 0/300.
- **Nie ein neuer Chat**: AdiCode legt keinen neuen Notion-Chat mehr an, auch nicht im Main-Tab. Fehlt der Chat "Willkommen bei Notion", stoppt der Check mit klarer Meldung.
- **Kein Auto-Verstecken**: Die drei Worker-Fenster bleiben beim Start sichtbar und werden erst ueber "Worker verstecken" ausgeblendet.
- **Richtige Version oben links**: Seitenleiste und Titel zeigen jetzt die echte Release-Version statt der internen Agent-Nummer.
- Veraltete Meldung "Tab laeuft nicht" wird beim Start sofort geloescht; das Dashboard aktualisiert sich nach dem Check mehrfach.
## Neu in 1.0.27
- **Browser zuerst waehlen**: Co-Work hat oben Schritt 1 "Browser waehlen" (Chrome, Edge oder Brave - nur installierte werden angeboten). Alle Tabs starten in diesem Browser; ohne Auswahl ist "Alle starten" gesperrt. Wechsel nur, wenn alle Tabs gestoppt sind.
- **Co-Work neu gestaltet**: drei Schritte (Browser -> Konten & Tabs -> Starten), Tabs als Karten statt Tabelle, alle Buttons mit Icons (Play, Stop, Speichern, Login, MCP, Kopieren ...). "Worker zeigen / verstecken" ist jetzt ein echter Umschalter.
- **Gmail-Login direkt in AdiCode**: "Gmail-Login" oeffnet ein Login-Fenster in AdiCode. Der echte Tab-Browser laedt die Google-Anmeldung versteckt im Hintergrund, AdiCode zeigt ihn live und leitet Klicks, Tippen, Einfuegen und Scrollen weiter. Die Cookies landen im Container des Tabs. Fallback: "Echtes Fenster zeigen".
## Neu in 1.0.26
- **Usage-Dashboard (live)** auf der Seite Co-Work: zeigt fuer alle Konten das KI-Kontingent (alle 30 s, auch ohne laufendes Co-Work). Usage kommt jetzt direkt aus der Notion-API des eingeloggten Kontos (Fallback Seitentext); aufgebrauchte Konten werden rot markiert.
- **Selbst schreiben** pro Co-Work-Tab: stoppt den laufenden Worker-Lauf, zeigt das Fenster und pausiert die automatischen Erinnerungen, damit du in Worker-Chats selbst tippen kannst. **Freigeben** gibt den Tab zurueck an Co-Work. Hat ein Worker-Fenster den Fokus, schickt AdiCode keine Erinnerung hinein.
- **MCP-Connector** auf der Seite Notion AI: Tab auswaehlen und **AdiCode eintragen** (Settings -> Connections -> Custom MCP) oder **Fuer neue Tabs automatisch eintragen** aktivieren - jeder laufende, eingeloggte Tab ohne Verbindung bekommt AdiCode dann automatisch.
## Neu in 1.0.25
- Co-Work: Der Main-Agent bittet dich nicht mehr, selbst weitere Notion-AI-Tabs zu öffnen. Die Worker-Tabs von AdiCode treten automatisch bei. Nur wenn nach 2 Minuten niemand beigetreten ist, kommt ein Hinweis.
## Neu in 1.0.24
- Fixes aus dem ersten echten Co-Work-Test mit allen Konten (Todo-App: Coder und Tester haben parallel gebaut, sich gegenseitig reviewt und Fix-Runden gedreht).
- Worker-Startprompt neu und natürlich formuliert: Notion AI hatte den alten Prompt als „Prompt-Injection“ abgelehnt. Der Prompt erklärt jetzt auch, dass die MCP-Verbindung in Notion „asf“ oder „AdiCode“ heißen kann.
- Erinnerungen von AdiCode sind als Erinnerung im Namen des Nutzers formuliert statt als „System“-Befehl.
- Ist das KI-Kontingent eines Kontos aufgebraucht, sagt AdiCode das klar. Die anderen Worker starten trotzdem und übernehmen dessen Aufgaben.
- `safe_write_file` kennt jetzt den Parameter `agent`. Bisher konnten Co-Work-Agents ihre eigenen gesperrten Dateien nicht sicher schreiben.
- Rollen-Bewertung: Die Aufgaben warten jetzt bis zu 90 s auf die Worker, die gerade beitreten (vorher 20 s, dann haben sie die Bewertung verpasst).
## Neu in 1.0.23
- **Co-Work neu durchdacht – Rollen sind Schwerpunkte, keine Grenzen.** Coder, Reviewer und Tester helfen bei allen Aufgaben. Ist der Experte für eine Aufgabe beschäftigt (oder wartet sie länger als 90 s), übernimmt ein anderer Agent. So arbeiten wirklich alle vier gleichzeitig.
- **Jeder Prompt wird aus Sicht jeder Rolle bewertet:** beim Start und bei jeder neuen Nachricht mitten im Projekt (neues Tool `cowork_prompt`). Coder: Umsetzungsplan und Risiken. Reviewer: Qualitätsrisiken und Akzeptanzkriterien. Tester: Testplan und Edge Cases. Fehlende Teilaufgaben schlagen die Worker direkt vor. Die Aufgaben starten, sobald alle Bewertungen da sind (spätestens nach 150 s).
- **Weckfunktion:** Notion-AI-Chats bleiben nach jeder Antwort stehen. AdiCode erkennt jetzt, wenn ein Chat untätig ist, obwohl Arbeit wartet, und schickt ihm einen kurzen Weckruf („[AdiCode Co-Work] …“). Der Main-Chat wird geweckt, wenn Bewertungen oder Blocker da sind und wenn alles fertig ist (dann fasst er zusammen).
- Main bekommt beim Co-Work-Start jetzt klare Anweisungen (parallele Tasks, `cowork_prompt` für neue Nachrichten, selbst mithelfen, am Ende eine zusammengefasste Antwort).
- Fix im Chat-Tab: Die Ladeanzeige und die doppelte Nachricht blieben hängen, wenn beim Senden ein neuer Notion-Chat entstand.
## Neu in 1.0.22
- Neuer Tab **Chat**: eine ChatGPT-artige Chatbox direkt in AdiCode. Sie läuft über deinen Co-Work-Main-Tab (Notion AI) – das Main-Fenster kann dabei versteckt bleiben (Knopf „Main-Fenster zeigen/verstecken“).
- Fotos, Videos und Dateien anhängen: **+**-Knopf, Drag & Drop oder Strg+V (Screenshots). Vorschaubilder mit Upload-Anzeige, bis 512 MB pro Datei.
- **Sprechen:** Mikrofon-Knopf = Diktat (Notion transkribiert, Text landet im Eingabefeld). Wellen-Knopf = **Sprachmodus** wie bei ChatGPT: reden → automatisch senden → Antwort wird vorgelesen → wieder zuhören (Esc beendet).
- Antworten mit Formatierung, Code-Blöcken (Kopieren), „Gedankengang“-Schritten, Vorlesen, Freigabe-Knöpfen (Allow/Deny) und Stopp-Knopf. Lokaler Chat-Verlauf in der Seitenleiste.
- **Co-Work aus dem Chat:** „Co-Work starten“ öffnet alle Konten; danach schreibst du die Aufgabe einfach in den Chat. Der Co-Work-Status (Coder/Reviewer/Tester verbunden) steht oben.
- Fix: Die Erkennung „AdiCode-MCP bereits verbunden“ in Notion griff wegen kaputter Regex-Zeichen nie – jetzt behoben (erkennt auch die Verbindung „asf“).
- Hinweis: Der 1.0.21-Build wurde auf GitHub ohne Runner abgebrochen; 1.0.22 enthält alles aus 1.0.21.
## Neu in 1.0.21
- Echte automatische MCP-Einrichtung an die aktuelle Notion-Oberfläche angepasst: Workspace-Menü → Settings → Connections → Installed/Discover → Custom MCP, zweistufiger URL-/Bearer-Dialog und echte CDP-Mausklicks für Notions sicherheitsrelevanten Connect-Knopf.
- Beim Chatstart wird pro Account bevorzugt der vorhandene Verlauf **Welcome to Notion / Willkommen bei Notion** geöffnet. Fehlt er, wird bewusst ein neuer Chat benutzt – nie mehr versehentlich ein anderes laufendes Projektgespräch.
- Worker-Prompts werden über den echten „Submit AI message“-Knopf gesendet; Ctrl+Enter bleibt Fallback. Dadurch funktioniert der Kickoff auch bei der Notion-Einstellung „Enter fügt eine neue Zeile ein“.

## Neu in 1.0.20
- Co-Work-Hotfix aus dem echten Vier-Account-Test: Der Worker-Startprompt bleibt im zuvor geöffneten „Willkommen in Notion“-Chat und findet den Composer jetzt anhand Sichtbarkeit, Position und Notion-AI-Platzhalter statt einer falschen Bildschirmhöhen-Grenze.
- Damit funktionieren auch mittig platzierte Eingabefelder und unterschiedliche Notion-Onboarding-Layouts; der Kopier-Button bleibt Fallback.

## Neu in 1.0.19
- Sicherer Co-Work-Start: Alle vier Konten werden zuerst einmalig nach **Chat / Willkommen in Notion** geführt. AdiCode prüft pro Tab, ob das sichtbare Notion-AI-Eingabefeld wirklich bereit ist. Erst wenn alle vier erfolgreich sind, werden Coder, Reviewer und Tester versteckt.
- Rollen-Handshake: Nach dem Main-Kickoff müssen Coder, Reviewer und Tester dem neuesten Board nachweislich über `cowork_join` beitreten. Der Live-Status zeigt den Fortschritt; bei Fehlern werden alle Fenster automatisch sichtbar statt still weiterzulaufen.
- Präzisere Worker-Schleife: feste Agent-Namen und Rollen, selbstständiges `claim → arbeiten → update → wait` sowie Blocker-Nachrichten an Main.

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
