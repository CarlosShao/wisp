# 票 305 判死件（305-r1 写腿，2026-11-16 14:0x +08）

**性质**：判死件，非修复件——**两枚症状在 HEAD 上均已绿，本腿零产码改动**（`git diff -- cmd/wisp internal/panel`＝空，全程只在工作树上做过两枚突变并逐字节还原）。

## 0. 台面（按票面三根轴交齐）

- **HEAD**＝`a683a6706`（派单锚，未前进）；工作树进场时已有 4 枚别腿修改件（`44-costmeter-c23-done.md` 等，本腿⛔ 碰）。
- **dist 第二轴**：`npm run build`（vite 8.3.0，2.94s）现造，`wc -c frontend/dist/index.html`＝**1044 B**，mtime **2026-11-16 13:52:45**，`git cat-file -s HEAD:frontend/dist/index.html`＝`exists on disk, but not in 'HEAD'`（HEAD 只有 `.gitkeep`）。跑完已删回 `.gitkeep`（见 §6）。
- **第三轴**：全部 Go 测试经载具 `scripts/wisp-cli-staged-tests.sh`（`D:\work\soft\Git\bin\bash.exe`，A1884），DLL staged sha256 ok，本机无 onnxruntime System32 影蔽（载具首枚守卫未触发）。每次跑前 `Get-Process wisp.test`＝0（串行单发，A1888）。

## 1. 三枚钉今日读数（复量，⛔ 抄票面 10-10 行号）

尺＝载具 `-test.run 'TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe|TestAC14AwaitedBindingReplyReachesThePage|TestAC14GoSideEvalPushReachesThePage'`（一发）：

| 钉 | 今日读数 | 页面自己的话（逐字） |
|---|---|---|
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | **PASS (1.82s)** | settle `"https://wisp.panel/index.html\|complete\|Wisp"`；entry-id report `"1"`；`AC#13 page answer (head a683a6706): "1" - 1 of 1 probe id(s) present in the live document` |
| `TestAC14AwaitedBindingReplyReachesThePage` | **PASS (1.56s)** | `AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)` |
| `TestAC14GoSideEvalPushReachesThePage` | **PASS (2.33s)** | settle `"https://wisp.panel/index.html\|complete\|Wisp"`；`AC#14 nail 2 (Eval push hop), page's own words: title="PUSHED-33R5-OK"` |

## 2. AC#1／AC#2 判死（两枚症状均已绿 ⇒ 具名"由哪提交绿"，⛔ 假装修）

- **症状①（AC13，入口交接）＝本来就好**。绿由 **`3b806876f`（299-r1，2026-10-13）** 引入的形状锁死：交接从 `SetHtml` 改成**导航到 `panelEntryURL`＋资源过滤器**（尺＝`git log -S panelEntryURL -- cmd/wisp` 首中 `3b806876f`；`git log -S AddWebResourceRequestedFilter -- cmd/wisp` 同枚）。绿凭据＝299-r1b 的绿发 `.scratch/wisp/probes/299/r1b/logs/run-ac13-green.txt`（2026-10-13 08:50，workdir 行逐字 `the green emit: both AC13 nails on the committed shape (3b806876)`，`--- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.81s)`），编排者当日收腿时自己复跑（票 299 收腿节、账 `A872`）。后续是**仪器硬化非行为改动**：`70b008856`（33-r10 无窗尺）、`5dee200c0`（447 settle 有界重问）、`8d4eb0e5c`（33-lane2，今日 13:05，not-built notice 也走过滤器带 CSP——同文件但⛔ 碰 AC13 那两跳的次序）。红句彼时那句 `after a real cold start the live document contains NONE of the 1 element ids...` 今天要再见到它，只有把次序拆掉（见 §3 M1）。
- **症状②（nail2，Go 侧 Eval 推入）＝也已绿，由 **`a07f685eb`（402-r1，2026-11-02）** 绿**。该笔给 `TestAC14GoSideEvalPushReachesThePage` 前面加了 `awaitSettledHandoverDocument(t, rp, h)`（有界重问，等最后一份交接文档落地再推；`git show a07f685eb -- cmd/wisp/panel_resident_windows_test.go`＝+56 行测试件，**零产码文件**）——即那枚"推不进"红量的是**推入横跨异步换文档**的测量病，产码推入跳本身没坏。修后中程报 6/6 PASS（commit message 自带）；今天 2.33s PASS 为本腿独立复量。
- **本票 AC#1/AC#2 两格的落地动作＝零产码**：哪枚都不需要修，判死即交付。⛔ 放宽任何判据（两句红句的判据原文一字未动，`git diff -- cmd/wisp`＝空为凭）。

## 3. AC#3 反恒真成对（两对，各 (a) 红逐字＋(b) 逐字节还原复绿）

### M1（症状①）：撤 AC13 修复＝`coldStartPageHandover` 两跳对调（探针页拿回最后一手）

- (a) 红（16.16s）逐字：
  `after a real cold start the live document contains NONE of the 1 element ids the embedded entry declares (the page itself answered "0"). That is ticket 33 AC#13: the document in the window is not the ...`（`panel_resident_windows_test.go:541` 判据原句，⛔ 未放宽；settle 读 `"about:blank|complete|"`、`no probe attempt over 714 attempt(s) saw a present id; last report: "0"`）→ `--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (16.16s)`。
- (b) `git show HEAD:cmd/wisp/panel_host_windows.go > cmd/wisp/panel_host_windows.go` 逐字节还原，`git diff --stat`＝空 ⇒ 复绿：`--- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (1.90s)`。

### M2（症状②）：撤 402-r1 修复＝摘掉 `awaitSettledHandoverDocument` 那一行

- (a) 红（1.39s）逐字：
  `Go's Eval push did not reach the document: the page reports its title as "", want "PUSHED-33R5-OK". This is the push dimension, separate from the awaited-reply dimension in TestAC14AwaitedBindingReplyReachesThePage - on...`（`:1084` 判据原句）→ `--- FAIL: TestAC14GoSideEvalPushReachesThePage (1.39s)`。**与票面 10-10 的历史红句逐字同形**。
- (b) `git show HEAD:cmd/wisp/panel_resident_windows_test.go > cmd/wisp/panel_resident_windows_test.go` 逐字节还原，`git diff`＝空 ⇒ 复绿：`--- PASS: TestAC14GoSideEvalPushReachesThePage (1.84s)`，`title="PUSHED-33R5-OK"`。

### 跨维度负控（票面 AC#3 外加那发）

本腿**零产码改动**，故不存在"修①的那笔让 nail2 变绿"要解释的情形；维度可分性由两枚突变各证一侧：M1（交接红）时推入跳产码原封不动、红句只谈内容不在；M2（推入红）时交接机制（资源过滤器＋导航）原封不动、红句只谈 push。票 305 `AC#0` 收腿节已裁"两枚症状部分同源、⛔ 合并"，本节与其一致、无新证据推翻。

## 4. AC#4 门禁与越界

- `go vet ./cmd/wisp/` ⇒ **rc=0**。
- `go vet -tags winlive ./cmd/wisp/ ./internal/ball/...` ⇒ **rc=0**。
- `sh scripts/d22scan.sh`（`D:\work\soft\Git\bin\bash.exe`）⇒ **rc=0**（control `PASS=46 FAIL=0 SKIP=0`；scan `clean - no D22 ban violations`，denominator 494）。
- 越界：本票唯一 commit 只落 `docs/evidence/s1/305-livedoc-verdict.md`（commit 后 `git show --name-only --format=` 自核，见 §7）；⛔ 碰 `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43／`design/**`（全程未编辑任何一处）；⛔ 放宽断言（两枚判据句逐字保留）；⛔ push；新协程 0 枚（未写任何 `go func`）；⛔ 反向动 C17 白名单；⛔ 提交 dist 字节。
- 整包成对（票面 AC#4 那两根）本程⛔ 跑：派单射程＝定向三枚＋两对突变＋vet＋越界；且本腿零产码 ⇒ "改前改后"是同一棵树，整包成对无信息量，如实记不跑。

## 5. AC#5 与票 33 对账（三行）

| 票 305 格 | 对应票 33 格 | 记账归属（⛔ 两本账同时记） |
|---|---|---|
| 症状①＝AC13 入口交接判死 | 票 33 `AC#13`（冷启动不盖真页面，`33-panel-host-c27.md:39`，已 `[x]`） | **只记本票**（本件 §1/§2/§3）；票 33 `AC#13` 已由编排者 **`A1900`（11-16 12:4x）** 翻勾（"AC#13 冷启动不盖真页面"复跑 6/6 PASS），本票⛔ 再翻 33 那格 |
| 症状②＝nail2 Eval 推入判死 | 票 33 `AC#14`（Go→页面这一跳，`:42`，已 `[x]`） | **只记本票**（本件 §1/§2/§3）；票 33 `AC#14` 已由 **`A1900`** 翻勾（"AC#14 真窗页自报回话到达"），本票⛔ 再翻 33 那格 |
| 落地动作＝零产码（两枚症状均已由 `3b806876f`／`a07f685eb` 绿） | 票 33 无对应格（它管"修没修"，本票管"今天为什么绿＋突变还红不红"） | 本票零产码 commit，凭据与对账全在本件 |

## 6. dist 清理

跑完按派单硬约束删回：`frontend/dist/index.html` 与 `frontend/dist/assets/*`（构建产物）已删，`frontend/dist` 回到 `.gitkeep`＋空 `assets/` 目录（与进场时同一状态）；dist 在 `frontend/.gitignore` 射程内，本票 commit 无任何 dist 字节（`git show --name-only` 为凭）。

## 7. 本件 commit 越界自核

本件首笔＝`f7909e628e14534e14ef2f3c3f9147819272de2f`；`git show --name-only --format=f7909e628` 现量＝**仅 `docs/evidence/s1/305-livedoc-verdict.md` 一行**。本节回填是第二笔（追加，⛔ amend），pathspec 同一枚件。

—— 305-r1（agent-84），commit `f7909e628`。
