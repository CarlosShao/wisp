# 票 235 v1（非实现者裁决腿 `235-v1`）逐发读数台件

- 起手锚点：`git rev-parse --short HEAD` ＝ **`98df640a`**（本腿自取，一切行号按此锚点现读）。
- 本腿身份：编排者之外的非实现代理（裁决腿）。**零实现码**：没有改过任何 `internal/**`／`cmd/**` 文件，
  所有突变只经 `go test -overlay` 施加。⛔ 票面四枚 AC 框一枚没碰（读数见 §6）。
- 环境注记：本机 self-hosted runner 上 23:2x 前后在跑 CI（含 `slo-full`），会与测量抢 CPU。
  本腿每一发红都做过隔离复量（见 §4 的两枚红归因）。
- 所有 `go test` 起手都带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`。

## 0. 改动面与合规（本腿现跑）

| 尺 | 命令 | 读数 |
|---|---|---|
| 改动面 | `git diff --name-only 25556879..HEAD -- internal cmd` | 只有 `internal/tools/subagent_197_test.go`＋`internal/tools/subagent_222_test.go` ⇒ **零产码**，与交件自述一致 |
| 197 件删除列 | `git diff -U0 25556879..HEAD -- .../subagent_197_test.go` | **8 行里 5 行删除全是注释行**（`-// Why a bigger pool is a lie…` ×5），13 行新增全是 `-//` 注释 ⇒ 函数体零字符变更 |
| 222 件删除列 | `git diff -U0 … -- .../subagent_222_test.go` | 删除仅 3 行＝`newH222` 里 `h:/name:/in:` 三行的 gofumpt 对齐重写，替换行同字段＋`arrived:` ⇒ **零断言丢失** |
| 历史句 | `grep -n "7/8-red shape ticket 197 leg A reported" internal/tools/subagent_197_test.go` | 命中 **`:407`** 逐字在位 |
| 被证伪那句 | `grep -rniE "bridge slot for its child\|whole life\|bigger pool is a lie" --include=*.go internal cmd` | `subagent_197_test.go` **零命中**（只剩三处无关件）⇒ 假理由确实从盘上消失 |
| 注释新占位 | `sed -n '398,411p'` | 新注释＝`:398-410`（13 行），函数名行 `:411` ⇒ 交件自述的行号成立 |
| 票面框数 | `^- \[ \]`／`^- \[x\]` 计数（锚点 vs `25556879`） | 票 235：4 未勾／0 已勾（前后同）；票 222：4 已勾／2 未勾（**前后同**，追加行是 `:36` 的子弹条，没动任何框） |

## 1. 我自己造的五枚突变（全走 `-overlay`，副本在 `mut/`，清单在 `overlay-*.json`）

| 别号 | 内容 | 现场件 |
|---|---|---|
| `m3-afterfix` | HEAD 的 `subagent_222_test.go` ＋ **两行**同时改指 `&fake197Dir{}`（`ParentTools` :302／`Tools` :312） | `m3-afterfix-v1.txt` |
| `m3-beforefix` | `1d2ad737`（AC#2 之前）那份 222 件 ＋ 同两行改指 fake（旧行号 :278／:288 ⇒ 票面现量 #4 的旧行号本腿独立复现） | `m3-beforefix-v1.txt` |
| `prereadoff-m3` | 终态件上把 3 枚 `h.awaitChildOnBridge222(t)`＋1 枚 `h.drainBridgeArrivals()` 注释掉，再叠 M3 | `ac3-1-prereadoff-m3-v1.txt` |
| `commentback` | 终态 `subagent_197_test.go` 的 `:398-410` 换回旧那 5 行（逐字取 `25556879`） | `ac3-2-commentback-v1.txt` |
| `m6-pool8` | 产码副本 `MaxConcurrentSubagents = 4` ⇒ `8`（一行，工作树未动） | `m6-pool8-v1.txt` |
| `pre235-baseline` | 两枚 `_test.go` 全部换回 `25556879` 版（本腿自己的"改动前"分母） | `pre235-baseline-v1.txt` |

## 2. 读数（四数只从 `-v` 量）

| 发 | RUN | 全名 | 顶层 PASS | `--- FAIL` | 红名册与秒数 |
|---|---|---|---|---|---|
| 终态 HEAD | 218 | 218 | 169 | **0** | 无（`ok 13.014s`） |
| `pre235-baseline` | 218 | 218 | 169 | **0** | 无（`ok 12.397s`）⇒ **票面"起手 161"在 235 动手之前就已过期**，不是本票改出来的 |
| `m3-beforefix` | 218 | 218 | — | **1** | `Test222SpawnConclusionArrivesThroughRealBridgeChildren (30.00s)`，红在 `:374` 护栏（"只等到 0/4 枚"）；另两枚 0.00s 绿 |
| `m3-afterfix` | 218 | 218 | — | **3** | 三枚 222 用例全红**全 0.00s**，红句三处同源＝`桥上的到达信号一枚都没有…`（`:441`／`:448`／`:455` 三条正文分别落在 `:444`／`:536`／`:604`），累计执行计数 0／0／6 |
| `prereadoff-m3` | 218 | 218 | — | **1** | 退回 `Test222SpawnConclusion… (30.00s)`，红在 `:449` 的护栏；另两枚 0.00s 绿 |
| `commentback` | 218 | 218 | 169 | **0** | `goexit=0`、`ok 12.444s`、红名册为空 |
| `m6-pool8` | 218 | 218 | — | **4** | `Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s, :413)`／`Test197SubagentPoolCapsAtBridgeCeiling (0.00s, :452)`／`Test197FullPoolRefusesNextSpawnWithReadableReason (0.00s)`／`Test222SpawnConclusionArrivesThroughRealBridgeChildren (3.00s, :488)` |

逐名对拉（本腿自己两发，不借交件的分母）：`v1-base-names.txt`（pre235）↔`v1-head-names.txt`（终态）＝218↔218，
`comm -23`／`comm -13` **双向差集皆空**。

## 3. 门禁四件（本腿现跑）

- `"$GOPATH/bin/gofumpt.exe" -l internal/tools/` ⇒ **零输出**。
- `go vet ./internal/tools/` ⇒ **clean**。
- `bash scripts/d22scan.sh` ⇒ `d22scan: clean - no D22 ban violations`（ban #8 射程 `internal/` 462 枚 Go 件含注释与 `_test.go`）。
- `-count=3 -run 'Test222|Test197SubagentPool|Test197SpawnDescription'` ⇒ **18 名全 PASS／0 红**（`stability-count3-v1.txt`，`ok 0.051s`）。
- `-race -count=2 -run Test222` ⇒ ⚠ **本腿这一发撞出 1 枚红**（`race-222-v1.txt`：`--- FAIL: Test222CeilingStillCapsExecutedCallsWhileParentsWait (0.00s)`，
  红句 `subagent_222_test.go:555: 只有 3/4 枚父任务到达等待点（父任务没派生成功…）`），`DATA RACE` 计数 **0**。⇒ 见 §4 归因。

## 4. 那枚红的归因（本腿的复量链，⛔ 新造真红＝0 枚）

1. 隔离复量 `-race -count=20 -run Test222` ⇒ **0 红**；`-race -count=100 -run Test222CeilingStillCaps…` ⇒ **0 红**；
   12 枚 CPU -spinner 压载下 `-count=150 -run Test222CeilingStillCaps…` ⇒ 终态 **0 红**、pre235 **0 红**（`flake-under-load-*.txt`）。
2. 决定性对照：同一枚命令 `-race -count=60 -run Test222`（三枚 222 用例一起）⇒
   **终态红 1 枚**（`flake-race60-head.txt`，红在 `:514`＝leg2 的 `parkParents` 调用点）、
   **pre235 基线也红 1 枚**（`flake-race60-pre235.txt`，红在 `:368`＝leg1 的同一枚 helper）。两发 `DATA RACE` 都是 0。
3. 机制（本腿现读产码得到的结论，不是猜）：`internal/tools/subagent_197.go` 里
   `bg := child.RunAsync(childCtx, prompt)`（`:347`）**先**于 `giveBackWhileWaiting(ctx)`（`:365`）与
   `onUpdate("task.spawn 已派生子代理 …")`（`:368`）。所以"孩子报 `started`"**不**蕴含"父的派生 delta 已落 `h.parked`"，
   而 `parkParents`（`subagent_222_test.go:366-372`）在等满 n 枚 `started` 之后**立刻**要求 `len(h.parked)==n`。
   `-race` 的调度延迟把这个既有边际假设撬开 ⇒ 红。**这段 helper 票 235 一字未动**，前置读数跑在它之后 ⇒ 不归本票。
4. 结论：**票 235 的 diff 没有制造红**；本腿新造真红 **0 枚**。但这条 flake 是真的，落在票 222 那族用例的地界（见裁决表 §6 交给编排者）。

## 5. 每发之后的还原自证（`git status --porcelain -- internal cmd`）

每一发（含 §2 七发、§3 两发、§4 六发复量，共 15 发）之后都跑过这一把尺：**每一次输出都是空**。
上表未逐行重复，此处一次性记账；命令与现场件名一一对应，可逐发核。

## 6. 票面 AC 框（本腿零写入的证据）

`235-…bridge.md`：`^- \[ \]`＝**4**、`^- \[x\]`＝**0**（与 `25556879` 同）。
`222-….md`：已勾 4／未勾 2（前后同）。本腿对 `.scratch/wisp/issues/**` **零写入**（`git diff --name-only 98df640a..HEAD` 里本腿的提交只含 `docs/evidence/s1/235-*.md` 与 `.scratch/wisp/probes/235/v1/**`）。
