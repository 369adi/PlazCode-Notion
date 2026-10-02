# PlazCode updater, Windows PowerShell 5.1. Accepts an optional local release ZIP.
param([string]$ZipPath, [string]$ExpectedSha256, [string]$ExpectedVersion, [switch]$ShowProgress)
$ErrorActionPreference = 'Stop'
$install = $PSScriptRoot
$stage = Join-Path ([IO.Path]::GetTempPath()) ('PlazCode-update-' + [guid]::NewGuid())
$backup = Join-Path $stage 'backup'
$written = New-Object 'System.Collections.Generic.List[string]'
$stopped = $false
function Get-ExtensionRoot([string]$Root) {
    $nested = Join-Path $Root 'PlazCode-Extension'
    if (Test-Path -LiteralPath (Join-Path $nested 'manifest.json')) { return $nested }
    return $Root
}
$extensionRoot = Get-ExtensionRoot $install
$splitInstall = $extensionRoot -ne $install -and !(Test-Path -LiteralPath (Join-Path $install 'manifest.json'))
try {
    if ($ShowProgress) {
        Add-Type -AssemblyName System.Windows.Forms
        Add-Type -AssemblyName System.Drawing
        $form = New-Object System.Windows.Forms.Form
        $form.Text = 'Updating PlazCode'; $form.Width = 460; $form.Height = 180
        $form.StartPosition = 'CenterScreen'; $form.ControlBox = $false
        $form.BackColor = [Drawing.Color]::FromArgb(9,21,34)
        $label = New-Object System.Windows.Forms.Label
        $label.ForeColor = [Drawing.Color]::FromArgb(255,192,82)
        $label.Location = New-Object Drawing.Point(24,24); $label.Width=400; $label.Height=60
        $label.Text = 'Verifying update...'
        $bar = New-Object System.Windows.Forms.ProgressBar
        $bar.Location = New-Object Drawing.Point(24,90); $bar.Width=400; $bar.Style='Marquee'
        $form.Controls.Add($label); $form.Controls.Add($bar); $form.Show(); [Windows.Forms.Application]::DoEvents()
    }
    Write-Host 'PlazCode updater' -ForegroundColor Yellow
    $current = [version](Get-Content (Join-Path $extensionRoot 'manifest.json') -Raw | ConvertFrom-Json).version
    $source = Get-Content (Join-Path $install 'update-source.json') -Raw | ConvertFrom-Json
    if (!$source.feedUrl) { $source.feedUrl = "https://raw.githubusercontent.com/stoveez/PlazCodeneww/main/latest.json" }
    New-Item $stage -ItemType Directory | Out-Null
    if (!$ZipPath -and $source.feedUrl) {
        if ($source.feedUrl -notmatch '^https://') { throw 'The release feed must use HTTPS.' }
        [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
        $release = Invoke-RestMethod -Uri $source.feedUrl -TimeoutSec 30
        $latest = [version]$release.version
        if ($latest -le $current) { Write-Host "Already up to date ($current)."; exit 0 }
        if ($release.url -notmatch '^https://' -or $release.sha256 -notmatch '^[a-fA-F0-9]{64}$') { throw 'Invalid release feed: expected version, HTTPS url and SHA256.' }
        $ZipPath = Join-Path $stage 'release.zip'
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
    if ($ExpectedSha256) {
        if ($ExpectedSha256 -notmatch '^[a-fA-F0-9]{64}$' -or (Get-FileHash -LiteralPath $ZipPath -Algorithm SHA256).Hash -ne $ExpectedSha256) { throw 'Package checksum mismatch. Installed files were not changed.' }
    }
    if ($ShowProgress) { $label.Text='Preparing verified update...'; [Windows.Forms.Application]::DoEvents() }
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
    [IO.Compression.ZipFile]::ExtractToDirectory($ZipPath, $unpacked)
    if (!$packageRoot) { throw 'Release ZIP is empty.' }
    $package = Join-Path $unpacked $packageRoot
    $packageExtension = Get-ExtensionRoot $package
    $next = [version](Get-Content (Join-Path $packageExtension 'manifest.json') -Raw | ConvertFrom-Json).version
    if ($ExpectedVersion -and $next -ne [version]$ExpectedVersion) { throw 'Release version does not match the verified feed.' }
    if ($next -le $current) { Write-Host "Already up to date ($current; selected package $next)."; exit 0 }
    foreach ($required in @('PlazCode.exe','WebView2Loader.dll','Start-PlazCode-Agent.cmd')) {
        if (!(Test-Path -LiteralPath (Join-Path $package $required))) { throw "Release is missing $required." }
    }
    foreach ($required in @('background.js','core/main.js')) {
        if (!(Test-Path -LiteralPath (Join-Path $packageExtension $required))) { throw "Extension is missing $required." }
    }
    if ($release -and $next -ne $latest) { throw 'Release feed version does not match the downloaded package.' }
    $preservedNames = @('config.json','config.local.json','plazcode-settings.json','bridge-pairing.json','memory.json','chat-history.json','checkpoints.json','catalog.json','update-source.json')
    $preservedFolders = '^(?:logs|backups|templates|runtimes|PlazCode\.exe\.WebView2)(?:[\\/]|$)'
    # Compatibility ZIPs also carry legacy root extension copies for old installers.
    # Do not add those duplicates to installations already using the split layout.
    $legacyExtensionFiles = '^(?:manifest\.json|background\.js|popup\.(?:html|js)|overlay\.css|icon\.png|ollama\.html|ollama-page\.js|(?:core|providers|ui)[\\/].*)$'
    $files = @(Get-ChildItem $package -File -Recurse | Where-Object {
        $relative = $_.FullName.Substring($package.Length + 1)
        $_.Name -notin $preservedNames -and $relative -notmatch $preservedFolders -and !($splitInstall -and $packageExtension -ne $package -and $relative -match $legacyExtensionFiles)
    })
    # Back up overwritten files before stopping the app or modifying the installation.
    foreach ($file in $files) {
        $relative = $file.FullName.Substring($package.Length + 1)
        $old = Join-Path $install $relative
        if (Test-Path -LiteralPath $old) {
            $copy = Join-Path $backup $relative
            New-Item (Split-Path $copy) -ItemType Directory -Force | Out-Null
            Copy-Item -LiteralPath $old -Destination $copy
        }
    }
    if ($ShowProgress) { $label.Text='Installing update. PlazCode will relaunch...'; [Windows.Forms.Application]::DoEvents() }
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
        if ($ShowProgress) { $label.Text='Waiting for PlazCode to close...'; [Windows.Forms.Application]::DoEvents() }
        Start-Sleep -Milliseconds 250
    } while ($true)
    foreach ($file in $files) {
        $relative = $file.FullName.Substring($package.Length + 1)
        $target = Join-Path $install $relative
        $written.Add($relative)
        New-Item (Split-Path $target) -ItemType Directory -Force | Out-Null
        Copy-Item -LiteralPath $file.FullName -Destination $target -Force
        if ($ShowProgress) { [Windows.Forms.Application]::DoEvents() }
    }
    Start-Process -FilePath (Join-Path $install 'PlazCode.exe') -WorkingDirectory $install
    Write-Host "Updated $current -> $next. Reload the extension, then refresh open AI chat tabs." -ForegroundColor Green
} catch {
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
    exit 1
}
if ($ShowProgress) { $form.Close() }
# Keep backup for recovery; downloaded/extracted staging files can be removed.
Remove-Item (Join-Path $stage 'unpacked') -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item (Join-Path $stage 'release.zip') -Force -ErrorAction SilentlyContinue
