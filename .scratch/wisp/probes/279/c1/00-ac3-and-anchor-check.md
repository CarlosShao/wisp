# 279/c1 — AC#3 终态复核（12 处仍全为注释行）＋ AC#2 触发条件登记（P39 未贴＋那枚行号锚 today 漂不漂）

- 腿：`279-c1`（只读复核腿）。起手 HEAD＝`dc65df28a12043c2f5ce3b6f3b9ba1be604445ca`（分支 `dev`，该提交时刻 `2026-10-09 10:42 +0800`）；全程 HEAD 未变（`git rev-parse --short HEAD` 起手与收工同为 `dc65df28`）。
- 票面＝`.scratch/wisp/issues/279-p39-comment-and-the-255-roster-citation-ride-together.md`（20 行／1647 字节，逐字读过；⚠ 派单里给的 slug `279-thirteenth-comment-p39-and-ticket-255-roster-anchor-linkage.md` **盘上不存在**，真名以 `ls | grep '^279'` 现取为准）。
- 取数纪律：一律 `git show HEAD:<path>`／`git grep ... HEAD --`／`git diff`／`git log`，**不拿工作树当证据**；零 `go build`／`go test`／`go vet`／`go list`／`go doc`、零 `scripts/*.sh`、零产码字节、零 push。本件只做 **AC#3 与 AC#2**，**AC#1 一格未碰**（理由见 §4）。

---

## 1. 落点名册：怎么找到的（尺命令＋枚数）

名册**不是**从 `comment-fix-prep-1`（那是料文，13 处的原文/替换文本）取的，而是从**落地腿自己的记录件**取的，再用 git 反查验真伪：

| 尺（命令） | 读数 |
|---|---|
| `ls .scratch/wisp/issues/ \| grep '^279'` | `279-p39-comment-and-the-255-roster-citation-ride-together.md` |
| `find .scratch/wisp/probes -maxdepth 1 -type d -iname '*comment*'` | 四枚：`comment-truth-2`／`comment-fix-prep-1`（料）／`comment-fix-check-1`（复核）／`comment-fix-land-1`（落地） |
| `cat .scratch/wisp/probes/comment-fix-land-1/01-landing-log.md` | 13 行逐处账：**12 行「已贴」＋ P39「未贴」**；commit 清单 `a1b19873`（锚）→`d7a3d62c`(P01-03)→`d51dde33`(P04,P06)→`8dbae00c`(P07-09)→`aec801ef`(P10,P11)→`aa6c1881`(P12,P13) |
| `grep -nE '^## P[0-9]+\|^## P39' .scratch/wisp/probes/comment-fix-prep-1/01-ready-to-apply.md` | **恰 13 段**（`P01 P02 P03 P04 P06 P07 P08 P09 P10 P11 P12 P13 P39`；编号**没有 P05**，不是漏抽） |
| `git log --oneline a1b19873^..HEAD -- <10 目标件＋cmd/wisp/run.go>` | 整族 **6 枚**：上述 5 枚落地提交＋`623a4d81`（编排者 A742 落账，代修 `internal/agent/approval/doc.go:27` 的 `:369`→`:370`，numstat `1 1`，纯注释） |
| `git show --numstat --format= <每枚>` | 逐枚增删行数见 §2 表 |

- **枚数写法（具名）**：12 处落点＝**整族 12 枚**（不是抽样），分布在 **10 个文件**（`config_reload.go` 占 3 处＝P07/P08/P09，其余各 1 处）；产码提交**整族 5 枚**＋编排版代修 1 枚。⚠ 派单里那句「12 处」若被读成「12 个文件」就错了：尺命令 `git show --numstat` 的受改件行数是 **10＋1（doc.go 代修）→ 去重 10**。

12 处落点（HEAD 现量，取自 `comment-fix-land-1/00-anchor.md` 表，逐枚与 `01-ready-to-apply.md` 段号对齐）：

| 处 | 文件（HEAD 现量块位） | 落地提交 |
|---|---|---|
| P01 | `internal/panel/composer_dispatch.go` | `d7a3d62c` |
| P02 | `internal/agent/approval/doc.go` | `d7a3d62c` |
| P03 | `cmd/wisp/panel_inbound.go`（HEAD 现读块 11-15） | `d7a3d62c` |
| P04 | `internal/panel/git.go` | `d51dde33` |
| P06 | `cmd/wisp/config_readers_255.go` | `d51dde33` |
| P07/P08/P09 | `cmd/wisp/config_reload.go`（同文件三块，按 P09→P08→P07 从后往前贴） | `8dbae00c` |
| P10 | `cmd/wisp/firstrun.go`（真身块 7-14，非料文自称 8-15） | `aec801ef` |
| P11 | `cmd/wisp/resident_approval_windows.go` | `aec801ef` |
| P12 | `cmd/wisp/approval_reply.go` | `aa6c1881` |
| P13 | `internal/risk/assessor.go` | `aa6c1881` |
| P39 | `cmd/wisp/run.go`（块 331-336） | **无提交**（未贴，§3） |

---

## 2. AC#3 读数：**「那 12 处全为注释行」今天仍成立**（非注释行改动＝**0**）

### 判据（我自己定义并具名）

- **R0 基数**：`git diff a1b19873^..aa6c1881 -- <10 目标件>` 抽 diff 正文行：`grep -E '^[+-]' \| grep -vE '^(\+\+\+|---)'`（排 `+++`/`---` 文件头，否则每枚文件白送 2 行"非注释"假读数）。
- **R1 注释行判据（主尺）＝整行注释**：`^[-+][[:space:]]*//` —— 符号后只允许 tab/空格，然后必须是 `//`。**块注释中间行不算注释**（本批 10 枚文件里没有 `/* */` 参与，见 R2），**行尾注释不算**（`cfg := rt.cfg // x` 会被判成非注释，这是**故意从严**：真把行尾注释混进产码行一定是产码改动，宁可报红也不放过）。
- **R2 形状探针**：`^[-+][[:space:]]*(/\*|\*/|\*)`（`/*`、`*/`、块注释续行 `*`）。
- **R3 形状探针**：`^[-+][[:space:]]*[^/[:space:]].*//`（符号后先是非斜杠字符、行内出现 `//`＝产码行带尾注释，或字符串里有 `//`）。

### 读数（同一次取数，⛔ 无截断）

| 量法 | R0 正文行 | R1 命中 | **非注释（期望 0）** | R2 | R3 |
|---|---|---|---|---|---|
| 逐枚累计（`a1b19873^..aa6c1881`，10 件）＝12 处批次 | 111 | 111 | **0** | 0 | 0 |
| 到今日 HEAD（`a1b19873^..HEAD`，同 10 件，含 `623a4d81` 代修） | 111 | 111 | **0** | — | — |

- 累计 `+65 / −46` ＝ 111 正文行，与 `comment-fix-land-1/01-landing-log.md` 自报「10 件合计 +65/−46」**逐字对上**。
- 逐枚 numstat（`git show --numstat --format=`）：`d7a3d62c` 5/4＋10/10＋5/3；`d51dde33` 4/3＋9/5；`8dbae00c` 12/8；`aec801ef` 7/5＋3/2；`aa6c1881` 6/4＋4/2 ⇒ 合计 **65/46** ✓。`623a4d81` 对 10 件只碰 `doc.go` **1/1**（注释行，累计读数不变）。
- 代修那笔今天仍在：`internal/agent/approval/doc.go` 的 `:369`→`:370` 引用已在累计 diff 内，且 R1 判它注释行。

### 最容易误判的那一形（具名＋本件怎么排掉的）

**把 raw string 字面量的中间行当成整行 `//` 注释**（Go 的 `` `...` `` 跨行串里没有注释语义，一行内容正好写成 `// xxx` 就会被 R1 误判成注释；次一形是**行尾注释与整行注释混淆**，R3 就是为此单设的）。
排法（本腿实测，非推理）：对 10 枚文件逐枚数「行内反引号数为奇数」的行 ⇒ 只有 `cmd/wisp/panel_inbound.go` 有 **4 行**命中：`:8`/`:9`（**注释里的行内 code span**，跨两行的一对反引号，不是串）与 `:285`/`:289`（真多行 raw string `panelInboundUsage` 的起止）。P03 的改动块在 `:11-15`，两处都不重叠 ⇒ 本批 111 行里**零**行落在 raw string 内部。其余 9 枚文件奇数行＝0。
第二道保险：R3＝0，说明 111 行里没有任何一行是「符号后先有非斜杠字符再出现 `//`」——即没有产码行带尾注释、也没有字符串内含 `//` 被改动。

**结论（AC#3）**：判据 R1 下非注释行增删＝**0**，且 R2/R3 双探针＝0；12 处「全为注释行」今天（`dc65df28`）**成立**。附带读数：`git diff --numstat a1b19873^..HEAD -- cmd/wisp/run.go` **空** ⇒ P39 至今未落地，与 §3 一致。

---

## 3. AC#2 触发条件登记（只登记，⛔ 不落地）

### ① 今天 P39 那块「还没贴」的证据（同一次取数）

- 待改文本在探针件里的位置（逐字现取）：`.scratch/wisp/probes/comment-fix-prep-1/01-ready-to-apply.md`
  - `:327` = `## P39 — cmd/wisp/run.go`
  - `:329` = `1. **目标**：`cmd/wisp/run.go:333`（现量＝`// Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does`；替换块＝**331–336**，行前有 **1 个 tab**）。`
  - `:330` = `2. **原文逐字**（块 331–336，逐字含前导 tab）：`，围栏 `:331` ```` ```go ```` ⇒ 原文六行在 `:332-337`，围栏收在 `:338`
  - `:339` = `3. **替换文本**（换 331–336；6 行 → 9 行，粘贴时保留前导 tab）：`，围栏 `:340` ⇒ 替换九行在 `:341-349`，收在 `:350`
- HEAD 上那一处的现读原文（`git show HEAD:cmd/wisp/run.go \| sed -n '328,338p'`，带真实行号）：

```
331		// modeWrites is ticket 114 AC#2's gate: the one handler a panel mode request
332		// is answered by, assembled with the SAME L2 leg the Store above was handed.
333		// Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does
334		// not exist in this tree (tickets 33/35) - so this is the gate's other half,
335		// not a path: the host that gets wired later cannot assemble a wider档 without
336		// a confirmation leg.
337		modeWrites *panel.ModeWriteHandler
```

- 判别位具名：替换版第 3 行（`:343`）是 `// Nothing reads this field yet: the "event -> ParseComposerRequest" hop now`，HEAD `:333` 仍是 `// Nothing calls it yet - the WebView2 ...` ⇒ **原文形状，未贴**。
- 三条旁证：① `run.go` 在 HEAD 共 **1427** 行（贴后应为 1430）；② `git diff --numstat a1b19873^..HEAD -- cmd/wisp/run.go` 空；③ `git log -1 -- cmd/wisp/run.go` ＝ `613606c0 2026-10-02 09:39`（**早于整个批次**，那之后没人动过 `run.go`）。
- ⇒ **触发条件（票面 `AC#2`／台账 `A742`）今天仍未到来**：`run.go` 自 `613606c0` 起零产码改动，「下一枚因实事而动 `cmd/wisp/run.go` 的腿」还没出现。本腿不动它，也不替编排者翻框。

### ② 「照贴会打红一枚行号锚」的因果关系——**今天仍在**（三处现量 file:line＋短语）

| 处 | HEAD 现量 file:line | 逐字短语 |
|---|---|---|
| 被引行（锚本体） | `cmd/wisp/run.go:991` | `	cfg := rt.cfg`（前导 1 个 tab；函数头 `:990` ＝ `func (rt *agentRuntime) execute(parent context.Context, task string) int {`） |
| 名册逐字引 | `cmd/wisp/config_readers_255.go:110` | `"agent": hotClaimSnapshotOnly + "cmd/wisp/run.go:991 [cfg := rt.cfg] - every cfg.Agent read is the boot copy, so this run keeps acting on the old value",` |
| 核它的用例 | `cmd/wisp/config_receipt_255_test.go:179`（断言体在 `:207-211`） | `:207` `got := strings.TrimSpace(lines[n-1])`／`:208` `if !strings.Contains(got, token) {`／`:209` `t.Errorf("row %q cites %s:%s for %q, that line now reads %q - the evidence drifted, so re-adjudicate this row ...` |

- 仪器形状（具名，非推理）：尺 `evidenceCite`＝`cmd/wisp/config_readers_255.go:229` 的 `([\w./-]+\.go):(\d+) \[([^\]]+)\]`；用例 `:179` 起 `for row, verdict := range hotRowClaims`（`:93` 定义的 map）**只扫 map 值**，逐枚 `os.ReadFile(root/rel)` → `lines[n-1]` → `Contains(token)`。**按行号取，不按短语搜索** ⇒ 行号一漂就红，不是"引用漂移探测器"的比喻而是字面机制。
- **静态推演（⛔ 未跑任何测试）**：P39 是把 `run.go:331-336` 的 6 行换成 9 行 ⇒ 净 **+3**；`331 < 991` ⇒ 其后每行整体 **+3**：`cfg := rt.cfg` 从 `:991` 漂到 **`:994`**。名册 `:110` 仍写 `run.go:991` ⇒ 尺读到的第 991 行＝原第 988 行，现读逐字 `// so a single task can be cancelled without touching the entry that produced it.`，不含 `cfg := rt.cfg` ⇒ `:208` 条件成立 ⇒ **那枚用例真的会红**（`t.Errorf`，rc=1）。旁及：`citedRowsFloor = 8`（`:220-222`）**不受影响**（`citesChecked++` 在 `Errorf` 之后仍执行），所以红只会是"漂移"这一枚，不会伪装成"引用数不足"。
- **同批必须一起改的不止那一枚**（具名，票面未提）：`config_readers_255.go` 的**注释散文里还有四枚 `run.go` 行号引用**（`:100/:106/:107/:108`），尺**不扫注释**⇒ 贴完不会红、但会**静默过期 3 行**：`:100` `// (cmd/wisp/run.go:435 [res := llm.NewResolver(cfg, st)]), so "已立即生效"`、`:106` `// at assembly (cmd/wisp/run.go:424 [rt.cfg = cfg]) and every use reads that`、`:107` `// copy (cmd/wisp/run.go:991 [cfg := rt.cfg]; the field is consumed at`、`:108` `// cmd/wisp/run.go:1014 [cfg.Agent.PerToolTimeoutMS]). plan() did overwrite`（HEAD 现读 `run.go:1014` ＝ `			PerToolTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS) *`，漂后应指 `:1017`）。⇒ 随行落地时四处（`:100/:106/:107/:108`）＋名册 `:110` 一起 +3，⛔ 别只改被钉住的那一枚。
- 本项**判不动的欠账**：无（三处 file:line＋短语＋漂移后行号均已现量到；唯一未证的是 rc，按派单禁跑测试，此处只要静态结论）。

---

## 4. AC#1 欠账具名登记（**按住的理由是派单约束，不是本腿失败**）

`AC#1 同笔联动`＝三发，本腿**一发未做**：

| # | 欠的那一发 | 为什么欠 |
|---|---|---|
| 1 | 贴 P39（`cmd/wisp/run.go:331-336` → 料文 `:341-349` 九行，净 +3） | 三发里第 2 发要跑 `go test`，此刻 `285-r1` 独占 Go 编译/测试面与本包突变面（本仓串行铁律：同一时刻只许一枚腿在 Go 导入图里动手） |
| 2 | 改名册引用 `config_readers_255.go:110` 的行号 `991`→`994`（`[cfg := rt.cfg]` 原文一字不动），并按 §3② 把 `:100/:106/:107/:108` 四处散文引用同批 +3 | 同上：改完必须靠那枚尺证明没放宽 |
| 3 | 复跑 `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 至 **rc=0** | 需要 `go test`，本腿禁跑；且落票面要求「随行 commit 具名」，本腿零产码改动 |

- ⛔ 本腿对票面/台账/产码**一字未动**（写面只有本件）。翻勾与 `AC#2` 落地登记归编排者。
- 禁区照抄确认：不许为变绿放宽那枚 roster 钉子；不许把 P39 标成"实际没改却已改"（本件 §3① 三条旁证就是防这一条的）；三枚冻结件一字不动；零 push。

---

## 5. 具名顶回（票面／派单与盘上原文不符处；⛔ 未迁就）

1. **派单给的工单 slug 不存在**：派单写 `279-thirteenth-comment-p39-and-ticket-255-roster-anchor-linkage.md`；盘上真名 `279-p39-comment-and-the-255-roster-citation-ride-together.md`。
2. **票面 `config_readers_255.go:109` 已漂**：HEAD 现量名册 `"agent"` 行在 **`:110`**（同文件 P06 +1 行所致；`comment-fix-land-1/00-anchor.md` 与 `01-landing-log.md` 早已预告这一漂）。串内容与逐字短语**未变**。
3. **票面 `run.go:991`（`cfg := rt.cfg`）今天仍准**，未漂（现量＝`\tcfg := rt.cfg`）。`config_receipt_255_test.go:179` 也仍准（＝用例函数声明行），但**断言本体在 `:207-211`**，票面把函数行当断言行指，按本件为准。
4. **票面「现量」节说待改文本在 `01-ready-to-apply.md` 的 P39 那一块**——对，但**段内行号**票面未给：本件给全（`:327` 段首／`:332-337` 原文／`:341-349` 替换）。
5. **「12 处」不等于「12 个文件」**：名册 12 处落在 **10 个文件**（`config_reload.go` 三处）；编号族**无 P05**（P01-P04＋P06-P13＋P39＝13）。派单让我"或从 `A7xx`/`A733`/`A742` 具名的 commit 找"——`A733` 的交件 sha 是 `a1f0d2ff`/`e7cd0c6e`（**复核件**，不含落地 sha），落地 sha 在 `A742`/`01-landing-log.md`（`a1b19873`→`e5b46742` 七笔）；`e5b46742` 是记录件一笔，产码只有 5 枚＋编排者代修 1 枚（`623a4d81`，**纯注释、票面未列**，已并入 §2 读数）。
6. 派单说「12 处已贴、第 13 处被按住」✓ 与盘上一致；但**盘上还有第 14 笔注释改动**（`623a4d81` 代修 `doc.go:27` 引用），它不在 12 也不在 P39 之列——AC#3 若按"累计到 HEAD 的 10 件"量，它已被计入且仍为 0 非注释行。
7. `A742` 自陈「我复跑三把尺」用的是 `grep -vcE '^[-+]\s*//'` 一把尺；本腿**独立**复现同一结论，并另加 R2/R3 两把形状探针＋raw string 重叠检查（`A742` 未做）。结论一致：非注释行＝0。
