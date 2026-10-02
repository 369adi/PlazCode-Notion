## PlazCode 1.18.104: Animated installer and Notion protocol files

- A rounded dark installer window with PlazCode orange gradients, an animated shimmer and swirl, stage text and real file-copy progress.
- The installer UI runs on its own STA message loop so download, verification, backup, extraction and process waits cannot freeze its animation. File replacement progress advances from actual completed file counts.
- Notion long protocol messages use one named .txt file containing the complete original text, plus a short native composer instruction. Upload readiness and one confirmed send are required; text is not truncated.
- Notion no longer retries a 33k+ text paste that the site converts into a Pasted text attachment. Existing draft attachments are preserved and receive a specific fresh-chat/remove-draft explanation. Stopped owned uploads are cleaned up.
- Short protocol messages retain the native composer path. Saved settings, MCP configuration, memory, templates, installation process scope, checksums and rollback behavior are preserved.

The refreshed helper window is used after this version is installed; the already-installed old helper can still appear during the update to 1.18.104. Notion remains experimental until live startup is confirmed. Open a fresh Notion AI chat if old Pasted text draft cards remain. Reload the extension AND refresh existing AI tabs after updating.

## PlazCode 1.18.103: Notion controlled paste and Co-work

- Notion sends complete startup and tool-result text through one literal controlled paste transaction, verifies retention through two renders, and commits at most once. No prompt truncation or repeated tiny inserts.
- Co-work uses Follow up in the native composer while enabled, including idle time and after Stop. Internal startup/tool injection retains its temporary input mask. Turning Co-work off restores the original placeholder.
- Notion no longer uses bulk execCommand or direct React-owned DOM replacement as the large-send fallback. Short native edits rely on Chromium input events instead of duplicate synthetic model events.
- Notion input cover targets only the editable, preventing a growing ancestor from becoming the click-blocking mask. Start on a populated AI landing first opens a fresh chat.

Notion remains experimental: simulated model/DOM tests cannot confirm live Notion behavior. After updating, reload the extension AND refresh all open AI tabs.

## PlazCode 1.18.102: Updates, Stop, Notion and Co-work

- The default GitHub updater checks the current repository revision every two seconds and fetches metadata by immutable commit. Unchanged metadata is reused; quiet automatic checks and the eight-second request budget remain. Custom feeds retain their configured path. Network delays and browser suspension still apply; this is not guaranteed instant push.
- Notion large startup drafts use one complete native beforeinput/insertText/input write instead of chunked clipboard pastes. Editable/readonly attributes are temporarily restored during the write, focus is acquired after unlocking, and the complete draft must survive two renders before Send.
- Co-work tracks placeholder nodes introduced by composer remounts, and ChatGPT Work literal placeholder labels display Follow up without changing user text or replacing the native input.
- Stop releases the composer immediately and prevents late work from restoring the Agent is working cover until a new user action resumes work. Tool cancellation and Safe/Immediate semantics remain unchanged.
- PlazCode Working/Worked transcript summaries and bar disclosure are disabled on ChatGPT. Activity records remain available to desktop; other providers retain their collapsible groups.
- Default update checks avoid the cached mutable branch feed when current revision discovery succeeds. Older cached metadata cannot replace a newer known release.

Notion remains experimental. Native draft checks and placeholder/group tests use simulated DOMs; live Notion startup still needs confirmation.

After updating, reload the extension AND refresh every existing AI tab.

## PlazCode 1.18.101: Notion startup and engine switching

- All Notion literal pastes, including startup chunks, use escaped code HTML to avoid rich-text Markdown conversion. Complete normalized text must be retained before Send; sampled near-match acceptance is removed.
- The browser bar shows Switching while AgentScript/RobloxScript reconnects, disables Start, ignores old-engine status and restores actual health on completion or after an eight-second bound.
- A rejected retained large paste no longer falls into repeated tiny inserts and full rewrites. One initial empty hydration rejection may retry; partial failures stop without sending.
- Intentional engine reconnection does not flash a bridge-down banner. A genuine failure remains visible after the transition timeout.

Notion remains experimental. Tests use simulated rich-text imports and provider DOMs; live Notion startup still needs confirmation.

After updating, reload the extension AND refresh every existing AI tab.

## PlazCode 1.18.100: Browser bar startup and placement fix

- Fixed the 1.18.99 shared UI initialization order. When a composer already existed, placeholder state was read before initialization and startup aborted, leaving the bar missing or unpositioned. This affects multiple providers and Chromium browsers.
- ChatGPT anchors to the actual bounded composer card, rather than a full-page form or its inner text scroller. Existing DeepSeek/Claude/Notion placement rules are preserved.
- Desktop browser status now reports bar visible, bar hidden or refresh AI tab separately from bridge connectivity. A worker connection alone does not establish that the page UI loaded.
- Full content-script startup tests cover preloaded ChatGPT, DeepSeek, Claude and Notion composers. Native Co-work and Stop behavior are preserved. Live tester layouts still require confirmation.

After updating, reload the extension in chrome://extensions or edge://extensions AND refresh existing AI tabs.

## PlazCode 1.18.99: Native Co-work, continuation and reliability

- Co-work uses the existing ChatGPT, Claude, DeepSeek and Notion composer. While work is active, its placeholder becomes **Follow up**. Requests are accepted immediately and delivered at the next safe AI message boundary. Native Stop keeps its role. Unsent drafts are protected during tool feedback.
- Notion long startup/results use bounded, verified paste chunks rather than repeatedly pasting a large prompt as a text-file attachment. Startup still requires the list/ready acknowledgement; failed delivery stays stopped.
- New chats can offer **Yes, continue.** or **No, I wanna start a new project.** Context is loaded only after successful startup. Captured history is bounded to 2,000 messages / 800,000 text bytes per chat, 20 chats / 8 MB total; continuation is a shorter excerpt and can omit unloaded messages. Inspect current state before resuming uncertain actions.
- Shared command/time budgets pause before the next command; zero disables a limit. Export continuation or explicitly resume a paused task. Automatic Studio visual checks are bounded, optional and do not prove runtime correctness.
- Restore an individual captured local file or Studio script in Tasks. Restoration checks current contents first and refuses conflicting later edits. Studio script snapshots cover direct multi_edit/script_set_source edits in Edit mode and the same bridge session; arbitrary Luau mutations are not selectively restored.
- The extension toolbar popup now follows the desktop/bar palette, glow and gradient settings.
- Each update now creates a GitHub Release with both ZIPs and change notes.

Update in the desktop Updates tab or run Update-PlazCode.bat, then reload the extension and refresh AI tabs. Automated checks and the Windows build are verified; live Windows/Studio and live provider behavior still need target-machine confirmation.

## PlazCode 1.18.98: Cancellation settling guard

Includes all 1.18.97 features below. Keeps Immediate Stop latched while an interrupted tool is settling and blocks further tool dispatch during that interval. Safe Stop waits for the current tool step. Reload the extension and refresh AI tabs.

PlazCode 1.18.97
Immediate Stop is the default; Safe Stop waits for the current tool step. External operations may continue and partial changes can remain. Do not promise rollback. Reload the extension AND refresh AI tabs after updating. Keep tray, media and work folding separate from provider parsing/startup contracts.

# PLAZCODE MAINTENANCE, DEBUGGING & RELEASE PROMPT

Release policy: every version must publish fresh and compatibility ZIPs, atomically update latest.json, and include matching change notes in README and the desktop update feed. Verify the GitHub release workflow succeeds and the published version release contains both ZIP assets with matching SHA-256 digests. Never call an update shipped solely because source changed. Keep private chat-history.json, memory, settings and credentials out of packages and preserve user data on update. Native Co-work must use provider-owned composers; do not add a replacement text box. Keep native Stop separate from Send, guard draft writes against changed text, and test long Notion prompts for attachment conversion, retained content, cancellation and rejection. No live startup claim without live evidence.



You maintain PlazCode: its browser extension, Windows desktop application, and local Rust bridge/agent. Perform requested work, verify it, and ship a complete update. Preserve everything that already works. Make the smallest safe change that fully satisfies the request; a requested visual redesign can change presentation substantially while preserving all behavior.

## Source of truth

At the time this guide was written, the prepared release is **1.18.88**. Never treat that number as permanently current. Before starting, inspect the actual files, `manifest.json`, `agent/Cargo.toml`, release metadata, and the user's newest supplied build. Use the newest verified project state rather than an older remembered ZIP.

Windows locations:
- Project: `C:\Users\plazm\OneDrive\Desktop\PlazCode`
- AgentScript workspace: `C:\Users\plazm\PlazCodeWorkspace`

Update repository: `https://github.com/stoveez/PlazCodeneww`
Stable release feed: `https://raw.githubusercontent.com/stoveez/PlazCodeneww/main/latest.json`

The repository currently hosts complete release ZIPs and update metadata. Inspect it before assuming it contains a checked-out source tree. If working elsewhere, use the actual available workspace paths and tools. Do not pretend to have accessed the user's Windows machine.

## Architecture and boundaries

Release layout: outer PlazCode/ holds the desktop app, launchers and updater; PlazCode-Extension/ inside it holds only Chrome assets. The working source may remain flat for existing include_str! paths. Build ZIPs with package_release.py, which stages this split and validates every manifest asset. Publish a compatibility update ZIP for older flat-layout installers, without renaming installed folders or deleting user data. Never rename the whole application folder PlazCode-Extension.

- Extension: `manifest.json`, `background.js`, `core/`, `providers/`, `popup.*`, `overlay.css`.
- `background.js`: service worker, routing, provider URLs, authenticated local bridge requests, desktop-to-tab actions, and settings synchronization.
- `core/main.js`: PlazCode bar, startup, agent/tool loop, shared controls, and coordination through the provider interface.
- `providers/*.js`: each site's composer, submission, transcript, generation state, conversation identity, and site-specific DOM behavior. Keep site-specific selectors and DOM work here, not in core.
- `core/cowork.js`, `core/memory.js`, `core/activity.js`: shared Co-work, memory/transfer UI and policy, and Working/Worked state. Inspect the implementations before extending them.
- Rust source: `agent/`. The crate is named `plazcode-agent`, but the current binary target and shipped desktop executable are **`PlazCode.exe`**, not `plazcode-agent.exe`.
- Embedded desktop UI: `agent/src/desktop.html`; native window and IPC: `agent/src/gui.rs`; authenticated APIs: `agent/src/main.rs`; memory: `agent/src/memory.rs`; updates: `agent/src/updater.rs`.
- The desktop embeds some files with `include_str!`, including the memory UI and release notes. Changes to embedded HTML/JS/JSON require a rebuilt executable even when no Rust logic changed.
- `Start-PlazCode-Agent.cmd` launches the application. `Update-PlazCode.bat` invokes `Update-PlazCode.ps1`. `update-source.json`, `latest.json`, `release-notes.json`, `UPDATE.txt`, and `SHA256SUMS.txt` participate in releases.

The extension and native agent are separate components and can update independently. Do not rebuild the agent for unrelated extension-only changes. Do not ship two competing desktop executables.

## Current capabilities to preserve

- The desktop starts its local bridge and presents connection status. Browser tabs recognize the bridge without a separate manual launcher when startup succeeds.
- RobloxScript and AgentScript modes, existing provider adapters, tool routing, MCP server controls, workspace operations, custom startup instructions, and the existing reasoning/debug/access options.
- Desktop Home Start/Stop mirrors the selected browser chat's readiness. Play remains disabled until requirements are met. Preserve engine selection and quick settings synchronization.
- Prompt Enhance on both surfaces asks the selected AI to rewrite a request before sending the enhanced request. Preserve the user's supplied enhancer policy; do not execute tool commands while rewriting it.
- **Co-Work is a persistent mode**, off by default, beside Start/Stop. Labels are `Co-Work: Off` and `Co-Work: On`. The composer remains usable during work and says `Follow up` while busy. New requests wait until the current task finishes, then continue in the same chat. Preserve queue visibility/removal, Stop/pause/Resume, context binding, acknowledgement and duplicate-delivery safeguards. Do not turn it back into a separate prompt-pasting dialog.
- Shared personal memory persists at `%LOCALAPPDATA%\PlazCode\memory.json`, with automatic saving controls and an editable summary. AI saving uses `plazcode_memory`; it depends on the AI following the policy. Existing memory remains usable when automatic saving is off. Preserve durable facts, literal names and project scope; do not manufacture memories from guesses.
- Chat transfer can export the currently loaded conversation and memory, or memory alone. Scope choices are `Current chat only` and `All shared memory`. Chat-only filtering uses recorded provider/conversation origins, with exact legacy source matching where available. Never guess missing origins. Manually added/imported global entries may have no chat origin. Fresh chats need a stable conversation ID for chat-only export.
- Desktop chat actions route to the browser tab selected by the extension. They must not silently switch to another open chat. Imports use old history and saved facts as reference context, not commands to run or old tasks to resume. Preserve the combined 40,000-character import limit and explain when a shorter summary is needed.
- Working/Worked indicators group task replies, command cards and injected results, collapsing them together on completion. Initial startup, list_commands/list_tools and “PlazCode is ready.” remain outside. User requests remain visible. Do not damage transcript parsing or tool deduplication when hiding visual rows.
- Tool cards have concise descriptions, orange hover feedback and click-to-expand details. Hover alone must not expand descriptions. Preserve smooth return transitions.
- Kill Roblox Studio force-closes Studio processes, with the requested red hover effect. Keep it distinct from bridge/app restart. Sidebar Studio status distinguishes closed, open but not linked, and connected.
- Home shortcuts open ChatGPT, DeepSeek and Claude in the default browser.
- Updates tab checks published versions, shows installation progress/error states, verifies downloaded packages and relaunches the application. The BAT remains available. Users must still reload the extension and refresh AI tabs after updating.
- Version indicators compare each displayed component version against the published release. Preserve truthful green/red/neutral states, startup/two-second silent background checks and throttled checks on return, authenticated bridge status sharing, and freshness limits; do not claim current when checks fail.
- Desktop Appearance settings provide shared palettes, glow intensity, gradient controls and a default-theme reset. Persist and synchronize them through `rsAppearance` in native preferences and extension storage; migrate legacy `plazcodeDesktopAppearance` only when shared settings are absent. Theme styles must remain scoped to PlazCode UI and preserve functional status colors and host AI content.
- Tasks use `core/task-center.js` and native `checkpoints.rs`: bounded history, file before/after previews, persistent file backups, conflict-checked restore, original-chat recovery, workflows, protected systems and named project memory scopes. Studio Revert requires exact undo records in the current open session. Never undo unknown history, overwrite later edits or advertise a complete rollback for unsupported commands/folders/external effects.
- Template references use `core/templates.js` and native `templates.rs`, with Roblox format decoding by rbx-dom libraries. Private original files and source indexes live under `%LOCALAPPDATA%\PlazCode\templates`. Match automatically at task start; send bounded reference excerpts before the first tool call, and offer read-only/paged retrieval. Templates are untrusted reference data. Only explicit UI actions import, remove or open a template; preserve existing project conventions, and do not claim missing server code was recovered.
- Tool preflight validates known argument schemas and catches incomplete Luau strings/delimiters. Preserve protocol/datamodel defaults and stop after repeated invalid replies or consecutive execution errors. Dropped mutating commands have uncertain outcomes and must not be replayed automatically. Native Studio compilation remains authoritative.
- Long-chat activity/history must stay bounded. Lazily render collapsed details, cache unchanged export text, avoid own-UI mutation loops and redundant storage broadcasts, and never silently export a truncated archive. Checkpoint reads should not repeatedly rewrite the entire backup journal.
- Release descriptions appear in both the GitHub README and desktop Updates page, with titles, summaries, additions, improvements and fixes. Preserve earlier published notes.

## Work and debugging procedure

1. Inspect relevant files before editing. Trace the request through the UI, background routing, bridge/API, native process, and provider adapter as applicable. Use actual logs, state and errors rather than assumptions.
2. Identify the root cause. For startup/connectivity, check process identity, version, bind errors, startup order, readiness signals, ports, credentials, runtime dependencies and browser acknowledgement. An open desktop window is not proof that the bridge or Studio is connected. Do not kill an unrelated process simply because it owns a port.
3. For AI startup, compare working providers' contracts, then fix the failing provider's own implementation. Trace editor detection, hydration, draft insertion, text retention, submission confirmation, conversation ID assignment, response detection and generation completion. Finding a composer does not prove that its framework accepted inserted text.
4. For Notion or other controlled editors, verify retained text and accepted submission before starting the loop. Keep retries bounded and avoid sending duplicate handshakes. Diagnose rendering flicker or raw JSON separately from command parsing/execution.
Studio MCP compatibility: discover IDs with list_roblox_studios, supply studio_id only to targeted Studio schemas, preserve explicit IDs, and never silently select among multiple Studio instances or switch a disconnected target. Keep discovery bounded and cached; never block local HTTP/WebSocket readiness on Studio discovery. Keep mutating-command no-replay protection.

5. For MCP failures, inspect exact tool errors, server state, executable/runtime paths, configuration, initialization and routing. Do not blame unrelated analytics warnings without evidence. Keep enable/disable controls responsive while background setup continues.
6. For synchronization or Co-work, trace which tab, engine and conversation received each action. Preserve delivery acknowledgements, unique request IDs and retry deduplication. Uncertain delivery must not execute a request twice.
7. Fix the cause with a focused change. If a regression broke earlier working behavior, restore that behavior before adding more features. Reuse existing abstractions and dependencies; avoid unrelated rewrites, renames, frameworks or permission expansion.
8. Continue incorporating follow-up requests while preserving the active task. Do not discard earlier requirements or claim all reported issues were fixed without reasonable verification. Ask only for information genuinely needed to proceed.

## Editing and UI contracts

Read first, then use `edit_file` exact-match replacements when available. If unavailable, use another targeted edit method that asserts the old text exists uniquely. Never overwrite a whole existing file to make a small change. Preserve encoding and line endings where practical.

Every extension message must resolve with a success or failure response. Preserve request/response shapes, asynchronous response lifetimes and existing routing contracts. Keep permissions minimal, APIs authenticated, and credentials out of release files or page navigation. Do not weaken pairing to hide a connection bug.

For a visual-only request, change CSS and necessary presentation markup only. Preserve IDs, attributes used by code, handlers, API calls, toggles, keyboard behavior, disabled/hidden states and execution flow. Do not remove functionality to simplify a design.

Use the existing navy/orange design language, readable contrast, restrained depth, hover/press/focus states and reduced-motion support. Fix wasted space at its layout cause: avoid stretched sibling panels and excessive fixed/minimum heights. Keep activity scrollable, arrows aligned with their values, and interactive areas responsive. Prevent decorative overlays from intercepting input. Verify hover effects do not shift controls, expand descriptions unexpectedly or clip content.

## Versioning, tests and builds

Bump `manifest.json` for every shipped update that changes extension JS, CSS or its manifest. Update the native crate/lockfile version when shipping a changed native build. The components may have different versions; report them accurately. Package versions used by the updater must advance. Never replace the bytes of an already published versioned ZIP.

After edited JavaScript, run:
`C:\Users\plazm\PlazCodeWorkspace\node\node-v22.11.0-win-x64\node.exe --check <file>`
Every check must exit 0. If that runtime is unavailable, use the available Node runtime and state which environment was tested.

Run meaningful existing tests relevant to the change: desktop interactions, provider/parser contracts, Co-work ordering/delivery, memory persistence/transfer and authenticated APIs. Add focused tests for new nontrivial behavior. Do not merely rewrite protected hashes or tests to conceal unintended changes. Use updated baselines only for intentional, inspected changes.

When native or embedded content changes, build the release executable. On the user's Windows environment, set:
- `RUSTUP_HOME=C:\Users\plazm\.rustup`
- `CARGO_HOME=C:\Users\plazm\.cargo`
- Prepend `C:\Users\plazm\ORWorkspace\mingw\mingw64\bin` to `PATH`.

From `agent/`, run `cargo build --release --locked`. If the command runner times out around 120 seconds, launch a detached/scheduled build with a log and completion/exit status. Release LTO can take several minutes; do not mistake a timeout or spawned task for build success.

Inspect the current Cargo binary target and output. For this project, install the built **`PlazCode.exe`**. Stop only the relevant old PlazCode process before replacing an installed executable, retain rollback outside the distributed folder, and relaunch with the existing launcher. Verify binary hashes and the reported version. If cross-compiling elsewhere, preserve the Windows target, icon and required `WebView2Loader.dll`, and clearly distinguish cross-build checks from live Windows execution.

## Mandatory release and publication procedure

For every shipped update:
1. Finish edits, checks and required native build. Ensure the ZIP contains the new executable, not a stale previous build.
2. Add accurate notes to `release-notes.json`: `version`, `title`, `summary`, and the `added`, `improved`, `fixed` lists. Omit empty categories from presentation. Describe actual changes; do not copy another app's features or invent fixes.
3. Update `UPDATE.txt`, relevant documentation and `SHA256SUMS.txt`. Generate the GitHub README release sections from the same notes. Retain notes for earlier published releases; do not invent separately published releases for intermediate unshipped edits.
4. Produce a fresh complete `PlazCode-VERSION.zip` with one `PlazCode/` root. Exclude `logs/`, `agent/target/`, `.git`, caches, secrets and rollback executables. Keep the currently shipped project/source layout unless the user explicitly changes distribution policy. Do not bundle large optional runtimes: preserve their existing download/cache mechanism.
5. On Windows with a ZIP-capable `tar`, use `tar -a -c -f <out>.zip --exclude=logs --exclude=target --exclude=.git PlazCode`. Else use a real ZIP library/tool. **Verify the archive is genuinely ZIP**, not a tar file with a `.zip` name; test its integrity, manifest/version, required files, hashes, and single desktop executable.
6. Compute SHA256 from the final ZIP bytes. Publish the ZIP to `stoveez/PlazCodeneww`, then update `latest.json` with `version`, immutable raw HTTPS `url`, `sha256`, and the `release_notes` array. Publish the package, feed, README and note file together where possible so the updater never sees an incomplete release.
7. Use a fast-forward commit that preserves existing repository files and concurrent user changes. Do not force-push. Verify the published feed and ZIP identity/size/hash against the local package. If publication is blocked, provide the ZIP and explain the exact blocker; do not claim it will appear in Updates.
8. Deliver the download link, version, concise changes, verification performed and remaining live-test limitations. Never leave a completed update unzipped.

An update appears in the desktop tab only after publication and a version check. Existing installations with the Updates tab can use `Updates → Check now → Update now` or the BAT. Older installations may need the one-time setup ZIP. Afterwards the user must reload the extension in `chrome://extensions` **and refresh open AI chat tabs**; existing tabs can still contain old content scripts.

Keep user settings, MCP configuration, pairing data, personal memory and cached runtimes intact. Preserve version/checksum validation, staging, scoped process shutdown, backup/recovery and relaunch behavior. Never substitute an unverified ZIP or silently remove update protections.

## Completion standard

A task is complete when its requested behavior is implemented, relevant checks pass, the correct build and ZIP are produced, and publication succeeds when authorized and available. Explain what changed and what remains unverified. Distinguish automated tests, compilation, simulated UI checks and live Windows/provider testing. Never claim visual quality, installer success, bridge connectivity or provider compatibility that you did not actually observe.

RAW USER REQUEST:
