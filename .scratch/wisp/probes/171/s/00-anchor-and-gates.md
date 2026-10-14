# 171s —— 票 171 剩余腿：r3 的「删样本不删腿」钉按票 178 的四堆（piles）口径重推

**结论先行**：票 178 改形（`31b5cd36`，只动 `probes/154/gate-clauses.sh` 一枚文件 ＋91/−25）**只打破了一枚钉**
—— `probes/171/r3/r3-sample-deleted.sh`（旧形在钉处 rc=2，逐字 `no 'want G6' declaration - no verdict`）。
就地把钉重推后 **rc=0**（D2／D1／D1c／D0 四态全绿）；自造反向变异 `SAMPLE_PIN_MUT=noempty` ⇒ **rc=1 并逐字点名 D0**；
还原 ⇒ rc=0。`probes/171/**` 32 枚可执行脚本逐枚现跑：**178 打破者 1 枚**，其余红/绿逐枚有归因（`logs/rc-roster.txt`）。
AC#1／AC#3／AC#5／AC#9 一字未碰；`probes/154/gate-clauses.sh` 零字节；⛔ 未 push、⛔ 未动票面 Status/框。

## 0. 锚与漂移（逐报）
- 起手锚 `1740fd48`；读尺/取数期间 HEAD 被别腿推进：`af983cc4` → `a84fec49` → `eb141e98` → …（每次现读）。
  所有钉都打在各自脚本自带的锚上（4f2a777b／f78d3cb／f1b99a70／19513cce／3df74960／合成树），与 HEAD 漂移无关。
- 供归因的两发现量：**HEAD 上** `sh .scratch/wisp/probes/154/gate-clauses.sh` ⇒ **rc=2**（G1b 基线 0 实测 1＝
  `internal/tools/subagent_197.go:591` 一枚 `.Execute`；G2 基线 1 实测 2＝`cmd/wisp/run.go:986/:992` 两枚真调用行）；
  **4f2a777b 上** ⇒ **rc=1**（只有 G1b——那笔账 **10-03 起就红着**，r3 自己的 `gate-post-at-anchor.txt`／
  `gate-at-lownhead.txt` 逐字记着）。两笔都**不是 178**（178 零生产码、零基线改动；其后无人再动那把尺）。

## 1. 先量全：全脚本 rc 名册（32 枚逐枚现跑；详表＝`logs/rc-roster.txt`）
- **【178 打破】1 枚**：`r3/r3-sample-deleted.sh`（已重钉，两态见 §2 与 `logs/repin-two-states.txt`）。
- **绿 11 枚**：r1-pair-shapes／r1-roster-diff／r3-trailing-comment／r3-roster-diff／重钉后的 r3-sample-deleted／
  v1 六枚（run-rulers／leg-diff／mutation／ruler-mutations／declaration-flips／sample-d）。
- **红 4 枚，逐枚归因（无一是 178）**：
  - `r2/r2-roster-diff.sh` rc=1 与 `v2/r2-roster-diff-copy.sh` rc=1＝**【r3 旧账】**：G7／G7neg 各退场 2 枚
    （`internal/observe/goroutine.go`／`internal/proc/shutdown.go`，都是 AC#1 分子口径改动的既有账，逐枚有 COMMENT-ONLY 解释）。
  - `r2/r2-run-legs-and-subtraction.sh` rc=1＝**【r3 旧账】**（它假设 G2 基线 2／态 1 恰在基线，r3 已核成 1 ⇒ 两态都在基线上方 ⇒
    "FIXED CODE did not react"）＋**【漂移】**（pristine/fixed 两尺都吃 HEAD 的 G1b/G2 ⇒ 要求 rc=0 的两格不再成立）。
  - `r3/flip-prefix-4f2a777b.sh` rc=1＝**【漂移】**：它**无锚钉**（锚在跑的那一刻的 HEAD）；baseline／restore／两枚同翻
    三格要 0 而 HEAD 上 rc=2。五个单向 flip（G2／G5neg／G7／G1／G7pos）全 ok；178 不涉（它的 flip 一枚不碰 G6/G6neg）。
- **门件 3 枚**（r3/doors/door-lie-pre／door-lie-post／door-honest-post）：本由 `r3-flip-columns.sh` 代跑（它是生成方＋调用方）。
  ⚠ 该 harness 本轮跑到 step(2) 被**外部杀**（起 20:22、杀 ~21:11，49 min；step(1) 已完整）：step(2)/(3) 改以
  `door-lie-post.sh`／`door-honest-post.sh` **单跑补到位**（与 harness 当步同字节：door-honest 与"从现 flip 源重新补丁"逐字节
  相同，现跑 cmp 过）。harness 自身**整 rc 因此无读数**（具名，见 §4）。
- 其余 11 枚（gate-prefix-f78d3cb／两枚 fixture／v2 两枚 gate／8 枚 v2 变异尺）：尺副本/夹具/变异件，
  rc＝那版尺打在今天的树上的读数（非钉；非 0 按设计）。
- `a1/commands.sh`：**非可执行**（`-rw-r--r--`、无 shebang、正文是命令转录）⇒ 未跑（具名）。

## 2. 重钉：就地改（不另起 `s/` 后继脚本）＋两态＋牙
- **落点理由**：钉的本体整个换了（翻哪两行副本登记、哪几态、哪条牙）；后继脚本得把合成树机器（~150 行）复制一份，
  且旧件会永远留一枚 rc=2 的破钉在盘上（下一程/CI 一跑就撞）。就地改只留一枚正面台件；旧读数用**新文件名**让位：
  `logs/samples-D*.txt`（10-03 旧口径）一字节未动，新读数＝`logs/samples2$tag-*.txt`。
- **重推**（读现尺：`want G6 quiet piles`＋`want_n 0`；判据只吃【② 只开】堆**文件**枚数）：
  - D2 两枚只开样本 ⇒ ②=2＝基线 ⇒ `ok 基线=2枚 实测=2枚`（表空、不进 BAD 名册）；
  - D1 摘一枚 ⇒ ②=1 ⇒ `SHR 基线=2枚 实测=1枚`＋STALE 表点名 G6、聚合退码不含它那一味；
  - D1c（171s 新加：一枚只开＋一枚**只关**）⇒ ②=1／③=1，实测只吃 ② ⇒ `SHR`、不响退码——把 178 的口径本体钉住；
  - D0 摘光 ⇒ ②=0，腿仍自称 ring ⇒ `BAD 因=空/非空那一味`、聚合退码跟着动。
- **两态＋变异两态原文** ＝ `logs/repin-two-states.txt`（逐字）。旧钉 rc=2 由 178-v1 §7 与本腿各复现一次。
- **牙**：`SAMPLE_PIN_MUT=noempty`（只把副本里『空/非空』那一味挪走）⇒ rc=1，逐字点名
  `FAIL D0: expected a BAD row naming 空/非空 for G6, got '# SHR …实测=0枚…'` 与
  `FAIL D0: the emptied leg never reached the BAD list - AC#8-1's other half is hollow`；还原（不带变量）⇒ rc=0。

## 3. 门禁与改动面
- 提交面（显式 pathspec）：`r3/r3-sample-deleted.sh`（M）＋`r3/fixtures/d_sample_c_close_only.go.txt`（新）＋
  `s/00-anchor-and-gates.md`＋`s/logs/rc-roster.txt`＋`s/logs/repin-two-states.txt`＋票 171 尾 1 bullet。
- `probes/154/gate-clauses.sh` **零字节**（`git status --porcelain -- .scratch/wisp/probes/154` 空；变异只落 $work 副本，
  每发以 cmp 自证"改了一个字节"才继续）。未动 `.github/**`／`scripts/**`／`internal/**`／`cmd/**`／台账。
- **复跑副作用**（未提交、未还原；先例＝171-r1 结案记录）：71 枚 tracked 读数被刷新（带今日时刻；r1/r2/r3/v1/v2 的 logs）
  ＋`r3/fixtures/gate-lie-columns.sh` 被 flip-columns 依当前尺重生成；旧字节仍在 HEAD，取数前快照另存仓外。
- 删除命令 0 次；权限拒绝 0 次；凭据零抄录。

## 4. 具名没做
- AC#1 的正面两发对照台件（把 run＋pair 两族分子一起收）未做（r3-trailing-comment 只侧面覆盖）——按派单押后。
- AC#3（CI 新步＝`.github/**`）／AC#5（契约轴）／AC#9 其余格：未做（按派单押后）；本腿只做"178 打破的那枚钉"的重推。
- `r3-flip-columns.sh` 本体**未跑完**（外部杀在 step(2)）：step(2)/(3) 以两枚 door 单跑补（读数在原样输出里）；
  harness 自身整 rc 无读数。
- `a1/commands.sh` 未跑（非可执行，见 §1）；`r1/roster-diff.sh` 按其用法带 09-27 两枚读数作参数跑（纯报告）。
