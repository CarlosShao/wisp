# 193 — 动画档 Go 侧落地腿 r1（字段 + 生效链 + 可证关闭档 + 门禁清单）

- 票：`.scratch/wisp/issues/193-lightweight-must-no-longer-constrain-frontend-animation-and-visuals-make-the-per-panel-gates-named.md`
- 派单：`.scratch/wisp/probes/orch/20261119/txt/193-animation-opt-dispatch.txt`
- 写面：`internal/config/**`、`internal/panel/**`、本件。⛔ 未碰：`cmd/wisp/**`、`tokens_fourway_test.go`、`docs/PLAN.md`、`thresholds.go`/golden/SLO 阈值、`frontend/**` 产码、他人在飞件（`frontend/fixtures/composer-states.html` 等 M 只报不动）。
- 时刻：2026-11-19 19:0x–19:4x +08。

---

## ① AC#0（只报不动）：那枚 owner 按过的测试，今天的读数与"红因"判定

**结论先说：今天它不红。派单预判的红因（`design/assets/tokens.css` 缺失/被删未 staged）已过时——该文件的删除早已随 commit 入库，测试本体已被票 380 重建为第四代四向、不再读该文件。**

四把尺（可逐条重跑）：

| 尺 | 今天读数 |
|---|---|
| `go test ./internal/panel -run "^TestTokensFourWay"`（派单原命令） | `ok github.com/CarlosShao/wisp/internal/panel 0.091s [no tests to run]`，rc=0 —— 该名字匹配零个测试 |
| 同包真测试名：`go test ./internal/panel -run "^TestC21DesignTokensFourWayAgree" -count=1` | `ok github.com/CarlosShao/wisp/internal/panel 0.065s`，rc=0（测试在 `tokens_fourway_test.go:631`） |
| `git status --porcelain design/` | 只有 untracked 新增（`design/doubao/**`、`design/old/`），**没有** ` D design/assets/tokens.css` 行 |
| `git log --oneline -3 -- design/assets/tokens.css` | `72b65c614`（删除发生在这枚 commit 里，`--diff-filter=D` 现跑可复）、`5d777f78a`（initial import）；`Test-Path design/assets/tokens.css` = False |

**为什么"删除已入库但测试还绿"**：`tokens_fourway_test.go:13-18`（原文逐字）："P1 used to be design/assets/tokens.css, compared against P4. That file was deleted in 72b65c614; only a copy survives as the ball package's own CSS fixture … The fourth generation's reference implementation is the TOKENS table inside the generator"。即票 380（`0a8920a15`）已把四向重指到 `frontend/scripts/gen-tokens.mjs`／`docs/evidence/s1/c21-native-tokens.md`／`internal/ball/tokens.go`／`frontend/src/styles/tokens.generated.css`。

**本腿动作：零。** 未 restore、未编辑 `tokens_fourway_test.go` 一字节；要不要"动 tokens 那枚"仍只有 owner 一句能启动，本腿不拿它当绿灯也不当拒绝前端的理由。

---

## ② AC#2 门禁清单（Go 侧现量普查 + 放开/保留逐枚判定）

**尺（派单原命令，git grep 等价实现）**：`git grep -nE "轻量|lightweight|D1\b|D32" -- "internal/*.go" "frontend/*.go" "tools/*.go" "scripts/*.go"`，剔除 `_test.go`。

**总读数：非测试 Go 源命中 115 行。其中 `轻量`／`lightweight`＝0 行；`frontend/` 下 Go 文件命中＝0 行（frontend 无 .go）。`D1\b`＝11 行，`D32`＝105 行（两者交集 1 行＝本腿新注释 schema.go:606）。**

判定口径（票面 AC#2 原文）：**放开＝把 D1/D32 的数当动画/过渡/视觉约束的条款引用**；保留＝零 emoji、ban #6 面板侧不得成批准/执行入口，以及一切非视觉射程的 D32 纪律。

### 逐枚判定

**(a) 放开候选：0 枚。** 115 行里没有任何一行拿 D1/D32 的数去约束前端动画/过渡/视觉复杂度。票 193-a1 普查找到的三句"拿 D32 当理由禁动画"在 `docs/PLAN.md`（`:3451`/`:3486`/`:3516`），那是文档不是 Go 产码，且其"D32 半句"已由台账 **A378**（2026-09-28，编排者落）裁定"从此不算理由"；本腿不改 PLAN.md 一字（AC#6），Go 树里没有可放开的引用。**A377/A378 的口径在本腿 Go 侧现量下成立：轻量级从未在 Go 产码里卡过前端动画。**

**(b) 误报（正则命中、非决策引用）：8 行，不计入放开/保留**——`internal/audio/aecspike/rig_windows.go:72`、`internal/audio/wasapi_windows.go:72`（GUID 字节 `0xD1`）；`internal/ball/d2d_windows.go:35,38`、`scripts/spike/common/winshell.go:27,53`（`d2d1.dll`/`D2D1CreateFactory`，Direct2D API 名）；`internal/ball/tokens.go:266,432`（色值 `0x246BD1`）。

**(c) 真 D1 引用：2 行，保留**——`internal/models/ads.go:13`、`internal/models/ads_windows.go:13`（D1 在这里是"模型对抗样本 ADS 的 NTFS 实测上下文"，与前端视觉无关）。

**(d) 真 D32 引用：104 行（存量）+ 1 行（本腿新注释），全保留。** 逐包计数（git grep 现跑可复）：

| 包 | 行数 | 内容 | 保留理由 |
|---|---|---|---|
| internal/watchdog/* | 20 | D32 行限/沉没规则/白名单 | 资源门（D32 的数，AC#6 ⛔） |
| internal/observe/* | 23 | SLO 采样/阈值/冻结行 | 资源门（`thresholds.go` 7 行 ⛔ 一字节不动） |
| internal/proc/* | 12 | 树内存/C30 Job | 资源门 |
| internal/ball/* | 14 | 静帧/零定时器 CPU 纪律 | **原生球**渲染纪律，射程非 `frontend/`；球侧动画可否可选＝a1 普查登记的 needs-approval 项（`:1899-1901`），**本票不解决，停手上报维持** |
| scripts/spike/* | 12 | 实测脚本引预算 | 测量记录 |
| internal/speech/* + statemachine | 10 | 语音链路/状态机 | 轻量级继续约束语音链路（票面口径原句） |
| internal/agent/scheduler/* | 4 | 无定时器/可观测排空 | 运行时 CPU 纪律 |
| internal/risk/provenance.go | 4 | 索引内存预算 | 资源门 |
| internal/session、memory、plugin、tools、config（schema.go:670 SLO 采样节拍） | 17 | 杂项纪律 | 非视觉射程 |

**(e) 零 emoji（ban #8）与批准面（ban #6）：保留，未动。** 本腿产码零 emoji（d22scan ban #8 现跑无 emoji finding，见 §⑥）；ban #6（`frontend/` 面板批准面）一条未碰——安全边界，前端功能不依赖它（票面原句）。`tools/d22scan/allowlist.txt` 未动。

### 拟落账 A## 条目（供编排者追加 `docs/reports/pending-and-issues.md`，本腿写面不含台账文件）

- **拟 A-193-a（口径登记）**：文件/行＝Go 产码全域（尺见 §②）；判定＝"轻量级卡前端动画"在 Go 侧**零引用可放开**（115 行现量，逐包表见本件 §②(d)）；边界＝①零 emoji 与 ban #6 保留不动；②`internal/ball` 的 D32 静帧纪律保留（原生球，needs-approval 另议）；③`thresholds.go`/golden/SLO 数一字节未动；④本件只登记引用口径，`PLAN.md:3451/3486/3516` 三句的正式改写仍需人工批准（AC#6）。**撤销口令：恢复"轻量级约束前端动画"的引用口径（A378 作废）。**

---

## ③ AC#3 字段与生效链

### 字段

`internal/config/schema.go` 的 `PanelSection` 新增（与 `font_size` 同节）：

```go
Animation string `toml:"animation" default:"reduce"`
```

- 三档词表：`full`｜`reduce`｜`off`（全开／减少／关闭），单一来源＝`validate.go` 的 `panelAnimationVocabulary`；未知值＝**载入错误**（`validatePanelAnimation`，不走回退——同 `risk.permission_mode` 形状），写门与载入门共用同一句拒绝（`panelAnimationEnumError`，65/413 单句形状）。
- **默认值＝待 owner 拍（⛔ 本件与 schema 注释均明写，不许当已定）**：票面原句"默认值由 owner 定"；派单令按最保守＝"减少"（reduce）占位。owner 拍板后：schema tag、`TestPanelAnimationDefaultIsTheOwnerPendingPlaceholder` 与本件此行同 commit 移动即可。

### 定档：**hot**，理由四条

1. 机器答案在册：`TierRegistry["panel"] = TierHot`（`internal/config/tiers.go:44`），[panel] 整节 hot，`panel.animation` 随整节走，**无需新增 per-key 行**（per-key 行是 [app]/[voice] 那种"一节两档"的形状，panel 不需要）。
2. D36 出处：[panel] 行钉 hot（票 193-a1 普查引 `PLAN.md:2743`；本腿不改 PLAN.md）。
3. 消费者是面板页：档位词随快照每次泵送（reader-gated 载体），重启档会把开关变成装饰品（票 180 同一枚病）——hot 是唯一不自欺的档。
4. Go 侧无 applier 是诚实形状：同 `panel.width`/`panel.font_size`（`config_write_wave.go` 的 tier 注释原文）——配置值落地即可被读者热读，页面消费归前端腿。

### 现量读者数（两把尺，覆盖派单的 grep）

| 尺 | 今天读数 |
|---|---|
| 派单尺 `PanelSection\.`（internal+cmd，非测试） | **1 行**：`cmd/wisp/panel_host_windows.go:482`（注释引用） |
| 强尺：`[panel]` 字段的生产读者（`.Panel.` 字段访问 + `Panel*Options` + 泵接线） | **6 枚字段里 5 枚有生产读者**：`Enabled/Width/Height/KeepAliveInSession`（`cmd/wisp/panel_resident_windows.go:231,258,276`）、`FontSize`（`internal/config/settings_appearance_wave.go:101,188` 写+读 + `cmd/wisp/run.go:1359-1363` 泵接线 + `panel_resident_pump.go:231-235` 常驻泵）；**`Scale`＝0 读者**；收据行 `config_readers_255.go:310` |

⚠ 票面 a1 普查的"`[panel]` 整节零读者"是 09-28 的读数，已被 255 接线与外观波超越——本件以今天现量为准（5/6 有读者，唯 Scale 仍是零读者存量项）。

### 本腿落的生效链

```
config.toml [panel].animation
  → loader/validate（未知词=载入错误）           internal/config/validate.go      （新）
  → Manager 热读（写后同进程可读）               settings_animation_193.go         （新）
     SetPanelAnimation（guarded writeOneKey，拒绝词=单句）                        （新）
     PanelAnimationOptions()（nil-安全读面，带词表）                             （新）
  → 快照载体（reader-gated，无读者=零键）        panel_animation_193.go（新）＋
     Snapshot.panelAnimation + PumpSources.PanelAnimation + 泵发射块    composer.go/pump.go
  → 对账（panel.ts 未声明 ⇒ 台账登记判定）       composerTSGapLedger["PanelSnapshot"]=["panelAnimation"]  （新）
  → 页面消费（panel.ts 声明 + 动画真停）         ⛔ frontend/**，前端腿清单见 §④
  → cmd/wisp 接线（run.go 泵源一行 + 常驻泵一行 + ApplySetting case）
                                                 ⛔ cmd/wisp/**，接线腿清单见 §⑤
```

**诚实边界（本腿禁区导致的两根未接线，均已具名，不许当已接）**：
1. **cmd/wisp 泵源接线未落**（`run.go:1359-1363` 旁 + `panel_resident_pump.go:231-235` 旁各一行）⇒ 生产快照今天**不发** `panelAnimation` 键——reader-gated 缺席，四键字节钉保持绿，绝不伪造档位词（`TestPanelAnimationCarrierWithNoReaderSendsNoKeyAndNeverInventsOff` 钉死）。
2. **panel.config.set 写路由未落**：第九名册（`panel_animation` 字段）**刻意不进** `knownWritableField`——因为 A881 并集钉在 `cmd/wisp/panel_config_313_test.go:225`（新名册必须在那把尺里追加一行），而 cmd/wisp 是本腿禁区；进了名册而 store 无 case＝"门 admitting、腿拒绝"的假路由。config 侧 `SetPanelAnimation` 已就位（接缝单源），cmd/wisp case＋名册＋并集钉一行是接线腿三步。

---

## ④ AC#4 关闭档可证（Go 半＝可重跑的尺；前端半＝清单）

### Go 半：三枚测试（`internal/panel/panel_animation_193_test.go`，可重跑）

| 尺 | 钉什么 |
|---|---|
| `TestPanelAnimationCarrierRidesTheReaderVerbatim` | off/reduce/full 三档逐字上 wire：`{"panelAnimation":{"tier":"off"}}`——泵不派生、不改写、不发明 |
| `TestPanelAnimationCarrierWithNoReaderSendsNoKeyAndNeverInventsOff` | 无读者 ⇒ 键缺席，且全包**不含任何 tier 词**——"未设置"永不冒充"off"（把"设了个 CSS 变量"式假证从 Go 侧就掐死） |
| `TestPanelAnimationCarrierSendsTheUnsetTierItRead` | nil-Manager 读数 `""` 逐字上 wire＝页面的"未设"，不是"off" |

配对（internal/config）：`TestPanelAnimationLoadsOffAndReadsHot`（off 落盘⇒热读 off）、`TestPanelAnimationLoadRefusesUnknownTierByTheSingleSentence`（loud/OFF/none/"" 全拒）、`TestSetPanelAnimationWritesAndRefusesByTheSingleSentence`（写入收据＝`panel.animation`，拒绝零字节移动）。Go 侧为该键**新增定时器/协程＝0**（载体一个结构体一根发射线，无 tick）。

**"真停 vs 还在跑只是看不见"的判定权在哪**：Go 侧只能保证"档位词逐字到达 wire、且无任何 Go 侧动画驱动"；页面是否真停（rAF 不注册、transition 不挂）只能在渲染层证——那是前端腿的尺，清单如下。

### 前端半：交接给界面腿的清单（本腿零行 frontend/** 产码）

1. **off 档的"真停"判据（可重跑尺形状，接 a1 普查 §⑤ 的两发）**：
   - 运行时观测：off 档下页面 `document.getAnimations()` 长度＝0、活跃 `requestAnimationFrame` 句柄＝0、无 `setInterval` 驱动的动画循环；对照 full 档 >0——差分为证，不是"看不见"。
   - 产物差分：off 档构建产物里动画类/transition 声明不挂载（元素上**没有** `transition`/`animation` 属性，而非 opacity:0 的"还在跑只是看不见"）。
2. **off ≠ prefers-reduced-motion**：媒体查询不是这枚开关的实现路径——用户选 off 必须无条件真停（系统"减少动态效果"开着关着都一样）。
3. **reduce 档**：过渡允许、无限循环动画不允许（A378 边界：`PLAN.md:3497/:3499` 的生命周期句不在放开范围——非空闲态才可有循环动画、面板关闭后不得有动画在跑，这与 owner 要的"最低档关掉就不卡"同向）。
4. **panel.ts 声明**：UI 腿声明 `panelAnimation: PanelAnimationSection`（`{tier: string}`）后，`composerTSGapLedger` 的 `"PanelSnapshot": {"panelAnimation"}` 登记会**故意转红（stale registration）**——届时删该登记、在 `composerTSNameMap` 加 `PanelAnimationSection` 翻译，走普通双向对账（外观波四键的同一条收尾路，4017962b3）。
5. **档位词来源唯一**：页面只认 wire 的 `panelAnimation.tier` 三词（`full/reduce/off`），`""`＝未设置（不得当 off 用）；禁止页面自造第四个词或从 CSS 变量反推档位。

---

## ⑤ 交接给接线腿（cmd/wisp，三步，本腿⛔禁区）

1. 泵源两行：`run.go` 泵装配处（`PanelWidth`/`PanelFontSize` 旁）加 `PanelAnimation: func() panel.PanelAnimationSection { return panel.PanelAnimationSection{Tier: mgr.PanelAnimationOptions().Tier} }`；`panel_resident_pump.go:231-235` 旁同形。
2. store 一 case：`panel_config_store.go` 的 `ApplySetting` 加 `case panel.FieldPanelAnimation:` → `s.mgr.SetPanelAnimation(v)`（tier 句照 config_write_wave.go 的 [panel] 两行的 Note 形状答 hot）。
3. 第九名册：新 roster 文件 `internal/panel/config_write_193.go`（`FieldPanelAnimation = "panel_animation"`）＋ `knownWritableField` 并集加一件＋ `cmd/wisp/panel_config_313_test.go:249-263` 的并集钉追加 `panel.WritableFields193()` 一行。**顺序注意**：名册与 store case 必须 same commit，否则门 admitting、腿拒绝。

---

## ⑥ AC#6 契约轴 + AC#7 门禁读数

- **契约轴**：`git diff --name-only HEAD -- docs/PLAN.md`＝空（一字未改）；`internal/observe/thresholds.go`、golden、SLO 阈值零字节未动；D36 三档语义沿用 TierRegistry 现量，未新造档位词。
- **d22scan**：`sh scripts/d22scan.sh` 末行（证据件落盘后复跑）：`d22scan: clean - no D22 ban violations; … ban #8 internal/＝1070, cmd/＝359, frontend/＝188, design/＝39`，scan rc=0，**0 findings**；positive control（`tools/d22scan ./...`）rc=0。
  - 过程记录（诚实记账）：第一次跑 1 finding——`internal/config/schema.go:621 [phantom-citation]` 引用了尚未落盘的本证据件路径（ticket 212 ban #9），写完本件即愈；非放宽、非 allowlist。
  - `scripts/gate-clauses.sh`：**不存在**（`Test-Path`＝False），派单条件"若存在给读数"不触发。
- **两包四数**（`go test -count=1 ./internal/config/ ./internal/panel/`，最终复跑）：config **ok**（PASS=全绿，含 4 枚新 PanelAnimation 测试）／panel **ok**（全绿，含 3 枚新载体测试＋契约对账 17 键/19 缺口登记）——精确秒数与四数以回话首行为准。
  - 基线（改动前）：config ok 4.252s、panel ok 32.120s——票面 AC#7 提到的"`internal/panel` 2 枚已知红"（09-28 a1 时点）今天**不存在**，两包改前已全绿。
- **gofmt**：本腿触碰文件 gofmt 干净；`gofmt -l internal/panel` 现列 `config_handlers.go`、`result_golden_200_36_test.go` 两件为**存量不齐**（HEAD 逐字验证，非本腿引入，未触碰）。

---

## ⑦ 拟 Progress log（供编排者追加票面，本腿写面不含 .scratch）

- 2026-11-19 19:4x `193-leg`（agent-156）交件：AC#0 判定"红因过时、今天绿、零动作"；AC#2 Go 侧放开 0/保留 107（逐包表）；AC#3 `panel.animation` 落 schema（默认 reduce＝owner 待拍占位，档 hot 四条理由，读者现量 5/6）；AC#4 Go 三尺＋前端清单；AC#6 契约零触碰；AC#7 d22scan 0 findings＋两包全绿。**未勾格建议**：AC#2 待编排者把拟 A-193-a 落台账后翻勾；AC#3/#4 的"生效链/可证"还差 cmd/wisp 接线与前端半（§④/§⑤ 清单），接线腿落地前不翻。未动：`cmd/wisp/**`、`frontend/**`、`tokens_fourway_test.go`、PLAN.md、thresholds。
