# 票 314 AC#3 — CI 暴露与推送前门（证据件）

- **腿**：314-ac3（实现腿）
- **时间**：2026-11-11 21:1x–21:4x +08（读数逐条带钟点）
- **树**：分支 `dev`，HEAD `14b73de13`（A1716 wave3 派发提交；含 `68275983e` 的 314 修复与 `6232efeb6` 的 302 回退）
- **原始日志根**：`D:\tmp\wisp\20261111\314ac3\`（只建不删）
- **票面逐字**（`.scratch/wisp/issues/314-cmdwisp-suite-process-death-in-window-zone.md:33`）：
  > **AC#3 CI 暴露与推送前门**：现读 `ci.yml:601-621` 一带确认哪一跳跑它＋**推送前**按本票形状重跑并把颜色写进推送前检查（⛔ 没跑过不许推）。⚠ 首推会让 CI 见生死——先在本机把死法定性再推。

---

## 1. ci.yml 那一跳的逐字读数（★票面行号已腐烂）

**票面说的 `ci.yml:601-621` 现在不是那一跳**：现 HEAD 上那段是 lint job 的注释块（"Attribution instrument self-test" 上方，`:601-610`）加 `test-core` 的注释头（`:612-633`）。**真正跑 cmd/wisp 测试的跳在别处**，grep 现读：

```
grep -n 'cmd/wisp|go test|wisp-cli-tests' .github/workflows/ci.yml
→ :888 step 名  "cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)"
→ :908 run:    bash scripts/wisp-cli-tests.sh
```

### 1.1 哪一跳

- **job**：`test-windows`（`ci.yml:812`），`runs-on: windows-latest`（`:813`），`timeout-minutes: 30`（`:821`，注释明说 30=实测 max 14.78min 的 2.0x，因为 cmd/wisp 的失败形状是"卡窗/控线挂死"不是慢）。
- **step**：`"cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)"`（`:888`），`shell: bash`、`if: ${{ !cancelled() }}`（`:906-907`），命令一行（`:908`）：

```yaml
run: bash scripts/wisp-cli-tests.sh
```

### 1.2 什么命令（链条三跳）

```
ci.yml:908  bash scripts/wisp-cli-tests.sh
  └─ scripts/wisp-cli-tests.sh:304   bash scripts/portable-tests.sh --scope=cli
       └─ scripts/portable-tests.sh:1007   sh tools/d22scan/runtests.sh ./cmd/wisp/ -count=1 -skip "<ledger 拼的正则>"
            └─ tools/d22scan/runtests.sh:211   go test -v -count=1 ./cmd/wisp/ -count=1 -skip "^(8枚 ledger 名)$"
```

- **实际 go 命令形状**：`go test -v -count=1 ./cmd/wisp/ -skip '^(TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|TestHelperProcess|TestLiveWasapiSmoke|TestAC297LivePacketByteFace|TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestSyncRegistryProbeLive)$'`
  （skip 名单＝`portable-tests.sh` 的 ledger 在 windows 平台的 8 枚 active 条目，本机 21:23 现读逐字，件 `rerun1-cli-form.txt`；**没有 `-timeout`、没有 `-tags`**——`go test` 全程零 `-timeout` 旗，job 级只有 GitHub 的 `timeout-minutes: 30` 兜底挂死。）
- **档**：**cli 档**（`cli_pin` 只含 `github.com/CarlosShao/wisp/cmd/wisp` 一包，`portable-tests.sh:404-406,465`）。**不是** core（ubuntu）、**不是** windows（那是不含 cmd/wisp 的 13 包档）。
- 前置：同 job 的 `actions/cache`（deps.toml 键）+ `build.ps1 -Env dev`（`:874-886`）把 sherpa DLL 铺进 `third_party/sherpa-onnx` 并进 PATH；脚本自己的 GUARD（`:112-138,251-288`）缺 DLL 时拒绝跑、报 setup 红。

### 1.3 CI 那一跳今天会不会跑到那 5 枚死亡点 —— **会，全部 5 枚都在射程内**

死亡点 5 枚（票 314 §现量＋更正后口径）：`TestAC13BringUpSurvivesAReusedThreadQuit`、`TestAC14GoSideEvalPushReachesThePage`（q1/q4 死点）、`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（r1b 死点族边界）、`TestAC14AwaitedBindingReplyReachesThePage`、`TestPanelThreadIsSTAAndExitsCleanly`（r1b 具名死点）。

- ★**302-r1 的"搬 AC14 两枚进 winlive"已被回退**：commit `6232efeb6`（2026-11-11 21:17:24 +080，"Revert 302: move the three page-reply rulers…"）⇒ 默认档今天**重新包含 AC14 两枚**。
- 现证（本机 21:26:25，`go test -list` 对默认档=**无 `-tags winlive`** 的测试二进制，rc=0，件 `list-death-points-default-tags.txt`）：5 枚名字**逐一列出**＝全部编进默认二进制；`panel_resident_windows_test.go` 头部只有 `//go:build windows`（`:1`），无 winlive 半边。
- ⇒ **CI 的 cli 档（`./cmd/wisp/`、默认 tags）与 5 枚死亡点的交集＝5/5**。CI 的 `-skip` 名单 8 枚里**没有一枚**是死亡点。也就是说：**下一次推送，`test-windows` 的这一步会在真窗区再次暴露生死**（票面 §Progress log 2026-10-15 那两发 CI 早停，run `38690995413`/`38687919508`，同一形状）。

---

## 2. 本机复跑"本票形状"（同命令形状，一发）

- **命令**（与 CI 那一跳完全同形，直走 `scripts/wisp-cli-tests.sh`，让它的 GUARD 链自己铺 DLL 与 GOTMPDIR）：

```
bash scripts/wisp-cli-tests.sh        # Git Bash D:\work\soft\Git\bin\bash.exe；PATH 前置 build\ 与仓根（DLL）与 go\bin
```

- **起跑**：2026-11-11 21:22:55 +08；GOTMPDIR＝`C:\Users\swq\AppData\Local\wisp-test-gotmp\cli-20261111-212255-1657`（ticket 310 行逐字在日志头）。
- **终局**：2026-11-11 21:30:03 +08（`FAIL github.com/CarlosShao/wisp/cmd/wisp 392.721s`），脚本 rc=**1**。
- **四数**（`portable-tests.sh` 四数行逐字）：`=== RUN=418  --- PASS=325  --- FAIL=3  --- SKIP=1`
- **★本发是"早停"，且死法就是本票的签名**（这是本次复跑的核心读数）：
  - **无 panic / 无 signal / 无 runtime error / 无 `exit status` 行**（四类逐字 grep 均 0 命中）——整进程安静消失，与票面 §现量"死法共同特征"逐字同形。
  - **死点**：最后一行 `=== RUN   Test145ArgsJSONComesFromTheBookedForensicsRow`（`cmd/wisp/panel_run_145_test.go:217`，票 145 run-packet 取证族）**无任何 verdict**，紧接 `FAIL …cmd/wisp 392.721s`——与 2026-10-15 那发 CI run `38687919508`（死在 `TestLegacySubShapeOneHookDiesInAReentryLoop` 之后、下一枚无 verdict）**同族**：死点都落在 `panel_*` 家族、都死在窗口区、都没有 verdict。
  - **终点锚没到**：`TestZZ315A1` 在 618 枚里排在**最后**，`=== RUN   TestZZ315A1` 在日志里 **0 命中** ⇒ 按 A976 口径，**本发四数不作覆盖面分母**。
  - **分母缩了多少**：`go test -list '.*'`（默认档）＝**618** 枚顶层用例（件 `list-all-top-level-default.txt`）；`-skip` 8 枚在本机全跳 ⇒ 期望分母 610；本发只 `=== RUN` 418 ⇒ **约 192 枚（含 3 枚 ZZ315A1 锚）这一发根本没跑**，四数行却在、rc 也只是普通红——现有 CI 步的判据**完全接不住**这种死法。
  - ★**5 枚死亡点全数"今天活着"**：`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`(1.08s PASS)、`TestAC13BringUpSurvivesAReusedThreadQuit`(1.03s PASS)、`TestAC14AwaitedBindingReplyReachesThePage`(1.19s PASS)、`TestAC14GoSideEvalPushReachesThePage`(1.45s PASS)、`TestPanelThreadIsSTAAndExitsCleanly`(0.92s PASS)——`68275983e`（去共享 profile）之后这 5 枚**全绿跑完**；本发的死点**换了一枚**（145/transport 族），与票面"死点只落在 2 枚测试"的旧形状不同 ⇒ 死亡**没有**被修掉，只是换了宿主（这一点与 408 落地腿 2026-11-08 的"残余共享态清单"提示一致）。
- **3 枚具名红**（都打了完整 verdict，属软形态、非本票死法）：`TestTicket198FirstRunCallerIsTheRunEntryOnly`（`firstrun_198_test.go:281`，caller census 多出 `runConfigSet65`）、`Test313WritableUnionEqualsTheTwoRosters`（`panel_config_313_test.go:252`，door union 与 roster 不等）、`TestLockedSuffixNeverTakesTheLockItself33r11`（`panel_locked_naming_33r11_windows_test.go:243/245`，`setPriorFocusLocked` 体 0 把锁）。
- **1 SKIP**：`Test423PunctuationSwitchGatesTheStage`（`config_wire_423_windows_test.go:286`，本机 System32 的旧 onnxruntime 遮蔽了 PATH 上那份——环境类、ledger 外，属本机噪声）。
- **对照读数**：与本票背景给的修后基线（`68275983e` 后 418/325/3/1、353.8s、"零 panic/signal/killed"）四数**逐字一致**——但四数一致≠跑到终态：本发现量了同一四数下的**早停形状**。⇒ 更正：**"修后整包跑到终态"的说法今天复核不成立，准确形状＝"修后 5 枚旧死点全绿跑完，但整包仍会早停在窗口区，死点从 AC13/AC14/STA 漂移到 panel_run/transport 族"**（修复只救活了那 5 枚具名死点，未让整包到锚）。
- ⚠ 读数方法注（防同形错）：四数行的 `SKIP=1` 与原始流的 2 行 `--- SKIP` 不矛盾——第 2 行是 `portable-tests.sh` 的 unaccounted-SKIP 报告**重印**同一枚（`config_wire_423_windows_test.go:286`），计数应认四数行。
- **日志**：`D:\tmp\wisp\20261111\314ac3\rerun1-cli-form.txt`（567KB，只建不删）；分母件 `list-all-top-level-default.txt`；死亡点 `-list` 件 `list-death-points-default-tags.txt`。

---

## 3. "推送前门"的落点建议（⛔ 本件只提建议，未改任何门禁形状）

**先答"为什么今天这扇门缺失"**：现在的编队推送前检查（`.scratch/wisp/probes/orch/gate-roster-304.md` §1/§2）覆盖进程名册、临时目录残渣、static scan、`internal/audio` 档——**唯独没有"cmd/wisp 的 cli 档今天跑到终态了吗"这一问**。票 314 的形状（整进程无声死亡、四数分母不可信）恰好是唯一逃过现有全套尺的死法：A976 口径（读四数前先证终点锚）在台账里，但没有一把推送前跑的尺执行它。

**建议落点（按优先序）**：

1. **首选：`scripts/portable-tests.sh` 的 cli 分支末尾加"终点锚断言"（收益最大、形状最主流）**。
   现状：cli 档跑完只有四数＋GUARD B（包级 `ok` 行在不在）。本票的死法是**整进程消失**，§2 的本机复跑已经现量出 CI 判据的洞有多大：rc=1、四数行在、GUARD B 的包级 `FAIL` 行也在——**每一条现有判据都各自满足、没有一个 catches 早停**；分母从 610 缩到 418 无人过问。主流 harness 的做法是给套件一个哨兵：要求该档**最后一个具名用例**（现成可用＝`TestZZ315A1*` 三枚，`cmd/wisp/zz315a1probe_test.go:466/525/646`，字母序天然垫底、文件头自述 "STANDING REGRESSION NAIL"）必须出现在 `-v` 输出里，否则该步红并写"early-stop、四数不作分母"——这正是 A976 口径的仪器化。改一处脚本、CI 的 `:908` 一字不动，且按本仓 GUARD 同款纪律落（写清理由、不动断言、不新增 skip）。
2. **次选：`scripts/wisp-cli-tests.sh` 的 verdict 汇总段（`:302-316`）**。它已经握着完整日志（`$cli_log`）加一个具名红形状（ticket 98 的 load-failure 判别），加一个"日志里有 `^--- FAIL` 但缺 `TestZZ315A1` 则输出 named finding 并保持 rc"的兄弟形状，改动半径最小。但它只护 CI 这一跳，护不了本机直发 `--scope=cli` 的腿。
3. **推送前人工步骤（编排者侧，立即可用、零代码）**：把"推送前跑 `bash scripts/wisp-cli-tests.sh` 一发，四数＋rc＋终点锚逐字抄进推送 note"写进编排者推送 checklist（`gate-roster-304.md` §2 的终态锚那一节加一行）。这满足票面"⛔ 没跑过不许推"的字面，但在有第 1 条之前它是人肉闸，会漏。
4. **不建议**：改 `ci.yml`（本票边界明写 ⛔ 动 `.github/**`；且 CI 侧加门治不了"推送前"，颜色永远晚一手）。

**门内容（建议的最小判据，照 A976 逐字；①③已由 §2 本机复跑各自取得正反读数）**：
- 判据①（已现量反例）：`-v` 流里 `^=== RUN TestZZ315A1` 至少 1 行（套件活到了尾部）；本发 0 行＋rc=1＋四数行在 ⇒ 正是必须变红的形状；
- 判据②：四数里 `FAIL>0` 时必须逐枚具名（禁"只有数字没有名字"的绿）；
- 判据③（已现量正例形状）：`^=== RUN` 后同枚无 `--- PASS/FAIL/SKIP` verdict、且后面紧跟 `FAIL <包> <秒>s` ⇒ 写"early-stop at <最后具名用例>，四数不作分母"，该步红；
- 判据④：red 用例的名单（§2 那 3 枚）**不进本门**——它们是断言面，归各自票；本门只管"跑完没跑完"。
- ⛔ 红线不碰：不为变绿放宽任何断言；skip 名单只许 ledger 机制改；不动阈值/golden；死亡点漂移（本发死在 145 族）⇒ 锚判据必须**只锚终点、不锚死亡点**，否则门会跟着死点漂。

---

## 4. 本件边界

- 写面：仅本件＋`D:\tmp\wisp\20261111\314ac3\`。⛔ 未动 `.go`/`.yml`/`scripts/**`/`internal/**`；⛔ 未改任何门禁形状（§3 是建议）；⛔ 未 push。
- 复跑用的命令、roster 前读数、`-list` 取证全部落在日志根，可逐字重跑。
