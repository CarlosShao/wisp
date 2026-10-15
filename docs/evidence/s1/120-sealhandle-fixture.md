# 120-r3 交件读数表（票 120 残余缺口 1：给 `sealHandle` 那半的 inline 复验补"能红"的具名台件）

件＝`.scratch/wisp/probes/120/r3/`（23 枚读数＋`tooth.sh`／`gates.sh` 两枚脚本）；写面＝`internal/winsec/**` 测试件三枚。
**本件是实现者自报读数，⛔ 判语**（裁决归非实现者；⛔ 翻 AC 框、⛔ 改票名/Status/正文、⛔ 改产码）。
起手锚 `HEAD=e57dca98`（派单 §0）；交件锚 `e1973191`（编排者在飞期间又落了 4 枚 `.scratch/**` 存档件，
`internal/winsec/winsec_other.go` 全程 `sha256=137442cd…` 未变）。Shell＝Git Bash，容器＝`golang:1.27`（docker 29.6.2）。

## 0. 一句话

缺口 1 的**两形都造出来了**，各自都有"单独去这一支的牙 ⇒ 只有它红"的反控：
`winsec_other.go:269`（名字变成链接）与 `:272`（名字指向另一枚真对象）今天**各有一枚确定性台件**
走得到，且**盲跑读数证明这世界不加任何接缝同步也天然发生**（2000 发里 25 发落 `:272`、11 发落 `:269`）。
未去牙整包 **rc=0 / RUN=62 / PASS=62 / FAIL=0 / SKIP=0**；去牙一发 **FAIL=2**，且**变色只发生在我这枚用例上**（逐名比对）。

## 1. 两形各自的世界（AC#1′）

| 形 | 世界 | 造法 | 走到的那支 |
|---|---|---|---|
| 名字指向另一枚真对象 | fd 开在真对象 A 上 ⇒ 用 `rename(2)` 把**第二枚真对象 B** 换到名字上，全程零链接 | 确定性，⛔ 接缝、⛔ 赛跑者：`pinPlacement` 按名钉到的就是 B，`sameObjectPinned` 因此**满足**，只有 `sealHandle` 自己的 `os.SameFile(pinned, now)` 能拒 | `:272` |
| 名字本身变成链接 | fd 开在 A ⇒ **pinned seal 跑完之后**名字被换成指向同目录兄弟对象的**相对链接** | 只能落在窗口里：经包自己声明的接缝变量 `applyDescriptor`（`winsec_other.go:14-23`，`private_fail_test.go` 为 AC#5 早换过它）**先调真的 `applyDescriptorPOSIX`、保留其判语**，判语为 nil 之后才动手 | `:269` |

★**派单 §1 那句"⚠ 必须 Root 会跟的形（相对兄弟链接）"这一条要更正为读数**（⛔ 我按它做、也⛔ 当它成立）：
**相对兄弟链接若种在调用之前，`sealHandle` 仍走不到 `:269`**——`Root` 会跟它，于是 `pinPlacement` 按名钉到 B、
`sameObjectPinned` 在 `:188` 的链接那一支当场拒（`applyDescriptor` 非 nil ⇒ `:257` 就返回），
这**正是缺口成因**。要逼出 `:269` 只有一条路：链接在 `applyDescriptor` 判语出来**之后**才出现。
这一形我没按派单的字面去"预先种相对链接"，而是**当场复现了一次成因**并具名成用例
`TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3`（见 §2 枚 3），读数行逐字：
`leg that spoke = pinPlacement sameObjectPinned link face (winsec_other.go:188)`。

## 2. 名册：新用例 7 枚（顶层，全 `TestAC1POSIX…120r3`；现量⛔ 抄票面）

| # | 用例（`internal/winsec/`） | 这一格量什么 | 未去牙的读数行（逐字，取自 `120r3-run1-teeth-on-body.log`） |
|---|---|---|---|
| 1 | `TestAC1POSIXSealHandleNameCarriesAnotherRealObject120r3`（`sealhandle_identity_120_other_test.go`） | `:272` 真对象形 | `:44 AC#1' world: fd and name = two different objects; name kind=file mode=-rw-rw-rw-` → `:48 leg that spoke = sealHandle inline SameFile face (winsec_other.go:272)` → `:53 AC#1' read-back: name's object=-rw-------, caller's object=-rw-rw-rw-` |
| 2 | `TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3`（同上） | `:269` 链接形（接缝开窗） | `:86 leg that spoke = sealHandle inline link face (winsec_other.go:269)` → `:91 read-back: name kind=link, caller's object=-rw-------, witness=-rw-rw-rw-` |
| 3 | `TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3`（同上） | **缺口成因的机械复现**：预种链接由 walk 回答，走不到 `:269` | `:120 world: the name carries a link at a relative sibling; the fd and the name hold the same object` → `:124 leg that spoke = pinPlacement sameObjectPinned link face (winsec_other.go:188)` |
| 4 | `TestAC1POSIXSealHandlePlainObjectStillSeals120r3`（同上） | 反方向：诚实路径必须仍 seal 成功（防"什么都拒"） | 无日志行＝`--- PASS`（断言是 `err==nil` 且名字上对象落到 `-rw-------`、witness 仍 `-rw-rw-rw-`） |
| 5 | `TestAC1POSIXPrivateFileRefusesWhenNameCarriesAnotherRealObject120r3`（`sealhandle_privateroute_120_other_test.go`） | `:272` 走**产码路由** `PrivateFile`（真 open＋真 unlink-on-refusal） | `:85 leg that spoke = sealHandle inline SameFile face (winsec_other.go:272)` → `:89 tree after the refusal: 1 readable entries, none carrying the payload` |
| 6 | `TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3`（同上） | `:269` 走产码路由；且拒后的 unlink 只删链接、不碰指向的对象 | `:112 leg that spoke = sealHandle inline link face (winsec_other.go:269)` → `:124 tree after the refusal: 1 readable entries, none carrying the payload` |
| 7 | `TestAC1POSIXPrivateFileWindowIsRacedWithoutTheSeam120r3`（同上） | **盲跑读数**：⛔ 接缝帮忙，赛跑者自己翻名字；分布只上报（120-v1 §2 裁③），硬判据只留"外站对象永不被收窄"＋raced 载荷不落进它 | `:209 sealHandle inline SameFile face (winsec_other.go:272) = 25`／`:209 sealHandle inline link face (winsec_other.go:269) = 11`／`floor=1781`／`walk-link=170`／`walk-SameFile=11`／`returned nil=1`／`other=1`（`:212` 逐字 `openat caller-object-120r3.bin: too many levels of symbolic links`，⛔ 判、只上报） |

共用夹具与判据在第三枚 `sealhandle_window_120_other_test.go`（抬头一次写完成因／两形／接缝诚实性／反控计划）。
名册现量：`-run 'TestAC[1-4]POSIX'` 顶层 **27 枚**（＝票里那族原有 20 枚＋本腿 7 枚），容器
**RUN=38 / PASS=38 / FAIL=0 / SKIP=0 rc=0**（`120r3-roster-ac1234.log`；20 那枚旧口径已被 108/119/125 扩过，
本腿⛔ 改票面措写）；只跑本腿 7 枚：`-run '120r3$'` **RUN=7 / PASS=7 / 0 / 0 rc=0**（`120r3-roster-own7.log`）。

## 3. 反控：**单独**去 `sealHandle` 这一支的牙（两发各去一支；⛔ 拿 `sameObjectPinned` 那条腿冒充）

台面与成对：同一快照目录 `/d/tmp/probe120r3/fixed`（`git archive HEAD`＝`6eb3d2e8`，先数 `.go=1076`、
`internal/winsec/*.go=36` 证非空）、同一镜像 `golang:1.27`、同一 `-count=1 -v ./internal/winsec/`、
同一 temp（镜像自己的 overlay `/tmp`＝C 卷那台，容器内⛔ GOTMPDIR 指宿主机——见 §6）。
脚本 `tooth.sh`；牙的**物理改写＋sha256 还原**，`.bak` 两份都留底未删：
clean `137442cdecbb77701c7029c42960593016a152fc965679613314b764c13e2313`
→ tooth-269 `e81a8bf7a4478490c464fd9142c4c0455ab4a921230df8fed74647b5814739c5` → 还原 `137442cd…` **RESTORE MATCHES**
→ tooth-272 `d859e67263508ca702127ca9d62df5b2af81f9fd50f848a848f96cd3d246faa5` → 还原 `137442cd…` **RESTORE MATCHES**。
★突变⛔ 用行号当锚（同文件 `sameObjectPinned` 在 `:188`/`:191` 有**一对同形判据**）：`tooth.sh` 先由
`func sealHandle` 的**函数体范围**里现找行，再落 `LANDS` 自证，并硬校 `false &&` 在全文件**恰好 1 处**——
逐字：`LANDS: tooth 269 targets line 269 …, inside func sealHandle (256..276), before: 	 if now.Mode()&os.ModeSymlink != 0 {`
／`LANDS: after: 	 if false && now.Mode()&os.ModeSymlink != 0 {`／`LANDS: isolated - exactly 1 line in the file carries false && (grep -n: 269)`。
仓内工作树全程未碰：`git status --porcelain -- internal/winsec` 只列我自己的三枚测试件，
`repo winsec_other.go sha256 = 137442cd…`（= clean，见 `120r3-control-driver.log` 末两行）。

### 3.1 去牙那发：红句逐字（文件:行＋文案）＋其余零变色

**tooth-269（`:269` 链接那一支摘掉）＝ `rc=1 / RUN=62 / PASS=60 / FAIL=2 / SKIP=0`，两枚红都是本腿的链接形：**
```
--- FAIL: TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3 (0.00s)
    sealhandle_identity_120_other_test.go:86: AC#1 RED name-becomes-a-sibling-link-after-the-pinned-seal: refusal came from sealHandle inline SameFile face (winsec_other.go:272), want the leg carrying "the name carries a link, so this descriptor is not the object the caller named"; verbatim: winsec: path is not provably resolved, refusing to seal: /tmp/TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3392843619/001/caller-object-120r3.bin: pinned object at /tmp/…/caller-object-120r3.bin is not what the name carries: this descriptor is not the object the name carries now, so the bytes would not land in what was just sealed (winsec: cannot apply a private security descriptor)
--- FAIL: TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3 (0.00s)
    sealhandle_privateroute_120_other_test.go:112: AC#1 RED private-file-route-name-becomes-a-sibling-link: refusal came from sealHandle inline SameFile face (winsec_other.go:272), want the leg carrying "the name carries a link, so this descriptor is not the object the caller named"; verbatim: winsec: path is not provably resolved, refusing to seal: /tmp/TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal3530650110/001/caller-object-120r3.bin: …（同一形状，路由那份）
```
★这一发的读数不止"有人红了"：**去 `:269` 之后这形仍然被拒，但拒因换成了 `:272` 那支**（`os.Lstat` 拿到链接的
stat，与 fd 上对象必不同 ⇒ SameFile 假）。⇒ 判据必须问"**是哪一腿说的**"，⛔ 只问"有没有拒"；
用例正是这么写的（`assertLegSpoke120r3` 带腿名册），所以它红、而 `:272` 那两枚用例全绿（未误伤）。
盲跑读数同发留下同一枚指纹：`sealHandle inline link face (winsec_other.go:269)` **这一行整行消失**（0 发命中），
`SameFile face` 升到 32。

**tooth-272（`:272` 真对象那一支摘掉）＝ `rc=1 / RUN=62 / PASS=60 / FAIL=2 / SKIP=0`，两枚红都是本腿的真对象形：**
```
--- FAIL: TestAC1POSIXSealHandleNameCarriesAnotherRealObject120r3 (0.00s)
    sealhandle_identity_120_other_test.go:48: AC#1 RED name-carries-another-real-object: the call returned nil, so the leg carrying "this descriptor is not the object the name carries now, so the bytes would not land in what was just sealed" never spoke - this is the shape of that leg with its teeth taken out
    sealhandle_identity_120_other_test.go:60: AC#1 RED: the caller's own object was narrowed to -rw------- while the refusal says the bytes would not land in what was just sealed - the face at winsec_other.go:272 is not speaking
--- FAIL: TestAC1POSIXPrivateFileRefusesWhenNameCarriesAnotherRealObject120r3 (0.00s)
    sealhandle_privateroute_120_other_test.go:85: AC#1 RED private-file-route-name-carries-another-real-object: the call returned nil, so the leg carrying "this descriptor is not the object the name carries now, so the bytes would not land in what was just sealed" never spoke - this is the shape of that leg with its teeth taken out
    sealhandle_privateroute_120_other_test.go:87: AC#1 RED: the refused write left the name carrying a file, want it unlinked; err=<nil>
```
`:60` 那一行是**盘上事实**而不是错误文案的复述：没牙时 `chmodPinned(f, …)` 直接落在 fd 的对象 A 上
（`caller's object=-rw-------`，而未去牙那发它是 `-rw-rw-rw-`）＝"seal 的是名字现在携带的对象、写的是另一枚"
这一形真的会静默发生。**路由那发的 `err=<nil>` 更硬**：产码 `PrivateFile` 会**报成功**并留下一个
不是它写的对象在名字上（`tree after the refusal: 2 readable entries`，比未去牙那发多一枚）。
同发指纹：盲跑里 `sealHandle inline SameFile face` 整行消失、`link face` 仍在（6 发）。

**零变色的尺**（⛔ 只看四数）：把 `--- PASS/FAIL/SKIP` 连名字逐枚排序比对——
`diff 120r3-teeth-on-verdict-by-name.txt 120r3-tooth-269-verdict-by-name.txt` 与 `…272…` 各自**只有 4 行**
（同一枚用例 PASS→FAIL，两枚用例各一对），名册本身（`*.names`）三发 `diff` **全等**。
即：其余 60 枚——含票 113/118/119/125 那族全部既有用例、我的成因枚 3 与反方向枚 4、以及盲跑枚 7——**零变色**。

### 3.2 未去牙那发（全绿）

`120r3-run1-teeth-on.log`：`go-test rc=0`／`RUN=62`／顶层 `PASS=40 FAIL=0 SKIP=0`／全行 `PASS=62 FAIL=0 SKIP=0`。
交件锚再确认一发（`e1973191`，编排者的 4 枚 `.scratch/**` 存档件之后，Go 字节未变）：
`120r3-run4-confirm-at-final-anchor.log` 同数 **rc=0 / RUN=62 / PASS=62 / FAIL=0 / SKIP=0**。
参照 120-v1 §2 的那发：同一支牙**去两支**跑整包是 `55/55/0/0 rc=0`（零枚变色＝今天没这世界）；
补完之后同一牙形必须出现红——本件两发各出 2 枚红，都是本腿用例。

## 4. 门禁（逐把带 `rc=`，日志在 `120r3-gates-driver.log` 与同名分件）

| 尺 | 口径 | rc | 读数 |
|---|---|---|---|
| `go build ./...` | 宿主 windows（`winsec_other.go` 不参编译，⛔ 当"已验"） | **0** | `120r3-build.log` |
| `GOOS=linux GOARCH=amd64 go vet ./internal/winsec/` | 真读 `!windows` 那一半（含三枚测试件） | **0** | `120r3-vet-linux.log` |
| `GOOS=darwin GOARCH=arm64 go vet ./internal/winsec/` | 派单点名的交叉尺 | **0** | `120r3-vet-darwin.log` |
| `gofmt -l`（三枚） | 空输出＝0 字节 | **0** | `120r3-gofmt.log`（`body=0 bytes`） |
| `/d/work/base/gopath/bin/gofumpt -l`（绝对路径） | 空输出＝0 字节 | **0** | `120r3-gofumpt.log`（`body=0 bytes`） |
| `sh scripts/d22scan.sh` | 注释豁免、字符串不豁免、零 emoji | **0** | `clean - no D22 ban violations`；分母 **无一下降**：`internal/=568`（120-v1 报 564，+4 含本腿 3 枚）、`cmd/=136`、`ban#7 internal/tools/=25`、`ban#6/#8 frontend/=116`、`ban#8 design/=39`、`examined 280` |
| `sh scripts/check-path-length-budget.sh` | VERDICT GREEN | **0** | `longest=185 chars relative`（仍是别腿那枚 252 的旧件）、`worst full path=229`；本腿最长新件 62 字符（`.scratch/wisp/probes/120/r3/120r3-teeth-on-verdict-by-name.txt`），⛔ 动名册 |
| 容器整包（成对三发） | teeth-on／tooth-269／tooth-272 | **0／1／1** | §3 逐字 |
| 容器名册（两发） | `TestAC[1-4]POSIX`／`120r3$` | **0／0** | 38/38/0/0 与 7/7/0/0 |

## 5. 没做完／没证的（具名，⛔ 摊平）

1. **接缝开窗不是"自然发生的证明"**：`:269` 那枚确定性台件靠 `applyDescriptor` 接缝把交换**落进窗口**
   （真腿先跑、判语保留）。自然发生的证据在枚 7 的盲跑读数里（同发 11 落 `:269`、25 落 `:272`，⛔ 硬判），
   但**枚 7 的红绿不取决于这两枚数**——按 120-v1 §2 裁③，不能承诺的命中数⛔ 升成门。
   ⇒ 若判者认为"自然 landing 必须可判"，需要的是**同一分布的多次读数**，⛔ 一条阈值；本腿只把一发分布落盘。
2. **`other failure = 1`（2000 发）**逐字是 `openat caller-object-120r3.bin: too many levels of symbolic links`：
   只上报、⛔ 判、⛔ 归因到某一行产码——那是赛跑者自己造出的链接链形，`spokeLeg120r3` 把它归"some other failure"
   并列进名册。本腿没查清它的具体交错顺序（没做完，具名在此）。
3. **R-322 那一形在这条路由上是活的**：枚 7 里名字带链接时 `privateFile` 的 `O_TRUNC` 会在 walk 之前把
   外站对象截空（读数行逐字 `the witness now holds ""`）。产码 `winsec_other.go:78-82` 自己写着这一句不做
   （"the O_TRUNC damage of that open is not undone"），票 322 已立未派（编排者 2026-10-15 22:0x 刚把 322-r1
   写成套子）。⇒ 本腿**只计数上报、⛔ 判、⛔ 顺手修**；枚 5/6 的确定性路由用例不受它影响（交换都发生在 open 之后）。
   ⚠ 一处判据因此是**三态**而不是两态（模式仍硬判／raced 载荷落进外站对象仍硬判／只被截空＝计数上报），
   这⛔ 是放宽既有断言——那条"外站对象内容不许变"的硬断言是本腿自己首发起的，**首发即红**，
   改的是**本腿用例的射程**，改的根据是产码注释＋票 322 的已立事实。判者若要，可攻这一点。
4. **Windows 侧零证据**：`winsec_other.go` 带 `!windows`，宿主 `go build ./...`／`go vet` 的 rc=0 证不了这一族
   任何东西（派单 §3.5）。`winsec_windows.go:328` 缺同一对象闸门是**票 120 的 Windows 半格**，非本单射程。
5. **票面措写**（`154/4000` 与"16→20→27"名册口径）⛔ 本腿动一字，只追加 1 条 bullet。

## 6. 台面故障实录（今天本腿现撞三次，⛔ 当读数）

1. **`awk` 替换串里的 `&`** 是"匹配到的文本"：`sub(/\tif /, "\tif false && ")` 造出
   `if false <TAB>if <TAB>if  now.Mode()…` ⇒ `build failed / syntax error: unexpected keyword if, expected {`、
   四数全 `0`（`RUN=0 PASS=0 FAIL=0 SKIP=0`）。**那是台面故障⛔ 是牙响**。改成"现找行号＋`sed -i "${line}s/if /if false \\&\\& /"`＋
   `false &&` 全文件恰好 1 处"的三重自证。**那两发的作废读数件没留档**：快照目录被后续 `rm -rf` 重建，
   日志随之消失（只建不删＝⛔ 删，这里是被我自己的重建步骤带走的，具名记此一错）。
2. **`MSYS_NO_PATHCONV=1` 泄漏给宿主尺**：`export` 之后 `sh scripts/d22scan.sh` 拿到
   `D:\d\work\workspace\projects plans\Wisp … missing go.mod`、`exit status 2`——⛔ 违规、⛔ 分母，是路径转换。
   改为只在 `docker run` 那一步带前缀 ⇒ `rc=0`（`120r3-d22scan.log`）。
3. **`tail -1 f1 f2 f3`**（此 shell 的 tail 不接受多文件＋`-1`）＋`set -e` ⇒ 驱动脚本在中途死掉、后面的尺没跑。
   改 `tail -n 1` 逐件＋host 段 `set +e` 自落 `rc=`。
4. 另两处按派单 §3 提前防掉：容器树放 `/d/tmp`（⛔ Git Bash `/tmp`＝C: Temp，Docker 挂进去是空目录）、
   宿主 `GOTMPDIR` 先 `mkdir -p` 且指 **C 卷**（`C:/Users/swq/AppData/Local/wisp-test-gotmp/120r3gotmp`）；
   成对三发用同一快照同一镜像同一 temp，卷名具名于此。
5. 首发（枚 7 的初版）`panic: close of closed channel`——显式 `close(stop)` 与 `defer` 各关一次，
   另加 `sync.Once`＋join。那一发的 `--- FAIL` 件同样随快照重建消失，读数在提交 `8beaed39` 正文里具名。

## 7. 复跑（可逐字重放；⛔ 抄本腿的数）

```sh
cd 'D:/work/workspace/projects plans/Wisp'
sh .scratch/wisp/probes/120/r3/tooth.sh    # 成对三发：teeth-on／tooth-269／tooth-272（含 LANDS 与 sha 还原自证）
sh .scratch/wisp/probes/120/r3/gates.sh   # 两发名册＋宿主六把尺，逐把自落 rc=
# 变色尺（⛔ 只看四数）：
L=.scratch/wisp/probes/120/r3
diff "$L/120r3-teeth-on-verdict-by-name.txt" "$L/120r3-tooth-269-verdict-by-name.txt"
diff "$L/120r3-teeth-on-verdict-by-name.txt" "$L/120r3-tooth-272-verdict-by-name.txt"
# 产码未被改动：
sha256sum internal/winsec/winsec_other.go   # 137442cdecbb…（与快照 clean 同一串）
```
盘上件（⛔ 删）：`.scratch/wisp/probes/120/r3/` 19 枚读数＋`tooth.sh`／`gates.sh`；
去牙底与快照留 `/d/tmp/probe120r3/`（`fixed/internal/winsec/winsec_other.go.120r3.bak` 两支各一份、
`winsec_other.go.bak-120r3-teeth-on`）。**⛔ push。**
