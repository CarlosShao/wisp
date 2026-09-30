# 244 — SPEC-11 要求"无参数＝GUI 子系统"，`docs/BUILD.md` 把这一切换**推迟给票 07**、票 07 结案时**没把它带走**：今天 `build/wisp.exe` 的 PE 子系统我读到的是 **CUI(3)**＝双击连一个黑控制台窗口一起起，而两行注释反过来声称"已经链接成 GUI 子系统"

- Status: **已立，未派**（09-30 16:1x，编排者立）。来源：只读普查腿 `243-c1` 交的表里第 22／23 行（裁决件 `.scratch/wisp/probes/243/c1/census.md` §3 把这两行列为"最危险"），**指认由我自己逐把尺复跑复认**（下面现量表全是我 16:0x 自己跑的，不是转抄腿）。
- 为什么单开一票而**不塞进票 228**：那两行注释**说的不是球的活**，是**构建产物的形态**；而这一格今天在**整个票池里无人认领**（尺见下表第 4 行）。票 243 是"只出表、零产码"，也不该接这活。
- ⚠ **这不是安全事件，照三行读**：① **现象在哪**＝出厂产物的外观（双击 `wisp.exe` 会连一个控制台窗口一起起，球旁边挂一个黑框）；② **有没有本机被入侵的证据**＝**没有**；③ **最坏后果是什么形状**＝**产品看起来像半成品** ＋ `cmd/wisp/console_windows.go` 那整条"attach 父控制台"的通路**今天永远走不到**（它自己的注释还以为正在用它）。⇒ 外观与一处死通路，不是被攻破。

## 现量（09-30 16:0x 编排者自己跑，别信行号、自己复算）

| 事实 | 读数 | 尺 |
|---|---|---|
| **规格要求它** | `docs/specs/SPEC-11-build-deploy-containerization.md:50` 逐字「CLI 与 GUI 同一二进制：无参 = GUI（`-H=windowsgui`）」 | `sed -n '50p'` 我现读 |
| **这笔推迟有明文记录** | `docs/BUILD.md:87` 逐字「windowsgui 子系统切换（`-H=windowsgui`）**推迟到票 07**；届时如需「缺 DLL 仍能弹出友好错误」，应改为运行时 LoadLibrary 包装，这是那票的设计题」；`:90` 又说「AttachConsole 通路已按本票代码预先验证：用 `go build -ldflags "-H=windowsgui"` **临时构建**后…」 | `sed -n '85,92p' docs/BUILD.md` 我现读 |
| **票 07 结案了、没带走它** | 盘上「07-ball-state-machine-core-done.md」（带 `-done`＝结案）；而 `docs/BUILD.md:90` 那句"验证过"是**临时构建**，不是构建链 | 我现跑 `ls .scratch/wisp/issues/ \| grep '^07-'` |
| ⛔ **今天无人认领这一格** | 全票池提到 `windowsgui` 的只有**两枚都已结案的票**：`01-build-chain-done.md:47`（记的是"临时构建下验证 AttachConsole"）与 `117-…-done.md`；**没有任何在开的票拥有这一切换** | `grep -rn windowsgui .scratch/wisp/issues/` 我现跑（⚠ 我上一发同类尺把中文分支和 `windowsgui` 拼进同一个模式串，**回了零命中＝坏尺**，换词重跑才抓到这两枚；教训见 `A473`） |
| **今天产物是控制台子系统** | `build/wisp.exe`（09-30 **11:21**，29,771,385 字节）→ `objdump -p` 读 **`Subsystem 00000003 (Windows CUI)`**；`build/balldebug.exe` 同为 CUI | `objdump -p build/wisp.exe \| grep -i subsystem` 我现跑 |
| **构建链里根本没有这一维** | `scripts/build.ps1:103-110` 的 `$ldflags` 六枚**全是 `-X`（版本／commit／日期／环境／两个 DLL 版本）、零枚 `-H`**；全仓 `grep -rn '\-H[= ]' scripts/ .github/ tools/` 唯一命中是一枚二进制 blob（`scripts/spike/bin/goja-caps.exe`） | 我现读＋现跑 |
| **两行注释声称相反的事实** | `cmd/wisp/console_other.go:7`「The Windows build **links as a GUI-subsystem binary** (ticket 07), which starts with no console at all」＋`cmd/wisp/console_windows.go:26`「the **windowsgui** subsystem (the final GUI build, ticket 07)」 | `sed` 我现读；⚠ `windowsgui` 在 Go 代码里**只出现在这两行注释里**（全仓尺同上） |

## 要建什么（四格，顺序有讲究：AC#1 单独做完会留下一枚"看着对、其实没人验过"的产物）

- [ ] **AC#0 先答"现在切还是等"**（归编排者排程，**不摆 owner**：SPEC-11:50 已经要求，这是**漏做**、不是契约变更）——⚠ 但要把**回归面**写清：切到 GUI 子系统之后，`wisp run`／`wisp doctor` 这些 CLI 腿**必须靠 `attachParentConsole` 才有输出**，而那条通路今天的真机凭据只有 `docs/BUILD.md:90` 那句**当时的临时构建**读数（＝过期读数，不可继承）。⇒ **AC#1 与 AC#2 不许拆成两批改**：只加 flag 不验 CLI＝把 CLI 那条腿盲切。
- [ ] **AC#1 构建链真带 `-H=windowsgui`**：改后产物 PE 子系统必须是 `00000002 (Windows GUI)`。尺＝`objdump -p build/wisp.exe \| grep -i subsystem`，**改前读数 CUI(3) 与改后读数都要写在交件里**。⚠ 具名口径：`build/**` 被 `.gitignore` 掉 ⇒ 这枚产物是**中间态、不可再生**，所以交件必须带**产生它的那条命令与 flag 全文**，否则下一位只能读到"当时"、读不到"现在"（本仓为这一形付过学费）。
- [ ] **AC#2 切换之后的 CLI 真机读数（本票的验收面，不许省）**：三条各一发真跑——① 在**父控制台里**跑 `wisp run "…"`，回复文本真出现在那个控制台；② `wisp doctor > out.txt` 重定向仍然有效（`BUILD.md:90` 当年自陈的一条）；③ **双击／`explorer` 拉起**（无父控制台）时**不再出现黑框**、且常驻腿照常起。⛔ 不许用"单测里 mock stdout"代替真机——那正是 `AGENTS.md` §1.3 禁的形状。
- [ ] **AC#3 那两行注释改成带条件的事实句**：说清"今天构建链带不带 `-H`、因此 `attachParentConsole` 这条分支今天走不走得到"。⛔ 不许留"links as a GUI-subsystem binary"这种**切换一没发生就变假话**的绝对句（与票 228 AC#7、票 197 那枚 `run_mode101_test.go` 同族）；⛔ **不新增扫注释票号的词面型仪器**。
- [ ] **AC#4 把"文档里的推迟"这一形报给票 225 那一族**（⚠ 本票不做，只登记）：票 225 查的是代码里的 `DEFERRED(D-xx)` 标记与 `SPEC-12 §5` 的双向 1:1；而 `docs/BUILD.md:87` 这种**写在文档正文里的"推迟到票 NN"**不在它的射程内——**票 07 一结案，那笔推迟就从盘上消失了**。⇒ 具名交回 225 的后续：**要不要一枚扫 `docs/**` 里"推迟到票 NN"并核 `-done` 的尺**（注意：这枚尺必须是**报表**，不能当门——票 243 AC#3 已证"交付动词＋指名处"这类句型判别是**抽样尺**，反例是 `internal/proc/boot_windows.go:46`：它引已结案的票 06，而那件活**真落地了**（`internal/proc/envfork.go:211 ApplyPortableOverride` 在场）⇒ 词面型门会把这条判反）。

## 禁区

- ⛔ 不动 `docs/PLAN.md`／`docs/specs/**`／`docs/BUILD.md` **一字**（本票是**照 SPEC-11:50 补做**，不是改规格；`BUILD.md:87` 那句"推迟到票 07"过期了也**不改它**，由本票在票面具名更正）。
- ⛔ **SLO 阈值／golden／`internal/observe/thresholds.go`／`allowlist.txt` 一字节不动**；三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）一字不动；⛔ 不许为变绿放宽任何断言。
- ⛔ 不许顺手把"球进常驻""审批门接线"做掉（那是票 228）；不许顺手把 GUI 的"缺 DLL 弹友好错误"做掉（`BUILD.md:87-89` 明写那是**设计题**，要动先摆出来）。
- `frontend/**`／`design/**` 零读零写零转述。
- git：只 commit 不 push；commit 必带显式 pathspec；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件。

## 排程与串行

- 动 `scripts/build.ps1`（构建链）＋`cmd/wisp/console_other.go`／`console_windows.go` ⇒ ⛔ **与 `228-r1` 串行**（它此刻正在 `cmd/wisp` 写，且新增了一枚 `cmd/wisp/resident_ball_windows.go`；同包并发＝互相洗读数）。
- ⚠ **改构建 flag 会改变 CI 产物形态**：`slo-full` 跑在本机 self-hosted runner 上、每次 push 自启抢 CPU ⇒ **取数期间不改构建链**；本票排在**票 228 主体落地之后**（球进常驻之后再看"双击起什么"，一次改对，免得为同一件事改两次构建链）。
- 派单前必做**撞钉预检**：先跑 `go test ./cmd/wisp/ -count=1`（带 sherpa PATH）把**今天绿的用例名**抄进派单，并逐枚读 `console_*`／`attachParentConsole`／`stdout` 相关的断言——本仓已知一枚"读盘的 AST 测试对 `-overlay` 失明"的坑，构建形态改动同样可能只被真机看得见。
