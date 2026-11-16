# 设置页「未接线」精确清单（只读普查 · 2026-11-16 20:1x +08）

> 审计腿 · 只读 · ⛔ 未改产码 / 未改票面 / 未 push
> 尺：`frontend/src/components/config-screen.tsx` 渲染的每一行；判据＝`live` 从哪来。
> 机主问题：「为什么设置页上一堆未接线」。

---

## 0. 头条结论（先说最硬的一条）

**`live` 根本不是 Go 侧的 wire 字段。** 它是前端**静态字面量**：

`frontend/src/components/config-screen.tsx:411`
```tsx
const live = row.blocked === "";
```

`blocked` 是 `APPEARANCE_ROWS` 表里的**硬编码字符串**（`config-screen.tsx:152-162`）。
⇒ **哪一行显示「未接线」，在编译期就定死了，与 Go 侧有没有读者、快照有没有推、有没有引擎，一点关系都没有。**

⇒ 用户在设置页看到的「一堆未接线」，**不是后端没接线的读数，是前端一张手写的分类表。**
⇒ 后端今天接没接，只能**逐行另查 Go** 才能知道（下文 §2/§3/§4）。

**判据出处（题目点名要读的三个文件，实际结论）**：

| 题目问的 | 实际在哪 | 结论 |
|---|---|---|
| `frontend/src/lib/config-view.ts` 里 `live` | **不在**。该文件 256 行只有 `config.get` 回执的解析器（`ConfigView` / `parseConfigReply` / `useConfigView`），零个 `live` 字面量 | 与设置页外观节**无关** |
| `internal/panel/config_handlers.go` 里焊 `live`/`immediate`/`tier` 的字段 | **没有 `live` 字段**。该文件有 `EffectiveTier`（`:328-335`：`now`/`next_task`/`restart`/`not_applied`）与 `tierSentence()`（`:603-616`） | tier 是**收据**上的字（写完才说的），**不是**徽章的来源 |
| 前端徽章的 `live` | `config-screen.tsx:411` | 静态字面量 |

`config_handlers.go` 的 tier 确实被前端用到了——但走的是**另一条路**：
`config-screen.tsx:572-585` 的 `EffectBadge` 的 `word` 参数是**硬编码**的
（`ballStyle`/`ballSize`→"能改"、`clickThrough`→"立即/重启"、其余三枚→"关窗/重启生效"），
**没有读 Go 的 tier**。真正的 Go tier 只出现在**写完之后的收据行**里
（`config-actions.ts` → `config-feedback.ts:155-166` `translateTieredWriteReply`）。

---

## 1【清点】设置页外观节的每一行

源：`frontend/src/components/config-screen.tsx:152-162`（`APPEARANCE_ROWS` 表），渲染在 `:410-589`。

| # | 行名（wire 字段） | desc | `blocked` | `control` | 徽章 `word` | 渲染成 |
|---|---|---|---|---|---|---|
| 1 | `app.theme` | 界面主题 | **"主题由原生宿主决定，面板只跟随"**（`:153`） | — | — | **未接线** |
| 2 | `panel.opacity` | 面板背景不透明度 | `""`（`:154`） | `alpha` | 无⇒**即时** | **即时** |
| 3 | `panel.animations` | 界面动画（关掉后全部静帧，性能最省） | `""`（`:155`） | `motion` | 无⇒**即时** | **即时** |
| 4 | `panel.font_size` | 面板正文字号（px） | `""`（`:156`） | `panelFontSize` | "关窗/重启生效" | 关窗/重启生效 |
| 5 | `ball.size` | 悬浮球直径（44-160） | `""`（`:157`） | `ballSize` | "能改" | 能改 |
| 6 | `ball.opacity_idle` | 空闲时悬浮球不透明度 | `""`（`:158`） | `ballOpacity` | "关窗/重启生效" | 关窗/重启生效 |
| 7 | `ball.click_through` | 空闲时鼠标穿透悬浮球 | `""`（`:159`） | `clickThrough` | "立即/重启" | 立即/重启 |
| 8 | `ball.style` | 悬浮球外观风格（玻璃 / Jarvis 粒子） | `""`（`:160`） | `ballStyle` | "能改" | 能改 |
| 9 | `panel.width` | 面板宽度（px） | `""`（`:161`） | `panelWidth` | "关窗/重启生效" | 关窗/重启生效 |

### 总数

- **共 9 行**
- **即时 2 行**（`panel.opacity` `panel.animations`）＝徽章字面就是"即时"（`EffectBadge` 的 `word` 缺省值，`:178`）
- **未接线 1 行**（`app.theme`）＝唯一 `blocked` 非空的那行

### ★ 但这个「1 行未接线」是**假读数**

| 口径 | 数 |
|---|---|
| 页面**徽章**上写「未接线」的 | **1 行**（`app.theme`） |
| 页面**徽章**上写「即时」的 | **2 行** |
| 页面**徽章**上写其它 tier 词的（能改/立即-重启/关窗-重启生效） | **6 行** |

⇒ 用户在设置页**看不到 6 行「未接线」**，看到的是 1 行。
⇒ 那「一堆未接线」**不是来自 `config-screen.tsx` 这一页**。见 §5。

---

## 2【归因】每行的「接没接」——现读 Go，读者在哪

判据口径（**不是看注释**）：`grep` 该键的**读者**——谁在**运行期**读它并用它做事。

| 行 | 写侧（config.set） | 读侧（快照读者） | 生产装配点 | 真正生效的引擎 | 判定 |
|---|---|---|---|---|---|
| `app.theme` | ⛔ **不在可写名册** | ⛔ **零读者** | — | — | **真未接线** |
| `panel.opacity` | N/A（页面本地） | N/A | — | CSS 变量 `--panel-alpha`（`frontend/src/styles/theme.css:37-38`） | **真即时**（页面本地，不落盘） |
| `panel.animations` | N/A（页面本地） | N/A | — | `html[data-motion]`（`frontend/src/lib/motion.ts:45`） | **真即时**（页面本地，不落盘） |
| `panel.font_size` | ✅ `internal/panel/config_write_wave.go:40` `FieldPanelFontSize` | ✅ `internal/config/settings_appearance_wave.go:179` `PanelFontSizeOptions` | ✅ `cmd/wisp/run.go:1272-1274`；`cmd/wisp/panel_resident_pump.go:141` | ⛔ **零消费引擎**（见下） | **写通／无引擎** |
| `ball.size` | ✅ `internal/panel/config_write_65.go`（`FieldBallSize`） | ✅ `internal/config` `BallSizeOptions` | ✅ `cmd/wisp/run.go:1249-1252` | ✅ `cmd/wisp/ball_size_413_windows.go:40` `ballSizeSeed413` → `cmd/wisp/resident_ball_windows.go:488` `SizePx:` | **已接**（重启档） |
| `ball.opacity_idle` | ✅ `config_write_wave.go:37` | ✅ `settings_appearance_wave.go:128` | ✅ `cmd/wisp/run.go:1262-1265` | ✅ `cmd/wisp/ball_idle_appearance_windows.go:56` → `resident_ball_windows.go:497` `OpacityIdle:` → `internal/ball/idle_appearance_windows.go:48-51` | **已接**（即时档，有球时） |
| `ball.click_through` | ✅ `config_write_wave.go:38` | ✅ `settings_appearance_wave.go:149` | ✅ `cmd/wisp/run.go:1266-1268` | ✅ `ball_idle_appearance_windows.go:56` → `resident_ball_windows.go:498` → `internal/ball/idle_appearance_windows.go:64` | **已接**（即时档，有球时） |
| `ball.style` | ✅ `internal/panel/config_write_65.go` | ✅ `internal/config` `BallStyleOptions` | ✅ `cmd/wisp/run.go:1236-1239` | ✅ `cmd/wisp/ball_style_65_windows.go:42` `ballStyleSeed65` → `resident_windows.go:444` `ball.SetStyle` | **已接**（即时档，有球时） |
| `panel.width` | ✅ `config_write_wave.go:39` | ✅ `internal/config` `PanelWidthOptions` | ✅ `cmd/wisp/run.go:1269-1271` | ✅ `cmd/wisp/panel_resident_windows.go:231` `cfg.Panel.Width` → `cmd/wisp/panel_host_windows.go` 建窗 | **已接**（关窗再开） |

### 2.1 `app.theme` —— 唯一真正「未接线」的一枚，且是最硬的那种（**零读者**）

- **读者数＝0（生产）**。`grep -rn "App\.Theme" --include=*.go` 全仓命中 **22 处**，
  逐条核过：**全部在 `internal/config` 包内自转圈**或测试里：
  - `internal/config/manager.go:378` `themeChanged := cur.App.Theme != fresh.App.Theme`（比较）
  - `internal/config/manager.go:383` `plan.set = append(plan.set, func() { cur.App.Theme = fresh.App.Theme })`（拷贝）
  - `internal/config/validate.go:51` `slices.Contains(themes, c.App.Theme)`（枚举校验）
  - `internal/config/tiers.go:59` `"app.theme": TierHot`（档位名册）
  - 其余全是 `*_test.go`
  - ⇒ **没有任何一个组件读它并据此改变任何行为。**
- **生产读者名册自己就这么写**（**不是注释当结论，是名册里的判定**）：
  `cmd/wisp/config_readers_255.go:323`
  ```
  "app.theme": hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.App.Theme ... so this row stays no-reader"
  ```
- **写侧也没有**：`app_theme` **不在** `acceptedWritableFields()`（`internal/panel/config_write_313.go:201-209`，
  ＝六册并集：248 七枚 ＋ 313 十三枚 ＋ 187 一枚 ＋ 65 两枚 ＋ wave 四枚 ＋ 215a 一枚）。
  ⇒ 页面即使想写 `app.theme`，`config.set` 会被 `config_handlers.go:463-468` 按名拒掉。
- **球侧缝在盘、但没人踩**：`internal/ball/renderer_windows.go:106` `func SetTheme(t Theme)`
  ——全仓**只有定义、零调用者**（`grep "SetTheme("` 只命中定义处与一条测试注释
  `internal/ball/liquid_test.go:281`）。`ball.Options.Theme`（`internal/ball/ball_windows.go:93`）
  字段在，**装配根不填**（`cmd/wisp/resident_ball_windows.go:488-500` 的 Options 字面量**没有 Theme 键**）。
- **面板侧主题是纯前端本地态**：`frontend/src/App.tsx:191-195` `toggleTheme()` 改的是
  `document.documentElement.dataset.theme`，**不读 config.toml、不落盘**。
  ⇒ 即便把 `app.theme` 接上，今天的面板主题**也不跟它走**。

### 2.2 `panel.font_size` —— 写通了、读者有、**但零消费引擎**（"半接线"）

- 写侧 ✅：`internal/panel/config_write_wave.go:40` 在名册；`cmd/wisp/panel_config_store.go:431-440` 有真 case 调 `s.mgr.SetPanelFontSize(v)`。
- 读侧 ✅：`internal/config/settings_appearance_wave.go:179-188` `PanelFontSizeOptions()`；
  快照读者 ✅ `cmd/wisp/run.go:1272-1274` ＋ `cmd/wisp/panel_resident_pump.go:141`。
- **引擎 ⛔ 零**：`grep -rn "Panel\.FontSize"` 全仓 **6 处命中**，全是
  `internal/config/settings_appearance_wave.go:101`（写）、`:188`（读）、**两枚测试**。
  ⇒ **没有任何一行把 `panel.font_size` 接到页面正文字号上**：
  - `cmd/wisp/panel_host_windows.go` 的 WebView2 建窗只带 `Width`/`Height`（`:380` 附近），
    **没有 FontSize**；
  - 前端 `panelFontSize.fontSize` 只被 `config-screen.tsx:557/564` 用来**填输入框**，
    **没有 `setProperty("--font-size", …)` 一类的消费**；
    `grep "fontSize" frontend/src` 命中 24 处，**没有一处是 `--panel-…` 字号变量**。
- ⇒ **写进去不生效**。这正是 `internal/config/schema.go:594-602` 自己写的：
  「the consumer (the panel page's font rendering) lands with the frontend wiring leg」——
  **那一腿没落地**。

### 2.3 `panel.opacity` / `panel.animations` —— 真即时，但**不是 config 键**

- `panel.opacity`：`config.toml` 的 `[panel]` 段**没有 opacity 键**
  （`internal/config/schema.go:583-604` 的 `PanelSection`：`enabled`/`width`/`height`/`keep_alive_in_session`/`scale`/`font_size`）。
  ⇒ 它读的是**页面自己的 CSS 变量** `--panel-alpha`（`config-screen.tsx:113` `ALPHA_VAR`），
  由 `frontend/src/styles/theme.css:37-38` 的 `color-mix(in srgb, var(--page) var(--panel-alpha), transparent)` 消费。
  ⇒ **真即时、不落盘**（关掉面板回到 token 默认 `100%`，`frontend/src/styles/tokens.generated.css:66`）。
- `panel.animations`：同理，**`[panel]` 段没有 animations 键**；消费面是 `html[data-motion]`
  （`frontend/src/lib/motion.ts:45` `MOTION_ATTR`）。
  ⇒ 两枚徽章写「即时」是**实话**，但它们**不是 config 键**，别被当成已接线的 config 行。

---

## 3【归票】每行归到哪张票

票池：`.scratch/wisp/issues/*.md`。判据：**票面正文点名该键/该功能**。

| 行 | 票号 | 票状态 | 证据 |
|---|---|---|---|
| `app.theme` | **票 39**（GUI 配置编辑器） | **open**（`39-config-editor-gui.md`，无 `-done` 后缀，`Status: ready-for-agent`） | `180-...-done.md:288`：「`app.language`/`app.theme`（→票 39 open）」；`439-...-done.md:81`：「`app.theme`｜**丙＋丁**｜GUI 半边票 39；球侧半边摆机主」；`83-...-done.md:129`：「📋 票 39」 |
| `app.theme`（球侧半边） | **无票认领** | — | `439-...-done.md:81` 明写「**球侧半边摆机主**」；`docs/reports/pending-and-issues.md:36331` 把它列成待人拍板项「①批"球吃 app.theme"→开小票接线（缝在盘）②批"theme 只归 GUI"→归票 39 ③维持现状」⇒ **至今没开票** |
| `panel.font_size` | **票 437**（间接） | **done**（`437-the-313-door-reports-seven-writable-fields-the-two-rosters-do-not-name-done.md`） | `437-...-done.md:47` 点名 `panel_font_size`；`:135` 逐字段真分支 `panel_font_size :377→:386`。⚠ 但 437 认领的是**门报名册**那一格，**不是「谁消费这个字号」**⇒ 消费引擎那半边**无票认领** |
| `panel.font_size`（消费引擎半边） | **无票认领** | — | 票池 `grep "panel_font_size\|font_size\|字号"` 全仓命中 8 处，**没有一张票点名「谁把 panel.font_size 接到页面字号上」**。`180-...-done.md:270` 只说「现含 `font_size`(12)」，是**键存在**，不是**谁消费** |
| `panel.width` | **票 180**（间接） | **done**（`180-panel-width-is-a-config-field-with-zero-production-readers-so-changing-config-toml-has-no-visible-effect-done.md`） | `:288`/`:310`/`:351` 逐跳记它的读者链；`config_readers_255.go:310` 也点名 |
| `ball.size` | **票 413 + 422** | **done**（`413-...-done.md`、`422-...-done.md`） | `413-...-done.md:38` 记写侧+读侧；`422-...-done.md:41` 记面板接线 |
| `ball.style` | **票 65** | **done**（`65-ball-glass-quality-rework.md`，无 `-done` 后缀但 `:192` 记了 UI 腿交付） | `65-...md:192`：「设置页「球风格」切换 UI 落地」 |
| `ball.opacity_idle` / `ball.click_through` | **票 437**（间接）＋外观波（**无独立票号**） | 437 **done**；外观波**无票** | `437-...-done.md:44` 点名 `ball_opacity_idle`；`cmd/wisp/appearance_wave_windows_test.go` 与 `internal/panel/config_write_wave.go` 都不带票号，文件头写「the appearance wave」。`grep "外观波" .scratch/wisp/issues` 只命中 180/437 的**提及**，**没有一张票以它为题** |
| `panel.opacity` / `panel.animations` | **无票认领** | — | 它们**不是 config 键**（`PanelSection` 无此二键），是页面本地 CSS/DOM 控件。票池零命中 |
| **左栏「成本」格**（§4.3 的假话来源） | **票 194**（Go 侧读口，**done**）＋ **票 40**（页面，**open**） | 194 **done**；40 **open** | `194-...-done.md:124`「AC#5 ★原"生产调用者零枚"**已被后续腿解决**…**七格全勾，票 194 -done**」；`194-...-done.md:110`「前端 `panel.ts` 不携带乙堆任何名字」⇒ **194 的正文自己就说了"前端不携带"**；票 40 `40-security-privacy-cost-pages.md` 标题即「Security/Privacy/Cost pages」，`:3` `Status: ready-for-agent`（**open**）、`:6` `Blocked by: 29-memory-l1l2, 35-panel-bridge-c17, 44-costmeter-c23` |
| **左栏「安全」格**（档位已推但未渲染） | **票 40** | **open** | 票 40 `:1` 标题含「Security」；`:2` purpose 写「权限档位、门控阈值、批量授权的边界」⇒ 正是这一格 |
| **左栏「技能」格** | **票 215b**（**open**，blocked-on-owner）＋ 215a（**done**，插件半） | 215b **open**；215a **done** | `215b-skills-inventory-needs-a-config-face.md:3`「Status: **blocked-on-owner（契约面：动 config.toml 的 section 树＝D36）**」；`:10` AC#2「技能清单可枚举…每枚至少一枚**生产读者**（grep 改前 0 ⇒ 改后 ≥1）」；`:17`「**要机主批准：给 config.toml 增加 [skills] section（D36 配置树扩一个 section）**」 |
| **左栏「记忆」格** | **票 29**（**blocked-by**，票 40 自己点名） | open | 票 40 `:6` `Blocked by: 29-memory-l1l2` |
| **左栏「MCP」格** | **票 215**（**done**，但明确**不接 MCP**） | 215 **done**；**本项永久无票** | `215-composer-plus-menu-skills-plugins-inventory-done.md:28` 逐字：「**MCP 不在本票**：要不要接 MCP 是 **D13 的取舍**，本票只保证"如果哪天接，菜单里有位置"——**这一条不许顺手做**」；`:40`「不接 MCP（D13 取舍）」；`50-tier1-manifest-plugins-done.md:35`「MCP (**RESERVED**)」；`225-...-done.md:32` 记 `internal/tools/registry.go:55` 的 `KindMCP` 带字样 `REJECTED: D13/16.9#7, interface slot only` ⇒ **架构上主动不接**，不是"忘了接" |
| **左栏「引导」格** | **无票认领** | — | 票池零命中。首启画面归票 41/14（票 39 `:32` 提到「First-run wizard (firstrun screen lands with 41/14 polish)」），但**首启屏 ≠ 设置页的引导格** |

### ★ 票面缺口

**「外观波」（ball.opacity_idle / ball.click_through / panel.width / panel.font_size 四枚）没有自己的票号。**
证据：四枚的 Go 侧载体文件全部无票号前缀——
`internal/panel/config_write_wave.go`（文件头只写 "the appearance wave"）、
`internal/panel/ball_idle_appearance_wave.go`、`internal/panel/panel_tuning_wave.go`、
`cmd/wisp/appearance_wave_windows_test.go`、`cmd/wisp/panel_appearance_wave_test.go`。
唯一沾边的是 437（只认领门报名册）与 180（只认领 `panel.width` 的读者链）。
⇒ **这一族没有一张票在盯「谁消费它」**。

---

## 4【判死】有没有「其实已经接了却仍显示未接线」？

### 4.1 就 `config-screen.tsx` 这一页：**0 行**

九行里唯一显示「未接线」的 `app.theme`，**确实没接**（§2.1：读者 0、写侧 0、球侧缝零调用者、面板侧纯本地态）。
⇒ **本页没有假免接线。**

### 4.2 但全仓有 **1 处**「已接却仍说未接线」的**假话**——而且它正在误导人

**`internal/config/unwired.go` 的守门表**，与票 180 的实测**已经矛盾**：

| 位置 | 现状 |
|---|---|
| `internal/config/schema.go:603` | `FontSize int \`toml:"font_size" default:"12"\`` —— 键**在 schema 里**，能加载、能落盘 |
| `internal/config/settings_appearance_wave.go:101` | `base.Panel.FontSize = px` —— **有真 setter** |
| `internal/config/settings_appearance_wave.go:179-188` | `PanelFontSizeOptions()` —— **有真 reader** |
| `cmd/wisp/run.go:1272-1274` | **装配根真接了快照读者** |
| `cmd/wisp/panel_config_store.go:431-440` | **写路由真接了 setter** |
| ⇒ | **四跳全通。** 这枚键**不是未接线**。 |

⇒ 与 §2.2 合起来：`panel.font_size` 是**「写通／读通／快照通／引擎缺」**——
**接线接完了四分之三，缺最后一跳（谁把它接到页面字号上）**，
而 `internal/config/schema.go:594-602` 的注释仍把这最后一跳当"待落地"。

### 4.3 ★★ 真正的「假免接线」在**另一个页面**：那才是「一堆未接线」的来源

`frontend/src/components/harness/settings-page.tsx:49-132` 的 `CATEGORIES` 表——
**每一格**都用 `waiting: "..."` 手写一句话说明「它还在等哪个数据源」，
`settings-page.tsx:202` 把这类格的 aria 名改成 **`${c.name}（还没有数据）`**，
`:229-234` 再画一枚灰点当「未就位」标记。

**12 格里 9 格是 waiting 格**（唯一不 waiting 的是 外观/模型/语音）：

| 类别 | waiting 文案（`settings-page.tsx` 行号） | Go 侧现读：有没有源 | 真/假 |
|---|---|---|---|
| 通用 | `:57` 「配置读口已经落地…这一类自己的键还没有节推上来」 | ⚠ **半句**：config.get 确实落地并在用（`internal/panel/config_handlers.go:429`）；但**这一格自己的键**（`app.language`/`app.autostart`/`app.single_instance`）确实无节——`internal/panel` 零 `[app]` 节 | 文案基本为真（措辞绕，事实成立） |
| 记忆 | `:89` 「还没收到记忆数据。D20 定义了记忆系统，但没有任何一条推送告诉面板」 | 无源（真）：`internal/panel` 零 memory 快照节 | 真 |
| 技能 | `:94` 「还没收到技能清单。技能要能被列出来，前提是它们从插件运行时注册上来」 | ⚠ **半句假**：**插件清单已有源**（`internal/panel/plugins_inventory_215.go:176-183` 的 `engineState`/`skillState` 节，reader 在 `cmd/wisp/run.go:1313` `Plugins:`）；**技能清单**确实无键——同文件 `:86` `SkillsNotInConfigModel215` 逐字「技能维今天没有配置键，无清单可枚举」 | 技能半真；插件半已有源 |
| MCP | `:99` 「还没收到 MCP 状态。D13 定了 MCP 的取舍，但没有字段把已连的服务端报给面板」 | 无源（真）：`git grep "MCP\|Mcp" internal/panel` ＝ **0 命中** | 真 |
| 安全 | `:104` 「还没收到策略快照。这一节最敏感…」 | ⚠ **半句假**：**档位已每包推**（`cmd/wisp/run.go:1187` `rt.modes.PermissionMode` → `internal/panel/pump.go:372-373`）；但档位/阈值/批量授权那一整面确实无快照节 | 档位已有源，其余无 |
| 隐私 | `:109` 「还没收到隐私配置」 | 无源（真）：`internal/panel` 零 privacy 快照节 | 真 |
| 成本 | `:114` 「还没收到用量数据。成本口径由 C23 定，但面板快照里没有 token 数」 | ⚠ **后句假**：**Go 侧读口已接**——`cmd/wisp/panel_cost_194.go:51` `CostDays` 读 `memory.Store.ListCostDaily`；`cmd/wisp/panel_inbound.go:354-355` 装配；`internal/panel/cost_handlers.go:86` 出 `wisp-cost-summary:` payload 行；票 194 `:124` 记「七格全勾，票 194 -done」。**但前端零调用**：`git grep "cost.summary\|costSummary" frontend/` ＝ **0 命中** | **后句假**（源有、页面没去取） |
| 引导 | `:119` 「还没收到引导进度」 | 无源（真） | 真 |
| 悬浮球 | `:130` 「球的大小、外观风格、空闲透明度、穿透都已能切…唤醒词与热键还没有数据源」 | 前半**真**（§2 已证四行已接）；后半「唤醒词与热键还没有数据源」——`git grep "Hotkey\|WakeWord" internal/panel` ＝ **0 命中** ⇒ 也真 | **真** |

⇒ **「一堆未接线」的本体＝设置页左栏那 9 枚灰点 + 「（还没有数据）」aria 名**，
**不是** `config-screen.tsx` 的「未接线」徽章（那只有 1 枚）。
⇒ 其中**至少 3 枚文案与事实有出入**：`安全`（档位已每包推）、`成本`（Go 侧读口已接、票 194 已结，前端没去取）、
`技能`（插件清单已接，技能半确实无键）。
⇒ **最值得跟机主讲的一条＝`成本`**：**数据已经在门口，只差前端去取**——
`cmd/wisp/panel_cost_194.go:51` 读真表、`internal/panel/cost_handlers.go:120` 出 payload 行，
而 `frontend/src` 里 `cost.summary` **零调用**。
这一格不是「没数据」，是「没去拿」——**接上只要一个前端调用，不是新造引擎。**

### 4.4 ★★★ 最硬的一条：**「一堆未接线」是真的，但那是 2026-11-11 之前的旧页面**

`git show <commit>:frontend/src/components/config-screen.tsx` 逐提交数「`blocked` 非空的行数」
（尺＝数 `blocked: "` 后面不是空串的行）：

| 提交 | 日期 | 总行数 | **未接线行数** | 备注 |
|---|---|---|---|---|
| `ad9b29fd9` | 2026-09-26 | 8 | **6** | 加 `panel.animations` 那次；只有 opacity/animations 两行是活的 |
| `a3fdb50dd` | 2026-11-07 | 9 | **6** | 加 `ball.style`（票 65）；其余六行仍 `blocked` |
| `17e210c33` | 2026-11-11 | 9 | **5** | 加 `ball.size`（票 422） |
| `4017962b3` | 2026-11-11 | 9 | **1** | 加外观波四行（`ball.opacity_idle`/`ball.click_through`/`panel.width`/`panel.font_size`）⇒ **6→1** |
| `HEAD` | 2026-11-16 | 9 | **1** | 未再变 |

⇒ **机主看到「一堆未接线」的那一屏，是 2026-11-11 之前的旧构建**（6 枚）。
⇒ **今天 HEAD 只有 1 枚**（`app.theme`），其余六枚在 **2026-11-11 一天之内**从「未接线」变成能改/立即-重启/关窗-重启生效。
⇒ **⚠ 但今天真机仍然看不到新页面**（下一条）—— 所以机主在真机上看到的仍是**旧构建**。

### 4.5 ★★★ 承接上一条：**真机今天根本看不到设置页**——页面没被打包进二进制

- `frontend/embed.go:19-20`：`//go:embed all:dist` —— 页面只能从 `frontend/dist/` 进二进制。
- **现读**：`git ls-files frontend/dist` = **1 枚**（`frontend/dist/.gitkeep`）；
  盘上 `frontend/dist/` 也**只有** `.gitkeep`（0 字节）。
- ⇒ `internal/panel/assets.go:54-60` `newAssets()` 找 `index.html`，**找不到 ⇒ `Built()==false`**；
  `cmd/wisp/panel_host_windows.go:1024-1025` 直接报错
  「panel host: embedded bundle not built (only the tracked placeholder is present)」。
- ⇒ **真机上 `wisp` 二进制里没有页面**。机主看到的「设置页」要么来自
  `?harness=2` 的开发服务器（`frontend/src/main.tsx`），要么来自 `design/doubao/demo/screens/config.js`
  （**那份 demo 里每一枚徽章都写「即时」**，`config.js:287/292/296/300/304/308`——**它比真页面乐观得多**）。

**⇒ 所以「设置页上一堆未接线」这个观察，只可能来自 `?harness=2` 的开发服务器或 demo，
而那两处的数据源与真机不同。** 这是本普查最该先跟机主讲清的一条。

---

## 5. 汇总数字

| 口径 | 数 |
|---|---|
| 设置页外观节总行数 | **9** |
| ├ 徽章字面「即时」 | **2**（`panel.opacity`、`panel.animations`） |
| ├ 徽章字面「未接线」 | **1**（`app.theme`） |
| └ 徽章字面 tier 词（能改/立即-重启/关窗-重启生效） | **6** |
| 其中**真**接了（读+写+引擎全通） | **5**（`ball.size`、`ball.style`、`ball.opacity_idle`、`ball.click_through`、`panel.width`） |
| 其中**写通但零消费引擎** | **1**（`panel.font_size`） |
| 其中**真·零读者** | **1**（`app.theme`） |
| 其中**不是 config 键**（页面本地） | **2**（`panel.opacity`、`panel.animations`） |
| **假免接线（本页内）** | **0 行** |
| **假免接线（全仓）** | **3 处**：①`settings-page.tsx:114` 成本格说"没收到用量数据"，但 Go 侧读口已接（`cmd/wisp/panel_cost_194.go:51` ＋ `cmd/wisp/panel_inbound.go:354`），前端零调用 ②`settings-page.tsx:104` 安全格说"没收到策略快照"，但档位每包都推（`cmd/wisp/run.go:1187`） ③`panel.font_size` 四跳已通却无消费引擎（`internal/config/schema.go:594-602` 仍称"待落地"） |
| 设置页左栏灰点（`settings-page.tsx` waiting 格） | **9 枚**（12 格里 9 枚）⇒ **这才是"一堆未接线"的本体** |
| 其中**文案与事实有出入** | **3 枚**（成本／安全／技能） |
| ★**「一堆未接线」的历史真相** | **2026-09-26→11-07 期间是 6 枚；11-11 一天之内 6→1**。机主看到的那一屏是**旧构建**（§4.4） |
| ★★**真机今天能不能看到新页面** | **不能**——`frontend/dist` 只有 `.gitkeep` ⇒ `panel.Assets.Built()==false`（§4.5） |

---

## 6. 复跑命令（可核）

```powershell
# 1) live 是静态字面量，不是 wire 字段
Select-String -Path frontend/src/components/config-screen.tsx -Pattern 'const live = row.blocked'
# ⇒ config-screen.tsx:411:  const live = row.blocked === "";

# 2) config-view.ts 里没有 live
Select-String -Path frontend/src/lib/config-view.ts -Pattern 'live'   # ⇒ 0 命中

# 3) config_handlers.go 里没有 live 字段（只有 EffectiveTier）
Select-String -Path internal/panel/config_handlers.go -Pattern '"live"|live\b'  # ⇒ 0 命中
Select-String -Path internal/panel/config_handlers.go -Pattern 'EffectiveTier|tierSentence'

# 4) app.theme 零生产读者（命中全在 internal/config 内 + 测试）
git grep -n "App\.Theme" -- "*.go"
git grep -n "SetTheme(" -- "*.go"      # ⇒ 只命中 internal/ball/renderer_windows.go:106 定义

# 5) panel.font_size 四跳通、引擎缺
git grep -n "Panel\.FontSize" -- "*.go"   # ⇒ 6 命中，全在 internal/config + 测试
git grep -n "FontSize" -- cmd/wisp/panel_host_windows.go   # ⇒ 0 命中

# 6) 页面没打包（⇒ 真机无页面）
git ls-files frontend/dist            # ⇒ frontend/dist/.gitkeep   （1 枚）
Get-ChildItem frontend/dist -Recurse  # ⇒ 只有 .gitkeep（0 字节）

# 7) 左栏 waiting 格（"一堆未接线"的本体）
Select-String -Path frontend/src/components/harness/settings-page.tsx -Pattern 'waiting: "'
```

---

## 7. ⛔ 边界（诚实登记）

- 本腿**只读**：未改任何产码、未改任何票面、⛔ 未 push、未建/删任何仓内件（本件是唯一的写入，且是 `docs/evidence/s1/` 证据件）。
- 行号均为 **2026-11-16 20:1x** 现读所得；本仓行号会随其它腿的提交漂移，引用前请按 §6 的尺现跑。
- §4.4 的「真机看不到设置页」是**由 `frontend/dist` 只有 `.gitkeep` 推得**，
  ⛔ 未在真机双击 `wisp.exe` 复现（本次审计不含起进程）。
- §3 的「无票认领」= 在 `.scratch/wisp/issues/*.md` 里 grep 该键/该功能**零命中**，
  **不是**断言「世界上没有票」——票池可能有以别的名字描述同一件事的票。

---

**撤销口令**：`删 docs/evidence/s1/census-unwired-20261116-2017.md`
