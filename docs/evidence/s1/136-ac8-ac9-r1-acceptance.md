# 136 — AC#8（判据②③）＋ AC#9（①②③）：对抗验收（第二格）

- 验收方：`acceptor-ticket136-ac8-ac9-r1`（**非实现者**；实现方＝`worker-ticket136-ac8-ac9`）
- 判据物：票 `.scratch/wisp/issues/136-*.md` 的 **AC#8②③**（按编排者 23:3x 的「措辞更正」那条性质判）与 **AC#9①②③**
- 本程**一个字未改** `internal/observe/**`（只读），仓内未建 worktree、未 checkout；所有变异落在仓外纯净快照树
- 本程**未翻任何 AC 勾**（AC#8／AC#9 的勾由编排者按本表翻）
- 证据文件渐进写：§0／§1／§2… 每裁完一节 commit 一次

---

## §0 锚点、树与基线（开工第一步自己量，不抄派单）

| 项 | 实测 |
| --- | --- |
| 开工第一次 `git rev-parse --short HEAD` | **`76662d8`**（＝本程锚点，所有快照都从它或它的父 commit 抽） |
| 那一刻的 `date -u` 原文 | `Wed Sep 23 15:16:02 UTC 2026` ⇒ 本机 **23:16 +08**（+8 手工换算，没塞进 `date` 的格式串） |
| 分支／开工 `git status --porcelain` | `dev`；工作树只有 1 枚未跟踪件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`（**未读、未提交、未改、未删、未据它改任何判据**） |
| 实现方自报锚点 | `1d38206` |

### 0.1 `1d38206..76662d8` 在 `internal/observe/` 的差集（有没有"别人动过被测码"）

```
$ git log --oneline 1d38206..76662d8 -- internal/observe/
2f291d0 test(136,AC#9①): 给 CheckSettle 的零可信样本面装钉（与 AC#1 同形，两腿）
79ddd49 test(136,AC#8②): 给 sampler_test.go:308 那枚越界直取加长度守卫（要红不要静默）
$ git diff --stat 1d38206..76662d8 -- internal/observe/
 internal/observe/sampler_settle_zerosample_136_test.go | 153 +++++++++++++++++++++
 internal/observe/sampler_test.go                       |  10 ++
```

⇒ 差集**只有实现方自己那两枚 commit**，没有第三方动过被测面。再往宽处核了一发（派单没要求，但"别人有没有在飞我的包"只能这样问）：

```
$ git diff --name-status 1d38206..76662d8 -- '*.go'      # 全仓 .go 差集
A	internal/observe/sampler_settle_zerosample_136_test.go
M	internal/observe/sampler_test.go
```

⇒ **两枚锚点之间全仓只动了这两枚 .go 文件**，都是本程被验对象自己写的。兄弟在飞的 `cmd/wisp`（票 133）与 `internal/winsec`（票 137）**在本差集里零命中** ⇒ 本程既没读也没被它们影响；两枚 sha 之间交错的都是 `.md`／票面／台账。

被验面的三态归属（自己 `git diff -U0` 复核，不采信报告）：

| 文件 | 1d38206..76662d8 | 结论 |
| --- | --- | --- |
| `internal/observe/sampler.go`（被测生产码） | **无输出** | 本程生产码零改动 ⇒ AC#9④"不许改 SampleState 那侧语义"成立 |
| `internal/observe/sampler_test.go` | `10 增 0 删`，唯一 hunk 在 `@@ -307,0 +308,10 @@`，加的是 `if len(rep.Samples) == 0 { t.Fatalf(...) }` | 只动那枚用例；`t.Fatalf` 不是 `t.Skip`；无删除行 |
| `internal/observe/sampler_zerosample_136_test.go`（AC#1 那枚钉） | **无输出** | AC#1 的断言一字未改（AC#1 已终裁 PASS，本程不重判） |

### 0.2 仓外纯净树（只建不删，全部带本会话后缀 `wisp136r2-ac8ac9-*`）

| 树 | 来源 | 用途 |
| --- | --- | --- |
| `/d/tmp/wisp136r2-ac8ac9-noguard` | `git archive 1d38206` | AC#8① 那一侧的**分母**（无守卫 ⇒ 同一发 M3 今天确实吞读数） |
| `/d/tmp/wisp136r2-ac8ac9-nailless` | `git archive 79ddd49` | AC#9② 的**改前对照**（有守卫、无 settle 钉） |
| `/d/tmp/wisp136r2-ac8ac9-anchor` | `git archive 76662d8` | 本格主读数面（守卫＋两枚新用例都在） |
| `/d/tmp/wisp136r2-ac8ac9-probe` | `git archive 76662d8` ＋ 我自己两枚探针文件 | §3 探哨兵、§4 判 `sampler_test.go:31` 可达性；**探针文件不在仓库里、不在被验版本里** |
| `…-pristine-{noguard,nailless,anchor,probe}` | 各树 `internal/observe/*.go` 的拷贝 | 每发变异**先无条件 restore** 的源 |
| `/d/tmp/wisp136r2-ac8ac9-mut.py` | 本程自写驱动器 | 锚点文本命中数≠期望值 ⇒ 直接退出并回滚（防叠发、防打偏；本程真拦住了一次，见 §2.5） |
| `/d/tmp/wisp136r2-ac8ac9-counts.sh`、`-battery.sh`、`-battery2.sh` | 计数与批量 | 四数只从 `-v` 量 |

工具链写明：`go version go1.27.1 windows/amd64`；`gofumpt` 用**本机既有二进制** `/d/work/base/gopath/bin/gofumpt.exe`（我实测其 `--version` ＝ **`v0.12.0 (go1.27.1)`**，与实现方写明的一致）；**本程没有跑过任何 `go install …@latest`**（派单明令：那会升掉宿主工具链）。容器 `golang:1.27` ＝ `go1.27.1 linux/amd64`（§5）。

### 0.3 基线（未变异，`go test -count=1 -v ./internal/observe/`）

| 树 | rc | RUN | PASS | FAIL | SKIP | panic |
| --- | --- | --- | --- | --- | --- | --- |
| `noguard`（1d38206） | 0 | 56 | 56 | 0 | **0** | 0 |
| `nailless`（79ddd49） | 0 | 56 | 56 | 0 | **0** | 0 |
| `anchor`（76662d8） | 0 | **58** | 58 | 0 | **0** | 0 |

⇒ 58 ＝ 56 ＋ 本程被验的两枚新用例，与实现方 §4.1 的"装钉后 58"逐位相同。

**两笔如实账（不粉饰）：**

1. `noguard` 基线**第一发**就是 `rc=1 / 56 / 55 / 1 / SKIP0`，红名 `TestNoopTaskReturnsToBaseline`、红点 `goroutine_test.go:33: PerTask mid-task = 2, want 3` —— 这是票面 **AC#11** 那枚既有 flake 在**本程锚点上的第三次独立命中**（前两枚是实现方量的）。重取一发才得到上表那行 56/56。它不在本格两格判据上（本格三态都判在"红名是本程仪器"的读数上），但下一位读数的人还会撞到。
2. 我有一发 `go test -count=2 -v` **取在了还没 restore 的变异树上**（读到 52 RUN/47/5＋panic），当场识破、`restore` 后重取（§5 G4 用的是重取那份）。归因：我自己 sequencing 失误，不是被验面的性质。
