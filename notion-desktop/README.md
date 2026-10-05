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

## Start
1. Original-PlazCode und die alte `PlazCode-Notion-Bridge.exe` beenden (Port 3000 / 8787 und ngrok-Domain werden sonst blockiert).
2. ZIP entpacken und `PlazCode-Notion\PlazCode.exe` starten.
3. Seite **Notion AI** öffnen → Server-URL und Token in Notion als benutzerdefinierten MCP-Server eintragen.
   Token und Domain der alten Bridge werden beim ersten Start automatisch übernommen.

Einstellungen: `%LOCALAPPDATA%\PlazCodeNotion\plazcode-notion.json`

## Build
GitHub Actions (`.github/workflows/build-notion-desktop.yml`) entpackt `PlazCode-source-1.19.36.zip`,
wendet `notion-desktop/plazcode-notion.patch` an, baut `agent` mit Cargo (MSVC) und ersetzt `PlazCode.exe`
im Original-Paket `PlazCode-1.19.36.zip`.
