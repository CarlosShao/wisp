# d22-49-clean — 清 d22scan 5 枚红（票 49 引入的写腿漏跑门）

- 腿：agent-108（写腿，只 commit 不 push；⛔ push、⛔ dist 字节、⛔ thresholds/golden/allowlist.txt/PLAN/specs/D43/design 一字未动）
- 起手锚：`04617f09b17652e00dd8ce5c7812466bb453cb54`（2026-11-17 10:43:55 +0800 现跑）
- 派单：`.scratch/wisp/probes/orch/20261117/txt/d22-49-clean-dispatch.txt`

## 1. 修复清单（5 枚红逐枚）

| # | 扫描读数（修复前） | 修法 | 语义判据 |
|---|---|---|---|
| 1 | `cmd/wisp/approval_reply.go:464` `[wallclock-timeout] unix timestamp on a timeout/deadline line` | `g.ExpiresAt <= time.Now().Unix()` → `time.Unix(g.ExpiresAt, 0).Before(time.Now())`，紧邻注释写明理由并引 D42#9 | `ExpiresAt` 是落盘 `approval_grant` 行的**绝对过期时刻**（`internal/session/grants.go` Record/RecordTarget/recordRow 写库，unix 秒，跨重启仍有效；`Now` 可注入，见 ticket49_lifecycle_test 的 `Now: func() time.Time`），不是进程内时长测量 ⇒ D42#9 禁的「wall clock 差值做超时判断」不适用：无差值运算，是持久化绝对时刻的 time.Time 比较；单调时钟表达不了跨重启的绝对时刻，落盘时间戳按 D42#9 本就该用 wall clock。判定语义与改前逐位等价（`> now` ⇔ `!Before(now)`，零值死行语义照旧） |
| 2 | `cmd/wisp/approval_reply.go:379` `[emoji]` ⚠(U+26A0) | 返回文案 `"⚠ 危险档：…"` → `"危险档：…"`（纯文字，零 emoji） | UI/返回文案零 emoji（D23/ban #8），字符串不豁免（Q-46(c)） |
| 3 | `internal/session/ticket49_lifecycle_test.go:172` `[emoji]` ≠(U+2260) | 断言消息 `(Word ≠ 终端)` → `(Word != 终端)` | 扫描器射程内（字符串不豁免）；注释里的 ≠ 豁免、未动 |
| 4 | `internal/tools/ticket49_limits_test.go:139` `[emoji]` 🔒(U+1F512) | `(🔒 安全配置放宽不适用会话授权)` → `(安全配置放宽不适用会话授权)` | 同上；断言本体一字未动，只改消息文字 |
| 5 | `cmd/wisp/ticket49_assembly_test.go:453` `[emoji]` ≠(U+2260) | `(Word ≠ 终端)` → `(Word != 终端)` | 同上 |

⯈ 全部 5 枚：断言条件、阈值、测试存在性零变动；没有删/弱化任何测试；allowlist.txt 未开过。

## 2. 三样现跑读数（2026-11-17 10:4x–10:5x +0800，本腿本人）

1. `sh scripts/d22scan.sh`（D:\work\soft\Git\bin\bash.exe）：**`d22scan: clean - no D22 ban violations`**，末行 **`d22scan.sh: OK - control rc=0, scan rc=0`**。分母（末行原文）：`bans #1-5 internal/=411, bans #1-5 cmd/=91, ban #6 frontend/=141, ban #7 internal/tools/=52, ban #10 internal/risk/=21, ban #8 design/=39, ban #8 frontend/=141, ban #8 internal/=1039, ban #8 cmd/=324`；另 `examined 502 production Go files under internal/ and cmd/`。修复前同尺：`d22scan: 5 finding(s)`、末行 `d22scan.sh: FAIL - control rc=1, scan rc=1`。
2. 载具 `scripts/wisp-cli-staged-tests.sh -test.run 'Test49|TestTicket49' -test.v`（跑前 `Get-Process wisp.test` 空、一发）：cmd/wisp **PASS / `ok github.com/CarlosShao/wisp/cmd/wisp 134.00s`，rc=0**；`go test ./internal/session -run TestTicket49` **ok 1.275s rc=0**；`go test ./internal/tools -run TestTicket49` **ok 0.127s rc=0**。
3. `go build ./...` **rc=0**。

## 3. 写面（本笔提交的显式 pathspec）

`cmd/wisp/approval_reply.go` · `internal/session/ticket49_lifecycle_test.go` · `internal/tools/ticket49_limits_test.go` · `cmd/wisp/ticket49_assembly_test.go` · 本证据件。工作树里其他腿的改动（`.scratch/**`、`frontend/fixtures/**` 等）一字未碰、不进本笔 pathspec。
