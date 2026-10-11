# 72-a2 — 交件清单与"下一枚落地腿开工前必须先有人点头的那一件事"

## 1. 件册（全在 `.scratch/wisp/probes/72/a2/**`，⛔ 新建 `.go`，⛔ `.out`）

| 件 | 是什么 |
|---|---|
| `00-anchor.md` | 起手三读＋票 72/307 现读＋7 枚名册住址我自己重跑 |
| `10-seven-cases.md` | ★七问 × 7 枚逐枚（每一枚的"红那一行"都是本腿自己两发台面量出来的，⛔ 猜） |
| `20-answer-A-cause.md` | 必答 A：与 `syncdirs.go` 因果＝0（四条尺）＋红落在哪＋★票面 `AC#5` 回执第 4 条已过期 |
| `30-answer-B-line.md` | 必答 B：一把可执行的分界尺＋逐形最小诚实落点与代价＋**四枚点名"这样改就变假绿"** |
| `40-nails-that-would-break.md` | 12 条既有钉，逐条带尺与射程目录 |
| `50-what-i-did-not-measure.md` | 本腿⛔ 量的 10 件事，每件配"下一程该用哪把尺" |
| `logs/aliased-7named.txt` | 别名台面 × 甲族 7 枚：`rc=1`／FAIL=7／PASS=0／SKIP=0（**7 枚红句逐字在这枚件里，第 3/6/9/12/15/18/21 行**） |
| `logs/default-7named.txt` | 默认台面 × 同 7 枚：`rc=0`／FAIL=0／PASS=7／SKIP=0 |
| `logs/aliased-all16-bothfiles.txt` | 别名台面 × 两枚 windows 文件全 16 枚：`rc=1`／FAIL=7／PASS=9／SKIP=0（★台风形证据） |
| `logs/t03-broken-surface-note.txt` | 本腿⛔ 成立的那一发（漏搭台面→panic）的原始读数，逐字保存 |
| `logs/ac5-byte-ruler.txt` | 票面 `:82` 那把逐字节尺的现跑（junction rc=0／anchor rc=1） |

Go 面发数＝**4**，全部 `go test ./internal/risk/ -count=1 -v -run '<具名>'`；⛔ `./...`、⛔ `cmd/wisp`、⛔ `d22scan`、⛔ 建树变异、⛔ push、⛔ 翻框。

## 2. 一句话交付

**7 枚红全在测试自己体内，分成四形（I 期望串／II 前提句／III `want` 计数／IV 拿短拼当 OS 担保拼写），每一形都有同包里现成的正确写法可抄；判级一枚都没被证伪——别名台面上"票 72 的三条不变量用例"自己是绿的。**

## 3. ★落地腿开工**之前**必须由编排者点头的一件事（⛔ 由写腿自裁）

**票 72 已勾的 `AC#5` 用"那枚测试文件逐字节⛔ 变"当取证尺，而甲族 7 枚里有 4 枚住在那枚文件里。**
- 尺与件＝`40-nails` 第 1 条＋`logs/ac5-byte-ruler.txt`（现跑 `git diff --quiet a04d3e2 HEAD -- internal/risk/pathresolver_junction_windows_test.go`＝**rc=0，那把尺现在还在成立**）。
- 要裁的只有两条：**(a)** 编排者在票 72 尾部追加一节（⛔ 改 `:82` 原句），把 `AC#5` 的取证方式从"逐字节相同"换成"改后红位集合⛔ 变＋两台面 16 枚同跑"，然后放落地腿进那枚文件；
  **(b)** 只在 `anchor_spelling` 文件里治 #1/#2/#3，那 4 枚（#4/#5/#6/#7）留在原地继续红 ⇒ 别名台面/CI runner ⛔ 会绿。
- 本腿的读码意见（⛔ 判语、⛔ 授权）：**(a)** 是唯一走得通的，因为 (b) 交不出"票 72 `AC#4` 那条 runner 真 run 变绿"这件事，而那正是这枚票剩下的两枚未勾框之一。

## 4. ★需要人拍板／转述给别家的三件（除上面那件之外）

1. **台账一处待更正**（⛔ 本腿动原句）：票 72 `AC#5` 回执第 4 条那句"非测试面 `Gate(` 命中只有 `blacklist.go:76` 定义自己 ⇒ 今天还⛔ 生产 caller 传 override 表"**在 HEAD 上⛔ 成立**＝`internal/tools/mode.go:91/:93/:98` 是一枚拿着 `b.confirmations()` 的生产 caller。⇒ 后果：动 `overrideApplies` 现在**会改生产分档**（`BlacklistNote.AlreadyUnlocked` ⇄ `Unlockable`）。出处＝`20-answer-A-cause.md` A.4。
2. **⛔ 属票 72 射程、但本腿撞见的一枚脆面**：`TestClassifySpellingInvariance`（`anchor_spelling:78`）对 `filepath.VolumeName` 返回值零守卫 ⇒ 一台让 `os.TempDir()` 交回无卷标路径的机器/环境下会 **panic 掉整个测试二进制**（本腿第 3 发实锤：16 枚只跑到第 2 枚就崩、`--- PASS`=0）。它⛔ 在那 7 枚里、⛔ 是判级问题，属**仪器健壮性**那一族。要不要立票归编排者裁（若立，票名⛔ 超 100 字符，AGENTS/README 规矩 9 那枚雷）。
3. **`307-v1` 已经把一枚读数指回票 72**：commit `ed1d2bad` 标题逐字含"(d) 3/5 无信号裁归票 72"。本腿⛔ 读那枚件（射程外＋是裁决件）。⇒ 落地腿开工前**必须先自己读 `probes/307/v1/**`**，否则"恒真／无信号"那一问会记成两本账（票 307 `AC#5` 那句"⛔ 同一件事记两次"）。

## 5. 你需要做什么

见上面 §3（一件必须点头的事）与 §4（三件待转述）。**⛔ 本腿翻任何框；⛔ 本腿 push；产码⛔ 动过一字节。**
