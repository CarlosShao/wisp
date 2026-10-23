# 票 362 收尾腿 `362-r3` —— gofumpt 步注释里的三个分母数：改选〔乙〕「不写数 + 一条现推命令」

- **腿**：`362-r3`＝写腿（实现者）。⛔ 翻框／⛔ 改票名·`Status:`·`-done`／⛔ 改台账／⛔ push／⛔ 动任何 `.go`（本笔 **零 `.go` 字节**）。
- **写面**：`.github/workflows/ci.yml`（**只**那道 gofumpt 步的注释）＋本件。临时件只建不删：`D:\tmp\wisp\20261023\362r3\**`。
- **锚**：读数钉在 **`0754fa7d9`**（`362-r2` 第三笔）的**仓外 `git archive` 导出树**；⛔ 主检出（有别人未提交改动）。起手 HEAD `b22225513`。

## §0 开工前三条硬检查（现跑，⛔ 信票面）

| 检查 | 逐字命令 | 读数 |
|---|---|---|
| 索引 | `git diff --cached --name-status` | 仅 `R072`（票 362 票面 → `-done`，**别的腿暂存的**）；**无 `D`／`A`** ⇒ 通过、不动它 |
| `ci.yml` 脏否 | `git status --porcelain -- .github/workflows/ci.yml` | **0 行**；`git diff`／`git diff --cached` 皆 0 ⇒ 无"别人的改动" |
| 行号（自找） | `grep -n 'name: gofmt (gofumpt)$' ci.yml` | HEAD **`:217`**、改后 **`:235`**；三个数在 **`:191-193`**（票面的 `:191` 是**数所在行**） |

## §1 改前那三个数（逐字，`git show HEAD:.github/workflows/ci.yml | sed -n '190,193p;212p'`）

```
      # tracked .go. Measured at the 362-r2 anchor (2026-10-23, re-derived with
      # `git ls-files`, not carried over): of 1162 tracked .go files, 802 sit in
      # those five directories (internal 623, cmd 149, tools 18, scripts 11,
      # frontend 1), 360 sit under .scratch/ and 0 sit at the repo root, and no
      # (.scratch/**), so the two steps agree on the bytes today (802 each) by two
```

⇒ 三数＝**`1162` ／ `802` ／ `internal 623`**（另附 `149/18/11/1`、`.scratch 360`、根 `0`）。**`802` 在 `:212` 还有第二处**（见 §8②）。

## §2 我在 `0754fa7d9` 树上现量的真数（逐字命令 + 输出）

```sh
unset GOTMPDIR TMPDIR TEMP TMP; export GOFLAGS=
git archive --format=tar 0754fa7d9 -o /d/tmp/wisp/20261023/362r3/0754fa7d9.tar   # rc=0
mkdir -p /d/tmp/wisp/20261023/362r3/tree-0754fa7d9 && tar -xf /d/tmp/wisp/20261023/362r3/0754fa7d9.tar -C /d/tmp/wisp/20261023/362r3/tree-0754fa7d9   # rc=0
cd /d/tmp/wisp/20261023/362r3/tree-0754fa7d9

find . -type f -name '*.go' | wc -l
# -> 1165                          （分母①：total tracked .go；导出树无 untracked）
find . -type f -name '*.go' | sed -e 's|^\./||' -e 's|/.*||' | sort | uniq -c | sort -rn
# ->  626 internal / 360 .scratch / 149 cmd / 18 tools / 11 scripts / 1 frontend
find internal cmd tools scripts frontend -type f -name '*.go' | wc -l
# -> 805                           （分母②：handed to gofumpt）
find internal -type f -name '*.go' | wc -l
# -> 626                           （分母③：internal）

/d/work/base/gopath/bin/gofumpt.exe -l internal cmd tools scripts frontend > gofumpt-0754fa7d9.txt 2>&1
rc=$?                              # ⛔ 不走管道
# -> rc=0 ; lines=0 ; bytes=0      （门本身在那棵树上绿）
```

- **三数齐全且与总分对上**：`626+149+18+11+1 = 805`；`805+360 = 1165`；根 `0`、其它首段 `0`。
- **独立第二把尺**（⛔ 同一算法自证）：`git ls-tree -r --name-only 0754fa7d9 | grep -cE '\.go$'` ⇒ **1165**；`^internal/` ⇒ **626**；五目录 ⇒ **805**、`^\.scratch/` ⇒ **360**。两把尺同值。
- **相邻那条 `20` 复量**（不属三个数、未动）：同树 `gofumpt -l .scratch` ⇒ **20 行**，`grep -vcE '^\.scratch[\\/]'` ⇒ **0**（反斜杠输出，按 `[\\/]` 判归属）⇒ 那句仍真。
- **射程声明**：只跑 `run:` 块里判据那行，⛔ `go install …@latest`（要网、会装"当时最新"）；用本机 `/d/work/base/gopath/bin/gofumpt.exe`。

**锚点归属（⛔ 别只做减法）**：直接量树，`938216b8e` 的树 ＝ **`1162 / 802 / internal 623`、`.scratch 360`** —— 与注释**逐字相同** ⇒ 那三个数**确实是从 `938216b8e` 量的**；而写它们的提交 `ade547ac8`（`ci.yml` blob 在这一笔由 `ff1711e1` 变 `5d856a83`）自己的树是 **`1165/805/626`**。另：`0754fa7d9` 的 `ci.yml` blob 与 HEAD 同为 `5d856a83` ⇒ 我改的正是 `362-r2` 出厂那一版。

## §3 选了〔乙〕：⛔ 写具体数、只写"怎么现推"，为什么它最有价值

**选〔乙〕**（派单括号里我倾向的那条也正好是乙，但我按尺重判了，不是顺着走）：

1. **那段话的职责是"定义"不是"规模"**。它要读者记住的是"这道门的分母＝五个产码目录，不是全部 tracked `.go`"；`1162/802/623` 只是**装饰性规模**。命令能同时**给出**三个数并**重新证明**"没有第三种首段"这个结构断言 —— 写死数只能**断言**它。
2. **写死数在本仓已连错两次、且第二次是刚写完就错**：`362-r1` 的 `1160/800/621` 被 r2 覆盖；r2 的 `1162/802/623` 在**它自己的提交上**已不成立（§2 量得 `1165/805/626`）—— 不是它造假，是它量的树比写它的提交旧 11 笔（§2）。**2/2 手抄数都错**，这就是这条路的实测分布，不是假设。
3. **甲的成本落在读者身上**：写"数 @ commit"要求读者**正好站在那个 commit** 才能复核；他若在别处，得到的数不同却会以为"门变了"—— 正是派单说的"错的数比没有数更坏"。**commit 号本身不朽（故保留 `938216b8e` 作失败证据），数不朽不了**。
4. **主流 harness 的做法**：注释/门禁写**不变式 + 复跑命令**，计数交给运行时输出与日志；把会漂的数印进 prose 是这仓反复付学费的形状（`362-r1`→`362-r2`→本笔，三笔只为改一句话）。

## §4 改后注释（逐字，现工作树 `sed -n '186,217p'`）

```
      # WHAT THAT DOES TO THE DENOMINATOR. "The denominator of the formatting gate
      # is the TRACKED set, not the working tree" is the claim the census step
      # below is built on, and it is no longer true OF THIS STEP as written: the
      # denominator here is now the PRODUCTION DIRECTORY SET named above, not
      # every tracked .go. The counts are deliberately NOT written here (ticket
      # 362-r3): a hand-copied count is honest about some tree and silent about
      # which one. Re-derive them, in the checkout you are actually reasoning
      # about, with one recipe, two lines:
      #
      #     git ls-files '*.go' | wc -l
      #     git ls-files '*.go' | cut -d/ -f1 | sort | uniq -c | sort -rn
      #
      # The first line is the denominator. The second is the same set split by
      # first path segment: the sum of internal, cmd, tools, scripts and
      # frontend is what this step is handed, and the `internal` line is the
      # number the previous version of this paragraph quoted on its own. That
      # split also re-proves, instead of repeating, the structural claim this
      # paragraph depends on - if a tracked .go file's first path segment is
      # anything else it prints as its own line (a file at the repo root prints
      # under its own name), so "the five directories are exactly production and
      # .scratch/ is exactly what leaves this denominator" is a reading you
      # take, not a sentence you trust. The reading is still taken over tracked
      # bytes on a checkout; what narrowed is which tracked bytes are counted.
      #
      # WHY NO NUMBERS: this paragraph has carried two hand-written count sets
      # (362-r1's, then 362-r2's), and the second does not reproduce on the
      # commit that wrote it - it matches, exactly, the tree of 938216b8e,
      # eleven commits behind that commit. Nothing about it was invented; it
      # was measured on a real tree. But a reader who tries to reproduce it at
      # the commit it claims gets different numbers and concludes the gate
      # moved. A count that rots is bad; a count that rots and then misleads is
      # worse than no count at all.
```

第二条改动（`:230`）：`(.scratch/**), so the two steps agree on the bytes today (802 each) by two` ⇒ **去掉 `(802 each)`**（`802` 是三个数之一，留着就是孤儿数）；其余全句未动。

## §5 YAML 解析 `rc=0`（**显式 `encoding='utf-8'`**）

```py
with open(p, encoding='utf-8') as f: doc = yaml.safe_load(f)      # 工作树
head = yaml.safe_load(subprocess.run(['git','show','HEAD:'+PATH],capture_output=True).stdout.decode('utf-8'))
```

```
YAML_PARSE_RC=0
HEAD jobs= 6 steps_total= 53 hits= 1      # job = lint, name = "gofmt (gofumpt)"
WORKTREE jobs= 6 steps_total= 53 hits= 1  # 同名同步、同 job、唯一命中
```

- **不写 `encoding='utf-8'` 会先死在 `UnicodeDecodeError`**（本机 python 默认 `cp936`，`ci.yml` 里有 UTF-8 中文）—— 那是**编码**的错、不是 YAML 的错；本件全部读数都显式写了。

## §6 "步名与 `run:` 行逐字未动"的 diff 读数

```
git diff --stat -- .github/workflows/ci.yml
#  1 file changed, 34 insertions(+), 16 deletions(-)
git diff -U0 -- .github/workflows/ci.yml | grep -E '^-' | grep -vE '^---' | grep -cE '(^| )\s*(run|name):'
#  0        （grep_rc=1）
git diff -U0 -- .github/workflows/ci.yml | grep -E '^-' | grep -vE '^---' | grep -nE 'run:|name:'
#  （无输出）
```

**更强的一条（语义而非文本）**：把 HEAD 与工作树各自 `yaml.safe_load`，取 `jobs.lint` 里 `name == "gofmt (gofumpt)"` 的那一步 ⇒ `RUN_EQUAL=True`、`IF_EQUAL=True`、`NAME_EQUAL=True`（`run` 逐字 `go install mvdan.cc/gofumpt@latest\nOUT="$(gofumpt -l internal cmd tools scripts frontend)"\n…`）；另 `grep -nE '(^|[^0-9])(1162|802|623|805|626|1165)([^0-9]|$)' ci.yml` ⇒ **0 命中**（`:619` 的 `360-minute` 是超时分钟数、不是分母，未动）。

## §7 门禁 `rc`（跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`；⛔ 走管道）

| 门 | 逐字命令 | rc | 读数 |
|---|---|---|---|
| d22scan | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations`（`ban #8` 实扫 `internal/`＋`cmd/`＋**`design/`＋`frontend/`** —— 派单写"只 `internal/`＋`cmd/`"**偏窄**；但"⛔ 扫 `.github/`"这一点不变，故本笔安全） |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `VERDICT GREEN - every over-budget tracked path is rostered by name…` |

**提交后复跑（在"我的树"上，⛔ 不靠共享工树）**：`git clone --shared` 到仓外 → `git checkout -q b833be20f` ⇒ d22scan **rc=0**、长度门 **rc=0／GREEN**（`tracked paths=17629`、`over-budget=61  covered by roster=61  not in roster=0`）。
★**共享工作树在 11:5x 那一刻是 rc=1（RED）**：两行 RED 都是别人**暂存未提交**的 `.scratch/wisp/issues/196-…-in-use-done.md`（`git status` 读 `R065`）⇒ **归因清楚、⛔ 我这两枚文件**；我提交前那次 GREEN 是改名前的树（`17627`）。**该改名随后由编排者 `1fdda4c47` 收编并补 roster ⇒ 12:0x 现工树复跑 `rc=0／GREEN`。**（这正是派单"⛔ 在主检出上量"的理由：同一棵工树，几分钟内能从 RED 变 GREEN，变的不是门。）

**尺的射程注**：d22scan 的 live scope 印 `internal/ 630`、`cmd/ 148` —— 那是**含未提交改动的工树**计数，**不是** §2 的 tracked `626/149`；两者不同对象、不同算法，⛔ 互相当核对。

## §8 两处要报的偏差（⛔ 阻塞，但不许含糊过去）

1. **派单/`362-v1` 的"旧 12 笔 / 基 `6e19f4a43`"对不上**（现跑）：`938216b8e..ade547ac8` ＝ **11**；`938216b8e..6e19f4a43` ＝ **9**；`ade547ac8^` ＝ **`8ebbc1491`**（**不是** `6e19f4a43`）。⇒ **那两枚数都该改**。但**结论不动**：三数**确系 `938216b8e`** —— 这次不是靠减法推的，是**直接量了 `938216b8e` 的树**（§2），尺从"数提交"换成"数树里的文件"。
2. **停手条件②命中，但我判它在本笔写面内**：`ci.yml` 里写这三个数的只有两处 —— `:191-193` 与 **`:212` 的 `(802 each)`**，**两处都在同一道 gofumpt 步的注释块里**（换行后 `:186-234`，步头 `:235`）。按写面"只那道步的注释"它在我面上；且 `802` 正是三个数之一，留着会成孤儿数 ⇒ **同笔一并去掉数（不改语义）**。**其它步／job 0 处**；另核 `.scratch/wisp/probes/161/r5/attrib.sh` 分母是**运行时算的**（`$TRACKED_GO`／`$TRACKED_GO_ALL`）、⛔ 写死数 ⇒ 仓里无第二处硬编码。★这条是我唯一越过字面三行的一步，**要收窄请裁、我再追加一笔**。
3. **停手条件①未命中**：那道门的口径**五个目录全印**、加 `.scratch` 与根 ⇒ 无"没印的第三个数"；`626+149+18+11+1=805`、`805+360=1165`。

## §9 落笔前五问（`ruler-scope-checklist.md`）

1. **这把尺量的是哪个对象？** ＝ `0754fa7d9` **提交里的 tracked 名册**（`git archive` 导出树 / `git ls-tree`），⛔ 主检出的工树。
2. **输出形状？** ＝ `find|sed` 正斜杠；`gofumpt -l` 在 Windows 上打**反斜杠**（故判归属用 `[\\/]`）。
3. **我数的是"代码"还是"注释"？** ＝ **文件枚数**（`.go` 文件名册），⛔ 行数、⛔ 注释。
4. **和总数对上了吗？** ＝ 是：`805+360=1165`、五目录内部加总 `626+149+18+11+1=805`，另用 `git ls-tree` 第二把尺同值。
5. **谁算的、什么算法、几位？** ＝ `git ls-files`/`ls-tree` 的**名册成员数**（十进制枚数），⛔ 哈希、⛔ 行数。

## §10 未动清单

零 `.go` 字节；⛔ 动 `run:`／`if:`／步名／其它步／`attrib.sh`／`.scratch/**`／台账／票名／`Status:`；⛔ `add -A`／`.`,⛔ amend/reset/rebase/stash/checkout ./clean；⛔ push。本笔只 `add` 两枚具名文件（`.github/workflows/ci.yml`、`docs/evidence/s1/362-r3-impl.md`）。
