# 票 99 对抗验收裁决表 — `scripts/d22scan.sh` 第一步（阳性对照）被 `go test` 缓存端过去

验收方：`acceptor-ticket99`（独立对抗验收代理，未参与实现，未看先前对话）。
登记时间：`date` 实测 `Mon Sep 21 18:27:42 CST 2026`（会话开始 `18:12:12 CST 2026`）。
仪器：`go version go1.27.1 windows/amd64`；`GOCACHE=C:\Users\swq\AppData\Local\go-build`（本机共享缓存）。
快照全部在 **仓外** `/tmp`，`git archive <sha> | tar -x`，仓内未建 worktree、未 checkout：

| 用途 | 目录 | 导出的 sha | 第一步实现 |
| --- | --- | --- | --- |
| AC#1（复现缓存回放） | `/tmp/wisp-99-ac1-ac99` | `d0d8782^` = `360efdf` | 改前：裸 `go test ./...` |
| AC#3（双向变异） | `/tmp/wisp-99-mut-ac99` | 建快照时 HEAD = `06906f7` | 改后：`runtests.sh` |
| AC#4（纯净门禁两跑） | `/tmp/wisp-99-gate-ac99` | 建快照时 HEAD = `06906f7` | 改后：`runtests.sh` |

⚠ 本报告的**每一枚绿/红都来自上表 /tmp 快照**；**主树（共享工作树）未跑门禁**——共树在飞的有
`internal/config`（票 95）、`scripts/portable-tests.sh`+`ci.yml`（票 93）的未提交改动，
那棵树里的红不是本票的账（详见 R-99-5）。

来源档位：**A=独立复现**（我自己敲的命令、自己的 rc）；**B=日志＋归档，我抽验**（文件/commit message 我读过并局部核过）；
**C=仅自述，不背书**（前实现代理自己的叙述，我无从复跑）。每行末标注。

---

## AC#1 — "先把缓存真能端过去量出来，不许只论证"

票面原句（`.scratch/wisp/issues/99-*.md:25-28`）：
> **AC#1** 先把"缓存真能端过去"**量出来**，不许只论证：在**仓外快照**里 ① 跑一次脚本让它绿；② 往树里种一个已知违规；
> ③ **再跑一次** ⇒ 报第二次的 rc 与结论原文。若第二次仍是绿 ⇒ 缺陷成立并记下这条序列；若第二次就红 ⇒ **如实写"未复现"**并说明差异。

我亲自敲的命令（`cd /tmp/wisp-99-ac1-ac99`，脚本 = `d0d8782^` 原样，`grep -n '^go test \./\.\.\.$' scripts/d22scan.sh` → `32:go test ./...`）：

| 腿 | 命令 | 真实 rc | 原文行 | 档 |
| --- | --- | --- | --- | --- |
| ① 跑一次让它绿 | `sh scripts/d22scan.sh` | rc=**0** | 第 2 行 `ok  	github.com/CarlosShao/wisp/tools/d22scan	10.352s`（真跑，无 `(cached)`，把缓存喂热） | A |
| ② 种已知违规 | `printf 'export function decide() { return approval.decide({ allow: true }); }\n' > frontend/src/app.js` | rc=0 | `grep -n . frontend/src/app.js` → `1:export function decide() { return approval.decide({ allow: true }); }` | A |
| ③ 再跑（整脚本） | `sh scripts/d22scan.sh` | rc=**1** | 第 2 行 `ok  	github.com/CarlosShao/wisp/tools/d22scan	(cached)`；第 13 行 `frontend/src/app.js:1: [panel-approval] ...`；第 14 行 `d22scan: 1 finding(s); D22 bans are not negotiable` | A |
| ③′ 再跑（**只**第一步＝阳性对照本体） | `cd tools/d22scan && go test ./...` | rc=**0** | 唯一一行 `ok  	github.com/CarlosShao/wisp/tools/d22scan	(cached)` | A |

**结论：通过（缺陷按票面形状量出来了，判定分两级，不是一句"复现"糊过去）。**
- **步骤级 = 复现**：种了真违规之后，承载"证明门禁能红"的那一步端的是 `(cached)` 的 `ok`、rc=0 ——
  `TestScannerSelfScanOfRealRepoIsGreen` 在测试**运行时**读仓库树，缓存判定里只有构建输入，它没再看那棵树一眼。
- **脚本级 = 未复现**（如实写，票面 AC#1 末句要求的情形）：整脚本第二次 rc=**1**，因为第二步 `go run . -root` 不受测试缓存管辖，
  这次违规是**它**抓的。⇒ 危害的准确形状：**回放的是阳性对照，兜住的是扫描**；判据不能建在"第二步恰好还在、恰好没命中缓存"上。
- 差异说明（票面要求的"Go 版本/缓存策略"）：`go1.27.1`，测试缓存键**包含运行目录**，
  所以我新建的每个 `/tmp` 快照第一次都是真跑（`10.352s` / `9.492s`），**同目录**内第二次立刻回放。
  ⇒ 这条回放只咬"同一个工作树里反复自检"这个形状（正是 A64② 要求每个代理收尾都跑本脚本的那条路），不是跨 checkout 全局回放。见 R-99-3。 A

---

## AC#2 — "脚本第一步与 `runtests.sh` 对齐成同一档，不许顺手改覆盖面"

票面原句（`:29-31`）：
> **AC#2** 修法：脚本第一步与 `runtests.sh` **对齐成同一档**（显式 `-count=1`，或干脆改调 `runtests.sh`，两处调用不再各持一套规矩）。
> ⚠ **不许**用 `GOFLAGS=-count=1` 之外的方式"全局加快速"，也**不许**顺手给脚本加 `-failfast`/`-short`/`-run` 之类**改变覆盖面**的开关。

改动只有一枚 commit 的一个文件（B）：`git show --numstat d0d8782` → `23 4 scripts/d22scan.sh`，全仓无第二个文件。
承载行为那一行我自己定位（行号会漂）：`grep -n '^sh "\$root/tools/d22scan/runtests\.sh" -C tools/d22scan \./\.\.\.$' scripts/d22scan.sh`
→ `51:sh "$root/tools/d22scan/runtests.sh" -C tools/d22scan ./...`；其对齐目标是 `tools/d22scan/runtests.sh:75` 的 `go test -v -count=1 "$@"`。A

### 本票最大可疑点的专门核查：覆盖面有没有被这次替换悄悄改变（三问）

**① 测试集合与旧那一步完全相同吗？——相同。** A
- 旧：`cd "$root/tools/d22scan"`（第 48 行，本次未改，`git show d0d8782` 里它是 context 行）→ `go test ./...`。
- 新：同一个 `cd` 之后 → `sh "$root/tools/d22scan/runtests.sh" -C tools/d22scan ./...`，`runtests.sh:41-42` 用自身位置
  推导 `root=$here/../..`，`-C tools/d22scan` 解析成 `$root/tools/d22scan`（`runtests.sh:53-56`）⇒ **同一个 module 根目录**。
- 分母我自己数：`(cd tools/d22scan && go list ./...)` → 只有 `github.com/CarlosShao/wisp/tools/d22scan` **1 个包**；
  仓里另有 `go.mod` 的 `scripts/spike`、`tools/mockllm` 在**改前改后都不在分母里**（旧那一步本来也已经 `cd` 进 tools/d22scan）。
  ⇒ 没有"少一个包 = 把被检查对象从门禁里删掉"。用例数：改后日志 `PASS=21 ... === RUN=31`；
  我在变异树上单独跑裸 `go test ./...`（冷缓存真跑）复现的是**同两个测试名**红 ⇒ 同一集合，B/A 双向一致。

**② 新加的"SKIP 不算过""PASS=0 就红"是加严还是顺手改语义？——是加严，覆盖面语义未动。** A（读 `runtests.sh:85-109`）
- 它只可能把 rc 从 0 变 1（`skipped!=0` → exit 1；`passed==0 && failed==0` → exit 1；`go test` 自己的 rc 原样透传），
  没有任何路径能把红洗成绿；包参数 `$@` 只有 `./...`，**没有** `-failfast`/`-short`/`-run`，也没有 `GOFLAGS` 注入。
- 唯一继承来的新风险（登记，不算缺陷）：本地门禁从此和 CI 一样吃"平台性 SKIP = 红"这条票 71 规矩，
  `tools/d22scan` 里将来若有 `_test.go` 在 Windows 上 `t.Skip`，第一步会红且原因不是禁令。补救动作：那是**正确**的红，
  按 `runtests.sh:100` 的提示改步骤 scope，**不许**回头放宽脚本。A/B

**③ `-C` 与 `$root` 推导在从任意目录调用时成立吗？——成立，我造了三种 cwd。** A

| 调用形状 | 命令 | rc |
| --- | --- | --- |
| 仓外绝对路径 | `cd /tmp && sh /tmp/wisp-99-gate-ac99/scripts/d22scan.sh` | **0**（末行 `d22scan: clean - no D22 ban violations; live scope work: ...`） |
| 树深处相对路径 | `cd /tmp/wisp-99-gate-ac99/frontend/src && sh ../../scripts/d22scan.sh` | **0**（同一行 clean 末行） |
| 裸调 `runtests.sh -C` | `cd /tmp && sh .../tools/d22scan/runtests.sh -C tools/d22scan ./...` | **0**，`runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0` |

`runtests.sh:58-62` 的 `go.mod` 守卫会把写错的 `-C` 直接判 exit 2（"a wrong -C is how a green no-op is born"），
所以"从任意目录调用"最坏是红，不会假绿。A

**结论：通过（修法与票面 AC#2 的第一选项逐字一致，覆盖面三问全数），但附带一条本票账外的条件：见 AC#2-④ / R-99-4。**

### ④ CI 那侧：现在同一台仪器在 CI 里跑了两遍

`.github/workflows/ci.yml` 我自己读（B→A 静态核对）：`:48` `bash tools/d22scan/runtests.sh -C tools/d22scan ./...`；
`:68` `run: sh scripts/d22scan.sh`，而该脚本第一步（现 `scripts/d22scan.sh:51`）就是**逐字同一条命令**
⇒ CI 每次多花一整遍 `tools/d22scan` 的 `-v -count=1` 测试（本机实测这一步 9.5s 冷跑 / 21 用例 / 31 个 `=== RUN`）。A
**判账**：这不是票 99 的缺陷——AC#2 正是**要求**两处并成同一台仪器，重复是票面判据的直接后果；
`ci.yml` 在票面上明确划归票 85（`:7` "改法归票 85"），且共树里 `ci.yml` 正被票 93 改着（`git status` 显示 `M .github/workflows/ci.yml`）。
⇒ 登记给票 **85/93**（R-99-4）。**"合一"的判据**：合一 = 把 CI 的**两步并成一步**
（保留 `scripts/d22scan.sh` 作为唯一入口、删/改 `:48` 那一步的重复调用，或让脚本第一步成为 CI 唯一的阳性对照），
**不是**把 `scripts/d22scan.sh` 改回裸 `go test ./...`（那是把回放绿请回来），也**不是**给脚本加 `-run`/`-short` 之类缩小分母的开关。B（票面 `next=` 与我的一致，我独立核对）

---

## AC#3 — 变异双向

票面原句（`:32-35`）：
> **AC#3** 变异双向：(i) 把我加的 `-count=1` 去掉 ⇒ AC#1 那条序列必须重新出现"种了违规还绿"；
> (ii) 阳性对照：种一个 `frontend/` 里的 `approval.decide` ⇒ 脚本 rc=1 并点名……锚点=承载行为那一行，同链 grep 证落地，
> 还原后 `git diff --quiet` 证干净，**编译失败不算变异**。

`-count=1` 由 `tools/d22scan/runtests.sh:75` 承载，而 `tools/d22scan/**` 本票禁改 ⇒ 变异只能做在**调用行**上（`scripts/d22scan.sh:51`），
我按"退回旧实现 = 换回裸 `go test ./...`"来做，这是票面 AC#3(i) 的等价形。A

### (i) 退回旧实现 ⇒ "种了违规还绿"回来了（快照 `/tmp/wisp-99-mut-ac99`）

落地证明（同一条 `&&` 链里 `grep -n` 打印被改后的整行，不是 `sed` 静默不匹配）：

```
$ sed -i 's|^sh .*runtests\.sh.* -C tools/d22scan \./\.\.\.|go test ./...|' scripts/d22scan.sh \
  && grep -n '^go test \./\.\.\.$' scripts/d22scan.sh && bash -n scripts/d22scan.sh && echo MUTANT-LANDED
51:go test ./...
MUTANT-LANDED
```

| 步 | 命令 | rc | 原文行 |
| --- | --- | --- | --- |
| 干净树喂热缓存 | `sh scripts/d22scan.sh` | **0** | 第 2 行 `ok  	github.com/CarlosShao/wisp/tools/d22scan	9.492s`（真跑） |
| 种违规 | `printf ... > frontend/src/app.js` | 0 | `1:export function decide() { return approval.decide({ allow: true }); }` |
| 再跑整脚本 | `sh scripts/d22scan.sh` | 1 | 第 2 行 `ok  	github.com/CarlosShao/wisp/tools/d22scan	(cached)`；红只来自第 14 行 `d22scan: 1 finding(s); D22 bans are not negotiable` |
| **再跑第一步本体** | `cd tools/d22scan && go test ./...` | **0** | `ok  	github.com/CarlosShao/wisp/tools/d22scan	(cached)` —— 违规在树上，对照却盖章通过 |

⇒ **AC#3(i) 通过**：把承载去掉，AC#1 那条序列**逐字回来**（且是在**修法曾存在的同一棵树**上做掉的，不是我另找一棵树）。A

⚠ 方法学登记（假绿抓我自己）：我头**两次** `sed` 都**静默不匹配**——模式 `^sh .root/...` 少算了 `"$root` 里的那个 `$`。
第一发被同链 `grep -n '^go test \./\.\.\.$'` 无输出打断 `&&` 才发现（那次的 run1 日志第 2 行是 `=== RUN TestScanDetectsAllSeededViolations`
= `-v` 输出 = 脚本**根本没被改**，读数已作废重跑）。⇒ 这正说明票面"锚点必须同链 grep 打印整行"是硬规矩，`sed` 单独跑就是一次假绿。A

### (ii) 修法在场 + 种 `frontend/` 的 `approval.decide` ⇒ 这条调用路径自己红

同一快照还原后（`diff -q scripts/d22scan.sh <主树>` rc=**0** 逐字节相同），只种同一个违规文件，不碰任何源码/测试：

```
$ printf 'export function decide() ...' > frontend/src/app.js && sh scripts/d22scan.sh ; echo rc=$?
AC3ii-script-rc=1
25:--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (0.47s)
72:--- FAIL: TestRealRepoLedgerIsHonest (0.49s)
24:    scan_test.go:269: repo HEAD violates: frontend/src/app.js:1: [panel-approval] `approval.decide` in frontend/ is banned (D33/F2: allow decisions are native-side only)
153:runtests.sh: go test exited 1 - packages=[./...] top-level: PASS=19 FAIL=2 SKIP=0, === RUN=31, '[no tests to run]'=0
```

- 红**来自第一步**：`grep -n 'd22scan.sh: scan of'` 在该日志里**零命中** ⇒ `set -eu` 让脚本在阳性对照就死了，第二步根本没跑
  ⇒ 证明的是**这条调用路径有牙**，不是扫描器本体有牙（本体有牙是票 88 的账，我没借它的口）。A
- 数 `=== RUN` 是因为非 `-v` 的 `go test` 看不到 PASS：本日志 `=== RUN`=31、`--- PASS`=19、`--- FAIL`=2、`--- SKIP`=0、`(cached)`=0。A
- 还原：`cp` 主树脚本后 `diff -q` rc=0；种子文件 `rm -f` 后 `ls` → `No such file or directory`；
  主树 `git status --porcelain -- scripts/d22scan.sh tools/d22scan allowlist.txt` 为空 ⇒ 我没有覆盖别人的东西。A
  （票面"还原后 `git diff --quiet` 证干净"我在快照里做，快照无 `.git`，等价证据是 `diff -q` 与主树逐字节相同。B）

**结论：通过（双向我自己都跑出来了）。**

---

## AC#4 — 门禁

票面原句（`:36-39`）：
> **AC#4** 门禁：`bash -n scripts/d22scan.sh` rc=0；在纯净快照 `sh scripts/d22scan.sh` 连跑**两次**都 rc=0
> （⚠ **HEAD 上现在它是 rc=1**，唯一命中是 `internal/winsec/winsec.go:126`，那是**票 94** 的账 ⇒ 若你开工时它还没修完，
> 你的判据是"**除那一条已知命中之外没有新命中**"……）**收尾前必跑 `sh scripts/d22scan.sh`**（A64②）。

| 判据 | 我敲的命令 | 真实读数 | 档 |
| --- | --- | --- | --- |
| 语法 | `bash -n scripts/d22scan.sh`（纯净快照 + **主树**各一次） | rc=**0** / rc=**0** | A |
| 纯净快照连跑两次 | `cd /tmp/wisp-99-gate-ac99 && sh scripts/d22scan.sh`（两次独立跑） | run1 rc=**0**、run2 rc=**0**，两次末行都是 `d22scan: clean - no D22 ban violations` | A |
| `(cached)` 自己数 | `grep -c 'cached' /tmp/ac4-run1.log /tmp/ac4-run2.log` | **0** 和 **0**（不是引用票面数字） | A |
| 第一步是真跑 | `runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0`（两次都有） | 21 用例真跑，SKIP=0 ⇒ 不是"SKIP 当过" | A |
| 台账 8 行 scope 齐全 | `awk '/d22scan.sh: scan of/,0' /tmp/ac4-run2.log \| grep -c 'd22scan: scope'` | **8**：`bans #1-5 internal/`=197、`bans #1-5 cmd/`=20、`ban #6 frontend/`=37、`ban #7 internal/tools/`=17、`ban #8 design/`=16、`ban #8 frontend/`=37、`ban #8 internal/`=342、`ban #8 cmd/`=26 | A |
| `ban #6/#8 frontend/` 不降 | 同上，对比票面 AC#4 报的快照读数 37/37 | **37 / 37，未降**；`ban #8 internal/` 340→**342**（票 103 的 `internal/winsec` 新文件进了树，只增不减，非本票造成的） | A |

**结论：数字全达标 = 通过；但票面判据本身有毛病**（"HEAD 上现在它是 rc=1，唯一命中 `internal/winsec/winsec.go:126`" 是**腐坏的引用**：
我在 `d0d8782^` 与 `06906f7` 两棵纯净树上都是 rc=**0**，`internal/winsec/` 我一行没读没改）
⇒ AC#4 的"除那一条已知命中之外没有新命中"这半句在今天的 HEAD 上没有对象，登记 R-99-2，判据应按实测 rc=0 写。A

---

## commit 归属核对（本票来历特殊，专门核）

| sha | 时间 | 实际内容（`--numstat`） | commit message 归属写法 | 核对结论 |
| --- | --- | --- | --- | --- |
| `d0d8782` | 18:09:54 | `23 4 scripts/d22scan.sh`（**代码修法本体**） | `fix(99,AC#2)`，正文首句"**由编排者代 agent-ticket99 落档**（它撞轮数上限死亡，改动留在共享工作树未提交…）"，并明写"AC#3 双向变异与 AC#4 的票面读数**仍缺**，由 agent-ticket99b 从断点续" | **代码是编排者代落的**，message 写得清楚，无冒名 |
| `9e00629` | 18:10:46 | `99 5 .scratch/wisp/issues/99-*.md`（**只改票面**） | 标题 `fix(99,AC#1-#4)`，全篇第一人称复述"修法：第一步从裸 `go test ./...` 改成 …"（= `d0d8782` 的内容），并交回 AC#1-#4 全部读数 | **有记账问题（R-99-1）**：见下 |
| `06906f7` | 18:13:19 | `13 0 .scratch/wisp/issues/99-*.md`（票面再追加，**无代码**；= 我建快照时的 HEAD） | 标题 `docs(99,交件后登记)`："我的脚本改动在我提交前 52 秒被 `d0d8782` 从未提交工作树带走"、"账分两枚：代码在 `d0d8782`，票面在 `9e00629`" | **拆分已被作者自己披露**（不是隐瞒），但"我的脚本改动"与 `d0d8782` 的"由编排者代 agent-ticket99 落档"两种口径并存 |

票面第 61 行的交件读数抬头写的是 `（agent-ticket99，2026-09-21 18:05 实测 date）`、第 3 行 Status 写"18:05 agent-ticket99 交件：四框全勾"，
但 `d0d8782` 的 message 自己声明 18:09 时 **AC#3/AC#4 读数仍缺**，而四框是在 18:10:46 的 `9e00629` 里勾满的
⇒ 那一整段里 **AC#3/AC#4 的读数实际出自 agent-ticket99b 的会话**（快照路径 `/tmp/wisp-t99b-agent-ticket99` 也印证是 b 的目录），
却整段署在 `agent-ticket99` 名下与 18:05 时间戳下；同时票面 line 126 的"**我没提交**"与仓库现状矛盾（修法早就落了，是编排者代落的）。
`06906f7`（99b 交件后自登记）把"代码在 `d0d8782`／票面在 `9e00629`"这笔**账分两枚**披露了 ⇒ **不是隐瞒**；
仍欠两笔：**(a)** 交件读数抬头与 Status 把 99b 测的 AC#3/AC#4 仍冠在 `agent-ticket99 @18:05` 名下；
**(b)** 谁是 `scripts/d22scan.sh` 那次改动的作者，`d0d8782`（编排者代 agent-ticket99 落档）与 `06906f7`（"我的脚本改动"）两种口径并存、没有唯一定稿。
⇒ 这不是读数造假（我逐条复现出来了），是**归属/时间戳记账混淆**。
**更正动作**：我在票面文末追加段逐枚点名 sha 与实际做事的会话（append-only，见 Progress log），不改写 99b 已写的正文，不改 Status、不改文件名；唯一定稿归编排者。A/B

---

## 登记（不在验收期动任何代码/契约）

- **R-99-1（记账，需更正，已在本报告与票面文末落档）**：`9e00629` 把编排者代落的 `d0d8782`（AC#2 修法本体）与前会话未完成的
  AC#3/AC#4 读数统一署为 `agent-ticket99 @18:05`，并留下"我没提交"这句与仓库状态矛盾的残留；Status 行同。补救：票面文末追加段
  按 sha 拆清三方（编排者代落代码 / 99b 交读数 / 99 未提交即亡）；`-done` 归编排者时把这段并进台账。A
- **R-99-2（判据本身有毛病）**：票面 AC#4 括号里"HEAD 上现在它是 rc=1，唯一命中 `internal/winsec/winsec.go:126`"已腐坏
  （两棵纯净树实测 rc=0）。补救：编排者在票面/台账把该句标为已过期，后续票的门禁判据一律写"纯净快照 rc=0 且 `(cached)`=0 且 8 行 scope 齐全"。A
- **R-99-3（危害形状要加限定语）**：`go test` 测试缓存键含**运行目录**——我的独立证据是不同 `/tmp` 快照之间互不回放
  （新目录第一步必然真跑 `10.352s`/`9.492s`），**同一目录**第二次立刻 `(cached)`。票面与脚本注释现在的写法（"加个违规 go test 就可能回放"）
  会让人误以为跨 checkout 也回放。补救：把限定语"同一工作树内反复自检才会回放（也就是 A64② 那条路）"补进 `scripts/d22scan.sh` 头部注释或票面；
  我**没动脚本**（禁改范围外也不在验收期顺手改）。A
- **R-99-4（账归票 85/93，非本票缺陷）**：`ci.yml:48` 与 `:68` 第一步现在是同一条命令 ⇒ CI 每轮多跑一整遍 `tools/d22scan` 的 `-v` 测试。
  判据见上文 AC#2-④：**合两步，不改回脚本**。B/A（静态核对，我未跑 CI）
- **R-99-5（共树噪声，登记不判账）**：我开工时 HEAD=`0717bf2`，建快照时已漂到 `06906f7`（票 103/95 在连续落档），
  主树 `git status` 里 `.github/workflows/ci.yml` 与 `internal/config/*` 是别人的未提交改动、`scripts/portable-tests.sh` 未跟踪。
  ⇒ 本报告所有绿都点名来自 `/tmp` 快照；**我没在共享树跑门禁**，也没跑整仓门禁（本机 `GOOS=linux go vet ./...` 常红，与本票无关）。A
- **R-99-6（仪器方法坑，给下一个验收代理）**：`sed` 对含 `$` 的锚点静默不匹配，两发变异差点以未改动的脚本充数；
  唯一防线是同一条 `&&` 链里的 `grep -n` 整行打印 + `bash -n`。票面 AC#3 已含此规矩，本报告把它变成了一次实测案例。A

## 总判

**通过（附条件）**。票面四框我逐条独立复现：AC#1 步骤级 `(cached)` rc=0 复现、脚本级如实未复现（rc=1）；
AC#2 修法与票面第一选项逐字一致，覆盖面三问全数（同 module / 同 1 包分母 / 只加严不改语义 / 任意 cwd 成立）；
AC#3 双向变异都在**同一棵树**上落地并逐行证明，红名 `TestScannerSelfScanOfRealRepoIsGreen` 出自第一步本身；
AC#4 `bash -n` rc=0、纯净快照两跑 rc=0/0、`(cached)` **0** 次、8 行 scope 齐全、`ban #6/#8 frontend/`=37/37 不降。
条件：R-99-1 的记账更正（我已随本票文末追加段落档）、R-99-2 判据过期、R-99-4 合 CI 两步归票 85/93。
`allowlist.txt` 与 `tools/d22scan/**` 在本次两枚 commit 里未被触碰（`git show --numstat` 只有两个文件，见下表），我没有"顺手修好"任何东西。
