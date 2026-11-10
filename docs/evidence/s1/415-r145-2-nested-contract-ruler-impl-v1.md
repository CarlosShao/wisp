# 415 (R-145-2) — 合同尺伸到嵌套 section：落地读数与 AC#5 交付清单

- 时刻：2026-11-10 17:0x +08；起手锚 `dd2974e73`；交付 commit **`551d2a16d`**（只 commit `internal/panel/composer_test.go`，390 行，未 push）。
- 写面：`internal/panel/composer_test.go`（本件＋票面 Progress log；`frontend/**`、`design/**` 零触碰）。

## AC#1 尺的射程（形状与理由）

- 旧尺 `pairs` 只配 6 顶层类型（Snapshot/ComposerState/ModeView/WorkspaceView/AttachmentRef/ResultChunk）；票 145 之后上的 9 个嵌套 section 全部够不着。
- 新形状＝**双侧类型递归走查**（不是再抄一张枚举表——枚举表本身就是这次翻车的缺口形状）：从 `Snapshot` 的 json 形状出发（剥 ptr/slice/array、`time.Time` 保持叶子），每个可达 Go struct 必须 ① 在 `composerTSNameMap` 有页侧接口翻译 ② 在 `frontend/src/lib/panel.ts` 有同名接口；TS 侧字段类型按同一规则跟随（`ComposerAttachment[]`、`RunErrorRowView | null` 解析；`CredentialState` 字面量联合＝叶子）。
- 双向 exact-set：可达但无翻译＝红（**这正是"明天新加 section"被尺接住的那一发**）；翻译了但走查没走到＝红（stale）；TS 字段类型不指认翻译名＝红。
- 为什么按主流 harness 判：quicktype / ts-rs / OpenAPI 族合同测试都是从类型**派生**合同，不是手抄清单；手抄清单正是 R-145-2 的病根。

## AC#3 先判后改：11 枚 Go-only 键的逐键判定（全部"该在 TS 侧存在"）

判定基准：每枚键都是 Go 产码真实上线的 wire 事实，Go 注释里各有其生产者与语义；页侧声明欠给 UI 腿（`frontend/**` 归外部 agent，Q-51）。尺上落成 **`composerTSGapLedger`＝exact-set、自翻红的注册表，不是"跳过"**——panel.ts 声明了注册键⇒红（stale，删注册）；Go 加了别的未声明键⇒红（unregistered）；Go 删了注册键⇒红（dead）。缺口从"所有仪器都看不见"变成"11 次具名＋两侧任一动即红"。

| # | Go 键（位置） | TS 接口缺它 | 该在的理由（一句话） |
|---|---|---|---|
| 1 | `ModelOption.display`（model_switch_187.go:51） | `ModelOptionView` | 145 AC#2b① 的人读名，picker 行要用它渲染 |
| 2 | `ModelSwitchSection.adapter`（model_switch_187.go:90） | `ModelSwitchSection` | 页面区分"模型没档"还是"adapter 没映射"需要它 |
| 3 | `ModelSwitchSection.adapterThinkingLevels`（model_switch_187.go:101） | `ModelSwitchSection` | 同上，空 offer 的两种解释靠它分开 |
| 4 | `ToolRunRow.argsJson`（run_feed_145.go:98） | `ToolRunRowView` | 已终结工具行的取证参数，forensics 表该显示 |
| 5 | `ToolRunRow.level`（run_feed_145.go:100） | `ToolRunRowView` | 行的风险档（L0/L1/L2），145 行 4 的风险维度 |
| 6 | `RunSection.running`（run_feed_145.go:167） | `RunSection` | 在飞调用行（子型 `ToolRunningRow` TS 无接口，TS 落键时尺会连带点名要翻译） |
| 7 | `RunSection.usage`（run_feed_145.go:172） | `RunSection` | 累计花费（子型 `RunUsageRow` 同上） |
| 8 | `RunSection.done`（run_feed_145.go:196） | `RunSection` | 终止词 llm.StopReason（子型 `RunDoneRow` 同上） |
| 9 | `RunErrorRow.retryable`（run_feed_145.go:128） | `RunErrorRowView` | class-policy 的逐字答案，false 是信息不是缺失 |
| 10 | `GitView.uncommitted`（git.go:178） | `GitView` | 未提交维度（子树 `GitUncommittedView/GitUncommittedCounts/GitChangedFile` TS 全缺，落键时尺点名要翻译） |
| 11 | `GitView.diff`（git.go:181） | `GitView` | diff 维度（子树 `GitDiffView/GitDiffFile` 同上） |

票面点名的 6 枚（display/adapter/adapterThinkingLevels/argsJson/running 等）全部在册；walk 首跑实测的完整缺口集就是上表 11 枚，没有第 12 枚。

## AC#5 交付清单（交外部 UI agent；本票写面不含 frontend/**）

请 `panel.ts` 的主人补 11 枚键声明（全部建议 optional，照该文件既有的 reader-gated 注释风格），并按需补 8 个子接口（映射名已在尺的期望里，落地即被普通双向对比接管）：

1. `ModelOptionView` += `display?: string`
2. `ModelSwitchSection` += `adapter?: string`、`adapterThinkingLevels?: string[]`
3. `ToolRunRowView` += `argsJson?: string`、`level?: string`
4. `RunSection` += `running?: ToolRunningRowView[]`、`usage?: RunUsageRowView`、`done?: RunDoneRowView`（新接口 `ToolRunningRowView{callId,name}`、`RunUsageRowView{tokensIn,tokensOut,cachedTokens?,micros,currency}`、`RunDoneRowView{stop?}`）
5. `RunErrorRowView` += `retryable: boolean`
6. `GitView` += `uncommitted?: GitUncommittedView`、`diff?: GitDiffView`（子树形状见 `internal/panel/git_readonly_exec.go:135/:200`）

落地顺序注意：每声明一枚注册键，`composerTSGapLedger` 对应条目必须同 commit 删除（否则尺报 stale 红）；子接口名照上表拼写，否则尺报"no page-side interface"红点名。

## AC#2 反向判据读数（两次）

- 造突变（对真实 panel.ts 文本删嵌套键 `toolsDropped`，Go 侧仍在）：**红 1 枚并点名**——`Snapshot.run emits toolsDropped and interface RunSection does not declare it - a nested key one side grew and the other never saw...`（`go test -run TestPlantedNestedContractDriftGoesRed -v`，`red as required: 1 violation(s)`）。
- 还原（真树）：**绿 0 枚**——`nested contract held: 26 Go view-model types translated and walked against panel.ts, 11 judged page-side gap(s) pinned by exact set`。
- 附加反向钉全绿：幽灵键（TS reads ghostRenderKey ... never sends）红并点名；账本吸新缺口红、藏已闭合缺口红；可达 Go 型剥翻译红点名 `GitWorktree`。

## AC#4 门禁读数

- `gofumpt.exe -l internal/panel/`：空（rc0）。
- `go vet ./internal/panel`：rc0。
- `go build ./internal/panel/... ./cmd/wisp/...`：rc0。⚠ `go build ./...` rc1＝**预存在**：未跟踪且被 .gitignore:14 忽略的 `build/mut_rvw.go`（他腿遗留，`git ls-files` 空、本腿零触碰；本 commit 是纯测试文件，不参与 go build）。
- `go test ./internal/panel -count=1`：合同面全绿（本腿 2 枚新/展测试＋其余全部 PASS）；唯一红＝**预存在 HEAD 红** `TestC21DesignTokensFourWayAgree`（其读的 `internal/ball/tokens.go`、`docs/evidence/s1/c21-native-tokens.md`、frontend css 在 git 里全部干净，`git status --porcelain` 该域只有本腿的 composer_test.go——与本写面无关，归口 65/239 族）。
- `sh scripts/d22scan.sh`：**clean**（452 产码 Go 文件，0 findings；positive control PASS=45 FAIL=0）。
- `sh scripts/check-path-length-budget.sh`：**VERDICT GREEN**（19262 tracked paths，62/62 全在名册）。

## 没做到的（诚实）

- `frontend/src/lib/panel.ts` 的 11 枚键补齐**没做也不能做**（写面禁区）；清单在上方，等外部 agent。
- `docs/evidence/s1/` 本件与票面 Progress log 的追加按派单"只 commit internal/panel/"留在盘上未 commit，请编排者收件时一并收。
