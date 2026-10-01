# PlazCode updates

This repository hosts the official PlazCode update packages and the automatic update feed.

## Install once

Download `PlazCode-1.18.80.zip`. On an existing 1.18.77 installation, run `Update-PlazCode.bat` and select that ZIP. For a new installation, extract it and run `PlazCode.exe`, then load its browser extension.

## Future updates

In the desktop app, open **Updates**, click **Check now**, then **Update now**. Download progress is shown, installation verifies the package, and PlazCode relaunches.


Run `Update-PlazCode.bat`. It checks `latest.json`, downloads a newer published ZIP, verifies SHA256, installs it, and restarts PlazCode. Reload the extension in `chrome://extensions` and refresh open AI chat tabs afterwards.

Memory, settings and cached optional runtimes are preserved. A custom non-empty release-feed URL takes precedence over the default repository.

## Publish an update

Upload a complete `PlazCode-VERSION.zip` with a `PlazCode/` top-level folder and update `latest.json` in the same commit. Its fields are `version`, `url` (the ZIP's raw HTTPS download URL), and `sha256` (the ZIP's complete SHA256 hash). Use a higher package/extension version and a new filename for every release. Do not modify published ZIP bytes in place.

The BAT can only install releases that have been published here. Source changes or chat attachments alone do not publish an update.

See `UPDATE.txt` and `UPDATER.txt` inside the ZIP for change details and recovery information. Windows updater execution still requires live validation.
