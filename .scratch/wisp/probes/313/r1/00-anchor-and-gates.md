# 313-r1 锚＋门禁（界面配模型·Go 侧）

判语：三格 Go 半（甲读面／乙写面／丙读面半）落 **5 枚产码件＋3 枚交付件**，门禁 6 枚全跑，**整包新增红 0 枚**（现树失败 ⊂ 纯净树基线）。
尺：`git log --oneline -1 -- internal/panel/config_handlers.go` · `wc -c .scratch/wisp/probes/313/r1/logs/*`

## 锚

- 起手 `git log -1` 逐字＝`logs/00-head.txt`（**HEAD `39f6341d`**，2026-10-13 18:09）。⚠ 本腿起手时 HEAD 已从 `b54f1740` 漂到 `39f6341d`（别的腿在提交，共享工作树）。
- 产码笔＝`89b3b550`（5 枚件，见下）。⛔ 零 push。

## 交件（产码 5 枚＋交付件 3 枚；件名短、⛔ 带 `§`）

| # | 件 | 行数 | 是什么 |
|---|---|---|---|
| 1 | `internal/panel/config_handlers.go` | +159/−7 | 读面：`SettingsView` 追加字段＋`VoiceSettingsView`＋`ProviderSettings.Models`＋`renderSettingsPayload`＋标记；门：两表并查、三列取并 |
| 2 | `internal/panel/config_write_313.go` | 123（新） | 第二张可写表六枚（`A877` ①②）＋并集/选择器助手 |
| 3 | `internal/panel/config_313_test.go` | 324（新） | 6 枚用例（含票 248 那枚钉的复述式正控、载荷键集钉、命名禁词/锁定族） |
| 4 | `internal/config/settings_313.go` | 222（新） | 六枚 setter：显式建行／云链两枚／realtime 两枚／chat pair |
| 5 | `internal/config/settings_313_test.go` | 273（新） | 5 枚用例（建成/拒绝/云链/成对/realtime 半对），每枚拒因都断言"文件逐字节不变" |
| 6 | `docs/evidence/s1/313-panel-ask.md` | 74 | 页面接口规格（转交机主 agent） |
| 7 | `.scratch/wisp/probes/313/r1/10-cmdwisp-landing-package.md` | 93 | 接续腿的精确落点包（`cmd/wisp`，本腿 ⛔ 未动） |
| 8 | `.scratch/wisp/probes/313/r1/00-anchor-and-gates.md` | 本件 | 锚＋门禁＋没做的格 |

## 门禁（每件自落 rc）

| 件 | 命令 | rc | 备注 |
|---|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | **0** | 树任何时刻可编译；每小步后小射程 `go build ./internal/config/ ./internal/panel/` 也各 rc=0 |
| 格式 | `$(go env GOPATH)/bin/gofumpt.exe -l internal/config/ internal/panel/` | **0** | 输出空；`settings_313.go`＋`config_handlers.go` 两枚曾红，已 `-w`（只动我加的块） |
| 配置包 | `go test -count=1 -v ./internal/config/` | **0** | `ok`；313 五枚用例 PASS，整包绿 |
| 面板包 | `go test -count=1 -v ./internal/panel/` | **1** | 见下：**6 枚失败全部是既存红**，313 六枚用例 PASS |
| 纯树扫描 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；internal/ 235 枚＋cmd/ 40 枚＋ban#8 541 枚全过（`logs/40-d22scan.txt`） |
| 台面 | `tasklist` 现跑 `wisp.exe`／`balldebug.exe` | **0 命中** | 读数前跑；`GOTMPDIR=C:\tmp\wisp313gotmp`（只 C: 卷） |

### 整包"新增红 0"的尺（改前/改后取交集）

- **改后（本腿工作树）**：`internal/panel` 失败 **6 枚**——`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／`TestComposerContractTypesMatchFrontend`／`TestNoGitSwitchCapabilityInThePanelSurface`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`。
- **改前（纯净树基线）**：`git archive HEAD`（`39f6341d`）导出**仓外** `/c/tmp/wisp313base` 跑同一把 ⇒ 失败 **8 枚**（上面 6 枚 ＋ `TestComposerRenderFixtureTellsTheTruth` ＋ `TestReadGitOnThisRepositoryIsSelfConsistent`）。
- ⇒ 现树失败集 **⊂** 基线失败集 ⇒ **新增红 0 枚**。六枚红的性质（逐名）：`frontend/`↔Go 类型漂移（票 197 那一族）· `design/assets/tokens.css` 在工作树里被删（脏件）· `frontend/git-branch.tsx` 存在（票 186 在飞）· 背压上界两枚＝票 35r8 既存红。**均与本票面无关**。
- 大日志留盘不入库：`logs/31-two-packages-verbose.txt` = **107,011 字节**（`-v` 原始，231 枚 `^--- PASS`、0 枚 `^--- SKIP`）。

## ⛔ 本腿没做／具名上报

1. **`cmd/wisp` 装配那一跳**（`ReadSettings` 填字段 ＋ switch 六 case）——派单补记 §1 明令该面归 `294-r1`，本腿 ⛔ 一字未动 ⇒ 精确落点包见 `10-cmdwisp-landing-package.md`。在它落地前，`config.get` 载荷恒为零值（"腿没答"）、六枚新写字段会被腿按"没有对应的写入腿"拒（**响亮拒，非静默成功**）。
2. **`role_chat_pair` 是本票唯一"一字段写两键"**（`roles.chat` 两种半态都被 `validateRoles` 拒，β 形结构上不存在）——与 `A877` ① 的 one-key-per-call 纪律的关系**具名待裁**：本腿判"无它则 `AC#2` 文本模型这一格界面到不了"，故落之并上报。
3. **清空云链无界面口**、**`voice.realtime.enabled` 不可写**（`A877` ① 名册未含）——两枚缺口已写进 `313-panel-ask.md` §5，⛔ 未自造 sentinel/新字段。
4. **`WritableFields()` 一字未动**（票 248 冻结钉 `config_route_248_test.go:271-277` 把 `WritableFields()` 钉在那 7 枚上）⇒ 313 六枚放**第二张表**、门并查两表。这是"只许追加／⛔ 改任何既有断言"下的唯一形，具名在此。
