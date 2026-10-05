# PlazCode Notion Web Agents (Chrome-Erweiterung)

Verbindet freigegebene KI-Chat-Tabs (ChatGPT, Claude, Gemini, DeepSeek, Qwen, Kimi, GLM, Arena, optional Notion AI) mit PlazCode Notion. Die Erweiterung navigiert nicht selbst und klickt keine beliebigen Befehle.

## Installation

1. Chrome 116 oder neuer.
2. `chrome://extensions` öffnen, **Entwicklermodus** aktivieren, **Entpackte Erweiterung laden** und den Ordner `%LOCALAPPDATA%\PlazCodeNotion\extension` (oder `notion-desktop/extension` aus diesem Repo) wählen.
3. PlazCode Notion starten → Seite **Notion AI** → Karte „Web-Agenten & Sicherheit“ → **Erweiterungs-Token** kopieren.
4. In den Optionen der Erweiterung Adresse `ws://127.0.0.1:8787/extension/ws` und den Token eintragen, speichern.
5. Im Erweiterungs-Menü sollte „connected“ stehen; in PlazCode zeigt die Karte „online“.

## Agent pro Tab freigeben

KI-Chat öffnen, anmelden, Erweiterungs-Menü → **Agent freigeben**. Nur so freigegebene Tabs sind in Notion über `web_agents` sichtbar.

**Direkt starten** lässt den Chat als Roblox-Agent arbeiten (braucht PlazCode Notion + verbundenes Roblox Studio, kein Notion). Lesen ist frei, jede Änderung erscheint in einem Freigabe-Panel auf der Seite. Während eine Direkt-Sitzung läuft, sind Notion-Schreibzugriffe auf Studio gesperrt; Stop gibt die Sperre frei.

## Nutzung aus Notion

- `web_agents` mit `{"action":"list"}` – freigegebene Tabs und `agent_id`
- `web_agents` mit `{"action":"run","agent_id":"…","task":"…"}` – Aufgabe senden und Antwort abwarten
- `web_agents` `inspect` / `read`, `web_chat` (`inspect`, `send`, `wait`, `read`) für den aktiven Tab, `web_sites` listet unterstützte Seiten

DeepSeek ist Beta, die anderen Seiten experimentell (sie können ihr DOM ändern). Herkunft: portiert aus Secretscript, siehe THIRD_PARTY_NOTICES.md.