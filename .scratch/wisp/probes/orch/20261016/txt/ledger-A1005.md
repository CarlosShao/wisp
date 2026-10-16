## A1005 [2026-10-16 12:0x +08] 编排者：收四枚交件（`326-r1` 落地／`270-v1`／`120-v2`／`165-c2`）⇒ 翻框 11 枚＋票 120 结案 `-done`＋代死腿 `220-r2` 记收腿一条；★六枚 `cmd/wisp` 红的归因**已闭**；并更正我自己那条"sealHandle 缺世界"的旧账

### ① 起手锚（全部现跑）

- 时刻＝`date` 现读 11:56；`git rev-parse --short HEAD`＝`ddffbf09`；`git rev-list --count origin/dev..HEAD`＝**ahead 95 枚**（写这条时我重跑一次，见本节末"推送"行的最新数；⛔ 推送——本波两枚腿仍在飞，推送只有编排者做、只推 `origin/dev`，机主授权 `A813`）；台账本件＝18722 行（本条追加前）。
- CPU 闸门尺（机主原话＝"只要 cpu 不高于 70% 就可以继续补位"）：11:34＝**62%** ⇒ 放 `326-r1` 之后的补位读数；11:40＝**83%** ⇒ **拒放**、先存档不启动（具名记在 `batch1007-send.txt`）；11:43＝**57%** 且 `go.exe`／`*.test.exe` 名册＝0 枚 ⇒ 放第九波两枚；11:56＝**42%** ⇒ 再补两枚（`326-v1`、`165-r3`）。
- ★⛔ 把"写了派单"当"放了腿"：每一枚的 agentId 只在"启动补记"一节里才算数（`batch1007-send.txt` 已补）。

### ② 收件（只认盘上三把尺；完成通知一律当未验证）

- **`326-r1`**（票 326 落地腿·写，会话里继续聊不会把它顶到前面）：尺① 五笔提交逐枚 `git log -1` 全在盘＝`86298e24`（产码＋用例）→`d3f98938`（读数）→`a60da735`（判语件）→`f2787393`（票尾一条 bullet）→`b4adfd0e`（名册自证）；尺② 逐枚 `git show --name-only` 并集 **31 枚文件**、全在五处写面内（`internal/memory/**`＋`internal/session/**`＋probes＋`docs/evidence/s1`＋票尾），禁区尺 0 命中；尺③ 证据件真在盘＝`docs/evidence/s1/326-append-updated.md` **172 行／17568 字节**、probes 14 枚＋logs 12 枚、**0 字节件 0 枚**（`find -size 0` 空）。它的红句逐字（未修码仓外树）＝`appending a turn did not reach session.updated_at: sess-a.updated_at = 1777593660, want the appended turn's stamp 1777593780 (before=1777593660, delta=0)`；成对＝`1777593660 -> 1777593780 (delta=120) … list = [sess-a sess-b]`。
  ⇒ **⛔ 我翻任何一格**：它是实现方，裁决者≠实现者。★**我自己复跑**的那两发（宿主台面、`GOTMPDIR` 先 `mkdir -p`）＝`internal/memory` **`RUN=78 PASS_TOP=41 FAIL_TOP=0 SKIP_TOP=2 rc=0`**、`internal/session` **`20/0/0 rc=0`**，与它自报同色。**新腿待派＝`326-v1`**（非实现者终裁六格）。
- **`120-v2`**（票 120 非实现者终裁）：判语**五格全「成立」**；件 `docs/evidence/s1/120-v2-verdict.md` **244 行／27735 字节**；两枚可重放突变脚本 `probes/120/v2/tooth.sh`（120 行）／`tooth2.sh`（44 行）。
- **`270-v1`**（票 270 非实现者终裁）：AC#1／#3／#4／#5「成立」、AC#2「不适用」；件 `270-assembly-guard.md` **109 行／19515 字节**；四笔 `13caf09d`／`fd21404e`／`7785e683`／`3c16dbc5`。
- **`165-c2`**（票 165 门铃仪器落点普查·只读）：件 `docs/evidence/s1/165-doorbell-census.md` **105 行／18148 字节**；两笔 `606bb2ff`＋`7770108d`，`.go` 改动 **0 枚**、日志 41 枚零字节 0。

### ③ 我自己复跑的关键读数（⛔ 引腿的数当"我量过"）

- ★**逐字重放 `120-v2` 的牙**（同镜像 `golang:1.27`、同挂宿主 mod 缓存、同 `GOPROXY=off`、台面＝`git archive HEAD` 的仓外树 `/d/tmp/orch-120v2-replay2/head`；突变打在物理字节、⛔ `-overlay`；读数件 `probes/orch/20261016/logs/orch-120v2-replay.txt` 与三枚原始日志同批入库）：
  `sealHandle` 摘链接面（`:300`）＝**71/41/2** ⇒ 红 `TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3`＋`TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3`；摘同一对象面（`:303`）＝**71/41/2** 红在那两枚 `…NameCarriesAnotherRealObject120r3`；**区分发**＝摘另一对同形（`sameObjectPinned` 的 `:219`＋`:222`）＝**71/39/4** 且红名**逐枚不同**（票 322 那枚＋120 的另三枚）⇒ "打的是哪一对"是读数；还原 `sha256=19f3c211…0d3e2` **三发 MATCHES**、复确认 **71/43/0/0 rc=0**。
  AC#3 那把尺我原样重放＝`-run 'TestAC[1-4]POSIX' -count=1 -v` ⇒ **39/28/0/0 rc=0**（顶层枚数尺带 `sort -u`）；AC#5＝整包 `-count=2 -v` ⇒ **142/86/0/0 rc=0**；AC#4＝`git diff 58de465a..HEAD -- internal/winsec/winsec.go` **2 hunk、+42/−1**、该区间该件只有 **1 枚 commit**＝`8baec646`（票 322、已登记 `A1001`）；能力类＝`sealHandle` 产码调用点 **1 枚**（`internal/winsec/winsec.go:158`，`PrivateFile` 与 `PrivateFileExclusive` 两条入口都经它）＋定义 2 枚＋测试 5 枚件。
- **票 270** 我自己现读产码＝`internal/tools/subagent_197.go` 的 `:265` guard／`:297` 甲形自保护／`:301` 唯一真调用点三处对号；**五句拒绝文案逐字各不相同**（`:254` 名册未接线／`:266`·`:299` 装配未接线／`:273` 没有宿主给的任务 id／`:304` 给不出装配快照／`:311` 没有审批准入钩子）⇒ AC#3"三句互不充当替身"我复认。**AC#0 由我单独裁**（裁决腿具名⛔ 射程）：料在盘＝`probes/270/a1/census.md`（208 行／41561 字节，三形代价表逐枚给"改动面 文件∶行／动不动文案／要不要新导出名"，并点明甲形**今天不可达**、且甲**不解决** `:293` 那一形）＋`probes/270/a2/pairing.md`（103 行／23759 字节）＋落地腿现跑复认顶红 **0 枚**（起手 `rc=0／PASS=235／FAIL=0`、改动后 `PASS=239`＝235＋本腿新 4 枚——★这两发是 `270-r1` 件里的读数、**⛔ 我复跑过**，我这一格只核到"两件普查文件真在盘＋三形表逐格齐＋枚数口径自洽"这一层）。
- ★**六枚 `./cmd/wisp` 红的成对归因（本条是这批最贵的一发，因为它决定后续程要不要去"修"别人的红）**：同台面两枚仓外导出树（改前锚＝`fd01e7962`＝接线 `2d10bd0f` 的父提交，改后锚＝HEAD），逐件同法补台面 ⇒ 读数 **6 枚 ⇄ 6 枚、名集合逐枚同、红句除 HWND 句柄数字外同句**（件 `orch-sixred-pair-pre.txt`／`…-post.txt`，提交 `6c20357b`）。**⇒ 票 220 那两跳接线一枚红都没引入。** 六枚归两形本机状态：① 活实例 `wisp.exe` **PID 21892**（命令行 `"D:\work\workspace\projects plans\Wisp\build\wisp.exe"`、起于 `2026-10-15 19:03`、父 PID 11728 已不在＝台账 `A988` 机主自己按重启档起的那一枚，⛔ 谁都不许杀）握热键与单实例锁打死四枚；② 前台被别的窗握着打死两枚真窗（红句自述 `the panel did not take the foreground on Show`／`the page reports its title as ""`）。`TestAC246DevLegIgnoresTheTestTaskInjection` 的红句**自己就说实话**＝`AC#7 RED (desktop state, not our code): a dev Wisp is already running in this session, so this leg handed its activation over and exited. Stop it and re-run.`

### ④ 翻框与结案（尺＝机械 sed `flip-1005.sh`，只改 `[ ]`→`[x]`、⛔ 动判据文字；翻后名册逐票现量）

- **票 120**：五格全翻 ⇒ `unchecked=0 / checked=5` ⇒ **改名 `-done`**（提交 `ddffbf09`）。残余逐枚**具名归口**、⛔ 抹平：Windows 那半格＝**票 322 射程**（现量：`322` 的 **AC#3 已勾**、只剩 **AC#4** 等整包终态读数）；枚数过期只具名记现量（票面"16"／`120-v1` 20／`120-r3` 27／`120-v2` 与我**都是 28**），改措写要人批；裁决腿那枚口径落差（票 113 那批对"换过去又换回另一枚真对象"失明）已随本票结案由 120 自家仪器＋`120-r3` 台件接住。
- **票 270**：六格全翻 ⇒ `unchecked=0 / checked=6`，**⛔ 改名 `-done`**（按住的理由＝它自己的交付缺陷未清：`:257` 那句注释 "the only site" 与同块 `:261`、`:294`–`:295` 互顶＝过期文字，微腿 `270-r2` 待派）。AC#2 翻勾的凭据是**"不适用"**、⛔ 是"验过"（本格按票面括号只在选乙时适用，本票裁甲）。
- **票 220**：五格**一枚没翻**（实现方＝死腿 `220-r2`，⛔ 自裁）；我**代它补一条票尾收腿记**（现量该件 `- [20` 开头的日期 bullet 只有 1 枚、本条＝第 2 枚；尺＝`grep -cE '^\- \[20'`）。
- **票 326**：六格**一枚没翻**（实现方自报现量 6 未勾／0 已勾，与我 `grep -c` 同色）。**票 165**：现量 **5 已勾／0 未勾**——但**⛔ 结案**：本条新开一格 **AC#6** 记"门铃仪器"（`165-c2` 量到的三处断口，AC 原文⛔ 覆盖"模型真的会不会调"这件事），新开 AC 属编排者排程动作、⛔ 改已冻结的 C/D 契约文字。
- 累计本波翻框＝**11 枚**（120 五枚＋270 六枚）。

### ⑤ 我裁的事（能自己定的⛔ 上机主的清单）

- **`ParentTools`／`BaseOptions` 那半边的弱牙不另立票**：`ParentTools` 半边已有票 236 那枚具名红作牙；`BaseOptions` 半边"甲形新分支今天不可达"是票面 §1.1／§2.1 已写明的**形状性质**，要装牙必须拆 `:256` 走乙形＝同批改掉票 236 **刚翻勾**的三枚凭据钉 ⇒ 属排程／契约取舍、不是漏账。撤销口令＝**「270 补 BaseOptions 半边的牙」**。
- **票 165 的门铃仪器选形＝甲**（新建 `internal/agent/testdata/golden/plan-present.sse`＋1 枚用例）：凭据＝`165-c2` 现量该目录**无孤儿钉**（`ReadDir` 三处全在 `t.TempDir()`）＋**已知钉 0 枚**＋⛔ 碰共享编译面；乙-1 顶名册钉（调用点恰 1 枚、只许 `run.go`）且占本波 `cmd/wisp` 唯一格 ⇒ ⛔；丙＝0 枚 `.go`、属〔仅本机可量〕＋需 key＋winlive（该腿严尺现量 **14 枚**＝`cmd/wisp` 9＋`internal/ball` 5；⛔ 拿台账旧枚数"8"，那枚复现不出）⇒ 只当补充读数、⛔ 当仪器。★它同时上交一枚**必须先裁的硬约束**：`internal/agent` 根 21 枚测试件逐枚 `package agent`、`agent_test` 全树 **0 枚**，而 `bridge.go:13` 已 import `agent` ⇒ **真桥进不去＝导入环**，落地腿只能走包内形状。
- **`165-v1` 那枚"56"我自己复跑过、结论要分开写**：56 是真读数但⛔ 恒量——`71002d0..232f309e` 32 枚逐跑＝45→54→55→**56**（10:49 `d646625e`）→…→**63**（该腿自己的锚）；"fixture 面 0 命中"这半句在全部 32 枚上成立 ⇒ 今后引"56 枚"必带锚号，⛔ 不带锚号引用。

### ⑥ 我自抓的错（记我，⛔ 抹平）

1. **导出树建树少一件**：六件套 `go.mod go.sum internal cmd tools scripts` **不够**，`internal/panel/assets.go:25` 引用 `github.com/CarlosShao/wisp/frontend` ⇒ 报 `no required module provides package …/frontend`／`FAIL … [setup failed]`，**长得像测试红**。定式＝**七件套带 `frontend`**。
2. **导出树跑不动 `cmd/wisp` 的三个缺件面**：`frontend/dist` 只有 **1 枚跟踪件**（真窗用例必然读到空白页）；`build/` 与 `third_party/sherpa-onnx/` **零跟踪件**（缺 DLL ⇒ `exit status 0xc0000135`、0.04 s、无 `--- FAIL`，或红句自述 `no native DLLs in ..\..\third_party\sherpa-onnx`）。定式＝跑之前手工补这三面并具名写补了哪几枚。
3. **容器树建在 `mktemp -d`＝假 rc=1**：Docker Desktop 把容器外那层的 `/tmp/…' 读成它自己 VM 的路径 ⇒ 挂载指到空目录、`go test` 连 `go.mod` 都看不见、且⛔ 落 `$tag.log`。定式＝**凡 rc 非 0 又拿不到日志，先证"日志该出现的那枚目录真被挂上"**；容器要挂的树建在 `/d/…` 这类它认得的卷。（第一次那发的**突变与还原**部分真跑成了 ⇒ 那几行算我的读数，四数行⛔ 算。）
4. **`./cmd/wisp` 整包默认 10 分钟超时不够**：我 11:2x 那一发跑到 `panic: test timed out after 10m0s`（600.19 s、`=== RUN` 已累计 430 枚）⇒ **那一发⛔ 算终态读数**。定式＝跑整包必须显式 `-timeout 30m` 并写清带了哪枚值；票面 AC#5"单跑不算终态"之外再加这一条。
5. **一把差两枚前缀的尺差点给我造出假 0**：`git grep -n 'BaseOptions' HEAD -- internal/tools/internal/tools/subagent_197.go`（变量里已含目录，我又补了一次）⇒ 输出空。**定式＝任何"0 命中"之前先看这条尺的路径拼法与 stderr，并跑一枚已知存在的正控**（本条我自己抓到、没落进任何判语）。
6. **拿一轮前的 `git status` 快照当现状**：我在 `220-v1` 派单里写"工作树 14 枚脏、含 `internal/**`＋`cmd/**`"——那是 11:3x 的快照、当时腿的在飞态已归位；11:4x 我现跑 `git status --porcelain -- internal cmd`＝**0 行**。⇒ 已在派单第三节更正，教训＝**进正文的工作树状态必须落笔前现跑**（同族旧错：过期前提）。
7. **另一枚协调教训**：我把成对台件的脚本放在 `probes/orch/20261016/txt/` 且**输出路径写死成我自己的证据件**，11:4x 起那两枚 `orch-sixred-pair-*.txt` 被在飞腿重跑覆盖（`pre` 24813 字节、`post` 只剩 46 字节表头）⇒ 我提交的副本（`6c20357b`）才是权威、工作树那两枚我⛔ 再提交。定式＝**台件脚本的输出路径要参数化，⛔ 让别人的复跑覆盖我已入库的读数件**。

### ⑦ 存量登记（⛔ 我动、⛔ 我提交、⛔ 我 revert；归属查不到就留）

- `go.mod` 工作树脏（`github.com/jchv/go-webview2` 从 `// indirect` 被升进直接 `require` 块），mtime **2026-10-10 22:28:50**＝六天前的存量；`go.sum` mtime 2026-10-01；仓根 0 字节未跟踪件 `-`（mtime 2026-10-03 10:20:56）。工作树现量＝**tracked 脏 28 枚＋untracked 1148 枚**（尺＝`git status --porcelain`），其中 `internal`／`cmd` 两面 **0 行**。
- 归属尺跑不出主人的都留着（同族规矩：先归属、只删可重建）。

### ⑧ 编队与队列（放腿闸门＝CPU ≤70%，机主那把尺）

- **在飞 2→4 枚**：`216-r1`（写，唯一产码面 `internal/panel/**`）agentId `ageneral-purpose-e7f81d89980842b9`；`220-v1`（裁决，票 220 五格＋`cmd/wisp` 整包那一格归它）agentId `ageneral-purpose-257b01002298979f`；本条之后补放 `326-v1`（裁决）与 `165-r3`（写，`internal/agent/**`）。
- ★**写面互斥现量**：`cmd/wisp/**` 本波**零写腿**（`220-v1` 只许新建 `cmd/wisp/*_220v1_test.go`）；`internal/panel/**` 唯一写腿＝`216-r1`；`internal/agent/**`＝`165-r3`；整包 `./cmd/wisp` 那一格只归 `220-v1`。
- **队列（具名释放条件）**：`270-r2`（过期注释微腿，写面只 `internal/tools` 注释面，按在 `216-r1` 交件后）；票 322 **AC#4** 等一枚 `cmd/wisp` 整包**终态**读数（`220-v1` 交回后我复跑）；票 216 的**"两份合一"**（`sanitizeReplyText` 与 `findMetachars`，写面撞 `cmd/wisp` ⇒ 下一波 `cmd/wisp` 空时另派）；`165-c1` 那枚 `-overlay` 幻影注释微腿＋门铃仪器第二形；`247-r1`（真麦克风电平 ⇒ `cmd/wisp`）；票 313 落点包；`259-v1`／`181-v3`／`253-v1` 的过期前提重验；票 252 残余三格；`135-r1`。
- ⛔ 推送：本波写腿未收完、ahead 已 91 枚 ⇒ 收完再推（只 `origin/dev`、⛔ force、推前核 fast-forward；推送＝抢本机 CPU，`slo-full` 自建 runner 会自启）。
