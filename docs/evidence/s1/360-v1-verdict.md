# 票 360 终裁腿 `360-v1` 裁决件（`AC#2` 实现格 · 非实现者裁）— part1/2

- 裁者＝`360-v1`（⛔ 实现者：本腿零 `.go` 字节、零翻框、零票面改动、零 push）
- 被裁＝`360-r1` 四笔 `d68ddfae`／`fc7c85cc`／`0d907103`／`7c42145c`
- 裁决时刻 2026-10-22 18:1x–19:0x +08｜裁定时 HEAD＝`55ec4037`（编排者翻三格那笔）
- ★**结论：`AC#2` ⛔ 翻**（理由见 §2：两态正控**可被整臂摘除而不变红**）
- ★每项另起尺、独立现跑；⛔ 采信 `360-r1` 自报、⛔ 采信编排者的翻框
- ★件分两半：本件＝§1–§3（前提／`AC#2` 裁／零产码）；`-part2` ＝§4–§9（门禁／红名册／写面／`attrib` 分解／汇总）

## 1. 票面前提新鲜度（硬闸）—— 成立

对 **HEAD blob** 逐条复量（`git show HEAD:<p>` 管道 `Select-String`／`git grep`），与票面逐字同：

| 票面断言 | 我的读数 | 判 |
|---|---|---|
| `statOwnWrite` 5 行＝1 定义＋1 调用＋3 注释 | 定义 `allowdirs.go:157`／调用 `writeguard.go:179`／注释 `allowdirs.go:150`、`writeguard.go:105`、`writeguard_226_test.go:186` ⇒ **5 行、同分层** | ✅ |
| 写回入口 3 枚汇聚 `mergeWrite` | `allowdirs.go:142`／`permmode.go:79`／`settings.go:303` 均 `m.mergeWrite(...)`；定义 `writeguard.go:111`、`:118` 读盘当 base、`:163 SaveFile` **从不回读** | ✅ |
| Adopt 全仓一条 | `manager.go:142 CheckAndReload` → `:160 LoadFile` → `:193 plan.commit`；生产驱动 `config_reload.go:123`／`:157` | ✅ |
| 两条路互不相连 | `git grep 'm\.cur = '` 全仓 **仅** `settings.go:304`（失败回滚），⛔ 任何 adopt 写 | ✅ |
| `permmode.go` 旧调用点已消失 | `statOwnWrite` 在 `permmode.go` **零命中** | ✅ |
| 今天能复现（`file=60 mem=56`） | §2.1 状态 1 逐字 | ✅ |

⇒ **前提今天全部成立**，不触发票面 `AC#0` 的停手上报条②。

## 2. ★★`AC#2` 判语：⛔ 翻

**判〔甲〕形是否真的钉住了那一跳 ⇒ 答案是"钉住了那一跳，但没钉住那枚用例"。**

### 2.1 两态读数我自跑（⛔ 抄）—— 成立

`go test -count=1 -run TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit -v ./internal/config/`：

```
writeguard_adopt_360_test.go:54: STATE 1 after writeback: file.ball.size=60 mem.ball.size=56
writeguard_adopt_360_test.go:82: STATE 2 after one Adopt pass: mem.ball.size=60 file.ball.size=60 report.Hot=[ball]
--- PASS: TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit (0.05s)
ok  	github.com/CarlosShao/wisp/internal/config	0.121s     RC=0
```

★**两态能分辨**：`60/56` ⇄ `60/60 [ball]` **不同** ⇒ 用例**不是**恒真、台子能分辨两态。★连跑 5 次全 `RC=0`（0.079s／0.109s／0.155s／0.106s／0.098s）⇒ 不抖。

### 2.2 ★我自造的反向判据（⛔ 只看它报的）

全部在**仓外导出树** `D:\tmp\wisp\20261022\360v1\tree`（`git archive HEAD` ⇒ **17572/17572** 件、与 `git ls-files` 同数）上跑；每发先 `grep -n` 证落地 ＋ `go build ./...` `rc=0` 再读数，**读完立刻单独还原**。

| 发 | 造法（落地证明） | build | 读数 |
|---|---|---|---|
| **原样绿** | 无变异 | 0 | `STATE 1 … 60/56`、`STATE 2 … 60/60 [ball]`、`--- PASS` `RC=0` |
| **M1 摘 Adopt 臂**（`sed -i '73,81d'`；`grep -n CheckAndReload` 只剩注释 `:19`） | 代码调用 **0** | 0 | `STATE 2 … mem=-1 file=-1 report.Hot=[]`／**红行 `:76`** `the Adopt control did not produce the other state: … want 60 / 60 / [ball]; without it the readings above cannot tell two states apart`／`RC=1` |
| **M2 产码侧 adopt**（`writeguard.go:169` 后插 `m.cur = deepCopyConfig(base)`；`grep -n MUTATION-360v1` ⇒ `:170`） | 命中 `:170` | 0 | `STATE 1 … 60/60`／**红 `:60`**／**红 `:65`**（`divergence = [], want exactly [ball.size]`）／**红 `:85`**／`RC=1` |
| ★**M3 单独恒真**（`:55/:59/:64` 改 `if false`；`grep -n "if false"` ⇒ 3 命中） | 3 命中 | 0 | `STATE 1 60/56`、`STATE 2 60/60`、`--- PASS` **`RC=0`**（★照实记账成立） |
| ★★**M3′ 单独把 Adopt 臂自己的断言改恒真**（`:84` 改 `if false`） | 1 命中 | 0 | **`RC=0` 仍绿** |
| ★★**M1＋M3′ 摘臂 ＋ 把该臂断言改恒真** | `CheckAndReload` 只剩注释、`if false` 1 命中 | 0 | `STATE 2 … -1/-1 []`、**`--- PASS` `RC=0`** |
| **M3＋M2**（状态 1 三断言恒真 ＋ 产码 adopt） | 各命中 | 0 | **红 `:85`** `RC=1` ⇒ 正控臂**独立成牙** |
| ★**M4（我新造，不在它三发内）＝clobber 回归**：写回不再拿文件当 base（`writeguard.go:141`） | 命中 `:141` | 0 | `STATE 1 … 56/56`／**红 `:56`** `the file no longer carries the hand edit: … want 60`／**红 `:65`**／**红 `:85`**／`RC=1` |

还原核（逐发后即刻）：
- `writeguard_adopt_360_test.go` ⇒ `481fb1626f4307fe8ddca4d9d67dc43a7d20aadcc7cc8787cab5ee5071462c7e`（＝HEAD blob）
- `writeguard.go` ⇒ `efe12661e71b940f1fcf99d722f1c25e4403431997e90041491d5eae337d47bd`（＝HEAD blob）

### 2.3 ★那枚"M3 恒真仍绿"的解释：**⛔ 成立**

它与编排者的解释都是：「判据原文是『摘掉 Adopt 臂**或**把断言改成恒真 ⇒ 必须转红』，**M1 已满足"摘臂"那一半**」。
★**我判这个解释不成立**，两条硬读数：

1. ★**`或` 在判据里是"两种造法各自都必须转红"的并列**，不是"二者之一转红即算"的析取豁免。理由＝判据写它的**目的**是钉"用例有牙"：`摘臂` 测**状态 2 那一臂**、`改恒真` 测**断言本身**——**测的是两个不同的东西**。用 M1 的红去顶 M3 的绿，等于**用"臂在不在"的牙去顶"断言咬不咬"的牙**，两者不可互代。⇒ 编排者那句"**M1 已满足"摘臂"那一半**"是**把并列读成了析取**。
2. ★★**我造的那发 `M1＋M3′` 是决定性反证**：**同时**摘掉 Adopt 臂**和**把该臂自己的断言改恒真 ⇒ **`RC=0` 仍绿**。
   即：★**这条用例的两态正控可以"整臂连断言一起被删掉"而 CI 依旧全绿**——这正是判据 `改恒真` 那一半要排除的形状。⇒ **那一半没有被任何别的发覆盖**。
   ⇒ 它报的"M3 单独恒真仍绿"不是"诚实记账即可过关"的边角，而是**判据的一半真的漏了**。

★**它做对的部分（⛔ 抹）**：M3 的绿它**没有藏、没有算成红**（交付件 §3 逐字"⛔ 声称三发全红"），并且 `M3＋M2` 那发证明正控臂**独立成牙**（红 `:85`）。★但"照实记账"只解决**诚实**，⛔ 解决**判据是否满足**。

### 2.4 ★缺什么（补法，供下一腿）

- **缺的是一发"断言级"的牙**：现在 `:84`（臂的三条值断）与 `:89`（两态不同）**重叠**——`:89` 只在 `adoptedMem == memSize` 时咬，而摘臂时 `adoptedMem=-1`、`memSize=56`，**`:89` 不咬**（这就是 `M1＋M3′` 仍绿的机理）。
- **最小补法**（仍属〔甲〕形、零产码）：把"该臂**真的跑过**"独立成一条**不可被同臂断言吞掉**的断——例如在摘臂后仍为真的**独立哨兵**（`ranAdopt` 布尔／`hot` 的 `nil` 与非 `nil` 区分／或把 `:89` 改成 `adoptedMem < 0 ⇒ Fatal("the Adopt arm did not run")` 的**前置**断，放在 `:84` 之前）。这样 `M1` 与 `M3′` 各自单独都必须红。
- ★**判据原文与 `AC#1` 的裁定文字 ⛔ 改**（那是"改契约文字"、属摆机主的三类之一）⇒ 本腿只判"未满足"、⛔ 代改。

## 3. 零产码复核 —— 成立（且比它报的更宽）

- 它报的五枚逐枚 `git rev-parse HEAD:<p>` 对 `git hash-object <p>`：`allowdirs.go ef54b37c`／`writeguard.go 0ebc415b`／`manager.go f53d8a07`／`settings.go c4e3b09b`／`permmode.go 6f385353` ⇒ ★**五枚我全部 SAME**（复现它的读数）。
- ★**我另加两把更宽的尺**（它没做）：
  - `internal/config` 全部 **17 枚**非测试 `.go` ⇒ `diff_count=0`（全 SAME）。
  - 全仓 `git diff --name-only d68ddfae^..7c42145c` 里非测试 `.go` ⇒ **NONE**。
  - `git diff --name-status d68ddfae^ 7c42145c -- internal/config/` ⇒ **仅** `A internal/config/writeguard_adopt_360_test.go`（唯一一枚新增件）。
  ⇒ ★**〔甲〕形"零产码"是真的**，且范围不止它报的五枚。
