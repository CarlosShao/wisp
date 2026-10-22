# 232-r1 落地腿证据件（票 232 AC#4）

角色＝写腿（实现者）。只 commit、未 push。产码零字节改动。

## 0 起手三把尺（现跑）

| 尺 | 读数 |
|---|---|
| `git status --porcelain -- cmd/wisp` | 0 行（cmd/wisp 无别人的未提交改动） |
| 起手 HEAD | `12db9bf4`（编排者派单存档：232-r1） |
| `md5`/`sha256` `cmd/wisp/config_reload.go` | `72762bd5…96c602c`，工作树与 `git show HEAD:` 逐字相同 |
| `TestTicket223` 族基线（`-v`） | 10 枚逐名全 PASS，rc=0（`logs/01-baseline-223-v.txt`） |

族内基线首发曾出现 `TestTicket223HandEditedFsLooseningCostsAnL2Card` FAIL 6.65s（红在 `:403`/`:406` 卡片未渲染），
同一份字节复跑全 PASS => 带载型计时红，非持久红，与票 223 r2 AC#5 记的那族同形。

## 1 种了什么（逐字）

票面 AC#4 逐字：**种一发"启动横幅里先含同样字样、之后才种改动"的文件／参数，断用例不许通过**。

- **参数**（不是产品文案）：`bannerSeed232Env = "WISP_232R1_BANNER_SEED"`。
- **写入点**：`newReloadRun223` 里 `seedBanner232(r.h.out)`，位置在 `live()` 装配运行栈**之前**，
  因此早于用例取 `mark`（`:621` 附近）与 `prePlant`（`:641` 附近）。
- **种子正文逐字**：
  `wisp run: 232-r1 seeded start-up banner: 这些段的改动本次运行不会生效（启动横幅先带字样，之后才种改动）`
- **关掉时**：`seedBanner232` 首行 `if os.Getenv(bannerSeed232Env) != "1" { return }`
  => 普通运行零字节、零断言变化。

为什么不改产品文案：票面禁区逐字写着"不改产品文案来让测试变绿（那正是本票要防的'文案瘦成空壳'）"。
本腿走的是反方向——**只有那几个字必须不够**。

## 2 改前读数（先量后写，`go test -overlay`，工作树当时未落地）

| 尺 | 形 | 读数 | 件 |
|---|---|---|---|
| 票 272 修前那版尺（blob `ab6c84fb`＝`8d30a862~1`） | M：横幅先带全字样，之后才种改动 | **PASS 4.80s**，rc=0 | `logs/22-M-on-true-pre272-ruler.txt` |
| 今天 HEAD 的尺 | 同一 M 形 | **FAIL 3.99s**，红在 `:644`（arm1） | `logs/20-M-on-HEAD-ruler.txt` |

=> **改前的真相有两条，必须分开写**：
1. 票 232 立案时那枚缺口（272 修前）**确实存在**，M 形绿着（4.80s）。
2. 票 272 已把**牙齿**装上（HEAD 上 M 形红在 arm1），但**仓里没有任何用例种下这一形**——它只活在
   `.scratch/wisp/probes/232|272/**` 日志与仓外 overlay 里。AC#4 要补的正是这一环。

## 3 改后读数（写清是哪一种）

**是"转红"，不是"变慢"。**

- 带种子那一发：子进程 `rc=1`，红句是 arm1 的
  `a copy nobody planted cannot be evidence that this run will not use the new values`。
- 用例显式**排除**"变慢"这条读法：若子进程红在
  `never carried … AFTER the plant within`（40s await 期限）=> 直接 `t.Fatalf`
  判"用例只是变慢了，不是本票要它察觉的那件事"。
- 所以本格的读法＝**转红**；"起点窗口化之后自然会红或变慢"这句话里，本腿落在**红**这一支。

## 4 新用例逐名 + 为什么今天响

| 用例 | 为什么今天响 |
|---|---|
| `TestTicket232R1ASeededStartupBannerCannotCarryTheRestartTier` | 用本二进制的真子进程跑 `TestTicket223RestartTierSaysItWillNotApply` 两次：(1) 种子关 => 必须 PASS；(2) 种子开 => 必须 FAIL 且红在 arm1。两个读数缺一个，种子就是装饰。 |

子进程面：`exec.Command(os.Executable(), "-test.run=^TestTicket223RestartTierSaysItWillNotApply$", …)`
——一个用例不能在自身进程里起另一个用例，而这里要问的恰恰是**那一枚**用例怎么答。

## 5 反向判据（摘掉新加的那一跳 => 必须红）

把 `newReloadRun223` 里**新加的那一行** `seedBanner232(r.h.out)` 整行摘掉
（`go test -overlay`，工作树未动）：

```
=== RUN   TestTicket232R1ASeededStartupBannerCannotCarryTheRestartTier
    config_reload_223_test.go:946: a start-up banner carrying the restart notice's own words
        BEFORE the plant did not move this case: it exited 0. The seed is back to being
        decoration - the shape ticket 232 AC#4 names is passing again. Output:
--- FAIL: TestTicket232R1ASeededStartupBannerCannotCarryTheRestartTier (9.96s)
FAIL	github.com/CarlosShao/wisp/cmd/wisp	10.179s
FAIL
rc=1
```

## 6 三读数逐字

| 读数 | 结果 | 件 |
|---|---|---|
| 变异红 | `--- FAIL: TestTicket232R1ASeededStartupBannerCannotCarryTheRestartTier (9.96s)` rc=1，红在 `:946` | `logs/31-mut-seed-jump-removed.txt` |
| 还原绿 | `--- PASS: TestTicket232R1ASeededStartupBannerCannotCarryTheRestartTier (10.49s)` rc=0 | `logs/33-asis-green.txt` |
| 原样绿 | `--- PASS: TestTicket223…` 族 10 枚逐名 PASS，`ok … 56.581s` rc=0 | `logs/34-asis-223-family.txt` |

三处逐字相同，均 `9a428907e2bb4ff8abce488bbf702ee3ac0f8147344d661fab8f5517478baa8e`、CR=0、len=45462：
(1) 提交版 `git cat-file blob HEAD:cmd/wisp/config_reload_223_test.go`
(2) 还原台＝工作树；(3) `git archive` 到仓外；另一发 `git show HEAD:` 直出同值。

## 7 `TestTicket223*` 两枚仍绿（逐名抄回）

- `--- PASS: TestTicket223RestartTierSaysItWillNotApply (5.38s)`（另一发 4.92s）
- `--- PASS: TestTicket223R2FailureSentenceRouting (0.73s)`（另一发 0.68s，9 子例全 PASS）
族内 10 枚全 PASS，rc=0（`logs/34-asis-223-family.txt`）。

## 8 零放宽

- 一条既有断言未删、未改成"或满足任一即可"。
- `restartSentence272` 与 `TestTicket223RestartTierSaysItWillNotApply` 内既有断言**行号未动**
  （`:642` 之前的一切原样；新字节只加在 `:47` 前后、`:86`、文件尾）。
- 种子只在**测试件**的宿主流上生效，产码零字节改动。

## 9 门禁 rc（本腿现跑）

| 门 | rc | 备注 |
|---|---|---|
| `go build ./...` | 0 | |
| `go vet ./cmd/wisp/` | 0 | |
| `sh scripts/d22scan.sh` | 0 | clean；ban #8 cmd/=146 Go 文件在射程内 |
| `sh scripts/check-path-length-budget.sh` | 0 | VERDICT GREEN；tracked paths=17547，over-budget=61 全在名册 |
| 格式门禁**甲尺**（判据）：`git archive` 到仓外跑 `gofumpt -l` | 0 | 空输出＝干净 |
| 格式门禁**乙尺**（诊断）：工作树 `gofumpt -l` | 0 | 空输出＝干净 |

两把尺都要报（判据＝tracked blob）。乙尺首量曾把载体列为待重排（文件尾多一空行），修正后两尺皆空。

## 10 整包终态（按名作差，逐发登记死点）

`go test -count=1 -v -timeout 30m ./cmd/wisp/`：**rc=1**，`FAIL … 322.668s`。
PASS 行 228、FAIL 行 1、SKIP 行 0（含子例计数）。

| 发 | 红名册（逐名） | RUN 无判决行（死点，逐发登记） |
|---|---|---|
| 改前（`12db9bf4`） | `TestAC4FocusReturnToPriorWindowGap33r5` | `TestAC14GoSideEvalPushReachesThePage` |
| 改后（`7c5074f9`） | `TestAC14GoSideEvalPushReachesThePage` | `TestBookedSnapshotReachesHostBitForBit329` |

**按名作差 = 0 枚新增红**：改前那枚红（`TestAC4Focus…33r5`）改后 PASS 0.82s；
改后那枚红（`TestAC14…`）**在改前就是死点**（RUN 无判决行），不是本腿引入。
`TestBookedSnapshot…329` 改后是死点，单独跑 **PASS 3.40s rc=0**（`logs/42`）=> 带载型。

死点确实**随会话漂**（派单已预告："每发恰 1 枚 RUN 无判决行、死点随会话漂"）：
改前死 `TestAC14`、改后死 `TestBookedSnapshot…329`，且 `TestAC14` 在这发里
从"死"变成"真红 0.57s"。**单独跑也红 1.34s rc=1**（`logs/41`）=> 它是既有红，
不在本腿写面上。

族内两枚必绿（见 sec 7）在**整包这一发里**同样 PASS，与单跑一致。

## 11 提交与越界自检（正控＝同尺打一笔含 `.go` 的真提交）

- 提交号：`7c5074f91a986737d047adbaa116d21f455e7bd2`。
- `git show --name-only` => 恰 1 枚：`cmd/wisp/config_reload_223_test.go`（`.go` 计 1）。
- 越界扫描（按上面同一把尺，禁面逐条前缀比对）=> **命中 0**。
- **正控**：同一把尺打 `294d4741` => 命中 **1**（`scripts/check-path-length-budget.sh`）
  => 尺是活的，0 不是"尺没跑"。
- 本笔新增 143 行 / 删除 0 行，只动测试面；产码 `config_reload.go`
  sha256 `72762bd5…` 与 `git show HEAD:` 逐字相同（见 sec 0 起手同值）。

## 12 未定义即停的两条，逐条核过，均未触发

1. "种一发"是否要动**产码**？**不需要**：种子落在测试件宿主流（`newReloadRun223` 的 `r.h.out`）上。
2. `AC#4` 与 `AC#5`"零放宽"是否**冲突**？**不冲突**：既有断言一枚未删、未放宽，新字节全是新增。

## 13 一处据实登记（编排者需知道）

票面 `## 编排者翻勾记录` 已把 AC#4 的"不许 **0ms** 通过"**就地更正为"不许假绿"**（`:45` 逐字：
"23 发读数里从不出现 0ms（最快 2.29s）=> 那三个字是我造的一个不存在的读数形状，本格按'不许假绿'读"），
并把残余缺口具名立成**票 272**（已 `-done`，落地 arm1＝种改动之前零份）。派单里"不许 0ms"若按字面读，
本仓**不存在**能被量到的读数形状 => 本腿按**编排者自己的更正**判读为"不许假绿"，交出**转红**这一支。
若编排者认为本格仍应按字面"0ms"判，则该判据需先被重写（本腿不自行假设）。
