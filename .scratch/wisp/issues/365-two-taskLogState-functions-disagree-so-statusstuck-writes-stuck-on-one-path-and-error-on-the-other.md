# 365 — 两处 `taskLogState` 各写各的：同一枚 `StatusStuck` 一边写 `stuck`、一边写 `error`

- Status: **ready-for-agent**
- 立案：2026-10-23 10:5x +08 编排者立
- ★**来路（实现者签名）**：`196-r1` §9-1（件 `docs/evidence/s1/196-r1-impl.md` **139 行**、提交 `6e19f4a43`／`e89d1fe98`／`2c3628679`）
- ★**编排者已现跑核实**（见下）
- 关联：票 **196**（已落地甲＋乙、`cmd/**` 不在它写面）· `cmd/wisp/**` · `internal/agent/**`

## 这是什么（人话）

**同一个任务状态，两条生产路径写进库里的词不一样。**

★**编排者现跑核实的原话形状**：
- `internal/agent/loop.go:1048` `func taskLogState(s Status) string` ⇒ ★**有 `case StatusStuck: return "stuck"`**。
- `cmd/wisp/run.go:1361` `func taskLogState(s agent.Status) string` ⇒ ★**没有 `StatusStuck` 那一支**，默认落 `:1370 return "error"`。
- ⇒ ★★**同一枚 `StatusStuck`：库里一边是 `stuck`、一边是 `error`。**

★**这就是票 196 那个"可观察分裂"的真身**：票 196 的甲（`internal/memory/taskstate.go` 的权威映射 ＋ DAO 校验）**已经盖住"写进库的那一跳"**，
★**但两个 `taskLogState` 函数本身还在各自映射** ⇒ **分裂的源头没除根**。

## 结案判据

- [ ] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：
      ①**逐条复核上面两处函数**（★带 `file:line` 现跑复认、⛔ 抄本票）；★**行号会不会漂要现跑**。
      ②★**把"分裂"量成一发真读数**：造一发 `StatusStuck` 走 `cmd/wisp` 那条、一发走 `internal/agent` 那条 ⇒ **逐字贴两行库里读到的 `state`**。
      ③★**穷举两处 `case` 集合的差集**：★**除 `StatusStuck` 外还有没有别的分歧**（★`StatusTruncated`／`StatusFailed`／`StatusControl` 各比一遍）。
      ④**正控**：造一发两边都该写**同一个词**的状态（如 `StatusCompleted`）⇒ ★**两处必须同词** ⇒ **证明尺不是全红**。
- [ ] **AC#1 裁形（编排者拍）**：三形各写代价＝
      〔甲〕**`cmd/wisp` 那处改成调用 `internal/agent` 的**（单一实现、`cmd` 侧只留一层薄壳）；
      〔乙〕**两处都改成走票 196 的权威映射 `internal/memory.TaskState`**（`taskstate.go` 已存在）；
      〔丙〕**只补 `cmd/wisp` 缺的那几支**（最小动作、但两处仍在各写各的）。
      ★**未裁形任何实现都⛔ 开工**。
      ★**注意**：★**〔甲〕〔乙〕都要动 `cmd/wisp/**`**（★共享编译面、每波只许一枚写腿）；★**〔乙〕还跨到 `internal/memory`**。
- [ ] **AC#2 落地＋牙（写腿）**：按裁定那一形做；★**每处改动配一枚今天响、不修不响的用例**；
      ★**反向判据**＝把新加的那一跳摘掉 ⇒ 用例必须**转红并逐名报红行**；三读数（变异红／还原绿／原样绿）逐字贴 ＋ **还原核**。
      ★★**验收判据是"两处对同一 `Status` 读出同一个词"**（★⛔ 只测"某处改了"）。
- [ ] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet` 动过的包／`go test -count=1`（**四数分开报、⛔ 相减**）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh` 各 rc。
      ★**格式门禁甲乙两尺都报**（甲＝入库 blob、**判据**；乙＝工作树、**只作诊断**）。
      ★★**`cmd/wisp` 整包四数 ⛔ 当判据**（未判决恒 1、红名册逐发漂移、截断点后整族 `RUN=0`）⇒ 只按名作差＋逐发登记死点。
      ★★**`go test ./cmd/wisp/` 有预存红 `exit status 0xc0000135`（STATUS_DLL_NOT_FOUND、测试进程一次没起）** ⇒ ★**⛔ 当本票的红**（票 196 用原样 HEAD 树复跑同样失败 ⇒ 环境）。
- [ ] **AC#4 纪律格**：⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 196／票 188 任何原句。

## ★硬提醒（⛔ 踩）

1. ★★**`cmd/wisp/**` 每波只许一枚写腿**、且 `go test ./cmd/wisp` 重包同一时刻只跑一个。
2. ★**⛔ 顺手把两处都改成 D43 词**：★**票 196 的名册已经定了 12 枚**（旧词 7 ＋ D43 里任务行真会取的 5 枚）⇒ ★**落地要跟那张名册对齐、⛔ 自造第三套**。
3. ★**票 196 的乙（CHECK 约束）已在库上生效** ⇒ ★**任何新词必须先进名册、否则写不进库**（★那是好事、也是约束）。

## 排程

`365-a1`（只读＝`AC#0`）→ 编排者裁 `AC#1` → `365-r1`（写腿＝`AC#2`）→ `365-v1`（非实现者裁）。

## Progress log

- [2026-10-23 10:5x +08] agent=`编排者` did=立案（来路＝`196-r1` §9-1；★两处函数编排者已现跑核实：`internal/agent/loop.go:1048` 有 `case StatusStuck: return "stuck"`、`cmd/wisp/run.go:1361` 无那一支且 `:1370` 默认返 `"error"`）next=派 `365-a1` 只读腿。
