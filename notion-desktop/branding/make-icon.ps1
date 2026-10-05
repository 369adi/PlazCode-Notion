param([string]$Out = $PSScriptRoot)
Add-Type -AssemblyName System.Drawing
function RoundRect([float]$x,[float]$y,[float]$w,[float]$h,[float]$r) {
  $p = New-Object System.Drawing.Drawing2D.GraphicsPath
  $d = $r * 2
  $p.AddArc($x, $y, $d, $d, 180, 90); $p.AddArc($x + $w - $d, $y, $d, $d, 270, 90)
  $p.AddArc($x + $w - $d, $y + $h - $d, $d, $d, 0, 90); $p.AddArc($x, $y + $h - $d, $d, $d, 90, 90)
  $p.CloseFigure(); return $p
}
function Draw([int]$s) {
  $bmp = New-Object System.Drawing.Bitmap $s, $s, ([System.Drawing.Imaging.PixelFormat]::Format32bppArgb)
  $g = [System.Drawing.Graphics]::FromImage($bmp)
  $g.SmoothingMode = 'AntiAlias'; $g.TextRenderingHint = 'AntiAliasGridFit'; $g.InterpolationMode = 'HighQualityBicubic'
  $g.Clear([System.Drawing.Color]::Transparent)
  $u = $s / 256.0
  $black = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::FromArgb(255, 25, 25, 25))
  $white = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::White)
  # outer black body (gives the Notion cube depth on left/bottom)
  $g.FillPath($black, (RoundRect (14*$u) (14*$u) (228*$u) (228*$u) (34*$u)))
  # white face, shifted up/right
  $g.FillPath($white, (RoundRect (40*$u) (26*$u) (190*$u) (190*$u) (22*$u)))
  # bold serif P
  $font = New-Object System.Drawing.Font 'Georgia', ([float](150*$u)), ([System.Drawing.FontStyle]::Bold), ([System.Drawing.GraphicsUnit]::Pixel)
  $fmt = New-Object System.Drawing.StringFormat; $fmt.Alignment = 'Center'; $fmt.LineAlignment = 'Center'
  $rect = New-Object System.Drawing.RectangleF ([float](40*$u)), ([float](30*$u)), ([float](190*$u)), ([float](190*$u))
  $g.DrawString('P', $font, $black, $rect, $fmt)
  $g.Dispose(); return $bmp
}
$sizes = 256,128,64,48,40,32,24,20,16
$pngs = @()
foreach ($s in $sizes) { $b = Draw $s; $ms = New-Object IO.MemoryStream; $b.Save($ms, [System.Drawing.Imaging.ImageFormat]::Png); $pngs += ,($ms.ToArray()); if ($s -eq 256) { $b.Save((Join-Path $Out 'plazcode.png'), [System.Drawing.Imaging.ImageFormat]::Png) }; $b.Dispose() }
$fs = [IO.File]::Create((Join-Path $Out 'plazcode.ico')); $bw = New-Object IO.BinaryWriter $fs
$bw.Write([UInt16]0); $bw.Write([UInt16]1); $bw.Write([UInt16]$sizes.Count)
$offset = 6 + 16 * $sizes.Count
for ($i = 0; $i -lt $sizes.Count; $i++) { $s = $sizes[$i]; $dim = if ($s -ge 256) { 0 } else { $s }; $bw.Write([byte]$dim); $bw.Write([byte]$dim); $bw.Write([byte]0); $bw.Write([byte]0); $bw.Write([UInt16]1); $bw.Write([UInt16]32); $bw.Write([UInt32]$pngs[$i].Length); $bw.Write([UInt32]$offset); $offset += $pngs[$i].Length }
foreach ($p in $pngs) { $bw.Write($p) }
$bw.Close()
