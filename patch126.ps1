$ErrorActionPreference='Stop'
$s='C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\build-src\PlazCode\agent\src'
$bk='C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\backups-tab5110be'
Copy-Item "$s\notion_team.rs" "$bk\notion_team.rs.126pre" -Force
Copy-Item "$s\notion_profiles.rs" "$bk\notion_profiles.rs.126pre" -Force
$enc=New-Object System.Text.UTF8Encoding($false)
function Rep([string]$txt,[string]$a,[string]$b,[string]$what){ $n=([regex]::Matches($txt,[regex]::Escape($a))).Count; if($n -lt 1){throw "nicht gefunden: $what"}; "ok $what ($n)"|Out-Host; return $txt.Replace($a,$b) }
# ---- notion_team.rs
$t=[IO.File]::ReadAllText("$s\notion_team.rs",$enc)
$t=Rep $t 'if id == main { why.push(format!("{short}: Haupt-Chat")); continue; }' 'if id == main && running.contains(&id) { continue; }' 'Haupt-Chat einladbar'
$t=Rep $t 'if open.contains(&id) { why.push(format!("{short}: offen im Fenster (Nutzer)")); continue; }' 'let _ = open.contains(&id);' 'offene Chats einladbar'
$t=Rep $t '    let (ws, _) = team_open(port, tid).await?;' '    let (ws, opened) = team_open(port, tid).await?;' 'deliver opened'
$t=Rep $t 'if ws_eval(&ws, "document.hasFocus()").await.ok().and_then(|v| v.as_bool()) == Some(true) { return Err("Nutzer ist gerade in diesem Chat".into()); }' 'if opened.is_none() && ws_eval(&ws, TEAM_USER_JS).await.ok().and_then(|v| v.as_bool()) == Some(true) { return Err("Nutzer ist gerade in diesem Chat".into()); }' 'deliver focus'
$t=Rep $t '    if !team_ready(&ws, tid).await { return fail(' '    if opened.is_none() && ws_eval(&ws, TEAM_USER_JS).await.ok().and_then(|v| v.as_bool()) == Some(true) { return Err("Nutzer ist gerade in diesem Chat".into()); }
    if !team_ready(&ws, tid).await { return fail(' 'invite focus'
$t=Rep $t 'const TEAM_FRESH_SECS: i64 = 45;' 'const TEAM_FRESH_SECS: i64 = 45;
/// Nutzer sitzt wirklich in diesem Chat: Tab sichtbar + fokussiert + Eingabe hat Text oder Fokus.
const TEAM_USER_JS: &str = "(()=>{if(document.visibilityState!==''visible''||!document.hasFocus())return false;const a=document.activeElement;return !!(a&&(a.isContentEditable||a.tagName===''TEXTAREA''||a.tagName===''INPUT''))})()";' 'TEAM_USER_JS'
$t=Rep $t '            lock(&TEAM_FAIL).insert(tid.clone(), std::time::Instant::now());' '            if !e.contains("Nutzer") { lock(&TEAM_FAIL).insert(tid.clone(), std::time::Instant::now()); }' 'kein Cooldown wenn Nutzer'
[IO.File]::WriteAllText("$s\notion_team.rs",$t,$enc)
# ---- notion_profiles.rs: nur Brave + mehr Tracker/Fingerprint-Hosts blocken
$p=[IO.File]::ReadAllText("$s\notion_profiles.rs",$enc)
$p=Rep $p 'const BROWSERS: [(&str, &str); 3] = [("Chrome", r"Google\Chrome\Application\chrome.exe"), ("Edge", r"Microsoft\Edge\Application\msedge.exe"), ("Brave", r"BraveSoftware\Brave-Browser\Application\brave.exe")];' 'const BROWSERS: [(&str, &str); 1] = [("Brave", r"BraveSoftware\Brave-Browser\Application\brave.exe")];' 'nur Brave'
$p=Rep $p 'pub fn browser() -> Option<(&''static str, PathBuf)> { let c = chosen_browser()?; installed_browsers().into_iter().find(|(n, _)| n.eq_ignore_ascii_case(&c)) }' 'pub fn browser() -> Option<(&''static str, PathBuf)> { let _ = chosen_browser(); installed_browsers().into_iter().next() }' 'browser immer Brave'
$extra=', MAP *.intercom.io 127.0.0.2, MAP *.intercomcdn.com 127.0.0.2, MAP www.googletagmanager.com 127.0.0.2, MAP *.google-analytics.com 127.0.0.2, MAP *.doubleclick.net 127.0.0.2, MAP connect.facebook.net 127.0.0.2, MAP *.sentry.io 127.0.0.2, MAP *.datadoghq.com 127.0.0.2, MAP *.splunkcloud.com 127.0.0.2, MAP *.fpjs.io 127.0.0.2, MAP *.fpnpmcdn.net 127.0.0.2, MAP *.fingerprint.com 127.0.0.2, MAP *.hotjar.com 127.0.0.2, MAP *.fullstory.com 127.0.0.2, MAP *.clarity.ms 127.0.0.2, MAP *.hs-scripts.com 127.0.0.2, MAP *.hs-analytics.net 127.0.0.2, MAP static.ads-twitter.com 127.0.0.2, MAP *.pendo.io 127.0.0.2, MAP *.statsig.com 127.0.0.2, MAP *.statsigapi.net 127.0.0.2'
$p=Rep $p 'MAP cdn.pendo.io 127.0.0.2"' ('MAP cdn.pendo.io 127.0.0.2' + $extra + '"') 'Tracker-Blockliste'
[IO.File]::WriteAllText("$s\notion_profiles.rs",$p,$enc)
'--- browser.txt'; Get-ChildItem "$env:LOCALAPPDATA\PlazCodeNotion" -Recurse -Depth 2 -Filter 'browser*' -ErrorAction SilentlyContinue | % { $_.FullName+' = '+(Get-Content $_.FullName -Raw) }
'--- TEAM tab1 (Einladen)'; Select-String -Path "$env:LOCALAPPDATA\PlazCodeNotion\profiles\timeline.log" -Pattern 'tab1 TEAM [^-]|TEAM-INVITE' | Select -Last 12 | % { $_.Line -replace '^(.{0,300}).*$','$1' }
$ErrorActionPreference='Continue'
Set-Location 'C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\build-src\PlazCode\agent'
cmd /c "cargo test notion_ > C:\Users\liket\PlazCode-Shared\work\t126.log 2>&1"
"TEST_EXIT $LASTEXITCODE"
Select-String -Path C:\Users\liket\PlazCode-Shared\work\t126.log -Pattern '^error|test result|FAILED' | Select -First 30 | % { $_.Line }
$l=Get-Content C:\Users\liket\PlazCode-Shared\work\t126.log; $e=Select-String -Path C:\Users\liket\PlazCode-Shared\work\t126.log -Pattern '^error\[' | Select -First 3; foreach($h in $e){ $l[($h.LineNumber-1)..($h.LineNumber+10)] }
