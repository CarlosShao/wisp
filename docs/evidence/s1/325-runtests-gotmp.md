# 票 325 AC#1-AC#5 落地腿读数表（腿＝`325-r1`，实现者本人）

- 起手锚＝`46007856`（派单锚，现验＝本腿 HEAD 的祖先，`git rev-list --count 46007856..HEAD`＝1，
  期间唯一新增＝编排者那笔"存档并放出四枚" `96fc2343`）
- 产码笔＝`428ef21e`（`tools/d22scan/runtests.sh`，一枚文件）｜读数件笔＝本腿下一笔（`probes/325/r1/`，22 枚）
- ★**本表不是裁决**：实现者自陈，判语权归非实现者验收腿；这里只交"哪一发跑在哪枚码上、响没响、逐字读数在哪"
- ⛔ 本腿未勾任何 AC 框、未动 `Status:`／票名／`-done`、未动票 317 原文、未 push

## 0. 台面（先说清每一发跑在哪，因为同批有三枚腿在飞）

共享工作树此刻被 `322-r1b` 写着 `internal/winsec/**`，而 windows 档的射程里就有 winsec ⇒
**一切整树/整档读数一律不在仓里跑**。两枚导出树（`git archive <锚> | tar -x -C <仓外目录>`）：

| 树 | 锚 | 六件套逐枚 `[ -e ]` | `find -name '*.go' | wc -l` | 那一行的形状 |
|---|---|---|---|---|
| `C:/tmp/wisp325r1/pre` | `46007856` | go.mod/go.sum/internal/cmd/tools/scripts 全 OK | **1090**（非空⇒才信它的 0） | 混形（基座 `\\`＋尾接 `/`） |
| `C:/tmp/wisp325r1/post` | `428ef21e` | 同上全 OK | **1090** | 原生形（`\\d22scan-`） |
| `C:/tmp/wisp325r1/revert` | ＝post 的副本 | — | — | 反控时**物理**退回混形，再还原 |

干净台面＝每发都 `env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT`（票 311 的 rc=2 射程拒绝一次没触发：
根始终在 `%LOCALAPPDATA%` 同一卷 C:）。唯一在仓内跑的一发＝长度门（它要 git 索引，`git archive`
不带），只读、名册见 §4，⛔ 整树 Go 读数一律没这么跑。

## 1. AC#1 改形（这一发在未修码上：门绿着，坑没人看见）

现读（起手三把尺，本腿自己跑，⛔ 抄票面）：

- `sed -n '92p' tools/d22scan/runtests.sh | od -c`（pre 树＝锚）＝
  `"${WISP_TEST_GOTMP_ROOT:-"$LOCALAPPDATA\\wisp-test-gotmp"}/d22scan-$(date +%Y%m%d-%H%M%S)-$$"`
  ⇒ 基座 `\\`、尾接 `/`＝**混形**（票面那句在此复核成立）
- 三处已修对照（同尺）＝`scripts/portable-tests.sh:155-156`（`gotmp_base` 折法＋`$gotmp_base\\portable-`）／
  `scripts/winsec-tests.sh:69`／`scripts/wisp-cli-tests.sh:82` ⇒ 三处尾接**都是 `\\`**
- 修后那一行（现 `:115`，注释把它推下 23 行）`od -c` 逐字＝`… w i s p - t e s t - g o t m p " } \ \ d 2 2 s c a n - $ ( d a t e …`
  ⇒ 与三处**同形**（对拉串：`}\\d22scan-$(d` ⇄ `}\\winsec-$(da` ⇄ `}\\cli-$(date`）

判据那一发（`sh scripts/d22scan.sh`，两发同一台面）：

| 发 | 树 | rc | 发射的根（逐字） | 正控制四数 |
|---|---|---|---|---|
| 改前 | pre | **0** clean | `C:\Users\swq\AppData\Local\wisp-test-gotmp/d22scan-20261016-093651-329` | `PASS=35 FAIL=0 SKIP=0, === RUN=77` |
| 改后 | post | **0** clean | `C:\Users\swq\AppData\Local\wisp-test-gotmp\d22scan-20261016-093745-842` | `PASS=35 FAIL=0 SKIP=0, === RUN=77` |

⇒ **这一格在未修码上"不响"**：混形照跑、门照绿。两发四数逐字同形＝改形对门本身零影响，
差别只在分隔符（这就是 §2 要证的"射程遮蔽"）。

护栏（本腿自觉守住的边界）：⛔ 未动 `scripts/**` 一字；⛔ 未动 `wisp_311_ref_path` 四处对照位
（`portable-tests.sh:818`／`winsec-tests.sh:174`／`wisp-cli-tests.sh:186`／`runtests.sh:155`；
最后一处因本笔注释**位移到 `:178`**，字节未变，编排者那把行号尺请按此重锚）。

注释写的不是风格：点名会炸的那族用例（`mklink /J` 把参数内部的正斜杠读成开关）＋两种 cmd 措辞
＋"折回 `/` 等于把活坑放回去"，全文在 `tools/d22scan/runtests.sh:92-114`。

## 2. AC#2 成对反控（这族尺的牙）＋射程那一句

先自证四枚真名存在（⛔ 拿不存在的名字跑＝`no tests to run`，不是读数）：
`GOFLAGS= go test ./internal/risk/ -run "^<名>$" -list "^<名>$"` 逐枚命中 **1/1**（含第五枚
`TestPathResolverJunctionWindows`）＝`logs/a2-93-list-four-names.txt`。

四枚＝`TestSyncAncestorActableLegFailsClosedOnMovedAncestor116`／
`TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116`／
`TestSyncAncestorLegReResolvesThroughJunctionToRealTree`／
`TestSyncAncestorLegNotExemptedStaysFailClosed`。

**这一发在未修码上响**（＝本票最硬的一枚读数，且它不是我构造的正斜杠、是**那一行自己生成的根**）：
pre 树（混形码）`sh tools/d22scan/runtests.sh ./internal/risk/ -run '^(…四名…)$'`
⇒ rc=**1**，`RUN=4 PASS=0 FAIL=4 SKIP=0`，发射行逐字 `GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp/d22scan-20261016-093432-450`
（`logs/s1-mixed-emitted.txt`）。红句四行**逐字**（行号与编排者、与 `317-v1` 全部同行同号）：

1. `syncdirs_ancestor_actable_leg_116_test.go:182: precondition: deepestExistingAncestor("…\\001\\%WISP116ANC\\B%\\new.md") = "…\\001\\%WISP116ANC\\B%", want "…wisp-test-gotmp/d22scan-20261016-093432-450\\TestSync…\\001\\%WISP116ANC\\B%"`
2. `syncdirs_ancestor_actable_leg_116_test.go:286: precondition: the ~ spelling must resolve under the relocated home "C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp/d22scan-20261016-093432-450\\TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116795367349\\002", got "C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp\\d22scan-20261016-093432-450\\TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116795367349\\002\\wisp116-tilde-probe\\new.md"`（★两侧同一条根，一侧 `/d22scan-`、一侧 `\d22scan-`＝分隔符在**同一次运行内**分家）
3. `syncdirs_ancestor_reparse_ticket105_windows_test.go:59: mklink /J C:\Users\swq\AppData\Local\wisp-test-gotmp\d22scan-…\TestSyncAncestorLegReResolvesThroughJunctionToRealTree…\002\proj C:\Users\swq\AppData\Local\wisp-test-gotmp/d22scan-…\TestSync…\001 failed: exit status 1: Parameter format not correct - "d22scan-…\TestSync…\001".`
4. 同上 `:116`（`TestSyncAncestorLegNotExemptedStaysFailClosed`）

⚠ **措辞要说准**：混形（尾接那枚 `/`）下 cmd 报的是 `Parameter format not correct - "d22scan-…"`；
编排者与 `317-v1` 那对**纯正斜杠根**报的是 `Invalid switch - "tmp".`／`Invalid switch - "Users".`。
两种文案＝同一个因（mklink 把参数里的正斜杠读成开关），⛔ 把本腿这发写成 `Invalid switch`。

反控链（在 post 树的副本 `revert` 上**物理**改字节，⛔ `-overlay`——这族尺运行时读盘，overlay 恒绿＝假读数）：

| 步 | 读数 |
|---|---|
| 提交态 sha256 | `ad6f6a8d72cbe2a3fc9b7036774b8e16ddf92bfd64ea6055cbf12dde55164bf6` |
| 退回混形后 | `6319065990c0fa2c374d5bb789d5fbc2702293fe57d5000f4366ca675e853240`＋**突变自证** `mixed-shape-hits=1 native-shape-hits=0`（⛔ 恒真：sed 没落就会留在这行） |
| 那发四枚 | rc=**1**（`logs/a2-revert-mixed-red.txt`） |
| 还原（`git show 428ef21e:tools/d22scan/runtests.sh` 的字节） | sha256 **回到 `ad6f6a8d…`**＝提交态那枚；`mixed-shape-hits=0 native-shape-hits=1` |
| 复跑 | rc=**0**，`PASS=4 FAIL=0 SKIP=0`（`logs/a2-restore-green-four.txt`；改后树同构造同样 4/0＝`logs/a2-post-green-four.txt`） |

**射程里今天真含那四枚吗＝不含（缺世界，具名）**：
- 经由本入口的唯一门＝`scripts/d22scan.sh:51` → `runtests.sh -C tools/d22scan ./...`，射程只有
  `tools/d22scan` 这一枚 module；`grep -rn mklink tools/`＝**0 命中**（全仓 mklink 只在 `internal/` 七枚文件）
  ⇒ 这一档今天**根本不调 `./internal/risk/`**，四枚不在分母里，坑活着、仪器看不见。
- CI 唯一直接喊这枚 runner 打 risk 的那步＝`ci.yml:831`，`-run TestPathResolverJunctionWindows`。
  本腿自己在 pre 树（混形码）跑了那一发的形状：rc=**0**、`PASS=1 FAIL=0 SKIP=0 === RUN=1`
  （`logs/a2-94-ci-shape-mixed-root.txt`）⇒ 那枚用例的 link/target 走 `filepath.Join`（Clean 归一化），
  **吃得下混形**；所以它绿不代表坑没了，只代表 `-run` 的名册没盖到那四枚。
- ⛔ 本腿**没有**为了凑红把 d22scan 的档位改成跑 `./internal/risk/`（那是改门、要人工批准）。
  红是从"把同一枚 runner 按 CI 那步的形状直呼一次"读出来的，档位一字未动。

⇒ 结论：AC#2 的"改前必红"成立（未修码真红），"射程遮蔽"同样成立并具名；票面那句
"今天没炸只因为射程里没有那几条"由本腿的两发读数（混形四红 ⇄ 同根下 CI 那枚用例绿）钉住。

## 3. AC#3 顺带交出票 317 的 AC#3 那一发

`bash scripts/portable-tests.sh --scope=windows`（post 树，干净台面）＝rc=**0**，全量输出
`logs/a3-post-scope-windows.txt`（1887 行）。三样齐：

① 那一行 `GOTMPDIR=` 逐字**原生形**：
`portable-tests.sh: ticket 310 AC#1: GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp\portable-20261016-094356-1108 - this run's t.TempDir() family lands here. Nothing deletes it: …`

② **逐包**四数（口径写清：`RUN=`＝`^=== RUN`；`PASS=`＝`^--- PASS` **只数顶层**，子例不进去；
两把之差＝本发 662−462＝**200** 枚子例行，与派单那句"约 11 枚"不同，那 11 是别的口径，此处以本发实测为准）。
★**派单/票面写"八包"＝过期计数**：现量 windows scope＝**11 包**（`portable-tests.sh:329-336`，
scope 行逐字见日志第 2 行）⇒ 本腿按现量逐枚报，不迁就那个数：

| 包 | RUN | PASS | FAIL | SKIP |
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

（合计行与脚本自己打印的那行逐字一致：`four numbers (all from -v output): === RUN=662 --- PASS=462 --- FAIL=0 --- SKIP=0`；
`internal/risk` 那行与 `317-v1` 的 `RUN=200 PASS=132` 同形＝跨腿复认。）

③ 红名册＝**空**（`--- FAIL` 顶层 0 枚、`FAIL` 包行 0 枚）。作差基准＝10-14 那两发 CI
（`38508311376` vs `38484004521`，账＝台账 A922）里带 `Invalid switch`／`want "C:/…"` 的名册＝
`test-windows` 5 增中的那**四枚** Sync/116 族（第五枚 `TestSessionListByWorkspace` 文案不是这句，另账）。
⇒ 新增红 **0 枚**；那四枚今天在这发里全绿（＝票 317 那三处的效果，账记 317，⛔ 记本腿）。

⛔ `--scope=cli` 未跑、⛔ `./cmd/wisp` 整包未碰（票 314 死亡区＋假分母，派单禁）。

## 4. AC#4 门禁（逐枚 rc，台面具名）

| 门 | 台面 | rc | 关键行 |
|---|---|---|---|
| `bash -n tools/d22scan/runtests.sh` | post 树 | **0** | — |
| `sh scripts/d22scan.sh` | post 树（干净台面） | **0** | `clean - no D22 ban violations`＋正控制 `PASS=35 FAIL=0 SKIP=0 === RUN=77` |
| `sh scripts/check-path-length-budget.sh --with-self-test` | ★**共享工作树**（只读；该门要 git 索引，`git archive` 不带） | **0** | `VERDICT GREEN - every over-budget tracked path is rostered by name …`＋`over-budget=60 covered=60 not-in-roster=0` |
| `GOFLAGS= go build ./...` | post 树 | **0** | 输出 0 字节（＝无诊断，非"没跑"；四数见 `a4-gates-stdout.txt`） |

读数纪律复核：全部读数里 `=== RUN` 均非 0；⛔ 出现 `walk`／`build failed`／`no tests to run`／
`0xc0000135`／`RUN=0` 的台面孔（本腿一次没触发；票 311 的 rc=2 射程拒绝也没触发——根始终在 C:）。

`gofumpt` 这笔：`which gofumpt`＝**rc=1**（本机 PATH 里没有），四处候选目录
（`~/go/bin` 空／`~/go/pkg/mod`／`$LOCALAPPDATA/go/bin`／`/d/work/base/go/bin`）零命中真身。
本腿写面＝**0 枚 `.go`**（只一枚 `.sh`）⇒ gofumpt 射程为 0，这里据实写"不适用＋找不到真身"，
⛔ 拿 `gofmt` 的 0 冒充它的 0，也⛔ 写"格式没问题"。

## 5. AC#5 账（本腿只动票 325 一处）

- 票 317 面上那句"同形、今天未爆"→"混形，射程遮蔽中的活坑"＝**由编排者重述**，本腿⛔ 未改票 317 一字、⛔ 未碰其 AC 框。
- 本腿在票 325 末尾**只追加一条** bullet（含四枚真名、两种 cmd 措辞之别、缺世界结论、"八包→11 包"的过期计数更正）。
- 出身两枚（`581dd252`／`76b9e214`）＝编排者自己的手，账记编排者，⛔ 记任何腿（派单已言明，本表复认）。
- 本表与读数件里凡"某格成立"字样都是**实现者自陈**，翻框权在非实现者验收腿。

## 6. 没做完的（具名，⛔ 含糊话）

1. **CI windows 档换色的那一发没取**：本票的 windows scope 读数全是**本机导出树**的；
   推送后要看的＝`test-windows` 步的 `GOTMPDIR=` 行是否 `C:\Users\RUNNER~1\…\portable-…` 原生形
   ＋`ci.yml:831` 那步是否仍绿。缺的那一发＝CI 日志，只有编排者推送后能取（⛔ 本腿 push 权）。
2. **`--scope=cli`／`./cmd/wisp` 整包＝派单禁跑的**，不是"忘了跑"；如果收口需要那两发的形状，
   得先解票 314 那条死亡区（本腿没有凭据可给）。
3. **纯基座正斜杠那一发本腿没重跑**：只有编排者（`probes/orch/20261016/logs/10-nat.txt`／`10-psx.txt`）
   与 `317-v1`（M1）的两份；本腿交的是**混形自己生成的根**那一发（更贴本票射程）。要第三种形状请验收腿自取。

## 7. 临时根名册（只建不删）

`C:\Users\swq\AppData\Local\wisp-test-gotmp\d22scan-20261016-093432-450`／`-093651-329`／
`-093745-842`／`-094307-1201`／`-094313-1243`／`-094931-1982`／`portable-20261016-094356-1108`
＋仓外导出树三枚 `C:/tmp/wisp325r1/{pre,post,revert}`；
`revert/tools/d22scan/runtests.sh` 留在"已还原＝与提交态同哈希"的状态（未删，留着复尺）。
