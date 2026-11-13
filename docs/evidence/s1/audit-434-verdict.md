# 434 收口核查（只读腿 · 裁决者 ≠ 实现者）

- 核查人：只读核查腿（非本票实现者、非本票收件编排者）
- 核查时刻：2026-11-13 10:39–10:42 +08
- 票面：`.scratch/wisp/issues/434-the-198-first-run-rule-gained-a-second-caller-so-runconfigset-now-creates-the-config-too-done.md`（`-done.md` 结尾，确认）
- 核查方式：**不信票面自述，全部现跑**

---

```
verdict: HONEST

commit: 9661bc6f4c4b52f3575a4252a1e37defba9d625e 存在（git cat-file -t ⇒ commit；且是 HEAD 的祖先，
  merge-base --is-ancestor rc=0）。改了 3 件（git show --stat / --name-only）：
    1) .scratch/wisp/issues/434-the-198-first-run-rule-gained-a-second-caller-so-runconfigset-now-creates-the-config-too.md  (+16)
    2) cmd/wisp/firstrun_198_test.go                                                                                            (+40 −5)
    3) docs/evidence/s1/434-firstrun-198-second-caller.md                                                                        (+248)
  合计 299 insertions(+), 5 deletions(-)。产物码（cmd/wisp/run.go、config_cmd_65.go）零改动，与票面自述一致。
  ⚠ 另核了票面点名的历史提交 ef00174ac（2026-11-07，票 65）：cat-file ⇒ commit，在；
  --name-only ⇒ 6 件：65-ball-glass-quality-rework.md、cmd/wisp/config_cmd_65.go(+236)、cmd/wisp/config_cmd_65_test.go(+215)、
  cmd/wisp/leg_dispatch_gate_133_test.go、cmd/wisp/main.go、docs/evidence/s1/65-ball-style-cli.md。它不在本票写面内，只是加那一跳的来源。

证据件: 在（docs/evidence/s1/434-firstrun-198-second-caller.md，248 行 / 19335B，盘上现读）。
  内含测试名 TestTicket198FirstRunCallerIsTheRunEntryOnly（§0.1 逐字红句、§5.1 PASS 读数），
  另具名 TestTicket198R2J1TheAssemblyRootStillCreatesNothing（行为半独立仪器）。

AC 框: 勾了 5 个（AC#0 只读取证 / AC#1 修 / AC#2 越界 / AC#3 门禁 / AC#4 与票 174 的关系，全为 [x]），
  没勾 0 个（grep '^- \[ \]' 零命中）。

复跑读数: 逐字行（命令：D:\work\soft\Git\bin\bash.exe -c "cd '/d/work/workspace/projects plans/Wisp' &&
  scripts/wisp-cli-staged-tests.sh '-test.run' '^TestTicket198FirstRunCallerIsTheRunEntryOnly$'"，rc=0）：
    === RUN   TestTicket198FirstRunCallerIsTheRunEntryOnly
    --- PASS: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.60s)
    PASS
  附带整族复跑（-test.run '^TestTicket198'，rc=0）：11 枚全 PASS，无一 FAIL；该族内本枚读到
    --- PASS: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.35s)

越界: 没有。git show 9661bc6f4 --name-only --format='' 三枚：票面 log、cmd/wisp/firstrun_198_test.go、
  docs/evidence/s1/434-firstrun-198-second-caller.md。docs/PLAN.md / docs/specs/ /
  internal/observe/thresholds.go / frontend/ / design/ 五项零命中（ef00174ac 那枚也零命中）。

不符: 票面与证据件自述"^TestTicket198 整族 10/10 PASS"，我现跑同一正则实测 11 枚全 PASS——只多数了一枚
  （TestTicket198AC3CreatedFileLandsPrivate），读数方向一致、不影响判定，属计数笔误而非隐瞒。
```

---

## 现跑凭据（逐条可复跑）

| 步 | 命令 | 读数 |
|---|---|---|
| 1 | `ls .scratch/wisp/issues/ \| grep '^434-'` | `434-the-198-first-run-rule-gained-a-second-caller-so-runconfigset-now-creates-the-config-too-done.md`（`-done.md`） |
| 2a | `git cat-file -t 9661bc6f4` | `commit` |
| 2b | `git cat-file -t ef00174ac` | `commit` |
| 2c | `git show 9661bc6f4 --stat` | 3 files, +299 −5（清单如上） |
| 2d | `git merge-base --is-ancestor 9661bc6f4 HEAD` | rc=0（提交真在 dev 历史里，不是游离件） |
| 3 | 读 `docs/evidence/s1/434-firstrun-198-second-caller.md` | 248 行，含 `TestTicket198FirstRunCallerIsTheRunEntryOnly` |
| 4 | `bash.exe -c "... scripts/wisp-cli-staged-tests.sh '-test.run' '^TestTicket198FirstRunCallerIsTheRunEntryOnly$'"` | `--- PASS: ... (0.60s)` / `PASS`，rc=0 |
| 4b | 同上，`'-test.run' '^TestTicket198'` | 11 枚全 PASS，rc=0 |
| 6 | `git show 9661bc6f4 --name-only --format=''` | 3 枚，禁止面零命中 |

## 补充现核（超出步骤、用于交叉验证）

- 盘上 `cmd/wisp/firstrun_198_test.go:310` 现读＝`if want := []string{"runConfigSet65", "runTextTask"}; !slices.Equal(callers, want)` ⇒ 断言确为"精确集合"形状，不是"不含常驻腿就算过"的谓词（修法形状与票面自述一致）。
- `git status --porcelain` 对该测试件与证据件**零输出** ⇒ 盘面干净、读数来自已提交态，无未提交改动撑绿。
- 票 198 的 AC J1 原文（`198-…-done.md:13`）逐字为"首建只挂在 `wisp run` 这一条「用户明确发起」的入口；⛔ 不许进 `assembleRuntime`；⛔ 不许改那两枚钉" ⇒ 判定乙（断言该承认第二个用户发起入口）与原文可对得上；禁止项 `assembleRuntime` 未被放宽（`tree-assemble` 变异仍红，见证据件 §3.3）。

## 撤销口令

「434 回甲」——把 `cmd/wisp/firstrun_198_test.go` 的 `want` 改回 `[]string{"runTextTask"}`，并删掉 `config_cmd_65.go:173-176` 那四行。⚠ 执行前先读证据件 §1 与 §5.1：甲会红票 65 自己的钉子 `TestConfigLegGetSetGuarded65`。
