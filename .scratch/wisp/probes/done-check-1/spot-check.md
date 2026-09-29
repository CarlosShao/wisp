# done-check-1 — 非实现者抽验：`done-fix-1` 标了「戊类·待非实现者抽验」的 9 枚格

- 起手锚点：`49c34a7e`（`git rev-parse --short HEAD`，分支 `dev`）
- 本腿身份：**非实现者抽验腿**，只裁不产码；⛔ 不跑 `go test`/`go build`/`go vet`/`gofumpt`/`d22scan`、不执行任何二进制（机器上有 `223-v2` 验收腿在编译与跑门，争用会洗掉它的读数）
- 允许用的尺：`Read`／`grep -rn`／`sed`／`wc`／`find`／`ls -l`／只读 `git log`／`git show`／`git cat-file blob <commit>:<path>`
- 结论只允许四种：①**可翻勾** ②**不能翻，缺的是<一句话>** ③**判不了：凭据在被禁读的目录里** ④**这格不是判据、是纪律句／待办行（形状写歪）**
- 双锚通则（`A442`）：引用别的表的行号时同时带上那一行开头的字串。

## 待验 9 格（位置＝`票文件:行号`，锚定字串 `` `done-fix-1` 追加（戊类 ``）

| # | 位置 | 判据 | 结论（未验＝〔待填〕） |
|---|---|---|---|
| 1 | `104-…-done.md:65` | AC#5 与票 89 第 4 条的分工写清 | **可翻勾**（已翻，`1d5678fb` 之后本腿 commit） |
| 2 | `110-…-done.md:44` | AC#4 `R-93-4` 一并收（步级证据或明写为何不该纳入） | **可翻勾**（走"或"的第二支；已翻） |
| 3 | `113-…-done.md:76` | AC#6 把 `R-108-2` 边界写进包文档（只增不减／不动判定分支／容器读数不受影响） | **可翻勾**（已翻原框，依编排者 `:85` 裁定） |
| 4 | `115-…-done.md:49` | AC#2 裁方向 A/B＋写理由代价＋明写是否推翻 104/105 | **可翻勾**（已翻；与票 230 的口径撞车已登记） |
| 5 | `115-…-done.md:53` | AC#3 修完四枚红＋新增行为用例（正/反两半） | **可翻勾**（已翻；"转绿"证据＝`ci-step-readings-2026-09-22.md` 非实现者读数） |
| 6 | `115-…-done.md:66` | AC#7 结案前必须有一枚远程 run id | **可翻勾**（已翻；run `35616790753` 两处独立盘上落点） |
| 7 | `89-…-done.md:262` | 〔待填〕 | 〔待填〕 |
| 8 | `89-…-done.md:310` | 〔待填〕 | 〔待填〕 |
| 9 | `89-…-done.md:345` | 〔待填〕 | 〔待填〕 |

---

## 第 1 格 — `104-sealfile-silently-drops-inherited-grants-done.md:65`（AC#5）

判据原文（票面 `:63-64`）：
> `- [ ] **AC#5** 与票 89 第 4 条的**分工写清**：票面/commit 正文里明说"本票只补继承那一半"，并核对 `89-...-done.md` 票面第 4 条的措辞**有没有被本票读成"已全部覆盖"**——若有，登记更正，**不改它的原文**。`

本腿现跑的尺（四把，全部本腿自己跑，零转述）：
1. `grep -n 'R-104-2' .scratch/wisp/issues/89-0600-is-decorative-on-windows-acl-for-private-data-done.md`
2. `git show --numstat --format=%h 9141d4cf -- <89 票>`
3. `git cat-file blob 9141d4cf^:<89 票> | sed -n '157p'` 与当前 `<89 票>` `sed -n '157p'` 对拉 `diff` + `md5sum`
4. `git show -s --format=%B 4d43447 | grep -n "继承那一半\|AC#5\|89"`；`grep -n "Inherited\b" internal/winsec/winsec_windows.go`

读数：
1. `:159` 命中，逐字「**⚠ 编排者更正（2026-09-21 21:0x，来源=`acceptor-ticket104` 的 `R-104-2`；票 104 的 AC#5 就卡在这一句）**」，更正段延伸到 `:170`。
2. `12 0`＝该 commit 给 89 票面加 12 行、**删 0 行**（纯追加）。
3. `diff` 空输出＋`SAME-157`，两行 `md5sum` 同为 `f4b608d555ab2de30c2ba36122627fc2`＝被更正的原句「为什么只报显式、不报继承来的」原文一字未动。
4. commit `4d43447` 正文第 4 行逐字「…本票**只补继承那一半**。」；本票 `:152` 亦写着「**AC#5 分工**：本 commit 正文与本票只主张"补继承那一半"」。代码侧 `internal/winsec/winsec_windows.go:66/:71`（`Inherited []string`）、`:138`/`:140`（两桶分支）、`:147`（`cleared_inherited` 字段）、`:317`（`noticeNarrowed(... Principals: explicit, Inherited: inherited)`）⇒ 更正段那句"通知面扩成两桶"是当前树的事实。

结论：**可翻勾**。判据两半（分工写清／更正已登记且原文未动）都由本腿在盘上现读到；`done-fix-1` 记的那句「实质已闭，接线的动作从没做」本腿认——缺的正是"由非实现者认一次"，本腿这次就是那一认。已按派单形状在原格下追加一行 `done-check-1 勾＝…`，并把行首 `- [ ]` 翻成 `- [x]`。
补充（不改变结论、只留痕）：表 `docs/evidence/s1/104-adversarial-acceptance.md:250` 对这格判的是「**不通过**（…→ R-104-2）」，那是登记**之前**的状态；本腿不抹那行。另 `R-104-2` 在 `docs/reports/pending-and-issues.md` 里 grep 为 0 命中（表 `:231` 当年写的是「写进 104 结案语**或** `pending-and-issues.md`」，落点选了 89 票面＋104 结案语 `:5`，两个或支里落了一个，本腿判这属"登记位置二选一已兑现"、不构成缺读数）。

## 第 2 格 — `110-no-ci-step-runs-internal-winsec-done.md:44`（AC#4）

判据原文（票面 `:43`）：
> `- [ ] **AC#4** `R-93-4` 一并收：给出"windows 腿确实把 `TestSyncRegistryProbeLive` 纳入分母"的**步级证据**，或明写它为什么不该被纳入。`

本腿现跑的尺：
1. `grep -n "AC#4 \`R-93-4\` 的答复" .scratch/wisp/issues/110-no-ci-step-runs-internal-winsec-done.md`
2. `grep -n "TestSyncRegistryProbeLive" scripts/portable-tests.sh`
3. `sed -n '17p;33p' docs/evidence/s1/110-adversarial-acceptance.md`（表侧裁语）
4. `grep -n "AC#7" .scratch/wisp/issues/111-ci-tests-20-of-33-packages.md`（那笔附条件账今天在哪）

读数：
1. `:90` 逐字「（agent-ticket110）**AC#4 `R-93-4` 的答复 = 它不该进 RUN 分母，但它该拿的步级证据我给了**」＋"纳入只有两种读数：`runtests.sh` 下 SKIP=1 的常红，或裸 `go test` 下的假 `ok`"。
2. `:326` 那行 ledger（`windows|fixture|…As of ticket 110 AC#4 the windows leg carries ./internal/risk/ in its scope…`）今天在树里（同 `:7` 那行注释也还在）。
3. 表 `docs/evidence/s1/110-adversarial-acceptance.md:17`「| AC#4 | …**通过（附条件：步级证据至今不存在）**…**"不该进 RUN 分母"的论证我认可**…⚠ **但它许诺的步级证据至今 0 次**…」，同表 `:33`「**通过（附条件）**…AC#4 成立但**附条件**（其步级证据至今不存在，条件 = R-110-2/R-110-4 的补救）」。
4. `:75` `- [ ] **AC#7** `R-110-4`：票 110 承诺过…` **未勾**，票 111 无 `-done` ⇒ 那发 step7 `go test -list` 读数仍欠着，本勾没有把它抹掉（同根账另见 `110-…md:113`「**R-110-4** AC#4 许诺的…至今**无一次步级读数**」）。

结论：**可翻勾**（走判据"或"的第二支）。凭据链＝理由文本写在票面 `:90`＋脚本 `:326`，且这一支**早已由非实现者 `acceptor-110-97` 认过**（表 `:17` "论证我认可"）——`done-fix-1` 交给本腿的那一问（"那行理由文本算不算 `R-93-4` 要的步级证据"）答案是：**不算步级证据，但判据第二支本来就不要求读数**；第一支缺的那发读数另记在 `R-110-4`／票 111 AC#7，仍开放。
⚠ **推翻 `done-fix-1` 分类之处（第 1 条）**：它把这格记为"戊·形状齐、裁决缺，要翻请连那发 `go test -list` 读数一起看"，本腿认为那发读数是另一张账的对象，本格按"或"句的第二支已闭合。**分歧已在票面追加行里写明，交编排者裁。**
三行定性：现象＝一枚 windows 专属用例被记在 ledger 而不进 RUN 分母；本机被入侵证据＝无（本腿只读文件与只读 git）；最坏后果形状＝若那条论证错了，真实注册表行为在 CI 永不被执行（覆盖缺口），不是对任何已交付码的放宽或授权。

## 第 3 格 — `113-posix-platformverifypplacement-has-no-link-leg-done.md:76`（AC#6）

判据原文（票面 `:61`＋`:66`）：
> `- [ ] **AC#6（编排者 20:2x 追加，只改文档不改语义）** 把 R-108-2 的边界写进 internal/winsec/doc.go 一段话：…`
> `判据：该文件 git diff 只增不减；不新增/不修改任何判定分支；容器内 AC#5 那三门读数不受影响。`

本腿现跑的尺与读数：
1. `git show --numstat --format="%h %s" 1499efe` → `1499efe8 docs(113,AC#6): R-108-2 的边界写进 winsec.go 包文档——守卫管拼写/祖先链/树归属，不管"这棵树归谁"`；两行 `18 0 internal/winsec/winsec.go`、`62 0 .scratch/wisp/issues/113-…-has-no-link-leg.md`（**删除列两枚都是 0**，且 commit 真在历史里）。
2. `git show 1499efe -- internal/winsec/winsec.go | grep '^+' | grep -v '^+++' | wc -l` ＝ **18**；同集 `grep -vc '^+[[:space:]]*//'` ＝ **0** ⇒ 18 行全是注释，零行可执行码被加或被删。
3. `Read internal/winsec/winsec.go:40-54`：`:40-42`（"None of them asks whose tree it is"）＝"管拼写/祖先链/树归属、不管这棵树归谁"；`:42-45`（"Given a foreign absolute path … SealFile succeeds - on Windows that strips an explicit S-1-1-0 grant"）；`:46-49`（"a ruled boundary, not a gap (ticket 113 AC#6, answering R-108-2) … belongs to the caller's data-root discipline (tickets 76/95)"）。`grep -n "S-1-1-0" internal/winsec/winsec.go`＝只有 `:44` 一处，且落在注释里。
4. 表侧对口：`docs/evidence/s1/113-adversarial-acceptance.md:18`「| AC#6 文档边界（`R-108-2` 落 `winsec.go` 头部） | 票面自述"还没做" | …**验收现场工作树里确有一段未提交的 +18 行注释**…⇒ 未 commit = 未交件」＋ `:110`「**AC#6 必须落成真 commit**…判据是 `git diff` 只增不减、不动判定分支⇒ 在此之前**票面 AC#6 那格不许勾**」。

第三条判据（容器读数不受影响）怎么办：本腿禁跑编译/门，**不复跑**；改用构造证——尺 2 量出"零行可执行码变动"、注释在 `package winsec`（`:55`）之上 ⇒ 不可能改变任何用例结果。票面 `:82-83` 那组 `16/16 与 32/32, rc=0` 是实现方自述，本腿**不当凭据**。

结论：**可翻勾**（原框 `- [ ]`→`- [x]`，删除列 1）。翻的是验收方自己给的那两条条件的兑现，不是绕过它们；编排者本票面 `:85`「**我的裁定：以"翻转原框"为准（框翻转不是抹内容，删除列 1 可接受）**」是这条动作的授权文本，`done-fix-1` 因怕越权未翻，本腿引的正是它自己点出的那条例文。
⚠ **行号漂（`A442` 那一类）**：`done-fix-1` 追加行里写"追加勾 `:66`""编排者 `:79` 裁定"，本腿现量分别是 `:71` 与 `:85`——内容逐字对得上、指针漂了 5/6 行。

## 第 4 格 — `115-…-done.md:49`（AC#2 裁方向 A/B）

判据原文（票面 `:47-48`）：
> `- [ ] **AC#2** 裁方向 A / B（上面那一刀），写清理由与代价，并**明写它会不会推翻票 104/105 已交的语义**；若你判"两者都不对、真正该改的是用例的比对面"，就照实说，**不许为省事把判定放宽**。`

本腿现跑的尺与读数：
1. 裁定与理由在盘：票面 `:103`「**AC#2 裁定 = 方向 B（按树不按拼写）**…」＋理由段 `:160-178`（三条"为什么不是 A"按硬度排；`:173` 把 B 对到票 106/112/risk 票 72 的既有先例；`:176` 落地形状）。
2. "不推翻 104"这一半（生产侧是纯加法）：`git show --numstat 391878a`＝`275 0 notice_attribution_115_windows_test.go`、`61 0 winsec_windows.go`（**删除列两枚 0**）；`Read winsec_windows.go:95-135` 现读 `noticeNamesTree` ＝ `ResolvePath(spelling)` 后 `sameTree(n.Path, resolved.String())`；`grep -c "filepath\.\(Clean\|Abs\)" winsec_windows.go`＝**0**（全包另有 `resolve.go` 1、`winsec.go` 1 两处注释提及与 `seam_probe_root_125_other_test.go` 4 处 POSIX 测试夹具）⇒ D22 ban #2 成立；桶仍在（`:71/:138/:140/:147/:317`）。
3. "不推翻 105"这一半（本腿自算，不抄票面）：`grep -rn "RewrittenRoots(" --include=*.go . | grep -v _test.go`＝生产侧只有 `internal/tools/paths.go:194`（定义）＋`internal/tools/bridge.go:1053`（读取）；票面 `:170` 写的是 `bridge.go:864` ⇒ 同一函数、行号漂 189 行，内容对、指针漂（`A442` 那一类）。115 三枚 commit（`391878a`/`527d303`/`c6dbbf9`）均未碰 `internal/tools/**`。
4. "不许为省事把判定放宽"这一半：`git show 527d303` 逐处读——六处全是把比对面换成 `noticeNamesTree`/`noticesAboutTree`，极性一字未动（`n > 1`、`n != 1`、`t.Errorf` 文案原样），`private_set_sid_windows_test.go` 那处换后仍 `cleared = append(cleared, n.Principals...)`。
5. 失败方向：`winsec_windows.go:103-107`「Failure direction is a false alarm rather than a false all-clear…」与读码一致（`ResolvePath` 出错 ⇒ `return false` ⇒ 不归属＝偏响）。

结论：**可翻勾**。⚠ **与票 230 的撞车已登记**：`230-…md:19`「票 115 要一张真裁决表：那 7 格必须由非实现者出一张 `docs/evidence/s1/115-*.md`…表出来之后才谈补勾」，而 `find docs/evidence -name '115*' | wc -l` 本腿现量＝**0**。本腿用的凭据不是"115 名下的表"，是同一目录里另一张非实现者读数表 `docs/evidence/s1/ci-step-readings-2026-09-22.md` ＋本腿在这张票的码/commit 上的现读。**要不要认这张表代替 115 名下的表＝编排者的裁，不是本腿的裁。**
三行定性：现象＝通知归属的比对面从调用方拼写改成树归属；本机被入侵证据＝无（只读文件与只读 git）；最坏后果形状＝若"按树"判错，两棵不同树的通知会被并成一条（归属丢失）⇒ 本腿专查了反半边用例在位（第 5 格）与失败方向偏响不偏松。

## 第 5 格 — `115-…-done.md:53`（AC#3 四枚红转绿＋新增行为用例）

判据原文（票面 `:51-53`）：
> `- [ ] **AC#3** 修完四枚红，并且**新增一条行为用例**钉住你选的语义："同一个文件夹的两种合法拼写（长名 / 8.3 短名）触发通知时，通知的**归属**不因拼写而丢失"。⚠ 反半边必须同时钉：**不同两棵树不能被当成同一棵**（这正是票 112 那道守卫的另一半，别把它改宽）。`

本腿现跑的尺与读数：
1. 四枚红→绿（非实现者读数表，逐枚点名）：`docs/evidence/s1/ci-step-readings-2026-09-22.md:241`（run `35603195107` step4 `FAIL=4`＋四枚红名逐字）→ `:243`（run `35608530583` `PASS +4 / FAIL −4，全部归 527d303(票115 AC#3 后半格)…剩下的 2 枚红正是那对新用例`）→ `:245`（run `35616790753` `RUN 不变、PASS +2、FAIL −2 = 恰好 c6dbbf9(票115) 把那两枚转绿`）。
2. 该表自陈口径 `:3`「只读取数。本表所有结论都指到 run id + job id + step 号 + step 名」＋ `:16` 的日志真伪四判据（run 35616790753：`HTTP/1.1 200 OK`、`Content-Length` 413232 ＝ body 413232、3044 行）。
3. 本腿**不引那个 42**：同表 `:250-257` 自己写死了「42 这个数在 CI 上已经永久不可复现 ⇒ 判票 115 结案不能盯 42，要盯的是 (i) `--- FAIL=0`；(ii) `TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree` 与 `TestNoticeAttributionKeepsTwoTreesApart` 在 `--- PASS` 列里——这两条在 run 35616790753（L288/L290）与 run 35737627205（L246/L248）里都成立」⇒ 本勾引 (i)(ii)。
4. 用例在盘（本腿自己读，不依赖任何表）：`grep -n "^func Test" internal/winsec/notice_attribution_115_windows_test.go`＝`:191 TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree`、`:240 TestNoticeAttributionKeepsTwoTreesApart`；正半体内含 `len(*got) != 1`⇒Fatal、`n.Path == ""`⇒Fatal、`strings.EqualFold(n.Path, tr.spelling[1])`⇒Fatal（仪器自证"没量到东西"就红）、`!answerNamesTree115(...)`⇒Errorf；反半边两棵树各种子、`len(*got) != 2`⇒Fatal、对外树逐拼写要求"只命中它自己那 1 条"（红点文案「spelling %q of the other tree matched %d notice(s) of this pair, want its own 1」）⇒ 票 112 那一半没被改宽。
5. 生产承载体：`winsec_windows.go:106 noticeNamesTree`、`:120 noticesAboutTree`。

结论：**可翻勾**。本腿未复跑这四枚（禁跑编译/门，`223-v2` 在飞），"转绿"那一半的证据全部出自上面那张读数表。同样受票 230 那张"115 名下表"的形式账约束，已在第 4 格登记。

## 第 6 格 — `115-…-done.md:66`（AC#7 结案前必须有一枚远程 run id）

判据原文（票面 `:68-69`）：
> `- [ ] **AC#7** 结案前**必须有一枚远程 run id**：这道门在 CI 的 test-windows step 4，本机绿不算结案证据（子代理不许 push ⇒ 你写完把"要看哪枚 run 的哪几步、期望什么"写进 next=，我 push 后你再取或让接续者取）。`

本腿现跑的尺与读数：
1. 非实现者读数表：`docs/evidence/s1/ci-step-readings-2026-09-22.md:69`「### 2.1 票 115 的结案 run —— **拿到，且逐字命中**」＋`:72`「**run 35616790753 / job `106389418059` / step4**（`##[group]Run bash scripts/winsec-tests.sh` 在该日志第 **186** 行）」＋`:79` 逐字贴第 **1456** 行「`winsec-tests.sh: four numbers (all from -v output): === RUN=82  --- PASS=42  --- FAIL=0  --- SKIP=0`」。
2. 编排者结案语（同一枚 id 的第二处独立盘上落点）：`git show -s --format="%h %an %ad%n%B" f5bbccd8` 现读 A101 第 1 句「① 票 115 结案依据 run `35616790753` / job `106389418059` / step4「Windows ACL sealing gate…」第 1456 行逐字 `82/42/0/0`，两枚目标用例 L288/L290 转 PASS（基线 `35608530583` 同一区间它们正是那 2 枚红）」。两条互不派生（一贴日志原文＋判据，一写结案语）。
3. `next=` 那一半在票面：`grep -n "⑦ AC#7 next=" 本票`＝`:272`，其下逐条写了"取新 run 的 `test-windows` job 看 step 4、期望两枚新名字进 `--- PASS`、期望那四枚仍红（当时）、step 5–8 仍 skipped 不是本票的账"。
4. 票面自陈过"当时还不翻"：`:391`「框仍不翻：…**35608530583** / step 4」⇒ 本机绿没被拿来顶账，本勾落在其后结案的那一枚。（以上票面行号＝本腿落笔时刻现量；每追加一行都会再往下漂，故一律同时带该行开头的字串。）

结论：**可翻勾**。⚠ **推翻 `done-fix-1` 分类之处（第 2 条）**：它把这格记为"缺一次把这枚 run 读数落到本票名下的确认（可以是 `docs/evidence/s1/115-*.md` 里的一行）"，本腿判语＝那是**形式落点账、不是证据账**：证据已由 `ci-reader-s22` 表与编排者结案语两处独立给出；若仍要 115 名下的表，那是票 230 AC#1 的账，本格翻勾不撤它。

## 第 7 格 — `89-…-done.md:262`

判据原文（票面）：〔待填〕
本腿现跑的尺：〔待填〕
读数：〔待填〕
结论：〔待填〕

## 第 8 格 — `89-…-done.md:310`

判据原文（票面）：〔待填〕
本腿现跑的尺：〔待填〕
读数：〔待填〕
结论：〔待填〕

## 第 9 格 — `89-…-done.md:345`

判据原文（票面）：〔待填〕
本腿现跑的尺：〔待填〕
读数：〔待填〕
结论：〔待填〕

---

## 汇总

- 可翻勾：〔待填〕
- 不能翻（缺读数）：〔待填〕
- 判不了（凭据在禁读目录）：〔待填〕
- 形状歪（纪律句／待办行）：〔待填〕
- 推翻 `done-fix-1` 分类的条目：〔待填〕
- 未做完／判不了：〔待填〕
