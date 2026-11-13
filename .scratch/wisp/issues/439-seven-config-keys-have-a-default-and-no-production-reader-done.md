# 439 — 配置树里还剩 **7 枚**"有 `default` 但生产码零读者"的键（`[panel] width` 的同族，本程普查从 30 收窄到 7）

- Status: **done（7 枚逐枚归因：0 甲 / 0 乙 / 4 丙 / 2 丁；3 件摆机主）**。
- 来路：票 **180** 的 AC#4 同族普查（本程 2026-11-13，件 `docs/evidence/s1/180-panel-width-reader.md`，提交 `5da3848e2`）。
- 关联：票 180（`[panel] width` 本身＝归口票，读者已由票 255 AC#4 接上）· 票 423（上一轮把 30 枚压到 7 枚的主力）· 票 147／139／97（"字段存在、无人读"那一族）· D36（配置模型三档生效级别）。

## 现量（票 180 腿本程跑，我复核过提交与件）

`internal/config/schema.go` 里带 `default` 标签的键共 **70 枚**，三判类：

| 判类 | 11-10 普查 | 本程现量 |
|---|---|---|
| R 有生产读者 | 33 | **48** |
| G 守卫拦截 | 7 | **8** |
| **Z 零生产读者** | **30** | **7** |
| 合计 | 70 | 70 |

**⇒ 上一轮（票 423 系列）已经把 30 枚压到 7 枚。本票收的就是这剩下的 7 枚。**

## 那 7 枚（逐字具名，含 default 与本程判读依据）

| # | 字段（toml 键） | default | 真读数 | 判读依据 |
|---|---|---|---|---|
| 1 | `app.language` | `zh-CN` | **0** | 唯一 config 侧命中＝`manager.go:379` 重启档比较＋tiers/名册字符串；`speech`/`web_search` 的 `.Language` 是别家字段 |
| 2 | `app.theme` | `dark` | **0** | 全部在 config 包内：`manager.go:378/383` 比较/拷贝、`validate.go:51` 枚举、`tiers.go:59`、roster 字符串 |
| 3 | `ball.hide_on_fullscreen` | `true` | **0** | 全树（含测试）仅 schema 定义＋`boundary_test` 赋值＋341 注释；`internal/ball` 无任何 fullscreen 读点 |
| 4 | `voice.wake_word.veto_words` | `取消,停下,别` | **0** | 仅 `manager.go:434/441` 段拷贝/比较；KWS 侧 `kwsVetoWords` **空表**（票 423 已记） |
| 5 | `voice.tts.provider` | `local-sherpa` | **0** | 仅 `manager.go:423` 段拷贝；云 TTS speak 半边未实现（423） |
| 6 | `voice.tts.speed` | `1.0` | **0**（**死读点**） | 唯一 config 读点 `internal/speech/cloudcatalog.go:119` `ep.Speed = cfg.Voice.TTS.Speed`，但下游没人用它 |
| 7 | `observe.slo_sample_interval_sec` | `5` | **0** | 全树（含测试）仅 schema 定义＋`boundary_test` 赋值；SPEC-03 §4.3 冻结 watchdog 1s 双职责共用 |

## 本票要做的

- [x] **AC#0 逐枚复核（只读，⛔ 产码）**：对上面 7 枚**每一枚**自己重跑那把尺（`grep -rn "\.<字段名>" --include=*.go internal/ cmd/`），逐枚给：命中数 | 其中"真读"的命中 | 与上表的差集。⚠ **⛔ 照抄上表**——那是别人的读数，你要自己复现。
- [x] **AC#1 逐枚归因**：每一枚判成下面四类之一，**逐枚写理由**：
  - **(甲) 该接读者**（规格要求它生效，只是没接上）⇒ 落地＋成对判据
  - **(乙) 该删字段**（规格不要求，是历史残留）⇒ 删＋证据
  - **(丙) 归口他票**（已有票认领，如票 423 的余量）⇒ 具名票号
  - **(丁) 规格待定**（D36 或 specs 里没写清）⇒ 摆机主
- [x] **AC#2 成对判据（对每一枚判"甲"的）**：正向＝设两个不同值 ⇒ 结果随之变（⛔ 只断"字段非零"）；反向＝**⛔ 用"改默认值"交差**。
- [x] **AC#3 契约轴**：⛔ 改 `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 一字。若某枚的归因需要改契约文字 ⇒ **摆机主批**，⛔ 自己改。
- [x] **AC#4 门禁**：`go test -count=1 ./internal/config/ ./internal/panel/` ＋ `gofumpt -l internal cmd tools scripts frontend`（无输出）＋ `scripts/d22scan.sh`（干净）。⚠ `./cmd/wisp/` 本机测不到（票 98），⛔ 算进你的红。
- [x] **AC#5 边界声明**：明确写"本票收哪几枚、⛔ 收哪几枚"，并把余下的留回台账。

## 写面
`internal/config/**` ＋（判"甲"时）该字段的消费者包。⚠ **开工前先 `git status --porcelain internal/config` 记基线**——同波可能还有别的腿。

## 判据提示（⛔ 跳过）
这 7 枚里有几枚**看起来该删**（`app.theme`／`app.language` 在 config 包内自转圈），有几枚**看起来该接**（`voice.wake_word.veto_words` 的 KWS 侧是空表，那是真功能缺口）。
**但⛔ 按我的猜测做**——逐枚复核后再判。**同族 30→7 说明这一族是真问题，不是理论洁癖。**

## Progress log（append-only）

- [2026-11-13T12:4x+08] agent=439-w1 did=AC#0＋AC#1 全量交付，**AC 框未动（归编排者）**。起手 HEAD `673da11df`；`git status --porcelain internal/config` 基线空（干净）。**AC#0**：票面那把粗尺 7 枚逐枚现跑＋裸词补尺（VetoWords/HideOnFullscreen/SLOSampleIntervalSec/ThemeLight|ThemeDark/Language）＋段内立尺（Voice\.TTS\.Provider/Speed），粗尺命中 13/22/3/10/2/4/3，真读**全部 0**，与 180 表**无一处矛盾**（一处精化：`voice.tts.speed` 的下游消费其实在盘——`cloudtts.go:140 [Speed: t.ep.ttsSpeed()]` 写进云端请求体；死读点死在 `ResolveCloudChain` 非测试调用者 0，名册 :333 同判）。**AC#1 归因：0 甲 · 0 乙 · 4 丙 · 2 丁**——`app.language`→**丙·票 39**（GUI 配置编辑器 open；restart 档机器在盘：planApp＋reportRestartPending＋255r2 测试钉名，缺的是"重启时定 UI 语言"的壳，i18n 在 SPEC-12 §5:70 DEFERRED 到 S8）；`app.theme`→**主判丙·票 39**＋球侧半边**附丁**（球缝在盘 `ball.Options.Theme`→`SetTheme`→`PaletteFor`，但二值词表与第四代 C21 主题（owner 2026-09-25 pivot，tokens.go:73-75）的关系规格没写清、球侧接线无认领票）；`ball.hide_on_fullscreen`→**丙·票 62/68**（open；功能半边不在：internal/ball 无任何全屏判定、Options 无此字段，接前先造"全屏判定"）；`voice.wake_word.veto_words`→**丁**（消费半边在盘且已接线：kws.go:291 分类→onKwsHit:1155-1159 打响；缺的是**数据**——kws.go:78-82 逐字说要"editing the keywords file"重生成转导表（模型资产＋manifest sha256 46de8c68 变更），config 喂值＝假绿，disposition §一#2 同判）；`voice.tts.provider`→**丙·票 61 收尾**（SPEC-12 §5:81 行 open；引擎/链/解析器全在盘：NewCloudTTSChain:144 实装 TtsEngine、KeyResolver 缝＝secret.Store.Resolve 生产现用；装配现状＝engines:436-448 无条件 matcha＋boot 门 :1866 只拦 ASR；完整收尾票非顺手一行，且 cmd/wisp 共享面每波一枚）；`voice.tts.speed`→**丙·随 61 搭车**（61 接上 ResolveCloudChain 即自动词活；本地半边 tts.go:34-48 是 15-b/15-c 测量纪律常数，不旋钮化＝终态）；`observe.slo_sample_interval_sec`→**丁**（SPEC-03 §4.3 冻结 1s tick＝config_reload.go:101 的 `configReloadPollInterval`，SustainTicks=60 语义按它定义；编排者已裁挂账 DEFERRED＋C1 不拆，台账 34012 逐字）。**AC#2**：零甲⇒无成对判据可跑（如实记，非漏跑）；反向自动满足——`git diff HEAD -- internal/config/schema.go` 空，70 枚 default 标签零字节动。**AC#3**：PLAN.md/specs/thresholds/golden/allowlist 零触碰；**三处摆机主**（①theme 皮跟谁走 ②veto 模型资产批不批 ③slo 字段语义三选一——选项与不答的代价逐条在证据件 §3 表）。**AC#4 门禁五条**（commit 前现跑，本程 Go 零改动⇒读数＝基线）：`go build ./cmd/wisp/` rc=0 无输出；`go vet ./cmd/wisp/` rc=0 无输出；`go test -count=1 ./internal/config/ ./internal/panel/` ⇒ `ok ...internal/config 2.1xxs`＋`ok ...internal/panel 15.xxs` rc=0 全绿；`gofumpt -l internal cmd tools scripts frontend` 空输出；`sh scripts/d22scan.sh` 末行 `d22scan: clean - no D22 ban violations; ...` rc=0。./cmd/wisp 包测试按票 98 不进红账。**AC#5 边界**：本票收 0 枚接线、0 枚删字段；7 枚去向全具名（39／39＋丁／62·68／丁摆机主／61 收尾／61 搭车／丁摆机主），留回台账。交付：证据件 `docs/evidence/s1/439-seven-config-keys.md`＋本 log＋commit（显式 pathspec，不 push）。写面：internal/config 零触碰、frontend/design 零触碰、build/mut_rvw.go 未删。
---

## ★结案（2026-11-13 12:4x，编排者翻框）— **7 枚全部复现零读者，0 枚该接、0 枚该删**

### 三把尺（我现跑）
- **提交真在**：`3cdc57faa`（2 files, +176）
- **件真在**：`docs/evidence/s1/439-seven-config-keys.md`（17306B）
- **自己复跑**：`go test -count=1 ./internal/config/ ./internal/panel/` 由腿跑 `ok config 2.582s` ＋ `ok panel 19.098s`

### AC#0：7 枚逐枚复核，**与票 180 的表无一处矛盾**
| # | 字段 | 粗尺 | 真读 |
|---|---|---|---|
| 1 | `app.language` | 13 | **0** |
| 2 | `app.theme` | 22 | **0** |
| 3 | `ball.hide_on_fullscreen` | 3 | **0** |
| 4 | `voice.wake_word.veto_words` | 10 | **0** |
| 5 | `voice.tts.provider` | 2 | **0** |
| 6 | `voice.tts.speed` | 4 | **0**（死读点） |
| 7 | `observe.slo_sample_interval_sec` | 3 | **0** |

**一处精化（#6）**：下游消费其实**在盘**——`cloudtts.go:140` 把 `ep.ttsSpeed()` 写进云端请求体；**死的是 `ResolveCloudChain` 全仓非测试调用者 0** ⇒ "接票 61 即活"钉得更死。

### AC#1：归因 **0 甲 · 0 乙 · 4 丙 · 2 丁**
| 字段 | 判 | 归口/理由 |
|---|---|---|
| `app.language` | **丙** | 票 **39**（restart 档机器全在盘，缺"重启时拿它定 UI 语言"的壳；i18n 在 SPEC-12 §5:70 DEFERRED 到 S8） |
| `app.theme` | **丙＋丁** | GUI 半边票 39；**球侧半边摆机主**（见下 §丁#1） |
| `ball.hide_on_fullscreen` | **丙** | 票 **62/68**（功能半边不在：`internal/ball` 无"前台窗口是否全屏"判定，想接都没对象） |
| `voice.wake_word.veto_words` | **丁** | 消费半边**已接线在盘**（`kws.go:291`），缺的是**数据**（veto 词必须进 keywords 转导表＝模型资产＋manifest sha256 变更） |
| `voice.tts.provider` | **丙** | 票 **61** 收尾（引擎/链/解析器全在盘；装配现状 `engines:436-448` 无条件 matcha） |
| `voice.tts.speed` | **丙** | 随票 61 搭车（接上 `ResolveCloudChain` 即自动活，零额外工作量） |
| `observe.slo_sample_interval_sec` | **丁** | **任何接法都动冻结面**（SPEC-03 §4.3 冻结 1s tick），字段语义没有一句契约文字 |

### ★腿纠正我的一处
**#4 `kws.go:291` 的 `Veto: kwsVetoWords[...]` 读的是 speech 包内 map，不是 config 字段** —— 这恰好是"消费半边已接线、数据半边缺"的证据。**不是"完全没接"。**

### AC#2：零枚判甲 ⇒ **无落地、无成对判据可跑**（如实记，非漏跑）
反向自动满足：`git diff HEAD -- internal/config/schema.go` ⇒ **空**（70 枚 default 标签零字节动）。

### AC#5：边界声明
**本票收：复核＋归因＋台账回填。⛔ 收：接线（0 枚该接）、删字段（0 枚该删）。**
4 枚丙已有具名票认领；2 枚丁 + 主题球侧半边**摆机主**（见 `A1826`）。

**撤销口令：「439 回红」**。

