# 313-r2 锚＋门禁＋自报（`cmd/wisp` 精确落点包·落地腿）

判语：按落点包逐字落两处＋`splitChainValue` 助手，另补一发端到端用例与 A881 ⓒ 并集等值钉；
门禁 6 枚全跑；新增红＝0（除 255 名册那两枚**因本腿落地而必须重裁**的具名红，见 §4）。

## 0. 锚

- 起手 `git log -1`＋date 逐字＝`logs/00-head.txt`（本腿起手时 HEAD `7be9573a`，2026-10-13 23:14:51；
  派单写 `044b9066` 是编排者落笔时的锚——共享工作树在呼吸）。
- 产码件两枚＋交付件见 §2；⛔ 零 push。落点包行号量于 `39f6341d`——**逐枚现跑复认**：`:96`（view 那行）、
  `:112-114`（`for range p.Models`）、`:218`（`FieldProviderCredential` case）、`:228`（`res.Tier`）、`:94`（不可读支）
  全部与包上一致，未漂。

## 1. 落点（逐字照落点包；三处禁区未动）

1. **落点 1** `ReadSettings`：`:96` 扩为完整 `panel.SettingsView{}` 字面量（`Presets`/`PresetsKnown`/
   `ChatProvider`/`Configured`/`ConfiguredKnown`/`Voice` 六味，逐字）；`:112-114` 改为 `modelNames`
   排序切片（`sort.Strings` 后 `ModelCount = len(modelNames)`、`Models = modelNames`，计数与名字清单同源）。
2. **落点 2** `ApplySetting`：`FieldProviderCredential` case 之前追加六 case（逐字，`voice_cloud_asr/tts_chain`
   过 `splitChainValue`）；文件尾加 `splitChainValue` 助手（逐字，注释含原文）。
3. ⛔ 三处一字未动：不可读支（现量 `:94`）／`res.Tier = panel.EffectiveRestart`（现量 `:228`）／
   `credentialValue` 那条路（`:53`、`:59-66`）。`git diff` 现读只含上述两处＋助手（`logs/` 里的 diff 尺：47+/3-）。

## 2. 交件

| # | 件 | 行数 | 是什么 |
|---|---|---|---|
| 1 | `cmd/wisp/panel_config_store.go` | 改 | 落点 1＋落点 2＋助手（gofumpt 已洗） |
| 2 | `cmd/wisp/panel_config_313_test.go` | 263（新） | 端到端用例＋A881 ⓒ 并集等值钉＋名册 spy |
| 3 | `.scratch/wisp/probes/313/r2/00-anchor-and-gates.md` | 本件 | 锚＋门禁＋自报 |
| 4 | `.scratch/wisp/probes/313/r2/logs/{00-head,10-gates,20-a881,20-newtests}.txt` | 5/54/41/21 | 小日志四枚 |

## 3. 完成判据（逐条现读）

- **六枚 case 各一发"落到真 setter"证据**（`Test313EmptyConfigToStepfunRowChainAndPair`，真
  `config.Manager`＋真 `ComposerDispatch`；回执键路径逐枚）：
  `provider_row_create`→`llm.providers.stepfun` · `voice_cloud_asr_chain`→`voice.cloud_asr_chain` ·
  `voice_cloud_tts_chain`→`voice.cloud_tts_chain` · `voice_realtime_provider`→`voice.realtime.provider` ·
  `voice_realtime_model`→`voice.realtime.model` · `role_chat_pair`→`llm.roles.chat.provider` **＋**
  `llm.roles.chat.model`（A881 ⓐ 两键都报）。每发回执含"重启"、⛔ 无"立即生效/保存即生效"；
  写后磁盘＋`config.get` 载荷双向复读。
- **端到端用例绿**：空配置（零 provider 行）→ 建 stepfun 行 → 写两条云链 → 写 realtime 两枚 → 写 chat pair；
  读面回执 `presetsKnown/configuredKnown` 真、`configured` false→true、行名/模型名/链/实时头全对位。
- **A881 三条件**（`logs/20-a881.txt` 逐字现读）：
  - ⓐ 第一张表与票 248 冻结钉**一字未动**：`allWritableFields` 块在 `0d87a681` 与 HEAD 间 `diff` 字节相同；
    `config_route_248_test.go` 末次改动＝`0d87a681`（313-r1 提交 `89b3b550` 对它零改动）；钉仍在 `:271`，
    七枚逐字。本腿对这两件 `git diff HEAD` 为空。
  - ⓑ 并集由单一 accessor：`acceptedWritableFields`（`config_write_313.go:104`）是唯一并集计算点，
    仅两处拒绝文案消费；`knownWritableField`＝两表查找之或；⛔ 无第三处列名册。
  - ⓒ 并集等值钉：313-r1 **未放**（其 `config_313_test.go` 只钉两表不相交），本腿补于新件
    `Test313WritableUnionEqualsTheTwoRosters`——半边＝门拒绝文案现读 == `WritableFields() ∪
    AdditionalWritableFields()`（排序全等）；半边＝并集每员过门恰达腿一次。现跑 PASS。

## 4. 门禁（rc 逐枚；`logs/10-gates.txt`）

| 门 | 命令 | rc | 备注 |
|---|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | 0 | 改前/改后各一发 |
| 格式 | `gofumpt -l` 我动过两件 | 0 | 输出空 |
| 包测试 | `go test -count=1 ./internal/panel/ ./internal/config/` | 1 | 改前=改后=同一 5 枚既有红；`internal/config` ok |
| cmd/wisp | `PATH=…sherpa-onnx:…build go test -count=1 ./cmd/wisp/` | 1 | 见下 255 段；AC14 改前红改后绿（真窗族漂色） |
| 纯树扫 | `sh scripts/d22scan.sh` | 0 | clean（含本腿新件） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 1 | RED 非本腿：`.scratch/wisp/issues/173-…-done.md`（2a8ab9c4，22:58 先于本腿改名引入）；本腿新件无一超预算 |

**255 名册具名红（A878 已有欠账，本腿落地使重裁面扩大——⛔ 别算本腿的账，但 313 收尾链/255-r2 清账时必看）**：

- `TestTicket255RosterStillMatchesTheActualReadSites`：改前红句＝voice 读面 `[resident_audio_windows.go
  internal/panel/config_handlers.go]`；改后多了 `cmd/wisp/panel_config_store.go`（落点 1 现读
  `cfg.Voice.CloudASRChain/…`）⇒ 重裁时 voice 行需补**两枚**新读者（r1 的 `config_handlers.go`＋本腿的
  `panel_config_store.go`）。
- `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`（改后新显）：`config_readers_255.go:103` 那行
  `llm` 证据钉在 `panel_config_store.go:96`，落点 1 把：96 那行换成结构体起始行、读点移到 `:98`
  ⇒ 该 receipt 行需按 255 机制重裁（行号＋引文）。
- 两枚的修面都在 `cmd/wisp/config_readers_255.go`／`config_receipt_255.go`，**不在本腿写面**（派单 §2 逐枚点名），
  故本腿只具名上报、⛔ 未动。`⛔ push 前必须清`（A878 逐字）。

## 5. ⛔ 没做完/没做的格（自报）

1. **255 名册重裁**（上节两枚）：越界未动，转 313 收尾链或 255-r2。
2. **`voice.realtime.enabled` 不可写、清空云链无界面口**：313-r1 已具名（`313-panel-ask.md` §5），本腿未自造。
3. 变异测试未用（派单"如用"）；大原始日志 `12-cmdwisp-pre.txt`（1002 行）/`43-cmdwisp-post.txt`（964 行）
   留盘不入库（同 r1 的口径）。
