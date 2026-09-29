# 票 223 非实现者对抗验收 `223-v1` —— CheckAndReload 热加载接线

- 工单：`.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md`
- 被裁的实现件：`docs/evidence/s1/223-hot-reload-wiring-r1.md`（实现腿 `223-r1` 死于 150 轮上限；AC#5 由编排者代跑；五枚文件未提交改动由编排者代落＝commit `248095d1`）
- 落点普查前案：`.scratch/wisp/probes/223/c1/census.md`
- 跨票前案：`.scratch/wisp/issues/226-always-answer-rewrites-the-whole-config-snapshot-and-hides-hand-edits.md`／`docs/evidence/s1/226-config-write-no-clobber-v2.md`／`internal/config/writeguard.go`
- 本腿起手锚点：`git rev-parse --short HEAD` = **`577ae8c7`**，分支 `dev`，`date` = 见下方起手时刻表
- 台件目录：`.scratch/wisp/probes/223/v1/`（原始输出全量落盘，绝不接 `| head`／`| tail`）
- 本腿性质：**只裁、只跑台件、不产码、不 commit 任何源码改动**
- ⛔ 零读零写：`frontend/**`、`design/**`（本腿对这两目录不读、不引、不转述）
- ⛔ 不碰：`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件

## 入库清单（本腿会 commit 的路径，逐枚具名）

| 路径 | 是什么 |
|---|---|
| `docs/evidence/s1/223-hot-reload-wiring-v1.md` | 本表 |
| `.scratch/wisp/probes/223/v1/` | 全部原始读数与台件（含 `full-v.txt`、突变日志、loader 对抗输入台件） |

---

## 起手读数（本腿现跑，不复用他人数字）

- 起手时刻：`date` = **Tue Sep 29 15:44:04 CST 2026**；锚点 `git rev-parse --short HEAD` = **`577ae8c7`**（骨架提交后＝`4fc03f21`，本腿所有读数都在 `4fc03f21` 这枚含代落产码的锚点上跑）。
- 在飞检查：`ps -W | grep -i -E "go\.exe|go-test|test\.exe"` 起手＝**空**；`git diff --cached --name-only` 起手＝**空**（索引干净）。
- 别人地界的脏改动本腿一律不碰：`.gitignore`、`design/**`（16 枚删除＋若干新增）、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`docs/evidence/s1/152-*.md`（都在起手 `git status --porcelain` 里现读到，未读其内容）。

| 尺（本腿现跑，原文见 `.scratch/wisp/probes/223/v1/`） | 读数 |
|---|---|
| `grep -rn "CheckAndReload()" --include=*.go cmd/ internal/ \| grep -v _test.go` | **2 枚生产调用者**：`cmd/balldebug/main.go:243` ＋ **`cmd/wisp/config_reload.go:153`**；另 `internal/ball/hotkey_reload.go:21` 是注释。（实现腿表里写 `:150`，现读 **153** ⇒ 行号漂移，见"票面写的 vs 我读到的"） |
| `grep -rn "ConfirmLocked = " --include=*.go cmd/ internal/ \| grep -v _test.go` | **1 枚**＝`cmd/wisp/config_reload.go:114`（实现腿写 `:111`，现读 114）；含测试则 11 枚命中，全在 `*_test.go` 与本腿台件目录 |
| `grep -rn "OnRestartPending = " --include=*.go cmd/ internal/ \| grep -v _test.go` | **1 枚**＝`cmd/wisp/config_reload.go:115`（实现腿写 `:112`）；测试命中 3 枚（`manager_test.go:82/94`） |
| `grep -rn "OnReload = " --include=*.go cmd/ internal/ \| grep -v _test.go` | **1 枚**＝`cmd/balldebug/main.go:244` ⇒ **`wisp run` 仍没有 reload 档的生产赋值点**（与实现腿"没做完"第 3 条一致，现读复认） |
| `grep -rnE "^[[:space:]]*go [a-zA-Z_]" --include=*.go cmd/ internal/ \| grep -v _test.go`（ban #1 连具名 `go f()` 一起扫） | **1 枚**＝`internal/observe/goroutine.go:281`，正是 d22scan 按文件路径豁免的那枚注册表实现 ⇒ **本票新起的 tick 不是裸协程** |
| `grep -n "time.Now\|\.Sub(\|time.Since\|Unix(" cmd/wisp/config_reload.go` | 只命中 `:37` 那一行**注释**（`//` 开头）。尺的判定：`tools/d22scan/main.go:759-770` 对每行 `TrimSpace` 后 `HasPrefix(code,"//")` 就 `continue` ⇒ 注释行不进 `wallclockRe`；非注释行**零命中** |
| `grep -n "time.Now\|\.Sub(\|time.Since" internal/config/manager.go` | **零命中** |

## AC#1 生产里真有人在轮询

〔待填〕

## AC#2 三档生效级别各有读数

〔待填〕

## AC#3 放宽必带 L2 复确认（D33 正控）

〔待填〕

## AC#4 不生效与读不到是四句话

〔待填〕

## AC#5 整包终态读数（任务一：带 `-v`）—— **不成立（终态有一枚本票自己的红，且它不是争用）**

- 尺（逐字）：先 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，再
  `go test ./cmd/wisp ./internal/... -count=1 -v -timeout 30m`。
  **起手 `date`＝`Tue Sep 29 15:46:22 CST 2026`／终态 `date`＝`Tue Sep 29 15:48:28 CST 2026`**（两个时刻都写进了原始文件首尾）。
- 原始输出全量落盘＝`.scratch/wisp/probes/223/v1/full-v.txt`（**7,057 行／841,116 字节，未接 `| head`／`| tail`**），起手锚点 `4fc03f21`（＝`577ae8c7`＋本腿的表；产码与 `248095d1` 逐字相同），`GATE_EXIT=1`。
- **起手前先自证没有并发门**：`ps -W | grep go.exe|test.exe` 空；跑期间该机只有这一发。
- 口径（**引用数字必带**）：顶层 `=== RUN` **1737**；顶层 `--- PASS` **1158**；顶层 `--- FAIL` **7**；`--- SKIP`（含缩进层）**7**；子测试 `--- PASS` **565**；子测试 `--- FAIL` **0**。
  包级：**26 枚包＝22 `ok`／4 `FAIL`**（`cmd/wisp`、`internal/ball`、`internal/panel`、`internal/risk`）。
- **红名册（用例名集合，7 枚，不用枚数）**：

| 用例名 | 包 | 归谁 |
|---|---|---|
| `TestC21TableColourRowsMatchTokensCSS` | internal/ball | 历史在册（别人地界） |
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | internal/panel | 历史在册 |
| `TestComposerContractTypesMatchFrontend` | internal/panel | 历史在册 |
| `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | internal/panel | 历史在册 |
| `TestC21DesignTokensFourwayAgree` | internal/panel | 历史在册 |
| `TestResolvePerCallBudget` | internal/risk | 争用型假红（本腿复量见下） |
| **`TestTicket223RestartTierSaysItWillNotApply`** | **cmd/wisp** | **本票新增＝AC#7 那一枚用例自己** |

- 本票 12 枚具名用例在这发里的逐名读数：`TestTicket223RunArmsTheReloadTick` PASS 2.22s／`…HandEditedFsLooseningCostsAnL2Card` PASS 2.22s／`…RefusedLooseningKeepsOldValues` PASS 5.30s／`…TighteningRaisesNoCard` PASS 2.25s／`…ModeLooseningChangesTheRunningModeAfterAllow` PASS 2.39s／**`…RestartTierSaysItWillNotApply` FAIL 2.10s**／`…FailureSentencesAreDistinct` PASS 8.62s（4 枚子例全绿）／`…PanelInboundSaysHotReloadIsDisabled` PASS 0.01s／`…PermissionDeniedSitsInItsOwnSentence` PASS 2.15s（真 ACL 这台上跑通）／`TestConfirmLockedRunsOutsideTheManagerLock` PASS 0.03s／`TestCheckAndReloadAsksOncePerFileChange` PASS 0.04s。

- **那一枚红的复量＝它不是争用**（尺：`go test ./cmd/wisp -count=5 -v -run TestTicket223RestartTierSaysItWillNotApply`，单包、安静机、**15:50:33→15:50:45**，落 `.scratch/wisp/probes/223/v1/restart-tier-recheck.txt`）：
  **4 PASS／1 FAIL**（FAIL 在 `config_reload_223_test.go:453`，内容逐字＝"the operator is not told the edit will not land this run"，且 dump 出来的 stdout 里只有"答复监听已接入"＋"配置热加载已接管"两行，**没有重启档那句**）。
  ⇒ 同一枚用例在**没有别的门在飞**时仍然 5 发里红 1 发 ⇒ **争用解释被排除，这是用例自身的读序竞态**：`cmd/wisp/config_reload.go:282` 与 `:284` 两行审计**先**写、`:288` 的 stdout 句子**后**写，而用例 `config_reload_223_test.go:438/:442` 用 `awaitAudit` 轮询 stderr 到第二行就返回，接着 `:452` 对 stdout **只读一次、不轮询**（同文件里别的用例读 stdout 都走 `awaitLive`/`awaitAudit` 之后再读）⇒ 存在真实窗口。**本腿不修它**（禁区＋非实现者腿不产码），登记为 AC#7 的注与"该翻勾前必须补的一发"。
- **`TestResolvePerCallBudget` 的处置照派单走完**（尺：`go test ./internal/risk -count=3 -v -run 'TestResolvePerCallBudget'`，15:52:20→15:52:25，落 `budget-recheck.txt`）：
  **3/3 PASS，`ok github.com/CarlosShao/wisp/internal/risk 3.839s`**，逐发读数 `248908/252504/246368 ns/op`（预算 1.000 ms、约 4900 样本）。⇒ 这枚是争用型假红，**两发读数都写进表**：本腿 `3.839s`；编排者 `3.804s`／`3.813s`。**`thresholds.go` 与 golden 一字节未动（本腿全程未打开过这两枚文件）**。

- **复核编排者那一发（AC#5 里"由编排者代跑补齐"那节）对不对**：
  ① 它说 **23 包 ok／3 包 FAIL 共 6 例** ⇒ 我这发是 **22 ok／4 FAIL 共 7 例**。**包计数与例数都不同**，差异可由**一枚**用例完全解释＝`TestTicket223RestartTierSaysItWillNotApply`：它红 ⇒ `cmd/wisp` 从 `ok` 变 `FAIL`（包计数 23→22、FAIL 包 3→4），例数 6→7。**编排者那一发的读数在它自己那一刻是自洽的**（那枚 flaky 没命中），但"6 例"不是稳定终态。
  ② 它说红名册＝5 枚历史红＋`TestResolvePerCallBudget` ⇒ 与我这发的集合**逐名相符**（那 5 枚名字一字不差）。
  ③ 它说复量两发 `ok 3.804s`／`ok 3.813s` ⇒ 与我 `3.839s` 同向，**归因（争用）复认**。
  ④ **缺 `-v` 会不会改变结论**：**不会改变红名册**（Go 不带 `-v` 也打 `--- FAIL`，它那发的 6 例逐名可读、我已复认），**但结论本身要改**——它缺的正是"逐名绿册"（我这发才采到 1158 顶层＋565 子测试 PASS），而且**带 `-v` 的终态多抓出一枚本票自己的红**，那枚红在它那发不存在。⇒ **它那发不足以定"终态干净"，这一条它自己已在表里具名承认，本腿复认它的诚实、但读数口径要按我这发替换。**
- **"拿掉实现会不会红"**：会，且不是靠推断＝编排者/实现腿的 m1（注释掉三枚接线）与**本腿任务二那发锁内突变**都会让 `cmd/wisp`／`internal/config` 转红（见 AC#6）。本格自身的不成立点是**用例有竞态**，不是"没有实现"。

## AC#4 不生效与读不到是四句话 —— **成立但带注（注＝有一句会说反，见推翻清单第 4 条）**

- 台件＝**`go test -overlay`**（物理件全在 `.scratch/wisp/probes/223/v1/overlay/`：`zz_v1probe_config_test.go`／`zz_v1probe_cmdwisp_test.go`／`overlay.json`；**`internal/config` 与 `cmd/wisp` 目录里不存在这两枚文件，本腿零源码改动**）。
  尺：`go test -overlay .scratch/wisp/probes/223/v1/overlay/overlay.json ./internal/config -count=1 -v -run TestV1ProbeLoaderClassification` ＋ 同 overlay 的 `./cmd/wisp -run TestV1ProbeProductionSentence`，15:54:12→15:54:17，原始输出 `.scratch/wisp/probes/223/v1/ac4-probe.txt`，两发 `EXIT=0`（台件只读数、不判等，所以它绿不代表实现绿）。
- **票面要求的四句话各有真实出口且互不共用**（全量门里逐名绿）：缺失（`os.Remove` 真文件）／语法错（写 `this is not toml [[[`）／权限不够（**真 ACL `icacls … /deny *S-1-1-0:(R)`**，`TestTicket223PermissionDeniedSitsInItsOwnSentence` PASS 2.15s）／热加载被禁用（第二宿主 `wisp panel-inbound`，`cmd/wisp/panel_inbound.go:211` 逐字 `auditf("%s", hotReloadDisabledPanelInbound)`，这枚宿主确实 `config.NewManager` 开同一个 `config.toml`（`:200`）且 `Confirm: nil`（`:220`）⇒ "它没法弹卡"这句是真的）。
- **`peekSchemaVersion` 的解析错误确实不再被丢掉**（这是＋72/−6 的核心断言，实测）：17 发对抗输入里凡是"读不出版本"的形状（`B2/C1 之外的`…逐条见原始件）产出的错误文本都以 `config.toml parse: toml: …` 开头（＝loader.go:86 那枚 `observe.Wrap(…, peekErr, "config.toml parse")`），生产句子＝`cause=syntax`。**改前**这一支不可能出现（错误被 `_ =` 丢掉 ⇒ ver=0 ⇒ 一律进 `applyMigrations` 报"cannot migrate from schema version 1"）。
- **但实现腿自述的边界规则被实测推翻了一半**：它写的是"**读不出版本＝语法错；读得出版本＝交迁移管线继续说话**"。现测——**"读得出版本"且版本恰等于 `SchemaVersionCurrent=2`（`internal/config/schema.go:28`）时，`loader.go:93` 的 `ver != SchemaVersionCurrent` 不成立，迁移管线根本没被调用**，错误由 `decodeStrict`→`formatDecodeError`（`internal/config/parse.go:104-113` 的 `toml.DecodeError` 分支）给出 `config.toml: line N, col M: expected character =`，再被 `describeReloadFailure`（`cmd/wisp/config_reload.go:352` 的 `HasPrefix("config.toml:")` 分支）归入 **`cause=invalid`**，而那一句中文逐字是「config.toml **语法没问题**，但内容被校验拒绝（值不合法或引用解不开）」。**文件真正的毛病就是语法错**（`expected character =`）⇒ **这句话在这个形状上是假的＝"语法错"与"内容不合法"说反**。命中发数（17 发中 4 发）：
  | 对抗输入 | declaredSchemaVersion | 生产句子 | 该不该是语法错 |
  |---|---|---|---|
  | C1 注释行以 `[` 开头（`# [fs] …`）＋版本 2＋坏行 | `2,true`（注释被正确跳过） | `cause=invalid` "语法没问题" | **说反** |
  | E1 CRLF＋版本 2＋坏行 | `2,true`（`\r` 被正确剥掉） | `cause=invalid` | **说反** |
  | G1 `schema_version=2` 无空格＋坏行 | `2,true` | `cause=invalid` | **说反** |
  | L1 行首缩进的版本行＋坏行 | `2,true` | `cause=invalid` | **说反** |
  | J1 版本 99＋坏行 | `99,true` | `cause=invalid`（detail 本身是"newer build"，诚实；句子前半句仍多说了"语法没问题"） | 半说反 |
  | H1 同名键两次 | `2,true`（取第一枚） | `cause=syntax` | 对（靠 `formatDecodeError` 的兜底 `Wrap("config.toml parse")` 恰好救回） |
  | A2 版本 1＋坏表头（`migrate_test.go:123` 那枚既有断言的形状） | `1,true` | `cause=migration` | 对，且既有断言一字未放宽、仍在绿包里 |
  | B1 版本写在 `[app]` 里 | 交迁移管线→`unknown key "app.schema_version" at line 4` | `cause=unknown-key` | 对 |
  | D1/D2 多行字符串里有一行以 `[` 开头 | `0,false`（扫描器撞假表头就停） | `cause=syntax` | D1 对；**D2 说反方向＝"文件声明了版本却因为假表头读不出"，本腿判它对**（它确实是语法错，且它没有把两句调换成迁移） |
  | F1 BOM＋版本 2 | `0,false`（BOM 让 `schema_version` 匹配不上） | `cause=syntax`（detail 逐字 `invalid character at start of key: ï`） | 对 |
  | I1 首行就是 `[fs]` | `0,false` | `cause=syntax` | 对 |
  | K1 版本被写成字符串 | `0,false` | `cause=syntax` | 对（严格讲该归"值不合法"，但本腿不据此加码） |
  | M1 合法 v1 | 走真迁移 | `APPLIED hot=[] reload=[] restart=[] locked=0`，文件被重写并留 `.bak-1` | 对（唯一一发改了文件，符合 SPEC-03 §4.4） |
- ⇒ **结论**：AC#4 的字面判据（四件各一句、不许合成一句）**成立**；`语法错` 这句话从"生产里不可能出现"变成"可达且有 9 发读数"，这是真修好。带注的部分＝决定句子归谁的**实际规则不是"能不能读出版本"，而是"读出的版本等不等于 2 ＋ go-toml 抛的是哪一种错误对象"**，于是一批"声明了当前版本又语法坏"的文件被告知"语法没问题"。这条具名进推翻清单第 4 条，**建议另立票**（修法要动 `describeReloadFailure` 的分类或 `formatDecodeError` 的 Detail 前缀，都不是本腿权限）。
- ⛔ 本腿**没有**为了让它绿去改判据，也没有碰 `internal/config/migrate_test.go:123` 那句既有断言（现读该行仍在，包仍绿）。

## AC#6 裁决不在锁内（任务二：专属突变）

〔待填〕

## AC#7 "重启后生效"那一档要有出口 —— **成立但带注（AC#5 里那枚红就是它的用例）**

〔待填〕

## 编排者代落＋代跑这两件事的独立裁决

〔待填〕

## 跨票攻击：票 226 的"手改／不认领"在真轮询路径下是否仍成立

〔待填〕

## 推翻清单（实现腿＋编排者）

〔待填〕

## 票面写的 vs 我读到的（不一致单列）

〔待填〕

## 没做完／判不了（具名清单）

〔待填〕

## 门禁与尺：时刻表

〔待填〕
