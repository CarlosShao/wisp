# 275-v1 — 非实现者对抗验收（形丁落在 `.scratch/wisp/probes/161/r5/attrib.sh`）

Leg: `275-v1`. Role: adversarial acceptance (NON-implementer). Code under review = leg `275-r1`
(commits `4dc6cab6` → `febca8f9` → `d10697eb` → `31261b4d`, orchestrator's own landing commit `79fbcd57`).

Write surface: `.scratch/wisp/probes/275/v1/**` only. Scratch/mutations OUTSIDE the repo under
`D:/tmp/wisp275v1/`. Evidence extension = `.txt` (⛔ not `.out`, `.gitignore:8` is a repo-wide `*.out`).
I did not edit `attrib.sh`, did not flip any ticket box, did not touch `.github/**`, `scripts/**`,
`cmd/**`, `internal/**`, or the two other tickets' stage pieces.

---

## §0 起手锚点（在任何长跑命令之前落盘；本节即第一笔 commit 的内容）

Raw file: `anchor-raw.txt` (same bytes copied to `D:/tmp/wisp275v1/anchor-raw.txt`).

```
$ date                                  Wed Oct  7 12:45:32 CST 2026
$ git rev-parse --short HEAD            79fbcd57
$ md5sum …/161/r5/attrib.sh             77ee11a8d0cc8dd8248c6a9b8e164a78
$ wc -l …/161/r5/attrib.sh              459
$ sh -n …/161/r5/attrib.sh              sh_n_rc=0
$ git status --porcelain                753 lines  (full list -> D:/tmp/wisp275v1/porcelain-anchor.txt)
$ git status --porcelain -- cmd internal scripts tools .github docs frontend
                                        0 lines    (PORCELAIN ANCHOR = EMPTY for these groups)
$ git ls-files '*.go' | wc -l           935
$ git diff --numstat febca8f9^ febca8f9 -- …/161/r5/attrib.sh
                                        19  0  .scratch/wisp/probes/161/r5/attrib.sh   (19 added / 0 deleted)
```

Restore proof compares verbatim to: (a) `md5sum attrib.sh` = `77ee11a8d0cc8dd8248c6a9b8e164a78`,
(b) `git status --porcelain -- cmd internal scripts tools .github docs frontend` = **0 lines**,
(c) the 753-line full-repo porcelain file.

Authority read on disk (my own reading, not the ticket paraphrase):
- `attrib.sh:341` `A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?`
- `attrib.sh:342` `if [ "$A_RC" -gt 1 ]; then`
- `attrib.sh:343` stderr `echo` "…the tracked tree is not readable by this ruler"
- `attrib.sh:344-353` the added comment block (A663 / shape 丁)
- `attrib.sh:354` `if [ -n "$A_OUT" ]; then`
- `attrib.sh:355` banner `== (A) per-file roster gofumpt emitted before it hit an unparseable file …`
- `attrib.sh:356-361` `while IFS= read -r roster_line` + `printf '   A-ROSTER\t%s\n' "$roster_line"`
- `attrib.sh:363` `exit 2`

The landing is exactly the shape `A663` :13149 authorized (only print `A_OUT` before `exit 2`;
no exit-code, guard, or denominator change).

Note on the ticket vs disk: ticket §2 row 5 quoted `attrib.sh:173-178`; the ticket's own
correction block (`:31-36`) says the real branch is `:338-344`. On disk (459-line version) the
guard is at `:342` and `exit 2` now sits at `:363` after the insertion. **Disk wins** — I cite
line numbers from the current 459-line file.

Method note (how I got a denominator without touching the shared tree): `git clone --no-hardlinks`
to `D:/tmp/wisp275v1/clone` (HEAD `959529c5` = my anchor commit, so the clone carries the SAME
landed instrument: `md5sum clone/…/attrib.sh` = `77ee11a8d0cc8dd8248c6a9b8e164a78`). The clone's
checkout is **LF** (`tr -cd '\r' < clone/cmd/wisp/models.go | wc -c` = **0** CR bytes while the
shared working tree holds stale CRLF there) ⇒ the clone is a genuine **CI-equivalent** shape, and
all index surgery (`git add` / `git rm --cached`) happened only inside the clone. My plant `.go`
files were written directly by heredoc (LF), ⛔ never via `git archive` (`.gitattributes` +
`core.autocrlf=true` would have rewritten their bytes — A664 ②, 第 118 条).
Per-run index state and outputs: `logs/run-ledger.txt`; raw stdout/stderr of every run in `logs/`.

---

## §1 Cell-by-cell verdicts

| 格 | 判 | 我的读数（一手） |
|---|---|---|
| **AC#3 零误伤** | **成立** | rc=2 形与 rc=1 形都点名了我另造的样本；两形名集**逐串相同** |
| **AC#4 恒真性进攻** | **成立**（两处进攻都红得起来；⚠ 但**没有任何自动门**为这支打印保票，见 §2 附注与 §5） | (i) 注释掉打印 ⇒ `rc=2`／stdout **1 行**／名册 **0** 枚；(ii) 坏件排第一 ⇒ 名册**不缩水**（5/5 全在） |
| **AC#5 门禁＋还原** | **本地半＝成立；CI 半＝〔待推送取数〕；`go vet ./cmd/wisp/` 半＝未做（派单禁根模块 go 窗口，与票面 `:55` 冲突，冲突具名上报）** | d22scan rc=0（正控 `TestScanDetectsAllSeededViolations` 先 PASS）；`--self-test` rc=0 `cases=8 failures=0`；致盲正控 rc=1 `failures=4`；`sh -n` rc=0；`attrib.sh` md5 起手＝终态；porcelain 七族 0 行、全仓 753 行仅 1 行差＝我自己的 `?? …/275/v1/logs/` |

Verbatim anchors for the readings above:
- `attrib.sh:341` `    A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?` — the
  denominator expression; unchanged by the landing (`git diff --numstat febca8f9^ febca8f9` = `19 0`).
- `attrib.sh:363` `        exit 2` — still the exit of that branch; my runs reproduce **rc=2** in
  both shapes (CI-equivalent and working tree).
- `attrib.sh:355` (banner, verbatim):
  `        echo "== (A) per-file roster gofumpt emitted before it hit an unparseable file (stdout; exit code unchanged, denominator unchanged):"`
- `attrib.sh:358` `                printf '   A-ROSTER	%s\n' "$roster_line"`
- ticket `:53` AC#3 verbatim requirement: 「另造一枚**语法正确但没格式化**的 tracked 样本（放仓外、经 `-overlay`／临时索引，⛔ 不原地编辑共享工作树），证明选形落地后**这类真回归仍然会被看见**」.
- ticket `:54` AC#4: 「判据换成反形也必须红——指名哪一发突变会让本票落地的判据响；如果**换形后仍全绿**＝那把尺对本件事不敏感」.
- A663 `:13146` 「**代价＝零**」；A667 `:13203` 「rc=2／stdout 35 行（其中 `A-ROSTER` 33 行）」.

### AC#3 detail (both shapes name the plant)

Plant: `plant/zz_ac3_plant.go` — syntactically valid, 4-space-indented body (gofumpt rewrites it ⇒
`gofumpt -l` lists it; verified standalone: `logs/` + the pre-check showing all 6 of my unformatted
samples flagged while both broken samples flagged `:5:1: imports must appear before other declarations`).

(a) **rc=2 形**（full denominator + plant, broken sample present; `logs/ac3-rc2-out.txt`）:
`tracked .go files handed to it: 936` ⇒ **rc=2**, stdout **23 行**, `A-ROSTER` **21 行**,
and the plant is named at stdout line 20 verbatim:
`   A-ROSTER	.scratch\wisp\probes\275\v1\plant\zz_ac3_plant.go`
stderr unchanged (1 line): `attrib.sh: (A) gofumpt exited 123 on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler`

(b) **rc=1 形**（same index minus ONLY `fs_broken.go`; `logs/ac3-rc1-out.txt`）:
`handed to it: 935` ⇒ **rc=1**, `TRACKED-DIRTY` **20 行**, `A-ROSTER` **0 行**（the new block is
only reachable when `A_RC>1`, so the normal path is untouched）, stderr **0 bytes**, plant named verbatim:
`   TRACKED-DIRTY	.scratch/wisp/probes/275/v1/plant/zz_ac3_plant.go	(last commit touching it: )	<== CI IS RED ON THIS`
then `attrib.sh: RED (tracked-only) - tracked-dirty=20 files_dirty=20. …`

(c) **名集相同**（this is the zero-false-injury proof proper）: normalizing both rosters (backslash→slash,
strip `:line:col: msg`, drop the `fs_broken.go` diagnostic row) ⇒ `logs/n-rc2.txt` **20 行** vs
`logs/n-rc1.txt` **20 行**, `diff` **clean**. So a red that cannot read the tree still names **exactly**
the same set as a red that can. Both lists include the two other tickets' stage pieces
(`logs/n-rc1.txt:13` `…/probes/241/v1/posctl/badly_formatted.go`, `:15` `…/probes/259/r1/gofumpt-negctl/probe.go`)
⇒ the gate did not swallow their "该被看见" either.

(d) ★Bonus cell the implementer explicitly did not measure: the clone run gives **20 = 19 CI-equivalent + my 1 plant**
(`grep -v zz_ac3_plant logs/n-rc1.txt | wc -l` = **19**), `outside .scratch` = **0**, `in probes` = 20.
⇒ **A663 `:13137` / A667 的那枚 19 由独立一手复现**（不需推送、不需碰共享树）。

### AC#5 detail (gates)

- `sh scripts/d22scan.sh` ⇒ **rc=0**（`logs/d22scan.txt`, stderr 0 bytes）。Positive control runs FIRST inside it
  and is green: `--- PASS: TestScanDetectsAllSeededViolations (0.04s)`, then
  `TestScanCleanRepoIsGreen` / `TestCheckRootRejectsBlindRoots` / `TestScanAloneIsNotAFalsifier` /
  `TestEmojiBanCoversGoSourcesNotJustDesign` / `TestBan8MathBandAndRemainingGaps` … all PASS, ending
  `d22scan: clean - no D22 ban violations`. (d22scan is its own module `tools/d22scan` ⇒ no root-module go window used.)
- `sh .scratch/wisp/probes/161/r5/attrib.sh --self-test` (landed copy, shared repo) ⇒ **rc=0**,
  `attrib.sh: --self-test cases=8 failures=0` (`logs/selftest-landed.txt`) — **re-run, not copied**; matches impl `:168`.
- Blinded-classifier positive control, MY mutation (different site than the leg's, same effect): copy out of repo,
  `attrib.sh:233` `    CL_RC=1` → `    CL_RC=0` ⇒ `sh -n` rc=0, run with `ATTRIB_ROOT=<shared repo>` ⇒
  **rc=1**, `attrib.sh: --self-test cases=8 failures=4` (`logs/selftest-blind.txt`) — carrier still has teeth for the classifier.
- `sh -n .scratch/wisp/probes/161/r5/attrib.sh` ⇒ **rc=0**（`gofmt -l` not applicable: not a `.go` file）.
- `md5sum attrib.sh` 起手 `77ee11a8d0cc8dd8248c6a9b8e164a78` = 终态同一串（我只用仓外拷贝做突变）。
- ⛔ **CI 侧那一半（这一步真跑过／什么颜色）＝〔待推送取数〕**：机主 今日不推 ⇒ 我没有、也不许用上面任何一发本机读数冒充 CI 读数（第 109 条）。
- ⚠ 冲突具名上报：票面 `:55` 要求 `go vet ./cmd/wisp/` rc=0，本派单 ⛔ 禁在根模块起 go 窗口（另一腿持有）。
  按"派单与票面冲突时盘上/派单赢并上报"的规矩，那一发**未跑**，登记在 §5，⛔ 不算成立。

---

## §2 两发恒真性进攻（transcripts）

### (i) 把新增的打印块整段注释掉 ⇒ 红退回"只有 1 行、零具名"

`cp clone/…/attrib.sh /d/tmp/wisp275v1/attrib-neutered.sh`（源件 md5 = 落地版 `77ee11a8`），
`sed -i '354,362s/^/#MUT /'`（正好罩住 `attrib.sh:354-362` 那 9 行，`exit 2` at `:363` 不罩），`sh -n` ⇒ rc=0。
分母＝同一发 936（含我的样本、含 `fs_broken.go`），`ATTRIB_ROOT` 指回 clone（instrument :127 那支就是为仓外拷贝存在的）。

```
MUTATED_RC=2
out_lines=1 roster=0
== (A) tracked tree: git ls-files -z '*.go' | xargs -0 D:\work\base\gopath/bin/gofumpt.exe -l   (tracked .go files handed to it: 936)
attrib.sh: (A) gofumpt exited 123 on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler
```
对照同一分母未突变那发：`RC=2 / out_lines=23 / A-ROSTER=21`。⇒ ** stdout 从 1 行变 23 行这件事，唯一来源就是新增的打印块 **；
退码、stderr、守卫、分母四处**一字未动**。这一发同时是"落地前形状"的复现（AC#1 那发行情的反向证据）。

**有没有牙？**——对"名册必须出现"这件事：**接受判据有牙（我这一读就红），本仓自动门没牙**。
证据＝同一份被删掉打印的拷贝跑票 161 的自检载具 ⇒
`NEUTERED_SELFTEST_RC=0 … --self-test cases=8 failures=0`（`logs/selftest-neutered.txt`）⇒
**载具对这支打印完全失明**；今天若有人把 `:354-362` 删掉，CI 那步照旧 rc=2 照旧红，**没有任何断言会为此变脸**。
这不是"换形后仍全绿所以不敏感"那种假凭据（我的读数本身会红），但它是**一支没有正控的落地**，
必须补（票 `:45` 为形丙提过同族要求："种一枚新坏文件 ⇒ 名册里必须出现它的具名行"；丁同样需要它）。

### (ii) ★完整性进攻：坏件排在最前时，名册会不会缩水

分母只留 6 枚（clone 索引 `git rm --cached` 全部 `.go` 后只 add 我的件），gofumpt 一次拿到全部参数。

```
--- ls-files order (run A, broken FIRST) ---
.scratch/wisp/probes/275/v1/plant/aaa_broken.go        <- 排第一
.scratch/wisp/probes/275/v1/plant/u1_unformatted.go
… u2 u3 u4 u5 …
A_RC=2 A_lines=8 A_roster=6
   A-ROSTER	.scratch\wisp\probes\275\v1\plant\aaa_broken.go:5:1: imports must appear before other declarations
   A-ROSTER	.scratch\wisp\probes\275\v1\plant\u1_unformatted.go
   A-ROSTER	…u2…  u3  u4  u5 …                 （5 枚需格式化的全在）

B（只差那一枚坏件，denom=5）: B_RC=1 B_dirty=5   —— u1..u5 全部 TRACKED-DIRTY

--- ls-files order (run C, broken LAST) ---
C_RC=2 C_lines=8 C_roster=6   （同样 5 枚 + 1 行诊断）
```

⇒ **A 的名册 == B 的名册（5/5），C 亦然。名册不随坏件的位置缩水。**
再加上全分母那一发（936 枚、要跨多轮 `xargs` 批次）的**名集逐串相同**（§1(c)），
两把（小集合定向＋真实集合对拉）都判"完整"。
⇒ **"红而不名"这件事在丁之下确实被消掉了；丁不是半截子工程。**

**CI 会看见哪一形，为什么**：CI 那步（`.github/workflows/ci.yml:235`，守卫在 `:234`）跑的是
`actions/checkout` 检出的 **LF 形**（`.gitattributes` 里 `*.go text eol=lf` 压过本机 `core.autocrlf=true`），
所以 CI 那一发 = 我这发的形状：**rc=2 / stdout 23 行 / 名册 21 行 = 20 枚需格式化 + 1 行 `fs_broken.go` 诊断**
（本机工作树那一发是 33 行 = 32 + 1，多出的 12-13 枚＝陈旧 CRLF 字节，A663 `:13139-13140`）。
两形都**完整具名**，⛔ 不存在"5/30 那种半瞎名册"。

⚠ 但我必须记一句**措辞**缺陷（不是判据缺陷）：那行 banner 自称
`gofumpt emitted before it hit an unparseable file`（暗示"撞上坏件就停了，所以只是前面那一截"），
而我上面量的是**它根本没截**；banner 又自称 `(stdout; …)`，而 `attrib.sh:341` 用的是 `2>&1`，
`A_OUT` 里**混的是 stderr 的诊断行**。⇒ 拿 `grep -c A-ROSTER` 数"有几枚要格式化"的人**必然多算 1 枚**
（真实集合 936：21 行 = 20 枚 + 1 诊断；工作树形：33 行 = 32 枚 + 1）。
rc=1 那一支有 `lines=/files=` 两个数自辨（`attrib.sh:377-378`），rc=2 这一支**没有**，所以只能靠人读格式。

---

## §3 实现者（`275-r1` impl.md）那些句子：我confirm / 我推翻

- **CONFIRM**（一手复跑）`impl.md:127` "stdout = 35 lines … 33 `A-ROSTER` rows"：我在共享工作树同一把尺
  ⇒ `WT_RC=2 WT_lines=35 WT_roster=33 WT_err_lines=1`（`logs/wt-rc2-out.txt`／`-err.txt`）。
- **CONFIRM** `impl.md:144-146` "A-ROSTER count in the rc=1 shape = 0 … normal red path byte-for-byte what it was"：
  我在 clone 的 rc=1 形量到 `A-ROSTER=0`、stderr 0 字节、`tracked-dirty=20`。⚠ 细节更正：它说
  "(last commit touching it: %s)" 那一段——我的样本没历史，打出来是**空串**而不是 `-`
  （`attrib.sh:372` 的 `|| echo '-'` 不触发，因为 `git log` 对无历史路径**退 0 且输出空**）。不是违规，是形状。
- **CONFIRM** `impl.md:168-173` 自测 8/0 绿与致盲正控 4 红：我换了突变位（`:233` `CL_RC=1`→`0`）也得到
  `cases=8 failures=4` / rc=1 ⇒ 它那句"载具仍有牙"为真，**但射程只有分类器**（见 §2(i)：删掉打印⇒载具全绿）。
- **CONFIRM + 它是它自己最准的一句** `impl.md:201-203` "the 丁 roster print has **NO assertion yet** —
  it prints; whether it prints is not itself gated"。我为它补了证据（`logs/selftest-neutered.txt`）。
- **OVERTURN（把它的"未做"结掉，不是推翻它的结论）** `impl.md:154/205` "CI 等价 LF 形的 19 枚那一发它没复现"：
  仓外 `git clone` 就是 CI 等价形（CR bytes=0），我没碰共享树就复现了 **19**（加我的样本＝20）。
  ⇒ 这枚数**不需要推送、也不需要 `-overlay`**，先前把它当成"取不到"是射程估计偏保守。
- **OVERTURN（措辞）** 我不同意它 §3 对 banner 的默认接受：banner 的两处自称（"before it hit an unparseable
  file"、"(stdout; …)"）都与实测不符，见 §2 末段。落地没坏，文案坏了。

---

## §4 编排者（票面 / A663 / A667）那些句子：我推翻什么

1. **推翻 A663 `:13146` 与票 `:93` 的「代价＝零」**：实测代价 ≠ 0，但 ≠ 大。三件具体的：
   ①rc=2 那一支的名册**没有 `lines=/files=` 计数**，`grep -c A-ROSTER` 会多算 1 枚（诊断行混在里面）；
   ②banner 文案对机制是**错的**（说会截、实际不截；说 stdout、实际是 2>&1 合并）；
   ③这支打印**零断言覆盖**（§2(i) 证据）。⇒ 正确说法＝"判据代价＝零；**可读性/文案代价非零**"。
   这一处不影响我判丁成立，但影响"以后照这段散文再造一支丁"的人。
2. **推翻 A667 `:13206` 的那句「载具仍有牙」的射程读法**：数字我复跑为真（8/0 与 4 红），
   但它被写在**本票**验收段里，容易被读成"这次落地的判据有牙"。实测：把打印块删干净 ⇒ 载具 **rc=0 全绿**。
   ⇒ 那句应加限定："**分类器**仍有牙；本票落地的那支打印**无牙**"。
3. **确认（不推翻）A667 `:13204` 的结论、但把它的机制说反了一半**：编排者/腿那枚"gofumpt 撞坏件时可能整份名册都不出"
   的预判**确实不成立**——我把它做成定向实验（坏件排第一 vs 排最后 vs 没有坏件）坐实了"名册完整"；
   但落地文案写的恰恰是那枚**被推翻的预判**（"emitted before it hit an unparseable file"）。
   ⇒ 判据赢在了预判的反面，而文案还在替预判说话。
4. **确认 A667 `:13203` 的 35/33**（我一手：`WT_RC=2 / 35 行 / 33 枚 A-ROSTER`），并把"落地前后只差名字看不看得见"
   这句做成因果突变（§2(i)）而不只是对比两次读数。
5. **确认票 `:87` / A663 `:13145` 那两枚台件仍被点名**（`logs/n-rc1.txt:13,15`），⛔ 我没动它们，
   且**门没把它们吞掉**——这是我这格最想要的反向证据之一。

---

## §5 我没做完／判不动的格（具名）

- **AC#5 的 CI 那一半**：〔待推送取数〕。⛔ 我没有用任何本机读数替代它。
- **AC#5 的 `go vet ./cmd/wisp/`**：未跑（派单禁根模块 go 窗口）。票面 `:55` 与本派单冲突，已具名上报。
- **票 275 AC#0/#1/#2 的原始那一把**（935 枚逐枚退出码名册、"恰好 1 枚不可解析"的全量证明）：
  我没重跑编排者已翻过的那三格；我给的是**等价形状**的独立读数（clone 全分母 + 只差一枚坏件的两把）。
  若有人要按"逐枚 rc 名册"再审一次，那一格仍属〔未做〕。
- **我没有给"打印名册"补断言**（不是我的面：我是非实现者，且 A663 的解冻只罩 `:338-344` 那一支的打印，
  不罩新增判据）。⇒ 需要另开一枚具名解冻/新票：要求"删掉打印 ⇒ 某枚载具必须红"。
- **形丙的那枚定向判据**（票 `:45`）依然不存在，我也没造。

---

## §6 最终自证（还原）

```
$ md5sum .scratch/wisp/probes/161/r5/attrib.sh
77ee11a8d0cc8dd8248c6a9b8e164a78   == §0 锚点（起手＝终态，我没动仪器；所有突变都在仓外拷贝）
$ sh -n .scratch/wisp/probes/161/r5/attrib.sh ; rc=0
$ git status --porcelain -- cmd internal scripts tools .github docs frontend
0 lines                            == §0 锚点（空）
$ git status --porcelain | wc -l
753 lines                          == §0 锚点的 753；逐行 diff 只有 1 行差：
  < ?? .scratch/wisp/probes/275/v1/                 (锚点：整目录未跟踪)
  > ?? .scratch/wisp/probes/275/v1/logs/            (终态：我的两件已 commit，只剩 logs/)
⇒ 差集 100% 落在我自己的写面 `.scratch/wisp/probes/275/v1/**`；别人的在飞物（8 枚 ` M probes/161/r6/logs/flip-*.txt`、
  3 枚未跟踪目录、`design/**` 那一整片）**逐枚仍在原位、行数一字未变**（A667 `:13211` 登记过同一形状，mtime 归它们自己）。
$ git ls-files '*.go' | wc -l        935  == 起手（分母没被我动过；我的样本只活在 clone 的索引里）
```

我的仓外产物（⛔ 不入库也不删）：`D:/tmp/wisp275v1/`（clone、`attrib-neutered.sh`、`attrib-blind.sh`、
`porcelain-anchor.txt`/`porcelain-now.txt`/`porcelain-final.txt`、每发的 raw）。入库件：
`.scratch/wisp/probes/275/v1/verdict.md`、`anchor-raw.txt`、`logs/*.txt`（21 件，第二笔 commit `git show --stat` 自数＝**23 files changed**）。
全部 commit 带显式 pathspec、消息先写在仓外文件、⛔ 无 push/amend/reset。

**收尾复核（两笔 commit 都落完之后现量，`porcelain-final.txt`）**：
```
$ md5sum .scratch/wisp/probes/161/r5/attrib.sh
77ee11a8d0cc8dd8248c6a9b8e164a78                == 起手（仪器一字未动）
$ git status --porcelain -- cmd internal scripts tools .github docs frontend | wc -l
0                                                == 起手
$ git status --porcelain | wc -l ; diff vs anchor
752  vs  753  —— 唯一差行＝锚点里的 `?? .scratch/wisp/probes/275/v1/`（我的两件＋logs 现已 commit，
                 该未跟踪条目自然消失）；其余 752 行**逐字相同**
$ git ls-files '*.go' | wc -l
935                                               == 起手（分母没被我动过）
```
⇒ 还原自证成立：差集 100% 由我自己的写面构成，别人的在飞物一枚没动。

### 裁决一览（给编排者，⛔ 我不翻任何框）

- **AC#3 成立**、**AC#4 成立**（两发进攻都响；完整性进攻判"名册不缩水"；附一枚必须补的洞：打印无断言覆盖）。
- **AC#5 本地半成立／CI 半〔待推送取数〕／`go vet` 半未做** ⇒ 本格整体**不该按"成立"归档**，按"三缺一待推"登记。
- 形丁本身：**成立、可保留**（判据语义零移动，红从"零具名"变成"全具名"）。
- 票的**第二桩病（`exit 1` 点名 19 枚）没被这次落地解决，也不可能被它解决**——见下面那段判断。

### 判断（不动手）：`exit 1` 那一桩能不能在不瞎的前提下变绿

**不能——按本票现在的写法＋这次落地的东西，那一步永远绿不了；而"能绿"的三支里只有两支不瞎，且都要新代价。**
现量支持：
- 我那一发 CI 等价 LF 形给的是 **rc=1 / tracked-dirty=20**（去掉我的样本＝**19**），`outside .scratch = 0`、
  `in probes = 20` ⇒ 19 枚全是别人工单的探针件，⛔ 没有一枚是产码；其中 **2 枚是票 241 的正控与票 259 的负控**（`logs/n-rc1.txt:13,15`）。
- 丁只改 `A_RC>1` 那一支；rc=1 那一支我量到 `A-ROSTER=0` ⇒ **丁对第二桩零作用**（票 `:94` 自己就这么写的，我 confirm）。
- **即使把两枚台件都豁免，也还是不绿**：剩 **17** 枚（19−2）照旧点名 ⇒ `TRACKED_DIRTY>0` ⇒ rc=1。
  这是新增的一手数字（编排者只给过 19/32/13 那几个）。

要真绿只有两条不瞎的路，代价我都量了边界：
- **丁 + "只豁免两枚已具名台件" + 逐枚格式化其余 17 枚**（编排者票 `:95` 那句"该是只把两枚台件加豁免＋其余照报"的完整体）。
  代价＝**碰 15 个别的工单的探针件**（17 枚来自 14 个不同 `probes/<NN>` 目录，见 `logs/n-rc1.txt`），
  每枚都要具名解冻；这是"绿"最贵的一支，但它是唯一**既绿又不瞎**的一支。
- **不做（对第二桩）**：⛔ 不是"没选项"（票 `:96` 的原话我照抄）。**我推荐这一支**，理由＝落地后的形状已经是
  "红 + 全具名"：CI 那一发在 rc=1 时本来逐枚打 `TRACKED-DIRTY … <== CI IS RED ON THIS`，在 rc=2 时现在也打名册 ⇒
  第二桩**已经有名字**，缺的只是"绿"，而买绿的价钱（甲＝对 286 枚探针件整体失明，票 §2 第 2 行那个数；
  丙＝拿假绿换安静且至今没有它要求的定向判据）比红着贵。
  我**不推荐甲**：它顺手把两枚台件的"该被看见"也抹了（票 `:95` 的判断我实测支持）。
  ⚠ 但我不假装它没用：**甲确实能让那一步绿**——我这手量到 19 枚里 `outside .scratch = 0`、`in probes = 20(含我的样本)`，
  也就是整个 `exit 1` 分母**全部住在 `.scratch/**` 里**，排除这一族就等于把整批名字闭掉。
  ⇒ 甲的代价不是我原先写的"只丢 286 枚里的真回归"，而是**这 19 枚全体＋未来任何探针件从此不可见**，
  而它换来的只是"绿"＝第 108 族那条：一门变绿的方式如果等于把尺闭上，就不该被当成修复。
  ⇒ 结论：**甲＝能绿但等于闭眼；乙＝禁区（改别人的凭据）；丙＝能绿但门失去"语法坏掉"这一支的牙且没补判据；
  丁＋17 枚格式化＝能绿且不掉牙，但要 15 张工单的解冻；不做（保持红、已具名）＝今天最便宜的诚实形状。**
  ⇒ 给编排者的建议：**保留丁；第二桩走"不做＋另立一票逐枚处置 17 枚探针件"，⛔ 不要甲、不要在没补定向判据时走丙。**
