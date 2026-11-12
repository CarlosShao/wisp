# 346-a2 - AC#4 读数：winsec 双平台成对＋两仪器安静档 3 发复跑

纪律：跑前本机移除 `GOTMPDIR`/`TMPDIR`（现跑确认＝空串）；容器档两变量**不设**、⛔ bind 宿主卷；`-count=1 -timeout 30m`；⛔ `./cmd/wisp`、⛔ `go test ./...`。

## ① 全包成对（同一 `-run` 名册＝全包，无 -run 过滤）

| 平台 | 台面 | rc | 顶格四数 | 包结果行 | 两枚仪器 |
|---|---|---|---|---|---|
| Windows 原生 | 共享工作树，go1.27.1 windows/amd64 | 0 | RUN=129 PASS=78 FAIL=0 SKIP=0 | `ok github.com/CarlosShao/wisp/internal/winsec 19.406s` | **0 命中**＝`//go:build !windows` 两件在本档⛔编译（结构性缺席，⛔ 读成"没跑"也⛔ 读成绿） |
| Linux 容器 | golang:1.27 linux/amd64，export 树 | 0 | RUN=76 PASS=48 FAIL=0 SKIP=0 | `ok github.com/CarlosShao/wisp/internal/winsec 3.038s` | 在册且全绿（读数见下） |

⛔ 两平台四数**混用/相减**——roster 结构不同（Windows 侧多 windows-only 件、Linux 侧多 !windows 件），只逐枚点名。

## ② 四数逐枚点名（全包那一发里，两枚弃权仪器各自）

**甲 `TestPOSIX120SealWindowLeafSwapIsClosed`（Linux 容器档；四数里适用者＝attempts/refused-at-act/sealed-own）：**

```
AC#1 total shape=leaf-swap: attempts=2000 hit-foreign=0 of-which-reported-nil=0 refused-at-check=948 of-which-refused-at-act=10 sealed-own-object=916 other=136
--- PASS
```

**乙 `TestAC1POSIXOpen…322r1`（Linux 容器档；适用者＝attempts/flips/leafLegRefusals）：**

```
AC#1 paired reading on one screen: attempts=3000 flips=1272 flipErrs=0 damagedWitnesses=0 leafLegRefusals=378 foreignsNarrowed=0 foreignsHoldingPayload=0
--- PASS
```

**Windows 档**：上述两枚计数字段在本档**不存在**（⛔ 编译⇒既⛔绿也⛔弃权，普查件 (i) 同向）——Windows 侧的成对读数只有包级四数（上表）。
Windows 日志内两枚仪器名命中＝0（尺＝`Select-String 'POSIX120SealWindowLeafSwapIsClosed|SwappedLeafLink'`，0）。

## ③ 两枚仪器各自 3 发独立复跑（Linux 容器、安静档、串行、⛔ 与整包并发）

**甲（3/3 PASS，窗口全部到达＝零弃权）**

| 发 | attempts | refused-at-check | refused-at-act | sealed-own | other | 判 |
|---|---|---|---|---|---|---|
| 1 | 2000 | 971 | 18 | 912 | 117 | REACHED |
| 2 | 2000 | 1067 | 13 | 822 | 111 | REACHED |
| 3 | 2000 | 1071 | 11 | 791 | 138 | REACHED |

**乙（3/3 PASS，窗口全部到达＝零弃权）**

| 发 | attempts | flips | flipErrs | leafLegRefusals | damaged | 判 |
|---|---|---|---|---|---|---|
| 1 | 3000 | 1334 | 0 | 271 | 0 | REACHED |
| 2 | 3000 | 1192 | 0 | 237 | 0 | REACHED |
| 3 | 3000 | 1268 | 0 | 359 | 0 | REACHED |

红⇄绿⇄弃权散布（安静档）＝**3 绿 0 红 0 弃权（每枚）**；散布不在安静档、在 stress 档出现——`-cpu=1` 五发＝乙 3 弃权红/2 绿（01 件逐发）、甲 3/3 弃权红（条件另带 ATTEMPTS=1）。两档逐发分日志、⛔ 混读。

## ④ 红名册作差

- 两发全包（Windows⇄Linux）顶格 `--- FAIL` ＝ **0 枚** ⇒ 红名册＝空，无"逐名归位在册票"可作。
- 本腿唯一出现过的红＝AC#2 stress 档 pristine 弃权红（甲 3 发、乙 3 发）＝**票 346 票面点名的仪器形状本身**，非新账、非产码红（同条件变异体同读数转绿＝02 件直证）。
