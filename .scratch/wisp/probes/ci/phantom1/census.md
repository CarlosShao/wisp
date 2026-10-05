# 普查件 ci-phantom-1 — HEAD 干净检出下 ban #9 `[phantom-citation]` 全名册（只读腿，⛔ 零 `go` 命令）

> 腿：`ci-phantom-1`｜派单目的：我 08:58 推上双远端之后，**新树上的 CI 会不会继续红在 `TestRealRepoLedgerIsHonest`**，以及那枚红的代价面。
> ⚠ **本件由编排者代落盘**（09:2x，与 `9ec258d9` 转录 `267-a1` 同形，今天第二次同坑）：派单又选错了腿型——只读会话**没有 Write/Edit、也没有 commit 能力**，所以它物理上写不出 `.scratch/wisp/probes/ci/phantom1/census.md`。正文逐字是腿的产出，编排者只加本节抬头与下面的复尺记录。**记我（第二次，同一条）**：给"要落文件"的腿派单前，先确认那个腿型有写工具；否则交付形状那一格从一开始就是不可满足的。
> ★ **编排者复尺（09:1x 同发现跑，含两处它顶回我的地方）**：
> ① **被引件 `:640` 那句是真读数**：`sed -n '638,642p' .scratch/wisp/probes/33/r6/fullpack2.log` ⇒ `:640` 逐字 `msg="panel host: pre-create message purge hit its cap before the queue was clean" cap=256 removed=256`，`:641` 是 `timed out after 15s waiting for the gesture's show request` ⇒ 注释里"purge 撞 256 上限、show 没走完"**成立**，缺的只是那枚文件没进版本库（修前 `git ls-files -- <该路径>`＝空）。
> ② **`testdata` 赦条逐字在盘**：`tools/d22scan/main.go:680-683`＝`if d.Name() == "testdata" || d.Name() == ".git" || s.ign.skip(path, true) { return filepath.SkipDir }` ⇒ 名册第 2 枚（`cmd/wisp/testdata/esclistener/main.go:9`）今天确实不在 ban #9 射程。⚠ 这条赦**只在 walkGo 源码里、ban #9 文档头 `:50-69` 没写**，且（它具名）没有 selftest 样本钉住它——⇒ 任何腿日后收窄 walkGo 射程就会让那行引用突然变红，这一格我接走（见台账 `A609`）。
> ③ **★顶回我第 1 处（数目）**：我在 `A608` 写的是"五步被 skipped"。它照 `ci.yml` 文本拉出来的是**七步**，并点名我漏的那枚正是**票 262 那道门**（`:136 Tracked path-length budget (ticket 262)`）。我 09:1x 现读 `ci.yml:74`/`:83`/`:136` 三处名字与位置复认为真 ⇒ **`A608` 那句"五枚"作废，按七枚读**（原句不抹，`A609` 里更正）。⇒ 后果比我写的更重：**票 262 的 AC#5（"该步在 CI 真出颜色"）被同一枚链头红压着**，它不是"等推送"，是"等这枚红消"。
> ④ **★顶回我第 2 处（行号）**：那枚白盒钉的真身现在是 `tools/d22scan/scan_test.go:2174`，不是我引用的 `:2207/:2211`——那两个号来自 run `37240161874` 在**旧树 `fd269de1`** 上的日志，两棵树各自的行号都对，但**我把它当常量抄了**（与票 267 那枚 `:310` 漂移同一族，定式＝`grep` 被指的字符、别抄行号）。
> ⑤ **它自报一枚形状越界**（我如实登记）：SEP3 那把全仓计数尺 `git grep -cFn <token> HEAD` 未带 pathspec ⇒ git 内部遍历过全树；产出的只是**文件名计数**、`design/**`/`frontend/**` 零枚命中、**未读任何内容**。此后我派单里的搜索尺一律钉双根。

---

## §0 起手锚与并发窗口

- 起手时刻 `2026-10-05T09:00:14+0800`；起手 `git rev-parse --short HEAD`＝`21bec8a1`；`git status --porcelain` 总行数＝**645**；branch＝`dev`。
- **HEAD 在本腿飞行期间动了两次**：09:02 落到 `7828c972`（ledger A608），09:08 复锚时已是 `9ec258d9`（`267-a1` 普查件转录）。终版枚举**钉死在 `9ec258d91b6ecf21d51ea82e75c12bb4bf6f2abd`** 复跑；`21bec8a1..9ec258d9` 共两枚 delta commit，`git diff --name-only` 过滤 `cmd|internal` 非 `_test.go` 的 `.go` 文件＝**0 枚**（0 枚＝没找到＝射程未动，好消息）。复跑 MISSING 名册三度一致。
- 09:08:27 复跑 `git status --porcelain | wc -l`＝**646**（脏面在长，共享树实况）。同机在飞：`267-a1`（写面 `.scratch/wisp/probes/267/**`，与本腿零交集；它的转录 commit `9ec258d9` 恰成我落不了件的先例）。
- 本节写下那一刻的行号只对 `9ec258d9` 树负责。

## §1 仪器真身复认（照 `tools/d22scan/main.go` 实际规则复刻）

判定主体＝`main.go:831-876` 的注释循环（在 `scanGoFile` 内），规则逐格：

- **射程**：ban #9 文档头 `main.go:50-69` 逐字——"internal/ 与 cmd/ 的生产 Go 文件 ONLY，`tools/**` 与 `_test.go` 在射程外"；执行面＝`Scan` 只 `walkGo(internal)`＋`walkGo(cmd)`（`main.go:264-269`）。
- **walkGo 的隐藏臂（名册第 2 枚的赦因，文档头没写）**：`main.go:682` `if d.Name() == "testdata" || d.Name() == ".git" || s.ign.skip(...) { return filepath.SkipDir }`——**`testdata/` 整树不进扫描**；`_test.go` 于 `:689` 被后缀滤掉。我的过滤器照 `264-269 + 672-702(含682) + 831-876` 复刻，不是照猜。
- **认什么形状**：注释文本取自 AST `f.Comments`（`//` 行一条一节点），token 由 `repoPathRe`（`main.go:890-891`）命中：前缀 ∈ {`docs`,`\.scratch`,`internal`,`cmd`,`tools`,`scripts`}＋`/`＋字符类 `[A-Za-z0-9_./\-一-鿿]*`＋尾 `[A-Za-z0-9_-]`（**CJK 在类内**——我的提取尺同步带 `\p{Han}` 跑了两遍，名册一致）。
- **什么算存在**：`os.Stat(filepath.Join(s.root, tok))`（`:858`），root＝仓根（`scripts/d22scan.sh` 派生）⇒ 干净检出下等价于"在 HEAD 树里"＝我的 `git ls-tree $PIN -- <tok>` 判据（文件与目录皆可判）。
- **豁免**：① 简写——`shorthandPathStarts`（`:932-942`，region 键＝region 起点，region 类比 repoPathRe 多 `*`/`…`，含 `...`/`…`/`*` 才免）；② API 引用——`symRefRe`（`:954-955`）`^[a-z][a-z0-9]*(/[a-z][a-z0-9]*)*\.[A-Z]`，`:843-855` 逐字点名 `internal/tools.Result.AppliedSteps`、`internal/proc.WithRegistry` 属"API citation 不 convict"；③ 去重——`citationReported` 每 token 每扫至多 1 枚 finding（`:840-841`）。
- **allowlist**：`ban-id<TAB>路径前缀<TAB>理由`（`:71-73`）；现 `tools/d22scan/allowlist.txt` 内 phantom 条目＝**0 枚**（grep -c 零命中＝没找到＝当前无既存豁免）。
- 钉真身：`TestRealRepoLedgerIsHonest` 在 `tools/d22scan/scan_test.go` **现树 :2174**（派单给的 `:2207/:2211` 已漂，A608 自己也登记过行号漂移）。

## §2 全名册（主表，钉 `9ec258d9`）

169 枚 unique token 全量过 `git ls-tree`，MISSING 5 枚，逐枚判：

| # | 引用所在 file:line | 被指名路径 | 磁盘有? | index 有? | 字节数（磁盘读得／index 无） | 简写尾缀形? | symRef/豁免 | **CI 判词** | 最小修法那一格 |
|---|---|---|---|---|---|---|---|---|---|
| 1 | `cmd/wisp/panel_host_windows.go:343` | `.scratch/wisp/probes/33/r6/fullpack2.log` | 有 | **无** | 197,213（被引件在盘；index 格＝无） | 否——token 后紧跟 `:640`，`:` 不在类 ⇒ 全拼 token、region 无标记 | 否 | **红，1 枚**（与 run 37240161874 同行同 token 逐字同形） | 三支见下：(a) commit 该 log（0 行产码）；(b) 改 :343 指向真实已跟踪路径（1 文件 1 行）；(c) allowlist 前缀条目（机制合法、属裁） |
| 2 | `cmd/wisp/testdata/esclistener/main.go:9` | `.scratch/wisp/probes/245/v1/esclistener/main.go` | 有（8,518） | 无 | 8,518 | 否（后随 `)`） | **射程外**——`main.go:682` testdata SkipDir；docs/.scratch 里的同名引用（`246-veto…md`×1、212 工作件）也不在 ban #9 walk | **静默，0 枚**（新推送 `0c1b32cf` 带来的新引用，被 testdata 规则赦） | 无需修；预警：任何腿若把这行引用搬出 testdata 即转红 |
| 3 | `internal/buildinfo/buildinfo.go:5`（`-X …/internal/buildinfo.Name` ldflags 注释） | `internal/buildinfo.Name` | 两者皆无 | 无 | — | 否 | **symRefRe 命中**（`.N` 大写）⇒ `:853` continue | **静默，0 枚** | 无需修（本就是 ldflags 路径非文件） |
| 4 | `cmd/wisp/resident_ball_windows.go:234` | `internal/proc.WithRegistry` | 两者皆无 | 无 | — | 否 | symRefRe 命中，main.go:845-846 逐字自认 | **静默，0 枚** | 无需修 |
| 5 | `internal/agent/approval/report.go:15` | `internal/tools.Result.AppliedSteps` | 两者皆无 | 无 | — | 否 | symRefRe 命中，main.go:945 逐字自认 | **静默，0 枚** | 无需修 |

**分母读数：新树 CI 会红＝1 枚**（#1）。其余 4 枚全部有具名赦因。两把提取尺（ASCII 类＋带 CJK 类）在三个 SHA（21bec8a1/7828c972/9ec258d9）下 MISSING 集恒等＝枚数对"我跑在哪个漂移点"不敏感。

## §3 与既有账的关系

- `ls .scratch/wisp/issues/ | grep '^212'`＝唯一一枚 `212-comments-cite-evidence-files-that-do-not-exist-done.md`。⚠ **票面里 grep 'fullpack2|GitHub|skipped|本机|rc=0' 零命中＝那两枚"CI 红／本机绿"具名句不在 212 票面里**（0 枚＝没找到，如实报）。这两句的真身在**台账 A608＝`docs/reports/pending-and-issues.md:11883`**（今天 09:02 由 `7828c972` 落账）：逐句①"本机该文件在磁盘上存在（197,213 字节），但不在 git index 里"、②"本机 `sh scripts/d22scan.sh` rc=0 clean"。票面可引的是 `:44`"要消红只能把被引用的那件真做出来，或把那句话改成实话"——正是 #1 两支修法的出处。
- 名册归属：#1 fullpack2 **非这次新出**——212 时代已登记（`probes/212/a3/work/classified.tsv` 3 处、r2/r3 各 log 1-2 处），A608 已具名。#2 esclistener **这次新出**（`0c1b32cf` 246-r1 引入，145 枚之内），但被 testdata 规则赦＝不入库任何账，本件是第一登记者。#3-5 是三枚**仪器自认**的 symRef 例形（main.go 注释逐字点名其中两枚），非账。
- 族产先例（供归口参照，非 212 内）：票 110 面 `:112`"CI runner 上真实红 3 条、本机全绿"、票 128 面"本机绿/CI 红"判定——A608 称 fullpack2 为"同族一枚"有前例链支撑。⛔ 212 一字未改（我也改不了）。

## §4 分母与后果

- **推送树（origin/cnb 的 `21bec8a1` 与新 `9ec258d9` 在此判据下同树同分母）：CI lint 会红 1 枚**，逐字仍是 `cmd/wisp/panel_host_windows.go:343`。**不比昨天多，也不比昨天少**——链头继续常红。
- skip 级联（逐名逐行抄自 ci.yml；链头＝`:74` "D22 scanner positive control"）：无 guard 必被 skipped 的是 **7 枚**——`:83` "D22 scanner self-test (ticket 161 AC#2…)"、`:108` "D22 seven-ban + emoji scan (tools/d22scan)"、`:136` "Tracked path-length budget (ticket 262)"、`:168` "gofmt (gofumpt)"、`:184` "gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)"、`:206` "go vet (module)"、`:209` "go vet (tools/d22scan module)"。带 `if: ${{ !cancelled() }}` 仍会跑：`:213` staticcheck（guard :263）、`:289` mockllm module vet（guard :300）。⚠ **与你派单"五步"及 A608 的名册差一枚**：`:136` path-length budget 也在 skip 链里（我照的是 ci.yml 文本，非 run 日志；run 日志面我无凭据复看，归 §5）。
- "把 fullpack2.log 纳进版本库"这一支的代价：文件 197,213 B，**无 ignore 规则挡**（`git check-ignore -v` rc=1）⇒ 普通 `git add` 即可，不需 `-f`；r6 目录本就有 6 枚已跟踪件（commit-*.txt×4、`.bak`、status-start.txt），**单 add 这一枚不牵出别枚**。对其他门的连带：路径名 41 字符＜262 预算；非 `.go` ⇒ gofmt/gofumpt 分母不动；ban #6/#8 walk 不含 `.scratch` ⇒ 分母不动；钉 `TestRealRepoLedgerIsHonest` 要的正是在树文件 ⇒ 转绿。代价＝197KB blob 永久进 history＋每次 checkout 携带＋票面语境的"证据文件入库"先例从此开闸（要不要开是你的裁，不是我的）。

## §5 我判不动／量不到的地方

- **没跑扫描器**（⛔ go 令遵）——全部判词是对 `main.go:264-269/672-702/831-876/890-955` 的手工复刻。〔推断，未跑扫描器〕；唯 #1 有外部实证背书：06:29 run 37240161874 逐字点名的就是同行同 token。
- #2 的 testdata 赦依只写在 walkGo 源码（`:682`），ban #9 文档头 `:50-69` 未提——**若任何腿收窄 walkGo 射程，它立刻转红**；该规则现无 selftest 样本钉（我没找到钉 testdata 豁免的样本，0 枚＝没找到）。归口：落地腿或后续账。
- `os.Stat` 与 `git ls-tree` 的等价性在符号链接、大小写折叠文件系统上有理论差；本名册 5 枚均不涉。
- fullpack2.log 的 `:640` 读数真伪、r7 evidence md 是否真载同一读数——我没打开 197KB 日志核对；修法 (b) 落地前须由落地腿复读，否则"改成实话"会造出新的假引用。**（⛔ 这一格编排者 09:1x 已替我复：`:640`/`:641` 逐字属实，见抬头复尺①。）**
- skip 级联以 ci.yml 文本为据（7 枚），与 A608 run 日志读数未逐 step 复核（gh 不在我尺面）。**（编排者补：run `37249563077` 的 lint job 步名册里 `:136` 那枚确在 skipped 之列，两边一致。）**
- 并发窗口：porcelain 645→646、HEAD 两次前移；本报告枚数只对 `9ec258d9` 及 `21bec8a1`（两树在此判据下同构）负责，推送前若有新 cmd/internal 注释 commit 落树，名册需重跑一把（成本≈一条 pipeline，我留了形状）。

## §6 交件判语

- 本腿**写文件 0 枚、commit 0 枚**：工具面严格只读，交付即本文（前例 `9ec258d9` 转录 267-a1，同形）。骨架先落的要求我物理不可达，不假报。
- 跑过命令名册：`date / git rev-parse / git branch / git status --porcelain / git grep / git ls-tree / git ls-files / git show / git log -S / git diff --name-only / git check-ignore / grep(-c/-o -P/-E/-F) / sed / wc / sort / head / ls`——**`go` 出现 0 次**（0 次＝没跑＝你的 SLO 争用前提未被我触犯，好消息）。
- AC 框、票面（含 212）：零改动（能力上亦不可能）。`tools/d22scan/**` 只读。
- frontend/design：未读任何内容。⚠ 一格越界自报：SEP3 那把全仓计数尺 `git grep -cFn <token> HEAD` 未带 pathspec，git 内部遍历过全树——产出的只是**文件名计数**，design/frontend 零枚命中行、我未读其任何内容；此后所有搜索尺均改钉 `-- cmd internal` 双根（违反形状的边缘情形，如实登记供你裁）。
- 冻结件：未碰（`docs/reports/pending-and-issues.md` 只 grep 行号读数一处，未读其余）。
