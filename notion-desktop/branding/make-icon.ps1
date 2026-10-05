param([string]$Out = $PSScriptRoot)
# AdiCode logo: rounded square with violet->cyan gradient and a bold white "A".
# Writes plazcode.png (256 px) and plazcode.ico (filenames kept for the build workflow).
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
  $g.SmoothingMode = 'AntiAlias'; $g.TextRenderingHint = 'AntiAliasGridFit'; $g.InterpolationMode = 'HighQualityBicubic'; $g.PixelOffsetMode = 'HighQuality'
  $g.Clear([System.Drawing.Color]::Transparent)
  $u = $s / 256.0
  $rect = New-Object System.Drawing.RectangleF ([float](8*$u)), ([float](8*$u)), ([float](240*$u)), ([float](240*$u))
  $c1 = [System.Drawing.Color]::FromArgb(255, 124, 92, 255)
  $c2 = [System.Drawing.Color]::FromArgb(255, 34, 211, 238)
  $grad = New-Object System.Drawing.Drawing2D.LinearGradientBrush $rect, $c1, $c2, ([float]45)
  $g.FillPath($grad, (RoundRect (8*$u) (8*$u) (240*$u) (240*$u) (56*$u)))
  # soft top highlight
  $hl = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::FromArgb(36, 255, 255, 255))
  $g.FillPath($hl, (RoundRect (8*$u) (8*$u) (240*$u) (112*$u) (56*$u)))
  # A monogram as polygon (crisp at every size)
  $white = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::White)
  $outer = New-Object System.Drawing.Drawing2D.GraphicsPath
  $pts = @(
    (New-Object System.Drawing.PointF ([float](106*$u)), ([float](52*$u))),
    (New-Object System.Drawing.PointF ([float](150*$u)), ([float](52*$u))),
    (New-Object System.Drawing.PointF ([float](208*$u)), ([float](204*$u))),
    (New-Object System.Drawing.PointF ([float](170*$u)), ([float](204*$u))),
    (New-Object System.Drawing.PointF ([float](158*$u)), ([float](170*$u))),
    (New-Object System.Drawing.PointF ([float](98*$u)), ([float](170*$u))),
    (New-Object System.Drawing.PointF ([float](86*$u)), ([float](204*$u))),
    (New-Object System.Drawing.PointF ([float](48*$u)), ([float](204*$u)))
  )
  $outer.AddPolygon([System.Drawing.PointF[]]$pts)
  $hole = @(
    (New-Object System.Drawing.PointF ([float](128*$u)), ([float](92*$u))),
    (New-Object System.Drawing.PointF ([float](147*$u)), ([float](140*$u))),
    (New-Object System.Drawing.PointF ([float](109*$u)), ([float](140*$u)))
  )
  $outer.AddPolygon([System.Drawing.PointF[]]$hole)
  $outer.FillMode = 'Alternate'
  $g.FillPath($white, $outer)
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
