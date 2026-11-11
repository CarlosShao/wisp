# 票 297 名册补腿（roster-fill）读数 — 2026-11-11 10:4x–10:5x +08

腿＝297 收口腿（CI 名册补漏）。写面＝`scripts/portable-tests.sh` 一枚名册行，零 Go 变更。

## 根因

- 票 297-l1 落的 `internal/audio/livepacket_297_windows_test.go:41` `TestAC297LivePacketByteFace`（//go:build windows，`WISP_LIVE_MIC != "1"` ⇒ 顶层 `t.Skip`，`:42-43`）。
- windows 档 scope（`scripts/portable-tests.sh:461`）含 `./internal/audio/` ⇒ 该用例编进 windows 测试二进制；CI 不设 `WISP_LIVE_MIC`（grep `.github/`＝0 命中）⇒ 顶层 `--- SKIP`；strict runner 规矩 2「any top-level --- SKIP is fatal」⇒ CI windows 档必红。
- 改前名册 grep 297＝0（无账）。⇒ 无账 SKIP 红，坐实 297-l2 的上报。

## 处置（两条路择一：选①）

- **①登记进名册（选）**：`ledger` 加一行 `TestAC297LivePacketByteFace|./internal/audio/|windows|fixture|…`，与同包同门先例 `TestLiveWasapiSmoke`（同 WISP_LIVE_MIC=1 门、同"hosted runner 无音频端点/麦克风"理由）同形。机制＝每轮 `go test -list` 对编译二进制复核存活＋锚定 `-skip` 静默过滤＋理由进 CI log；改名/删除/挪 tag ⇒ 条目变 stale ⇒ 该步红。
- ②挪 test build tag（否）：测试在 `internal/audio` 不在本腿 aecspike tag 写面；且藏掉真仪器不如记账诚实。

## 名册前后读数

- 改前：`grep -c "297" scripts/portable-tests.sh`＝**0**（rc=1）；`grep -c LIVE_MIC`＝1；ledger 行数 `grep -c '^\s*"Test'`＝**11**。
- 改后：`grep -c "297"`＝1 行（`grep -o "297" | wc -l`＝6 处）；`grep -c LIVE_MIC`＝2；ledger 行数＝**12**（12 entries, 8 active on windows）。

## 隔离探针（不经名册脚本，坐实 Test52 与本腿零关）

- 尺＝`go test -v -count=1`（strict runner 同款 13 包、无 `-skip`），逐字 `.scratch/wisp/probes/297/l3/no-script-13pkg-probe.log`，rc=0。
- 读数：Test52 家族 17 枚**自身条件性 SKIP**（门＝`internal/plugin/command_52_test.go:41`「t.TempDir() lives under a TEMP root … the D46 stability fixtures need a stable tree」与 `:263-265` TEMP 缺省门）；同探针里 `TestAC297LivePacketByteFace` 打出真顶层 `--- SKIP`（改前 CI 红的形状复现）＋`TestLiveWasapiSmoke`／`TestSyncRegistryProbeLive`／`TestHelperProcess` 各 1 枚（即名册在记账的那些）。
- ⇒ 名册脚本整档跑下 Test52 两枚子测试随机红＝**D46 fixture 在本机全档并发下的预存 flake**（单跑 ok 7.985s／5.461s；无脚本时整族自 SKIP），与本腿一行名册零关（该行只改 `-skip` 模式，且被跳名与 Test52 无交集）；CI 判定归 CI 真跑，本腿不碰 `internal/plugin`（越写面）。

## 验证（整档实跑，非推演）

- 尺＝`bash scripts/portable-tests.sh --scope=windows`（本机 windows 档真跑，13 包）。
- 运行日志自证：`portable-tests.sh: 12 ledger entries, 8 accounted on this platform` ＋新行理由全文打印；`-skip pattern` 含 `TestAC297LivePacketByteFace`；全档 `-v` 流里 **TestAC297LivePacketByteFace 零 `=== RUN`／零 `--- SKIP` 行**（被在册 -skip 静默过滤，即 l2 报的"无账 SKIP 红"的根已除）；同包另两枚 AC297 用例（KnownAmplitudePair / ConvertPacketSynthetic）真跑 PASS；`internal/audio ok 22.766s`。
- 每轮存活性探针：`go test -list '^TestAC297LivePacketByteFace$' ./internal/audio/` ＝ 列出该名（windows 主机现跑）⇒ 条目不会变 stale。
- 完整档复跑（CI 形状：TEMP 在、GOTMPDIR 不设）逐字节日志＝`probes/297/l2-roster-verify-fullwindows.log`：四数字 `=== RUN=1264 --- PASS=818 --- FAIL=2 --- SKIP=0`、rc=1（FAIL=2＝上节 Test52 预存 flake，SKIP=0＝无任何无账 SKIP）；第一发（TEMP unset 环境件）四数字 `RUN=1264 PASS=817 FAIL=2 SKIP=1`（SKIP=1＝Test52 的 TEMP 门）。

## 环境件归因（诚实记，⛔ 不是树缺陷、⛔ 与名册行无关）

第一发整档跑按铁律 `unset TEMP TMP` ⇒ `internal/plugin` Test52 家族 3 枚红/SKIP（`Test52SchemaViolationsAreRejected`、`Test52CommandKindValidates` ＋ SKIP `Test52AbsoluteTempPathIsRefusedForTheTempReason`，command_52_test.go 依赖 TEMP）。A/B 对照：TEMP 在 ⇒ 整包 `ok 5.461s`；TEMP 无 ⇒ 该 3 枚红/SKIP。CI 恒有 TEMP ⇒ 不出现在 CI。plugin 包单跑（TEMP 在）＝`ok 7.985s`。

## 门禁读数

- `bash -n scripts/portable-tests.sh` rc=0。
- `sh scripts/d22scan.sh` rc=0 clean（462 文件全量）。
- `sh scripts/check-path-length-budget.sh` rc=0 VERDICT GREEN。
- go 变更＝0 ⇒ go vet/gofumpt 不适用（写面纯 shell 一行名册）。
- `wisp.exe` 起手 tasklist＝0。

## 未做到／边界

- CI 上的最终确认只能由编排者推送后看 CI 真跑（本腿零 push）；本机读数＝同一脚本同一路径的整档实跑。
- 本腿未勾票 297 任何 AC 框（AC#0..#4 归非实现者验收腿）。
