# 307-r1 · 50-ac5-ledgers —— 三本账对账：逐枚写清"⛔ 翻谁的格"

盘上现量（2026-10-11 10:3x，本腿自己跑；⛔ 引用前必重跑）：

| 账 | 尺 | 现量 | 本腿动过它没有 |
|---|---|---|---|
| 票 307 | `grep -cE '^[[:space:]]*- \[x\]'` ⇄ `'- \[ \]'` | **1 勾／5 未勾**（`AC#0` 已由编排者翻） | **⛔**（`git status --porcelain -- .scratch/wisp/issues docs` ⇒ **0 行**；本腿写面内无票面文件，理由＝`40-ac4-gates.md` §6 那处两条相撞） |
| 票 72 | 同上 | **4 勾／2 未勾** | **⛔ 一枚框** |
| 票 252 | 同上 | **4 勾／2 未勾** | **⛔ 一枚框** |
| 台账 `A817`／`A833` | `grep -n 'A817' docs/reports/pending-and-issues.md` ⇒ `:15313` 起；`A833` ⇒ `:15807` | `A817` 欠账⑥"12 枚 `internal/risk` CI 红仍〔待归因〕"逐字在账（`:15451`/`:15487`/`:15595` 三处提及） | **⛔**（台账只追加不删，且**追加归编排者**：`docs/reports/**` ⛔ 本腿写面） |

## 1. ① `A817` 那"12 枚〔待归因〕"＝**票 307 ＋ 票 72 合销**，⛔ 本票独吞

本腿**跑出来的分族证据**（⛔ 抄 `risk-red-attrib-1`／编排者的读数，两台面各两发，件＝`logs/full-aliased-pre.txt`／`logs/full2-aliased-pre.txt`）：

- 别名台面 12 枚红＝**甲 7 ＋乙 5**，两发**逐字同集合**。
- **乙 5 枚**＝`TestSyncFixtureFallbackAndMatch` / `TestSyncFallbackNotDisarmableByWeakRoot` /
  `TestSyncUnverifiedRootKeepsFallback` / `TestSyncSuspectFallbackIsComponentBounded` / `TestSyncSuspectFallbackWhenUndetectable`
  ⇒ 修复后**全绿**（`logs/full-aliased-post.txt`），本机台面**⛔ 动色**。⇒ 这 5 枚**由本票销**。
- **甲 7 枚**＝`TestClassifyAnchorSpellingIsNotVerdict` / `TestCanonicalInputGainsNoSecondForm` /
  `TestAListWinsWhereBothTablesHit` / `TestPathResolverShortNameAListDenied` / `TestPathResolverUNCAListDenied` /
  `TestPathResolverExtendedLengthPrefixAListDenied` / `TestBListDefaultDenyAndOverride`
  ⇒ 修复后**仍全红、枚数⛔ 减（7→7）、名字⛔ 换**（同一把尺逐字对得）。
  ⇒ **⛔ 本票的账**：它们红在"用例期望侧写死了 `t.TempDir()` 的拼写"（票面追加节 §2 甲族那句），
  **归口＝票 72**，本腿⛔ 动它们一字。⇒ 台账那一格的**合销＝票 307 结乙 5 ＋ 票 72 结甲 7**，两本账⛔ 各记一次同一件事。
- 本腿⛔ 新增红（两台面、四发、`comm -13` 尺＝0 枚）⇒ 归因表⛔ 因本票而需要改写；
  但"12 → 7"这一枚**颜色变化本身**属本票凭据，翻账动作⛔ 归编排者（`docs/reports/**` ⛔ 本腿写面）。

## 2. ② 票 252＝**只引当形状先例**

- 本票引它的**唯一用途**＝"一条路径两种拼法分家"这一**形状**的先例（`syncdirs.go` 里的两支比较是它的**第二枚调用点**，非同一枚消费者）。
- **⛔ 在 252 身上翻任何一格勾**（现量 4 勾／2 未勾，见上表尺）——本腿⛔ 打开过那枚文件写字。
- **⛔ 扩 252 的面**：本票⛔ 把它没闭的 `AC#1`/`AC#6`（非实现者读数）与 `AC#4`/`AC#5`（`cmd/wisp` 空）拖进自己的收口。
- **方向相反那一枚必须写死**：票 252 那族＝根内写读成越界＝**多报**；本票＝资料目录下的写读成不可疑＝**少报**。
  本票的修复**⛔ 改变 252 的任何后果面**（本腿⛔ 动存值、⛔ 动 `finalize()`／`canonical` 分级、⛔ 碰 `resolveTarget` 那两句 fail-closed），
  凭据＝`git diff --stat` 那 20 行改动全在 `match()` 的两枚比较上（`internal/risk/syncdirs.go | 20 +-`，见 `git show --stat 51a0a0f2`）。
- 编排者在追加节 §2 第 2 条已经把"⛔ 并回 252"裁死并给了撤销口令（「307 并回 252」）⇒ 本件只是**照记**，⛔ 新裁。

## 3. ③ 票 72 那 7 枚的**期望侧**⛔ 进本票

- 本票对甲 7 枚做的事＝**零**（⛔ 改期望、⛔ 改 tag、⛔ 搬档、⛔ 让它们"顺手变绿"）。
  尺＝修复后别名台面上那 7 枚**照样红**（第 1 节），⇒ 本腿⛔ 用任何一种方式替票 72 收口。
- 本票**新增**的两枚夹具（`syncdirs_alias_307_windows_test.go`）⛔ 属票 72 的名册、⛔ 改它分母：
  它们是 **`//go:build windows`** 新文件里的**新名字**，那 5 枚指名用例仍住⛔ tag 的 `syncdirs_test.go`
  （尺＝`head -3 internal/risk/syncdirs_test.go` ⇒ `package risk` 起头，无 tag ⇒ POSIX／core 档分母未动，票 301 包级 guard D 未撞）。
- 票 72 现态 4 勾／2 未勾、⛔ 动框（票面"边界"节＋追加节 §2 都写死），由它下一程自己读这一节。

## 4. 一句话交账

**⛔ 同一件事记两次**：乙 5 枚的凭据只落本票（`10-ac1-colours.md` ＋ `20-ac2-shapes.md` ＋ `30-ac3-fixture.md`），
甲 7 枚的账⛔ 由票 72 结，形状先例的账⛔ 由票 252 结；台账 `A817`／`A833` 的**翻格动作本腿⛔ 执行也⛔ 请求**。
