# 票 267 — `[risk].confirm_timeout_sec`／`l1_window_sec` **没有任何上下界校验**就进了审批门：种 99999 ⇒ L2 卡可挂 **27 小时 46 分**；种 10 ⇒ `Timeout()-WarningLead()` 为负 ⇒ **C18"超时前醒目提示"在那一发常驻进程里静默消失**

**立票时刻**：2026-10-04 15:4x +08，锚点 HEAD `c528035f`（`dev`）
**来路**：验收腿 `256-v1`（非实现者）在验票 256 时**现跑种子量出来的**（`.scratch/wisp/probes/256/v1/verdict.md` §2②），编排者 15:4x 自己复认了机制两枚（见下面 §现量 2/3，⚠ 编排者复认的是**读码**，种子读数归 `256-v1`）。相邻不是重复：**票 256** 管"常驻腿读不读得到 `[risk]` 这两枚键"（已由 `256-r1` 读到、`256-v1` 判 AC#1/AC#2 未闭合），**票 255** 管"配置说了生效而没人读"（方向相反：那一枚是**没读者却声称生效**，这一枚是**有读者但不问范围**）。**本票只管一格：值域没有门。**

## 现量（⚠ 引用前先重跑，别把这几行当常量）

1. **默认值就是 300**：`internal/config/schema.go:450` 逐字 `ConfirmTimeoutSec int \`toml:"confirm_timeout_sec" default:"300"\``；同节 `l1_window_sec` 的 default 是 `2`（`256-v1` §2 那枚钉 `resident_approval_risk_256_windows_test.go:198-208`/`:253-255` 断的是"等于 schema 默认"，⛔ 不是"随配置变"）。
2. **没有值域校验**（编排者 15:4x 现跑）：尺＝`grep -n "ConfirmTimeout\|L1Window\|Risk\." internal/config/validate.go` ⇒ **只有 `:110` 一枚命中**，那是 `risk.ParseMode(c.Risk.PermissionMode)`，**与这两枚数值键无关** ⇒ 这两枚键**今天零校验**（无上界、无下界、无 0/负数防护）。
3. **不钳位、原值进门**（编排者 15:4x 现跑）：`cmd/wisp/resident_approval_windows.go:310` 逐字 `time.Duration(c.Risk.ConfirmTimeoutSec) * time.Second`（`run.go:616` 同形）；门那一侧 `internal/agent/approval/gate.go:528` 逐字
   `if lead := g.q.Timeout() - g.q.WarningLead(); lead > 0 {` ⇒ **`Timeout() ≤ WarningLead()` 时 `warn` 保持 nil**，C18 那句提示**不发**（`gate.go:33-35` 注明 `WarningLead` 就是"提前多久醒目提示"；`gate.go:562` 那枚 `Remaining: g.q.WarningLead()` 是同一枚常量的另一处用法）。
4. `256-v1` 的种子读数（**归它，编排者未复跑**）：种 `99999` ⇒ 门 `Timeout()=27h46m39s`；种 `10` ⇒ `lead<0` ⇒ 提示消失（它具名"45s 时仍在"＝只有低于 `WarningLead` 那一支才消失）。⚠ 顺带它的对照读数：**L1 窗口那一枚是钳位的**（99→3s、1→2s、0/-1→3s＝带内最大），**超时这一枚不钳位**——同一次落地里两枚键行为不一致。

## 为什么会发生（一句机制）

不是记账坏了：**"能不能配"这件事先落了地，"配成什么算离谱"从来没落地**。写侧把一枚裸 `int` 直接乘成 `time.Duration` 交给门，而门只问"这个时长减掉提示提前量还剩多少"，剩负数就当没有提示——**它没有拒绝、也没有报警，只是安静地把保护撤了**。

## 要建什么（⛔ 形不在票面写死，归 AC#0 的普查）

- **ⓐ 值域进门**：在 `internal/config/validate.go` 给这两枚键配上界／下界（越界＝拒载并说清哪一枚键超出哪个范围），语义上＝"配置不合法"而不是"程序自己决定无视配置"。
- **ⓑ 消费侧钳位**：在门／装配那一侧钳（`gate.go` 或 `resident_approval_windows.go:310`）。⚠ 这一形会把"配了但不生效"重新变成常态——**正是票 255 立案要治的形状**，所以走 ⓑ 必须先回答"这算不算第二枚主动误报"。
- **ⓒ 只补可见性**：不改值域，只在超出 `WarningLead` 时**响亮地**告诉用户"这一发没有提前提示"。⚠ 它治不了 27 小时那一支。
- ⚠ **三条我都没跑作用面**，一律〔待验〕，⛔ 不许据此排除任何一条（第 74／83 条）。

## 判据（AC 框由编排者翻，产码腿一枚都不许碰）

- [x] **AC#0（只读普查，第一格）**：现跑答四问——① `WarningLead` 今天由谁给、生产里是不是恒为那个常量（具名 `file:line`）；② 全仓还有哪几枚数值型 `[risk]`/`[hotkey]`/`[app]` 键**裸着进消费者**（枚数自己拉，尺用 `| wc -l` 收尾）；③ 既有钳位那形（L1 窗口）是在哪一层做的、能不能复用同一层；④ 有没有任何一枚既有尺在断"配置值域"（⛔ 不许照抄本票"零校验"那句，自己拉）。
- [x] **AC#1**：种一发**越界值**（上界与下界各一发）⇒ 指名用例必须红，且红句要**逐字念出是哪枚键、超出哪个范围**；⛔ 判据换成反形（随便换个越界值就不响）还全绿＝不敏感，不许当凭据。
- [x] **AC#2**：**C18 那枚提示的存活**要有一发正向钉：`confirm_timeout_sec` 取一个**小于 `WarningLead`** 的合法值时，要么响、要么按所选形把提示保住——⛔ 不许"静默把保护撤了"继续是合法终态。
- [x] **AC#3**：⛔ **不许顺手把 C18 的 300s 常量改掉**（`docs/PLAN.md` 的 C18 属冻结契约面）；本票只处理"可配之后值域没人管"。**"配置值 vs C18 写死 300s 谁优先"这一格归 `Q-77`（待机主一句话）**，任何腿不许自行裁定。
- [x] **AC#4**：越界检查——`git diff` 出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` 任一路径 ⇒ 直接退回。卫生四门读数不扩大，红名集合逐名比对并**带取数时刻**。

## 排程与禁区

- **现在不按排**：`internal/config` 此刻由产码腿 `257-r1` 在飞（票 257 的形 ⓒ），**同包一律串行**，不接受"改的是不同文件"这种推理。AC#0 那枚**只读普查腿**可以现在派（⛔ 禁跑任何 go 命令）。
- ⛔ 本票不许动会话授权那条链（票 265／224 的地界）；不许动票 245 那条"稳态不绑、只在确认那两三秒借"的时机与存在性。
- ⚠ `internal/config/unwired.go:121` 那句 `"risk.confirm_timeout_sec": "consumed: cmd/wisp/run.go builds the approval timeout from it"` 现在**只说了 `run.go`**——票 256 落地后**常驻腿也消费它**了 ⇒ 这一行属过期指认，**落地腿同批改释**（⛔ 不单独立票，见 `A603`）。

## 编排者更正（2026-10-05 09:1x）：§现量 3 那枚行号**已漂**（⛔ 原句不抹，按这一条读）

- §现量 3 逐字写着 `cmd/wisp/resident_approval_windows.go:310` 是 `time.Duration(c.Risk.ConfirmTimeoutSec) * time.Second` 的落点——**现在的树上 `:310` 落在注释块里**（09:0x 现读 `:308-312`＝票 256/265 那段 "…two of the ten fields are now passed FOR [risk]" 的注释）。真身＝**`:460`**（尺＝`grep -n 'ConfirmTimeoutSec) \* time.Second' cmd/wisp/resident_approval_windows.go cmd/wisp/run.go` ⇒ `:460` 与 `run.go:616`）。⇒ **`run.go:616`／`gate.go:528`／`gate.go:562` 三处未漂**，只有本票自己那一枚漂了。
- ★ **为什么会漂（记我）**：票 265（ⓐ-Ⅰ）在同一段文件里加了 100 多行 holder ⇒ 行号整体下移，而我 15:4x 写这行时是照 `256-v1` 的读数抄的、**没有走 `grep` 定式**。⇒ 定式补一条：**票面里凡是"`file:行号`＝某行代码"的指认，一律用"被指的那串字符"当尺去 grep，别抄行号**（与 09-30 那条短引坑、10-04 那条"walk 调用点不读注释"同族）。

## 编排者裁定（形）＝**ⓐ 值域进门**，2026-10-05 09:1x，锚点 HEAD `21bec8a1`

⛔ 这一节是我（编排者）裁的，普查腿 `267-a1` 按派单"只摆料不裁形"，正确地没裁（凭据＝`.scratch/wisp/probes/267/a1/census.md`，81 行／12,579 字节／占位 0；该件由我**代落盘**——派单选错了腿型，只读会话没有写工具）。

- **裁 ⓐ 的理由（三条，都出自它的现量，不是我另造的）**：① 票 255 整张票立的就是"配置说了生效而没人读／配了不生效"这一族，ⓑ（消费侧钳位）＝**把同一形状再造一遍**，⛔ 我不在本票里给自己立新欠账；② ⓐ 落在**已有先例那一层**——`internal/config/validate.go` 今天已断 3 枚数值域（`ball.size [44,72]`／`ball.opacity_idle [0,1]`／`cost.alert_threshold [0,1]`），且有同名测试尺 `validate_test.go:155 TestValidateBallSizeRange` ⇒ 复用现成形状，不新造依赖边、`internal/agent/approval/**` 零改动；③ ⓒ（只补可见性）**治不了 99999 那一支**（票面自己写的），而 27 小时 46 分那一支才是本票标题里最贵的那格。
- **带的形状（⛔ 不是我把数钉死，是把不变式钉死）**：下界**必须严格大于** `approval.DefaultApprovalWarning`（现量＝30s，`queue.go:109`）——这一条不是审美：`queue.go:88` 的兜底 `warnBefore <= 0 || warnBefore >= timeout` 在"生产者从没传过 `WarningLead`"（`267-a1` §1 实测：全仓非测试**零枚**生产者）时恒把 warn 赋成 30s，于是 **`timeout ≤ 30s ⇒ gate.go:528 的 lead ≤ 0 ⇒ C18 那句提示不发**"。⇒ 只要下界 > warn，"静默把保护撤了"那一支在**配置层就不可达**，AC#2 自然兑现。
  - 具体取值：**下界 31s、上界 3600s**（默认 300s 落在带中央）。⚠ **为什么不是 60s**：把下界抬到 60s 会让票 256 那枚既有钉 `resident_approval_risk_256_windows_test.go:251-254`（种子 **45s**）在加载层就被拒 ⇒ 一枚本来绿的用例被打红。**要抬到 60 必须先落一枚具名解冻（A##＋撤销口令），改那两枚种子值**——本票**不做**，列为残余（见下）。
- **⛔ 落地腿不许顺手动 `l1_window_sec`**：它今天**已经有**双向钳（`gate.go:143-151`，`[2s,3s]`，`256-v1` 的 99→3s／1→2s 读数与该实现逐支对得上），而 `:263-266` 那枚钉断的就是"钳后必落在带内"；在本票里给它加 validate ＝把"越界被钳"改成"越界拒载"＝**改写那枚钉的语义**，不在本票射程。⚠ 它的**不对称本身**（一枚拒、一枚钳）是账不是 bug，见下面残余。
- **残余（本票不做，具名归口，⛔ 不许任何腿自行"顺手"补）**：① 下界 60s 那支＝要改 256 的种子钉 ⇒ 待编排者具名解冻；② `confirm_timeout_sec` 用"拒载"而 `l1_window_sec` 用"钳位"这枚**同一次落地里的行为不一致**（`267-a1` §3 逐字确认）⇒ 归口给下一枚治理票；③ `WarningLead` 生产里**恒零值、靠兜底**这一事实（`267-a1` §1）＝"配了但装配根从不传"的同族 ⇒ 归票 255 的账；④ 附录里那 6 枚无上界的数值键（`llm.timeout_ms` 全无守卫等）＝⛔ 不在本票，本票只处理 `[risk]` 那两枚键的题面射程。
- **AC 框状态**：**AC#0 已翻**（09:1x，凭据＝上述普查件＋我同发复跑的七把尺）；AC#1–AC#4 归落地腿产出读数、**仍由我翻**，⛔ 腿一枚都不许碰。Q-77（配置值 vs C18 写死 300s 谁优先）**继续待机主**，⛔ 任何腿不许据此自裁。

## 编排者收件与 §7.1 裁定（2026-10-05 09:5x，账 `A611`；本节覆盖上面那条"AC 框状态"子弹——现态＝**AC#0／AC#1／AC#2／AC#3 已翻（`9e24d178`），AC#4 按住**）

**1. 落地腿 `267-r1` 核过为真**：证据件 `.scratch/wisp/probes/267/r1/evidence.md`＝**294 行／30,943 字节／八节填实**（我第一次读到的 §1–§5"未判"是中间态，它 09:41:40 那笔 `0c2445d1` 补完了；我那条"占位符＝没交件"的尺在这枚件上现量**零枚真占位**——⚠ 但它自己的占位尺写成 `grep -c "未判|待填|填写中"`（ERE 的 `|` 在基本正则里是字面竖线）＝**那把尺是瞎的**，恒返 0；我的复跑用了 `\|` 才看见真相）。六笔 commit 名册我自己核过：**产码面只有 `internal/config/**` 五枚**（`validate.go`／`validate_test.go`／`validate_267_test.go` 新件／`schema.go` 仅注释／`unwired.go` 仅 `:121` 一行释文）。越界尺同发（09:39）：`git diff --name-only f761a017..HEAD` 里 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／`internal/agent/approval/**` **全部零字节**；C18 那两枚 300 s 现量未动（`schema.go:460` 仍 `default:"300"`、`queue.go:107` 仍 `DefaultApprovalTimeout = 300 * time.Second`）。

**2. 三格为什么敢翻、一格为什么不翻**：AC#1 的凭据不是"绿了"而是**上界与下界各自单独可打红**（M2 松上界→只有 3601/99999 红；M3 松下界→只有 30/10 红，且红的是**由 `approval.DefaultApprovalWarning` 现算的那条断言**＝同源守卫真有牙）；腿还主动登记了一处不敏感（`TestMinimumLegal…` 摘带仍 PASS，它断的是正向存活）并写明"不许替自己圆"——这条我认下来当凭据用。AC#2 的读数是**真 `approval.Queue`**（喂法逐字照 `run.go:616`／`resident_approval_windows.go:460` 那两枚装配根，生产者不传 lead ⇒ 兜底 30 s）＋种 30 必拒载 ⇒ 按 ⓐ 这形"静默撤保护"在配置层不可达；**gate 级"warn 真响一次"那半格写面外，见下面第 4 条——它由落地种子那一发免费补上**。**AC#4 不翻**：票面那句"红名集合逐名比对"这一发有实测答案了，而且是**红的**——`cmd/wisp` 此刻＝`PASS=209／FAIL=17／SKIP=1`（09:46:00→09:50:2x，246 s，PATH 带 `third_party/sherpa-onnx`），而我在 10-03 同口径量过它是 `PASS=240／FAIL=0／SKIP=0` ⇒ **本票这把 band 把下游打红了**。逐名归因在 `A611`，名册台件＝`.scratch/wisp/probes/267/gate/cmdwisp-HEAD.log`。⛔ 这不是"AC#4 判失败要退回整票"，而是**AC#4 的凭据要等种子迁移之后那一发复跑才齐**（越界路径那半已经零命中，卫生四门那半在 `internal/config` 未扩大）。

**3. ★§7.1 裁定＝出口**甲（抬种子进带内），⛔ 乙／丙／丁 逐支具名否掉（料＝只读腿 `267-a2` 普查件 `.scratch/wisp/probes/267/a2/census.md`，它**顶回了我派单里的名册**：带外种子不是 3 枚文件 8 处，是 **5 枚文件 12 处**——我漏了 `cmd/wisp/approval_seam_201_test.go:50/:138` 与 `cmd/wisp/ticket224_assembly_test.go:81/:234`；它同时量到**生产路径零枚带外种子**：`cmd/wisp/firstrun.go:82` 落盘的是 `NewDefaults()` 反射来的 `default:"300"`＝带内）。
- **甲的代价被我自己的尺核过，只有三笔墙钟**：`approval_reply_201_test.go:439`（全场唯一一枚"观察对象就是超时本身"——`:467` 逐字 `This is the timeout path, so the wait is the C18 clock`）＋29 s；`panel_pump_test.go:136` ＋29 s **且必须同把 `:112` 那枚 30 s ctx 一起抬到 >31 s**，否则观察对象从"C18 超时"变成"上下文作废"＝**这一枚是我派单里没预见、腿自己抓出来的配套件**；`run_mode101_test.go` 那族在失败路径上每次多堵 30 s（今天先响的是 `:544` 那枚 300 ms ctx，1 s 超时从来没轮到）。其余 8 枚答完卡就走＝**零增量**。先例现成：这棵树上今天已有 **15 枚用例用 ≥40 s 的种子**在做"卡开了再合上"。
- **丙（fake clock）＝今天不可达，不是我推理**：`cmd/wisp` 里 `Clock` 引用**产码 0 枚、测试也 0 枚**（09:53 我自己复跑两把 grep 都是 0），`run.go:612-619` 的 `approval.Options` 没传 `Clock` ⇒ 要让假时钟进到那两发得先动 `internal/agent/approval`（本票禁区＋C18 面）。
- **乙（测试窄缝）＝会把集成腿变成另一扇门**：注入 gate 走 `run.go:600-608`，那里逐字 `rt.ui = nil // no console surface in a process whose cards go to another host` ⇒ `panel_pump_test.go:169` 读的 `rt.ui.shown()` 与快照泵宿主面直接消失；`approval_seam_201`／`ticket224_assembly` 两族的主题恰恰是"**装配根自己**那套 reply/gate/ledger"。
- **丁（`l1_window_sec` 那枚合法小窗口）＝只对得上 3 枚**（`#7`/`#8`/`#12`，它们的 `confirm_timeout_sec` 种子对本发是**死重**——L1 走 `gate.go:291` 的 `g.window`）；替不了 L2 那两枚，且文件自己写了为什么：`panel_pump_test.go:160-162` 逐字"verdict is L2, so the item goes into the approval QUEUE (**an L1 window does not**)"，而极性还相反（`:545` 逐字 `the L1 timeout polarity changed (this is the D4 contract…)`——L1 到点**执行**、L2 到点**拒绝**）。

**4. 甲顺手买到一枚别人买不到的东西**：种子抬到 31 s 后，`approval_reply_201_test.go:439` 那一发的 `lead = 31 − 30 = 1 s > 0` **第一次为真** ⇒ 它自然成为 `267-r1` §7.2 说"没人能读"的那半格（C18 提示端到端到达）的**免费正控**。⇒ 落地腿的判据里要写死这一发：**种 31 ⇒ 那句提示真到达**，不许只改种子颜色。

**5. ★新抓到一枚本票射程外的形状，具名立成票 268**（`267-a2` §1.1★／§7.3，我自己复读了代码）：常驻腿 `residentRiskGateValues` 也走 `config.LoadFile`，但它对"加载失败"的处理是**回落编译常量 300 s 并写 provenance `"defaults (config.toml unreadable)"`**——不是 run 腿那样退码 2 响亮拒绝。⇒ band 之后，**"文件存在但越界"与"文件不存在"在常驻腿里折成同一枚 provenance**，而票 256 那一族五枚字面（45/90，全带内）里没有一枚把这两种"读不到"分得开。安全结论不受影响（300 s 仍 > 30 s ⇒ 提示仍武装），坏的是**说实话**：用户写了个数、系统按另一个数跑，且不说。⛔ 本票不动那一支（写面属常驻腿），归票 268。

**6. 落地腿排程**：种子迁移那一腿**必须在 `cmd/wisp` 上跑完整包**（它要的就是这一发颜色），⛔ 不许与 `257-r2`（同为 `cmd/wisp`）并发，⛔ 也不许在我自己那发 `cmd/wisp` 复跑期间起飞。迁移名册＝上面那 10 枚带外喂值点（5 枚文件）＋`panel_pump_test.go:112` 那枚 ctx；解冻走 `A611` 具名（只到 `_test.go` 里的种子字面量与那一枚 ctx，⛔ 不许动断言、不许动阈值、不许 `t.Skip`）。撤销口令「**267 撤**」。

## 结案（编排者，2026-10-05 10:5x，锚 `c7bb02be`；凭据全部编排者现跑，账 `A615`）

**五格全勾 ⇒ 本票改 `-done`。** AC#4 的凭据（⛔ 不是腿的自述，是我这两小时自己跑的）：

1. **越界尺**：`git diff --name-only 21bec8a1..HEAD`（＝本票全链，含 `267-r1` 的 band 与 `267-r2` 的种子迁移）里，`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／`internal/agent/approval/**` **命中 0 行**〔10:39 现量〕；该范围内非 `.scratch` 改动只有 **11 枚文件**＝`internal/config/**` 五枚（`validate.go`／`validate_test.go`／`validate_267_test.go`／`schema.go` 注释／`unwired.go` 一行）＋`cmd/wisp/**_test.go` 五枚＋本台账。C18 那枚 300 s **未动一字**（`queue.go:107` 仍 `300 * time.Second`、`schema.go:460` 仍 `default:"300"`）。
2. **卫生四门**（10:50–10:53 我自己重跑，⛔ 没引用腿的读数）：`sh scripts/d22scan.sh`＝**rc 0、clean**（live 分母 `bans #1-5 internal/=228`／`cmd/=38`／`ban #6 frontend/=85`，与旧基线不同口径、按当发引）；`sh scripts/check-path-length-budget.sh`＝**rc 0、VERDICT GREEN**；`go vet ./cmd/wisp/ ./internal/config/`＝**rc 0**；`gofumpt -l internal/config cmd/wisp`＝只剩预存的 `cmd/wisp/models.go`（CRLF，**不在本票写面**⇒ 不算扩大）。
   ⚠ **记我一枚仪器坑（第 70 条同源）**：我第一次跑 gofumpt 用了 `~/GOPATH/bin/gofumpt`，而本机 `GOPATH` 实际是 `D:\work\base\gopath` ⇒ 那次"**空输出**"**不是干净，是尺没跑到**（两个分支的 stderr 都被我 `2>/dev/null` 吞了）。换真路径才读到 `models.go` 那一行。**"空输出先怀疑仪器"这条我今天又需要一次才算长进。**
3. **红名逐名比对（判据＝只减不增）**：基线＝**16 枚**（09:46:00→09:50:2x，`cmdwisp-HEAD.log` 1,575 行，逐枚都带 band 的拒载句）；我这发＝**10:37:0x→10:49:15、508.882 s、`--- FAIL` 共 1 枚**＝`TestAC1ResidentLegInstallsItsLogListenerOnDisk`（`resident_sink_nail_127_windows_test.go:433`／`:498`"the child never reached the event loop"），**与那 16 枚名册不相交** ⇒ **16 枚全消、新增 0 枚**。
   - ★**两发各红一枚、且不是同一枚**（腿那发＝`TestAC14GoSideEvalPushReachesThePage`；我这发＝上面那枚）。两枚**都早已具名在册**：AC14 那枚在 `docs/evidence/s1/33-panel-host-c27-v2.md:183` 的"本机真建窗 11 枚"名册里、且台账 `:10727` 与 `:11899` 记过它在 CI 侧翻色；AC1 那枚＝`A607` §7 我亲自核过归口的"**票 127 钉 `:497` 不 poll 竞态（隔离 ×4 绿、带载偶发）**"。⇒ **判语＝本票没造成任何新增红；这两枚属带载偶发的既有账，⛔ 不许我为了变绿去放宽任何断言，也不许把它们划成本票的伤。**
   - ⚠ **我这发的 RUN／PASS 数拿不到**：我跑的是不带 `-v` 的整包，Go 只打失败行 ⇒ 我这条尺天生读不到 `--- PASS`（`^=== RUN`＝0 命中）。四数（RUN=323／PASS=226／FAIL=1／SKIP=0）**归 `267-r2` 那一发台件**（`cmdwisp-after.log`，腿用 `-v` 跑的），我这边成立的只有"红名一枚、与 16 枚不相交"。**引用四数时不许挂在我的时刻上。**
4. **墙钟账（⛔ 这条是我派单的错，不是腿的错）**：腿实测 `cmd/wisp` **246 s → 508.9 s（我这一发）＝+262.9 s**，我派单里预估的是 **"+58 s 属预期"**——**低估了 4.5 倍**。逐枚归因（**时刻与秒数出自腿的 `-v` 台件 `cmdwisp-after.log`；我这发没带 `-v`，只有包级 508.882 s 这一枚是我的**）：`TestTicket101SessionGrantDoesNotCrossRestart` 81.70 s（重启两发各等满）、`TestTicket101ManualSwitchSurvivesRestart` 42.24 s、`TestRunBooksWithASnapshotOfItsLiveQueue` 32.60 s、`TestUnansweredL2CardTimesOutIntoRejectNeverExecution` 32.44 s——**这四枚的等待就是被测对象本身**，票 267 的判据明写不许压种子、不许 `t.Skip`，所以这 262.9 s 是**本票的诚实代价，我认下并写进账**。
   - ⚠ 顺带一枚可省的账（⛔ 不在本票做，归后续）：`run_mode101` 那族抬的是 **40 s**，而它的最小合法档是 **31 s**（那一族的观察对象是 L2 卡、ctx 先响），三枚重启类改 31 可省约 **54 s／发**。要动就是新的一枚 `cmd/wisp` 写腿，⛔ 不许顺手夹在任何票里。
5. **`267-r2` 的 §4 那一格我认**：种 31 时 `[warning] 审批将在 30 秒后自动拒绝，请尽快确认` **真到达装配根 stdout**（链＝`gate.go:528`→`:557-568`→`run.go:611`→`:1394`），并带两发突变（M1 把 needle 改错⇒红并倒出真实 stdout；M2 把种子压到 30⇒加载层响亮拒载）；"掐发送路"那一形**写面在 `internal/agent/approval` 之外、本票禁区不许进**，腿具名归口给编排者 ⇒ **这一格按 AC#2 已翻的口径成立，残余那半格记在 `A615`，不抹进"已证"**。

**残余（具名，⛔ 任何腿不许"顺手"补）**：① `l1_window_sec` 用钳位而 `confirm_timeout_sec` 用拒载的**行为不一致** ⇒ 归下一枚治理票；② 下界抬到 60 s 那支要动票 256 的种子钉 ⇒ 待具名解冻；③ `WarningLead` 生产里恒零值、全靠兜底 ⇒ 归票 255 的账；④ **常驻腿把"越界拒载"读成"文件读不到"** ⇒ **票 268**（料已由 `268-a1` 量齐，见 `.scratch/wisp/probes/268/a1/census.md`）；⑤ `Q-77`（配置值 vs C18 的 300 s 谁优先）＝**待机主一句话**，⛔ 本票任何腿不许自裁。

---

## 凭据归口（搬运笔，非新验收）

1. 本节＝**搬运**：凭据早已在盘，只是没归口到本票名下；⛔ 本节不新增任何验收读数、不改变本票任何一枚勾的状态。
2. 出处＝只读普查腿 `evidence-close-6`（件 `.scratch/wisp/probes/evidence-close/6/backfill-worklist.md`，§4.3 名册）＋编排者裁定。
3. ⛔ 本票**不该被本节视为已验收**。

归口的凭据路径（⛔ 本节不替它们写任何判语，判语在原件里）：

- `.scratch/wisp/probes/256/v1/verdict.md`（`:1` 自称"非实现者验收表"）里 `:12` 引的那两发种子读数

存在性复量（搬运腿 `backfill-7a`，2026-10-06 现跑 `test -f`／`wc -l`／`Read`，⛔ 不引该表的任何判语、不替本票五格裁任何东西）：
该件**在盘、223 行**（与编排者给的 223 一致），`:1` 逐字含「**非实现者验收表**」（与编排者那句自述一致）。
⚠ **一处行号对不上，具名登记（⛔ 不照抄、不自己找替代品填）**：编排者给的落点是 `:12`，
本腿复量该件 `:12` ＝**分隔线 `---`**，不含任何种子读数；那两发种子读数（种 `99999` ⇒ `gate_timeout=27h46m39s`、种 `10` ⇒ `warning_lead=30s` 使 `c18_warning_scheduled=false`，另带一发 `45`）
今天实际落在**`:62-66` 那张三行表**里（正文描述在 `:60`），与本票票面 `:4`／`:12` 所引的读数是同一批。
⇒ 本节按编排者的**语义落点**（"`:12` 引的那两发种子读数"）归口，并把"盘上这一批读数的真实行号＝`:60`／`:62-66`、不是 `:12`"这一处如实报给编排者；
⛔ 本腿不改动原件任何一行、不替它重新编号、不据这一处判任何格。
