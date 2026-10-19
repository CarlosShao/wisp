# 346-a1 件00 = CI 取日志尺与全量四数（本腿现跑，非引用他人）

## 取尺（照派单硬规矩）
1. `gh run view 39244425666 --repo CarlosShao/wisp --json jobs` → 数字 `databaseId`：
   test-core = 118032776497｜test-windows = 118032776659｜lint = 118032776906
2. `gh run view 39244425666 --repo CarlosShao/wisp --log-failed --job <数字ID>`（★⛔ 喂作业名）
   - 首发本腿带错 `--no-pager` → rc=1＋0 字节；错误已具名，改正后 rc=0。
   - lint（118032776906）本腿三次重试均 0 字节＝api.github.com 连接超时（errbytes=319）；
     判读＝**未知**（⛔ 读成"没有红"）；lint 那格不在本票必答射程，弃。
3. 两枚旧 run 全量日志本腿自取 sha256（同法归一后量）：39124122121 / 39101060753 / 39071959151。

## 归一尺（★四数只能取在顶格行上，本腿现量这一步的必要性）
GitHub `--log-failed` 每行 = `job<TAB>step<TAB>2026-...Z<payload>` ⇒ payload 绝不在行首。
尺 = `D:\tmp\wisp\20261019\346a1\norm.ps1`：剥两个 TAB 字段＋剥首个 `Z` 后时间戳＋剥一个前导空格。
★自证：归一**前**跑 `'^--- FAIL'` = 0 命中（正控 'FAIL' 子串命中数百）＝派单 `A1054`③ 那形本腿复现。
归一后文件＝本腿盘上件（仓外，`D:\tmp\wisp\20261019\346a1\`）：

| 归一文件 | 来源 | 字节 | sha256 前 12 | RUN | PASS_TOP | FAIL_TOP | SKIP |
|---|---|---|---|---|---|---|---|
| core.norm.log | 39244425666 / 118032776497 | 686396 | 760D9541D294 | 1833 | 1239 | 7 | 2 |
| win.norm.log | 39244425666 / 118032776659 | 683457 | 77C254BD6A8D | 1041 | 692 | 5 | 2 |
| old39124.norm.log | 39124122121（旧 run 全量） | 1389927 | CA9D1A225F6E | 2991 | 2037 | 13 | 6 |
| old39101.norm.log | 39101060753 | 1132935 | 7C032836115B | 2755 | 1858 | 14 | 4 |
| old39071.norm.log | 39071959151 | 1475365 | DE5382017C37 | 2849 | 1917 | 7 | 2 |

原始件字节（gh 直出）：1489719 / 1240950；sha256 前 12 ＝ CAAF2E9A8643 / A25AC9537009。
四数尺＝`grep -c '^=== RUN'`／`'^--- PASS'`／`'^--- FAIL'`／`'^--- SKIP'`（顶格锚定、只顶层、⛔ 混用⛔ 相减）。
⚠ 射程＝单作业文件、非整仓；上表 39124122121 是全量文件（含三作业）⇒ 四数比单作业大发，如实记。

## 本票三枚名的落位（全部尺＝`^(=== RUN|--- (PASS|FAIL|SKIP)): <名>` 于归一文件）
- `TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1` ⇄
  `TestPOSIX120SealWindowLeafSwapIsClosed`：只在 **test-core（Linux）** 现身；
  test-windows（win.norm.log）两枚 **0 命中**＝`//go:build !windows`（件03 逐文件尺）。
- `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`：只在 **test-windows 的
  cmd/wisp CLI 步骤** 现身（core.norm.log 0 命中）。
