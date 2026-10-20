# 票 296 验收裁腿 296-v2 判语（★非实现者·只读·裁 `AC#2` ⇄ `AC#5` 两格）

- 起手锚（本腿第一笔 `git log -1` 现认）＝`b9b60985`（分支 `dev`，2026-10-20 09:1x +08）。
- 被裁实现＝`296-r1`（甲案落地）＋`296-r2`（`AC#5` 三处注释改准，5 笔 `1d6b0cf9`→`280150ca`）。
- 本腿＝裁决者≠实现者：零 `.go` 字节、零翻框、零 push。散件与全部尺＝`.scratch/wisp/probes/296/v2/**`。
- 四格未重开：`AC#0`/`AC#1`（`296-v1` 已判成立）、`AC#3`/`AC#4`（编排者亲跑）＝复跑未撞到与既有判语矛盾的事实。

## AC#5 判语（三处注释改成说实话那一格）＝**成立**（逐格 ≤8 行）

- 尺①（对提交本身，⛔ 工作树）＝`git show -U0 1d6b0cf9 -- cmd/wisp/resident_windows.go`：单 hunk `@@ -40,5 +40,11 @@`，变动 16 行逐枚首字＝`+//` 11／`-//` 5，`^[+-]` 而不以 `+//`/`-//` 开头＝**0 枚**；`--name-only` 越界仅 `cmd/wisp/resident_windows.go` 一枚、`.go` 仅此 1、`internal/`／`_test.go`／winsec＝0。
- 尺②累计（起手锚→最终 HEAD）＝`git diff -U0 887dd356 280150ca -- 同文件`＝同为 16 行、非注释变动 **0** ⇒ 别腿未在此文件落过字、`296-r2` 也没多落。正控（`fc4aedaf` 同尺）＝非注释变动 15+14＝尺灵。
- 内容锚三条逐枚对（凭据＝本腿现读 `git show HEAD:internal/ball/hotkey_reload.go`：`:87 if r.src == nil` ⇄ `:93 r.exists && reflect.DeepEqual(r.applied, cfg)`，两条独立、全文件无 all-empty 守卫）：见下逐处。

| 处 | 改后（HEAD blob `:40-50`） | 判 |
|---|---|---|
| ① | `Check() has two skips, neither of them an all-empty guard: a nil src keeps whatever is live, and a set DeepEqual to what was applied is left alone` | 逐枚点名两条＋具名"两条都不是 all-empty 守卫"；⛔ 旧的"internal/ball has no all-empty guard"全仓否定句已删（该串 HEAD 0 命中）＝说实话 |
| ② | `dead until a read answers a DIFFERENT set - reading the same empties again only hits the DeepEqual skip` | 与 `:93` 一致：同组空串只撞 DeepEqual 跳过、须读到不同集才重绑＝恢复条件没再说轻＝说实话 |
| ③ | `Whether an empty [hotkey] slot can be a user's way of turning a key off is Q-82, still on the owner's list; this closure takes no position on it and only says what it does today - it answers every empty slot with the compiled default before the bridge ever sees the value` | 只述今天这一支怎么处理空串（`ApplyHotkeyDefaults` 于 `hotkey_windows.go:95/98/101/104` 四支填默认，本腿现读核实）＋具名 `Q-82` 未拍、闭包不表态＝⛔ 对 `Q-82` 作裁决＝说实话 |

## AC#2 判语（落地那一发·改后是否只说行为、⛔ 新造承诺）＝**成立**（逐格 ≤8 行）

- `296-v1` 已判修法成立（甲案两支 `return` 都过合并）；本腿复核其仍在位＝`HEAD` blob `:57` 与 `:61` 两枚 `return ball.ApplyHotkeyDefaults(...)` 逐字在、`1d6b0cf9` 未动它们。
- 头句 `NO SLOT IS EVER ANSWERED EMPTY ... Both branches merge through ball.ApplyHotkeyDefaults`＝兑现：两支都过合并，且默认四枚非空（`hotkey_windows.go:75`）⇒ 无空槽漏出。
- 全块**没有新造一句未兑现的承诺**：① 不再写"只有一种跳过"、② 恢复条件写成"读到不同集才重绑"、③ 不裁决 `Q-82`、④ "the bridge does NOT keep the current bindings on an empty answer"与 `:87`/`:93` 两条出口逐字相符、⑤ "an unchanged file costs zero Win32 calls" 那半段未被本程改写（`:33-34` 原句，DeepEqual 支坐实）。
- "four empty strings from here would ... leave ... dead"＝虚拟语气陈述**旧回归成因**，非承诺今天仍会产出空串；与头句不矛盾。

## 台面那几枚的归因（本腿自己坐实、⛔ 记成 296-r2 新造红）

- 起手 `Get-Process wisp`＝**PID 21892** `build\wisp.exe`（机主常驻）占 `Ctrl+Alt+Q/M/P`；本腿⛔ 停他进程＝`AC#4` "⛔ 活着的 wisp.exe"起手前提今天不满足（与派单/编排者现量同）。
- `-run 'Test296|Test258'` 五发（改前 treeA×2、改后 treeB×3，仓外导出树、POSIX 前缀 PATH＋绝对 `go.exe`、`GOTMPDIR` 反斜杠形、⛔ `cygpath -wp` 整条）＝四数逐字相同 `RUN=13／PASS_TOP=8／FAIL_TOP=3／SKIP=2`，逐名作差**新增红＝0 枚**：
  - FAIL×3＝`Test258BridgeMutationNoSrcKeepsOldBinding`(`:361`)／`Test258OccupiedCombinationNamesTheNewValue`(`:448`/`:453`)／`Test258V1ProbeSummonEditRebindsLiveBall`(`v1probe:108`)＝逐字含 `occupied by another app`/`Hot key is already registered`＝台面被占。
  - SKIP×2＝`Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots`/`...WhenConfigUnreadable`＝`:172`/`:213` `SKIP-LOUD` 逐字点名三枚组合键被占＝**abstention-red 仪器，⛔ 当已闭环**（`A1041` 那族）。
- winsec 仪器件（`342-r2`）＝`package winsec`⛔ 进 `cmd/wisp` 测试二进制、其 `TestAC*342r2` 名⛔ 匹配 `Test296|Test258`；五发无点名 winsec 的红；本腿⛔ 动 `internal/winsec/**` 一字节、⛔ 跑其整包。

## 门禁并排两把尺（本腿复跑，跑前 `unset GOTMPDIR TMPDIR`）

- 工作树 `gofmt -l cmd/wisp`＝**1 枚**＝`cmd\wisp\models.go`（CR 334⇄HEAD blob CR 0＝已知 CRLF 存量，本腿⛔ 洗）；判据尺 `git archive HEAD cmd/wisp`→仓外同把＝**0 枚**（139 枚 `.go` 非空）；正控 planted 一枚格式坏件＝该尺命中 1＝尺灵。
- `gofumpt`：工作树 1 枚（同一枚 `models.go`）⇄ 干净 archive 副本 0 枚。`GOFLAGS= go build ./...` rc=0／`sh scripts/d22scan.sh` rc=0（clean，ban#8 cmd/ 143）／`go vet ./cmd/wisp/` rc=0／长度门 rc=0（GREEN）。

## 票 296 框终量（尺＝`grep -cE '^- \[[x ]\]'` 全文）

- 总量 **6** 格（`AC#0`..`AC#5`）＝已勾 **4**／未勾 **2**＝未勾正是本腿被裁的 `AC#2`⇄`AC#5`；本腿一枚未翻（翻勾归编排者凭此判语）。

## 具名上报（交编排者另裁、⛔ 本腿自扩射程）

1. 函数体内 `resident_windows.go:59-60` 两行 `are "unset"`／`disabled by config` 不在 `AC#5` 三枚锚内，`296-r2` 一字未动、本腿维持原判⛔ 动它（是否另格归编排者）。
2. `AC#5` 票面判据写"`go test ... -run 'Test296|Test258'` **rc=0**"，而该 rc 的可达前提是机主常驻 `wisp.exe` 不在场；今天 `AC#4` 起手前提不满足 ⇒ 本腿按"新增红＝0 枚、逐名作差"裁，rc=1 归台面、⛔ 记腿失败。同 `A1041`/`A1042` 那族。

**结论**：`AC#5`＝**成立**；`AC#2`＝**成立**（改后只说行为、无新造未兑现承诺）。两格可翻勾之凭据已具上，翻勾归编排者。
