# 284 交付件：`Gate.Replay` 注释补齐〔零产码调用者／非死代码勿删／C18 一键重放的一半〕

- 锚＝`e60b76877`（起手 `git rev-parse HEAD` 现认，2026-11-05 16:0x +08）。
- 写面现量：本件 1 枚 ＋ `internal/agent/approval/gate.go` 1 枚（**只增注释 9 行，产码语句零改动**）。`frontend/**`／`design/**`／`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`internal/risk/**` **零读零写**。
- 结论＝票 284 `AC#1` 落地（三要素齐），`AC#3` 的两发自证与四道门禁全绿；`AC#2` 的**调度前提**本腿自报为"没搭上别人的 `gate.go`"（见 §5），框一枚没翻（⛔ 归编排者）。

## 1. 票面现读（起手，删歧义）

`.scratch/wisp/issues/284-gate-replays-comment-omits-the-zero-production-caller-claim.md` 存在、**无 `-done` 变体**（`ls .scratch/wisp/issues/ | grep '^284'` 只有一枚）；票面无 `Blocked-by` 字段、无 owner 保留句、无真机项、无前端项 ⇒ **不触 STOP**。

票面"现量"四发逐条重跑（行号已漂，票面写 `:751-754/:755`，HEAD 实在 `:872-875/:876`；票面自己警告过这类漂）：

| # | 票面声称 | 本腿现跑 | 判 |
|---|---|---|---|
| 1 | `[A-Za-z0-9_]\.Replay\(` 在非测试产码里零命中 | `git grep -nE "[A-Za-z0-9_]\.Replay\(" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'` ⇒ `rc=1`（空） | 成立 |
| 2 | 唯一调用点在 `queue_test.go:283` | 同一把尺不剥测试 ⇒ 只多出 `internal/agent/approval/queue_test.go:283` 与 `:285` | 成立 |
| 3 | 唯一下游 `Queue.replay`＋`replayOf` | `queue.go:609`（`func (q *Queue) replay`）／`queue.go:69-70`（`replayOf` 字段） | 成立 |
| 4 | 入口结构性不存在：两侧接口面都不声明它 | `NativeAPI`（`ui.go:156-175`：`Allow`／`AllowSession`／`Reject`）与 `PanelAPI`（`ui.go:180-184`：`Reject`／`Head`／`View`）**均无 `Replay`**；控制台动词 switch（`cmd/wisp/approval_reply.go:556-581`）九个动词 `yes/session/no/veto/panel-no/panel-yes/always/head/view` ＋ `help/quit` 里**没有 `replay`** | 成立 |

## 2. `AC#1` 一句话落地（三要素）

改后注释块 `gate.go:872-884`（逐字）：

```
// Replay re-displays a refused or expired L2 request under a fresh
// correlation id (C18 一键重放). It requires the task to be admitted again
// (D47 still applies to a replay), and it is NOT an answer: the new item must
// still be approved with a newly minted grant.
//
// Today it has no production caller, and that is a state to record, not a
// defect to tidy away: do not delete it as dead code. It is the queue half of
// C18's 一键重放, already implemented here - Queue.replay mints the fresh
// correlation id and links the new item back through replayOf. The missing half
// is the answer-side verb, and neither NativeAPI nor PanelAPI declares one, so
// no card, no panel and no console line can reach this method yet. Adding that
// half is a wiring decision with an owner, which is why the caller side reads
// the way it does rather than because somebody forgot to call it.
```

三要素对账：

| 要素 | 落在哪一句 |
|---|---|
| ① 今天零产码调用者 | `Today it has no production caller` |
| ② 非死代码、勿删 | `that is a state to record, not a defect to tidy away: do not delete it as dead code` |
| ③ 替 C18 一键重放留的 | `It is the queue half of C18's 一键重放 … The missing half is the answer-side verb` |

票面两条禁忌的现跑核对：

- **不嵌行号**：`sed -n '872,884p' gate.go | grep -nE '\.go:[0-9]+|:[0-9]{2,3}[^0-9]'` ⇒ `rc=1`（空）。指路只用具名符号：`Queue.replay`／`replayOf`／`NativeAPI`／`PanelAPI`。
- **不写会被自己扫的现在时计数**：全文**没有一个数字**（`grep -c '[0-9]'` 对该块＝0），"零／0 次"这类会自我失效的 token 一枚没有；句子把"今天"写在最前，读法是**某一天的状态记录**，不是永远为真的不变量。

## 3. RED-before / GREEN-after（文本读数，**不是测试**）

票 284 的性质行逐字写着「⛔ 不接线、不动行为、**不新立仪器**」，禁区又重申「⛔ 不新立"零调用者"负向钉」（票 278 `AC#2` 已裁）⇒ 本腿**没有新增任何 `_test.go`**，「红→绿」由同一把 grep 尺在改前／改后各取一发，读数逐字如下（`sed -n '872,875p'` 改前 / `'872,884p'` 改后）：

| 尺 | 改前（RED） | 改后（GREEN） |
|---|---|---|
| `grep -inE "no production caller\|zero production caller"` | `rc=1`（无此行） | `:877: // Today it has no production caller, and that is a state to record, not a` |
| `grep -inE "dead code\|do not delete"` | `rc=1`（无此行） | `:878: // defect to tidy away: do not delete it as dead code. It is the queue half of` |
| `grep -n "C18"`（契约半件，改前就应有） | `:873`（`rc=0`） | 2 命中 |
| 块字节数 | 275 B（4 行） | 见 §2 逐字块（13 行注释） |

诚实边界：这是一对**读数**，不是一对测试。本票要的正是"注释说了实话"，而"注释有没有说"只有尺能读、没有行为能红；若日后有人把本注释删掉，**不会有任何绿/红变化**——这正是票面裁掉负向钉时接受的代价（票 278 `AC#2`：负向钉在真接线那天必须改期望值）。

## 4. `AC#3` 自证 ＋ 四道门禁（全部现跑，同一环境）

环境：Git Bash；`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；`unset GOTMPDIR TMPDIR TEMP TMP`；**⛔ 未设 `CGO_ENABLED=0`**；`go1.27.1 windows/amd64`。

| 门 | 命令 | 读数 |
|---|---|---|
| vet | `go vet ./internal/agent/approval/` | `rc=0` |
| d22 | `sh scripts/d22scan.sh` | `rc=0`（`d22scan: clean - no D22 ban violations`；正控 `PASS=35 FAIL=0 SKIP=0`，`=== RUN=77`；扫 418 枚产码 Go＋ban#8 internal/ 893 枚含注释与 `_test.go`） |
| build | `go build ./cmd/... ./internal/... ./tools/...` | `rc=0` |
| test | `go test ./internal/agent/approval/ -count=1 -timeout 240s` | `ok … 0.829s`，`rc=0` |
| gofumpt（本文件） | `gofumpt -l internal/agent/approval/gate.go` | **空** |
| gofumpt（整包） | `gofumpt -l internal/agent/approval/` | 列 `internal/agent/approval/pending_read.go` —— **工作树行尾污染，HEAD blob 是 LF**，见下方说明 |

⚠ `pending_read.go` 的读数是**基线，不是本腿造成的**，而且它**在 `git status` 里看不见**：

```
git ls-files --eol internal/agent/approval/pending_read.go
  i/lf  w/crlf  attr/text eol=lf
git ls-files --eol internal/agent/approval/ | 同包其余 36 枚全部 i/lf w/lf（只有这 1 枚 w/crlf）
core.autocrlf=true
```

即：`.gitattributes` 要 LF、index 里是 LF（`git hash-object` ＝ `git rev-parse HEAD:…` ＝ `a7db9b104489b4e3b1a6b2e2e619cac8865f04e1`、`git status --porcelain` 空），只有**工作树那一份**被写成 CRLF（131 行）——autocrlf 在比对时把它归一掉，所以 `git status` 读不出这枚差异。⇒ gofumpt 报的是这枚**行尾**，不是内容。台账 `A756` 那张表里已经有同一个数（`pending_read.go`：工作树 CR 字节 **131**／HEAD blob **0**／gofumpt -d **0**）⇒ 已知基线，**不在票 284 射程内，本腿一字未动**（⛔ 不顺手改别人的文件，哪怕 git 看不见）。

## 5. `AC#2`（随批）——本腿必须自己报的那一格

票面 `AC#2` 逐字：⛔ 不许为这一行单独开一条腿去改 `gate.go`，落地时机＝下一枚**本来就动** `gate.go` 的腿。

本腿的自报：**起手 `git status --porcelain` 对整棵树的非 `??` 行＝0**（`modified-count=0`），`git diff --stat -- internal/agent/approval/` 空 ⇒ **盘上没有别人的 `gate.go` 改动可搭**。本腿是**被专门派来做票 284 的一发**，因此从字面上看它**正是 `AC#2` 想避免的那条腿**。

两件只能说清、不能装作不存在的事：

1. **代价是真的**：为了 9 行注释，本腿花掉一次 `gate.go` 的编译面（build 全量＋`approval` 包测试＋vet＋d22 全量 52.9s）。
2. **编排者那一侧的前提这一发是成立的**：`A754` 立票时定的调度律是"⛔ 不为一行注释单独开腿"、并把它排在"下一枚本来就动 `gate.go` 的腿"；`A757`／`A760` 的编队现状两条都写着"Go 编译面此刻空 ⇒ 下一枚写腿可派"，而本腿起手时编译面确实空（上表 `modified-count=0`、`git diff --stat -- internal/agent/approval/` 空）。⇒ 这笔代价买到的不是"一次白跑的编译"，而是**一格挂在票池里 27 天未落地的诚实缺口**（`A754` 立票 2026-10-09，本腿落地 2026-11-05）。

⇒ 给编排者的**具名待办**（不在本腿翻框）：若认定 `AC#2` 未被满足，处置是"本腿读数记进 `AC#3`、`AC#2` 保持不勾、票不 `-done`"，本腿**不反对**；本腿只保证不把它悄悄记成"随批"。

## 6. 本腿**没做**什么（逐条可核）

| 禁区 | 现量 |
|---|---|
| ⛔ 不接线 | `git diff -- internal/agent/approval/gate.go` 只有注释行；`Gate.Replay` 的函数体零字节变化；两侧接口面（`ui.go`）零字节 |
| ⛔ 不动 C18 契约文本 | `git status --porcelain docs/PLAN.md docs/specs` 空 |
| ⛔ 不新立仪器 | 本次提交 `_test.go` 零枚（`git diff --name-only` 只两枚：`gate.go` ＋ 本件） |
| ⛔ 三枚冻结件 | `internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go` 零字节（不在 pathspec 里） |
| ⛔ 零翻框 | 票面 3 格保持未勾，`-done` 未加（`git status --porcelain .scratch/wisp/issues/284*` 空） |
| ⛔ 零 push | 只 commit，未 push |
