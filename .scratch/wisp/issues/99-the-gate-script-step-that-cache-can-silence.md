# 99 — `scripts/d22scan.sh` 第一步是**裸 `go test ./...`** ⇒ 本机测试缓存能把一次真违规"端"成绿

**Status:** ready-for-review（原 open：2026-09-21 16:0x 编排者建；来源=`acceptor-ticket88` 交件时点名的残留，它明写"非本票账"。2026-09-21 18:05 agent-ticket99 交件：四框全勾，读数见文末"交件读数"）
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
