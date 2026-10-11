# 表①——`A817` 那 12 枚的归属与包名册（分母）

台面 HEAD＝`ec87076f`（本腿第 1 笔 `63bc6335` 之前现量）。

## 0. 这把尺（逐字，可重跑）

```
$ for n in <12 枚名>; do grep -rn "func $n(" --include=*_test.go . ; done
```

- 射程＝**仓根 `.` 全树**（⛔ 只数 `internal/risk`）；`--include=*_test.go`。
- 形＝`func <名>(` ⇒ **注释行不可能命中**（命中面只有函数声明行），也⛔ 匹配子测试。
- 枚数结果：**12 枚名 × 每枚恰好 1 行命中**＝12 行；**0 枚多命中／0 枚零命中**。⇒ 名册**没有**"排除 N 枚"这一格（N＝0，名册行数＝12＝CI 名册枚数）。
- 顶层名 vs 子测试尺：CI 那把尺（`r9-ci-after-red-roster.txt`）数的是 `--- FAIL: <名>` 的**顶层名**；本腿 12 枚在本地日志里也全是顶层行（`grep -c '^--- FAIL'`＝12，见 `logs/t02-shorttemp-12.txt`）⇒ 两把尺同一分母，⛔ 拿子测试充数。

## 1. 名册（12 行＝12 枚）

| # | 用例 | 包（import path 尾段） | 文件 | 内容锚（这枚钉的是哪一句判据） | 声明行（快照） |
|---|---|---|---|---|---|
| 1 | `TestClassifyAnchorSpellingIsNotVerdict` | `internal/risk` | `pathresolver_anchor_spelling_windows_test.go` | `:59` `Classify(res.Canonical) != ClassA`（junction 拼写的锚下 A 表必须仍赢）；`:63` `Gate` 不得 Allow/NeedL2 | `:29` |
| 2 | `TestCanonicalInputGainsNoSecondForm` | `internal/risk` | `pathresolver_anchor_spelling_windows_test.go` | `:139` `len(f.spellings()) != 1`（票 72 的爆炸半径上界：已解析拼写⛔ 长出第二种形） | `:125` |
| 3 | `TestAListWinsWhereBothTablesHit` | `internal/risk` | `pathresolver_anchor_spelling_windows_test.go` | `:237` `Gate(out, {lower(out):true})` 必须 Allow（**B-only 正控**：A 的赢面不是"整棵树都被拒"） | `:204` |
| 4 | `TestPathResolverShortNameAListDenied` | `internal/risk` | `pathresolver_junction_windows_test.go` | `:135` `EqualFold(res.Canonical, long)`＝8.3 必须展开成长形；`:138` Classify=ClassA；`:141` 不得 granted | `:123` |
| 5 | `TestPathResolverUNCAListDenied` | `internal/risk` | `pathresolver_junction_windows_test.go` | `:164` 三种 UNC 拼写都必须归一到 drive 长形；`:167` 不得 granted | `:148` |
| 6 | `TestPathResolverExtendedLengthPrefixAListDenied` | `internal/risk` | `pathresolver_junction_windows_test.go` | `:182` `\\?\` 拼写必须归一；`:185` 不得 granted | `:174` |
| 7 | `TestBListDefaultDenyAndOverride` | `internal/risk` | `pathresolver_junction_windows_test.go` | `:288` B 表单文件 override 必须 Allow；`:291` 必须落一条日志 | `:272` |
| 8 | `TestSyncFixtureFallbackAndMatch` | `internal/risk` | `syncdirs_test.go` | `:161` fixture 根下的写必须 sync（**根成员判定**）；`:169`/`:172` fallback 不得被弱档解除 | `:152` |
| 9 | `TestSyncFallbackNotDisarmableByWeakRoot` | `internal/risk` | `syncdirs_test.go` | `:210`/`:213` 四种 `Source`（default/fixture/options/空）都⛔ 算 confirmed、under-profile 必须仍 armed | `:198` |
| 10 | `TestSyncUnverifiedRootKeepsFallback` | `internal/risk` | `syncdirs_test.go` | `:230`/`:233` C26 验不了的 registry 根⛔ 解除 suspect 网 | `:224` |
| 11 | `TestSyncSuspectFallbackIsComponentBounded` | `internal/risk` | `syncdirs_test.go` | `:354` under-profile 必须 suspect；`:357` 同名前缀的兄弟目录 `profileevil` ⛔ 被扫进（N-5） | `:344` |
| 12 | `TestSyncSuspectFallbackWhenUndetectable` | `internal/risk` | `syncdirs_test.go` | `:401` under-profile 必须 `Sync && Root.Provider=="sync-suspect"`；`:405` profile 外⛔ suspect（`outOfProfile` 自证前提） | `:390` |

⇒ **包归属＝12/12 全在 `internal/risk` 一枚包里**（文件 3 枚）。
⚠ 派单让我别默认"12 枚都在 `internal/risk`"——**这一枚默认在这颗 HEAD 上成立**，`A817` §5 那句"internal/risk 那 12 枚"是准确的，⛔ 需要更正。

## 2. 编译面差（`A817` 没写、归口会用到的那一半）

- 名册 1–7 ＝ 文件带 `//go:build windows`（`pathresolver_anchor_spelling_windows_test.go:1`、`pathresolver_junction_windows_test.go:1`）⇒ **只在 windows 档进编译面**。
- 名册 8–12 ＝ `syncdirs_test.go` **无 build tag** ⇒ POSIX 也编译、也求值。
  ⇒ 后果：`syncdirs` 那 5 枚的修法**⛔ 能写成"给文件加 windows tag"**（那会把它们从 `--scope=core` 的 ubuntu 分母里也一起摘掉）；`A817` §3 只说它们是"under-profile fallback 家族"，没写它们**跨平台在跑**这一条。本腿把它补进分母描述。

## 3. 与 `A817` 的不符处（具名报回，⛔ 改台账）

1. ⚠ **`A817` §3 说"`TestAListWinsBothTablesHit`、`TestBListDefaultDenyAndOverride` ＝同文件同族"——盘上⛔ 同文件**：前者在 `pathresolver_anchor_spelling_windows_test.go:204`，后者在 `pathresolver_junction_windows_test.go:272`。**同包、不同文件、不同主题**（一枚是 A∩B 归口的 B-only 正控，一枚是 B 档 override 行为）。⇒ "同族"这一半成立、"同文件"这一半过期；本腿按名册重归族（两枚同因＝override 键的拼写面，见 `20-mechanism-per-red.md` #3/#7）。
2. ✅ `A817` §3 给的 `syncdirs_test.go` 五个红点行号（`:162/:213/:233/:355/:402`）在本腿**复跑全中**（现量逐字相同），⛔ 照抄＝本腿自己跑出来的（`logs/t02-shorttemp-12.txt`）。
3. ⚠ `A817` 把 12 枚统称"`internal/risk` 那一族（A/List 表判定与 under-profile fallback）"＝**两族合一枚名册**；本腿的名册把它们分成三族（拼写面 7 枚／override 键面 2 枚／sync 归属面 5 枚），因为表② 的判定**⛔ 同判**（前两族＝甲、第三族＝乙）。
