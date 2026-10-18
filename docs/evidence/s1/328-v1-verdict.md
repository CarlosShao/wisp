# 328-v1 判语件 —— 票 328 五格（非实现者终裁腿，⛔ 翻框／⛔ 改判据／⛔ push）

腿＝`328-v1`（非实现者终裁腿·第二程；第一程死于日限额，只交出起手笔 `c59d0354`）。
裁决锚＝**`5db0b538`**（入场时 10:0x `git log -1` 现认 `c11163c8`，10:13 前又进两笔 ⇒ 我**固定**裁决锚并把树与只读尺钉到同一份字节；尺＝`tree_roster_blob=anchor_blob=2da6c596369c1ee38672ae2656a7d5a3ecbaabfb`，逐字 `logs/build-trees-v2.log`＋`logs/attribute-driver-raw.log:2`）。
落地腿＝`328-r1`（末笔 `46d815d9`、Go 锚 `f705b082`）；代表尺现跑 `f705b082..HEAD -- internal cmd`＝**0 行** ⇒ 我裁的就是腿交那批字节。
写面＝`.scratch/wisp/probes/328/v1/**`＋本件＋票 328 尾一条 bullet；`.go` 命中 **0**（尺＝逐笔 `git show --name-only` 全集，见 §六）。
★落点新铁律（机主 10-18 10:1x 亲口＝盘符根目录零污染，覆盖派单台面配方两行）⇒ 两棵仓外树新建在 `D:\tmp\wisp\20261018\328v1\{before,after}`，`GOTMPDIR` 同根下；旧树全部迁 `D:\tmp\wisp\legacy\`（编排者 10:3x），我一枚没复用、一枚没删。

## 一、AC#2 判语：**过**（附两处具名界）

1. **正控独立复跑成立**（腿那句"经真路径"我逐跳现量，⛔ 抄）：`callCorr` def `internal/agent/loop.go:603` ⇄ 唯一生产调用 `:676`（`ToolRequest{` 构造 `:675`，非测试全仓唯一）⇐ `executeCalls` 调用点 `:504` ⇐ `run` `:346`；搬运 `panel/pump.go:98 → approval.go:72`；名册 `subagent_roster_197.go:196 → :207`；第三臂 `:312`。颜色＝`after-agent-positive: RUN=2 PASS_TOP=2 FAIL_TOP=0 SKIP=0 rc=0`（两发逐名 PASS）。
2. **★本票关键敏感点由我的读数钉住**（派单点名那一发，落在我自己树内的物理字节，内容锚＝函数体内那枚 token）：`strings.HasPrefix(key, task+"#")` → `strings.HasPrefix(key, task)` ⇒ `mut-V1-boundary: RUN=7 PASS_TOP=0 FAIL_TOP=1 SKIP=0 rc=1`，红的正是 `subagent_blocked_220_test.go:142` 那枚 `child-1x`（红句逐字 `:147: row "child-1" blockedOnApproval=true，期望 false； packets={"pending":[{"correlationId":"child-1x",…`），**同场其余五枚全绿**。⇒ 「带边界」这个限定词有**单独**仪器敏感性，丙-带边界被钉住；不需触发派单那句"退回 AC#1 那一半"。
3. 成对第二发落点不同枚（`:313 set[task] = true` → `_ = task`）⇒ agent 正控红（`:124` 带真 UUID corr）＋ `:165` 两支红；panel 侧只红 `ⓑ` 那枚（`child-1#9`／`#147`／`other-1#9`／`child-1x` 全绿＝数字臂没被碰、两臂各自承重）。两发**当场还原**，前后各断言一次 pristine：`sha256=450e1d1c47…bd87028`／`gitblob=2da6c596…`＝锚那枚（`restored_same_as_pristine=YES`，同一 shell）。
4. **★具名界①（派单问腿没答那一层，我按盘上把话说清）**：corr 那半边＝生产铸造（上面那条链）；**任务号（"在册"）那半边＝夹具**。尺＝`TaskRosterState{` 非测试构造点 **0 枚**（含 stderr 0 字节；正控＝全仓 10 命中逐枚点名全是 `*_test.go`），生产真装配者是 `cmd/wisp/panel_pump.go:191 taskRosterState()`（`:233 rt.tasks.Look` ⇄ `:240 Descendants` 塞行），而正控用 `corr_roster_join_328_test.go:43 rosterWithTwoRows(res.TaskID, …)` 代劳。⇒ 判语一律写作「**corr 半边生产真形＋在册半边夹具代劳**」，⛔ 写成"整条链是生产真形"（腿件①那句"⛔ 测试里硬写"仍成立：夹具硬写的是**在册**那半边，⛔ corr）。
5. **★具名界②**：`cmd/wisp/**` 本波零写腿但我⛔ 跑那包（唯一持跑者＝`228-r3c`）⇒ 那包三枚 `false` 断言（现量行号 `subagent_blocked_197_test.go:190/:223`、`subagent_carrier_197_test.go:311`）**没有改后颜色读数**。我补的是静态界（量的）：`git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' 5db0b538 -- cmd/wisp`＝2 命中、逐枚都是注释里的 `"AC#0"`／`"AC#3"`，卡片 corr 带 `#` 者＝**0** ⇒ 第三臂不可能在那一包点亮新行；**这是推算的界，不是颜色**（正控＝同尺族 `internal/panel` 9 行）。

## 二、AC#3 判语：**过**（判定与批准语义零字节，四面我自己打）

- 那 4 枚 `.go` 逐枚现量（`git diff -U0 --numstat 648eda04..5db0b538`）＝`subagent_roster_197.go` **+66/-10**（新非导出 `rosteredTaskIDs:276`＋`indexWaitingKey` 加一枚 `known` 入参＋第三臂三行 `:311-314`＋两处散文）／`subagent_blocked_220_test.go` **+12/-4**（`:141` 一枚表项翻向＋表头散文）／两枚 `*_328_test.go` **新增** +89／+171。
- ★**数字臂一字未动**不是我抄的：`splitClashSuffix` **函数体** sha16 两侧同一串 `7fec699e1a0f00a7`（起手锚 ⇄ 我的锚），def `:333`／doc `:318`／唯一读者 `:308`。
- 门控档位⇄L2 许可面 blob 对撞（起手 `648eda04` ⇄ 腿末笔 `46d815d9` ⇄ 我的锚 `5db0b538`）＝`risk/mode.go c21beeea`／`risk/provenance.go d634aae0`／`approval/gate.go 78e2d28a`／`approval/queue.go 66fec7ae`，四枚**三锚同一串＝ZEROBYTE**；整面 `diff -U0` 行数＝`internal/risk` 0／`internal/agent/approval` 0／`internal/statemachine`（D43 那面）0。
- 零命中带正控：三枚冻结件各 0 行（正控字节 19175／109415／15659）；`thresholds.go`／`testdata/**`／`go.mod`／`go.sum`／`ci.yml` 各 0 行（正控＝ci.yml 64794 字节、golden 名册 66 枚）。
- ⛔ 新导出名／⛔ 新 JSON 键（我这把尺）：panel 导出 `func` 名册两侧**作差 0**；四枚 `+` 行里唯一的顶层导出标识符＝五枚 `Test*` 用例名；`json:` 标签名册**作差 0**；`ApprovalCardView` 里 `TaskID` 命中仍 **0**（卡片仍 10 枚 `json:` 字段）⇒ 没撞语言层、没开契约面。`⓭` 那句自陈现取逐字已改到说实话（"stopped being true the day loop.go:603 started minting them"）。

## 三、AC#4 判语：**过**（成对复跑＋逐名作差＋门禁四数）

- 四数逐枚点名（⛔ 混用⛔ 相减），两把临时根台面同配方同机 `-count=1 -timeout 30m`：`./internal/panel` 改前 **257/160/7/0 rc=1** ⇄ 改后 **268/163/7/0 rc=1**；`./internal/agent` 改前 **112/88/0/0 rc=0** ⇄ 改后 **114/90/0/0 rc=0**。⇒ **腿那四个数我逐枚复现成立**（＋11 RUN／＋3 PASS_TOP＝它新加的用例数，⛔ 当判据）。
- ★红名册**逐名**作差（派单点名看名字不看枚数）：`grep '^--- FAIL' | sed 去时长 | sort -u | diff` ⇒ **REDROSTER_DIFF_RC=0，两侧各 7 枚同名**（`logs/red-{before,after}-panel-rerun2.txt`）。七枚逐名＝`approval_test.go:129` 键族／`composer_test.go:78` 键族（两枚起手即红、与卡片无关，本票⛔ 加键 ⇒ 没多出第二句，复认）／两枚 35r8 堆量尺／`composer_test.go:873`／`git_test.go:600`（缺 `.git`）＋`tokens_fourway_test.go:441`（缺 `design/`＝机主外部 agent 写面）＝两枚台面局限，改前改后同红、作差互抵（四枚逐枚复认）。
- ★本腿抓到并归位一枚红名抖动：改后侧**第一发** `FAIL_TOP=8`，多出 `TestFloodRulerFiresOnAnUnboundedSink35r8`（`:369` 堆量句）。归位四把尺＝① 那文件两侧 blob 同串 `3e5c0e81`（`diff_lines=0`，⛔ 328 射程）② 单跑两侧各两发全绿 rc=0 ③ 改后整包**复跑**该枚 `--- PASS` ④ 复跑两侧名册逐名作差＝0 ⇒ 判**同场堆量分配的台面抖动**，⛔ 计入 328 判据、⛔ 据此退腿。
- 门禁四数（只在共享工作树，逐枚 `rc=`）：`gofumpt.exe -l <腿那四枚>`＝**rc=0**（无输出）；`sh scripts/d22scan.sh`＝**第一发 rc=2＝GUARD 射程拒绝**（offender＝我的 `GOTMPDIR=D:/tmp/…` 与 `%LOCALAPPDATA%` 不同卷＝票 311 AC#1 shape B 那道闸；⛔ 读成红）⇒ 按其自陈 unset `GOTMPDIR`/`WISP_TEST_GOTMP_ROOT` 重打＝**clean rc=0**（射程现量 ban#8 internal 586／cmd 143／design 39＝非空）；`check-path-length-budget.sh`＝**VERDICT GREEN rc=0**（分母 15217 tracked paths，腿报 15035 随仓在长，两数都带尺）；`go vet ./internal/panel/ ./internal/agent/`＝**rc=0**。跑门禁前现量：共享树两包与我的锚**差 0 文件**、`status -- internal/panel internal/agent cmd/wisp`＝空。

## 四、必答①（队列撞名"叠改"那半边只有夹具形，算不算本票欠修）

**裁：不算票 328 射程内的欠修**，另立新票。凭据：票面 `AC#2` 那句话要求的是「经真 `callCorr` 拼出的 corr 走真名册点亮那一行」——那一跳本票已钉（§一 2、3）；叠改形 `T#<callID>#<seq>` 的**第三臂侧**已由 `subagent_corr_task_prefix_328_test.go` 实测点亮（我复跑逐名 PASS：`arm B with the queue clash rewrite stacked…`），缺的是**队列侧真产那一形**这一格，而它需要跨包把两枚非导出铸造点链起来（`loop.go:603 callCorr` ⇄ `queue.go:156 Queue.push`）＝**要新导出名**＝语言层越界，⛔ 属本票判据（票面边界那条明写"⛔ 新导出名 ⇒ 停手上报"，腿如实停手）。
建议新票：**标题**「队列撞名改写的 corr 叠改形（`<corr>#<seq>`）全仓无真产用例——把 `callCorr` 与 `Queue.push` 链成一发端到端读数」；**那一格判据**「一发经真队列改写出 corr 的用例：corr 逐字节含两枚 `#`、剥出的段仍是生产 corr，且名册点亮原任务行；自带去牙＝把 `queue.go:168` 的 `%s#%d` 改回不叠 ⇒ 该用例必须红」；**为什么要新导出名**＝两枚铸造点分属 `internal/agent` 与 `internal/agent/approval`、都非导出，测里要拿到"真队列改写真 corr"只有开一枚跨包入口（或把该接缝做成测试注入口），那是**契约/语言层动作**，须人工批准。★⛔ 把这条读成"叠改形已被钉住"：腿 M2 里 `child-1#9`／`#147` 仍绿＝**数字臂**在承重，⛔ 是"队列改写＋第三臂"合起来那形（我的 M2 复跑同读数）。本腿⛔ 动票池（写面不含新票件）。

## 五、必答②（那一句定性词不重于证据）

**裁：AC#2 正控买到的是 Go 侧接线，⛔ 用户可见真相。** 尺与射程（我的锚 `5db0b538`，量的）：`git grep -l 'blockedOnApproval' -- frontend`＝**0 枚文件**、`-- frontend/src`＝**0**、大小写不敏感亦 0；★正控＝同尺族 `correlationId` 命中 **11 枚文件**（`frontend/src` **7 枚／24 处**）——腿对票面与普查件"4 枚"那处的更正我复跑成立。旁证两枚自陈现量行号：`cmd/wisp/panel_snapshot_deliver_329.go:28`（`window.wispOnPanelSnapshot does not exist yet`）＋`cmd/wisp/subagent_carrier_197_test.go:598`（逐字 `blockedOnApproval is never asserted on a production packet`）。⇒ 屏幕上那一行**今天仍不会亮**；那半格＝接线票（票 329／前端会话），`frontend/**` 归机主外部 agent＝本编队零写面，⛔ 记在票 328 头上。

## 六、剩下的零读数／未证之事（具名，⛔ 当已做）

- `cmd/wisp` 那包无改后颜色（§一 5，只有推算的静态界）。
- `splitClashSuffix` 仍**无任何测试直呼**（我的锚现量＝def＋doc＋唯一读者三命中，测试直呼 0；正控＝`indexWaitingKey` 8 行）⇒ 行为级仪器（`child-1x`／`other-1#9`／四形表）是唯一凭据，直呼那格仍零。
- 叠改形真产＝零（§四，尺＋正控已带）。
- 陪聊全双工／realtime 到不了那一格＝普查件 `approval/gate.go:229-241` 的 fail-closed 散文，本腿**没重跑**那一把尺（只读腿 337-a1 在飞 risk 链），⛔ 复用其读数当我的凭据。
- 面板侧 `BlockedOnApproval` 在真 packet 上的仪器＝零（归票 329）；改后整包复跑那一枚堆量抖动**未消除**，只归位（下轮若再抖，仍⛔ 算进 328）。

## 七、两处对编排者的具名上报（⛔ 我自己填假设／⛔ 据此退腿）

1. **复认第一程 `c59d0354` 那处派单前提与盘上不符**（我自己现跑，⛔ 抄它）：派单 `AC#2③` 说"把第三臂换成 `strings.HasPrefix(key, task)` 这一发＝腿没发的反形"，而腿 `logs/mutate.sh:31` 的 **M1 替换串与那一发逐字节同形**（`sed 's/strings.HasPrefix(key, task+"#")/strings.HasPrefix(key, task)/'`）。我仍自己现跑（§一 2）＝判据按我的读数成立；⛔ 据此退腿那一发。⇒ **建议派单模板改一句**：这类"补反形"先对 `git show <腿末笔>:…/mutate.sh` 现量再写。
2. **新落点铁律与共享树门禁相冲**：`GOTMPDIR` 只许在 `D:` ⇒ 共享工作树里的 `scripts/d22scan.sh` 触发票 311 AC#1 shape B 的 GUARD（`rc=2`，射程拒绝，⛔ 读成红）。本腿用"跑门禁前 unset 两枚 env"绕开（逐字 `logs/gates-recheck-nogotmp.log`），但那是**每台机每次都要记的一手**，属编排者/机主那一层的定案（⛔ 属本票）。

## 八、交件自证三把尺（逐字，入库笔见下）

1. **写面名册**：尺＝逐枚 `git show --name-only <我的每笔>`，全集 ⊆ {`.scratch/wisp/probes/328/v1/**`, `docs/evidence/s1/328-v1-verdict.md`, 票 328 尾一条 bullet}；★`.go` 命中 **0**＝同一把尺数出来（逐笔文件名以 `.go` 结尾者枚数 0）。
2. **证据件字节／行数**：`wc -c`＋`wc -l` 逐枚现量入库（`logs/deliverable-sizes.txt`），本件 ≤130 行。
3. **空槽尺（自避分词形）0 命中**：逐字命令＝`grep -nE '([T]ODO|[F]IXME|[T]BD|待[补]|占[位]|X{3})' docs/evidence/s1/328-v1-verdict.md`（射程＝判语件一枚；写成本文件正文里这个形状、⛔ 匹配自己）＝**命中 0**；★正控＝同一串尺打到 `internal/agent/approval/ticket146_liveapprovals_backing_test.go` 命中 `:296`（一枚八字母 X 串的夹具）⇒ 先证这把尺会响，再信它的 0。逐字读数＝`logs/deliverable-sizes.txt` 末段。

本腿提交（严格前缀尺＝`git log --grep='^票 328 终裁腿 328-v1·' --oneline`）＝起手 `ce16d7bf`／AC#2 `c8e741a8`／AC#3 `8ba2ae09`／AC#4 `45b474ec`／本件＋票尾另两枚。⛔ push／⛔ 翻框／⛔ 改票名／⛔ `add -A`／⛔ `commit -- .`／⛔ amend·reset·rebase·stash·checkout .·clean／⛔ 仓内 worktree；临时件只建不删。
**判语汇总**：`AC#2` 过（具名界①corr 半边生产／在册半边夹具，界②cmd/wisp 无颜色）·`AC#3` 过·`AC#4` 过·必答①＝不算本票欠修＋建议新票·必答②＝买到 Go 侧接线⛔ 用户可见真相。**五枚框一枚没翻**（翻框归编排者）。
