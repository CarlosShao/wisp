# gw-1 量色腿 · 味①＝票 321 `AC#2` 五把门对今天 HEAD 重跑（★只测不写·⛔ 落任何产码/测试字节）

尺＝本腿 `sh/01-gates.sh`（第 1/2/3/5 把）＋`sh/06-attest-gofumpt.sh`⇄`sh/07-vet-posctrl.sh`⇄`sh/08-gofumpt-ci-shape.sh`⇄`sh/11-blobform-fixed.sh`（第 4 把三形＋正控）。
射程＝**共享工作树**（第 1⇄5 把只在共享树跑＝公共块"门禁四数只在共享工作树"）；第 4 把的"blob 形"另打在仓外 `blobcheck`（＝`git cat-file blob` 逐枚取出的入库字节＝CI checkout 看到的同一形）。
★量的（跑出来的退码/到达数），⛔ 算的。跑前逐把 `unset GOTMPDIR TMPDIR TEMP TMP`＋`GOFLAGS=`（带 env＝票 311 GUARD rc=2 射程拒绝、⛔ 读成红）。
退码取形＝先重定向到文件、`echo RC=$?` 紧跟裸命令（⛔ 管进 tail、⛔ PowerShell 内联 `$?`）。

## 锚（三味各自起手⇄交件，`git rev-parse HEAD`）

| 味 | 起手锚 | 交件锚 | 漂没漂 |
|---|---|---|---|
| ① 五把门 | `a19d5bf2c915b24b836016fa80be47dba4e76abc` 22:44:16 | `a19d5bf2…` 22:53:45 | ⛔ 未漂＝同一枚锚，其间无提交 |
| 派单存档锚 | `51aa922eb3a7838e44bf6f1531082c4a8bf48d49`（编排者 22:31 那笔） | — | 本腿起手的 HEAD 已含其后两笔台账/派单笔（`620e01b5`⇄`a19d5bf2`＝**非 `.go` 字节**：尺＝`git diff --name-only 51aa922e..a19d5bf2 -- "*.go"`＝**0 枚**，量的） |

## 一、五把逐枚 rc＋末行逐字

| 把 | 命令（逐字） | rc | 输出字节 | 末行逐字 |
|---|---|---|---|---|
| 1 | `GOFLAGS= go build ./...` | **0** | 0 | （0 字节＝无行可引；`GATE1_BYTES=0`） |
| 2 | `sh scripts/d22scan.sh` | **0** | 23774 | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=244, bans #1-5 cmd/=42, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=607, ban #8 cmd/=143; ban #8 emoji coverage: …`（件内逐字另存 `logs/gate2-d22scan.txt`） |
| 3 | `sh scripts/check-path-length-budget.sh` | **0** | 1462 | `check-path-length-budget.sh: unit cross-check: the byte view and the character view select the SAME 60 paths on this tree, recomputed just now in both units`；判据行＝`VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`（denominator `tracked paths=17331  over-budget=60  covered by roster=60  not in roster=0`） |
| 4 | 见下面"第 4 把三形"＝`gofumpt` ⇄ `gofmt` 并排 | **1** ⇄ **1**（工作树形）／**1**（CI 仪器形 A） | 55⇄49⇄35 行 | `attrib.sh: RED (tracked-only) - tracked-dirty=35 files_dirty=35. …` |
| 5 | `GOFLAGS= go vet -tags winlive ./cmd/wisp/` | **0** | 0 | （0 字节；复量发 `logs/vet-sharedtree-repeat.txt` 同 rc=0／0 字节＝两发同号） |

★四把＝1/2/3/5 全 rc=0；第 4 把（格式门）今天⛔ 绿＝**存量**，逐名见下，本腿⛔ 洗一枚、⛔ `-w` 一次。

## 二、第 4 把的三形（⛔ 一把尺通吃＝并排三形＋仪器自己的牙）

### (a) 权威仪器＝本仓 CI 逐步执行的那两枚（`ci.yml:172-181` gofumpt 步 ⇄ `ci.yml:230` = `sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`）

| 发 | 命令 | rc | 读数 |
|---|---|---|---|
| 牙 | `attrib.sh --self-test` | **0** | `SELF-TEST GREEN - the classifier rings on a line it cannot attribute and stays silent on a line a real ticket owns …`；`--self-test cases=9 failures=0`（★两方向都响＝这把尺没瞎） |
| 形 A（CI 真跑形＝判据） | `attrib.sh --tracked-only` | **1** | `(A) tracked-tree  lines=35 files=35  rule: MUST be empty`（ handed 1141 枚 tracked `.go`） |
| 形 A⇄B | `attrib.sh`（both） | **1** | `(B) working-tree  lines=55 files=53  tracked=35 attributed_tickets=[161 162 174 224 226 33] unattributable=0` ⇄ `RED - tracked-dirty=35 unattributable=0 forms-disagree=0` |

⇒ ★`unattributable=0`＝55 行**逐行都归得到一枚真票**；★`forms-disagree=0`；⇒ 今天这枚 rc=1 的成因＝在册的存量 bench 样本，⛔ 是新造 injury、⛔ 是别的腿的字节。

### (b) 35 枚名册（＝形 A 的 `TRACKED-DIRTY` 逐名列，工作树字节那形）

31 枚住 `.scratch/wisp/probes/<NN>/…`＝票 163/174/183/185/197(×2)/212/220/222/224/231/235/241(×2)/253/257(M1..M8 共 8 枚)/259/263/287(M1..M3 共 3 枚)/33(q1,q2,q3 共 3 枚)；生产面另 4 枚＝`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`。逐名底片＝`logs/attest-tracked.txt`（35 行逐枚带最后触碰笔号＝⛔ 此处省略，件在盘）。

### (c) 同一批字节对**入库 blob** 重量＝"CI checkout 那形"（`sh/10`⇄`sh/11`，射程＝`blobcheck` 35 枚）

| 尺 | 命令 | rc | 名册 |
|---|---|---|---|
| gofumpt（blob 形） | `/d/work/base/gopath/bin/gofumpt.exe -l .`（在 `blobcheck` 根） | 0 | **19 枚仍脏**＝逐名 `logs/blob-dirty.names3`：`probes/163/a1/main.go`、`174/c2/zz174c2_wiring_pair_windows_test.go`、`183/r2/mut/task-boxset-off.go`、`185/r1/mut-m1/hostpath_185.go`、`197/r1c/pre/subagent_197{,_test}.go`、`212/v1/mut/main-noq9.go`、`220/r1/prechange/prechange220_probe_test.go`、`222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go`、`224/v2/dialect_probe_test.go`、`231/v1/mutation/config_reload_branch_after_prefix.go`、`235/v1/mut/prereadoff-m3_222_test.go`、`241/v1/posctl/badly_formatted.go`、`241/v1/probe_v1_readings_test.go`、`253/v2/denoms253v2.go`、`259/r1/gofumpt-negctl/probe.go`、`33/p1/q{1,2,3}/main.go` |
| gofmt（blob 形，独立第二把） | `gofmt -l .`（同射程） | 0 | 15 枚；两把对撞＝`ONLY_GOFUMPT=4`（`174/c2/…`、`197/r1c/pre/subagent_197{,_test}.go`、`220/r1/…`＝gofumpt 比 gofmt 严的那族）⇄ `ONLY_GOFMT=0` ⇄ `BOTH=15` |
| 只在本地工作树脏（＝blob 干净） | 作差 `comm -23` | — | **16 枚**＝12 枚 bench（`257/M1..M8`、`263/v1/src/main.go`、`287/M1..M3`）＋生产面那 **4 枚**＝本地 `autocrlf` 造的 CRLF 形：尺＝逐枚 CR 枚数（对**入库 blob** 用 `tr -cd`）＝`cmd/wisp/models.go` 工作树 334／blob **0**、`pending_read.go` 131／**0**、`tools.go` 225／**0**、`provenance.go` 1115／**0**；正控＝同一条 `tr` 尺对自造 2 枚 CR 件读出 **2** |
| 生产面 blob 脏枚数 | 尺＝`grep -v '^\.scratch/'` 打在 blob 名册上 | — | **0 枚**（量的，射程＝`blobcheck` 35 枚） |

⇒ ★量的：生产面（`internal/**`⇄`cmd/**`）的字节在 **blob 形**上格式门**干净（0 枚）**；今天第 4 把的 rc=1 里，blob 形仍脏的 19 枚**全部**住 `.scratch/wisp/probes/<NN>/**`。
⇒ ★算的（⛔ 当读数）＝"CI 的 gofmt⇄`attrib.sh --tracked-only` 两步在 checkout 上会红"这一推论：CI checkout 因 `*.go eol=lf` 取到的正是 blob 形⇒按 (c) 那 19 枚名册该步会打印非空；★本腿⛔ 证这一味＝今天⛔ 有一次 CI 跑在这枚 HEAD 上（`ecffb93b..origin/dev` 的祖先尺与 run 底片都⛔ 属本腿射程）。⛔ 圆判、⛔ 翻框＝只把两形（本地工作树 ⇄ blob）的逐名册都摆出来交裁腿。

## 三、第 5 把的正控（★证那枚 0 是验过的读数）

| 发 | 台面 | rc | 逐字 |
|---|---|---|---|
| 负档 | 仓外副本 `gwL-vetmut/cmd/wisp/`，坏件先叫 `.go.txt`（⛔ 编译面） | **0**（0 字节） | ＝证明那枚 0 ⛔ 被新件搅出来 |
| ★正控 | 同副本改名成 `.go`＝坏 `printf`（`fmt.Printf("%s %s\n", 1)`） | **1** | `cmd\wisp\vet_posctrl_gw1.go:8:20: fmt.Printf call has arguments but no formatting directives`（93 字节） |
| 还原 | 副本里改回 `.go.txt` | **0**（0 字节） | `diff -rq gwL/cmd/wisp ⇄ 副本`＝只 `Only in …: vet_posctrl_gw1.go.txt` 一行＝**逐枚内容差异 0** |
| 共享树复量 | `GOFLAGS= go vet -tags winlive ./cmd/wisp/` | **0** | 0 字节（与首发同号） |

⇒ 那枚 `VET=0` 是有牙的读数＝同一把尺在同射程打坏 `printf` 时逐字转红、还原后回到 0。全程⛔ 动共享工作树一个字节（突变只落仓外副本）。

## 四、味① 台面登记（台损计数）

- `0xc0000135`＝**0 枚**（味① 不含跑色；跑色那两味的计数见味②③ 件）／`build failed`＝**0**／`walk`＝**0**／弃权族＝**0**。
- 跑门 preflight（`logs/gates.txt`⇄`gates-r2.txt` 的 `DIRTY_LINES`/`DIRTY_GO`）＝tracked 脏 **1 枚**（`.scratch/wisp/issues/115-…-done.md`，别腿在飞面）⇄ **`.go` 脏 0 枚**＝共享树⛔ 处半保存态＝门禁读数可用。
- ★本腿自抓两把坏尺（照实入账、⛔ 藏）：① `sh/01` 首发把 `gofumpt -l ./...` 当"全仓尺"＝`./...` 是 **Go 包模式**⛔ 是文件路径⇒工具按路径走盘、rc=2 且 offender 逐字 `open ...\.git: The system cannot find the path specified.`＝**用法坏**（原始失败留档 `logs/gate4a-wrongshape.txt` 复现同串）；改正＝改用 CI 在册仪器 `attrib.sh`＋并排打在目录 `.` 上。② `sh/10` 第 3⇄4 段那把 `comm` 名册尺：两侧路径不同形（`D:/tmp/…` ⇄ 仓相对）⇒作差把 35 枚**全**判成"只在本地脏"＝假读数；随后我在 pwsh 内联补的尺又因 `sed` 表达式被 shell 吃掉而整条失败（stderr 两行 `unterminated 's' command`⇒名册文件 0 行＝`BLOB_DIRTY=0` 属**尺坏**、⛔ 是读数）。改正＝`sh/11-blobform-fixed.sh` 落仓外脚本文件跑＋归一到同一仓相对形＋★正控先行（造一枚脏件→名册 20 行、删掉→19 行＝尺跟着还原走）＝上面 (c) 那张表才是数。
- ⛔ 洗已知脏件、⛔ `-w`、⛔ 放宽任何断言、⛔ 翻框、⛔ push、⛔ 动任何 `.go`/票面/台账字节。
