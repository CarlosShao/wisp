# 111-r5 — 把已存在的门自检脚本接进 CI 一步（GUARD D 正控读数的载体）

腿号 `111-r5`。零 go 命令。写面只有 `.github/workflows/ci.yml`（加一步）＋本件＋票面追加一行 Progress log。
本件不宣称 AC#3 完成：**"接进 CI 之后它到底响没响"只有推送后的 run 能回答**，见 §4。

---

## §0 起手锚

| 尺 | 读数 |
|---|---|
| `git log -1 --format='%h %ad %s' --date=format:'%m-%d %H:%M'` | `15fbf18f 10-06 13:24 probe(pool-validity/4c): 骨架 —— 15-26 这 8 枚最老零勾票的名册与四档尺` |
| `date '+%m-%d %H:%M'` | `10-06 13:25` |
| `git status --porcelain -- scripts .github` | **0 行**（`LINES=0`，与编排者 13:1x 的现量一致 ⇒ 没有别的写腿在这面） |

票面框尺 `grep -n '^- \[ \]' .scratch/wisp/issues/111-ci-tests-20-of-33-packages.md` ＝ **4 枚**（`UNCHECKED=4`、已勾 `CHECKED=6`），行号与名：

```
27:- [ ] **AC#3** `session`/`watchdog` 这种"空分母"要**响亮**：scope 校验加上
31:- [ ] **AC#5** 门禁：`bash -n` 改动脚本 rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 不降；
75:- [ ] **AC#7** `R-110-4`：票 110 承诺过"`-run TestSyncRegistryProbeLive` 经 `runtests.sh` 接进 windows 腿（step7）"，
77:- [ ] **AC#8** `R-110-3`：包匹配式缺前缀锚定（`"winsec"` 宽松串曾让我把 0 命中说成 18 次）
```

⛔ 这 4 枚框一枚都不归本腿，本腿一枚不勾。

自检基线尺（第 3 把）：

```
sh scripts/portable-tests-selftest.sh > /tmp/selftest-111r5-first.txt 2>&1; echo rc=$?
rc=0
```

输出末尾两行原文（`/tmp/selftest-111r5-first.txt:224-226`）：

```
portable-tests-selftest.sh: 32 case(s) ran, 0 assertion(s) failed
portable-tests-selftest.sh: GREEN - every seeded anomaly was refused, and the clean scope passed.
```

⇒ **32 枚 case、0 条断言失败、rc=0**，与编排者今早的读数**逐字一致**（同一锚 `15fbf18f`，无差异需具名）。
同发另两行名册读数（`first.txt:4-5`）：

```
portable-tests-selftest.sh: 254 split universe = core_pin (27) + github.com/CarlosShao/wisp/internal/winsec/acl (28 lines)
portable-tests-selftest.sh: 111 GUARD D universe = core_pin (27) + github.com/CarlosShao/wisp/internal/carrier111unclaimed (28 lines, claimed by no tier)
```

---

### §0b 本腿续作锚（`111-r5b`）

前一枚腿（`111-r5`）132 次调用后死于服务故障，本节由 `111-r5b` 在**同一目录**续写（未新建 `r5b/`）。起手复量：

| 尺 | r5b 起手读数 | 与上面骨架节的关系 |
|---|---|---|
| `date '+%Y-%m-%d %H:%M %z'` | `2026-10-06 14:19 +0800` | r5 记 `10-06 13:25`，本腿晚 54 分钟 |
| `git rev-parse --short HEAD` | `3b8873c9`（起手瞬间）；写件期间 HEAD 又推进两枚：`b1b7a770`（14:23）、`ee22fafc`（14:34） | r5 的锚 `15fbf18f` **是** `1309757b` 的父提交；`git merge-base --is-ancestor 1309757b HEAD` ⇒ 真，骨架已入库 |
| `git status --porcelain -- .github` | **0 行**（复量，与编排者 14:2x 现量一致） | ⚠ 这把尺的真含义由本腿澄清，见 §1.0：0 行＝"ci.yml 没有**未提交**改动"，**不等于**"CI 里还没有那一步" |
| `sh scripts/portable-tests-selftest.sh` | **rc=0 ／ `:224` ＝ `32 case(s) ran, 0 assertion(s) failed`** | **与 r5 逐字一致** ⇒ 判据未动、无差异需具名（"末尾两行"那一处 drift 另说，见 §4 d） |

票面框尺复量：`grep -c '^- \[ \]'` ＝ **4**、`grep -c '^- \[x\]'` ＝ **6**，与 r5 同值。⛔ 这 4 枚框一枚不归本腿，本腿一枚不勾。

---

## §1 现状：两处 grep 的改前／改后读数，以及起手先量到的一件事

### ★1.0 先记一件比"我打算怎么改"更要紧的事实：那一步**已经在库里**

派单写"ci.yml 此刻干净 ⇒ 前一枚腿没来得及改 CI，改 CI 这一步归你"。**这条推断盘上可证为不成立。**
`1309757b` 标题写着"骨架"，但它的 `git show --numstat` 是三行：

```
55	0	.github/workflows/ci.yml
1	0	.scratch/wisp/probes/111/r5/msg-skeleton.txt
49	0	.scratch/wisp/probes/111/r5/selftest-step.md
```

⇒ 那枚腿把 CI 那一步连同骨架**一起提交了**，且 ci.yml 上**删除列＝0**（只加）。本腿三个方向核过，全对上：

| 尺 | 读数 |
|---|---|
| `git diff --quiet HEAD -- .github/workflows/ci.yml` | **成立** ⇒ 工作树 ci.yml 与 HEAD 逐字节相同（＝"干净"读数的真含义） |
| `git diff --numstat 1309757b HEAD -- .github/workflows/ci.yml` | **0 行** ⇒ `1309757b` 之后无人再动 ci.yml；那一步的落地形态＝`1309757b` 交付的形态 |
| `grep -c portable-tests-selftest .github/workflows/ci.yml`（工作树＝HEAD） | **5**（＝票面要的"改后"） |
| `git show 1309757b^:.github/workflows/ci.yml \| grep -c portable-tests-selftest` | **0**（＝"改前"） |
| `git show 1309757b:.github/workflows/ci.yml \| grep -c portable-tests-selftest` | **5** |

5 处命中里**只有 1 处是可执行那一行**（`:357`），另 4 处落在该步自己的注释里 ⇒ **命中枚数不能当步数**，步数由 §3 第三把的两把尺定。

⇒ **本腿对 `.github/workflows/ci.yml` 的净改动＝0 行。** 本腿**没有**再钉一遍：同一作业加第二枚跑同一载体的步，会把 1 枚步变 2 枚、把 `lint` 作业步数尺从 11 推到 12，而票面要的只是"载体在 CI 上可跑"。
本腿射程于是收在两件事：**(1) 独立复核已落地那一步的形状与安全性，(2) 补齐门禁四把与本地两向读数**——即原指令除"落这一步"之外的全部要求。
⚠ 这不是"我没干活"的托辞：指令与盘上事实不一致时，本腿按盘上事实办事，并把**两边原文都留下**（§4 b）。

### 1.1 改前读数（两处 grep，本腿亲手复量）

```
$ git show 1309757b^:.github/workflows/ci.yml | grep -c portable-tests-selftest
0
$ grep -c portable-tests-selftest .github/workflows/slo-fresh.yml
0
```

⇒ 落地前**整个 CI 没有任何一步跑该自检**（`ci.yml` 0、`slo-fresh.yml` 0）。票面 `AC#3` 那一格记的 `grep -c … ci.yml`＝**0**，在 `15fbf18f` 那枚锚上是**现量、不是外推**。
落地后：`ci.yml` **5** ／ `slo-fresh.yml` **0**。⛔ 本腿没动也没该动 `slo-fresh.yml`——它不在"lint 作业加一步"这一格里。

### 1.2 那一步现在长什么样、落在哪个作业

| 尺 | 读数 |
|---|---|
| 作业归属 | `lint`：作业定义在 `ci.yml:65`，`:66` 即 `runs-on: ubuntu-latest` ⇒ **放 ubuntu，与指令一致** |
| 同作业 d22scan 三步 | `:74` positive control／`:83` self-test／`:108` seven-ban scan ⇒ 同作业同 runner，"三步同 runner"是现量不是推测 |
| ★`staticcheck` 那一步的 `if:` | **`ci.yml:263`**——与指令"约 `ci.yml:263`"**逐字对上**；本步抄的正是这一枚写法 |
| 新步位置 | `- name:` 在 **`:304`**、`if:` 在 **`:356`**、`run:` 在 **`:357`**（作业内最后一步；`:358` 空行，`:359` 起是 `test-core` 的注释头） |
| 新步键集合（python yaml 回读） | **恰好 `['if','name','run']`**——无 `working-directory`、无 `shell`、无 `continue-on-error` |
| `grep -c 'Portable tests carrier self-test'` | **1** ⇒ 同名步没有第二枚（"不许加两枚"的锁） |
| 是不是"写了但作业根本不跑"那一形 | `lint` 作业**无**作业级 `if:`（`:67` 直接是 `steps:`），触发器 `pull_request`／`push` 在 `:11-12` ⇒ 不是那一形 |
| ★该步在 **CI 日志**里的编号 | **step 14**（YAML 第 13 枚 ＋ GitHub 自己的 `Set up job` ＝ 14）。**换算律由票面自己的读数钉住**：票面 `AC#2 翻` 那一条写 `lint step8 gofmt／lint step12 staticcheck`，本腿按 `YAML+1` 换算改前名册 ⇒ `gofmt (gofumpt)`＝YAML 7→**8**、`staticcheck`＝YAML 11→**12**，**两处逐字对上** ⇒ 换算律成立，那一步落地后的 `13+1=14` 可直接引给取数腿 |

⚠ **引用尺子的行号会漂，本腿自己就是肇事者之一**：本腿在那张票 `:307` 追加了一行 Progress log ⇒ 上面那两条"票面自己记的"读数，行号从 **`:306`（HEAD）→ `:307`（追加后）** 那一段往后各移 1 行（含 `AC#2 翻` 那条：HEAD＝`:313`／追加后＝`:314`；编排者那一节标题 HEAD＝`:308`／追加后＝`:309`）。
⇒ **本件引用票面一律以文字为准、行号只作辅读**；追加点上方的 `:305-306` 那两行（`lint 的 10 枚步名里只有 staticcheck 与 mockllm module vet 带 !cancelled()`）**逐字未动**，本腿实测 `diff` HEAD 与工作树的 `sed -n '305,306p'` ＝ **0 行差异**。⛔ 追加前后 `- [ ]`＝4／`- [x]`＝6 未变，且 `:306` 之后**没有任何一行框**被本腿挪动（`awk 'NR>306 && /^- \[[ x]\]/'` ＝ 0 命中）。

新步可执行三行原文（含缩进；`:305-355` 是那 51 行注释，本腿**一字未碰**，全文见 `git show 1309757b -- .github/workflows/ci.yml`）：

```yaml
        if: ${{ !cancelled() }}                    # :300  （mockllm module vet，前一步）
        run: go vet ./...                          # :301
        working-directory: tools/mockllm           # :302
                                                   # :303  空行
      - name: "Portable tests carrier self-test (ticket 111 AC#3 - GUARD D's positive control)"   # :304
        ...                                        # :305-355  本步注释块（51 行，逐行 `        #` 起）
        if: ${{ !cancelled() }}                    # :356
        run: bash scripts/portable-tests-selftest.sh   # :357
                                                   # :358  空行
  # ----------------------------------------------------------- test-core ----   # :359
```

⇒ 6 空格 `- name:` ＋ 8 空格 `if:`／`run:`，与同作业其余各步一致；yaml 能整份吃下（§3 第三把）。

**一枚独立佐证（不是本腿量的，是票面自己记的）**：`111-ci-tests-20-of-33-packages.md:305-306` 写着"lint 的 **10 枚**步名里只有 `staticcheck` 与 `mockllm module vet` 带 `!cancelled()`"。
本腿的改前步数尺读到的正是 **10**（§3 第三把 c），而带 `!cancelled()` 的那两枚在 `:263`／`:300`——**两把不同来源的尺对上了同一枚数**；那一步落地后 `lint` 变 11 枚、带 guard 的变 3 枚。

### 1.3 ★指令里那把 go 尺：静态命中数与"会不会 exec 真 go"是两回事

照抄指令的尺跑了一遍：

```
$ grep -nE '(^|[^a-zA-Z])go (test|list|build|vet|run)' scripts/portable-tests-selftest.sh
9 行命中（grep -cE 同尺亦 9）
```

**命中 9 行，不是 0**。逐条分类（判据＝剥掉前导空格后是否以 `#` 起）：

| 行 | 类别 |
|---|---|
| 6, 8, 25, 288, 396, 583, 775 | 注释（7 行），如 `# WHY IT EXISTS. The defect was: `go list` progress text…` |
| **371, 381** | **不是注释，是断言字符串**：`has "GUARD C - go list's stdout carried 1 line\(s\) that are not"` 及其 `10 line\(s\)` 同形 |

⇒ 该尺真读数：**7 行注释＋2 行字符串，0 行是被执行的 go 命令**。指令那句"正文命中只有注释"在 `:371/:381` 上**不成立**（差 2 行）。
更要紧的是这把尺**问的不是 CI 上要命的问句**：它量"源码里有没有 go 这个词组"，而 CI 上的风险是"**运行时会不会 exec 出真工具链**"。这两件事在这份脚本上恰好**反着**——正文里 `go` 这个词出现 **61 次**（`grep -oE '\bgo\b' … | uniq -c`），却一次真 go 都不碰。本腿换了一把能问对问题的尺：

**★绊线实测**（把一枚会自报家门的假 `go` 塞在 PATH 最前，再用**与 CI 那一步逐字相同**的命令跑一遍）：

```
$ cat /tmp/golog-111r5b/go
#!/usr/bin/env bash
echo "TRIPWIRE-EXECUTED $(basename "$0") $*" >> /tmp/golog-111r5b/hits.txt
exit 123
$ PATH="/tmp/golog-111r5b:$PATH" bash scripts/portable-tests-selftest.sh
rc=0 elapsed_s=170
=== tripwire hits ===
NO HITS FILE = the tripwire go was never executed
```

⇒ **rc=0、`:224` 仍是 `32 case(s) ran, 0 assertion(s) failed`、末尾仍是 GREEN**，而那枚绊线 `go` **一次都没被 exec**（连 `hits.txt` 都没被创建）。
"接进 CI 不会跑测试、不会和写腿抢 go"这条**从此不再靠 grep 词面外推**。机制也读到，与读数对得上：

- 载体在 `:118-120` 把 `scripts/testdata/portable-tests/go` 拷成 `$work/bin/go`，**每一处**跑被检脚本的子 shell 都自带 `PATH="$work/bin:$PATH"`——`:158 :250 :335 :555 :600 :608 :743 :755 :846` 共 **9 处，无漏网**；
- 假 go：`cmd=${1-}` 在 `:41`、`case $cmd in` 在 `:44`、`test)` 支在 `:148`、兜底 `*)` 在 `:176` → **`exit 99`（`:178`）**、`esac` 在 `:180`。`test)` 支只从 `$FAKEGO_PIN_FILE` 打印合成的 `ok \t<pkg>\t0.01s` 行；**全脚本 `exec `／`command -v`／`which go` 命中数＝0** ⇒ 没有任何回落真工具链的路径；认不出的子命令一律响亮 `exit 99`，不静默放行；
- 唯一会下真 `go test` 的那行是 `scripts/portable-tests.sh:694`（`sh "$strict" …` → `tools/d22scan/runtests.sh:75` `go test -v -count=1 "$@"`），而 **`--scope=census` 在 `portable-tests.sh:424` 就 `exit 0`，走不到 `:694`**；本自检调到的 scope 只有 `census/core/winsec/winsecfoo`（`7/5/9/2` 枚），`core`／`winsec` 两支又都跑在假 go 之下。

---

## §2 本腿的改动

**净改动：`.github/workflows/ci.yml` ＝ 0 加 0 删。** 理由与证据在 §1.0：那一步已由 `1309757b` 落地，且 `1309757b..HEAD` 无人再动 ci.yml。
本腿**没有**新增第二枚步、**没有**挪位、**没有**改那 51 行注释、**没有**碰 `gofmt (gofumpt)` 那枚常红步（`:168`／`:184`，票 236 AC#6 并案地界）。

"为什么这一步该放在 `lint` 作业这个位置"——这一格由落地那枚腿在它自己的注释 `:342-355` 里答了。本腿核过它讲的三件事**都真**：

1. `if: ${{ !cancelled() }}` 抄的是同作业 `staticcheck`（`:263`）与 `mockllm module vet`（`:300`）已有的写法，⛔ 不是 `always()`、⛔ 不是 `continue-on-error`（`ci.yml:246-250` 那一段把"`!cancelled()` 是票 85 裁出来唯一不改自己红/绿、只保证这一步答话的形状"写在门的正文里，本步与它逐字同形）；
2. `lint` 作业**今天就有会红的步在上面**（`gofmt (gofumpt)` `:168`、`staticcheck` `:213` 是本仓在册常红），没带 `if:` 的新步会**出生即无 CI 读数**——正是票 111 AC#5 点名的洞；
3. 用 `bash` 不用 `sh` 是**量过的**：本腿独立复跑 `dash -n scripts/portable-tests-selftest.sh` ⇒ **rc=2，报在 `:246` `Syntax error: "(" unexpected (expecting "}")`**，与注释写的行号与报错形状**逐字一致**；同尺下 `scripts/portable-tests.sh`、`scripts/winsec-tests.sh` 亦 **rc=2**，只有 `tools/d22scan/runtests.sh` 是 **rc=0**。
   本机 `/bin/sh` 恰是 bash（`sh --version` ⇒ `GNU bash, version 5.2.37(1)-release (x86_64-pc-msys)`），所以 `sh <载体>` 在本地也绿——**这就是为什么形状必须拿 dash 直量**；ubuntu-latest 的 `/bin/sh` 是 dash，写 `sh` 会让这一步**死在语法、一枚 case 都判不到**。

⚠ 有一处**本腿不跟随该注释的措辞**（不改它，只留读数）：`:339` 写"被检的（portable-tests.sh、winsec-tests.sh、runtests.sh）是真码在跑"。真实形状是 `portable-tests.sh:694` 用 **`sh`** 起 `runtests.sh`——**恰好只有这一枚 dash-clean**，所以方向安全；反过来说，若哪天有人把 `runtests.sh` 写成 bash-only，这一步会在 ubuntu 上以**语法错**红（而不是以它想报的那个 finding 红）。属未被本步覆盖的一处**形状耦合**，记 §4 g，本腿不修（修它要动 `scripts/**`，越界）。

本腿对"这一步该被相信"提供的增量证据＝§3 四把尺 ＋ §1.3 那次绊线。

---

## §3 门禁四把（全部本腿亲手跑，原文抄）

### 第一把：改前／改后各一发 `sh scripts/portable-tests-selftest.sh`

⚠ 先讲清这把尺在"本腿净改动 0 行"前提下**怎么读**：派单要的形状是"同一把尺在两向各量一次且对得上"。本腿把**四发**一次摆全（含 r5 那枚腿的基线），**判据一字未动、未改任何尺的措辞**。

| # | 起跑锚 | 启动器（命令原文） | rc | `:224` 计数行 | 末行 |
|---|---|---|---|---|---|
| 1 | `15fbf18f`（r5 骨架，＝改前内容） | `sh scripts/portable-tests-selftest.sh` | 0 | `32 case(s) ran, 0 assertion(s) failed` | GREEN |
| 2 | `3b8873c9`（本腿起手） | `sh scripts/portable-tests-selftest.sh` | **0** | `portable-tests-selftest.sh: 32 case(s) ran, 0 assertion(s) failed` | GREEN |
| 3 | `b1b7a770`（ci.yml 内容与 `1309757b` 逐字节同） | `bash scripts/portable-tests-selftest.sh`（＝CI 那一步逐字命令） | **0** | 与 #2 逐字一致 | GREEN |
| 4 | 同上 | 同上＋绊线假 `go` 置于 PATH 最前 | **0** | 与 #2 逐字一致 | GREEN |

第 2／3／4 发末尾三行原文（三次输出**逐字节相同**；本腿实测 `diff` #2 与 #3 只在 `== case <name>: rc=N log=/tmp/tmp.<X>/….log` 那几行的**临时目录名**上报行，计数行与 verdict 行不进 diff）：

```
portable-tests-selftest.sh: 32 case(s) ran, 0 assertion(s) failed
portable-tests-selftest.sh: carrier scratch kept at /tmp/tmp.<X> (rules: temp files are created, not deleted)
portable-tests-selftest.sh: GREEN - every seeded anomaly was refused, and the clean scope passed.
```

⇒ **两向一致：rc=0、32 枚 case、0 条断言失败**，与 r5 骨架记的数**逐字相同 ⇒ 无差异需具名，未改判据凑数**。
⚠ "末尾那两行"要打折读：现量末尾两行是 `scratch kept` ＋ `GREEN`，**case 计数行在 `:224`**；r5 骨架写的"`:224-226`"里 224 正是那枚计数行，它把三行说成了两行。drift 成因本腿只核到一半 ⇒ **§4 d**，不替它编解释。

### 第二把：`sh scripts/d22scan.sh`（纯净快照 rc=0、各 scope 不降）

| 跑在哪 | 命令 | rc | 结论行 |
|---|---|---|---|
| 工作树（脏树。全仓 `git status --porcelain \| wc -l`＝**747 @14:19 起手 / 744 @14:5x 复量**——⚠ 这把尺**在漂**，因共树里别的腿正在提交；本行只是"这不是干净快照"的注记，⛔ 不是判据，两把原文都留下） | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations` |
| **纯净快照** `git archive HEAD \| tar -x -C /tmp/wisp-111r5b-clean` | `sh scripts/d22scan.sh` | **0** | 同上（AC#5 点名要的正是"纯净快照"这一发） |

各 scope 两发对照：

```
bans #1-5 internal/=228   bans #1-5 cmd/=38   ban #6 frontend/=85   ban #7 internal/tools/=23
ban #8 design/ = 39 工作树 / 30 快照          ban #8 frontend/=85   ban #8 internal/=513   ban #8 cmd/=104
```

⇒ **不降**：`internal/`=228、`cmd/`=38、`frontend/`=85、`internal/tools/`=23、ban8 `internal/`=513、ban8 `cmd/`=104 —— 两发**逐字相同**。
唯一差集是 `ban #8 design/`（工作树 39 ／快照 30，差 9 枚文本文件）。**这一枚与本腿无关、也不是"降"**：本腿没动 `design/**`（本件只记这一枚计数，⛔ 不转述其内容）；`git archive HEAD` 只装已跟踪内容，未跟踪件不进快照 ⇒ 方向是**工作树扫得更多（39 > 30）**，覆盖更宽而非门更松。
另留一发**本仓在册的机器相关读数**：快照那发多印一行
`d22scan: gitignore rules NOT APPLIED - git cannot be consulted in …: fatal: not a git repository (or any of the parent directories): .git`
（原因：`git archive` 不带 `.git`）。这是 d22scan **自报的响亮方向**，原文紧接 `every path in every scope is being scanned … This is the loud direction: a scanner that cannot ask git which paths are tracked does not get to skip any (A207's machine-dependent denominator)`，rc 仍 0；工作树那发不报此行。⛔ 本腿不据此判任何事，只留读数（它对应台账 `A207` 那一族）。

### 第三把：★YAML 形状自证（本仓为"CI 一步写坏但静默不跑"付过学费）

| # | 尺 | 命令原文 | 读数 |
|---|---|---|---|
| a | 有没有未提交的 ci.yml 改动 | `git diff --numstat -- .github/workflows/ci.yml` | **0 行** ⇒ 本腿对 ci.yml **无改动**，"删除列必须为 0"这条**天然满足**（ci.yml 根本不出现在 numstat 输出里，没有可停之物） |
| b | 落地那一笔的删除列 | `git show --numstat 1309757b` 的 ci.yml 行 | **`55  0`** ⇒ 删除列＝**0**，与"只加一步"相符（⚠ 加的是 55 行但**步只有 1 枚**，由 c/d 定） |
| c | ★**步数尺**（"步数没写坏"那把） | `awk` 数 `lint` 作业区间内 `^[[:space:]]*- name:` 的枚数；改前喂 `git show 1309757b^:.github/workflows/ci.yml`，改后喂工作树 | **改前＝10 枚** ／ **改后＝11 枚** ⇒ **差 1 枚，正是本线程涉及的那一步**；改前名册 `:74 :83 :108 :136 :168 :184 :206 :209 :213 :289` **逐行保留**，改后原样再加 `:304` ⇒ 次序未重排、无既有步被吞 |
| d | yaml 吃不吃得下＋解析后步数 | `python` `yaml.safe_load` 后数 `jobs.lint.steps` | 解析**通过**；`lint.runs-on == ubuntu-latest`；**13 枚 step** ＝ 11 枚有名步 ＋ `:68 uses: actions/checkout@v4` ＋ `:70 uses: actions/setup-go@v5`（后两枚**没有** `name:` ⇒ 与 c 的 11 **不矛盾**，两把尺量的是不同东西）；`name.startswith('Portable tests carrier self-test')` 的枚数＝**1** |
| e | 有没有把门拆了 | python 回读新步 | keys＝`['if','name','run']`；`if`＝`'${{ !cancelled() }}'`；`run`＝`'bash scripts/portable-tests-selftest.sh'`；**`'continue-on-error' in step` ＝ False** |
| f | 载体语法（AC#5 的 `bash -n`） | `bash -n scripts/portable-tests-selftest.sh`／`scripts/portable-tests.sh`／`scripts/winsec-tests.sh` | **rc=0 / 0 / 0** |
| g | checkout 在不在（没它 `scripts/` 是空的） | `sed -n '68,72p'` | `:68 - uses: actions/checkout@v4`、`:70 actions/setup-go@v5` ⇒ 新步在同一作业内跑，载体文件必然在场 |

⛔ 本腿**没有**用 `always()`、**没有**用 `continue-on-error`、**没有**动 `gofmt (gofumpt)`、**没有**动 ci.yml 里任何一行既有注释。

### 第四把：这一步会不会拖时间／CI 上缺不缺工具

| 尺 | 读数 |
|---|---|
| ★**耗时**（命令与 CI 那一步逐字同） | `bash scripts/portable-tests-selftest.sh` ⇒ **rc=0、158467 ms（≈2 分 38 秒）**；绊线那发 `elapsed_s=170`。同时间本机另有别的腿在跑 ⇒ ⚠ **这不是 CI 侧读数**，只是"本地不是一秒钟的门"的量级证据 |
| 时间敏感形状 | `grep -n '\bsleep\b\|\btimeout\b'` 在载体／`portable-tests.sh`／`runtests.sh` **三处全 0 命中**；载体也不用 `SECONDS`／`date +%s` ⇒ 没有"用墙钟时间差实现超时"（D22 禁的那一枚） |
| ★**外部命令依赖**（剥掉注释行后跨 5 个文件取词：载体＋`portable-tests.sh`＋`winsec-tests.sh`＋`runtests.sh`＋假 go） | 出现的只有 `printf / grep / sed / awk / tee / cat / wc / sort / uniq / tr / cut / mktemp / bash / sh / cp / mv / chmod / mkdir / rm / dirname / tail / diff / go` ⇒ **纯 coreutils＋bash**。其中 `go` 那 87 次命中全部落在**假 go 的覆盖范围或注释**里（见 §1.3 三条机制）；`git` 只在**注释**里出现（`selftest:48 :49 :74`，全是给人看的 `git show d253703a^:…` 复现指引），**运行路径上零依赖**（`portable-tests.sh`/`winsec-tests.sh`/`runtests.sh` 里 `\bgit\b` 命中 **0**）；`timeout / tput / sqlite3 / python / jq / curl / wget / dot` 全 **0** |
| 文件模式（ubuntu checkout 只认 git 里那一位） | `git ls-files -s` ⇒ `scripts/portable-tests-selftest.sh`、`scripts/portable-tests.sh`、`scripts/testdata/portable-tests/go` **三枚都是 `100644`（无 x 位）**。这一步安全：`run: bash <脚本>` 不需要 x 位；载体在 `:118-120` 先 `cp` 假 go 再 `chmod +x "$work/bin/go" 2>/dev/null || true`，且**全仓没有任何一处直接 exec `scripts/testdata/portable-tests/go`**（`grep -nE 'testdata/portable-tests/go' scripts/*.sh` 只命中注释与 `shim=` 赋值行）⇒ 不存在"ubuntu 上 Permission denied"那条路 |
| 绝对路径耦合 | `grep -n '/d/work\|/home/runner\|C:/' scripts/portable-tests-selftest.sh` ⇒ **0 命中**；载体 `:86-88` 用 `dirname $0`＋`cd` 自算 `$root`，断言里没有一枚机器绝对路径 ⇒ 换 runner 不脱靶 |
| **本腿结论** | 没在这一步里看到"runner 缺工具"的形状：依赖面只有 coreutils＋bash，ubuntu-latest 两样必有；不跑真 go、不下模块、不建测试二进制。⚠ 但这一格的**正解不在本腿手里**，见 §4 a／f：**"本地判不出事"和"ubuntu-latest 上真跑绿"是两件事**，本腿不合并它们。 |

---

## §4 判不动的地方（每条具名，⛔ 不自行拍）

**a) ★AC#3 第二半"在 CI 上真响过一次"——本腿给不出这行读数，也不宣称 AC#3 完成。**
载体（组 23 的标题注释起在 `selftest:766`，`if selected census-unclaimed-package-goes-red` 在 `:791`，`want_rc 1` 在 `:795`，`has 'GUARD D - 1 package\(s\) compile a test file for GOOS='` 在 `:796`，`has 'census totals: packages=… unclaimed-with-tests=1'` 在 `:800`）此刻**在 CI 上可跑**了——这是本腿唯一能交的东西。
它**跑没跑、跑成什么颜色**只有推送之后的 run 能回答 ⇒ 归编排者取数。
票面那句"缺的那一行读数＝**CI 里一步真跑该 selftest 并让它为这枚正控红/绿各一次**"**仍未收**；`AC#5`（原文 `这一格依赖 AC#3 那枚正控进 CI`）与 `AC#8`（原文 `与 AC#3 同因（selftest 不在 CI）`）同判：**变的是"不在 CI"那半，没变的是"CI 上响过"那半**。⛔ 本腿一枚框不勾、不改口径。

**b) ★指令与盘上事实冲突：那一步不是我加的。**
派单："ci.yml 此刻干净 ⇒ 那枚腿没来得及改 CI，改 CI 这一步归你"。盘上："工作树干净 **并且** ci.yml 里那一步已在 `1309757b` 落地"。两把尺原文在 §1.0。
本腿处理＝不加第二枚、不改那一枚、独立复核实。
留给编排者裁定的是：**已死腿的 `1309757b` 该按"骨架"记还是按"落地步"记**——它的 commit 标题写"骨架"，它的 numstat 写 `55 0 .github/workflows/ci.yml`，**名与实不符**。⛔ 本腿不改写已入库历史（AGENTS.md §1.4：要更正只能追加），也不重述那一笔的归属。

**c) ★同一把尺的第二遍读数与指令预期不同，两把原文都留下，⛔ 不挑一枚为准。**
指令："正文命中只有注释"。实测：该尺 **9 行命中**，`:371`、`:381` 两行**不是注释而是断言字符串**。
第二把尺（能问对问题的）：绊线假 go 置 PATH 最前、跑逐字同命令 ⇒ **rc=0、32 case、0 失败、`hits.txt` 从未被创建**。
⇒ 两把尺**指向同一个结论**（这一步不 exec 真 go），但**不是同一把尺**。本腿不把第一把"修正"成 0 命中，也不因第一把的字面不成立就否认第二把。
**建议（不是本腿能定的）**：今后引用这一步的安全性时以**绊线那次**为准，词面 grep 只作辅读。⛔ 判据本身一字未动，这一条只涉及"该用哪把尺说话"，请编排者按**尺子写法**登记，不要读成判据变更。

**d) "末尾两行"在两枚腿之间不一致，本腿只追到一半，⛔ 不替它编因。**

| 腿 | 时刻／锚 | 它记的 `tail` | 本腿现量同一发应读到什么 |
|---|---|---|---|
| r5 | 13:4x／`15fbf18f` | `32 case(s) ran…` ＋ `GREEN…` 两行 | — |
| r5b | 14:2x 起／`3b8873c9`、`b1b7a770` | — | `carrier scratch kept…` ＋ `GREEN…`，**计数行在 `:224`** |

已核的一半：那行 `carrier scratch kept` 由 **`6307e369`（10-02 09:19，`ticket 250 AC#1/AC#2: keep go list's stderr out of portable-tests.sh's denominator`）**引入，`git log -S 'carrier scratch kept'` 只回这一枚，且 `6307e369` **早于**两枚腿的起手时刻 ⇒ **不是 r5 或本腿动脚本造成的**。
⛔ 未核的一半：按这个日期，r5 那发**本该**也把 `scratch kept` 读在倒数第二位——它没有。本腿**没有** r5 那发的原始日志（`/tmp/selftest-111r5-first.txt` 是否仍在盘上，本腿未查，且它不是本腿该写的面），无法区分三种可能：① 它引的是 `sed -n '224,226p'` 而非 `tail -2`（它自己写的行号 `:224-226` 更像这一种）；② 它跑在另一枚内容不同的锚上；③ 别的。
⇒ **两把原文都留在上表**。不受影响的结论：**case 数与断言数（32 枚／0 条）两枚腿逐字一致**，第一把尺的两向对齐不靠这三行成立。⛔ 本腿不改写 r5 那节的任何字。

**e) ★AC#5 要求的"非实现者复跑件"不在本腿射程，本腿不冒充它。**
票面 `AC#5` 原文：`bash -n／sh scripts/d22scan.sh 那两把我自己的门禁尺今天**没有非实现者的复跑件**……等 111-r5 一并取`。
本腿交的 §3 那两把是**写码腿自证**，形状上正是票面说"缺非实现者复跑"的那一枚东西——它**不能**被读成 AC#5 已闭。⇒ AC#5 需要的是**另一枚腿**复跑 §3 第二把与第三把 f 行；本腿不替它勾，也不因为它没做就把自己的读数写轻。

**f) ★"这一步会不会缺工具"：本腿在本地判**不出红**，但这是**本地**读数，⛔ 不是 ubuntu 读数。**
派单要求："若你判断该自检在 CI runner 环境会因缺工具而红：**不修**，把这条具名写进判不动节——'接上去但环境不满足'和'根本没接'是两种不同的债。"
本腿的判断是**倾向不会红**，理由全部在 §3 第四把（依赖面＝coreutils＋bash；`git` 零运行依赖；无 x 位需求；无绝对路径；`bash -n` 与 `dash -n` 两向都量过；假 go 覆盖全部子命令、未知即 `exit 99`）。
⛔ 本腿**不**把这条写成"CI 上不会红"。三处只在未来才响的残余，具名留这儿：
  1. **`/bin/sh` 的分布**：本步写 `bash …` 已经避开主风险，但 `portable-tests.sh:694` 用 `sh "$strict"` 起 `runtests.sh`——今天 `runtests.sh` 是唯一 dash-clean 的一枚（§2 第 3 条实测），一旦有人把它写成 bash-only，这一步会在 ubuntu 上以**语法错**红，而**这枚红不是它想报的那个 finding**。本腿不修（要动 `scripts/**`，越界）。
  2. **≈2 分 38 秒**（本地）落进 `lint` 作业：本腿**没有** ubuntu 侧耗时读数，不能断言它在 CI 上也在这个量级；若编排者取数发现它显著更慢，那是新问题，不归本腿现在改判据。
  3. **`mktemp` 落点**：**三处**同形回退——`:116`（`work=$(mktemp -d 2>/dev/null || echo "$root/.portable-selftest.$$")`）、`:106` 与 `:223`（`shadow=$(mktemp -d 2>/dev/null || echo "$root/.portable-shadow.$$")`）。若 runner 上 `mktemp` 失败，回退路径落在**仓库根**里（`.portable-selftest.$$` / `.portable-shadow.$$`），会在工作树留残留件（载体自己那条"temp files are created, not deleted"的规矩更是如此）。⛔ 本腿**没有**判它会不会发生（本机 `mktemp` 显然可用，四发都写在 `/tmp/tmp.*`），只记这条形状存在，以及"若在 CI 上真发生，它会以**未跟踪残留件**出现而不是以红出现"——这一枚是取数时才响的。

**g) ★一处本腿读不出来的东西：`design/` 的 9 枚差集。**
§3 第二把里 `ban #8 design/` 两发不同（39／30）。本腿**只读计数、不读内容面**（`design/**` 零读零写零转述是派单地界），因此**无法**自证这 9 枚全是未跟踪件、也不怀疑这个数——它由"工作树比 `git archive HEAD` 多出未跟踪内容"这一条机制解释，且方向（39 > 30）与"覆盖更宽"同向。⇒ 若编排者要把它当"纯净快照的 scope 降了"来处理，**先复量**，不要以本件为凭。

**h) ★票面 `:306` 那句"注释写在门的正文里"的出处，与本腿被引到的位置不同。**
派单说"票面 `:306` 具名'注释写在门的正文里'"。本腿用该词面全 `issues/` 搜，**唯一命中在 `236-six-cells-that-only-surface-at-the-reading-layer.md:96`**，原文：
`⛔ 注释写在门的正文里，改它会被读成"动门的形状"（先例逐字在 ci.yml:146-151）。⇒ 我裁：AC#4 与 AC#6 错开、由非实现者写腿各做各的`
——那句话**确实归属票 236 AC#4**，与本腿被引到的地界**判断一致**，只是**行号不在票 111 的 `:306`**（票 111 `:305-306` 是另一枚内容：`lint` 作业那 10 枚步名里只有两枚带 `!cancelled()` 的读数，本腿已在 §1.2 末把它当佐证用了）。
本腿先按"中文词面"量到 0 命中、**没有**就此下结论（这正是本仓今天那次"用中文词面 grep 纯英文名册"的错形），改成搜词面真名后才定位到 236:96。⛔ 本腿不改任何票面行号引用，只把两边出处都留下，请编排者按 236:96 认这条地界。

---

## §5 本腿动了哪些面（逐枚点名，⛔ 别的面一字未写）

| 面 | 本腿动作 |
|---|---|
| `.github/workflows/ci.yml` | **0 写**（0 加 0 删；那一步属 `1309757b`，见 §1.0／§4 b） |
| `.scratch/wisp/probes/111/r5/selftest-step.md` | **续写本件**（在 r5 同一目录，⛔ 未新建 `r5b/`；r5 已写的 §0 内容**一字未改**，本件只把最后那行"（骨架节，§1–§4 后续补齐。）"换成 §0b 起的实内容） |
| `.scratch/wisp/issues/111-ci-tests-20-of-33-packages.md` | **只在末尾追加 1 行 Progress log**；⛔ 框一枚未动（追加前后 `^- [ ]`＝**4**、`^- [x]`＝**6** 不变） |
| `scripts/**`／`docs/**`／`cmd/**`／`internal/**`／`tools/**`／台账 `pending-and-issues.md` | **零写**；`internal/**` 也**未读工作树内容面**（本件里所有 `scripts/**` 引用都是**行号＋单行原文**尺，不是通读） |
| `frontend/**`／`design/**` | **零读零写零转述**（§4 g 只有 `design/` 的**计数**，无内容） |
| `.scratch/wisp/probes/236/**` | **未读**（派单禁面）；本件不含其任何结论 |
| go 命令 | **零次**（`go test/build/vet/run/list/env` 均未出现；本腿跑的是 `sh scripts/*.sh`、`bash scripts/*.sh`、`dash -n`、`bash -n`、`git …`、`awk/grep/python` 尺） |
| `-tags winlive` | **未加**（owner 未批） |

## §6 一句话：这一步为什么值钱

票 111 **AC#3** 逐字要两半：① GUARD D 守卫本体（已落地：`scripts/portable-tests.sh:360`／`:396`／`:407-424`，`:409` 打印 `GUARD D - $guardd package(s) compile a test file for GOOS=$goos and`、`:422` 那一支 `exit 1` 真在）；② "**并人为抽掉一个包证明它会红**"——第②半的载体**盘上一直有**（组 23 `census-unclaimed-package-goes-red`），⛔ 但 **CI 从不跑它**（改前 `grep -c portable-tests-selftest .github/workflows/ci.yml`＝**0**，本腿 §1.1 复量）。
⇒ 本件所核的那一步把**载体接到 CI 上**了：同一枚原因还卡着 **AC#5** 与 **AC#8**（票面自己写了这两格与 AC#3 同因）。
⛔ **本腿不宣称 AC#3 完成**——交付物只有"**载体在 CI 上可跑**"＋本地两向读数；"它在 CI 上真响过一次"归编排者推送后取数（§4 a）。

## §7 编排者下一步（本腿不替你做）

1. 取数：推送后的 run 里 `lint` 作业那一步的**颜色与日志**——CI 日志编号应为 **step 14**（YAML 第 13 枚 ＋ `Set up job`；换算律由票面 `AC#2 翻` 那一条的 `step8 gofmt`／`step12 staticcheck` 两枚读数钉住，该行 **HEAD＝`:313`／本腿追加后＝`:314`**，见 §1.2 末行的漂移注记），步名逐字
   `Portable tests carrier self-test (ticket 111 AC#3 - GUARD D's positive control)`。**只有这个能收 AC#3 的第二半。**
2. 裁定 §4 b（`1309757b` 名实不符）与 §4 c（该用哪把尺说话）；§4 d 那处 tail drift 若要追因，需要 r5 那发的原始日志文件。
3. AC#5 的"非实现者复跑件"仍未闭（§4 e）：请派一枚只读腿复跑 §3 第二把／第三把 f 行。
4. ⛔ 框仍归编排者：本腿未勾任何一枚，票面 4 枚零勾框（AC#3／AC#5／AC#7／AC#8）数量与状态未变。
