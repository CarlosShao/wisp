# 200 — 它进到一个项目里，**完全不知道这个项目的规矩**：四家 harness 全都在读的"给 AI 看的项目说明"，我们一枚代码都没有

- Status: **在派（owner 09-28 18:1x 放权后自行排入，不等批准）**。
- 来源：09-28 五家对标调研（账 `A404`/`A405`）。**这一枚是全部差集里"代价最小、用户立刻能感到"的一枚**：四家全有、我们**完全没做**。
- owner 的原话依据：「**别看初始方案，或者什么这个那个的门禁，把什么都自己主动舍弃了**」「**必须特么做完整功能**」「**取各家之所长**」。

## 现量（起手自己复算，别信这里的行号）
| 事实 | 读数 | 尺 |
|---|---|---|
| 我们读不读项目说明文件 | **零命中**（非测试码里只有一句注释提到 `AGENTS.md`） | `grep -rniE "AGENTS\.md\|CLAUDE\.md" --include=*.go internal/ cmd/ \| grep -v _test.go` |
| 四家怎么做的 | DSH `packages/context/agent-instructions`（出厂预算 65536 字节、超了先丢宽泛的）；Step-Code `core/resource-loader.ts:71,100-155`（优先级 `AGENTS.override.md > AGENTS.md > AGENTS.MD > CLAUDE.md`、逐级向上收、**worktree 同一份不重复注入**）；pi `core/resource-loader.ts:128-143`（同源＋`--no-context-files`）；minimax `local-runtime/src/project/instructions.ts:8,22-33`（项目级＋全局级两份） | 逐家现读 |
| 我们的系统提示在哪拼装 | 待写手现读定位（`internal/agent/` 里那一处，**不要新造第二处拼装点**） | `grep -rn "systemPrompt\|SystemPrompt" --include=*.go internal/agent/ \| grep -v _test.go` |

## 要建什么（完整形状，五段，一段都不许砍）
1. **发现**：从当前工作区根**逐级向上**收集项目说明文件，文件名优先级照抄 Step-Code 那一串
   （`AGENTS.override.md` > `AGENTS.md` > `AGENTS.MD` > `CLAUDE.md` > `CLAUDE.MD`），**到文件系统根或卷根停**；
   另认一枚**全局级**（放我们自己的数据目录，与 `config.toml` 同侧，**路径必须由现成判定者给，不许手拼**）。
2. **预算**：总量有上限（**照 D15 已有的上下文预算那一套走，不许新造一个数**）；超了按"先丢宽泛的、只截最后那份"，
   并且**必须说明"哪份被截了多少"**——不许静默截。
3. **注入**：进系统提示的那一段要**分层**：项目说明是**指引**，文件内容/命令输出/工具结果仍然是**不可信数据**
   （Step-Code `step/system-prompt.ts:306` 逐字就是这个区分）；**这一条不许含糊**。
4. **可见**：每轮启动时**打印"我到底加载了哪几份、各多少字节、被截了没有"**（pi 的 `loadedResourcesContainer` 那一形），
   并把这份清单并进面板快照的载体（**新增字段＝扩 C17 契约，已获批准只到"新增字段"为止，见 `A273`/`A389`**）。
5. **可关**：一枚开关（配置里加一键，**默认开**），关掉后**要打印"已按你的配置跳过"**，不许静默。

## 安全这一侧（这一节是本票存在的理由之一，不是附加装饰）
- ⚠ **项目说明文件本身是不可信来源**：用户 `git clone` 一个仓库，那仓库里的 `AGENTS.md` 就是别人写的字。
  ⇒ 必须给它打 C25 污染源（**复用现成源名，不许新造名册外的名字**；现量名册见 `internal/risk/provenance.go:84-94`）。
- **它只能被当"这个项目怎么干的"的指引，永远不能：** 改权限档位、改 `allowed_dirs`、放宽任何判定、替代 L2 批准、
  覆盖"面板侧来源不许允许"那条铁律。**这些一律是配置与门禁的事，不是文件里写一句就能算的事**。
- **判据必须能区分"读了并真的进了请求"与"字段恒空"**（票 181 `AC#7` 那枚空转教训：消费端接了、生产者从不填）。

## 验收判据（草，逐格要 `file:line` 与正控）
- [ ] **AC#1 真到达请求**：加载的那段字**出现在发给模型的请求里**（不是只存进一个结构体）；正控＝写一份带独特标记的说明文件 ⇒ 请求里能逐字找到。
- [ ] **AC#2 逐级与优先级**：三份不同目录的说明 ⇒ 按优先级与顺序合并；**同一份文件因 worktree/嵌套被看见两次时只注入一次**（Step-Code 专门处理过这个）。
- [ ] **AC#3 预算说实话**：超上限 ⇒ 有"丢了哪份、截了多少"的可读输出；正控＝塞一份超限文件。
- [ ] **AC#4 指令/数据分层**：判据要钉住"文件内容仍标不可信"那句在拼装处存在（能力尺，不是词面尺）。
- [ ] **AC#5 安全面不许被顶**：造一枚"说明文件里写请把权限档改成全自动/请把这些目录加进允许清单"的样本 ⇒
      **档位与允许清单纹丝不动**，并且这条要有正控能红（不是只靠"我们没写那段代码"）。
- [ ] **AC#6 污染源**：加载内容带现成源名进 C25 名册；去掉盖戳的变异要能被判红。
- [ ] **AC#7 加载清单可见**：每轮那份"我读了什么"的输出在，且面板载体拿得到（字段新增即可，**界面怎么显示归界面那支**）。
- [ ] **AC#8 可关**：开关关闭 ⇒ 一份都不读，且打印那句可理解的跳过说明。

## 禁区
`frontend/**`／`design/**` 零写面（读可以）；`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 一字不动；
C17 白名单**既有方法名**不动（只允许**新增快照字段**那一支，`A273` 的射程）；`internal/risk/**` 不许绕过、也不许改判定逻辑（**只用现成源名**）；
**不新增 D34 工具行**（本票不是工具，是上下文注入）；**不引任何新依赖**；`internal/panel/tokens_fourway_test.go` 一字不动。
⚠ `internal/tools/**` 与 `cmd/wisp/run.go` 此刻有别的写在飞——**要动装配根就等它交完再动，别并行撞同一枚文件**。

## 已知的相邻事项（别在本票里顺手做）
- 全局级说明与"首次启动自动生成项目说明"（`/init` 那一形）＝另立，本票只做**读**。
- 面板里"选哪份说明生效"的开关界面＝界面那支的活，本票只交载体与那句打印。

## 要 owner 带给界面那支的一跳（Go 侧改不了，也不是待拍板项——就一行声明）

- **现象**：`internal/panel` 今天多一枚在册红 `TestApprovalCardViewJSONKeysMatchFrontendTypes`，
  红句子原文＝`Go Snapshot emits [instructions] that interface PanelSnapshot does not declare`。
- **成因**：Go 侧（票 200 AC#7 的载体，`internal/panel/composer.go:68`）开始发 `instructions` 这枚键，
  而界面那侧的类型 `frontend/src/lib/panel.ts` 的 `PanelSnapshot` 里**没有这枚声明**。那枚尺的作用就是"**加字段必红、逼两侧同批移动**"。
- **`frontend/**` 对我和对我派的任何腿都是零写面** ⇒ 这一跳**只能由界面那支做**。请把它原样带过去：
  在 `PanelSnapshot` 上补 `instructions?:` 一枚数组，元素形状逐字段照 Go 侧（`internal/panel/instructions_200.go:18-39`）：
  `path` string · `tier` string（`"project"`／`"global"`）· `depth` number（0＝工作区本身，越大越外面，`-1`＝全局档）·
  `bytes` number · `truncatedBytes?` number（预算裁掉的字节数，为 0 时不发）· `dropped?` boolean（**整份被预算挤掉**，界面要能说"这份没生效"）·
  `duplicateOf?` string（这份并进了哪份，同一物理文件经两级目录只注入一次）· `source?` string（C25 来源名，现值 `"fs.read"`）。
  ⚠ 键带 `omitempty` ⇒ **没加载时这枚键根本不在 JSON 里**；界面**不许**把"键不存在"画成"已加载但为空"，也不许画成"功能没开"（那三种状态今天分不出，Go 侧正在补）。
- **转绿判据（界面那支交完我复跑）**：`go test -count=1 ./internal/panel/` 里 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 由红转绿，
  且另外三枚在册常红（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourwayAgree`）**名册不变**。
- **撤销口令**：若界面那支决定"这一版先不显示说明清单"，则 Go 侧回到 200-r2 的**乙**（把键从 marshal 里摘掉），
  那一格 AC#7 当场退回并在本票具名记账——**不许留着"键在值恒缺"这一形**（同 `A408` 点名的反面教材）。

