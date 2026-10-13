# 票 173 × 验收腿 `173-v1` —— 非实现者裁决表（`AC#1`..`AC#5` 逐格；翻框归编排者）

- 裁决腿＝`173-v1`（**非实现者**：⛔ 写过票 173 任何产码/用例件）；起手 HEAD＝`68c294f3`（`2026-10-13 22:27:07 +0800`）；分支 `dev`；本腿全程 `internal/tools` 树内零写入（源件 sha256 起手/收手同为 `e07b6692…d94`，用例件 `f1b4622d…eb5`）。
- 被验提交：`ec2f0864`（`internal/tools/bridge.go`＋新件 `audit_allowlist_scope_173_test.go`）＋`0f8114db`（票尾 bullet＋7 枚 probes；编排者 `1e842d63` 把误落仓库根的 `probes/173` 整树 `git mv` 回 canon，纯 rename 0 行）。
- 读数件＝`.scratch/wisp/probes/173/v1/logs/{00-head,01-unfixed-repro,02-controls-unchanged,03-mutations-v1,04-consumers-and-contract,05-fixed-green-and-gates}.txt`（6 件 ≤200 行；变异一律 `-overlay`＋仓外 C: 卷副本，无一件 `.go`/`.out` 入 probes）。
- **本表⛔ 翻任何框、⛔ 动 `Status:`、⛔ push。**

大白话一句：这一味三态改动的牙我亲手验过——把仓里 `bridge.go` overlay 回未修版，两枚"没判过"的子用例立刻红在逐字 `in_allowlist_scope=true`；修后同两发 `=unknown`；两发对照（真在范围/真越界）修前修后**一字不变**；我再自造四株同族 overlay 突变，逐个把四个子用例各弄红一次（源件盘上零写入，sha256 前后同）。

| 格 | 判语 | 本腿自跑的尺（一行；件名在尾） |
|---|---|---|
| **AC#1** 未判定不再印 `true` | **成立（有牙；红/绿两向＋一枚独立旁证）** | overlay 回 `ec2f0864^` 的 `bridge.go`（仓外副本 sha `e6c685da…b611`）⇒ `-run .../unjudged` rc=1，两枚子用例红且逐字含 `… tool=fs.edit risk=L2 … rules_hit=[R1] in_allowlist_scope=true …` 与 `… tool=probe.nopath risk=L0 … rules_hit=[] in_allowlist_scope=true …`（＝先例 `01-unfixed-run.txt` 同形）；修后 rc=0、同两发 `=unknown`。旁证：旧件 `TestFSEditRefusesAPathC26CannotCanonicalize`（同一空白 path 夹具、独立于新用例）现印 `=unknown`〔01、05〕 |
| **AC#2** 对照不许放宽 | **成立（有牙；修前修后逐字不变）** | 同两发子用例在未修/修后两版上各跑一发：真在范围 ⇒ `=true`（`fs.read … decision=allow outcome=success`）、真越界（R2 族、人工拒）⇒ `=false`（`… risk=L2 decision=reject … rules_hit=[R2] in_allowlist_scope=false …`）；仅把 `t.TempDir` 随机数字换 `NUM` 后 `diff` **空**（rc=0）⇒ 修前修后一字不变。牙：v1-m2/v1-m4 逐支见下行〔02、03〕 |
| **AC#3** 射程自查（未定义即停） | **成立：非契约面成立，未停手是对的** | 自跑：`in_allowlist_scope` 在 `docs/PLAN.md`＋`docs/specs/**` 命中 **0**；`docs/` 面全部命中只剩 `evidence/`＋`reports/`（＋票面）；`SPEC-02 §3` 的 `tool_call` DDL 列表里没有它（也没有 `rules_hit`）；"审计行" 在冻结面 0 命中；"renders it verbatim" 只存在于 `bridge.go:1136-1140` 的**代码注释**、指 `Decision` 的规则列表进卡片，不是这行日志。该行由本包 `book()` 经 `Options.Logf` 拼（`bridge.go:1150-1157`），落库列 `decision/outcome/risk_level/error_class` 本改一个字节未动〔04〕 |
| **AC#4** 分母（别只修这一枚） | **成立（分母已交；"不改"裁定复核成立）** | 第二枚（审批超时 class）不改是对的：`PLAN.md:2790` 逐字 `| user_rejected | 审批拒绝/超时 | …`，D37(a)（`PLAN.md:2776`）"必须是枚举" ⇒ 拆两枚＝改 D37，须人批准（`AGENTS §0.2` 同款）；两路共用同一实参现量＝`bridge.go:1070-1071`。**"超时一味三处可读"抽一处现跑**：`Test173ApprovalTimeoutKeepsTheTimeoutFlavor` 修后绿，读数＝回执 `"审批超时未确认，已自动拒绝"`（含"超时"）＋审计行 `kind=approval-timeout … decision=timeout`；第三处 `tool_call.decision='timeout'` 有旧尺 `cmd/wisp/approval_reply_201_test.go:507-511`（本腿只静态核，见"没跑"）。另登记项抽查：`fs.nope`/坏参数照印 `risk=L0＋rules_hit=[]` 与 `bridge_test.go:666` 既有 L0 钉一致；`book()` 行 11 栏里"判定形"四栏（`risk`/`rules_hit`/`in_allowlist_scope`/`grant_id`）都进了登记表，分母读法我复核合理〔04、05〕 |
| **AC#5** 契约轴 | **成立（逐提交 0 命中）** | `ec2f0864` 文件面＝2（`bridge.go`＋新用例）；`0f8114db`＝8（票面＋7 logs）；两提交各自对禁列（`docs/PLAN.md`/`docs/specs/`/`tools/d22scan/`/`allowlist.txt`/`thresholds.go`/golden/`frontend/`/`design/`）**0 命中**、无新增豁免。票面 diff 到字节：原句一字未动，仅末行补 `\n`＋尾追加 1 条 bullet（numstat 3+/1-，那个 1- 是同一文本"无 `\n`→有 `\n`"）〔00、04〕 |

## 自造同族突变四株（⛔ 未抄 r1 的 m1/m2/m3；逐株只红在指名子用例；源件 sha256 前后同＝`e07b6692…d94`）

| 株 | 变异（仓外副本＋`-overlay`） | 红在哪里 |
|---|---|---|
| v1-m1 | `token()` 的 `return "unknown"` → `return "false"`（票面明确拒绝的第二句假话形） | 两枚未判定子用例，红在 `:75`/`:97`"must render as unknown"，读数 `=false` |
| v1-m2 | `token()` 的 `case scopeOut: return "false"` → `"true"` | `judged_out_of_range_stays_false`，红在 `:135`，读数 `=true` |
| v1-m3 | 守卫收窄 `b.paths==nil \|\| len(dec.Paths)==0` → `b.paths==nil`（空路径从未判定漏回"在范围"） | 两枚未判定子用例，读数 `=true`（＝原缺陷形状借另一处编辑复现） |
| v1-m4 | `token()` 的 `case scopeIn: return "true"` → `"yes"`（对照①给牙） | `judged_in_range_stays_true`，红在 `:113`，读数 `=yes` |

⇒ 四个子用例各自都有"弄坏必须红"的独立证明；正读（`=unknown`）还被第二断言钉死——将来若有人把这一栏**删掉不写**（票面备选形），`=unknown` 断言同样会红，所以"选形"本身在牙上。

## 消费面普查（攻点 4）

- 生产面（`cmd internal tools scripts`，去 `_test.go`）：唯一命中＝写方 `bridge.go:1154` 与其注释 `1221`——**无第二生产读者**。
- 测试面消费者只有 substr 形：`fs_edit_ac5_gate_r4_test.go:206/335`（要求 `=false` 的越界形状）、新用例 `:71/74/93/96/112/134`。
- `parse|strconv|unmarshal|split|fields|bool` 与该字段交叉命中 **0**；全仓非 Go/非 docs/非 probes 面命中 **0**；`scripts/`、`tools/` 无解析 `tools: call` 行的程序；`cmd/wisp` 的 `auditf` 是纯写方（`run.go:127-131`、`resident_approval_windows.go:536-540`）。
- ⇒ **没有解析器把 `unknown` 按 `true/false` 二分消费**；`unknown` 落地后今天无消费者会被它骗。

## 发现（具名；不阻断翻框，但须登记）

1. **同族旧注释已腐烂**：`internal/tools/fs_edit_ac5_gate_r4_test.go:350-356` 仍写 "booked line carries … `in_allowlist_scope=true`, which is 'nothing was judged about a path'"——修后**同一夹具**实测 `=unknown`（05 B 段现跑）。这句注释现在说了一句不再成立的话，正是本票打击的形状；但它是注释、不是判据，那个测试对该字段本来无牙（`r4AuditLine(..., "risk=L2")` 不含 scope 子串）。建议：由下一次 touching 该文件的腿改这一行，或编排者把它登记进票尾/台账。
2. **r1 日志行号比提交版少 1**（`probes/173/r1/logs/01` 记 `:69/:71/:74/:91/:93/:96/:110/:132`，提交版实测 `:70/:72/:75/:92/:94/:97/:111/:133`；读数逐字一致）。r1 的 01/03 取自提交前瞬时工作树；本腿在**提交版**上复跑取得同形读数，故判据成立。只作溯源注记，⛔ 扣格。

## 哪里没牙／哪里没跑（具名，⛔ 掩盖）

1. **没跑**：`cmd/wisp/approval_reply_201_test.go:507-511`（`tool_call.decision='timeout'` 那处活体）——并行腿 `186-r1` 正在写 `cmd/wisp`＋`internal/panel`，本腿⛔ 进去跑包；该处只做静态核（codepath＋既有判据）。
2. **没跑**：`d22scan` 本腿未自跑（r1 已跑 rc=0；本腿产物只有 probes 文本＋本表，不进扫描面）。
3. **没跑**：票 147 offset 那一枚的"已修"复核（r1 引 `cmd/wisp/slo_windows.go:571-618/737`；不属本票写面，只按在册接受）。
4. **无"没牙"格**：AC#1/AC#2 的四个子用例各有本腿自造突变覆盖（上表四株）；AC#3/AC#4/AC#5 是 0 命中/清单核对形，非"空集通过"（每条都给了自跑的计数）。
5. 攻点 5 的"三处可读"第三处、攻点 4 的 `cmd/wisp` 面均是**静态核**，已在上面具名。

## 翻框建议（翻框归编排者；本腿零翻）

- **建议 `AC#1`..`AC#5` 全翻**。附两条件：① 随翻把"发现 1"（旧注释腐烂）登记进票尾或台账，或派一枚一行注释修；② "发现 2"只注记、不派活。AC#4 的"不改"是裁定的正确答案（拆 class＝改 D37），翻框时应写成"分母已交、登记不修"。
- ⛔ 本腿未碰 `bridge.go`（除 overlay 瞬态，源件零写入）、⛔ 未改动件、⛔ 未翻框、⛔ 未 push。
