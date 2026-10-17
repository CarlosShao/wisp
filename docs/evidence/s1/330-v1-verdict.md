# 票 330 · 非实现者终裁腿 330-v1 · 八格逐格判语

- 被裁对象＝落地腿 `330-r2`（`docs/evidence/s1/330-r2-g2-teeth.md` ＋ `.scratch/wisp/probes/330/r2/**` ＋
  票尾一条 bullet）。起手锚 `feab45e3`、交件尖 `b9069dfa`；`git diff --stat feab45e3..b9069dfa -- internal/agent`
  ＝只两枚测试件（`plan_present_crowded_world_330r2_test.go` ＋272／`plan_present_doorbell_165r3_test.go` ＋3−1），
  `b9069dfa` 之后 `internal/agent` 零枚别腿字节（`git log b9069dfa..HEAD -- internal/agent` 空）。
- ★读数一律取具名提交树 `b9069dfa`：本腿自建**仓外导出树** `D:\tmp\w330v1`（`git archive b9069dfa`
  八件套＋`docs`＋手工三件＋`GOTMPDIR=D:\tmp\w330v1tmp`），两枚目标文件 sha256 逐位等提交内基线
  （doorbell `6542773d…`、crowded `891fbf5c…`）。⛔ 读工作树在飞字节、⛔ 改跟踪件。突变全打在仓外树、
  逐发 `finally` 还原并核回基线哈希（`matches_committed_baseline=True` 全程真）。
- 凭据独立性：`330-r2` 交回的成对两发**编排者没复跑**，故本腿**自己重放**（见 `probes/330/v1/logs/`），
  ⛔ 抄它的读数。下面凡引数皆标"哪把尺＋射程＋量/算"。

## AC#1 · G1 对接钉＋正控两发 —— 不适用

`330-r2` 未交付此格（其证据 §8 与票尾 bullet 都写"AC#1/#2 仍空，等 `330-r1` 在 `cmd/wisp/**` 落对接钉"）。
现量：`grep -rn planPresentName165r3 cmd/wisp`＝**零命中**、`cmd/wisp` 无任何 `*330*` 对接件（唯一含 roster
字样的是票 253 旧件）。⇒ 去"重放正反两发"暴露的是**这枚对接钉今天不存在**，⛔ 判成立、也⛔ 判不成立
（不成立要有钉可打红）；本格属 `330-r1` 的 `cmd/wisp` 面，本裁对象够不着，具名判**不适用**。

## AC#2 · ⛔ 新名册假对接 —— 不适用（并记一条正向事实）

同 AC#1：对接钉未落地，"名册多一枚／少一枚今天这枚钉会不会响"之问对**不存在的钉**无对象，判**不适用**。
正向事实（防"恒真第三形"预置）：`330-r2` 的 crowded-world 夹具⛔ 抄第二份镜像，而是组合
`planPresentDirectory165r3()`（`crowded_world:106`），且头部注释（`:66-70`）明确声明它**只测驻留臂机制、
⛔ 声称对接真名册**（真名册那半在 `internal/tools`、`agent` 包因 `internal/tools/bridge.go:13` 已引 `agent`
而够不着）。⇒ 它**没有**制造 AC#2 要防的"假对接"；AC#2 仍待 `330-r1` 落那枚真读装配名册的钉。

## AC#3 · G3 那句注释说实话 —— 成立（且经本腿独立突变背书，非注释自述）

`:131` 现量三行逐字：

    // REDUNDANT NON-ZERO GATE, not THE JUDGMENT: 165-v2 measured that feeding
    // sinkCount165r3 a constant keeps this green, because the tooth of this case
    // is the independent row walk at line 145 below (rowOfTool165r3, :95-105).

本腿⛔ 让"注释改口"替"牙在别处"背书，而是把那一发自己打了一遍：在仓外树把
`n := sinkCount165r3(rows, planPresentName165r3)` 换成常量 `n := 1`，单跑
`TestPlanPresentCallBooksTheToolCallRow165r3` ⇒ **rc=0 仍绿**（`logs/v1-ac3-sinkcount-constant.log`）。
⇒ 那扇 `n==0` 非零门确属**冗余**，真牙在 `:145` 的 `rowOfTool165r3`（`:95-105` 自带独立 `for` 遍历、不复用计数、
无行即 `t.Fatalf`）。注释与那扇门真实射程**一致、说实话**。
⚠ 票面 AC#3 原文写"牙在 `:143`"＝改注释前（1 删 3 加＝净 +2 行）的旧号；现文件真实行是 `:145`，
注释也写 `:145`＝对准现量、并非用旧号糊弄。★产码字节零改动：`git diff --name-only feab45e3..b9069dfa -- internal/agent
| grep -vE '(_test\.go|testdata/)'`＝空、rc=1；正控＝克隆树里把产码 `inject.go` 的字面量 `MaxRetrievedTools=5`
改 `5+1`，同一把尺输出 `internal/agent/inject.go`、rc=0 ⇒ 这把尺两种形状都分得出、⛔ 失明。

## AC#4 · G2 归口具名（缺世界已造出＝真验过档位）—— 成立

`330-r2` 交回的是"缺世界已造出（多候选夹具）"。本腿独立核它**真把两臂世界造出来了**（逐发重放，
`logs/v1-m1..m5`，尺＝`go test ./internal/agent/ -count=1 -timeout 30m -v -run <三枚点名用例>`）：

- 正控（未突变、点名三枚）：rc=0，三枚全 PASS。
- 去牙发 **M2**（长尾候选 `Resident: false→true`）⇒ rc=1，**唯一**红在 `:247`（`tier leak: resident arm =
  [plan.present ledger.sweep330r2] …`）；`TestCrowdedWorldStillOffersTheBrake330r2` 同发 PASS。
- 反控发 **M3**（查询换与长尾**无共同字** `"frobnicate widget quietly"`）⇒ rc=1，**唯一**红在 `:254`
  （`retrieval arm = [], want "ledger.sweep330r2" in it …`）。
- ★M2 与 M3 **落点行号不同（247 ⇄ 254）、红句不同** ⇒ 这枚用例响的是"档位门"与"检索门"各自、⛔ 恒响。
- G2 那一发重放 **M1**（镜像 `Resident: true→false`）⇒ rc=1，红在 `:188`（第一轮请求体里那扇门整个不见了）
  ＋ `:240`（驻留臂 `[]`），同发 165 那枚用例仍 PASS。⇒ 驻留臂与检索名额**同时**被拨动＝两臂世界确实存在。

⇒ 判**成立**：不是"写了夹具"就当"验过档位"，档位此时确实决定谁被注入。

## AC#5 · 台面与颜色 —— internal/agent 半格成立；cmd/wisp 半格不适用

- internal/agent（本裁对象面）：仓外树配方齐（八件套＋`docs`＋手工三件、目录名 8 字符、`GOTMPDIR` 先建指 `D:`、
  取出 `.go`=750 证非空）。★本腿独跑整包 forward（`-count=1 -timeout 30m -v`，尺＝逐枚 `grep -c`）＝
  **rc=0／RUN 112／PASS_TOP 88／FAIL_TOP 0／SKIP 0／PANIC 0**（`logs/v1-forward-fullpackage.log`）。
  ⛔ 与 `330-r2` 自己那台或编排者 10-17 复跑的数**相减**——三把台面不同；此处只报**本腿这一台**的数。
- cmd/wisp 整包那一格＝**不适用**：属 `330-r1`／`cmd/wisp` 面，`330-r2` ⛔ 碰那面一字，本腿⛔ 跑 `./cmd/wisp`
  （归 `228-r3`／`329-v1`）。

## AC#6 · 门禁四数 —— 成立（只算 internal/agent 这一面，逐枚 rc，⛔ cmd/wisp）

只在共享工作树跑（`logs/v1-gates-and-rulers.log`）：

| 尺 | rc | 具名 |
|---|---|---|
| `gofumpt.exe -l <两枚交付件>` | 0 | 输出空＝本腿交付件干净；`-l internal/agent/` 列 `approval/pending_read.go`／`tools.go`＝本锚点既存漂移、非 `330-r2` 字节 |
| `sh scripts/d22scan.sh` | 0 | 现量全清；交付当时那条 `cmd/wisp/panel_snapshot_deliver_329.go:14` 属并行腿 329 字节、⛔ 本面 |
| `sh scripts/check-path-length-budget.sh` | 0 | VERDICT GREEN |
| `go vet ./internal/agent/` | 0 | 输出 0 字节 |

## AC#7 · 时序格（本票须排在票 329 之后）—— 不适用

这格量**排程纪律、不是字节**。票面 AC#7 的前提是"甲形落 `cmd/wisp`、与 329 同枚共享编译面须串行"。
本裁对象 `330-r2` **只写 `internal/agent`**（与 329 的 `cmd/wisp` ⛔ 同枚共享编译面），起手锚 `feab45e3`＝第十四波
派单，当时 `internal/agent` 测试面无别腿在写。★今天前提已变：`329-r1` 于 09:58 撞 150 轮帽、编排者 10:0x
代提其遗产（`7535d83b`）、`cmd/wisp` 写面已交 `228-r3`。⇒ 判**不适用**，依据＝`330-r2` 从未与 329 争同一枚
`cmd/wisp` 写面；本判语⛔ 出现"与在飞腿并发写同面"这类今天不存在的事。

## 必答题（能力类格一律要问：钉今天被产码读到还是只被测试读到）

- **G2 钉住的机制被产码读到**：`Inject` 两支臂切分＝`internal/agent/inject.go:47`（`if t.Resident`）；
  检索削名额＝`inject.go:55`（`bm25Select(..., MaxRetrievedTools, b.SectionToolIndex)`）与 `:238`
  （`if used+cost > budget && len(out) > 0`）。产码读者＝`internal/agent/loop.go:573`
  （`inj := Inject(latestUserText(l.History()), dir, l.b)`，逐枚 file:line 已核）。⛔ 用"用例绿"背书"接线真"——
  本腿是另指 `loop.go:573` 这枚产码调用点、⛔ 靠 `go build`/`vet` 的 rc=0（Go 不为未调用方法报错）。
  crowded-world 夹具与其断言本身＝测试面（测试件天生不被产码读）。
- **G1 对接钉今天不存在** ⇒ 谈不上被产码读：真名册名＝`internal/tools/plan_present.go:72`（`Name()→"plan.present"`）、
  镜像名＝`internal/agent/plan_present_doorbell_165r3_test.go:38`（const `planPresentName165r3`）；
  二者今天**零枚用例对接**（`cmd/wisp` 无对接件），正是 AC#1/AC#2 待 `330-r1` 落的那枚钉。

## 交件自证（三把尺）

- 提交：探针面 `0db57b8e`（`git show --name-only` 只落 `probes/330/v1/**` 11 枚）；本件＋票尾 bullet 各另计。
- 证据件＝`docs/evidence/s1/330-v1-verdict.md`；原始读数 `.scratch/wisp/probes/330/v1/logs/*.log`（9 枚，皆非 0 字节）；
  驱动器 `.scratch/wisp/probes/330/v1/replay-330v1.py`。分词形未填标记自尺＝0 命中（见交件回执）。
- 决定性尺逐字（起手 ⇄ 交件、成对两发红句、还原哈希、改前改后四数、门禁 rc）全在 `logs/`，
  ★成对两发＝本腿自己重放、⛔ 抄 `330-r2`。⛔ 翻框／⛔ 改判据文字／⛔ push。
