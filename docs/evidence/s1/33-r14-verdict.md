# 33 — Panel host C27, verdict leg (non-implementer) on implementation leg 33-r14, commit `721b4b00f`

**Verdict leg:** 33-r14-verdict（非实现者；实现者≠裁决者，本表全部读数来自本腿自己的新鲜复跑，不采信实现腿的任何自述）
**Commit under review:** `721b4b00f`（33-r14 面板宿主腿；8 枚文件：产码 5 枚＋判据件 2 枚＋证据件 1 枚）
**HEAD at verdict time:** 起手 `7393cc604`（= `1cc67664b` 之后再一枚台账 commit，`git diff --stat 721b4b00f HEAD -- cmd/wisp` ＝ **空** ⇒ 台账那枚不触本票写面，在 `721b4b00f` 上验＝在 HEAD 上验）；收笔时仍 `7393cc604`。
**Method:** 规则 21 共享克隆（主树 `cmd/wisp` 有 323 腿两枚未跟踪在飞件 `panel_resident_323_windows_test.go`／`panel_resident_assembly_323_test.go` ⇒ 不许在主树跑）。克隆 `D:\tmp\33r14-verdict\wisp` checkout `721b4b00f`（porcelain ＝ 0 行）＋ 基线对照克隆 `D:\tmp\33r14-verdict\wisp-base` checkout `e3abe0ebb`（`721b4b00f` 的父）。两枚克隆都在仓库目录**外**；运行脚本在 `D:\tmp\33r14-verdict\*.sh`，零仓内 worktree/checkout。
**Environment:** Git Bash（`D:\work\soft\Git\bin\bash.exe`），`unset GOTMPDIR TMPDIR TEMP TMP`，PATH 前置 `$REPO/build:$REPO/third_party/sherpa-onnx`＋mingw64/bin，**CGO_ENABLED 全程未设**（gcc `E:\work\base\msys64\mingw64\bin\gcc.exe` 在场，go1.27.1）。`gofumpt`＝`/d/work/base/gopath/bin/gofumpt.exe`。零 GUI、零 `wisp.exe`、零 `*winlive*`、零 push、零 `frontend/**`／`design/**`／`internal/panel/**`／台账写面、票面一字未动。

## 一、新鲜复跑读数（全部本腿自取，非转述）

| # | 闸门（克隆 @`721b4b00f`） | 读数 |
|---|---|---|
| 1 | `gofmt -l cmd internal tools` | **空**，rc 0（比派单口径宽：含 tools/） |
| 2 | `gofumpt -l cmd/wisp internal/panel` | **空**，rc 0 |
| 3 | `go vet ./internal/panel ./cmd/wisp` | **rc 0**，无输出 |
| 4 | `go test ./cmd/wisp -run '<7 新尺逐名>' -count=2 -timeout 20m -v` | **14/14 PASS**（7 枚各两轮），逐名：`TestAC6CSPHeaderIsOnEveryServedResponse33r14`／`TestAC6CSPHeaderHasOneBuildSite33r14`／`TestAC8SnapshotClampBranchBooksThroughThePumpEntry33r14`／`TestAC8SnapshotClampBranchDropsIdsAndKeepsTheDigest33r14`／`TestAC8SnapshotShortPacketKeepsItsIds33r14`／`TestAC8TryEntryRefusesWhileTheBookingLockIsHeld33r14`／`TestAC8BlockingEntryStillBooks33r14`，`ok cmd/wisp 0.190s` |
| 5 | `go test ./internal/panel -count=2 -timeout 20m` | FAIL 恰 **3 名 × 2 轮**＝`TestComposerContractTypesMatchFrontend`／`TestComposerRenderFixtureTellsTheTruth`／`TestC21DesignTokensFourWayAgree`（red 句逐字与基线克隆同） |
| 6 | `sh scripts/d22scan.sh` | **rc 0 clean**（正控 `PASS=45 FAIL=0` 先行）：`bans #1-5 internal/=366, cmd/=58, #8 internal/=915, cmd/=219`（分母＝文件枚数）——与实现腿自述一致 |
| B | 基线对照克隆 @`e3abe0ebb`：`go test ./internal/panel -count=2` | FAIL **同名 3 枚 × 2 轮** ⇒ 三枚红在 `721b4b00f` 的父上就在，**零新红**（本腿写面从未触 `internal/panel/**`，对照实证） |

**第 4 项的两发环境教训（如实记）**：克隆里 `cmd/wisp` 测试二进制首次跑挂 `exit status 0xc0000135`（DLL not found）——`objdump -p wisp.test` 现量导入表＝`KERNEL32.dll / msvcrt.dll / sherpa-onnx-c-api.dll`，缺的是**未跟踪**的 `third_party/sherpa-onnx/*.dll`（共享克隆不带未跟踪件）。处置＝从主树把 `build/*.dll`＋`third_party/sherpa-onnx/*` 四枚 DLL **拷进克隆**（运行时 fixture，零跟踪文件影响，事后 `git status --porcelain` ＝ 0 行）；补齐后 14/14。这与实现腿诚实注记的"克隆缺未跟踪模型/DLL 产物"同一成因，**不是**本票代码缺陷。

**红色名册归因（给编排者）**：`TestComposerRenderFixtureTellsTheTruth` 在基线 `e3abe0ebb` 同红，但它**不在**票 33 10-01 那张已知红名册上（那张列 approval/composer/frontend_hygiene/tokens_fourway 四枚）⇒ 名册过期，需非实现者/编排者重钉基线。本腿只报不修。

## 二、AC#6（CSP 响应头）独立验证 ＝ **PASS**

1. **逐字比对（绕开判据件自己的常量，双源现量）**：`panelCSPHeaderValue`（`panel_control_windows.go:101`）与 `docs/specs/SPEC-06-security-gatekeeping.md:126` 反引号内的策略串做字节级比较：**174 字符 == 174 字符，EXACT-EQUAL = True**（`connect-src 'none'` 在串内）。`wantCSPPolicy`（判据件 `panel_csp_33r14_windows_test.go:29`）同串，双向对账成立。
2. **响应 HEADER、不是 `<meta>`**：`panelCSPHeader`（`:107`）＝ `"Content-Security-Policy: " + panelCSPHeaderValue`，经 `panelResponseHeaders`（`:513`）→ `respondWithEmbeddedAsset`（`:494`）→ `panelBuildResponse`（`:538`）→ `env.CreateWebResourceResponse(data, status, reason, headers)`（`:545`）。装配在 `AddWebResourceRequestedFilter`（`:468`）这一支——本树唯一的供页缝，浏览器把它当**响应头**执行，页面改自己文档剥不掉。
3. **200 与 404 都带**：`panelResponseHeaders` 200 臂＝`Content-Type: …\r\n`＋CSP 行；miss 臂＝CSP 行（404 不发明 Content-Type）。404 在生产**真实可达**：`panelAssetLookup`（`:563-576`）对 nil assets／Resolve 出错返回 `(nil, "", 404, "Not Found")`，同一调用点走同一构造点。
4. **单一构造点（四重现量）**：① 全 `cmd/` 产码 `Content-Security-Policy` 字面量**恰 1 处**（`panel_control_windows.go:107`，grep 现量，其余 5 命中全在判据件）；② `panelResponseHeaders` 定义恰 1（`:513`）、调用点恰 1（`:494`）；③ `respondWithEmbeddedAsset` 是过滤器唯一应答路径（`:466`）；④ 判据件第二枚尺（AST 扫 string literal，注释豁免）钉"恰 1"——M2 变异下它保持绿（数的是字面量不是策略），牙分工正确。
5. **变异证明（本腿独立重做，非采信实现腿）**：
   - **M2** `connect-src 'none'` → `'self'` ⇒ `TestAC6CSPHeaderIsOnEveryServedResponse33r14` **红**（三牙同咬：逐字比对牙＋200 头牙＋404 头牙，红句逐字落在盘上 `mut.log`），`HasOneBuildSite` 保持绿；还原后 `md5sum -c` 两枚 **OK**。
   - **M3** miss 臂 `return panelCSPHeader` → `return ""` ⇒ 同尺**红且只咬 404 牙**（`the 404 header set carries no CSP line: ""`），200 臂不受累；还原 **OK**。
   - 变异体零枚进提交/暂存；事后克隆 porcelain ＝ 0 行。

## 三、AC#8ⓐ（摘要钳制支装牙）独立验证 ＝ **PASS**

1. **起手零覆盖现量（在 `aa7159b29` 复量）**：`git grep summaryClamp aa7159b29 -- cmd/wisp` ＝ 恰 4 命中：产码 3（分支 `panel_pump.go:422`／注释 `:432`／常量 `:435`）＋ 判据件 **1 枚且是注释**（`panel_config_248_test.go:181`）⇒ 实现腿"分支起手零覆盖"成立（票面引 `:186-192` 已漂，腿具名漂移＝如实）。
2. **改名零行为变化（逐字现读）**：`bookPanelSnapshotBlocking`（`panel_pump.go:393-399`）函数体与改名前逐语句同（锁→记账→记 summary→return nil）；生产挂载 `run.go:958` 的 `Out:` 随名指向它 ⇒ 生产路径不变。`tryBookPanelSnapshot`（`:407-415`）＝ `TryLock` 输者 `return false, nil` 不记账不落日志；**生产调用者 0 枚**（grep：定义＋run.go 注释＋判据件两处调用，余零）——实现腿已如实自陈，归票 35 AC#5 的合并/背压消费地界，本腿不记缺陷。
3. **5 尺形状（精确记账）**：**3 枚走 Publish 生产入口**（`BooksThroughThePumpEntry`／`ShortPacketKeepsItsIds`／`BlockingEntryStillBooks`，`Out` 用产码同名 `bookPanelSnapshotBlocking`，与 `run.go:958` 同值）；`DropsIdsAndKeepsTheDigest` 直调 `panelSnapshotSummary`（过界包 3×120 字符 id，本腿现量 2347 字节、钳后行 `depth=3 pending=- … sha256=` 且 ≤440）；`TryEntryRefusesWhileTheBookingLockIsHeld` 钉 TryLock 语义本身（占锁 ⇒ (false,nil) 零落账；释放后 ⇒ 记账成功）。派单口径"5 尺全走 Publish"按字面是 **3/5**；另 2 枚分别直调被钉函数／钉 TryLock 语义——判据价值不减（直接调被钉函数正是"问分支自己的输出形状"），但口径如实记。
4. **变异证明（本腿独立重做）**：**M1** `if len(line) > summaryClamp {` → `if false && …` ⇒ `DropsIdsAndKeepsTheDigest` **红两牙**（`pending=` ids 全文 spelled out＋`the clamped line is 496 runes`），而 `ShortPacket` 负控、`BooksThrough`、`TryEntry`、`Blocking` 四枚**全绿** ⇒ 尺分得开"分支"与"函数拼写"；还原 `md5sum -c` **OK**。

## 四、逐 AC 判语表（票 33 全 14 格；勾框仍归编排者按本表翻，本腿一字未动票面）

| AC | 判语 | 依据（本腿现量） |
|---|---|---|
| **AC#6**（CSP 头） | **PASS（本轮 33-r14 落地，本腿独立复证）** | §二全项：逐字 174==174、200/404 双臂、单构造点四重现量、M2/M3 变异各红后逐字节还原 |
| **AC#8ⓐ**（钳制支） | **PASS（本轮 33-r14 落地，本腿独立复证）** | §三全项：起手零覆盖钉死、生产入口不变、M1 变异红＋负控绿 |
| AC#8 的 Hide 读数顺带项 | **落地成立（零行为变化）** | `panel_host_windows.go:795-811`：`Hide()` 补 `panel window hidden (focus returned to hwnd 0x…)`／`(no prior window…)` 两行，Win32 调用与簿记未动（diff 现读） |
| **AC#8ⓑ**（EvToolStart） | **REMAINING（真缺口，不是本腿欠账）** | 全仓 `.go` 现量恰 **3 命中**＝`internal/agent/sink.go:24`（注释）／`:25`（定义）／`cmd/wisp/run.go:2170`（消费 case；实现腿写 `:2166`，4 行漂，同一 case）⇒ **零发射者**。装牙＝造假事件（F-PUMP-6 明拒）、删 case＝产品形状决策（console 工具行＋`changed=true` 脉冲在发射者真落地那天是对的）⇒ **必须问 owner**（见 §六 Q2） |
| **AC#7**（快照双瞬间） | **REMAINING（触发条件未如期发生＋指名修法证不了实值）** | ① `.Marshal()` 调用点现量 **13 枚全在 `_test.go`**（定义 `internal/panel/pump.go:349` 非测试）⇒ "生产首调 Marshal"从未发生；生产消费全字节走 `Publish → src.Out(bookPanelSnapshotBlocking)`（`pump.go:360-380`、`run.go:958`）＋票 329 推宿主。补一枚假调用＝F-PUMP-4 之罪。② `Snapshot()`（`pump.go:238` 起）逐源现读：Verdicts/Results/Mode/Workspace/Git/Files/Model/Credential/Tasks/Instructions/L1Windows 各在**各自外部锁**下读，`pump.mu` 只护 publish 计数（`:389-397`）⇒ 持泵锁跨 `Snapshot()` 对真源**证不了双瞬间一致**，红→绿翻转不可达 ⇒ **必须问 owner 定契约**（见 §六 Q1） |
| AC#1–AC#5 | 归旧腿／终裁，**非本轮射程** | 宿主 r1–r13 已落（focus 那枚已知负载敏感红、AC#2 时延读数等都在旧表），本轮七文件 diff 未触它们 |
| AC#9 | 归 33-r2/r3（CLI 入向缝已落地在册） | 非本轮射程 |
| AC#10 | 已由 255 族腿落地（schema→panelWidthPx/HeightPx＋reshow SetSize＋两枚判据件） | 非本轮射程，旧表有据 |
| AC#11 | 已勾（编排者 10-01 翻） | — |
| AC#12 | 前件闭合（出厂链把 dist 填上，A794），勾归非实现者终裁 | 非本轮射程 |
| AC#13 | 次序支已落（`coldStartPageHandover`＋33-r10 判据件），"屏上有内容"那维被 A794 白屏实证否证 → 归票 299 | 非本轮射程 |
| AC#14 | 编排者 2026-10-16 裁归**票 329**（出向推快照那一跳） | 非本轮射程 |

## 五、诚实边界（本腿没做的，逐条具名）

- **整包 `go test ./cmd/wisp` 未跑**（派单四闸门不含它；实现腿克隆整包两轮的环境性红——15e 模型产物族、resident 族——本腿未复跑未归因，**不采信也不否证**）。本轮判语只覆盖派单四闸门＋本腿自加的 gofumpt／panel count=2／基线对照／三发变异。
- **零真窗验证**：AC#6 的"响应头"凭据是过滤器构造链＋判据件＋变异，**不是**真 WebView2 里抓到的响应头（真窗属 `*winlive*`／GUI，本腿明令禁止）。若 owner 要"眼见为实"那一级，需一次机主在场的真窗会话（同 A794 那次形状）。
- 本腿未跑 `staticcheck`（旧账：本机版与 CI 钉版不同，A500 裁定〔未复认〕照旧）。
- 实现腿证据件两处行号小漂如实记：`run.go:2166`→现量 `:2170`；`panel_pump.go:422`（起手）→现 HEAD `:444`（判据件加行所致，内容锚不受影响——M1 变异照红）。
- 克隆两份＋脚本＋pristine 拷贝留盘 `D:\tmp\33r14-verdict\`（零删除，供复核）；主树 porcelain 始终只含 323 腿那两枚未跟踪件，本腿零污染。

## 六、必须摆给 owner 的（编号清单；不答的代价写明）

1. **AC#7 的契约形状（卡死格）**：票面字母的触发条件（"生产首调 `Marshal()`"）与现实（生产走 `Publish→Out` 全字节）已脱节，且"持泵锁"修法证不了实值。**owner 要答的就一句：哪两个瞬间必须一致？**（例：T1 时刻的 verdicts 深度与 T2 时刻的 results 是否必须同帧？mode 变更与 results 变更要不要同帧？三个外部判官＝gate 队列锁／stream 锁／mode 锁，答案不同修法不同。）**不答的代价**：这一格永远空着，且任何腿都不能合法替他填——它是"改契约文字"类，不是腿本地可编辑项。
2. **AC#8ⓑ 的二选一**：`case agent.EvToolStart:`（`run.go:2170`）保留等真发射者（登记"零发射者"为在册事实），还是删掉（产品形状决策，等真有发射者再随发射者一起回来）。**不答的代价**：一个不可达分支继续"看似在守什么"，与 F-PUMP-6 的本意相悖；但删除是一个不可逆性低、可随发射者回炉的动作，倾向**删**，等 owner 一句话。
3. **红色名册重钉（编排者侧，低利害）**：`internal/panel` 在册已知红需从 3 枚（Contract/RenderFixture/TokensFourWay）重钉并落台账——`RenderFixture` 那枚今天在任何名册上都找不到（实现腿已点名，本腿基线对照复证同红）。**不答的代价**：下一枚碰 `internal/panel` 的腿会把这 3 枚红误记成自己造成的。

**本腿判语汇总**：`721b4b00f` 的两格主判据（AC#6、AC#8ⓐ）**全部 PASS**——四闸门新鲜全绿、三发变异各红后逐字节还原、基线对照证明零新红。票 33 剩余的真缺口恰是两枚**等 owner 答题**的格子（AC#7 契约形状、AC#8ⓑ 删留），其余 AC 归旧腿与 33-v2 终裁，本轮不欠。
