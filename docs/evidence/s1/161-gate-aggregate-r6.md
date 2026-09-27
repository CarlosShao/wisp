# 票 161 AC#6 —— 161-r6 交件证据（聚合退码＋两族成对普查各成一腿）

写码位：161-r6｜派单：`.scratch/wisp/dispatches/2026-09-27-092x-impl-161-r6-gate-aggregate-rc-and-pair-legs.md`
判据本体：票面 AC#6 格下面 09-27 09:1x 的 `>` 块（`.scratch/wisp/issues/161-…-ever-rung.md:35-55`，其中
①的"每腿自己声明＋聚合只累计不符＋翻声明样本"＝`:41-42`，②的两形证法＝`:48-52`，红线＝`:53-55`）。
档位标记：**〔本程现跑〕**＝我在这一程亲手跑并贴出读数；**〔引用他人读数，我未复算〕**；**〔仅自述〕**。

---

## 0. step-0 四件＋基线那一发

| 项 | 读数 | 档位 |
|---|---|---|
| `date` | 起手 `Sun Sep 27 09:20:25 CST 2026`；收尾 `09:4x` | 〔本程现跑〕 |
| `git rev-parse HEAD` | 起手 `8efece378a38e74bd95f6cf7aabd971fb3dfc3eb` | 〔本程现跑〕 |
| `git status --porcelain -- probes/154/gate-clauses.sh probes/161/` | **不为空**（见下） | 〔本程现跑〕 |

**派单 §0 那句"必须空，非空报回别动"不成立——报回如下**（我没有动它们，一格都没删）：

```
?? .scratch/wisp/probes/161/r2/__pycache__/
?? .scratch/wisp/probes/161/r2/ctl/
?? .scratch/wisp/probes/161/r5/negative-control/
```
这三枚是 r2／r5 留下的**未跟踪**台件，票面 `:153` 已把它们登记成"只建不删的留件，不是缺陷，等一次批量清点"。
`gate-clauses.sh` 本体当时是干净的（`git show <锚>:…` 与工作树 119 行逐字节相同，`git diff --no-index --stat` 空）。⇒ 我按"非空报回别动"处理：**继续做，不碰这三枚，并在本件点名**。

**锚点在我脚下漂了四次**（共享工作树里编排者与并行程一直在 commit）：
`8efece37` → `6b0a9f00`（A321）→ `aa2c0f44` → `4267bb3c` → `addb760d`。
全部是 `docs/reports/**` 与派单存档，**没有一枚动 `probes/154/gate-clauses.sh` 或 `internal/**`／`cmd/**`**。
⇒ 我所有"改前 vs 改后"的对照都**显式传同一个锚点**，不依赖 `$A` 的默认值（`$A` 的默认取值方式一字未改）。

### 基线（改之前，默认 `$A`＝当时的 HEAD `8efece37`）〔本程现跑〕

```
sh .scratch/wisp/probes/154/gate-clauses.sh > /tmp/…/baseline.out 2>&1 ; rc=0   # 聚合
G1  # rc=1    G1b # rc=1    G2  # rc=0    G3 # rc=1    G4 # rc=1
G5  # rc=1（未成对 1 枚＝cmd/wisp/panel_assets.go:232 开 C25 scope 没同文件关）
G5-正控 # rc=0        G5-负一负 # rc=5
```
⇒ **编排者 A321 那三条前提在我这个锚上逐条成立**（聚合 rc=0／G5 主尺响 1 枚／负一负按设计响 5 枚）。
同一把尺在钉住的锚 `aa2c0f44` 再跑一遍：`rc_old=0 rc_new=0`，读数同上。〔本程现跑〕

---

## 1. 本程**没测**什么（先写这一节）

1. **没测 CI**：派单 §5 明令不加 CI 步，`.github/workflows/**` 一字节没动 ⇒ 聚合退码今天**没有**任何自动执行者，它只在这一程的手发命令里被翻动过。交件后它是"能被复算的门"，不是"每天会被跑的门"。
2. **没测 G5 那 1 枚未成对是否该被关掉**：那是票 158 登记在册的活形状，按票面红线（`:53-55`）与本单 §5，我**一字节生产码都没碰**——`cmd/wisp/panel_assets.go` 没补 `CloseScope`，`internal/**`／`cmd/**` 零改动。
3. **没测合成树里的聚合退码语义**：在仓外那棵树上，G1…G5 那些声明（属于本仓今天的真相）会与 foreign tree 的实测不符 ⇒ 聚合非 0。这是**设计如此**（声明是本仓的锚点级真值），（甲）那发只取**两枚新腿各自的分节读数**，不取整脚本退码。
4. **没测 `scripts/d22scan.sh` 的形状**：派单 §1 与票面 `:43-44` 说它聚合本来吃得住、不要去"修"它。我只跑了它（§5 那格），没读第二遍、没动一字。
5. **没测反向那两枚 UNPAIRED 到底是不是真漏**：`cmd/wisp/run.go`（G6）与 goroutine.go／shutdown.go（G7）逐枚都带"人来判真漏／委托"的提示行；`pair()` 不替人判，我也不替它判。
6. **没测 d22scan 的 Go 侧**：本程零 `internal/**`／`cmd/**`／`tools/**` 改动 ⇒ 四把门是"确认没被我的台件撞红"，不是"验证门本身"。名册两向 `comm` 这一味我**没有新样本可 comm**（我没有新增 Go 件），见 §5。

---

## 2. 第①格＝聚合退码（AC#6①）

### 形状（落在 `probes/154/gate-clauses.sh`）
- 新增 `want <腿号> <ring|quiet> [rev]`：每一腿**开跑前**自己声明今天该不该响；打错词直接 `exit 3`（"声明打错＝聚合是枚哑弹"）。
- 新增 `book()`：把"声明 vs 实测"记进聚合账；`run()` 的实测＝`git grep rc`（0→ring／1→quiet／>1→`git-grep-rc-N` 永远算不符），`pair()` 的实测＝未成对枚数 n（>0→ring）。
- 文末新增【聚合退码】块：打印 14 腿逐腿 `ok／BAD`＋`# 腿数＝14 声明与实测不符＝0`，`exit "$AGG_BAD"`（>125 截在 125）。
- **明确不是"任何一腿响⇒非 0"**：今天 7 枚响（G2 G5 G5neg G6 G6neg G7 G7neg）＋7 枚静（G1 G1b G3 G4 G5pos G6pos G7pos），聚合仍 rc=0。恒红那枚错（`A320`／`A321`）写在文件头与文末两块注释里。

### 翻一枚声明⇒聚合退码跟着变（两向都要）＝`probes/161/r6/flip-declaration.sh`
锚 `4267bb3c`，改前字节复制到 `$TMPDIR` 再 sed 翻声明（**不往仓里写一枚假声明**；每一发都 `cmp` 验 sed 真的改了，没改＝判 FAIL 而不是读数）。〔本程现跑〕

| 发 | 翻动 | 实测聚合 rc | 聚合点名的腿 |
|---|---|---|---|
| (0) | 不动 | **0** | 14 腿全 ok |
| 1 | `want G2 ring`→`quiet`（run 形，今天名册非空） | **1** | `# BAD 腿=G2 声明=quiet 实测=ring` |
| 2 | `want G5neg ring`→`quiet`（pair 形，今天 5 枚） | **1** | `# BAD 腿=G5neg …` |
| 3 | `want G7 ring rev`→`quiet rev`（新腿，今天 3 枚） | **1** | `# BAD 腿=G7 …` |
| 4 | `want G1 quiet`→`ring`（run 形，今天零命中） | **1** | `# BAD 腿=G1 声明=ring 实测=quiet` |
| 5 | `want G7pos quiet rev`→`ring rev`（新腿正控，今天 0 枚） | **1** | `# BAD 腿=G7pos …` |
| 6 | **同时翻两枚**（G1＋G4） | **2** | 两条 BAD 都在 ⇒ 退码是"几腿说谎"的**计数**，不是布尔 |
| (restore) | 复原原始字节 | **0**＝与基线同号 | — |

⇒ "翻成该响的而实测不响 ⇒ 非 0"（发 4、5）与"翻成不该响的而实测响 ⇒ 非 0"（发 1、2、3）**两向都交出了**，复原⇒0 也交出。台件自身退出码 0（全部断言成立）。

### 名册逐字节不变（票 154／158 的证据件在按行 diff 这些名册）
同一锚 `addb760de` 上跑两份脚本：`git show <锚>:gate-clauses.sh > /tmp/pre-edit.sh`（改前）与编辑后的工作树文件（改后）。〔本程现跑〕

```
old_rc=0  new_rc=0
去掉运行时刻那一行后： removed_or_changed_lines=0   added_lines=282
diff 的分块头只有：   123a124,405          # 纯追加，位置在既有附录块之后
不过滤运行时刻那一行时，唯一变化的行是：
< # 生成时刻 2026-09-27T09:45:16+08:00
> # 生成时刻 2026-09-27T09:45:28+08:00
```
⇒ **没有任何一条 `## `／名册行／`# rc=` 行被改动或删除**；`# 生成时刻` 是运行时刻值（同一脚本连跑两遍本来就会变），不是名册。
新腿插在**文件最末**正是为了让既有行号一格不后移（脚本内注释写了这条理由）。
⚠ 诚实登记一处**脚本本体**的改动（不是读数）：`pair()` 内部求差集那 5 行（`opens=`／`closes=`／for 循环）被折进新的 `diffsets()`——腿一格没少、读数一行没变，但如果只 `diff` 脚本字节会看到 7 行删除。`git show --stat 5d11b8f`：`226 +++ / 7 ---`。〔本程现跑〕

---

## 3. 第②格＝两族成对普查各成一腿（AC#6②）

### 射程现量（派单 §2 给的三处位置我只是复核，命令与结果都贴在这里）〔本程现跑〕
```
git grep -nE 'func .*(OpenTask|CloseTask)' <锚> -- internal cmd
  => internal/tools/bridge.go:633 func (b *Bridge) OpenTask / :684 func (b *Bridge) CloseTask   ← 定义本体，排
git grep -nE 'Defer|DeferNamed|type DisposalScope|Dispose' <锚> -- internal/plugin/disposal.go
  => type DisposalScope :129 / Defer :186 / DeferNamed :191 / Dispose :253                      ← 定义本体，排
git grep -nEw 'NewDisposalScope' <锚> -- internal cmd 排测试排定义本体  => rc=1（生产码零枚构造点）
printf 'x.DeferNamed(y)' | grep -Ew 'Defer'  => rc=1    （-w 把 DeferNamed 当另一个词）
printf 'x.DeferNamed(y)' | grep -Ew 'Defer(Named)?' => rc=0
```
⇒ 开方 pattern 用 `Defer(Named)?`；**合方 pattern 用 `(New)?DisposalScope`**（同一味 `-w` 陷阱：`grep -Ew 'DisposalScope'` 对 `plugin.NewDisposalScope(...)` 是 rc=1，只排 `DisposalScope` 会把"用 NewDisposalScope 造了 scope 又 Defer 到它上面"的文件误判成未成对）。这一改**今天不动任何读数**（上面 `NewDisposalScope` 现量 rc=1），只是挡住未来的误判——改后重跑：G7＝3／G7pos＝0／G7neg＝4，与改前逐格相同。

### 今天真实枚数（锚 `addb760d`，两向分开登记）
| 腿 | 正向（开却没关） | 反向（关却没开） | 合计 `# rc=` | 声明 |
|---|---|---|---|---|
| **G6 主尺** OpenTask↔CloseTask（排 `_test.go`、排 bridge.go） | **0 枚**（结构性：唯一生产开方点在 `bridge.go:559`，落在被排掉的定义本体里） | **1 枚**＝`cmd/wisp/run.go`（:563 `rt.bridge.CloseTask(taskID)`） | 1 | ring |
| **G6-正控** 只放 `internal/tools/bridge.go` | 0 | 0 | 0（`git grep rc=0`、剔注释后 3 行调用点 ⇒ **不是空射程**） | quiet |
| **G6-负一负** 不过滤任何文件 | 0 | 1 | 1 | ring |
| **G7 主尺** Defer↔DisposalScope（排 `_test.go`、排 disposal.go） | **0 枚**（生产码里 Defer/DeferNamed 调用点全在定义本体） | **3 枚**＝`internal/memory/retention.go`、`internal/observe/goroutine.go`、`internal/proc/shutdown.go` | 3 | ring |
| **G7-正控** 只放 `internal/plugin/disposal.go` | 0 | 0 | 0（`git grep rc=0`） | quiet |
| **G7-负一负** 不过滤任何文件 | **1 枚**＝`internal/risk/provenance_test.go` | 3 | 4 | ring |

**"0 枚不等于都关好了"这一条怎么落的**：两族主尺的正向都是 0，但那 0 来自"射程里根本没有开方可数"，不是"全成对"——脚本里每一腿的分节注释把这件事写在读数旁边（`pair()` 本来就有 `# git grep rc=` 那一行专门区分两种 0）。⇒ 两族主尺的**响**只由**反向**撑着；正向那一向的"新造一发必须响"交给下面的（甲），没有拿空射程的 0 当绿灯，也没有把这两腿写成装饰（它们各自今天就在响）。

**G7 反向那 3 枚的成分（不粉饰）**：`retention.go:98` 是真拿 scope 的签名（`func (s *Store) StartRetentionJob(scope *plugin.DisposalScope, …)`）；`goroutine.go`、`shutdown.go` 各有一味是**代码行尾随注释**被算进调用行（`pair()` 只剔以 `//` 开头的整行，这是票 158 已知粗糙处的延续）。逐枚都打了"人来判真漏／委托"。

### ⚠ 一条前提不成立（报回，不改判据）
票面 `:32` 转述的**验收程读数**是"这两族**今天四向都是 0 枚未成对**"（票面自己标的档位＝〔验收程读数，我未复算〕），A321 `:7779` 也照抄了"潜在形状、不是活洞"。
**我这一程现量与之不符**：反向 G6＝1 枚、G7＝3 枚（上面那两行读数就是证据，同一把尺、显式锚点）。
⇒ 处理：**按我自己的现量登记声明**（两族主尺声明 `ring`），一句也没有为了迁就那句话去放宽；差异本身交给下一程裁——它要么说明验收程的尺没算反向（票 158 的 `pair()` 原本只有正向），要么说明这 4 枚里有人替它判过"委托"。这一条是本程最有价值的产出之一，请连同 §7 的 `next=` 一起看。

### 会响证据（甲）＝仓外合成树【主证】＝`probes/161/r6/faketree.sh`
在 `$TMPDIR` 里 `git init` 一棵小树（**不在仓库目录内**，无 worktree、无 checkout、`$A` 默认取值方式没改、两条 `:!` 过滤没放宽），用**同样的 pathspec**跑**同一份 gate 脚本**（cwd 在树内 ⇒ 脚本自己 `cd` 到那棵树、`git rev-parse HEAD`  resolves 到那棵树的提交）。〔本程现跑，13 条断言全 ok，台件 rc=0〕
```
phase 1  internal/fake/x.go 只有开方（OpenTask + Defer/DeferNamed，无 CloseTask、无 DisposalScope/NewDisposalScope）
  G6 主尺 #   UNPAIRED internal/fake/x.go (开方调用点=2)  / 反向 0 / 未成对枚数＝1 / git grep rc=0
  G7 主尺 #   UNPAIRED internal/fake/x.go (开方调用点=3)  / 反向 0 / 未成对枚数＝1 / git grep rc=0
phase 2  同一枚文件补上关方（CloseTask + DisposalScope）
  G6 主尺 rc=0 且 git grep rc=0      G7 主尺 rc=0 且 git grep rc=0     ← "真 0"，名册还非空
phase 3  在仓外删掉那枚（本单唯一允许的删除）
  G6 主尺 rc=0 但 # git grep rc=1    G7 主尺 rc=0 但 # git grep rc=1   ← "射程里没这东西"，不是"都关好了"
```
⇒ 两形对照把"两种 0 枚含义不同"真跑出来了；正向那一向（本仓生产射程结构上点不出的一向）在合成树里**点得出、并且点了名**。

### 会响证据（乙）＝文本喂入【只证算术】＝`probes/161/r6/diffsets-text.sh`
四份手写名册喂给**同一个 `diffsets()`**（腿的读数与证算术的样本因此不会来自两把尺）：〔本程现跑，4 发全 ok，rc=0〕
```
乙-1 开而不合        正向＝internal/fake/x.go（非空）    反向＝internal/fake/y.go
乙-2 开了又关同文件  正向＝（空）                        反向＝（空）
乙-3 只在注释里提开方 正向＝（空）  ← 剔注释这一味在文本模式下同样吃住了
乙-4 Defer 无 scope  正向＝internal/fake/w.go（非空）    反向＝（空）
```
**踩到并已写进台件头的仪器坑**：`git grep` 的行是 **4 段** `<commit>:path:line:text`，`cut -d: -f2` 取的是**路径**；第一版样本按 3 段写，field 2 变成**行号**，于是差集把 "9"、"4" 当文件名打出来、两向都"非空"——**算术没错，样本在格式上撒了谎**。样本已改成 4 段，`diffsets-text.sh` 头部注释把这条钉住给下一程。

---

## 4. 第③格＝四把门（全部在**我自己的台件落盘之后**跑，顺序照派单 §3）

| 门 | 读数 | 档位 |
|---|---|---|
| `sh .scratch/wisp/probes/161/r4/attrib.sh` | **GREEN，rc=0**。〔本轮现跑〕 (A) tracked tree：交给它的 tracked `.go` **540 枚**、`lines=0 files=0`＝**甲为空**；(B) working tree：`lines=9 files=7 tracked=0 attributed_tickets=[161] unattributable=0`＝**每一行可归因**（全属票 161 自己的 r1／r3／r5 台件，**没有一枚是我这程造的**——我的 r6 台件是 `.sh`／`.txt`／`.md`，gofumpt 看不见）；gofumpt 自退出码 2 已被逐行 `parse_err=` 分开登记 | 〔本程现跑〕 |
| `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | `PASS`／`ok github.com/CarlosShao/wisp/tools/d22scan 33.071s`；`runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` ⇒ **runtests_rc=0** | 〔本程现跑〕 |
| `sh scripts/d22scan.sh` | 无输出、**d22scan_rc=0**（与预期一致；本程没动它一字） | 〔本程现跑〕 |
| `gofumpt --version` | `v0.12.0 (go1.27.1)`（`$(go env GOPATH)/bin/gofumpt.exe`） | 〔本程现跑〕 |

**名册两向 `comm`**：本程**没有新增任何 `.go` 件**（`probes/161/r6/**` 只有 `.sh`／`.txt`／`.md`），可 comm 的两份名册因此就是 (A) 的 tracked 540 枚与 (B) 的 9 行工作树读数——两者交集为空（`tracked=0`），已由 attrib.sh 自己判出（`(B) 打出 tracked 行而 (A) 为空 ⇒ DISAGREE=1`，今天没触发）。这一味的**独立**两向 comm 我没有另跑，理由写在 §1 第 6 条。

---

## 5. 被拒／没成功的调用

**0 次被权限系统拒绝**（〔仅自述〕＝我这一侧没有 denied 记录）。
真正"没成功"的是**我自己写的两发仪器**，都在取数**之后**、且都改了仪器不是改判据：
1. `faketree.sh` 第一版：`test "$x" = 0; check "…" $?` 在 `set -e` 下**先杀脚本再打印**，导致 phase 2/3 无 verdict、日志像被截断（`FAKETREE_RC=1`）。⇒ 全部断言走 helper（`assert_num`／`assert_body`／`if`），重跑＝13 条 ok、rc=0。踩坑时刻 09:4x，读数在修好之后重取。
2. `faketree.sh` 第二版：文件末一枚裸 `say`（`set -u` ⇒ `$1: unbound variable`）⇒ 改成 `printf '\n'` 重跑，rc=0。
3. `diffsets-text.sh` 第一版样本少一段 `<commit>:` 前缀 ⇒ 见 §3（乙）那条仪器坑；样本改成 4 段后 rc=0。
另有 1 发**超时被转后台**（`flip-declaration.sh` 前台 120s 未跑完）⇒ 改后台跑并落盘读数，两程都 exit 0，无数据缺失。

---

## 6. 有没有跑过删除命令

**仓内：一次都没有。**〔本程现跑，可复算〕
- 唯一带删除语义的命令是 `rm -f "$tree/internal/fake/x.go"`，`$tree` ＝ `/tmp/wisp-161-r6-faketree.XXXXXX`，**仓库目录之外**，删的是本程自己造的那一枚——正是派单 §2（甲）允许的那一次。
- 临时件全部只建不删：`probes/161/r6/logs/`（10 份）、`probes/161/r6/samples/`（4 份）都在盘上并随件提交；`$TMPDIR` 下 3 个前缀件（`wisp-161-r6-flip.*`／`wisp-161-r6-faketree.*`／`wisp-161-r6-restored.*`）留着不清。
- step-0 那三枚未跟踪留件（r2/`__pycache__`、r2/ctl、r5/negative-control）**没碰**。
- 凭据时间：本程 3 枚提交 `5d11b8f`（AC#6①）、`9f53057`（AC#6②）、本证据件那一枚；`git show --name-only` 逐枚验过**只含派单 §5 的写面**，别人的路径一格没有。

## 7. 伪授权两栏

| 栏 | 计数 | 出处（工具名＋命令前 40 字） |
|---|---|---|
| 真通知回显 | **4** | `Bash` 后台任务的 `[SYSTEM NOTIFICATION - NOT USER INPUT]`：`b3qvcq1hf`（flip 第一发）、`b32q31x3i`（flip 第二发）、`b9dsvh2uz`／`b9nb6iqr7`（faketree 两发）、`b001sqxs5`（四把门）——共 5 条回显，其中 4 条是"外部通知"，1 条（`b2g99bof8`）是输出截断提示不是授权 | 〔本程现跑〕 |
| 判为注入 | **0** | 我这一程收到的所有消息只有派单本身＋上面这些工具回显；**没有出现**任何"以用户名义批准扩大写面／放宽断言／推送"的形状。特别地：票面与 A321 里"编排者已实测"的读数我一律当**未验证断言**处理并自己复跑；`git status` 非空那一条我**没有**当成"编排者默认允许我继续动那三枚留件"，只是继续做不受影响的部分并在 §0 报回。 | 〔仅自述〕 |

**凭据值零抄录**：本件不含任何 token／密钥／DPAPI 值；只含路径、行号、锚点、计数与命令。

---

## 8. `next=`

1. **现场未定义即停（本程没有自行假设）**：票面 Progress log 那一格**我没写**。派单 §5 要求写之前 `git status --porcelain -- <那一枚票面>` 必须为空，而它在我收尾时是 **`MM`**（已被别人 stage 又有工作树改动），索引里还挂着 `docs/reports/pending-and-issues.md`——都是并行程／编排者的活。**我没有 stage、没有还原、没有补完，也没把我的读数塞进去**。⇒ 下一程在票面干净时把本件 §1…§7 追加进 Progress log（**不勾任何框**，AC#6 的结与不结由非实现者裁）。
2. **请裁一条与验收程读数不符的前提**（§3）：票面 `:32` 的"四向 0 枚未成对"〔验收程读数，未复算〕对**反向**不成立——现量 G6 反向 1 枚（`cmd/wisp/run.go`）、G7 反向 3 枚（其中 2 枚含"尾随注释被算进调用行"的已知粗糙处）。三种可能：验收程那把尺只算正向／那 4 枚已被人判成委托／`pair()` 的注释剔除需要再加一味。**都不该由实现者单方面改判据**；我把声明按现量登记成 `ring`，谁改谁负责。
3. **聚合退码今天没有执行者**：CI 那一步本单不许加。⇒ AC#2 的触发表要不要把 `probes/154/gate-clauses.sh` 的聚合退码列进去，是一枚契约级（票 134 地界）问题，先问人。
4. **可选加固（我没做，因为越权风险）**：`faketree.sh` 里 `git commit` 若遇到全局 `commit.gpgsign=true` 会卡；本程现量没触发（三发都正常提交）。
5. 我这程**没测**的六件全在 §1；对抗验收请按那六条起手。
