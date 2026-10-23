# 363 — 生产日志每一发 L0 都在说假话：`MODE-REDLINE … still asks` 而实际放行

- Status: **DONE（2026-10-23 10:4x 编排者收口＝五格全闭＋两格判语到位、文件名加 `-done`）** ~~**ready-for-agent**~~
- 立案：2026-10-22 23:2x +08 编排者立
- ★**来路（非实现者签名）**：`179-r1`（写腿，件 `docs/evidence/s1/179-r1-impl.md`、4 笔 `32c8eeac`／`18e3f670`／`0aa30ede`／`8956ece7`）§7① 具名上报
- ★**编排者已现跑核实四条链**（见下）

## 这是什么（人话）

**生产日志在每一发 L0（最低风险级）上都打印"仍然会问"，而实际上 L0 直接放行、一张确认卡都不弹。**

⇒ ★**这是向操作员说假话**：他读到"还会问你"，于是**以为系统会拦一下**，实际**已经执行完了**。

## ★四条链（编排者现跑核实，逐条带 `file:line`）

1. **`internal/tools/bridge.go:311`** 逐字 `sil := mode.Screen(verdict)` ⇒ **L0 的判定被交给 `mode.Screen`**。
2. **`internal/risk/mode.go:163-165`** 逐字：
   ```
   if d.Level != L1 && d.Level != L2 {
       return Silenced{Level: d.Level, Kept: "无法识别的风险级，fail-closed 保留询问", Mode: m...}
   }
   ```
   ⇒ ★**L0 ≠ L1 且 L0 ≠ L2 ⇒ L0 落进这一支**，被当成"**无法识别的风险级**"。
3. **`internal/tools/bridge.go:324`** 逐字：`case sil.Kept != "": b.log("tools: MODE-REDLINE mode=%s still asks tool=%s rules=%v because=%s", …)`
   ⇒ ★**只要 `sil.Kept` 非空就打印 "still asks"**、**⛔ 区分那是不是真的会问**。
4. **`internal/tools/bridge.go:352`** 逐字注释 `// --- routing: L0 pass, L1 window, L2 approval, Deny never reaches a gate ---` ＋ `b.route(ctx, &dec, sil, grantID)`
   ⇒ ★**`route` 对 L0 直接放行** ⇒ **0 张卡**。

⇒ ★★**链条闭合**：L0 → 被 `Screen` 误判为"无法识别的风险级" → 日志说 "still asks" → 实际 `route` 放行、**0 卡**。

## 结案判据

- [x] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：把这条假话量成可复跑读数——
> **2026-10-22 23:5x 编排者翻框＝成立（凭只读普查腿 `363-a1`，件 `docs/evidence/s1/363-a1-census.md` **116 行**、提交 `7866ad4e`＋`7825d31e`）**
> **① 四条链逐条复认（它现跑、我抽核）**：`bridge.go:311` **不漂**／`mode.go:163-165` **逐字同号**／`bridge.go:323` 条件＋`324-325` 日志 call **不漂**／
>   ★**唯一偏差**：票面把 `// --- routing: L0 pass…` 那条注释并进了 352 ⇒ ★**实测注释在 351、`b.route` 在 352**（我现跑复认）。
>   ★**`route` 的落点 ＝ `bridge.go:432-435`**：`switch sil.Level { case risk.L0: dec.DecisionColumn = agent.DecisionAllow; return true, "" }` ⇒ ★**L0 根本不碰 gate**（我现跑复认）。
> **② ★★它造出了**真读数**（本格核心，⛔ 只读码）**：真 `fs.read`（真 D34 声明＋真 R1–R9 判定器）、**三个档各一发**，逐字同一句：
>   `tools: MODE-REDLINE mode=<m> still asks tool=fs.read rules=[] because=无法识别的风险级，fail-closed 保留询问`
>   ★**同一发**：`CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0` ＋ `kind=success decision=allow outcome=success` ＋ `EXEC text="hello"`（**已执行完并回结果**）
>   ⇒ ★★**"会问"与"0 卡"在同一发并存 ＝ 假话的直接证据**。
>   ★**且与档无关**（L0 支排在模式判断之前）⇒ ★**票面定性「每一发 L0」成立**。
> **③ ★正控（硬闸）成立**：`L1/ask_every_step` ⇒ **1 张窗口卡且根本没有那行**；`L1/ask_high_risk` ⇒ 0 卡走 `MODE-SILENCE`；
>   `L2/ask_high_risk` 与 `L2/auto_approve` ⇒ **各 1 张审批卡且 "still asks" 说的是真的** ⇒ ★**⛔ "所有档都假"**。
> **④ ★它判「`L0` 是认得的级」——我现跑复认、成立**
>   - `assessor.go:41` 逐字 `L0 < L1 < L2 < Deny: fusion takes the max, and Deny is terminal.`／`:45` 逐字 `// L0: read-only, no side effects — execute directly.`／`:46` `L0 Level = iota`（**枚举零值**）。
>   - ★★**最硬的一条**：**`Screen` 自己在 `mode.go:175`／`188` 返回 `Silenced{Level: L0, Silenced: true, Mode: m}`**（我现跑复认两行逐字）
>     ⇒ ★★**同函数把 `L0` 当**输出级**认得、却在**入口**把 `L0` 判成"无法识别的风险级"** ⇒ **出入口不一致、范畴错误**。
>   - ★`mode.go:155-156` 自述的 "unknown level" 被实现成 "not L1/L2"；★**`Kept` 契约（`mode.go:143-146`：只在**红线拒绝静默**时非空）在 `L0` 上是范畴错误**。
> **⑤ ★★两条会改变裁形的额外实测（它主动交的、我采信）**
>   - ★**A：`still asks` 在 `Deny` 上**同样假**（反方向）**：`fs.read` 打 `%APPDATA%\wisp\config.toml` ⇒ `level=Deny`、`CARDS TOTAL=0`、
>     `text="该目标被绝对禁止访问（R3 A 档…）"`，而日志**仍打 `still asks`**（`~/.ssh/id_rsa` 同形）
>     ⇒ ★★**乙（`Screen` 认 `L0`）修不到这一条** ⇒ ★**"日志不再说假话"这个目标需要**甲**独立承担**。
>   - ★**B：`Decision.ModeKept` 也在每一发 `L0` 上带假话**（`bridge.go:314` 写入、`gate.go:48-50` 自述"names the red line"），
>     ★**但实测读取方只有测试**（`ModeKept` 全仓 12 处、`cmd/**` 与 `internal/panel/**` **零处**）⇒ ★**今天唯一操作员可见面就是那行日志**；**行为面无回归**（L0 仍放行）。
> **⑥ ★尺与射程（它把"落笔前五问"逐条答了）**：仪器＝`go test -overlay`、探针在**仓外**、经**虚拟路径**参与编译 ⇒ ★**仓内 `.go` 字节 0 枚**（`git status --porcelain -- internal cmd` 首尾均空）。
>   ★**组合＝生产组合同形**（`cmd/wisp/run.go:764-788`）、唯一差异＝`Logf` 指捕获器。
>   ★**卡数怎么数**：`Gate` 接口＝两张卡的面（`gate.go:110-117`：`PendingWindow`＝L1 窗口、`PendingApproval`＝C18 队列）⇒ **数这两个回调被调用几次**、**⛔ 数 UI/tick**。
>   ★⚠**它如实标明一段**未实测**：日志到操作员＝`run.go:979-983` → 滚动 JSONL（`logsink.go:129`）——★**这是读码、它没跑整条 `wisp run`**（要 golden SSE）⇒ **件里已标明**。
> **⑦ 门禁**：`d22scan 0`（clean）／长度门 `0`（`tracked=17601 over-budget=61 covered=61 not-in-roster=0`、GREEN）。★环境：`GOFLAGS=` 空、**未设 `GOTMPDIR`**、⛔ 未用 `CGO_ENABLED=0`。
> **⑧ 纪律**：`.go`＝**0**（★同尺正控 `ec520750` ⇒ **2 枚** `.go`、**尺有牙**）；⛔ 动任何 `.go`／票面／台账／`.github/**`；⛔ 翻框、⛔ push。
> **2026-10-22 23:5x 编排者翻框＝成立（凭只读普查腿 `363-a1`，件 `docs/evidence/s1/363-a1-census.md` **116 行**、提交 `7866ad4e`＋`7825d31e`）**
> **① 四条链逐条复认（它现跑、我抽核）**：`bridge.go:311` **不漂**／`mode.go:163-165` **逐字同号**／`bridge.go:323` 条件＋`324-325` 日志 call **不漂**／
>   ★**唯一偏差**：票面把 `// --- routing: L0 pass…` 那条注释并进了 352 ⇒ ★**实测注释在 351、`b.route` 在 352**（我现跑复认）。
>   ★**`route` 的落点 ＝ `bridge.go:432-435`**：`switch sil.Level { case risk.L0: dec.DecisionColumn = agent.DecisionAllow; return true, "" }` ⇒ ★**L0 根本不碰 gate**（我现跑复认）。
> **② ★★它造出了**真读数**（本格核心，⛔ 只读码）**：真 `fs.read`（真 D34 声明＋真 R1–R9 判定器）、**三个档各一发**，逐字同一句：
>   `tools: MODE-REDLINE mode=<m> still asks tool=fs.read rules=[] because=无法识别的风险级，fail-closed 保留询问`
>   ★**同一发**：`CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0` ＋ `kind=success decision=allow outcome=success` ＋ `EXEC text="hello"`（**已执行完并回结果**）
>   ⇒ ★★**"会问"与"0 卡"在同一发并存 ＝ 假话的直接证据**。
>   ★**且与档无关**（L0 支排在模式判断之前）⇒ ★**票面定性「每一发 L0」成立**。
> **③ ★正控（硬闸）成立**：`L1/ask_every_step` ⇒ **1 张窗口卡且根本没有那行**；`L1/ask_high_risk` ⇒ 0 卡走 `MODE-SILENCE`；
>   `L2/ask_high_risk` 与 `L2/auto_approve` ⇒ **各 1 张审批卡且 "still asks" 说的是真的** ⇒ ★**⛔ "所有档都假"**。
> **④ ★它判「`L0` 是认得的级」——我现跑复认、成立**
>   - `assessor.go:41` 逐字 `L0 < L1 < L2 < Deny: fusion takes the max, and Deny is terminal.`／`:45` 逐字 `// L0: read-only, no side effects — execute directly.`／`:46` `L0 Level = iota`（**枚举零值**）。
>   - ★★**最硬的一条**：**`Screen` 自己在 `mode.go:175`／`188` 返回 `Silenced{Level: L0, Silenced: true, Mode: m}`**（我现跑复认两行逐字）
>     ⇒ ★★**同函数把 `L0` 当**输出级**认得、却在**入口**把 `L0` 判成"无法识别的风险级"** ⇒ **出入口不一致、范畴错误**。
>   - ★`mode.go:155-156` 自述的 "unknown level" 被实现成 "not L1/L2"；★**`Kept` 契约（`mode.go:143-146`：只在**红线拒绝静默**时非空）在 `L0` 上是范畴错误**。
> **⑤ ★★两条会改变裁形的额外实测（它主动交的、我采信）**
>   - ★**A：`still asks` 在 `Deny` 上**同样假**（反方向）**：`fs.read` 打 `%APPDATA%\wisp\config.toml` ⇒ `level=Deny`、`CARDS TOTAL=0`、
>     `text="该目标被绝对禁止访问（R3 A 档…）"`，而日志**仍打 `still asks`**（`~/.ssh/id_rsa` 同形）
>     ⇒ ★★**乙（`Screen` 认 `L0`）修不到这一条** ⇒ ★**"日志不再说假话"这个目标需要**甲**独立承担**。
>   - ★**B：`Decision.ModeKept` 也在每一发 `L0` 上带假话**（`bridge.go:314` 写入、`gate.go:48-50` 自述"names the red line"），
>     ★**但实测读取方只有测试**（`ModeKept` 全仓 12 处、`cmd/**` 与 `internal/panel/**` **零处**）⇒ ★**今天唯一操作员可见面就是那行日志**；**行为面无回归**（L0 仍放行）。
> **⑥ ★尺与射程（它把"落笔前五问"逐条答了）**：仪器＝`go test -overlay`、探针在**仓外**、经**虚拟路径**参与编译 ⇒ ★**仓内 `.go` 字节 0 枚**（`git status --porcelain -- internal cmd` 首尾均空）。
>   ★**组合＝生产组合同形**（`cmd/wisp/run.go:764-788`）、唯一差异＝`Logf` 指捕获器。
>   ★**卡数怎么数**：`Gate` 接口＝两张卡的面（`gate.go:110-117`：`PendingWindow`＝L1 窗口、`PendingApproval`＝C18 队列）⇒ **数这两个回调被调用几次**、**⛔ 数 UI/tick**。
>   ★⚠**它如实标明一段**未实测**：日志到操作员＝`run.go:979-983` → 滚动 JSONL（`logsink.go:129`）——★**这是读码、它没跑整条 `wisp run`**（要 golden SSE）⇒ **件里已标明**。
> **⑦ 门禁**：`d22scan 0`（clean）／长度门 `0`（`tracked=17601 over-budget=61 covered=61 not-in-roster=0`、GREEN）。★环境：`GOFLAGS=` 空、**未设 `GOTMPDIR`**、⛔ 未用 `CGO_ENABLED=0`。
> **⑧ 纪律**：`.go`＝**0**（★同尺正控 `ec520750` ⇒ **2 枚** `.go`、**尺有牙**）；⛔ 动任何 `.go`／票面／台账／`.github/**`；⛔ 翻框、⛔ push。
      ①**逐条复核上述四条链**（每条带 `file:line` 现跑复认，⛔ 抄本票）；
      ②★**造一发真读数**：跑一发 L0 工具调用，**逐字贴出那行 `MODE-REDLINE … still asks`**，并**同一发里证明确认卡数 ＝ 0**（★**这就是"假话"的直接证据**）；
      ③**正控**：造一发 **L1 或 L2**（那两档 `Screen` 认得）⇒ ★**它的日志必须 ⛔ 说 "still asks"**（或说的是**真的**会问）⇒ **证明这条链只在 L0 上假**；
      ④**票面前提新鲜度**：行号漂不漂要现跑。
- [x] **AC#1 裁形（编排者拍）**：三形各写代价＝〔甲〕**只改日志文案**（把 "still asks" 改成能说实话的措辞，例如带上 `sil.Kept` 的原文、或对 L0 单独一支）；
> **★★2026-10-22 23:5x 编排者拍形＝〔丙〕两者都做**——依据＝`363-a1` 的两条"会改变裁形"的实测
> **① ★为什么⛔ 选〔甲〕只改日志文案**
> - ★〔甲〕能修 `L0` 那句假话，★**但修不到 `Deny`**——`363-a1` 实测：`fs.read` 打 `%APPDATA%\wisp\config.toml` ⇒ `level=Deny`、`CARDS TOTAL=0`、
>   `text="该目标被绝对禁止访问（R3 A 档…）"`，★**日志仍打 `still asks`** ⇒ ★**"日志不再说假话"这个目标，甲一个人扛不住两个方向**。
> **② ★为什么⛔ 选〔乙〕只让 `Screen` 认 `L0`**
> - ★〔乙〕能修**范畴错误**，★**但修不到 `Deny` 那条**（`Deny` 是 `Screen` 认得的级、它的问题是**日志把"禁令"说成"会问"**）。
> - ★★且〔乙〕**只动 `Screen` 不改日志** ⇒ `L0` 那句假话**从"说错原因"变成"不说"**、★**但 `Deny` 那句照旧假** ⇒ **目标未达**。
> **③ ★为什么选〔丙〕（按 `AGENTS.md` ⓘⓘⓘ 铁律第 2／3 条）**
> - **①最有价值**：本票的题目是「**生产日志每一发都在说假话**」⇒ ★**判据是"日志不再说假话"**，★**而那需要在 `L0` 与 `Deny` 两个方向上都对**。
>   〔丙〕是**唯一同时覆盖两个方向**的形。
> - **②"复杂"买得到真价值吗？买得到**：〔丙〕比单形多做的事＝**修一个真范畴错误**（`Screen` 出入口不一致）。
>   ★**那个错误今天是"只影响日志"，但它是一个**认得的级被当成不认识**的 bug** ⇒ ★**修它属正确性、⛔ 美化**（`363-a1` 逐字："乙属正确性、非美化"）。
> - **③⛔ 把"动作大"当降级理由**：〔丙〕要动 **`internal/risk/mode.go` ＋ `internal/tools/bridge.go` 两处**（两包、可能撞别的腿）——
>   ★**那是代价、写清即可，⛔ 用来换小方案**。
> **④ ★〔丙〕的**硬边界**（⛔ 越）**
> - ★**⛔ 为了让日志"好看"而放宽任何判定**：`L0` **仍必须放行**（`route:432-435` 不动放行语义）、`Deny` **仍必须拒**。
> - ★**⛔ 把 `Screen` 改成"认得 L0 之后就不管了"**——★`L0` 要**明确**走出「不询问」这一支（`Silenced{Level: L0, Silenced: true}` 那个出口已经在、要走对）。
> - ★**`Deny` 那条**：日志要**说真话**（它是**禁令**、⛔ "会问"），★**但 `Deny` 的 `Kept` 语义按 `mode.go:143-146` 是"只在红线拒绝静默时非空"** ⇒ ★**落地腿必须先读清 `Kept` 契约再动**。
> - ★★**两包顺序**：`internal/risk`（`Screen` 一形）与 `internal/tools`（日志一形）**必须串行或同一腿**（★**同包只一枚写腿、跨包也要避免与别的腿撞**）。
      〔乙〕**`mode.Screen` 认 L0**（让它返回 `Kept=""`、不进"无法识别"支）；〔丙〕**两者都做**。★**未裁形任何实现都⛔ 开工**。
- [x] **AC#2 落地＋牙（写腿）**：按裁定那一形做；★**每处改动配一枚今天响、不修不响的用例**；
> **2026-10-23 10:4x 编排者翻框＝`AC#2` 成立、`AC#3` 成立（凭**非实现者** `363-v1`，件 `docs/evidence/s1/363-v1-verdict.md` **157 行**／0 CR、提交 `fede8ab88`＋`062030af3`、**零 `.go` 字节**）**
> **① 它自造的读数（⛔ 采信自报：仓外 `-overlay` 探针、真 D34 声明／真 R1–R9／真 C26，卡数＝`Gate` 两回调）**
> - **`L0` 三档**：★**无任何 `MODE-` 行**、`kept=""`、`CARDS 0/0`、`decision=allow`、`text="hello"`、`ModeKept=""`／`ModeSilenced=false`。
> - **`Deny` 两档**：`MODE-REDLINE … refuses … (a Deny never reaches a gate: no card is raised)`、`0/0`、`decision=reject`、`isError=true`。
> - ★**正控**：`L1/ask_every_step` 窗口＝1 且**没有那行**；`L1/ask_high_risk` 0 卡走 `MODE-SILENCE … not an allow`；`L2` 两档各 **1 张审批卡**且 `still asks` **是真的**；
>   ★★**它另加 `L2/ask_every_step` ⇒ 1 卡、无 `MODE-REDLINE`（没人扣住它）** ⇒ ★**证明不是"把日志全关了"**。
> **② 四档无回归（各一发）**：`L0` 放行／`L1` 窗口 1 ＋ **写入真落盘**／`L2` 审批 1／`Deny` 拒；
>   ★`bridge.go` **`+13/−0`（全文零删除）**、`mode.go` **`+25/−1` 且唯一的 `−1` 是注释** ⇒ ★**`route` 的 `case risk.L0`（现 446–448）一字未动、两文件纯插入** ⇒ ★**这是"没动行为"最硬的证明形状**。
> **③ ★`Kept` 契约未动（判语逐字核）**：`mode.go:143-146` 逐字未变、`Deny` 的 `Kept` 文本逐字未变；
>   ★**`L0` 由"无法识别…"变空 ＝ 把**违反契约的旧值**改回契约**；★`ModeKept`／`ModeSilenced` 在 `cmd/**`、`internal/panel/**` **零读取**。
> **④ ★★三条具名残余（判语 §9，我写进票尾）**
>   1. ★`mode.go:141` 那句 `Silenced is true when the caller must NOT ask` 与 `L0` 的 `Silenced=false` 在**字面**上相冲 —— ★**但修前 `L0` 也是 `false`** ⇒ **既存口径差、⛔ 本笔引入**；★**改那句属"改契约文字"⇒ 须机主批**。
>   2. ★两枚新用例注释里的**活指针行号是改前的**（`433-435` → 现 `446-448`；`437-439` → 现 `450-452`）。
>   3. ★`L2` 被扣住但 **gate 未接线**（`NoGate`）的组合里，`still asks` 仍会打在 **0 卡**上 —— ★**属"审批层没接线"**、⛔ 本笔。
> **⑤ 反向判据（它自造、⛔ 只看实现腿的）**：`M1` 摘 `mode.go:184-186` ⇒ `rc=1`（`ScreenRecognizesL0`(4)／`FourLevelsKeepTheirRoute`(3)／`L0ReadNeverClaimsAQuestion`(15)）；
>   `M2` 摘 `bridge.go:323-335` ⇒ `rc=1`、**只红** `DenyIsRefusedNotAsked`(6)；
>   ★★**互不代偿**双向各证一次**（★`M1` 下 `Deny` 仍绿、`M2` 下 `L0` 仍绿），★**且它自己的探针在同一对突变体下也各自红**（★**两把尺互证**）。
>   ★`git hash-object`＝**SHA-1/40 位**、四枚 **4/4 同串**（导出树同 4/4）。
> **⑥ ★`AC#3` 门禁（判语现跑、先 `unset GOTMPDIR`）**：`build 0`／`vet 0`／两包 `go test` `rc=1`（**两发同形 `RUN=399 PASS=397 FAIL=1 SKIP=1`，唯一红＝既存 `TestC26…`**）／
>   `./internal/perm/ ./internal/agent/approval/` **0**／★**`d22scan` 主检出 `rc=0`**／导出树 `0`（`ban #8 internal/=623`＋`ban #7=25`＋`cmd/=148`）／路径长度 **0**（GREEN、17612→提交后 17617）。
>   ★格式**甲（blob 树）：20 行全 `.scratch/**`、非 `.scratch` ＝ 0**；**乙（工作树）：0 行**、四枚 `CR=0`。
>   ★★**它报 `d22scan` 主检出今天是 `rc=0`、⛔ 派单预期的 `rc=1`，并归因对**：`docs/evidence/s1/196-r1-impl.md` **已落盘并入库**（我现跑核：在 HEAD 里、`12979` 字节），
>   引用只在票 196 腿**未提交**的 `internal/memory/schema.go:26/:157`（`git show HEAD:` **无命中**）⇒ ban #9 是"注释引的路径不存在"、**文件一在即消**
>   ⇒ ★**这是树的状态、⛔ 票 363 的功劳**（★**它主动划清、⛔ 抢功**）。
> **⑦ ★新增红 0（它自造基线）**：同一批测试文件、只把 `mode.go`／`bridge.go` 用 overlay 换回 `d59af9a^` ⇒ `RUN=399 PASS=393 FAIL=5 SKIP=1`（红＝`TestC26…`＋**4 枚新用例**）
>   ⇒ ★**终态 1 红是基线的**真子集** ⇒ 新增红 0**。★**它明确说实现腿报的 `396/2` 它没复现**（`TestResolvePerCallBudget` **逐发漂**、它两发都绿）⇒ ★**按红名集合判、⛔ 按四数判**（`A1126` 那条）。
      ★**反向判据**＝把新加的那一跳摘掉 ⇒ 用例必须**转红并逐名报红行**；三读数（变异红／还原绿／原样绿）逐字贴 ＋ **还原核**。
- [x] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet` 动过的包／`go test -count=1`（**四数分开报、⛔ 相减**）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh` 各 rc。
      ★**格式门禁甲乙两尺都报**（甲＝入库 blob、**判据**；乙＝工作树、**只作诊断**）。
      ★★**`cmd/wisp` 整包四数 ⛔ 当判据**（未判决恒 1、红名册逐发漂移、截断点后整族 `RUN=0`）⇒ 只按名作差＋逐发登记死点。
- [x] **AC#4 纪律格**：⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 179/183/185 任何原句。

## ★写面（初判，裁形后可能收窄）

- ★**最可能落在** `internal/tools/bridge.go`（日志）与 `internal/risk/mode.go`（`Screen` 的 L0 支）——★**两处都在别家地界** ⇒ 落地时**必须一次只碰一处**、且**与那两包的其他腿串行**。
- ⛔ `docs/PLAN.md`／`docs/specs/**`（契约文字须人工批准）、`thresholds.go`／golden／`allowlist.txt`／`tools/d22scan/**`／`frontend/**`／`design/**`。

## ★两条硬提醒（⛔ 踩）

1. ★**⛔ 为了让日志"好看"而放宽任何判定**：本票要的是**日志说实话**、⛔ 把 L0 改成会问。
2. ★**⛔ 只改文案就交差**——若 `mode.Screen` 把 L0 判成"无法识别的风险级"这件事本身是错的（★**L0 是认得的级**），那它**在别处也会假**（不止日志）⇒ 裁形时要判这一点。

## 排程

`363-a1`（只读＝`AC#0`）→ 编排者裁 `AC#1` → `363-r1`（写腿＝`AC#2`）→ `363-v1`（非实现者裁）。

## Progress log

- [2026-10-22 23:2x +08] agent=`编排者` did=立案（来路＝`179-r1` §7①；★四条链编排者已现跑核实：`bridge.go:311`／`mode.go:163-165`／`bridge.go:324`／`bridge.go:352`）next=派 `363-a1` 只读腿。

## ★残余（2026-10-23 10:4x 编排者按 `363-v1` §9 写进票尾，⛔ 随 `-done` 蒸发）

1. ★**`mode.go:141` 字面口径相冲**：那句 `Silenced is true when the caller must NOT ask` 与 `L0` 的 `Silenced=false` 在**字面**上相冲
   —— ★**但修前 `L0` 也是 `false`** ⇒ **既存口径差、⛔ 本笔引入**；★**改那句属"改契约文字"⇒ 须机主批**。
2. ★**两枚新用例注释的活指针是改前行号**：`ticket363_screen_l0_test.go:55` 的 `433-435` → 现 `446-448`；`ticket363_test.go:269` 的 `437-439` → 现 `450-452`。
3. ★**`L2` 被扣住但 gate 未接线（`NoGate`）** 的组合里，`still asks` 仍会打在 **0 卡**上 —— ★属"**审批层没接线**"、⛔ 本笔。

★**另记一条本轮学到的**：`d22scan` 在主检出今天 `rc=0`、⛔ 派单预期的 `rc=1` —— ★**因为 `196-r1-impl.md` 已落盘**（`git show HEAD:` 对 `internal/memory/schema.go` 无命中 ⇒ ban #9 是"注释引的路径不存在"、**文件一在即消**）
  ⇒ ★**这是树的状态、⛔ 本票的功劳**（★`363-v1` 主动划清、⛔ 抢功）。
