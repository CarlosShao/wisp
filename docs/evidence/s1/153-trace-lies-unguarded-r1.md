# 153 — 那枚压缩痕的两条谎法 · 实现件 r1（写码程）

票面：`.scratch/wisp/issues/153-the-compression-trace-has-two-unguarded-ways-to-lie-swapping-the-ran-guard-for-need-keeps-all-four-nails-green-and-the-trace-carries-no-taskid.md`
来路：票 139 对抗验收件 `docs/evidence/s1/139-compression-leaves-no-trace-r1-accept-r1.md` §6.2 第 2、3 件。
派单正文：`.scratch/wisp/dispatches/2026-09-26-091x-r152-r153-r151.md` 派单 A。

两格**没有互相抵账**：AC#1 交的是判据（少一形），AC#2 交的是字段（少一枚归因键），
各自有各自的用例、各自的红句、各自的承重答句（§2 / §3 / §4）。

原始读数（逐枚 `go test -v` 全文、驱动脚本、探针源件）落 `.scratch/wisp/probes/153/`；
所有变异只发生在**仓外**快照 `C:\Users\swq\AppData\Local\Temp\wisp153-snap-{1,3,probe}`，
逐发还原 + `cmp` 自证（见 §2.4 / §4.3）。本程未删任何临时件、未在仓库内建 worktree。

**票面四框本程不自勾**；`-done` 后缀未加。

---

## 0. 锚点与工作树现量（开工第一发）

```
$ git rev-parse --short HEAD            # 开工第一发
86b0161                                 ← 与派单共同锚点相同
$ git log --oneline -1 86b0161
86b01613 feat(frontend 10 首启引导): FirstrunScreen props 化 + showcase 第十节
```

⇒ **我这枚锚上是前端会话的提交**（不是编排者那枚台账枚 `f7478d3`）。开工期间 HEAD 又走过：
`5365cb2`（编排者：派单存档落盘）、`ff550f3`（本程 AC#1）、`cdf2471`（前端会话：动画开关机器）。
按派单 §0：**这不算漂移**，只需点名锚上是谁——已点名。本程全部"改前"读数一律钉在 `86b0161`，
不用"当前 HEAD"当改前。

**先确认派前那句"票 139 的实现在 HEAD 上"**（按 sha 取，不读脏工作树）：

```
$ md5sum <快照>/internal/agent/{compress.go,compress_trace_test.go,loop.go}   # 86b0161 取树
9c5015dd23d29060b956502aeb954cac  compress.go
7bf2fbf79f78f935afac6683f906b467  compress_trace_test.go
70f8fb1068dd6541358e0b2ce646d5ca  loop.go
```

三枚逐枚等于 139 验收件 §2 开头记的那三枚（`compress.go 9c5015dd… / loop.go 70f8fb10… /
compress_trace_test.go 7bf2fbf7…`）⇒ `internal/agent/**` 自被验版本 `d949d9c5` 起**零漂移**，
痕确实已经在 HEAD 上，本票是给它补牙、不是重做它。

工作树脏的部分**全是别家的**，本程一律不碰、不还原、不算进任何"零命中"宣称：
`design/**`＋`frontend/**`（前端会话在写）、`cmd/wisp/run.go`＋`cmd/wisp/task_scope_close_151_test.go`＋
`internal/tools/bridge.go`（派单 C＝票 151 的程）。

**本格改了哪些文件**：无（只读）。

---

## 1. 派单带进来的两条"未验证前提"——现量结果：一条成立、一条只成立一半

### 1.1 前提①「M5 那形在未修码上今天就响」——**成立（两半都复跑过）**

先复跑"逃逸"这半（＝139 验收 §2.2 的原句：四枚钉子全在场、整包全绿），在 `86b0161` 的仓外快照上施
`m5.py apply`（`compress.go:170` `if rep.Ran {` → `if c.Need(hist) {`，落地有打印）：

```
[R1]  go build rc=0 | go test ./internal/agent/ -count=1 → rc=0
      RUN(all)=80 PASS_all=80 FAIL=0 SKIP=0 panic=0     ← 四枚 CompressionTrace 用例逐枚 PASS
      logs/r1-m5-4nails.txt
```

⇒ 逃逸复算成立，**不是**背它的数。再复跑"响"那半（同一枚变异 + 本票新增那一发）见 §2.2 的 `F2`。

### 1.2 前提②「普通收尾形 0 枚、重复阶梯退化形 1 枚」——**前半本程亲量成立；后半只在进程内复现，未复跑到盘那一发**

- **前半（普通形 0 枚）成立**：本票新增的那一发就是在量它——1 枚 user 轮、1 006 token 对阈值 384，
  `Need()` 真、`rep.Ran` 假、痕 **0 枚**（`logs/f1-delivered-unmutated.txt`，用例 PASS）。
  根因与本程复算一致：`groupRounds`（`compress.go`）只在 `RoleUser` 开一轮，折叠条件
  `if len(raw) <= c.b.KeepRawRounds { break }` 对 1 枚原始轮直接 break。
- **后半（退化形 1 枚）本程未复跑到盘那一发**：那需要"三条独立重复阶梯"喂真 provider + 读
  `<data>\logs\wisp-*.jsonl`（139 验收 §7.1 的探针形状，落点在 `cmd/wisp`，不是本票地界）。
  本程只在**进程内真 Loop** 上复现了同形的枚数：`buildRoundHistory(4, 400)` 经
  `Loop.Run` 折叠 ⇒ 痕**恰 1 枚**，且走的是进程默认 logger（`logs/ext-A-delivered.txt`，
  见 §3.4）。⇒ 本程引它是引 §7.2/§8 的**两形对照**，未引它 §1（那节被同一枚程自己在 §5 推翻，票面明令不引）。

⚠ 本程没有因为"前提① 成立"就跳过 §2 的改前红句，也没有因为"前提② 后半未复跑"就否认它——
两句各自给了读数与未读数。

---

## 2. AC#1 — 补上"过阈值但折不动 ⇒ 痕 0 枚"那一枚用例（判据形状）

### 2.1 判据与形状

新增 `TestCompressionTraceSilentWhenNothingFoldableOverThreshold`（`internal/agent/compress_trace_test.go:426`）。
它断四件事，缺一就不是本票要的那一形：

| # | 断言 | 为什么必须有 |
|---|---|---|
| 1 | `len(rawRoundIndexes(groupRounds(hist))) == 1` | 钉住"折不动"是**轮数地板**造成的，不是"没过阈值"造成的 |
| 2 | `c.Need(hist)` 为**真**（否则 `t.Fatalf`） | 这正是 139 那四枚的洞：`TestCompressionTraceSilentWhenNothingFolded` 在 `Need()` 为真时**自己 Fatalf 掉了那一支**（票面点名的 `:270-271`），于是"过阈值"这一形在四枚里不存在 |
| 3 | `rep.Ran == false` 且历史逐字节没动（msgs / tokens 两端相等） | 排除"其实折了但没打痕"这一支混进来 |
| 4 | `recs.with(traceMsg)` 长度 **== 0**；不为 0 时先把那条谎痕 `flat()` 打出来再判失败 | 判据本体；且红句里直接带读数，不用二次取证 |

阈值一律从 `BudgetsFor(4096).HistoryCompressTokens` 读（384 是读数不是字面量），
与本文件既有的"不拿 12000 当触发点"口径一致；夹具用 `compress_test.go:23 buildRoundHistory`
（本程未动那个文件，见 §5）。

### 2.2 改前红句原文（**未修码上**，`86b0161` 的 compress.go + 交付版用例 + M5 已落地）

```
$ python m5.py apply <snap-1>
M5 applied: `if rep.Ran {` -> `if c.Need(hist) {`
  landed at compress.go:170: if c.Need(hist) {
$ go build ./internal/agent/ && go test ./internal/agent/ -count=1 -v      # logs/f2-delivered-m5.txt
=== RUN   TestCompressionTraceSilentWhenNothingFoldableOverThreshold
    compress_trace_test.go:458: over-threshold pass that folded nothing left a trace: agent: history compressed compressed_msgs=0 kept_raw_rounds=1 history_changed=false tokens_before=1006 tokens_after=1006 threshold=384 msgs_before=3 msgs_after=3
    compress_trace_test.go:461: "agent: history compressed" records = 1, want none: a record reading "compressed" while nothing folded is a false positive, and 1 raw round(s) against the KeepRawRounds floor of 3 is exactly the case where no fold can happen (all records: [trace-capture-ruler-control agent: history compressed])
--- FAIL: TestCompressionTraceSilentWhenNothingFoldableOverThreshold (0.00s)
```

```
RUN(all)=81 PASS_all=80(63 top + 17 sub) FAIL=1 SKIP=0 panic=0
```

⇒ 三件事同时读到：**只有这一发红**（既有四枚在 M5 下依旧逐枚 PASS，与 §1.1 的 R1 一致）；
**谎痕长什么样**（`tokens_before == tokens_after == 1006`、`compressed_msgs=0`、
`history_changed=false`，正是 139 验收 §2.2 末段推断"比没有痕更误导"的那一形——本程现在把它印出来了，
不再只是推断）；**新路径真执行过**（它红在这发上，不是恒红、也不是恒绿）。

### 2.3 改后转绿（同一发，未施任何变异；工作树交付态）

```
$ go test ./internal/agent/ -count=1 -v      # logs/f1-delivered-unmutated.txt（仓外快照）
--- PASS: TestCompressionTraceBooksCountsOnSuccess          (0.00s)
--- PASS: TestCompressionTraceSilentWhenNothingFolded       (0.00s)
--- PASS: TestCompressionTraceSurvivesTheLoopWiring         (0.02s)
--- PASS: TestCompressionTraceDoesNotAlterTheFold           (0.00s)
--- PASS: TestCompressionTraceSilentWhenNothingFoldableOverThreshold (0.00s)
RUN(all)=81 PASS_all=81(64+17) FAIL=0 SKIP=0 panic=0     → rc=0, ok internal/agent 1.9s
```

```
$ go test ./internal/agent/ -count=1 -v      # logs/ac5-post-count1.txt（工作树，AC#2 也在里面的终态）
RUN(all)=83 PASS_all=83 FAIL=0 SKIP=0 panic=0 → rc=0
```

同一发再还原变异复跑一次（`logs/f3-restored.txt`）：81 / 0 FAIL ⇒ 红绿之差**只由那一枚守卫造成**。

### 2.4 本格自证：快照还原

```
$ python m5.py revert <snap-1> && restore from git show 86b0161:…
RESTORED IDENTICAL internal/agent/compress.go
RESTORED IDENTICAL internal/agent/compress_trace_test.go
RESTORED IDENTICAL internal/agent/loop.go
```

`m5.py` 全程用 `newline=""` 读写：本程第一版没用，导致"还原"后 `cmp` 报了 byte 14 差异
（是 CRLF，不是码），已改脚本并**用改后的脚本重跑了 F1/F2/F3**——上面那三行 RESTORED 是重跑后的读数。

**本格改了哪些文件**：`internal/agent/compress_trace_test.go`（＋71 行，纯追加用例＋文件头第三形说明），
commit `ff550f3`。`compress.go` / `loop.go` **本格零字节**（AC#1 交的是判据，不是码的谎法；
把守卫换成"更保险的双层判"不改变 M5 的可造性，本程不走那条路）。
