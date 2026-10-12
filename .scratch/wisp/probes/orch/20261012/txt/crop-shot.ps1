param(
    [Parameter(Mandatory=$true)][string]$In,
    [Parameter(Mandatory=$true)][string]$Out,
    [Parameter(Mandatory=$true)][int]$X,
    [Parameter(Mandatory=$true)][int]$Y,
    [Parameter(Mandatory=$true)][int]$W,
    [Parameter(Mandatory=$true)][int]$H
)
Add-Type -AssemblyName System.Drawing
$src = [System.Drawing.Image]::FromFile($In)
$bmp = New-Object System.Drawing.Bitmap($W, $H)
$g = [System.Drawing.Graphics]::FromImage($bmp)
$g.DrawImage($src, (New-Object System.Drawing.Rectangle(0,0,$W,$H)), (New-Object System.Drawing.Rectangle($X,$Y,$W,$H)), [System.Drawing.GraphicsUnit]::Pixel)
$bmp.Save($Out, [System.Drawing.Imaging.ImageFormat]::Png)
$g.Dispose(); $bmp.Dispose(); $src.Dispose()
"saved=$Out"
"bytes=$((Get-Item $Out).Length)"
