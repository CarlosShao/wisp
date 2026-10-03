# 212-v1 验收裁决（第二任骨架 `234e7301`＋第三任取数满稿；非实现者）

- 验收者：212-v1 第三任（第二任死于 §2 探针中途；其骨架、`mut/`、`fix/`、基线复用，**"待填"一律由本程自取数填满，前任读数不当作数**）。
- 实现者：编排者本人（A544 预授权代笔，台账 A578，锚 `5e8748b3`）。按 D22 双角色本件以非实现者攻判据；勾框零碰，产码零改。
- 本程起手锚 `de82ec76`；date 起手 2026-10-03 21:40＋0800、交件 22:3x＋0800。⚠ 共享树多腿并发在飞（258-v1、260-a2），本程读数带时限，差异逐处具名（§3/§6）。
- ⚠ **三实例碰撞（归编排者裁，本程只记事实）**：本票被并发派了至少两枚 v1 实例——"第一实例"（21:08 起手、未死，全量裁决落同目录 `verdict-first-instance.md`，其尾指针 22:09 加进本文件）与"第二任"（21:19 骨架、`234e7301`）。A578 只记派一枚。基线 `gate-baseline.txt` mtime 21:14 早于第二任起手、与第一实例作者声明一致 ⇒ "前任死于服务错误留基线"的骨架叙述与树不符（§6 #6）。两份裁决独立取数、缺陷集合一致；**本文件是派单指定的落盘面**，第一实例档案完整保留、不删不改。
- 纪律：frontend/**、design/** 零读零引；跟踪树零改动（探针走仓外夹具＋`-overlay`）；只 commit 证据件、不 push。

## §0 起手锚与基线

- 复用件：`gate-baseline.txt`（mtime 21:14:23；34 PASS/76 RUN＋全仓 clean，八枚 ban 读数 internal/=228 cmd/=38 frontend/=85 tools(=#7)=23 design/=39 internal/#8=498 cmd/#8=96）。
- 第二任门禁第一发 `v2-gate-run1.txt`（21:24）：rc=0，34 PASS/0 FAIL/76 RUN，与基线逐名相等。
- 本程门禁第三发 `r3-gate-run3.txt`（21:54）：rc=0，34 PASS/0 FAIL/76 RUN（读数见 §3）。第二发未由本程重跑，以第一、三发夹逼。
- ⚠ 夹具探针件 `p10.go`（缩写形）系**第一实例** 21:41:09 所建（与 `fix/scan-ellipsis.txt` 同秒，作者归属由此落定，非来历不明）；第二任 `fix/scan-pristine.txt`（4 findings，21:31）跑在 p10 落盘之前，对 p10 无覆盖力，由本程 `r3-fixture-run1.txt`（5 findings）与第一实例 `scan-ellipsis.txt`（同为 5）双实例一致取代。

## §1 恒真两问（第二任跑两发；本程全量复核②、源级复核①）

- ① **摘发射点**（`mut/main-noq9.go`＝删 `for _, cg := range f.Comments {` 至 `return nil` 整段）：本程源级复核——变异源 `f.Comments` 残留恰 1 处（ban #9 的块注释），发射点确已整段删除。读数（第二任 `noq9-out.txt` 全仓 clean／`noq9-selftest.txt` 35/36，唯一 FAIL＝ban #9 ring 向；第一实例另喂突变名册得名册审计 FATAL×2 HOLE，`mut/noq9-selftest-mutroster.txt` 本程读过：断肢自报）——三份互证：ban #9 只有这一处发射点，摘门即哑且**缺牙自己会叫**；而摘门后全仓照绿 ⇒ "绿"不证明牙活，必须配②。
- ② **反形判据**（`mut/main-inv.go`：`if exists{continue}` 翻成 `if !exists{continue}`）——**本程 `-overlay` 全重跑**（`r3-inv-out.txt`）：真仓 rc=1、**160 枚 phantom-citation**，首屏红句 approval.go:24 引 `internal/llm/ratelimit.go`——**该文件实存**（ls 在案），反形把活引注定罪；`diff mut/inv-out.txt r3-inv-out.txt` 仅两行差（ban #8 cmd/ 分母 96→97＝§3 具名并发件；尾行 exit status），160 枚红句逐字相等。self-test 向（第二任 `inv-selftest.txt`）：ring FAIL＋silent FAIL（35/36），红句与本程净版 `go run . -self-test`（36/36 OK）的样本逐字可对。⇒ 判据是**存在性双向**：反形把"存在的引注"全数定罪＝self-scan 与 seeded-citations 两向都被咬，正向判据不是恒真。
- 结论：恒真两问均不成立——判据有牙、牙朝正确方向、断肢自报。

## §2 symRefRe 三发探针（本程亲手复跑夹具；读数 `r3-fixture-run1.txt`，与第一实例 `fix/scan-ellipsis.txt` 独立同读数）

夹具 `fix/r2probe`：12 枚产码 Go 文件（internal/ 11＋cmd/ 1），rc=1、5 findings、其余全静默。

- **(a) p1** `// probe a: api citation internal/tools.Result.AppliedSteps` → **静默**。推演：repoPathRe 抽出 `internal/tools.Result.AppliedSteps`，symRefRe `^[a-z][a-z0-9]*(/[a-z][a-z0-9]*)*\.[A-Z]` 以 `internal/tools`＋`.`＋`R` 命中 → 排除。实测＝推演（两实例一致）。
- **(b) p2** `// probe b: the error class internal/tool` → **红**（p2.go:1）。推演：token 无 `.[A-Z]` 尾，symRefRe 不匹配 → 不排除 → stat 不存在 → 红。实测＝推演。⇒ 该 token 的消红**靠 5e8748b3 改注释**（tools.go:89），不靠正则；同时坐实 main.go:802-805 把它列为"symRefRe 排除例"的仪器自述与行为不符（§6 #3）。
- **(c) p3** `// probe c: see internal/Build/x.go for details` → **红**（p3.go:1）。推演：`internal/Build` 第二段大写 B，symRefRe 目录组只吃小写 → 不匹配 → 照常 stat → 红。**实测＝推演＝实现者点名的嫌疑格不成立**（不误排除）。对照 p4（全小写 `internal/build/x.go`，同样没种）也红——排除只由 symRefRe 决定、存在性照常判。
- **控制四发**：p6/p7（夹具种下的 `tools/d22scan/main.go`、`internal/tools/bridge.go`）**静默**＝存在性如实放行；p8/p9 干净件静默。
- **(d) p10（第一实例所种，本程自取读数）** `// readings in docs/evidence/s1/152-...-accept-r1.md` → **红**（p10.go:2）。推演：repoPathRe 字符类含 `.`，`...` 照吞、全拼 token 入抽、stat 失败 → 定罪。**与名册第 9 行/finding 文案自称的 "shorthand … NOT a violation" 相撞**（§5 AC#3）。该形唯一活体标本在 `cmd/wisp/slo_report_144_windows_test.go:851`＝`_test.go`，ban #9 搭的 `walkGo` 对 `_test.go` 显式 return ⇒ 今日全仓 clean 是射程豁免垫着，不是判据放行。
- **射程钉**：ban #9 只经 `walkGo → scanGoFile`（main.go:686 唯一调用点；walkGo :680 一带跳 `_test.go`）⇒ 射程＝产码注释面，与名册第 9 行一致。
- **假排除普查**（`symref-sweep.txt`）：5129 枚 tracked 路径中 symRefRe 全局命中仅 `docker/builder.Dockerfile`、`docker/mockllm.Dockerfile`——`docker/` 不在 repoPathRe 前缀表 ⇒ **可达误排除面＝0 枚**（今日零实伤的已量边界）。

## §3 AC#5 八枚 ban 对账

- 本程第三发：rc=0；PASS=34/FAIL=0/RUN=76；bans #1-5 internal/=228、cmd/=38；ban #6 frontend/=85；ban #7 internal/tools/=23；ban #8 design/=39、frontend/=85、internal/=498、**cmd/=97**；verdict clean。
- 与 A578 记录（…internal/#8=498 cmd/#8=96）及基线逐名比：**唯一差＝ban #8 cmd/ 96→97**。定因具名：并发 258-v1 腿 21:45:13 落的 untracked 测试件 `cmd/wisp/resident_hotkey_v1probe_test.go`（mtime 在案；A578 记于 20:5x、基线捕于 21:14、第二任发于 21:24，皆在其前）。ban #8 分母数 `_test.go` 且本机扫描连 untracked 一起读（A207 机器相依分母口径）——**非实现件所致**。
- 正控：TestScanDetectsAllSeededViolations PASS（run3 在案）；`cd tools/d22scan && go test ./` → ok 16.187s（本程自跑）；`go run . -self-test` → **36/36 clean（20 ring/16 silent）**，ban #9 两向在内（ring＝引 `docs/evidence/s1/212-citation-ruler.md` 未种必响；silent＝种 `docs/readings.md`＋`internal/probe/roster.md` 后静默＋散文式短引不问罪）。⇒ 既有 8 枚 ban 的牙未被 ban #9 挤掉。**AC#5 成立**（96→97 已具名并发件）。

## §4 注释更正逐枚复核（`git show 5e8748b3` diff 逐 hunk 读；**212 本票集＝8 文件 9 hunk**，下表 #1-#8；#9/#10 是 258-r1 残局搭车 hunk、标不适用）

| # | 文件:处 | 更正前→更正后 | 真话？ | 凭据 |
|---|---|---|---|---|
| 1 | internal/risk/pathresolver.go:28 | `scripts/check-pathclean-ban.sh`→`tools/d22scan 的 pathresolver-bypass ban; see main.go` | **真** | 全仓 find 该脚本＝0 枚（原话确是幻影）；ban #2 tag 在 main.go:15 名册＋:726 发射点俱在 |
| 2 | internal/agent/approval/pending_read.go:43 | `ticket 17's frozen-contract note in tools/gate.go`→`…note (enforced by tools/d22scan's pathresolver-bypass ban since the scanner landed)` | **疑（实质半句真、路标丢失）** | 冻结注真身＝**internal/tools/gate.go:13**（"renders RulesHit and Reason VERBATIM (ticket 17's frozen-contract note)"）——原引只缺 `internal/` 前缀、一词可修；新文却换成"由 pathresolver-bypass ban 执行"，而该 ban 管 Clean/Abs（main.go:15/:726），**不执行**票 17 逐字显示契约 ⇒ 新断言不成立＋导航路标没了 |
| 3 | internal/agent/approval/queue.go（两 hunk） | `tools/bridge.go`→`internal/tools/bridge.go` | **真（两处）** | 文件在；orDefault(req.CorrelationID, req.TaskID) 于 bridge.go:519/:526，reject 分支 orDefault(why,…) 于 :445-:464 |
| 4 | internal/agent/tools.go:89 | `class internal/tool`→`error class internal-tool` | **疑（类名不实）** | `internal-tool` 不在 D37 17 类枚举（PLAN §16.7 表＝observe/errors.go:39 `allClasses`）；实测 booking：bridge Execute 裸错误经 loop.go:694 `ErrorClassOfTurnError`→guard.go:274 回退 **ClassInternal**。幻影 token 消了、"host 失败≠tool 失败"方向对，名目仍是造的 |
| 5 | cmd/wisp/models.go:20 | `internal/engines/ directory`→`engines directory under internal/` | **真** | internal/ 下无 engines（ls 实证）；新句正说"没有"，与盘相符 |
| 6 | internal/ball/statevisual.go:68 | `docs/evidence/s1/62-*`→`the ticket-62 evidence tables (…under docs/evidence/s1)` | **真** | `62-` 前缀 7 枚在；96 DPI 实见于 62-adversarial-acceptance.md／62-visual-spec-draft.md。⚠ 括号内同短语重复两遍＝编辑瑕疵，真而不顺 |
| 7 | internal/risk/provenance.go:63-64 | `reaches tools/agent/cmd`→`reaches the plugin agent command` | **真（带名注）** | tools/ 下无 agent/cmd；承重的是负断言＋引 `(160-c1 §2.2)`——160-design-core-c1.md :99-102 现量"生产码里没有任何一枚 *plugin.DisposalScope 走到桥、走到 loop、走到组合根"，核实相符。"the plugin agent command" 是松称、不可机读锚定（与第一实例存疑标签的分歧点，缺陷集合不变） |
| 8 | internal/tools/bridge.go:252 | `internal/provider`→`internal-provider` | **疑（类名不实）** | 同 #4：不在 D37 枚举；实测裸 Execute 错误 booking＝ClassInternal。方向（REJECT 不是故障）与 bridge.go:255 起正文相符 |
| 9 | cmd/wisp/config_readers_255.go hotkey 行＋白名单 | `:257 Hotkeys: ball.DefaultHotkeys()`→`:269 Hotkeys: cfg` | ——【258-r1 随行件·非本票证据】 | 顺手核了：resident_ball_windows.go:269 逐字 `Hotkeys:  cfg,`（awk 在案）；链 `residentBallHotkeyChain258`(:183)←resident_windows.go:175/:210、reloader 桥 :306；Esc borrow 残余归票 245 有名有姓（ball_windows.go:805-812）。真，但**翻勾权归 258-v1** |
| 10 | cmd/wisp/resident_approval_(live_)246_windows_test.go | `startResidentBall(reg,…)` 两参→四参 | ——【258-r1 随行件·非本票证据】 | 现签名 :249＝`(reg, onCancelEsc, hotCfg, hotReload, hooks…)`；4 处调用点与 arity 一一对应 |

- **计数：212 本票 8 文件 9 hunk＝5 真（#1/3/5/6/7）＋3 疑（#2/4/8）**；#2 实质半句真、名目半句假。三枚疑的共同形状：幻影 token 都消了（ban #9 不再红），但**新断言**各自不成立——票面禁区要求"把那句话改成实话"，这三枚改后仍不是实话，具名上报，修法归编排者（#2 诚实形＝`internal/tools/gate.go`；#4/#8 诚实形＝`internal`）。
- **与第一实例判表对读**：缺陷集合一致（#2/#4/#8）；严重度标签两处分歧（#2 疑↔真但带条件、#7 真↔存疑），不影响任何 AC 判语。
- **票 258 表处置（派单指令项）**：骨架 #9/#10 两行系 258-r1 残局随 5e8748b3 搭车（commit message 自把它们记在【258-r1】名下）——**判"跑错文件带的料"，选"标不适用"而非删行**（骨架保留＋禁删前人件纪律类推）；真话性顺手答了（都真），不计入本票 AC 凭据。

## §5 AC 格判语（勾归编排者）

- **AC#2（落点由非实现者裁）＝成立但带条件**。落点＝tools/d22scan 独立 ban #9（名册第 9 行、独立 tag/发射点、auditSelfCases 连续编号）。本验收腿追认该落点：与既有 8 枚同属产码注释面门禁、复用 parser＋walk 基建零新分母；§3 实测背书"逐名不变"；§1 证明有牙且断肢自报。**条件（程序瑕疵）**：票面原文"由非实现者裁"，而落点初裁（A574 翻勾节）与实现（A578 代笔）同为编排者一人——本腿现以非实现者身份追认补正。附带 cosmetic：ban #9 搭 bans #1-5 的 walk，scope 标签 "bans #1-5 internal/" 字面不再穷尽该 walk 服务的 ban 集。
- **AC#3（形状）＝成立但带条件**。②「声称读数在某件里而该件不存在⇒红」**全成立**：控制 p4/p5 红、pre-image 真②全数更正后全仓 clean、反形证明判据是存在性而非反向。**③「缩写先规定、别一上来就当违规」对省略号-全拼路径形不成立**：该形被 repoPathRe 照抽照定罪（§2(d) p10），与名册第 9 行、scanGoFile 注释（main.go:791-793）、finding 文案三处自述"shorthand … is NOT a violation"矛盾；self-test silent 样本只练散文式短引，恰好漏过这形。今日零实伤——唯一活体在 `_test.go`（射程外）。与第一实例"不成立"标签所指同一点，severity 归编排者裁。**修法二选一（仪器射程变更须走批准，本腿零产码改动）**：regex 对含 `...` 的 token 降③不问罪，或改注释/文案＋补该形 self-test 样本；判据本身不须回滚。
- **AC#5（不许伤）＝成立**（§3：除具名并发件所致 cmd/ 96→97 外逐名相等；正控＋self-test 36/36＋包测试 ok 16.187s）。

## §6 推翻清单

1. 票面「今天加牙不会打红任何代码（0 枚违反）」——维持推翻（census 翻勾节在案）；本程 `preimage-phantom-census.txt` 复算背书：`5e8748b3^` 266 产码 Go 文件、**11 distinct phantom token**。
2. commit/A578 计数「首发 10 枚分类」＋「10 枚真②全部改注释」——**与树不符**：11 token＝8 真②＋3 API 形（symRefRe 排除）；更正 diff 实为 **8 文件 9 hunk**（queue.go 双 hunk 同 token 两址）。两个 "10" 都对不上。
3. main.go:802-805 symRefRe 注释把 "class internal/tool" 列为排除例——**自述与行为不符**：该 token 无 `.大写` 尾、不被 symRefRe 匹配（p2 实测红）；本票治的就是撒谎注释，仪器自己的注释也得守规矩。
4. 名册第 9 行／finding 文案／silent 样本 note 三处「shorthand is NOT a violation」——对**省略号-全拼路径形**不成立（p10 红）；clean 靠 `_test.go` 射程豁免垫着、非判据放行（§2(d)/§5）。
5. 「十枚注释诚实化＝更正后都是真话」——**3 枚存疑**（#2 enforcement 嫁接＋路标丢失；#4/#8 枚举外类名）。
6. 第二任骨架「前任死于宿主服务错误、基线系其所留」——**与树不符**：基线 mtime 21:14 早于第二任起手 21:19，与第一实例作者声明一致；A577 所记"死于服务错误"是写腿 212-r1 非验收实例。重复派单＋死活叙述归编排者对账。
7. A578／基线「cmd/#8=96」——本程 **97**：非读数错、是**带时限读数**（并发 untracked 测试件 21:45 落，mtime 在案）；跨腿比分母必须带此刻注。
8. 「symRefRe 会误排除大写段目录」（实现者自点名嫌疑格）——**不成立**（p3 红）。

## §7 判不动／量不到

1. 反形逐字对账基于"除已知并发差外逐字等"（§1②）：第二任反形跑于 21:24、本程重跑于 21:53，258-v1 探针件（21:45）横在中间，同刻同树不可重建。
2. #9 收据行「wisp run builds no ball／读者在常驻进程」的结构断言：核到装配侧代码＋引用行逐字手核，**未起真进程复跑 stdout**（真进程档归 258-v1/winlive 惯例）。
3. cmd/wisp 整包 go test 未跑（纪律⛔全仓 go test；TestTicket255* 钉族不在本程配额）——:269 引用真伪由 awk 逐字手核补上（§4 #9），测试侧背书待 258-v1。
4. 摘门①未由本程重跑二进制（复核配额给②）；源级锚点＋两实例三份读数（noq9-out／noq9-selftest／mutroster HOLE×2）互证留痕。
5. 正控 want map 的逐种子点数（bare-goroutine 2/…/emoji 1）本程未逐名重数，凭 TestScanDetectsAllSeededViolations PASS＋self-test 36/36 两向。
6. Windows/NTFS 大小写不敏感 Stat 残洞：注释引 `internal/Build/x.go` 而小写件实存时 stat 会命中 ⇒ 静默；今日 sweep 交集 0、实害量不到，记已知形状（承第一实例）。
7. "首发 10 枚"的原始口径（r1 死腿会话内）无法复原；最硬现量＝§6 #2 的 11/8/9 三数。

## 终态

- status：**§0-§7 全满、零"待填"残留；AC#2 成立但带条件／AC#3 成立但带条件（缺陷具名：省略号-全拼路径形被问罪）／AC#5 成立；十枚复核＝212 本票 8 文件 9 hunk（5 真＋3 疑）＋2 行 258 随行件标不适用；推翻 8 条；判不动 7 条；证据件随本 commit 交件**。不翻勾、不 push、不改产码一字。三实例碰撞与重复派单归编排者裁撤对账。
