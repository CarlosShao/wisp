# 票 296 验收腿 296-v2 — 证据散件索引（只读裁决腿，非实现者）

起手锚（本腿第一笔 `git log -1` 现认）＝`b9b60985`（2026-10-20 09:1x +08，分支 `dev`）。
被裁实现＝`296-r1`（甲案落地）＋`296-r2`（`AC#5` 三处注释改准，5 笔 `1d6b0cf9`→`280150ca`）。
本腿零 `.go` 字节、零翻框。判语正文＝`docs/evidence/s1/296-v2-verdict.md`。

## logs/ 逐件（尺名＋射程＋读数）

| 件 | 尺 | 关键读数 |
|---|---|---|
| `00-show-U0-commit.txt` | `git show -U0 1d6b0cf9 -- cmd/wisp/resident_windows.go`（对提交本身，⛔ 工作树） | 单 hunk `@@ -40,5 +40,11 @@`；变动 16 行逐枚首字＝`+//` 11／`-//` 5；非注释变动行＝0 |
| `01-cumulative-887dd356-280150ca.txt` | `git diff -U0 887dd356 280150ca -- <该文件>`（起手锚→最终 HEAD 累计尺） | 同为 16 行、`+//`11／`-//`5、非注释 0＝别腿未在此文件落字 |
| `02-head-resident_windows.txt` | `git show HEAD:cmd/wisp/resident_windows.go`（对象层，blob `5011e5d0`） | 改后文档块 `:36-50`；三处改前串 0 命中 |
| `03-before-resident_windows.txt` | `git show 887dd356:cmd/wisp/resident_windows.go`（blob `37915a5a`） | 三处过头话原文 `:40`/`:42`/`:44` |
| `04-head-hotkey_reload.txt` | `git show HEAD:internal/ball/hotkey_reload.go` | `Check()` 两条出口 `:87 if r.src == nil` ⇄ `:93 r.exists && reflect.DeepEqual`；全文件无 all-empty 守卫 |
| `10-pre1.txt` `11-pre2.txt` | 改前台面 treeA（`git archive 887dd356`）`go test ./cmd/wisp -run 'Test296\|Test258' -count=1 -v` | 两发 `RUN=13 PASS_TOP=8 FAIL_TOP=3 SKIP=2`，rc=1 |
| `12-post1.txt` `13-post2.txt` `14-post-extra.txt` | 改后台面 treeB（`git archive 280150ca`）同把尺 | 三发 `13/8/3/2`、rc=1；逐名名册与改前逐字相同 |
| `20-gates.txt` | 门禁并排（共享工作树⇄`git archive HEAD cmd/wisp`→仓外） | 见下 |
| `21-build.txt` `23-vet.txt` | 空捕获注记（raw 无 stdout/stderr） | rc 在 `20-gates.txt` G3/G5 |
| `22-d22scan.txt` | `sh scripts/d22scan.sh` | clean、ban#8 cmd/ 143 枚 Go files |
| `24-pathlen.txt` | `sh scripts/check-path-length-budget.sh` | VERDICT GREEN |

## 门禁四数与两把尺（尺名＝`20-gates.txt`，跑前 `unset GOTMPDIR TMPDIR`）

- G1 工作树 `gofmt -l cmd/wisp`＝**1 枚**＝`cmd\wisp\models.go`
- G2 判据尺 `git archive HEAD cmd/wisp`→仓外同把 `gofmt -l`＝**0 枚**（139 枚 `.go` 非空）
- G2 正控＝同尺下 planted `zz_296v2_posctrl.go` 命中 **1 枚**＝尺灵
- G7 工作树 `gofumpt -l cmd/wisp`＝**1 枚**（同一枚 models.go）；判据尺干净 archive 副本 `gofumpt -l`＝**0 枚**
- G3 `GOFLAGS= go build ./...`＝**rc 0**；G4 `sh scripts/d22scan.sh`＝**rc 0**；G5 `go vet ./cmd/wisp/`＝**rc 0**；G6 长度门＝**rc 0**

## 台面成因（`10`~`14` 逐字，本腿坐实、⛔ 记成 296-r2 新造红）

起手现量 `Get-Process wisp`＝**PID 21892** `build\wisp.exe`（机主常驻，本腿⛔ 停他）。三枚红⇄两枚弃权同一成因＝它占着 `Ctrl+Alt+Q/M/P`：
- `resident_hotkey_258_windows_test.go:361` `the old binding did not survive a no-rebind host (VK=0x0): the premise moved`
- `:448` `the boot verdict lost the live-count shape: … live 1/4: … mute … occupied by another app … panel … occupied by another app`
- `:453` `live slots after the rebind = 1, want 2 (summon+mute; cancel is standby, panel is taken)`
- `resident_hotkey_296_windows_test.go:172`/`:213` `SKIP-LOUD: boot could not register the three default combinations on this host (summon=occupied … mute=occupied … panel=occupied …)`
- `resident_hotkey_v1probe_test.go:108` `258-v1 PROBE: boot verdict lost the live-count shape …`

⇒ 票面 `AC#4` "⛔ 活着的 wisp.exe" 起手前提今天不满足；`296` 那两枚＝abstention-red 仪器，⛔ 当已闭环（`A1041` 那族）。改前⇄改后逐名作差＝**新增红 0 枚**。

## 可比性 seed（本腿现跑）

- treeA(887dd356) ⇄ treeB(280150ca) 在 `cmd/` 目录整树逐文件 sha256 差异＝**仅 `cmd/wisp/resident_windows.go` 一枚**（其余逐字相同）。
- 该文件 blob：`37915a5a`→`5011e5d0`＝`1d6b0cf9` 唯一改动；`280150ca`⇄`HEAD` 该 blob 恒为 `5011e5d0`（工作树亦逐字相等）。
- winsec 仪器件（`342-r2` 新增 `ancestorwindow_342r2_windows_test.go` 等）＝`package winsec`，⛔ 进 `cmd/wisp` 测试二进制；其 `TestAC1…`/`TestAC2…` 名⛔ 匹配 `Test296|Test258`；五发定向读数无点名 winsec 的红。

## 探针件卫生

裸码位 U+00A7 命中＝0（blob 副本内该字符已写成字面 `U+00A7`）；CR 字节＝0；⛔ `.out`；取 blob 落仓外 `D:\tmp\wisp\20261019\296v2\`，副本⛔ 叫 `.go`（`.txt`）。只建不删。
