# 票 325 AC#1-AC#5 非实现者验收表（腿＝`325-v1`，裁决者）

- 起手锚＝`71002d00`（＝现 HEAD，`2026-10-16 10:16` 那笔台账 A1001/A1002），起尺时刻 `2026-10-16 10:17:12 +0800`
- 实现者＝`325-r1`（`428ef21e` 改形／`73131d29` 读数／`c7190b27` 证据件／`1e53de91` 票尾 bullet／`a916bf6f` 补档）；
  它的 `docs/evidence/s1/325-runtests-gotmp.md` 与 `.scratch/wisp/probes/325/r1/**` **每一句都当未验证**，下面凡"复认"＝本腿自己现跑对上，凡"更正"＝它那枚读数在 HEAD 上已过期或找错地方
- 票框现量（起手尺）＝`grep -c '^[[:space:]]*- \[ \]'`＝**5 未勾**／`- \[x\]`＝**0 已勾**；本腿⛔ 翻任何框
- **本表不是翻框**：五格各给"成立／不成立／缺哪一发"，判语权在本腿，落框归编排者

## 0. 台面（每一发跑在哪一枚字节上）

导出树＝`git archive HEAD | tar -x`（⛔ 共享工作树当改前那发、⛔ 仓内 worktree）：
`C:/tmp/wisp325v1/head3`（正形＝提交态）与 `…/mixed3`（同码副本，突变只在这里做）。
六件套 `[ -e ]` 逐枚在；`find -name '*.go' | wc -l`＝**1094 枚非空**（仓内 `git ls-files '*.go'`＝1091，差的 3 枚在 `.scratch` 探针码里，尺的口径不同、非缺件）。
`third_party/sherpa-onnx` 的三枚 DLL **不被 git 跟踪** ⇒ `git archive` 带不出来，本腿手动放进导出树，每发 go 前缀 `PATH=$tree/third_party/sherpa-onnx`。
干净台面＝每发 `env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT`；`LOCALAPPDATA=C:\Users\swq\AppData\Local`（根始终在 C: 同卷⇒票 311 的 rc=2 射程拒绝全程零命中，见 §4）。

★**我自己两次导出缺件的红，具名不作判据**（尺在件里，⛔ 抹平）：
`logs/v1-ac3-scope-windows.txt`＝导出只取六件套 ⇒ `internal/projctx [setup failed]`（缺 `frontend/` 那枚 embed 包）；
`logs/v1-ac3-scope-windows-full.txt`＝补了 frontend 仍缺 `design/` ⇒ `TestC21TableColourRowsMatchTokensCSS` 读不到 `design/assets/tokens.css` 而红（`--- FAIL=1`）。
两发都是 `RUN`/`PASS` 尺能看见的**台面故障**，与 `walk`/`build failed`/`no tests to run`/`0xc0000135`/`RUN=0` 同处理 ⇒ ⛔ 计入任何一格。第三发（head3，全树）才是判据。

`runtests.sh` 字节尺（起手先对，⛔ 信 r1 的哈希）：`428ef21e` 的 blob、`HEAD` 的 blob、工作树、两枚导出树，
`sha256` **五处逐字同值** `ad6f6a8d72cbe2a3fc9b7036774b8e16ddf92bfd64ea6055cbf12dde55164bf6` ⇒ 本票射程内没有第六枚手。

## 1. AC#1 改形：**成立**

现跑尺＝`v1-ac1-shape-ruler.sh`，两发（仓内工作树 `logs/v1-ac1-shape-repo.txt`、导出树 `logs/v1-ac1-shape-head3.txt`，读数一致）。
形状一律 `grep -F` 固定串，⛔ 正则形；正控＝基座那枚 `'\\wisp-test-gotmp'` 四处各 **1** 命中（尺能看见字节才谈 0）。

| 处 | `gotmp_run=` 行 | `od -c` 尾接字节 | 固定串 `}\\<name>-$(date` |
|---|---|---|---|
| `tools/d22scan/runtests.sh` | 115 | `" } \ \ d 2 2 s c a n - $ ( d a t e` | 1 |
| `scripts/winsec-tests.sh` | 69 | `" } \ \ w i n s e c - $ ( d a t e` | 1 |
| `scripts/wisp-cli-tests.sh` | 82 | `" } \ \ c l i - $ ( d a t e` | 1 |
| `scripts/portable-tests.sh` | 156 | `e } \ \ p o r t a b l e - $ ( d a t e` | 1 |

⇒ 与三处已修**逐字同形**（本票那枚的尾接＝两枚反斜杠＋`d22scan-`；`portable` 走 `gotmp_base` 变量，但尾接分隔符同形）。
混形尾接 hunt（`'}/d22scan-'` 等四枚）＝**0**，且同屏给了正则假 0 示范：八枚反斜杠的 `grep -c` 读成 0，固定串读成 1（`logs/v1-ac1-shape-head3.txt` 第四段）。

票面判据那两发：**改后**＝AC#4 里 `sh scripts/d22scan.sh`（导出树）rc=**0** clean，发射根逐字
`C:\Users\swq\AppData\Local\wisp-test-gotmp\d22scan-20261016-104730-496`；
**改前**＝`v1-ac1-pre-gate.sh` 把那一行物理退回混形（带门，突变态 sha=`6319065990c0…3240`）后同一枚门 rc=**0** clean，发射根逐字
`…\wisp-test-gotmp/d22scan-20261016-104634-390`；**两发正控制四数逐字同形 `PASS=35 FAIL=0 SKIP=0, === RUN=77`**
⇒ 复认 r1 那句"这一格在未修码上不响＝门绿着、坑没人看见"。
注释 23 行（`:92-114`）点名会炸的那族用例＋两种 cmd 措辞，⛔ 只写风格 ⇒ **合**；⚠ 唯一越界倾向＝它把"AC#2 记为缺世界"写进了产码注释，那一格当时还没被验收腿裁（本腿已裁，见 §2，不构成退回理由，只记形）。

## 2. AC#2 成对反控：**成立**（但本腿第一发作废、具名重跑）

★**先自证一次假读数**（这族尺的牙最容易崩在拿尺的人身上）：`logs/v1-ac2-three-shapes-pair.txt` 里我的
`sed -i "${mline}s/}\\\\d22scan-/…"` 写在**双引号**里 ⇒ shell 吃一枚、sed 再吃一枚 ⇒ 只剩一枚字面反斜杠 ⇒
对不上文件里的两枚 ⇒ **突变静默 no-op**（`after mutate` 那行 sha 未变、`mixed_tail_hits=0`），
而脚本只"打印"了自证、没**拦住**，于是那一栏"B-mixed"跑的其实是正形 ⇒ 读数无效（件留着，⛔ 删）。
该文件里 A 正形绿／C 纯正斜杠四红／D 还原绿三段仍有效。第二枚仪器把自证改成**硬门**（不落就 exit 3、不跑测试）。

第二发＝`logs/v1-ac2-mixed-tail-pair2.txt`（四形同屏，`./internal/risk/ -run '^(四枚)$'`，四枚名先 `-list` 逐枚自证 1/1）：

| 形 | 码／根 | rc | 顶层四数 | 关键红句（逐字摘录） |
|---|---|---|---|---|
| A 正形 | 提交态 `ad6f6a8d…` | **0** | `RUN=4 PASS=4 FAIL=0 SKIP=0` | — |
| B 混形（本票的牙） | 物理退回 `}/d22scan-`，突变态 `6319065990c0…3240` | **1** | `RUN=4 PASS=0 FAIL=4 SKIP=0` | `…reparse_ticket105_windows_test.go:59`／`:116`＝`mklink /J … failed: exit status 1: Parameter format not correct - "d22scan-…`；`…actable_leg_116_test.go:182`／`:286`＝`precondition:` 文案里同一条根两侧分隔符分家（`…wisp-test-gotmp/d22scan-…` ⇄ `…wisp-test-gotmp\\d22scan-…`） |
| C 纯正斜杠根 | 提交态码＋`GOTMPDIR=C:/tmp/wisp325v1/psx2-…` | **1** | `RUN=4 PASS=0 FAIL=4 SKIP=0` | 同四枚名，`:59`／`:116` 措辞＝`Invalid switch - "tmp".` |
| D 还原 | `git show 428ef21e` 灌回 ⇒ sha 复尺 `ad6f6a8d…` | **0** | `RUN=4 PASS=4 FAIL=0 SKIP=0` | — |

⇒ **两形同屏已交**（B 与 C 都炸、措辞不同；A 与 D 都绿）。B 那发的发射根逐字
`GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp/d22scan-20261016-104116-180`＝**那一行自己生成的混形**，不是本腿手搓的字符串。
行号 `:182/:286/:59/:116` 与 r1、与编排者 `10-nat/10-psx`、与 `317-v1` **全部同行同号**；突变态 sha 与 r1 的 `a2-91-sha-reverted-mixed.txt` 逐字同值 ⇒ 跨腿复认。

**回答派单点名那一格（"仪器自红"还是"判据红"）**：r1 的证据件与读数件里**没有出现这两个词**（尺＝`grep -rn '自红\|判据红' .scratch/wisp/probes/325/r1 docs/evidence/s1/325-runtests-gotmp.md`＝**0 命中**，正控＝同一把尺在 `docs/evidence/s1/` 全域命中 `自红` **17 行**／`判据红` **16 行**，说明这对立说法是本仓的常用词、⛔ 它只是换了说法），它只写了"假红／the ruler malformed its own path"。本腿据现跑判语：
**＝仪器自红**（假红）。四枚判据本身一字未动，同一枚根只差分隔符就 4/4 绿（A/D 两发）；`:182`／`:286` 两枚红句**开头的字面就是 `precondition:`**，即用例自己的前置条件在拒绝一条被写坏的根，`:59`／`:116` 两枚红的是 `cmd /c mklink` 的 **exit status 1 文案**（造具环节），没有一枚落在产品断言上。⇒ ⛔ 把 B 那四红记成 risk 模块的判据红，也⛔ 拿来当"325 修好了什么功能"。

**射程那一格＝缺世界，成立（本腿独立现跑，不抄 r1）**：`logs/v1-ac2-range-probe.txt`
① 本入口唯一的门 `scripts/d22scan.sh:51`＝`sh "$root/tools/d22scan/runtests.sh" -C tools/d22scan ./...` ⇒ 射程只有 `tools/d22scan` 一枚 module；
② `tools/` 里**可执行的** `mklink` 调用＝**0 枚**；`grep -rnF mklink tools` 现量＝**4 行，全在 `428ef21e` 自己新加的注释里**（`runtests.sh:97/:98/:103/:106`）⇒ ⚠ **r1 那句"grep -rn mklink tools/＝0 命中"在 HEAD 上已过期**，结论不变、读数要重述（见 §6 更正①）；`internal/` 里真用 `mklink` 的＝**14 枚文件**；
③ CI 唯一直呼本 runner 打 risk 的那步＝`ci.yml:831`（`-run TestPathResolverJunctionWindows`）；本腿在**混形码**上按那步形状直呼一发 ⇒ rc=**0**、`PASS=1 FAIL=0 SKIP=0, === RUN=1`，发射根逐字混形（`logs/v1-ac2-range-probe.txt.ci-shape.raw`）⇒ 那枚用例吃得得住混形、它绿⛔ 等于坑没；
④ 本腿也**没有**为凑红去改 d22scan 的档位（那是改门、要人工批准）。

## 3. AC#3 票 317 那一发：**成立**（11 包逐包，⛔ 总数交差）

`bash scripts/portable-tests.sh --scope=windows`（head3、干净台面、sherpa PATH 已置）＝**rc=0**，全量 `logs/v1-ac3-scope-windows-head3.txt`（1897 行）。
先按派单把"八包"当**过期计数**处理：现量尺＝`sed -n '329,336p' scripts/portable-tests.sh` 的 `windows` 数组＝**11 枚包**（proc／secret／config／risk／ball／perm／plugin／`cmd/llmrecord`／session／projctx／audio），本腿逐枚报。

① 发射行逐字**原生形**：
`portable-tests.sh: ticket 310 AC#1: GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp\portable-20261016-103259-518 - this run's t.TempDir() family lands here. …`

② **11 包逐包四数**（口径写死：`RUN=`＝`^=== RUN` 行数；`PASS_TOP=`＝`^--- PASS` **只数顶层**；`FAIL_TOP=`＝`^--- FAIL`；`SKIP=`＝`^--- SKIP`；单位＝顶层用例枚数）：

| 包 | RUN | PASS_TOP | FAIL_TOP | SKIP |
|---|---|---|---|---|
| internal/proc | 48 | 40 | 0 | 0 |
| internal/secret | 29 | 21 | 0 | 0 |
| internal/config | 185 | 104 | 0 | 0 |
| internal/risk | 200 | 132 | 0 | 0 |
| internal/ball | 75 | 69 | 0 | 0 |
| internal/perm | 14 | 14 | 0 | 0 |
| internal/plugin | 7 | 7 | 0 | 0 |
| cmd/llmrecord | 2 | 2 | 0 | 0 |
| internal/session | 36 | 19 | 0 | 0 |
| internal/projctx | 10 | 10 | 0 | 0 |
| internal/audio | 56 | 44 | 0 | 0 |
| **合计** | **662** | **462** | **0** | **0** |

归因法与自证＝`logs/v1-ac3-perpackage-four-numbers.txt`：Go 把一个包的 `-v` 输出连续打完再以 `^ok/FAIL` 结果行收尾 ⇒ 按结果行切块；
**切块之和 662/462/0/0 ＝ 脚本自打印的总计行**（`BLOCKS=11`，⛔ 并行打乱）；再拿三枚小包**独立单跑**复尺＝`perm PASS=14`／`llmrecord PASS=2`／`plugin PASS=7`，与切块值逐枚对上。
**11 枚全部与 r1 的 `a3-per-package-four-numbers.txt` 逐枚同值**（含 risk 200/132 那枚，即 `317-v1` 也复认过的同一行）⇒ 跨腿三方同读数。

③ 红名册＝**空**：顶层 `--- FAIL` 行 **0** 枚、包级 `FAIL` 行 **0** 枚（正控＝同一把尺在同一文件里数出 `--- PASS` **462** 枚）；
`Invalid switch` 与 `Parameter format not correct` 两串在这发里各 **0** 命中（正控同上）。
⇒ 与 10-14 那两发 CI 的 `Invalid switch` 名册（账＝A922，那四枚 Sync/116 族）作差＝**新增 0 枚**。
⛔ `--scope=cli`、⛔ `./cmd/wisp` 整包（票 314 死亡区＋本波那档归写手）。**这一发够票 317 AC#3 的门**：三样齐、整档跑完、无缺包。

## 4. AC#4 门禁：**成立**（四枚 rc 自落；两枚台面雷按派单处理）

`logs/v1-ac4-gates.txt` ＋ 五枚捕获：

| 门 | 台面 | rc | 关键行 |
|---|---|---|---|
| `bash -n tools/d22scan/runtests.sh` | head3 | **0** | 0 字节输出＝无诊断 |
| `sh scripts/d22scan.sh` | head3、干净台面 | **0** | `clean - no D22 ban violations` ＋正控 `PASS=35 FAIL=0 SKIP=0, === RUN=77`（与 r1 同值） |
| `sh scripts/check-path-length-budget.sh --with-self-test` | ★仓内只读一发（该门要 git 索引，`git archive` 不带） | **0** | `VERDICT GREEN`、`tracked paths=13967 over-budget=60 covered=60 not-in-roster=0`、`control 2/3 ok` |
| `GOFLAGS= go build ./...` | head3 | **0** | 0 字节输出 |

雷①票 311 rc=2 射程拒绝：**零命中**（`grep -rnF 'shape B: GUARD'` 五枚捕获＝0，正控＝同尺同目录命中 `clean - no D22 ban violations` 2 次）。
雷②`gofumpt`：`which gofumpt` rc=**1**（不在 PATH，r1 这一枚复认），但★**真身在盘上**＝`/d/work/base/gopath/bin/gofumpt.exe`（v0.12.0，`-version` 现跑）——r1 报"四处候选目录零命中真身"是因为它没看 `GOPATH=/d/work/base/gopath/bin` 这一处。本腿按绝对路径跑了那把更严的：`gofumpt -l tools scripts`＝**0 枚**；对照`gofmt -l` 同域也 0 枚（这里两把同值，但本腿交的是更严那把的 0，⛔ 顶包）。本票 diff（`428ef21e --stat`）里 `.go`＝**0 枚**（只一枚 `.sh`），且本腿写面 0 枚 `.go` ⇒ gofumpt 射程为 0，格上写"不适用＋真身已定位"。
台面孔扫描（`walk`/`no tests to run`）命中 2 处，逐枚读过＝**均非故障**：`d22scan` 自检输出里的注释性文本（`independent walk: no index to consult`、`walks nothing`）与严格 runner 自己的摘要 `'[no tests to run]'=0`；`build failed`/`0xc0000135` 零命中。

## 5. AC#5 账：**验毕**（翻框与重述票 317 归编排者）

尺＝`logs/v1-ac5-gotmpdir-inventory.txt`（只跑在导出树上，⛔ 把在飞写腿的字节算成第五处）。
生产面（`scripts tools cmd internal`）里提到 `GOTMPDIR`/`WISP_TEST_GOTMP_ROOT`/`gotmp_run=`/`gotmp_base=` 的行共 **61 行**，逐枚归类后只有两类：

- **当值用（真流向 `export GOTMPDIR`）＝4 枚**：`scripts/portable-tests.sh:155-158`（`gotmp_base`＋`\\portable-`）、`scripts/winsec-tests.sh:69-71`、`scripts/wisp-cli-tests.sh:82-84`、`tools/d22scan/runtests.sh:115-117`；`export GOTMPDIR` 行数现量＝**4**。四枚尾接**全是反斜杠** ⇒ 本票修的是第四处、也是最后一处。
- **只当展示用＝4 枚（有意为之、⛔ 算病）**：`wisp_311_ref_path=${LOCALAPPDATA//\\//}`＝`portable-tests.sh:818`／`winsec-tests.sh:174`／`wisp-cli-tests.sh:186`／`runtests.sh:178`（⚠ 票面写的 `:155` 已被 325-r1 那 23 行注释推到 `:178`，字节未变＝本腿现量复认）。它们只喂票 311 的**音量比较**与拒绝文案，⛔ 进 `GOTMPDIR`。
- 另有一枚同族**记载**（不是赋值位）＝`internal/winsec/openleaf_guard_322_windows_test.go:95-97` 的注释：那族用例用 `filepath.Join` 规避，并写下 `Invalid switch - "Users"` 的实测。

混形尾接 hunt（`}/d22scan-`、`}/winsec-`、`}/cli-`、`}/portable-`、`/wisp-test-gotmp` 五串，`grep -rnF` 覆盖 `scripts tools cmd internal .github`）＝**逐串 0**；
正控＝同尺同域对 `'${LOCALAPPDATA//\\//}'` 命中 **4**（尺看得见行）。`.github/workflows/ci.yml` 里**没有任何** `GOTMPDIR`／`WISP_TEST_GOTMP_ROOT` 赋值行（现量 0）⇒ CI 那两腿吃的是 entry 脚本自己的默认根。

**该不该立新票＝不立"第五处同形"的票**（现量 0，第四处已由本票收口，`scripts/**` 三处已由票 317 收口）；
但本腿具名上报两枚**要编排者拍板、不由验收腿自立**的缺口：
① 这处修复今天**不改变任何一枚门的颜色**（§2 的缺世界），要让它真被门看见，得把某档射程推到 `./internal/risk/`＝**改门、人工批准**；
② `ci.yml:831` 那步的 `-run` 名册只盖住一枚"吃得住混形"的用例 ⇒ 想让 CI 也测到这族，要在 `.github/**`（禁区）加/改一步＝同样人工批准。
⇒ 建议：票 325 按"改形＋反控＋317 那一发"收口，上面两枚另记台账（`Q##` 级），⛔ 混进本票的 AC。

## 6. 对 `325-r1` 证据件的更正（三枚，全部现跑）

1. §2 那句 **`grep -rn mklink tools/`＝0 命中** ⇒ HEAD 上现量＝**4 行、全在它自己新加的注释**（`runtests.sh:97/:98/:103/:106`）。结论（射程不含 risk）不变，但这把尺的读数要重述成"可执行 0 枚／注释 4 枚"，否则下一枚腿会把它当"尺坏了"。
2. §4 那句 **`gofumpt` 四处候选目录零命中真身** ⇒ 真身在 `d:\work\base\gopath\bin\gofumpt.exe`（`GOPATH=/d/work/base/gopath`，v0.12.0）。"本腿写面 0 枚 `.go`"那一半复认。
3. §1 的对照位行号 **`runtests.sh:155`→`:178`** 位移复认（字节未变）；它已在票面 bullet 里自报，本腿独立现跑到同一枚 `:178`。
（其余：五笔提交的写面逐枚核过＝`428ef21e` 只碰 `tools/d22scan/runtests.sh`、`73131d29`/`a916bf6f` 只碰 `probes/325/r1/**`、`c7190b27` 只碰 `docs/evidence/s1/`、`1e53de91` 只碰票尾＋`probes/325/r1/commit-evidence.txt` ⇒ **零越界**；`scripts/**` 与四处对照位一字未动。）

## 7. 没裁的／缺的（具名）

- **CI 换色那一发**：推送权在编排者，本腿同样只有本机导出树读数 ⇒ §3 的 11 包是 dev 主机导出树的，⛔ 等于 `test-windows` 步的读数。
- **`--scope=cli`／`./cmd/wisp` 整包**：派单禁跑（票 314 死亡区＋本波那档归写手 `220-*`），⛔ 算遗漏。
- **AC#5 的"立新票与否"是判语不是动作**：本腿没立票、没动票 317 一字、没翻任何框。
- 我把 `internal/winsec/**` 在飞写腿的字节排除在读数之外（一律以 `git archive HEAD` 的树为准）；工作树里 `220-*` 若已落了新码，那部分⛔ 在本表射程。

## 8. 临时根名册（只建不删）

`C:\Users\swq\AppData\Local\wisp-test-gotmp\`：`portable-20261016-103259-518`（AC#3 整档）／
`d22scan-20261016-103551-922`／`-103554-961`／`-103557-998`（AC#3 的三枚独立小包）／
`d22scan-20261016-103804-307`／`-104116-180`（混形那发）／`-104125-267`（还原那发）／`-104234-1715`（CI 形状那发）／`-104634-390`（改前门）／`-104730-496`（改后门）。
仓外树：`C:/tmp/wisp325v1/{head,head2,head3,mixed,mixed3,sedtest}` ＋ `C:/tmp/wisp325v1/psxroot-103814-270`／`psx2-104118`（纯正斜杠根）／`wisp325v1/psx2-*`。
`mixed3/tools/d22scan/runtests.sh` 停在**已还原＝与提交态同哈希 `ad6f6a8d…4bf6`** 的状态（未删，留着复尺）。
`head2` 是缺 `design/` 的那枚半成品台面，一并留着当"我这把台面也会错"的物证。
