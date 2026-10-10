# 301-r1 — 交回：双锚／越界自查／门禁 rc／两枚纪律缺陷（具名自报）

## 1. 双锚（新规矩，两把读数都留）

| | 起手（`00-anchor.md` 同批） | 交回（本文件） |
|---|---|---|
| `git rev-parse --short HEAD` | **`64ceaacd`** | **`42617eba`**（＝本腿第 3 笔） |
| `git branch --show-current` | `dev` | `dev` |
| `git status --porcelain \| wc -l`（整仓） | **810** | **815**（+5，⛔ 是本腿造的，见 §1.1） |
| `git status --porcelain -- scripts internal .github` | **0** | **0** |
| `git status --porcelain -- .scratch/wisp/probes/301` | 0 | **0**（本腿件全部已提交） |
| `date` | 2026-10-10 14:26:28 +0800 | 2026-10-10 14:51:25 +0800 |

★**本文件所属那笔＝`32e8aa29`（第 4 笔）**，它⛔ 在 §1 那两列里——上表「交回」列的 `42617eba` 是本文件落笔**之前**现量的 HEAD（一笔自引用⛔ 存在，⛔ amend）。
⇒ **真·交回锚＝`32e8aa29`**，四笔名册＝`74eb032c`(锚) → `0a0f62ef`(AC#1 两枚 token 同一笔) → `42617eba`(AC#2 名册 34 枚) → `32e8aa29`(本文件)。
交回后再量的两枚尺记在下面 §1.2。

### 1.2 真·交回读数（`32e8aa29` 落笔后现量）

- `git rev-parse --short HEAD` ＝ **`32e8aa29`**
- `git status --porcelain -- .scratch/wisp/probes/301` ＝ **0**（本腿件全落仓，无遗漏）
- `git status --porcelain -- scripts internal .github` ＝ **0**（写面干净：改动已提交，工作树⛔ 残留）
- `git log --format='%h' 64ceaacd..HEAD | wc -l` ＝ **6**：本腿 **4** 枚（`74eb032c`/`0a0f62ef`/`42617eba`/`32e8aa29`）
  ＋编排者 **2** 枚（`98e0b81f`/`6c2dd788`）⇒ **HEAD 在两锚之间漂过别人 2 枚 commit**，写面交集＝**0 枚**（尺＝§1.1 那一把，
  输出只 `scripts/portable-tests.sh` 一枚＝本腿自己的改动）。
- `git diff --stat 64ceaacd..HEAD -- scripts/portable-tests.sh` ＝ **`1 file changed, 2 insertions(+), 1 deletion(-)`**
  ⇒ 三笔探针 commit 对 `scripts/` 零字节，本票在 `scripts` 上的全部净效果＝那两枚 token，可一枚尺读完。

### 1.1 两锚之间 HEAD 漂过别人的 commit —— 具名列出并证名册⛔ 含别人的东西

`git log --format='%h %an %ad' --date=format:'%H:%M' 64ceaacd..HEAD` 名册（5 枚）＝
`74eb032c`(本腿) → `98e0b81f`(CarlosShao＝编排者, 14:3x) → `0a0f62ef`(本腿) → `6c2dd788`(CarlosShao, 14:4x) → `42617eba`(本腿)。
⇒ **漂过 2 枚别人的 commit**（`98e0b81f`／`6c2dd788`）。

尺＝`git diff --name-only 64ceaacd..HEAD -- scripts internal .github` ⇒ 输出**只有 `scripts/portable-tests.sh` 一枚**（＝本腿那枚改动）。
⇒ 那两枚别人的 commit 在本腿写面（`scripts`／`internal`／`.github`）上**零字节**，
本腿读过的每一枚 `scripts`／`internal` blob 与我量过时同形；`internal/audio` 更硬的证＝
`git diff --stat 64ceaacd..98e0b81f -- internal/audio` 输出**空**（0 行）。
⇒ 名册里⛔ 混进别人的东西（除 §3.1 那枚落点事故，那枚是**我的**缺陷，⛔ 是读数污染）。
整仓 porcelain 810→815 那 +5 行＝别人的 `.scratch`／仓根散件，射程⛔ 触及本腿写面（`probes/301` 专用尺＝0）。

## 2. 三句尺口径（逐字，本腿每个数都挂这三句）

1. **顶层名册尺 R1**＝`sed -n 's/^--- \(PASS\|FAIL\|SKIP\): \([^ ]*\).*/\2/p'` ＋ `sort -u`，
   判据＝`--- ` 落在**第 0 列**；Go 的子测试结果线缩进四格 ⇒ **⛔ 含子测试**；日志里
   `portable-tests.sh: `／`runtests.sh: ` 前缀线与 `    ` 缩进的 t.Log 行都⛔ 匹配 ⇒ **⛔ 含注释行**。
2. **含子测试名册尺 R2**＝同串但前置改 `^[[:space:]]*`，即
   `sed -n 's/^[[:space:]]*--- \(PASS\|FAIL\|SKIP\): \([^ ]*\).*/\2/p'` ＋ `sort -u`，子测试保留 `Parent/sub`。
   ★R1／R2 **两把都交，⛔ 混成一枚数**（windows：R1 41 枚／R2 47 枚）。
3. **包内名册尺 R4**（⛔ 编译面）＝对 `internal/audio/*_test.go` 用
   `sed -n 's/^func \(Test[A-Za-z0-9_]*\).*/\1/p'`；tagged 判据＝`git show HEAD:<path> | head -3 | grep -c 'go:build windows'`
   （逐字模式串 `go:build windows`），射程＝**HEAD blob**，⛔ 工作树。
   每枚数在本报告里都写「哪把尺＋射程目录＋blob 还是工作树」，词面尺逐字抄模式串。

## 3. `git show --stat`／`--name-only` 名册逐名（三笔）

| 笔 | 号 | 名册 | 越界 |
|---|---|---|---|
| 1 | `74eb032c` | `00-anchor.md`（1 枚） | 0 |
| 2 | `0a0f62ef` | **28 枚**＝`scripts/portable-tests.sh` ＋ `probes/301/r1/**` 26 枚 ＋ `probes/301/orch/2026-10-10-baseline-and-delta-prediction.md` 1 枚 | **1 枚＝见 §3.1** |
| 3 | `42617eba` | **34 枚**，枚枚在 `probes/301/r1/**`（`20-ac2-differential.md` ＋ `logs/` 17 ＋ `rosters/` 16） | 0 |

★`AC#1` 那笔（第 2 笔）＝**两枚 token 改动在同一笔里**（档清单 `./internal/audio/` ＋ `win_pin` 补导入路径），
`git diff --stat` 逐字＝`scripts/portable-tests.sh | 3 ++-`，`1 file changed, 2 insertions(+), 1 deletion(-)`。

### 3.1 ★纪律缺陷甲（**记我，⛔ 记编排者）＝pathspec-less commit 吞了别人的 staged 件

编排者在 `6c2dd788` 里把这件事写成「补位件写在 `probes/301/orch/` 被在飞腿的 pathspec `probes/301/**` 一起提交」。**这一句的机制读反了，本腿顶回**：

- 本腿第 2 笔的命令形＝`git commit -q -m "$(cat <<'EOF' … EOF)"`，**尾部⛔ 带 `-- <pathspec>`**；
  而 `git add` 那一步用的是**逐枚显式路径**（27 枚，⛔ 含 `orch/` 那枚，尺＝当刻 `git diff --cached --name-only` 数得 27 枚）。
- ⇒ 吞进去的成因⛔ 是「wildcard pathspec 太宽」，**是 `git commit` 缺省提交整个 index**：
  编排者在我 `git add` 与 `git commit` 之间把 `orch/` 那枚 stage 进**共享 index**，我的无-pathspec commit 就把它一起提了。
- ⇒ 违反的是派单逐字那条「**commit 必带显式 pathspec**」，责任在**本腿**，⛔ 在编排者的落点选择。
  定式更正＝`git add` 显式路径 **+** `git commit -- <同一批显式路径>`（两枚都要，一枚⛔ 够）；
  以及提交前跑一发 `git diff --cached --name-only | grep -v '^<我的面>'` 应为空（本腿第 3 笔前跑了这一发，得 0 枚）。
- **无损证明**：那枚件的内容在我 commit 里的 blob 与工作树 `cmp` ⇒ 退码 0＝逐字节相同（§现量），
  编排者那枚件**已进仓**，只是署在本腿的 commit 号下。历史⛔ 改写（⛔ `--amend`／`rebase`／`reset`），按规矩**追加**更正＝本节。
- 附带一条：`AC#4` 那句名册「只含 `scripts/portable-tests.sh` ＋ `probes/301/**`」按字面**仍成立**
  （那枚件确在 `.scratch/wisp/probes/301/**` 之下）⇒ 本腿⛔ 拿「字面过关」搪塞，缺陷照记。

### 3.2 ★纪律缺陷乙（记我）＝第 3 笔的 commit message 被我自己的 shell 形吃掉了

第 3 笔 `42617eba` 的 message **不是**我写的提交说明，而是 `20-ac2-differential.md` 的全文（137 行，首行
`# 301-r1 — AC#2 让它真的开口…`）。成因＝我把 `-- <pathspec>` 写进了 `$(cat <<'EOF' … )` **括号内部**，
于是那三枚路径成了 **`cat` 的参数**而⛔ 是 `git commit` 的：`cat` 收到文件参数就⛔ 读 stdin（heredoc 被丢），
改读那枚 `.md` 并把两枚目录报 `Is a directory`（stdout 里可见）。⇒ **本意那段英文提交说明此刻盘上⛔ 存在**（尺＝
`git log -1 --format=%B 42617eba | grep -c 'Denominator: the ticket sizes'` ＝ **0**），本节 §4 是它唯一的复活落点。
影响面＝**只伤 message，⛔ 伤内容**：名册 34 枚逐枚在 `probes/301/r1/**`、`.md` 136 行工作树==blob 都对。
⛔ `--amend` 修（派单逐字禁），改法＝**追加**：本文件 §4 逐字收录那条丢失的 message，第 4 笔带正经 message。
定式入册＝**heredoc 提交必须把 `--` pathspec 放在 `$( … )` 之外**（`git commit -m "$(cat <<'EOF' … EOF)" -- <paths>`）。

## 4. 第 3 笔本该带的提交说明（逐字复活，⛔ 被 shell 吃掉的那段）

> probe(301-r1): AC#2 differential readout - audio's 42 top-level tests measured before/after in --scope=windows
>
> The contrast AC#2 asks for, measured by actually running the tier rather than by `go test` on this machine:
> `windows-before.txt` mentions `internal/audio` zero times in both stdout and stderr, and after the change
> `comm -13` on the top-level ruler returns 41 names, `comm -23` returns 0, `comm -3` returns 41 (all one-sided).
> Those 41 are a subset of `internal/audio`'s 42 top-level names with exactly one absence — `TestLiveWasapiSmoke`,
> the ledger fixture row — so the blast radius is that package and nothing else.
>
> Denominator: the ticket sizes the cost at 10 windows-tagged tests, but the tier unit is the package, so what
> actually enters the denominator is all 42 (10 tagged + 32 untagged).
>
> SKIP branch measured empirically: `--- SKIP = 0` in the after run, and `TestLiveWasapiSmoke` appears 3 times,
> none of them a result line (ledger banner, skip-pattern echo, runtests.sh package echo). So `-skip` filters
> silently and does not mint a SKIP, which means landing audio cannot redden `test-windows` through the
> runtests.sh skip rule. The `:593` row was left byte-identical.
>
> Named new reds = 0 on both scopes. census rc=0 before and after with the totals line byte-identical
> (`unclaimed-with-tests=0`), and its whole diff is one roster row: `internal/audio` moves from `core` to `core windows`.
>
> Gates: `sh scripts/d22scan.sh` rc=0, `bash -n` rc=0, `bash scripts/portable-tests-selftest.sh` rc=0
> (32 cases, 0 assertions failed, script not modified).

## 5. 越界自查（每枚一条尺，全 0）

| 靶 | 尺 | 结果 |
|---|---|---|
| `frontend/**` | `git diff --name-only 64ceaacd..0a0f62ef -- frontend/ \| wc -l` | **0** |
| `design/**` | 同上 `-- design/` | **0** |
| `internal/**`（含 `internal/audio/**`） | 同上 `-- internal/` | **0** |
| `.github/workflows/ci.yml` | 同上 `-- .github/` | **0** |
| 三枚冻结件（`panel/tokens_fourway_test.go`／`panel/l2_grant_boundary_test.go`／`perm/ticket90_persist_test.go`） | 逐枚 `-- <path>` | **0／0／0** |
| `internal/observe/thresholds.go`（SLO 阈值） | `-- <path>` | **0** |
| `tools/d22scan/allowlist.txt` | `-- <path>` | **0** |
| `.gitattributes`／任何 git 配置 | `-- .gitattributes` ＋ `git config` 零调用 | **0** |
| golden／D43 转移表／C1–C32／D1–D47（落点＝`docs/**`） | `git diff --name-only 64ceaacd..0a0f62ef -- docs/ \| wc -l` | **0** |
| `scripts/portable-tests.sh` 里**除那两枚 token 外**的每一行 | `git diff --stat`＝`2 insertions(+), 1 deletion(-)` 且 `git show` 全文只两行带号 | 只 2 行 |
| `:593`（现 `:594`）ledger 行 | `git show HEAD:…\|grep TestLiveWasapiSmoke` vs 工作树同 grep，`cmp` | **退码 0＝逐字节相同** |
| `portable-tests-selftest.sh`／`d22scan`／`runtests.sh` | 名册（§3）⛔ 含它们 | **0** |
| 放宽任何断言 | §3 名册只两枚 token ＋ `docs/` 0 ＋ `internal/` 0 | **0 枚** |

⛔ 未 push（三笔＋本笔全本地；`git push` 零调用）。临时件只建不删：本腿建的所有件都在
`.scratch/wisp/probes/301/r1/**`（已提交）；`portable-tests-selftest.sh` 自己的 carrier scratch
留在 `/tmp/tmp.8XmdckHnIR`（它 stdout 逐字声明「temp files are created, not deleted」），本腿⛔ 删过任何东西。
`git reset`／`--amend`／`rebase`／`stash`／`checkout .`／`clean`／`add -A`／`add .` **零调用**
（撤销 stage 的正解＝⛔ `reset`，本腿压根没 stage 过要撤的东西；查换行属性用的是
`git check-attr text eol -- scripts/portable-tests.sh` ⇒ `text: set`／`eol: lf`）。

## 6. 门禁 rc 逐行（⛔ 0 字节＝那格没交，具名）

```
windows-before.rc.txt = rc=1   (151,922/236-byte log pair, 5-byte rc)
windows-after.rc.txt  = rc=1
core-before.rc.txt    = rc=1
core-after.rc.txt     = rc=1
census-before.rc.txt  = rc=0
census-after.rc.txt   = rc=0
d22scan.rc.txt        = rc=0   (sh scripts/d22scan.sh)
bashn-after.rc.txt    = rc=0   (bash -n scripts/portable-tests.sh)
selftest.rc.txt       = rc=0   (bash scripts/portable-tests-selftest.sh, 32 cases / 0 failed)
```
0 字节的件具名＝`logs/census-before.err`／`logs/census-after.err`（各 0 字节）＝census 那两发⛔ 有 stderr
（GUARD D⛔ 咬、STALE 腿 0 行）；`logs/d22scan.err`／`logs/selftest.err` 0 字节同理。
`rosters/census-*.{A,B}-*.txt` 5 枚 0 字节＝那一发⛔ 跑 `go test`，故无用例名册可抽
（尺＝`grep 'four numbers' logs/census-*.txt` ＝无命中）⇒ ⛔ 是「没交」，是那格射程本来如此。

## 7. 本腿⛔ 同意的派单／票面预设（必答，逐条；先例＝票 300／301 的腿各顶回过 3–4 处）

1. **§3.1 那处机制顶回**（＝⛔ 同意编排者 `6c2dd788` 的归因句「被在飞腿的 pathspec `probes/301/**` 一起提交」）：
   盘上＝本腿第 2 笔**无 pathspec**，`git add` 是 27 枚显式路径；尺＝当刻 `git diff --cached --name-only` 数与
   `git show --name-only 0a0f62ef` 名册 28 枚作差。责任在本腿，更正见 §3.1 定式。
2. **分母 10 vs 42**：派单逐字「audio 那一族应是 10 枚顶层用例——⚠ 分母按 10⛔ 按 9」。
   本腿实测**新增求值 41 枚＝42−1**：分档单位是**包**，进分母的是 audio 全部 42 枚顶层用例（10 tagged＋32 无 tag，尺 R4）。
   10 那枚数只回答「tagged 用例」这一问，⛔ 回答「这次进分母几枚」。★本腿⛔ 改票面一字，只做具名＋本文件 §1.1 落笔。
   （编排者 14:3x 独立件已自抓到 42 并预测 `A-evaluated 399→440`＝**本腿真跑逐字对上**，见 §8。）
3. **⛔ 跑对照前那几发之前先 `tasklist` 现量**：照做（起手 14:26 `wisp.exe`＝0／`balldebug.exe`＝0，记于 `00-anchor.md` §6）。
   ⚠ 但派单⛔ 规定**改后那几发之间也量**——本腿在 `windows-after`(14:35) 与 `core-after`(14:39) 之间**漏量了这一发**，
   具名报为**本腿的过程缺口**（⛔ 静默）。补救读数＝交回时 14:5x 现量一次：`wisp.exe`=**0**／`balldebug.exe`=**0**／`go.exe`=**0**，
   且事后核两把旁证：`core-before` 与 `core-after` 的 `four numbers` 行**逐字相同**（`RUN=1866 PASS=1257 FAIL=8 SKIP=0`）、
   `internal/audio` 的 `TestPinnedThreadStable10s` 在两发同名同形 ⇒ 那两发之间⛔ 有真窗插进来的读数证据为无（尺＝名册，⛔ 颜色）。
   定式建议入册＝**每发长跑前各量一次 tasklist 并把读数并排落进同一枚 rc 件**（⛔ 靠事后旁证补）。
4. **探针落点**：派单写 `probes/301/r1/logs/…`（仓根 `probes/`），本腿按盘上先例落 `.scratch/wisp/probes/301/r1/**`
   （尺＝同族腿 `300-r1`／`301-a1`／`301-a2` 逐笔 `git log --name-only`；仓根 `probes/` 只含 1 枚未跟踪目录、
   `git ls-files probes/`＝0）。`AC#4` 那句 `probes/301/**` 作后缀匹配成立。口令下达后本腿整批平移。
5. **GUARD C／D 归属**：照派单更正执行（那两句在 GUARD D 退码文本 `:407`–`:423`，C 从 `:542` 起、同义句 `:553`–`:558`，
   两枚都本腿 `sed` 现量）。本腿⛔ 在 GUARD C 里补任何一句。
   ⚠ 一处**新读数**给编排者：本改动后 `--scope=windows` 的 **GUARD C 判绿**（`grep -c 'GUARD [ABCD]'` 于 after 日志＋err ＝ 0／0）
   ⇒ 「档清单↔pin 同一笔」这件事的**机器判据**确实长在 C，与「同一笔」那句要求的**文字**长在 D 是两回事，
   票面 14:0x 那句「引文归属≠判红者」被本腿真跑坐实。
6. **`:593` 行号已漂 1**（现 `:594`）：漂因＝本腿在 `win_pin` 插了一枚 token（在它**之上**），⛔ 内容改动。
   ⇒ 任何后续按 `:593` 引用的人都读到⛔ 存在的东西。属票 255 `AC#6` 那一族，本腿⛔ 新造判据，只做具名。
7. **一枚派单⛔ 预见的形状**＝`TestResolvePerCallBudget` 是**墙钟负载敏感**用例（before 1.384 ms/op 判超 1.000 ms 预算／
   after 0.745 ms/op 判过，样本 902→1405）。它让「⛔ 拿本机跑绿当凭据」那条禁区多一枚含义：
   **本机负载会同时污染 before 与 after 的红绿集合**，反方向（假新增红）会让 `AC#2` 的「0 枚新增红」假失败。
   本腿处置＝⛔ 动阈值、⛔ 为变绿放宽断言，只在 `20-ac2-differential.md` §3 具名上报＋建议编排者复核名册作差时
   **把这枚用例单列为已知噪声**（⛔ 本腿自造豁免，⛔ 改 `runtests.sh`）。

## 8. 与编排者 14:3x／14:4x 两枚独立预测的对拉（⛔ 引腿的数，也⛔ 引编排者的数：两把都现跑）

| 预测（`probes/301/orch/2026-10-10-baseline-and-delta-prediction.md` ＋ `6c2dd788` message） | 本腿真跑 |
|---|---|
| `A-evaluated 399→440` | **399→440** ✓ 逐字 |
| `A-skip 0` | **0**（before 0／after 0）✓ |
| `A-pass = 41 − 具名红` | after `A-pass=439`＝before 397 ＋ 41 ＋ 1（`TestResolvePerCallBudget` 基线红翻绿，§7-7 噪声）⇒ 本腿⛔ 读成「41−0」的干净形，**具名偏离 1 枚** |
| `RUN +41／SKIP 0／PASS 41−具名红`（14:3x 那发） | R1 新增**41** ✓；`=== RUN` 计数 577→624＝**+47**（＝R2 含子测试那一把，⛔ 与 41 混）⇒ 她那句按 R1 成立、按 `=== RUN` 那枚数成立到 +47 |
| `撞名尺 comm -12 现量＝0 重合` | 本腿独立复算：`comm -13` 的 41 枚**全部** ∈ audio 的 42 枚名册，`comm -23` 空 ⇒ ⛔ 一枚来自别的包 ✓ |
| `分档单位是包 ⇒ 分母 42（10 tagged／32 无 tag）` | 本腿 R4 现跑＝42／10／32 ✓ 逐枚对上 |

## 9. 本票射程内**未完成**的事（⛔ 由本腿裁，具名交接）

- `AC#3`＝编排者已裁「本格⛔ 落任何一形，整枚射程搬票 111 `AC#12`」⇒ 本腿零动作。
- `AC#5`（`ci.yml:400-401` 那两行过期注释）＝**单独一笔**，⛔ 归本腿这一程（派单逐字「`AC#5` 不归你这程」）。
  现量备好＝`ci.yml` 那两行在 HEAD 上逐字仍是
  `  # Windows-only packages (ball GUI, cgo speech) are out of this job's scope` ＋
  `  # by platform, not skipped: they run in test-windows / slo jobs.`（`sed -n '400,401p'`），
  而 `AC#1` 落地后前半句「runs in test-windows」对本包**开始成真**，后半句「slo jobs」仍⛔ 真。
- ★本腿⛔ 翻任何框：票面 `AC#0`／`AC#3` 之外，`AC#1`／`AC#2`／`AC#4` 三格**一枚不勾**（勾＝非实现者裁，`SPEC-12 §4.3` #1）。
- ⛔ 零 push。
