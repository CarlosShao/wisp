# 295-v1 · 复核格（AC#1／AC#3／AC#5 的独立尺 ＋ 台面归属）· 落件 2026-10-12 16:3x +08

起手 HEAD＝`a92cd10f`；读数期间别的腿入库到 `f1155198`／`d164e50e`／`8ef6b672`；
★尺＝`git diff a92cd10f..HEAD -- cmd/wisp`＝**0 行**、`git diff 421469b7..HEAD -- cmd/wisp`＝**0 行**（两把各 `rc=0` 无声）
⇒ 本腿三发的代码面自始至终＝票 295 落地那四笔的**最终形状**（`421469b7` 那版），HEAD 漂动全在别的面。

## AC#1 只动了 ⓓ 那批点＋两组注释；⛔ 第二枚生产读者被漏；`verdict` 没被顺手改

| 问 | 判 | 尺（逐字＋射程＋provenance） |
|---|---|---|
| 落点枚数对不对 | **对＝12 枚** | `git grep -n -e 'started atomic.Bool' -e '\.started\.Store(' -e '\.started\.Load()' -- cmd/wisp` ⇒ **12 行**（声明 1＋生产 Store 2＋生产 Load 1＋测试读点 8）；逐名已落 `claims-rulers.txt` 同族尺。〔验收腿现跑〕 |
| 残形（裸赋值／裸读） | **0／0** | `git grep -n -e '\.started = ' HEAD -- cmd` ⇒ **stdout 0 行／stderr 0 行／rc=1**（我先读 stderr 才报 0；`rc`＝grep 自己的码，⛔ 用 `| cat` 后的码，那次我用错过形并当场标〔作废〕） |
| 除了 ⓓ 还改了谁 | **只 4 枚件** | `git diff --numstat 0615c6c9..HEAD -- cmd/wisp` ⇒ `1 1` live ＋ `3 3` 247 ＋ `10 10` 产码 ＋ `4 4` mute_290 ＝ **18/18**；行数 425⇄425。〔验收腿现跑，⛔ 引编排者复跑〕 |
| 断言有没有被放宽 | **⛔ 一枚都没有**（这把是本腿新尺） | 把新件里每枚 `ra.started.Load()` 复原成 `ra.started` 再与改前 blob `diff`：`resident_audio_247_windows_test.go` ⇒ **rc=0（零差）**；`resident_audio_247_live_windows_test.go` ⇒ **rc=0**；`resident_mute_290_windows_test.go` ⇒ **只剩 1 行差＝`:142` 的 `open=%v started=%v`→`open=%t started=%t`**（实参两版都是 `ra.started`）⇒ 全部差异＝形状＋两枚动词，⛔ 条件极性／阈值／删断言。读点枚数尺：旧 `.started` **4/3/1** ⇄ 新 `.started.Load()` **4/3/1** |
| `posture()` 之外第二枚生产读者 | **没有＝1 枚** | `git grep -n 'posture()' HEAD -- cmd internal` ⇒ 16 行，剥 `_test.go` 后 **3 行**：`:103`（注释）＋`:375`（定义）＋`cmd/wisp/resident_windows.go:385`（★唯一生产调用＝那句 `fmt.Printf("wisp: 采集腿：%s\n", …)`）。再补一把更宽的：`git grep -n '\.started' HEAD -- cmd internal` 剥测试 ⇒ 本字段只有 **3 处**（Store `:168`／Store `:323`／Load `:379`），其余命中是别家的同名物（`internal/ball/sta_windows.go` 的 channel、`internal/llm/openaichat/toolasm.go` 自己的 struct）。⚠ 票 293 那枚 `trayMuteState()`（`:420`）现读＝只回 `ra.gate.Muted()`，**⛔ 读 `started`** ⇒ 它不是被漏的第二读者〔验收腿现跑〕 |
| `verdict` 有没有被顺手改 | **⛔ 零字节** | `git diff 0615c6c9..HEAD -- cmd/wisp \| grep -cE '^[+-].*verdict'` ⇒ **0**（0 旁边站分母＝那次 diff 的 `+`/`-` 共 36 行）。AC#0 前件独立复跑：`toggleMute` 体内（`:156-:196`）`verdict` 命中＝**1 行**且是 `:189` 那句**文案字符串**里的"采集腿 verdict"字样，⛔ 赋值；七枚 `ra.verdict =` 全在 `assembleCapture`（`:263/:284/:298/:313/:320/:324`）＝boot 线程 ⇒ "手势路径不写 verdict"**成立**，票面"若写就同批"那一支**未被触发**＝腿只改 `started` 是对的〔验收腿现跑〕 |
| ⛔ 新造规矩／加锁／加协程 | **成立** | `sync/atomic` 两版都在 `:74`（`git diff … \| grep -E '^[+-].*(import\|"sync)'`＝**0 行**）；`atomic.Bool` 改前锚 `0615c6c9` 的 `internal cmd` 射程＝**5 处**、HEAD＝**6 处**＝5＋本票那一枚；`go func(`／`sync.Mutex` 在本件＝**0 命中**（尺同下 d22scan 那一层，编排者已跑 rc=0，本腿⛔ 重复）〔验收腿现跑＋编排者现跑〕 |

## AC#3 越界：换一把尺复算 ⇒ **成立（十枚禁列路径 0 命中）**

- 腿交的尺＝`git show --name-only --format=`；本腿换＝逐笔 `git diff-tree -r --name-only --no-renames --format= <c>`（四笔分别跑，⛔ 一次传四枚——那一形本腿试过，**输出 0 行＝哑尺**，已当场标〔作废〕改回逐笔）。
- 并集＝**38 枚**不同路径；其中 `.scratch/wisp/probes/295/r1/**` 34 枚 ＋ **非 `.scratch` 仅 4 枚**＝`resident_audio_windows.go`／`resident_audio_247_windows_test.go`／`resident_audio_247_live_windows_test.go`／`resident_mute_290_windows_test.go`。
- 十枚禁列逐枚过（`grep -c -F` 对那 38 枚）：`frontend/`／`design/`／`PLAN.md`／`docs/specs/`／`internal/observe/thresholds.go`／`golden`／`tools/d22scan/allowlist.txt`／三枚冻结件 ⇒ **每枚 hits=0，分母 38**（0 全部站在分母旁边，且 `grep` 的 0 与"件不存在"用 `wc -l` 双尺分开）。
- 逐笔（⛔ 只报并集）：`d5ef9cb9` 19 枚（15 probes＋4 件）／`35744137` 6 枚（全 probes）／`1c1d7df3` 4 枚（3 probes＋1 件）／`421469b7` 15 枚（全 probes）⇒ 每笔的非 probes 面都在那 4 枚 `cmd/wisp` 件之内，与派单"只碰它派单那 4 枚"同值。

## AC#5 race 代价表：素材复跑通过、★**结论句越界**（具名，⛔ 我改它一字）

- 件＝`logs/60-ac5-race.txt`（**360,724 字节／3,073 行**，腿自落 `RACE_RC=1`）。
- 腿的尺 `grep -c 'resident_audio'`＝**0**（分母 3,073 行）⇒ 复跑同值。本腿再扩四把词尺（同一件）：`toggleMute`＝**0**／`muteGesture`＝**0**／`atomic.Bool`＝**0**／`posture`＝**6**（★这 6 枚⛔ 在栈里：逐行看＝`:86/:87/:89/:530/:542/:554` 的 **slog 属性名** `"posture"`，其中三枚是 `audio: capture device unavailable` 那几行的字段）。
- 块级尺（本腿的独立形＝按 `WARNING: DATA RACE` 切块，python 逐块判）：**23 块**；含 `resident_audio` 的块＝**0/23**；含 `toggleMute`＝**0/23**；含 `posture`＝**0/23**；含 `audio`＝**0/23**；含 `atomic`＝**0/23**。
- ⚠ 两块本腿新尺要具名给编排者，免得"0 命中"被读成"什么都碰不到"：① 23 块里 **6 块**含 `gate.go` 帧，逐字看全是 **`internal/agent/approval/gate.go:536/:554/:631/:727/:729`**（审批门，⛔ `internal/audio/gate.go`）；② 4 块含 `started`，逐字全是日志属性 **`started_at=`**（memory schema migration 那几行），⛔ 本票那枚字段。⇒ 腿的"与本票 `started` ⛔ 相关"方向**对**，但它的尺只测了 `resident_audio` 一个词。
- 栈顶口径（小瑕，具名）：腿写"23 块栈顶全在 `auditf()`"。本腿量到**最内层**帧＝`bytes.(*Buffer).tryGrowByReslice()` **11 块**／`bytes.(*Buffer).Write()` **10 块**／`runtime.slicecopy()` **2 块**；`auditf` 则是**每块都有的最外 cmd/wisp 帧＝23/23**（含 `auditf` 的块 23、不含 0）。⇒ 那句话按"cmd/wisp 那一层全在 auditf"读＝成立；按字面"栈顶"读＝与量到的形状差一层。
- 红名册：`--- FAIL` 顶层 **8 枚**（腿写 8＝同值），逐名＝真窗族 3 枚（`TestAC4FocusReturnToPriorWindowGap33r5`／`TestAC13ColdStart…`／`TestAC14GoSideEvalPush…`）＋ `race detected during execution of test` 引出的 4 枚（`TestRunPacket…` ×3＋`Test197Subagent…`）＋ `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`（腿自己在 §2 已具名为审批超时 flake）。
- ★**本格唯一的直接答复**（"改形状这件事到底被仪器看得见吗"）：**⛔ 看不见**——23 块 race 里 0 块碰 `resident_audio`/`toggleMute`/`posture`/`atomic`；`-race` 那一发⛔ 造不出"手势线程写 ⇄ boot 读"的重叠（今天没有任何用例把这枚字段放到两枚线程上）。⇒ 能咬这枚字段的只有**形状变异仪器**（本腿 M1/M2 那两发）。
- ★**越界的结论句**（具名报回，⛔ 我改）：`30-gates.md` §3 反面一句写"本票那枚 `started` 的形状**在 race 探测器前不需要新用例来证明**"。前件（"今天没有仪器看得见它"）本腿复跑＝真；后件⛔ 由它推出——票面 `AC#2` ⓒ 说的是"**若**本票要造两线程对拉 ⇒ 必须自带正控"，那是条件式，⛔ "不必造"这一读法；而本腿量到的正是"⛔ 任何仪器看得见"＝**造一枚自带正控的对拉用例仍是唯一能让这格被门禁照到的路**。判语：那张表的**素材层成立、结论层越界一句**。归编排者裁是否要我把这句退回给落地腿改（本腿⛔ 改别人件一字）。

## 台面与归属（起手三问 ＋ 家族尺）

- 起手三问 16:04–16:11 现跑：① `git status --porcelain -- cmd/wisp`＝**0 行 rc=0**（16:28 又复跑一遍＝同值）；② 别腿⛔ 跑那一包的颜色＝`testexe_count=0`（分母 `all_process_count=349`／`with_executablepath_count=184`）、`goexe_count=0`；③ 编排者自己那发已入库（`orch-295r1-cmdwisp-shot1.txt`）⇒ 台面空。⚠ 三问的尺我自己用错过一形：`Get-CimInstance` 那两枚 filter 被 bash 的引号吃掉单引号 ⇒ **Invalid query＋`.Count` 打 0＝假 0**，stderr 明写报错；改成品名册脚本 `-File` 形后才有可信读数（件＝`procnames.ps1`，我把它留在专用根，⛔ 进仓）。
- 进程归属尺（按派单 §5 的更正用 `@(…)` 再 `.Count`，且**根名在 `ExecutablePath`**）：我停掉一枚在飞的变异发时，凭 `pid=28360 exe=C:\Users\swq\AppData\Local\wisp-test-gotmp\295v1-…\w-pfull\go-build…\b001\wisp.test.exe` 认定＝**我自己的**；同屏别腿在跑 `./internal/winsec/…`／`-list ./internal/…`（⛔ 碰 `cmd/wisp` 颜色）⇒ ⛔ 记到我头上、⛔ 抢那一格台面。
- 专用根＝`C:\Users\swq\AppData\Local\wisp-test-gotmp\295v1-20261012-160949-1256`（C: 卷，⛔ 换卷；每发先 `mkdir -p` 再指 `GOTMPDIR`——第一族三发因根不存在而 `go: creating work dir: GetFileAttributesEx …` **rc=1**，当场标〔作废〕，落点见 `gates-restore-proof.txt`）。
- ⛔ push、⛔ 翻框、⛔ 动票面原句；变异＝`-overlay` 载仓外副本（⛔ 动过任何一字节跟踪文件：四枚件 `git status --porcelain -- cmd/wisp`＝**0 行**、工作树 sha256 ⇄ `git show HEAD:` sha256 逐枚相等）。
