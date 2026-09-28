# 187 — 聊天框那一排的**模型**与**思考档位**要能当场改（owner 09-28 明确要；`Q-64` 改判＝做，不再走"只读显示"那一版）

- Status: **ready-for-agent（先只读普查，再落地）**。⚠ **功能票**。
- 来路：owner 09-28 14:3x 原话（逐字见台账 `A368`）＋截图里那两枚控件（模型名＋"极高"那一档，以及旁边的语音/发送）。票 181／145 AC#2b 那一版只承诺"**显示**"，本票把它**改成"能改"**——他这句改判优先于我上轮的推荐。
- 关联：**票 186**（同一片界面地界、同一枚 C17 白名单话题，两枚不同物）· **票 145 AC#2b**（快照载体：清单与当前值先要送得出去）· **票 90／R20**（"用户侧只能调严、调松要显式勾选"那条口径）

## 现量（锚 `72c76d42`，09-28 14:3x 编排者本程现跑）

| # | 断言 | 尺 | 读数 |
|---|---|---|---|
| 1 | 思考档位的**枚举与默认值都在** | `grep -n "thinking_intensity\|thinking_levels" internal/config/schema.go` | `:291 default:"off"`、`:359 ThinkingLevels []string`（枚举逐枚校验在 `internal/config/validate.go:201`） |
| 2 | 它今天**真的被用**（不是死字段） | `grep -rn "ThinkingIntensity" --include=*.go internal/ \| grep -v _test.go` | `internal/agent/prompt.go:113-114`＋`:225` 一路带进出模型的请求 |
| 3 | 这一轮用哪个模型是**按角色**定的 | `grep -n "Role is one llm.roles" -A8 internal/config/schema.go` | `:280-287`＝`provider`＋`model` 成对引用目录里的模型 |
| 4 | 面板今天**没有任何写入口**改这两样 | `grep -n '"panel\.' internal/panel/bridge.go` | 恰好四枚（`:42-45`），里面没有 model／thinking 相关方法 |
| 5 | 配置有三档生效级别（热改／重载／重启） | `grep -n "Hot-tier\|reload" internal/config/schema.go \| head` | `:38`、`:206`、`:242` 等——**"面板改档位算哪一档"是本票第一问** |

⇒ **缺的不是数据源，是"从界面写回去"那一条腿**（与票 145 AC#2b 的判断一致，只是本票要的不止显示）。

## AC（每格都要答"这一发在**未修码**上响不响"；先测→再写→再提交）

- [ ] **AC#1 先答"改了之后什么时候生效"**：把 D36 那三档生效级别对这两枚字段逐一定档——① 换**思考档位**（`thinking_intensity`）是热改、重载还是重启？② 换**模型／角色**呢（它牵 `C5 LlmProvider` 与上下文预算 `D15`）？⚠ 每支给一条现跑：找到那枚字段的**读者**（`grep -rn "ThinkingIntensity\|roles\." --include=*.go internal/agent/ internal/config/ cmd/ | grep -v _test.go`），看它是**每次调用现读**还是**启动时读一次**——**"每次现读"才可能当场生效**，不许按"应该可以"写。
- [ ] **AC#2 写入口走哪条通道，并把它写成契约变更入档**：候选＝复用 `panel.mode.request` 那条"请求→处理器→写者"的现成形状（`internal/panel/composer_handlers.go:111 HandleModeRequest` 是**同类先例**：它已经在做"面板发起、宿主判定、只许调严／调松要确认、每笔写审计"）／新起 `panel.model.request`＋`panel.thinking.request`。**判据＝改动面最小**；若结论是要新增 C17 白名单条目，⚠ 先具名报我、我落 `A##` 再写（C17 定稿是契约面，owner 那句"别老说契约冻结"是**方向授权、不逐条点名**）。
- [ ] **AC#3 只许"变严不ask、变松要显式"这条口径要现量核对**（R20/M4 的形状）：换到**更贵／更强**的模型、把思考档位从 `off` 抬到 `high`，算不算"放宽"？⇒ **给一枚判据钉住"面板改这两样绝不可能顺带放宽权限档位"**（这两枚字段与 `risk.Mode` 必须互不相干；尺：证明写腿里没有任何一处碰 `Modes.Set`）。⚠ 若你判"必须碰权限腿才能生效"＝**停手上报**，那是 `Q-49` 那一族刚钉过的门。
- [ ] **AC#4 清单要真、不许造**：面板下拉里能选的那几枚模型，**必须**来自目录与角色表里真存在的项（`llm.providers.<p>.models.<m>`），思考档位**必须**只列那枚模型**声明支持**的档（`ThinkingLevels`，可能少于四档）。⇒ 判据：给一枚模型只声明 `off|low`，面板不许出现 `medium|high`；空目录／没配 provider 那一形要**显式报"没有可用模型"**（宁缺毋造，票 92/145 同形）。
- [ ] **AC#5 生效范围要具名**：这一改是**只作用于当前这一轮／这个会话／整个进程**？逐支答，并给一条判据钉住"改完之后界面显示的那一枚＝运行时真用的那一枚"（票 146 那族"值保真"的教训：显示与真用分账会静默漂）。
- [ ] **AC#6 覆盖面：谁还会读到旧值**——语音链路（`voice.*` 是 reload-tier）、`roles` 的其它使用者、上下文预算 `D15` 的 `BudgetsFor(<窗口>)`（换模型＝换窗口！尺 `grep -rn "BudgetsFor" --include=*.go internal/ cmd/ | grep -v _test.go`）。⚠ 换模型把上下文窗口换小会让**正在跑的那一轮**超限——这一支必须给结论（拒绝切换／等到轮次边界／还是照切并如实报）。
- [ ] **AC#7 契约轴**：`docs/PLAN.md`（D36 那一整节＋D8/D15 近邻）、`docs/specs/**`、`thresholds.go`／golden／审批超时常量／`allowlist.txt` 一字节不许动；**凭据面**：任何 API 密钥的值不许进快照、面板与读数文件（只许出现 `api_key_ref` 这类变量名）。
- [ ] **AC#8 门禁**：逐包 `go test -count=1 ./internal/panel/ ./internal/config/ ./internal/agent/`；CLI 那一面走 `-overlay` **且带在册 PATH 前缀**（不带＝`0xc0000135`，票 98）；`sh scripts/d22scan.sh`（基线 `ban #8 internal/` **examined=433**）；`gate-clauses.sh` **比红腿名册**（在册只有 `G6neg`）；⚠ 终态读数取在你自己最后一枚 commit 之后；⚠ **不许跑 `probes/161/r6/flip-declaration.sh`**；⚠ 重跑既有台件前先看它的写出路径（`A367` 那枚坑）。
- [ ] **AC#9 界面那一半不在本票写面**：`frontend/**`／`design/**` 归另一枚会话。⇒ 本票只交 Go 侧的读面＋写腿＋快照字段，另出一张"面板要画哪几枚控件、点下去发哪条请求"的需求表（`docs/evidence/s1/187-panel-ask.md`）交编排者转前端。

## 本票**不**解决

- 不做模型下载／删除／凭据录入（那是 `wisp models` 与票 26 那一族；凭据只能由 owner 在界面自助录入，我不代填）。
- 不做语音侧模型切换（`voice.*` 是另一条 reload 路径）。
- 不碰权限档位三档的语义（R20 已定，本票只保证"不渗进去"）。

## Progress log

- 09-28 14:3x 编排者立票：owner 改判 `Q-64`＝做（`A368`）。上面 5 把尺本程现跑（锚 `72c76d42`）。⚠ 与票 186 **分程派**（186 是 git 那一块、本票是模型与档位，两枚写面不同物，但都会碰 `internal/panel/bridge.go` 的白名单那一段 ⇒ **同批只派一枚碰 `bridge.go`**，另一枚按住）。未派。
