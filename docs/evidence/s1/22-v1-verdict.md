# 票 22 对抗验收裁决表 —— 验收腿 22-v1（裁决者≠实现者）

`s1/` 证据 · 非实现者验收 · 2026-11-02 21:3x–22:0x +08 · 起手锚点 HEAD `61046b65a`（编排者 A1429 派单枚＝盘上最新，≥`a18ff4546` ✅）。
被裁交付＝`547d9da00`（票 22 落地腿 22-r1，24 files +5014，全部新增）。
原始读数件（本腿现跑，只建不删）＝仓外 `D:\tmp\wisp\20261102\22v1\`：`scoped-run.txt`／`test22-full-v.txt`／`test22-race.txt`／`whole-package.txt`／`mut-x1…x5*.txt`。

## 〇、起手三核

| 核 | 读数 |
|---|---|
| HEAD 现盘最新 | `61046b65ad87…`＝派单枚（A1429），晚于 `a18ff4546` ✅ |
| 索引干净 | `git status --porcelain` tracked-dirty 枚数＝**0**（`??` 未跟踪为 .scratch 杂物，一律未碰）✅ |
| 票面在盘 | `.scratch/wisp/issues/22-web-tools-d30.md`（51 行，AC 六格全未勾——勾不勾由编排者翻，本腿只裁成立与否）✅ |

## 一、真跑（★scope 只到 internal/tools，未跑 `./...`——50-r1c 在写 internal/plugin，全仓会假红）

### 1.1 派单字面尺（逐字发货）

```
go test ./internal/tools -run 'Web|Guard|Fetch|FirstVisit|Matrix|Net' -count=1 -timeout 20m
→ ok github.com/CarlosShao/wisp/internal/tools 2.199s   exit=0   FAIL=0   （21:4x 现跑）
```

⚠ **尺的缺口（具名，不指控交付）**：该正则不锚定，实测只命中 Test22 名册 46 枚本票用例中的 **20 枚**；漏掉 26 枚＝`Test22AppLaunch*` 五枚（AC4 全族）、`Test22SSRFDeniedBeforeConnection`／`Test22DNSRebindingBlockedPerHop`（AC1 核心）、`Test22Redirect*` 三枚（AC2）、`Test22FileOpen*`／`Test22DenyBeforeAllowOrder`／`Test22UnknownSchemeDefaultsToDeny`（AC3）、`Test22Search*`／`Test22AnySearch*`（AC6）、`Test22URLLengthCap`／`Test22UnbracketedV6LiteralsRefused`／`Test22DeclsRosterMatchesD34Rows`／`Test22PrivateTargetRefusedByTheBridgeLayerToo`。名册里另有 14 枚 `Test221*/Test222*/Test229*` 属票 221/222/229，被前缀顺带扫进，全绿。⇒ 本腿补跑整名册（1.2）。

### 1.2 补跑整名册 ＋ race ＋ 整包回归

| 尺 | 读数 |
|---|---|
| `go test ./internal/tools -run 'Test22' -count=1 -timeout 20m -v` | 顶层 **60 枚：PASS 60／FAIL 0／SKIP 0**（46 枚本票＋14 枚他票顺带），2.9s（21:4x 现跑） |
| `go test ./internal/tools -run 'Test22' -race -count=1` | **ok 5.543s exit=0** |
| `go test ./internal/tools -count=1`（整包） | **ok 21.253s exit=0**，FAIL 名册＝空集（21:47±现跑） |
| AC4 notepad 是否真跑非 skip | `--- PASS: Test22AppLaunchNotepadUnderRealJob (0.15s)` 现读名册＝真跑（windows 本机） |

⚠ 对照交付腿自报：整包曾有 `TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop` 环境红（跨卷 staging 写 `C:\` 根 Access denied）。**本腿整包该枚绿**＝本机环境态而非回归；两态均非 22-r1 引入（该测试不在 24 枚文件内、交付 commit 纯新增 0 deletions），不记账。

## 二、逐格 AC 裁决（AC＝票面六格；派单五面映射在每格点名）

### AC1 — SSRF suite … **PASS**

尺＋读数（全部绿，见 1.1/1.2）：`Test22SSRFDecodeMatrix`＝25 私有行（七枚冻结网段＋inet_aton 全族：十进制/十六进制/八进制/混合 dotted/两段/三段/链路本地两形/v4-mapped）＋5 公有反控双向断；`Test22SSRFDeniedBeforeConnection`＝真 TCP 收方 accept 计数、7 种字面/编码 URL＋两枚名字解析（→127.0.0.1、→169.254.169.254）**receiver hits=0** 断死；`Test22DNSRebindingBlockedPerHop`＝同名字 hop1 公有 hop2 翻 10.99.99.99⇒拒、第二收方 0 命中、lookups≥2（每跳重解析）；`Test22GuardAgreesWithR5OnLiterals`＝R5↔netGuard 九枚字面双向一致＋四枚编码形 net 层必接；`Test22URLLengthCap`＝2048 过形闸、2049 连接前拒；`Test22SchemeMatrixAtGuard`＝ftp/file/mailto/javascript/data/urn/无scheme 全拒。
产码现读（非转述）：`webnet.go` `dialContext` 解析→判→**钉已检 IP 再连**（hostname 永不交 dialer＝rebinding 在 pin 处关门）；`blockedAddr` 缝生产构造器不设＝composition root 无开关可关守卫。

### AC2 — Redirect chain ＋ Markdown ＋ spill … **PASS**

`Test22RedirectDownsgradeBlocked`（TLS 起源→http，降级目标 hits=0）／`Test22RedirectToNonHTTPAndPrivateBlocked` 四子行／`Test22ChainLength`（cap 切断＋两跳短链 Origin＝最终 URL）全绿；Markdown＝`Test22FetchHTMLToMarkdown`（script 体整体丢弃、可见注入句作为数据保留＝票 19/25 形状、实体解码不复活活标签）＋`webmd.go` 现读：`dropEntirely` 表＋`skipUntilClose` 子树消费＋F2 belt（解码后拼出标签的文本再转义）⇒ **F2「绝不以 HTML 呈现」按构造成立**；大响应＝`Test22FetchHugeResponseCappedForSpill`（1 MiB 硬上限复用 `refRawOutputBytes` 标度＋喂真 `agent.Spiller.Prepare` 落盘 Spilled=true）；`Test22FetchStatusNoBodyEcho`＝状态码不回显正文、审计零正文（body never logged 有尺）。

### AC3 — Protocol matrix 16.5.3 穷举 … **PASS**

`Test22ProtocolMatrix`＝**17 行×两门**（web.open 只 http/https；file.open 白名单＋默认拒），每行正反双查（允许⇒恰 1 次 launch 且 kind 对；拒绝⇒launches=0＝"未达 launcher"断死）；`Test22DenyBeforeAllowOrder`＝六枚具名理由逐字（CLSID/脚本宿主/脚本协议/命名空间/明确拒绝）——注册≠授权（search-ms:/mshta: 在注入名册里"已注册"仍拒）；`Test22FileOpenPathArms`＝真 C26：目录内开/目录外拒/`file:///C:/` 驱动位/远程 host 拒；`Test22UnknownSchemeDefaultsToDeny`＝9 枚发明协议全落拒侧（白名单形非黑名单形）。

### AC4 — app.launch under Job；AUMID/name/path 三路 … **PASS**

`Test22AppLaunchNotepadUnderRealJob` 本机真跑（path 臂＋name 臂，TreePIDs≥2，Close 后轮询 pid 退场＝kill-on-close 亲验 0.15s）；AUMID 正形收／6 枚畸形拒；跨平台形状五枚全绿（argv 向量逐字钉⛔字符串拼接＝D33、Job 失败上抛、**Jobs 未接时结果行必说"未纳入"＝诚实姿态**）；AUMID 载体入 Job、UWP 本体归 OS 激活宿主＝平台事实直写结果文字（`launcher_windows.go` 现读，非粉饰）。`AppResolver` 缝留给票 23 `search.apps`＝票面并行槽 B 范围内、未越界。

### AC5 — Taint→L2（R4）＋ 首访提示一次、可记住 … **PASS**

攻击链本尊 `Test22AttackChainFetchThenTaintedSearchEscalates`：fs.read 打 taint→web.search query 携带片段⇒卡片 L2、RulesHit∋R4、Reason 点名来源、**provider 调用计数=0**（效果未发生）；`Test22PrivateTargetRefusedByTheBridgeLayerToo`＝双层（bridge R5 先拒＋工具直调再拒）。首访七枚全绿：一次／per-domain／remember=false 续弹／拒＝连接不发（收方 0）／**通道未接＝fail-closed 全拒**（NoGate 同姿态）／IP 字面不弹域提示／JSON 档跨进程持久／坏档多弹不默许。
⚠ 存储落点＝`WebDeps.DomainStorePath` 数据目录 JSON 文件，票面原文"config/db"二选一均未取——理由成立（`[net]` 冻结 schema 与 SPEC-02 冻结 DDL 无此位，加键/加列＝改契约文字须人工批，D22），登记为**待人拍板**，不构成 AC 不成立。

### AC6 — web.search 实现路径 … **PASS（带注）**

票面"OPEN decision"已由台账 **R4** 裁定（`pending-and-issues.md:50` 本腿现读逐字：「H4 web.search 路径已定：搜索 API…票 22 的 `SearchProvider` 接口按 anysearch 兼容实现为默认 provider」；接口面逐字段对 **R9** 登记面核过）；`Test22AnySearchWireContract`（路径/字段/匿名不发 Authorization/max_results 钳 10/results+items 两形/R7 key 不进输出）＋`Test22SearchProviderSwitchDoesNotTouchCallSites`（换 provider 不动调用点＝票面硬约束）＋`Test22ScrapeProviderIsExplicitlyUnlanded`（抓结果页＝显式工具级错误，**不塞死码**）。注＝AC6 框的翻勾动作归编排者；真网黄金流留票 12 接线后（残档§四.2 已具名）。

### 零回归（派单④）… **PASS**

整包 `-run ''`（全 299 枚）ok exit=0；Test22 名册 race ok；`go vet ./internal/tools` exit=0；gofumpt v0.12（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`）对 24 枚 commit 内 `.go` 文件 `-l` 输出空。写面纯新增 0 deletions ⇒ 既有 253 枚他票用例同树同绿＝回归作差为空。

### 纪律（派单⑤）… **PASS**

| 核 | 读数 |
|---|---|
| commit↔题面 | `git show --stat 547d9da00`＝24 files +5014，题面点名面全在（web*.go／webnet.go／launcher／docker/mockweb／证据件），`git log 547d9da00..HEAD -- internal/tools`＝**空**（无夹带追加） |
| 写面边界 | 全仓 `a18ff4546..73ab5da6b` 改动面＝22-r1（24 枚）＋50-r1c（internal/plugin 五枚＋台账＋其证据件，**他腿地界、本腿零评**）⇒ 22-r1 写面＝`internal/tools`＋`docker/mockweb`＋compose 追加＋证据件，⛔越界零枚 |
| thresholds/golden/cmd/wisp | `git diff --name-only` 两把尺均 **0 files**（逐字核） |
| compose 既有服务 | `git diff` 对旧行 **0 删除**＝纯追加；mock-web/ssrf-target 落 `manual` profile，CI `test` profile 零变化 |
| d22scan | `tools\d22scan\d22scan.exe -root <repo>` 本腿实扫＝**clean**（358 production files；ban#7 internal/tools 34 files 含本票新文件）exit=0 |
| git 纪律 | 单枚 commit、无 amend/add -A 痕迹（name-status 与题面逐名相符）；⛔push（HEAD 未推＝编排者地界） |
| 派单存档差异 | 交付腿如实登记 `22r1-dispatch.txt` 盘上不存在——本腿复核 `.scratch/wisp/probes/orch/20261102/txt/` 亦无该件 ⇒ 诚实登记成立，记**编排者地界**一注（先存档再发规矩的缺口），不记交付腿账 |

## 三、反向判据（摘防护⇒点名用例红）＋还原核 —— 五发全咬

做法＝`go test -overlay` 把仓外变异文件（`D:\tmp\wisp\20261102\22v1\mut\`）换进编译，**仓内一字节未动**；基线九枚点名用例先行全绿（21:5x）。

| # | 摘掉的防护 | 结果 |
|---|---|---|
| X1 | `ssrfBlockedAddr` 恒放行（＝"放行私网"） | **8 枚红**：DecodeMatrix／DeniedBeforeConnection／DNSRebinding／GuardAgreesR5／UnbracketedV6／URLLengthCap／RedirectToNonHTTPAndPrivate／PrivateTargetRefusedByBridge ⇒ 本票核心判据"摘 SSRF 守卫⇒点名用例红"**咬** |
| X2 | `decodeNumericHost` 不认编码形 | **6 枚红**：DecodeMatrix／GuardAgreesR5／UnbracketedV6／RedirectToNonHTTPAndPrivate／PrivateTarget／FirstVisitIPs ⇒ inet_aton 解码臂有牙非装饰 |
| X3 | checkRedirect 的 https→http 降级拒摘除 | 首发编译红（`prev` declared-and-not-used，读数件 `mut-x3.txt` 留档）；修正后二发 **恰 1 枚红**＝`Test22RedirectDownsgradeBlocked`，断言分域不伤他行＝交付腿 M2 复现 |
| X4 | `explicitDenySchemes` 命中即拒摘除 | **恰 1 枚红**＝`Test22DenyBeforeAllowOrder`＝交付腿 M3 复现（具名理由逐字才咬，泛"拒绝"不咬——交付腿自抓那课为真） |
| X5 | 首访 ask==nil fail-closed 改默许 | **恰 1 枚红**＝`Test22FirstVisitUnwiredFailsClosed`（0.63s＝真拨号尝试发生＝防护在场证据）＝交付腿 M5 复现 |

还原核＝撤 overlay 复跑同批九枚 **ok exit=0**＋Test22 全名册 **ok 3.088s**；三枚被摘文件 `git hash-object` 对 `git rev-parse HEAD:<path>` **逐枚 same=True**（webnet `2c9dbb88`／web_open `a38247f9`／web_firstvisit `1f06be8b`）；工作树 tracked-dirty＝**0**＝本腿零污染。
未复刻 M4/M6（launcher Job 臂、webmd script 臂）：X1–X5 已覆盖派单点名的 SSRF/降级/矩阵/首访四面；M4 属会真起 notepad 的桌面侧（交付腿已自报孤儿进程事故并清），本腿在无人值守窗不重演——具名为**未复核格**，不留白谎报。

## 四、残档核（交付腿§四八条逐条对盘，不得无声消失）

1. `[net].allowlist` 生产者半程归 `cmd/wisp`（禁区）⇒ unwired 守卫行一字未动＝诚实——本腿现读 compose/config 面未越界，成立。
2. anysearch 真网黄金流留票 12＝成立（R7 缝已就位）。
3. 首访原生卡归票 21/37、未接＝全拒已测——成立。
4. UWP Job 边界＝平台事实直写——成立（§AC4）。
5. 编码形/字面量 R5↔net 分工钉成测试——成立（X2 反向亦咬）。
6. CGNAT 100.64/10 按面放行＋公有反控五枚在册＝未擅自扩冻结清单，方向正确（fail 朝清单）。
7. mock-web/ssrf-target manual profile、Go 套件才是牙齿——与 AC1/AC2 真跑互证，成立。
8. `doc.go` 那句过期文字：交付腿留给裁决腿随结票一并正——⚠ 改 `internal/tools/doc.go` 在写面内、但属翻勾同批动作，**随结票由编排者定**，本腿零改码只具名。

## 五、缺陷与缺口清单（具名，无则零）

| # | 项 | 定性 | 归口 |
|---|---|---|---|
| D1 | 派单字面尺 `Web|Guard|…` 漏扫 26/46 枚本票用例（含 AC4 全族、AC1 核心两枚） | 派单尺缺陷，非交付缺陷；本腿已补跑整名册零红补洞 | 编排者（今后点名尺建议 `^Test22[A-Z]`） |
| D2 | 证据件 §五 末行"Commit：见下方枚字（落地后回填）"**未回填** | 文档完整性小缺陷（历史枚字可 `git log` 现取＝本腿已代取：单枚 `547d9da00`） | 交付腿随结票补一行或作废该行 |
| D3 | `Test22DNSRebindingBlockedPerHop` 注释说 hop2 翻 `192.0.2.1`，代码实翻 `10.99.99.99` | 文字漂移，测试更强非更弱 | 交付腿，微 |
| D4 | 首访存储 JSON 文件 ≠ 票面"config/db"字面 | 待人拍板（契约文字面须人工批，交付腿不擅动正确） | 机主/编排者一句裁定 |
| D5 | 交付腿 M4 事故自报：notepad 孤儿进程疑误杀窗（已 Stop-Process＋面板复零）＋重建文件过程披露 | 如实登记合格；机主如需核一眼 config.json 可自行确认 | 机主知情 |
| D6 | 派单存档 `22r1-dispatch.txt` 不存在（先存档再发规矩缺口） | 交付腿如实顶回，成立 | 编排者 |
| D7 | 本腿未复刻 M4/M6 两发突变（§三末具名） | 复核缺口自报，不影响五发已咬面 | 如需可补 |

## 六、结论

**是否同意结 22：同意——六格 AC 逐格判据成立（AC6 带 R4 具名注），零回归、纪律、反向判据（五发全咬＋还原核 blob 逐枚 same）全部 PASS；D1–D7 七条均为注记/归口件，无一构成"AC 未满足"。** 翻勾与推送归编排者；D4 建议随结票一并摆机主一句裁定（JSON 落点是否即"config/db"之满足），其余按归口带办。

——验收腿 22-v1 · 裁决者≠实现者（SPEC-12 §4.3）· 只 commit 不 push · 本腿写面＝本文件一枚（显式 pathspec）· 临时件全在仓外 `D:\tmp\wisp\20261102\22v1\`（只建不删）
