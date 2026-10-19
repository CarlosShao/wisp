# 346-a1 件04 = blob 同一性尺 ＋ 两枚仪器落位（射程尺）

## 一、blob 同一性尺（本腿现跑 `git rev-parse <rev>:<path>`，5 枚路径 × 3 锚）
| 路径 | `a42b35af` | `a5080612` | `6f808a83` | 判 |
|---|---|---|---|---|
| `internal/winsec/openleaf_322_other_test.go` | `dfaddac8…` | `dfaddac8…` | `dfaddac8…` | **三锚同串** |
| `internal/winsec/posix_seal_window_120_other_test.go` | `0b0e5173…` | `0b0e5173…` | `0b0e5173…` | **三锚同串** |
| `internal/winsec/winsec_other.go` | `e039da10…` | `e039da10…` | `e039da10…` | 同串 |
| `internal/winsec/winsec_windows.go` | `63714a83…` | `63714a83…` | `63714a83…` | 同串 |
| `internal/winsec/winsec.go` | `67783ae0…` | `67783ae0…` | `07300bff…` | ★**在 `a5080612`→`6f808a83` 之间变过** |

★三枚 run ⇄ 锚（票 346 现量表第 2 行逐字）＝39101060753 = `a42b35af`（双绿）｜39124122121 = `a5080612`（双红）｜
39244425666 = `6f808a83`（双绿）；更早 39071959151 = `adec1db1`＝甲枚 PASS、乙枚**缺席**
（本腿现跑 `git cat-file -e adec1db1:internal/winsec/openleaf_322_other_test.go` ＝ **rc=128**）。

⇒ 本腿由此得两句（★⛔ 写成票面那句通则）：
① 两枚**仪器自身**字节⛔ 变 ⇒ "blob 一字未动却翻脸"在仪器这一层成立；
② 但 `winsec.go`（被测产码）在红发锚 ⇄ 下一个绿发锚之间**变过** ＝票 342/348 乙形落地那批 ⇒
**"翻脸⛔ 来自任何字节"⛔ 成立**（红发 `39124122121` 在落地之**前**，本腿据两锚现量＋台账时序）；
⇒ 该结论＝交编排者裁的**票面补强**（票面第 2 行那句只钉了两枚仪器＋`winsec.go`/`winsec_other.go` 两锚）。

## 二、两枚仪器的落位 ＝ 本票最要紧的一格（射程尺）
- 尺＝`^(=== RUN|--- (PASS|FAIL|SKIP)): <名>` 于五份归一日志（件00 第一节）＋本机整包日志：
  两枚名**只在 `test-core`（Linux）**现身；`test-windows`（`win.norm.log`）各 0 命中；
  本机整包 `D:\tmp\wisp\20261019\orch227-b\full.log`（1544211B，四数 2317/1529/29/8）各 **0 命中**。
- 机制（读码，⛔ 运行）＝两枚文件第 1 行都是 `//go:build !windows`
  （尺＝逐枚 `Get-Content -TotalCount 3`；★同包对照＝`winsec_windows.go`＝`//go:build windows`）。
  ⇒ 在 Windows 主机 ⇄ CI Windows 档**根本没参与编译**＝⛔ 绿、⛔ 红、⛔ 弃权。
  ⇒ 派单要的"容器⇄本机""D 卷⇄C 卷"两轴对这两枚**结构上⛔ 存在**；本腿⛔ 用"本机没红"当"本机绿"。
- 正控（⛔ 尺失明）＝同尺在 `core.norm.log` 命中；`probes/342/v1/logs/20-host-winsec-baseline.log`
  （192842B）`^--- PASS`＝61 枚；`probes/342/a1/logs/20-container-winsec-baseline-1.keyrows.txt`
  自陈容器档四数 `RUN=71 PASS_TOP=43 FAIL_TOP=0 SKIP=0`＝容器确实跑得到这两枚。
