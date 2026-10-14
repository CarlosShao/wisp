# 255-v3 验收腿锚 — 票 255 AC#6／A878 两处重裁的非实现者判语（锚＋自造反证＋还原证）

本腿＝验收腿（非实现者；实现者 `255-r7`）。判语正文＝`docs/evidence/s1/ticket255-255-v3-accept-v1.md`；本件只放锚、窗口、还原证与件清单。⛔ 不翻票框、⛔ 不动 `Status:`、⛔ 零 push。

- 起手 HEAD＝`94398165`（进本腿时；`dcaff922`/`9bd5ff21`/`94398165` 三笔已在树）。**在跑期间编排者顺路落 `edf643c9`→`9e9964f8`→`2bb6394a`→`07a4d5ba`（frontend/**＋docs/reports/**＋probes/orch/**，与判面零重叠）**；判面两枚文件全程停在交付态 sha（每发窗口前后复核）。
- 交付态 sha：roster `6b840071c183b3ed6470935508247037c416ed3787b22c1b50f839a56f2b55f2`／test `2195e669aab19be80643d0040be67c3fdbd716d082b4e04de25d60caf68dfa4e`（与 `9bd5ff21` 件内自报一致）。
- 环境：win32 + Git Bash，go1.27.1 windows/amd64；`PATH` 带 **MSYS 形** `third_party/sherpa-onnx`；`GOTMPDIR=C:\tmp\wisp255v3gotmp`。判据＝`go test -run 'TestTicket255'`（⛔ 未跑整包，314 死区）。

## 门禁（本腿自跑）

| 门 | rc | 读数 |
|---|---|---|
| `go test -count=1 -v -run 'TestTicket255' ./cmd/wisp/`（09:27:09→09:27:39） | 0 | 20 PASS／0 FAIL；`ok … 23.186s`（log `logs/gate-test255-v3.txt`，24370B） |
| `gofmt -l <两枚>` | 0 | 空 |
| `python verify-roster-255v3.py`（自写复刻扫描，⛔ 不 import 被测件） | 0 | 14 cite 全 unique＋行号==现读；14 段全 match；尺①27/20、尺②23/22、测试侧 2（`logs/static-verify.json`，8541B） |

## 自造反证六发（逐发原文在 `logs/`；窗口内零他腿写入判面）

| 发 | 突变 | 窗口（+08） | rc | 判句（摘） |
|---|---|---|---|---|
| P1 | `internal/observe/logging.go` cited 行上插一行 | 09:29:14→09:29:30 | 0 | Logf 漂读 51→52 ⇒ ⓐ 绿 |
| P2 | `cmd/wisp/logsink.go:149` token 改掉 | 09:29:41→09:29:54 | 1 | `row "observe" … token is now nowhere in that file (recorded line reads …)` ⇒ ⓑ 红逐名 |
| P3 | `logging.go` 加二义注释（2 命中） | 09:30:04→09:30:14 | 1 | `not unique … hits lines [51 52]` ⇒ >1 红 |
| P4 | 名册 llm cite `:98`→`:96`（错读数） | 09:30:30→09:30:44 | 0 | 仅 Logf ⇒ `:LINE` 只当读数 |
| P5a/b | `internal/ball/ball_windows.go` 插行：新尺／同窗换旧尺 | 09:30:55→09:31:19 | 0／1 | 新尺 Logf 64→65；旧尺 `that line now reads "// probe…" - the evidence drifted` |
| P6 | 父态（`dcaff922^` roster＋test）跑旧尺 | 09:33:43→09:33:52 | 1 | `row "llm" … :96 … now reads "view := panel.SettingsView{"` ⇒ 交付前实红缺陷面 |

## 还原证（`sha256sum` before==after；`git status --porcelain <面>` 事后空）

- `internal/observe/logging.go` `ae6ddfe2…0a189e`（P1/P3）；`cmd/wisp/logsink.go` `f06e2a5b…ae62c11`（P2）；`internal/ball/ball_windows.go` `89c54360…5984d1`（P5）；roster `6b840071…b55f2`（P4）；test `2195e669…dfa4e`（P5 换回）。备份＝`C:/tmp/wisp255v3backup/`（含 P5 装入的旧尺 `5701ce4c…67fd275`、P6 的父 roster `367d4a43…278cb884`）。
- 突变窗口逐发起止已列上表；在飞 `237-r1` 面＝`internal/tools`（本轮它落 `9dade451` 已在起手之前），与本腿零文件重叠。

## 件（入库＝本 commit）

- `docs/evidence/s1/ticket255-255-v3-accept-v1.md`（判语，11180B）
- `.scratch/wisp/probes/255/v3/00-anchor-and-gates.md`（本件）
- `.scratch/wisp/probes/255/v3/verify-roster-255v3.py`（独立复量器，7145B）
- `.scratch/wisp/probes/255/v3/logs/probes-summary-v3.txt`（汇总，9933B）

**留盘不入库（点名＋字节数）**：`logs/gate-test255-v3.txt`(24370B)／`logs/static-verify.json`(8541B)／`logs/probe-P1-insert-logging.txt`(837B)／`probe-P2-token-change-logsink.txt`(1013B)／`probe-P3-nonunique-logging.txt`(994B)／`probe-P4-wrong-reading-roster.txt`(849B)／`probe-P5-old-ruler-red.txt`(1741B)／`probe-P6-parent-state-red.txt`(1274B)。

## 具名没做 / 下一位要点

- ⛔ 未跑整包、`go vet`、d22scan、长度门（不在本腿门禁表；本腿零 Go 产码新增）；⛔ 未动三冻结件/thresholds.go/golden/PLAN/specs；⛔ 未动别人脏件；⛔ 未 push。
- 判语要非实现者裁的形已裁（对照件 §1 表）：AC#6＝成立；削尖＝可以接受（具名残留：llm/ball 两枚含 gofmt 对齐空格，格式重排即红，fail-loud）；A878 两处＝各自成立；零放宽＝成立（一处具名降级＝L445 断言去行号）。翻框手续归编排者。
