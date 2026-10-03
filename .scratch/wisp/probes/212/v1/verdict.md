# 212-v1 验收裁决（第二任；非实现者）

- 验收者：212-v1 第二任（前任死于宿主服务错误，留基线 `.scratch/wisp/probes/212/v1/gate-baseline.txt`——**本裁决复用该基线，具名声明**）。
- 实现者：编排者本人（A544 预授权代笔，台账 A578，锚 `5e8748b3`）。按 D22 双角色本件以非实现者攻判据。
- 起手锚：`dffd9456`（HEAD，含 5e8748b3 全部实现件）；date 起手 `2026-10-03 21:19 +0800`。
- 纪律：勾框零碰；产码零改；frontend/**、design/** 零读零引。

## §0 起手锚与基线

- 复用件：`gate-baseline.txt`（前任死于服务错误前的门禁全量读数，34 PASS/76 RUN + 全仓 clean，八枚 ban 读数 internal/=228 cmd/=38 frontend/=85 tools(=#7 internal/tools/)=23 design/=39 internal/#8=498 cmd/#8=96）。
- 本任自跑门禁第一发（`sh scripts/d22scan.sh`，存 `v2-gate-run1.txt`）：rc=0，PASS=34 FAIL=0，=== RUN=76，八枚 ban 读数与基线**逐名相等**，verdict clean。AC#5 第一向成立（详见 §3）。

## §1 恒真两问

（探针进行中，本节先落骨架；探针件全部建在仓外临时目录，跟踪树零改动，完成后逐条回填读数。）

- ① 摘发射点：待填
- ② 反形判据：待填

## §2 symRefRe 三发探针

（同上，先落骨架。）

- (a) 待填
- (b) 待填
- (c) 待填

## §3 AC#5 八枚 ban 对账

- 门禁第一发 `v2-gate-run1.txt`：bans #1-5 internal/=228、cmd/=38；ban #6 frontend/=85；ban #7 internal/tools/=23；ban #8 design/=39、frontend/=85、internal/=498、cmd/=96；verdict clean。**与 A578 记录及前任基线逐名相等，零偏移。**
- 正控两向：positive control（种子违规必红）＝ runtests.sh PASS=34 里含 TestScanDetectsAllSeededViolations（9 种子→8 ban 逐枚点数，want map 逐名核过：bare-goroutine 2/pathresolver-bypass 1/plaintext-key 1/wallclock-timeout 1/mirror-hash 1/internal-artifact-tool 1/panel-approval 1/emoji 1，且断言"无多余 finding"）；self-test 两向＝36 案（20 ring/16 silent）全过。⇒ 既有 8 枚 ban 的牙未被 ban #9 挤掉。

## §4 十枚注释复核表

（逐枚读 diff 进行中，先落骨架；每枚判"更正后的话是不是真话"。）

| # | 文件:处 | 更正前→更正后 | 真话？ | 凭据 |
|---|---|---|---|---|
| 1 | internal/risk/pathresolver.go:28 | `scripts/check-pathclean-ban.sh`→`tools/d22scan 的 pathresolver-bypass ban; see tools/d22scan/main.go` | 待填 | |
| 2 | internal/agent/approval/pending_read.go | `in tools/gate.go`→`enforced by tools/d22scan's pathresolver-bypass ban` | 待填 | |
| 3 | internal/agent/approval/queue.go（两处） | `tools/bridge.go`→`internal/tools/bridge.go` | 待填 | |
| 4 | internal/agent/tools.go:89 | `class internal/tool`→`error class internal-tool` | 待填 | |
| 5 | cmd/wisp/models.go | `internal/engines/ directory`→`engines directory under internal/` | 待填 | |
| 6 | internal/ball/statevisual.go | `docs/evidence/s1/62-*`→`the ticket-62 evidence tables (under docs/evidence/s1)` | 待填 | |
| 7 | internal/risk/provenance.go | `reaches tools/agent/cmd`→`reaches the plugin agent command` | 待填 | |
| 8 | internal/tools/bridge.go | `internal/provider`→`internal-provider` | 待填 | |
| 9 | cmd/wisp/config_readers_255.go（hotkey 行+Hotkey 白名单） | `:257 Hotkeys: ball.DefaultHotkeys()`→`:269 Hotkeys: cfg`（258 形 A） | 待填 | |
| 10 | cmd/wisp/resident_approval_*246*_test.go（两枚 3 处） | `startResidentBall(reg, …)` 签名跟 258 改形 | 待填 | |

## §5 AC 格判语

- AC#2：待填（判语归验收者；勾归编排者）
- AC#3：待填
- AC#5：待填

## §6 推翻清单

（票面、A574/A578、派单每一句都是待验断言；推翻一条记一条。）

- 待填

## §7 判不动／量不到

- 待填

## 终态

- status：骨架已落（进行中）
