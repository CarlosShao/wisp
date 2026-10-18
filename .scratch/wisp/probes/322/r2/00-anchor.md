# 322-r2 起手件（锚＋落点合规＋台面自证）

## 0. 起手锚（现跑，⛔ 抄派单）

- `git log -1` = `c11163c896b09e7b6863d347dadbcb65702585cc`，`Sun Oct 18 10:02:56 2026 +0800`，
  主题＝「接手自检：§1 八条全跑＋差异登记…」＝与编排者 10-18 10:0x 所见同一枚。
- 起手尺复量票框：`.scratch/wisp/issues/322-…md` `- [ ]` 命中 **1** 枚（AC#4）、`- [x]` 命中 **5** 枚
  ＝与派单「1 未勾／5 勾」相符，无第三腿动票。

## 1. ★落点合规（机主 2026-10-18 10:1x 新铁律，覆盖原台面配方）

**已改落点＝`D:\tmp\wisp\20261018\322r2`**（盘符根零污染）。

- 导出树＝`D:\tmp\wisp\20261018\322r2\head.RJNBgS`（`mktemp -d` 唯一新名）；
  8.3 短名别名实测在位＝`dir /x` 列出 `HEAD~1.RJN`，`for /d %d … %~sd` 解出
  `D:\tmp\wisp\20261018\322r2\HEAD~1.RJN` ⇒ 短名族判据不需要退回盘符根。
- `GOTMPDIR`（D 档）＝`D:\tmp\wisp\20261018\322r2\gtD`，跑前 `mkdir -p` 建过。
- 本腿在 `D:\` 与 `C:\` **盘符根下零建物件**（自检尺＝`dir /x D:\` 名册里本腿名下 0 枚；
  `D:\` 根 10-17 那批 `W334R1*`／`wisp228r3b` 等**不是本腿建的，本腿⛔ 清**——「⛔ 顺手清理 D 根旧树」照办）。
- ⚠ 本腿第一版台件确实建在 `D:\tmp\probe322r2\`（`gotmpwinsecdir`／`shortnametest1234567890`／`w322fix`）＝
  **在 `D:\tmp` 下、不在盘符根下**，但仍不是新铁律具名的唯一允许落点 ⇒ 已全部作废、改用
  `D:\tmp\wisp\20261018\322r2`；旧三枚目录**只建不删**留在盘上，⛔ 被后续任何读数引用。

## 2. 台面自证（导出树字节 ⇄ 起树锚）

| 文件 | 导出树 sha256（本机 Get-FileHash） | 起树锚 `c11163c8` 的 winsec 树对象 |
|---|---|---|
| `internal/winsec/winsec.go` | `41ba9a9b8fcf7050…`（＝`322-r1b` 记录的 pristine 同串） | `8ddcb8cb64530cf024ffbc2bcf34a084cac740d9` |
| `internal/winsec/winsec_other.go` | `19f3c211055011cc…`（同上） | 同一枚树对象 |
| `internal/winsec/winsec_windows.go` | `5cdef88183deb576…`（同上） | 同一枚树对象 |
| `internal/winsec/sealhandle_privateroute_120_other_test.go` | `fb6c54ed02fb8b43…` | 同一枚树对象 |

- 漂移尺（照 `322-v1` §0 那把）：`git rev-parse c11163c8:internal/winsec` =
  `git rev-parse HEAD:internal/winsec` = **`8ddcb8cb`** 同一枚 ⇒ 别腿在飞的提交没动过 winsec，
  本腿读数对起树锚与对现 HEAD 同时有效。
- 导出树非空尺＝`GOCOUNT=760` 枚 `.go`（八件套＋docs＋手工补 `build/*.dll`、`third_party/sherpa-onnx/*.dll`、
  `frontend/dist`）；`go.mod` 头行＝`module github.com/CarlosShao/wisp`。
- 共享工作树写面尺＝`git status --porcelain -- internal cmd` 起手 **0 行**（⛔ 本腿写这两面）。

## 3. 本腿自己一笔台面故障（自抓，⛔ 计入判据）

第一版取证脚本 `host-winsec.sh` 把仓路径写成 MSYS 形 `/d/work/workspace/projects plans/Wisp`，
本机 `sh` 的 `git -C` 当场 `fatal: cannot change to …: No such file or directory`、rc=128、**零读数**
⇒ 判台面故障，件按「只建不删」留在 `host-winsec.sh`（文件头已标〔故障件〕），有效版＝`host-winsec.ps1`
（原生 PowerShell，⛔ MSYS 路径形）。同一条坑也解释了 `container.ps1` 为什么用 `D:/` 形挂载路径。
