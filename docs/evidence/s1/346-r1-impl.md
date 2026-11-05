# 346-r1 落地件：AC#3 归因尺落盘（票 346）

- 腿＝`346-r1`（实现）；起手 HEAD `2f7f9de16`；2026-11-05 10:1x +08。
- 本件只结 **AC#3**（唯一一枚本机现刻可闭环的判据）；其余 AC 的停手理由见第四节，⛔ 伪造。
- 写面＝`.scratch/wisp/probes/346/a1/**`（票面 `AC#6` 逐字）＋本件（派单白名单允许
  `docs/evidence/s1/346-*.md`）；⛔ 动 `internal/winsec/**`／`scripts/**`／`.github/**`／票面文字／
  任何 `.go`／任何冻结件。

## 一、起手三核（⛔ 抄票面，逐枚现跑）

| 核 | 现读 | 尺 |
|---|---|---|
| (a) 两枚仪器仍在 HEAD | 都在；首行 `//go:build !windows` | `git rev-parse HEAD:<path>` |
| 两枚 blob 与票面同串 | `dfaddac8...`／`0b0e5173...`；HEAD、`a42b35af`、`a5080612` 逐枚同串 | 4 次 `git rev-parse` 对撞 |
| 更早锚 `adec1db1` 无 openleaf | `cat-file` 非 0（与票面同向） | `git cat-file -e` |
| (b) 红句 file:line | `openleaf_322_other_test.go:206` 逐字在；`posix_seal_window_120_other_test.go:233`＋`:240` 逐字在；★票面写的 **`:370` 是过期读数**（同一 blob 在两锚就是 `:233`／`:240`） | `grep -n` |
| (c) 有没有别人在飞的改 | 无活写者。`probes/346/a1/**` 是**本腿上一趟**（09:52–09:58 起手、死在写本件之前）的遗产：`D:\tmp\wisp\20261105\346r1\` 里 `baseline.sh`／两枚原始 CI log／`live.out` 逐枚在盘；`baseline.sh` 与 `346-c1-census.md` 第三节内联版 **sha256 同串**（`db1f1694...`）；10:04 现查无该腿进程 | `Get-CimInstance Win32_Process`／`sha256sum` |

★票面自己写着"行号只当**读数**、⛔ 当判据"，所以 (b) 的 `:370` 不致 STOP；但下一程引用前按本行改正。

## 二、AC#3 落地：RED-before ⇄ GREEN-after

**落件**＝`.scratch/wisp/probes/346/a1/ci-red-attribution.sh`（201 行）＋`selftest.sh`（34 检）＋
`README.md`＋`.gitattributes`＋`fixtures/**`（9 枚）。基线＝`docs/evidence/s1/346-c1-census.md` 第三节内联版
（sha256 `db1f1694...`）。

### RED-before（基线）

1. 用本件 selftest 打基线：`TOOL=<baseline> bash selftest.sh` ⇒ **8/16 FAIL、rc=1**（基线没有
   `--issues/--ledger/--func-index/--posctl` 离线缝，把 flag 当 run id 去调 `gh`，主跑即 exit 3，
   后 18 检被短路）。逐名失败：
   - `main run exit code (expected [0], got [3])`
   - `empty roster exit code is 5 (expected [5], got [3])`
   - `empty roster is refused out loud (missing [REFUSING])`
   - `missing positive control exits before attribution (expected [1], got [3])`
   - `unreadable log exits 4 (expected [4], got [3])`
   - 另 3 条＝主跑没有产出的连带断言。
2. 用环境变量喂同一批夹具（绕开 flag 缺口），两处行为缺口照样现形：
   - 空名册：基线逐字 `unique red names=0  UNREGISTERED=0`、**rc=0** ⇒ 把"没量到"报成绿；
   - 日志缺失：基线照样 `logs=1  顶格 FAIL (job x name)=0`＋`unique red names=0`、**rc=0** ⇒
     缺文件不拒答。

### GREEN-after（本件）

- `bash selftest.sh` ⇒ **PASS (34 checks)、rc=0**。
- 两枚真 run 现跑（`--log` 两枚原始 log 39101060753 ⇄ 39124122121）：
  `logs=2  顶格 FAIL (job x name)=16`、`POSITIVE CONTROL ok`、
  `unique red names=14  UNREGISTERED=0`、rc=0。票面两枚逐字：

  ```
  name=TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1
    jobs=test-core  file=openleaf_322_other_test.go
    判据在册票=348  具名文件在册票=322,346,369  散文提及=2枚票  台账=—
    => REGISTERED
  name=TestPOSIX120SealWindowLeafSwapIsClosed
    jobs=test-core  file=posix_seal_window_120_other_test.go
    判据在册票=348  具名文件在册票=120,322,342,346,369  散文提及=2枚票  台账=—
    => REGISTERED
  ```

- 反向形（selftest 内逐条 rc）：空名册 **5**／缺失 log **4**／`gh` 失败 **3**／正控不在池 **1**／
  无输入 **2**；`^--- FAIL` 顶格尺在 gh 前缀日志上 **0 命中**（正控证明尺灵）。

## 三、门禁与写面自证

- `bash -n` 三枚 shell（工具／selftest／fixture `gh`）＝ **rc 0/0/0**。
- `sh scripts/d22scan.sh` ＝ **rc=0 clean**（internal 871＋cmd 203 枚 Go 文件；本腿 0 枚 `.go` 改动）。
- go build／go test／gofumpt ＝ **N/A**（本腿 0 枚 `.go`；`git status` 逐枚核过）。
- 字节尺（probes 全部 17 枚文件）：U+00A7 ＝ **0**（README 原一处已改英文 `section 3`）、CR ＝ 0、
  零字节 ＝ 0、`.out` ＝ 0。
- 换行属性：落件前 `git check-attr` ＝ `text: auto`／`eol: unspecified` ⇒ Windows 检出
  （`core.autocrlf=true`）会把夹具 log 变 CRLF、把 `bin/gh` 的 shebang 打断；本腿落
  `.scratch/wisp/probes/346/a1/.gitattributes`（`* text eol=lf`，只罩本子树）⇒ 落件后
  `fixtures/gh-prefixed.log`／`fixtures/bin/gh`／`fixtures/funcs.tsv` 逐枚 `text: set`／`eol: lf`
  （`git check-attr` 现读）。
- ⛔ 动：`internal/winsec/**` 0 字节、`scripts/**` 0 枚、票面 0 字、台账 0 字、冻结件／golden／
  `thresholds.go`／`go.mod`／`go.sum`／`.github/**` 0 枚；⛔ push。

## 四、STOP-notes（本机现刻不能闭环的判据，逐条具名）

1. **AC#1 三形表已在**（`346-c1-census.md` 第二节：甲／乙／丙逐形给"动几枚文件、顶不顶
   322・120 的既有钉、CI 那一步分母变不变"，并排表在 `:137-143`）。**但票面逐字"⛔ 选形"**
   （`:34`）⇒ 选形⛔ 本腿，须非实现者裁；且 **甲** ＝ `.github/workflows/ci.yml` ＝ 人工批准
   （`AC#5` 逐字禁）；**乙／丙** ＝ 动两枚仪器，须对齐本票"为什么现在不派"那句
   "与票 342 的落地腿一起排写面"，而 **票 342 今日 Status 仍是"已立，未派"**（342 文件 `:3` 现读）
   ⇒ 未满足。
2. **AC#2**（去牙反形：摘弃权守卫 ⇒ 假绿）以"已选形"为前提；且两枚仪器 `//go:build !windows`，
   摘守卫那一发要在仓外导出树的 **Linux** 档跑（本机 Windows 结构性不编译）⇒ 本机不可闭环。
3. **AC#4**（Linux 容器 ⇄ 本机 Windows 成对＋各 ≥3 发）同理：Windows 档不编译两枚仪器；台账
   `A1407` 已具名"346 暂缓，其绿读数须 CI，本机不可闭环" ⇒ 本腿⛔ 伪造成对读数。
4. **AC#5／AC#6**：`winsec.go`／`winsec_other.go`／`winsec_windows.go` 0 字节；票
   113／116／119／120／322 既有钉 0 处改；框 0 枚翻、`Status` 未改、未加 `-done`；票尾 0 字
   （票面 `AC#6` 允许一条 bullet，但本腿派单白名单不含票面 ⇒ 留给编排者／裁腿）。
5. **写面差异（具名交回）**：票面 `AC#6` 写面 ＝ `probes/346/a1/**` ＋ 票尾一条 bullet；本腿派单
   白名单另含 `internal/winsec/**`／`scripts/**`／focused `_test.go`／`docs/evidence/s1/346-*.md`。
   本腿只用了**两边交集**（`probes/**`）＋派单白名单的证据件；`internal/winsec/**`／`scripts/**`
   0 枚（因 `AC#1` 未选形＋`AC#5` 禁字节）。

## 五、本件自己的射程

- 只证 AC#3 的"尺有牙、能量、能拒答"；⛔ 证"两枚仪器已稳"、⛔ 证"CI 不会再红"（那属
  AC#2／AC#4，未选形前无可证对象）。
- 工具两条已知极限照抄 `README.md`：只答"这枚红有没有票管"、⛔ 答"是不是本腿新红"；`--log`
  只见 failed-step log，步数是下界；散文提及不算归属。
