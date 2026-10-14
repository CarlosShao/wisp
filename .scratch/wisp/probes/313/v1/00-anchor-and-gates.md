# 313-v1 锚与门禁（验收腿；非实现者）

- **锚**：`bfb62488`（2026-10-14 10:07:35 +08 终跑；验收窗口 09:44–10:07）。被裁＝`89b3b550`（r1，internal/panel＋internal/config）／`2f8964b9`＋`9e3013e4`（r2，cmd/wisp）；页面侧 `9e9964f8` **只读**。
- **树呼吸记录**：窗口内 HEAD 由 `159a416e` → `d8a4d913` → `f3dc1389` → `bfb62488`（237-r1 两笔＋编排者 A904）；与 `internal/panel`／`internal/config`／`cmd/wisp/panel_config_store.go` **零文件重叠**。
- **载具**：`zz313v1probe_test.go` 经 `overlay-probe.json` 注入 `cmd/wisp/`（`go test -overlay`），装配与 `run.go:431`／`panel_inbound.go:279` 同形：真 `config.NewManager`＋真 `secret.NewStore`＋`newConfigStore`＋真 `panel.ComposerDispatch`；⛔ 无 mock 顶替被裁对象。三方键对拉＝`compare-keys-313v1.py`（Go 载荷现跑／`313-panel-ask.md` §1 表／`config-view.ts` 接口，机械抽取）。

## 门禁 rc（现跑，log 见 `logs/10-run-and-keys.txt`、`logs/20-mutations-and-ledger.txt`）

| 尺 | 命令 | 结果 |
|---|---|---|
| 构建 | `GOFLAGS= go build ./...` | rc=0 |
| 格式 | `gofumpt -l` 七枚被裁件 | 空（rc=0） |
| d22scan | `sh scripts/d22scan.sh` | clean rc=0 |
| 长度门 | `sh scripts/check-path-length-budget.sh` | VERDICT GREEN；58=58 rc=0 |
| 三包 | `-run 'Test313' ./cmd/wisp/ ./internal/panel/ ./internal/config/` | 三包全 ok（终跑 10:07） |
| 探针 | `-overlay … -run 'TestV1Probe' ./cmd/wisp/` | ok（四发全绿，终跑 10:07） |
| 248 哨兵钉 | `-run 'TestAC2CredentialSentinelAppearsInNoArtifact'` | PASS（亲跑） |
| 255 名册 | `-run 'TestTicket255' ./cmd/wisp/` | 全 PASS（r2 记的两处欠账已由 255-r7 `dcaff922` 清） |

## 突变窗口（真实文件编辑；仓外备份＋sha256 还原证）

| 突变 | 改件 | 窗口起/止 | 还原 sha256 |
|---|---|---|---|
| B1 并集截一员 | `internal/panel/config_write_313.go` | 09:55:05–09:55:28 | `932bd7a9…d8e4` 逐字节回 |
| B2 表删一员 | 同上 | 09:56:14–09:56:28 | `932bd7a9…d8e4` 逐字节回 |
| A case 换脚 | `cmd/wisp/panel_config_store.go` | 09:57:50–09:58:02 | `43acff43…9a9d` 逐字节回 |
| B 分隔符改坏 | 同上 | 09:59:07–09:59:28 | `43acff43…9a9d` 逐字节回 |

四发还原后均复跑绿（logs 内逐条）。备份在仓外 `D:/tmp/wisp-313-v1/backup/`。

## 具名没做的格

- ⛔ 整包 `./cmd/wisp/`（派单：314 死区 157–215s 进程死）；一律 `-run` 小射程。
- ⛔ `-race`（派单未要求；被裁件非并发面）。
- ⛔ 页面侧渲染/交互判语（机主前端会话面；本腿只做只读读键对拉）。
- ⛔ 未 push；⛔ 未动 `Status:`／票名／AC 框。
