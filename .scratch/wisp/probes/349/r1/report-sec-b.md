# 票 349 落地腿 `349-r1` 交件报告｜第 2 件＝第 6 节至第 10 节（8 枚牙／成对四数／门禁／纪律面／实话栏）

承 `report.md`（第 0 节至第 5 节）。本件只装读数，⛔ 重复判据。

## 6. 票 336 那 8 枚牙（复跑＝便宜保险）⇒ **8/8 全红**

台件 `.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh` **一字⛔ 动**：blob `454da3ee5feabf39cb63ef1bea95178e8ea2c432`
⇄ 我交件 HEAD 同一枚（＝与 `336-v1` 交付锚 `c315fa4c` 那枚同值，我逐枚 `cat-file -e`/`rev-parse` 现验）。
env 指**我自己 `git archive` 重导**的 `treeDel`（锚 `db2f9eab`；树里 `blacklist.go` blob＝`964ee75a…`＝锚侧逐字对撞）；
⛔ `-overlay`、⛔ 复用别人的树。逐发验：`落点命中枚数 needle_before=1 needle_after=0`（8 枚全中）＋
`branch_removals_on_disk=1`（恰一枚操作数失去回落支）＋ `changed_lines=2`。
★还原哈希＝`964ee75a4befc7a51424f409a95745296e369da3`，`match=YES` 共 **11 次**（start／after-pristine／
8×after-MU／final），⛔ 任何一发 NO；台件跑完 `shared_worktree_dirty_internal_cmd=0`。

| 发（落点行） | 牙 roster `TestAnchorFallbackTeethPerAnchor336` 四数 | 三枚 legacy 仪器四数＋红名 |
|---|---|---|
| pristine（⛔ 落刀正控） | `rc=0` RUN=17 PASS_TOP=1 FAIL_TOP=0 SKIP=0（8 枚 operand 全 `class=A`） | `rc=0` RUN=10 PASS_TOP=3 FAIL_TOP=0 SKIP=0 |
| MU-1a（`:223` git-credentials） | `rc=1` 17/0/**1**/0，红子名 `/git-credentials`（`class=B`） | `rc=0` 10/3/0/0 |
| MU-1b（`:227` wisp-config-under） | `rc=1` 17/0/**1**/0，红 `/wisp-config-under` | `rc=0` 10/3/0/0 |
| MU-1c（`:228` wisp-config-eq） | `rc=1` 17/0/**1**/0，红 `/wisp-config-under`＋`/wisp-config-eq`（顶层 1 枚、子名 2 枚） | `rc=0` 10/3/0/0 |
| MU-1d（`:230` ssh） | `rc=1` 17/0/**1**/0，红 `/ssh` | `rc=0` 10/3/0/0 |
| MU-1e（`:232` aws） | `rc=1` 17/0/**1**/0，红 `/aws` | `rc=1` 10/2/**1**/0，红 `TestAnchorMissReasonNamesCandidateAndMissingAnchor/absent-aws-leaf-accuses-aws-and-says-so` |
| MU-1f（`:234` kube） | `rc=1` 17/0/**1**/0，红 `/kube` | `rc=1` 10/2/**1**/0，红 `…/installing-aws-moves-the-named-anchor-to-kube` |
| MU-1g（`:238` dpapi-protect） | `rc=1` 17/0/**1**/0，红 `/dpapi-protect` | `rc=0` 10/3/0/0 |
| MU-1h（`:240` local-credentials） | `rc=1` 17/0/**1**/0，红 `/local-credentials` | `rc=0` 10/3/0/0 |

⇒ 8 枚逐名各自兑现（每发都有 ≥1 枚顶层红＋操作数名点名），⛔ 一枚被我改红的常驻测试。
逐发原始 stdout＝`.scratch/wisp/probes/349/r1/logs/teeth/ac2-*.log`（18 枚件）＋`ac2-sweep-runner.log`。
⚠ 四数⛔ 混用：`FAIL_TOP`＝顶层 `--- FAIL`，子名另数（MU-1c 顶层 1 枚／子名 2 枚＝台件 `red roster` 那行
`teeth_red=3` 把两行都算了＝公共块"四数⛔ 相减"那一形，具名摊出）。

## 7. AC#4 成对四数（两平台，同一台面，`-count=1 -timeout 30m`，⛔ 混用⛔ 相减）

| 平台 | 侧 | RUN | PASS_TOP | FAIL_TOP | SKIP | rc |
|---|---|---|---|---|---|---|
| Windows（导出树 `treeBase`＝起手锚 7b633fa9） | base | 247 | 136 | 5 | 1 | 1 |
| Windows（导出树 `treeDel`＝交件锚 db2f9eab） | del | 247 | 136 | 5 | 1 | 1 |
| Linux（docker `golang:1.27`，同样两棵树） | base | 195 | 115 | 1 | 1 | `GO_TEST_RC=1` |
| Linux | del | 195 | 115 | 1 | 1 | `GO_TEST_RC=1` |

红名册**逐名作差＝两侧同集**（Windows⇄Linux 各自 `comm` 两向都空 ⇒ **本腿新增红＝0 枚**）：
Windows 5 枚＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（＝票面唯一合法存量红，`A380`／票 200／216），另 4 枚＝
`TestSyncAncestorActableLegFailsClosedOnMovedAncestor116`、
`TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116`、
`TestSyncAncestorLegReResolvesThroughJunctionToRealTree`、
`TestSyncAncestorLegNotExemptedStaysFailClosed`。
★后 4 枚⛔ 是我造的：票 343 终裁腿 `343-v1` 在同包同台面读到的是**同一集**（`RUN=247 PASS_TOP=136 FAIL_TOP=5 SKIP=1`、
红名册 5 枚逐字同名，见 `.scratch/wisp/probes/343/v1/logs/pair/pair-summary.log`），我这把尺从零复现＝同值。
⇒ 这 4 枚是**平台／台面相关存量**（`//go:build windows`＋junction/重解析点在本机真树上的行为），
`A380` 那枚是唯一"票面承认的合法红"，其余 4 枚请编排者按存量归口（我⛔ 动任何测试字节、⛔ 翻任何框）。
Linux 侧红名册＝只 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 1 枚（SKIP＝`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`）。
`TestResolvePerCallBudget`＝Windows 两侧都 **PASS**（base 1.19s／del 1.21s）⇒ 本波⛔ 撞争用、⛔ 需复量 3 发、
⛔ 计入我的新增红。台面故障名册＝`none`（`walk`／`build failed`／`no test files`／`setup failed`／
`creating work dir`／`0xc0000135` 四份日志全 0 命中）；两侧 `RUN` 名集⛔ 用作"我改了 N 枚测试"的凭据（公共块明令）。

## 8. 门禁四数（★只在共享工作树、跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）

| 门 | rc | 读数 |
|---|---|---|
| `gofumpt -l` 三枚被动文件 | **0** | 打 **blob 形**（LF）时 base⇄del 两侧都 **0 行**＝我的新字节 gofumpt 干净；打**工作副本**时 `provenance.go` 被标脏＝该文件磁盘副本 CRLF 的**存量**（起手锚的 blob 抽出到磁盘同样被标、⛔ 我引入、我⛔ 顺手重排别腿文件）＝公共块"本机 6 枚既有件因 CRLF 被标脏"那一形，具名登记 |
| `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；examined 283 production Go files（internal 241＋cmd 42），ban #8 internal 592 枚含注释与 `_test.go` |
| `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`；tracked paths=16549／over-budget=60／rostered=60／not in roster=0 |
| `go vet ./internal/risk/` | **0** | 空输出（＝无告警） |

原始读数＝`.scratch/wisp/probes/349/r1/logs/gates/g1.txt`／`g1-before.txt`／`g1-after.txt`／`g2.txt`／`g3.txt`／`g4.txt`。
`g2.txt:191` 那句 `ban #8 scope frontend/ examined 0 files` 是 d22scan **自己的注释文本**（含 `walk` 字样），
⛔ 是红句、⛔ 台损＝具名登记，别腿拿 `grep -nE 'walk'` 扫 g2 会撞上它。

## 9. 纪律面自证（尺⛔ 我说）

- 名册尺（锚定行首）：起手 **0 枚**（＝第 0 节）→ 交件时＝本腿全部提交（逐枚主题起头 `票 349 落地腿 349-r1…`）。
- ⛔ 翻框：票 349 盘上现量 `^- [ ] **AC#`＝**7 枚未翻**、`^- [x]`＝**0**（＝七格一枚没翻）。
- ⛔ 改票名／`Status:`／`-done`：`git ls-tree HEAD --name-only -- .scratch/wisp/issues/ | grep -c 'issues/349-'`＝**1**、
  文件名与 `Status:` 行逐字未变（`git diff` 该票只多我那一枚 bullet 一行）。
- 写面⛔ 越界：逐笔 `git show --name-only` ⊆ {`internal/risk/blacklist.go`, `internal/risk/provenance.go`,
  `internal/risk/syncdirs.go`, `.scratch/wisp/probes/349/r1/**`, 票 349 尾一行}；`.go` 命中恰 **3**；
  ⛔ `internal/winsec/**`（`342-r1` 独占，我起手现量它名册 0 枚＝尚未启动，我仍⛔ 碰那一包一字）、
  ⛔ `cmd/wisp/**`／`internal/panel/**` 写、⛔ `go test ./cmd/wisp` 整包、⛔ `go test ./...`。
- 临时件全部落 `D:\tmp\wisp\20261019\349r1\`（⛔ 盘符根）、**只建不删**（⛔ `rm`，中途一次 `rm -f .tar` 我已在
  脚本里改成"留着不删"＝`build-tree.sh` 那行注释）。
- 证据件⛔ `.out` 扩展名：`find .scratch/wisp/probes/349 -name '*.out'`＝**0 枚**；`-size 0`＝**0 枚**；
  逐枚 `wc -c` 非 0（尺＝`sh/hygiene-scan.sh` 读数）。
- 被禁字面字符 U+00A7＝写面逐枚 **0 命中**（正控＝同一把尺打 `docs/reports/HANDOVER.md`＝304 非零）；
  ⛔ 与被禁字符混用间隔号 U+00B7（本腿主题里合法，摊出来证明两把尺没混）。
  ⚠ 我第一版卫生尺用**字节** 0xA7 计数＝会误报（三字节 UTF-8 序列的续字节可能是 0xA7），已改成
  `-CSD` 码位形并把字节形具名摊出＝"两把尺矛盾时先假定新那把坏"那一形我自己踩到。
- `add`／`commit` 两处都逐枚字面 pathspec（⛔ 从共享 index 反推、⛔ `-A`／`.`／通配）；新件先 `git add` 才进 pathspec；
  中文正文 Write 成文件后 `git commit -F`，`-F` 一律 **Windows 形**路径；每次提交后现跑 `git log -1`＋`git status` 对回。
  ⛔ push（全程）。
- 两次 `commit -F` 失败＝`rc=128`（`-F` 路径写成 bash 形／写成不存在的临时路径）：⛔ 任何副作用、
  我按现跑对回后才重打，⛔ 拿"命令没报错"当已落库。

## 10. 本程没测什么／未决（实话栏）

1. ⛔ 跑 `go build ./...`、⛔ 跑 `cmd/wisp`／`internal/panel` 任何测试（派单禁）⇒ 我这三枚注释改动对**包外**的
   影响只由"剥注释后可执行字节三把尺同串"这把结构尺承担，⛔ 其它包的成对读数。
2. Linux 侧成对＝docker `golang:1.27` 里跑我自己的两棵树，**⛔ 8 枚牙**：那批 operand 测试文件带
   `//go:build windows`，Linux 侧结构上跑不到 ⇒ 牙只在 Windows 侧复跑，两平台各自摊开（⛔ 拿 Linux 少 4 枚红当"我修好了"）。
3. ⛔ 复核 `343-v1` 判语的其余格（那不是我的票）；我只复跑与"注释改动会不会碰到那 8 枚牙"相关的那一把＋同包四数。
4. 第 3 节 那 4 枚候选（两枚 SC 假阳＋两枚 `_test.go` 面）具名交回＝**归口请编排者裁**；我⛔ 扩射程、⛔ 动 `_test.go` 一字。
5. 未决坑一枚（⛔ 我拍）：`internal/risk/provenance.go` 的工作副本是 **CRLF** 而 HEAD blob 是 **LF**
   （`.gitattributes *.go text eol=lf` ＋ `core.autocrlf=true`）。我这枚改动没引入它（起手侧同一状态、同一把
   `gofumpt -l` 同样标脏），但它让"拿磁盘文件比剥注释字节"这类尺**天然假红**、`gofumpt -l` 也长期把它报脏
   ⇒ 建议编排者立案（存量、⛔ 本票射程）。
6. ⛔ 造门（AC#3 只出结论＝第 4 节）；票 350 已在那儿承接有界词表形。
