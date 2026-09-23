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
| AC#1 五发（N-3／X4／X8／X12／X14 两拍，门关着仍红） | 票面 AC#1 + 归属硬线 | （待填） |
| AC#2 覆盖面主张（清单式 vs 图式＋"删掉它哪条用例会红"） | 票面 AC#2 | （待填） |

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

（待填）

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
