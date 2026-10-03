"""Create a versioned GitHub Release from an already published PlazCode feed.
No secrets or local install data are packaged. Existing published releases are
never edited; an interrupted draft can be completed by rerunning the workflow.
"""
import base64, hashlib, io, json, os, pathlib, re, subprocess, urllib.request, zipfile

def gh(*args):
    return subprocess.check_output(['gh',*args],text=True)
def fetch(url):
    with urllib.request.urlopen(url,timeout=60) as r:
        data=r.read(32*1024*1024+1)
    if len(data)>32*1024*1024: raise ValueError('Release download exceeds 32 MB')
    return data

def notes(feed):
    version=feed['version'];entry=next(x for x in feed['release_notes'] if x['version']==version)
    lines=['# PlazCode '+version+': '+entry['title'],'',entry['summary'],'']
    for key,title in [('added','Added'),('improved','Improved'),('fixed','Fixed')]:
        if entry.get(key):lines+=['## '+title,'']+['- '+x for x in entry[key]]+['']
    lines+=['## Update','','- In the desktop app: **Updates → Update now**, or run **Update-PlazCode.bat**.','- Reload the extension in **chrome://extensions** and refresh open AI tabs.','- Download **PlazCode-'+version+'.zip** for either a fresh installation or updating an existing one.','- Existing settings, memory and enabled MCP servers keep their data locations.','','See VALIDATION.txt inside the ZIP for checks and live-test limitations.']
    return 'PlazCode '+version+': '+entry['title'],'\n'.join(lines)+'\n'

def validate_archive(data,version):
    with zipfile.ZipFile(io.BytesIO(data)) as archive:
        if archive.testzip() is not None:raise ValueError('ZIP checksum failure')
        manifest=json.loads(archive.read('PlazCode/PlazCode-Extension/manifest.json'))
        if tuple(map(int,manifest['version'].split('.')))!=tuple(map(int,version.split('.'))):raise ValueError('ZIP version does not match feed')
        if 'PlazCode/PlazCode.exe' not in archive.namelist():raise ValueError('Desktop executable missing')

def main():
    repo=os.environ['GITHUB_REPOSITORY'];commit=os.environ['GITHUB_SHA']
    if not re.fullmatch(r'[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+',repo) or not re.fullmatch(r'[a-f0-9]{40}',commit):raise ValueError('Invalid repository or commit')
    root='https://raw.githubusercontent.com/'+repo+'/'+commit+'/'
    feed=json.loads(fetch(root+'latest.json'));version=feed['version']
    if not re.fullmatch(r'\d+\.\d+\.\d+',version):raise ValueError('Invalid release version')
    tag='v'+version;title,body=notes(feed)
    try:existing=json.loads(gh('release','view',tag,'--repo',repo,'--json','isDraft'))
    except subprocess.CalledProcessError:existing=None
    if existing and not existing['isDraft']:
        print('Version is already published; leaving it unchanged.');return
    assets=[]
    for name in ['PlazCode-'+version+'.zip']:
        data=fetch(root+name)
        content=json.loads(gh('api','repos/'+repo+'/contents/'+name+'?ref='+commit))
        actual=hashlib.sha1(b'blob '+str(len(data)).encode()+b'\0'+data).hexdigest()
        if actual!=content['sha']:raise ValueError('Repository ZIP hash mismatch')
        if hashlib.sha256(data).hexdigest()!=feed['sha256']:raise ValueError('Update SHA-256 mismatch')
        validate_archive(data,version);pathlib.Path(name).write_bytes(data);assets.append(name)
    pathlib.Path('release-notes.md').write_text(body)
    if not existing:gh('release','create',tag,'--repo',repo,'--target',commit,'--draft','--title',title,'--notes-file','release-notes.md')
    gh('release','upload',tag,*assets,'--repo',repo,'--clobber')
    current=json.loads(gh('api','repos/'+repo+'/contents/latest.json?ref=main'))
    current=json.loads(base64.b64decode(current['content']))
    gh('release','edit',tag,'--repo',repo,'--draft=false','--latest='+str(current['version']==version).lower())
    print('Published https://github.com/'+repo+'/releases/tag/'+tag)
if __name__=='__main__':main()
