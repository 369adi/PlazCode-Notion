const fs=require('fs');const {chromium}=require('playwright');
(async()=>{
 const original=fs.readFileSync('agent/src/desktop.html','utf8');
 const themeScript='function Byid(id){return document.getElementById(id)};'+original.slice(original.indexOf('  var Desktopthemes ='),original.indexOf('  function Saveappearance()'));
 const html=original.replaceAll('__PLAZCODE_VERSION__','1.19.25').replace(/<script>[\s\S]*?<\/script>/g,'');
 const browser=await chromium.launch({headless:true});const page=await browser.newPage();
 await page.setContent(html);
 await page.addScriptTag({content:themeScript});
 await page.evaluate(()=>Applyappearance({theme:'default'}));
 fs.mkdirSync('visual-checks',{recursive:true});
 for(const width of [1440,1024,760]){
  await page.setViewportSize({width,height:960});
  for(const name of ['home','settings','tools','creators']){
   await page.evaluate(name=>{document.querySelectorAll('.page').forEach(n=>n.classList.toggle('active',n.id==='page-'+name));document.querySelectorAll('.navbtn[data-page]').forEach(n=>n.classList.toggle('active',n.dataset.page===name));},name);
   await page.screenshot({animations:"disabled",path:`visual-checks/${name}-${width}.png`});
   const overflow=await page.evaluate(()=>{const n=document.getElementById('content');return n.scrollWidth>n.clientWidth+1});
   if(overflow)throw Error(`${name} overflows at ${width}`);
  }
 }
 await page.setViewportSize({width:1440,height:960});
 for(const theme of ['ocean','copper','aurora','orchid','solar']){
  await page.evaluate(theme=>{Applyappearance({theme});document.querySelectorAll('.page').forEach(n=>n.classList.toggle('active',n.id==='page-home'));},theme);
  await page.screenshot({animations:"disabled",path:`visual-checks/theme-${theme}.png`});
 }
 for(const deviceScaleFactor of [1,2]) {
  const sample=await browser.newPage({viewport:{width:1440,height:960},deviceScaleFactor});
  await sample.setContent(html);await sample.addScriptTag({content:themeScript});
  for(const theme of ['default','ocean','orchid']) {
   for(const active of [false,true]) {
    await sample.evaluate(({theme,active})=>{Applyappearance({theme});document.querySelectorAll('.navbtn[data-page]').forEach(n=>n.classList.toggle('active',active&&n.dataset.page==='settings'));},{theme,active});
    const setting=sample.locator('.navbtn[data-page="settings"]');
    await setting.screenshot({animations:'disabled',path:`visual-checks/settings-icon-${theme}-${active?'selected':'idle'}-${deviceScaleFactor}x.png`});
   }
  }
  await sample.close();
 }
 for(const name of ['home','templates','settings']){
  const icon=page.locator('.navbtn[data-page="'+name+'"] svg');
  await icon.evaluate(n=>{n.style.width='96px';n.style.height='96px'});
  await icon.screenshot({path:`visual-checks/icon-${name}.png`});
 }

 await page.evaluate(()=>document.querySelectorAll(".navbtn svg").forEach(n=>{n.style.width="";n.style.height="";}));
 await page.setViewportSize({width:1440,height:1100});
 await page.addScriptTag({content:['core/headless-builder.js','core/creator.js','core/creator-ui.js'].map(file=>fs.readFileSync(file,'utf8')).join('\n')+'\nwindow.CreatorUI=PlazCodeCreatorUI;'});
 await page.evaluate(async()=>{window.CreatorWorkspace=CreatorUI.mount(document.getElementById('creatorWorkspace'),async()=>({ok:true,records:[],feedback:[]}),async()=>({ok:true}),async()=>({ok:true}),async()=>({ok:true}),()=>({canTaskSend:false}));await CreatorWorkspace.refresh();});
 for(const width of [1440,1024,760]) {
  await page.setViewportSize({width,height:1100});
  for(const theme of ['default','ocean','orchid']) {
   await page.evaluate(theme=>{Applyappearance({theme});document.querySelectorAll('.page').forEach(n=>n.classList.toggle('active',n.id==='page-creators'));document.getElementById('content').scrollTop=0;},theme);
   await page.screenshot({animations:'disabled',path:`visual-checks/guided-creator-${theme}-${width}.png`});
   const overlap=await page.evaluate(()=>{const nodes=[...document.querySelectorAll('.pc-create-editor button')];return nodes.some((n,i)=>i>0 && n.getBoundingClientRect().top<nodes[i-1].getBoundingClientRect().bottom+5) && innerWidth>1000;});if(overlap)throw Error('Creator buttons overlap or lack spacing');
   const overflow=await page.evaluate(()=>{const n=document.getElementById('content');return n.scrollWidth>n.clientWidth+1});if(overflow)throw Error('Guided creator horizontal overflow');
   await page.evaluate(()=>{document.querySelectorAll('.page').forEach(n=>n.classList.toggle('active',n.id==='page-settings'));document.querySelector('.settings-guide').open=true;document.getElementById('content').scrollTop=0;});
   await page.screenshot({animations:'disabled',path:`visual-checks/settings-guide-${theme}-${width}.png`});
   const control=await page.locator('#savedChatsSelect').evaluate(n=>({background:getComputedStyle(n).backgroundColor,scheme:getComputedStyle(n).colorScheme}));if(control.background==='rgb(255, 255, 255)' || control.scheme!=='dark')throw Error('Unstyled saved-chat selector');
   await page.locator('[data-pref="rsToolBudget"]').scrollIntoViewIfNeeded();await page.screenshot({animations:'disabled',path:`visual-checks/execution-help-${theme}-${width}.png`});
  }
 }
 const cards=await browser.newPage({viewport:{width:720,height:580}});await cards.setContent('<html><body style="margin:24px;background:#11151c;color:#e7edf6"><div id="rs-menu"><div class="rs-prompt-field"><label>Commands per run</label><input type="number" value="0"><select><option>Saved chat</option></select><button>Save continuation file</button><button>Resume paused task</button></div></div><p>This explanation should remain visible after a tool result.</p><div class="rs-chip"><div class="rs-chip-head"><span class="rs-chip-ic">✓</span><span class="rs-chip-tx">execute_luau</span><span class="rs-chip-dt">complete</span></div></div><div class="rs-chip"><div class="rs-chip-head"><span class="rs-chip-ic">✓</span><span class="rs-chip-tx">execute_luau · result</span></div></div></body></html>');await cards.addStyleTag({content:fs.readFileSync('overlay.css','utf8')});await cards.addStyleTag({content:'#rs-menu{position:static!important;opacity:1!important;transform:none!important;filter:none!important;visibility:visible!important;inset:auto!important;width:100%!important;max-width:none!important;margin:0 0 24px!important}'});await cards.screenshot({path:'visual-checks/extension-controls-command-spacing.png'});const geometry=await cards.locator('.rs-chip').evaluateAll(nodes=>nodes.map(n=>n.getBoundingClientRect().toJSON()));if(geometry[1].top<geometry[0].bottom+8)throw Error('Command cards overlap');await cards.close();
 await page.emulateMedia({reducedMotion:'reduce'});
 const motion=await page.locator('.page.active').evaluate(n=>getComputedStyle(n).animationName);if(motion!=='none')throw Error('Reduced motion was ignored');
 await browser.close();console.log('Visual checks: three widths, four pages, five new palettes, reduced motion passed.');
})().catch(e=>{console.error(e);process.exitCode=1});
