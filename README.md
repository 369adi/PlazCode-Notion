# PlazCode Notion Desktop / AdiCode (eigenständige Weiterentwicklung, GPL-3.0-Basis)

Steht in keiner Verbindung zu den PlazCode-Autoren oder zu Notion Labs. Alle Original-Lizenz-,
Copyright- und Branding-Hinweise (LICENSE, BRANDING-NOTICE.txt, THIRD-PARTY-NOTICES) bleiben erhalten.

## Was ist anders?
- Name „AdiCode" (technisch PlazCode Notion), dunkles Design.
- Neue Seite **Notion AI**: eingebauter MCP-Server (127.0.0.1:8787/mcp, Token via Header oder Pfad) + ngrok-Tunnel mit fester Domain — Notion AI direkt, ohne Browser-Erweiterung.
- Notion sieht `roblox_*`, `pc_*`, `browser_*` und alle weiteren aktivierten MCP-Server.
- Automatische Updates über den eigenen Launcher (Original-Update deaktiviert).

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

## Neu in 1.0.34–1.0.23
(fixiert in 1.0.29: lokale schnelle Releases; Kompaktversion der alten Einzelnachrichten)

## Neu in 1.0.22 und älter
(siehe HANDOUT.md Changelog — unverändert übernommen)
