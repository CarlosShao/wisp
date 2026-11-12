# 近完成后端票 · 只读分诊（triage-nd，2026-11-12）

- 派单：`.scratch/wisp/probes/orch/20261112/txt/triage-nd-dispatch.txt`。只读分诊、零产码；本件是唯一交付物。
- 待诊 12 张：23, 174, 216, 220, 230, 246, 247, 256, 274, 290, 291, 346。
- ⛔ 票面引用的行号/锚号一律以本腿现读为准并标"现读"；引用前请重跑。AC 框数为本腿现读逐格清点（`grep -c '^(- \[[ x]\]) \*\*AC#'` 同族尺）。

## 0. 总表

| 票 | chk/unc | 归类 | 一句话卡点（逐格未勾 AC） | 可自办写面 / 波次 |
|---|---|---|---|---|
| 23 | 4/1 | **F**（混 B 味） | AC#3 残三枚：R1 真 fixture app 300 字符注入＋R2 真提权窗 UIPI 实机（真机窗口期）、R3 确认卡三选项 UI＝`frontend/**`（机主外部 agent）；23-verdict §8 终裁"归属 frontend/owner-open-window 一次收口" | 无后端写面；R3 前端＋R1/R2 真机同窗口批发 |
| 174 | 7/2 | **A** | AC#4＝artifacts 目录入授权根（owner 批准已落 A1571；施工图=174-screen）；AC#5＝Q-59 分账结题说明（编排者结题动作，落地腿附草稿） | `cmd/wisp/run.go` 授权根装配一处＋两向判据；约 2 波 |
| 216 | 5/1 | **G** | AC#5ⓑ 被 2026-10-17 编排者裁"归口 S6 DEFERRED（票 225 射程）、⛔ 本票翻"；ⓐ＝d22scan def-use 仪器工程，账 A1010"队列尾、本波不放、谁也不许排前" | 本票结构性闭合不了；凭据 A1010/A1013 |
| 220 | 4/1 | **E** | AC#5 要 `go test ./cmd/wisp` 整包终态读数；220-ac5 三发（40m/60m/60m）全被球窗 WebView2 UDF 自挤掐死（根因已归票 408，Q-85 在案） | 被 cmd/wisp 整包死区卡 |
| 230 | 4/1 | **E** | AC#5 整包终态；230-cells ③发 cmd/wisp `-timeout 25m` 仍包级超时 1500s（真窗/热键活体测试挂起，归活体测试票面） | 被 cmd/wisp 整包死区卡 |
| 246 | 8/1 | **A** | AC#8 修已落盘（246ac8：`signal.Notify` 提前，`-count=10` 10/10 绿，件 `docs/evidence/s1/246-ac8-boot-exit-notify.md`），只剩复认＋翻勾＋归属裁定 | 只读复认腿 1 波；零产码 |
| 247 | 8/2 | **B** | AC#2 归口票 297/291（电平尺定义）；AC#6 欠"无权限/禁用设备"两形（要机主动 Windows 设置）＋"设备被占"自造形（编排者欠账）＋"降级时球与任务管线仍在跑"仪器 | 机主在场真机窗口（与 290/291 同批） |
| 256 | 2/2 | **A** | AC#1 实现已由票 265 完成（A1666 现读 `resident_approval_windows.go:401-407` 三字段全在），欠非实现者复判后"并入 265"翻勾；AC#2＝唯一真剩余（正控＋禁现定形，定向测试） | `cmd/wisp/resident_approval_risk_256_windows_test.go`（已存在）；约 2 波 |
| 274 | 9/3 | **C** | AC#8＝SPEC-11 §2.2 归属裁决（改 `docs/specs/**` 文字＝人工批准）；AC#0/AC#2 的读数已被 2026-11-06 `docs/evidence/s1/274-ac0-counts.md` 用真 CI run 42205868481 逐项量齐，只欠非实现者翻勾 | 翻勾半＝收口腿 1 波；契约文字那一格待人批 |
| 290 | 5/1 | **B** | AC#3 形②"用户主动开门后 SetAudioLevel 随声音变"＝机主在场对麦说话那发（形① 装配级已验）；与 247/291 同一次窗口批发 | 机主在场真机窗口 |
| 291 | 1/3 | **B** | AC#2 四形读数：291-mic（11-08）已证形②真人窗采不到（子代理无法触发人工交互）、形③音调麦未听到（要机主先查扬声器音量/麦增益）；AC#3/AC#4 落地按在 AC#2 之后＋曲线问摆机主 | 机主在场真机窗口 |
| 346 | 5/2 | **A** | AC#4＝winsec 双平台成对＋两弃权仪器各 3 发复跑（本机 Windows 原生 ⇄ Docker linux/amd64，174-r6 先例）；AC#2＝三句弃权红句仍在＋去牙反形（仓外 Linux 导出树物理字节）；collect-r2 注"AC#2 选形先裁" | 本机＋Docker 台面；约 1–2 波 |

**计数：A＝4（174, 246, 256, 346）；B＝3（247, 290, 291）；C＝1（274）；D＝0；E＝2（220, 230）；F＝1（23）；G＝1（216）。合计 12。**

---

## 1. A 类票补完施工图（下一批写腿照此派）

### 1.1 票 174 — AC#4（artifacts 入授权根）＋ AC#5（结题分账）

**依据**：台账 `A1571`（2026-11-06 机主 B 桶全批，逐字"改契约文字那步照批"）＋ 普查腿件 `docs/evidence/s1/174-screen.md`（2026-11-07，施工形状已裁）。

- **写面**：`cmd/wisp/run.go` 授权根装配处——174-screen 现读＝授权根仅 `run.go:622-623` 一路（`cfg.FS.AllowedDirs` 拷贝）、spill 落盘目录在 `:1682`/`:1836` 一带未入根；**行号落地前现算**。
- **做法**：给 `allowed` 加第二来源，**只加 `<dataDir>\artifacts` 这一枚目录**；`secrets\` 与 `config.toml` 留根外（`run.go` 注释逐字写着这两个不能进根）；`internal/config`、`SPEC-03:35`、默认值零字节；PLAN/specs 零字节（AC#4 推荐两行契约文字不改，A1571 使"要不要改文字"降为文书格）。
- **判据（两向＋正控，缺一向＝没量完）**：
  1. 默认配置（`allowed_dirs=[]`）发一发真 spill（超 4000 token）⇒ 同一发 `fs.read` **能读回全文**（L0、零卡）；
  2. 越界路径（artifacts 之外的根外路径）⇒ 仍 R2/L2 拒（现读逐字拒绝文本）；
  3. 安静正例（授权根内真实普通文件）⇒ 回执"注意"句 0 次、指针句仍在。
- **落点**：接缝级判据落 `internal/tools/**`（174-c1 表的形状，`.scratch/wisp/probes/174/c1/` 可重放）；端到端可选（174-r8 的 harness `cmd/wisp/task_output_pointer_notice_174r8_test.go` 形状，铺 sherpa DLL 后本机可跑端到端）。
- **AC#5**：落地腿交件附两行分账草稿（Q-59＝偏移参数待 owner 那一支 vs 本票＝授权根已由 A1571 批、本格落地状态），编排者结题时合入翻勾（A1279 先例：不需 owner）。
- **同批顺手格**：AC#2b 残半句"或批一张 L2 卡"——台账 `:20099` 已裁 `task.go:299/:349` 是文档注释非模板钉，不涉契约面，可一并补。
- **波次/互斥**：`cmd/wisp` 共享编译面每波一枚写腿；与 256（同面）、246 收口只读不冲突。预计 **2 波**：波1 落地腿；波2 非实现者验收＋编排者结题。

### 1.2 票 246 — AC#8（boot 时序 flake 归属）＝收口复认，零产码

**事实**：2026-11-07 落地腿 `246ac8` 已真修（`cmd/wisp/resident_windows.go` 把 `signal.Notify(bootExit…)`＋Stop defer 提到 `runResident` 第一句；修前正控 `-count=5`＝5 发 1 红，修后 `-count=10`＝10/10 绿零 `0xc000013a`；门禁 gofumpt/vet/build/d22scan 全绿；件 `docs/evidence/s1/246-ac8-boot-exit-notify.md`）。

- **要做的**：非实现者复认腿——① 亲跑 `-run '^TestAC228ExitRequestDuringBootStillLeavesThroughD38E$' -count=10 ./cmd/wisp/`（带 `third_party/sherpa-onnx`＋`build` PATH harness；跑前 `tasklist` 核 `wisp.exe`/`balldebug.exe`＝0）核 10/10 绿；② 核改动面＝仅 `resident_windows.go` 一枚产码＋票 228 的 `residentBootWindowExit228`/托盘挂载消费点零触碰；③ 翻 AC#8＋落归属裁定（名义地界属票 228，修已由 246ac8 改道完成——按票 316 并票先例在票面注明"并入 246 翻勾"或走撤票流程，二选一写明）。
- **判据**：无 skip、无放宽断言、无加大超时；红句归位可复现（修前形状用 `-overlay` 种回旧行验证仍红）。
- **波次**：**1 波**（只读复认）。与 174/256 的写腿不同波无冲突（本腿只读，但复跑会占 cmd/wisp 测试色 ⇒ 与整包读数腿错峰）。

### 1.3 票 256 — AC#1 收口（并入 265）＋ AC#2（唯一真剩余）

**依据**：编排者自我更正 `A1666`（2026-11-09）：§9 那枚批准（A1665）基于过期前提——现读 `cmd/wisp/resident_approval_windows.go:401-407` 逐字 `UI/Channels/Window/ApprovalTimeout/Logf/Grants: ra.grants` 三枚全在；票 265 已 `-done`（5 格全勾、ⓐ-Ⅰ 形、产码 `a04a095f` 等、四发突变全红全还原、仪器 `cmd/wisp/resident_grant_writer_265_windows_test.go` 在盘）。

- **AC#1（复判收口，只读腿）**：现读该 Options 六字段字面＋票 265 四发突变证据 ⇒ 按 A1666 指令"非实现者复判后按并入 265 处置"翻勾并注明归口；同时按 §6 建议把 A1665 批准记作废（append-only）。⛔ 编排者/实现者直接翻。
- **AC#2（写腿，唯一真剩余）**：
  - **写面**：`cmd/wisp/resident_approval_risk_256_windows_test.go`（**已存在**——起手先读它缺什么，别重造）；⛔ 碰 `cmd/wisp/run.go`（P7：第 424 行以后任何位移红 `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`；起手/收尾各跑一发 `go test -count=1 -run 'TestTicket255Roster' ./cmd/wisp/` 定向尺，读数进证据件）。
  - **正控（唯一有效形状）**：种 `[risk] confirm_timeout_sec`/`l1_window_sec` 种子值启动一发 ⇒ 断言 `ra.gate.Queue().Timeout()` 随种子变（`gate.go:169`＋`queue.go:126` 一带，行号现算）。**⛔ 用 `Window()` 当正控**——被 `MaxL1Window=3s` 钳死会读出恒绿假象，只可当"钳位仍在"的反控。
  - **禁现格（两形先定形再断言）**：只禁第一形（`g.grants==nil` 那支）——钉 `gate.go:675` 一带那句独有词组「本机没有接入会话授权记账」；**⛔ 不许钉前缀 `approval: GRANT-DROPPED`**（第二形"卡片主题在答复前已离开队列"与 Grants 无关、必须仍能出现；既有仪器 `internal/agent/approval/ticket224_reply_grant_test.go:283` 用前缀、不可复用也不许改它）。
  - **形状级事实写进判据**：`g.window`/`q.timeout` 是构造期定值、全仓无 re-apply 路径 ⇒ AC#2 只钉"带种子值启动这一发"，**不许写成"改配置活进程立刻跟着变"**。
- **波次/互斥**：`cmd/wisp` 写面 ⇒ 与 174 不同波；起手核 A1658 降载令是否仍生效。预计 **2 波**：波1 AC#1 复判（只读，可与 246 收口腿同批不同时）＋AC#2 写腿；波2 非实现者验收。

### 1.4 票 346 — AC#4（双平台成对读数）＋ AC#2（守卫保持＋去牙反形）

**依据**：collect-r2（2026-11-09）已翻 AC#0/1/3/5/6；剩 AC#2/AC#4，读数需求具名在 `docs/evidence/s1/346-collect-r2.md`。

- **AC#4（测量，可本机自办）**：
  - `./internal/winsec/` 两平台成对：**本机 Windows 档原生** ⇄ **Linux 档 Docker Desktop linux/amd64**（先例＝174-r6：本地 golang 镜像、仓只读挂载实跑 PASS）；同一 `-run` 名册、`-count=1 -timeout 30m`、四数逐枚点名 ⛔ 混用 ⛔ 相减。
  - 两枚弃权仪器（`internal/winsec/openleaf_322_other_test.go`／`posix_seal_window_120_other_test.go`，均 `//go:build !windows`）各自 **≥3 发独立复跑**，逐名点名红⇄绿⇄弃权散布如实报；红名册作差按**名字集合**比（⛔ 按枚数，`TestResolvePerCallBudget` 间歇那族教训）。
  - 门禁四数只在共享工作树、跑前 `unset GOTMPDIR TMPDIR`（票 311 形 B GUARD）。
- **AC#2（守卫保持＋去牙反形）**：
  - 三句弃权红句逐字仍在（现读锚：openleaf `:206` 一带、posix `:370` 一带两句；行号现算），只在"真到不了窗口"时响。
  - 去牙反形＝**仓外 Linux 导出树物理字节**摘弃权守卫（⛔ `-overlay` 不够——票面逐字"只打仓外导出树物理字节、每发前后各验 pristine 哈希"）⇒ 仪器必须从"弃权红"变"假绿"且这一维能被点名。
  - ⚠ **前置**：collect-r2 报"AC#2 需 Linux 档导出树＋**选形先裁**"——编排者须先在甲（台面侧）/乙（仪器侧＝动 `internal/winsec` 产码注入点＝⛔ 本票射程、另票＋人工批准）/丙（判据侧，弃权照红但同屏打四数读数）里选形；**选形不阻 AC#2/AC#4 的读数半与去牙反形**，落地形状等裁。
- **禁区（复述）**：⛔ 放宽/删弃权红句、⛔ `t.Skip`、⛔ 拉重试到总能过、⛔ `continue-on-error`；⛔ 动 `winsec.go/winsec_other.go/winsec_windows.go`、票 113/116/119/120/322 既有钉、冻结件/thresholds/golden/allowlist/go.mod/go.sum/.github/**。
- **波次**：**1 波**台面腿（Docker＋本机成对读数＋去牙反形，选形批注随派单）；若选甲/丙需动 `scripts/abstention-census.sh`（346-a1 已落）或 CI ⇒ 另波另批。

---

## 2. 非 A 类逐张凭据（诚实优先，别把 D/E 硬说成 A）

- **23（F，混 B 味）**：23-verdict（`docs/evidence/s1/23-verdict.md`）终裁 AC#3 PARTIAL，残余 R1/R2/R3 具名；编排者 2026-11-09 翻框节逐字"残余 R1–R3 归属 frontend/owner-open-window，等真机窗口期一次收口"。R3 写面在 `frontend/**`＝F；R1/R2 真机窗口与 R3 同批＝B 味。后端余量（L2 声明/UIPI 拒绝序/绑定键）已闭，无 A 可做。
- **216（G）**：2026-10-17 编排者裁定节：AC#5ⓑ"归口那行 DEFERRED 的 S6 子项、⛔ 在本票翻它、⛔ 另立接线票"（`SPEC-12:77` + `PLAN.md:182` 推迟面）；AC#5ⓐ"没牙"登记为 A1010 仪器工程（def-use 级 d22scan 禁令），"本波不放，谁也不许把它排到 `internal/risk` 那枚用例前面"、"此后任何人⛔ 把 AC#5ⓐ 读成已钉"。⇒ 本票剩余格在票内结构性闭合不了（归口＋按住），非"待派活"。
- **220 / 230（E）**：两票 AC#5 字面都要 `go test ./cmd/wisp`（整包）终态读数。实据：220-ac5（`docs/evidence/s1/220-ac5-endstate.md`）三发 40m/60m/60m 分别死在 `TestAC294MuteHotkey…`/`Test296Bridge…`/`Test258BridgeMutation…`（球窗 WebView2 UDF 固定共享路径 `%LOCALAPPDATA%\wisp\orb` 自挤，2026-11-08 408 腿定性为独立根因）；230-cells（`docs/evidence/s1/230-cells-closure.md`）③发 `-timeout 25m` 包级超时 1500.217s、两发两死点。`./internal/...` 两侧都达终态，卡点**只在 cmd/wisp 整包**＝314/408 死区。任何"补读数"写腿今天都交不出终态，故 ⛔ A。
- **247/290/291（B）**：三票未勾格全部指向同一次机主在场真机窗口（290 处置节具名"与票 247 三格＋票 291 AC#2 四形同一次机主在场窗口批发"）：247-AC#2 归口 297/291、AC#6 无权限形要机主动 Windows 设置；290-AC#3 形②要机主对麦说话；291-AC#2 形②真人窗 291-mic 实证"子代理无法触发人工交互"、形③要机主先查扬声器音量/麦增益。都不是改测试/改 Go/补证据能闭合的。
- **274（C）**：AC#8 的裁决对象是 `docs/specs/SPEC-11` §2.2 那两句（`--with-frontend` 开关兑现 vs 承认被 `frontend/dist`＋CI 取代）——"取代＝要动 spec 文字＝人工批准"；274-ac0-counts §5 表 #6 同判"契约文字——待人批"。AC#0/AC#2 的读数已由 2026-11-06 只读腿用真 CI run 量齐（表 #1–#3 判"已闭"），只欠非实现者翻勾＝A 味收口半，但结题卡在 AC#8 ⇒ 主类 C。D＝0（本批没有"要机主拍设计/花钱选形"的新项；机主拍板面均已由 A1571/A1549/A1665 先例覆盖或已归口）。

---

## 3. 派单时注意（本腿现读到的排程事实）

- `cmd/wisp` 共享编译面每波一枚写腿：174 与 256 同面 ⇒ 必须分波；246 收口复认（只读但要占测试色）与整包读数腿错峰。
- A1658 降载令（2026-11-09）曾按住 cmd/wisp 大测试——派单前现核是否仍生效。
- 346 的 Docker linux/amd64 路线有 174-r6 先例（本地镜像、仓只读挂载），⛔ 不 push。
- 真机/整包读数前先 `tasklist` 核 `wisp.exe`/`balldebug.exe`＝0（A785 定式）。
