const fs=require('fs');const {chromium}=require('playwright');
(async()=>{
 const original=fs.readFileSync('agent/src/desktop.html','utf8');
 const themeScript='function Byid(id){return document.getElementById(id)};'+original.slice(original.indexOf('  var Desktopthemes ='),original.indexOf('  function Saveappearance()'));
 const html=original.replaceAll('__PLAZCODE_VERSION__','1.19.20').replace(/<script>[\s\S]*?<\/script>/g,'');
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
 await page.emulateMedia({reducedMotion:'reduce'});
 const motion=await page.locator('.page.active').evaluate(n=>getComputedStyle(n).animationName);if(motion!=='none')throw Error('Reduced motion was ignored');
 await browser.close();console.log('Visual checks: three widths, four pages, five new palettes, reduced motion passed.');
})().catch(e=>{console.error(e);process.exitCode=1});
