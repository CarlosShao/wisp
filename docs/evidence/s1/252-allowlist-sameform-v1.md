# 票 252 · AC#2「两侧同形」修法 · 非实现者对抗验收表（腿 `252-v1`）

被验收件：实现腿 `252-r1` 的 `internal/tools/paths.go`（`Canonicalize` ＋新函数 `sameFormOfUnresolved`）
与载具 `internal/tools/paths_shortname_252_r1_test.go`；提交链 `bbbb3593`→`34d89191`→`2b1a3071`→`ebd5da2b`。
本腿＝对抗验收腿，**不翻任何 AC 框**；判语只写"成立／不成立／附条件成立"。
本腿读过的全部自述（`.scratch/wisp/probes/252/r1/impl.md`）在表里都标了"自述"或"本腿复量"，两者不混。

---

## §0 起手锚（同发取，逐字）

| 项 | 读数 |
|---|---|
| 取锚时刻 | `2026-10-03 09:38:07 +0800` |
| `git rev-parse --short=8 HEAD` | `1fea1e97`（**不是**被验收的 `ebd5da2b`；见下一行） |
| HEAD subject | `145-c2: S4 reconfirm/overturn against r1/r2/r3 (fresh line reads)` |
| 被验收链与 HEAD 的关系 | 四枚 `git log -1` 逐枚命中；`git merge-base --is-ancestor ebd5da2b HEAD` ⇒ **YES**（别的腿在本票之后又推了 `1fea1e97`→`aaf94ee1`→…） |
| 全仓 porcelain 行数 | `405`（`.scratch/wisp/probes/252/v1/start-porcelain-fullrepo.txt`，全是他人台件与 `design/` 脏面；同内容另有一份早于本腿 `v1/` 目录建立的 `.scratch/wisp/probes/252/v1-tmp-porcelain.txt`，未删，见 §7 第 8 条） |
| 本腿写面基线 `git status --porcelain -- internal tools` | 起手 `0` 行（`.scratch/wisp/probes/252/v1/start-porcelain-internal-tools.txt`）；突变开工时该面出现 `1` 行 `?? internal/panel/inbound_roster_253_test.go`＝**`253-r5` 的活，非本腿**，本腿全程未碰 |
| 分支 | `dev`（`git rev-parse --abbrev-ref HEAD`） |
| 被验收的 `internal/tools/paths.go` 是否仍是 HEAD 版 | **是**：`git diff 2b1a3071 HEAD -- internal/tools/paths.go` 空；blob@2b1a3071＝blob@HEAD＝工作区，md5 都是 `7a86da7aeba8420639491cd252a47a0c`（`.scratch/wisp/probes/252/v1/q1-paths-go-lineage.txt`）⇒ **本腿盘上读的就是被验收那份** |

骨架（§0＋§1 表头）先落盘：commit `3138e012`。

本腿自己复量的门禁终值（全部突变还原之后，`.scratch/wisp/probes/252/v1/final-*.log`）：

| 门禁 | 本腿读数 |
|---|---|
| `GOFLAGS= go build ./...` | `BUILD_EXIT=0` |
| `GOFLAGS= go test ./internal/tools/ ./internal/risk/ -count=1` | `ok internal/tools 15.881s`／`ok internal/risk 5.297s`，`--- FAIL` 枚数 **0** |
| `go test -count=1 -v -run 'TestTicket252R1' ./internal/tools/` | `--- PASS` **18**／`--- FAIL` **0**／SKIP **0**（`base-252r1-roster.log`）＝自述"18/0/0"**复量为真** |
| `go test -count=1 -v -run 'TestTicket252P1' ./internal/tools/` | `--- PASS` **9**／`--- FAIL` **0**／SKIP **0**＝自述复量为真 |
| `$(go env GOPATH)/bin/gofumpt.exe -l internal/tools internal/risk` | 输出 **0 行** |
| `./tools/d22scan/d22scan.exe` | `D22_EXIT=0`，末行 `d22scan: clean - no D22 ban violations`。⚠ 该跑里另有一行 `skipped as git-ignored: 1 file(s) … [frontend/dist/assets/]`（本腿终态跑同样命中 1 次）——按派单"SKIP 算红"的口径本腿**不替它判无害**，只具名它存在，见 §6 甲-9 |
| 载具 `t.Skip` 枚数 | **0**：`grep -nE '(^|[^.[:alnum:]_])t\.Skip\(' internal/tools/paths_shortname_252_r1_test.go` 命中 `0`（唯一的 `t.Skip` 字样是 `:39` 的注释自陈） |

⛔ 本机 staticcheck 产假绿，本腿未跑。⛔ 本腿未跑 `cmd/wisp`／`internal/config`／`internal/panel` 任何一枚测试，未跑全仓 `go test ./...`。

---

## §1 逐格判语（每格三选一＋凭据）

| # | 格 | 判语 | 凭据（本腿自己取数，未抄自述） |
|---|---|---|---|
| 1 | 「两次包含没合一」 | **成立**（附一条仪器告警，见下） | `internal/tools/paths.go` 现读：`InAllowlist`＝`:189-219`，内含 `rootsContain` **3 次**——`:196`（词法腿，对 `p.roots`）／`:209`（解析腿，对 `p.roots`）／`:215`（票 92 的 workspace 收窄，对 `p.workspace`，**不是**第三次根包含）。⇒ 对 `p.roots` 的包含仍是**两段**。本腿自算段体哈希：`2b1a3071^` 的 `:133-163` 与 `HEAD` 的 `:189-219` 都是 `e4aac2d53fe6bca6f43e5505b03e4202`，且 `diff` 输出 `IDENTICAL`；连注释块（`:129-163`／`:185-219`）都是 `b4a7e9bb6f8cd015c1701eea89b4d83b`。`git diff -U0 2b1a3071^..2b1a3071 -- internal/tools/paths.go` hunk 枚数 **1**＝`@@ -127 +127,57 @@`（全文 `.scratch/wisp/probes/252/v1/q1-inallowlist-body.txt`、`q1-paths-go-lineage.txt`） |
| 1b | ⚠ 那条硬禁**还有没有仪器** | **不成立（指仪器侧）** | 本腿亲手种 **M-merge**：把 `:196-198` 那一段词法包含整块删掉（＝票面"绝不许"的那件事实做出来），整包 `./internal/tools/` 跑＝**257 PASS／0 FAIL**，`./internal/risk/` 亦 `ok`（`mut-MERGE-run.log`、`mut-MERGE-risk.log`）。⇒ **修完之后，盘上没有任何一把尺能区分"两段"与"一段"**（原因见 §2：折叠保证了同形，`measure252` 的 `final == leg1 && leg2` 恒等式因此恒成立）。这条不改判第 1 格，但要单独报编排者 |
| 2 | fail-closed「解析不到＝未授权」没被改软 | **成立** | 那一支原样在 `internal/tools/paths.go:208-211`（`rf, ok := resolvedForm(canonical)`／`if !ok || !rootsContain(...)`／`return false`），且第 1 格已证段体逐字节未动。**M2 本腿自己复跑**（`if !ok \|\| …` → `if ok && …`）：`--- FAIL` **3**／`--- PASS` **254**，红句逐字见 §4 的 M2 行；其中**既有件** `TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`（`internal/tools/paths_ticket107b_probes_test.go:180`）确实红 ⇒ 那把尺对"改软"敏感，不是形式上保留。机制也复认：M2 日志里 `N7` 读数是 `ok=false rf=""`，而票 107 那枚探针在 M2 下的 `resolvedForm` 同样走 `ok=false`（`paths.go:300-302` 自陈"there, but I could not look through it"），拒绝确实来自 `!ok` 那一支 |
| 3 | ★落点偏差算不算遵守裁的形 | **附条件成立** | 见 §2（一句话裁决＋凭据行） |
| 4 | 存在性分岔有没有被改出新岔 | **成立**（没改出新答案岔；实现路径有一岔但输出同形） | 前一枚只读腿定的分岔＝存在性（`internal/risk/pathresolver.go:138` 缺失叶子照抄原拼法）。本次**没动 `internal/risk` 一个文件**：`git diff --name-only bbbb3593^..ebd5da2b -- internal/risk` ＝ **0 行**。本腿用现有载具的 `-v` 日志做四格对照（`.scratch/wisp/probes/252/v1/q4-existence-fork-readings.txt`，逐字）：<br>·存在的目录＋短名 ask（p1 的 `P2`/`P3`）⇒ `Canonicalize = "D:\\work\\workspace\\projects plans\\Wisp"`（长）<br>·缺失叶子＋短名 ask（p1 的 `Q1`/`Q3b`、r1 的 `A1`/`A3`）⇒ `Canonicalize = "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt"`（长前缀＋逐字尾巴）<br>·缺失叶子＋长名 ask（`P1`/`A2`/`A4`）⇒ 同一串长名<br>⇒ **两支的输出同形、`InAllowlist` 同答（`true`）**，两形四组合（A1–A4）也同答。新增的岔在 `paths.go:145-148`（`res.Resolved` 真⇒原样返；假⇒走折叠），那是**代码路径的岔、不是答案的岔**；票 252 要的"同一物理路径同一个答复"在判定面成立 |
| 5 | 正控的恒真自查（换形还绿不绿） | **成立（载具有牙，不是恒真）** | 三发都在 §4：**M5c**＝把 `paths_shortname_252_r1_test.go:126`（两侧同形串比较）与 `:130`（两侧同答）两枚断言换成恒真开关，**同时**种回缺陷（M1）⇒ 仍然 `--- FAIL` **9**，且 `TestTicket252R1BothSpellingsAnswerTheSame/{Wisp,docs,tools}` 依旧红，红句换到 `:136`（`AC#2 RED: both spellings must be authorized inside their own root, got LONG=true SHORT=false`）与 `:148`（`AC#1's split is still live: the two containments disagree on one input`）⇒ 恒真化之后同一格仍抓得住，**不是靠那两枚断言撑着**。<br>**M5d**＝在**修好的**生产码上把 `:126` 倒置成 `==` ⇒ 立刻 4 枚红（母格＋3 子格）⇒ 该断言在读真值、可失败，非空转。**M5b**＝恒真化＋生产码修好 ⇒ 4 PASS（与基线一致，说明开关没把测试挖空成永真失败） |
| 6 | 有没有越界 | **成立（零越界）** | 见 §3 名册：四枚 commit 合计 20 枚文件，**非 `.scratch` 的只有 2 枚**（`internal/tools/paths.go`、`internal/tools/paths_shortname_252_r1_test.go`），点名禁改清单命中 **0**；票面 6 枚 AC 框在盘上仍是 `- [ ]`（`:29-34`），票面最后一次改动是编排者的 `f25f9572`，实现腿一枚未碰 |

---

## §2 ★第 3 问：偏差裁决（算不算遵守编排者裁的形）

**裁决（一句话）**：**算遵守——但只在"生产面"这一半算，另一半（把裁语字面塞进 `InAllowlist`）本腿实测会打红票 252 自己的 AC#1 仪器，因此落点判**附条件成立**，条件＝把"进 `InAllowlist` 的串必先过 `Canonicalize`"这条不变量补一把尺钉住。

凭据与推理，逐条带行号：

1. **裁语的字面落点本来就是 `Canonicalize`**。票面 `.scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-short-name-stays-short-long-name-folds-long-so-in-root-writes-read-as-out-of-bounds-r2-l2.md:83` 逐字＝"**把进 `InAllowlist` 之前那一步补齐成'两侧同形'**——asked 与 roots 一律走同一条已存在的折叠"。落点在 `internal/tools/paths.go:119-149`（`Canonicalize` 尾部 `:145-148`）＋新函数 `:169-184`（自述写的是 `:151-186`，那是含注释块的跨度；函数体实跨 `169-184`，`:183` 才是 `return rf`——**行号小漂移，内容对得上**，另见 §6 甲-7）。
2. **roots 那一侧本来就已同形**，所以裁语里"asked 与 roots 两侧"在实现上只剩一侧要补：`internal/tools/paths.go:56`（根也过 C26）＋`:93`（`p.roots = append(p.roots, foldPath(res.Canonical))`）；`252-p1` 已量到"存在的树必折长、含 roots"（票面 `:82`），本腿在 `Roots()` 读数里复见（载具 `paths_shortname_252_r1_test.go:247` 也钉着这条）。
3. **残余不对称在生产面不活**——这是实现腿自己没敢定的那一枚（`impl.md:187-189` §5 甲-2），本腿追到底：三个判定方**全部**先 `Canonicalize` 再 `InAllowlist`：
   - `internal/risk/rules_gateway.go:37` → `:45`；
   - `internal/tools/task.go:834` → `:838`；
   - `internal/tools/bridge.go:290`（`dec.Paths = b.displayPaths(rawPaths)`）→ `bridge.go:928-942` 的 `displayPaths` 里 `:934` 逐枚 `b.paths.Canonicalize(p)` → `bridge.go:1145` 的 `b.paths.InAllowlist(p)`；失败支以 `" (无法规范化: …)"` 后缀入串并被同一行的 `strings.Contains(p, "无法规范化")` 单独拦成不在册。
   - `cmd/` 里非测试码对 `InAllowlist` 的调用＝**0 枚**（`grep -rn InAllowlist cmd --include=*.go \| grep -v _test` 计数 0）。
   ⇒ 它自述"把短名字面直接递给 `InAllowlist` 这条路只剩测试面"**为真**（测试面的确还留着：`paths_shortname_252_probe_test.go:114-126` 的 `measure252` 用词法串重算 leg1、`paths_ticket107b_probes_test.go:85-93` 的 `judge107b` 在 `Canonicalize` 报错时回落成词法串）。
4. **它隐含的那句否证，本腿复跑了实际作用面才入账**。种 **M3**＝按"字面喂给两条腿"的另一种落点实现（`Canonicalize` 摘掉折叠、把 `sameFormOfUnresolved` 塞进 `InAllowlist` 第一段 `:190`）：
   - 结果 `--- FAIL` **12**／`--- PASS` 245，红句逐字（`mut-M3-run.log`）：`paths_shortname_252_probe_test.go:252: leg breakdown does not explain the boolean: InAllowlist=true but leg1=false leg2=true`（同一句还在 `:272` 与 `paths_shortname_252_r1_test.go:98` 各响若干枚）＋ `paths_shortname_252_r1_test.go:127: one physical path still reaches the judge as two shapes`。
   - ⇒ **"两段塌成一段"在盘上的真实表现＝第一段不再是"词法腿"**，被 AC#1 那把恒等式尺（`final == leg1 && leg2`）当场抓住。**这一支复认成立**，可以入账。
   - 但它另半句"且新增一次词法在根外／解析落回根内的授权"——**本腿量不到**：M3 之下 `N1..N7` 七枚负控全 `--- PASS`、`TestTicket107bProbeA/B/C` 三枚全 `--- PASS`、`AC#3 RED` 在整包日志里计数 **0**。现有载具里没有"根外经入向链接折回根内"的载体，所以那半句是**论证、不是读数**，⛔ 不许当凭据引用（§5 建议里请编排者若要它成立就派一枚入向 junction 载具）。
5. **两条硬禁都没被绕**：合一＝没做（§1 格 1，段体逐字节未动）；改软＝没做（§1 格 2，M2 复跑到红）；判级策略（R2→L2）与本票射程外的 `internal/risk`＝零文件。
6. **反过来说，选这个落点也不是零代价**：正因为折叠在 `Canonicalize`，`InAllowlist` 才第一次出现"两段包含的参数值恒等"的状态 ⇒ 第 1b 格那条"绝不许合一"从此**无人值守**。这不是实现腿违令，是裁语选形的必然后果，但**它必须被记账**，否则下一枚腿把 `:196-198` 删掉时 CI 是绿的。

---

## §3 越界名册（`git diff bbbb3593^..ebd5da2b` 逐枚）

⚠ 方法说明（本腿自己踩到的坑，具名写下）：**按区间取数会被别的腿污染**——`bbbb3593^..ebd5da2b` 的 `--name-only` 列出 60+ 枚，里面有 `docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`、票 242/244/253/259 的票面、`docs/evidence/s1/242-grant-binding-v1.md` 等**都不是本票四枚写的**（四枚提交在时间线上不连续，区间把别人的活圈进来了）。⇒ 判越界只能用**逐枚 `git show --name-only`**，本腿两者都取了：

- 逐枚并集去重＝**20 枚文件**（`.scratch/wisp/probes/252/v1/four-commit-files.txt`），其中**非 `.scratch` 只有 2 枚**：
  `internal/tools/paths.go`（1 枚 hunk，见 §1 格 1）／`internal/tools/paths_shortname_252_r1_test.go`（新增）。
- 点名禁改清单逐条比（`cmd/wisp`／`internal/config`／`internal/panel`／`internal/agent`／`docs/PLAN.md`／`docs/specs`／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／`.github/workflows`／`frontend/`／`design/`）⇒ 命中 **0**。
- `internal/risk` 在区间内改动文件数＝**0** ⇒ SLO／阈值／判级面本票没碰。
- 票面 6 枚 AC 框：`git log --oneline -- .scratch/wisp/issues/252-*` 最新一枚是编排者 `f25f9572`，实现腿四枚均未带该路径 ⇒ **框未碰**（盘上 `:29-34` 逐枚仍 `- [ ]`）。
- 载具文件的改动分布（本腿逐枚 `git diff` 取出）：`34d89191` 新建 282 行；`2b1a3071` 改 24+/12− **全是 N1–N6 结构体字面量的 gofumpt 换行**（无判据变化）；`ebd5da2b` 改 3+/1−，把 `:252` 那条红句文案从 `carrier moved under the test: …` 改写成 `AC#2 RED: the short spelling … the same-form step is not in effect (leg1=%v leg2=%v folded=%q)`——**条件式未变**（仍 `if !lShort.final`）。⚠ 一枚"docs 提交里带测试改动"的形状本腿如实登记，不替它圆场（见 §6 甲-7）。

**结论：零越界。**

---

## §4 突变自证表（种什么→哪枚必须红→红句逐字→还原复跑终值）

纪律执行记录：每发前都从 `git show HEAD:<path>` 取备份并记 md5（`.scratch/wisp/probes/252/v1/backup/backup-md5.txt`，起手 prod `7a86da7aeba8420639491cd252a47a0c`／carrier `cc9e50e4533fc5d0c37b7826ac982586`）；跑完**立刻**`git cat-file blob HEAD:<path> > <path>` 还原；两发之间量 `git status --porcelain -- internal tools`。⛔ 全程未删任何文件（备份与日志只建不删）。起手态里那 1 枚 `?? internal/panel/inbound_roster_253_test.go` 是 `253-r5` 的，本腿从第一发到最后一发都没让它出现在自己的 pathspec 里。

| 号 | 种的形（file:line） | 突变后 md5 | 哪枚**必须**红 | 红句（逐字，只截本腿判定所需那段） | 还原读数 |
|---|---|---|---|---|---|
| **M2**（派单点名要本腿复跑） | `internal/tools/paths.go:209` `if !ok \|\| !rootsContain(p.roots, foldPath(rf))` → `if ok && !rootsContain(...)`（＝把"解析不到＝未授权"改软） | `1dd3a51a654b936d94367da74f346b6c` | 既有件 `TestTicket107bProbeCLinkInsideAllowedRootStaysOutside` **必须红**（不红＝那把尺对本修法不敏感） | `paths_ticket107b_probes_test.go:180: AC#3 RED: InAllowlist("C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside670927220\\001\\proj\\esc\\marker.txt") = true: it leaves the allowed root "…\\001\\proj" through the link "…\\001\\proj\\esc" and lands on "…\\001\\proj\\esc\\marker.txt"` <br>＋ `paths_shortname_252_r1_test.go:234: AC#2 RED (widening): an unresolvable path inside the root was authorized (leg1=true leg2=false ok=false rf=""): "cannot be resolved means not authorized" must stay` | `--- FAIL` **3**／`--- PASS` **254**；还原 md5 `7a86da7a…`、`git status --porcelain -- internal/tools`＝**0 行**、定向复跑 `TestTicket252R1\|TestTicket252P1\|TestTicket107b` 退 0（`restore-check-after-merge.log` 同形） |
| **M-merge**（票面第一号硬禁有没有牙） | 删掉 `internal/tools/paths.go:196-198` 整段词法包含（＝把两段并成一段），`f` 仍被 `:215` 使用 ⇒ 编译通过、非"编译失败当红" | `438fe50ae700e87d878901588eff80d8` | **没有任何一枚红**＝本腿要量的那一格 | 整包 `./internal/tools/`＝`--- FAIL` **0**／`--- PASS` **257**；`./internal/risk/`＝`ok … 5.100s`。⇒ **结论：修完之后"绝不许合一"这条硬禁在可跑面上无人守**（`measure252` 恒等式因两侧同形而失去鉴别力） | 还原 md5 `7a86da7a…`；定向复跑退 0 |
| **M1**（自述复量：同形那一步摘掉） | `internal/tools/paths.go:183` `return rf` → `return canonical`（`rf`/`ok` 仍被 `:180` 使用） | `cac483cee6b43221ec83348fefbcfc63` | 依赖折叠的 9 枚 | 名册＝`TestTicket252R1ShortSpellingOfNewFileIsAuthorized`(母)＋`/A1`＋`/A3`、`TestTicket252R1BothSpellingsAnswerTheSame`(母)＋`/Wisp`＋`/docs`＋`/tools`、`TestTicket252R1AlignmentAddsNoAuthorization`(母，尾判)、`TestTicket252R1CanonicalStillNamesOneTree`。红句逐字：`paths_shortname_252_r1_test.go:98: POSITIVE CONTROL RED for A1 roots=LONG asked=SHORT+missing-leaf: same-shape ask was refused, so this probe's carrier is broken and no negative reading from it can be trusted`；`…:127: one physical path still reaches the judge as two shapes: "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-round-trip.md" vs "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252r1-round-trip.md"`；`…:252: AC#2 RED: the short spelling of a new file was still refused under its own root ("D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252r1-out-of-scope.md"): the same-form step is not in effect (leg1=false leg2=true folded="d:\\work\\worksp~1\\projec~1\\wisp\\q252r1-out-of-scope.md")`；`…:276: two spellings of one path still canonicalize to two strings: …` | `--- FAIL` **9**／`--- PASS` **248**＝**自述枚数复量为真**；同发 `TestTicket252P1*` 9 枚仍 `--- PASS`（它是读数件，不持期望）⇒ 与自述一致。还原 md5 `7a86da7a…` |
| **M3**（★裁语另一落点的作用面） | 两枚一起改：`paths.go:148` `return sameFormOfUnresolved(res.Canonical), nil` → `return res.Canonical, nil`；`:190` `f := foldPath(canonical)` → `f := foldPath(sameFormOfUnresolved(canonical))` | `e2a269e33b6b95890ff82ca222bd624f` | 第一段不再是词法腿 ⇒ AC#1 仪器必须红；**同时**看负控会不会被放宽 | `--- FAIL` **12**／`--- PASS` 245。红句逐字：`paths_shortname_252_probe_test.go:252: leg breakdown does not explain the boolean: InAllowlist=true but leg1=false leg2=true (the workspace leg at paths.go:159 is inert in this probe - if it is not, the carrier changed under us)`（同句在 `:272` 与 `paths_shortname_252_r1_test.go:98` 再响）＋`paths_shortname_252_r1_test.go:144: leg breakdown no longer explains the boolean: LONG true/true/true SHORT false/true/true`＋`…:148: AC#1's split is still live: the two containments disagree on one input (LONG leg1=true leg2=true, SHORT leg1=false leg2=true)`。<br>⚠ 反向那一半没抓到放宽：`N1..N7` 七枚全 `--- PASS`、`TestTicket107bProbeA/B/C` 三枚全 `--- PASS`、`AC#3 RED` 计数 **0**；`./internal/risk/` 仍 `ok` | 还原 md5 `7a86da7a…`；定向复跑退 0；`git status --porcelain -- internal tools`＝回到 `1` 行（那枚 `internal/panel` 他人件），非本腿 |
| **M5b/c/d**（正控恒真自查） | 载具 `paths_shortname_252_r1_test.go:126`／`:130` 两枚断言各加运行时开关 `&& os.Getenv("WISP252V1") != "neutralize"`（用非常量防 vet 常量折叠）；M5d 再把 `:126` 的比较**倒置**成 `foldPath(cLong) == foldPath(cShort)` | carrier `cb92ac9293632a3aa4b5d8e800a5ec96` | 恒真化后若仍全绿＝该判据不敏感 | **M5b**（恒真化＋生产码已修）：`--- PASS` 4 枚（母格＋`/Wisp`＋`/docs`＋`/tools`）⇒ 基线一致。<br>**M5c**（恒真化＋**再种 M1**）：`--- FAIL` **9**／`--- PASS` 18，且 `BothSpellings/{Wisp,docs,tools}` 仍红，红句变成 `paths_shortname_252_r1_test.go:136: AC#2 RED: both spellings must be authorized inside their own root, got LONG=true SHORT=false (roots=[D:\work\workspace\projects plans\Wisp])` 与 `…:148: AC#1's split is still live: the two containments disagree on one input (LONG leg1=true leg2=true, SHORT leg1=false leg2=true)` ⇒ **换形不是恒真、载具仍有牙**。<br>**M5d**（`:126` 倒置＋生产码已修）：4 枚 `--- FAIL` ⇒ 该断言能失败、在读真值 | carrier 还原 md5 `cc9e50e4…`＝起手值；prod 全程 `7a86da7a…`（除各发突变自身）；末态 `git status --porcelain -- internal tools`＝**0 行** |

发与发之间的量面（逐发记）：`internal tools` 回到 0 行（起手那枚 `internal/panel` 他人件在 M3 还原之后已不在名册里＝`253-r5` 自己 commit 掉了，与本腿无关）。
**终态复跑**（所有突变还原之后）：`go build ./...` 退 0；`./internal/tools/`＋`./internal/risk/` 两枚 `ok`、`--- FAIL` 0；`gofumpt -l` 空；`d22scan` 退 0。

### 4-乙 自述与盘上不符之处（只列本腿复量到的）

1. **枚数全中**：`InAllowlist` 前后哈希（本腿自算同值 `e4aac2d53fe6bca6f43e5505b03e4202`）、`-U0` 一枚 hunk、改前 9 FAIL／改后 18 PASS、M1 的 9/248、M2 的 3/254、18/9/0 SKIP、gofumpt 空、d22scan clean、0 枚 `t.Skip`——**全部复量为真**。
2. **M2 第三枚红的命名错了**（`impl.md:174`）：它写"红在 `N7` ＋**母格尾判**＋既有件 `TestTicket107b…`"。本腿复跑：`:252` 那条尾判红句在 M2 日志里**一次都没出现**（`grep -c 'paths_shortname_252_r1_test\.go:252:' mut-M2-run.log`＝**0**）；第三枚 `--- FAIL` 是 `N7` 的**母格聚合行**，不是尾判在响。⇒ 计数 3 对、归因错一格。（同一枚 `:252` 在 **M1** 下发里确实响，见 §4 M1 行——两发行容易看串。）
3. **行号小漂移**（不改实质）：`impl.md:81` 的"`sameFormOfUnresolved` 在 `:151-186`"含注释块，函数体实跨 `:169-184`；`impl.md:173` 说 M1 改的是"`paths.go:185` 的 `return rf`"，盘上 `return rf` 在 **`:183`**；`impl.md:108` 说 `unifySeparators` 在"`:349` 附近"，盘上在 **`:338`**（`:349` 落在 `foldPath` 头上）。
4. `impl.md:24` 起手写"本腿写面全程保持空"，盘上 `cmd/wisp` 4 枚脏件／`internal/config` 1 枚新增＝他人——与本腿读数**不冲突**（本腿同样一枚没跑）。

---

## §5 我建议编排者怎么处置

1. **可以翻**：AC#2、AC#3。凭据＝§1 格 1／2／4／5／6 ＋ §2 ＋ §4 的 M1/M2/M3/M5 四发。AC#3 那一格本腿复见到位：真短名载具、本机可判、0 枚 `t.Skip`、未挂〔仅 CI 可量〕，且改前 9 枚红／改后 18 枚红转绿是**同机同载具**（`base-252r1-roster.log` 与 §4 M1 行）。
2. **不许翻**：AC#1 之外的那三格——**AC#4／AC#5／AC#6 一枚未做**（`impl.md:212-215` 具名原因是 `cmd/wisp` 写面被占，本腿认同该理由成立，因为它没跑那包是派单硬禁）。AC#1 本身早在 `252-p1` 就结清（票面 `:80`），本票这一格不必等本腿。
3. **两枚缺的读数，建议指名去量**（都不该记在 `252-r1` 名下，也不该由本腿补）：
   - 甲：**"绝不许把两次包含改成一次"现在没有仪器**。本腿 M-merge 的 257 PASS／0 FAIL 就是证据。建议派一枚结构尺（钉"`InAllowlist` 第一段是词法、第二段是解析形"，或钉 `measure252` 恒等式对"合一"敏感），否则下一枚腿删掉 `:196-198` 时 CI 绿。
   - 乙：`impl.md:123-124`／§9(a) 的后半句（"折叠挪进 `InAllowlist` 会新增一次词法在根外、解析落回根内的授权"）**没有载具**——本腿 M3 量了，`AC#3 RED` 计数 0。若编排者要把这句话写进结案，需要先有人造一枚**入向 junction**（根外→折回根内）载体；不需要的话就该明说"该句为论证，非读数"。
4. **AC#4 那枚 301 秒 witness 的归口仍未结**（`cmd/wisp` 面，本腿与实现腿都没跑）⇒ 别把本表读成那笔账有了答案。
5. ⛔ 提醒一条判据边界：本表的加分只落在"拼法一致性"射程内。`allowed_dirs` 仍是**判级输入**、不是执行时硬边界（既有裁定），本腿**没有**把任何"变严／顺手做成硬边界"的迹象算成加分，也不建议把这一格往那个方向结。
6. 建议对实现腿的态度：**结 AC#2／AC#3，不退回**。它自请裁定的那一格（`impl.md:252` (a)）本腿判"附条件成立"，条件见 §2 末段与上面第 3 条甲。

---

## §6 我不确定的地方（自我对抗，不空着）

甲-1｜**`cmd/wisp` 一枚没跑**（写面被 `255-r2` 占）。⇒ AC#4／AC#5／AC#6、那枚 301 秒 witness、以及 `in_allowlist_scope=` 审计行的真实形状，本腿**全部无从判**；本腿只在本包的 `fswrite_silentloss_ac1_test.go:193` 日志里见过一例 `in_allowlist_scope=true`，那不是跨面证据。
甲-2｜**M-merge"全绿"的射程只到 `internal/tools`＋`internal/risk`**。若 `cmd/wisp` 里藏着直接喂字面串给 `InAllowlist` 的载具，本腿看不见。旁证（非结论）：非测试码里 `InAllowlist` 只有三枚调用点、`cmd/` 命中 0；但面板侧经 `PanelBridge` 走到 `bridge.inScope` 这条链本腿**没逐字追**（只追到 `bridge.go:290→934→1145`）。
甲-3｜**POSIX"逐字节不变"是推断不是读数**（`paths.go:170` 的 `filepath.Separator != '\\'` 编译期常量分支；派单禁全仓跑、本机无 POSIX）。同一条实现腿也自陈没量（`impl.md:190` 甲-3）。
甲-4｜**`res.Resolved == false` 的真实覆盖面**：ACL 挡句柄、离线卷、`ERROR_INVALID_NAME` 这三类"存在却没解析出来"，本腿造不出真形。本腿只证到守卫 `paths.go:173` 用 `os.Lstat` 把"存在"挡回原串、以及 `N7` 那一枚量到 `ok=false rf=""`。
甲-5｜**存在性竞态**：折叠读的是"问的这一刻链上存不存在"。`Canonicalize` 与真 `open` 之间树被改动时形状可能不同。本腿没测（实现腿 `impl.md:195` 甲-5 也登记未结）。
甲-6｜**大小写／磁盘拼法副作用**：折叠后的形取 OS 自己的拼法。判定面因 `foldPath` 折小写不受影响，但 `Canonicalize` 的**输出串**是要真开的那串、也进审计行 ⇒ CI 上 `C:\Users\RUNNER~1` 那族若磁盘大小写与 `os.TempDir()` 给的不同，输出会变（`impl.md:197` 甲-6 同一件事）。**CI 那一发本腿没有读数。**
甲-7｜**"docs 提交里改测试"**：`ebd5da2b` 改了载具 `:252` 的红句文案（条件式未变）。本腿的 18 PASS 是在**当前盘上版本**复量的，所以凭据仍成立；但"改前 9 枚红"那份名册是 `34d89191` 版载具的读数，两版之间文案不同名册相同——本腿没去逐枚比对 9 枚红句在 `34d89191` 版的原文（那需要一次跨提交复跑，本腿选择不用它当任何判语的唯一凭据）。
甲-8｜**§2 那句"roots 侧本来就已同形"依赖"根必然存在"**。若 `allowed_dirs` 配了一台不存在的目录：`NewPathCanonicalizer`（`paths.go:56-59`）会把它丢进 `unusable` 并**不**入册 ⇒ 那一侧根本不参与比较。所以"两侧同形"是"在册根都被 OS 确认过"前提下的同形，本腿没有量过"根不存在但配置里有"那一发在修法后的行为（`paths_ticket107_portable_test.go`／`pathshape_portable_test.go` 里可能有相近形状，本腿没逐枚读）。
甲-9｜**d22scan 的 `skipped as git-ignored: 1 file(s) … [frontend/dist/assets/]`**：本腿终态跑同样命中 1 次。按派单"SKIP 算红"的字面口径，本腿**不替它判"与测试无关因此无害"**，只具名：它不是本票四枚带进来的（`git log --oneline -- frontend/dist` 不在四枚名册里），且 `frontend/**` 内容本腿未读未引。
甲-10｜**票面 §25「其余 21 枚共红里有多少同因」**：本腿一枚未归因，⛔ 别把本表读成那 21 枚有了答案。

---

## §7 判语

1. **AC#2 的修法成立**：允许根判定对同一物理路径的两种拼法给出同一个答案，落点在裁语写的"进 `InAllowlist` 之前那一步"（`internal/tools/paths.go:119-149`＋`:169-184`，复用已有 `resolvedForm`（`:303`））；票面两条硬禁**都没碰**——两段包含仍是两段（`InAllowlist` 段体改前改后逐字节相等，本腿自算），"解析不到＝未授权"仍是 `return false`（`paths.go:208-211`）。
2. **★偏差＝附条件成立**（§2）：不算绕过裁语，因为它落的那一步正是裁语点名的那一步、且生产三枚判定方都先 `Canonicalize`（`rules_gateway.go:37→45`／`task.go:834→838`／`bridge.go:290→934→1145`，含实现腿自陈没追到底的 `bridge` 那一枚）；条件＝"合一"这条硬禁修完之后失去仪器（M-merge＝257 PASS／0 FAIL），须另派尺补钉，否则不变量只靠约定。
3. **它上交的那半句否证**："两段塌成一段"本腿已用 M3 复跑到红（AC#1 的恒等式尺抓住"第一段不再是词法腿"）⇒ 可入账；"会新增一次授权"那一半**本腿量不到**（M3 下 `N1–N7`＋`107b A/B/C` 全绿、`AC#3 RED` 0 枚）⇒ 本表**不为它背书**，§5 第 3 条乙具名请编排者定夺。
4. **AC#3 载具有牙**（§1 格 5）：恒真化两枚同形断言＋种回缺陷仍 9 枚红；把 `:126` 倒置在修好的码上立刻红；真短名、0 枚 `t.Skip`、未挂〔仅 CI 可量〕。
5. **门禁复量全中**：build 0／两包 `ok`／0 FAIL／gofumpt 空／d22scan clean（那枚 git-ignore 行按 §6 甲-9 具名不解释）。
6. **越界＝零**（§3）；`internal/risk` 零文件；阈值／golden／`allowlist.txt`／`ci.yml`／三枚点名既有判据件一字未动（本腿也只读内部判射程，未引用其内容）。
7. **未做与不该翻**：AC#4／AC#5／AC#6 一枚未做（`cmd/wisp` 写面被占），那枚 301 秒 witness 的归口、"其余 21 枚共红"的同因计数，本腿同样**没有**答案。
8. **纪律**：**AC 框一枚未碰**（盘上 6 枚仍 `- [ ]`，本腿与实现腿都没动过票面）；**未 push**（本腿只在本地 commit，`git remote` 未动）；本腿在 `internal/tools` 里的全部改动都是**临时突变且已还原**（prod `7a86da7a…`／carrier `cc9e50e4…`，`git status --porcelain -- internal tools` 终态 0 行），**没有顺手补任何一枚尺**；台件只建不删，落在 `.scratch/wisp/probes/252/v1/**`。⚠ 一处如实登记：`v1/` 目录建立之前，本腿为取起手锚在上一层写过一枚只读台件 `.scratch/wisp/probes/252/v1-tmp-porcelain.txt`（内容＝全仓 porcelain，无本腿产码），按"仓内绝不删文件"它留在原地，同内容已另存进 `v1/start-porcelain-fullrepo.txt`。
