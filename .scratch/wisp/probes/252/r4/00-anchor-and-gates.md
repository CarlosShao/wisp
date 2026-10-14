# 252-r4 结算补尺件：把"两段必须各自可达"做成会响的尺——现读、两向读数、反形对照

- 腿：`252-r4`（票 252 结算补尺；写腿；面＝`internal/tools`；**仪器面：零产码、零新尺件**）
- 锚（起手现取）：HEAD `037bfb7f`（`2026-10-14 15:31 +08`）；`internal/tools/paths.go` sha256
  `28890e6370e3719fd8faad5919c968a14fa35c0c16f1f53d8e091816c96db76a`（每发突变后由仓外备份逐字节还原，终态复量同值）
- 起手/收尾名册：`git status --porcelain internal/tools internal/risk`＝**0 行**（起手与收尾各一发，均为 0）
- 原始日志留盘 `C:\tmp\wisp252r4\logs\`（6 枚，sha256 见 §7；件内摘录 5 枚在 `logs/`）

## §1 现读（行号全为本次现取，票面/旧件行号已漂）

**1.1 `internal/tools/paths.go` 的 `InAllowlist` 两段结构（现读）**
- `:201` 函数头；`:202` `f := foldPath(canonical)`。
- **词法腿** `:208-210`：`if !rootsContain(p.roots, f) { return false }`。
- **fold 腿** `:220-223`：`rf, ok := resolvedForm(canonical)` ＋ `if !ok || !rootsContain(p.roots, foldPath(rf)) { return false }`。
- 工作区腿 `:227-229` 问的是 `p.workspace`（另一个本），不属授权本；S1 尺的设计正是只数"对授权本取用的包含调用"（r2 §3.3）。
- 票面曾引的 `:189/:196/:208-209` 是 10-03 时点；r2 交付后今值如上（票 12:3x 追加已具名过这一类漂移）。

**1.2 生产调用图（非测试；显式根 `cmd internal tools scripts`，本次现 grep）**
- `internal/risk/rules_gateway.go:45` `if !ctx.canon.InAllowlist(canonical) {` —— `rulePathAllowlist`（`:32`），R2⇒L2 判级现场。
- `internal/tools/task.go:840` `case !d.Paths.InAllowlist(canon):` —— `pointerNotice`（`:831`），提示文案的判定。
- `internal/tools/bridge.go:1265` `strings.Contains(p, "无法规范化") || !b.paths.InAllowlist(p)` —— `scopeVerdictOf`（`:1260`）。
- `internal/agent/spill.go:231` `if !s.judge.InAllowlist(canon) {` —— 生产装配 judge=nil（v1 判语），今天仍不可达（如实列，不冒充调用点）。
- 接口声明两处：`internal/risk/assessor.go:143`／`internal/agent/spill.go:70`。

**1.3 判据尺现状（落点裁定的依据）**
三枚已在树，出处 commit `70a935ce`（2026-10-03 10:59；`git merge-base --is-ancestor 70a935ce HEAD` 本腿现验为真）：
- **S1**＝`internal/tools/paths_twocontainments_252_r2_test.go`＝`TestTicket252R2BothContainmentsReachable`（可携；按 `(receiver 类型, 方法名)` 认 `InAllowlist`，不硬编码 `rootsContain/foldPath/resolvedForm`；四断言：a 数到恰 2 枚、b 每枚能自拒、c 两键不同源、d 恰一 clean 一 tainted）。
- **W1**＝`..._r2_windows_test.go`＝`TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses`（词法腿值证；8.3 分歧输入，逐格带正控）。
- **W2**＝同首件＝`TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses`（fold 腿值证；根外链接分歧输入）。
- 既有牙（非本腿所造）：N7＝`paths_shortname_252_r1_test.go:223-238`；107b probe C＝`paths_ticket107b_probes_test.go:179-183`。

## §2 落点裁定（现读后裁）与对派单前提的推翻

- **判据①（拆任一段⇒尺红且红句具名）与②（两段在⇒绿）今天已由在树三枚尺满足**（读数＝§3）⇒ 本腿**不新造尺**（新造＝复制在树件，违查重），也**未动 `paths.go` 一字**（sha256 同值）。派单"造尺（能力形，两向）"的能力形＝调用图尺＋两枚分歧输入值证，已不在词表/恒等式那一类；本腿交的是"把它做成会响的尺"的**结算读数**：每发突变量到哪一句、具名说哪一段。
- **推翻 A921／收口那句（逐字）**：「`"绝不许把两段并成一段"` 硬禁**今天零仪器**（M-merge 实测 257 PASS/0 FAIL）」——
  (a) "257 PASS/0 FAIL"是 **10-03 09:49**（v1 一程、尺入库前）的旧读数；尺 `70a935ce` 落于 **10-03 10:59**；
  (b) 今天现跑：pristine 13 顶 PASS（含三枚判据尺）rc=0；删任一段⇒判据尺红且红句具名该段（§3）。
  ⇒ "零仪器"作废；票 10-03 12:3x 追加第 3 条（"那把尺已进树"）才是现态，收口节这一行是过期句。

## §3 两向读数（`GOTMPDIR` 只反斜杠形 `C:\tmp\wisp252r4gotmp`；每发 `-run 'TestTicket252|TestTicket107b'`）

| 发 | 突变（各自单摘） | rc | 判据尺红句具名哪一段 | 具名红句出处 |
|---|---|---|---|---|
| 正控 | 无（pristine，发前＋还原后各一发） | 0 | 全绿：13 顶 PASS／22 子 PASS／0 FAIL／0 SKIP | logs/10 |
| M-lex | 删词法腿（`:208-210` 三行） | 1 | **词法腿**（W1 红句逐字 "lexical containment missing…although the lexical leg refuses it"；S1 断言 a 数到 1≠2，found=`[paths.go:218 key=foldPath(rf)]`） | logs/20 |
| M-fold | 删 fold 腿（`:220-223` 四行） | 1 | **fold 腿**（W2 红句逐字 "resolvedForm leg missing"；S1 断言 a found=`[paths.go:208 key=f]`；既有 N7＋probe C 同红） | logs/30 |
| M-neuter-lex | 腿在位但不参与（`false &&`） | 1 | **词法腿**（只有 W1 值证抓；S1 结构上仍数到 2——如实登记） | logs/40 |
| M-neuter-fold | 腿在位但不参与（`false &&`） | 1 | **fold 腿**（W2＋N7＋probe C 抓；S1 静默） | logs/40 |

- M-lex 同格值读数（逐字）：`LEG1 rootsContain(roots,folded): false`／`resolvedForm(canonical): …q252r2-lexical-leg-only.md ok=true`／`LEG2 …: true`／`InAllowlist(canonical): true`／`FIRST FAILING LEG: LEG 1 lexical`——分歧输入形如设计。
- 两向都**不是静默绿**：删任一段时"两段并成一段"当场红且点名（判据①）。
- 附加两发**不是判据**，是射程具名：S1 只抓"删/改"，抓不到"在位却 dead"；后者由 W1/W2 值证抓（Windows 全有）；⚠ **POSIX 上词法侧 dead 形无值证**（W1 是 windows 件）＝本腿具名的那一格（§6 第 2 条）。

## §4 反形对照（236-v2 先例：反形不许当凭据）

- 旧形＝`measure252` 系的 `final == (leg1 && leg2)` 恒等式（`paths_shortname_252_probe_test.go:192`）与 R1 的
  `:143/:147` 同形断言；它们一律只喂 `Canonicalize` 之后的**同形串**。
- **同一两棵突变树上**（与 §3 的红出自同一发日志）：P1 三枚＋R1 四枚（除 N7 外）**全绿**（逐字 PASS 名册 logs/50）。
- ⇒ 旧形对"删段"照绿＝**不可作凭据**；同一两棵树上红凭据只由判据尺（S1/W1/W2）供。这就是派单第 3 条的读数。

## §5 门禁（终态现跑，全部在盘）

| 尺 | 命令 | 读数 |
|---|---|---|
| build | `GOFLAGS= go build ./...` | rc=0（0 字节输出） |
| gofumpt | `$(go env GOPATH)/bin/gofumpt -l internal/tools` | 0 行（本腿未动任何 Go 件，仍跑） |
| `-run` | 同上 `-run 'TestTicket252|TestTicket107b'`（终态发） | rc=0（13 顶 PASS／22 子 PASS） |
| d22scan | `sh scripts/d22scan.sh` | rc=0，`clean - no D22 ban violations` |
| 长度门 | `sh scripts/check-path-length-budget.sh` | rc=0，`VERDICT GREEN`：tracked=12779／over=60=roster／not in roster=0（与派单"现 60=60"一致） |

## §6 具名没做／量不到

1. ⛔ 未跑整包 `internal/tools` 全名册、未跑 `internal/risk`、⛔ 未跑 `cmd/wisp`（314 死区按令不碰）、未跑 CI、未 push。
2. 本机是 Windows：**POSIX 上"词法腿在位但 dead"没有值证**（S1 抓不到、W1 不编）；本腿只量到 Windows 全牙（§3）。这一格要真牙需一枚 POSIX 侧词法值证（镜像链接形：拼法在外、解析在内）；本机现量 `filepath.EvalSymlinks` 对 junction 不解析（`EvalSymlinks(outside\lnk\f.txt)` → err="The system cannot find the path specified."）⇒ 本机不可造，如实不做。
3. 未把 S1 红句改成"具名缺的那一段"的更强文案（那会改写 v1 已逐字对过账的红句）——点名，留给后续腿裁。
4. 未翻任何 AC 框、未动票名／`Status:`（本票已 `-done`，本腿只追加票尾 bullet）。
5. 临时件只建不删：突变件与备份在 `C:\tmp\wisp252r4\`，全数在盘。

## §7 件与哈希

- 原始日志（留盘，不随提交）：`C:\tmp\wisp252r4\logs\`＝
  `00-pristine-252verbatim.log`(52070B, sha256 `f956d67107ef60333f8899184d3ad5011657ab82258b0f385b860dfdbaa82ce9`)／
  `10-mlex-all.log`(60042B, `ede57a85f80aa039f8f58160f542214d2343e705165531de8b1ffd6b976f9150`)／
  `20-mfold-all.log`(60412B, `b7a1be58b5e1eb6a08657d09a5cb62ce5296bd12628cf57cdfd4fc8abaf81b94`)／
  `30-mneuterlex-all.log`(59447B, `5fe6611c7f69a1510f527813147b2759a3afb1610a34ca8f048bbdde93a9d442`)／
  `40-mneuterfold-all.log`(59867B, `c18263170c657730bfd071f28e1316419365a97293f5fc23585052f3c9fcba5f`)／
  `50-pristine-after-restore-all.log`(57925B, `3d5b798342fb6847a4fc5bc918383e67f7598b4017ace3b832c58b11ff8cb96a`)。
- 备份 `C:\tmp\wisp252r4\backup\paths.go.pristine` sha256 `28890e63…`；变形件 `mlex-paths.go`(`bca225ce…`)／`mfold-paths.go`(`a4c997e2…`)。
- 本件族（随提交）：`00-anchor-and-gates.md`＋`logs/10-poscontrol.txt`／`20-mlex-reds.txt`／`30-mfold-reds.txt`／`40-neuter-two-way.txt`／`50-antiform-green.txt`。
