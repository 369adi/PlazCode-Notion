# PlazCode Notion Desktop / AdiCode (eigenständige Weiterentwicklung, GPL-3.0-Basis)

Steht in keiner Verbindung zu den PlazCode-Autoren oder zu Notion Labs. Alle Original-Lizenz-,
Copyright- und Branding-Hinweise (LICENSE, BRANDING-NOTICE.txt, THIRD-PARTY-NOTICES) bleiben erhalten.

## Was ist anders?
- Name „AdiCode" (technisch PlazCode Notion), dunkles Design.
- Neue Seite **Notion AI**: eingebauter MCP-Server (127.0.0.1:8787/mcp, Token via Header oder Pfad) + ngrok-Tunnel mit fester Domain — Notion AI direkt, ohne Browser-Erweiterung.
- Notion sieht `roblox_*`, `pc_*`, `browser_*` und alle weiteren aktivierten MCP-Server.
- Automatische Updates über den eigenen Launcher (Original-Update deaktiviert).

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
## Neu in 1.0.34–1.0.23
(fixiert in 1.0.29: lokale schnelle Releases; Kompaktversion der alten Einzelnachrichten)

## Neu in 1.0.22 und älter
(siehe HANDOUT.md Changelog — unverändert übernommen)
