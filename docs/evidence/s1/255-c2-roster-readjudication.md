# 票 255 · C2 名册重判（写腿 `255-c2r1` 交件）

**日期**：2026-11-05 12:2x–12:4x +08
**锚**：HEAD `05e02c4df`（工作树起手 `cmd/wisp`／`internal/config` porcelain 只有一枚未跟踪 `zz_probe_temp_test.go.txt`；期间另一条 133 腿在写 `cmd/wisp/leg_dispatch_gate_133_test.go`／`dispatch_hop_selfproof_133_test.go`，本腿一个字节未碰）
**来路**：`docs/evidence/s1/98-attribution.md` §2 **C2**（ubuntu 81 红里的那一枚「两条腿都红」）＋ 台账 `A1538`「98 普查撞出 C2 真 bug ⇒ 建议 255/133/119 各开腿」。
**性质**：**实现腿交件**。本文件是唯一新增件；判语要**非实现者**另起一发（裁决者 ≠ 实现者）。

## 0. 一句话

`cmd/wisp/config_readers_255.go` 的名册把第 44／61 票落地后新长出来的生产读者判丢了：**3 枚 `match=false`**（`cost`／`llm`／`voice`）在**干净 HEAD 快照上两条腿同红**（票 98 §4 在 windows 干净快照上实测逐字同三条），所以 `scripts/wisp-cli-tests.sh`（CI `test-windows` 步）在这棵树上**不可能绿**。本腿按 44／61 的真实读点**重判名册**，14/14 段 `match=true`，并把重判钉进测试（判据留在数据里，不靠下一波人记得）。

## 1. RED-before（改前，本机现跑，逐字）

```
$ export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
$ go test ./cmd/wisp/ -run "TestTicket255RosterStillMatchesTheActualReadSites" -count=1 -v
SECTION 5/14  section="cost"  field=Cost  roster=[] scan=[internal/agent/costbook_44.go internal/agent/quota_44r4.go] match=false
SECTION 7/14  section="llm"   field=LLM   roster=[cmd/wisp/panel_config_store.go cmd/wisp/panel_models_320.go cmd/wisp/providers.go cmd/wisp/run.go internal/llm/resolver.go] scan=[同上五枚 + internal/agent/costbook_44.go internal/agent/quota_44r4.go internal/speech/cloudcatalog.go] match=false
SECTION 14/14 section="voice" field=Voice roster=[cmd/wisp/panel_config_store.go cmd/wisp/resident_audio_windows.go cmd/wisp/resident_voice_windows.go internal/panel/config_handlers.go] scan=[同上四枚 + internal/speech/cloudcatalog.go] match=false
--- FAIL: TestTicket255RosterStillMatchesTheActualReadSites (0.92s)
rc=1
```

三条红句逐字同票 98 §2 C2 的描述（`cost`／`llm`／`voice`）。

## 2. 重判的三件事（每件都有可核读点）

| 段 | 旧判 | 新判 | 凭据（本腿现读） |
|---|---|---|---|
| `[cost]` | `no-reader: 扫描零命中 … the C23 budget knobs (ticket 44) are not landed` | **`consumed:`** | 44 票的装配早就落了：`execute()` 每任务现读活配置（`cmd/wisp/run.go:1524 [cfg = rt.mgr.Config()]`），把同一枚 `cfg` 交给 `agent.CostWiring`（`run.go:1570`），后者用 `cfg.Cost.DailyBudget/MonthlyBudget/AlertThreshold/OverBudget` 建 budget（`internal/agent/costbook_44.go:183`）；`agent.QuotaWiring`（`internal/agent/quota_44r4.go:135`）读同样四枚 cap。⇒ 热加载改 `[cost]` ⇒ **下一个任务**的 meter／quota gate 按新值走，无需重启（`cmd/wisp/costwiring_44r3_test.go:118` 早就把这条写成"the live-config half"）。 |
| `[llm]` | `consumed:`（读数不变） | **`consumed:`（不变）** | 活的读者仍是 settings leg（`cmd/wisp/panel_config_store.go:98`）＋ 320 fetch 口；44／61 的三枚是"**别人把 cfg 递给它**"的读点，按 255 AC#1 的口径不改变这一段"是否立即生效"的判语，只进 `sectionReadSites` 的钉。 |
| `[voice]` | 四枚 hot 键全 `no-reader` | **`no-reader`（其中 `voice.tts.speed` 重写证据）** | 61 票的 `internal/speech/cloudcatalog.go` 确实读 `cfg.Voice.*`（⑤ 进 `sectionReadSites["Voice"]`），其中 `:119 [ep.Speed = cfg.Voice.TTS.Speed]` 读的正是 hot 键 `voice.tts.speed`。**但**那枚读点住在 `ResolveCloudChain` 里，而该函数全仓**非测试调用者 0 枚**（`grep -rn "ResolveCloudChain" --include=*.go cmd internal tools` 只命中定义与注释）⇒ 「读点存在、没人调用」，所以这行**留 `no-reader`**，但把过期的「扫描零命中」换成具名点名那枚死读点。这是 C2 同一病的第二层：不重判这一行，名册就会在文件内部自相矛盾（`sectionReadSites["Voice"]` 收了 cloudcatalog.go，而行文说"没人读"）。 |

**边界（明写不做）**：`voice.punctuation`／`voice.wake_word.thresholds`／`voice.wake_word.veto_words` 三行原样不动（那三枚键没有任何读点）；`internal/panel/config_handlers.go` 的 `v.Voice.*` 是 `panel.SettingsView` 镜像，按 313 的既有裁定继续算读者（本腿不重开那一格）。

## 3. 改了哪三件

1. `cmd/wisp/config_readers_255.go`
   - `hotRowClaims["cost"]`：`no-reader` → **`consumed:`**，cite 为内容锚 `internal/agent/costbook_44.go:183 [b, err := observe.BudgetFrom(cfg.Cost.DailyBudget, cfg.Cost.MonthlyBudget,]`（token 在该文件**唯一命中 1 次**，行号只当读数——票 255 AC#6）。
   - `hotRowClaims["llm"]`：行文补三枚"被递 cfg"的读点（44 两枚＋61 一枚），**判语不动**。
   - `hotRowClaims["voice.tts.speed"]`：换掉「扫描零命中」，点名 `internal/speech/cloudcatalog.go:119 [ep.Speed = cfg.Voice.TTS.Speed]` 这枚死读点，**判语仍 `no-reader`**（无运行路径按它做事）。
   - `sectionReadSites`：`Cost: nil` → `{internal/agent/costbook_44.go, internal/agent/quota_44r4.go}`；`LLM` 追加 3 枚（44×2＋61）；`Voice` 追加 `internal/speech/cloudcatalog.go`；三处各带一段"为什么"的注释。
   - 四枚 `[voice]` 键行里过期的「the four files sectionReadSites pins」→「the files …」（名单从 4 变 5，硬写数字会立刻再烂）。
2. `cmd/wisp/config_receipt_255_test.go`
   - `TestTicket255SplitOnlyClaimsSectionsWithALiveReader`：`claimable` 期望 `"llm agent"` → **`"llm agent cost"`**（cost 现在真有活读者；**不是放宽**——枚数从 2 变 3 但每一枚仍必须由名册的 `consumed` 支撑）；`quiet` 计数 `len(hot)-2` → `len(hot)-len(claimed)`。
   - `TestTicket255ReceiptSentenceAssemblyIsFiltered`：把 `[cost]` 加进"必须在「已立即生效」句里"的正控，并加进"不许同时出现在诚实句里"的互斥检查。
   - **新立** `TestTicket255RosterReJudgedAfterTickets44And61`：钉住两条判断（`cost`=consumed＋cite、`voice.tts.speed`=no-reader 但点名死读点、`sectionReadSites["Cost"]` 精确两枚、`splitHotTier([cost]).claimable=[cost]`、`splitHotTier([voice]).claimable=[]`）。
   - 顺带把 `citedRowsFloor` 上方那段"13 of 17 rows / 三枚 marker-only"的过期计数改成现读（12 of 17，marker-only 是 `memory`／`app.theme`／三枚 `[voice]` 键行）。
3. `docs/evidence/s1/255-c2-roster-readjudication.md`（本文件）。

**没碰**：`docs/PLAN.md`、`docs/specs/**`、`internal/config/**`、`thresholds.go`、golden、`cmd/wisp` 整包、`frontend/**`、`design/**`、票 44／61 的夹具文件、票 255 的票面（`-done` 文件不在本腿写面；翻框归编排者）。

## 4. GREEN-after（改后，本机现跑）

```
$ go test ./cmd/wisp/ -run "TestTicket255" -count=1 -timeout 200s
ok  github.com/CarlosShao/wisp/cmd/wisp  10.470s          rc=0

$ go test ./cmd/wisp/ -run "TestTicket255RosterStillMatchesTheActualReadSites|TestTicket255RosterReJudgedAfterTickets44And61" -count=1 -v
SECTION 5/14  section="cost"  roster=[internal/agent/costbook_44.go internal/agent/quota_44r4.go] scan=[同一对] match=true
SECTION 7/14  section="llm"   roster=[8 枚] scan=[8 枚] match=true
SECTION 14/14 section="voice" roster=[5 枚] scan=[5 枚] match=true
--- PASS: TestTicket255RosterStillMatchesTheActualReadSites (1.11s)
--- PASS: TestTicket255RosterReJudgedAfterTickets44And61 (0.00s)
rc=0
```

**产品面读数（票 255 AC#1 的正控/反控仍成立，且多了一枚真读者）**：

```
IMMEDIATE "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[llm] [agent] [cost]"
HONEST    "…（票 255 AC#1：这一半不许说成「已立即生效」…）：[ball] [session] [audio] [privacy] [memory] [panel] [models] [observe] [hotkey] [app] [voice]"
```

`[panel]`／`[ball]` 仍在诚实句（正控不变）；`[llm]` 仍在立即句（反控不变）；`[cost]` 是新增的第三枚立即档，凭据是上面那枚每任务活读。

**内容锚复跑**（`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`）：本腿新增的两枚 token 都**唯一命中**（`cfg.Cost.DailyBudget` 家族 1 次、`ep.Speed = cfg.Voice.TTS.Speed` 1 次）；另有 15 条**既有**读数漂移（`agent :1263→:1524`、`panel :207→:208` 等）被同一把尺以 `t.Logf` 记下——那是 AC#6 明确允许的"行号只是读数"，本腿**没有顺手刷新别人的读数**（不在本腿射程，也不该拿它充数）。

## 5. 门禁（每件自落 rc，⛔ 无管道吞退码）

| 门 | 命令 | rc／读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | **0** |
| vet | `go vet ./cmd/wisp/` | **0** |
| 本票全尺 | `go test ./cmd/wisp/ -run "TestTicket255" -count=1 -timeout 200s` | **0**（`ok 10.470s`） |
| 相邻配置面 | `go test ./cmd/wisp/ -run "TestTicket255\|TestTicket324\|TestTicket223\|TestTicket233\|TestTicket258\|TestTicket64" -count=1 -timeout 200s` | **0**（`ok 46.760s`） |
| gofumpt | `gofumpt -l <本腿两枚 .go>` | **空** |
| d22scan | `sh scripts/d22scan.sh` | **0**（`d22scan: clean`；`runtests.sh: OK … PASS=35 FAIL=0 SKIP=0`） |

**两枚具名例外（都不是本腿的账，⛔ 不许读成"绿"）**：

1. `go build ./...` 这一形 rc=1，红在 `build/mut_rvw.go`（`undefined: residentBall` 等）——`build/` 在 `.gitignore:14` 里，是**别人留下的变异体临时件**（本腿起手就在，mtime 11-04），不进仓、不进 CI；本腿用上面那条按真实包树的 build 尺。
2. `gofumpt -l cmd/wisp` 报 `leg_dispatch_gate_133_test.go`——本腿 12:33 改完时它是干净的，12:36:11 那枚 mtime 是**在飞的 133 腿**刚写的；`git status` 里它与 `dispatch_hop_selfproof_133_test.go` 同属 133 腿的改动，本腿一个字节未碰。
3. 整包 `go test ./cmd/wisp/ -count=1 -timeout 200s` rc=1，读数是 `panic: test timed out after 3m20s`、**`--- FAIL` 行数 = 0**；这是本包既有的"整包死区"（255-r1 已量过同一形状：8 分钟帽 panic），与本次三件改动无关，本腿**不主张整包绿**。

## 6. 复跑口令

```bash
cd "/d/work/workspace/projects plans/Wisp"
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
go test ./cmd/wisp/ -run "TestTicket255" -count=1 -timeout 200s
go test ./cmd/wisp/ -run "TestTicket255RosterStillMatchesTheActualReadSites|TestTicket255RosterReJudgedAfterTickets44And61" -count=1 -v
/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/config_readers_255.go cmd/wisp/config_receipt_255_test.go
sh scripts/d22scan.sh
```

## 7. 没建立的事（宁缺毋造）

- **没有在 ubuntu 腿复跑**：C2 是两条腿同红（票 98 §4 已实测），但本腿只在本机 Windows 上取绿；`config_readers_255.go` 与这枚测试都没有 build tag，扫描面同样覆盖 `internal/**`，所以 linux 上应当同绿——**这是推理，不是读数**，要 CI 那一发才算。
- **没有跑整包 `cmd/wisp`**（见 §5 第 3 条，200s 帽内必超时，且是在飞 133 腿的共享编译面）。
- **没有改票面、没翻 AC 框、没写台账**（不在本腿写面；`A1538` 的收口归编排者）。
- **没有重判 44／61 票自己的行**：那两票的读者是"生产事实"，本腿只改它们在 255 名册里的登记；票 44／61 的测试文件一个字节未动。
