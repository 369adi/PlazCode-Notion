#!/bin/bash
set -euo pipefail
root=$(cd -- "$(dirname -- "$0")" && pwd)
zip=${1:?Update ZIP path required}
expectedhash=${2:?Checksum required}
version=${3:?Version required}
pid=${4:?Running process required}
[[ "$expectedhash" =~ ^[a-fA-F0-9]{64}$ && "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ && "$pid" =~ ^[0-9]+$ ]] || exit 2
[[ "$(/usr/bin/shasum -a 256 "$zip" | /usr/bin/awk '{print $1}')" == "$expectedhash" ]] || { echo 'Update checksum mismatch. Installation unchanged.' >&2; exit 3; }
# Reject traversal, absolute paths, symlinks and unexpected package roots before extraction.
/usr/bin/unzip -Z -1 "$zip" | /usr/bin/awk 'BEGIN{bad=0} /(^\/|\\|(^|\/)\.\.(\/|$)|:)/{bad=1} !/^PlazCode\//{bad=1} END{exit bad}'
if /usr/bin/unzip -Z -l "$zip" | /usr/bin/grep -q '^l'; then echo 'Symlinks are not accepted in updates.' >&2; exit 4; fi
stage=$(/usr/bin/mktemp -d "${TMPDIR:-/tmp}/PlazCode-update.XXXXXX")
backup="$stage/backup";mkdir -p "$backup"
/usr/bin/ditto -x -k "$zip" "$stage/extracted"
source="$stage/extracted/PlazCode"
[[ "$(/usr/bin/plutil -extract version raw -o - "$source/PlazCode-Extension/manifest.json")" == "$version" ]]
[[ "$(/usr/bin/plutil -extract CFBundleShortVersionString raw -o - "$source/PlazCode.app/Contents/Info.plist")" == "$version" ]]
/usr/bin/file "$source/PlazCode.app/Contents/MacOS/PlazCode" | /usr/bin/grep -q 'Mach-O'
/usr/bin/codesign --verify --deep --strict "$source/PlazCode.app"
[[ "$(/bin/ps -ww -p "$pid" -o comm=)" == "$root/PlazCode.app/Contents/MacOS/PlazCode" ]] || { echo 'Running installation changed. Retry from the app.' >&2; exit 5; }
/bin/kill -TERM "$pid"
for ((attempt=0;attempt<100;attempt++));do if ! /bin/kill -0 "$pid" 2>/dev/null;then break;fi;/bin/sleep .1;done
if /bin/kill -0 "$pid" 2>/dev/null;then echo 'PlazCode did not exit. Installation unchanged.' >&2;exit 6;fi
rollback(){
 for name in PlazCode.app PlazCode-Extension;do
  if [[ -e "$backup/$name" ]];then
   if [[ -e "$root/$name" ]];then /bin/mv "$root/$name" "$stage/failed-$name";fi
   /bin/mv "$backup/$name" "$root/$name"
  fi
 done
 /usr/bin/open "$root/PlazCode.app" || true
 echo "Update failed; restored installation. Diagnostic files: $stage" >&2
}
trap rollback ERR
for name in PlazCode.app PlazCode-Extension;do
 [[ ! -e "$root/$name" ]] || /bin/mv "$root/$name" "$backup/$name"
 /bin/mv "$source/$name" "$root/$name"
done
# Only owned launchers/docs are refreshed; settings, MCP configuration and user data are preserved.
for name in MacOS_Setup.command Start-PlazCode.command Update-PlazCode.command UPDATE.txt README.md release-notes.json;do
 if [[ -f "$source/$name" ]];then /bin/cp "$source/$name" "$root/$name";fi
done
/bin/chmod 755 "$root/PlazCode.app/Contents/MacOS/PlazCode" "$root/Start-PlazCode.command" "$root/MacOS_Setup.command" "$root/Update-PlazCode.command"
/usr/bin/open "$root/PlazCode.app"
trap - ERR
echo 'Update installed. Reload the extension and refresh open AI tabs.'
