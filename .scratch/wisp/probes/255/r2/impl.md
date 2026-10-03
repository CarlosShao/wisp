# 票 255 · 实现腿 `255-r2` · 证据件

写面＝`cmd/wisp/`＋`internal/config/`（本腿只做 AC#1 主格＋AC#3 残余补尺）。
本件由实现腿自己写，**判语不归自己**：AC 勾选框一枚未碰。

---

## 0. 起手锚（同一发命令取数）

```
$ date
Sat Oct  3 09:10:04 CST 2026
$ git log -1 --format=%h
8a3790f0
$ git status --porcelain -- cmd internal tools docs .scratch | wc -l
352
$ git status --porcelain -- cmd internal
（空＝本腿写面起手干净）
```

- 起手时刻 `2026-10-03 09:10 +08`，锚点 HEAD `8a3790f0`。
- `cmd internal` 两枚目录 porcelain＝**0 行** ⇒ 工作树在写面上与 HEAD 同字节（AC 判据里"工作树＝HEAD"的那把尺成立）。
- 全仓 porcelain 352 行主要来自 `.scratch/`（别的腿的临时件，规则＝只建不删），与写面无冲突。
- 开工前已完整读票面：`.scratch/wisp/issues/255-config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers.md`（55 行，含编排者 2026-10-02 11:00 裁定小节与 13:1x 行号更正）。

**AC#1 判据原文（逐字抄自票面第 17 行）**：

> AC#1 那句"已立即生效"只许说真话：判据＝stdout 里"这些段已立即生效"那一段**只列真有生产读者的段**；**正控**＝手改 `[panel] width`  ⇒ 那句话**不许出现 `[panel]`**（要么不出现、要么换成诚实形如"值已换、但当前无组件应用它"）；**反控**＝真有读者的段（逐名给出证据）**仍要出现**。⛔ 不许把整句删掉——那是把热加载的可见性摘了，本仓定式：⛔ 摘尺不许当修尺。

**AC#3 残余补尺判据（逐字抄自票面第 19 行编排者翻勾段末句）**：

> 范围披露：`[app]` 新增键不在键级尺射程（固定名单非反射走查）——我裁＝记范围不算欠账，补尺（一行循环）归 `255-r2` 顺手做

本腿 commit 序列：`c8480bce`（§0 骨架）→ `ae60a87c`（AC#1 名册＋回执＋AC#3 尺）→ `a4906b6c`（逐名方括号＋句子装配钉）→
`17ff058d`（gofumpt 归零）→ `13dd60b4`（证据件 §1—§7 终值）→ 本件最后一枚（A549 判给本腿的 `restartTierKeys` 幻影注释清理，§2.4＋§4 m3）。

---

## 1. 现量（回执那几行 · `rep.Hot` 的来源 · 逐名"谁读它"）

### 1.1 回执那几行（起手锚 HEAD `8a3790f0`，逐字）

```
cmd/wisp/config_reload.go:170		if len(rep.Hot) > 0 {
cmd/wisp/config_reload.go:171			fmt.Fprintf(rt.stdout, "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：%v\n", rep.Hot)
cmd/wisp/config_reload.go:172		}
```

`rep.Hot` 被**原样**印出。票面现量第 1 条说的就是这一行：它报的是"内存里的值换掉了"，不是"有人按新值做了事"。

### 1.2 `rep.Hot` 的三处生产者（`internal/config/manager.go`，本腿未改此文件）

| 处 | 行 | 逐字 | 判定条件 |
|---|---|---|---|
| 段级热应用循环 | `:303` | `rep.Hot = append(rep.Hot, s.name)` | `:301` `if !reflect.DeepEqual(s.oldr, s.newr)` ⇒ **只问值换没换** |
| `planApp` | `:359` | `rep.Hot = append(rep.Hot, "app")` | `:353` `themeChanged := cur.App.Theme != fresh.App.Theme` |
| `planVoice` | `:425` | `rep.Hot = append(rep.Hot, "voice")` | `:415-419` thresholds/veto_words/speed/punctuation 任一不同 |

段级循环的表在 `manager.go:275-292`（12 枚：ball hotkey session audio llm agent privacy memory panel cost models observe），
同源守卫在 `:294-300`（不在 `TierRegistry` 登记成 `"hot"` 就 `panic`）——**本腿一字未动，也没削弱**。

### 1.3 `config.TierOf` 的读者枚数（起手零 ⇒ 现一枚）

```
$ grep -rn --include=*.go "TierOf(" cmd internal tools scripts
cmd/wisp/config_readers_255.go:34:   ← 注释里引用（不是调用）
cmd/wisp/config_readers_255.go:178:  if tier, ok := config.TierOf(name); ok {   ← 生产调用者，一枚
internal/config/tiers.go:93:         func TierOf(path string) (tier string, ok bool) {   ← 定义
```

起手读数（同一命令，`dd92bb92` 之后、本腿之前）＝**只有 `tiers.go:93` 那一枚定义，零枚调用者**；
现在生产调用者＝**1 枚**（`cmd/wisp/config_readers_255.go:178`，在 `hotRowsFor()` 里）。
`go build` 的 `rc=0` 不算"接上了"的证明（Go 不为未使用方法报错），所以这里给的是调用点枚数。

### 1.4 逐名"谁读它"（本腿现跑扫描，非抄 180-a1）

扫描式（`cmd/wisp/config_receipt_255_test.go` 的 `scanReadFiles`，每次跑测试重量一遍）：
`cmd internal tools scripts` 四根，排除 `*_test.go`、`internal/config/`、`testdata/`、注释行、
以及本票名册文件自身（它的裁决字符串是**引用**代码，不读配置）；
命中形状＝`.<段字段>.<大写字段>` 或 `Config().<段字段>` / `Config.<段字段>`。
下表"生产读者文件"就是扫描输出，由 `TestTicket255RosterStillMatchesTheActualReadSites` 与名册逐字比对（不等即红）。

| hot 行（`TierRegistry` 的路径） | 生产读者文件（实测） | 具名 file:line（源码回读已钉住） | 裁决 |
|---|---|---|---|
| `llm` | cmd/wisp/panel_config_store.go, providers.go, run.go, internal/llm/resolver.go | **`cmd/wisp/panel_config_store.go:92` `cfg := s.mgr.Config()` → `:96` `cfg.LLM.Roles.Chat.Model`**＝每次 `ReadSettings` 都问活配置 | `consumed:`（唯一一枚可说"已立即生效"） |
| `agent` | cmd/wisp/run.go | `cmd/wisp/run.go:424` `rt.cfg = cfg`＋`:991` `cfg := rt.cfg`＋`:1014` `cfg.Agent.PerToolTimeoutMS` ⇒ 读者吃的是**装配快照** | `snapshot-only:` |
| `models` | cmd/wisp/models.go | `cmd/wisp/models.go:163` `cfg.Models.Mirror`；调用者是 `wisp models` 子命令（`:209/:243/:273`），各自 `:183` `func loadModelConfig` fresh LoadFile ⇒ 别的进程 | `other-process:` |
| `hotkey` | cmd/balldebug/main.go | `cmd/balldebug/main.go:238` `h := mgr.Config().Hotkey`＝旁路调试程序；出货宿主 `cmd/wisp/resident_ball_windows.go:171` `Hotkeys:  ball.DefaultHotkeys(),` | `debug-host-only:` |
| `panel` | （零） | `cmd/wisp/panel_host_windows.go:304` `Width:  420,`（`:305` `Height: 260,`）；`cmd/wisp/panel_host_windows.go:175` `func NewPanelManager(disp *panel.ComposerDispatch, assets *panel.Assets, dataPath string) *PanelManager`＝**无配置参数** | `no-reader:` |
| `ball` | （零） | `internal/ball/ball_windows.go:64` `SizePx   int`＝球用自己的 option，不经 config | `no-reader:` |
| `privacy` | （零） | `internal/observe/logging.go:50` `RedactPaths`（注释自称 mirrors `[privacy] redact_paths`，两枚 `InitLog` 都自己填） | `no-reader:` |
| `observe` | （零） | `cmd/wisp/logsink.go:149` `Level: logSinkLevel`（sinks 自选级别） | `no-reader:` |
| `session` `audio` `memory` `cost` | （零×4） | 扫描零命中 | `no-reader:` |
| `app.theme` | （零） | 扫描零命中；`manager.go:353/:358` 只比较与拷贝 | `no-reader:` |
| `voice.tts.speed` `voice.punctuation` `voice.wake_word.thresholds` `voice.wake_word.veto_words` | （零×4） | 扫描零命中（`cfg.Voice` 在 `internal/config` 之外无人取） | `no-reader:` |

`llm` 那枚 caveat（本腿自己声明，没藏）：活读者是**设置腿**（面板问一次答一次），跑着的 provider 链仍是装配时
`cmd/wisp/run.go:435` `res := llm.NewResolver(cfg, st)` 建的那条。所以 `[llm]` 那句"已立即生效"的射程＝
"下一次设置读取答新值"，**不是**"在飞的模型换了"。这句写在 `hotRowClaims["llm"]` 的字符串里，也写在该文件头注释里。

### 1.5 "真有生产读者"的两种读法，本腿取哪一种（写明以便复核）

- 读法 **(A)＝180-a1 的 R 类**："生产代码读它的值"（不分何时读）⇒ 按这条 `agent`、`models` 也算有读者。
- 读法 **(B)＝票面现量第 1 条那句"有人按新值做了事"** ⇒ 只有 `llm`（活读）满足。

本腿按 **(B)** 落地（本票的缺陷陈述就是"报的是内存换了、不是有人做了事"），
(A) 类那两枚**没有消失**：它们在诚实句里被点名，名册里逐枚记着读者文件:行。
若编排者裁的是 (A)，改动是两枚前缀：`hotRowClaims["agent"]`、`["models"]` 换成 `hotClaimConsumed` 并把 cite 换成读者行
——**这一条同时进 §5**。

### 1.6 修好之后生产代码实际印出的两句（`go test -v` 从产品流里 log 出来，逐字）

```
wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[llm]
wisp run: 配置热加载：这些段的值已换进本进程内存，但本宿主没有会按新值做事的读者，本次运行不会因此改变行为
（票 255 AC#1：这一半不许说成「已立即生效」；逐段的读者判定见 HOT-RELOAD-READER 行）：
[ball] [session] [audio] [agent] [privacy] [memory] [panel] [cost] [models] [observe] [hotkey] [app] [voice]
```

台账逐字（同一次运行，`config: HOT-RELOAD-READER`）：

```
[audit] config: HOT-RELOAD-READER section=panel tier_row=[panel] claims=["panel":no-reader: nothing outside internal/config reads cfg.Panel - cmd/wisp/panel_host_windows.go:304 [Width:  420,] is hard-coded and NewPanelManager receives no config (票 255 AC#4 owns that break)]
[audit] config: HOT-RELOAD-READER section=llm tier_row=[llm] claims=["llm":consumed: cmd/wisp/panel_config_store.go:96 [cfg.LLM.Roles.Chat.Model] - configStore.ReadSettings calls s.mgr.Config() per call]
[audit] config: HOT-RELOAD-READER section=app tier_row=[app.theme] claims=["app.theme":no-reader: 扫描零命中：nothing outside internal/config reads cfg.App - manager.go's planApp compares and copies it, and no component re-skins from it]
[audit] config: HOT-RELOAD-READER section=voice tier_row=[voice.punctuation voice.tts.speed voice.wake_word.thresholds voice.wake_word.veto_words] claims=[...四枚键级行逐枚...]
```

⇒ 手改 `[panel] width` 那一发（单段场景，正控用例实测）：立即档那**一句整行不出现**，
诚实句为 `...：[panel]`，同一 tick 的台账给出 `no-reader` 与 `panel_host_windows.go:304` 的 cite。

---

## 2. 两格各自的修法与代码

### 2.1 AC#1（主格）：回执只说真话

**落点＝`cmd/wisp/config_readers_255.go`（新增）＋`cmd/wisp/config_reload.go` 的 `reportReload`（改造）。**

名册结构（`config_readers_255.go`）：

| 件 | 是什么 |
|---|---|
| `hotRowClaims map[string]string` | **17 枚** hot 行（12 段级 + `app.theme` + 4 枚 voice 键级）逐枚裁决；值的前缀五选一：`consumed:` / `snapshot-only:` / `no-reader:` / `other-process:` / `debug-host-only:` |
| `sectionReadSites map[string][]string` | 每段**允许**出现读者的文件集合（实测值），扫描对不上即红 |
| `hotRowsFor(name)` | 经 `config.TierOf(name)` 解析；段级不命中时收 `TierRegistry` 里 `name+"."` 前缀且 tier=="hot" 的键级行 ⇒ `[app]`/`[voice]` 的段名不再被当成"一枚档" |
| `splitHotTier(hot)` | 三分：`claimable`（该名字背后**每一枚** hot 行都是 `consumed:`）/ `quiet`（热应用但无活读者）/ `disagree`（plan() 报了名、登记表解释不了） |
| `bracketed(names)` | 名单渲染成 `[panel] [llm]`（每段自带括号）。理由：AC#1 的判据是"那句话里不许出现 `[panel]`"，而 `%v` 多元素时印成 `[ball session panel]`——段名在、括号不在，人和 grep 都会漏看，谎报能在抽查里活下来 |
| `verdictFor(name)` | 给台账用：这段的判定落在哪几枚行、每行的裁决原文 |

回执改造（`cmd/wisp/config_reload.go:178-207`，`reportReload`）：

- `:184` 每个 hot 名字一行台账：`config: HOT-RELOAD-READER section=%s %s`（票面问"谁读它"时那是**读台账**，不是猜）
- `:187` 原句**逐字保留**，只是列表换成 `bracketed(split.claimable)`：`"wisp run: 配置热加载：这些段已立即生效（D36 立即档）：%s\n"`
- `:190-194` 新增诚实句（票面允许的第二个形状："值已换、但当前无组件应用它"）
- `:196-205` `disagree` 另有一句并一条 `state=tier-disagreement` 台账——登记与实现对不上时**声张**，不静默

⛔ **整句没删**（AC#1 明文禁止摘尺当修尺；`:187` 那行的措辞与 ticket 223 的钉 `config_reload_223_test.go:274`（旧行号）仍逐字相容）；
⛔ **没去接线 AC#4**：`panel_host_windows.go:304` 的 `420` 一个字节没动、`NewPanelManager` 签名没动、装配根没加新参数——
那一格要"改 width ⇒ 真窗口宽度随之变"，属只有本机可量那一族，编排者另排；
⛔ `manager.go:294-300` 的同源守卫没删没弱化。

三把名册自身的尺（`cmd/wisp/config_receipt_255_test.go`）：
1. `TestTicket255HotRowRosterCoversTheRegistry`——`TierRegistry` 每枚 hot 行必须有裁决（新增段没裁决＝红）、
   名册不许留非 hot 行（漂移＝红）、裁决必须带 `file.go:LINE [token]` 或"扫描零命中"标记，
   `consumed:` 那类**必须**带 cite（不许用扫描零命中糊过去）。
2. `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`——名册里每一处 `file:line [token]` **回读源码**：
   文件要能读、行号在范围内、那一行现在必须仍含那个 token；`citesChecked` 下界 8，防止集体退化成散文。
3. `TestTicket255RosterStillMatchesTheActualReadSites`——反射走查 `config.Config` 的 14 枚 hot 段字段，
   实测读者文件集合与 `sectionReadSites` **逐字相等**（少一枚多一枚都红）。
   ⇒ "无读者段被判成有读者"与"读者出现了而名册还说没有"两个方向都会响。

两枚 stdout 判据 + 一枚装配钉：
- 正控 `TestTicket255ReceiptOmitsPanelFromTheImmediateSentence`：进程活着时手改 `[panel] width 640→641` ⇒
  台账 `hot=[panel]`、活配置确实 641（诚实句不是在说"什么都没发生"）、**逐行**扫 stdout：含"这些段已立即生效"的行不得含 `[panel]`、
  诚实句必须含 `[panel]`、台账 `section=panel` 必须含 `no-reader` 与 `panel_host_windows.go:304`、卡片数 0。
- 反控 `TestTicket255ReceiptStillNamesTheLiveReadSection`：手改 `[llm]` provider 的 `api_key_ref` ⇒
  "这些段已立即生效"那行必须仍含 `[llm]` 且不含 `[panel]`/`[ball]`、诚实句不得含 `[llm]`；
  当场调 `r.rt.settings.ReadSettings(ctx)` 拿到**新值**、同时 `r.rt.cfg`（启动快照）仍是旧值 ⇒ "有活读者"是读数不是散文。
- 装配钉 `TestTicket255ReceiptSentenceAssemblyIsFiltered`：把 14 枚热档名喂给生产函数 `reportReload`，
  钉立即行只许 `[llm]`、诚实行须逐枚点到其余各段、台账须列出 `[voice]` 背后的四枚键级行。
  它存在的理由已写进注释：退回 pre-255 形状时**判据那一行本身就红**，而不是靠"缺句"的 40s 超时兜住（§4 m1 的实测正是走了超时那条，见下）。

### 2.2 AC#3 残余补尺：`[app]` 新增键

落点＝`internal/config/tiers_app_255r2_test.go`（新增；未改 `tiers_255_test.go` 任何一行）。

- `TestEveryAppKeyIsRegisteredInTierRegistry`＝主尺：反射走查 `AppSection` 每枚叶子键，`app.<leaf>` 必须在
  `TierRegistry` 有行，档位必须是 hot/reload/restart 之一；**反向**也查：名册里每枚 `app.*` 行必须仍对应一枚真实键
  （删键留行＝红）。`AppSection.Portable` 带 `toml:"-"`（schema.go 自陈"never by config.toml"），两边都不许出现，注释里具名。
- `TestAppKeyWalkSeesTheRosterItClaims`＝防尺本身被改窄：走查到的叶数 ≥ 名册 app 行数 ≥ 4。
  **刻意是下界不是等式**：等式会让"补登记后不响"这条正控失败（登记完还得回去改一个数字），而票面要的正是"补登记后不响"。
  等式那把保险在段级已有（`tiers_255_test.go` 的 `TestRegistryCoversSchemaSectionsGreen`），本腿没动。

### 2.3 第三格（派单没写、台账写了）：`restartTierKeys` 的幻影注释

台账 `A549`（commit `69c9094c`，编排者收 `255-v1` 时裁）逐字：
「判不动五条我裁两条＝… **restartTierKeys 幻影注释归 `255-r2` 清**」。票面 §1.1 的 T2 行也早记着同一件事：
`cmd/wisp/config_reload.go` 那句注释自称 "and by a test"，而 `grep -rn restartTierKeys --include=*_test.go cmd internal`＝**0 命中**。

本腿起手复量（`2026-10-03`，还原后仍成立）：

```
$ grep -rn --include=*.go "restartTierKeys" cmd internal tools scripts
cmd/wisp/config_reload.go:320 / :326 / :329 / :332      ← 四枚全在同一文件，测试侧零枚
```

**清法＝把那句谎变成读得到**（⛔ 不是把注释里 "and by a test" 删掉了事——那是摘尺）：
新增 `cmd/wisp/restart_tier_keys_255r2_test.go` 的 `TestTicket255RestartTierKeysAreBackedByATest` 做四件事：

1. `restartTierKeys` 逐名钉死（`app.language app.autostart app.single_instance`）；
2. 每枚键必须在 `config.TierRegistry` 里登记成 `restart`（第二枚 `TierOf` 生产消费点，T2 词表与 ⓑ 登记表不许各说一套）；
3. **行为形**：逐枚只改那一枚键，走真的 `config.SaveFile → NewManager → CheckAndReload`，
   要求落进 `Report.Restart` 且**不落进** `Report.Hot`；
4. 反面正控：`app.theme` 单独改 ⇒ 必须落 `Report.Hot`、不落 `Report.Restart`。

同批把注释里那枚不存在的符号 `applyApp` 改成真身 `planApp`（`grep -rn applyApp cmd internal`＝除该行本身外零命中），
并把新读者的名字写进注释。产品行为一字未改（该文件只有注释与名册文字变化）。

### 2.4 我动了别人的一张既有断言（具名，不藏）

`cmd/wisp/config_reload_223_test.go` 的 `TestTicket223RunArmsTheReloadTick` 原三行：

```go
		if !strings.Contains(r.h.out.String(), "这些段已立即生效") {
			t.Errorf("no user-visible line names the immediately-effective tier; stdout:\n%s",
				r.h.out.String())
		}
```

它要求**一次 `[ball] size` 手改**印出"这些段已立即生效"——`ball` 是零读者段，这条断言钉住的正是本票要禁的那句谎。
改写为**双向**（不是 OR、不是放宽）：诚实句必须出现且必须点名 `[ball]`；且 stdout 里"这些段已立即生效"必须**不出现**。
该用例其余断言（`hot=[ball]` 台账、活配置 64、启动快照未动、卡片数 0）**逐字未动**；
`config_reload_223_test.go:495` 那枚反向钉（restart 档不许冒充立即档）也未动、仍绿。

---

## 3. 门禁四数（真实读数，还原突变后的终值）

PATH 前置照派单：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`；构建用 `GOFLAGS= go build ./...`。

| # | 命令 | 终值读数 |
|---|---|---|
| 1 | `GOFLAGS= go build ./...` | `BUILD_RC=0`（无输出） |
| 2 | `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp internal/config` | **空**（`GOFUMPT_LINES_ABOVE_RC=0`，零行输出） |
| 3 | `GOFLAGS= go test ./internal/config/ ./cmd/wisp/ -count=1` | `ok github.com/CarlosShao/wisp/internal/config 1.306s` ／ `ok github.com/CarlosShao/wisp/cmd/wisp 311.859s` |
| 4 | `./tools/d22scan/d22scan.exe`（独立模块，只跑已构建好的 exe；`tools/d22scan/main.go` mtime 09-26 23:09 早于 exe 23:33＝非过期构建） | `D22_RC=0`；末行逐字：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=227, bans #1-5 cmd/=37, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=489, ban #8 cmd/=89; ...` |

- 门禁 #4 首行另有一句 `d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`——
  那是它内建的 git-ignore 处理（`frontend/.gitignore` 决定），**起手锚 HEAD 上同形**，不是我这一格引入的 SKIP；
  它不是"某个扫描范围被跳过"那一族的 SKIP。
- 本机没有可用 staticcheck（会产"干净的绿"），⛔ 未跑。
- 两包 `-count=1` 之外**没跑过全仓 `go test ./...`**（同机并发会互相洗读数）。
- 全仓跑批里我撞到两枚**非本腿代码**的偶发红，两次都换成了不同的用例，单跑各 3/3 与 2/2 全绿；具名与判读见 §6 第 1 条。
- ⚠ 终值之前还有一发是 `FAIL github.com/CarlosShao/wisp/cmd/wisp [build failed]`——红因**不在我这两枚包里**：
  同一时刻 `git status --porcelain -- cmd internal` 显示他腿在飞的 `M internal/agent/approval/gate.go`、
  `M internal/panel/pump.go`、`M internal/panel/subagent_roster_197.go`（`cmd/wisp` 链接这两包），
  而我写的文件在 `GOFLAGS= go vet ./cmd/wisp/` 下零错误；他腿收口后重跑即 `ok github.com/CarlosShao/wisp/cmd/wisp 311.859s`。
  读数与形状记在这儿，不去猜别人的码写完没有。
- 不设 PATH 的形状：派单说会以 `exit status 0xc0000135`、无 `--- FAIL` 出现。我用"只去掉 sherpa 目录"的 PATH 复现时拿到的是
  另一种非断言形状：`build constraints exclude all Go files in .../sherpa-onnx-go-windows@v1.13.8` ＋ `[setup failed]`。
  **共同点才是判据**：两者都不是代码断言失败；判绿只认 `--- FAIL`/`--- PASS` 行。

---

## 4. 变异自证表（种什么形 → 哪枚必须红 → 复跑终值）

备份/还原式：突变前 `git show HEAD:<path> > .scratch/wisp/probes/255/r2/backup/<名>.before` 并记 md5；
还原用 `git cat-file blob HEAD:<path> > <path>`；还原后 `md5sum` 与备份逐字相等；
每发之间量 `git status --porcelain -- cmd internal`。
（期间他腿在共享树里留下的只有 `?? internal/panel/inbound_roster_253_test.go` 与 `M internal/tools/paths.go`，均非我写面。）

| 发 | 种的形 | 备份 md5（种前） | 必须红的用例 | 红句（逐字） | 还原后 |
|---|---|---|---|---|---|
| **m1**（AC#1 正控） | `cmd/wisp/config_reload.go` 的 `reportReload` 退回 pre-255 形状：`if len(rep.Hot) > 0 { Fprintf(...，rep.Hot) }`（split/台账/诚实句/`disagree` 整段撤掉） | `config_reload.m1.before`＝`a7094e7b91ee75d2ab57fcb094ac5218` | `TestTicket255ReceiptOmitsPanelFromTheImmediateSentence`、`TestTicket255ReceiptSentenceAssemblyIsFiltered`、`TestTicket223RunArmsTheReloadTick` 三枚全红（0.00s 判据之外，产品流里实际印出的谎是：`wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[panel]`） | `config_receipt_255_test.go:402: stdout never carried "值已换进本进程内存" within 40s; full stdout:` ／`--- FAIL: TestTicket255ReceiptOmitsPanelFromTheImmediateSentence (42.63s)`；`--- FAIL: TestTicket255ReceiptSentenceAssemblyIsFiltered (41.51s)`；`config_reload_223_test.go:280: stdout never carried "值已换进本进程内存" within 40s`＋`--- FAIL: TestTicket223RunArmsTheReloadTick (42.85s)`；dump 里逐字含 `wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[panel]` | md5 复量＝`a7094e7b91ee75d2ab57fcb094ac5218`（与备份一致）；porcelain 我写面 0 行 |
| **m1b**（AC#1"无读者段被判成有读者"） | `cmd/wisp/config_readers_255.go` 里 `hotRowClaims["panel"]` 的前缀换成 `hotClaimConsumed`（cite 仍是真行，所以证据回读尺**不会**响——响的是判据） | `config_readers.m1b.before`＝`d798e1057d57b6534772caeee104e5d3` | `TestTicket255SplitOnlyClaimsSectionsWithALiveReader`（0.00s，纯判据）＋`TestTicket255ReceiptSentenceAssemblyIsFiltered`（1.19s，stdout 装配）＋正控用例（43.37s） | `config_receipt_255_test.go:76: claimable = [llm panel], want exactly [llm]: it is the only hot section this host reads live (cmd/wisp/panel_config_store.go:92 takes s.mgr.Config() per call). A section with no live reader must never be in this list.` ／`config_receipt_255_test.go:522: ticket 255 AC#1: the 已立即生效 sentence claims [panel], which has no live reader in this host: "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[llm] [panel]"` ／`config_receipt_255_test.go:528: the honest sentence omits [panel], so the edit vanished from the receipt: "..."` | md5 复量＝`d798e1057d57b6534772caeee104e5d3`；porcelain 我写面 0 行 |
| **m2a**（AC#3 补尺） | `internal/config/schema.go` 的 `AppSection` 种入假新键 `PortraitMode bool \`toml:"portrait_mode" default:"false"\``，**不登记** | `schema.m2.before`＝`4edbe60979ef0b3a234616a16a10a2a0` | `TestEveryAppKeyIsRegisteredInTierRegistry` 必须红；同时三枚既有尺（`TestEverySectionHasATierInRegistry`／`TestPerKeySectionsAreRegisteredPerKey`／`TestVoicePerKeyRowsMirrorPlanVoiceWalks`）实测**全绿**＝255-v1 披露的射程缺口成立 | `tiers_app_255r2_test.go:49: app key "app.portrait_mode" (AppSection.PortraitMode) has no tier row in TierRegistry: ticket 255 AC#3's residual - a NEW [app] key is invisible to the fixed-name list, so it must be caught by the walk (register it in tiers.go, or delete the key)` ＋`--- FAIL: TestEveryAppKeyIsRegisteredInTierRegistry (0.00s)`；`--- PASS: TestEverySectionHasATierInRegistry`／`--- PASS: TestPerKeySectionsAreRegisteredPerKey`／`--- PASS: TestVoicePerKeyRowsMirrorPlanVoiceWalks` | — |
| **m2b**（AC#3 正控） | 在 m2a 之上补登记：`internal/config/tiers.go` 加 `"app.portrait_mode": "restart",` | `tiers.go` 未预备份（流程失手，见 §5 第 6 条），还原用 `git cat-file blob HEAD:` 并当场留底 `backup/tiers.go.head-restored` | 两枚 app 尺都必须**不响** | `--- PASS: TestEveryAppKeyIsRegisteredInTierRegistry (0.00s)`／`--- PASS: TestAppKeyWalkSeesTheRosterItClaims (0.00s)`／`ok github.com/CarlosShao/wisp/internal/config 0.033s` | `internal/config/tiers.go`＝`bb077e6269ab54ba829004f2d0229b97`＝`backup/tiers.go.head-restored`；`internal/config/schema.go`＝`4edbe60979ef0b3a234616a16a10a2a0`＝`schema.m2.before`；porcelain 我写面 0 行 |

复跑终值（全部还原、格式化之后）：门禁四数见 §3；本腿点名的七枚用例一次跑齐
`TestTicket255Split…`／`TestTicket255HotRow…`／`TestTicket255RosterEvidence…`／`TestTicket255RosterStillMatches…`／
`TestTicket255ReceiptOmitsPanel…`／`TestTicket255ReceiptStillNames…`／`TestTicket255ReceiptSentence…` ＋ `TestTicket223RunArmsTheReloadTick`
＝`--- PASS` ×8、`ok cmd/wisp 9.198s`（该发只跑这 8 枚）；internal/config 侧
`TestEveryAppKeyIsRegisteredInTierRegistry`／`TestAppKeyWalkSeesTheRosterItClaims` ＝ `--- PASS` ×2。

m1 的诚实自我批评：种回 pre-255 形状时，三枚用例都是**先**在"缺诚实句"的 40s await 上红的（`awaitStdout` 的超时把整段 stdout
连谎句一起 dump 出来，谎句 `...：[panel]` 在 dump 里逐字可见）。这正是我给 `TestTicket255ReceiptSentenceAssemblyIsFiltered`
写注释时说的形状——判据那一行本身没有先响。m1b 才是判据直响（`claimable = [llm panel]`、`the 已立即生效 sentence claims [panel]`）。
两发都红＝两种"把谎放回产品里"的路都被堵；但"装配钉能在判据行上先红"这件事目前**只对 m1b 成立**，已记进 §5 第 2 条。

---

## 5. 我可能写错的条目（自我对抗，逐条可反证）

1. **判据口径可能选严了**：我按 (B)"读者会在新值上做事"裁决（§1.5），于是 `agent`（装配快照读者）、`models`（别的命令进程读盘）
   被放进诚实句那一半。若编排者按 180-a1 的 R 类（(A)"生产代码读过它的值"）来验反控，
   这两枚应当出现在"已立即生效"那一句里 ⇒ 需要改的只是 `hotRowClaims` 两枚前缀＋各自 cite。
   **风险方向**：我的形不会说谎，但可能让反控看起来"少列了两枚"。
2. **m1 的紅因不是我想要的那个**：种回全集时三枚用例都先在 40s await 红（见 §4 末段）。若验收要求"判据行本身先响"，
   装配钉对 m1 这一形没做到——它只对 m1b 做到。补法（一行）：把装配钉里的 `awaitStdout` 换成
   `awaitStdout(t, "wisp run: 配置热加载：")` 再扫两句，但那样会引入"两句只到其一"的新窗口，我没在剩下轮次里把它做严。
3. **`hotRowsFor` 对"段级不命中就收键级行"的推断**：今天只有 `app`/`voice` 是键级段，靠的是 `TierRegistry` 的前缀形状。
   若将来出现一枚段既登记了段级行又登记了键级行，我的分支会走段级行、忽略键级行 ⇒ 判定可能偏松（`plan()` 的守卫只挡段级不一致）。
   这条我今天**量不到**（没有这种段），也没写尺。
4. **`consumed` 要求"每一枚背后行都是 consumed"**：偏保守方向（宁少说不多说），但它意味着——
   若将来 `voice` 的 4 枚热键里只有 3 枚接上线，回执会整段沉默地走诚实句。这是刻意的，票面没裁反方向。
5. **扫描尺的命中形状是文本级**（`\.Field\.[A-Z]` / `Config().Field` / `Config.Field`）。它看不见
   `x := cfg.Panel` 之后再读 `x.Width` 这种两跳形状，也看不见反射/`toml` passthrough。
   180-a1 对它自己的普查也声明了同类边界（U3 裸字面量未扫）。⇒ `no-reader` 那 11 枚的"零"是**这套形状的零**，不是绝对零。
6. **流程失手一处**：m2b 我改了 `internal/config/tiers.go` 而**没有先做 pre-mutation 备份**（只在该文件还原后留了
   `backup/tiers.go.head-restored` 并比对 md5）。还原本身可证（与 HEAD 逐字节同），但顺序不对，具名认。
7. **`sectionReadSites` 的排除清单**：排掉了 `internal/config/`、`*_test.go`、`testdata/`、注释行与本票名册文件。
   排掉名册文件是因为它的裁决字符串**引用**代码（如 `cfg.Models.Mirror`）。如果有人往这个文件里写真读配置的代码，
   这条排除会把它藏住——目前它只做字符串比对与反射，且排除是按精确文件名（不是模式）。
8. **我改了他票的一枚断言**（§2.4）。若裁"任何既有断言都不许动"，这一处要退回；退回去的结果是
   `TestTicket223RunArmsTheReloadTick` 在 AC#1 修复下必红（它钉的正是被禁的谎）。这条冲突需要编排者裁，不是我裁。
9. **`bracketed()` 改了产品文案的形状**（`[ball session]` ⇒ `[ball] [session]`）。若有我没找到的下游/测试按 `%v` 形状解析那句，
   它会红。我在 `cmd internal tools scripts docs` 里搜过"这些段已立即生效"，除我改的两枚之外只有 `config_reload_223_test.go:495`
   那枚反向钉（不匹配即绿，形状改动不影响）。

---

## 6. 判不动的地方（量不到就具名写"量不到"）

1. **同机并发下的偶发红，本腿无法归因**。全仓两包跑批三次里，两次各有一枚**非本腿代码**的用例红：
   - `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`（`config_reload_223_test.go:454` "the card does not name risk.permission_mode"）：
     单跑 3/3 `ok`（2.627s / 3.005s / 2.341s）。该路径上 `rep.Hot` 为空（[risk] 走 `planLocked`），我的代码**不产出任何一句**，
     且那枚断言是"awaitCard 之后单次读 stdout"的形状——正是 223 r2 在别的用例里用 `awaitStdout` 修掉的争用形状（该文件 `:475-483` 注释具名）。
     ⛔ 我没去修它：那是他票的钉，且修它需要跑 `cmd/wisp` 终裁程串行（台账 `A##` 里编排者已排过同族）。
   - `TestAC4FocusReturnToPriorWindowGap33r5`：单跑 2/2 `ok`（1.342s / 1.124s）。它是真窗口焦点用例（只有本机可量那一族），
     争用敏感；台账 `docs/reports/pending-and-issues.md:5708`/`:5743` 已具名"抢 CPU 会造出那一族偶发红"。
   ⇒ 终值那一发（§3 第 3 行）两包全绿。**判不动的正是"偶发与我的改动有没有关系"这件事**：
   我能证的是两次红的用例不同、单跑都绿、且其中一条走的代码路径我的改动不产出输出。
2. **AC#4 今天判不动**（不是本腿的格）：`[panel] width` 真生效要"改 width ⇒ 真窗口宽度随之变"，只有本机可量，编排者另排。
   本腿只让回执不说谎，没接线。
3. **`[app]` 新键"登记了但 planApp 不处理"这一形今天量不到**：m2b 里我把假键登记成 `restart`，
   而 `planApp`（`manager.go:348-367`）的条件式里根本没有它 ⇒ 它会既不进 `rep.Hot` 也不进 `rep.Restart`，
   即"登记了档位、实现里却无人分档"。我的尺**不响这一形**（它只管登记与键存在性）。
   能响它的尺要把判定做成行为形（逐键改一次、看落进哪条切片）——那是"值不值一枚尺"的裁量，交编排者，不在本腿顺手做。
4. **`OnReload` 出货进程没听众**这一格归票 42（票面禁区节明文不许并进本票），我没量。
5. **不设 PATH 的 `0xc0000135` 形状我没复现到**（实测到的是 `[setup failed]` 那一种，见 §3 末），
   共同结论仍成立：两种都不是断言失败。要一条干净的"0xc0000135 且无 `--- FAIL`"读数，本腿**量不到**。
6. **面板侧（TS/界面）那一半不判**：票面与 `AGENTS §1.2` 都写明 `frontend/**`／`design/**` 既不读也不引，
   所以"读侧没有出口"这句只到 Go 侧为止。

---

## 7. 交件判语

- **AC 勾选框一枚未碰**：票面 `.scratch/wisp/issues/255-...-tier-has-zero-readers.md` 的 `- [ ] AC#1` /
  `- [ ] AC#4` 两框一字未改（`AC#2`、`AC#3` 两框是编排者已翻勾的 `[x]`，我也没动）；
  翻勾归编排者凭**非实现者**验收表来做，本腿不自裁。
- **未 push**：本腿 5 枚 commit（`c8480bce`／`ae60a87c`／`a4906b6c`／`17ff058d`／本件终值 commit）全部只 `git commit`
  带显式 pathspec，**一次 `git push` 都没跑**；推送是编排者的动作。
- 写面只在 `cmd/wisp/`＋`internal/config/`；`internal/panel/`、`internal/agent/**`、`internal/ball/`、`internal/tools/`、
  `internal/risk/` 一字未写；`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、
  `tokens_fourway_test.go`、`ticket90_persist_test.go`、`.github/workflows/ci.yml` 一字未动。
- 未放宽任何既有断言、无 `t.Skip`、无把 SKIP 读成通过；仓内**没有删任何文件**（备份、日志、msg 全部只建不删）。
- 请验收者**特别核**三处：①§1.5 的口径选择（(A)/(B)）；②§2.4 我改了 223 的一枚断言；③§4 m1 的红因是"缺句超时 + dump 里的谎句"，
  而 m1b 才是判据行直响。
