# 票 336 落地腿 336-r1 · 交件报告（仪器补牙，2026-10-18 第一波）

写面＝三处：`internal/risk/*_336*_test.go`（2 枚，全测试面）＋`.scratch/wisp/probes/336/r1/**`＋票 336 尾一条 bullet。
⛔ 动 `blacklist.go` 任何非注释字节（尺：逐笔 `git show --name-only` 全集⊆三处、`.go` 非测试命中 0，见 第 7 节）。
起手锚现认＝`c11163c8`（10-18 10:02:56）；台面锚具名＝`ed8ff551`（放出笔），`internal/risk` 在两号下 `git diff --stat`＝空 rc=0、
`blacklist.go` blob 两号同为 `c1799cd24afce19e6accef23db23a7d7ee746dec`（＝334-v1 交件 blob）。交件侧 HEAD 见 第 7 节。
★新铁律已生效于本腿：落点＝`D:\tmp\wisp\20261018\336r1\{tree-before,tree-mut,tree-comm,tree-after,gotmp,msg}`（⛔ 盘符根）；
8.3 别名现量＋"无一枚判据需要根目录"的读数＝`logs/08-alias83.log`（`dir /X`＋`cygpath -d` 逐枚：`TREE-B~1`/`TREE-A~1`/`TREE-C~1` 三枚有别名字，
`tree-mut`(8)/`gotmp`(5)/`msg`(3) 名已 ≤8 字符⇒短名＝长名；本腿⛔ 因此停手上报）。

## 1 AC#0 起手现量＝票面四发串行重放，逐字复认（`logs/ac0-runner.log`＋`ac0-0[0-4]-*.log`）

台面＝`tree-mut`（`ed8ff551` 导出，`.go` 748 枚证非空）；突变只打盘上物理字节、⛔ `-overlay`、⛔ 改跟踪文件；
逐发「从具名提交重取字节覆盖」＋跑前跑后各验一次哈希＝八次全 `match=YES`（`c1799cd2…`）；末态 `shared_worktree_untouched=0 行`。

| 发 | 落点（内容锚行差） | 三枚在册仪器四数（rc） | 与票面 |
|---|---|---|---|
| 基线 | — | `10/3/0/0 rc=0` | 同 |
| M-V2a `%LOCALAPPDATA%\Microsoft\Credentials` | `:237` aUnder→anchorMatches | `10/3/0/0 rc=0` 全绿 | 复认 |
| M-V2b `~/.ssh/**` | `:227` 同上 | `10/3/0/0 rc=0` 全绿 | 复认 |
| M-V2c `~/.kube/config` | `:231` aEq→anchorMatches | `10/2/1/0 rc=1` 红＝`installing-aws-moves-the-named-anchor-to-kube` | 复认（＝ⓑ 反证） |
| M-V3 去祖先界（变严） | `:366` return true | `10/1/2/0 rc=1` 红＝AC#1 两枚＋票 216 那枚；★闸门格 `TestAnchorMissKeepsNonOverridableClassAndNeedL2` 逐字 `--- PASS` | 复认＝票面"闸门格恒不红"成立 |

★具名更正一处（本腿自己写错的，⛔ 改判据原文）：提交 `f67048d9` 的正文说 M-V3 "闸门格在红名册里"——现量逐字为 `--- PASS`，
红的是 AC#1 的 reason 两枚与票 216 那枚。票面那句成立，我这句写错，在此与票尾 bullet 更正。

## 2 AC#1 名册（尺＝`git grep -n -e aEq -e aUnder ed8ff551 -- internal/risk/blacklist.go`，具名提交树）

原始命中 12 行＝2 行注释（`:193`/`:265`）＋2 行定义（`:268`/`:270`）＋8 行用法 ⇒ **8 枚回落支操作数跨 7 枚 A 档锚**：
`:220` `.git-credentials`(eq)、`:224`＋`:225` `%APPDATA%\wisp` 这一枚 case **AND 两枚操作数**（under＋eq）、
`:227` `~/.ssh`(under)、`:229` `.aws`(eq)、`:231` `.kube`(eq)、`:235` `Microsoft\Protect`(under)、`:237` `Microsoft\Credentials`(under)。
不落回落支的 A case＝`:222` `.git/config`、`:233` browser store（`aAnyForm`）＝本票射程外。
今天有谁钉（⛔ 抄票面"七枚"那句）：只有 `.aws`（`TestAnchorMiss*` 两枚）与 `.kube`（ⓑ 那枚子用例）活在回落支上；
其余五枚锚在册用例里只被**叶子真装上了**那一形读到（`TestAListDenyAndUnoverridable`／`TestTierAAnchorsMatchNativeSpelling`／
`TestClassifyAnchorSpellingIsNotVerdict`／`TestReparseExceptionExplicitPathOnly`，home 全是假主目录）＝`anchorMatches` 命中支，⛔ 回落支。
⇒ 票面标题那句"只有 `.aws`＋`.kube` 会红"在**突变现量**这一维复认（第 1 节／第 3 节 的 LEG 尺逐名），本腿把射程从 2 枚扩到 8 枚操作数。

## 3 AC#2 牙＝`internal/risk/anchor_fallback_teeth_336_windows_test.go`（`c55e32b1`）＋逐操作数突变扫描（`8c36c23c`）

每枚操作数一枚用例：装齐**其余**每一枚叶子（装上的叶子解析到自己⇒与候选无祖先关系⇒诚实 miss、⛔ 抢 switch），只缺被测那一枚；
钉四件事：`class=A` ＋豁免键买不到放行 ＋ reason 点名**本锚**字面串（switch 有序⇒只看 class 是盲的）＋ `anchorMatches=false 且 uncertainAnchorMiss=true`；
另带同世界链外臂＝读 B（防"整台机器都算 A"那一形冒充成功）。
成对正控（未突变）＝`teeth-01-pristine-dshape.log`：**`rc=0 RUN=17 PASS_TOP=1 FAIL_TOP=0 SKIP=0`**，八枚逐枚 `matched=false clause=true`。
★本腿自抓一名册错并具名更正（`teeth-00` 首发 `rc=1 FAIL_TOP=1`）：`wisp-config-under` 从 `<profile>\work` 出发是 **CERTAIN miss**
（`<APPDATA>\wisp` 的 root 与 work 无祖先关系，回落支压根不参与⇒摘它不变判据）；该行落点改正＝候选放 `<APPDATA>\Roaming`。

去牙扫描（`run-ac2-sweep.sh`；一发只摘一枚、内容锚、`落点命中枚数=1` 逐发、每发前后 pristine `match=YES`）：

| 发／落点 | 我的牙（TEETH） | 在册三枚（LEG） |
|---|---|---|
| MU-1a `:220` | `rc=1 17/0/1/0` 红＝`git-credentials` | `10/3/0/0 rc=0` |
| MU-1b `:224` | `rc=1` 红＝`wisp-config-under`，★只咬 `:138` 那一句（reason 点名本锚）——`class` 仍 A，由**更后面**的 `dpapi-protect` 经自己的回落支接住 | rc=0 |
| MU-1c `:225` | `rc=1` 红＝`wisp-config-under`（:138）＋`wisp-config-eq`（:129/:132/:135/:138/:141 五句） | rc=0 |
| MU-1d `:227` | `rc=1` 红＝`ssh` | rc=0 |
| MU-1e `:229` | `rc=1` 红＝`aws` | `10/2/1/0 rc=1`＝`absent-aws-leaf…` |
| MU-1f `:231` | `rc=1` 红＝`kube` | `10/2/1/0 rc=1`＝`installing-aws-moves…` |
| MU-1g `:235` | `rc=1` 红＝`dpapi-protect` | rc=0 |
| MU-1h `:237` | `rc=1` 红＝`local-credentials` | `10/3/0/0 rc=0`＝票面 M-V2a 那半格逐字复认 |

⇒ 摘任意一枚的回落支必红（8/8），且 LEG 尺只接住 2 枚＝票 334 `AC#2` 那半格在本腿名下补齐；红名册逐名对撞见上表，⛔ 只看枚数。

## 4 AC#3 两形代价（现量，⛔ 定案；选不选＝裁腿＋编排者）

甲＝本腿已落（第 3 节）。乙＝把 334-v1 那把 42 行三值对撞尺搬成仓内断言，代价现量（`logs/ac3y-runner.log`，尺原件 blob `69a70fa4…` 逐字同串）：
未突变 `rc=0 RUN=22 PASS_TOP=4 FAIL_TOP=0 SKIP=0`，首发 wall 20.48s（含编译）⇄ 复用缓存 4.26s；
★摘 MU-1h 那一枚后这把尺 **`rc` 仍＝0**（它 13 发 `t.Log`＋7 发 setup `Fatalf`＝⛔ 断言），只在 shape 文件显形
`chain:case9-local-creds class=A→B／rule=A:Microsoft-Credentials→B-tier` ⇒ 搬进仓内必须逐枚补 want＝**18 行 chain ＋ 9 行 offchain＝27 枚具名期望**，
248 行一件、一装进跟踪面就吃 CI 两平台射程（票 339 那一形）。甲在同一发 `rc=1 17/0/1/0`、wall 1.85s、红名册逐名点名。
⇒ 材料交上，⛔ 本腿定案。

## 5 AC#4 颜色与门禁

成对两把临时根台面、同一条 `-run` 名册（我的牙＋票 334 三枚在册仪器）、改前 `ed8ff551` ⇄ 改后 `8c36c23c`（`logs/ac4-runner.log`）：
ci-shape（`%LOCALAPPDATA%` 组法）前 `10/3/0/0 rc=0`（pkg 0.727s）⇄ 后 `27/4/0/0 rc=0`（0.541s）；
dshape（`gotmp/`，先 `mkdir -p`＝`mkdir_gotmp_rc=0`）前 `10/3/0/0 rc=0`（0.168s）⇄ 后 `27/4/0/0 rc=0`（0.351s）＝两形同色、⛔ 混用⛔ 相减。
整包 `-run .`（同台面各两发）：前 `230/135/5/1 rc=1` ⇄ 后 `247/136/5/1 rc=1`；红名册逐名作差 `diff_rc=0`（五枚同名，
`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`＋四枚 `TestSyncAncestor*/TestSyncFirstActable*116`），SKIP 两发同一枚 `TestSyncRegistryProbeLive`。
★台面敏感具名（⛔ 计入本票判据）：那四枚 116/junction 族只在 `--dshape`（GOTMPDIR 落 D 卷）红，ci-shape 整包红名册只有 1 枚＝既有台面形状，归位交裁腿。
必答(ii) 时长差＝我的牙顶层 0.21s＋16 发子用例逐枚 0.01–0.02s（`RUN` 230→247＝＋17 枚逐名对上）；整包包级秒数在台面噪声内摆动（前 12.593/4.258 ⇄ 后 3.477/3.100）＝⛔ 当差值读数。
门禁四数只在共享工作树（`logs/gates-runner.log`，跑前 `unset GOTMPDIR TMPDIR`＝lead 提醒）：
`gofumpt.exe -l` 对我两枚新增件**列出 0 枚 rc=0**；包级 `-l` 只列存量 `internal\risk\provenance.go`（CRLF 假脏，⛔ 动它一字；★第一发我把自己 `logs/*.sh` 也喂给 gofumpt＝尺用错、`rc=2` 是 `illegal character U+0023` 噪声，本版改正并具名登记）；
`sh scripts/d22scan.sh` **rc=0** clean（ban#8 internal=589 含我两枚 .go＝注释豁免那一维也过）；长度门 **rc=0** VERDICT GREEN（分母 15469）；`go vet ./internal/risk/` **rc=0**。

## 6 必答题（逐字）

(i) 现量枚数＝**8 枚回落支操作数／7 枚锚**；票面与 334-v1 那句"七枚 aEq/aUnder 锚"＝**按锚计数**，差在我这把尺把
`%APPDATA%\wisp` 那一枚 case 的 **AND 两枚操作数**分开算（`blacklist.go:224` 与 `:225`）——尺＝`git grep -n`（12 行原始命中，2 注释＋2 定义＋8 用法）。
(ii) 见 第 5 节：定向名册两形台面各自前后（0.727→0.541／0.168→0.351，pkg 级）；我的牙自身 0.21s 顶层＋16 发 0.01–0.02s；整包 RUN＋17 枚逐名对上，包级秒数＝噪声、⛔ 当差值。
(iii) **缺世界具名＝整个 POSIX 平台**：`pathresolver_other.go` 的 `resolveHandle` 是 `("", false)` 桩 ⇒ 候选与每一枚锚 `root=""`，`uncertainAnchorMiss` 恒 false，
"摘这一枚⇒必红"在 POSIX 造不出来。容器实测（`golang:1.27`，`logs/ac4-08-linux-tripwire.log`）逐枚 `cand_root=EMPTY`、我的绊线 `--- PASS`。
⇒ ⛔ 假绿、⛔ 恒红：新落 `anchor_fallback_teeth_336_other_test.go`（`c87df16d`，`//go:build !windows`）把这件事量成**绊线**（POSIX 解析器一落地 root 非空即当场红并点名"把牙搬过来"），
并具名交 CI 射程账＝**票 339 债⛔ 由本票新增**（本腿两平台各自绿；Linux 档现存两枚红＝票 334 字节 `anchor_miss_reason_334_test.go:131/:156/:193/:197/:201`，早于本腿锚，⛔ 本腿引入、⛔ 本腿动它一字）。
另具名一枚非平台缺口＝**枚数对不上＝8 发 8 枚不同落点，但红名册只有 7 枚唯一名**：`MU-1b`（`:224`）与 `MU-1c`（`:225`）都咬 `wisp-config-under`
＝同一枚 case 的 AND 两枚操作数在同一批世界（`:224` 与 `:225` 的 `omit` 都是 `wispCfgLeaf334`）上重叠，"一发只摘一枚"在这一枚 case 上做不到红名册单枚。
★且 MU-1b 只咬 `:138`（reason 点名本锚）这一句＝`class` 没降：那枚候选改由**更后面的** `dpapi-protect` 经自己的回落支接住——
这正是本票 第 3 节 那句"只看 class 是盲的"的实测形状。⇒ 若裁腿要逐枚独立，需给那枚 case 再造第三个世界（装回 `config.toml` 叶子、只缺 `wisp` 目录那一形），属**加件**、本腿⛔ 自行扩面。
(iv) 成对两发落点不同自查（逐名列，尺＝逐发 `落点命中枚数=1`＋`diff` 行差逐字在 `ac2-sweep-runner.log`）：
发⇄物理行＝MU-1a `:220`／MU-1b `:224`／MU-1c `:225`／MU-1d `:227`／MU-1e `:229`／MU-1f `:231`／MU-1g `:235`／MU-1h `:237`
＝**8 发 8 枚不同物理行，⛔ 任何两发落在同一枚行**；每发红句咬的是被测锚自己的子用例名（＋自己的字面串）。
★唯一一处名册交叉具名＝MU-1b ⇄ MU-1c 都含 `wisp-config-under`（同一枚 case 的两枚 AND 操作数，见 (iii) 末段），属**枚数**问题、⛔ 落点重合。
红句 `file:line` 逐发＝除 MU-1b 只 `:138` 一句外，其余七发各咬 `:129`/`:132`/`:135`/`:138`/`:141` 五句（＋`:148` 读数行）＝六处尺面，逐名取自 `ac2-MU-*-teeth.log`。

## 7 纪律格与交件自证

逐笔（严格前缀尺 `git log --all --grep='^票 336 落地腿 336-r1'`＝**8 笔**，越界 0）：
`ff56a9e0` 台件 1 枚／`f67048d9` AC#0 8 枚／`c55e32b1` AC#2 牙 1 枚（`.go` 测试面）／`feac8f77` AC#2 正控 3 枚／
`8c36c23c` AC#2 突变名册 21 枚／`c87df16d` POSIX 半格 2 枚／`90dafe77` AC#4 20 枚／`28da3a7e` AC#3 7 枚。
写面尺＝逐笔 `git show --name-only` ⊆ `^internal/risk/[^/]+_test.go` ∪ `^\.scratch/wisp/probes/336/r1/` ∪ 票 336 ⇒ `out_of_scope=0` 逐笔（第 2 节 命令可重跑）；
`.go` 非测试命中 **0**；`blacklist.go` blob 交件侧仍 `c1799cd2…`（＝起手同串）；票面七格现跑**未勾 7／已勾 0**、票名与 `Status:` ⛔ 动。
形状尺：`anchor_fallback_teeth_336_windows_test.go` 160 行／`..._other_test.go` 92 行，均 CRLF 0、非 ASCII 0、`exported_non_test=0`、`t.Skip=0`。
具名遗留＝`logs/ac4-full-before-dshape-a.rednames` 等 **4 枚 0 字节**＋`teeth-01-pristine-dshape.log.redlist` 0 字节＝
★尺坏 artefact（我在双引号里写 glob＝失效，`redlist=0` 是假读数，真名册＝`*.log.rednames` 逐枚非 0）；按"临时件只建不删"⛔ 删、⛔ 当凭据。
⛔ push、⛔ 翻框、⛔ 改票名；`go test ./cmd/wisp` ⛔ 跑（唯一持跑者＝`228-r3c`／`322-r2`）；本腿⛔ 碰 `internal/panel/**`、`cmd/wisp/**`、`frontend/**`、`design/**`。
