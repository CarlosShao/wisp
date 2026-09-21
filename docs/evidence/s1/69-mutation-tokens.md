# 69 — C21 表↔码双向断言的变异检验原始输出（AC#2）

`s1/` 证据 · 2026-09-21 · agent=ticket69 · 被测用例：`internal/ball/tokens_table_test.go`

纪律（仓内血的教训，A24-D4 与票 66/71 都栽过）：**每一次变异都先用 `grep` 证明它真的落盘**
（`sed` 静默不匹配会造出假绿），再跑测试；跑完立刻还原，并证明
`grep -c MUTATION` = 0 **且** `git diff` 对那两个文件为空。

三次变异（码侧 1 次、表侧 2 次），全部还原。命令与原始输出一律照抄。

## 基线（变异前，`internal/ball` 全绿）

```
$ gofmt -l internal/ball                       # 空
$ go vet ./internal/ball/                      # rc=0
$ go test -count=1 -v ./internal/ball/ | grep -E "^--- (PASS|FAIL)"
PASS=48 FAIL=0
$ go test -count=2 ./internal/ball/
ok  	github.com/CarlosShao/wisp/internal/ball	0.555s
```

新增的 4 条用例（本票的全部新断言都在这 4 条里）：

```
--- PASS: TestC21TableColourRowsMatchCode (0.00s)
--- PASS: TestC21GeometryRowsMatchCodeConstants (0.00s)
--- PASS: TestC21TokenConsumerReport (0.02s)
--- PASS: TestC21CodeTokensAreTabledOrExempt (0.19s)
```

## M1 — 码侧：`internal/ball/tokens.go` 的几何常量改一位

`tokens.go:361` `DockTriggerPx = 16` → `DockTriggerPx = 17 // MUTATION-69`

变异落地证据：

```
$ grep -c MUTATION internal/ball/tokens.go
1
$ grep -n "DockTriggerPx = " internal/ball/tokens.go
361:	DockTriggerPx = 17 // MUTATION-69
$ git diff --stat -- internal/ball/tokens.go
 internal/ball/tokens.go | 2 +-
 1 file changed, 1 insertion(+), 1 deletion(-)
```

跑测试（`-count=1`）：

```
$ go test -count=1 -run TestC21GeometryRowsMatchCodeConstants -v ./internal/ball/
=== RUN   TestC21GeometryRowsMatchCodeConstants
    tokens_table_test.go:838: docs/evidence/s1/c21-native-tokens.md:155: DockTriggerPx = 17, but this row states no such number (its numbers are: 160ms, 0.42, 17px, 96dpi, 32) - table and code drifted
--- FAIL: TestC21GeometryRowsMatchCodeConstants (0.02s)
FAIL	github.com/CarlosShao/wisp/internal/ball	0.120s
```

⚠ 顺带把 `16px` 记进了"无人认领的数字"清单（21 条），也就是说**表侧那个孤立数字**同样看得见：

```
NUMBERS NO NAMED CONSTANT CLAIMS (21): ... 16px @ docs/evidence/s1/c21-native-tokens.md:155, ...
```

全套件在同一变异下只有这一条红：

```
$ go test -count=1 -v ./internal/ball/ | grep -E "^--- FAIL"
--- FAIL: TestC21GeometryRowsMatchCodeConstants (0.00s)
$ go test -count=1 -v ./internal/ball/ | grep -c "^--- PASS"
47
```

⇒ **既有测试对这条几何漂移全盲**（含 `dock_test.go:118 const trigger = DockTriggerPx`，
它是从常量推出来的，改值不会红）。这正是 A24-D4 说的空档。

还原：

```
$ grep -c MUTATION internal/ball/tokens.go
0
$ grep -n "DockTriggerPx = " internal/ball/tokens.go
361:	DockTriggerPx = 16
$ git diff --quiet -- internal/ball/tokens.go && echo "tokens.go diff: EMPTY"
tokens.go diff: EMPTY
```

## M2 — 表侧：look 色改一位

`docs/evidence/s1/c21-native-tokens.md:218` `hex(0xFDBA74,0.50)` → `hex(0xFDBA75,0.50)`
（标记 `MUTATION-69` 放在该行的 CSS 列，不污染 Go 值列的解析）。

变异落地证据：

```
$ grep -c MUTATION docs/evidence/s1/c21-native-tokens.md
1
$ grep -n "solar\].Glow" docs/evidence/s1/c21-native-tokens.md
218:| （无，原生自有）MUTATION-69 | —— | `looks[solar].Glow` | `hex(0xFDBA75,0.50)` |
$ git diff --stat -- docs/evidence/s1/c21-native-tokens.md
 docs/evidence/s1/c21-native-tokens.md | 2 +-
 1 file changed, 1 insertion(+), 1 deletion(-)
```

跑测试：

```
$ go test -count=1 -run TestC21TableColourRowsMatchCode -v ./internal/ball/
    tokens_table_test.go:602: docs/evidence/s1/c21-native-tokens.md:218: looks[solar].Glow drifted - the table states hex(0xFDBA75,0.50) (rgba(253,186,117,0.500000)) and internal/ball/tokens.go holds rgba(253,186,116,0.500000)
    tokens_table_test.go:672: verified 113/114 colour rows against tokens.go (40 palette fields checked in code -> table, 36 look colours, 4 looks)
--- FAIL: TestC21TableColourRowsMatchCode (0.00s)
FAIL	github.com/CarlosShao/wisp/internal/ball	0.044s
```

全套件同样只有这一条红（实测：`--- FAIL` 计数 1、`--- PASS` 计数 47、`-v` 总行数 48）。
`113/114` 这个自报数字本身就是防"静默少测"的：分母是表里真实存在的配色行数。

还原：

```
$ grep -c MUTATION docs/evidence/s1/c21-native-tokens.md
0
$ git diff --quiet -- docs/evidence/s1/c21-native-tokens.md && echo "M2 restore: diff EMPTY"
M2 restore: diff EMPTY
```

## M3 — 表侧：几何行的数字改一位（码不动）

`docs/evidence/s1/c21-native-tokens.md:155` 「距边 **16**px 内开始挤压」→ 「**17**px」。
这条是 M1 的镜像：**改表不改码**。

变异落地证据：

```
$ grep -c MUTATION docs/evidence/s1/c21-native-tokens.md
1
$ git diff --stat -- docs/evidence/s1/c21-native-tokens.md
 docs/evidence/s1/c21-native-tokens.md | 2 +-
 1 file changed, 1 insertion(+), 1 deletion(-)
```

跑测试：

```
$ go test -count=1 -run TestC21GeometryRowsMatchCodeConstants -v ./internal/ball/
    tokens_table_test.go:838: docs/evidence/s1/c21-native-tokens.md:155: DockTriggerPx = 16, but this row states no such number (its numbers are: 160ms, 0.42, 17px, 96dpi, 32) - table and code drifted
--- FAIL: TestC21GeometryRowsMatchCodeConstants (0.01s)
FAIL	github.com/CarlosShao/wisp/internal/ball	0.120s
$ go test -count=1 -v ./internal/ball/ | grep -E "^--- FAIL"
--- FAIL: TestC21GeometryRowsMatchCodeConstants (0.00s)
$ go test -count=1 -v ./internal/ball/ | grep -c "^--- PASS"
47
```

还原（三次变异之后，全仓 `internal/`+`docs/` 再无任何 `MUTATION-69` 痕迹）：

```
$ git status --short -- internal/ball/tokens.go docs/evidence/s1/c21-native-tokens.md
(空)
$ git grep -n "MUTATION-69" -- internal docs
(无输出，rc=1)
```

## 覆盖与判定范围（本票到底锁住了什么）

| 方向 | 判据 | 今天被检查的行/字段数 |
|---|---|---|
| 表 → 码（配色） | 每行的 `hex()/rgba()` 构造按 `tokens.go` 同一套 `hex()`/`rgba()` 求值后**逐通道相等** | 40 暗 + 38 亮 + 36 look = **114 行** |
| 表 → 码（几何/动效） | 行点名的每个常量必须存在，且码里的值必须是该行**真写出来的数字**（单位按名字承诺：`*Ms` 允许 `1.6s` 这类秒写法，`*Px` 要 px，`*FPS` 要 fps，比率/透明度/增益要裸数） | **25 行 / 56 个常量**（55 个有数值：51 个带单位标记、**4 个只能裸匹配** = `BallSizeMinPx`/`BallSizeMaxPx` ← 「可配 44–72」，`FontSizeMonoPx`/`FontSizeMicroPx` ← 「12.5 / 11」，这 4 个由测试日志显式列出；`FontFamily` 是字符串无数值） |
| 码 → 表 | `Palette` 40 个 Color 字段、`looks` 4×9 个 look 色、`tokens.go` 56 个导出常量，逐一要求**恰好一行**声明；亮表缺的字段必须真的沿用 `:root`（级联断言：`BallHalo`、`OnSolid`） | 135 个码侧 token |
| 零消费者 | 只报告不判红（见票面清单） | 报告 **29 条** |
| 表外码侧常量 | 既不在表里也不在 `c21OutTableExempt` 显式豁免里 ⇒ 红 | 豁免 **15 条**（含理由） |

## 本票没能证明的部分（如实登记）

1. **表 ↔ `design/assets/tokens.css` 仍无机器检查**。表头那句「权威真相源 = `tokens.css`」
   今天依旧只有票 12 AC#7 的人工对账背书（A24-D6 把 80 条面板侧声明显式判为域外）。
   本票锁的是**表 ↔ 码**，不是**表 ↔ CSS**。
2. **几何行是"值出现在行内数字里"的判定，不是双射**。同一行里两个数撞号可以互相掩盖：
   把 `BallSizeMaxPx` 从 72 改成 44，行内仍有 44 ⇒ 该断言不红。
   兜住它的是"无人认领数字"那份日志（改完之后 72 会掉进无人认领列），**那是报告不是门**。
3. **`FontFamily` 的值没有任何断言**（表只写 `--font-sans` 引用，不写字符串值）。
   它被显式列进 `c21StringTokens`，名字双向仍检查，值不检查。
4. **零消费者报告按字段名匹配 selector**：若哪天 `Visual` 也长出一个叫 `Glow` 的字段，
   `looks[*].Glow` 会被误判为"有消费者"（今天实测无同名冲突：`Visual` 用的是
   `GlowColor`/`RingColor` 这类全名，`grep` 过）。报告方向偏乐观，不会偏严。
5. **`looks[*].Deep` 的零消费者是本票新发现**，但"哪一套 look 在被渲染"由运行时
   `SetLook` 决定，测试只能锁住四套的**值**，锁不住"owner 到底要哪套"。
