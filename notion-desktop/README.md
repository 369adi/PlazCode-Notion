# PlazCode Notion Desktop / AdiCode (eigenständige Weiterentwicklung, GPL-3.0-Basis)

Steht in keiner Verbindung zu den PlazCode-Autoren oder zu Notion Labs. Alle Original-Lizenz-,
Copyright- und Branding-Hinweise (LICENSE, BRANDING-NOTICE.txt, THIRD-PARTY-NOTICES) bleiben erhalten.

## Was ist anders?
- Name „AdiCode" (technisch PlazCode Notion), dunkles Design.
- Neue Seite **Notion AI**: eingebauter MCP-Server (127.0.0.1:8787/mcp, Token via Header oder Pfad) + ngrok-Tunnel mit fester Domain — Notion AI direkt, ohne Browser-Erweiterung.
- Notion sieht `roblox_*`, `pc_*`, `browser_*` und alle weiteren aktivierten MCP-Server.
- Automatische Updates über den eigenen Launcher (Original-Update deaktiviert).

## Neu in 1.0.129

- Co-Work stoppen beendet jetzt wirklich alle Co-Work-Tabs (auch wenn sie in Edge/Chrome liefen) und schliesst die Notion-Tabs in allen offenen Browsern.
- Alle starten startet alle Co-Work-Tabs in Brave (fremde Browser auf den Co-Work-Ports werden vorher beendet).
- Jetzt aktualisieren wirkt sofort: der Launcher laedt neue Versionen schon im Hintergrund vor, kein 30-Sekunden-Warten mehr mit verstecktem Fenster.
- Chats verbinden sich automatisch mit dem AdiCode-MCP und machen nach Verbindungsabbruechen selbst weiter.

## Neu in 1.0.128

- Weniger MCP-Abbrueche bei vielen Chats: Notions Dauer-Verbindung (GET /mcp) wird jetzt als ruhiger SSE-Stream mit Keep-Alive gehalten statt alle paar Sekunden abgelehnt (halbiert die Tunnel-Last).
- ngrok-Fehler landen im Log, Abbrueche sind nachvollziehbar.
- Co-Work laeuft nur noch in Brave: alte Edge/Chrome-Fenster mit Co-Work-Konto werden automatisch beendet und in Brave neu geoeffnet.
- Co-Work-Login oeffnet im selben Coworker-Brave, kein zweites Fenster mehr.
- Jedes Konto wird automatisch als Chat und Coworker erkannt.
- Update-Knopf zeigt an, wenn ein Update vorgeladen und sofort bereit ist.

## Neu in 1.0.127
- **Chats holen sich immer Hilfe**: Jede Aufgabe mit mehreren unabhaengigen Teilen wird sofort per Co-Work auf mehrere Chats verteilt; grosse Teilaufgaben werden automatisch weiter zerlegt.
- **Stabil mit vielen Chats gleichzeitig**: Add-on-Pool (pc, codebase, files, git ... bis zu 3 parallele Instanzen), Builds laufen mit niedriger Prioritaet und halben Kernen, ngrok-Tunnel mit hoher Prioritaet.
- **Weniger Tokens**: Co-Work-Schemas 58 % kleiner, kompakte Ausgaben fuer codebase/files/git.
- **Helfer-Chats**: warten bis zu 40 s auf Bereitschaft statt abzubrechen.
## Neu in 1.0.122
- **Notion-Tabs immer auf Chat**: AdiCode stellt die Seitenleiste oben links in jedem Notion-Tab automatisch auf *Chat* zurueck.
- **MCP-Hinweis fuer neue Konten**: Hat ein Konto AdiCode noch nie benutzt, beginnt sein Co-Work-Prompt mit "Verbinde dich mit dem AdiCode MCP". Ab dem ersten AdiCode-Aufruf faellt der Hinweis weg.
- **Login mit jeder E-Mail**: Neben Gmail gehen jetzt auch andere Adressen (z. B. Proton oder Firmen-Konten). AdiCode traegt die E-Mail im Notion-Login ein und klickt *Weiter*; du gibst nur noch den Code aus deinem Postfach ein. Login-Art und Sitzung bleiben im Konto gespeichert.
- **Co-Work-Login-Port-Fix**: Fester Port pro Tab, damit die Login-Vorschau kein fremdes Konto zeigt.
- **Updates**: Vor einem Update werden alle aktiven Notion-Tabs benachrichtigt, danach wird aktualisiert und neu gestartet.

## Neu in 1.0.121
- **Co-Work nutzt alle freien Chats**: Jedes laufende Konto setzt automatisch alle bestehenden, freien Notion-Chats als Co-Work-Chats ein (nacheinander, mit RAM-Check). Chats, die gerade arbeiten, werden uebersprungen; es wird nie ein neuer Chat erstellt.
- **Fester CDP-Port pro Tab**: Port haengt an der Tab-ID statt an der Reihenfolge - Erinnerungen und Weckrufe landen nicht mehr im falschen Konto.

## Neu in 1.0.120
- **Lokales Modell im Chat**: Im Chat-Kopf waehlst du jetzt *Notion AI* oder ein lokales Ollama-Modell (z. B. Qwen 3.5). Ollama wird bei Bedarf automatisch im Hintergrund gestartet.
- Das lokale Modell bekommt automatisch die AdiCode-Tools (PowerShell, Dateien, Programme, Screenshot, Hintergrund-Jobs, Tool-Suche fuer alle weiteren Add-ons) plus Websuche und Webseiten-Abruf - schlanke Auswahl, spart RAM.
- **Rueckfrage nur bei Unsicherheit**: Das Modell arbeitet selbststaendig und fragt nur dann mit *Erlauben / Ablehnen*, wenn es selbst unsicher ist, ob eine Aktion riskant ist. Stopp und Neuer Chat jederzeit moeglich.

## Neu in 1.0.119

- Schwarze leere Fenster der Notion-Tabs (Edge-Hilfsfenster) werden automatisch ausgeblendet.

## Neu in 1.0.118
- **Live-Vorschauen optional**: Tab-Previews sind standardmaessig aus und lassen sich in den Einstellungen einschalten (weniger Last).
- **Tabs verstecken entfernt**: Notion-Fenster verschwinden nicht mehr von selbst.
- **MCP-Rechte schlanker**: Nur Tab `Notion Agent` -> Write auf `Run automatically`.
- **Einzelschritte**: Jedes Check-Haekchen ist klickbar und fuehrt nur diesen Schritt aus; Checks lassen sich stoppen.
- **Schnellerer, klarer Check**: Kuerzere Wartezeiten, Anzeige `Schritt X/7 ... laeuft seit N s`, gemerkte Schritte werden uebersprungen.
- **Login-Fenster** ist wieder sichtbar und bedienbar.
## Neu in 1.0.117
- **Co-Work Auto-Join**: Neue Tabs, die noch keinem Projekt beigetreten sind, werden automatisch eingeladen, dem neuesten Projekt mit offenen Tasks beizutreten, und uebernehmen danach Tasks (max. 3 Erinnerungen, mind. 90 s Abstand, Main-Chat nie).
- **Weniger Update-Hinweise**: Der Hinweis auf neue Versionen erscheint nur noch selten.
## Neu in 1.0.116

- **AdiCode erstellt nie einen neuen Notion-Chat:** Co-Work und jeder Chat bedienen nur bestehende Chats. **+ Chat** holt einen freien *bestehenden* Chat des Kontos dazu (zuerst Tab-n-Kontext-/Willkommen-Chats, dann zuletzt genutzte).
- Harte Sperren im Code: Tab oeffnen und Senden nur, wenn die Seite ein bestehender Chat ist (Link mit t=). Sonst wird nichts gesendet.
- Prompt-Enhancer ueber Notion AI entfernt (legte jedes Mal einen neuen Chat an), jetzt nur lokal. Knoepfe *Neuer Chat* / *Neu + Zusammenfassung* im AdiCode-Chat entfernt.

## Neu in 1.0.115

- **Co-Work ohne Rollen und ohne Main:** Die festen Rollen (Main, Coder, Reviewer, Tester ...) sind komplett raus und ueberschreiben beim Start nichts mehr. Der Chat, in dem du zuerst schreibst (beliebiger Notion-Tab oder AdiCode-Chat), leitet das Co-Work; die anderen Tabs treten automatisch bei und helfen bei allem.
- **Nie mehr ein neuer Chat:** AdiCode schreibt nur noch in einen bestehenden Chat (Willkommen bei Notion). Die automatische Nachricht "Verbinde dich mit AdiCode MCP Server" am Ende des Refresh entfaellt.
- **Aktualisieren = sofort:** Der Knopf schliesst AdiCode direkt und installiert das Update. Chats werden dabei nicht mehr gestoppt, nach dem Neustart wird kein "weiter" mehr geschrieben.
- **Fertig-Ton nur, wenn der Chat wirklich fertig ist:** Der Ton kommt erst, wenn der Chat einige Sekunden lang nicht mehr arbeitet (kurze Pausen zwischen Tool-Aufrufen zaehlen nicht).
- **Kein Flackern mehr auf der Update-Seite:** Das gruene "installiert"-Badge bleibt bei jeder Pruefung ruhig; bei GitHub-Fehlern bleibt die Liste stehen.

## Neu in 1.0.114

- **Mehrere Chats pro Konto im Co-Work:** Knopf **+ Chat** holt einen weiteren Chat desselben Kontos dazu (eigener Hintergrund-Tab, eigener Agent-Name wie Tab1-K2). Bis zu 8 Chats pro Konto; neue Chats hoechstens alle 20 s (Notion verwirft sonst neue Chats), RAM-Waechter (mind. 1,2 GB frei). **- Chat** nimmt den letzten Zusatz-Chat wieder raus.
- **Live-Chat-Zaehler:** Alle Chats eines Kontos werden direkt aus den Notion-Daten gelesen (alle 10 s): Anzahl, wie viele gerade laufen und welche als Tab offen sind (Tooltip mit Titeln).
- **Live-Usage aus dem echten Limit:** Notion-KI-Kredite im 6-Stunden-Fenster (z. B. 44/100) inkl. Reset-Zeit; Usage-Balken nutzt diesen Wert.
- Getestet: 10 bestehende Chats eines Kontos laufen gleichzeitig stabil; mehr als ca. 5 *neue* Chats kurz hintereinander verwirft Notion; pro Chat-Tab ca. 350-400 MB RAM.

## Neu in 1.0.113
- **Weniger Tokens, schnellere Antworten**: Ungenutzte bzw. doppelte Add-ons entfernt (Memory, Sequential Thinking, Serena, SQLite, Time, PostgreSQL, Sentry) - auch aus bestehenden Konfigurationen. Code-Navigation laeuft ueber codebase, Gedaechtnis ueber project_memory. Anleitung und Toolbox entsprechend gekuerzt.
- **Notion-Fenster immer sichtbar**: Co-Work- und Haupt-Tabs werden nicht mehr versteckt; von aelteren Versionen versteckte Fenster erscheinen beim Start wieder.
- **1.0.111 wiederhergestellt**: Codebase immer sichtbar, Auto-Index und ADR im Projekt-Gedaechtnis waren in 1.0.112 versehentlich verloren gegangen.

## Neu in 1.0.112
- **Co-Work-Tabs**: Beim Sichtbarmachen (Alle zeigen) erscheint keine "Wiederherstellen?"-Leiste mehr - das Profil wird vor dem Start als sauber beendet markiert.
## Neu in 1.0.111
- **Codebase immer sichtbar**: Die codebase-Tools (schnelle Code-Suche und -Lesen) sind jetzt in jedem Tool-Profil direkt verfuegbar statt versteckt.
- **Auto-Index**: project_memory_open und project_memory_update indexieren das Projekt-Repo (repo_path) automatisch inkrementell in codebase neu - der Index ist immer aktuell.
- **Codebase zuerst**: Die Anleitung fuer Notion AI sagt jetzt klar: Code immer zuerst ueber codebase_* lesen und suchen, pc/files nur zum Schreiben.
- **ADR statt Pflicht-HANDOUT**: project_memory_open haengt die Projekt-ADR aus codebase an; Entscheidungen werden per codebase_manage_adr gespeichert. HANDOUT.md ist nur noch optional.

## Neu in 1.0.110
- **Co-Work-Dashboard**: Status- und Usage-Abfragen laufen mit zufaelligem Abstand (+-20 %), werden bei Fehlern schrittweise langsamer (bis 60 s) und pausieren, solange das AdiCode-Fenster minimiert oder verdeckt ist.
## Neu in 1.0.109
- **Co-Work nur noch ueber die Bridge**: Die Browser-Erweiterung ist komplett entfernt (Code, Tools, Einstellungen). Alle Tabs laufen direkt ueber AdiCode auf diesem PC.
- **Main-Fenster-Check**: cowork_start warnt, wenn das als Main eingetragene Fenster gar nicht gestartet ist oder Fenster offline sind, und liefert einen Einfuegetext mit Projekt-ID. Der eigentliche Lead-Chat wird nicht mehr faelschlich als Worker eingetragen.
- **Ein Konto, ein Fenster**: Dasselbe Notion-Konto kann nicht mehr in zwei Fenstern laufen - das vorhandene Fenster wird fokussiert; doppelte E-Mails werden beim Speichern gemeldet.
- **Release-Zug**: Jeder Tab darf ohne Nutzer-OK veroeffentlichen. Andere Tabs haengen sich mit action=ready an; arbeiten noch Tabs an der Version, faehrt der Release nach hoechstens 10 min ab, Nachzuegler kommen automatisch in die naechste Version. Bereits veroeffentlichte oder gerade gebaute Versionen leiten automatisch weiter - es laeuft immer der neueste Stand.

## Neu in 1.0.108
- **Blender**: Das doppelte Add-on blenderwright ist entfernt - AdiCode nutzt nur noch MCP for Blender (blender-mcp), keine Port-Konflikte/30-s-Timeouts mehr.
- **Blender-Helfer 1.3**: Kamera und Lichter zielen korrekt aufs Motiv (kein graues Erst-Render mehr), Kamera rahmt auch hohe Objekte komplett, Boden/Hintergrund wird bei Framing und Lichtstaerke ignoriert (keine Ueberbelichtung), Holz-Material sieht nach Holz aus, neue Material-Presets (Messing, Gold, Kupfer, Stahl, Chrom, Glas, Kunststoff, Gummi, Keramik, Holz, Stoff, Beton ...), Grundformen mit Radius/Tiefe.
- Blender-Skills korrigiert (Vorschau-Render, Presets, Kamera-Hinweise).## Neu in 1.0.107
- **Hintergrund-Jobs**: Braucht ein Tool laenger als das Soft-Limit (Standard 20 s, einstellbar unter Notion AI), antwortet AdiCode sofort und der Befehl laeuft weiter - Notion wartet nicht mehr, sondern arbeitet parallel weiter. Das Ergebnis haengt sich automatisch an die naechste Tool-Antwort; gezielt mit dem neuen Tool adicode_job (id, wait_s).
- Notion-AI-Anweisungen: unabhaengige Tool-Aufrufe parallel starten, bei Hintergrund-Job nicht warten oder neu starten.## Neu in 1.0.106
- **Refresh laeuft komplett durch**: Der MCP-Manage-Schritt klickt jetzt selbst auf Manage in der AdiCode-Zeile und stellt Read- und Write-Tools auf Run automatically. Du musst nicht mehr eingreifen.
- **Kein Warten mehr**: Schliesst sich der Google-Login, macht der Refresh automatisch mit Zum Chat weiter. Ein bereits eingeloggter Tab wird sofort erkannt.
- **Chats heissen Tab N Kontext 1/2/3**: hoechstens so viele Tabs, wie Co-Work-Fenster offen sind. Bereits richtig benannte Chats bleiben unveraendert.
- **Letzter Schritt**: In den Chat wird Verbinde dich mit AdiCode MCP Server. geschrieben.
- **Neue Co-Work-Haken**: Run automatically, Chats benannt, MCP gestartet.
## Neu in 1.0.105
- **Chaterkennung**: Kein Chat wird nur noch angezeigt, wenn Notion meldet: KI fuer diesen Workspace deaktiviert (EN/DE). Jeder vorhandene Chat zaehlt, nicht nur Willkommen bei Notion; es wird nie ein neuer Chat angelegt.
- **Co-Work-Haken stimmen**: MCP configured und Modell nutzen die gemerkten Checks, erledigte Schritte stehen nicht mehr als offen da.
- **Live-Usage**: Usage jedes laufenden Tabs wird ca. alle 30 s still direkt von Notion gelesen (ohne Einstellungsfenster), Anzeige live HH:MM:SS. Pausierte Tabs fragen nicht ab.
## Neu in 1.0.104
- **Blender-Screenshots ohne Fenster**: Neues Tool blender_capture rendert Ansichten der Szene (Perspektive, vorne, rechts, oben ...) als Kontaktbogen mit kleiner Vorschau direkt in den Chat. Laeuft Blender (auch minimiert), wird das genutzt, sonst startet AdiCode unsichtbar ein Hintergrund-Blender. Blender muss nicht offen sein.
## Neu in 1.0.103
- **Release-Zug fuer mehrere Tabs**: Neues Tool release_queue. Hat ein Tab fertige, getestete Aenderungen, waehrend ein anderer gerade baut, haengt er sie an dessen naechstes Release an statt parallel zu releasen. Konflikte (gleiche Datei inzwischen geaendert) werden erkannt, ein Release startet erst, wenn alle wartenden Beitraege uebernommen sind.
- **Robustere Releases**: Kein doppelter Upload mehr durch den automatischen GitHub-Build, halb hochgeladene Dateien werden ersetzt und der Upload bis zu 3x versucht.
## Neu in 1.0.102
- **MCP-Check haengt nicht mehr**: Nach dem Einrichten wechselt AdiCode selbst auf Browse Connections, kein manueller Klick mehr noetig.
- **Run automatically zuverlaessig**: Der Tab Manage/Verwalten wird erkannt, dann AdiCode -> Notion Agent -> alle Berechtigungen auf Run automatically.
- **Chats pro Konto benannt**: Die 3 Willkommen-Chats heissen jetzt Tab N, Tab N Kontext 2, Tab N Kontext 3. Jedes Konto bekommt eine eigene Nummer (kleinste freie), keine doppelt.
## Neu in 1.0.101
- Logger fuer 'Notion AI was disabled in this workspace': steht der Hinweis in einem Konto-Tab, schreibt AdiCode einmal einen Bericht mit allem, was kurz vorher passiert ist (Zeitprotokoll, Tool-Aufrufe aller Chats, App-/Launcher-/MCP-Logs, Screenshot) nach %LOCALAPPDATA%\PlazCodeNotion\logs\ai-disabled

## Neu in 1.0.100
- **Updates nur noch per Klick**: AdiCode aktualisiert sich nicht mehr von selbst (weder beim Start noch im Hintergrund). Neue Version -> in AdiCode auf Aktualisieren klicken.
- **Sicheres Update**: Vor dem Neustart prueft AdiCode alle Notion-Tabs. Steht noch ungesendeter Text im Chat, wird das Update abgebrochen (erst abschicken/loeschen). Laufende Antworten werden sauber gestoppt.
- **Automatisch weiter**: Nach dem Neustart schreibt AdiCode in genau die gestoppten Chats "weiter", sobald die Verbindung wieder steht und der Chat frei ist.
## Neu in 1.0.99
- **Co-Work-Sperren repariert**: Es werden nur noch echte Dateipfade gesperrt. Pseudo-Ziele wie "0", "&1" oder "nul" (aus Umleitungen wie 2>&1) blockieren nicht mehr jeden PowerShell-Befehl anderer Chats.
- **Serena versteckt**: Serena laeuft weiter im Hintergrund (sonst fehlen die serena_*-Tools), oeffnet aber kein Browser-Dashboard und kein Log-Fenster mehr.
- **Blender nach Referenzbild**: neue Helfer measure_ref, ref_camera, compare (Umriss-Aehnlichkeit + Vergleichsbild), ortho_views; der Blender-Skill arbeitet in Runden bis der Umriss passt.
- **Blender-Profil automatisch**: Laedt die KI den Blender-Skill, zeigt AdiCode nur noch die passenden Tools (ohne Roblox, blenderwright, Code-Navigation) - weniger Kontext, bessere Ergebnisse.
## Neu in 1.0.98
- Alle offenen Notion-Chats sind immer verbunden, auch ohne Co-Work: cowork_message geht an jeden Chat (to = Name, z. B. Main1), freie Tabs werden mit der Nachricht geweckt
- Jeder Chat bekommt beim ersten AdiCode-Aufruf eine Uebersicht: welche Chats offen sind, was sie gerade tun und wartende Nachrichten
- Nachrichten werden nicht mehr nach 120 Zeichen abgeschnitten

## Neu in 1.0.97
- Gleicher Stand wie der zweite 1.0.96-Build (lean Tool-Profil, pc_job, Add-on-Watchdog, cowork_release) – neu nummeriert, damit alle mit dem ersten 1.0.96 das Update auch bekommen
- Co-Work-Fixes: Main trägt Gruppen-Tabs automatisch ein, Prompts werden zuverlässig abgeschickt, Wecker/Keeper laufen auch nach Neustart weiter
- Release-Status für alle Chats sichtbar (adicode_status + einmaliger Hinweis): wer gerade welche Version baut/veröffentlicht
- Release-Schutz: schon veröffentlichte Versionen können nicht mehr versehentlich überschrieben werden

## Neu in 1.0.96
- Konto-Check bleibt nicht mehr bei der MCP-Einrichtung hängen: pro Tab nur ein Check gleichzeitig, Zeitlimits, zweiter Versuch, offene Dialoge werden geschlossen
- Nach dem Hinzufügen von AdiCode werden alle Berechtigungen (Manage, Agent, Write skill, Read/Write tools …) auf „Run automatically“ gestellt – auch bei schon verbundenen Konten

## Neu in 1.0.89
- Chat-Spiegel robuster (Notion hat interne Strukturen geändert), Live-Anzeige der Tabs, präzisere Datei-Reservierung

## Neu in 1.0.88
- Datei-Reservierung ohne Timer: eine Datei ist belegt, solange ein Chat live daran arbeitet

## Neu in 1.0.87
- Live-Anzeige auf Deutsch im AdiCode-Chat; Tabs sehen sich gegenseitig (Co-Work-Bridge immer an)

## Neu in 1.0.86
- Fix Senden-/Stopp-Knopf-Erkennung (Notion hat test-IDs geändert); Chat-Design v2

## Neu in 1.0.85
- Co-Work-Bridge (Datei-Locks, Arbeitspakete)

## Neu in 1.0.84
- Neues UI-Design v2 (Animationen, Glow, Navigation; respektiert reduzierte Animationen)

## Neu in 1.0.83
- Notion ohne Chat-Leiste:自动 wieder „Willkommen bei Notion“ öffnen (statt 10 s warten: direkt neu laden)

## Neu in 1.0.82
- Bis zu 8 Co-Work-Tabs, Check pro Konto, Account-Cache

## Neu in 1.0.81
- Co-Work direkt aus dem Chat starten (Unterhaltung „Co-Work" im Willkommen-Chat)

## Neu in 1.0.80
- „Neuer Chat" im AdiCode-Chat startet eine neue Unterhaltung im Chat „Willkommen bei Notion" (nie ein neuer Notion-Chat)

## Neu in 1.0.79
- Chat-Liste wie bei ChatGPT; Co-Work-Kontext im Main-Chat; alles grün = Fenster weg

## Neu in 1.0.78
- Kein endloses Laden (Watchdog lädt nach 20 s neu)

## Neu in 1.0.77
- Chat auf Vollbild umstellen; Verbindungen-Seite in allen Layouts

## Neu in 1.0.76
- Notion immer auf Englisch (US); „MCP schon vorhanden" gilt als verbunden

## Neu in 1.0.75
- Einstellungen zuverlässiger öffnen; MCP-Dialog autom. Name + Token

## Neu in 1.0.74
- ngrok-Token bleibt beim Wechsel auf eigenes ngrok erhalten

## Neu in 1.0.73
- ngrok immer frisch installiert statt alter PATH-Version; Authtoken maskiert

## Neu in 1.0.72
- Store-ngrok wird erkannt; Schnellstart ohne Token-Vorabsperre

## Neu in 1.0.71
- ngrok-Schnellstart: nur Authtoken + Start (installation, Gratis-Domain, Tunnel, Test automatisch)

## Neu in 1.0.70
- Worker-Fenster bleiben versteckt (Fenster-Wächter); ngrok-Fehler lesbar

## Neu in 1.0.69
- Token in allen Formaten akzeptiert (Bearer, mit/ohne Präfix, x-api-key, ?token=); URL auch ohne /mcp

## Neu in 1.0.68
- MCP-URL mit Schlüssel (funktioniert auf jedem PC); Custom-MCP-Menüpunkt; Chat-Modell-Knopf

## Neu in 1.0.67
- Onboarding-Hänge-Fix (Skip-to-content); sichtbares Google-Konto

## Neu in 1.0.66
- Tab aussetzen/mitmachen; Zeitprotokoll timeline.log

## Neu in 1.0.65
- Check ohne Cache; Chat-Modell wird umgestellt; Token-Modus startet Check; Abo-/Freischaltungs-Status erkannt und automatisch verwaltet (KI-Konnektoren)

## Neu in 1.0.64
- Token-Modus-Dropdown sichtbar

## Neu in 1.0.63
- Aufwecken-Button (nur auf Klick); schnellerer Google-Login; kurze Memory-Antworten

## Neu in 1.0.62
- Keine automatischen Erinnerungen mehr; nur Willkommen-Chat; Token-Modus (wenig/viel)

## Neu in 1.0.61
- Co-Work deutlich schneller (kürzere Gates, sofortiges Aufwachen, parallele Startprompts)

## Neu in 1.0.60
- Google-Fehler 400 behoben (Login über notion.so/login); deutsches Onboarding; nur ein Willkommen-Chat

## Neu in 1.0.59–1.0.46
- Auto-Login-Kette: Google-Konto zur E-Mail wählen, Popup-Blocker aus, Cookie-Dialoge, nicht eingeloggte Konten im Check automatisch anmelden; E-Mail-Ändern speichert automatisch

## Neu in 1.0.45 / 1.0.44
- Chat-Uploads (Ordner, Auto-Bereinigung); Edge: kein zweites Fenster

## Neu in 1.0.43
- Onboarding wählt „For work"

## Neu in 1.0.42
- Auto-Einrichtung nach Session-Erneuerung (Onboarding neu)

## Neu in 1.0.41
- Live-Gedankengang von Main im AdiCode-Chat

## Neu in 1.0.40
- Refresh öffnet danach Gmail-Login mit derselben E-Mail (Session-Erneuerung)

## Neu in 1.0.39
- Usage-Check robuster (CDP-Zeitlimit); Refresh-Button auch bei Konten ohne Chat; Update-Notizen sauber untereinander

## Neu in 1.0.38
- Chat-Check erkennt Konten ohne Chat („Kein Chat")

## Neu in 1.0.37
- „Prüfe Konten" hängt nicht mehr bei 58 %

## Neu in 1.0.36
- Co-Work-Prompts gekürzt (weniger Token)

## Neu in 1.0.35
- Usage-Fix: echter Monatsverbrauch wird gelesen

## Neu in 1.0.90
- **Alle Fenster live**: Neuer Knopf "Alle Fenster" im Chat zeigt alle offenen Co-Work-Tabs gleichzeitig (2x2). Klicken, Scrollen und Tippen geht direkt in jeder Vorschau.
- **Main frei waehlbar**: In der Tab-Leiste im Chat werden alle Tabs angezeigt; per "Als Main" (Stern) wird ein Tab zum Main-Chat.
- **Prompt-Enhancer**: Der Zauberstab-Knopf im Chat schreibt deinen Prompt per KI verstaendlicher um (Fallback lokal). Shift+Klick stellt den Originaltext wieder her.
- **Co-Work starten**: Der Knopf ist waehrend des Starts ausgegraut und schickt keinen Prompt mehr automatisch.
- **Sofort-Update**: "Jetzt aktualisieren" installiert ein verfuegbares Update sofort (Knopf zeigt die neue Version an). Laeuft der Launcher nicht, laedt AdiCode das Update selbst.
- Neuer Chat denkt nicht mehr sofort los; Chat allgemein stabiler.
## Neu in 1.0.91
- Mehrere Mains: jeder Tab kann Main sein und leitet ein eigenes Projekt. Tabs werden per Auswahlfeld im Fensterkopf einem Main zugewiesen (z. B. Main 1 + Tab 1, Main 2 + Tab 2).
- Keine festen Rollen mehr (Coder/Reviewer/Tester entfernt): jeder Tab hilft bei allem und schickt Ergebnisse und Kritik an seinen Main, der so lange neue Runden ausgibt, bis das Ergebnis perfekt ist.
- Fenster-Wall: alle Tabs gleichzeitig live sehen und steuern, Groessen per Trennlinien ziehen, Doppelklick zum Vergroessern, Vollbild-Modus, Groessen werden gespeichert.
- Neuer Chat im AdiCode-Chat startet einen eigenen Notion-AI-Chat nur fuer AdiCode - Co-Work-Chats der Tabs bleiben unberuehrt.
## Neu in 1.0.92
- **Blender Pro**: Drei neue Experten-Skills (adicode_skill blender, blender-model, blender-anim) mit Arbeitsablauf, Qualitaets-Massstab und Pflicht-Kontrolle per Render/Screenshot.
- **Helfer-Bibliothek adicode_blender**: wird automatisch in Blender installiert (import adicode_blender as A). Studio-Licht, Kamera-Framing, Hard-Surface/Organic-Finish, prozedurale PBR-Materialien, Mesh-Qualitaetscheck, Scatter per Geometry Nodes, Bounce/Loop/Shake/Orbit/Follow-Path-Animation, Rigify-Rig + Auto-Weights, Mixamo/BVH-Import, GPU-Render, Video- und glb/fbx-Export.
- **Neues Tool blender_pro**: doctor (Blender, Addon-Port, GPU, Extensions pruefen), setup (LoopTools, Bool Tool, Extra Objects, Rigify, Node Wrangler, Cycles OptiX), install_helpers, helpers.
- **MCP-Add-on blenderwright** (191 Blender-Tools) als Alternative im Katalog.
- **project_memory_delete**: Projekt-Kontexte lassen sich dauerhaft loeschen (werden nicht wieder importiert).
## Neu in 1.0.93
- Fix: Notion-Tabs (z. B. trinix1337) wurden alle paar Sekunden neu geladen. Ursache: Die neue Notion-Sidebar hat keinen Chat-Tab mehr, AdiCode hielt das fuer ein kaputtes Layout und lud endlos neu. Jetzt wird nur noch neu geladen, wenn wirklich keine Chat-Oberflaeche da ist, hoechstens 2x in 10 Minuten und nie, waehrend du im Fenster arbeitest.
- Fix: Die Seite MCP Servers wurde jede Sekunde komplett neu gezeichnet. Jetzt aendert sie sich nur noch, wenn sich wirklich etwas geaendert hat.
- Co-Work: Steht bei einem Tab rot Chat fehlt, fuehrt AdiCode Zum Chat jetzt selbst aus (pro Tab hoechstens alle 2 Minuten).
- Enthaelt alles aus 1.0.92: Projekt-Chatnamen, Fertig-Ton, groessere Vorschau ohne schwarze Balken, echtes Vollbild, Ansichtsmodi, Gruppenfarben.
## Neu in 1.0.123
- **Co-Work erkennt Chats, die nicht schreiben koennen**: Tabs und Zusatz-Chats mit dem Hinweis *KI ist fuer diesen Workspace deaktiviert*, geschlossene Fenster und geschlossene Zusatz-Chats werden nicht mehr eingetragen oder geweckt. Ihre Aufgaben werden fuer andere Chats wieder frei. Ist der Hinweis weg, machen sie automatisch wieder mit.
- **Jeder offene Chat hilft bei jedem Co-Work**: Egal, was vorher im Chat stand - jeder freie Chat (ausser er arbeitet gerade) tritt neuen Co-Work-Projekten bei und holt sich Aufgaben. Solange Aufgaben offen sind, holt jedes Konto alle 15 s einen weiteren freien Chat dazu.
## Neu in 1.0.124
- **Co-Work verteilt Aufgaben endlich auf mehrere Chats**: Laufende Konto-Tabs wurden teils als *Tab laeuft nicht* erkannt, deshalb wurden nie Zusatz-Chats geholt. Ein Tab gilt jetzt auch als laufend, wenn sein Fenster antwortet - freie Chats bekommen sofort Aufgaben.
- **Tabs ohne AdiCode-Verbindung werden ausgelassen**: Ist in einem Konto die AdiCode-MCP-Einrichtung fehlgeschlagen, wird es nicht mehr eingetragen oder geweckt, statt Aufgaben zu blockieren.
## Neu in 1.0.125
- **Co-Work-Team neu geschrieben**: Eine einzige Schleife pro Konto ersetzt die alten, sich ueberschneidenden Mechanismen. Freie Chats (auch ganz neue ohne AdiCode-Verlauf) werden automatisch zum offenen Co-Work eingeladen und verbinden sich selbst mit AdiCode.
- **Sichtbarer Status pro Konto**: Im Konto steht jetzt, wie viele Chats im Team sind und warum ein Chat nicht mitmacht (z. B. gerade offen im Fenster oder KI deaktiviert).
- **Kein falscher Alarm mehr bei KI deaktiviert**: Erkannt wird nur noch der echte Notion-Hinweis, nicht mehr eigener Chat-Text.
- **Co-Work hilft weiter**: Agents ohne freien Task bieten Hilfe an; Co-Work-Nachrichten laufen direkt ueber die Live-Bridge.
## Neu in 1.0.126
- **Alle Chats eines Kontos machen mit**: Auch der Start-Chat und Chats, die gerade im Fenster offen sind, werden zum Co-Work eingeladen. Ausgelassen wird nur ein Chat, der gerade generiert oder in dem du wirklich gerade tippst.
- **Kein falsches 'Nutzer ist gerade in diesem Chat' mehr**: Hintergrund-Tabs gelten nicht mehr als 'vom Nutzer benutzt'.
- **Notion-Tabs nur noch in Brave**: Chrome/Edge werden fuer Konto-Tabs nicht mehr verwendet. Jedes Konto hat weiterhin einen eigenen Profilordner mit eigenen Cookies und bleibt eingeloggt.
- **Weniger Lag**: Tracking-, Analyse- und Fingerprint-Skripte (Intercom, Google Tag Manager, Sentry, Datadog, Splunk, FingerprintJS, Hotjar, Statsig, Pendo u. a.) werden in den Notion-Tabs geblockt.
## Neu in 1.0.130
- **Zwei Konten arbeiten zusammen**: Chats aus verschiedenen Notion-Konten koennen demselben Co-Work-Projekt beitreten, Tasks uebernehmen und sich Nachrichten schicken. Chat-Namen (z. B. AdiCode-Tab3-K2) werden nie mehr auf einen anderen Chat umgebogen.
- **Co-Work-Tab nur noch Brave**: Browser-Auswahl zeigt nur Brave, Edge/Chrome-Texte und -Icons sind raus. Fehlt Brave: 'Brave nicht gefunden - bitte installieren'.
- **Neuer Knopf 'Aktualisieren'**: laedt Tabs, Konten und Ports im Co-Work-Tab sofort neu.
- **Aufgeraeumt**: Edge/WebView2-Altlasten im Co-Work-Backend entfernt.
## Neu in 1.0.34–1.0.23
(fixiert in 1.0.29: lokale schnelle Releases; Kompaktversion der alten Einzelnachrichten)

## Neu in 1.0.22 und älter
(siehe HANDOUT.md Changelog — unverändert übernommen)
