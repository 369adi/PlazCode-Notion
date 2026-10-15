$ErrorActionPreference='Stop'
$root='C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\build-src\PlazCode\agent\src'
$bk='C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion\backups-tab5110be'
New-Item -ItemType Directory -Force $bk | Out-Null
foreach($f in 'notion_profiles.rs','notion_cowork.rs'){ Copy-Item (Join-Path $root $f) (Join-Path $bk ($f+'.bak_cw123')) -Force }
$files=@{}
function Rep($file,$old,$new,$count=1){
  if(-not $files.ContainsKey($file)){ $files[$file]=[IO.File]::ReadAllText((Join-Path $root $file)) }
  $s=$files[$file]
  $old=$old -replace "`r`n","`n"; $new=$new -replace "`r`n","`n"
  if($s.Contains("`r`n")){ $old=$old -replace "`n","`r`n"; $new=$new -replace "`n","`r`n" }
  $n=([regex]::Matches($s,[regex]::Escape($old))).Count
  if($n -ne $count){ throw "$file : erwartet $count Treffer, gefunden $n fuer: $($old.Substring(0,[Math]::Min(90,$old.Length)))" }
  $files[$file]=$s.Replace($old,$new)
  "ok $file : $($old.Substring(0,[Math]::Min(60,$old.Length)))"
}

# ---------------- notion_profiles.rs
Rep 'notion_profiles.rs' @'
static AI_OFF: Lazy<Mutex<BTreeMap<u32,bool>>> = Lazy::new(|| Mutex::new(BTreeMap::new()));
'@ @'
static AI_OFF: Lazy<Mutex<BTreeMap<u32,bool>>> = Lazy::new(|| Mutex::new(BTreeMap::new()));
/// Co-Work (1.0.123): Zusatz-Chats, die gerade nicht schreiben koennen (Tab zu oder "KI fuer diesen Workspace deaktiviert").
static SLOT_OFF: Lazy<Mutex<HashSet<String>>> = Lazy::new(|| Mutex::new(HashSet::new()));
static UNAVAIL: Lazy<Mutex<BTreeMap<String,(std::time::Instant,bool)>>> = Lazy::new(|| Mutex::new(BTreeMap::new()));
/// Tab-Fenster laeuft nicht (kein CDP-Port).
fn tab_offline(i:usize)->bool{std::net::TcpStream::connect_timeout(&std::net::SocketAddr::from(([127,0,0,1],cdp_port(i))),Duration::from_millis(200)).is_err()}
/// Konto-Tab zeigt "Notion KI ist fuer diesen Workspace deaktiviert".
pub(crate) fn tab_ai_off(id:u32)->bool{lock(&AI_OFF).get(&id).copied().unwrap_or(false)}
/// Co-Work: Agent kann gerade nicht schreiben (Fenster zu, uebersprungen, KI im Tab oder im Zusatz-Chat deaktiviert)
/// -> wird nicht eingetragen, nicht geweckt, seine Tasks werden wieder frei. Agents ohne AdiCode-Tab (freie Namen) gelten als verfuegbar.
pub(crate) fn agent_unavailable(agent:&str)->bool{
    if cfg!(test){return false;}
    if lock(&SLOT_OFF).contains(agent){return true;}
    if let Some((at,v))=lock(&UNAVAIL).get(agent).copied(){if at.elapsed()<Duration::from_secs(5){return v;}}
    let tabs=load();let base=agent.split("-K").next().unwrap_or(agent);
    let v=match tabs.iter().enumerate().find(|(_,t)|agent_of(t)==base){Some((i,t))=>skipped(t.id)||tab_ai_off(t.id)||tab_offline(i),None=>false};
    lock(&UNAVAIL).insert(agent.to_string(),(std::time::Instant::now(),v));v
}
'@

Rep 'notion_profiles.rs' @'
    if !hit{if was{tl(&format!("tab{} AI-DISABLED-Hinweis ist wieder weg",t.id));}return;}
'@ @'
    if hit!=was{mark_account(t.id,"ai_off",json!(hit));}
    if !hit{if was{tl(&format!("tab{} AI-DISABLED-Hinweis ist wieder weg",t.id));}return;}
'@

Rep 'notion_profiles.rs' 'if n%10==1{ai_off_check(t,i).await;}' 'if n%3==1{ai_off_check(t,i).await;}'

Rep 'notion_profiles.rs' 'if cfg!(test){return false;}let text=with_mcp_hint(agent,text);' 'if cfg!(test)||agent_unavailable(agent){return false;}let text=with_mcp_hint(agent,text);'

Rep 'notion_profiles.rs' '.flat_map(slot_agents));v}' '.flat_map(slot_agents));v.retain(|a|!agent_unavailable(a));v}' 2

Rep 'notion_profiles.rs' 'if targets.is_empty(){for j in 0..sl.len(){crate::notion_live::set_generating(&slot_agent(t,j+2),false);}continue;}' 'if targets.is_empty(){for j in 0..sl.len(){crate::notion_live::set_generating(&slot_agent(t,j+2),false);lock(&SLOT_OFF).insert(slot_agent(t,j+2));}continue;}'

Rep 'notion_profiles.rs' @'
                    crate::notion_live::set_generating(&agent,false);
                    // Co-Work oeffnet NIE selbst
'@ @'
                    crate::notion_live::set_generating(&agent,false);
                    lock(&SLOT_OFF).insert(agent.clone());
                    // Co-Work oeffnet NIE selbst
'@

Rep 'notion_profiles.rs' @'
                match ws_eval(&ws,SLOT_GEN_JS).await
'@ @'
                // Chat offen, aber "KI fuer diesen Workspace deaktiviert": er kann nicht schreiben -> aus Co-Work nehmen, bis der Hinweis weg ist.
                let off=ws_eval(&ws,AI_OFF_JS).await.ok().map(|v|v.is_object()).unwrap_or(false);
                {let mut so=lock(&SLOT_OFF);let was=so.contains(&agent);if off&&!was{tl(&format!("{agent}: KI fuer diesen Workspace deaktiviert - aus Co-Work genommen"));}if off{so.insert(agent.clone());}else{so.remove(&agent);}}
                if off{crate::notion_live::set_generating(&agent,false);continue;}
                match ws_eval(&ws,SLOT_GEN_JS).await
'@

Rep 'notion_profiles.rs' 'if crate::notion_cowork::project_count()==0{return;}' 'if crate::notion_cowork::project_count()==0||crate::notion_cowork::open_task_count()==0{return;}'

Rep 'notion_profiles.rs' 'match crate::notion_cowork::nudge_for(&agent_of(t)){Some((sig,_)) if !sig.starts_with("done:")=>{},_=>continue}' 'if tab_ai_off(t.id){continue;}'

Rep 'notion_profiles.rs' 'if last_fill.elapsed()>Duration::from_secs(40){' 'if last_fill.elapsed()>Duration::from_secs(15){'

# ---------------- notion_cowork.rs
Rep 'notion_cowork.rs' @'
pub(crate) fn joined_agents() -> Vec<String> { let b=BOARD.lock().unwrap_or_else(|e|e.into_inner());b.projects.values().flat_map(|p|p.agents.keys().cloned()).collect() }
'@ @'
pub(crate) fn joined_agents() -> Vec<String> { let b=BOARD.lock().unwrap_or_else(|e|e.into_inner());b.projects.values().flat_map(|p|p.agents.keys().cloned()).collect() }
/// Offene (nicht erledigte) Tasks ueber alle Projekte - solange > 0, nutzt jedes Konto alle freien Chats.
pub(crate) fn open_task_count() -> usize { let b=BOARD.lock().unwrap_or_else(|e|e.into_inner());b.projects.values().map(|p|p.tasks.iter().filter(|t|t.status!="done").count()).sum() }
'@

Rep 'notion_cowork.rs' '    p.agents.get(agent).map(|a| now() - a.last_seen < stale(p)).unwrap_or(false)' '    p.agents.get(agent).map(|a| now() - a.last_seen < stale(p)).unwrap_or(false) && !crate::notion_profiles::agent_unavailable(agent)'

Rep 'notion_cowork.rs' @'
    if agent.is_empty() { return None; }
    let mut board = lock_board();
'@ @'
    if agent.is_empty() || crate::notion_profiles::agent_unavailable(agent) { return None; }
    let mut board = lock_board();
    // Jeder freie Chat hilft bei JEDEM Co-Work (Chat-Kontext egal): eigenes Projekt fertig, neueres hat offene Tasks -> dort beitreten.
    let joined_open = board.projects.values().any(|p| p.agents.contains_key(agent) && p.tasks.iter().any(|x| x.status != "done"));
    if !joined_open && board.projects.values().any(|p| p.agents.contains_key(agent)) {
        if let Some(p) = board.projects.values().filter(|p| !p.agents.contains_key(agent) && p.tasks.iter().any(|x| x.status != "done")).max_by_key(|p| p.created) {
            return Some((format!("join:{}", p.id), format!("Kurze automatische Erinnerung von AdiCode (von mir eingerichtet): Neues Co-Work-Projekt {} \"{}\" laeuft und braucht Hilfe (egal, was vorher in diesem Chat stand). Lade adicode_skill name=cowork, rufe cowork_join (project {}, agent {agent}) auf und arbeite dann in Schleife: cowork_claim -> erledigen -> cowork_update -> cowork_claim ..., bis nichts mehr frei ist; dann cowork_wait.", p.id, short(&p.title, 80), p.id)));
        }
    }
'@

foreach($k in $files.Keys){ [IO.File]::WriteAllText((Join-Path $root $k),$files[$k],(New-Object Text.UTF8Encoding $false)); "geschrieben: $k" }
