# reconcile-8｜只读对账：`docs/reports/pending-and-issues.md` 的 `A751`–`A756` 承重数复跑

腿＝`reconcile-8`（只读）。写面＝本件一枚。⛔ 零产码改动、零票面改动、零台账改动、零 push、零 `go`/`scripts`/`gh` 执行。

## 0. 尺具名（本腿每一张表都用这四把）

- 取数时刻：`2026-10-09 10:24–10:5x +0800`。起手 HEAD＝`7789f953`；腿跑动中 HEAD 漂到 `3263eca4`（`242-v2`／`285` 在飞）。
- 代码面锚定：所有产码行号一律 `git show HEAD:<path>`；`7789f953..3263eca4` 之间 `cmd/`、`internal/`、`.github/`、`docs/PLAN.md` 只新增过 `internal/agent/approval/ticket285_route_denial_name_rulers_test.go` 一枚 ⇒ 本腿量到的产码行号在两枚 HEAD 上都成立。
- 历史面锚定：枚数/差异类一律钉到台账自己点名的 ref（`a99f9a39`/`58a4b2fe`/`c39e2853`/`542d7661`/`2bdf385d`/`039ec93c`），不用漂移中的 HEAD。
- 框尺（全仓统一）：`grep -cE '^[[:space:]]*- \[ \]'`（未勾）与 `grep -cE '^[[:space:]]*- \[x\]'`（已勾），件名走 `git ls-tree -r --name-only <ref> -- .scratch/wisp/issues/ | grep -E "issues/<NN>-"` 全路径取。
- 行尾尺（全仓统一）：工作树 `tr -cd '\r' < <path> | wc -c` vs HEAD blob `git show HEAD:<path> | tr -cd '\r' | wc -c` 两列并排；格式尺 `git show HEAD:<path> | gofmt -d | wc -l` 与 `… | gofumpt -d | wc -l`（stdin，零临时件、零 `-w`）。`gofmt`＝`/d/work/base/go/bin/gofmt`（在 PATH），`gofumpt`＝`/d/work/base/gopath/bin/gofumpt.exe` `v0.12.0 (go1.27.1)`（不在 PATH）；`git config core.autocrlf`＝`true`，`.gitattributes:4`＝`*.go text eol=lf`。
- CI 色：本腿只读 `.github/workflows/ci.yml` 的 HEAD 快照，⛔ 未跑 `gh`、未联网 ⇒ 凡"跑过／红过"只核它附没附 run／job／step 号。

判定用三形：`对上`／`对不上（真值 X，尺写法 Y）`／`复跑不了（缺什么）`。

## 1. 汇总表

| 节 | 本腿复跑的读数 | 对上 | 对不上 | 复跑不了 |
|---|---|---|---|---|
| A751 | 9 | 8 | 1 | 0 |
| A752 | 10 | 8 | 1 | 1 |
| A753 | 12 | 11 | 1（措辞级） | 0 |
| A754 | 20 | 16 | 2 | 2 |
| A755 | 8 | 6 | 1 | 1 |
| A756 | 12 | 11 | 1 | 0 |
| 合计 | 71 | 60 | 7 | 4 |

## 2. 对不上的逐条（按承重排序）

### F1（最要命）A751 仍写着"已勾 14→24"，且票 281 面自称这条已在 A751 更正

- 台账原文（`docs/reports/pending-and-issues.md:14511`）：``AC#4` 复跑＝乙类未勾 **14→4**、已勾 **14→24**``。
- 真值：乙类 8 张（`07/92/97/104/105/110/113/115`）BASE `c39e2853` ＝未勾 **14**／已勾 **26**；HEAD 侧＝未勾 **4**／已勾 **36**。⇒ `已勾 14→24` 两头都错（起点是把未勾列的 14 抄进已勾列，终点少算 12）。A756 写的 `26→36` 才是对的。
- 复跑尺：`for t in 07 92 97 104 105 110 113 115; do git show <ref>:<件> | grep -cE '^[[:space:]]*- \[ \]'; … - \[x\] …; done`，`<ref>` 取 `c39e2853`, `58a4b2fe`, `542d7661`；逐张＝BASE `1/2,3/3,1/3,2/3,1/4,2/3,2/5,2/3`，HEAD `1/2,1/5,0/4,0/5,0/5,0/5,0/7,2/3`。
- 加重的一笔：`git log -L 14509,14515:docs/reports/pending-and-issues.md` 只有 `bf9974c2` 一枚 ⇒ **A751 那节自落下起一字未改**，`sed -n '14509,14515p' | grep -c '更正|不采用'`＝**0**。而票 281 面 `:25` 那句写的是"（已在本节与台账 `A751` 更正，不采用其数）"——**票面这一句对台账的指认现在是假的**（后续程按它读 A751 会读到未更正的 24）。票面本腿不许动，落笔归编排者。

### F2 A755 的"261 长名除本门脚本外只有台账一枚"＝射程少算 14 枚

- 台账原文（`:14551`）：``261` 那个长名还被台账引过一次（`git grep -l` 尺：除本门脚本外只有台账一枚）``。
- 真值：`git grep -l "261-the-model-enabled-key-promises" HEAD`＝**16 枚文件**＝`scripts/check-path-length-budget.sh` ＋ 台账 ＋ **14 枚 `.scratch/wisp/probes/**`**（`212/a2/work/den-md.txt`、`261/{a2,p1,v1}`、`262/v1/logs/{b.txt,m6-with-self-test.txt,mutants-log.txt}`、`done-class-a-1/00-anchor.md`、`done-key-sweep-1/{01,02}`、`expired-premises/{7b/logs/ruler-4f-full.txt,7b/logs/ruler-hits-compact.txt,7b2/logs/scan-sentences.txt}`）。
- 后果：那条"再改一次名会让那处引用漂掉"的定式**方向对、量少报一个数量级**（改名会漂 15 处不是 1 处）。尺未具名排除 `.scratch`，写法与结论面都对不上。

### F3 A754 的"票 284 名长 62 字符"＝66（含 `.md`）／63（不含）

- 台账原文（`:14545`）：`票 284 已立（名长 62 字符，⛔ 不再踩 A755 那条路径雷）`。
- 复跑尺：`git ls-tree -r --name-only <ref> -- .scratch/wisp/issues/ | grep 'issues/284-'` 取名，`printf %s "$name" | wc -c`＝**66**；剥 `.md`＝**63**。三枚 ref（`23924be3`/`7789f953`/HEAD）同名未改。
- 结论面不受伤（两值都 < README `rule 9` 的 100 帽，`issues/README.md:57` 逐字确认该帽），但**这枚数本身就是抄错的**。

### F4 A754 的"residentCard 5 枚命中全在同一文件"＝两枚文件

- 台账原文（`:14540` 末）：`差的形状＝residentCard 在产码里零构造点（5 枚命中全在同一文件）`。
- 真值：产码零构造点这半**对**（`git grep -c "residentCard{" HEAD -- ':(exclude).scratch' ':!*_test.go'`＝空、rc=1）；括号那半**错**——`git grep -n "residentCard{" HEAD -- ':(exclude).scratch'`＝**5 命中分在两枚文件**：`cmd/wisp/resident_approval_246_windows_test.go`（`:97`/`:161`＝2）＋`cmd/wisp/resident_approval_live_246_windows_test.go`（`:115`/`:220`/`:311`＝3）。
- 这一条是本腿按题目要求"枚数写成〈类型全名〉＋〈调用形状〉"复跑出来的：原写法把两枚文件的命中说成一枚，后续程拿它当"集中在一个文件"的判据就会错。

### F5 A752 的"整链尺|回退形无仪器 全仓 1 命中"＝射程写错（数是台账的数）

- 台账原文（`:14523`）：`台账 `整链尺|回退形无仪器` 全仓 **1 命中**＝`A749` 那条缺口登记`。
- 真值：**台账面 1 命中**成立——`git show c39e2853:docs/reports/pending-and-issues.md | grep -nE '整链尺|回退形无仪器'`＝**1 行**＝`:14497`（正是 `A749`）属实。但"全仓"这把尺今天量到 **9 枚文件／22 行**：`git grep -l -E "整链尺|回退形无仪器" bf9974c2`＝票 283 面（5）＋`probes/281/r1/00-anchor`＋`probes/282/v1b/10-verdict`＋`probes/283/r1/00-anchor`＋`probes/283/r1/10-positive-controls`（6）＋`probes/parking-3/01-section`＋`HANDOVER.md`＋台账＋`internal/tools/ticket283_corr_identity_rulers_test.go`（2）。
- 出处对照：腿件 `.scratch/wisp/probes/283/r1/10-positive-controls.md` 的"③ 台账"那一发本来就把尺限定在 `docs/reports/pending-and-issues.md` ⇒ 编排者落账时把"台账"改写成"全仓"，**数的射程被换了**。

### F6 A756 的"07 票面 `:14` 写的落点名 `pending-human-review.md`"＝行号张冠李戴

- 台账原文（`:14577` 末）：`另具名一条不改：07 票面 `:14` 写的落点名 `pending-human-review.md` 盘上不存在，实际落点＝本台账 §待人工审核`。
- 真值：**结论两头都对**（`git ls-files | grep -c "pending-human-review.md"`＝0；台账 `:5`＝`## 待人工审核（pending-human-review）`）。**错的是"07 票面 `:14`"这个出处**：
  - `HEAD:.scratch/wisp/issues/07-ball-state-machine-core-done.md` 的 `:14` 逐字是 `DesignTokens, plus the `statemachine` module implementing the 20-state / 40-transition table`，不含那个文件名；
  - 那个文件名在 `docs/evidence/s1/07-adversarial-acceptance.md:14`（尺＝`git show HEAD:docs/evidence/s1/07-adversarial-acceptance.md | awk 'NR==14'`，逐字含 `已登记 docs/reports/pending-human-review.md`）；
  - 票 07 面**自己记着这条 MINOR 的位置是 `:146`**（`git show HEAD:<07 面> | grep -n "pending-human-review"`＝只回 `:146`，且那句正是在转述"裁决表 line 14 引了个死文件引用"）。
- 尺：`grep -n` 逐枚，见上。⇒ 建议改写为"07 **裁决表** `:14`（票面 `:146` 已把这条记为 dead file reference）"。

### F7 A753 的"照台件 `:40-41`"那句命令原文不是逐字

- 台账原文（`:14529`）：`命令原文（照台件 cmd/wisp/resident_task_source_live_246_windows_test.go:40-41）＝PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -tags winlive ./cmd/wisp -run '…' -count=1 -v`。
- 真值：台件 `:40-41` 逐字是 `//	PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -tags winlive ./cmd/wisp \` ＋ `//	  -run 'TestLive246ResidentPipeline|TestLive246ExitCancels' -v`——**没有 `-count=1`**。编排者实跑那发（`00-own-rerun.txt`）确实带了 `-count=1`，所以读数不受伤；伤的是"照台件逐字"这个说法。尺＝`git show HEAD:<件> | awk 'NR>=38 && NR<=41'`。

## 3. 复跑不了的（具名缺哪一发读数）

| 条目 | 缺的读数 | 本腿为什么拿不到 |
|---|---|---|
| A752 `go test ./internal/tools/ -run 'Test283'`＝rc=0／2 PASS／0.129s | 今天再跑一发的 rc/时长 | 题目禁 `go test`；本腿只核了留证件 `.scratch/wisp/probes/orch-283r1/{00-baseline,01-seedA,02-seedC,03-restored}.txt` 的原文（PASS/FAIL 枚数、时长 `0.129s`、红句行号 `:40`/`:194`/`:206`/`:209` **逐枚对上台账**）⇒ 只有"当时跑过"这一层无法独立复现 |
| A754 死腿 `278-r1`／`280-r1` 的 26／21 次调用、528s／565s | 会话侧调用计数与时长 | 盘上无该读数件；本腿无 `gh`、无联网 ⇒ 只能采〔腿报〕 |
| A754 "链接器里到底留不留 `Gate.Replay`" | `go build` 之后的符号表 | 禁 `go build`；台账自己也写明这一格判不动（与 F 无关，登记为射程边界） |
| A755 `sh scripts/check-path-length-budget.sh --with-self-test`＝rc=0／VERDICT GREEN | 脚本执行输出 | 禁跑 `scripts/*.sh`。本腿改做**静态等价复算**：`git ls-files` 全树 `tracked paths`＝`2bdf385d` 8446／`542d7661` 8451／`7789f953` 8453／HEAD `3263eca4` 8458；`basename>100` 的跟踪路径＝**57**；`R["…"]` 键＝**57**；键里**零枚**已不对应跟踪路径（`git cat-file -e HEAD:<key>` 全过）属实 ⇒ 台账那四个 57 与 "not in roster=0" 在静态面全部站得住 |

## 4. 对上但射程不足的（后续程别当已验）

1. **A754 的 `ci.yml:655` 连锁推论**（"删 `AskOnTaskRoot` ⇒ 三枚用例编不过 ⇒ 那道 `go vet -tags winlive` 立刻红"）：三枚调用点 `resident_approval_live_246_windows_test.go:115/220/311`、该件首行 `//go:build windows && winlive`、`ci.yml:655` 都验到了；**"立刻红"是编译面推断，本程零 run 观测**，不得升格成"CI 会红"的实测句。
2. **A756 的 CI 那格**：台账写"那扇门今天已因 `.scratch/**` 里的故意坏样本常红（票 111 `:310`/`:394`、票 269、票 276 三处既有账）"。三处既有账**都在**（票 111 `:310` 逐字含 step8 `gofmt (gofumpt)` failure ＋ **job id `112053738745`**＝唯一附了 run 凭据的一处；`ci.yml` 里 `:172`/`:173`/`:175`/`:176` 与 `:189`/`:239`/`:240` **逐行对上**属实；`gofumpt -l .scratch` 工作树现量 **39 枚**命中＝"坏样本在"这一条有盘上凭据）。但票 269/276 的名册引行（`:168`/`:171`/`:204`）**今天已漂到** `:172`/`:176`/`:240` ⇒ 引行号前先重跑。⚠ "常红"这个现在时仍只有票 111 那一枚 job id 撑，**本腿不许跑 `gh` ⇒ 未证今天仍红**。
3. **A754 的 `AskOnTaskRoot` 注释"已到位"**：`:689-693` 逐字写了 `today its only callers are this package's own cases`（＝零产码调用者那一半到位）属实；但票 284 的待补形状是三件"零产码调用者／非死代码／勿删"，本腿 `awk 'NR>=680 && NR<=700' | grep -iE 'dead|do not remove|not dead|勿删'`＝**零命中** ⇒ "已到位"只对到三分之一，与 `Gate.Replay` 那枚（`:751-754` 同样零命中）的**差比台账写的小**。尺＝上面那条 awk+grep。
4. **A753 的"CI 编译 winlive 档、但不执行其用例"**：yaml 面（`git show HEAD:.github/workflows/ci.yml | grep -nE winlive` 只有 `:606` 步名、`:608-:626` 注释、`:655` 一条 `go vet`，**零 `go test -tags winlive`**）；那道步**在真实 run 里求值过没有**＝台账自己也没敢写，本腿同样只能采"未观测"。
5. **A751 的"名册改动只落在 6 枚票面＋本腿探针件"**：票面 6 枚（`git diff --name-only a99f9a39 58a4b2fe -- .scratch/wisp/issues/`＝恰 `92/97/104/105/110/113`，`07`/`115` 零出现）；但同一区间还动了 `probes/parking-3/{00-anchor,01-section}.md` 两枚**非本腿**探针件 ⇒ "只落在 6 枚票面＋本腿探针件"这句在 diff 射程内不严（不伤结论，伤措辞）。
6. **A755 编队现状两格（`:14553`／`:14578` 在飞枚数、`porcelain` 空、`242-v2` 自报 90 PASS）**：台账自己已标"共享树读数会漂，引前重跑"；本腿实测 HEAD 在 25 分钟内从 `7789f953` 漂到 `3263eca4`（新增 `internal/agent/approval/ticket285_route_denial_name_rulers_test.go`、票 `282` 面**掉了 `-done` 后缀**、`-done` 名册 98→**103**）⇒ 这类句子的有效期比一节台账短，不判对错。

## 5. 对上清单（逐条给尺，供编排者背书用）

| # | 台账句（行） | 真值 | 尺（全部 `git show`／`git grep` HEAD 快照，除非注明 ref） |
|---|---|---|---|
| 1 | A751 翻框 10 对 10、逐对同文、零反向（`:14512`①） | 10／10，剥 `[x]→[ ]` 后两份名单 **set-equal**；反向 0/0 | `git diff a99f9a39 58a4b2fe -- .scratch/wisp/issues/` 的 `^-…- \[ \]` vs `^+…- \[x\]` |
| 2 | A751 AC 分布 3×2／4×2／5×4／6×1／7×1 | 逐枚＝104`AC#3,AC#4`、97`AC#5`、110`AC#3,AC#5`、105`AC#5`、113`AC#4,AC#5`、92`AC#6,AC#7` ⇒ 恰该分布 | 逐件 `git diff … -- <件> | grep -E '^\+.*- \[x\]'` 取行首标号（⚠ 不得 `grep -oE 'AC#[0-9]+'` 全文抓，会把句内引用的 `AC#1` 计进来） |
| 3 | A751 新增 Status 行恰 2 条＋逐票 BASE↔HEAD | 全 diff `^\+…\*\*Status`＝**2**；`92/97/104` `1→1`、`105/110` `1→2`、`113` BASE 已是 2 | `grep -cE '\*\*Status'`，ref 取 `a99f9a39`／`58a4b2fe` |
| 4 | A751 `105:4`/`110:4` 留档标记、现行值 `:3`、留档原句 `:5` | 两件的 `:4` 都是"（下箭头符号＋自下一行起整块＝**原句逐字留档**…"那枚横幅行（箭头条形符此处不复制，语义无差）、`:3` 新 Status、`:5` 旧 Status | `git show 58a4b2fe:<件> | grep -nE '\*\*Status'` |
| 5 | A751 `-done` 名册仍 98 | `a99f9a39`＝98、`58a4b2fe`＝98（今日 HEAD 已 103） | `git ls-tree -r --name-only <ref> -- .scratch/wisp/issues/ | grep -c -- '-done\.md$'` |
| 6 | A752 `039ec93c`＝1 file/227 insert/零删除 | 逐字对上，件行总数亦 227 | `git show --stat --format='' 039ec93c` |
| 7 | A752 种点 `cancel.go:69`/`task.go:708` | `:69`＝`return h.taskID`、`:708`＝`caller := TaskID(ctx)`（同一次取数拿到） | `git show HEAD:<file> | grep -n '<原句>'` |
| 8 | A752 三枚 hash 回基线 | `cancel 87ac0624`／`task 3ec8d498`／`subagent_197 8c9266a7` 三枚在 HEAD **逐枚等值** | `git rev-parse HEAD:<path> | cut -c1-8` |
| 9 | A752 尺件行号 `:40`/`:194`/`:206`/`:209`/`:218-226` | 逐行内容与用途相符（`:40` 单元尺 Fatalf、`:194`/`:206` 整链尺、`:209` 身份拒、`:218-226` 按 `CorrelationID` 形状拦）；seedA 留证件红句集＝{`:40`,`:194`,`:206`}、seedC＝{`:206`,`:209`} 且单元尺 PASS／整链尺 FAIL ⇒ **台账两发的"三处红/两处红"逐枚对得上** | `git show HEAD:<件> | awk 'NR==…'` ＋ `grep -oE 'ticket283_…test.go:[0-9]+' probes/orch-283r1/0{1,2}-*.txt | sort -u` |
| 10 | A752 `归口` README 零命中 | 0、rc=1 | `grep -c "归口" .scratch/wisp/issues/README.md` |
| 11 | A752 `TaskID(ctx)` 同名尺只回新件 | 只回 `internal/tools/ticket283_corr_identity_rulers_test.go` | `git grep -l "TaskID(ctx)" HEAD -- 'internal/tools/*_test.go'` |
| 12 | A753 卡行 `resident_approval_windows.go:886`、审计行 `:881`、`:885`＝channels attr | `:886`＝`fmt.Printf("wisp: 卡片挂起：…")`、`:881`＝`slog.Info("approval: 常驻进程显示一张确认卡片",`、`:885`＝`"channels", channelRosterText(p.Channels))` ⇒ `:885→:886` 的"漂 +1"是真的 | `git show HEAD:<件> | awk 'NR>=876 && NR<=892'` |
| 13 | A753 `ci.yml:606`/`:655` winlive | 逐字对上（`:606` 步名含 `ticket 111 AC#11`、`:655`＝`run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/`）；落档 `351e5a5e`＝2026-10-08、只碰 `ci.yml`、标题含 `票 111 r6 AC#11` | `git show HEAD:.github/workflows/ci.yml | grep -nE winlive`；`git show --name-only 351e5a5e` |
| 14 | A753 `runConsoleLoop` 台件零命中 | rc=1 空；正控＝不带 `:*_test.go` 时同符号在生产件 `resident_task_source_windows.go:348/374` 有命中 ⇒ 尺不空转 | `git grep -c "runConsoleLoop" HEAD -- '*_test.go'` |
| 15 | A753 自己那发读数（`rc=0`／`PASS` 2／`RUN` 2／`FAIL·SKIP` 0／17.48s／25.49s／corr `900bb089…`） | 留证件逐枚对上；与腿那发 corr `0455ee9f-5fc7-…` 不同 ⇒ 非缓存回放；腿件另具名剔掉 run2＝`ok … (cached)`（同 corr／同 `10.95s`/`21.60s`／同临时目录名）属实 | `grep -c -- '--- PASS'`／`'=== RUN'`／`grep -c -- FAIL`＝0／`grep -c -- SKIP`＝0；`01-a-道-readings.md:28-30` |
| 16 | A753 配方 §4 七句红/判不了 0 命中 | 七句在 `00-own-rerun.txt` **逐句 0 命中**（`never raised a card`／`任务入口未启用`／`任务管线未装配`／`不被受理`／`卡片无处呈现`／`somebody else owns`／`observer rig is blind`） | `grep -cF '<句>'` 逐枚 |
| 17 | A753 B 道三件（无 `config.toml`／无 secrets、exe 旧 31h23m、`%APPDATA%\wisp` 只有 `logs/`） | `ls` 现量＝`Roaming/wisp` 仅 `logs`（mtime 10-07 09:33）属实；`build/wisp.exe` mtime `2026-10-07 11:57:46`；`cmd/wisp` 最末提交 `aa6c1881 2026-10-08 19:20:48` ⇒ 差 **31h23m02s**（台账的 31h23m 对） | `ls -l --time-style=full-iso`；`git log -1 -- cmd/wisp/` |
| 18 | A753/A754 `Gate.Replay` 那组行号 | 定义 `gate.go:755`（旧抄 `:749` ⇒ **漂 +6** 与台账一致）；注释块 `:751-754` 逐字含 `C18 一键重放`、`NOT an answer`；台账 `:11139` **确实**是那句"尺只命中定义 `gate.go:749`"（本腿逐行读出该子句，指认无误）；`gate.go` 零 build tag | `git show HEAD:<件> | awk 'NR>=748 && NR<=760'`；`sed -n '11139p' | grep -oE '.{0,90}(Replay|:749).{0,90}'` |
| 19 | A754 入口结构性不存在 | `NativeAPI`＝`ui.go:143-162`（件在 `internal/agent/approval/`，台账只写 `ui.go` 未写目录）、`PanelAPI`＝`:167-171`，两面值内 `Replay` 零出现；动词条 `cmd/wisp/approval_reply.go:566-584`＝yes/session/no/veto/panel-no/panel-yes/always/head/view/help，`replay` 零出现 | `git grep -n "type (NativeAPI|PanelAPI)"` ＋ 逐行 awk ＋ `grep -niE '\breplay\b'` |
| 20 | A754 `PLAN.md:1368` C18 行 | 该行**是** C18 契约行且含"拒绝后**任务 root ctx 不取消，可一键重放**"；台账引文剥掉了原文的两个 `**` 加粗标记（严格"逐字"差两个记号，语义无误） | `git show HEAD:docs/PLAN.md | awk 'NR==1368'` |
| 21 | A754 两把承重负向尺 | `\.AskOnTaskRoot\(` ＝**rc=1 空**；`\.Replay\(` ＝**rc=1 空**。正控：剥掉 `:!*_test.go` 后前者回 3 枚（那三枚 live 用例）、后者回 `queue_test.go:283` ⇒ 尺不空转；再放宽到 `\bReplay\(` 全形，非测试只剩定义 `gate.go:755`（`g.q.replay` 是小写队内法，另一枚符号）⇒ **"零产码调用者"两句今天仍成立** | `git grep -nE "[A-Za-z0-9_]\.AskOnTaskRoot\(" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'`（`Replay` 同形）＋本腿两条放宽尺 |
| 22 | A754 那句歧义注释 | 真身 `resident_windows.go:249`（`This call is the caller: it reads`）属实、`:248`＝过去式 `still had zero product`、`:260`＝`src := startResidentTaskSource(rt, ra)` ⇒ 漂 +1 与"甲读法不成立/乙读法成立"的两段写法都被产码支持；boot 行 `:269-270` 确含 `任务来源：%s` | `git show HEAD:cmd/wisp/resident_windows.go | awk 'NR>=245 && NR<=275'` |
| 23 | A754 stale-claim 那枚腐烂率 | `probes/stale-claim-1/00-ruler-and-inventory.md:29`＋`04-unjudgeable-and-blind.md:47` 逐字＝产码状态级断言 **39 枚 ⇒ 已过期 15／判不动 9／仍成立 15** 与台账同数 | `grep -noE '.{0,70}(39 枚|15 枚).{0,70}' .scratch/wisp/probes/stale-claim-1/` |
| 24 | A754 278 两件正文 158／90 行 | HEAD 快照 `10-ac1-two-shapes.md`＝**158**、`20-ac2-ac3.md`＝**90**（两件均已跟踪，`2bdf385d` 代提属实） | `git show HEAD:<件> | wc -l` |
| 25 | A755 仓根 `probes/` | `git ls-files probes`＝**空**；`8ccc9568` 逐字＝`A probes/280/r1/00-anchor.md`（仓根）属实；现 HEAD 无该路径。⚠ 工作树仍有**未跟踪空目录** `probes/280/r1/`（"只建不删"下正常，但台账那句只承诺"无跟踪件"，别读成"目录已消失"） | `git ls-files probes`；`git show --name-status 8ccc9568`；`ls -R probes` |
| 26 | A755 两桩路径长度病 | 跟踪名 `258-…-balldebug-done.md` 相对长 **165**、`261-…-map-entries-done.md` **184** 且 184＝**全树最长**（次长 180/178）属实；`2bdf385d` 上跟踪件已是 `-done` 名而名册键是**老名**（`…ldebug.md`/`…ntries.md`）⇒ 病 A＋病 B 的形状复现；作废的短名 `258-hotkey-tier-vs-resident-ball-defaults` 现**零跟踪**（确已改回）属实 | `git ls-files | awk '{print length, $0}' | sort -rn`；`git show <ref>:scripts/check-path-length-budget.sh | grep -oE 'R\["…(258|261)[^"]*"'` |
| 27 | A755 同族惯例（153/152/151/184/250/155） | 六枚都在名册且都是"保留原名＋`-done`"形状；理由句 `closed, name doubles as the -done anti-double-claim key` 在脚本内出现 **20 次**（逐字）属实；README `rule 9`＝100 字符帽（`issues/README.md:57`） | `grep -oE 'R\["[^"]+"\]' scripts/check-path-length-budget.sh | grep -E '/(153|152|151|184|250|155)-'` |
| 28 | A756 那张 7 枚格式表（全列） | `gofmt -d` 行数 0/0/0/0/0/**14**/**14**、`gofumpt -d` 行数 0/0/0/0/0/**14**/**101**、工作树 CR **334/131/225/1115/1284/0/0**、HEAD blob CR **全 0** ⇒ 七行**逐格对上**（stdin 取 HEAD 快照，零临时件） | `git show HEAD:<p> | gofmt -d | wc -l`／`… | gofumpt.exe -d | wc -l`；`tr -cd '\r' < <p> | wc -c` vs `git show HEAD:<p> | tr -cd '\r' | wc -c` |
| 29 | A756 "票 280 前提只成立 2 枚"＋工作树射程 | 工作树 `gofmt -l cmd internal`＝**恰那 7 枚、逐名对上票 280 面 `:10` 的名册**（含 5 枚纯 CRLF artifact）属实；`gofumpt -l cmd internal`＝**7** ⇒ "两把尺吃工作树字节全点红、吃 HEAD blob 只剩 2 枚"这一整套说法复跑成立 | `/d/work/base/go/bin/gofmt -l cmd internal`；`gofumpt.exe -l cmd internal`；票面名册 `sed -n '10p' .scratch/wisp/issues/280-*.md` |
| 30 | A756 "同 hunk `@@ -63,8`＝纯 gofmt 风格" | `inbound_guards_35r3_test.go`：gofmt 与 gofumpt **各自唯一 hunk 都是 `@@ -63,8 +63,8 @@`**（同 hunk 同声，14 行）；`transport_35r2_test.go`：gofmt 1 hunk `@@ -804,8`（14 行）vs gofumpt **5 hunks／101 行**（`@@ -400,18 @@ -421,12 @@ -698,8 @@ -804,16 @@ -821,6`）属实 ⇒ "多出 gofumpt 独有"这一档也对 | `git show HEAD:<p> | gofmt -d | grep -oE '^@@ [^@]+ @@'`（gofumpt 同） |
| 31 | A756 `ci.yml` 那 7 个行号 | `:172` 步名 `gofmt (gofumpt)`／`:173` `if: ${{ !cancelled() }}`／`:175` `go install mvdan.cc/gofumpt@latest`／`:176` `gofumpt -l . tools/d22scan tools/mockllm`／`:189` 步名（tracked 分母）／`:239` `if: ${{ !cancelled() }}`／`:240` `run: sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` ⇒ **七枚逐字对上**，且 scope 确为仓根全树（含 `.scratch/**`、含 `_test.go`）属实 | `git show HEAD:.github/workflows/ci.yml | awk 'NR==…'` |
| 32 | A756 未勾 14→4／已勾 26→36／净＋10 | BASE `c39e2853`＝14/26、HEAD 侧＝4/36 与"补勾 10、撤名 0"自洽（`git diff -M --summary a99f9a39 58a4b2fe` **零 rename 零 delete**） | 框尺（见 §0） |
| 33 | A756 六张 `--numstat` | `92 4/2`、`97 2/1`、`104 4/2`、`105 4/1`、`110 6/2`、`113 4/2` **逐枚对上**；且删除列剥掉"框行/Status 行"后**零枚既有判据句被删**（逐件 `grep -E '^-[^-]' | grep -vE '\- \[[ x]\]|\*\*Status'`＝空） | `git diff --numstat a99f9a39 58a4b2fe -- <件>`；删除行分类尺见左 |
| 34 | A756 "97 那格真值 3、三 ref 一致" | `c39e2853`／`226d3717^`／`6d22dd6c` 三枚 ref 上锚定尺＝**3**、松散 `[x]`＝**3**、未勾＝1 ⇒ 逐枚一致；错数来路 `done-class-b-1/13-summary.md:14` **确实写着 `228 / 1 / 4`**（该行就是 97 那行）属实 | 框尺＋`sed -n '14p' .scratch/wisp/probes/done-class-b-1/13-summary.md` |
| 35 | A756 "四格/三格全勾并 `-done`" | `283`/`278`/`280` 各 `0 未勾/3 已勾`、`281` `0/4`，四件名册都带 `-done`；`283` 面 `Status: done` | 框尺＋`ls .scratch/wisp/issues/` |
| 36 | A756 "一字节格式改动都没落地" | `git diff --name-only c39e2853 542d7661 -- cmd internal`＝**只回 `internal/tools/ticket283_corr_identity_rulers_test.go` 一枚**（新尺件，非格式改动）属实 | 该条 diff |
| 37 | A756 残余 4 枚未勾（07 待人／92 AC#5／115 两枚） | 今日 HEAD：`07 un=1`、`92 un=1`、`115 un=2` ⇒ 合计 **4**；`111/112/140` 三件**仍无 `-done`** | 框尺＋`git ls-tree … | grep -c -- '-done\.md'` |

## 6. 一句话结论

六节里本腿复跑 **71 处**可独立复现的读数：**对上 60**、**对不上 7**（F1 最要命＝A751 的"已勾 14→24"未更正、且票 281 面自称已在 A751 更正；F2＝`.scratch` 射程少算 14 枚引用；F3 名长 62 真值 66/63；F4 "5 枚命中全在同一文件"真值 2 枚文件；F5 "全仓 1 命中"实为"台账 1 命中"；F6 "07 票面 :14"张冠李戴（真出处＝裁决表 `:14`／票面 `:146`）；F7 "照台件 :40-41"的原文缺 `-count=1`）、**复跑不了 4**（`go test` rc／死腿调用计数／符号表／门禁脚本执行面，均已给静态等价复算并逐条具名缺口）。A756 那张 7 枚格式表与 `ci.yml` 七枚行号是本腿复跑里最硬的一块，七格全对。
