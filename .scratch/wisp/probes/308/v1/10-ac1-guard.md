# 308-v1 · 10-ac1-guard（`AC#1`：守卫落点＋三条硬约束＋"⛔ skip 换绿"＋分母＋覆盖率⛔ 悄悄搬家）

判据原文（票 `AC#1`）＝无卷标那一形下这枚用例**以可解释的颜色结束**（⛔ panic、⛔ 整枚二进制没了）；守卫落点⛔ ＝`t.Skip` 整枚；只许可**具名去掉 `unc-admin` 那一枚拼法**且**分母要打印**；硬约束＝⛔ 改判级逻辑／⛔ 动 `pathresolver.go`／`pathresolver_windows.go`／⛔ 把 `t.TempDir()` 换成手写常量。

## 1. ★产码⛔ 动（硬约束，本腿三把尺各自独立）

1. `git show --numstat --format=` 逐笔（五笔全量，件＝本腿 `logs/v05-numstat.txt` 口径＝直接抄工具输出）：只有 `1ad48333` 动 `internal/`，且**三枚全是 `_test.go`**：
   `18 17 internal/risk/pathresolver_anchor_spelling_windows_test.go`／`12 1 …junction_windows_test.go`／`120 0 …volumeless_spelling_308_test.go`。其余四笔只动 `.scratch/wisp/probes/308/r1/**`。
   ⇒ 尺＝`git show --name-only --format= $c | grep -vE '_test\.go$|^\.scratch/wisp/probes/308/r1/'` 五笔**各 0 行**。
2. **blob hash 全列**（尺＝`git rev-parse <ref>:<path>`，射程＝`510c95ca^..HEAD` 七枚 ref）：
   - `internal/risk/pathresolver.go` ＝ **`15349e28baf70da3bd7aa1d101ebf61cee20fe35`** ＠ `510c95ca^`／`510c95ca`／`1ad48333`／`ee0a93e4`／`79242df7`／`1270e36e`／`HEAD`（七枚同 hash）
   - `internal/risk/pathresolver_windows.go` ＝ **`7b00b022cf437638794e82eaac79e5fa92845a4a`**（同七枚同 hash）
   - `git diff --stat 510c95ca^ -- <这两枚>`＝空，`git status --porcelain -- <这两枚>`＝空 ⇒ **工作树⇄HEAD⇄落地前三者同一**。
3. **守卫行逐字复跑**（尺＝`git show HEAD:internal/risk/pathresolver.go` 自 250 行数，命中行＝**262**，切片行＝**263**）：
   `262: if len(drive) == 2 && drive[1] == '$' && drive[0] >= 'a' && drive[0] <= 'z' || len(drive) == 2 && drive[1] == '$' && drive[0] >= 'A' && drive[0] <= 'Z' {`
   `263:     return strings.ToUpper(drive[:1]) + ` + 反引号 `:\` + rest`
   ⇒ 长度守卫**仍在 `[:1]` 之前**；票面 `:262-263` 那两行号今天⛔ 漂。
4. ★**全仓 `[:1]` 驱动器位普查**（尺＝`grep -rn --include=*.go 'vol\[:1\]\|drive\[:1\]' internal cmd tools`，⛔ 注释行另筛）＝**只剩两枚求值位**：产码 `pathresolver.go:263`（有守卫）＋新 helper `…volumeless_spelling_308_test.go:27`（`len(vol)==0 ⇒ ¬ok`）。`internal/tools/volume_windows.go:28/:54` 读 `VolumeName` 但走 `len(vol)!=2` 一支，⛔ 同款切片 ⇒ 本票射程外且⛔ 崩形。
   ⇒ **判级逻辑⛔ 动**（Classify/Gate/normalize 全在 hash 未变的产码里）。

## 2. 无 `t.Skip` 换绿（`AC#1`/`AC#2` 边界）

- 尺＝`grep -cE '^[[:space:]]+t\.Skip\(' <pre blob> ⇄ <工作树>`：anchor **0⇄0**、junction **1⇄1**（那一枚是 `TestPathResolverShortNameAListDenied` 里 8.3 关闭时的既有 skip，票 18 具名许可、且已在 `scripts/portable-tests.sh` 的 ledger 里挂号，⛔ 本票引入）。
- 五笔里 `+` 行含 `t.Skip` 字样的共 **3** 行，逐枚读过＝**全是注释**（anchor:110 "does NOT t.Skip"／junction:156 "do NOT t.Skip"／新件:37 "⛔ t.Skip"）。⇒ ⛔ 新增 skip、⛔ 造 `--- SKIP`。
- 实测色：工作树整包 `--- SKIP`＝**1**（`TestSyncRegistryProbeLive`，既有、ledger 在册）⇄ 改前导出树也是 **1** ⇒ skip 枚数⛔ 因本票而变。

## 3. ★只按名去掉 `unc-*`，且分母打印出来（裁语 ⓒ）

- 逐字（新件 `:51`）＝`return keep, []string{"unc-admin", "unc-extended"}` ⇒ 去的是**两枚**具名拼法（⛔ 整枚用例；票写"只许具名去掉 `unc-admin` 那一枚"，helper 把 `unc-extended` 一起去了——见 §5 与本格判语里的具名让步）。
- 分母打印＝尺＝跑新用例的 `t.Logf` 原句（件 `logs/v-target-case.txt`，工作树 Windows 台面）：
  `ticket 308 guard on volume-less input "\\tmp\\wisp-308-no-volume\\.ssh\\id_testkey": 7 volume-free spellings still evaluated [dot-and-dotdot extended long lower not-yet-there slashes upper]; 2 UNC spellings dropped by name [unc-admin unc-extended]`
  ⇒ **7 留／2 去，名册逐项印出**，"denominator is zero reports clean for a tree it never read" 那条标准正面满足（分母⛔ 是 0，⛔ 是没印）。
- linux 台面同一句＝`... volume-less input "/tmp/wisp-308-no-volume/.ssh/id_testkey": 7 ... `（件 `logs/v26-docker-linux-full.txt`）⇒ 同一分母跨 OS 一致。

## 4. ★覆盖率⛔ 悄悄搬家：改前⇄改后两把名册并排

| 台面／尺 | 改前（`1ad48333^` 导出树） | 改后（`HEAD` 导出树） | 作差 |
|---|---|---|---|
| windows default，顶层名册 `^--- (PASS\|FAIL\|SKIP): ` | rc=0 **132** 行（131 PASS＋1 既有 SKIP），RUN=200 | rc=0 **133** 行（132 PASS＋1 SKIP），RUN=201 | **＋1**＝`PASS TestVolumeInvariantSpellingsVolumeLess308`；⛔ 少名、⛔ 变色（件 `logs/v12`/`v13`/`v42`… ） |
| windows alias（8.3 那把） | rc=0 **132** | rc=0 **133** | ＋1，同一枚 |
| 子层级 `^    --- PASS` | **68** | **68** | 0 |
| **linux（新件⛔ tag⇒新纳入的这一档）** | rc=0 **108** 行，RUN=168 | rc=0 **109** 行，RUN=169 | ＋1，同一枚 |
| 整包求值枚数（工作树，守卫在位） | — | RUN=201 PASS=132 FAIL=0 SKIP=1 PANIC=0 | — |

**拼法级分母**（⛔ 止于"数一下 map 的 key"，两把尺并写、各标量法）：
- 改前＝**算的**：`git show 1ad48333^:…anchor…` 第 106–116 行 map 字面量 key 枚数＝**9**（`long upper lower extended unc-admin unc-extended slashes dot-and-dotdot not-yet-there`）＋`8.3-short` 条件位。
- 改后＝**量的**：仓外导出树里加一枚探针用例（⛔ 仓内⛔ 动交付件；件 `logs/v45-spelling-denominator-probe.txt`）：
  `V1PROBE volume-qualified: keep=9 [dot-and-dotdot extended long lower not-yet-there slashes unc-admin unc-extended upper] dropped=[]` ⇄ `V1PROBE volume-less-literal: keep=7 [...] dropped=[unc-admin unc-extended]`
  ⇒ 可达台面（卷标齐，＝本机与 CI windows 档唯一形态）上 **9 ⇄ 9、名字全同**；只有喂无卷标字面量才 9→7。⛔ 静默搬家。
- ⚠ 尺的极限（具名）：`TestClassifySpellingInvariance` 的拼法循环⛔ 是 `t.Run` 子用例 ⇒ 顶层名册看⛔ 见拼法枚数，"名册没变"这件事本身⛔ 能证拼法没少 ⇒ 所以上面那把探针尺是必需的，⛔ 拿名册差＝＋1 顶掉它。

## 5. 本格判语

**成立**（射程＝工作树⇄HEAD blob⇄仓外导出树；单位见 00-anchor）。三条硬约束逐条有尺、⛔ 一条破口。
两处⛔ 致命但具名的偏移：
- ⓐ 票面只许可去掉 `unc-admin` **那一枚**，交付去的是 `unc-admin`＋`unc-extended` **两枚**。判＝⛔ 越界（`unc-extended` 与 `unc-admin` 同样由 `vol` 构造，`\\?\UNC\localhost\` + share + rest 在无卷标时同样造不出；留它＝留一枚必然红的假拼法），且分母把两枚都印了名 ⇒ 属"具名去掉无法构造的拼法"的同一件事。⛔ 当⛔ 未发生的偏移，记此。
- ⓑ `TestClassifySpellingInvariance` 里 `if len(dropped) > 0 { t.Logf(...) }` 那一支，在**任何可构造台面上都⛔ 到达**（实测哨兵 `dropped the two volume-derived` 命中：工作树 0／导出树 0／别名 0／linux 0，因为那两枚件在 linux ⛔ 编译）。⇒ "无卷标那一形下这枚用例以可解释的颜色结束"这半句在 windows-tagged 用例上是**形状保证**（走同一枚 helper、helper 那一支已被输入面用例钉住），⛔ 是"那枚用例真在⛔ 卷标台面上跑过一回并绿了"。这一条⛔ 抹掉，登记为 `Q-308-v1-B`，判语口径见 20-ac2-teeth.md §4。
