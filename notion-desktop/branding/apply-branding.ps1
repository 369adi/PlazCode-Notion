param([Parameter(Mandatory = $true)][string]$Source)
# Replaces the PlazCode icons and sidebar logo with the PlazCode Notion branding.
$ErrorActionPreference = 'Stop'
$b = $PSScriptRoot
Copy-Item (Join-Path $b 'plazcode.ico') (Join-Path $Source 'agent/assets/plazcode.ico') -Force
Copy-Item (Join-Path $b 'plazcode.png') (Join-Path $Source 'agent/assets/plazcode.png') -Force
$html = (Resolve-Path (Join-Path $Source 'agent/src/desktop.html')).Path
$t = [IO.File]::ReadAllText($html)
$b64 = [Convert]::ToBase64String([IO.File]::ReadAllBytes((Join-Path $b 'plazcode.png')))
$img = '<img class="pcn-logo" alt="AdiCode" src="data:image/png;base64,' + $b64 + '">'
$pattern = '<svg class="logo-svg"[\s\S]*?</svg>'
$n = ([regex]::Matches($t, $pattern)).Count
if ($n -lt 1) { throw 'logo-svg not found in desktop.html' }
$t = [regex]::Replace($t, $pattern, $img)
$css = '<style id="pcn-branding">.brand-mark{border:none!important;background:none!important;box-shadow:none!important;padding:0!important}.pcn-logo{width:100%;height:100%;display:block;object-fit:contain}</style>'
if ($t.IndexOf('</head>') -lt 0) { throw '</head> not found in desktop.html' }
$t = $t.Replace('</head>', $css + '</head>')
[IO.File]::WriteAllText($html, $t, (New-Object System.Text.UTF8Encoding $false))
"Branding applied ($n logo)"
