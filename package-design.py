from pathlib import Path
import json,zipfile,hashlib,io,plistlib,struct
repo=Path.cwd();version='1.19.17';root=repo/'build/PlazCode'
entry={'version': '1.19.17', 'title': 'Visible desktop revamp and persistent automatic updates', 'summary': 'A more visible upgrade to the existing desktop workspace, paired with automatic installation for updates detected after launch and accurate desktop build detection.', 'added': ['Larger stacked AI launch cards, a stronger session panel, four separate status cards, and framed Quick Access and Recent Activity panels.', 'A clearer navigation rail, larger headings, more readable card text, improved input controls, and consistent spacing throughout the desktop.', 'Separate running desktop and extension version reporting, with native build requirements in both platform release feeds.'], 'improved': ['Smooth CSS hover, press, and navigation transitions retain existing controls and workflows. Reduced-motion preferences remain respected.', 'All twelve themes remain available, including Oceanic, Copper Atelier, Aurora, Orchid Noir, and Solar Dusk.', 'Responsive layouts adapt the upgraded workspace to wide, medium, and compact windows.', 'Automatic installation stays enabled after an up-to-date launch check, so releases detected while the app remains open can install.'], 'fixed': ['Verified same-version Windows packages can repair an older executable without bypassing checksum/version verification or permitting package downgrades.', 'Extension manifest updates no longer hide an older running desktop executable or make its older interface appear current.', 'Background release checks are spaced thirty seconds apart, replacing the two-second network loop. Transient check failures retry; installer failures remain visible and avoid repeated install loops.'], 'notes': ['Rebuilt Windows and universal Intel/Apple Silicon macOS desktop executables; separate normal and macOS ZIP downloads.', 'The desktop design retains existing navigation routes and control handlers, apart from the explicit version-reporting correction. AI, Studio, and creator behavior are unchanged in this update.', 'Engram memory and large-template improvements remain separate work and are not included in this release.', 'Authenticated native tests use the real published release download and checksum with an installer test helper; full GUI installation and relaunch are not exercised.']}
notes=json.loads((repo/'release-notes.json').read_text());notes=[entry]+[x for x in notes if x['version']!=version]
intro='PlazCode 1.19.17 — Visible desktop revamp and persistent automatic updates\n\n'+entry['summary']+'\n\n'+'\n'.join('- '+x for key in ['added','improved','fixed'] for x in entry[key])+'\n\n'
validation='PlazCode 1.19.17 validation\nApplicable JavaScript regressions, preserved navigation/control parity, theme propagation, and running-native version rendering are checked. Hosted Chromium verifies Home, Settings, Tools and Creators at 1440, 1024 and 760 pixels, five newer themes, and reduced motion. Visual changes use event-driven CSS without new animation libraries, timers, or observers. No end-to-end performance benchmark is inferred.\nWindows and macOS native tests and authenticated headless smoke checks must pass before publishing. Automatic-update checks cover outdated launch installation and an initially up-to-date app remaining open before its installation becomes outdated. Both use the real published feed/download with SHA256 verification and an isolated installer test helper. Full GUI installation/relaunch is not exercised. Unit tests cover extension metadata masking an older native build and extension-only releases that do not require rebuilding the desktop.\n\n'
assert (repo/'artifacts/native-Windows/PlazCode.exe').is_file()
appzip=repo/'artifacts/native-macOS/PlazCode-app.zip';assert appzip.is_file()
with zipfile.ZipFile(appzip) as app:
 plist=plistlib.loads(app.read('PlazCode.app/Contents/Info.plist'));assert plist['CFBundleShortVersionString']==version
 binary=app.read('PlazCode.app/Contents/MacOS/PlazCode');assert binary[:4]==bytes.fromhex('cafebabe')
 count=struct.unpack('>I',binary[4:8])[0];assert {0x01000007,0x0100000c}.issubset({struct.unpack('>I',binary[8+i*20:12+i*20])[0] for i in range(count)})
 appentries={f'PlazCode/{i.filename}':(i,app.read(i)) for i in app.infolist() if not i.is_dir()}
for name in ['README.md','UPDATE.txt','MAINTENANCE.md']:
 (repo/name).write_text(intro+(repo/name).read_text());(root/name).write_bytes((repo/name).read_bytes())
(repo/'Update-PlazCode.ps1').write_bytes((root/'Update-PlazCode.ps1').read_bytes())
(repo/'VALIDATION.txt').write_text(validation+(repo/'VALIDATION.txt').read_text());(root/'VALIDATION.txt').write_bytes((repo/'VALIDATION.txt').read_bytes())
(repo/'release-notes.json').write_text(json.dumps(notes,indent=2)+'\n');(root/'release-notes.json').write_bytes((repo/'release-notes.json').read_bytes())
for p in [root/'macos/Info.plist']:
 data=plistlib.loads(p.read_bytes());data['CFBundleShortVersionString']=data['CFBundleVersion']=version;p.write_bytes(plistlib.dumps(data))
source={f'PlazCode/{p.relative_to(root)}':p.read_bytes() for p in root.rglob('*') if p.is_file() and not any(x in p.relative_to(root).parts for x in ['target','.git','node_modules','__pycache__','PlazCode.app'])}
source['PlazCode/PlazCode.exe']=(repo/'artifacts/native-Windows/PlazCode.exe').read_bytes()
original=repo/'PlazCode-1.19.16.zip';metadata=[]
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
feeds={x['platform']:{'version':version,'desktop_version':version,'url':f"https://raw.githubusercontent.com/stoveez/PlazCodeneww/main/{x['file']}",'sha256':x['sha256'],'release_notes':notes} for x in metadata}
feeds['windows']['platforms']={'macos':{k:feeds['macos'][k] for k in ['url','sha256']}}
(repo/'latest.json').write_text(json.dumps(feeds['windows'],indent=2)+'\n');(repo/'latest-macos.json').write_text(json.dumps(feeds['macos'],indent=2)+'\n')
(repo/'SHA256SUMS.txt').write_text(''.join(f"{x['sha256']}  {x['file']}\n" for x in metadata));(repo/'release-metadata.json').write_text(json.dumps(metadata,indent=2)+'\n')
(repo/'release-description-1.19.17.txt').write_text(intro+'\nValidation\n'+validation)
print(json.dumps(metadata,indent=2))
