# 票 295 × 落地腿 `295-r1` — 验收腿 `295-v1` 的三件主菜（变异在最终形状上重跑）· 10-12 16:4x +08

判语在前、尺在中、件名在后。⛔ 本腿翻任何 AC 框、⛔ 动票面原句、⛔ 加 `-done`、⛔ 改生产/测试码一字。
读数面＝`421469b7` 的最终形状（尺＝`git diff --numstat 421469b7..HEAD -- cmd/wisp`＝**0 行**，16:44 现跑；起手 HEAD `a92cd10f`，三发在 `8ef6b672`⇄`a323b149`，`cmd/wisp` 全程零差）。

## 主菜① 删写变异在**最终形状**上重跑 ＋ 整包红名册差集 ⇒ **成立**（并补上 ⓔ-2 第一次实测）

- 定向那一发：M1＝把 `toggleMute` 体内 `ra.started.Store(gate.Open())` 换成等高的 `	_ = gate.Open()`（＝票面 ⓑ 那一形"boot 写过、手势路径不再更新"）⇒ 指名用例逐字红：
  `resident_mute_290_windows_test.go:125: mute-hotkey: the leg's derived mirror says nothing is running while the gate is open`＋`--- FAIL`＋`RC=1`；还原形（pristine）同命令 ⇒ `--- PASS (0.01s)`＋`ok … 0.102s`＋`rc=0`。
- ★整包差集（这才是票面 ⓔ-2 那句话**第一次被跑**）：同一台件、同一条命令、同码面两发（`go test -count=1 -v ./cmd/wisp/ ./internal/audio/`，pristine ⇄ M1），名册各 **455** 状态行（含子测试那把尺；只数顶层那把＝`PASS 328／FAIL 3／SKIP 2` ＋子项 122 行，两把并排⛔ 相加），逐行 `diff`＝**4 行**＝只有两枚用例翻色：
  `TestAC290BothMuteGesturesTurnTheGateAndBack` **PASS→FAIL**（变异该咬的）⇄ `TestAC4FocusReturnToPriorWindowGap33r5` **FAIL→PASS**（台账 `pending-and-issues.md:14768` 具名的真窗族 flake，与变异⛔ 相关）。
  ⇒ **"其余 7 枚读点对这枚变异不敏感"＝成立**：#1/#2/#3（`TestAC247*` 三枚＋表驱动 3 子项）三发全 PASS、#8（`TestAC290Outcome…`）在 M1 下仍 PASS、`TestAC290MutedDefaultStillDecidesTheBoot` 三发全 PASS；**#4（live mic）三发全 SKIP＝仍是零读数**，腿的欠账 `Q-295-r1-live-skip` ⛔ 被本腿销掉（真机归机主，派单⛔ 本腿设 `WISP_LIVE_MIC`）。
- 判语口径按派单 §5 的更正写：⛔ 说腿"作假"——腿的 M1 红句是真的（我在最终形状上复现到**同一枚断言、同一句逐字**），缺的只有"一把差集尺"；本腿交的就是那把尺。
- 件＝`logs/roster-diff-m1-p-m2.txt`（名册／差集／逐名表）＋`logs/gates-restore-proof.txt`（命令逐字、rc、还原证明、五族〔作废〕具名）。原始三发日志（2,252／2,257／2,255 行，sha256 逐枚记在件里）留专用根 `C:\Users\swq\AppData\Local\wisp-test-gotmp\295v1-20261012-160949-1256\logs\full-{p,m1,m2}.txt`。

## 主菜② 形状变异（`Store(true)` 恒真）到底红不红 ⇒ **成立：⛔ 恒真的假钉**（本票最该被抓的那枚⛔ 存在）

- 定向发：M2 红在第二轮那一枚断言，逐字 `resident_mute_290_windows_test.go:142: mute-hotkey: capture stayed handed over after muting (open=false started=true)`＋`RC=1` ⇒ `if ra.gate.Open() || ra.started.Load() {`（`:141`）**会**因恒真而进分支，⛔ 恒真。
- 整包发（新读数，腿⛔ 量过）：M2 的红名册 **5 枚** ⇄ pristine **3 枚**＝变异翻红 **2 枚**——除上面那枚，还有 `TestAC290OutcomeIsReadOffTheGateNotOffTheRequest`，红句 `resident_mute_290_windows_test.go:175: the leg's derived mirror says capture is running over a refused device`。
  ⇒ 腿那句"实际只有 1 枚有牙"**只对"删写"方向成立**；对"恒真"方向有牙的是 **2 枚**（`#5/#6/#7` 那一枚用例＋`#8` 那一枚）。这条差集只有整包发能看见，腿那两发是 `-run` 定向形。
- ⓓ-2 `#7` 的 `t.Fatalf` 参数位（派单指名要顺手核那一问）＝**用 `.Load()`、⛔ 把 `atomic.Bool` 丢进 `%v`**：`:142` 现量逐字 `t.Fatalf("%s: capture stayed handed over after muting (open=%t started=%t)", gesture, ra.gate.Open(), ra.started.Load())`；M2 那一发的实际输出＝`open=false started=true`（裸布尔）；花括号尺＝`grep -c -e 'started={' -e 'started= {' logs/full-m2.txt`＝**0 命中**（分母 2,255 行；失败日志⛔ 撒谎）。
- 件＝同上两枚（M2 那一列）。

## 主菜③ 注释里那三句事实主张 ⇒ **三句都成立**，★但**它们在 HEAD 已不存在**（派单的前提要更正）

- ★先更正口径（⛔ 我改注释一字）：那三句逐字在 `d5ef9cb9` 的 blob 里；`1c1d7df3`（腿自己那笔"压回原行数"）把它们**删掉了 19 行**（尺＝`git show --numstat 1c1d7df3 -- <件>`＝**6/25**；三枚 blob 行数 `0615c6c9`=425 ⇄ `d5ef9cb9`=**444** ⇄ HEAD=425）。⇒ HEAD 只剩 `ui-sta writes, boot reads: 295.`（`:104`）与 `It crosses threads (this body is ui-sta, boot reads it), so it is a Store.`（`:167`）。所以本判射程＝"腿当时写下的主张真不真"，⛔ "HEAD 注释里还有没有这句话"。
- 主张①「`-race` is not in the roster」⇒ **成立**。票面原尺 `grep -rn --include=*.sh --include=*.ps1 --include=*.yml -e '-race' scripts/ .github/`＝**0 命中 rc=1、stderr 0 行**，分母＝射程内 **15 枚件**（逐名列在件里）；正控＝同尺换 `go test` ⇒ **7 枚**；再补三把我扩的射程：`tools/d22scan/runtests.sh`（**223 行**，CI `run:` 第 82 行真在调用它）⇒ **0 命中**；`grep -rn GOFLAGS` 于 CI+scripts ⇒ **0**；`ci.yml` 的 **37 枚** `run:` 步骤里含 test 的 **13 枚**逐枚看 ⇒ ⛔ 一枚带 race。反向正控＝不加 `--include` 时命中 **7 枚** `scripts/spike/bin/*.exe` 二进制（与票面那句逐字对上）。
- 主张②「neighbour `levels` has been an atomic.Uint64 all along」⇒ **成立**。射程＝该件**全部 8 个 blob**（⛔ 抽样；`git log --follow --format=%H`＝8 枚，从引入笔 `5d407e77`）⇒ 8/8 枚都逐字 `	levels atomic.Uint64`；反向尺 `git log -S 'levels uint64' --oneline -- <件>`＝**输出 0 行**。
- 主张③「手势体跑在球的 ui-sta 线程上（见 `resident_ball_windows.go`）」⇒ **成立**（机制链四段现跑，⛔ 只信注释）：`:327`／`:331` 两枚手势闭包 → `muteGesture`（`:593`）`:605` **同步**调 `fn()` → 那枚 fn＝`resident_windows.go:342` `rb.attachMuteGate(raudio.toggleMute)` → 线程由 `internal/ball/ball_windows.go:170` `Registry.Spawn("ui-sta", …)` 命名。跨线程的另一半也现跑：挂载 `:342` **早于**唯一生产读 `:385`。⚠ 票面"现量"那四枚行号（`:245-246`／`:482`／`:501`／`:542`）在 HEAD 已偏 **+6／+6／+78／+83**（同一把内容尺现量到 `:251-252`／`:488`／`:579`／`:625`）＝票面自己声明过"行号是快照"，⛔ 伤主张、具名给编排者。
- ★本格最要紧的一句（"改形状这件事被仪器看得见吗"的唯一直接证据）：**⛔ 看得见**。腿的 `logs/60-ac5-race.txt`（360,724 字节／3,073 行）23 块 `WARNING: DATA RACE` 里，含 `resident_audio`／`toggleMute`／`posture`／`audio`／`atomic` 的块＝**0/23**（块级尺，本腿按切块 python 数，⛔ 全文 `grep -c` 一把尺）；`gate.go` 帧 6 块全是 `internal/agent/approval/gate.go`（审批门），`started` 帧 0 块（4 块里的 `started_at=` 是 migration 的日志属性）。⇒ 今天能让这枚字段"被看见"的只有**形状变异仪器**，也就是本腿这两发。
- ★腿的**结论层越界一句**（具名，⛔ 我改）：`30-gates.md` §3 写"本票那枚 `started` 的形状在 race 探测器前⛔ 需要新用例来证明"。前件真、后件⛔ 成立——票面 `AC#2` ⓒ 是条件式（"**若**要造两线程对拉 ⇒ 自带正控"），⛔ 等于"不必造"；而本腿量到的恰是"⛔ 任何仪器看得见"。建议：素材层按"成立"翻，结论句那句请编排者裁要不要退回给落地腿改措辞。
- 件＝`logs/claims-rulers.txt`。

## 顺带（派单 §2 的复核格，⛔ 重复腿的活，只挑差集与尺口径）

AC#1／AC#3／AC#5 三格的独立尺与两处口径瑕（`栈顶` 措辞 ⇄ 我量到的最内层帧；腿的修复笔自述"净增 18 行"⇄ 量到 **＋19**）＝件 `20-review-cells.md`。★本腿新尺一枚：把三枚测试件里每枚 `ra.started.Load()` 复原成 `ra.started` 再与改前 blob `diff` ⇒ 两枚 **rc=0 零差**、第三枚只剩 `%v`→`%t` 那一行 ⇒ "⛔ 放宽任何一条断言"是**量出来的**，⛔ 读码推的。AC#4 的色＝本腿 pristine 那一发红名册与腿 `logs/44`（D 发）**逐字同三枚**、多 0 少 0；E 多的 `TestTicket223HandEditedFsLooseningCostsAnL2Card` 本腿三发全 PASS＝腿自己具名的审批卡片 flake。
