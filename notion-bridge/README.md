# PlazCode Notion Bridge

> **Inoffizieller Fork** von [PlazCode](https://github.com/stoveez/PlazCode). Nicht vom PlazCode-Projekt veröffentlicht oder unterstützt. Lizenz: GNU GPL-3.0 (siehe `LICENSE` im Repository-Stamm).

Eine einzelne Windows-`.exe`, **ohne Browser-Extension**. Sie bündelt beliebig viele lokale MCP-Server hinter **einem** MCP-Endpunkt mit fester Adresse. Notion (Desktop-App, Browser, Mobil) wird **einmal** verbunden – danach werden Server nur noch in der App verwaltet.

```text
Notion (Custom MCP „PlazCode Notion“, einmal verbunden)
        │  HTTPS + Bearer-Token
        ▼
https://<deine-domain>.ngrok-free.dev/mcp   (feste ngrok-Domain)
        ▼
PlazCode-Notion-Bridge.exe  (127.0.0.1:8787/mcp)
        ├── pc_*      PC full access (windows-mcp 0.8.7)
        ├── roblox_*  Roblox Studio MCP
        ├── browser_* Brave Browser (Playwright MCP, eigenes Profil)
        └── <id>_*    weitere stdio-/HTTP-MCP-Server aus config.json
```

## Voraussetzungen

- Windows 10/11 x64
- [ngrok](https://ngrok.com/download) + einmalig `ngrok config add-authtoken <TOKEN>`
- [uv](https://docs.astral.sh/uv/) (für PC full access)
- Git + Node.js/npm (für Roblox Studio MCP und Brave-Browser-MCP)
- Brave (alternativ Edge/Chrome; Pfad über `browserPath` änderbar)

## Einrichtung

1. `PlazCode-Notion-Bridge.exe` starten. Beim ersten Start wird `%LOCALAPPDATA%\PlazCodeNotion\config.json` mit eigenem Bridge-Token angelegt.
2. Optional: **Roblox MCP installieren** → **Roblox-Plugin-Token kopieren** → Roblox Studio → Plugins → MCP → einfügen → Connect (einmalig).
3. **Alles starten**.
4. **Notion-Einrichtung (Anleitung)**: In Notion unter *Settings → Connections → Add custom MCP*
   - URL: `https://<domain>/mcp`
   - Auth-Header: `Authorization: Bearer <Bridge-Token>` (Button „Bearer-Token kopieren“)
   - Ohne Header-Feld: „URL mit Token kopieren“ (`/k/<token>/mcp`) verwenden.

Adresse und Token bleiben dauerhaft gleich. Nach PC-Neustart einfach die App starten (oder „Mit Windows starten“ aktivieren).

## Weitere Server

„Server-Konfiguration bearbeiten“ öffnet `config.json`. Beispiel:

```json
{ "id": "files", "name": "Dateisystem", "kind": "stdio", "enabled": true,
  "command": "npx", "args": ["-y", "@modelcontextprotocol/server-filesystem", "C:\\Projekte"] }
{ "id": "remote", "name": "Remote", "kind": "http", "enabled": true,
  "url": "http://127.0.0.1:9000/mcp", "bearerToken": "..." }
```

`id`: nur `a-z0-9`, max. 16 Zeichen. Danach „Konfiguration neu laden“ und „Alles neu starten“.

## Sicherheit

- Die Bridge lauscht nur auf `127.0.0.1`; öffentlich erreichbar ist sie ausschließlich über ngrok **mit** Bridge-Token.
- PC full access kann den Rechner vollständig steuern. Token nicht weitergeben; bei Verdacht „Bridge-Token erneuern“.
- Alle Kindprozesse laufen in einem Windows-Job-Objekt und werden beim Schließen der App beendet.
- Roblox-Publishing ist standardmäßig deaktiviert (`robloxAllowPublish`).

## Bauen

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -H windowsgui" -o PlazCode-Notion-Bridge.exe .
```

Unter Linux/macOS startet `go run . -dir ./data -no-ngrok` einen Headless-Testmodus.
