"""Run the real updater's preflight without reaching process termination or install."""
from pathlib import Path
import hashlib, shutil, stat, subprocess, tempfile, zipfile

source = Path(__file__).with_name('Update-PlazCode.command')
text = source.read_text()
start = text.index('version_order(){')
order = text[start:text.index('\n}\n', start) + 3]
for left, right, expected in [('1.19.13', '1.19.14', '-1'), ('1.19.09', '1.19.9', '0'),
                              ('1.19.100', '1.20.0', '-1'), ('1.20.0', '1.19.99', '1')]:
    result = subprocess.run(['/bin/bash', '-c', order + '\nversion_order "$1" "$2"',
                             'version-fixture', left, right], capture_output=True, text=True, check=True)
    assert result.stdout.strip() == expected
with tempfile.TemporaryDirectory(prefix='plazcode-updater-preflight-') as temporary:
    root = Path(temporary)
    script = root / source.name
    shutil.copy2(source, script)
    preserved = root / 'memory.json'
    preserved.write_bytes(b'unchanged user data')
    cases = ['checksum', 'traversal', 'absolute', 'backslash', 'wrong-root', 'symlink']
    for case in cases:
        archive = root / (case + '.zip')
        name = {'traversal': 'PlazCode/../outside.txt', 'absolute': '/outside.txt',
                'backslash': 'PlazCode/..\\outside.txt', 'wrong-root': 'Other/app',
                'symlink': 'PlazCode/link'}.get(case, 'PlazCode/file.txt')
        with zipfile.ZipFile(archive, 'w') as output:
            entry = zipfile.ZipInfo(name)
            if case == 'symlink':
                entry.create_system = 3
                entry.external_attr = (stat.S_IFLNK | 0o777) << 16
            output.writestr(entry, '../outside.txt' if case == 'symlink' else 'fixture')
        checksum = hashlib.sha256(archive.read_bytes()).hexdigest()
        if case == 'checksum':
            checksum = '0' * 64
        result = subprocess.run(['/bin/bash', str(script), str(archive), checksum,
                                 '1.19.14', '99999999'], capture_output=True, text=True, timeout=10)
        expected = 3 if case == 'checksum' else 4 if case == 'symlink' else 1
        assert result.returncode == expected, (case, result.returncode, result.stdout, result.stderr)
        assert preserved.read_bytes() == b'unchanged user data'
        assert not (root / 'outside.txt').exists()
        assert not (root / 'PlazCode.app').exists()
    print('PASS updater checksum, traversal, absolute/backslash paths, package root and symlink rejection; user data unchanged')
