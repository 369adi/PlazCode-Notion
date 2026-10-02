# VoidScript 6.0.1 compared with PlazCode 1.18.97

This audit reads the supplied VoidScript ZIP and PlazCode source. It distinguishes implemented browser/bridge code from desktop features described only in release notes: the VoidScript desktop executable is included, but its full native source is not. There are 119 archive entries and 78 provider-related JavaScript files; file counts do not establish working provider coverage. These are recommendations, not features newly added to PlazCode by this audit.

## Best additions to prioritize

1. **UI and model previews before insertion.** Extend PlazCode's existing GUI/model tools with editable previews, saved assets, and revisions. This offers a clear workflow improvement without replacing working tool execution.
2. **Per-script change review.** Expand desktop Tasks/checkpoints with individual diffs, conflict checks and selective restore. Keep Tasks out of the browser bar.
3. **Context compaction and handoff.** Preserve unfinished goals, verified state and exact paths before a long chat runs out of context. Avoid repeating a mutation whose result is uncertain.
4. **Beginner-friendly tools and project profiles.** Schema-driven tool forms, reusable named workflows, and place/workspace-specific instructions.
5. **Visual verification and safety dashboard.** Reuse screenshots, diagnostics and remote auditing; add evidence-based checks and a bounded automatic visual review.

## Full comparison

“Partial” means PlazCode already has a related capability; the difference is the specific workflow described below. Desktop-only claims require testing against the executable before treating them as confirmed implementation details.

| Capability in VoidScript | PlazCode today / gap | A better PlazCode implementation |
|---|---|---|
| Model generator: constrained part JSON, 3D orbit/zoom preview, revisions, Luau export and insertion (`modelkit.js`, `modelview.js`) | Partial: model/image and Blender helpers; no equivalent dedicated editable preview library | Stable part IDs, selected-part revisions, weld/collision checks, reusable models and Studio visual validation |
| UI builder: bounded GUI JSON, HTML preview, four styles, revisions and insertion (`uikit.js`) | Partial: GUI tools and ten styles already exist | Editable preview at multiple screen sizes, layout validation, drag/resize and reusable UI library; preview alone does not prove Luau behavior |
| Seeded terrain generator: six biomes, lakes and region sizes (`toolkit.js`) | Partial: terrain fill/clear | Bounded region preview, chunked generation, progress, cancellation and checkpointed restore |
| Manual lighting mood picker | Lighting presets already exist as tools; desktop picker is the gap | Preview existing presets and expose their controls without duplicating tools |
| Creator Store asset browser with insertion | AI asset search/import already exists; manual browser differs | Filters, source IDs and embedded-script inspection before insertion |
| Manual Luau editor/console with Play/Stop | Luau execution exists; desktop Terminal is primarily output | Explicit target, clear execution result, change preview and meaningful Play controls |
| Manual tool tester (“Try in terminal,” desktop release notes) | Tools catalog exists; schema form runner differs | Required-field forms, read/write labels, progress and exact tool errors |
| Health/security dashboard with grades and suspicious-code checks (`toolkit.js`) | Partial: diagnostics, lint and remote auditing | Unified dashboard with file/line evidence and false-positive handling; never blindly delete flagged code |
| Whole-place disk backup and restore (`bridge.py`) | Partial: file checkpoints and Studio undo | Verify a save before backup; restore to a separate copy. VoidScript's disk copy does not capture unsaved Studio state |
| Per-script Changes and individual Undo | Task-level checkpoints exist | Individual diffs and selective restore with conflict detection |
| Named Luau macros (`save_macro`, `run_macro`, `list_macros`) | Templates are references, not named executable macros | Typed parameters, project scope, validation and versioned reusable actions |
| Many named prompt snippets | Partial: built-in cards and a custom workflow | User collections with tags, search and clear beginner labels |
| Named settings profiles and per-place instructions | Global startup instructions and project memory exist | Stable place/workspace profiles with explicit overrides; exclude secrets from exports |
| Settings JSON export/import | Updates preserve settings; memory transfer exists | Portable validated settings export, merge preview and credential exclusion |
| Shareable build recipes encoded in links | No comparable guided recipe importer | Inspect a recipe before applying its settings or executing anything |
| Three-question build wizard | Prompt enhancer and starter cards exist | Short guided questions with useful defaults and automatic template matching |
| Context compaction/handoff (`compactNow`) | Context warning and manual memory/history transfer exist | Proactive bounded handoff with unfinished goals and verified state. VoidScript still requires a new chat; it is not fully automatic chat migration |
| Up to five command calls collected in one feedback | PlazCode uses one command per reply | Optional validated batches for independent reads, with ordered writes and per-step results. VoidScript's execution is sequential, not parallel |
| Inline tool catalog for startup (5.4 release notes) | Existing list/ready handshake is protected | Consider catalog hashes/deltas only after live ChatGPT, Claude, DeepSeek and Notion regression testing |
| Automatic visual checks after visible changes | Screenshots, automatic debugger and Motion QA exist | Bounded, deduplicated screenshot review after relevant changes; avoid unnecessary uploads/context growth |
| Command budget with Pause/Resume | Error guards exist; explicit runtime/call budget differs | Transparent call/time limits and loop detection with a resumable task state |
| Granular Studio trust levels | AgentScript permissions and task approvals already exist | Extend previews and action policies specifically to Studio operations |
| UI translation and AI reply language | Mainly English UI | Localized labels while preserving code, paths, commands and error literals |
| Voice dictation through browser speech recognition | No equivalent integrated workflow identified | Optional dictation with editable preview before sending |
| Browser completion notifications and spoken finish; desktop notification center documented | Sounds exist | Quiet hours, completion notifications, optional speech and reopening the correct chat |
| Sound volume and preview (desktop release notes) | Sound toggle exists | Volume, preview and separate completion/error preferences |
| Personal provider leaderboard and usage export | Session/activity information exists | Measured recent reliability by model/task type, without unsupported “intelligence” rankings |
| Browser session picker | Selected active browser-agent state exists | Explicit provider/chat/workspace selection and busy-state protection against wrong-chat sends |
| Broader provider adapter set | Thirteen supported providers; more files do not prove compatibility | Add selected providers with genuine startup/tool/media tests. Do not restore Copilot or Meta AI, which the user removed |
| Broader MCP setup catalog: 29 entries, 28 unique | Six native presets plus custom MCP already exist | Guided Figma, Notion, GitHub, Sketchfab or Playwright setup with schema validation, health checks and clear credential handling |
| macOS launcher (`Start.command`, Darwin launcher code) | Windows desktop target | A separately tested native distribution. Firefox and desktop API-model claims in documentation are not proven by the inspected browser files alone |
| OS-login startup install/uninstall/status (`bridge.py`) | Opening PlazCode starts the bridge; OS-login option differs | Optional per-user login startup with visible status and clean uninstall |
| Configurable hotkeys and runnable/copyable Luau snippets | Limited shortcuts exist | User-defined shortcuts with collision checks and clearly scoped execution |
| Session timeline export/archive and auto summary | Activity, logs and memory/history export exist | Searchable session archive with retention controls, exact operation results and separate concise summaries |
| Quick steering chips: Fix, Undo, Retry, Keep going, Explain | Co-work already supports additional active requests | Context-aware shortcuts that preserve unfinished goals and inspect uncertain state before retrying |

## Things not to count as finished advantages

The inspected recording/playback UI contains buttons and an empty recorder hook, but no populated recording pipeline was identified: its recorded array is never filled by a functioning operation recorder. Do not claim verified replay support from that UI. A real PlazCode replay feature would need validated operations, checkpoint mapping, dry-run review and conflict detection.

A model preview is not proof that collisions or joints work in Studio. An HTML GUI preview is not proof that LocalScripts, remotes or input behavior work. A security grade is not proof that a place is safe. Broad provider counts and desktop release notes need live verification.

## Existing PlazCode strengths to preserve

PlazCode already includes selectable themes, personal/project memory and transfer, automatic template matching, Co-work, work/thinking/plan controls, specialist helpers, automatic debugging/fixing, AgentScript file and terminal tools, modern Studio-ID routing, update notes, work cards, GUI helpers, lighting/terrain tools, command validation, Blender/animation/VFX support and Motion QA. Extend these systems rather than introducing duplicate frameworks or replacing working contracts.

## Performance requirements for any additions

Keep provider DOM access inside provider modules. Use bounded queues and caches, incremental tracking, debounce expensive work, and avoid repeatedly scanning an entire long transcript. Pause unnecessary background work in hidden tabs. Parse commands from original text, never folded presentation. Limit images, previews and exported history explicitly; report omitted content. Keep tool validation before execution, preserve exact literals, and never silently resend an uncertain mutation. Test long chats, navigation, Stop, startup and update preservation whenever a change crosses these boundaries.
