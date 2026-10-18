# 347 — POSIX 侧 `resolveHandle` 那枚**恒 `("",false)` 的桩**让黑名单的 fail-closed A 那一支在本仓**受支持的平台**上结构不可达＝产码缺口（⛔ 属票 55 的 macOS 移植；CI 的 `test-core` 今天就在 Linux 跑这个包）

- Status: **已立，未派**（2026-10-18 22:5x 编排者立；来源＝票 339 终裁腿 `339-v1` 必答 **(iv) 判"该立票"＋把分母量死**（判语件 `docs/evidence/s1/339-v1-verdict.md:60-64`）＋编排者自己复跑坐实（台账 `A1042`）；触发链＝票 339 那两枚仪器红的**真根因**）
- 类型：**安全必修邻接**（`D33` 那一族，C26 的 POSIX 半边）；★**产码活**＝不是仪器活 ⇒ 落地腿必须是写腿，且**开工前须机主拍一次"这一支现在要不要真落地"**（见 AC#0 前置那节）
- 射程＝`internal/risk/pathresolver_other.go`（桩本体）＋`internal/risk/pathresolver.go`（`formsOf` 那一跳）＋三枚受它影响的测试件

## ★⛔ 并入党（我裁：⛔ 并进票 55，与裁腿的建议分歧，理由三条）

`339-v1` 必答 (iv) 建议"与票 55 `S8 macOS port` 的 `AC :28`（`PathResolver macOS branch (realpath + lstat) passes red-team port.`）**同体**⇒建议并入或显式引用"。我现读三处后**不采纳并入**，本票独立成票：

| 尺 | 读数 | 后果 |
|---|---|---|
| 票 55 的阻塞条件（现读 `55-s8-macos-port.md` 第 5–7 行逐字） | `Blocked by: 54-s7-acceptance + self-use week + explicit user go-ahead (S8 starts only pre-open-source, D23)` | 并入＝这活**推到开源前才允许启动**；而 `test-core` 现在每天在 Linux 跑这个包 |
| `D7 平台范围`（`docs/PLAN.md:184` 现读逐字） | `Windows 一等公民 + macOS 二等 + Linux 不入官方支持矩阵` | ★macOS **今天在支持矩阵内**＝"macOS 那一支的安全语义没落地"＝**现存的契约缺口**，⛔ 只是远期移植项 |
| 桩本体自陈（现读 `pathresolver_other.go` 首 10 行） | `// DEFERRED(macOS/Linux): SPEC-06 §4 requires realpath + lstat symlink reject for non-Windows platforms. Until that ticket lands, handle-based resolution is unsupported here …` | ★它自陈 `macOS/Linux` **两平台**；票 55 那条只写 `macOS branch`＝**射程比那一格宽**（Linux 那一半没主人） |
| `scripts/portable-tests.sh` 现读 | `internal/risk` 在 `core` 档射程内（该文件 `:266`/`:282` 两处逐名点名）＋`:7-8` 逐字警告"Go 把 SKIP 记成 `ok`" | CI 今天就在 POSIX 跑这个包＝这一支可达不可达**不是纸面问题** |

⇒ **去向**：本票**显式引用**票 55 `AC :28`（同一枚 SPEC-06 §4 要求、同一族落点，⛔ 两票同波动 `internal/risk/pathresolver*.go`）；⛔ 把本票折进它的阻塞条件里。

## 现量（编排者 2026-10-18 22:5x 自己跑；引用前先重跑）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| 桩本体 | `internal/risk/pathresolver_other.go:10` 逐字 `func resolveHandle(p string) (string, bool) { return "", false }`，文件首行 `//go:build !windows` | `git show HEAD:…` 现读 |
| 唯一赋值点 | `pathresolver.go:345` 的 `formsOf` 里 `f.root` **只在** `:354-357` 的 `if final, ok := resolveHandle(cut); ok` 成功分支内被赋值 ⇒ POSIX 上 `pathForms.root` **恒空** | `339-v1` 用 `git grep '\.root = '`（量的、包内非测试仅此一枚）＋编排者复跑同锚 |
| 结构不可达 | `blacklist.go:359-361` 的 `uncertainAnchorMiss` **首行**即 `if cand.root == "" || anchor.root == "" { return false }` ⇒ A 档（R17 fail-closed）在 POSIX **永不触发** | 现读＋`339-v1` 探针 `guard-expr=true`（返回就在首行） |
| ★假主目录那侧⛔ 是原因 | `anchor_miss_reason_334_test.go:85-88` 同发设四枚 env（`USERPROFILE`／`HOME`／`APPDATA`／`LOCALAPPDATA`）＋探针 `fakeHomeEffective=true` ⇒ 票 339 票面第 9 行那句"`USERPROFILE` 在 Linux 不生效"＝作废（编排者已写错一次，账 `A1040`） | `339-v1` 复现；⛔ 新腿再按那句做诊断 |
| 落地后的分母（裁腿量死、⛔ 我复跑） | `stub→realpath` 最小预览形后 Linux 整包 **196/112/5/0 rc=1**；红 5 枚逐名＝`TestFallbackBranchHasNoWorldOnThisPlatform336`(`:76`)／`TestAnchorMissHasNoWorldOnThisPlatform339`(`:140/:143`)／`TestSyncNoGradeIsConfirmedOnPosix`(`syncdirs_other_test.go:64`)／`TestSyncMembershipDecidesOffProfilePosix`(`:95/:100`)／`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（存量 `A380`）；另 `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet` 由 **SKIP→PASS**＝守卫自解 | `339-v1` 的 `logs/v1-posix-resolver-lands-preview.log`（34276 字节）；★**落地腿起手必须自己重打这一发**，⛔ 抄这行 |

## ★前置：这一支今天到底要不要真落地＝**待人拍板**（`Q` 级，我先摆给机主、⛔ 替他选）

三条候选，代价各不相同（我只登记，⛔ 选）：
- **一＝落解析器**（本票 AC#1 甲形）：POSIX 上真做 `realpath`＋`lstat` 拒 symlink ⇒ A 档那一支在 macOS/Linux 变得**可达**⇒ 票 339/336 那两枚"此平台无世界"的绊线**要搬**（⛔ 删）。⚠ 代价＝`339-v1` (iv) 明写：★POSIX 上 `ESC`／`U+202E` 是**合法文件名字符**⇒334 那形靠"Windows 拒绝解析"造 partial 的机制**在新台面未必复现**⇒⛔ 承诺"搬牙即绿"。
- **二＝只把桩的自陈与票面口径钉准**（零产码）：不动解析器，只把 `DEFERRED(macOS/Linux)` 那枚标记与 `SPEC-12 §5` 登记表**逐字对撞**（★硬规矩＝标记与登记表必须 1:1 双向对得上，`AGENTS.md` §1.1）＋把票 339/336 那两枚绊线保持会响。⇒ 敞口照旧，但"没人知道它开着"这件事被消掉。
- **三＝不做**：维持现状。⇒ 后果具名＝`D7` 说 macOS 在支持矩阵内，而 C26 的 POSIX 半边（SPEC-06 §4 那两行要求）**今天零实现**，只有两枚会响的绊线在替它作证。

## 判据（逐枚可机读）

- [ ] **AC#0 前置**：编排者把上面那一／二／三摆给机主、拿到具名回答**之后**才许开工（⛔ 任何腿自己挑一形）。★起手还要现跑：① 上表五行逐枚重跑（一律取具名提交树）；② `DEFERRED(macOS/Linux)` 那枚标记 ⇄ `SPEC-12 §5` 登记表（`docs/specs/SPEC-12-roadmap-governance.md:66` 那行）现读对撞，具名说清"这一枚标记今天有没有登记面"（⛔ 新增 `DEFERRED(D-xx)` 标记却不登记＝直接违规）。
- [ ] **AC#1 三形代价表（⛔ 未出表不许动码）**：**甲**＝`resolveHandle` 落 `realpath`＋`lstat` 拒 symlink（SPEC-06 §4 原文口径）；**乙**＝只在本包内把 `f.root` 换成**另一条不依赖句柄**的解析路径（例：词法链＋逐组件 `Lstat`）⇒ 同样让 A 可达，⛔ 动 `os.ReadDir`／跨包；**丙**＝零产码（＝上面"二"）。每形现量：改动面（文件:行）／POSIX 与 **Windows 各自**的实测读数（⛔ 一发代两发）／会不会顶红票 336 那 **8 枚牙**与票 216／334／339 的既有仪器（逐名读断言）。
- [ ] **AC#2 ★⛔ 单独摘守卫**（`339-v1` (vi) 的补强，我原句搬进本票当判据）：**⛔ 只把 `blacklist.go:359-361` 那枚首行守卫摘掉而不落地解析器** ⇒ `isUnder("","")` 恒真会把 POSIX **反向全量 fail-closed 成 A**＝把"读不到"变成"全都拒"。⇒ 判据＝去牙那一发必须**点名这一形**（若腿的落地形在缺世界时退化成"全 A"，仪器要红、红句要归位）。
- [ ] **AC#3 搬牙不删牙**：票 336 那两枚仪器（`anchor_fallback_teeth_336_windows_test.go`／`..._other_test.go`）＋票 339 那枚绊线（`anchor_miss_reason_339_other_test.go:140/:143`）＝**设计好的"搬我"信号**⇒ 新台面下要**改其射程**，⛔ 删除、⛔ 让它静默绿。★`name_claim_two_arm_216_test.go:89` 逐名引用夹具 `candidateOffAnchorChain334`（`anchor_miss_reason_339_other_test.go:36-43` 的注释已具名该合并条件）⇒ 合件时必须**保名**。
- [ ] **AC#4 颜色与门禁**：`./internal/risk/` **两平台成对**（Linux 容器 ⇄ 本机 Windows 档，同一 `-run` 名册、四数逐枚点名⛔ 混用⛔ 相减、改前⇄改后同台面）＋红名册**逐名**作差逐名归位（★上面那 5 枚红逐枚点名，⛔ 拿"枚数差不多"交差）；门禁四数只在共享工作树、★跑前 `unset GOTMPDIR TMPDIR`（票 311 AC#1 shape B GUARD，带 env 时 rc=2 射程拒绝⛔ 读成红）。⛔ `./cmd/wisp` 整包、⛔ `go test ./...`。
- [ ] **AC#5 ★零新导出名＋零契约移动**：⛔ 新导出名／新 JSON 键／新 `panel.*` 方法名；⛔ 动 `C26` 的**入口语义**（本票补的是 POSIX 半边实现，⛔ 改判据＝改契约须人工批准）；⛔ 动判级档位／豁免键语义／`FILE_ATTRIBUTE_REPARSE_POINT` 那半（Windows 侧一字不动）；三枚冻结件／`internal/observe/thresholds.go`／golden／`allowlist.txt`／`go.mod`／`go.sum`／`.gitignore`／`.github/**`／`ci.yml`／D1–D47／C1–C32／R1–R9／D43 转移表；⛔ `frontend/**`／`design/**`；⛔ 为变绿放宽票 334／336／339 的任何既有断言。
- [ ] **AC#6 纪律格**：⛔ 翻任何框（本票七格由非实现者裁）、⛔ 改票名／`Status:`／加 `-done`；写面＝`internal/risk/**`＋`.scratch/wisp/probes/347/r1/**`（短名、⛔ U+00A7 那枚字符、逐枚非 0 字节、⛔ `.out` 扩展名）＋本票末尾**一条**只追加 bullet；`add` 与 `commit` 两处 pathspec 逐枚手写字面路径（⛔ 从共享 index 反推＝`A1027`／`A1029`／`A1040`）；中文正文⛔ 走 PowerShell 内联（`A1030`）；⛔ push。

## 禁区（补充）

⛔ 与票 55 同波动 `internal/risk/pathresolver*.go`（同一枚 SPEC-06 §4 要求，两票并发＝互相洗脏读数）；⛔ 与票 342（`internal/winsec` 的祖先那一支）混记＝**不同包不同机制**，但★两件事是同一族"POSIX 只护到某一处就停"，落点腿读数要互相引用一次；⛔ 回头改票 339 的判据文字（它已 `-done` 或将 `-done`，要更正只追加）。

**撤销口令**：「撤 347」（⇒ 本票退回未派；两枚绊线与桩字节不受影响＝本票⛔ 动过码）。

## 为什么现在不派

①**AC#0 那一道待人拍板今天没开**（机主未答"这一支要不要现在落地"）；②`internal/risk/**` 此刻由本波在飞腿占着（`341-r1` 交件后还要派 `341-v1` 裁，同包⛔ 两腿抢读数）；③与票 342 的落地腿**同族不同包**，要一起排宽度帽≤8。

## Progress log (append-only, newest last)

- [2026-10-18 22:5x +08] 编排者立票：来源＝`339-v1` 必答 (iv) 判"该立"＋分母已量死；★我**否决并入票 55**（三把尺现读：票 55 阻塞条件／`D7` 平台范围／桩自陈射程含 Linux）并改为显式引用；把 (vi) 的补强"⛔ 单独摘守卫"搬成 `AC#2` 硬判据。名 94 字符≤100、★开票前跑过两门（长度门 VERDICT GREEN／d22scan rc=0 clean）＋点名路径逐枚 `git cat-file -e HEAD:<path>` 现验（9 枚里 8 枚 rc=0；★`docs/specs/SPEC-06-security-gating.md` 不存在＝真身是 `SPEC-06-security-gatekeeping.md`，我照例先验再写）。
