# 228-r3c 交件：`R-228-3b-1` 那两发突变在仓外导出树上**串行**重做成彼此独立的两发读数

腿＝`228-r3c`（票 228 突变复跑，只仓外导出树）。派单存档＝`.scratch/wisp/probes/orch/20261018/txt/228r3c-dispatch.txt`（56 行，已 Read 全文）；
另读 `.scratch/wisp/probes/228/r3b/report.md`（59 行）＋票 228 `AC#11`（`:91-102`＋`:97-100` 拆钉段）＋台账 `A1019` 第③节。
起手锚（现跑 `git log -1`）＝`c11163c896b09e7b6863d347dadbcb65702585cc  Sun Oct 18 10:02:56 2026 +0800`（branch dev）。
逐字件＝`start.md`（起手三把尺＋落点更正＋短名尺）＋`logs/**` 10 枚（逐枚 `wc -c`>0，0 字节枚数＝**0**）。

## 0. ★落点更正（机主 10-18 10:1x 新铁律覆盖派单台面配方）

⛔ 盘符根零污染 ⇒ 导出树／`GOTMPDIR`／TMP 全在 `D:\tmp\wisp\20261018\228r3c\{tree,gotmp,tdst}`；`228-r3b` 的 `D:\wisp228r3b` 形状⛔ 复用、⛔ 删。
短名尺现跑三把（逐字在 `start.md` §2）：`fsutil 8dot3name query D:`＝`Error 5 Access denied`（本腿非管理员⇒尺读不到，⛔ 提权）；
`dir /X` 对新落点现建的长名目录吐出 `PROBE-~1`＝**8.3 短名在该处确实生成**；`cygpath -d …\tree`＝原样返回（每段 ≤8 字符）⇒ 不需"专用树根"例外，已回 lead。

## 1. pristine 基准（现算，⛔ 抄旧数）

尺＝起手锚上那枚 blob 两把同看：`git rev-parse c11163c8:internal/proc/boot_windows.go`＝blob **`531060a55c080ebd535b426b8c60dcd9011785fc`**；
`git cat-file blob … | sha256sum` ⇄ 导出树物理字节 sha256 ＝逐字同串 **`3e2e83fab184ce91640ad7dbc06f8dd116baffe325bd47aa46af2527aad73fe7`**（`bytes=12916 lines=313`）。
⚠ 两笔具名：① 这串与 `A1019`/`228-r3b` 引用的旧串相同＝**我这发导的是同一枚 blob**（该件末笔 `f4d4fda9` 早于本锚、之后无新笔），数是我现跑的，相同属构造使然；
② 我第一把尺写成 `git show … | Out-File` 读出 `c3890203…`＝**PowerShell 把 LF 换成 CRLF＝我那把尺坏**，按"两把尺矛盾先假定新那把坏"换 `git cat-file` 才对上（该件 `.gitattributes`＝`text eol=lf`，`git ls-files --eol`＝`i/lf w/lf`）。

## 2. 台面（八件套＋docs/＋手工三件）

`git archive c11163c8 go.mod go.sum internal cmd tools scripts frontend models docs | tar -x`，＋`build/*.dll`(3)／`third_party/sherpa-onnx/*.dll`(3)／`frontend/dist`。
取出后 `.go` 数＝**760**（尺＝`find $TREE -name '*.go' -type f | wc -l`，证非空，⛔ 读数 0 当"没件"）；`docs/PLAN.md` 在位。
全程 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（防 `0xc0000135`）；⛔ 在共享工作树跑 `go test`／`go build`；⛔ `-overlay`；⛔ `go test ./...`；⛔ 跑 `cmd/wisp` 整包（本腿只定向 `-run`）。

## 3. pristine 正控一发（未突变；定向四枚 × 两把临时根台面成对）

`-count=1 -timeout 30m -v -run '<四枚 \$ 锚定形>'` ⇒ **c `RUN=4 PASS_TOP=4 FAIL_TOP=0 SKIP=0` rc=0 ⇄ d 同 `4/4/0/0` rc=0**
（件＝`logs/estate-pristine-c.log`／`…-d.log`；四数不混用、不相减）。逐枚点名全 PASS：
`TestAC11StopBitIsBuiltAtBoot`(ⓐ-first)／`TestAC11StopRequestDuringTheBootWindowIsNotSwallowed`(ⓐ空窗)／
`TestAC11TrayExitExecutorIsMountedAtTheBallHost`(派单额外挂载钉)／`TestAC11TrayExitTriggerWalksTheTenSteps`(ⓑ-1，逐字 `trail=scheduler-close > … > exit`、`reason="stop request (tray Exit item selected in this process)"`)。
⇒ ⛔ 恒红当牙：这台子先证明四枚本底是绿的，突变才可能是牙。

## 4. ★两发突变（串行，各发内还原，发间只认哈希对撞）

病根＝`A1019` 第③节 `R-228-3b-1`：上一程跨 shell 的 `trap`/`finally` 竞态 ⇒ 导出台一度同时残留 M1＋M2 ⇒ 那两发红句取于残留窗口＝彼此不独立。
本腿的改法＝**一发＝一次 shell 调用**（`mutate-shot.sh`，⛔ 用 trap），脚本内顺序执行：验 pristine → 落突变 → 跑点名用例 → `git archive` 还原 → **还原后再验一次哈希** → 复跑同用例。
★凡 echo 的"已还原"不当证据；哈希对不上即 `exit 7`＝下一发⛔ 开跑（串行靠这道闸，不靠人记得）。

| | 内容锚（整行 token，⛔ 绝对行号） | 突变形 | 点名用例 | 三段哈希（sha256 ⇄ blob） | rc |
|---|---|---|---|---|---|
| **M1** | `stopSignal := rt.StopSignal()`（`RunEventLoop` 体内，`case <-stopSignal:` 唯一喂料者） | → `(<-chan struct{})(nil)`（nil 通道永不 ready） | `TestAC11TrayExitTriggerWalksTheTenSteps`(ⓑ-1) | 前 `3e2e83fa…`/`531060a5…`(12916) → 后 `3856c24aa1c1ee25b1a351d99fb91ec8ba7bb475cabeea793796944cfcd57fb3`/`54e0b8801d025be8d35659eeeb56f9e8c980f832`(12969) → 还原 `3e2e83fa…`/`531060a5…` | 突变 `rc=1` 红；复跑 `rc=0` 绿 |
| **M2** | `rt.stopReq = newStopRequest()`（`Boot` 体内唯一建位处） | 整行注释掉（`StopSignal()` 答 nil） | `TestAC11StopBitIsBuiltAtBoot`(ⓐ-first) | 前 `3e2e83fa…`/`531060a5…`(12916) → 后 `2708ff75b40a3219344c6e06b4740bf9f5526d3bd568195375b6b2c0f2c2225c`/`e62b05be126939e267a18a010e448a8ce239796e`(12963) → 还原 `3e2e83fa…`/`531060a5…` | 突变 `rc=1` 红；复跑 `rc=0` 绿 |

两发红句逐字（各取于**本发单发窗口**，件＝`logs/mutation-M1-b1.log`／`logs/mutation-M2-stopbit.log`）：
- M1 → `resident_tray_stop_228_windows_test.go:310: AC#11 RED: RunEventLoop never returned, although a stop request is recorded on this runtime (waited 5s). internal/proc's loop select is the only thing in this package that turns a recorded request into a return, so this is the clause that reads that check.` ＋ `--- FAIL: TestAC11TrayExitTriggerWalksTheTenSteps (5.01s)`
- M2 → `resident_tray_stop_228_windows_test.go:137: AC#11 RED: a runtime built through proc.Boot has no stop-request bit (StopSignal() answers nil), so the tray's Exit item has nothing to record into before the event loop exists` ＋ `--- FAIL: TestAC11StopBitIsBuiltAtBoot (0.00s)`

**两发互不污染自证**（每发一把"另一枚落点未动"的尺，整文件哈希＋该行逐字节哈希，⛔ 只 diff 一行）：
- M1 期间 `rt.stopReq = newStopRequest()` 行 sha256 前＝后＝`6268b19405ba7584d966c2a21c226185c50ef7a0a902be0b7a4f98f418e3afb8`、hits=1；
- M2 期间 `stopSignal := rt.StopSignal()` 行 sha256 前＝后＝`d0a029bc8ab8b2c274776c6914400182681a8c4717e231b91c29d4a185859950`、hits=1。
⇒ 两发各自动一枚落点、另一枚全程未动＝**独立的两发读数**（成对两发落点不同✓）。

**收干净复验**（两发之后，同一棵树，`logs/estate-postboth-c.log`／`…-d.log`）＝定向四枚 **c `4/4/0/0` rc=0 ⇄ d `4/4/0/0` rc=0**；
终态对撞（`logs/gate2-residue.log` [A]）＝`cmp-rc=0`（anchor blob ⇄ 导出树逐字节 IDENTICAL）、`RESIDUE-MUTATION-in-export-tree=0`（正控先跑：同尺对已知 token 命中 1、对我脚本里的标签命中 3 ⇒ 这把尺会响）。

## 5. 门禁（逐枚 rc；只在共享工作树跑共享尺）

`sh scripts/d22scan.sh`＝**rc=0 clean**（ban#1-5 `internal/`=241＋`cmd/`=42、ban#6 frontend=116、ban#7 internal/tools=25、ban#8 design=39/frontend=116/internal=586/cmd=143）。
`sh scripts/check-path-length-budget.sh`＝**rc=0 VERDICT GREEN**（tracked 15255、over-budget=60、roster=60、worst 229）。
`go vet ./cmd/wisp/`＝在**导出树**跑 **rc=0**（输出 0 字节）；⛔ 在共享工作树跑＝派单禁在共享树跑 go 编译面，具名不在此处。
`gofumpt -l`＝**本腿零 `.go` 跟踪件 ⇒ 空列表**（`gofumpt-tracked-go-count=0`），⛔ 当"过了"（无凭据，同 `228-r3b` 口径）。
本机存量具名（⛔ 我动）：工作树 6 枚既有件因 CRLF／注释对齐被标脏＝存量（台账 `A1024` 第③节在册）。

## 6. 写面与纪律自证

笔次＝读数笔 5 枚已入库（严格前缀尺 `git log --all --grep='^票 228 突变复跑腿 228-r3c'`，逐枚＝`953f8966` 起手／`c544bb3d` pristine 正控／
`52f24441` M1／`266c6463` M2／`f250b532` post-both 复验）＋本笔（report，＝第 6 枚）＋末笔（票 228 尾一条 bullet，＝第 7 枚）。
枚数尺按构造数不到"写它的那一笔"，故这里报"已入库 5 ＋ 本笔 ＋ 末笔"，⛔ 自吹成 7。
路径全在 `.scratch/wisp/probes/228/r3c/**` ＋票 228 尾一条；**`.go` 命中＝0**（`logs/gate2-residue.log` [D] 段逐枚 `--name-only` 尺）；
⛔ 落 `docs/evidence/s1/`（那是裁决者的面）；⛔ 翻框／改票名／`Status:`／`-done`／push ⇒ **票 228 未勾仍 11 枚、`AC#11` 十一格里一枚没翻**（尺＝`grep -c '^[[:space:]]*- \[ \]'`＝11，现跑）。
`cmd/wisp/**`／`internal/panel/**` 零写腿自证：`git status --porcelain -- cmd internal tools scripts go.mod go.sum frontend docs` 现量只吐别腿在飞字节（`internal/winsec/…`、`internal/session/manage.go`＝`322-r2`／`166-r1` 写腿），⛔ 我一枚。

## 7. 枚数尺就地更正（⛔ 改写票面，按"追加＋具名"）

票 228 尾那条 bullet（笔 `7b1b22d8`）里写 `logs/`＝**10 枚**，交件收笔时现量＝**11 枚**（`log-count=11`、0 字节枚数＝0）：
多出的那枚＝`logs/closing-final.txt`＝**写它的那一笔晚于票尾 bullet**，枚数尺按构造数不到自己（同 `A1020` §④ 那一形，⛔ 判腿谎报）。
票面一字不改（⛔ 删行／⛔ 改写），以本节为准；该件里另具名两把本腿写坏的尺：
① `for f in $T/logs/*` 未加引号 ⇒ 被含空格路径词分割＝吐 "No such file"，改正后逐枚 `wc -c` 见该件尾段；
② `df -BG /d | tail -1 | awk '{print $2}'` 读 501G＝取到 Size 列，与 PowerShell `(Get-PSDrive D).Free`＝**198.3 GB free** 矛盾 ⇒ 按"两把尺矛盾先假定新那把坏"，以 PowerShell 那把为准（与 `A1024` 第④节"170→199GB"同向），本腿⛔ 据那枚假读数下任何结论。
终态尺（同件）＝导出树 `sha256=3e2e83fa…`／`blob1=531060a5…`（与每发 log 的还原哈希同串、独立现跑）、`RESIDUE=0`（正控 token 命中 1＝这把尺会响）、
票 228 `unchecked=11 / checked=5`、本腿 `.go` 命中＝**0**、`D:/wisp228r3b exists=NO`（＝机主 10-18 清理已过账 `A1024` 第④节；本腿⛔ 建、⛔ 删任何盘符根件，全程只写 `D:\tmp\wisp\20261018\228r3c\**`）。

## 7. 必答 (i)–(iv)

**(i) `R-228-3b-1` 是否已被本腿重做成两发独立读数＝成立。**
凭据（逐字见 §4 表与两枚 log）：M1 前/后/还原三段＝`3e2e83fa…/531060a5…` → `3856c24a…/54e0b880…` → `3e2e83fa…/531060a5…`；
M2 三段＝`3e2e83fa…/531060a5…` → `2708ff75…/e62b05be…` → `3e2e83fa…/531060a5…`；
且每发另一枚落点行哈希前后同串（`6268b19405ba…`／`d0a029bc8ab8…`）⇒ 两发红句各取于"只有本发落点被改"的窗口＝**彼此独立**。

**(ii) 两发是否都"有牙"＝成立。** M1 ⇒ ⓑ-1 `:310` 点名红（`rc=1`）＋还原后同用例 `--- PASS … (0.01s)` rc=0；
M2 ⇒ ⓐ-first `:137` 点名红（`rc=1`）＋还原后 `--- PASS … (0.00s)` rc=0。pristine 两台面底绿（§3 `4/4/0/0`×2）⇒ 红只可能来自那一发落点。

**(iii) `228-r3b` §4 那句"两枚未做反形⇒不盖章"（挂载钉 `TestAC11TrayExitExecutorIsMountedAtTheBallHost`、空窗那枚 `TestAC11StopRequestDuringTheBootWindowIsNotSwallowed`）本腿有没有补＝没补。**
派单只授权两发（M1／M2），⛔ 我顺手盖章：这两枚在 `A1019`/`228-r3b` 与本腿口径下仍是**零反形读数**＝有绿（§3、§4 收干净复验）但**无牙凭据**。要钉它们有牙须另发两枚突变（挂载钉那枚的落点在 `cmd/wisp/resident_windows.go` 装配处，⛔ 本波 `cmd/wisp/**` 零写腿、且那是产码面）。

**(iv) 有没有任何一发落在残留窗口里取的读数＝没有。**
两发＝两次独立 shell 调用，发内还原＋发间哈希对撞；M2 的开闸第一把尺（`[1] 本发突变前哈希`）读回 pristine `3e2e83fa…/531060a5…` ⇒ M1 未残留；
两发之后终态 `cmp-rc=0`＋`RESIDUE=0`。⚠ 具名一处本腿尺伤（⛔ 计入判据）：`gate-and-residue.log` [6] 段因脚本 `cd` 到导出树后未 `cd` 回仓 ⇒ `git cat-file` 撞 `not a git repository`＝**我那把尺写坏**，已在 `gate2-residue.log` [A] 重跑并逐字改正（`cmp-rc=0`）。另一处＝`run-targeted.sh` 第一形让 `grep` 读写同一件 ⇒ `input file is also the output`＝台面故障，已改正后重跑，入库件用的是改正后的尺。
