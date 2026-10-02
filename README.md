## PlazCode 1.18.92: Tasks close button and live Co-work follow-ups

Close Tasks reliably and send follow-ups during agent work without waiting for the entire task to finish.

- Co-work delivers new requests alongside the next tool result or continuation, preserving the current task and the order of follow-ups. Streaming responses and running commands are not interrupted.
- Follow-up status confirms receipt and explains delivery at the next safe message boundary. Unconfirmed sends remain paused for review rather than being blindly retried.
- Tasks has a dedicated sticky header and accessible close button, separated from content, with explicit button type and click handling.

# PlazCode — Roblox Studio + AgentScript AI agent

Turn any major AI chat (**DeepSeek, ChatGPT, Google Gemini, Kimi, GLM, Qwen, Arena, Crax GPT, or Ollama running locally**) into an autonomous development agent. Three switchable engines:

| Engine | Toggle | Target | Port |
|---|---|---|---|
| **Roblox** (RS) | — | Roblox Studio via its built-in MCP server | ws://127.0.0.1:17613 |
| **AgentScript** (AS) | — | A local project folder — files + terminal | ws://127.0.0.1:17615 |
| **Animation** (AN) | — | Roblox Studio scoped to the motion workflow | ws://127.0.0.1:17613 |

Describe what you want in plain English and the AI builds instances, writes Luau/code files, sculpts terrain, tunes lighting, generates UI, runs builds and tests, and audits your project — inside Studio or directly on disk.

No API keys, no monthly fees. Chromium browsers (Chrome, Brave, Edge, Thorium). The extension and desktop app use the same dark navy + warm gold PlazCode theme.

---

## Engines

- The bar above every supported chat composer carries a segmented **RS / AS / AN** toggle. Switching engines wipes tool caches so commands cannot cross engines.
- **RS** drives Roblox Studio through StudioMCP (stdio JSON-RPC spawned by `plazcode-agent`).
- **AS** gives the AI full control of ONE local folder ("the workspace") through native Rust tools — sandboxed paths, exact-match diff editing, glob/content search, and terminal execution with hard timeouts.
- **AN** rides the same Roblox bridge as RS but steers the system prompt into the animation_* workflow.

Large `execute_luau` scripts are auto-chunked around 24 KB so Studio's parser never hits the ~64 KB wall. Each chunk is still one NDJSON/JSON-RPC line.

---

## Setup

1. Open `chrome://extensions` → Developer mode → **Load unpacked** → this folder (`manifest.json`).
2. Double-click **`PlazCode.exe`**. `PlazCode.exe --headless` runs without the desktop window. It starts:
   - HTTP API on `http://127.0.0.1:3000`
   - WS bridges on `17613` (RS/AN) and `17615` (AS)
   - Workspace folder for AgentScript (`PLAZCODE_WORKSPACE_ROOT` / `--workspace` / `%USERPROFILE%\PlazCodeWorkspace`)
3. **RS/AN:** Roblox Studio → Assistant AI → ⋯ → Manage MCP Servers → Enable Studio as MCP Server.
4. Open a supported chat and click **Start agent**.

---

## Agent

- Native crate: `agent/` (`plazcode-agent` 1.18.86). Desktop control center, MCP helper spawn, outbound WS channel so ping/status keep flowing during a 20 s `execute_luau`.
- Service worker skips stale-socket reconnect and MCP heal while a `call_tool` is in flight (the 25 s stale window used to kill long tools).
- 30 Studio skills, a 24-command animation suite, AgentScript file/terminal tools.
- Personas (Builder / Scripter / Animator / Fixer), Extra Thinking, Forge GUI, Image → Model, auto-fix playtest errors.

---

## Testing

```bash
node test-skills.js
node test-parser.js
node test-chatgpt.js
node test-animlib.js
node test-v111.js
node test-v112.js
node --check core/main.js && node --check core/config.js && node --check background.js
cd agent && cargo test
```

`test-bridges.js` is a live smoke test: start `PlazCode.exe` first. It checks HTTP `:3000` and WS `17613` / `17615` (no Unreal port).

## Privacy

Everything runs locally. The extension talks only to `127.0.0.1`. No telemetry. AgentScript stays inside the workspace root unless you flip FULL ACCESS.

Optional MCP runtimes are downloaded on first use if no installed runtime is found. They are cached under %LOCALAPPDATA%\PlazCode\runtimes and reused across app updates. The release ZIP contains no Node.js/npm/uv bundle. First use requires internet access.

## Release history

### PlazCode 1.18.91: Desktop work details and responsive status

Expanded desktop work details now show AI replies and commands, with cleaner Home spacing, accurate icons and faster automatic update detection.

- Templates now use Choose a Roblox file and Add template, with short numbered directions. Advanced script search is optional and collapsed.
- Desktop template imports support Roblox files up to 128 MB using direct binary uploads, avoiding base64 copies. Browser imports remain limited to 32 MB; script-index bounds still apply.
- Memory transfer separates Save to a file from Load into another chat, explains saved notes versus messages, and uses clear create/download/load buttons.

- **Added:** Desktop Working/Worked details show AI replies, command text and results from the selected browser chat, loaded only while the panel is open.
- **Improved:** Automatic release checks run every five seconds, with desktop/browser status propagation every two seconds. The default GitHub feed uses a five-second cache key; network/cache delays can still apply. Checks do not overlap.
- **Improved:** DeepSeek and Claude shortcuts use vector logos from the supplied references. Tools uses the supplied icon, and Settings has a symmetric centered gear.
- **Fixed:** Task checkpoints & workflows has a clear gap below the Home cards; the Working panel also has consistent spacing.
- **Fixed:** Expanded work details preserve unchanged content and scrolling across refreshes, collapse once on completion and remain reopenable. Retention is bounded; oversized/older omitted details receive an explicit notice.
- **Fixed:** The updater swirl rotates normally and uses a slower continuous rotation when reduced motion is enabled, instead of freezing.

### PlazCode 1.18.90: Grouped replies and extension status

Task replies and commands collapse together, with extension status under the bar name and a corrected ChatGPT shortcut logo.

- **Added:** The browser bar shows its installed extension version below PlazCode, with green (Up to date), red (Outdated - New update available), or a neutral unavailable/unchecked status.
- **Improved:** The desktop ChatGPT shortcut uses a vector logo traced from the supplied reference instead of a circle character. Its existing browser-opening action is preserved.
- **Fixed:** Working/Worked now includes task prose, final responses, command cards and injected results. Completed groups collapse together and reopen on click. Startup, list_commands/list_tools and PlazCode is ready. remain outside; user requests remain visible.
- **Fixed:** Activity identity survives streaming text changes and transcript remounts, and remains scoped to the chat.

- **Added:** ChatGPT: Highly intelligent in all aspects. DeepSeek: Recommended for free users. Claude: Expert at scripting and game development. Each description appears below its name and above Open in browser.

- The updater explicitly preserves native settings, MCP configuration/enabled servers, memory, pairing, templates, update source and WebView user data, even if a package contains saved-data files.
- The Tools page includes both engine skill libraries, imported Studio helpers, asset_bridge_import and developer-product tools with descriptions.
- Working/Worked uses a centered vector chevron that points right when collapsed and rotates down when expanded.

### PlazCode 1.18.89: Home engine switch

Switch between RobloxScript and AgentScript directly above Start Agent on the desktop Home page.

#### Added

- A themed, keyboard-accessible RobloxScript / AgentScript button group under the active AI status and above Start Agent. The selected engine is highlighted.

#### Improved

- The Home buttons use the existing saved engine preference and browser synchronization, and stay synchronized with the Settings engine selector. Repeated clicks are disabled while saving.

#### Fixed

- Start Agent is temporarily disabled while the engine change is saving or the selected browser chat still reports a different engine, preventing a start on the previous engine.

### PlazCode 1.18.88: Automatic template references

Import Roblox place/model files and let PlazCode find relevant reference systems automatically before a task’s first tool call.

#### Added

- Templates tab in the desktop app and a Template library inside browser Tasks. Import .rbxl, .rbxlx, .rbxm and .rbxmx files up to 32 MB, add a systems description, browse/read scripts and remove saved references.
- Automatic matching uses cached script paths, identifiers and template descriptions. Up to three bounded excerpts are supplied before the first agent tool executes; no explicit “reference this template” request is required.
- Read-only plazcode_templates tools support explicit matching, script lists and paged source reading. Duplicate script paths have separate IDs.
- Open in Studio explicitly opens the retained original through the Windows Roblox file association; select the desired Studio session in PlazCode afterwards.

#### Improved

- Template code is reference data, not executable instructions. Existing project conventions take priority. Matching is a relevance hint rather than guaranteed semantic understanding.
- Template files/indexes persist in your private app data, survive updates and are excluded from distribution ZIPs. Source indexes and excerpts have explicit size limits to avoid flooding long chats.

#### Fixed

- Malformed, unsupported and source-free files report clear import errors. Missing scripts in copied/decompiled games are not fabricated or recovered. Binary and XML imports preserve the original file bytes.

### PlazCode 1.18.87: Task checkpoints, project memory & reliable tools

Review and recover agent tasks, keep project facts separate, and reduce repeated tool errors and long-chat overhead.

#### Added

- Tasks panel in the desktop Home page and browser bar: bounded task history, file before/after previews, conflict-checked Revert and resume in the original chat.
- AgentScript file checkpoints persist across app restarts. Studio Revert is offered only when exact undo records are confirmed in the current bridge/Studio session. Unsupported operations disable whole-task automatic restore.
- Protected paths with explicit task approval; project memory scopes shared by named projects; reusable built-in and custom workflows.

#### Improved

- Long-chat activity retains at most 60 groups, builds collapsed details only when expanded, ignores its own UI mutations and caches unchanged export text. Retention overflow produces an explicit export error.
- Automatic release checks run on launch and every minute; returning to an AI tab or desktop window triggers a throttled check. Update labels refresh without pressing Check now. Network/cache delays can still apply.

#### Fixed

- Tool preflight checks JSON argument shapes against available schemas, detects incomplete Luau strings/delimiters and preserves datamodel defaults. Three consecutive malformed replies or execution failures pause the agent.
- Studio mutation commands are not blindly replayed after a dropped helper connection; uncertain results require inspecting Studio before retrying.
- Unchanged memory and release snapshots no longer generate repeated storage broadcasts. Workflow delivery acknowledgements are deduplicated.

### PlazCode 1.18.86: Reliable bridge port startup

Agent launch reserves every bridge endpoint before reporting ready and serializes overlapping launches.

#### Improvements

- Windows startup is serialized across PlazCode launches; an existing instance of the same version is focused, and older versions are replaced.
- All three required bridge listeners are bound before background services start. Port release is checked with bounded retries instead of fixed delays.

#### Bug fixes

- Removed broad process-name cleanup that could kill another launching PlazCode instance. Recovery only targets a PlazCode process holding a required port.
- Port ownership uses exact endpoint matching, avoiding matches such as port 30000 when checking 3000.
- Startup errors now identify the actual blocked port and, when available, its owning process instead of the generic bridge-bind dialog. Legacy listener failures are logged.

### PlazCode 1.18.85: Wait for app shutdown before updating

The updater waits for PlazCode to exit and release its executable before installing.

#### Improvements

- The updater checks both app executable paths, waits for stopped processes to exit, and briefly retries file-lock checks. Only processes from the selected installation are stopped; Roblox Studio is left running.

#### Bug fixes

- Update installation no longer starts copying immediately after Stop-Process. If the executable remains locked, it aborts before copying any update files and retains recovery files.

### PlazCode 1.18.84: Clickable palettes & friendly readiness

Select desktop themes directly from their color swatches, and keep Notion startup acknowledgements readable.

#### New additions

- Desktop startup instructions accept up to 25,000 characters, with an input counter and explicit save validation.

#### Improvements

- All six desktop color swatches are clickable, keyboard accessible, show the selected theme and use the same saved appearance settings as the dropdown. Glow and gradient preferences are preserved.
- Chat export retains messages captured while scrolling during the current page session. Idle Notion exports briefly load older history and restore the previous scroll position. Exports report their captured message count and remaining scope limitations.
- Download export is styled as a theme-colored button in desktop and browser memory settings.

#### Bug fixes

- All startup instruction variants explicitly request “PlazCode is ready.” instead of a generic readiness sentence.
- If Notion still returns the legacy PLAZCODE_READY token from earlier chat context, its visible acknowledgement is normalized to “PlazCode is ready.” without resending startup or modifying Notion-owned DOM.
- DeepSeek rich-text composers remain discoverable while input-locked, and agent writes temporarily enable editing then restore the lock. Explicit Send labels take precedence over stop-icon guesses.

### PlazCode 1.18.83: Clearer update indicators

Outdated version badges now explicitly say that a new update is available.

#### Improvements

- Desktop, browser bar and popup version badges show “(Outdated - New update available)” when a newer published version is detected. Red highlighting and existing current/unavailable states are preserved.

### PlazCode 1.18.82: Desktop overhaul & smoother Notion sessions

A redesigned desktop workspace, complete color themes and a polished animated updater, with focused Notion and Co-Work improvements.

#### New additions

- A distinct desktop layout: compact navigation, a dedicated agent-session card, AI launch cards, a slim connection strip, and redesigned Tools, MCP, Settings, Terminal and Updates surfaces.
- Animated SVG update swirl with layered rotating arcs, a pulsing center and a redesigned progress dialog. Reduced-motion preferences are respected.
- Expandable Working/Worked summaries in the chat show elapsed time and tool activity; completed details collapse automatically and reopen on click.

#### Improvements

- Themes now color panels, navigation, borders, inputs, buttons, connection surfaces, terminal and update dialogs as well as accents.
- Startup asks the AI to say “PlazCode is ready.”; older readiness acknowledgements remain accepted.
- Co-Work follow-up text uses the website composer’s text color, font, size and alignment. Only the placeholder is grey. Existing queue behavior is preserved.
- Copilot and Meta AI removed from supported-site choices, routing and extension injection.

#### Bug fixes

- Notion can settle a stable readiness acknowledgement without waiting through its nine-second generation tail; live Stop or workflow progress still blocks completion.
- Notion pastes an escaped literal HTML representation alongside plain text so rich-text Markdown conversion does not consume command/result markers. Complete-draft retention and send-confirmation checks remain in place.
- Notion reserves space inside the composer for the PlazCode bar, keeping response text and the native editor clear without modifying React-owned attributes.

### PlazCode 1.18.81: Themes, release notes & desktop refresh

See what each PlazCode update adds, improves and fixes before installing it.

#### New additions

- Release descriptions on the desktop Updates page, including expandable notes for earlier releases.
- The GitHub README now lists release titles, summaries, additions, improvements and fixes.
- An updated maintenance prompt documents architecture, current features, debugging, tests, builds and release publishing.
- Appearance settings with Orange, Amethyst, Polar Cyan, Rose, Emerald and Graphite palettes, glow strength, gradient controls and a default-theme reset. Selections apply immediately and are remembered on this device.
- Version indicators in the desktop, browser bar/menu and extension popup show green “Up to date”, red “Outdated”, or a neutral unavailable/not-checked status. Release checks run at startup and every five minutes.

#### Improvements

- Published update metadata carries the same release notes used in the README.
- Installed release notes remain available before checking online.
- Home places agent controls beside the AI shortcuts on wide windows, reducing unused space; smaller windows stack them neatly.
- The installed-to-latest version arrow is larger, vertically aligned with the version values and spaced more closely.
- Visual-only desktop polish: smoother navigation, card depth, hover and press feedback, keyboard focus, toggles and modal backgrounds. Existing actions and workflows stay the same.
- Quick Actions no longer stretches to match a tall Recent Activity panel; activity stays scrollable.
- A denser desktop layout reduces empty space in status cards and quick actions, with refreshed typography, surfaces and subtle orange highlights.

### PlazCode 1.18.80: Updates, memory transfer & AI shortcuts

Update PlazCode from the desktop app and carry chat context between supported AI sites.

#### New additions

- Desktop Updates tab with installed/latest versions, Check now, Update now, download progress and relaunch.
- ChatGPT, DeepSeek and Claude shortcuts below Home’s connection displays.
- Current-chat memory export and a memory-only export option.

#### Improvements

- Transparent Enhance and Co-Work controls; labels now read Co-Work: Off and Co-Work: On.
- Styled desktop memory import file picker.
- Saved AI memories retain conversation origins after edits, including facts learned in more than one chat.

#### Bug fixes

- Memory-only imports now include saved facts in the destination AI’s context.
- Oversized combined history and memory are rejected before saving or sending.
- View on GitHub opens in the default browser instead of navigating the embedded app.

### PlazCode 1.18.78: Automatic update setup

A stable repository feed lets the update BAT fetch published packages.

#### New additions

- Configured HTTPS update feed hosted in stoveez/PlazCodeneww.

#### Improvements

- Version and SHA256 checks, staged installation, backups and local settings preservation.

#### Bug fixes

- Blank feed settings from 1.18.77 now fall back to the configured repository.

### PlazCode 1.18.77: Shared memory & activity

PlazCode can remember lasting preferences across chats and transfer loaded conversation context.

#### New additions

- Shared personal memory with automatic saving controls and an editable summary.
- Chat history and memory import/export.
- Working/Worked indicators with expandable activity.
- Studio status below the desktop bridge indicator.
- Update-PlazCode.bat with a local ZIP fallback.

#### Improvements

- Saved memory is available across supported provider chats.
- Completed tool activity collapses while final replies remain visible.

