from pathlib import Path
import json,zipfile,hashlib,io,plistlib,struct
repo=Path.cwd();version='1.19.16';root=repo/'build/PlazCode'
entry={'version':version,'title':'Workspace design upgrade and new themes','summary':'A visual refinement built on the existing workspace layout, with smoother interaction tweens, clearer controls and five new shared themes. Startup update detection also receives dedicated macOS feed freshness protection.','added':['Five shared themes: Oceanic, Copper Atelier, Aurora, Orchid Noir and Solar Dusk. Each colors the complete desktop surface system, browser bar and extension popup.','Named theme preview tiles with visible selection and keyboard focus states.'],'improved':['Workspace card spacing, rounded surfaces, session hierarchy, navigation feedback and form styling build on the existing desktop design.','Short CSS hover, press and page-entry tweens add interaction feedback without new animation libraries, timers or DOM observers.','Reduced-motion settings disable desktop transitions and animations, including the update spinner. Existing glow and gradient controls still apply.','Normal app launch and reopening an existing app instance request automatic installation when a newer published version is detected.'],'fixed':['New palettes stay synchronized across desktop preferences, browser overlays and popup appearance.','The dedicated macOS release feed now uses the same commit-pinned freshness lookup and cache-busting fallback as the normal feed. Platform feed cache entries cannot reuse another platform payload.'],'notes':['Desktop native executables rebuilt for Windows and universal Intel/Apple Silicon macOS.','All existing UI IDs, navigation routes and desktop JavaScript handlers remain unchanged outside the theme definitions. Design changes do not alter AI/provider/tool behavior.','The normal ZIP retains macOS compatibility for older shared-feed updaters; a separate macOS ZIP is also supplied.','No live signed-in AI chat or real Studio workflow was exercised for this design update. Automated launch-update verification uses the real published download and checksum with an isolated installer test helper; full GUI installation/relaunch is not exercised.']}
notes=json.loads((repo/'release-notes.json').read_text());notes=[entry]+[x for x in notes if x['version']!=version]
intro='PlazCode 1.19.16 — Workspace design upgrade and new themes\n\n'+entry['summary']+'\n\n'+'\n'.join('- '+x for key in ['added','improved','fixed'] for x in entry[key])+'\n\n'
validation='''PlazCode 1.19.16 validation
56 applicable JavaScript regressions pass, including desktop behavior, supported providers, creator workflows, tool handling, Stop, Co-work and continuation. Twelve shared palettes are accepted by desktop, native preferences and browser overlays; popup propagation covers every palette. Desktop script parity confirms existing handlers and navigation routes are preserved outside theme definitions.
Hosted Chromium layout checks cover Home, Settings, Tools and Creators at 1440, 1024 and 760 pixels, all five new themes, and reduced motion. Screenshot review corrected theme-control alignment and number-input styling. New motion uses short event-driven CSS only; no new timers, observers, animation libraries or continuous decorative loops were added. No end-to-end performance benchmark or guarantee is inferred from these checks.
Windows and macOS native tests and authenticated headless smoke checks must pass before this package is published. Universal macOS binaries are combined with lipo and the app bundle is ad hoc signed and verified. Launch-update tests use a temporary outdated installation, real published feed/download and SHA256 verification, and record the installer version/checksum handoff. A test helper replaces full GUI installation/relaunch. Mac feed migration and updater archive/checksum rejection fixtures pass. Live signed-in AI chats, real Studio operations and real GUI update/relaunch are not exercised.

'''
assert (repo/'artifacts/native-Windows/PlazCode.exe').is_file()
appzip=repo/'artifacts/native-macOS/PlazCode-app.zip';assert appzip.is_file()
with zipfile.ZipFile(appzip) as app:
 plist=plistlib.loads(app.read('PlazCode.app/Contents/Info.plist'));assert plist['CFBundleShortVersionString']==version
 binary=app.read('PlazCode.app/Contents/MacOS/PlazCode');assert binary[:4]==bytes.fromhex('cafebabe')
 count=struct.unpack('>I',binary[4:8])[0];assert {0x01000007,0x0100000c}.issubset({struct.unpack('>I',binary[8+i*20:12+i*20])[0] for i in range(count)})
 appentries={f'PlazCode/{i.filename}':(i,app.read(i)) for i in app.infolist() if not i.is_dir()}
for name in ['README.md','UPDATE.txt','MAINTENANCE.md']:
 (repo/name).write_text(intro+(repo/name).read_text());(root/name).write_bytes((repo/name).read_bytes())
(repo/'VALIDATION.txt').write_text(validation+(repo/'VALIDATION.txt').read_text());(root/'VALIDATION.txt').write_bytes((repo/'VALIDATION.txt').read_bytes())
(repo/'release-notes.json').write_text(json.dumps(notes,indent=2)+'\n');(root/'release-notes.json').write_bytes((repo/'release-notes.json').read_bytes())
for p in [root/'macos/Info.plist']:
 data=plistlib.loads(p.read_bytes());data['CFBundleShortVersionString']=data['CFBundleVersion']=version;p.write_bytes(plistlib.dumps(data))
source={f'PlazCode/{p.relative_to(root)}':p.read_bytes() for p in root.rglob('*') if p.is_file() and not any(x in p.relative_to(root).parts for x in ['target','.git','node_modules','__pycache__','PlazCode.app'])}
source['PlazCode/PlazCode.exe']=(repo/'artifacts/native-Windows/PlazCode.exe').read_bytes()
original=repo/'PlazCode-1.19.15.zip';metadata=[]
for platform,name in [('windows',f'PlazCode-{version}.zip'),('macos',f'PlazCode-macOS-{version}.zip')]:
 output=repo/name;remaining=dict(source)
 with zipfile.ZipFile(original) as old,zipfile.ZipFile(output,'w',zipfile.ZIP_DEFLATED,compresslevel=9) as new:
  seen=set()
  for info in old.infolist():
   path=info.filename
   if path.startswith('PlazCode/PlazCode.app/'):continue
   if platform=='macos' and Path(path).suffix.lower() in ['.exe','.dll','.bat','.cmd','.ps1']:continue
   data=remaining.pop(path,old.read(info))
   if path=='PlazCode/update-source.json' and platform=='macos':data=json.dumps({'feedUrl':'https://raw.githubusercontent.com/stoveez/PlazCodeneww/main/latest-macos.json'},indent=2).encode()
   new.writestr(info,data);seen.add(path)
  for path,data in remaining.items():
   if path in seen or path.startswith('PlazCode/PlazCode.app/'):continue
   if platform=='macos' and Path(path).suffix.lower() in ['.exe','.dll','.bat','.cmd','.ps1']:continue
   info=zipfile.ZipInfo(path);info.compress_type=zipfile.ZIP_DEFLATED;info.external_attr=(0o100755 if Path(path).suffix=='.command' else 0o100644)<<16;new.writestr(info,data)
  for path,(oldinfo,data) in appentries.items():
   info=zipfile.ZipInfo(path);info.compress_type=zipfile.ZIP_DEFLATED;info.external_attr=oldinfo.external_attr;new.writestr(info,data)
 with zipfile.ZipFile(output) as z:
  assert z.testzip() is None and len(z.namelist())==len(set(z.namelist()))
  assert json.loads(z.read('PlazCode/manifest.json'))['version']==version
  assert json.loads(z.read('PlazCode/PlazCode-Extension/manifest.json'))['version']==version
  assert ('PlazCode/PlazCode.exe' in z.namelist())==(platform=='windows')
  assert z.getinfo('PlazCode/PlazCode.app/Contents/MacOS/PlazCode').external_attr>>16&0o111
  with zipfile.ZipFile(original) as old:
   for path in old.namelist():
    if '/providers/' in path:assert old.read(path)==z.read(path)
 raw=output.read_bytes();assert len(raw)<=32*1024*1024
 metadata.append({'file':name,'platform':platform,'bytes':len(raw),'sha256':hashlib.sha256(raw).hexdigest()})
feeds={x['platform']:{'version':version,'url':f"https://raw.githubusercontent.com/stoveez/PlazCodeneww/main/{x['file']}",'sha256':x['sha256'],'release_notes':notes} for x in metadata}
feeds['windows']['platforms']={'macos':{k:feeds['macos'][k] for k in ['url','sha256']}}
(repo/'latest.json').write_text(json.dumps(feeds['windows'],indent=2)+'\n');(repo/'latest-macos.json').write_text(json.dumps(feeds['macos'],indent=2)+'\n')
(repo/'SHA256SUMS.txt').write_text(''.join(f"{x['sha256']}  {x['file']}\n" for x in metadata));(repo/'release-metadata.json').write_text(json.dumps(metadata,indent=2)+'\n')
(repo/'release-description-1.19.16.txt').write_text(intro+'\nValidation\n'+validation)
print(json.dumps(metadata,indent=2))
