$ErrorActionPreference='Stop'
$f='C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\build-src\PlazCode\agent\src\notion_profiles.rs'
$bk='C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\backups-tab5110be'
New-Item -ItemType Directory -Force $bk | Out-Null
Copy-Item $f "$bk\notion_profiles.rs.124pre" -Force
$u=New-Object Text.UTF8Encoding($false)
$s=[IO.File]::ReadAllText($f,$u)
function Rep($old,$new,[switch]$All){ $c=([regex]::Matches($script:s,[regex]::Escape($old))).Count; if(($All -and $c -lt 1) -or (-not $All -and $c -ne 1)){ throw "Treffer $c fuer: $old" }; $script:s=$script:s.Replace($old,$new); "ok ($c): " + $old.Substring(0,[Math]::Min(60,$old.Length)) }

Rep 'fn running_ids(tabs: &[Tab]) -> Vec<u32> { let cmds=browser_cmds(); tabs.iter().filter(|t| {let d=dir(t.id).display().to_string();cmds.iter().any(|(_,c)|cmd_uses(c,&d))}).map(|t|t.id).collect() }' ('fn running_ids(tabs: &[Tab]) -> Vec<u32> { let cmds=browser_cmds(); tabs.iter().filter(|t| {port_open(port_for_id(t.id))||{let d=dir(t.id).display().to_string();cmds.iter().any(|(_,c)|cmd_uses(c,&d))}}).map(|t|t.id).collect() }' + "`n" + '/// 1.0.124: Tab gilt auch als laufend, wenn sein CDP-Port antwortet (die WebView2-Kommandozeile enthaelt den Profilordner nicht immer -> Co-Work holte nie Zusatz-Chats).' + "`n" + 'fn port_open(port:u16)->bool{ if cfg!(test){return false;} std::net::TcpStream::connect_timeout(&std::net::SocketAddr::from(([127,0,0,1],port)),Duration::from_millis(150)).is_ok() }')

Rep 'pub(crate) fn tab_ai_off(id:u32)->bool{lock(&AI_OFF).get(&id).copied().unwrap_or(false)}' ('pub(crate) fn tab_ai_off(id:u32)->bool{lock(&AI_OFF).get(&id).copied().unwrap_or(false)}' + "`n" + '/// 1.0.124: Konto-Tab ohne AdiCode-MCP (Einrichtung fehlgeschlagen) kann keine Co-Work-Tools aufrufen -> nicht eintragen, keine Zusatz-Chats.' + "`n" + 'pub(crate) fn tab_no_mcp(id:u32)->bool{ if cfg!(test){return false;} let a=lock(&ACCOUNT); let Some(v)=a.get(&id) else {return false}; if v["mcp_saved"].as_bool()==Some(true){return false;} v["mcp"].as_str().map(|s|s.starts_with("Fehlt")||s.starts_with("Auto-Einrichtung fehlgeschlagen")).unwrap_or(false) }')

Rep 'Some((i,t))=>skipped(t.id)||tab_ai_off(t.id)||tab_offline(i)' 'Some((i,t))=>skipped(t.id)||tab_ai_off(t.id)||tab_no_mcp(t.id)||tab_offline(i)'
Rep 'if tab_ai_off(t.id){continue;}' 'if tab_ai_off(t.id)||tab_no_mcp(t.id){continue;}' -All

[IO.File]::WriteAllText($f,$s,$u)
'PATCH124_OK'
Set-Location 'C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\build-src\PlazCode\agent'
$ErrorActionPreference='Continue'
cargo test --release --locked notion 2>&1 | Select-String -Pattern 'error|warning: unused|test result|FAILED|panicked' | Select -First 40 | % { $_.Line }
'TEST_EXIT ' + $LASTEXITCODE
