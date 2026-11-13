# 439 — 配置树里还剩 **7 枚**"有 `default` 但生产码零读者"的键（`[panel] width` 的同族，本程普查从 30 收窄到 7）

- Status: **ready-for-agent（先逐枚归因，再按枚定修法）**。
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

- [ ] **AC#0 逐枚复核（只读，⛔ 产码）**：对上面 7 枚**每一枚**自己重跑那把尺（`grep -rn "\.<字段名>" --include=*.go internal/ cmd/`），逐枚给：命中数 | 其中"真读"的命中 | 与上表的差集。⚠ **⛔ 照抄上表**——那是别人的读数，你要自己复现。
- [ ] **AC#1 逐枚归因**：每一枚判成下面四类之一，**逐枚写理由**：
  - **(甲) 该接读者**（规格要求它生效，只是没接上）⇒ 落地＋成对判据
  - **(乙) 该删字段**（规格不要求，是历史残留）⇒ 删＋证据
  - **(丙) 归口他票**（已有票认领，如票 423 的余量）⇒ 具名票号
  - **(丁) 规格待定**（D36 或 specs 里没写清）⇒ 摆机主
- [ ] **AC#2 成对判据（对每一枚判"甲"的）**：正向＝设两个不同值 ⇒ 结果随之变（⛔ 只断"字段非零"）；反向＝**⛔ 用"改默认值"交差**。
- [ ] **AC#3 契约轴**：⛔ 改 `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 一字。若某枚的归因需要改契约文字 ⇒ **摆机主批**，⛔ 自己改。
- [ ] **AC#4 门禁**：`go test -count=1 ./internal/config/ ./internal/panel/` ＋ `gofumpt -l internal cmd tools scripts frontend`（无输出）＋ `scripts/d22scan.sh`（干净）。⚠ `./cmd/wisp/` 本机测不到（票 98），⛔ 算进你的红。
- [ ] **AC#5 边界声明**：明确写"本票收哪几枚、⛔ 收哪几枚"，并把余下的留回台账。

## 写面
`internal/config/**` ＋（判"甲"时）该字段的消费者包。⚠ **开工前先 `git status --porcelain internal/config` 记基线**——同波可能还有别的腿。

## 判据提示（⛔ 跳过）
这 7 枚里有几枚**看起来该删**（`app.theme`／`app.language` 在 config 包内自转圈），有几枚**看起来该接**（`voice.wake_word.veto_words` 的 KWS 侧是空表，那是真功能缺口）。
**但⛔ 按我的猜测做**——逐枚复核后再判。**同族 30→7 说明这一族是真问题，不是理论洁癖。**
