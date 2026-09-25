# 146 — `Queue.LiveApprovals()` 的注释说"返回副本"，实现是**按值拷贝一个含引用字段的结构**：那张 map 与那三枚切片都和队列共用底层

- Status: in-progress（**已派**，issues/README 规则 1）
- Claimed by: `agent=ticket146`（实现程，写码）
- Last update: `2026-09-25T11:00Z`（＝`2026-09-25 19:00 +08`）
- Parallel slots: 本票只此一枚写手；`internal/panel/l2_grant_boundary_test.go` 与 `frontend/**`/`design/**` 零接触
- 来源：票 35 快照泵收尾程 `docs/evidence/s1/35-panel-snapshot-pump-fix-r1.md` 留给下一程的空格 ②；
  它**没改行为、也没在两个方向任一处下断言**，所以这一格今天仍是空的（原话："这格今天仍是空的"）。
  ⚠ **本票面把它写大了**：它只点了两枚切片，编排者 18:2x 逐字段读被拷贝的结构后是**一枚 map ＋ 三枚切片**（见下）。
- 关联：票 35（`LiveApprovals` 就是它那一批里新增的出口）、票 33（泵将来要往页面送这份东西）、
  票 145 `AC#3`（"当前屏"那一维的归属裁定在这里会被引用到）

## 现量的形状（三条命令可复算）

- `internal/agent/approval/pending_read.go:19-20` 的注释逐字写着
  **"It returns copies of verdicts the assessor already reached"**（"copies" 那一个字在本行行首）。
- 被返回的结构是 `LiveApproval{CorrelationID string; Decision tools.Decision; Position int}`
  （`pending_read.go:30-34`），实现里那一行是 **`Decision: it.Dec`——按值拷贝整个 `tools.Decision`**。
  ⇒ 但 `tools.Decision`（`internal/tools/gate.go:15-24` 起）**自己就带引用类型字段**：
  | 字段 | 类型 | 按值拷贝之后 |
  |---|---|---|
  | **`Params`** | **`map[string]any`**（`gate.go:18`） | **同一个 map**——写它会直接改掉队列里那条记录，**这比切片更硬** |
  | `RulesHit` | `[]risk.RuleID`（`gate.go:21`） | 共用底层数组 |
  | `Args` | `json.RawMessage`（`gate.go:19`，本质是 `[]byte`） | 共用底层数组 |
  | `Paths` | `[]string`（`gate.go:23`） | 共用底层数组 |
  | `Capabilities` | `[]Capability`（`gate.go:24`） | 共用底层数组（元素本身还有没有引用字段，本票没量） |
  ⇒ **"copies" 许诺的是深拷贝，实现给的是"浅一层"**。
- **今天读它的是谁（全量现量，含测试）**：`LiveApprovals(` 共 **11 处**＝1 处定义
  （`pending_read.go:42`）＋ **1 处生产调用**（`cmd/wisp/panel_pump.go:61`）＋ 9 处本包测试（`pending_read_test.go`）。
  而 `grep -rn "RulesHit =\|RulesHit\["`（剔 `_test.go`）在全仓只命中 3 处，都在别的地方构造/搬运 verdict、
  **没有一处是"拿到 `LiveApprovals()` 返回值之后就地写它"**。
  ⇒ **所以今天不是缺陷，是一个还没被触发的形状**——这句要写死在结论里，别读成"有个 bug"。
  ⚠ 上面那三条是我 18:2x 在本机工作树量的；**下一程一律现量，别抄我的数**。

## 为什么现在不许顺手"改成深拷贝"就完事

这一格有两种合法收法，**方向相反、代价不同**，本票要你**先判再改**：

- **ⓐ 改实现**（把那枚 map 与那三枚切片真拷一层）——买的是"注释说的字成立"；代价＝每次泵动一下都多一张 map＋三份底层数组，
  而**这张表是要频繁重算的**（`wisp run` 每有面板状态变动就现算一份），所以它**撞 `PLAN.md` 的 D32 资源预算那一族**。
  ⇒ 选 ⓐ 必须**同时量一次**：多拷这些东西，`wisp slo` 的 WorkPeak 那一档数字动没动。没量过就不许写"代价可忽略"。
- **ⓑ 改注释**（把 "copies" 换成"逐条复制结构体本身；**`Params`／`RulesHit`／`Args`／`Paths`／`Capabilities` 与队列共享底层，调用方不许就地写**"）
  ——买的是"注释不再许诺它没做的事"；代价＝把纪律放在调用方身上，**而调用方里有 `internal/panel` 与将来的页面**。

⚠ **本票的默认方向是 ⓑ，但 ⓑ 单独存在不够**：一句"调用方不许就地写"是**没有仪器的纪律**，
而这一族东西一旦被谁忘一次，症状是**队列里那条已决断记录被悄悄改掉**——那正好落在
"面板侧不许变成授权来源"那条硬禁线的旁边。⇒ **无论选 ⓐ 还是 ⓑ，都必须有一枚会响的检**（AC#2）。

## AC

- [x] **AC#1 先判 ⓐ／ⓑ，判不了就停手上报**（D22 闸门③，不许自行假设）。判据是**可核的一句**：
      今天有没有**任何**一枚调用方（含测试）在拿到返回值之后对那枚 map 或那三枚切片做 改元素／append／就地写／取地址后传出去？
      现量命令自己写进证据件。⇒ 全为"无"时 ⓑ 是诚实且零代价的一支；只要有**一枚**为"有"，ⓑ 就不许单独收，必须 ⓐ。
- [x] **AC#2 一枚会响的检（本票的硬核心）**：造一发**只有"就地改掉返回值里的 `Params` 或那三枚切片"才能触发**的用例，
      断言队列里那条原始记录**没有**跟着变。判"承重"按本仓操作定义答一句：
      **把 AC#1 选定的修法摘掉，这一发是不是从此打不红？** 答"仍不红"＝这枚检是装饰，本票作废并登记为什么。
      ⚠ **不许**用"永远绿"的那种断言凑数（本仓那族缺陷有名字：**恒真判据是一类新假绿**）。
- [x] **AC#3 若走 ⓐ，资源那一问不许空着**：给出"多拷那枚 map 与那三枚切片"在 WorkPeak 档的前后读数
      （`wisp slo` 或 `scripts/slo-check.ps1` 的同形发法），**阈值与 golden 一字节不许动**，
      数字没动要写"未观察到差异"而**不是**写"无代价"。
      > **⚠ 这枚框按"实质凭据"勾，不按字面勾（09-25 21:2x 编排者注，`A261` 有账）。**
      > **字面产不出来的原因**：WorkPeak 档的 subject 是 skeleton（`goroutines_max=1`），**那 5 秒里 `LiveApprovals` 一次都没被调用**
      > （唯一的生产接线在 `cmd/wisp/run.go` 的 `wisp run` 腿上，常驻腿不 import `internal/panel`）
      > ——两程独立现量一致（实现件 `146-liveapprovals-r1.md` §AC#3、验收件 `-accept-r1.md` AC#3 格）。
      > **所以"原句不删、框照勾"买的到底是什么**：两发**代替凭据**——①双臂各 3 次同形跑、均值差**小于组内极差** ⇒ 按票面写"未观察到差异"，
      > 没写成"无代价"；②一发确定性 `-bench`（仓外）：满深度一次调用 ＋4736 B／＋72 allocs。**阈值与 golden 零字节。**
      > **从未被字面复算过的那一支**：真在 WorkPeak 档里驱动 `LiveApprovals` 再取前后读数——**那一支归票 148 `AC#2`**，
      > 它的触发条件是具名的：**谁让 WorkPeak 的 subject 真跑起工具调用，谁就必须回来重量这一格**。
- [x] **AC#4 契约轴**：`internal/risk/**`、`rules_gateway.go`、`thresholds.go`、golden、`allowlist.txt`、
      `docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`、`tools/d22scan/**`、
      `internal/panel/l2_grant_boundary_test.go` **一律零字节**。
      ⚠ 这枚函数**不在 `PanelAPI`／C17 名册上**、不能回答从页面来的请求（`pending_read.go:14-18` 自己就这么写）⇒
      **本票不许把它变成一条入站路由**；那一变要重新走契约批准（`Q-51` 那一族）。
- [x] **AC#5 门禁（一律逐包单跑）**：`go test ./internal/agent/approval/` 改前改后各一次
      （四数之外**点名册差集**；本目录前八枚测试文件是 `package approval_test`，
      而 `7040493` 那枚是 `package approval` **白盒**——两种都有分母，别把它们当成同一棵尺）。
      ⚠ **两包合跑会凭空造红**（`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 报 `0xc000013a`、
      单跑就过、那文件与本票无关），这条是本仓 09-25 现量的仪器形状。`gofmt -l` 空、`go vet` 带你动过的包、
      `sh scripts/d22scan.sh` rc=0（独立 module；`set -eu`，**第一步是正控，正控红则真扫描根本不跑**）。

## 两条已知读数，别当结论用（原样交给下一程）

1. 收尾程在**锚点归档副本**里读到两枚 `TestSecret*` 红，**工作树单跑是 PASS**——它**没归因**，
   并明确写"别当'改前就红'引用"。⇒ 你若在副本里再见到它们，先做 `git diff <锚>..HEAD -- <那个包>` 对照再判。
2. 行号会漂：那份简报与验收件写 `cmd/wisp/run.go:741`，现树同一句 `changed = true` 在 **740**。
   ⇒ 本票引用的任何行号都**现量**，别从上面两枚文件抄。

## 规矩

显式 pathspec、禁 `add -A`/`.`/`--amend`/`reset`/`stash`/`checkout .`/`clean`、禁 push、
**临时件只建不删**（变异落仓外副本）、台账与勾归编排者、简报里任何一句前提复算不符 ⇒ 写进证据件并报回，
**不要按我认为的样子改**。

## Progress log（append-only；每个 commit 一行 `- [ISO-UTC] agent=... did=... next=...`）

- [2026-09-25T11:00Z] agent=ticket146 did=**认领 + 现量（本 commit 零代码改动）**。三组读数，命令原文与落盘副本见
  `docs/evidence/s1/146-liveapprovals-r1.md`：
  ① `LiveApprovals(` 在 `*.go` 里 **10 处**（1 定义 `pending_read.go:42` + 1 生产 `cmd/wisp/panel_pump.go:61` + 8 本包测试）
  ⇒ **票面说 11，我量到 10**，差 1 枚在测试侧（票面记 9 处测试，现树 8 处）；票面 :30 自己写了"下一程一律现量"，按现量走。
  ② **被拷贝的结构的引用字段是 8 枚、不是票面表里的 5 枚**：`tools.Decision`（`internal/tools/gate.go:15-53`）除
  `Params`/`Args`/`RulesHit`/`Paths`/`Capabilities` 之外还有 `Blacklist BlacklistNote`（`gate.go:43`），而
  `BlacklistNote`（`internal/tools/mode.go:66-81`）自己带 `Absolute`/`Unlockable`/`AlreadyUnlocked` **三枚 `[]string`**
  ⇒ 票面表漏计三枚（它只在 `Capabilities` 那行写了"元素本票没量"，没往 `Blacklist` 看）。形状仍然成立（确实共用底层），
  只是比票面更大，所以**修法是同一族、只是范围宽一格**，不必停手；这一条按"复算不符"登记并上报。
  ③ **AC#1 那句可核的判据：现量＝0 枚调用方（含测试）就地写过返回值里的 map/切片**（普查命令见证据件 §1.3），
  `dec.Params` 在生产里**确有真值**（`internal/tools/bridge.go:267` 每次调用现解一份 JSON），所以那不是死字段、风险是活的。
  next=**裁 AC#1**（ⓐ／ⓑ）。当前倾向 **ⓐ**，理由与"这跟票面默认方向 ⓑ 相反"的分歧一起写进证据件 §2 交编排者复核；
  若我判不下去就停手，不自己填。AC 的框**一枚都不勾**——本票 :81 明写"台账与勾归编排者"。

- [2026-09-25T11:11Z] agent=ticket146 did=**裁 AC#1 ＋ 装 AC#2 的检（`22f7b1a`，两枚 .go，均在 `internal/agent/approval/`）**。
  **AC#1 判 ⓐ（真拷一层），与票面默认 ⓑ 相反**——依据三条，全部在票内文字里，证据件 §2：
  ①普查（§1.3）确实答"0 枚就地写"，ⓑ 今天不违法；②但票面 AC#2 把断言极性写死成"队列里那条原始记录**没有**跟着变"，
  这句在浅拷贝之下永假 ⇒ 选 ⓑ 要么交不出一枚绿的检，要么把危险形状钉成规格（票面 :54 那族"恒真判据＝新假绿"）；
  ③"把修法摘掉这一发是否打不红"对**一句注释**恒答"仍不红"，正中票面 :53 的作废条款，而 ⓑ 唯一可能的另一类仪器是静态扫、
  AC#4 又把 `tools/d22scan/**` 划成零字节。⇒ **两枚 AC 只有 ⓐ 同时收得下**。若编排者改判 ⓑ：那要动 AC#2 的极性＝改票＝人工批准。
  **代码**：`cloneDecision`＋`cloneBacking`＋`cloneParamsMap`，把 §1.2 现量出的 **8 枚引用槽位（1 map＋7 slice，
  含票面漏计的 `Blacklist` 里三枚）**各自分配底层；`queue.go` 一字节未动（判 ⓐ 不需要它）。
  **检**：两发。probe 1 反射点名"返回值里哪几枚槽位仍与队列共用底层"（槽位名册按集合断言，walker 失效即红）；
  probe 2 = AC#2 的字面：只通过返回值就地写全部 8 枚，再从三条路读回队列记录（白盒 `qitem.Dec`／投影 `view()`+`head()`／
  第二次 `LiveApprovals()`）要求原样。
  **承重那句答"会红"，且是先写检、后装修法量出来的**：改前 `rc=1`、两行 `--- FAIL:`、14 条 `AC#2 RED` 明细，
  probe 1 逐字点名全 8 枚（`D:\tmp\wisp-146-agent-a\ac2-before-fix.txt`）；装上 ⓐ 后同一对 PASS。
  门禁（逐包单跑）：`./internal/agent/approval/` rc=0，PASS 32→34／FAIL 0→0／SKIP 1→1（那枚 SKIP 改前就在），
  点名册差集＝只多我这两枚、零删除；`gofumpt -l` 空、`go vet` 通。
  **没买到的东西也写在代码注释里**：拷贝只一层深，`Params` 的 value 里那层 `[]any`／嵌套 map 仍共用、两发检都不探它。
  next=**AC#3 的 WorkPeak 前后读数**（先确认 `wisp-selfhosted-01` `busy:false` 才敢取数）。

- [2026-09-25T11:35Z] agent=ticket146 did=**量完 AC#3 ＋ 收 AC#4/AC#5 ＋ 证据件写满（`0d12661` + 本 commit）**。
  **AC#3 走成了**（ⓐ 才需要它，我判了 ⓐ）：取数前先确认仪器空着——`gh api …/actions/runners` → `busy:false`
  （取数前后各一次）、`gh run list` 全 completed。**不在工作树上建 `cmd/wisp`**，改从两棵 `git archive` 纯净快照建
  （before=`630c218`、after=`22f7b1a`，唯一变量是本票那枚 commit），因为当时 `cmd/wisp/slo_windows.go` 被别人改着
  且不能编译（3 枚未用 import）——**我一字节未动、也没去救**。
  `wisp slo -state WorkPeak -seconds 5` 双臂各 3 次交替：before `mem_median` 均值 4134229（组内极差 139264）、
  after 4086443（极差 114688）⇒ 均值差 **−47787 B 小于组内噪声** ⇒ 按票面 :56 的字写"**未观察到差异**"，**不写"无代价"**；
  阈值/golden 一字节未动。**同时把更要紧的一句写死**：WorkPeak 的 subject 是"skeleton＋sit still"（`goroutines_max=1`），
  而 `LiveApprovals` 唯一生产接线在 `cmd/wisp/run.go:422` 的 `wisp run` 路上 ⇒ **被采的那 5 秒里这枚函数一次都没被调用**，
  两臂同形是"仪器与改动不同路"的证据、不是代价的证据。机制账另给确定性一份（仓外副本 `-bench`，**未进仓**）：
  满深度一次调用多 **4736 B／多 72 allocs**（3456→8192 B/op、1→73 allocs/op）。
  **AC#4**：本程 3 枚 commit 一共只碰过 4 枚文件（票面／证据件／`pending_read.go`／新测试），过一遍 AC#4 名单
  ⇒ **交集空**；`queue.go` 亦零字节；收尾时我那三枚路径 `git status --short` 为空；删除列逐行点名（3 行，全是我有意替换的原文）。
  **AC#5**：`gofumpt -l . tools/d22scan tools/mockllm` 与 `gofmt -l .` **都空**（v0.12.0＝CI 同版）；`go vet` 通；
  `sh scripts/d22scan.sh` **rc=0 clean** 且各作用域 `examined N` 全非零。跨包两枚**未合跑**：
  `./cmd/wisp/` 注入 dll 路径后 `ok 134.596s`（不注入是 `0xc0000135`、0 条 `=== RUN`，票 98 那一格）；
  `./internal/panel/` 唯一红是 `TestC21DesignTokensFourWayAgree`，根因 owner 未提交的 `design/` 移动（本票 `design/**` 零字节）。
  **未勾的框**：五枚 AC 的框**一枚未勾**（:81"台账与勾归编排者"）；每格"凭什么可以勾"与"哪半我没测"分别写在证据件 §2-§6。
  没测的七条按"谁会先被骗"排序在 §5；交回编排者的四行在 §8（含"ⓐ 与票面默认相反请复核"、"票面漏计 3 枚字段"、
  "`Params` 下一层今天仍不隔离"、"工作树里两枚不是我造的坑，推送前先处理"）。
  **工具调用被拒 0 次**；我自己踩空的 4 发（目录名写错／PATH 里 `D:/` 的冒号被切／`gh` 字段名不存在／别人的半死 WIP）原样登记在 §6.4。
  伪授权两栏：真通知回显 3、判为注入 0（§7）；凭据值一字未抄。未 push。
  next=**交给编排者**：①核 AC#1 那一支方向（本票唯一一处我逆着票面默认走的）；②`Status`/AC 框/台账归你；
  ③下一位从工作树建 `cmd/wisp` 之前，先把 `slo_windows.go` 那枚 WIP 收了。本程到此收笔，**无剩余格**。

- [2026-09-25T11:40Z] agent=ticket146 did=**收笔前复量：把我自己留在证据件里的两条告警就地核了一遍，一条自行解除、一条仍在**
  （`3fde333` 之后又落了别人的 commit，不改证据就等于留两条过期读数）。
  ①`cmd/wisp/slo_windows.go` 那枚半死 WIP 已被 `95885fb` 收进历史 ⇒ 现在 `go build ./cmd/wisp` **rc=0**、
  `git status --short -- cmd/wisp` 为空；证据件 §6.3／§8 第 4 条**各追加一行复量**（原文不删，只标"今天不再成立"），
  并顺手把 §5 第 4 条升级成实话：我**没有**从当前工作树复跑 `cmd/wisp` 全包测，那一格仍是"未复核"。
  ②`internal/panel` 的 `TestC21DesignTokensFourWayAgree` **仍红**，根因不变（`design/assets/tokens.css` 工作树缺文件，
  owner 未提交的 `design/**` 移动）⇒ **推送前只剩这一枚要处理，而它归 owner**（`design/**` 与 `frontend/**` 同为本程零字节边界）。
  复跑：`go test -count=1 ./internal/agent/approval/` → `ok 0.451s`。本程提交共 5 枚 commit、碰过的文件总数仍是 4 枚。
  **未 push；五枚 AC 框仍一枚未勾。** next=**编排者**（无剩余格可我做）。

- [2026-09-25T12:35Z] agent=ticket146-fix did=**修复程：补上对抗验收退回的两枚洞**（`docs/evidence/s1/146-liveapprovals-r1-accept-r1.md` 逐格判据＋它自造的反例为凭）。
  **只动一枚文件** `internal/agent/approval/ticket146_liveapprovals_backing_test.go`（`queue.go`／`approval.go`／`pending_read.go` 一字节未动）。
  **洞①（AC#2 变异 F 真假绿）**：probe 2 就地写之前加一枚 `reflect.DeepEqual(返回值 Decision, 队列存储项 .Dec)`（新尺 `:285`）——把 "copies" 的**值保真**那一半钉住。
  **洞②（进攻④ "第 9 枚会静静进来"）**：选**甲（穷尽普查、非降级文字）**——`refSlots`/`slotPaths`（按值走、只认 Map/Slice/Struct、nil 槽与 Ptr/Interface/Chan/Func/私有字段静默跳过）
  换成 `declaredRefSlots`/`typeCensus`（按 `reflect.TypeOf` 数声明、逐枚 `NumField`、命名 Map/Slice/Ptr/Chan/Func/Interface、下降 Struct/Array、**未知 kind 一律 `t.Fatalf` fail-closed**）＋ 一枚 fixture 完整性守卫（`:245`）。
  **判"穷尽普查是测试侧那把尺、还是生产少拷一枚"＝前者**：`cloneDecision` 现有 8 枚今天一枚没少（验收件 §1 独立普查同判），故 `pending_read.go` 不动。
  **先证洞真在**：纯净副本复现变异 F（两发全绿、逐包 rc=0）与变异 M-a（新 map 槽静静进来、两发全绿）。**再证修完会红**：改后重跑——
  F→probe2 值保真红、M-a→probe1 名册 9 vs 8 红（这就是派单单要的"新增一枚不拷的引用槽位会红"凭据）、
  E→不再 panic 吞读数（`=== RUN` 恒 54）、**合体 M-a+F→32/2/1**（旧尺下是 34/0/1 与"改后健康态"一模一样 ⇒ **"四数一模一样"不再成立**）。
  名册差集改前/改后两向 comm 皆空（未增删用例函数）。门禁逐包单跑于隔离纯净副本（避别人在改的 cmd/wisp/frontend）：
  `gofumpt -l . tools/d22scan tools/mockllm` 空、`gofmt -l internal/agent/approval/` 空、`go vet` rc=0、`sh scripts/d22scan.sh` rc=0 clean 且 examined N 全非零。
  **两处文字更正**：实现件 §2 之后追加 `>` 更正（理由②"ⓑ 摘掉注释任何用例都不会红"、理由③"ⓑ 拿不到任何仪器"＝**假**，验收件 §6.3 造出删 239 字节纪律句就红的 doc-pin、且两枚 ⓑ 支仪器都在本包内跑通、`tools/d22scan/**` 一字节未碰；
  ⓐ 方向仍成立、由理由①独自撑住；正确取舍句＝"ⓐ 是票面两支里唯一同时收得下 AC#2 字面的那一支"）；台账 `A255`（编排者自写的"第三条理由也成立"）**本程未碰、未在任何新文字里重复那句假理由**。
  **AC 框一枚未勾**；提交时 dev 已被另一程推进 3 枚（`1e94672`/`1b3fccc`/`1f76c06`，全为票 144/台账 docs、与本包交集空）。
  本程共 3 枚 commit（`b694378` 代码 / `73e0d7be` 两枚 docs / 本枚进度），碰过的文件＝测试件 + 两枚证据件 + 本票面 = 4 枚。权限系统拒绝 0 次（一次自伤 pathspec 手误已如实登记，未换路子绕过）。
  next=**编排者**：①AC#1/AC#2/进攻④ 的勾按"第二轮验收"定；**若编排者认为 AC#2 此刻可勾**，理由已备——验收件给的最小闭合集合（`reflect.DeepEqual` 值保真，一次收 E/F/M-a）本程已落地并复跑证红。
  ②"返回值继续往外送 `Params` 是否收窄"（ⓓ）＝下一票。③`Params` 一层深、`replay` 同族出口＝下一票（验收件 §12）。
- [2026-09-25 21:2x] **编排者落笔：五枚框全勾，勾的依据是两份非实现者表，不是任何一枚自勾。**
  `AC#1`／`AC#2` ← 第二轮验收表 `docs/evidence/s1/146-liveapprovals-r1-accept-r2.md`（两格均判**成立**，未动用"退回"档；
  五发进攻里 ②③④ 三发都没打穿）；`AC#3`／`AC#4`／`AC#5` ← 第一轮 `…-accept-r1.md`（三格成立）。
  ⚠ **`AC#3` 是按"实质凭据"勾的，票面字面那支从未复算过**——原因、两发代替凭据、以及"谁该回来重量这一格"
  已用 `>` 追加在原句之后（**原句一字不抹**），触发改成具名的跨票依赖 ⇒ 归**票 148 `AC#2`**。
  ⚠ **第二轮顺手改写了三句归属**（都不是缺陷、是措辞射程）：
  ①"一条 DeepEqual 一次收掉 E／F／M-a" **过宽** ⇒ 真实形状是**三味分管三形**：名册未长＝类型普查（`:223`）、
  名册已长而少拷＝probe 1 的 backing 比对（`:255`，此处 **DeepEqual 必然绿**）、全新空 map＝**只有 DeepEqual**（`:287`）；
  ②"第 9 枚不会静静进来" ⇒ 精确射程是**会被点名、map/slice 可测到、其余四支是"拒绝测"不是"测到"**
  （撤掉那五支各自仍响是因为红挪到 `default`；**只有撤 `default` 那一发出现全绿**，而那一发本程用仓外探针当场钉出"就地写返回的 `*[]string` 真改掉了已决断记录"）；
  ③"`=== RUN` 恒 54" **写宽了**（`v-nilpaths` 那一形仍 `RUN=5`）。
  ⚠ 另两枚**归第一轮件自己**要改的：那两句"旧尺／新尺"世代没标清（修复件 §1.1）、以及 `v-ma` 那种红是**逼同步评审**而不是测出少拷。
  ⇒ **结论**：本票可结案，两枚新洞（ⓓ 收窄、`replay` 同族出口）**不塞进本票**、另开票 148；
  那两枚的**前置都是人工批准**（一个要改写 `AC#2` 措辞、一个要放开 `queue.go` 禁改面），⇒ **148 立为 blocked、不派**。




