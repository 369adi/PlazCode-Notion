$ErrorActionPreference='Stop'
$s='C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\build-src\PlazCode\agent\src'
$bk='C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\backups-tab5110be'
Copy-Item "$s\notion_profiles.rs" "$bk\notion_profiles.rs.125team-pre" -Force
Copy-Item "$s\notion_cowork.rs" "$bk\notion_cowork.rs.125team-pre" -Force
$enc=New-Object System.Text.UTF8Encoding($false)
# --- notion_profiles.rs
$l=[System.Collections.Generic.List[string]]([System.IO.File]::ReadAllLines("$s\notion_profiles.rs",$enc))
$a=-1;$b=-1
for($i=0;$i -lt $l.Count;$i++){ if($a -lt 0 -and $l[$i].TrimStart().StartsWith('/// Zusatz-Chats: Generier-Status')){$a=$i}; if($l[$i].StartsWith('pub(crate) fn ensure_multichat')){$b=$i;break} }
if($a -lt 0 -or $b -lt 0 -or $b -le $a){throw "Bereich nicht gefunden a=$a b=$b"}
$removed=$b-$a
$l.RemoveRange($a,$removed)
"entfernt: $removed Zeilen (slot_monitor + auto_fill)"
$idx=$l.FindIndex([Predicate[string]]{param($x) $x.StartsWith('pub(crate) fn ensure_multichat')})
$old=$l[$idx]
if($old -notmatch 'tokio::spawn\(slot_monitor\(\)\)'){throw 'slot_monitor-Aufruf nicht gefunden'}
$l[$idx]=$old.Replace('tokio::spawn(slot_monitor())','tokio::spawn(team_loop())')
$l.Add('include!("notion_team.rs");')
[System.IO.File]::WriteAllLines("$s\notion_profiles.rs",$l,$enc)
"ensure_multichat -> team_loop, include ok"
# --- notion_team.rs: ensure_team als allow(dead_code)
$t=[System.IO.File]::ReadAllText("$s\notion_team.rs",$enc)
$t=$t.Replace('pub(crate) fn ensure_team()','#[allow(dead_code)]' + "`n" + 'pub(crate) fn ensure_team()')
[System.IO.File]::WriteAllText("$s\notion_team.rs",$t,$enc)
# --- notion_cowork.rs: active_project
$c=[System.Collections.Generic.List[string]]([System.IO.File]::ReadAllLines("$s\notion_cowork.rs",$enc))
$j=$c.FindIndex([Predicate[string]]{param($x) $x.StartsWith('pub(crate) fn joined_agents')})
if($j -lt 0){throw 'joined_agents nicht gefunden'}
if(-not ($c -match 'fn active_project')){
 $c.Insert($j+1,'/// Co-Work-Team: neuestes Projekt mit offenen Tasks (id, Titel) - dorthin werden freie Chats eingeladen.')
 $c.Insert($j+2,'pub(crate) fn active_project() -> Option<(String, String)> { let b=BOARD.lock().unwrap_or_else(|e|e.into_inner()); b.projects.values().filter(|p|p.tasks.iter().any(|t|t.status!="done")).max_by_key(|p|p.created).map(|p|(p.id.clone(),p.title.clone())) }')
}
[System.IO.File]::WriteAllLines("$s\notion_cowork.rs",$c,$enc)
"active_project ok"
Set-Location 'C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\build-src\PlazCode\agent'
cargo test notion_ 2>&1 | Select-String -Pattern '^(error|warning: unused|warning: function .* never used|warning: constant .* never used)|-->|test result|FAILED|panicked' | Select -First 60 | % { $_.Line }
"TEST_EXIT $LASTEXITCODE"
