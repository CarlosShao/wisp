# 票 289 / 腿 289-v1 — AC#3 判语：三数并排＋那一枚"转绿"的独立攻击（a／b／c）

**判语：成立（新增红 0 枚这一条本腿独立复算＋自跑双证），但带三处具名残余**：① r1 的 `raw-pre.md` 里**没有 rc 行**（`grep -n "^rc=" raw-pre.md` 零命中；`raw-post.md:2426` 才有 `rc=1`）⇒ 改前那一发的 `rc=1` 只在 `00-anchor-and-baseline.md` 的叙述里，证据件本身没带；② 那一枚"转绿"的**直接原因未定**（本腿 21 发一发都没复现出红），r1 的"同包串跑时序"归因本腿**不收为结论**；③⚠ 本腿把"改后整包串跑"跑了**两发**，量到 **`cmd/wisp` 的红名册在同一棵已入库的树上两发之间自己会变**（±1 枚，见 §6）⇒ "三数并排"这一把尺今天带噪声，AC#3 的实质结论只有压在 §3 那把语义面零的尺上才站得住。

## 1. 本腿自己跑的那一发（改后态＝当前 HEAD，产码笔已入库 `73c30dfe`）

命令逐字（CWD＝仓根）：

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -v ./internal/panel/ ./cmd/wisp/ -count=1
```

原始输出＝本目录 `raw-v1-post.md`（2426 行），末两行逐字：`FAIL` ／ `rc(go-test-full-post)=1`。
本腿自己的计数（尺＝`grep -c '^--- X: '`，顶层）：**PASS=390 / FAIL=11 / SKIP=2**；包级行逐字 `FAIL github.com/CarlosShao/wisp/internal/panel 6.271s`、`FAIL github.com/CarlosShao/wisp/cmd/wisp 451.539s`；环境红尺 `grep -c '0xc0000135' raw-v1-post.md` ＝ **0** ⇒ DLL harness 生效、用例真跑了。

## 2. r1 那两个数本腿**复算**（不是复述）——⛔ 改前那一发本腿无法复跑，具名说做不到

对 r1 落盘的 `raw-pre.md`／`raw-post.md` 用**同一把尺**自跑计数：

| 件 | 顶层 PASS/FAIL/SKIP（本腿 `grep -c` 自数） | 含子测试的同一把尺 | `0xc0000135` |
|---|---|---|---|
| `raw-pre.md`（r1 改前） | **389 / 12 / 2** | 573 / 12 / 2 | 0 |
| `raw-post.md`（r1 改后） | **390 / 11 / 2** | 574 / 11 / 2 | 0 |
| `raw-v1-post.md`（本腿改后） | **390 / 11 / 2** | — | 0 |

红名册逐名作差（本腿自己 `grep -E '^--- FAIL: '` → `sort` → `comm`，件＝`roster-fail-pre.txt`／`roster-fail-post.txt`／`roster-fail-v1post.txt`）：

- **新增红（post ⊖ pre）＝0 枚**：`comm -13` 输出空（`rc(comm)=0`）。
- **反向作差（pre ⊖ post）＝1 枚**：`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`（`rc(comm2)=0`）。
- 本腿改后那一发的红名册与 r1 改后名册**逐字节相同**：`diff roster-fail-post.txt roster-fail-v1post.txt` 输出 **0 行**（`rc(diff-vs-r1post)=0`）⇒ r1 的"改后三数"在**已入库的产码笔**上可复现，不只是它工作树里那一刻的读数。
- 时间戳旁证（⛔ 不是本腿跑的，是两份 raw 里的日志时间）：pre 首末 `time=2026-10-09T14:07:53`→`14:15:35`，改后 `14:19:46`→`14:27:28`，而产码笔 `73c30dfe` 的作者时间戳是 `Fri Oct 9 14:26:18 2026 +0800` ⇒ 改前那一发确实**早于改动落笔**，不是事后补的。

⛔ **做不到的一条具名**：本腿不能复跑"改前基线"——复跑要把那 4 枚件回到 `fd692f7e` 的内容，那是改源码/动工作树（硬约束禁止，共享树里也会吞别人的活）。所以改前的 12 枚红名册＝**本腿对 r1 原始输出的独立复算**，可信度低于自跑；它的"新增红 0 枚"结论另有 §1 那把不依赖它的尺兜着（本腿改后名册 ⊆ r1 改前名册，且逐名可归因）。

## 3. (a) 一把能答"注释-only 能不能改变测试行为"的尺

见 `20-ac2-line-neutrality-and-rot.md` §3：四枚件剥掉整行注释后 `fd692f7e` vs `73c30dfe` **md5 逐枚相同**（`57b80471…`/`c62a823e…`/`38eed8e1…`/`9369cb5f…` 两态一致），`git show -U0` 非注释行＝**0**（28/28），编译指令形注释＝**0**（`grep -cE '^[+-][[:space:]]*//[[:space:]]*(go:|line |export|extern|go:build|\+build)' raw-v1-diff-u0.md` → `0`，rc=1＝零命中）。
⇒ **两态喂给 `go test` 的 Go 语义体逐字节相同 ⇒ 那一枚转绿不可能是本票造成的。**这一条判"成立"是硬的。

## 4. (b) 靶向复跑：本腿跑了 21 发，**一发红都没复现出来**

件＝`raw-ticket223.md`（隔离 3 发 + `-count=5`）与 `raw-ticket223-probes.md`（组串跑 + `-cpu=1` ×2）与 `probe-conc-1..5.md`（5 枚进程并发抢 CPU），每发自己的 `rc`：

| 形状 | 命令逐字（省略共同前缀 `PATH=… go test ./cmd/wisp/`） | 结果 | rc |
|---|---|---|---|
| 隔离 ×3 | `-run '^TestTicket223ModeLooseningChangesTheRunningModeAfterAllow$' -count=1 -v` | PASS 2.02s／2.22s／2.04s | 0／0／0 |
| 同二进制重复 | 同上，`-count=5` | 5 枚全 PASS（2.01–2.05s） | 0 |
| 组内串跑 | `-run '^TestTicket223' -count=2`（本文件全部 223 用例串跑两遍） | `ok … 53.764s` | 0 |
| 压调度 | 隔离 `-count=1 -cpu=1 -v` ×2 | PASS 2.03s／2.00s | 0／0 |
| 压负载 | 5 枚进程**同时**各跑一发 `-count=1 -v` | 5 枚全 PASS，耗时 **2.82–3.00s**（隔离是 2.0s 级） | 0 ×5 |
| in-situ（整包串跑） | `… go test -v ./cmd/wisp/ -count=1` 本腿那一发 | 该用例 **PASS 2.14s**（`raw-v1-post.md:749`） | 包级 rc=1（别的红） |

⇒ 命中率：本腿自跑的 21 发里**红 0 发**；历史 3 发整包串跑里 r1 改前 1 发红、r1 改后与本腿改后各 1 发绿。负载确实改变了它的耗时（2.0s → 2.8–3.0s），但没把它压红。

## 5. (c) 判到哪一步、⛔ 不替它编解释

**本腿能证的**：那一枚红**不是本票造成的**（§3 的语义面零尺，硬）。
**本腿不能证的**：它当时为什么红——所以按派单的 (c)：**"仅够说没新增红，转绿原因未定"**。
可以具名的是**这一枚断言天生时序敏感的形状**（盘上现读，⛔ 不是解释那一次的结果）：

```
542:		card := r.awaitCard(t, configReloadTool)
543:		shown := r.h.out.String()
544:		if !strings.Contains(shown, "risk.permission_mode") {
545:			t.Errorf("the card does not name risk.permission_mode: %s", shown)
```

而 `awaitCard`（`:252-259`）轮询的是 **gate 侧的注册表** `r.rt.liveCards.pending()`（该注册表在 `cmd/wisp/approval_reply.go:502 rt.liveCards.bind(rt.gate, …)` 绑到 gate），`:543` 读的是**人眼那面的 stdout 缓冲**——两个面之间这一格里没有任何同步边；同一份 harness 里明明备着"等 stdout"的那把尺（`:148 func (r *reloadRun223) awaitStdout(t *testing.T, needle string) string`），`:543-545` 没有用它。r1 那次改前的失败读数也相容：`raw-pre.md` 里那一行打出来的 `shown` 只有两条启动横幅（`wisp run: 答复监听已接入（卡片上给了编号）…`／`wisp run: 配置热加载已接管…`），**卡片正文一行都还没有**。
⇒ 这一条**支持"时序敏感"这一族的形状存在**，⛔ 不构成"那次红的成因已被证实"，更不构成"本票修了它"。r1 的归因按派单要求**不收下**，写成〔未定〕。

## 6. 第二发整包串跑（同一棵树，in-situ 尺）——⚠ 本腿量到"名册尺自己会漂"

命令逐字（同 harness，只跑一枚包）：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -v ./cmd/wisp/ -count=1`
原样＝`raw-v1-post-run2.md`（1934 行），末两行逐字 `FAIL github.com/CarlosShao/wisp/cmd/wisp 449.314s`／`rc(go-test-cmdwisp-run2)=1`；计数（本腿的尺）＝**PASS 265 / FAIL 6 / SKIP 2**，`0xc0000135`＝**0**。

那一枚 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 在这一发仍是 **PASS 2.11s**（`raw-v1-post-run2.md:249`）⇒ in-situ 累计 3 绿 1 红（红只有 r1 改前那一发）。

**但这一发翻出了一枚前面三发都没红的用例**（名册尺＝本腿 `sort`＋`diff roster-fail-v1post.txt roster-fail-v1run2.txt`，`rc(diff-run2)=1`）：

- `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`（定义 `cmd/wisp/resident_sink_nail_127_windows_test.go:565`）——
  在 `raw-pre.md`／`raw-post.md`／本腿 `raw-v1-post.md` 三发里逐字都是 `--- PASS`，在 `raw-v1-post-run2.md:1354` 是 `--- FAIL (2.90s)`。失败读数逐字（`:1352-1353`）：
  `resident_sink_nail_127_windows_test.go:605: … holds a pipeline file with no record in it, so the ordering has nothing to be read out of; the leg exited exit code 3221225786 (exit status 0xc000013a)`
  ／`resident_sink_nail_127_windows_test.go:571: resident leg exit: exit code 3221225786 (exit status 0xc000013a)`。
  ⚠ 那一枚终止码是 **`0xc000013a`（Windows 的 CTRL_C/CTRL_BREAK 终止）**，⛔ 不是派单点名的 `0xc0000135`（那一枚在四发里的计数都是 0，本腿逐发核过）——本腿不给它编成因，只具名记形状：这一发跑的时候机器正被本腿自己的 5 进程并发抢 CPU 压着（§4 最后一行）。

⇒ 三条结论写给编排者：① 这一枚新红**与被验物无关**（同一棵已入库的树，两发之间自己翻，产码笔 `73c30dfe` 早已落定）⛔ 不算本票账、也不算"新增红"（AC#3 的尺是 post ⊖ pre，这一枚在 pre 是绿的、在两发 post 里一绿一红 ⇒ 它落在噪声带里，不在差集里）；② **"整包串跑一次的红名册"不能作为单值尺**，今后这类"改前/改后逐名作差"要么各跑 ≥2 发取交集，要么把结论压到语义面尺上；③ 本票恰好两把尺都有：名册尺给出"新增红 0 枚"，语义面尺（§3）给出"不可能新增"——**判语靠后者立得住**。
