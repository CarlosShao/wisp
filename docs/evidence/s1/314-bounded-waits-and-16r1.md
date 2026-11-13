# 票 314 建议落地 ＋ 16r1 premortem 新红定因（腿 314b，2026-11-13）

派单：`.scratch/wisp/probes/orch/20261113/txt/314b-bounded-waits-dispatch.txt`。
写面：`cmd/wisp/**` 三个测试件（本波唯一允许写 `cmd/wisp` 的腿）；`models/**`、`thresholds.go`、
golden、`internal/**`、`scripts/**`、`.github/**` **一字未动**（可核：`git show --stat`）。

---

## 事一：裸通道接收 → 有界等待（票 314 诊断腿推荐 (b) 的执行）

### 改了哪几处、上界、数的出处

| # | 裸接收（改前） | 文件:行（改后） | 上界 | 数的出处（文件:行） |
|---|---|---|---|---|
| 1 | `<-done` | `cmd/wisp/panel_pump_test.go:221`（`waitDoneBounded314`，helper :107-148） | **46s** | 被等的 fs.read 自带两道内钟：Execute ctx 45s（`executeOn145`，`panel_pump_test.go:155`）＋ L2 卡窗 31s（`nonDefaultConfig145` 写 `confirm_timeout_sec = 31`，`panel_pump_test.go:181`）⇒ 上界＝45s（最长内钟）＋1s 收尸边距 |
| 2 | `<-w.done`／`<-w.reply` ×3 枚用例 | `cmd/wisp/approval_always_219ac2_test.go`（`waitRoundTripBounded314`，helper :427-479；调用点 :184/:242/:279） | **46s**（对 reply 与 done **各一次有界接收，共用同一绝对 deadline**——reply 收完后的 `<-w.done` 本身也是裸接收，若不一并加界就是漏网） | 夹具自定的 40s L2 窗（`newReplyHost` 的 l2Wait，`approval_reply_201_test.go:154`，票 201 起每枚 replyRoute 用例的值）＋第二卡 15s 等待（`waitWidenCard187` 调用，`approval_always_219ac2_test.go:135`）⇒ 上界＝40s（最长内钟）＋6s 收尾边距 |

**为什么不是「刚好让本机过」的魔法数**：46s 不是调出来的，是**引用被等工作已有的最长内钟再加收尸边距**——
正常路径里这个 select 几乎立即返回（fs.read 到期在 31s），只有两道内钟都失灵时才会走到 46s。两处 helper 的
注释里都写了出处；失败语里也把出处逐字点名。

### 失败语点名内容（四要素全齐）

两处 helper 的 `t.Fatalf` 逐字点名：**在等什么**（channel done / w.reply / w.done 及其语义）、
**等了多久＋上界多少**（46s，内钟出处行号）、**当时可判读数**——
- pump 侧：卡上采样包的 pending 张数、最新包 pending 张数、gate 实显卡数（`rt.windowCount()`，`run.go:1606`）、pump 发布数（`rt.snapshotCount()`，`panel_pump.go:529`）；
- 219 侧：fs.write 是否已返回（done 闭没闭）、放宽是否已入账（stored 字节非空＝WIDEN-APPLIED 已 booked）。

### ★正控：注入「故意不等」⇒ 红（逐字读数）

方法：`-overlay`（同仓 `Test219AC2MutationRevokeWriteMustPersistForBytesToComeBack` 的既 sanctioned 形状），
变异只落 `D:\tmp\wisp\20261113\314b\`（**零还原动作**——仓内文件从未被碰）。变异内容：
- pump：`case <-time.After(46*time.Second)` → `50ms`（该 helper 在 fs.read 仍在 ctx 内时就会停在 select，小上界自己就响）；
- 219：`const bound = 46s` → `50ms` ＋ `defer close(w.reply)` → 不关（注一个真挂起——健康轮里两个 close 先于 run() 返回，光缩小上界红不了）。

**红读数逐字**（rc=1，二进制为 `-overlay` 变异编译件）：

```
panel_pump_test.go:222: the fs.read goroutine never closed done: waited on channel done for 46s; the work it awaits is bounded by the Execute ctx at 45s (panel_pump_test.go:155) and the L2 card window at 31s (confirm_timeout_sec=31, panel_pump_test.go:181), so an expiry here means an inner bound failed to fire. At expiry: the packet sampled while the card was up held 1 pending card(s), the latest packet the pump knows holds 1, the gate displayed 1 card(s), and the pump published 2 packet(s)
--- FAIL: TestRunBooksWithASnapshotOfItsLiveQueue (2.28s)
```
（读数解释：变异把等砍到 50ms，fs.read 还在 45s ctx 里没回来，所以四数如实报「卡还开着、泵发过 2 包」——
这就是原来 8m31s 零输出时什么都没有的那份信息。）

```
approval_always_219ac2_test.go:242: the reply driver never returned: waited on w.reply 50ms (the round's whole bound); the run's own clocks bound the work it awaits at the 40s L2 window (confirm_timeout_sec, approval_reply_201_test.go:154) + the 15s second-card wait (approval_always_219ac2_test.go:135), so an expiry here means a step between them never happened. At expiry: the fs.write side DID return (done closed), the widening reached the store=true - the driver's last typed lines are in the run's stdout/stderr above
--- FAIL: Test219AC2RevokeRoundTripByteForByte (2.27s)
```

### 正控还原 ⇒ 绿

还原＝nothing to restore（变异只在 tmp overlay）。原树重跑同一批（真 carrier，`scripts/wisp-cli-staged-tests.sh`，
CGO=1，Dll staged，夹具树在盘）：

```
--- PASS: Test219AC2WidenAppliedAuditLineCarriesTheFourFields (3.47s)
--- PASS: Test219AC2RevokeRoundTripByteForByte (2.61s)
--- PASS: Test219AC2SecondRevokeOfTheSameRuleIsRefused (1.93s)
--- PASS: TestRunBooksWithASnapshotOfItsLiveQueue (32.85s)
PASS
ok  github.com/CarlosShao/wisp/cmd/wisp	41.00s   （rc=0；RUN=4 PASS=4 FAIL=0 SKIP=0）
```

---

## 事二：`Test16r1PremortemModelSha256CorruptRefusesTheHandOff` 新红定因

### :523/:525/:538 到底断言什么

| 行 | 内容 | 断言什么 |
|---|---|---|
| :523 | `CORRUPT walk exit=%d took=%v` | **t.Logf**（信息行，非断言） |
| :525 | `CORRUPT walk_log %s` | **t.Logf**（信息行，非断言） |
| :538 | `took < 5*time.Second` ⇒ Fatalf | **断言**：走查耗时必须 ≥5s，以「重新下载约 80MB 不可能 5 秒内完成」推断走查不是 cache-hit 形 |

它读的 sha256 从哪来：**不是常量、不是期望字面值**。该用例自己把已安装 encoder 翻转一字节
（`blob[len(blob)/2] ^= 0xFF`，:518），断言① `VerifyInstalled` 必须报 `sha256 mismatch`（期望值=报错文本形状），
断言② 走查后 exit=0 且修复后 store 再验干净，断言③（:538）耗时 ≥5s。**全案没有任何一处把 sha256 值写死。**

### 两发 CI 对照（逐字，`gh run view --log` 现抓）

| | 绿发 43425302222 @ `285a5c310`（重钉那发） | 红发 43431925841 @ `10f98f52b` |
|---|---|---|
| 检测半 | `CORRUPT detected verify_error="...encoder.int8.onnx: sha256 mismatch (pins 81a70226a893..., got ab9962e4c1a2...)"` | **逐字相同**（同 pin `81a70226a893...`，同误 `ab9962e4c1a2...`） |
| 走查 stdout | `ticks=9847，末次进度 100%，5.01s` | `ticks=11514，末次进度 100%，4.094s` |
| :523 | PASS（`took=5.01s`） | `CORRUPT walk exit=0 took=4.094s` |
| :525 | `hand-off OK` | 同 |
| :538 | 5.01 ≥ 5 ⇒ 过 | **4.094 < 5 ⇒ 红**：`the walk returned exit 0 in 4.094s: a cache-hit shape would mean the tampered bytes were trusted` |
| 判定 | **PASS** | **FAIL (9.00s)** |

**两发的 manifest pin 逐字相同**（encoder pins `81a70226a893...`）——重钉 commit `285a5c310` 只改 matcha/TTS
条目（`git show 285a5c310 -- models/manifest.json` 逐行核对：改的是 matcha 归档 sha/size、
`model-steps-3.onnx` 成员、3 个新成员行、条目级聚合值——与 `docs/evidence/s1/111-matcha-repin.md` §1 的
新旧对照表逐字一致；**ASR 条目一字未动**）。
且 `285a5c310..10f98f52b` 的全 diffstat（10 文件）**零 models/**、零 16r1 测试文件、零 internal/models**。

### 定性：**(a) 测试期望值过期**（且「期望」是时钟猜想，不是 pin 值）

依据三条：
1. **产品在两发都做对了**：红发里走查 ticks=11514（真下载进度事件，基线 cache-hit 基线 ticks=1）、
   exit=0、修复后 store 再验干净——「self-heal 重取＋验证」合同两条都履行了。若把 5s 防的
   「cache-hit 信任脏字节」真发生，`VerifyInstalled`（:567）会先红，轮不到 :538。
2. **:538 的 5s 是墙钟猜想，与 CI 网速赛跑**：同代码同 pin，绿发测得 5.01s（贴着线过 0.01s）、
   红发测得 4.094s——0.9s 的网络差把一个真前提判成假。这是仪器问题，不是产品问题。
3. 本机复跑（真 carrier，修复前读数）补刀：重下载实测 **17m18s**（`took=17m18.349s`、ticks=12970、
   PASS）——慢网下猜想永远不响，快网 cache-hit 下反而可能冤枉。**双向都不成立。**

「旧 pin」的假设本身被 CI 逐字证伪：两发 pin 相同，不存在「期望值绑旧 pin」；绑的是「下载 ≥5s」这个
**过期的世界假设**（重钉后 CI 冷跑夹具变化让重取路径更快触达，让这个一直存在的竞态第一次显形——
上一发没红只是贴线 0.01s 的运气）。

### 改了什么（修复，非放宽）

`cmd/wisp/voice_e2e_demo_16r1_windows_test.go`：
- **删**：`if took < 5*time.Second { t.Fatalf(...) }`（墙钟猜想）；
- **换**：翻转字节时记 `corruptAt := time.Now()`（:528），走查后读被改文件的 **mtime**
  （:583-589）：`!st.ModTime().After(corruptAt)` ⇒ Fatalf。逻辑等价且更强——走查唯一的修复路径是
  重取＋重写该文件；cache-hit 会留下翻转字节（:567 的干净验证先红）；mtime 未前移说明文件从未被
  重写。**cache-hit 回归在任何网速下都红，慢速重取在任何网速下都绿**，且零魔法数、零放宽、
  零 skip、零产品改动。
- 修复后真 carrier 定向跑：**PASS**（rc=0，逐字关键行）：
  ```
  voice_e2e_demo_16r1_windows_test.go:545: CORRUPT detected verify_error="...sha256 mismatch (pins 81a70226a893..., got ab9962e4c1a2...)" class=model
  wisp models ensure: asr-streaming-paraformer-zh-en 已交还（state=FirstRun，ticks=12970，末次进度 100%，1038.35s，...）
  voice_e2e_demo_16r1_windows_test.go:556: CORRUPT walk exit=0 took=17m18.349s (baseline pristine walk was ~2s)
  voice_e2e_demo_16r1_windows_test.go:594: CORRUPT healed store verifies clean; the tampered encoder was rewritten by the walk (mtime +462575ms after the flip)
  --- PASS: Test16r1PremortemModelSha256CorruptRefusesTheHandOff (1046.40s)
  ```

---

## 门禁（逐条如实报）

```
CGO_ENABLED=1 go build ./cmd/wisp/        → BUILD-OK (rc=0)
CGO_ENABLED=1 go vet ./cmd/wisp/          → VET-OK (rc=0)
CGO_ENABLED=1 go vet -tags=winlive ./cmd/wisp/ → WINLIVE-VET-OK (rc=0)
sh scripts/d22scan.sh                     → clean (rc=0，"d22scan: clean - no D22 ban violations")
go run mvdan.cc/gofumpt@latest -l ./cmd/wisp/ → 空 (rc=0)
```
（gofumpt 本机未独立安装，用 `go run mvdan.cc/gofumpt@latest` 现拉——尺一致。）

## 四数（定向真 carrier）

- 事一四枚（`TestRunBooksWithASnapshotOfItsLiveQueue`＋3×`Test219AC2*`）：**RUN=4 PASS=4 FAIL=0 SKIP=0**，包级 `ok ... 41.00s`，rc=0。
- 事二一枚（`Test16r1PremortemModelSha256CorruptRefusesTheHandOff`）：**RUN=1 PASS=1 FAIL=0 SKIP=0**，rc=0（该枚本体 1046.40s——慢网重下载 17m18s 属实记录，非挂起）。

## 变更面与撤销口令

- 跟踪写面恰 3 件：`cmd/wisp/panel_pump_test.go`、`cmd/wisp/approval_always_219ac2_test.go`
  （＋它挂的 `cmd/wisp/approval_always_201_test.go` 结构体一字段）、`cmd/wisp/voice_e2e_demo_16r1_windows_test.go`。
- 临时件全部在 `D:\tmp\wisp\20261113\314b\`（变异件、overlay JSON、包装脚本、独立 stage 目录）——
  **零还原动作**（仓内文件从未被碰）。
- 撤销口令：回「撤 314b」⇒ revert 本笔 commit 即复原三件测试件（无其它落盘面）。

## 诚实边界

- 定向跑未含整包回归（`cmd/wisp` 整包色本波被另一条腿的 portable-tests.sh 占用，AGENTS.md 并发铁律
  「整包同一时刻只跑一个」）；三枚改动件的家族（panel_pump / approval_always_201/219 / 16r1）定向跑全绿＋
  全门禁绿为替代读数。
- 第一次正控跑读数（22:44，`RUN=4 PASS=4`）取自 job 流；随后 `stage-cli/wisp-cli-tests.log` 被并行腿的
  portable 跑覆写（发现于 22:50 复查，log 顶部出现 `portable-20261113-224814-1153` 的路径）——
  上述四数以 job 流为准，log 覆写不影响判定。
