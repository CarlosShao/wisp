# 330-ac5 台面配方与量尺（setup gauges，2026-11-07 12:4x-13:0x +08）

- 锚：`58eca3738fccda15662be5e1b640290aab776e4f`（`58eca3738`，2026-11-07 12:46:26 +0800，本腿起跑时 HEAD 现认）。两棵树同锚。
- 树：终态＝`D:\tmp\wisp\20261107\330ac5\export`，起手＝`D:\tmp\wisp\20261107\330ac5\export-base`。
  建法＝`git clone --shared "<repo>" <tree>` 后 `git checkout 58eca3738…`（⛔ 共享工作树；共享树 `cmd/wisp`
  有 2 枚别人未跟踪件 `panel_resident_323_windows_test.go`／`panel_resident_assembly_323_test.go`，
  两棵导出树 `git status --porcelain -- cmd/wisp` ＝ 0 行＝零污染）。
- 配方：clone 自带八件套（`go.mod go.sum internal cmd tools scripts frontend models`）
  ＋拷仓根 `docs/`＋手工三件（`build/*.dll` 3、`third_party/sherpa-onnx/*.dll` 3、`frontend/dist` 4 件）。
- 非空尺：`.go` ＝ 1559（两棵同数）。
- 8.3 短名：`cygpath -d` ＝ `D:\tmp\wisp\20261107\330ac5\export`（自带短形）／`D:\tmp\wisp\20261107\330ac5\EXPORT~1`。
- 跑前 env：`unset GOTMPDIR TMPDIR TEMP TMP`；`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；
  本腿另设 `GOCACHE=/d/tmp/wisp/20261107/330ac5/gocache`（两发共用，只影响编译缓存、不影响判据）。
- 命令：`go test ./cmd/wisp/ -count=1 -timeout 30m -v`（⛔ `./...`；`-v` 为四数点名所加）。
- 工具链：`go1.27.1 windows/amd64`；解析尺＝`parse-roster.sh`（纯 awk；口径：RUN＝全部 `^=== RUN` 含子测试，
  PASS/FAIL/SKIP_TOP＝仅零缩进判决行，顶层＝零缩进 RUN 且名字不带 `/`，死点＝顶层 RUN 无判决行，
  PANIC＝`panic:` 行）。★go1.27 的子测试 `=== RUN` 行也在零缩进（`logs/330ac5-final-full.log:817` 实证），
  顶层判定靠名字不带 `/`，第一版尺（按缩进判）作废，判据以 v2 为准。
- 实跑读数：终态 `rc=1`（test 292.281s / real 5m42.067s）；起手 `rc=1`（test 291.021s / real 5m8.757s）。
- ★LF 披露：两枚整发日志各含恰好 1 个 CR 字节（终态 offset 206154／起手 offset 206355），位置都在
  `TestTicket198AC3*` 用例捕获的**子进程原样输出**里（`Successfully processed 1 files; Failed processing 0 files\r\n`），
  是被捕获的数据字节、非本腿行尾形状（两文件行尾全 LF）。本腿选择字节忠实保留、⛔ 改原始读数。

