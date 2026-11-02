# 票 22 落地腿 22-r1 —— web 工具 D30：SSRF／重定向链／协议矩阵／首访提示

`s1/` 证据 · 写腿 `22-r1` · 2026-11-02 19:1x–21:3x +08 · 起手锚点 HEAD `5873aa3e8`（A1419 派单枚）
派单存档＝`.scratch/wisp/probes/orch/20261102/txt/22r1-dispatch.txt` —— ⚠ **该存档盘上不存在**（当日 txt 目录只有 244/337/340/344/345/346/395 各腿；A1419 的提交文字即派单正文）。本腿按票面＋派单文字执行，此差异如实登记，不自行补造存档。

## 〇、起手三核（现跑）

| 核 | 读数 |
|---|---|
| HEAD 现盘最新 | `5873aa3e8 编排者A1419：…派22-r1(web D30)补第3写腿` ✅ 派单枚即 HEAD |
| 索引干净／写面无他腿同包脏 | porcelain：`M internal/agent/{cost,loop}.go`、`?? internal/observe/*`、`?? internal/memory/costdao_44*` ＝ **44-r1 在飞面**；`internal/tools/` 起手零脏（31-r1 提醒面未落盘，其 D34 行落 `internal/tools` 属 S4，将来若同包只新增不相干文件名，两腿显式 pathspec 互不吞）✅ 非同包冲突，不停手 |
| 票面点名的串/行/文件还在（第19处） | `19-provenance-c25-done.md`／`20-host-bridge-fs-tools-done.md` 两枚 -done 键在盘 ✅；票面引用 `SPEC-07 §3`（`SPEC-07-tools-and-plugins.md:32`）、`SPEC-06 §6`（`SPEC-06-security-gatekeeping.md:85`）、`16.5.3`（`PLAN.md:2577`）、D30 五层（`PLAN.md:1063` 起）、F2（`PLAN.md:2409`）逐枚现抽对位 ✅ |
| 引文出处对不对（第20处） | D30 行号：AGENTS.md 索引写 `:1057`，**现跑＝`:1063`**（薄索引过期，票面只指 D30 号⇒不构成冲突，按 AGENTS.md 自己的"先现跑那把尺"规矩取现值）。`16.5.3` 行号现跑 `:2577` 与派单引用一致。web.search 裁决＝台账 **R4**（`pending-and-issues.md:50`）＋ **R9** 接口面（`:57-63`）✅ |

## 一、逐格 AC ＋ 尺 ＋ 读数

### AC1 — SSRF suite（每种私有网段形态含十进制/八进制/十六进制编码＋DNS rebinding，连接前拒；收方断零到达）

落点：`internal/tools/webnet.go`（`decodeNumericHost` inet_aton 全族解码：`2130706433 / 0x7f000001 / 017700000001 / 0x7f.0.0.1 / 127.1 / ::ffff:127.0.0.1 / 裸 v6 恢复路径`；`ssrfBlockedAddr`＝netip 判据与 R5 同源：`IsPrivate(10/8,172.16/12,192.168/16,fc00::/7) ∨ !IsGlobalUnicast(127/8,169.254/16,::1,fe80::/10,0.0.0.0,保留/多播)`；`dialContext` **解析→判→钉已检 IP 再连**（hostname 永不交给 dialer＝rebinding 在 pin 处关门）；`normalizeGuardedURL` 连接前拒编码/字面私有）。

尺（2026-11-02 21:0x 现跑）：
- `Test22SSRFDecodeMatrix`＝24 私有行（含全部编码族）+5 公有反控（`172.32.0.1/169.253.0.1/8.8.8.8/2001:4860::/100.64 CGNAT` 明记为"不在冻结清单⇒放行"残档）全 PASS
- `Test22SSRFDeniedBeforeConnection`＝真 TCP 收方 accept 计数：7 种字面/编码 URL + 名字解析到 127.0.0.1 + 名字解析到 169.254.169.254 ⇒ 全部拒、**receiver hits=0** PASS
- `Test22DNSRebindingBlockedPerHop`＝同名字 hop1 答 203.0.113.7（公有，veth  Below the judge 送达本地替身）、hop2 翻到 10.99.99.99 ⇒ "SSRF 防御"拒、第二收方 hits=0、每跳重解析（lookups≥2）PASS
- `Test22GuardAgreesWithR5OnLiterals`＝9 枚字面 R5↔netGuard 双向一致；4 枚编码形 **R5 按文档边界不解析、net 层必须接**（`rules_network.go:92-96` 注释的分工钉成测试）PASS
- `Test22UnbracketedV6LiteralsRefused` / `Test22SchemeMatrixAtGuard` PASS

### AC2 — Redirect chain（https→http、→file/私有 阻断；→Markdown；超限走 spill）

`checkRedirect`：链长上限、每跳 2048 复核、非 http(s) 拒、**https 起源→http 降级拒**、编码/字面私有跳转目标在下一拨前拒。
- `Test22RedirectDownsgradeBlocked`（TLS 起源→http，降级目标 hits=0）PASS
- `Test22RedirectToNonHTTPAndPrivateBlocked` 四子行 file/169.254/127.0.0.1/十进制 PASS
- `Test22RedirectChainLength`（cap=3 切断；两跳短链成功且 **Origin=最终 URL**＝C25 诚实）PASS
- Markdown：`Test22FetchHTMLToMarkdown`（h1/链接/列表/fence 转换✓；`<script>` 体文本整体丢弃✓；可见注入句**作为数据保留**✓（票 19/25 形状）；实体解码后不复活活标签✓——text run 出 `&lt;` 带一道）；`Test22FetchJSONAndText`/`Test22FetchBinaryRefuses` PASS
- 大响应：`Test22FetchHugeResponseCappedForSpill`＝工具端 1MiB 硬上限（复用 spill 的 `refRawOutputBytes` 标度，`budgets.go:35`）+ Truncated 标记，再喂**真 `agent.Spiller.Prepare`** 落 spill 文件（Spilled=true、Path 存在）PASS
- 状态码不回显正文＋审计零正文：`Test22FetchStatusNoBodyEcho`（404/500 错误文本与 logf 采集行都无 marker）PASS

### AC3 — Protocol matrix 16.5.3 穷举

`web_open.go::classifyOpenTarget`＝显式命名拒绝（search-ms/CLSID·mshta/脚本宿主·javascript·vbscript·shell·tel）**先于**任何注册检查；允许侧 `file:`(仅授权目录内+经 C26)、`mailto:`、`ms-settings:`、注册过的 `ms-*`；**其余一切默认拒绝**（白名单形，非黑名单）。
- `Test22ProtocolMatrix`＝17 行全表（两门×各协议，拒绝行断"未达 launcher"）PASS
- `Test22FileOpenPathArms`＝真 C26：目录内开、目录外拒、`file:///C:/…` 驱动位形、远程 `file://otherhost` 拒 PASS
- `Test22DenyBeforeAllowOrder`＝六枚**具名理由逐字**（CLSID/脚本宿主/脚本协议/命名空间/明确拒绝）——不是泛"拒绝"二字（M3 探针证明这样才咬，见§三）PASS
- `Test22UnknownSchemeDefaultsToDeny`＝9 枚发明协议全落拒侧 PASS

### AC4 — app.launch notepad fixture under Job；AUMID/name/path 三路

`launcher_windows.go`：path/name ⇒ `exec.Cmd` argv 向量 ⇒ `WebDeps.Jobs.StartInJob`（**C30**＝`proc.JobScope`，KILL_ON_JOB_CLOSE）；AUMID ⇒ `explorer.exe shell:AppsFolder\<AUMID>` 载体入 Job，**UWP 本体由系统激活宿主承载、不在 Job 树内＝平台事实**（姿态文字直写结果文本，不粉饰；名未落地≠失败）。name 解析默认走 App Paths 注册表（HKLM→HKCU，`name`→`name.exe`），`AppResolver` 缝留给票 23 `search.apps`。
- `Test22AppLaunchNotepadUnderRealJob`（windows 现跑）：path 臂+name 臂双开，TreePIDs≥2，`Close()` 后轮询 pid 退场 **PASS 0.13s**＝kill-on-close 亲验
- `Test22AppLaunchAUMIDClassifiesToWindowsArm`：正形收、6 枚畸形拒 PASS
- 跨平台形状（linux 同跑）：`Test22AppLaunchNameResolvesThroughSeam`（argv 向量逐字钉，⛔ 字符串拼接，D33）/`PathAndMissingShapes`/`JobFailureSurfaces`/`NoJobsHonest`（**未接 Job 的结果行必说"未纳入"**）/`AUMIDIsWindowsOnlyHere` PASS
- `Test22DeclsRosterMatchesD34Rows`：五枚工具名/Declared/Caps/Provider 与 D34 行一一对（search/fetch=L0 net、open=L1 net、file.open/app.launch=L1 shell）PASS

### AC5 — Taint-in-query→L2（R4 集成）＋ 首访提示一次、可记住

R4 循环接线（票 19 移交残余①的 tools 侧）：bridge `assessorFor` 已把 `prov.Detector(taskID)` 喂 R4，本票交付真实消费者——
- `Test22AttackChainFetchThenTaintedSearchEscalates`＝票面攻击链本尊：`fs.read` 明文凭据（打 taint）→ `web.search` query 携带 ≥8 字符片段 → 卡片 **L2、RulesHit∋R4、Reason 点名 "fs.read …git-credentials"** → **provider 调用计数=0**（效果未发生）；同 query 在新任务（先 `OpenTask`，= 装配根职责）跑绿＝scope 诚实。**未接线的第二半（每任务开头谁调 OpenTask）仍是票 12 接线面**，本件如实登记：测试显式 OpenTask，生产由 AdmitTask 边界拥有（bridge.go:896 注释区）。
- `Test22WebFetchResultIsMarkedAndBlocksFollowOn`＝web.fetch 结果被 C25 标记（`ScopeTaints` 见 `web.fetch`+URL origin），页面片段再上 search query→L2；`web.fetch` 已在 `sensitiveSourceTools`（票 19 面）零改动生效
- `Test22PrivateTargetRefusedByTheBridgeLayerToo`＝双层：bridge R5 先 L2 拒（NoGate 姿态、收方零到达），工具直调再 SSRF 拒
- 首访（`web_firstvisit.go`）：`Test22FirstVisitAppearsOncePerDomain`（一次、per-domain、跨工具共享）/`RememberFalseKeepsAsking`/`RefusalStopsTheCall`（拒=连接不发，收方 0）/`UnwiredFailsClosed`（**无提示通道⇒拒**，NoGate 同姿态）/`IPsAreNotDomains`（IP 字面不弹域提示）/`MemoryPersistsAcrossProcesses`（JSON 文件跨实例）/`StoreCorruptionStillPrompts`（坏档⇒多弹不默许，fail 朝拒）全 PASS。存储落点＝数据目录侧 `WebDeps.DomainStorePath`（票面"config/db"二选一都不越契约：`[net]` 冻结 schema 与 SPEC-02 冻结 DDL 均无此位，加列/加 key＝改契约文字须人工批——本腿不擅动，见§四残档）。

### AC6 — web.search 实现路径（R4 裁定 anysearch）

票面框文"OPEN decision"已由台账 R4 裁定（`pending-and-issues.md:50`，2026-09-20 用户拍板）：`SearchProvider` 接口 + **anysearch 兼容 provider 为默认**（R9 接口面：POST {base}/search · `query/max_results(1..10)/zone/language` · Bearer **可选**）；抓结果页降为备用，本票**不塞死码**——选择它得到显式"尚未落地（R4 裁定的备用路径）"工具级错误。
- `Test22AnySearchWireContract`＝路径/字段/匿名不发 Authorization 头/max_results 越界钳 10/results+items 两形解析＋**R7 卫生**（key 不进输出）PASS
- `Test22SearchToolRendersMarkdown`（provider EgressHost 弹一次首访）/`RefusalAndEmpty`/`SwitchDoesNotTouchCallSites`（换 provider 不动调用点＝票面硬约束）/`ScrapeProviderIsExplicitlyUnlanded`/`DefaultProviderWithoutBaseNeverDials`/`AnySearchErrorShapes` PASS
- AC6 框本身勾不勾＝裁决腿/编排者的勾（裁决者≠实现者）；本件给足可勾依据。

## 二、整包尺与红名册作差

- `go test ./internal/tools/ -count=1`（整包，windows）＝**红名册 {TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop} 一枚**。
  ⚠ 起手即该核而我漏核到中途才做——补记经过：把本腿 18 枚新文件整体挪离盘再单跑该测试＝**仍红**（`fs.move` 跨卷 staging 写 `C:\` 根 Access denied，本机环境），且 CI ledger 在册（`.scratch/wisp/probes/111/a5b/readings-…txt` 的 portable skip 名单含同名）⇒ **非本腿引入**，落点在他票面（票 20 家族），不越界修。
  写面冲突：`internal/tools` 全程仅此一枚红且作差为空——改动前后红名册相同。
- `-run Test22 -race -count=1`＝**ok**（5.0s）。Test22 顶层名册 **46 枚**（`grep -c '^func Test22'` 尺）。
- 双平台：`go build` windows-amd64 / GOOS=linux 均 rc=0；`go vet ./internal/tools/` rc=0。
- gofumpt（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`）：18 枚文件 `-l` 输出空。
- d22scan：正控 `go test ./...`（tools/d22scan 模块，-count=1）ok；实扫 `go run . -root <repo>` ＝ **clean**（354 production Go files；ban#7 internal/tools 34 files 含本腿新文件）。
- compose：`docker compose -f docker/compose.test.yml config --quiet` rc=0（mock-web/ssrf-target 两枚新服务在 manual profile，⛔ 不动 CI 的 test profile）。

## 三、反向判据（摘防护⇒点名用例红）＋还原核 —— 五发全跑

| # | 摘掉的防护 | 点名红 | 还原 |
|---|---|---|---|
| M1 | `decodeNumericHost` 编码族分支（只留字面量） | DecodeMatrix 10 行（decimal/hex/octal/dotted×2/two-part/three-part/link-local×2）+Redirect decimal 子行+GuardAgrees R5 层 | hash `73BE9333…D6B8` 逐字回原 ✅ |
| M2 | checkRedirect 的 https→http 降级行 | `Test22RedirectDownsgradeBlocked` FAIL（其余 redirect 行不伤＝断言分域） | 同上 ✅ |
| M3 | `explicitDenySchemes` 命中即拒的那步 | `Test22DenyBeforeAllowOrder` FAIL——**第一发不咬**：原断言只查泛"拒绝"，默认拒尾巴措辞也含"拒绝"；把测试改成具名理由逐字后二发咬中。测试变强，非放宽 | hash `45567DD7…929D7` ✅ |
| M4 | `startInJobOrBare` 的 Jobs 分支 | Notepad-under-real-job **双断言红**（姿态文字"未纳入"+Job 树空）| hash `799DD9BE…F97A`（重建件）✅ 见下★事故 |
| M5 | 首访 ask==nil 的 fail-closed 拒 ⇒ 改默许 | `Test22FirstVisitUnwiredFailsClosed` FAIL（1.27s＝真拨号尝试发生＝防护在场证据） | hash `AD78B42C…A466A` ✅ |
| M6 | webmd 的 script/style drop-entirely 臂 | `Test22FetchHTMLToMarkdown` 双断言红（script 体进正文+样式文本漏） | hash `4CFD4540…1D7FD` ✅ |

★ **M4 过程事故两条，如实登记**：① 第一发 M4 正则未匹配实际签名（我记成 `bool, error`，实际 `string, error`），变异未发生而"还原"行误用 webnet 备份覆盖了 `launcher_windows.go`——当场以 write 重建该文件（内容与覆盖前逐字一致：同文件其余测试+双平台构建+Test22 名册全绿为证；该文件从未入库，git 历史零暴露）。重建后二发才是真 M4（表列读数）。② 真 M4 那次 notepad 在 Job 外裸启成孤儿（"config.json - Notepad"，PID 24224，启 21:04:25 与变异窗口吻合，疑为 Win11 记事本会话恢复带出的标签页），已 Stop-Process 清掉并面板复零。**若机主当时正手工编辑某个 config.json，误杀即发生在本腿**，请机主核一眼。

## 四、残档与未接（不得无声消失）

1. **`[net].allowlist` 生产者半程**：`WebDeps.NetAllowlist→Facts→R5` 本票接通（消费者存在了），但"config 值→WebDeps"的写入点在 `cmd/wisp`（⛔ 禁区）⇒ `internal/config/unwired.go` 的 guard 行与 `unwired_test.go` 名册**一字未动**（现在写 allowlist 仍报错而非静默吞＝诚实）。整行撤除＝票 12 接线腿的活（届时 `lockedKeyDisposition` 同改），本件即其依据。
2. **anysearch 真网黄金流**：本腿测的是 R9 登记的接口面形状 + httptest；真实 key/真实响应留待票 12 接线后跑（R7：key 只从 secret 缝进，测试注入 `SearchKey func()` 已就位）。
3. **首访原生卡**：`FirstVisitAsk` 缝已定形（allow+remember 二值），渲染/一键记住属票 21/37 面；未接=全拒（fail-closed），已测。
4. **UWP Job 边界**：激活载体在 Job、应用在 OS 激活宿主——平台事实（同 402 那形），结果文字与证据双登记；如需硬包 UWP 属契约对话。
5. **编码形与 R5 的分工**：R5 不解析 inet_aton（其文件头注释的既有边界），net 层接——`Test22GuardAgreesWithR5OnLiterals` 把这条线钉成测试；若要在 verdict 层也拦编码形，动的是 `internal/risk`（越本票写面），留给编排者裁。
6. **CGNAT 100.64/10、6to4/NAT64 内嵌 v4**：不在票面冻结清单⇒按面放行，测试+注释双登记（fail 朝清单，不擅自扩契约清单）。
7. **compose 服务**：mock-web/ssrf-target 落 `manual` profile（Go 套件才是验收牙齿；CI test profile 零变化）。
8. **doc.go 的 "web/system/doc/search tools are 22-24"**：本票落地 web 组后此句过期一半；改它=动他人在读文字，留给裁决腿随结票一并正。

## 五、写面清单（18 枚，全部新增，显式 pathspec 逐枚点名 commit）

internal/tools/{web.go, webnet.go, webmd.go, web_fetch.go, web_search.go, web_open.go, web_firstvisit.go, launcher_windows.go, launcher_other.go,
web_guard_test.go, web_fetch_test.go, web_matrix_test.go, web_firstvisit_test.go, web_search_test.go, web_bridge_test.go, web_launch_test.go, app_launch_windows_test.go, web_testhelp_test.go}
＋ docker/mockweb/{README.md, nginx.conf, fixtures/clean.html, fixtures/injection.html}（新增）
＋ docker/compose.test.yml（追加两枚 manual 服务，既有服务零 diff）
＋ 本件 docs/evidence/s1/22r1-impl.md

⛔ 未碰：thresholds.go／golden／SLO／cmd/wisp／frontend／design／internal/config／internal/risk／票面／台账；无 add -A／amend／reset／push。

Commit：见下方枚字（落地后回填）。
