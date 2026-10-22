# 票 360 补强腿 `360-r2` 交付件（`AC#2` 补强：把"该 Adopt 臂真的跑过"独立成一条断）

- 时刻＝2026-10-22 19:2x–19:4x +08｜起手锚＝`dac97d702`（派单存档）｜本腿只 commit、⛔ push。
- 形仍是〔甲〕**只改一枚测试件**、**零产码**（五枚 `rev-parse` 全 SAME，见 §6）。
- 写面＝`internal/config/writeguard_adopt_360_test.go`（**只改这一枚**）＋本件。临时件全在仓外 `D:\tmp\wisp\20261022\360r2\**`。
- ★本腿**只做补强、⛔ 重做**：上一枚 `360-r1` 的四笔 `d68ddfae`／`fc7c85cc`／`0d907103`／`7c42145c` 未被本腿触碰或改写。

## 1. 开工前（派单「⛔ 跳步」三条，逐条现跑）

1. **`git status --porcelain` 写面**：`-- internal/config docs/evidence/s1` **零行** ⇒ 写面内无别人的未提交改动，可开工。
   （仓内别处有 `M docs/reports/pending-and-issues.md` 等既有脏件，**均在写面外**、本腿不碰。）
2. **票面前提逐条现跑**（`git grep` 对 HEAD，⛔ 抄 `360-a1`／`360-v1`）：
   - `statOwnWrite` **5 行 ＝ 1 定义 ＋ 1 代码调用 ＋ 3 注释**：定义 `allowdirs.go:157`／调用 `writeguard.go:179`／注释 `allowdirs.go:150`、`writeguard.go:105`、`writeguard_226_test.go:186` ⇒ **与票面逐枚同**。
   - 写回入口 **3 枚** → `mergeWrite`：`allowdirs.go:142`／`permmode.go:79`／`settings.go:303`；定义 `writeguard.go:111`。
   - Adopt 全仓一条：`manager.go:142 CheckAndReload`。
   - **两条路互不相连**：`git grep 'm\.cur = '` 全仓**仅** `settings.go:304`（失败回滚），⛔ 任何 adopt 写。
   ⇒ **前提今天全部成立**，不触发派单停手上报条②。
3. **复跑上一枚的两态读数**（⛔ 抄它的数字；件 `logs/01-pristine-before-edit.txt`）：
   ```
   writeguard_adopt_360_test.go:54: STATE 1 after writeback: file.ball.size=60 mem.ball.size=56
   writeguard_adopt_360_test.go:82: STATE 2 after one Adopt pass: mem.ball.size=60 file.ball.size=60 report.Hot=[ball]
   --- PASS: TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit (0.05s)
   ok  github.com/CarlosShao/wisp/internal/config  0.207s   RC=0
   ```

## 2. ★补了哪一行前置断（逐行）

**补法＝`ranAdopt` 哨兵 ＋ 一条独立前置断，放在 `:98` 那条臂断言之前。**

★**行号口径（必读）**：本件凡写"补强后"的行号，一律指**补强后的 107 行件**——`:72` 哨兵声明／`:75-84` Adopt 臂／`:85` 状态 2 打印／**`:93` 前置断的 `if`、`:94` 前置断的 `Fatalf`**／`:98` 臂的值断／`:103` 两态守卫。
★**变异发里报出的红行 `:84`** 是**摘臂后**（`sed -i '75,84d'` 删掉 10 行）的行号：前置断的 `Fatalf` 由 `:94` 上移到 `:84`。⇒ **同一行代码，两个号**：原件 `:94`、摘臂台 `:84`。★上一枚 93 行件的 `:84`／`:89` 是**另一套旧号**，本件引用旧号处一律注明。

改动 A（`ranAdopt` 声明，插在 `adoptedMem, adoptedFile := -1, -1` 之前）：

```
ranAdopt := false
```

改动 B（在 Adopt 臂的 `else` 体内置位，`CheckAndReload` 成功返回即证明臂体执行过）：

```
} else {
	ranAdopt = true
	if rep != nil {
```

改动 C（**本格的核心：前置断**，位置在 `t.Logf("STATE 2 …")` 之后、臂断言之前）：

```
// THE PRE-CONDITION (ticket 360-r2). Every assertion below reads the control
// arm's output, so first prove the arm ran at all. This is a separate statement
// with its own Fatalf and it reads nothing but the sentinels above, so it sits
// before the value assertions on purpose: neutering those assertions, or letting
// their two-state guard fall silent, cannot silence this one. Remove the arm and
// ranAdopt stays false while both sentinels stay at -1, so this fires first.
if !ranAdopt || adoptedMem < 0 || adoptedFile < 0 {
	t.Fatalf("the Adopt arm did not run: ranAdopt=%v mem.ball.size=%d file.ball.size=%d; "+
		"the positive control never executed, so the state 1 readings above are unbacked",
		ranAdopt, adoptedMem, adoptedFile)
}
```

- **为什么它"⛔ 能被同臂断言吞掉"**：它是**独立语句、独立 `Fatalf`**，只读 `ranAdopt`／`adoptedMem`／`adoptedFile` 三个哨兵，**不读** `:98` 那三条值断的任何中间量，也**不依赖** `:103` 的两态守卫。
  摘臂后 `ranAdopt=false` 且两哨兵停在 `-1` ⇒ **本条先于一切值断开火**。
  把 `:98` 改恒真、把 `:103` 守卫弄哑，**都改变不了 `ranAdopt` 的取值** ⇒ 本条仍红。
- **零产码**：只动测试件，`internal/config` 的非测试 `.go` 一个字节未改（§6）。

## 3. ★★`M1＋M3b`（摘臂 ＋ 该臂断言恒真）—— 本格唯一判据，我自跑

**台**：仓外导出树 `D:\tmp\wisp\20261022\360r2\tree`（`git archive HEAD` ⇒ **17575/17575** 件 ＝ `git ls-files` 17575），再把本腿改后的测试件覆盖进去。
**造法**：`sed -i '98s/.*/\tif false {/'`（把臂断言 `:98` 改恒真）＋ `sed -i '75,84d'`（摘掉整条 `if rep, err := m.CheckAndReload() … }` 臂）。
**落地证明**：`grep -n CheckAndReload` ⇒ 只剩注释 `:19`、**代码调用 0**；`grep -n "if false"` ⇒ **1 命中 `:88`**；`go build ./...` **rc=0**、`go vet` **rc=0**。

逐字读数（件 `logs/04-M1M3b-red.txt`）：

```
    writeguard_adopt_360_test.go:54: STATE 1 after writeback: file.ball.size=60 mem.ball.size=56
    writeguard_adopt_360_test.go:75: STATE 2 after one Adopt pass: mem.ball.size=-1 file.ball.size=-1 report.Hot=[]
    writeguard_adopt_360_test.go:84: the Adopt arm did not run: ranAdopt=false mem.ball.size=-1 file.ball.size=-1; the positive control never executed, so the state 1 readings above are unbacked
--- FAIL: TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit (0.04s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/config	0.170s
RC=1
```

⇒ ★**上一枚在这一发上 `RC=0` 仍绿；本腿在这一发上 `RC=1`、红行逐名 `writeguard_adopt_360_test.go:84`** ⇒ **判据 `改恒真` 那一半补上了**。

### 3.1 ★更狠的一发（我自己加的，⛔ 派单要求）：**摘臂 ＋ 把所有别的断言都弄哑**

若只证明"摘臂＋臂断言恒真"会红，还可以被质疑"是不是靠 `:89` 那守卫或状态 1 三条撑着"。所以我又造一发：
摘臂 ＋ `:55`／`:59`／`:64` 三条状态 1 断与 `:98` 臂断言**全部**改成 `if false && …`（编译得过的恒真形）。

```
    writeguard_adopt_360_test.go:84: the Adopt arm did not run: ranAdopt=false mem.ball.size=-1 file.ball.size=-1; the positive control never executed, so the state 1 readings above are unbacked
--- FAIL: TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit (0.08s)
FAIL	github.com/CarlosShao/wisp/internal/config	0.185s
RC=1
```

⇒ ★**整条用例只剩这一条前置断还在咬，它仍然红** ⇒ **"不可被同臂断言吞掉"是实测成立、⛔ 声称**（件 `logs/08-M1M3aM3b.txt`）。

## 4. 其余变异发（逐发独立、每发先还原到 sha 核过的原件）

| 发 | 造法（落地证明） | build | 读数 |
|---|---|---|---|
| **M1 单独摘臂** | `sed -i '75,84d'`；`CheckAndReload` 代码调用 **0** | 0 | `STATE 2 … -1/-1 []`／**红 `:84`** `the Adopt arm did not run: ranAdopt=false …`／`RC=1`（`logs/05-M1-alone.txt`） |
| **M3b 单独恒真**（只把臂断言 `:98` 改 `if false`） | 1 命中 | 0 | `STATE 1 60/56`、`STATE 2 60/60 [ball]`、**`--- PASS` `RC=0`**（`logs/06-M3b-alone.txt`） |
| **M2 产码侧 adopt**（`writeguard.go:169` 后插 `m.cur = deepCopyConfig(base)`；`grep -n MUTATION-360r2-M2` ⇒ `:170`） | 命中 `:170` | 0 | `STATE 1 … 60/60`／**红 `:60`**／**红 `:65`**（`divergence = [], want exactly [ball.size]`）／**红 `:99`**（臂值断，报 `report.Hot=[]`）／`RC=1`（`logs/07-M2-adopt.txt`） |
| ★**M1＋M3b**（唯一判据） | 见 §3 | 0 | **红 `:84`**／`RC=1` |
| ★**M1＋M3a＋M3b**（更狠） | 见 §3.1 | 0 | **红 `:84`**／`RC=1` |

★**M3b 单独恒真仍绿：照实记账，⛔ 藏、⛔ 算成红。** 机理＝单独把 `:98` 弄哑时臂**真的跑过**（`ranAdopt=true`、两哨兵 `60/60`），前置断因此正确地**不咬**；此时 `:103` 的两态守卫（`adoptedMem == memSize` ⇒ `60 == 56` 不成立）也仍不咬。
★但它**不构成判据缺口**：判据问的是"**摘臂**或**改恒真** ⇒ 必须转红"，而 M3b 与摘臂**同时**发生的那一发（§3）**已经红**；单独恒真这一发测的是"臂在、断言哑"这个本就不该红的组合（臂真的跑了，前置断没有理由开火）。
★**这正是补强的边界、我写清而不抹**：前置断管的是"臂跑没跑"，⛔ 管"值对不对"；值对不对仍由 `:98`＋`:103` 承担，且 M2 已证那两条**独立成牙**。

## 5. 三读数（逐字）＋ 还原核 `sha256`

| 读数 | 命令/台 | 逐字末行 |
|---|---|---|
| **变异红** | 导出树 `tree`，M1＋M3b | `--- FAIL: TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit (0.04s)` / `RC=1`，红行 `:84` |
| **还原绿** | 导出树 `tree`，拷回原件后 | `STATE 1 … 60/56`、`STATE 2 … 60/60 [ball]`、`--- PASS` / `ok … 0.227s` / `RC=0`（`logs/09-restored-green.txt`） |
| **原样绿** | 主检出（未变异树） | `--- PASS … (0.04s)` / `ok … 0.190s` / `RC=0`（`logs/02-pristine-after-edit.txt`） |

**还原核 `sha256`（三处逐字相同）**：

| 件 | (1) **入库 blob**（`git show 095699f2:…`） | (2) 主检出工作树 | (3) 仓外导出树 |
|---|---|---|---|
| `internal/config/writeguard_adopt_360_test.go` | `52a1aaee3bfe1d2b5b6eace8af062b1611ac91f6c381725ea81d3fbffe00eeb8` | 同串 | 同串 |
| `internal/config/writeguard.go`（**零改动**） | `efe12661e71b940f1fcf99d722f1c25e4403431997e90041491d5eae337d47bd` | 同串 | 同串 |

★`git rev-parse 095699f2:internal/config/writeguard_adopt_360_test.go` ⇒ `7077bc2b353f0009e593b5cebd3aa3a35d13bb4c` ＝ `git hash-object` 工作树同串 ⇒ **提交内容与工作树逐字相同**。
★上一枚那枚 93 行原件的 `sha256`＝`481fb162…62c7e`（HEAD blob），本腿改后为 `52a1aaee…eeb8`、**107 行**、CR 计数 **0**（`*.go text eol=lf` 守住）。
★`gofumpt -l` 单文件 ⇒ **空、`rc=0`**；`gofumpt -d` ⇒ **空**。

## 6. 零产码复核（五枚 `rev-parse` 全 SAME）

```
internal/config/allowdirs.go   HEAD=ef54b37c094ffd1349bcc8b27752313046a9b968 WT=ef54b37c… SAME
internal/config/writeguard.go  HEAD=0ebc415b05651104fa13d6e71c76e61ba8fdd93c WT=0ebc415b… SAME
internal/config/manager.go     HEAD=f53d8a0749cd5ceda1cc6700c7210a2ff019e063 WT=f53d8a07… SAME
internal/config/settings.go    HEAD=c4e3b09b21316a897412d7f04aff9e2406ff4061 WT=c4e3b09b… SAME
internal/config/permmode.go    HEAD=6f385353ec366ed1576512855bbd80ff97cd5c97 WT=6f385353… SAME
```

★**更宽一把尺**：`git ls-files 'internal/config/*.go'` 里**全部非测试** `.go` 逐枚 `rev-parse HEAD:` 对 `hash-object` ⇒ **DIFF 行 0**（全 SAME）。
★写面 `git status --porcelain -- internal/config docs/evidence/s1` ⇒ **仅** ` M internal/config/writeguard_adopt_360_test.go`（本件入库前另加本件）。

## 7. 门禁各 rc（主检出现跑；`GOFLAGS=`／unset `GOTMPDIR TMPDIR TEMP TMP`／未设 `CGO_ENABLED`）

- `GOFLAGS= go build ./...` ⇒ **rc=0**
- `go vet ./internal/config/` ⇒ **rc=0**
- `go test -count=1 -v ./internal/config/` ⇒ **rc=0**；**四数分开报**：`RUN=196`／`PASS=115`（顶层）／`FAIL=0`／`SKIP=0`，另 `^panic:` **0**（⛔ 相减）
- `sh scripts/d22scan.sh` ⇒ **rc=0**，末行 `d22scan: clean - no D22 ban violations`（`ban #8 internal/ 617 Go files` 是**分母**、⛔ 违规数）
- `sh scripts/check-path-length-budget.sh` ⇒ **rc=0**，`VERDICT GREEN`、`tracked paths=17575  over-budget=61  covered by roster=61  not in roster=0`
★**五枚在提交后 `095699f2` 上原样复跑**（`logs/32..36`）⇒ **rc 全 0、四数不变**（`RUN=196 PASS=115 FAIL=0 SKIP=0`、`panic=0`）。

### 格式门禁：甲乙两尺都报（裁定件 `pipe-format-gate-reads-tracked-blob-not-worktree.md`）

- **甲＝入库 blob 尺（判据）**：本腿新增件单文件 `gofumpt -l` ⇒ **空、`rc=0`**（本腿未新增 `.go`，只改既有件）。
- **乙＝工作树尺（只作诊断）**：整树读数属**既有债**（票 362 的 20 枚真脏 ＋ 票 361 的 15 枚 CRLF 假红，见 `360-v1-verdict-part2.md` §7）；★**乙尺跟踪集口径是 35、⛔ 55**（`A1134⑦a`）。本腿改动件在两尺上**均 0 命中**。★**件分两半**：本件＝§1–§7；`-part2` ＝§8–§12。
