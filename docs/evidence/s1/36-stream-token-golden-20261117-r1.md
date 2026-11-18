# 36 (leg 2, agent-126) — AC#4：真 token 计数上线 + 200 块 golden 回放（按 10-08 重裁的截断形写判据）

- 时刻：2026-11-17 23:xx +08；起手锚 `eb94c46d5`（=派单存档笔）；产码笔 `f8b0bc2a8`；写面＝
  `internal/panel/**`＋`cmd/wisp/run.go`＋`frontend/src/**`＋`frontend/scripts/*`＋本件。未 push。
- **判据的形**：票 36 `:44-47` 那句 "renders with merge" 与盘上已重裁规则方向相反——产码
  `internal/panel/pump.go` 的 StreamLog 块逐字 `Overflow TRUNCATES and never merges`（票 197 leg B 重裁，
  批准 `A406`，代码 `94ea50f7`）。本格全部判据按**截断**那形写；**pump.go 的截断规则零字节改动**。
  另两处同名不同链（`cmd/wisp/panel_pump.go:397` 与 `internal/panel/pump.go:397`）已按票面要求具名区分，
  本腿引用的是 internal/panel 这一枚。

## 1. token 计数上真线（wire 半＝Go，页半＝panel.ts，同 commit）

| 件 | 落点 | 读法 |
|---|---|---|
| `ResultChunk.Usage *RunUsageRow json:"usage,omitempty"` | `internal/panel/composer.go:227-243` | 任务自己的 Cost 累计行，**EvUsage 事件原值**；omitempty＝无事件的块**不带键**（显式未知），⛔ 假 0 |
| `StreamLog.Usage(key, row)` + Chunks 附着 + 逐出同删 | `internal/panel/pump.go`（usage 台账／`Usage` 方法／`Chunks`／`enforceBoundLocked` 的 `delete(s.usage, ...)`） | 行可在 delta 前/中/后落地都附着；键被硬顶逐出时 usage 行随键删（有界内存口径不变）；**从不用文本推 token** |
| 事件喂值 | `cmd/wisp/run.go` EvUsage 支（c.stream.Usage 分支） | 同一事件的五值逐字进 result 块／runFeed／taskLive 三个读者，零派生 |
| 页半声明 | `frontend/src/lib/panel.ts` `ResultChunkView.usage?: RunUsageRow` | 合同尺双向走的 page half，同 commit 落地 ⇒ 普通双向比较从此判这枚键（无注册间隙） |
| 页面显示 | `frontend/src/components/result-history.tsx` `TokenCounter` | 有 usage ⇒ `已用 N tokens`（N=in+out，页侧渲染决定，title 带 verbatim in/out/cached）；无 usage ⇒ `token 用量未知（本任务未发布 usage 事件）`；流式/完成两态都显示 |
| 合同尺/钉同步 | `internal/panel/subagent_stream_197_test.go:139-146` | 原"恰 3 键"钉改成"无 usage 事件时恰 3 键"（带注释：第 4 键只随 usage 事件到）——**不是放宽**，wire 扩了枚键、钉跟着枚键走 |

`TestComposerContractTypesMatchFrontend` 全绿（`ResultChunk <-> ResultChunkView: 4 JSON keys reconciled`，
嵌套走 `RunUsageRow` 翻译已在 map）。

## 2. 200 块 golden（fixture 两半仪器共用）

- fixture：`internal/panel/testdata/stream200-golden-36.json`（生成器 `frontend/scripts/gen-stream200-golden-36.mjs`，
  确定性：块 i 文本含独有标记 `tok<i>` 且**每枚 delta 都带标记**＝每个中间包都能验身份；偶数块带 usage 行、
  奇数块 null；task-000 的 usage 在 delta 前落地验顺序无关；另有 10,000 rune 长流（head/mid/tail 标记定点、100 枚 delta）。
- **为什么存活键是 task-136..task-199**：逐出的是最旧键——追加 key 64 时逐出 key 0，…追加 key 199 时逐出 key 135
  ＝共逐出 136 键（0..135），存活＝最后 64 键，命名台账留最后 64 枚（key 72..135）。

## 3. 无 jank 判据＝可判事实（⛔ 墙钟）

Go 仪器（`internal/panel/result_golden_200_36_test.go` `TestResultGoldenReplay200Chunks...36`）逐发布（Append/Usage/Close 一步一包，
即生产 sink 的形状）断言：
1. **增量渲染批次数**：900 发布＝900 步（逐步实数 fixture 算出），无一步被吞或合并；
2. **每包块数上限**：每个中间包 ≤ `DefaultStreamKeys×StreamKeyHardCeilingMultiple`＝64 行（日志自身硬顶）；
3. **逐块身份**：每个包的每行恰含自己块的标记、零外块标记——合并在发生的当步就红（不止终点）；
4. **消失必有账**：`DroppedCount()=136`、`DroppedKeys()`＝key 72..135（最后 64 枚逐名），消失只走逐出台账。

## 4. 截断（对齐 pump.go:397-409 的裁语）

- 200 块跨 32 键界：存活 64 行**按到达序**、文本与 golden 逐字相等（小块 < 512 rune，截断窗 256 头+尾不咬文）；
- 逐出 136 键的标记**零出现在任何存活行**（"nobody is showing this one"，不是"someone folded it"）；
- 长流部分（低于键界）：`head 4096 + "...[truncated: 1808 runes elided]..." + tail 4096` 逐字相等、
  mid 标记不回流、`TruncationFor` 带 ElidedRunes=1808（"truncated 说得出丢了多少"）；
- 页面仪器（`frontend/scripts/render-stream200-36.tsx`，npm 名 `render:stream20036`）读**同一份 fixture** 过产码
  `ResultHistory`：中间包 40 行＋40 枚 caret（静态渲染下 RevealText 未 reveal＝行内文本身份归 Go 仪器，具名在件内注释）、
  终包 64 行每行恰一枚标记、逐出块标记零出现、截断行渲染具名省略句 verbatim。

## 5. token 对账（逐块）

- Go 侧：32 行 usage 五字段与 golden 逐字段相等；32 行**wire 上无 usage 键**（`len(row)==3`）——未知≠假 0；
- 页侧：32 行 `已用 N tokens` 用 golden 自己的数（task-136 ⇒ `已用 392 tokens`，title `tokensIn 236 / tokensOut 156`）；
  32 行显式未知句；全 markup **零 `已用 0 tokens`**。

## 6. 牙（变异全部在仓外导出树 `D:\tmp\wisp36mut\p`＝`git archive f8b0bc2a8`，⛔ 未动真仓字节）

| 发 | 变异（导出树逐字） | rc | 红名/红句 |
|---|---|---|---|
| 正控 | 未变异导出树 | **0** | 两枚 golden 仪器全 PASS（900/900、64/136/64、32+32） |
| MUT-A | `enforceBoundLocked` 换回**旧 fold**（两最旧块并一行，MUTATION-36A 注释锚 `:931`，build rc=0） | **1** | `result_golden_200_36_test.go:148: step 145: row task-001 carries foreign marker "tok000" - two chunks merged into one row, the exact shape ticket 197 leg B re-cut away` |
| 还原 | `git show f8b0bc2a8:internal/panel/pump.go` 覆写回（MUTATION 命中 0） | **0** | golden 绿；同树整包 `go test ./internal/panel` ok 24.559s |
| MUT-B | `Chunks` 的附着换成**字符数估算** `TokensOut: len([]rune(c.Text)), Currency:"EST"`（MUTATION-36B 锚 `:912`，build rc=0） | **1** | `:283: chunk task-136 usage = &{TokensIn:0 TokensOut:92 ... Currency:EST}, want the golden's &{TokensIn:236 TokensOut:156 CachedTokens:3 Micros:2360 Currency:CNY} verbatim`＋`:368: a stream with no usage event carries usage &{TokensIn:0 TokensOut:8228 ...} - a fabricated zero's cousin` |
| 还原 | 同上覆写回 | **0** | golden 绿 |

（过程注记：MUT-B 第一发曾在真仓工作目录误跑 `go test`＝台面错、读数作废，已在导出树重跑并记录于此；导出树
tar 展开时 `.scratch` 里 16 枚含异形字符的历史件解包失败＝与 Go 测试无关，1906 枚 `.go` 全数在树。）

## 7. 门禁读数（全部现跑）

| 门 | rc | 备注 |
|---|---|---|
| `go test ./internal/panel -count=1`（真仓） | 0 | ok 26.020s |
| `go build ./internal/panel ./cmd/wisp` | 0 | |
| 载具 `scripts/wisp-cli-staged-tests.sh`（`go test ./cmd/wisp`，一次一发，跑前 `Get-Process wisp.test`＝空） | 见 §8 |  |
| `npm run build` / `typecheck` / `lint` | 0 / 0 / 0 | lint 0 errors / 45 warnings（存量；本腿文件零命中——唯一近邻是 `result-history.tsx:196` 既有 ToolIcon16 警告） |
| 新仪器 `render:stream20036` | 0 | |
| 相关既有仪器 `render:history36 / stream / historypage36 / receiver / roster / composer / shell77` | 全 0 | |
| `sh scripts/d22scan.sh` | 0（control 0, scan 0） | 顺手修了一枚**别人已提交的过期注释引用**：`internal/panel/bridge.go:200` 引 `.scratch/wisp/issues/38-command-palette-tasks.md`，该票结案改名 `-done` 后断链（d22scan phantom-citation 点名）；只改注释一行指向 `-done` 件，语义零变 |

## 8. 载具读数（`scripts/wisp-cli-staged-tests.sh`，跑前 `Get-Process wisp.test`＝空，一次一发）

- **第 1 发（作废）**：`D:\tmp\wisp36mut\logs\30-carrier-cmdwisp.txt`——夜里机主机器休眠把 `-test.timeout=90m`
  的墙钟打断（超时 panic 落在语音测试 `Test16r1PremortemModelSha256CorruptRefusesTheHandOff` 运行中，
  `panic: test timed out after 1h30m0s`），整包没有跑完 ⇒ 按 A1889 形状规则**判无效**，不引数字。
- **第 2 发（有效，`31-carrier-cmdwisp-retry.txt`）**：**rc=1**（包行 `FAIL  github.com/CarlosShao/wisp/cmd/wisp 1275.00s`），
  **RUN=964 / PASS=721 / FAIL=1 / SKIP=13，零 panic 零超时**；唯一红＝
  `TestAC4FocusReturnToPriorWindowGap33r5 (4.83s)`＝**台账 A1876 名册上的已知环境红**（前台被占），
  非代码层红。本腿改的 `cmd/wisp/run.go` EvUsage 支随整包绿（964 枚含其编译面）。

## 9. 诚实边界

- AC#4 原文三半（merge 那半已按 10-08 重裁换判据、token 计数、无 jank）里，**"no jank" 只用可判事实**（批次数/每包上限/逐块对账/消失有账）——真帧率仍需真宿主窗，与票 77/AC#3 同类环境边界，本腿不造读数；
- `ResultChunkView` 的 `usage` 键经合同尺双向判（page half 同 commit）；ticket 36 票面勾框照例留给非实现者；
- 变异台件（导出树＋日志）在 `D:\tmp\wisp36mut\**`，只建不删。
