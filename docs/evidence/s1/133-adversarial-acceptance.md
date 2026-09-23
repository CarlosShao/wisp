# 票 133 AC#1 / AC#2 —— 非实现者验收方的裁决表

**执行方**：`acceptor-ticket133-r1`（**非实现者**，只读被验代码；不复用实现方的任何判据）
**被验版本锚定**：`git rev-parse --short HEAD` = **`52cf311`**（本验收方开工时刻的工作树 HEAD，
仅供对照）；**被验 sha = `c5f140c`**（= `c5f140c7a1c2be66ea16f1fcb09ae56737ea1ccf`）。
**纯净树取法**（不在仓库内建 worktree、不 checkout）：

```
$ mkdir -p /d/tmp/wisp133-acc-r1
$ git archive c5f140c | tar -x -C /d/tmp/wisp133-acc-r1        # 有仪器的被验树
$ mkdir -p /d/tmp/wisp133-acc-r1-noins
$ git archive c720494 | tar -x -C /d/tmp/wisp133-acc-r1-noins  # 无仪器的"①栏"树（前任的基线树）
$ cp third_party/sherpa-onnx/*.dll <树>/third_party/sherpa-onnx/   # 三枚 DLL，缺即 0xc0000135
```

`git diff --stat c720494..c5f140c -- cmd/wisp/` ⇒ 只有 `leg_dispatch_gate_133_test.go` +1377 与
`main.go`/`panel_assets.go`/`slo_windows.go` 的 +19/+8/+6 行（裁决注释），与前任 §4.5 同形。

开工时刻 `date -u` = 2026-09-23 12:11:01z（本地 20:11 +08）。

## 0. 骨架（本节以下逐节填，每裁一格 commit 一次）

| 格 | 判据出处 | 判定 |
| --- | --- | --- |
| AC#1 五发（N-3／X4／X8／X12／X14 两拍，门关着仍红） | 票面 AC#1 + 归属硬线 | **PASS**（§2：六栏三态独立复现，③栏红名逐发点到 `TestAC1AC2DispatchHopGate133`；①栏"今天不红"这一步的语义已失效，裁定见 §5，不构成本格的减分） |
| AC#2 覆盖面主张（清单式 vs 图式＋"删掉它哪条用例会红"） | 票面 AC#2 | **退回**（§3：覆盖判据 (d) 今天可被一枚**永不运行的同名方法**满足，红被洗掉；另两处同族洞。§6：§3 那句"摘掉本尺五发零红"与它自己 §2.0.1 矛盾） |

## 1. 独立基线（未变异，四枚读数）

驱动器 `/d/tmp/wisp133-acc-r1-run.sh`：`go test -count=1 -v ./cmd/wisp/`，`PATH` 前置该树自己的
`third_party/sherpa-onnx`（三枚 DLL）。四数只从 `-v` 输出量（`^=== RUN` / `^--- PASS` /
`^--- FAIL` / `^--- SKIP`）。日志 `/d/tmp/wisp133-acc-r1-out/A-base-*.log`。
测量 `date -u` 12:16:32z–12:20:1xZ（本地 20:16–20:20 +08）。

| 读数 | 树 | 附加旗标 | rc | RUN | PASS | FAIL | SKIP | 对前任 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `A-base-ins-full` | `c5f140c`（有仪器） | 无 | 0 | 101 | 54 | 0 | 0 | §2.0 `S0-ins-pristine` 101/54/0/0 逐字同 |
| `A-base-noins-full` | `c720494`（无仪器） | 无 | 0 | 100 | 53 | 0 | 0 | §1 基线 100/53/0/0、§2.0 `S0-base-pristine` 逐字同 |
| `A-base-ins-doorclosed` | `c5f140c` | `-skip '^TestAC4EveryLegIsNailedOrRuled$'` | 0 | 100 | 53 | 0 | 0 | §2.0 `S0-ins-doorclosed` 100/53/0/0 逐字同 |
| `A-base-ins-skipself` | `c5f140c` | `-skip '^TestAC1AC2DispatchHopGate133$'` | 0 | 100 | 53 | 0 | 0 | §4.3 基线独立复现 100/53/0/0 逐字同 |

**"变绿"与"被跳过"分清**（关门的读数是"少跑一枚"而不是"跑成一枚 SKIP"）：
`-skip` 后 RUN 101→100、PASS 54→53、**SKIP 仍 0**（`--- SKIP` 行不产出），
且被验尺自己在那份日志里仍有 `=== RUN   TestAC1AC2DispatchHopGate133` 一行
（`grep -c '^=== RUN   TestAC1AC2DispatchHopGate133$' A-base-ins-doorclosed.log` = 1）。
⇒ 门关着的三枚读数里，本尺是**真跑了再判**的，不是被一起跳掉才绿。

**门禁同形**（本验收方自己量，见 §7）：`c5f140c` 树上 `go build ./...` rc=0、`go vet ./cmd/wisp/` rc=0。

## 2. AC#1 五发独立复跑（①无仪器／②门开着／③门关着）

本验收方自己写的变异台 `/d/tmp/wisp133-acc-r1-mutate.py`（每处改动先 `assert` 锚点唯一、改完打印
落地行 + `go build ./...` rc，然后才读数；还原＝从 `.acc-r1.bak` 拷回 `main.go`、种件改名 `*.off` 不删），
驱动器 `/d/tmp/wisp133-acc-r1-run.sh`，日志 `/d/tmp/wisp133-acc-r1-out/*.log`。
门关着的做法：`-skip '^TestAC4EveryLegIsNailedOrRuled$'`（131 那扇门不跑＝不存在，它的文件一个字没动）。

**三轮种法的账（先说清，因为前两轮的读数与前任不符，原因在本验收方的种法，不在被验的树）**：

| 轮 | 种法 | 后果 |
| --- | --- | --- |
| v1 | 四条新腿都用 `if` 早退分发，且种件里调 `resolveDataDir` | 除 131 的门以外，还多点红出票 128 的门（`cmdXxx131 解析了数据根却没有腿驱动它`），且 128 的门 `t.Fatalf` 掉自己 5 枚子用例 ⇒ RUN 分母从 100 缩到 95 |
| v2 | X12／X14 改用票面明写的 `case "…":` 进 `switch`；X4 保留早退 `if`（那就是 X4 的形） | 四数与前任逐字对上；X14 两拍仍多一枚红：种件里那行 `slog.Info` 走的是**进程级听众**，131 的 R-117-1 判据（记了账却没装听众、又没裁决句）据此判红 |
| v3 | 在 v2 之上把种件里的 `slog.Info` 去掉（复验方当年那条账本原文就是 `install=false records=false`） | 六栏读数与前任 §2.1–§2.5 逐字对上 |

**三轮都留盘**（`shots.txt`＝v1，`shots-v2.txt`，`v3-probes.txt`＝v3＋仪器探针），因为前两轮的
红名差本身就是读数：**同一形在不同种法下会多点红不同门**，这对"归属硬线"的判断有影响（见 §5）。

### 2.1 五发三态总表（采用与票面形状同形的 v2/v3 轮；N-3／X8 只有 v1＝v2 同种法）

| 发 | 轮 | ①无仪器 `c720494` | ②门开着 `c5f140c` | ③门关着（`-skip` 131 的门） | ④还原 |
| --- | --- | --- | --- | --- | --- |
| N-3 摘掉 `cmdModels` 的调用者 | v1 | rc=1 100/52/1/0 红=`TestAC4EveryLegIsNailedOrRuled` | rc=1 101/52/2/0 红=本尺＋131 的门 | rc=1 100/52/1/0 红=**只有 `TestAC1AC2DispatchHopGate133`** | difflines=0 |
| X4 早退 `if` 的 `--diag` 腿 | v2 | rc=1 100/52/1/0 红=131 的门 | rc=1 101/52/2/0 红=本尺＋131 | rc=1 100/52/1/0 红=**只有本尺** | difflines=0 |
| X8 `case "slo":` 换命名常量 | v1 | rc=1 100/52/1/0 红=131 的门 | rc=1 101/52/2/0 | rc=1 100/52/1/0 红=**只有本尺** | difflines=0 |
| X12 install 经包级函数值别名 | v2 | rc=1 100/52/1/0 红=131 的门 | rc=1 101/52/2/0 | rc=1 100/52/1/0 红=**只有本尺** | difflines=0 |
| X14 一拍 字段初值装听众 | v3 | **rc=0 100/53/0/0 零红（门 PASS）** | **rc=1 101/53/1/0 红=只有本尺** | **rc=1 100/52/1/0 红=只有本尺** | difflines=0 |
| X14 二拍 拆掉 install | v3 | **rc=0 100/53/0/0 零红（门 PASS）** | **rc=1 101/53/1/0 红=只有本尺** | **rc=1 100/52/1/0 红=只有本尺** | difflines=0 |
| 还原后的整包读数 | — | — | — | — | `Z2-restored-ins` rc=0 101/54/0/0 |

落地证明（每发读红名之前先跑的两条，原文节选）：

```
LANDED-MAIN x14:            main.go:80:  case "sfx131:"     sfx131x.go:16: var holder131 = sinkHolder131{open: installLogSink}
                            sfx131x.go:23: if _, err := holder131.open(args[0]); err != nil {
BUILD 1v3noins-x14c rc=0 / BUILD 2v3ins-x14c rc=0        （X14 一拍）
LANDED-MAIN n3:  cmdModels 的调用从 main.go 消失，case "models" 那三行仍在（脚本按"调用为 0"判落地）
```

真二进制活性（票面 X4／X14 的前提"不是死代码"，本验收方自己量的那一发）：

```
LIVE x4   cmd="wisp --diag <tmp-root>"  rc=0   jsonl: <root>/logs/wisp-20260923-001.jsonl
LIVE v3x14 cmd="wisp sfx131 <tmp-root>" rc=0   jsonl: <root>/logs/wisp-20260923-001.jsonl
```

### 2.2 ③栏红名原文（门关着，即"不许蹭 131 的门"那一格）

N-3（`3ins-n3-doorclosed.log`，四条，与前任 §2.1 逐字同）：

```
leg_dispatch_gate_133_test.go:169: AC#1 RED: cmdModels (models.go:104) reaches installLogSink (or is a production entry of this package's dispatch), and no chain from func main reaches it any more.
leg_dispatch_gate_133_test.go:169: AC#1 RED: modelStore.handOffModel (models.go:302) reaches installLogSink ...
leg_dispatch_gate_133_test.go:169: AC#1 RED: modelsEnsure (models.go:262) reaches installLogSink ...
leg_dispatch_gate_133_test.go:175: AC#1 RED: leg "models" (main.go:93) is claimed by nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" on entry "cmdModels", and the dispatch no longer reaches that symbol.
```

X4（`3v2ins-x4-doorclosed.log`）：`AC#1 RED: leg "--diag" (main.go:79) reaches installLogSink on this path: dispatch -> installLogSink` ＋ usage 双向差一枚。
X8（`3ins-x8-doorclosed.log`）：四条——非字面量标签 / 合成腿未覆盖 / `slo_windows.go:186` 的裁决句指向不在册的腿 / usage 写着 `slo` 而分发不再走它。
X12（`3v2ins-x12-doorclosed.log`）：`AC#1 RED: leg "fake131" (main.go:80) reaches installLogSink on this path: dispatch -> sinkAlias131 -> installLogSink` ＋ usage 一枚。
X14 一拍（`3v3ins-x14c-doorclosed.log`）：`AC#1 RED: leg "sfx131" (main.go:80) reaches installLogSink on this path: dispatch -> holder131.open -> installLogSink` ＋ usage 一枚。
X14 二拍（`3v3ins-x14b2c-doorclosed.log`）：`AC#1 RED: leg "sfx131" (main.go:80) is dispatched by func main and covered by nothing: no nail ... no test case ... and no WISP-LEG-COVERAGE-RULING: sentence naming it.` ＋ usage 一枚。

⇒ **五发的③栏红名都点到 `TestAC1AC2DispatchHopGate133`**，不是包级 `[build failed]`（六发 `go build ./...` 全 rc=0），
也不是 131 的门（③栏里它的名字在日志中出现 0 次，见 §1 的做法段）。**AC#1 的归属硬线：本验收方复跑成立。**

### 2.3 与前任读数的分歧（两边都留着）

| 处 | 前任 §2 读数 | 本验收方 | 归因 |
| --- | --- | --- | --- |
| §2.2 X4 ①②③ | 100/52/1、101/52/2、100/52/1 | v1 轮 95/51/2、96/51/3、95/51/2；**v2 轮逐字相同** | v1 种件调了 `resolveDataDir` ⇒ 票 128 的门也红并 `t.Fatalf` 掉自己 5 枚子用例（分母缩小，不是变绿）。v2 起本验收方改种法，与前任命题同形 |
| §2.4 X12 ①②③ | 100/52/1、101/52/2、100/52/1 | v1 轮 95/51/2 …；**v2 轮逐字相同** | 同上（v1 用 `if` 早退分发＋`resolveDataDir`） |
| §2.5 X14 两拍 ① | **rc=0 100/53/0/0 全绿** | v1/v2 轮 rc=1 100/52/1（131 的门红）；**v3 轮逐字相同（rc=0 全绿）** | 种件里那行 `slog.Info` 走进程级听众，触发 131 的 R-117-1 判据"记了账没装听众"。复验方当年的账本原文是 `records=false`，v3 去掉该行即同形 |
| §2.5 X14 ②栏 | 101/53/1/0 | v3 轮 101/53/1/0 同 | 一致 |

⇒ 结论：**每一处分歧都能被本验收方自己的种法细节解释掉**，改到与票面形状同形后六栏逐字复现。
不以本验收方的前两轮机件故障充抵被验方读数，也不改前任任何一行。

## 3. 首要攻击：这枚仪器自己的洞（同名方法／helper 名）

（待填）

## 4. 同形排查（first-match／按下标配对／只认前缀／只数不核）

（待填）

## 5. §2.0.1 语义裁定（票面"先证明现在不红"还剩几发）

（待填）

## 6. AC#2 覆盖面主张裁定＋"删掉它哪条用例会红"

（待填）

## 7. 四数账与门禁抽查复算

（待填）

## 8. R-133-x 清单

（待填）

## 9. 总判

（待填）
