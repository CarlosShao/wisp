# 权威分母尺：go test -list（编译后的包名册）⇄ 两发实际跑到的 ^=== RUN 名集，逐名作差。
# 起因＝本腿先用 grep '^func Test' 数分母＝356，那把尺⛔ 感知 Go 词法（反引号字面量里的夹具 func 也数了进去）。
$ErrorActionPreference = 'Continue'
$tree = 'D:\tmp\wisp\20261019\227v2\deliv-tree-ac5'
Set-Location $tree
$env:PATH = "$tree\third_party\sherpa-onnx;$tree\build;$env:PATH"
$env:GOTMPDIR = 'C:\Users\swq\AppData\Local\wisp-test-gotmp\227v2-l'
New-Item -ItemType Directory -Force -Path $env:GOTMPDIR | Out-Null
& 'D:\work\base\go\bin\go.exe' test ./cmd/wisp -list '.*' *> 'D:\tmp\wisp\20261019\227v2\list-cmdwisp.log'
"RC=$LASTEXITCODE"
# 取出 -list 的名册（只取 Test 起头的行），剥 CR 后按名排序
Get-Content 'D:\tmp\wisp\20261019\227v2\list-cmdwisp.log' | Where-Object { $_ -match '^Test' } |
  ForEach-Object { $_.Trim() } | Sort-Object -Unique |
  ForEach-Object { $_ -replace "`r", '' } | Set-Content -NoNewline -Encoding ascii 'D:\tmp\wisp\20261019\227v2\list-cmdwisp.txt'
# 名册 ⇄ 两发的实际跑到名集（⛔ 直接比文件，先各剥 CR）
$auth = Get-Content 'D:\tmp\wisp\20261019\227v2\list-cmdwisp.txt' | ForEach-Object { $_.Trim() } | Where-Object { $_ }
foreach ($pair in @(@('mine', 'D:\tmp\wisp\20261019\227v2\full-B.log'), @('orch', 'D:\tmp\wisp\20261019\orch227-b\full.log'))) {
  $tag = $pair[0]; $L = $pair[1]
  $ran = Get-Content $L | Where-Object { $_ -match '^=== RUN   ' } |
    ForEach-Object { ($_ -replace '^=== RUN   ', '') -split '[ /]' | Select-Object -First 1 } | Sort-Object -Unique
  $miss = @($auth | Where-Object { $ran -notcontains $_ })
  "tag=$tag authoritative=$($auth.Count) ran_unique_top=$($ran.Count) not_reached=$($miss.Count)"
  $miss | Select-Object -First 8 | ForEach-Object { "   MISS $_" }
}
