# 票 33 面板宿主 —— 对抗验收表 `33-v2`（非实现者终裁）

**本腿代号**：`33-v2`（对抗验收腿，非实现者；产码与测试文件一字未改，见 §9 自证）
**起手锚点（本腿现取，`2026-10-01 16:47:37 +08`）**：HEAD＝`3216ba6da5e53b5239e6dfb3d1f4ed9f3e54959e`，分支 `dev`。
**票面**：`.scratch/wisp/issues/33-panel-host-c27.md`（340 行）。本腿现读复选框＝**13 未勾／1 已勾**（`grep -c "^- \[ \]"`＝13、`grep -c "^- \[x\]"`＝1；已勾那枚是 AC#11，编排者 10-01 11:1x 翻的）。
**归属规矩**：⛔ 勾与不勾一律归编排者。本表只出判语（成立／不成立／无法判），不翻任何框、不往票面写任何东西。

---

## 0. 本腿的取数口径（钉死，引用本表任何读数都带上它）

- 整包命令（唯一口径）：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1 -v -timeout 25m`。
- 红绿只认 `--- FAIL`／`--- PASS`／`--- SKIP` 行（`t.Logf` 也带 `file:line:` 前缀，不当失败）。
- 前人读数一律具名，且注明本腿是否复跑。⛔ 本腿不复跑前人同格读数当凭据（派单 §0），但为差集归因所必需的那一发除外，且会写明。
- `staticcheck` 本机版与 CI 钉版不同 ⇒ 本表凡涉及它一律标〔未复认〕，不拿它出结论。
- 三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）本表若引用，一律标「这是射程判断，不是把冻结件内容当产品凭据」。
- `frontend/**`／`design/**`：⛔ 未读、结论未引（两层禁令）。

## 1. 逐格判语表（票面 13 枚未勾，一枚不落）

| 格 | 票面行 | 判据原文摘要 | 判语 | 凭据（判据行＋我跑了什么＋逐字读数） | 若这格是假的，最坏会被谁误信 |
|---|---|---|---|---|---|
| AC#1 | `:72-73` | Show/hide/destroy 生命周期：会话内二次 show 复用同扇窗（进程树子进程数稳定）；会话 dispose 拆窗＋WebView 子进程 ≤2s 退净 | （待填） | （待填） | （待填） |
| AC#2 | `:74` | 时延：冷 ≤1500ms／热 ≤200ms（10 发 P50/P95 进 SLO appendix） | （待填） | （待填） | （待填） |
| AC#3 | `:75` | 内嵌资源离线可给；不存在监听套接字（netstat 断言） | （待填） | （待填） | （待填） |
| AC#4 | `:76` | 焦点往返：编辑器 → 面板 → 隐藏 → 焦点回编辑器（自动＋手动） | （待填） | （待填） | （待填） |
| AC#5 | `:77` | 运行库缺失 fixture（改名/遮罩加载器）→ 降级事件＋App 存活＋原生 L2 标志开 | （待填） | （待填） | （待填） |
| AC#6 | `:78` | 所服务文档带 CSP（响应侧检），含 `connect-src 'none'` | （待填） | （待填） | （待填） |
| AC#7 | `:79-86` | `SnapshotPump.Snapshot()` 跨两瞬间的形状要有会响的检（`-race` 今天干净） | （待填） | （待填） | （待填） |
| AC#8 | `:87-94` | 两枚零执行分支（`panel_pump.go` 摘要超 440 字符丢 ids／`EvToolStart`）各要一发"拿掉就红"或用例或删 | （待填） | （待填） | （待填） |
| AC#9 | `:96-124` | `ComposerDispatch.Handle` 被非 test 文件真调用＋本机端到端留读数 | （待填） | （待填） | （待填） |
| AC#10 | `:245` | `[panel]` 的 `width`／`height`／`scale` 有生产者（真去设窗口边界，`hot` 生效） | （待填） | （待填） | （待填） |
| AC#12 | `:37` | "能构建"≠"有页面可发"：判据要区分"embed 只匹配到占位文件"与"真有一包页面产物" | （待填） | （待填） | （待填） |
| AC#13 | `:39-40` | 冷启动的往返探测不许把真页面盖掉（次序/过滤器＋会响的"最终文档含 embed 入口真内容"＋反控） | （待填） | （待填） | （待填） |
| AC#14 | `:42-44` | "Go→页面"这一跳由谁投递；一枚会响的用例＝真页面调绑定并 await，断回话真到达页面 | （待填） | （待填） | （待填） |

## 2. 门禁读数（本腿自己取，逐条带命令）

### 2.1 写面与锚点卫生

| 尺 | 命令（逐字） | 读数 |
|---|---|---|
| 起手 HEAD | `git log -1 --format=%H` | `3216ba6da5e53b5239e6dfb3d1f4ed9f3e54959e`（`2026-10-01 16:47:37 +08`），分支 `dev` |
| 产码自我锚点起未动 | `git diff --numstat 096aafad..HEAD -- cmd internal` | **空**（⇒ 本腿测的 `cmd`／`internal` 与编排者 16:37 那一发**逐字节同码**，两发读数可直接对表） |
| 工作树写面 | `git diff --stat HEAD -- cmd internal` | **空** |
| 本腿名下提交 | `git log --oneline 096aafad..HEAD` | `8fb99755`（骨架）＋后续逐节提交；`3216ba6d` 是编排者的 |
| 票面复选框 | `grep -c "^- \[ \]"`／`grep -c "^- \[x\]"` | **13／1**（与派单一致；本腿一枚未碰） |
| 台件目录未被占 | `ls docs/evidence/s1/ \| grep 33-panel-host`；`ls .scratch/wisp/probes/33/` | 本腿到达时名册＝`r1 r4 r5 r6 r7 v1`／`a1 a2 h1 p1 r1..r7`，**无 `v2`** ⇒ 代号可用 |

### 2.2 整包（本腿自跑一发，独占机器）

命令逐字：
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1 -v -timeout 25m`
起跑 `16:49:45`／终 `16:53:22`，逐字末两行：`ok  	github.com/CarlosShao/wisp/cmd/wisp	203.831s` ＋ `RC=0`。
台件＝`.scratch/wisp/probes/33/v2/fullpack-1.txt`（1,331 行）。

| 口径 | 本腿 | 编排者 16:37（`33r7-head-roster-1.txt`） | 编排者 14:07（`33r5-head-roster.txt`） |
|---|---|---|---|
| `grep -c '^--- PASS'`（**只顶层**） | **159** | 159 | **156** |
| `grep -c -- '--- PASS'`（顶层＋子项） | **240** | 240 | 237 |
| `^--- FAIL`／`-- '--- FAIL'` | 0／0 | 0／0 | 1／1 |
| `^--- SKIP` | **0** | 0 | 0 |
| 末行 | `ok ... 203.831s` rc=0 | `ok ... 238.545s` rc=0 | `FAIL ... 235.183s` rc=1 |

**⚠ 派单 §4(丁) 那句"156 涨到 240，口径＝`grep -c -- '--- PASS'`，同 `-v`"经复量＝口径混了，差集不是 84 枚而是 3 枚。**
同一份文件两把尺各量一次：`33r5-head-roster.txt` 顶层 **156**／含子项 **237**；`33r7-head-roster-1.txt` 顶层 **159**／含子项 **240**。
⇒ 156 那发用的是 `^--- PASS`（顶层），240 用的是含子项那一把。**涨的真实枚数＝顶层 156→159（＋3），含子项 237→240（＋3）**——不是 84 枚新用例，也不是有用例被洗掉。
（这一条不影响"逐名对过"的必要性，只影响差集该有几行；逐名见 §2.3。）

### 2.3 名册逐名差集（本腿自建的三枚名册，同一命令口径）

尺：`grep '^--- PASS\|^--- FAIL\|^--- SKIP' <log> | sed 's/^--- \([A-Z]*\): \([A-Za-z0-9_]*\).*/\2/' | sort -u`
台件：`.scratch/wisp/probes/33/v2/rosters/top-33r5.txt`（157 名）／`top-33r7.txt`（159 名）／`top-mine.txt`（159 名）。

- `diff top-33r5 top-33r7` ＝**恰两行新增、零行删除**：
  - `TestAC13BringUpSurvivesAReusedThreadQuit`（33-r6 落，`cmd/wisp/panel_resident_windows_test.go:445`）
  - `TestAC13BringUpRefusesAThreadWithAQueuedClose`（33-r7 落，同文件 `:522`）
- `diff top-33r7 top-mine` ＝**空**（本腿这发与编排者那一发逐名同集）。
- **静悄悄消失＝零枚**；**静悄悄多出来＝上面两枚**，都在票 33 射程内（AC#13 那族顺序依赖），归因见 §4。
- ⚠ 33-r6 还删过一枚 `TestAC13BringUpSurvivesAStaleCloseOnAReusedThread`（派单 §4(甲) 点名要复核的那枚）。本腿现量：`grep -rn "StaleCloseOnAReusedThread" cmd/ internal/` ⇒ **只在注释里剩指认**，函数名在树里 **零枚**，且它**从未进过任何一发名册**（r5 那发的 157 名里没有它）⇒ 删除不产生名册缺口，但它是 33-r6 自己刚加的，见 §4(甲) 判语。

### 2.4 其余门禁（终态）

| 门 | 命令 | 读数 |
|---|---|---|
| d22scan | `sh scripts/d22scan.sh` | **rc=0 clean**；分母：bans #1-5 `internal/=224` `cmd/=34`；ban #6 `frontend/=85`；ban #7 `internal/tools/=23`；ban #8 `design/=39` `frontend/=85` `internal/=476` `cmd/=81`（⚠ 分母＝文件枚数，不是违规数）。另记：`skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`——**这行本身证明工作树的 `frontend/dist` 里真有一包 assets**（我只读 git/d22scan 的元数据，⛔ 未读 `frontend/**` 任何内容） |
| gofumpt | `"$GOPATH/bin/gofumpt" -l cmd/wisp/ internal/panel/`（GOPATH＝`D:\work\base\gopath`） | **空输出** ⇒ 名下无未格式化件 |
| staticcheck | 未跑 | 〔未复认〕：本机版与 CI 钉版不同，派单 §0 禁止拿它出结论 |
| `go build ./...` / `go vet` | 未跑 | 口径偏离由编排者补跑销账（派单 §0）；本腿只用整包那一发作凭据 |
| `go mod tidy` / `go get` | **未跑** | 派单 §0 禁（HEAD 上 tidy exit 1）；相关债务裁见 §7 |
| 桌面卫生（起手 16:49:45 ／ 终态 16:55:23） | `tasklist //FI "IMAGENAME eq msedgewebview2.exe" //FO CSV \| grep -c msedgewebview2` | **12 → 12**（零增量）。逐枚按 `--user-data-dir` 归因（PowerShell `Get-CimInstance Win32_Process`）：5 枚 `MicrosoftWindows.Client.CBS_*`（系统自带）＋6 枚 `com.clipsync.desktop`（第三方应用）＋1 枚 CBS 子项＝12；**零枚指向本仓**（既无 `%TEMP%\wisp-33r5-panel-profile` 也无 `%APPDATA%\wisp`）。⛔ 不指控谁泄漏，只报数 |
| 孤儿子进程 | `tasklist //FI "IMAGENAME eq wisp.test.exe"`／`eq go.exe` | **0／0** |

### 2.5 ⛔ 本腿写面自证

本腿**未改**任何产码或测试文件。所有定向突变（§5）都用 `git cat-file blob HEAD:<path> > <path>` 还原，
还原前后各一次 `certutil -hashfile <path> MD5`，且每次还原后核 `git diff --stat HEAD` 为空。逐枚读数在 §5 每节末尾。

## 3. 判不动的地方

（填写中）
