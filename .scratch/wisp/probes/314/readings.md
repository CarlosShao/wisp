# 票 314 测量腿 — 读数根（2026-11-13）

- 腿：314 测量/定性腿（⛔ 非写腿）
- HEAD：`397cd9e9db1ebb898f629671c408795f48898d6f`（派单写的 `02939b8a2` 是其父提交）
- 原始日志：`D:\tmp\wisp\20261113\314\`（只建不删）
- 副本：`logs/`

## 逐发读数

| 发 | 命令形 | rc | 耗时 | RUN | PASS | FAIL | SKIP | 最后一行 | 死点 |
|---|---|---|---|---|---|---|---|---|---|
| run1 | carrier 全量 | 无（卡死后我 kill） | 8m31s 后杀 | 11 | 10 | 0 | 0 | `scope disposed …` | 卡在 `Test219AC2RevokeRoundTripByteForByte` |
| iso1 | carrier `-test.run '^(5 枚)$'` | 0 | 9s | 5 | 5 | 0 | 0 | `ok … 6.00s` | 无 |
| iso2 | 同上 | 0 | 18s | 5 | 5 | 0 | 0 | `ok … 7.00s` | 无 |
| iso3 | 同上 | 0 | 13s | 5 | 5 | 0 | 0 | `ok … 8.00s` | 无 |
| iso219 | carrier `-test.run '^Test219AC2RevokeRoundTripByteForByte$'` | 0 | 5s | 1 | 1 | 0 | 0 | `ok …` | 无 |
| run2 | carrier 全量 `-test.timeout=10m` | **2**（撞 test 超时 panic） | **604s** | 380 | 295 | 2 | 0 | `panic: test timed out after 10m0s` | 卡在 `TestRunBooksWithASnapshotOfItsLiveQueue`（`panel_pump_test.go:177` `<-done`）；终点锚 `TestZZ315A1` 命中 0 |

## 桌面状态（每发前后）

| 时刻 | input_desktop | logonui(session1) | 前台 | CPU |
|---|---|---|---|---|
| 20:56:12 | Default | 0 | Weixin | 100 |
| 21:07:16 | Default | 0 | Weixin | 67 |
| 21:07:46 | Default | 0 | Weixin | 100 |
| 21:08:12 | Default | 0 | Weixin | 97 |
| 21:09:40 | Default | 0 | Weixin | 71 |

判据：`OpenInputDesktop` 名字（`Default`＝未锁屏）；session 1 的 `LogonUI` 计数（2/3/4 是断开 RDP 会话，各自 1 枚，与本机桌面无关）。
