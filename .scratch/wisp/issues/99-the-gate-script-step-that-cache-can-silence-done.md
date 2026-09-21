# 99 — `scripts/d22scan.sh` 第一步是**裸 `go test ./...`** ⇒ 本机测试缓存能把一次真违规"端"成绿

**Status:** ready-for-review（原 open：2026-09-21 16:0x 编排者建；来源=`acceptor-ticket88` 交件时点名的残留，它明写"非本票账"。2026-09-21 18:05 agent-ticket99 交件：四框全勾，读数见文末"交件读数"）
→ **accepted-done**（2026-09-21 18:3x 由 `acceptor-ticket99` 判**通过（附条件）**，裁决表
   `docs/evidence/s1/99-adversarial-acceptance.md`（229 行、删除列 0），交件 commit `4b4cfc9`，快照全在 `/tmp` 仓外、**主树未跑门禁**）。
   ⚠ **结案必须先了清一枚账**（这是编排者自己造成的，见台账 `A77③`：我把这张票**双派**了一次）：
   AC#1/AC#2 归 `agent-ticket99`（它先交 `1f8d212` 的 checkpoint），**脚本代码由编排者代落档 = `d0d8782`**，
   **AC#3/AC#4 有两个独立会话各自跑过一遍** —— `agent-ticket99`（交件 commit `06906f7`，PASS=19/FAIL=2）
   与续跑代理 `agent-ticket99b`（发现前任已死亡后补做，票面 commit `9e00629`，被叫停、未 push）。
   ⇒ 两遍读数**彼此一致**（同一红名 `TestScannerSelfScanOfRealRepoIsGreen`、同一 `(cached)` 原文），
   所以本票的判据不是"某人说了算"而是**被跑过两次**；代价是票面抬头把两组读数冠在了一个名字下 ——
   **验收代理已把这处记账问题登记为 `R-99-1`，更正以本行为准**（本仓规矩：保留原文、文末追加，不覆盖）。
   ⚠ 验收另外教了一条**判据要加限定语**的点（`R-99-3`）：Go 的测试缓存键**包含运行所在目录**，
   所以准确说法不是"缓存会端过去"，而是"**同一目录内**、包与测试源码没变时，`(cached)` 会端过去"——
   这句话在本票的每一次复现里都成立（改前脚本干净树 run1 rc=0 `ok 10.352s` → 种违规 → run2 `(cached)`）。
   覆盖面三问（AC#2 是否偷偷换了仪器）验收亲答：**同 module 同 `cd`、`go list ./...` 分母同为 1 包；
   "SKIP 不算过 / PASS=0 就红"是纯加严（rc 只 0→1）；任意 cwd 三发实测 rc=0**。
   `R-99-4`（CI 的 `:48` 与 `:68` 因此重复跑同一台仪器）**已转票 85**，判据写死"合两步、不许把脚本改回裸 `go test`"。
**Type:** 门禁完整性（票 71 / 票 96 同族：**仪器自己会偷偷不跑**）
**Blocks:** nothing · **Blocked by:** nothing
**Packages:** `scripts/d22scan.sh`（就那个 `-count=1` 的差别）。**禁改**：`tools/d22scan/**`（扫描器本体）、
              `allowlist.txt`（**5 行非注释，只许变短**）、`.github/workflows/ci.yml`（改法归票 85）、
              任何 ban 的文本、任何测试的断言。

## 现场

票 88 的验收代理实测（它的话："另登记 `scripts/d22scan.sh` 第一步裸 `go test ./...` 可被本机缓存端过去"）：
脚本第一步跑 `tools/d22scan` 的**台账/阳性对照测试**，而 `tools/d22scan/runtests.sh:75` 那一侧是**强制 `-count=1`** 的。
⇒ 同一个仓里两套调用方式，一处严一处松（这与票 93 的"SKIP 一侧拒、portable 不拒"是**同一种结构缺陷**）。

**为什么这不是"理论风险"**：`TestScannerSelfScanOfRealRepoIsGreen` 是在**测试运行时**读仓库树的。
Go 的测试缓存**看不见运行时读的文件**——它在缓存判定里只记 import/构建输入。
⇒ 往 `internal/` 或 `frontend/` 加一个真违规（比如一个 `approval.decide`），
若**包与测试源码本身没变**，`go test` 完全可能直接回放上一次的 `ok`。
CI 是干净 runner 所以每轮真跑；**本地开发与代理自检拿到的可能就是回放**——
而我们现在要求每张票收尾前都跑这个脚本（A64② 新立的规矩），**等于让一条会缓存的步骤承担门禁职责**。

## AC（1:1，裁决表 `docs/evidence/s1/99-*.md` 由验收方出）

- [x] **AC#1** 先把"缓存真能端过去"**量出来**，不许只论证：在**仓外快照**里
      ① 跑一次脚本让它绿；② 往树里种一个已知违规；③ **再跑一次** ⇒ 报第二次的 rc 与结论原文。
      若第二次仍是绿 ⇒ 缺陷成立并记下这条序列；若第二次就红 ⇒ **如实写"未复现"**并说明差异（Go 版本/缓存策略），
      然后**仍然**做 AC#2——因为判据不能依赖"缓存恰好没命中"。
- [x] **AC#2** 修法：脚本第一步与 `runtests.sh` **对齐成同一档**（显式 `-count=1`，或干脆改调 `runtests.sh`，
      两处调用不再各持一套规矩）。⚠ **不许**用 `GOFLAGS=-count=1` 之外的方式"全局加快速"，
      也**不许**顺手给脚本加 `-failfast`/`-short`/`-run` 之类**改变覆盖面**的开关。
- [x] **AC#3** 变异双向：(i) 把我加的 `-count=1` 去掉 ⇒ AC#1 那条序列必须重新出现"种了违规还绿"；
      (ii) 阳性对照：种一个 `frontend/` 里的 `approval.decide` ⇒ 脚本 rc=1 并点名（票 88 已证明扫描器本身有牙，
      本票要证的是**这条调用路径**也有牙）。锚点=承载行为那一行，同链 grep 证落地，还原后 `git diff --quiet` 证干净，
      **编译失败不算变异**。
- [x] **AC#4** 门禁：`bash -n scripts/d22scan.sh` rc=0；在纯净快照 `sh scripts/d22scan.sh` 连跑**两次**都 rc=0
      （⚠ **HEAD 上现在它是 rc=1**，唯一命中是 `internal/winsec/winsec.go:126`，那是**票 94** 的账 ⇒
      若你开工时它还没修完，你的判据是"**除那一条已知命中之外没有新命中**"，并如实登记，**不要去碰 `internal/winsec/`**）。
      **收尾前必跑 `sh scripts/d22scan.sh`**（A64②）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；**不 push**；
不在仓内建 worktree（A38④，快照目录带会话后缀）；票面 append-only（改行前先读；标题前插段落要重抄标题，
`git diff --numstat` 删除列必须 0）；四种假绿逐条点名；数字不达标写 FAIL 附数字。
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。
⚠ 共树：`internal/winsec/`（票 94）、`internal/config`+`internal/agent`（票 90）、`tools/d22scan/`（票 96）都有人在写，
**本票只碰 `scripts/d22scan.sh` 这一个文件**，别扩界。

## Progress log（append-only）

- 2026-09-21 16:0x（编排者）：建票。票 88 的验收代理交件时把它判为"非本票账"并交回给我 ⇒ 按规矩落成票，不留口头。
  我把它排在很后：它**只在本地自检路径上有害**（CI 是干净 runner），
  但**恰恰是我们现在要求每个代理都走的那条路**（A64② 刚立的"收尾必跑此脚本"）⇒
  一条会回放的步骤当门禁，等于给每张新票发一台可能说谎的自检仪。
  ⚠ AC#1 我特意写成"**量出来，不许只论证**；量不到也照修"——因为这类"缓存会不会命中"的说法，
  本仓已经抓到过好几次"听起来对、实测没复现"，判据不能建在猜测上。
  next= 排队（它小、独立，可与票 93 同批做：那两张都是"两套仪器一严一松"的形状）。

## 交件读数（agent-ticket99，2026-09-21 18:05 实测 `date`）—— Status: **ready-for-review**

改动**只有一个文件**：`scripts/d22scan.sh`（`git diff --numstat` = 23 增 4 删，全在第一步那两行 + 头部注释那一段）。
`tools/d22scan/**`（扫描器本体）、`runtests.sh`、`allowlist.txt`（仍 **5 行非注释**，我一行没动）、`.github/workflows/ci.yml`、
任何测试断言与任何阈值：**一字未动**。变异/复跑全在 /tmp 的 `git archive HEAD` 快照里做
（`/tmp/wisp-t99-agent-ticket99`=AC#1、`/tmp/wisp-t99b-agent-ticket99`=AC#3/AC#4），仓内未建 worktree、未 checkout。
仪器：`go version go1.27.1 windows/amd64`。

### AC#0 开工先量 HEAD（因为票面那句旧读数很可能是腐坏的引用）

主树 `sh scripts/d22scan.sh`（脚本=HEAD 原样）：rc=**0**，末行 `d22scan: clean - no D22 ban violations`，
`ban #6 frontend/=40`、`ban #7 internal/tools/=17`、`ban #8 internal/=340`、`ban #8 cmd/=26`。
⇒ **票面"HEAD 上现在它是 rc=1，唯一命中 `internal/winsec/winsec.go:126`"确认是腐坏的引用**（票 94 已修完），
本票的门禁判据按实测写成 **rc=0**，不是"除那一条已知命中之外没有新命中"。我没有碰 `internal/winsec/`。
⚠ 同一次运行里就看本票要修的东西：第一步打印 `ok  	github.com/CarlosShao/wisp/tools/d22scan	(cached)` —— 门禁那一步在主树上就是回放。

### AC#1 缓存端不端得过去：**量出来了**，判定分两级（不是一句"复现"糊过去）

快照 `/tmp/wisp-t99-agent-ticket99`（`git archive HEAD` 导出，脚本=HEAD 原样=裸 `go test ./...`）：

| 腿 | 命令 | 原文读数 |
| --- | --- | --- |
| ① 跑一次让它绿 | `sh scripts/d22scan.sh` | rc=**0**；第 2 行 `ok  	github.com/CarlosShao/wisp/tools/d22scan	3.524s`（真跑，无 cached）；末行 clean，`ban #6 frontend/=37` |
| ② 种一个已知违规 | `printf 'export function decide() { return approval.decide({allow: true}); }\n' > frontend/src/app.js` | `grep -n .` → `1:export function decide() { return approval.decide({allow: true}); }`（ban #6 既有测试的种子形状，零 emoji；只种在快照） |
| ③ 再跑（整脚本） | `sh scripts/d22scan.sh` | rc=**1**；第 2 行 `ok  	github.com/CarlosShao/wisp/tools/d22scan	(cached)`；第 13 行点名 `frontend/src/app.js:1: [panel-approval] ...`；第 14 行 `d22scan: 1 finding(s); D22 bans are not negotiable`；第 15 行 `exit status 1` |
| ③′ 再跑（**只**第一步＝门禁那一步本体） | `cd tools/d22scan && go test ./...` | rc=**0**，原文 `ok  	github.com/CarlosShao/wisp/tools/d22scan	(cached)` |

**诚实判定**：
- **步骤级=复现**：种了违规之后第一步仍打 `(cached)` 的 `ok`、rc=0 ——
  `TestScannerSelfScanOfRealRepoIsGreen`（测试**运行时**读仓库树，缓存判定里只有构建输入）**没再看那棵树一眼**，
  "证明门禁能红"的那一步端了一棵它没检查过的树。这条序列已记下（AC#3(i) 又独立跑出一次，见下）。
- **脚本级=未复现**：第二次整体 rc=**1**。原因说清楚，不藏：第二步 `go run . -root` 每次真跑、不受测试缓存管辖，这条违规是**它**抓的。
⇒ 危害的准确形状是：**回放的是"阳性对照"那一步，兜住的是"扫描"那一步**。只要第二步被改坏／被 `\|\| true`／被跳过
（票面 D22 run-away mode 6 那一族，也正是票 67/71/96 反复抓的形状），第一步再也没有能力说"不"；
而且它今天在主树上就已经是 `(cached)`。判据不能建在"第二步恰好还在、恰好没命中缓存"上 ⇒ 按票面 AC#1 末句，仍然做 AC#2。

### AC#2 修法：**改调 runtests.sh**，两套规矩并成一套（不是再抄一份 `-count=1`）

HEAD 的第 31-32 行（改前）：

```sh
echo "d22scan.sh: positive control - go test ./... (tools/d22scan)"
go test ./...
```

现第 50-51 行（改后）：

```sh
echo "d22scan.sh: positive control - runtests.sh -C tools/d22scan ./..."
sh "$root/tools/d22scan/runtests.sh" -C tools/d22scan ./...
```

- 对齐到的那一侧是 `tools/d22scan/runtests.sh:75` 的 `go test -v -count=1 "$@"`，
  也正是 CI 的 "D22 scanner positive control" 步（`.github/workflows/ci.yml:48`）**逐字**在用的调用形状
  ⇒ 本地脚本与 CI 从此走同一台仪器，"一处严一处松"并成一处。
- **没有**加 `-failfast`/`-short`/`-run` 任何一个改变覆盖面的开关；`-v`/`-count=1`/"SKIP 不算过" 是 runtests.sh 自带的既有约束，
  不是我新立的阈值（我没改 `runtests.sh` 一个字）。
- 头部注释里"两步都是 load-bearing"那一段重写了，把 AC#1 的 `(cached)` 原文钉进脚本（零 emoji，`bash -n` rc=0）。

### AC#3 变异双向（快照 B；锚点=承载行为那一行本身，不是我记忆里的行号）

| 发 | 动作 | 落地证明（同一条 `&&` 链里） | 读数 |
| --- | --- | --- | --- |
| (i) 去掉 `-count=1` 的承载 | 把第 51 行那句 `sh "$root/tools/d22scan/runtests.sh" -C tools/d22scan ./...` 换回裸 `go test ./...`。`-count=1` 由 runtests.sh:75 承载而 `tools/d22scan/**` 禁改 ⇒ 变异做在**调用行**上（快照内 `sed -i`） | `grep -n "^go test ./...$"` → `51:go test ./...`；`bash -n scripts/d22scan.sh` rc=**0**；`cd tools/d22scan && go build ./...` rc=**0**（本票一行 Go 码都没动，"编译失败不算变异"这条不成立即不适用） | 干净树跑 mutant rc=**0**，第 2 行 `ok  	github.com/CarlosShao/wisp/tools/d22scan	3.345s`（真跑，把缓存喂热）⇒ 种 `frontend/src/app.js` ⇒ 再跑：第 2 行 **`ok  	...	(cached)`**，整体 rc=1 只因为第二步；单独再跑第一步 `go test ./...` → **rc=0 `ok  	...	(cached)`** ⇒ AC#1 那条"种了违规还绿（在承载对照的那一步上）**重新出现** |
| (ii) 阳性对照（修法在场） | 同一快照、同一个违规文件，脚本=修法版 | `grep -n "^sh \"$root/tools/d22scan/runtests.sh\""` → `51:sh "$root/tools/d22scan/runtests.sh" -C tools/d22scan ./...` | `sh scripts/d22scan.sh` rc=**1**；第一步**真跑**（整份日志 0 个 `(cached)`）并点名：`scan_test.go:269: repo HEAD violates: frontend/src/app.js:1: [panel-approval] ...`、`--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (0.28s)`、`--- FAIL: TestRealRepoLedgerIsHonest (0.28s)`、`runtests.sh: go test exited 1 - packages=[./...] top-level: PASS=19 FAIL=2 SKIP=0, === RUN=31` ⇒ **这条调用路径有牙**（票 88 证的是扫描器本体，这里证的是脚本） |
| 还原 | 从主树 `cp` 回脚本 + `rm -f frontend/src/app.js` | `diff -q` 快照那份 vs 主树那份 → rc=**0**（逐字节相同）；`ls frontend/src/app.js` → No such file | 主树 `git status --porcelain` 只剩别人的东西（`M .scratch/.../102-*.md`、`?? docs/evidence/s1/102-*.md`、`?? internal/winsec/seam_guard_windows_test.go`）⇒ **我没提交、没还原、也没覆盖**别人的改动 |

### AC#4 门禁

- `bash -n scripts/d22scan.sh` rc=**0**。
- 纯净快照 `/tmp/wisp-t99b-agent-ticket99`（`git archive HEAD` + 我这版脚本）`sh scripts/d22scan.sh` **连跑两次**：run1 rc=**0**、run2 rc=**0**，
  两次末行都是 `d22scan: clean - no D22 ban violations`，台账 `ban #1-5 internal/=197, ban #1-5 cmd/=20, ban #6 frontend/=37, ban #7 internal/tools/=17, ban #8 design/=16, ban #8 frontend/=37, ban #8 internal/=340, ban #8 cmd/=26`；
  run2 第一步 `runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0` —— PASS=21 是真跑出来的，不是 cached。
- 收尾前按 A64② 在**主树**用改后的脚本再跑一次：`sh scripts/d22scan.sh` rc=**0**，
  第一步 `runtests.sh: OK - ... PASS=21 FAIL=0 SKIP=0`，末行 `clean`、`ban #6 frontend/=40`、`ban #8 internal/=340`、`ban #8 cmd/=26` ⇒ **台账不降**。
- `ban #6 frontend/` 37（快照）vs 40（主树）差异**已点名归因**，不糊：`git ls-files -- frontend`=**37**、`find frontend -type f`=6566、
  `git status --porcelain --ignored -- frontend`=**3**（`frontend/dist/assets/`、`frontend/dist/index.html`、`frontend/node_modules/`）
  ⇒ 那 3 个是**构建产物**（git-ignored，`git archive` 不导出），**不是**共树代理的未提交源码，也**不是**本票造成的。
- 本机既有坑：`GOOS=linux go vet ./...`（整仓）永远 rc=1 —— 与本票无关，未跑整仓门禁（共树：票 101/102 在验收、票 103 在写 `internal/winsec/`、票 93 在写 `scripts/` 新 runner 与 `ci.yml`）。

### 四种假绿逐条点名

1. **空仪器**：第一步现在自带"PASS=0 且 FAIL=0 就红 / 有 SKIP 就红"两条断言；我没删步骤，两步仍 `set -eu` 串联、无 `|| true`、无 continue-on-error。
2. **回放绿**：本票主题。AC#1 用 `(cached)` 原文钉死，AC#3(i) 证明把承载去掉它就回来；改后日志里 `(cached)` 出现 **0** 次。
3. **数字被抹/被挑**：所有 rc 逐次单报（AC#1 三次、AC#3 四跑、AC#4 两跑 + 主树一跑），无平均、无"挑运气好的那次"；台账差异归因见上。
4. **借他人之口**：AC#3(ii) 的红是**这条调用路径**自己跑出来的（`--- FAIL: TestScannerSelfScanOfRealRepoIsGreen`），不是引用票 88 的结论。

⚠ 本代理这一程的工具输出末尾**没有**出现自称"编排者备注/停手/冻结某包"的附加文本，故无原文可登记；
未改任何契约文本的语义（D22 未触发）。`internal/winsec/zz_diag_test.go` 开工时在 `git status` 里、交件时已不在，是票 103 自己的动作，与本票无关，未动。

next= 等验收方出 `docs/evidence/s1/99-*.md` 裁决表。编排者的账（**不在本票**，本票一行没改）：
`.github/workflows/ci.yml:48` 与 `:68` 现在跑的是同一台仪器的同一次测试（CI 多花一遍 tools/d22scan 的 `-v` 测试，本机实测第一步 3.3-3.5s + 21 用例），
要不要把两步并成一步归票 85/93；若并，判据是"两步合一"而不是"把 `scripts/d22scan.sh` 改回去"。

## 验收读数（acceptor-ticket99，`date` 实测 2026-09-21 18:27:42 CST）—— 总判：**通过（附条件）**

裁决表：`docs/evidence/s1/99-adversarial-acceptance.md`（四框 1:1，每行标来源档位）。快照全在仓外 `/tmp`
（`/tmp/wisp-99-ac1-ac99`=`git archive d0d8782^`、`/tmp/wisp-99-mut-ac99`、`/tmp/wisp-99-gate-ac99`=建快照时 HEAD `06906f7`），
仓内未建 worktree、未 checkout；**主树未跑门禁**（共树在飞：票 95 `internal/config`、票 93 `ci.yml`+`scripts/`）。仪器 `go1.27.1 windows/amd64`。

- **AC#1 〔独立复现〕**：改前脚本（`grep -n '^go test \./\.\.\.$' scripts/d22scan.sh` → `32:go test ./...`）在快照里 run1 rc=**0**
  第 2 行 `ok  	github.com/CarlosShao/wisp/tools/d22scan	10.352s`（真跑喂热缓存）→ 种 `frontend/src/app.js` 的 `approval.decide` →
  run2 整脚本 rc=**1**，但第 2 行是 `ok  	github.com/CarlosShao/wisp/tools/d22scan	(cached)`；**单独再跑第一步本体**
  `cd tools/d22scan && go test ./...` → rc=**0** 同一行 `(cached)` ⇒ **步骤级=复现／脚本级=未复现**（红是第二步 `go run .` 抓的），与票面判定一致。
  差异说明：`go test` 缓存键含**运行目录**，所以新 `/tmp` 快照第一次必然真跑、同目录第二次回放 ⇒ 危害只咬"同一工作树反复自检"这条 A64② 的路（见 R-99-3）。
- **AC#2 〔独立复现〕覆盖面三问**：① **同一测试集合**——改前改后都 `cd "$root/tools/d22scan"`（第 48 行是 context 行、未动），
  `go list ./...` 在该 module 只有 **1 个包**，`scripts/spike`/`tools/mockllm` 改前改后都不在分母 ⇒ 没把被检查对象删出门禁；
  ② "SKIP 不算过 / PASS=0 且 FAIL=0 就红"是**纯加严**（只会把 0 变 1，rc 原样透传），`$@` 只有 `./...`，
  无 `-failfast`/`-short`/`-run`、无 `GOFLAGS` ⇒ 覆盖面语义未改；继承来的新风险已登记（`tools/d22scan` 将来若有平台性 SKIP，第一步会红且原因不是禁令，那是**正确**的红，不许回头放宽脚本）；
  ③ **任意 cwd 成立**：`cd /tmp && sh <绝对路径>/scripts/d22scan.sh` rc=**0**；`cd <snap>/frontend/src && sh ../../scripts/d22scan.sh` rc=**0**；
  `cd /tmp && sh <snap>/tools/d22scan/runtests.sh -C tools/d22scan ./...` rc=**0**（`runtests.sh:53-56` 把 `-C` 解析成自身推导的 `$root/…`，`:58-62` 的 `go.mod` 守卫让写错的 `-C` 只能 exit 2、不能假绿）。
  ④ **CI 那侧**：`:48` 与 `:68` 第一步现在是逐字同一条命令 ⇒ CI 每轮多跑一整遍（本机冷跑 9.5s/21 用例/`=== RUN`=31）。**不是本票的账**（AC#2 正是要求并成同一台仪器），
  登记给票 85/93；合一的判据=**合 CI 两步**（保留 `scripts/d22scan.sh` 为唯一入口），**不是**把脚本改回裸 `go test`。
- **AC#3 〔独立复现〕双向**：(i) 在**有修法的同一棵树**上把承载行 `51` 换回裸 `go test ./...`（`MUTANT-LANDED` + 同链 `grep -n` 打印 `51:go test ./...` + `bash -n` rc=0），
  干净树 rc=0 `ok ... 9.492s` → 种同一违规 → 整脚本 rc=1 而第一步 `(cached)`、**第一步单独跑 rc=0 `ok ... (cached)`** ⇒ "种了违规还绿"**回来了**。
  ⚠ 我头两发 `sed` 因锚点漏算 `"$root` 里的 `$` 而**静默不匹配**，被同链 `grep -n` 无输出打断才发现（那次 run1 日志第 2 行是 `=== RUN` = `-v` 输出 = 脚本没被改），读数已作废重跑——票面"同链 grep 证落地"这次救了我。
  (ii) 修法在场 + 同一个 `frontend/` 违规：脚本 rc=**1**，红名 `--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (0.47s)`（`scan_test.go:269: repo HEAD violates: frontend/src/app.js:1: [panel-approval]`）
  + `--- FAIL: TestRealRepoLedgerIsHonest (0.49s)`，`runtests.sh: go test exited 1 - packages=[./...] top-level: PASS=19 FAIL=2 SKIP=0, === RUN=31, '[no tests to run]'=0`；
  该日志里 `d22scan.sh: scan of` **零命中** ⇒ `set -eu` 在第一步就死，第二步根本没跑 ⇒ **这条调用路径自己有牙**（没借票 88 的口）。还原后 `diff -q` 与主树逐字节相同（rc=0）、种子文件已删。
- **AC#4 〔独立复现〕**：`bash -n scripts/d22scan.sh` rc=**0**（快照与主树各一次）；纯净快照 `sh scripts/d22scan.sh` 连跑两次 rc=**0**/**0**，
  两次 `runtests.sh: OK - ... PASS=21 FAIL=0 SKIP=0, === RUN=31`；`(cached)` 我自己数 = **0** 与 **0**；台账 `d22scan: scope` 行 **8** 行齐全：
  `bans #1-5 internal/`=197、`bans #1-5 cmd/`=20、`ban #6 frontend/`=**37**、`ban #7 internal/tools/`=17、`ban #8 design/`=16、`ban #8 frontend/`=**37**、`ban #8 internal/`=342、`ban #8 cmd/`=26
  ⇒ `ban #6/#8 frontend/` **不降**（342 vs 票面 340 是票 103 新增 `internal/winsec` 文件，只增不减，非本票造成）。

**归属核对（本票来历特殊）**：`d0d8782`（`23 4 scripts/d22scan.sh`，**修法本体**）message 首句自证"由编排者代 agent-ticket99 落档"，并明写 AC#3/AC#4 读数仍缺、由 99b 续；
`9e00629`（`99 5` 票面，**只改票面**）把四框读数整段署在 `agent-ticket99 @18:05` 名下（且 line 126 留了"我没提交"这句与仓库状态矛盾的残留），
而 AC#3/AC#4 的快照目录 `/tmp/wisp-t99b-agent-ticket99` 是 99b 的；第三枚 **`06906f7`**（18:13:19，`13 0` 票面，= 我建快照时的 HEAD）
已由 99b 自登记"账分两枚：代码在 `d0d8782`，票面在 `9e00629`" ⇒ **拆分已披露，不是隐瞒，读数本身不假**（我逐条复现出来了）。
仍欠两笔（登记 **R-99-1**）：(a) 交件读数抬头 line 61 与 Status line 3 把 99b 测的 AC#3/AC#4 仍冠 `agent-ticket99 @18:05`；
(b) `06906f7` 的"我的脚本改动被 `d0d8782` 带走"与 `d0d8782` 的"编排者代 agent-ticket99 落档"两种口径并存，脚本改动的作者没有唯一定稿（归编排者定）。
本段为追加，未改写 99b 正文、未改 Status、未改文件名。
- 2026-09-21 18:3x（**编排者，结案并把你要求的那处更正落定**）：
  ① **同票双派的根因在我**：我先收到一枚"轮数上限死亡"通知，据通知**代落档**脚本（`d0d8782`）并派出续跑代理；
    而原代理 `agent-ticket99` **没死**，随后自己跑完并交件（`06906f7`）⇒ 台账登记为 **A77③**，
    并已把规则改写成"**死亡通知≠已停手；落档前先查 `git log`+`git status` 里它自己的痕迹**"。
  ② **归属定稿（就是你说"没有唯一定稿"的那处）**：AC#1/AC#2 = `agent-ticket99`（checkpoint `1f8d212`）；
    脚本代码 = 编排者代落档 `d0d8782`（内容与你逐字节相同的那份由原代理在工作树里持有，我先提交了，**作者记原代理**）；
    **AC#3/AC#4 = 两个独立会话各跑一遍** —— `agent-ticket99`（`06906f7`，PASS=19/FAIL=2）
    与 `agent-ticket99b`（发现死亡后补做，票面 `9e00629`，被叫停、未 push）。
    ⇒ 两遍读数**彼此一致**（同红名 `TestScannerSelfScanOfRealRepoIsGreen`、同 `(cached)` 原文）：
    **本票的判据不是"某人说了算"，是被独立跑过两次**。你与验收都指出的 `R-99-1` 以票头 `accepted-done` 段为准，原文一律不删。
  ③ `R-99-3` 我接受并已进台账：**"缓存会端过去"必须加限定语——"同一目录内、包与测试源码未变"**（缓存键含运行目录）。
    `R-99-4`（CI `:48`/`:68` 因此双跑同一台仪器）归票 85，判据写死"合两步、**不许**把脚本改回裸 `go test`"。
    文件名已挂 `-done`（防重领键归编排者）。

**登记（验收期一律没动代码/契约）**：R-99-1 记账更正（本段落档）· R-99-2 票面 AC#4 括号"HEAD rc=1 / `internal/winsec/winsec.go:126`"是**腐坏的引用**
（我在 `d0d8782^` 与 `06906f7` 两棵纯净树都是 rc=0，`internal/winsec/` 一行没读没改），判据按实测写 · R-99-3 缓存键含运行目录 ⇒ 脚本头部注释与票面该补"同一工作树内反复自检才会回放"这个限定语（我没动脚本） ·
R-99-4 CI 两步重复跑，账归票 85/93 · R-99-5 共树 HEAD 在我会话中从 `0717bf2` 漂到 `06906f7`，本报告每份绿都点名来自哪棵 `/tmp` 快照 ·
R-99-6 `sed` 静默不匹配两发 = 方法坑，同链 `grep -n` 整行是唯一防线。

**第三档（仅自述，不背书）与补救**：99b 自述的 `/tmp/wisp-t99*` 目录内日志、"主树收尾 rc=0"、"allowlist 仍 5 行非注释"中属于它当时那次运行的部分，
我无从复跑 ⇒ 不背书。补救动作：AC#1/AC#3/AC#4 我已用**自己的**快照独立复现（上列读数即我的，不是抄的）；"主树收尾 rc=0"改由票 93/103 收尾时在共享树跑一次并点名每条命中的归属；
`allowlist.txt` 行数我按 `git show --numstat` + 现树计数抽验（见裁决表末段）。

next= 交编排者：`-done` 时把 R-99-1 的 sha 拆账并进台账、按 R-99-2 把 AC#4 的旧括号标为过期、按 R-99-3 补脚本注释限定语；R-99-4 派给票 85/93。本代理未 push、未改 `tools/d22scan/**`/`allowlist.txt`/`ci.yml`/任何断言与阈值。

- 2026-09-21 18:12（agent-ticket99，交件后立刻登记的一条**提交形状异常**，如实写不藏）：
  我这枚路径限定提交 `9e00629` 的 `--numstat` 里**只有票面一个文件**（`99 5`），因为 `scripts/d22scan.sh` 的改动
  在我提交前 52 秒（18:09:54）已被另一枚提交 **`d0d8782`**（`fix(99,AC#2): d22scan.sh 第一步改走 runtests.sh`，`23 4`）
  从我未提交的工作树里带走了 —— `git show d0d8782` 的那 27 行**与我写的字节字节相同**（`index bd6db7c..f64f2f0`），
  不是我提交的（我的提交动作只发生过一次，就是 `9e00629`）。⇒ **净结果无损**：HEAD 上脚本内容 = 我 AC#3/AC#4 实测的内容
  （`git diff --quiet -- scripts/d22scan.sh` 为空），票面在 `9e00629`。但**账要分两枚看**：
  代码那半在 `d0d8782`，交件那半在 `9e00629`；共树里有人（或另一个 99 的执行者）会把我未提交的改动扫进它自己的提交，
  这是编排者的账（A38④ 的共树家族），**本票不动历史**（禁 `--amend`/`reset`/`rebase`）。
- 同一次登记附**收尾复跑**（A64②，committed HEAD `9e00629`，主树，18:12）：`sh scripts/d22scan.sh` rc=**0**，
  整份日志 `(cached)` 出现 **0** 次，第一步 `runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31`，
  末行 `clean - no D22 ban violations`（`ban #6 frontend/=40`、`ban #7 internal/tools/=17`、`ban #8 internal/=340`、`ban #8 cmd/=26`）⇒ 台账不降。
  **未 push**。next= 验收方出 `docs/evidence/s1/99-*.md`。
