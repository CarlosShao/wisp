# 票 255 · 落地腿 `255-r8` · 锚＋现量＋门禁 rc（本腿只做 AC#1 ＋ AC#4 两格）

身份＝实现者（判语归非实现者；⛔ 未翻框、⛔ 未动 `Status:`/票名）。件内一切行号**现读于本件所记 HEAD**。

---

## 0. 起手锚与读数时刻（全本腿自跑）

| 时刻 | HEAD | 读数 |
|---|---|---|
| 17:17 | `f6fda72d` | 派发锚 `7e8ac3fa` 已漂：其间 `318-r1` 两笔（`internal/agent/approval`，不在本面） |
| 17:18 | `0d45d6f4` | AC#1 名册全尺 `-run TestTicket255` rc=0（窗口 17:18:05→17:18:44） |
| 17:22 | `93452d45` | 真窗读数①（重建路径，r3 既有尺）PASS（窗口 17:22:25→17:22:37） |
| 17:37 | `605838b0` | 真窗读数②（甲形 reshow，本腿新尺）PASS（窗口 17:37:00→17:37:15） |
| 17:42 | `c3d1fe4f` | 门禁五连＋AC#1 名册复跑 rc=0 |

**同源凭据**：`git log --oneline f6fda72d..HEAD -- cmd/wisp internal/config` ＝**空** ⇒ 上述 HEAD 之间本面（两 AC 的全部产码/仪器）**逐字节未动**，五发读数同源。
写面起手/收手态：`git status --porcelain cmd/wisp internal/config scripts` 起手空、其间只多本腿新尺一枚（见 §4）。
别人的在飞脏件（⛔ 未动）：`go.mod` 工作树有 `M`（go-webview2 indirect→direct，非本腿造，未 commit 它）。

## 1. 结论先行

- **AC#1＝现跑复认成立，本腿零产码改动**：落地（分拣回执＋名册）是 `255-r2` 的，现盘上句子的两半逐字见 §2；正控（手改 `[panel] width` ⇒ `已立即生效` 句里无 `[panel]`）与反控（真有活读者的 `[llm]` 仍在）**现跑双绿**。两格一致性按派单「自然成立」那支成立、**未翻转**（§3 末）。
- **AC#4＝功能已在盘（r3 建窗链＋r1 甲形 reshow），本腿补上它自述缺的那一发**：真窗实测**做了**，且是**两发**——重建路径 `640→900`（r3 尺，现跑复核）＋**甲形 reshow（活窗不销毁）**：同 HWND、客户区 `624→900`、外框 `640→916`、**帧差 16px 全程不变**（§3）。甲形这一发此前**没有任何真窗读数**（r1 自己文件逐字写着 "no case here pretends to have one"），现补上并 PASS。
- ⛔ 未 push、⛔ 未翻框、⛔ 未动 `docs/PLAN.md`／`docs/specs/**`／thresholds/golden/allowlist。

## 2. AC#1 —— 那句"已立即生效"只许说真话（现跑）

**落地出处（现读）**：`cmd/wisp/config_reload.go:182 reportReload` 用 `splitHotTier`（`cmd/wisp/config_readers_255.go:295`）按名册（`hotRowClaims`，同件 `:100`）把 `rep.Hot` 分成 `claimable/quiet`；`claimable` 进 `已立即生效` 句（`:191`），其余进诚实句（`:194-197`），第三类 `disagree` 单独点名。名册逐行判词另走 `HOT-RELOAD-READER` 审计行（`:188`）。

**逐字（产品的两句话，从生产函数输出、非手抄，`-v` 落 log；原文均为单行，此处不折）**：

```
IMMEDIATE "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[llm]"
HONEST    "wisp run: 配置热加载：这些段的值已换进本进程内存，但本宿主没有会按新值做事的读者，本次运行不会因此改变行为（票 255 AC#1：这一半不许说成「已立即生效」；逐段的读者判定见 HOT-RELOAD-READER 行）：[ball] [session] [audio] [agent] [privacy] [memory] [panel] [cost] [models] [observe] [hotkey] [app] [voice]"
```

**正控（`[panel] width`，真进程 `wisp run` 合成根＋真 tick＋手改）**：`TestTicket255ReceiptOmitsPanelFromTheImmediateSentence` **PASS**——种 `width = 641` 后：任何含 `已立即生效` 的行不含 `[panel]`（`:442`）；诚实句必须点名 `[panel]` 且含 `没有会按新值做事的读者`／`本次运行不会因此改变行为`（`:447-454`）；内存真的换了（`mgr.Config().Panel.Width=641`，`:456`）；审计行 `section=panel` 判 `other-process` 且 cite `panel_resident_windows.go` ＋ `[cfg.Panel.Width]`（`:461-476`）；0 张卡。
**正控二形（"哪儿都没读者"那半，`[session]`）**：`TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere` **PASS**（`warm_timeout_sec=91` ⇒ 句子不点名 `[session]`、判词 `no-reader`）。
**反控（真有活读者仍要出现）**：`TestTicket255ReceiptStillNamesTheLiveReadSection` **PASS**——`[llm]` 改 `api_key_ref` 后：`已立即生效` 行**必须含 `[llm]`**（`:550`，缺＝摘尺判红）、不得含 `[panel]`/`[ball]`；活腿 `settings.ReadSettings` 当场读到新 ref 而启动快照仍是旧 ref（`:560-573`，证明"有人按新值做事"而非"内存换了"）；审计行判 `consumed:`。
**句子装配（全 14 枚 hot 名一起）**：`TestTicket255ReceiptSentenceAssemblyIsFiltered` **PASS**（上面 IMMEDIATE/HONEST 逐字即它打的）。

名册全尺：`go test ./cmd/wisp/ -run TestTicket255 -count=1` 两发（17:18、17:42）**均 rc=0、20 顶层 PASS / 0 FAIL**（`logs/t255-verbose.txt`，含上列逐字与 LEDGER 行；本腿**零改动**本件任何一行、未放宽任何断言——⛔ 也没有需要改的缺口）。

## 3. AC#4 —— `[panel] width` 真生效（装配链逐跳＋真窗两发）

**装配链逐跳（全现读，`logs/chain-readings.txt`）**：

| 跳 | 处（现读） | 现量 |
|---|---|---|
| 0 键 | `internal/config/schema.go:542` | `Width int \`toml:"width" default:"640"\``（写死 420 已不在唯一来源） |
| 1 装配根 | `cmd/wisp/panel_resident_windows.go:198-209` `panelGeometrySource` | **每次建窗**现跑 `config.LoadFile(cfgPath)` ⇒ `cfg.Panel.Width, cfg.Panel.Height`；失败响亮回 `0,0` |
| 2 递入 | 同件 `:273` | `NewPanelManager(..., withGeometrySource(panelGeometrySource(dataDir)))`（唯一生产装配点） |
| 3a 建窗 | `cmd/wisp/panel_host_windows.go:236 windowOptions` ⇒ `:262 Width: uint(width)` ⇒ `:391 newPanelControl(m.windowOptions(), m.dataPath)` | 每次建窗**恰读源一次**；`geometry==nil` 时回宿主常量（`420×260` 只剩"配置缺省时的值"身份） |
| 3b 甲形 | 同件 `:532 m.requestGeometryOnReshow()`（`Show` 已建窗支）⇒ `:610` 方法 ⇒ `:627 w.SetSize(width, height, webview2.HintNone)`（走 `Dispatch`） | 活窗**不销毁**拿新值；⌐建窗＝外框语义 vs `SetSize`＝客户区语义（`panel_control_windows.go:637/654/655` AdjustWindowRect＋SetWindowPos）并列具名 |
| 4 泵 | `cmd/wisp/panel_resident_windows.go:319 w.Run()` | Dispatch 队列的**唯一**读者（测试用生产 `startResidentPanel`，同一拓扑） |
| 边 | `TestTicket255HostStillDoesNotParseConfigItself` PASS | 宿主件 import 块无 `internal/config`、注释外零 `config.toml`/`configFileName`/`config.LoadFile` ⇒ ⛔ 未新开 `panel→config` 依赖边、宿主⛔ 不自己读盘 |

**真窗读数①（重建路径，r3 既有尺，现跑复核 @`93452d45`）**：`asked 640x400 -> GetWindowRect 640x400；after dispose+show asked 900x400 -> 900x400；implied scale 1.0000；geometry source read 2 time(s) for 2 creates`＋`TEARDOWN NAIL: windows=0`，**PASS**（`logs/winlive-recreate.txt`）。
**真窗读数②（甲形 reshow，本腿新尺 @`605838b0`，逐字）**：

```
REAL WINDOW 255r8: config 640x400 -> create hop: outer 640x400 / client 624x361 (frame delta 16 px)
hwnd=0x2601144; config rewritten to 900x400, hide+show with no destroy -> outer 916x439 /
client 900x400 (frame delta 16 px) hwnd=0x2601144 (same hwnd: true); implied scale 1.0000;
geometry source read 2 time(s)
--- PASS: TestTicket255r8RealWindowFollowsTheConfigAcrossAReshow
```

⇒ 改 `[panel] width` 640→900 后，**同一扇活窗**（hide+show，不销毁）客户区**恰好 900**、外框 916、**帧差 16px 两发恒等**——这就是 "改 width ⇒ 真窗口宽度随之变" 在 10-08 裁定那个形（"面板重新显示时取一次并应用"，甲形）上的**本机实测**。
**两格一致性（未翻转，具名）**：保存那一刻窗口**不会**自己动——轮询 tick 不在常驻腿（`tiers.go:34` 把 `[panel]` 登记 `hot` 而 `OnReload` 只对 Reload 档开火；`startConfigReload` 唯一非调用者在 `run.go`），新尺寸在**下一次显示请求**到达 ⇒ `AC#1` 正控**照原支**成立（`[panel]` 不进"已立即生效"句），无需按"确实生效了"重写。节内两跳语义差（外框 vs 客户区）**没有**被任何文案说成"屏幕宽度＝`[panel] width`"。
**落地途中两发红，如实登记**（都是本腿自己的件，已修、逐字留档）：① 第一发红在**我断言的口径**（`client 增长 276 vs 我误写 260`——276＝260＋16 帧差，是正确答案；读数本身自始正确）；② 第二发红在 `t.TempDir()` RemoveAll 与 WebView2 profile 收尾竞态（`EBWebView\Default: The directory is not empty`，线程已 clean exit）⇒ profile 挪系统 temp（家族定式）后绿。

## 4. 产码（与本件分开的一笔 commit：`c88356c1`）

`cmd/wisp/panel_reshow_255r8_winlive_test.go`（`//go:build windows && winlive`，新立；查重：起飞前 `-tags winlive` 家族里**无**任何 reshow/SetSize 真窗尺）。它把 r1 自述缺的那一发量掉：真 config.toml → `panelGeometrySource`（生产闭包，只外包一层计数）→ 生产 `startResidentPanel` 线程＋`Run()` 泵 → 生产请求面 `RequestToggle`(hide)＋`RequestShow`(reshow)，读 `GetClientRect`/`GetWindowRect` 两发比。⛔ 未碰任何产码、⛔ 未碰 r1/r3/33 族既有件、⛔ 未动断言。

## 5. 门禁 rc 表（每件自落 rc，无管道吞退码；`logs/gates.txt`）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `go build ./...` | **0** | — |
| vet(winlive) | `go vet -tags winlive ./cmd/wisp/` | **0** | CI 的 winlive 类型检查步 |
| 格式 | `gofumpt -l cmd/wisp/panel_reshow_255r8_winlive_test.go` | **0** | 空输出 |
| 相关族① | `go test ./cmd/wisp/ -run TestTicket255 -count=1` | **0** | ok 25.5s；20 PASS/0 FAIL（-v 见 `logs/t255-verbose.txt`） |
| 相关族② | `go test ./internal/config/ -count=1` | **0** | ok 3.6s |
| d22scan | `sh scripts/d22scan.sh` | **0** | clean（正控 PASS=35 FAIL=0 SKIP=0） |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | VERDICT GREEN，over-budget=60 / roster=60 / not in roster=0（60=60） |

## 6. 具名没做（留编排者/机主）

- ⛔ **未跑整包 `go test ./cmd/wisp/`**（314 死区禁令）；本腿只跑 `-run TestTicket255` 与两发 winlive 单跑并报窗口时刻（§0）。
- **DPI 逻辑→物理那层**：仓里既无实现也无仪器（`A700` §2 已登记）；两发读数都发生在 implied scale **1.0000** 的这台机上，断言按第一发反推的 scale 表述，⛔ 没替它猜规则。
- **"保存那一刻自动变"**：没做，也不在本票裁定射程（排在热加载那一跳之后）；文案与注释通篇未宣称。
- **`[panel] height=0=auto`**：未实现的具名登记未动（`panel_geometry_255_test.go` 两半都在钉）。
- **真窗实测**：**做了**（§3 两发）；`SetSize` 客户区语义的**帧差真值 16px** 现已落读数（此前只有 a2 的"仅本机可量"）。
- 本件未写 `impl.md`（本腿活＝复认＋补一发真窗尺，§1/§3 即结论）；临时中间件（首/次红发）留 `logs/`，只建不删。
