# Lokaler AdiCode-Build -> erscheint im Update-Tab als "(lokal)". -Install tauscht sofort (Bootstrap).
param([string]$Notes = "Lokaler Build", [switch]$Install)
$ErrorActionPreference = 'Continue'
$repo = 'C:\Users\liket\PlazCode-Shared\work\PlazCode-Notion'
$agent = "$repo\build-src\PlazCode\agent"
$base = "$env:LOCALAPPDATA\PlazCodeNotion"
$log = "$base\logs\local-build.log"
function L($m) { $t = "$(Get-Date -f 'HH:mm:ss') $m"; Add-Content $log $t; Write-Output $t }
New-Item -ItemType Directory -Force "$base\logs", "$base\local-update" | Out-Null
Push-Location $agent
L 'cargo build --release ...'
cmd /c "cargo build --release 2>&1" | Select-Object -Last 15 | ForEach-Object { L $_ }
$code = $LASTEXITCODE
Pop-Location
if ($code -ne 0) { L "BUILD FEHLER ($code)"; exit 1 }
$exe = "$agent\target\release\PlazCode.exe"
$ver = (Get-Content "$repo\notion-desktop\VERSION" -Raw).Trim()
$build = Get-Date -f 'yyyyMMdd-HHmmss'
Copy-Item $exe "$base\local-update\PlazCode.exe" -Force
$meta = [ordered]@{ version = $ver; build = $build; date = (Get-Date -f 'yyyy-MM-dd'); notes = "$Notes`n(lokal gebaut $(Get-Date -f 'dd.MM. HH:mm'))" }
$meta | ConvertTo-Json | Set-Content "$base\local-update\local-update.json" -Encoding UTF8
L "BEREIT: lokales Update $ver build $build ($((Get-Item $exe).Length) B)"
if (-not $Install) { exit 0 }
# Bootstrap: alte App kennt lokale Updates noch nicht -> direkt tauschen.
$app = "$base\app"
$p = Get-CimInstance Win32_Process -Filter "Name='PlazCode.exe'" | Where-Object { $_.ExecutablePath -like "$app*" } | Select-Object -First 1
$swap = @"
Start-Sleep 3
`$app='$app'
Remove-Item "`$app\PlazCode.exe.old" -Force -EA 0
Rename-Item "`$app\PlazCode.exe" 'PlazCode.exe.old' -EA 0
Copy-Item '$base\local-update\PlazCode.exe' "`$app\PlazCode.exe" -Force
`$len=(Get-Item "`$app\PlazCode.exe").Length
'{"version":"$ver","build":"$build","exe_len":'+`$len+'}' | Set-Content "`$app\local-version.json" -Encoding ASCII
if ($($p.ProcessId)) { Stop-Process -Id $($p.ProcessId) -Force -EA 0 }
Start-Sleep 6
if (-not (Get-Process PlazCode -EA 0)) { Start-Process "`$app\PlazCode.exe" -WorkingDirectory `$app }
Add-Content '$log' "`$(Get-Date -f 'HH:mm:ss') installiert $ver build $build"
"@
$swap | Set-Content "$base\local-swap.ps1" -Encoding UTF8
$r = Invoke-CimMethod Win32_Process -MethodName Create -Arguments @{ CommandLine = "powershell.exe -NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File `"$base\local-swap.ps1`"" }
L "SWAP gestartet rc=$($r.ReturnValue)"
