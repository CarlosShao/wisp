# 票 316 · AC#1 三选一裁定（腿 `316-r2` · 2026-11-07）

- 腿＝写腿（编排者派单，四步拆步落盘）；锚点＝`git rev-parse --short HEAD` ＝ **`ea602a61b`**（2026-11-07 15:5x 现读）。
- 产码 md5 `cmd/wisp/config_reload.go` ＝ `0934ee7a88cd1c6e7cbd3dd0ce123971`；尺 `cmd/wisp/config_reload_223_test.go` ＝ `11d42757a64f167dbd60d8d44b9128a6`。
- ★先交代身份事实：票面 Status 已是 **withdrawn**（编排者 2026-10-14 13:38 撤，判与票 272 重复，台账 `A924`；实现腿 `316-r1` 全套读数在 `.scratch/wisp/probes/316/r1/`，交件 `8cc028f8`）。本腿按派单四步走，**独立现跑重核撤票的事实基础今天是否仍成立**；⛔ 不翻 Status、不翻 AC 勾。
- 起手三核③：`git status --porcelain -- cmd internal` 有**他腿在飞未跟踪件 3 枚**（`cmd/wisp/panel_resident_323_windows_test.go`、`cmd/wisp/panel_resident_assembly_323_test.go`、`internal/plugin/manifestslot/command_138_cancel_test.go`）——本腿零碰；家族测试编译包含它们且全绿 ⇒ 本件全部红绿读数非在飞件造成。

## 前半（步1）：三候选今天各自的真实形状（全部本腿现跑／现读）

### (a) `strings.Count(win, 该句)==1`（窗口内恰一次）—— **已在盘（重复落地＝零增量）**

- 载体：`cmd/wisp/config_reload_223_test.go:748`：`if win := r.stdoutSince232(t, mark); strings.Count(win, restartSentence272) != 1 {`（**arm2**）。
- 来路：票 272 落地（`8d30a862`，2026-10-07 09:33）。
- 对 M 形的牙：**零**——M 的窗与正常一发逐字节相同（横幅在 mark 之前、窗外），窗内怎么数都是 1（r1 引 272-r2 甲-alone 尖读数：M PASS 2.15s）。它保留的牙是对"窗内 ≥2 份"。

### (b) 标记取在 `state=armed` 横幅之后 —— **空操作**

- `mark` 现址 `:654`（`mark := r.h.out.String()`），取在测试体内；测试体跑在 `rtHook`（assembleRuntime 返回）之后；启动横幅＋`state=armed` 都由 `startConfigReload`（`config_reload.go:118`；调用点 `run.go:1149`）**同步**写进 syncWriter ⇒ mark 本来就在全部启动输出（含 armed 行）之后，再往后挪＝一句可挪出去的字节都没有。
- r1 引 272-r2 尖读数：乙在 M PASS 2.88s / cur PASS 2.15s ＝ **既不变严也不变松**。

### (c) 要求该句偏移晚于 `state=applied` —— **结构性误伤（落下去必假红）**

- 今日产码次序实测（`cmd/wisp/config_reload.go`）：`rt.mgr.OnRestartPending = rt.reportRestartPending`（`:128`）→ reload 轮询 `CheckAndReload()` 内经 OnRestartPending 发操作员句（`reportRestartPending` `:346`）→ `rt.reportReload(rep)`（`:193`）→ `auditf("config: HOT-RELOAD state=applied …")`（`:211`）。
- ⇒ **操作员句天然早于 applied 行**；"要求晚于"＝正常一发确定性 FAIL（r1 引 272-r2 丙：cur FAIL 42.20s 超时）＝把真绿判红。

### M 形今天现跑（本腿以 HEAD 为 base 重建突变；全走 overlay，共享树零编辑）

- 手法：`probes/316/r2/config_reload.M.go` ＝ HEAD 产码副本＋在启动横幅（`config_reload.go:148` 的 Fprintf）后追加一行含 `本次运行不会生效` 的横幅句（注入文本逐字承 `232-v1`/`316-r1`）；overlay 件 `probes/316/r2/prod-M.json`。
- 读数：`go test -count=1 -run TestTicket223RestartTierSaysItWillNotApply -overlay .scratch/wisp/probes/316/r2/prod-M.json ./cmd/wisp` → **rc=1，`--- FAIL: TestTicket223RestartTierSaysItWillNotApply (2.85s)`，红＝`config_reload_223_test.go:676`（arm1）**，红句 dump 第三条＝M 注入的横幅同名句逐字（log：`probes/316/r2/logs/M-rerun-ea602a61b.txt`，窗 15:59:0x→15:59:22）。
- 正控 cur（无突变）同发：**PASS，ok 3.928s，rc=0**（`probes/316/r2/logs/cur-rerun-ea602a61b.txt`）。
- ⇒ **判据①（复造 M ⇒ 必须红，红句逮住横幅同名句）今天在 HEAD `ea602a61b` 上现跑成立；让 M 红的是 arm1（票 272 的枚举外臂），不在本票三候选菜单内。**

### 行号漂移备注（三核①）

- 票面引 arm1 `:595/:596`、arm2 `:668` 是 10-07 那版尺的行号；其后 351-r2（`7c5074f91`）在同文件 +48 行 ⇒ 今日实读：arm1＝`:675-:680`（红句在 `:676`）、arm2＝`:748`、mark＝`:654`。**串都在，行号过期**——引用以今日 sed/grep 现读为准。

### 探针目录来历声明（三核②补记）

- `probes/316/r2/` 里 **8 枚件不是本腿建的**：`logs/T1-family.txt`、`T2-seedM.txt`、`T3-pin232.txt`、`T4-overlayM.txt`、`T6-d22scan.txt`＋`mut/{config_reload.go,insert.awk,overlay.json}`，时间戳 2026-11-07 14:51–14:54，属**上一轮返回 null 的腿的遗留半成品**。本腿按「临时件只建不删」原样保留并随步1 commit 一并入库，但其读数**本腿未采信、全部关键读数自己重跑**（15:59–16:05）。仅作旁证登记：其 `T4-overlayM.txt` 也量到 M 红（`FAIL 3.225s`，红句同一 arm1 形状、横幅注入措辞略异）——与前腿独立互证「M 今天已红」。

（后半部分：步2 裁定、步3 落地结果、步4 门禁，见下方追加节。）

## 步2 裁定（本腿拍；铁律 ⓪⓪⓪：判据＝最有价值＋主流 harness，⛔ 不是动作最小）

**裁定：三支全不落地。该缺口的闭合实体是票 272 的 arm1（枚举外臂），已在盘、今天现跑有牙（步1 节 M FAIL 2.85s 红 `:676`）。** 逐支具体代价：

- **不选 (a)「窗口内恰一次」**：已在盘（`:748` arm2）。重复落地代价＝同一断言两份、维护面翻倍而**新增牙为零**（对 M 零牙——横幅在窗外，窗内怎么数都是 1）；且再动这个文件会跟 351-r2 刚落的函数体搅在一起。主流 harness 对"已在盘的断言"的做法是复证咬合（mutation 重放），不是复制一份。
- **不选 (b)「标记取在 armed 之后」**：空操作。`mark` 与 armed 行出自**同一个同步写手** `startConfigReload`（写完才返回），而 mark 取在 rtHook 之后的测试体里——两者之间不存在任何可再排他的输出字节。落地代价＝**伪修复**：给台账留一条"加严了"的假账，后人排障会去信一个什么都不约束的改动；这是负价值（污染证据链），比不作为更差。
- **不选 (c)「该句偏移晚于 applied」**：直接负价值。产码真实次序（`:346` 操作员句 → `:193` reportReload → `:211` applied）是句子**早于** applied；落地即把正常一发判红（r1 引 272-r2 丙 cur FAIL 42.20s 超时形），家族测试必红，唯一"变绿"的路是放宽产码次序或放宽断言＝违判据④（不许为绿放宽）。
- **选什么**：**维持 arm1 为该缺口的收口**（272 已交付，`8d30a862`；本腿 M 突变现跑 FAIL 2.85s 复证其牙；r1 六发读数 M/MK/SA 全红、cur/MIRROR 真绿不误伤）。这是"最有价值"解：三形全堵、零新代码、零断言放宽。
- **残留（如实报，不在本票三选一菜单内）**：`272-v2` 验收腿（10-14 14:2x，件 `docs/evidence/s1/ticket272-272-v2-accept-v1.md`）量到"**单份 post-mark 欢迎语仍假绿**"（arm1 只堵种前份、arm2 只堵 ≥2 份；其 AC#7 mutant PASS 4.22s）。其钉死被 `272-v2` 判为"需产码侧因果锚点＝禁区＋人工批准"——机理：同一 run 的后续 reload 合法地会在首条 applied 之后再打操作员句，测试侧全局计数/次序断言会把合法流判红；要区分"没人种的句"与"下一轮合法的句"只能靠产码在句里带因果锚（属于产码接缝＋契约面）。⇒ **按未定义即停，本腿不自选落**，留编排者/机主裁。

## 步3 落地结果（如实）

- **零行新断言落地**。理由＝步2 裁定：三候选无可落——(a) 已在盘（重复）、(b) 空操作（伪修复）、(c) 负价值（判红真绿）；唯一同时满足「① 复造 M ⇒ 红 ∧ ② MK/SA/正例不回归」的形状是 arm1（票 272 枚举外臂），**已在盘**。
- **牙证（突变后必须红）＝本腿现跑，非引旧读数**：M 突变（HEAD `ea602a61b` 为 base，overlay 重建）→ **FAIL 2.85s，红 `:676` arm1，红句 dump 第三条＝横幅同名句逐字**（步1 节＋`probes/316/r2/logs/M-rerun-ea602a61b.txt`）；正控 cur PASS 3.928s rc=0；家族不回归见步4。
- 写面终态：`cmd`／`internal` 已跟踪件 **0 行改动**（本腿全部突变走 overlay，产码与测试尺零字节触碰）。

## 步4 门禁（本腿亲跑，2026-11-07 15:59–16:05；log 在 `probes/316/r2/logs/`）

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/config_reload.go cmd/wisp/config_reload_223_test.go`（绝对路径，派单要求） | **空，rc=0**（`logs/gates-build-vet-gofumpt-ea602a61b.txt`） |
| build | `go build ./cmd/wisp/` | **rc=0** |
| vet | `go vet ./cmd/wisp/` | **rc=0** |
| 家族最小集 | `go test -count=1 -run 'TestTicket223' -v ./cmd/wisp` | **rc=0，`--- PASS`×10 / `--- FAIL`×0，ok 32.666s**（`logs/family-TestTicket223-ea602a61b-full.txt`；与 r1 的 PASS=10 同值） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0，clean**；正控 runtests.sh 先行（含 `TestDeferredRegistryIsBidirectional` PASS）；`cmd/=226`、`internal/=941` 全 examined（`logs/d22scan-ea602a61b.txt`） |
| 写面终态 | `git status --porcelain -- cmd internal`（已跟踪件） | **0 行**；在飞未跟踪件 3 枚系他腿（见文首） |

- 环境：每发测试前 `unset GOTMPDIR TMPDIR TEMP TMP`＋`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`（逐发执行）。
- 编码：本件与 `probes/316/r2/**` 全部 LF（`sed -n "/\r/p" | wc -l`＝0）。

## 没做到的（如实）

1. **未落任何新断言／未加用例**——派单步3 的字面动作（"落地你选的那条＋加用例"）没有对应物可落：三候选现读分别为已在盘／空操作／判红真绿；落任何一个都违判据②③④之一。本腿以"现跑复证 arm1 咬合"替代（M overlay FAIL 2.85s）。
2. **未跑整包**（314 死区在册；派单只要求相关包 `./cmd/wisp`）。
3. **票面 Status 与 AC 勾零触碰**（⛔ 纪律）；撤票事实基础经本腿独立现跑**今天仍成立**——编排者 2026-10-14 的撤票裁定（与 272 重复）在 HEAD `ea602a61b` 上站得住。
4. 前腿遗留 8 枚件随步1 commit 一并入库（pathspec 目录级收进，非本腿所建），来历与处置已声明（三核②补记）。



