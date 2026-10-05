# PlazCode Notion (inoffizieller Fork)

**PlazCode Notion** ist ein inoffizieller Fork der PlazCode-Desktop-App von [stoveez/PlazCode](https://github.com/stoveez/PlazCode) (GPL-3.0).
Er steht in keiner Verbindung zu den PlazCode-Autoren oder zu Notion Labs. Lizenz- und Branding-Hinweise: `LICENSE`, `BRANDING-NOTICE.txt`.

## Download
Die fertige App gibt es unter **[Releases](https://github.com/369adi/PlazCode-Notion/releases)** → `PlazCode-Notion-<version>-windows.zip`.
Entpacken, Original-PlazCode beenden, `PlazCode-Notion\PlazCode.exe` starten.

## Was ist anders als im Original?
- Notion-Design und Name „PlazCode Notion“
- Seite **Notion AI**: eingebauter MCP-Server + ngrok-Tunnel → Notion AI steuert Roblox Studio, PC und Browser ohne Browser-Erweiterung
- Katalog-Einträge „PC (Windows-MCP)“ und „Browser (Playwright)“ unter MCP Servers
- Automatische Updates deaktiviert (sonst würde das Original den Fork überschreiben)
- Neu in 1.0.14: Seite Co-Work (mehrere Notion-Tabs als getrennte Container mit eigenem Gmail-Login, Alle starten). Fuer neue Chats: siehe `HANDOUT.md`
- Neu in 1.0.10: Seite Updates (Versionen + Notizen, Jetzt aktualisieren), ngrok-Ersteinrichtung in der App (installieren, Authtoken, Domain/URL)
- Neu in 1.0.9: aufgeräumte Oberfläche (ohne Model/UI Builder, Toolkit, Templates, Updates), getrennte Notion-Fenster mit eigenem Login pro Profil
- Neu in 1.0.8: Web-Agenten (Chrome-Erweiterung), Direkt-Modus mit Freigabe, `roblox_studio`/`roblox_workflow`, Projekt-Gedächtnis, Snapshots, sicheres Datei-Schreiben, Shell-Sicherheitsmodus

## Repository-Inhalt
| Datei | Zweck |
|---|---|
| `notion-desktop/plazcode-notion.patch` | Alle Änderungen gegenüber PlazCode 1.19.36 (Quellcode) |
| `PlazCode-source-1.19.36.zip` | Original-Quellcode (GPL) – Basis für den Build |
| `PlazCode-1.19.36.zip` | Original-Paket – darin wird nur `PlazCode.exe` ersetzt |
| `.github/workflows/build-notion-desktop.yml` | Baut die exe auf GitHub (Windows, Rust/MSVC) und veröffentlicht das Release |
