# gw-1 量色腿 交件汇总（★读数腿＝只交付成对读数＋门禁 rc、⛔ 落任何产码/测试字节、⛔ 翻框、⛔ push）

## 一、提交名册（三把尺同号，尺＝`sh/13`；枚数随 HEAD 漂⇒只认逐枚 `%s` 起头）

| 笔 | 提交号（短） | 时间 | 内容 | `--name-only` 枚数 ⇄ `.go` ⇄ 越界 |
|---|---|---|---|---|
| ① | `3e370700` | 23:04 | 味① `01-gates.md` ＋尺 `sh/01..sh/11` | 12 ⇄ **0** ⇄ 0 |
| ② | `8f86e936` | 23:45 | 味② `02-five-packages.md` ＋尺 `sh/12..sh/16` | 6 ⇄ **0** ⇄ 0 |
| ③ | `247d7f9e` | 23:46 | 味③ `03-winsec.md` | 1 ⇄ **0** ⇄ 0 |
| ④ | `54963e61` | 23:47 | `logs/` 六枚驱动转写（≤160 行那批） | 6 ⇄ **0** ⇄ 0 |
| ⑤ | 本笔 | — | 本汇总件 | 1 ⇄ **0** ⇄ 0 |

名册三把（本笔落笔**前**现量＝4⇄4⇄4 同号；★正控三枚＝`^票 139 落地腿 139-r1`＝11、`^票 321 终裁腿 321-v1`＝6、`^票 353 落地腿 353-r1`＝6＝尺⛔ 瞎）＝
① 锚定行首打在 `git log --all --format='%s'` 的导出（射程＝只有主题）＝锚定尺⛔ 吃正文引用；② `git log --all --grep="^gw-1 量色腿 gw-1"`；③ 单次 `%H%x09%s` 导出＋内存 `index($2,PFX)==1`。非锚定对照两把＝3⇄3（★⛔ 当名册）。
★写面自证正控＝同一把尺（`git show --name-only` 后 `grep -c '\.go$'`）打在 `3a509e50`（353-r1 笔①）给 `total=8 go=1`＝尺灵⇒本腿四笔 `.go`＝**0 枚**是真 0。

## 二、三味各自〔起手锚 ⇄ 交件锚〕＋判据结论

| 味 | 起手锚 | 交件锚 | 漂没漂 ⇄ 硬尺 | 结论（★只报数） |
|---|---|---|---|---|
| ① 五把门 | `a19d5bf2` 22:44:16 | `a19d5bf2` 22:53:45 | ⛔ 未漂 | 四把 rc=0⇄格式门 rc=1（存量 35 枚名册＋blob 形 19 枚名册，`unattributable=0`） |
| ② 五包成对色 | `a19d5bf2` 22:49 | 复量段 `a0124803`→`4fb6c088` | 漂＝**非 `.go`** 笔；尺＝`git diff --name-only a19d5bf2..HEAD -- "*.go"`＝**0**＋七目录 tree hash `SAME=YES` | ★新增红 **0 枚**（十格逐格 0；容器 config 那处单发换位经 3 发并集归因＝两侧并集同批名） |
| ③ winsec 两平台 | `a0124803` 23:25 | 复量段 `4fb6c088`⇄交件 `35fffa29` | ⛔ 未漂（复量 6 发起手＝结束同串）；同上两把硬尺 | 对甲新增红 **0 枚**、消失红 **1 枚**＝353 自己那枚仪器（基线 4/4 红⇄落地 4/4 绿）；container 两侧红册各 0 |

## 三、味① 五把门（逐枚 rc＠起手锚 `a19d5bf2` ⇄ 交件侧收线复量＠`35fffa29`＝两锚都报、⛔ 混时点）

| 把 | 命令 | rc＠`a19d5bf2` | rc＠收线 `35fffa29` | 末行／凭据 |
|---|---|---|---|---|
| 1 | `GOFLAGS= go build ./...` | **0**（0 字节） | **0**（0 字节） | — |
| 2 | `sh scripts/d22scan.sh` | **0** | **0** | `clean - no D22 ban violations; … ban #8 internal/=607, cmd/=143`（收线＝同串、23771 字节⇄23774＝分母随别腿入库件漂） |
| 3 | `sh scripts/check-path-length-budget.sh` | **0** | **0** | `VERDICT GREEN …`；denominator 起手 `17331/60/60/0` ⇄ 收线 `17403/60/60/0`（★tracked 数随别腿笔漂＝⛔ 当通则、⛔ 判红） |
| 4 | `attrib.sh --tracked-only`（＝CI 真跑形） | **1** | **1** | `RED (tracked-only) - tracked-dirty=35 files_dirty=35`；形 B `unattributable=0`⇄`forms-disagree=0`；`--self-test cases=9 failures=0`＝牙在 |
| 5 | `GOFLAGS= go vet -tags winlive ./cmd/wisp/` | **0**（0 字节） | **0**（0 字节） | ★正控＝仓外副本坏 `printf`→**rc=1** 逐字 `cmd\wisp\vet_posctrl_gw1.go:8:20: fmt.Printf call has arguments but no formatting directives`、负档⇄还原都 0、`diff -rq` 内容差 0 枚 |

第 4 把的三形⇄35 枚⇄19 枚逐名册＝`01-gates.md` 第二节；生产面（`internal/**`⇄`cmd/**`）blob 形脏枚数＝**0**（量的）。
★时点口径＝味① 那五把在起手锚 `a19d5bf2`（⇄交件同串＝未漂）跑过一次、又在交件侧 `35fffa29`（＝含本腿四笔读数件之后；★其间 `.go` 差异仍 0 枚）收线复量一次＝两锚都报、⛔ 拿任一发的号顶另一发。

## 四、味②⇄③ 成对四数（量的，`^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP`，⛔ 混用⛔ 相减）

五包（基线 `0fb399c3` ⇄ 落地 `a19d5bf2`）：host＝`secret 0/29/21/0/0`⇄同／`config 0/195/114/0/0`⇄同／`memory 0/80/43/0/2`⇄同／`agent 0/114/90/0/0`⇄同／`panel 1/268/163/7/0`⇄同；container＝`secret 0/4/4/0/0`⇄同／`config 1/196/112/1/0`⇄同（★红名册并集同批＝摆动族）／`memory 0/79/42/0/2`⇄同／`agent 0/113/89/0/0`⇄同／`panel 1/268/163/7/0`⇄同。
winsec（对甲 `3a509e50`⇄HEAD）：host 四发＝`121/67/3/0`（基线）⇄`121/67/3/0`首发、`121/68/2/0`×3（落地）；container＝`74/46/0/0`⇄同。对乙（`0fb399c3`⇄HEAD）＝host `119/65/3/0`⇄`121/69/1/0`；container `74/46/0/0`⇄同。★`RUN` 那 2 枚差由名集 symdiff 解释＝`ONLY_LAND=2`/`ONLY_BASE=0`＝正是 `nametopup_353_windows_test.go` 那两支。
逐名红册／具名连带逐枚色（五枚判定八发全 PASS、`openleaf_322`⇄`placement_symlink_113` 两档色、panel 存量 7 枚逐名）＝`02` 一⇄四节、`03` 二⇄四节。

## 五、台损枚数具名（★⛔ 计入判据）

| 形 | 计数 | 尺／射程 |
|---|---|---|
| `0xc0000135` | **0 枚** | 每发原始件 `grep`；全部 34 发（味② 20＋味③ 14）`SHAPE=OK`／`VOIDKIND=none`（量的，尺＝`sh/05`⇄`sh/14`⇄`sh/15` 写进 `four-counts.tsv`） |
| `build failed`／`[build failed]` | **0 枚** | 同上 |
| `walk`／`setup failed`／`no test files` | **0 枚** | 同上 |
| 结构性弃权（`!windows`/`unix` 档⛔ 编译＝ABSENT） | 味② host 侧 **10 名**（`openleaf_322` 1＋`placement_symlink_113` 8＋`overwrite` 那枚在 container 侧 0）⇄ container 侧 **2 名**（`openleaf_guard_322` 两枚）；味③ 名集面差 host 121⇄cont 74＝`only_cont=52`／`only_host=99` | ★既⛔ 绿也⛔ 记弃权＝票 322 `AC#1` 那形；⛔ 拿任何一枚的绿当凭据 |
| 热键族（机主常驻 `wisp.exe`） | 在场＝**现量 1 枚 pid 21892**（`Get-Process`，StartTime 10-15 19:03:14）；红名册里热键族命中＝**0 枚** | 弃权、⛔ 归因给字节 |
| `…116` 族 4 枚 ⇄ `TestTicket223ModeLoosening…` | 命中 **0**（同尺正控读 1＝尺⛔ 瞎）⇒ 真身分别住 `internal/risk/syncdirs_ancestor_actable_leg_116_test.go`⇄`cmd/wisp/config_reload_223_test.go:535`＝**⛔ 本腿三味射程** | ★缺哪一味＝具名报，⛔ 圆判：`…116` 要 `./internal/risk/`（⛔ 派单射程）；223 那枚要整包 `./cmd/wisp`＝派单逐字⛔ 跑、归编排者那窗 |

台面：台根＝`%LOCALAPPDATA%\wisp-test-gotmp\gw1-<ts>`（★⛔ `D:\tmp\wisp`＝109-r1 那味；现量 `icacls D:\tmp\wisp`＝带 `BUILTIN\Users:(I)(RX)`＋`(I)(OI)(CI)(IO)(GR,GE)` 两条继承、`%LOCALAPPDATA%` 那条根只有三条 `(I)(OI)(CI)(F)`）；`GOTMPDIR`/`TMPDIR`＝`cygpath -w` 反斜杠形＋先 `mkdir -p`；跑 panel/winsec 带 `PATH=<tree>\third_party\sherpa-onnx;<tree>\build`（⛔ `cygpath -wp` 转整条 PATH）；树＝八件套＋`docs/`＋手工三件、`.go` 非空尺 781/778/777。

## 六、tracked 脏（逐枚列名＋mtime，★⛔ 只报计数）⇄ 冻结件

- 交件侧现量＝**1 枚**：`M .scratch/wisp/issues/115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md` mtime `2026-10-20 20:24:44`（＝`115-v1` 面、⛔ 本腿；本腿全程⛔ 读其内容⛔ add⛔ 卷进任何一笔）；`.go` 脏＝**0 枚**；本腿对 tracked 脏净贡献＝**0 枚**（复算＝`git status --porcelain -- .scratch/wisp/probes/gw/1` 交件时 0 行）。
- 味① 跑门前后的另两枚＝`M .scratch/wisp/probes/353/v1/00-ruler-register.md` mtime `23:38:08`（＝`353-v1` 在飞，本腿只 stat⛔ 读）＝复量发逐名已入 `logs/delivery-rulers.txt`⇄`logs/delivery-rulers-final.txt`。
- 仓根冻结件＝`go.mod` 1026B/19:50:08、`go.sum` 6237B/10-01、`.gitignore` 1404B/09-27、`ci.yml` 64794B/10-15＝四枚 mtime 全早于本腿起手、blob CR 全 0、⛔ 0 字节（`A1069`③ 那形的自证）；全程⛔ 任何输出重定向落仓根件、⛔ `go mod edit/tidy/get`。

## 七、具名上报（⛔ 圆判那几味）

1. ★第 4 把今天 rc=1＝**存量**；35 枚名册里 31 枚归得到真票（形 B `unattributable=0`）＋生产面 4 枚只在本地工作树脏（blob CR=0⇄工作树 CR 334/131/225/1115＝`autocrlf`）＝CI checkout 那形⛔ 本腿能证（今天无 CI 跑在这枚 HEAD 上）＝⛔ 本腿能替这格盖章。
2. ★票面把 `private_fail_test.go:37/50/61/72/75` 算五枚判定，仪器只印四枚子测试名（`:72`⇄`:75` 同住 `/seal_file_and_dir_direct`）＝八发全 PASS；本腿⛔ 改口径、⛔ 圆判＝两形都摆着。
3. ★`internal/agent` 只打**包自身**（⛔ `/**`）＝`internal/agent/approval/**` 是 `230-r1` 写面那族（今天两侧 tree hash SAME⇒本腿读数未卷它字节）。
4. ★本腿自抓三把坏尺（照实入账）：① `gofumpt -l ./...` 用法坏（`./...`＝包模式⛔ 路径⇒rc=2 `open …\.git`）；② `comm` 名册两侧路径不同形⇒35 枚全误判"只在本地脏"；③ 钝 CR 尺 `tr -cd '\r' | wc -l`（只数换行⇒读 0）——本汇总里所有 CR 数都取自 `wc -c` 形＋正控（自造 2 枚 CR 件读出 **2**）。
5. ★这里⛔ 一枚读数能替 352/321/353 三张票的 AC 盖章＝本腿是读数腿：⛔ 翻框、⛔ 动票面/台账、⛔ push。★CI 那味现量两把（量的、只报数）＝`git merge-base --is-ancestor ecffb93b origin/dev`＝**rc=1** ⇄ `--is-ancestor 54963e61 origin/dev`＝**rc=1** ⇒ 落地字节与本腿五笔读数件都仍⛔ 在 `origin/dev` 之外＝"CI 从未见过落地字节"那一味照旧成立（⛔ 本腿处置＝推送只由编排者做、且机主按住）。
6. ★临时件只建不删＝`D:\tmp\wisp\20261020\gw1\{gwB,gwW,gwL,gwL-vetmut,blobcheck,posctrl,gocache-win,gocache-col-*,logs}`（162 枚原始底片留仓外 `logs/`＝帽口径先例＝139-r1 `04` 第 B 节；★可逐字重跑的尺 16 枚已入库 `sh/01..sh/16`）。

## 八、台损那一行的精确形（★纯追加＝⛔ 抹上面原句；起因＝本腿复量时抓到自己第五枚钝尺）

第五节那张表里"34 发"⇄"walk＝0 枚"两格⛔ 是精确读数，改正如下（两把尺都现跑、都带正控）：

| 问 | 精确读数 | 尺 |
|---|---|---|
| 跑色发数 | **40 发**＝味② 20（10 host＋10 container）＋味③ 8 成对＋6 复量（`ws3`）＋6 复量（`cfg`） | `wc -l` ⇄ `ls *.raw.txt \| wc -l`＝**40 枚**原始件（量的，射程＝仓外 `logs/`）；第五节写"34"＝我只数了带 `SHAPE=` 字段的那 34 行（28 行 `SHAPE=OK VOIDKIND=none`＋6 行 `ws3` 只带 `SHAPE=OK`）而 `sh/14` 那 6 发 `cfg` 复量的行形⛔ 带该字段＝**计数尺的射程⛔ 覆盖全部发**＝本腿第五枚钝尺，照实入账 |
| 零到达（台损那一形） | **0 枚** | 逐枚 `grep -c '^=== RUN'` 打在全部 40 枚原始件上＝`ZERO_RUN_COUNT=0`；★正控＝同一条尺在 40/40 枚里都读得出 `=== RUN` 行（`grep -rl` 命中 40/40＝尺⛔ 瞎） |
| `0xc0000135`／`build failed`／`setup failed`／`no test files` | **各 0 枚** | `grep -rl` 打在 40 枚原始件上＝命中 0；★先正控（同尺读 `=== RUN`＝40/40） |
| `walk` | 台面故障形＝**0 枚**；★⛔ 拿朴素 `grep -i walk` 的命中当读数＝它命中 **30 枚**原始件 | 朴素那把钝尺⛔ 用（＝公共块"名册尺⛔ 非锚定子串"同形）：30 枚命中**逐字全是用例名与日志正文**（`TestAC4SealedWalkSkipsLinks`／`TestAC4DirectorySymlinkInsideSealedTreeIsNotWalked`／`leg that answered = pinned walk`／`pinnedWalkRefusals=61`）；精确那把＝`grep -E "walk [A-Z]:\|walk .*: The system"`（＝d22scan GUARD  offender 那形）命中 **0 枚** |

⇒ 本节只更正"我数了几发、`walk` 那格用的哪把尺"（发数 34→**40**；`walk` 的尺从朴素子串换成台面形那把＝仍 0 枚）＝★判据读数⛔ 变＝三味四数⇄逐名红册⇄新增红枚数全部出自 `RUN>0` 的 40 枚原始件（逐枚量过）；⛔ 改任何一枚色、⛔ 翻任何框、⛔ 落任何 `.go` 字节。
