# 191-stream-tools-r1 — 票 191 甲支第一波（工具侧增量交付）· 2026-11-20

**腿**：写面 `internal/tools/**`；⛔ cmd/wisp、⛔ internal/panel、⛔ internal/risk 全程零字节（见 §5 突变面核）。
**现状基线**：`docs/evidence/s1/191-readonly-survey-r1.md`（读协程在、marker 整段取、零 onUpdate）。

## 1. 交付了什么

| 文件 | 改动 |
|---|---|
| `internal/tools/shell_session.go` | ①新类型 `SessionOptions{OnChunk func(chunk string)}`（票面指名的注入形）；②`run()` 增第五参 `opts SessionOptions`，**在写命令行之前**装 sink（回显字节也覆盖）、`defer` 清 sink（每条返回路径）；③`guardedText` 加只读第二面 `setSink`/sink——`takeAll`/缓冲语义一字未变，`append` 在锁外把新增字节交给 sink；④`Execute` 把既有 `onUpdate` 接缝接进 `SessionOptions{OnChunk: onUpdate}`，**Execute 签名零改动** |
| `internal/tools/shell_session_stream_191_test.go` | 新牙三枚＋分片器/洪水器（§3） |
| `internal/tools/shell_exec.go` | **一字未动**（一次性输出语义保持） |

**chunk 契约**（写死在 `SessionOptions` 注释里）：读协程（唯一一条，`observe.Spawn` 名册＋recover，无新裸 goroutine）每读一发就把**原始流**（回显命令行＋命令输出＋宿主 marker 行）按读序交给 OnChunk；marker 剥离仍只发生在最终 `Result`（`findMarkerCode` 不动）。回调 nil ＝ pre-191 行为逐字节。桥侧 `bridge.go:732-738` 本就把 onUpdate 转发到宿主 runFeed——工具侧增量即走这条**既有**管子，本波未放宽任何边界；面板/接线是下一波。

**上限与节流**：每发 chunk 累计帽＝`sessionPendingCap`（64MiB，**沿用现有限制**），与搜索循环对缓冲的那把帽同源同值；过帽停发、由循环既有的帽路径回收会话并拒绝该发。无新增缓冲、无新增裸 `go func(`（d22scan ban #1 clean，§4）。

## 2. 四发牙（含变异红读数）

| # | 牙 | 测试 | 读数 |
|---|---|---|---|
| ① | 回调 nil ⇒ 与今天逐字节一致 | `Test191NilCallbackByteIdentical`（nil 与非 nil 成对跑、逐字节相等）＋整套 163 既有套件全绿（全部 nil onUpdate） | PASS |
| ② | 回调非 nil ⇒ 分片按序拼接 == 整段输出 | `Test191StreamChunksConcatenateToRawStream`（分片器逐段写管、逐段记账；marker 行拆三段跨读；拼接==记账原流；另断言 Result 仍无 marker、`len(chunks)>=2`） | PASS |
| ③ | 读协程改成只发一次 ⇒ 红 | **变异 M3**：临时改 `guardedText.append` 只投递首个分片 → `Test191StreamChunksConcatenateToRawStream` **红**（红读数：`concatenated chunks must equal the raw stream: got "ok: b0c4…0\r\n" want "ok: echo hello 2>&1\r\nw1637a59…0\r\n"`），还原后绿 | 红已证 |
| ④ | 去掉上限 ⇒ 红 | **变异 M4**：临时同时摘掉 sink 计数帽＋循环 `len(pending) > sessionPendingCap` 拒绝块 → `Test191StreamChunkCapIsTheBufferCap` **红**（`over-cap command must take the cap refusal: 这条命令没在 C22 预算内回来…失去同步`，60s 到点才红＝帽路径消失的形状），两处还原后绿 | 红已证 |

牙④的常驻尺：`Test191StreamChunkCapIsTheBufferCap` 用洪水器灌 `sessionPendingCap+1MiB` 且永不给 marker，断言①仍走「单发上限」拒绝路径（不是靠 deadline）②消费端累计 `≤ sessionPendingCap`（增量面不开第二个无帽缓冲）。

## 3. 四数点名（`go test ./internal/tools -count=1 -v`，本机 2026-11-20 21:2x）

**RUN=601 / PASS=600 / FAIL=0 / SKIP=1**（`ok 26.516s`）。
SKIP 点名＝`Test23LiveClipboardRoundTripThroughRealOS`（真机剪贴板，存量既有 Skip，非本波引入）。
本波新牙逐枚点名：`Test191StreamChunksConcatenateToRawStream` PASS（0.00s）／`Test191NilCallbackByteIdentical` PASS（0.00s）／`Test191StreamChunkCapIsTheBufferCap` PASS（0.29s）；163 全族 16 枚 PASS。
`-race -run 'Test191|Test163'`：**ok 4.928s**（sink 面＝reader 协程单写、`guardedText.mu` 互斥、`delivered` 仅 reader 触碰，race 净）。

## 4. 门禁

| 门 | 命令 | 读数 |
|---|---|---|
| gofmt | `gofmt -l internal/tools` | `internal\tools\gate.go`、`internal\tools\ticket49_limits_test.go` 两枚——**HEAD blob 侧同脏**（`git show HEAD:<f> | gofmt -d`＝33/12 行，来自 `5d1f1e6df`），**非本腿字节、他人已入库件，只报不动**；本腿两件（`shell_session.go`／`shell_session_stream_191_test.go`）两把尺全净 |
| gofumpt | `D:\work\base\gopath\bin\gofumpt.exe -l internal/tools`（v0.12.0 go1.27.1） | 同上两枚存量，本腿件零输出 |
| go vet | `go vet ./internal/tools` | **rc=0** |
| d22scan | `sh scripts/d22scan.sh` | **clean**；末行逐字：`d22scan.sh: OK - control rc=0, scan rc=0`；`ban #7 internal/tools/=52`、`bans #1-5 internal/=424` |

## 5. 起手⇄收手突变面核

- `git diff --stat -- internal/tools` ＝ 仅 `shell_session.go`（+82/−4）＋新测试件（未跟踪）；`shell_exec.go` 零字节。
- 变异 M3/M4 已逐字节还原（还原后 `-count=1` 全包 `ok 31.126s`）。
- 禁区核：`cmd/wisp/**`、`internal/panel/**`、`internal/risk/**`、thresholds/golden/allowlist——零字节。

## 6. 只报不动（本腿观察到、未修）

1. **重置窗口的旧读者窄缝（存量，163 起）**：`ensureLocked` 在回收后重建 `s.pending = newGuardedText()`，而被回收 carrier 的旧读协程持 `&s.pending` 字段地址、退出前最后一次 `append` 与重建赋值之间无同步——理论竞态窗口，今天同样存在于纯缓冲路径（`-race` 下未现形）。191 的 sink 让旧读者若真落进这个窗口会把旧 carrier 字节交到新命令的 sink——同属该既有窗口，本波不改会话生命周期语义，只登记。
2. §4 那两枚 HEAD 存量格式脏（`gate.go`／`ticket49_limits_test.go`）——归 `5d1f1e6df`，要清另立票。
3. 桥侧事实（供下一波）：`bridge.go:732-734` 的 onUpdate 闭包**恒非 nil**（内部再判 `b.onUpdate`），故 shell.session 经桥调用时 `OnChunk` 恒非 nil、宿主未接线时增量被闭包吞掉（零缓冲、零可见变化）——「回调 nil ⇒ 逐字节一致」的直接调用者口径成立；宿主接线形态属下一波 cmd/wisp 面。

## 7. 未做（下一波）

面板/`cmd/wisp` 接线（把 chunk 真画进面板终端）、C25 名册缺口（Q-457 待机主批，本波⛔自行加）。
