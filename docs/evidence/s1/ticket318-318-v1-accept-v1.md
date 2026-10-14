# 票 318 验证腿 318-v1 判语（非实现者）——`WideningRule` 入口整卡拒收守卫

- 腿身份＝**非实现者**（实现者＝318-r1；本腿零产码，判断只落在本件与票尾 bullet）。⛔ 未翻框、⛔ 未动 `Status:`/票名、⛔ 未 push。
- 起手锚（现读）＝`b38567a9`＝**派发锚，零漂移**；被检两笔 `8707a7ed`（产码）＋`f6fda72d`（交付）均为其祖先；`f6fda72d..HEAD` 对 `internal/agent/approval/` 改动 **0** 笔（`git log` 现跑）。
- 方法＝仓外 `git archive` 快照两枚：修前＝`8707a7ed^`（`C:\tmp\wisp318v1\snap-pre`）、修后＝`f6fda72d`（`snap-post`）＋**我自造的复现尺**（`v318v1_repro_test.go`，从票面文字写、不读交付尺；真 `config.Manager.AddAllowedDir` 写盘）＋**我自己的变异两发**（仓外 snap-mut，perl 原地改后逐字留 diff）。⛔ 未跑 318-r1 的 overlay 抵账（其 3 枚 JSON 只作被检对象，一眼未执行）。
- 环境：`PATH=…/third_party/sherpa-onnx:$PATH`（MSYS 形）；`GOTMPDIR=TMP=TEMP='C:\tmp\wisp318v1gotmp'`（先自建）；判据全 `-run` 小射程 `-count=1`；⛔ 未跑整包 `cmd/wisp`（该面本波归 133-r1）。
- 仓内审批包面：`git status --porcelain -- internal/agent/approval/` ＝**空**（本腿对仓零改动；下述读全在仓外快照）。

## 0. 结论先行

| 格 | 判语 | 一句话读数 |
|---|---|---|
| ① 复现形（AC#1 判据①） | **成立** | 修前：碎片真写进 `allowed_dirs`（rc=1）；修后：同一形 `ok=false`、零写入（rc=0）；摘守卫还原后又红 |
| ② 对照不回归（AC#1 判据②） | **成立** | 无标记的同一路径/真链越界读两树都仍推真父目录（`dir="…/allowed/jn"`、`…/outside`） |
| ③ 自造突变（AC#1 判据③） | **成立（两发）** | M1 摘守卫 ⇒ `ticket318_…:58/:61` 逐字红＋zz`:191` 红；M2 逐 path 跳过 ⇒ `:84` 逐字红；两发 build rc=0、仓内零突变、还原后 9 PASS |
| ④ 被改 201-a1 复现件翻面 | **诚实翻面（附条件）** | 同一条链、同步骤在；断言翻向非变松；legs 4-5 删除＝主题不可存在；条件＝该半跳无修后测试直驱生产消费方（票面已登记不可达） |
| ⑤ AC#2 边界登记 | **成立（三条全核到）** | 空路径=两树已量；8.3=不属本修（邻域族我现跑 11 PASS）；多路径=已量＋跳过语义反向开口我独立复到 |
| ⑥ 既有断言零放宽 | **成立** | 名册双向差集＝+4 新增＋1 改名，恰为交付所载；全包两树 rc=0；SKIP 1→1 同枚既有、零新增 |

题外观察两条（**不否决**，供编排者登记）：
1. 修后生产侧「一直」路径的拒绝文案（`cmd/wisp/approval_always.go:85-91`）对**单路径**标记卡会打出「这张卡涉及 1 个目录…目录不唯一就不存」——字面不准（原因是标记不可存，不是目录多），但 fail-closed、无写入；属 `cmd/wisp` 面（本修未触，票亦未要求）。**未复跑该包，阅读数为源码现读**。
2. 翻后的 zz 件头注两处 `replies.go:503` 是**加注释前的旧行号**（现 `WideningRule` 起 `519`、守卫在 `521`）——行号腐烂，文档级瑕疵，无功能后果。

## 1. 复现形三态（AC#1 判据①）—— 成立

### 1.1 我自造的最小等价形（卡级直造＋真写盘；尺＝仓外 `v318v1_repro_test.go`）

等价箭头与 201-a1 同：`ReplyCard{Paths:[<raw> + " (无法规范化: …)"]}` → `WideningRule()` →（ok 时）真 `config.Manager.AddAllowedDir(dir)` → 读回真 `config.toml`。

- **态1 修前红**（`snap-pre`，rc=1；全量 749 行在仓外，摘录 `.scratch/wisp/probes/318/v1/logs/10-prefix-repro-red.txt`），逐字：
  - `A1 reading: ok=true dir="C:/proj/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction" rule="[fs] allowed_dirs += \"…(junction\"" writeAttempted=true writeErr=<nil>`
  - 真写盘：`allowed_dirs = ['C:/proj/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction']`（在 A2 的 config dump 里逐字；`writeAttempted=true, writeErr=<nil>` ⇒ 两道闸都过）
  - 空路径形：`B1 reading: ok=true dir=" (无法规范化: tools: empty path)" …` → `allowed_dirs = ['(无法规范化: tools: empty path)']`
- **态2 修后绿**（`snap-post`，rc=0；`.scratch/wisp/probes/318/v1/logs/20-postfix-repro-green.txt`）：
  - `A1 reading: ok=false dir="" rule="" writeAttempted=false writeErr=<nil>`；`A2 … (fragmentLanded=false)`、`[fs] allowed_dirs = []` 原样
  - `B1 reading: ok=false dir="" rule="" writeAttempted=false writeErr=<nil>`；B2 同空
- **态3 还原**（把修后树守卫摘掉＝回到修前产码形状；`.scratch/wisp/probes/318/v1/logs/30-mutation-noguard-red.txt`）：
  - 同尺同形即回到 `ok=true dir="…(junction"`（rc=1）——证明红/绿由这 5 行守卫驱动，非环境
  - 拷回 pristine `replies.go`（`diff -q` 逐字节同）后 rc=0（`.scratch/wisp/probes/318/v1/logs/40-restore-and-fullpack.txt`）

### 1.2 复合面 201-a1 真链（真 NTFS junction；`mklink /J`，本腿现跑）

把**修后的** zz 件原样放到修前树跑同一把尺（`.scratch/wisp/probes/318/v1/logs/50-zz201a1-flip-evidence.txt` §a）：
- 修前红：`READING-2b … dir="C:/tmp/…/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction" ok=true`，`--- FAIL: TestZZ201A1JunctionMarkerCardIsRefusedForWidening`
修后绿（§b）：`READING-2` 标记仍在卡面逐字、`READING-2b dir="" ok=false`、`READING-4/5/6 ok=false 空`、`READING-11 ok=true 真父目录`、rc=0。

## 2. 对照不回归（AC#1 判据②）—— 成立

- 卡级（**两树同读数**，logs 10/20）：`C1`＝把 A 的路径**去掉标记** ⇒ `ok=true dir="C:/proj/allowed/jn"`、rule 逐字；`C2` 普通过界读 ⇒ `dir="C:/proj/outside"`；`D3` 两干净 path 同一目录 ⇒ `ok=true dir="C:/proj/outside"`。
- 真链面（logs 35）：`READING-11 … dir="…/001/outside" … ok=true`，且件内 `EqualFold(cdir, 真父目录)` 为硬断言（对照腿比修前**加强**）。
- 修前树旧 zz 件的对照腿同样绿（全包 rc=0，§6）⇒ 控制组在两树都成立，守卫不误伤干净卡面。

## 3. 自造突变（AC#1 判据③）—— 成立（两发，均仓外）

- **M1 摘守卫**（diff/build 逐字在 `probes/318/v1/logs/30-mutation-noguard-red.txt`）：
  - 变异落地：`diff -u` 恰删 5 行守卫块；同树 `go build ./internal/agent/approval/` **rc=0**（逐字 `build rc=0`）
  - 指名的钉**红且逐字点名**：`ticket318_widening_marker_rulers_test.go:58: junction card as 201-a1 measured it (forward slashes): marked path still offered as a widening base: dir="…(junction" …` 与 `:61 … a refusal must clear dir and rule, got (…)`；同发 zz `:191` 红、`--- FAIL: TestTicket318AMarkerCardOffersNoWideningTarget`
- **M2 换判词＝逐 path 跳过**（`…/logs/31-mutation-perpathskip-red.txt`）：
  - 落地：删守卫＋派生循环首行 `continue`（diff 逐字）；build rc=0
  - 只有"整卡拒"那条钉红：`ticket318_widening_marker_rulers_test.go:84: one marked path must refuse the whole card, got dir="C:/proj/outside" rule="[fs] allowed_dirs += \"C:/proj/outside\""`
  - **我的混合卡尺 D1 同树独立复到"反向开口"**：`ok=true dir="C:/proj/outside" writeAttempted=true`（修前它也是 `ok=false`——跳过语义把修前拒的卡变成推一条规则）⇒ 交付选"整卡拒"的理由成立
- 两发后均拷回 pristine、复跑 `9 PASS rc=0`；仓内 `git status -- internal/agent/approval/` 全程为空。

## 4. 被改的 201-a1 复现件（+56/-72＝128 行）判语 —— **诚实翻面（附条件）**

判据（削牙会命中的形状）逐条核对：断言由 `t.Fatalf` 变 `t.Logf`／腿数减少但主题仍存在／对照被删弱／改名断掉 `TestZZ201A1` 前缀／旧读数被就地改写——**无一项命中**。逐腿：

| 腿 | 修前 | 修后 | 判定 |
|---|---|---|---|
| 前置（真 junction／正控读真字节） | 有 | **逐字未动**（diff 纯上下文行） | 在 |
| leg1（真 bridge 拒 ⇒ 卡面带标记，READING-1/2） | 有 | **逐字未动**（仅注释） | 在 |
| 2b（反斜杠拼法） | `!bsOK\|\|!HasSuffix(dir,"(junction")` ⇒ 要求碎片 | `bsOK\|\|bsDir!=""` ⇒ 要求拒 | **翻向**（同构镜像，非变松） |
| leg2（真 Gate/Queue/Replies、ledger 记账，READING-3） | 有 | **逐字未动** | 在 |
| leg3（READING-4/5/6，THE MEASUREMENT） | 5 条碎片断言（ok、非真父目录、含标记、尾碎片、rule 逐字） | 2 条拒断言（`if ok` 即败＋`dir/rule` 必须双空） | **翻向**；旧断言在 `ok=false` 下全部无事可做，新断言方向更窄（只许拒） |
| leg4（碎片进真 config.toml） | 有（READING-8/9） | **删除** | 主题不可存在（守卫后 `dir==""`，无物可写）；修前读数原名原字保留在 `probes/201/a1/logs/run3-full.txt`（READING-8/9＋config dump 垃圾行，本腿现读） |
| leg5（第二张卡带碎片） | 有（READING-10） | **删除** | 同上；生产 `approval_always.go:85-91` 在 `AddAllowedDir:134` 前折返 ⇒ 下半跳不可达（源码现读） |
| 对照腿（READING-11/12） | 有（弱：只查不含标记） | **加强**：`!cok` 即败＋`EqualFold` 真父目录硬断言 | 在且更强 |

- **翻面真实性**：同一（修后）文件在修前树红（2b，逐字见 §1.2）、修后树绿（rc=0）——确证是"同一形翻面"而非"删掉旧形"。
- **前缀保留**：`TestZZ201A1…` 前缀未断（改名 `TestZZ201A1JunctionMarkerCardIsRefusedForWidening`），`-run TestZZ201A1` 仍命中。
- **无新跳过/无环境门**：件内零 `t.Skip(` 调用（仅头注提到 d22scan 对 env-gated skip 的禁令）；`//go:build windows` 为原有。
- **附条件（具名残余面）**：修后无测试**直驱生产消费方**（cmd/wisp「一直」路径）遇到标记卡这一跳——leg5 的删除对象在测试面无法重造，只能靠源码在线上（`:84 → :85-91` 折返）与 318-r1 的 `logs/44` 邻域读数；票面 §"最坏后果"已登记该半跳自本修起不可达。**不构成削牙**，但"下游不再开口"的活体读数止于包边界。
- 头注措辞小瑕：`legs 3–5 below flipped their assertions` 里 legs 4-5 实为删除、非翻；逐字见 §0 观察 2 的行号腐烂。

## 5. AC#2 边界登记核 —— 成立

| 边界 | 落地者登记 | 我的独立读数 | 判语 |
|---|---|---|---|
| 空路径标记 | 已量（修前整条标记当目录） | **成立**：修前 `dir=" (无法规范化: tools: empty path)"` ＋真写 `allowed_dirs = ['(无法规范化: tools: empty path)']`；修后拒（logs 10/20 的 B1/B2） | 已量核到 |
| 8.3 短名 | 不属本修＋反证钉 | **成立**：`D5` 两树同读 `ok=true dir="C:/PROGRA~1/notes"`（守卫盲于无标记串）；上游 C26 折叠（`paths.go:139-158` 现读）不产标记；我现跑 `internal/tools -run 'Junction\|ShortName'` **11 PASS rc=0** | 不属本修成立 |
| 多路径 `ok=false` | 已量（混合=全拒；修前也拒） | **成立**：`D1` 混合两树 `ok=false`、`D2` 全标记两树 `ok=false`、`D4` 两干净两目录 `ok=false`；"整卡拒 vs 逐 path 跳过"的差别由我的 M2 独立复到（`:84` 红＋D1 反向开口） | 已量核到 |
| `Aggregate`（另一枚 `dirOf` 消费方） | 具名不属本修 | **成立**：`batch.go:35-37` 现读 `d.Level != risk.L1 ⇒ nil`；标记 ⇒ R2 fail-closed ⇒ L2（我 log 33 的 `READING-1 level=L2 rules=[R2 R3]`）⇒ 标记进不了其射程 | 具名成立 |
| 已知假阳面（文件名含「无法规范化」四字） | 登记为沿既有口径 | **认可**：`grants.go:264`、`bridge.go:1273` 现读同为 `strings.Contains(p, "无法规范化")`，本修不新造更细判据 | 登记如实 |

## 6. 既有断言零放宽 —— 成立

- **名册双向差集**（`go test ./internal/agent/approval/ -list '.*'`，剔除两树同挂的我方 `TestV318V1*`；`probes/318/v1/logs/60-roster-boundary-census.txt`）：修前 **79** → 修后 **83**；差集**恰**＝`+TestTicket318AMarkerCardOffersNoWideningTarget`、`+…MarkedPathRefusesTheWholeCard`、`+…TheSamePathWithoutTheMarkerStillDerivesItsParent`、`+…ShortNamesAreNotTheMarker`、`TestZZ201A1JunctionMarkerBecomesAWideningFragment → …CardIsRefusedForWidening`（改名）。**无其它任何增减**。
- **全包**（同命令、`-skip TestV318V1` 排除我方尺）：修前 `ok 1.953s` rc=0／修后 `ok 2.159s` rc=0；`--- PASS` 74→78（+4 新增）；`--- SKIP` 1→1，同枚既有 `TestDefaultDeadlineWallClockMeasurement`（ticket 84 设计内 env 门）⇒ **零新增 SKIP**。
- **改动面**（`git diff --name-status 8707a7ed^ f6fda72d`，非 docs/probes 行）：`M internal/agent/approval/replies.go`、`A …/ticket318_widening_marker_rulers_test.go`、`M …/zz_201_a1_junction_marker_windows_test.go`——**恰 3 件**；`replies_201_test.go` 等既有尺一字未动（不在名册差集、不在改动表）。

## 7. 具名没做

- ⛔ 未跑整包 `cmd/wisp`（本波面归 133-r1，且仓内 `cmd/wisp/main.go` 有他腿在飞改动）；生产侧只用源码现读（`approval_always.go:84/85-91/134`）＋采信 318-r1 `logs/44`，**未复跑**。
- ⛔ 未复跑 r1 的 session grants 先例族（本修物理不触 session——`approval` 无该依赖方向；先例口径只做 `grants.go:264/193-197`、`bridge.go:1273`、`bridge.go:1028` 逐字现读）。
- ⛔ 未端到端真跑 `wisp run` 遇真 junction（我的"真链"止于包内：真 NTFS junction→真 bridge→真 Gate/Replies→真 WideningRule＋真 `config.Manager` 写盘；上层 UI 面未触）。
- ⛔ 未用 overlay；未动 `tools.Decision.Paths`／`displayPaths`／`cmd/wisp/**`／`thresholds.go`／golden／`PLAN.md`／`specs/**`／`frontend/**`／`design/**`。
- ⛔ 未删任何临时件；仓外快照与突变副本留盘（`C:\tmp\wisp318v1\`：snap-pre/snap-post/snap-mut＋work/ 里 `v318v1_repro_test.go`、`mut-noguard-replies.go`、`mut2-perpathskip-replies.go`、`zz-201a1-{pre,post}.go`）。

## 8. 读数名册

本腿件（`.scratch/wisp/probes/318/v1/logs/`，7 枚，全 ≤200 行）：`10-prefix-repro-red.txt`(31) · `20-postfix-repro-green.txt`(38) · `30-mutation-noguard-red.txt`(51) · `31-mutation-perpathskip-red.txt`(44) · `40-restore-and-fullpack.txt`(30) · `50-zz201a1-flip-evidence.txt`(158) · `60-roster-boundary-census.txt`(41)。
仓外全量原日志 21 枚在 `C:\tmp\wisp318v1\logs\`（含 10/32 各 ~24KB 全量、53/54 全包 verbose）；快照三枚与突变副本见 §7。
