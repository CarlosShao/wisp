# 346-a2 - AC#2 前半：pristine（未变异）弃权红读数

被测件＝仓外导出树 `export` 内的两枚测试件，**零字节改动**（哈希前后同，见本件末）。

## 甲 `TestPOSIX120SealWindowLeafSwapIsClosed`（posix_seal_window_120_other_test.go）

- 安静档（3 发）：全绿、窗口到达——totals 逐发
  `attempts=2000 hit-foreign=0 refused-at-check=971 refused-at-act=18 sealed-own=886`（注意第三发 refused-at-act=11/791 等，全表在 03 件）。
- 弃权档＝`-cpu=1`＋`WISP_120_ATTEMPTS=1`（3/3 红，rc=1），样本逐字（shot1；shot2/3 同句同数）：

```
posix_seal_window_120_other_test.go:365: AC#1 sample 1/4 shape=leaf-swap: attempts=1 hit-foreign=0 of-which-reported-nil=0 refused-at-check=0 of-which-refused-at-act=0 sealed-own-object=1 other=0
posix_seal_window_120_other_test.go:370: AC#1 total shape=leaf-swap: attempts=4 hit-foreign=0 of-which-reported-nil=0 refused-at-check=0 of-which-refused-at-act=0 sealed-own-object=4 other=0
posix_seal_window_120_other_test.go:370: AC#1 RED instrument is dead shape=leaf-swap: refused=0, sealed-own=4 - a zero window reading here proves nothing
--- FAIL: TestPOSIX120SealWindowLeafSwapIsClosed (0.00s)
```

- 机制：`-cpu=1` 下 racer goroutine 只进 runnext、主循环 4 发在首个 10ms 抢占点前跑完 ⇒ racer 一次翻面都没打 ⇒ refused=0（或 sealed-own=0 皆触发）⇒
  "真到不了窗口"的形状与 CI 观测（`refused=1998, sealed-own=0`）同句同族。

## 乙 `TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1`（openleaf_322_other_test.go）

- 安静档（3 发）：全绿、窗口到达（flips=1334/1192/1268、leafLegRefusals=271/237/359）。
- 弃权档＝`-cpu=1`（5 发＝3 红 2 绿，rc 如 00 件表）。**红发逐字**（shot1；shot3/4 同句、flips=82/106）：

```
openleaf_322_other_test.go:192: AC#1 paired reading on one screen: attempts=3000 flips=16 flipErrs=0 damagedWitnesses=0 leafLegRefusals=0 foreignsNarrowed=0 foreignsHoldingPayload=0
openleaf_322_other_test.go:206: AC#1 RED instrument: across 3000 raced attempts the leaf leg ("too many levels of symbolic links") answered 0 and the racer landed 16 flips (0 errors, first ""), so the damaged count above is not evidence about a window that was never reached
--- FAIL: TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1 (0.57s)
```

- **绿发（同条件）**：shot2 `flips=96 leafLegRefusals=1`、shot5 `flips=140 leafLegRefusals=1` ＝ 叶腿答上了 1 次 ⇒ 弃权支不响 ⇒ PASS。
  弃权判据就是 `liveness==0` 这一行——同条件只差"叶腿答没答上"，红⇄绿逐字可对。
- 与票面 CI 观测同族对照：CI run 39124122121 ＝ `answered 0 and the racer landed 6 flips`；本台面 ＝ `answered 0 and the racer landed 16 flips`（同句、分子不同）。

## 梯 2 的反向一味（如实记）

`-cpu=1`＋容器内 32 个空转进程：乙＝PASS（`flips=584 leafLegRefusals=9`）——负载把 syscall 拉慢、P 交出、窗口反而到达。
⇒ "更挤＝更弃权"⛔ 是直觉、⛔ 是本台面读数；弃权只认到达数（普查件 (iv)③ 同向）。

## pristine 哈希（每发前后同值）

| 文件 | sha256（export 树＝真实仓库，逐字同） |
|---|---|
| `internal/winsec/posix_seal_window_120_other_test.go` | `e8223d01f05629af7ec0d2c3e8e0c4b8167eed544333ae2839cac2ef2901311a` |
| `internal/winsec/openleaf_322_other_test.go` | `0e6aba36ed64849d76b6db4aa7cd4e36d3f136fc5a9a58482cdca7bef2e44774` |

- blob 同一性（真实仓库）：`git rev-parse HEAD:<path>` ＝ `0b0e5173169dc9cb10f9f8794212dfa9fe5e0d32` / `dfaddac8d0a5bd76154b45771794ebb19d4707ab`
  ＝ 票面 AC#0 记录的三锚同串值 ⇒ 本腿起手到收尾两枚仪器 blob 一字未动。
