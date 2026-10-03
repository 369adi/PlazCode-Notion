# PlazCode updater, Windows PowerShell 5.1. Accepts an optional local release ZIP.
param([string]$ZipPath, [string]$ExpectedSha256, [string]$ExpectedVersion, [switch]$ShowProgress)
$ErrorActionPreference = 'Stop'
$install = $PSScriptRoot
$stage = Join-Path ([IO.Path]::GetTempPath()) ('PlazCode-update-' + [guid]::NewGuid())
$backup = Join-Path $stage 'backup'
$written = New-Object 'System.Collections.Generic.List[string]'
$stopped = $false
$uiReady = $false
function Set-UpdaterProgress([int]$Percent, [string]$Title, [string]$Detail) {
    if ($uiReady) { [PlazCode.UpdateProgress]::Set($Percent, $Title, $Detail) }
}
function Get-ExtensionRoot([string]$Root) {
    $nested = Join-Path $Root 'PlazCode-Extension'
    if (Test-Path -LiteralPath (Join-Path $nested 'manifest.json')) { return $nested }
    return $Root
}
$extensionRoot = Get-ExtensionRoot $install
$splitInstall = $extensionRoot -ne $install -and !(Test-Path -LiteralPath (Join-Path $install 'manifest.json'))
try {
    if ($ShowProgress) {
        try {
            Add-Type -AssemblyName System.Windows.Forms
            Add-Type -AssemblyName System.Drawing
            Add-Type -Path (Join-Path $install 'Updater-Progress.cs') -ReferencedAssemblies ([Windows.Forms.Form].Assembly.Location),([Drawing.Color].Assembly.Location),'System.dll'
            [PlazCode.UpdateProgress]::Open($ExpectedVersion)
            $uiReady = $true
        } catch { Write-Host ('Progress window unavailable: ' + $_.Exception.Message) -ForegroundColor Yellow }
    }
    Set-UpdaterProgress -1 'Checking release' 'Reading the verified update information.'
    Write-Host 'PlazCode updater' -ForegroundColor Yellow
    $current = [version](Get-Content (Join-Path $extensionRoot 'manifest.json') -Raw | ConvertFrom-Json).version
    $source = Get-Content (Join-Path $install 'update-source.json') -Raw | ConvertFrom-Json
    if (!$source.feedUrl) { $source.feedUrl = "https://raw.githubusercontent.com/stoveez/PlazCodeneww/main/latest.json" }
    New-Item $stage -ItemType Directory | Out-Null
    if (!$ZipPath -and $source.feedUrl) {
        if ($source.feedUrl -notmatch '^https://') { throw 'The release feed must use HTTPS.' }
        [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
        $feedUri = [UriBuilder]$source.feedUrl
        $query = $feedUri.Query.TrimStart('?')
        $feedUri.Query = ($query + $(if ($query) { '&' } else { '' }) + 'plazcode_check=' + [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())
        $release = Invoke-RestMethod -Uri $feedUri.Uri.AbsoluteUri -Headers @{ 'Cache-Control' = 'no-cache, max-age=0' } -TimeoutSec 30
        $latest = [version]$release.version
        if ($latest -le $current) { Write-Host "Already up to date ($current)."; exit 0 }
        if ($release.url -notmatch '^https://' -or $release.sha256 -notmatch '^[a-fA-F0-9]{64}$') { throw 'Invalid release feed: expected version, HTTPS url and SHA256.' }
        $ZipPath = Join-Path $stage 'release.zip'
        Set-UpdaterProgress -1 'Downloading update' 'Downloading the release package. This step depends on your connection.'
        Write-Host "Downloading $latest..."
        Invoke-WebRequest -UseBasicParsing -Uri $release.url -OutFile $ZipPath -TimeoutSec 180
        if ((Get-FileHash $ZipPath -Algorithm SHA256).Hash -ne $release.sha256) { throw 'Download checksum mismatch. No installed files were changed.' }
    }
    if (!$ZipPath) {
        Write-Host 'No release feed configured. Select a downloaded PlazCode release ZIP.'
        Add-Type -AssemblyName System.Windows.Forms
        $picker = New-Object System.Windows.Forms.OpenFileDialog
        $picker.Filter = 'PlazCode release (*.zip)|*.zip'
        if ($picker.ShowDialog() -ne 'OK') { Write-Host 'Update cancelled.'; exit 0 }
        $ZipPath = $picker.FileName
    }
    $ZipPath = (Resolve-Path -LiteralPath $ZipPath).Path
    Set-UpdaterProgress 15 'Verifying package' 'Checking the download before any installed files are changed.'
    if ($ExpectedSha256) {
        if ($ExpectedSha256 -notmatch '^[a-fA-F0-9]{64}$' -or (Get-FileHash -LiteralPath $ZipPath -Algorithm SHA256).Hash -ne $ExpectedSha256) { throw 'Package checksum mismatch. Installed files were not changed.' }
    }
    Set-UpdaterProgress 25 'Inspecting package' 'Validating the release contents and installation layout.'
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $unpacked = Join-Path $stage 'unpacked'
    $archive = [IO.Compression.ZipFile]::OpenRead($ZipPath)
    $packageRoot = $null
    try {
        foreach ($entry in $archive.Entries) {
            $name = $entry.FullName.Replace('\','/')
            if ($name -notmatch '^(PlazCode|PlazCode-Extension)/' -or $name -match '(^|/)\.\.(/|$)|:|(^|/)(logs|target|\.git)(/|$)') { throw "Unexpected ZIP entry: $name" }
            $root = $name.Split('/')[0]
            if (!$packageRoot) { $packageRoot = $root } elseif ($root -ne $packageRoot) { throw 'Release ZIP contains mixed installation folders.' }
            $destination = [IO.Path]::GetFullPath((Join-Path $unpacked $name))
            if (!$destination.StartsWith([IO.Path]::GetFullPath($unpacked) + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) { throw 'Invalid ZIP path.' }
        }
    } finally { $archive.Dispose() }
    Set-UpdaterProgress 35 'Unpacking update' 'Extracting the verified files to a temporary folder.'
    [IO.Compression.ZipFile]::ExtractToDirectory($ZipPath, $unpacked)
    if (!$packageRoot) { throw 'Release ZIP is empty.' }
    $package = Join-Path $unpacked $packageRoot
    $packageExtension = Get-ExtensionRoot $package
    $next = [version](Get-Content (Join-Path $packageExtension 'manifest.json') -Raw | ConvertFrom-Json).version
    if ($ExpectedVersion -and $next -ne [version]$ExpectedVersion) { throw 'Release version does not match the verified feed.' }
    # A verified native download may repair an older executable at the same extension version.
    if ($next -lt $current -or ($next -eq $current -and !($ExpectedVersion -and $ExpectedSha256))) { Write-Host "Already up to date ($current; selected package $next)."; exit 0 }
    foreach ($required in @('PlazCode.exe','WebView2Loader.dll','Start-PlazCode-Agent.cmd')) {
        if (!(Test-Path -LiteralPath (Join-Path $package $required))) { throw "Release is missing $required." }
    }
    foreach ($required in @('background.js','core/main.js')) {
        if (!(Test-Path -LiteralPath (Join-Path $packageExtension $required))) { throw "Extension is missing $required." }
    }
    if ($release -and $next -ne $latest) { throw 'Release feed version does not match the downloaded package.' }
    $preservedNames = @('config.json','config.local.json','plazcode-settings.json','bridge-pairing.json','memory.json','chat-history.json','checkpoints.json','catalog.json','creations.json','update-source.json')
    $preservedFolders = '^(?:logs|backups|templates|runtimes|PlazCode\.exe\.WebView2)(?:[\\/]|$)'
    # Compatibility ZIPs also carry legacy root extension copies for old installers.
    # Do not add those duplicates to installations already using the split layout.
    $legacyExtensionFiles = '^(?:manifest\.json|background\.js|popup\.(?:html|js)|overlay\.css|icon\.png|ollama\.html|ollama-page\.js|(?:core|providers|ui)[\\/].*)$'
    $files = @(Get-ChildItem $package -File -Recurse | Where-Object {
        $relative = $_.FullName.Substring($package.Length + 1)
        $_.Name -notin $preservedNames -and $relative -notmatch $preservedFolders -and !($splitInstall -and $packageExtension -ne $package -and $relative -match $legacyExtensionFiles)
    })
    Set-UpdaterProgress 45 'Saving recovery copies' 'Backing up installed files. Your preferences and templates are preserved.'
    $backedUp = 0
    # Back up overwritten files before stopping the app or modifying the installation.
    foreach ($file in $files) {
        $relative = $file.FullName.Substring($package.Length + 1)
        $old = Join-Path $install $relative
        if (Test-Path -LiteralPath $old) {
            $copy = Join-Path $backup $relative
            New-Item (Split-Path $copy) -ItemType Directory -Force | Out-Null
            Copy-Item -LiteralPath $old -Destination $copy
        }
        $backedUp++
        Set-UpdaterProgress (45 + [int](15 * $backedUp / [Math]::Max(1, $files.Count))) 'Saving recovery copies' ("Preparing files: $backedUp of $($files.Count).")
    }
    Set-UpdaterProgress 60 'Closing PlazCode' 'Waiting for this installation to release its files. Roblox Studio stays open.'
    # Stop only agent processes belonging to this installation. Leave Studio running.
    $agentPaths = @('PlazCode.exe','plazcode-agent.exe') | ForEach-Object { [IO.Path]::GetFullPath((Join-Path $install $_)) }
    $deadline = [DateTime]::UtcNow.AddSeconds(20)
    do {
        # Process.Path works for our same-user app even if CIM omits ExecutablePath.
        $agents = @(Get-Process -Name 'PlazCode','plazcode-agent' -ErrorAction SilentlyContinue | Where-Object {
            try { $_.Path -and $agentPaths -contains [IO.Path]::GetFullPath($_.Path) } catch { $false }
        })
        foreach ($agent in $agents) {
            $stopped = $true
            try { Stop-Process -Id $agent.Id -Force -ErrorAction Stop }
            catch { if (Get-Process -Id $agent.Id -ErrorAction SilentlyContinue) { throw } }
            # Stop-Process can return before Windows releases the executable image.
            Wait-Process -Id $agent.Id -Timeout 5 -ErrorAction SilentlyContinue
        }
        $locked = $false
        foreach ($path in $agentPaths) {
            if (!(Test-Path -LiteralPath $path)) { continue }
            $handle = $null
            try { $handle = [IO.File]::Open($path, [IO.FileMode]::Open, [IO.FileAccess]::ReadWrite, [IO.FileShare]::None) }
            catch [IO.IOException] { $locked = $true }
            finally { if ($handle) { $handle.Dispose() } }
        }
        if (!$locked) { break }
        if ([DateTime]::UtcNow -ge $deadline) { throw 'PlazCode.exe is still locked. Close other PlazCode windows and retry. No update files were copied.' }
        if ($ShowProgress) { Set-UpdaterProgress 60 'Waiting for PlazCode to close' 'Finishing process shutdown before replacing files.' }
        Start-Sleep -Milliseconds 250
    } while ($true)
    $copied = 0
    Set-UpdaterProgress 65 'Installing update' 'Replacing application and extension files.'
    foreach ($file in $files) {
        $relative = $file.FullName.Substring($package.Length + 1)
        $target = Join-Path $install $relative
        $written.Add($relative)
        New-Item (Split-Path $target) -ItemType Directory -Force | Out-Null
        Copy-Item -LiteralPath $file.FullName -Destination $target -Force
        $copied++
        Set-UpdaterProgress (65 + [int](30 * $copied / [Math]::Max(1, $files.Count))) 'Installing update' ("Replacing files: $copied of $($files.Count).")
    }
    Set-UpdaterProgress 98 'Relaunching PlazCode' 'Opening the updated desktop app.'
    Start-Process -FilePath (Join-Path $install 'PlazCode.exe') -WorkingDirectory $install
    Set-UpdaterProgress 100 'Update complete' 'Reload the extension and refresh your open AI chat tabs.'
    Write-Host "Updated $current -> $next. Reload the extension, then refresh open AI chat tabs." -ForegroundColor Green
} catch {
    Set-UpdaterProgress -1 'Restoring installation' 'The update failed. Restoring recovery copies before reporting the error.'
    foreach ($relative in $written) {
        $target = Join-Path $install $relative
        $saved = Join-Path $backup $relative
        try {
            if (Test-Path -LiteralPath $saved) { Copy-Item -LiteralPath $saved -Destination $target -Force }
            else { Remove-Item -LiteralPath $target -Force -ErrorAction SilentlyContinue }
        } catch { Write-Host "Could not restore $relative. Backup retained at $backup" -ForegroundColor Red }
    }
    if ($stopped) { Start-Process -FilePath (Join-Path $install 'PlazCode.exe') -WorkingDirectory $install -ErrorAction SilentlyContinue }
    if ($ShowProgress) { [Windows.Forms.MessageBox]::Show(('Update failed: ' + $_.Exception.Message + [Environment]::NewLine + 'Recovery files: ' + $stage), 'PlazCode updater') | Out-Null }
    Write-Host ('Update failed: ' + $_.Exception.Message) -ForegroundColor Red
    Write-Host "Recovery files: $stage"
    if ($uiReady) { [PlazCode.UpdateProgress]::Close() }
    exit 1
}
if ($uiReady) { [PlazCode.UpdateProgress]::Close() }
# Keep backup for recovery; downloaded/extracted staging files can be removed.
Remove-Item (Join-Path $stage 'unpacked') -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item (Join-Path $stage 'release.zip') -Force -ErrorAction SilentlyContinue
