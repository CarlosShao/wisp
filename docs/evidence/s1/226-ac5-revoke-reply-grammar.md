# 226 AC#5 撤销触发面（wave59 写腿 W226-v3 自查件，2026-11-07 15:1x +08）

> ⚠ **实现者自查件**（非对抗验收；按 `SPEC-12 §4.3` 该表的裁决仍归非实现者腿）。
> 起手锚点 `bee208958`（比派单锚点 `6ae44b1fb` 新——其间仅台账/裁决件追加，无产码交叉；
> 派单要求的三核见 §⑨）。

## ① AC#5 两半边的现跑判定（起手三核）

| 半边 | 派单时的现读 | 尺 |
|---|---|---|
| 落库半边 | **已在**：`SetAllowedDirs`（`internal/config/allowdirs.go:109`，注释自陈「the 撤销 half of AddAllowedDir」）＋ 2 枚用例（`writeguard_226_test.go:308/:313/:350`）| `grep -rn SetAllowedDirs --include=*.go` |
| 触发半边 | **零入口**：`replySurface.handle` 动词表＝`yes/session/no/veto/panel-no/panel-yes/always/head/view`，无撤销动词；`SetAllowedDirs` 生产调用者＝0 | 现读 `approval_reply.go:615-643` |
| 票 219 侧 | **未做**（仍只有 `grants.go:43` 与 `allowdirs.go:95` 的注释指名该入口），与本票 Progress log 11-06 收口普查一致 | 全仓 grep `撤销\|revoke` 产码 |

⇒ 判定＝**接线**（票面预判成立）：`always` 的对偶动词落进同一张答复语法表，两端各补有牙用例。

## ② 交付（本腿改/增的文件）

| 文件 | 动作 | 内容 |
|---|---|---|
| `cmd/wisp/approval_revoke.go` | 新增 | `revoke` 动词本体：`revokeRule`（拒绝/装配/在账/写/复读五支）＋ `revokeDirOf`（两种合法写法）＋ `dirInList`/`listWithoutDir` |
| `cmd/wisp/approval_reply.go` | 3 处编辑 | ① `replySurface` 方法表注加 revoke 行；② 开机横幅加一句；③ `handle` 加 `case "revoke"` ＋ help ＋ 未知词报错串 |
| `cmd/wisp/approval_revoke_226_test.go` | 新增 | 2 枚 E2E（组装栈、脚本答复流，§④）＋ 1 枚单元 |
| `internal/config/writeguard_revoke_226_test.go` | 新增 | 落库侧"下一启动视角"效果用例（§⑤） |

禁区遵守：零接触 `internal/ball/**`、`frontend/**`、`design/**`、`thresholds/golden/allowlist.txt`、三枚冻结件；
未翻票面任何 AC 勾或 Status；票 226/219 工单文件**一字未动**（Progress log 一行由本腿落库时提交）。

## ③ 撤销动词的形状（为什么是规则文本，不是卡片编号）

- `always` 流程把规则印在屏上（`要存的规则：[fs] allowed_dirs += "D:/x"`）＋ 审计（`WIDEN-APPLIED dir=...`）；
  卡片一被答复即 `Forget`（`replies.go:333` 等）⇒ 编号在需要撤销时已不在账，**锚规则文本才可达**。
- 两种写法：`revoke [fs] allowed_dirs += "D:/x"`（照抄屏幕/审计）或 `revoke D:/x`（裸目录）。
- **无卡**（不触 L2）：`SetAllowedDirs` 只能收紧；放宽仍只有 `always`＋第二重 L2 卡那一条门（`PLAN.md:1645-1646` 非对称，与 mode 换档同口径）。
- 无 card 号 ⇒ `record` 的 corr 恒为 `""`，审计行 `REVOKED corr="" tool=config.allow_dir route="native/revoke" dir=...`。
- 影响口径如实说：撤销与「一直」同轨，**下一次启动生效**（运行中 canonicalizer 不重建，`config_reload.go:268` 的既有句子）。

## ④ cmd/wisp 用例（组装栈＋脚本答复流，接缝 §1.3）

- `TestReplyRevokeRemovesTheStoredRuleForTheNextStart`：**always 存规则 → 第二张 L2 卡 yes → WIDEN-APPLIED 落盘 → 撤销行（趁首卡未答、监听器可证在读）→ 断言**：
  ① 撤销前 config.toml 确含该目录（存量对照）；② 撤销后文件**不含**该目录；③ `config.LoadFile`（下一启动的解析器）读出 `allowed_dirs` 恰＝ `[dir]`（手编目录还在）；④ 审计 `REVOKED` 行在、无 `REVOKE-FAILED/WIDEN-REFUSED`。
- `TestReplyRevokeRefusesARuleTheProcessDoesNotHold`：撤销不存在的规则 ⇒ 控制台"本进程账上没有"、审计 `REFUSED`、**文件字节不变**（撤销前快照对照）。
- `TestRevokeVerbRefusalsDoNotTouchTheConfigWriter`：空参／无写面两形单元（`replySurface{}`），断言报错文案、永不摸 Manager。

## ⑤ internal/config 用例（下一启动视角的效果读数）

`TestAC5RevokedRuleIsGoneForTheNextStart`：起手对照（start 内容 fresh load 确认两目录都在）→ `SetAllowedDirs(["D:\\notes"])` →
`LoadFile` 断言恰剩 `D:\notes` ＋ fresh `NewManager`（下一启动本体，其视图即 run.go:624 建 canonicalizer 的输入）同读数 ⇒
**"撤销后规则真的不生效"** 钉在效果上，不是"函数被调用"。

## ⑥ 门禁读数（跑前 unset GOTMPDIR/TMPDIR/TEMP/TMP；PATH 补 build 与 third_party/sherpa-onnx）

| 门禁 | 读数 |
|---|---|
| gofumpt（`-l`，绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`） | 4 件点名件＝**空**（第一轮 `approval_reply.go` 被点名、`-w` 修复后复绿） |
| `go vet ./cmd/wisp ./internal/config` | **ok** |
| `go build ./cmd/... ./internal/...` | **ok**（`go build ./...` 的红＝`build/mut_rvw.go`，未跟踪他腿探针件，§⑨③） |
| `go test ./internal/config -count=1` | **ok 2.537s** |
| `go test ./cmd/wisp -run` 三枚新用例 | **全 PASS** |
| `go test ./cmd/wisp -count=1` 全包 | 第 3 轮终态见 §⑧（FAIL 45 枚，**本腿 0 枚**；名册波动与归因见 §⑦） |
| `tools/d22scan/d22scan.exe` | **clean**（exit 0） |
| 牙检（突变） | 撤销写跳过（sed 突变）⇒ happy 用例 **FAIL 41.87s**（"revoked rule is still in config.toml"），恢复后复绿＋gofumpt 空 ⇒ **牙真** |

## ⑦ 归因声明（既有红不是我的；含本腿自己的一枚操作失误，如实记）

cmd/wisp 全包在本机的名册跑间波动极大（第 1 轮 fg+bg 20+ 红 → 第 2 轮 9 红 → 第 3 轮 45 红），
波动全部落在**音频／常驻／wire／panel 焦点家族**（15e/41a2r1/26P/392/381r2/390r1/394R1/26R3/382W/WireR1/12r2/256/247/AC4Focus/223），
与答复语法面零交集。逐项归因：

1. **本腿两枚 E2E 在第 2 轮全包的红**：根因＝`build tools/mockllm: exit status 1`——**go build 缓存损坏**（逐字：`could not import cmp … open E:\work\data\caches\go-build\…-d: The system cannot find the path specified`，另见 `package crypto/subtle is not in std` 怪象）。重建 `tools/mockllm`（`go build .` ok）后，rerun7 定向复跑：本腿 2 枚＋375/L1Veto/NativeHost×2 全部转绿 ⇒ 假红坐实、已修。
2. **`TestTicket223HandEditedFsLooseningCostsAnL2Card`（第 3 轮全包红）**：时序型假红——它断言"账本出现卡后"的控制台文本（`config_reload_223_test.go:434-436`，ledger Record 与 stdout 渲染之间有窗口期）；单跑复量＝**ok 2.636s**；且第 1 轮全包（本腿改动同样在飞）它是绿的 ⇒ 非本腿横幅文案所致。
3. **既有稳定红（与本腿无关）**：`TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`／`TestAC247HandingTheLevelToTheSeamIsNotVisibility`（HEAD 态对照同红）＋`TestAC4FocusReturnToPriorWindowGap33r5`（Win32 前台，三跑＋单跑均红；本会话桌面无前台权限的环境形状，票 33r5 家族）。
4. **音频家族大面积波动**：设备/负载依赖（本机共享、编队在飞），第 1 轮 fg 跑（本腿产码已落、但其红全在音频家族）即已出现 ⇒ 既有环境形状，非本腿。
5. ⚠ **本腿自己的一枚操作失误，如实记**：第 2 轮全包转后台后，我在它仍在飞时串发了三次对照/突变 `go test`（并发 go 进程）——这加重了第 2 轮的 mockllm 缓存损坏与名册噪声。第 3 轮全包＝暖缓存＋独跑；223 的单跑复量也是独跑。此后归因读数全部取自独跑。

## ⑧ cmd/wisp 全包终态名册（第 3 轮，暖缓存＋独跑，落盘件 `.scratch/wisp/probes/226/v3/cmdwisp-full-run2.log`）

- **FAIL 45 枚**，逐名归类：**本腿 0 枚**（两枚 E2E 在该轮为 PASS）；256/247/AC4Focus＝既有稳定红（§⑦3）；其余 42 枚＝音频/常驻/时序环境家族（§⑦4，与第 1 轮 fg 跑的红名册同族、跑间互有涨落）。
- **本腿用例终态**：`TestReplyRevokeRemovesTheStoredRuleForTheNextStart`／`TestReplyRevokeRefusesARuleTheProcessDoesNotHold`／`TestRevokeVerbRefusalsDoNotTouchTheConfigWriter`＝**定向 PASS**＋第 3 轮全包 **PASS**＋rerun7 **PASS**；突变牙检红 41.87s（§⑥）。
- 佐证读数：rerun7＝`.scratch/wisp/probes/226/v3/rerun7.log`（8 枚中 7 绿，唯一红＝AC4Focus）；223 单跑 ok＝2.636s。

## ⑨ 派单三核的现跑回执

1. **票面点名串仍在**：`SetAllowedDirs` 定义 `allowdirs.go:109`、用例 `writeguard_226_test.go:308/:313/:350` 现跑 grep 命中；"零调用者"半句仍真（生产调用者接完本腿＝1）。
2. **引文出处核对**：`grants.go:33-44`（invalidation side 无码＋撤销入口指名 SetAllowedDirs）、`replies.go:519`（`WideningRule` 规则文本形状）、`config_reload.go:268`（下一启动口径）——逐一现读，行号对得上。
3. **红/绿归属**：`git status --short -- "*.go"` 起手＝跟踪件**零改动**；未跟踪仅 `.scratch/**` 探针＋`cmd/wisp` 两枚 323 在飞测试件（他腿，未碰）；`build/mut_rvw.go` 为他腿未跟踪突变探针件（`git ls-files --error-unmatch` 报 pathspec 不匹配），`go build ./...` 的红归它，非本腿。

## ⑩ 边界与未尽

- **未翻任何格**：AC#5 勾与翻不翻＝非实现者裁决腿的事，本件只报读数。
- 票 219 AC#2 的"审计可查"读数（谁/何时/加了哪行）＝审计家族既有 `REPLY` 行＋`WIDEN-APPLIED`，本腿新增 `REVOKED` 行——219 验收腿可现跑取。
- 219 票面"长期要留下可撤销的一行"的完整判据表（含批量撤销 187/193 面）不在本票射程。
