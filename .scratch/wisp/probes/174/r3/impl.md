# 174-r3 证据件 — AC#2d 第二半（泄漏断言）＝ Canonicalize 返错那支的注意句不许转述毒值

## §0 起手锚与名册

- 时间锚：起手 `2026-10-03T17:28:50+08:00`（`date -Iseconds`）。
- HEAD 锚：`8199a644f9e203faf4266cb10c4400642e8617b7`（`git log -1 --format=%H`）。
- 起手 `git status --porcelain internal/tools internal/risk`＝**空**（一行输出，无条目）。
- 起手基线（改前）：`go test -count=1 ./internal/tools/`＝`ok github.com/CarlosShao/wisp/internal/tools 13.756s`，存档 `logs/baseline-gotest-before.txt`。
- md5 锚（改前）：`internal/tools/task.go`＝`4138177e29ffff427776a73eb0d61a9b`；`internal/tools/task_output_pointer_notice_test.go`＝`54cbbb4db55a2662d7c1a06c0a351347`（`logs/md5-anchor-before.txt`）。
- **行号已漂（推翻票面行号）**：票面 AC#2d 写 `task.go:313-314`，现锚该支在 **`internal/tools/task.go:836-837`**（`case err != nil:` → `notes = append(notes, "注意：…C26 连规范化都没通过（"+err.Error()+"）")`）。五支全貌在 `:829-851`（`pointerNotice`）。
- **写面锚（改后）**：`internal/tools/task_output_pointer_notice_test.go` 两枚新判据＝`:420`（泄漏钉）＋`:497`（合法路径正控）；毒值常量 `:392`/`:397`；junction 造法 `:404`。md5 改后＝`7179c7597ffb6c669045ebcaa9960f99`（`logs/md5-anchor-after.txt`）。
- **唯一改动文件**＝`internal/tools/task_output_pointer_notice_test.go`（diff 形状：只追加，无一行删除——`git diff -U0` 的 `-` 侧仅 `---` 头，零净删）。产码 `task.go` md5 改后仍＝`4138177e29ffff427776a73eb0d61a9b`（零字节）。

## §1 钉的形状（断言逐条＋毒值设计）

**背景（r2 半格止步处的接续）**：174-r2 已钉"说了没有"（`task_output_canonicalize_fail_174_test.go` 两枚：接线臂必须说 `连规范化都没通过` 且不说 `未接线`，`Paths:nil` 臂反之——MUT 摘 `task.go` 该支整句 ⇒ 恰两枚红）。但它那句"欠的第二半"止步于"钉它必须动产码（把 `err.Error()` 原文剥掉）"⇒ 停手上报。本程派单给了另一条路：**把毒值埋进路径字面**——"路径原文/根列表"两半就在不动产码的前提下变绿着钉；只有"C26 内部状态词"那半今天真的在漏，按派单 §6"要动产码＝另一格"上报（见 §6）。

**毒值设计**（三枚，全埋在 junction 造出的真实文件系统形状里）：

1. `leakProbeSegment174r3 = "ALLOWLIST-LEAK-PROBE-174r3"`（`test:392`）——junction 的**名字**。任何对"宿主填的那条路径"的转述（原样、Clean、Abs、短名、任何大小写折叠）都会把它带出来。
2. `leakProbeResolved174r3 = "RESOLVED-LEAK-PROBE-174r3"`（`test:397`）——junction **target** 的名字。堵"修复顺手把 junction 解析成真树告诉模型"那路。
3. `rootA`/`rootB`（两枚 `tempCanonical(t)`）——授权根列表，两枚判据的桥都配 `NewPathCanonicalizer([]string{rootA, rootB}, nil)`，使"根列表"禁令是一张列表而不是一枚目录。

**为什么用真 junction（`mklink /J`，`test:404`，照抄 `bridge_junction_windows_test.go` 的 `mkRealJunction`，按票 238 先例不抽公共）**：`pointerNotice` 的 `case err != nil` 支（`task.go:836`）在**今天**可达的输入只有两类——空白路径（`paths.go:120-122` 返 `tools: empty path`，已被 r2 的 fixture 占用）与非豁免 reparse 穿越（`risk.Resolve` 返 `ErrReparseDenied`）。空白路径里埋毒值没意义（它根本不是路径形状），**只有真 junction 能让"一条含毒值的路径字面"走进这支**。junction 是 NTFS/Windows-only：本判据平台无 tag，`mkJunctionForLeakProbe174r3` 失败时 `t.Skipf`（注明"前置条件缺失：造不出真 NTFS junction"，不许字符串假装 junction——抄 r2 先例的措辞方向）。`risk` 侧 junction 测试先例＝`internal/risk/pathresolver_junction_windows_test.go:95`（`errors.Is(err, ErrReparseDenied)`）。

**判据一 `TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3`（`:420`，泄漏钉）**，fixture：`rootA` 下造 junction `rootA/ALLOWLIST-LEAK-PROBE-174r3` → target `rootA外的tempDir/RESOLVED-LEAK-PROBE-174r3`；roster 按链接拼写原样登记（`ArtifactPath: filed`＝junction 字面——canonicalize 它恰好就是必须返错的那次调用，没有别的形态可登记）。断言逐条：

- 正向（r2 已钉那半，在本臂上重钉使判据独立成立）：`!IsError && Truncated`；含 `注意：`；含 `读不到`；含 `连规范化都没通过`。
- 指针仍给（PLAN.md:2564）：含 `全文见 `+filed 逐字（junction 含毒值的名字**合法地**印在指针句里——那是冻结模板的所有物）。
- **泄漏扫（牙）**：把注意句**切出来**（第一个 `注意：` 起、到 `，全文见 ` 止——恰是"拼进去的那段"的字面射程）逐项断 `strings.Contains(notice, banned)` 不成立，banned∈{filed 原文、`ALLOWLIST-LEAK-PROBE-174r3`、`RESOLVED-LEAK-PROBE-174r3`、rootA、rootB、`读得回来`/`随时可读`/`可以读回`、`不在你被授权的目录范围内`、`未接线`、`不可找回`}。
- **为什么扫注意句而不是全回执**：指针句合法含路径（`:299`/`:349` 两枚模板冻结钉系着，派单写面禁止动它）；全回执扫＝自相矛盾的尺，必恒红。

**判据二 `TestLeakProbeTokenOnALegalPathStaysSilent174r3`（`:497`，正控/尺活证）**：同一枚毒值 token 放在**根内合法普通文件**名上 ⇒ 走健康臂：零 `注意：`、指针照印（印的正是同一 token）、同桥 `fs.read`＝`isError=false level=L0` 全文。⇒ 尺按**支**区分，不按 token 是否在场区分。

## §2 种 X 必响读数

 Junction 臂实跑（未变异产码，`-run 'TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3|TestLeakProbeTokenOnALegalPathStaysSilent174r3' -v`，两枚 PASS，0.07s/0.02s）：

- 回执 announce 逐字（判据一，`logs` 无档但读数同 MUT-A 红句前缀；逐字存于 `mutA-leak-criterion-run.txt` 对照样张）：
  `[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，C26 连规范化都没通过（risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions）；注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状），全文见 C:\\…\\001\\ALLOWLIST-LEAK-PROBE-174r3…]`
- **说到了**：`C26 连规范化都没通过` 在场（正向半响）；junction 确实把路径带进 `case err != nil` 支。
- **没漏毒值**：注意句（切出射程）＝`注意：这条路径现在读不到，C26 连规范化都没通过（risk: path traverses …）；注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状）`——`ALLOWLIST-LEAK-PROBE-174r3`/`RESOLVED-LEAK-PROBE-174r3`/rootA/rootB **零次出现**（两枚毒值 token 与两枚根都在回执里存在——毒值在指针句里、根在桥的配置里——但都**不在注意句里**）。泄漏扫绿。
- 顺带量出：junction 臂同时触发**第二支**（`它存在但不是一般文件`——junction 的 Lstat 模式不是 regular file，`os.Stat` 会解析 target 但 `Mode().IsRegular()` 为 false），故 `注意：` 计数＝2。这是形状事实，不是缺陷；两支都是"读不到"语义。
- 判据二（正控）：announce＝`[……全文见 C:\\…\\ALLOWLIST-LEAK-PROBE-174r3.txt…]`，零 `注意：`，`fs.read` L0 全文。**同 token 两支分立**＝尺活。

## §3 变异自证（红句逐字＋还原）

全用 `-overlay`（跟踪树零改动；每发后 `git status --porcelain internal/tools internal/risk`＝空）：

- **MUT-A（摘掉那支的整句文案；`mut/task-MUTA.go`，sed 把 `notes = append(notes, "注意：这条路径现在读不到，C26 连规范化都没通过（"+err.Error()+"）")` → `_ = err // MUT-A`）**：`-run 'TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3|TestLeakProbeTokenOnALegalPathStaysSilent174r3'` ⇒ rc=1，红句逐字（`logs/mutA-leak-criterion-run.txt`）：
  `task_output_pointer_notice_test.go:451: the reply must name normalization as the step that failed, got "[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状），全文见 C:\\…\\001\\ALLOWLIST-LEAK-PROBE-174r3…]"`
  （正向句先红——同时证明摘掉该支后注意句无本钉的其他违规。）
- **MUT-B（`err.Error()` 换成自造词 `MUTB-LEAK-PROBE-174r3-leak of filed path`，不在禁表）**：rc=0 绿 ⇒ 刻意留的口径证词：扫的是**禁表里的具体毒值**（派单要求的"负向断言＋种的具体毒值"），不是任意字符串。
- **MUT-B2（把禁表毒值本体 `ALLOWLIST-LEAK-PROBE-174r3` 种进注意句，`mut/task-MUTB2.go`）**：rc=1，红句逐字（`logs/mutB2-poison-token-run.txt`）：
  `task_output_pointer_notice_test.go:486: the canonicalize-error notice must not relay "ALLOWLIST-LEAK-PROBE-174r3": notice "注意：这条路径现在读不到，C26 连规范化都没通过（risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions 处理时查过 ALLOWLIST-LEAK-PROBE-174r3 这一项）；注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状）"`
  ⇒ **泄漏扫有牙**：毒值进注意句必红。
- **还原**：`-overlay` 本身零改动；`task.go` md5 终态复归 `4138177e29ffff427776a73eb0d61a9b`（与起手同值），`git status --porcelain internal/tools internal/risk`＝` M internal/tools/task_output_pointer_notice_test.go` 一行（本程唯一写面）。
- 变异台件只建不删：`mut/task-MUTA.go`、`mut/task-MUTB.go`、`mut/task-MUTB2.go`、`mut/overlay-MUTA.json`、`mut/overlay-MUTB.json`、`mut/overlay-MUTB2.json`、`mut/test-flipped.go`、`mut/overlay-flipped.json`。

## §4 恒真自查（硬要求）

把判据换成反形（`if strings.Contains(notice, banned)` → `if !strings.Contains(notice, banned)`，即"注意句**必须**含毒值"，`mut/test-flipped.go`）跑在**未变异产码**上 ⇒ rc=1 红（`logs/flipped-criterion-run.txt`），红句逐字（节选）：
`task_output_pointer_notice_test.go:486: the canonicalize-error notice must not relay "ALLOWLIST-LEAK-PROBE-174r3": notice "注意：这条路径现在读不到，C26 连规范化都没通过（risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions）…"`
`task_output_pointer_notice_test.go:486: the canonicalize-error notice must not relay "RESOLVED-LEAK-PROBE-174r3": …`
⇒ 反形会红、正形会绿＝**尺对泄漏敏感，不是恒真尺**。只断言"有注意二字"的恒真尺不存在于此：MUT-A 证明正向断言有牙、MUT-B2 证明泄漏断言有牙、反形探针证明它不是怎么写都绿。

## §5 门禁四数（终态，现跑）

1. `go vet ./internal/tools/` ⇒ **rc=0**（`VET-OK`）。
2. `sh scripts/d22scan.sh` ⇒ **rc=0**，`clean - no D22 ban violations`；`ban #8 internal/ examined 498`（`logs/d22scan.txt`）。
3. `gofumpt -l internal/tools/` ⇒ **空**（v0.12.0 (go1.27.1)；`logs/gofumpt-l.txt` 零行）。
4. `go test -count=1 ./internal/tools/` ⇒ **rc=0**，`ok github.com/CarlosShao/wisp/internal/tools 13.966s`（`logs/gotest-final.txt`）；名册（`-v`，`logs/gotest-final-verbose.txt`）：顶层 `--- PASS` **191**／`--- FAIL` **0**／`--- SKIP` **0**（`=== RUN` 含子测试 264）。
- **名册 comm**（`git grep -h '^func Test' HEAD -- internal/tools/` vs 工作树，双向）：added＝`TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3`、`TestLeakProbeTokenOnALegalPathStaysSilent174r3`（恰好两枚，即本程新钉）；removed＝**空**。
- ⛔ 未跑全仓（派单：`cmd/wisp` 有 `261-r2` 写腿在飞、`internal/llm` 刚落门）。`go vet`/`d22scan` 只对 `./internal/tools/`；d22scan 脚本自身扫 internal/+cmd/ 皆 clean。
- 终态 `git status --porcelain internal/tools internal/risk`＝` M internal/tools/task_output_pointer_notice_test.go`（一行，唯一写面）。

## §6 判不动（含"产码文案要不要剥路径"那格——只上报不判）

1. **"产码文案要不要剥路径"＝判不动，上报**：`task.go:837` 今天把 `err.Error()` 原文拼进模型可见文本，junction 臂实测拼进的是 **`risk.ErrReparseDenied` 的完整原文**（`risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`）——C26 内部状态词**真的在漏**（形状无界：`Canonicalize` 将来任何错误文本都原样进上下文）。但**本格不许修**：那句受 r2 判据正向句（`连规范化都没通过` 字样）与 `:299`/`:349` 模板冻结钉系着，改它＝契约轴（须编排者＋owner 一句裁"算不算动冻结文字"，同 `Q-63` 边界）。**本钉的射程刻意不含它**：若今天钉"注意句不含 `risk:` 前缀/C26 词汇"，名册带着红 commit＝把"今天在漏"写成"已修"的假象；按 174-r2 先例（A375：停手上报、不写会红的断言），那一半写成已验证的钉放 probes（MUT-B2 证牙）、判不动写在这里。**具名风险**：`risk.ErrReparseDenied` 文本含 `(junction/symlink)` 与 `reparse_point_exceptions`——后者是配置键名，泄给模型等于把授权面词汇表端进上下文。
2. **junction 臂的平台射程**：真 junction 只在 NTFS 可造。本判据平台无 tag，非 NTFS/非 Windows 上 `mkJunctionForLeakProbe174r3` 会 `t.Skipf` ⇒ 那些宿主上"泄漏钉"**不执行**（不是绿着骗人——SKIP 计数在名册里可见；本轮读数 SKIP=0＝本机 NTFS 上两枚都真跑了）。空白路径臂（r2 的 fixture）不受此限，全平台跑。
3. **r2 判据与本钉的重叠**：判据一的正向三句（Truncated/注意/连规范化）与 r2 `TestCanonicalizeFailureFailsClosedInReply174` 在空白路径臂上同型，但 fixture 不同（junction vs 空白）且两枚 r2 判据仍独立由 MUT-A 红过（r2 交件时证）。本钉的正向半是让泄漏钉**独立成立**（摘文案的 MUT 红它），不是复写 r2。

## §7 推翻清单（票面、174-v1 表、本派单每一句都是待验断言）

1. **行号漂移**：票面 AC#2d 写 `task.go:313-314`——现锚该支在 `task.go:836-837`（五支全貌 `:829-851`）。派单 §3.1 说"看 task.go:313-314 的分支条件自己造"——条件本身没漂（`Canonicalize 返错`支），漂的只是行号。174-r2 证据件里的 `:325-326` 同漂。
2. **r2 交件里"要钉它必须动产码"的判断——对一半**：对"C26 内部状态词"那半（见 §6.1）成立；对"路径原文/根列表"那半**可推翻**：junction fixture 让毒值进了路径字面，那两半在**不动画着的那两行产码**的前提下就能绿着钉死。本程即据此落钉。
3. **派单 §0"`grep -c '规范化' …＝0`"**：那把尺扫的是 `task_output_pointer_notice_test.go` 单文件。现量：本程落地后它＝**1**（判据一的断言字符串 `连规范化都没通过` 含 `规范化`）；对照 r2 文件 `task_output_canonicalize_fail_174_test.go` 自 r2 落地起＝**6**。⇒ 票面那把 grep 尺从 r2 落地那刻起就没再是"零覆盖"的量度（票面格文字未改、照抄的旧读数），本程落地后它进一步失效——"零覆盖"的证词请用名册函数名或 MUT-V5/MUT-A 那类变异读数，别再用那把 grep。
4. **派单 §0"MUT-V5 把那支整体摘掉 → rc=0、124 枚全绿"**：方向仍真（r2 前的史实），但**数字过期**：r2 落地后两枚 r2 判据红；本程落地后判据一也红（MUT-A 现跑 rc=1）。MUT-V5 那发今天再跑不会再绿。
5. **派单 §0"它是唯一把 err.Error() 原文拼进模型可见文本的一支"**：**限 task.output 的 pointerNotice 五支内成立**；放大到全文件不真——`task.go:509`（参数解析失败：+err.Error()）与 `:519-521`（查不到任务）也拼上游文本（前者拼的是本地 json 解析错误、后者是 roster 自己的查无文案，无泄漏面）；**全仓面**更不真（fs.go/bridge.go 各有 err.Error() 入文点）。票面句子的射程是"五支之内"，照此理解。
6. **`mkRealJunction` 会 Fatalf 不可抄**：r2 时代已如此；本程抄的是造法、降级成 Skipf（平台无 tag 的判据文件不能 Fatalf 一个非 NTFS 宿主）。这不是推翻 r2 的用法（那文件平台 tag 只在 Windows 跑），是写面位置不同导致的前提不同。
7. **派单 §3.2 的正控措辞"把毒值换成合法路径时回执走的是另一支"**：照做且响（判据二）；但现量补充：junction 臂的回执是**两支并响**（`连规范化都没通过`＋`不是一般文件`）——`os.Stat` 对 junction 的 `IsRegular()` 为 false。这不影响尺（泄漏扫射程是注意句整体），但"走了哪一支"的措辞在本臂上是"走了两支"。
8. **票面"只建不删"对 probes 目录**：照做；`mkdir -p` 建出 `probes/174/r3/{logs,mut}`，目录本身不算"临时件删除"问题，全部留存。

## §8 量不到的格子

1. **`ErrReparseDenied` 那半泄漏的"修后形态"**：本钉刻意不钉"注意句不含 C26 词汇"（§6.1）——那一半"修后应响"的读数**没量**（会红），只有"今天在漏"的形状读数（junction 臂回执逐字含 `risk: …` 原文）与 MUT-B2 的牙证。修后判据长什么样、要不要连 `（junction/symlink）` 这种解释性词一起剥，归下一程。
2. **非 NTFS 宿主上的 junction 臂**：本机只量了 NTFS（Windows）。POSIX 宿主上这两枚新判据的实况（SKIP）没量——本仓生产平台是 Windows（D7），junction 判据的 SKIP 行为按 stdlib 语义写、未在真 POSIX 宿主上跑过。
3. **`fs.read` 对照臂在 junction 上的读数**：判据一没像 r1 那样对同一发跑 `readThroughBridge` 对照（junction 路径过桥会被 `risk.Resolve` 拒、而 `fs.go open` 只 Canonicalize 不查根——这条对照在 junction 臂上会量出"fs.read 也进 Canonicalize 返错"＝第三支形状，与本格的回执文案判据不同轴）。没量；登记为形状事实缺口，不影响本钉任何断言。
