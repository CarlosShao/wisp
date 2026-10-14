# 114-a5 — 票 114 九格拆分表 ＋ "干净子集"三格落地

> 身份＝实现者腿（判语归非实现者）。只 commit 不 push。⛔ 未翻任何 AC 框、未动 `Status:`／票名、未写台账。
> 锚：起手 HEAD `037bfb7f`（2026-10-14 15:3x +08）；交件时 HEAD 已由他腿推进到 `88c7e752`（253 收口等），本腿只动两枚文件。
> 落地面＝`internal/panel/composer_test.go`（AC#4＋AC#9）＋`internal/panel/bridge_test.go`（AC#8）——两枚都是非冻结测试件、非契约、非 `frontend/`、非 `.github`。
> 读数全在 `logs/01…06`（每枚含原命令行；被并进合并件的原始分片留在同目录**未提交**，内容逐字包含）。"改前"一律用**仓外** overlay／温度拷贝造，仓里产码未被任何突变碰过。
> ⚠ **纪律偏差自报**：产码 commit `d8dd4b38` 除两枚测试件外，**还带上了编排者先前已暂存未提交的** `.scratch/wisp/issues/253-panel-inbound-three-ruler-holes.md` **删除**——那是 253 改名 `-done` 收口的未完成半边（`git diff --cached --name-only` 现读＝3 行，我照单全提了）。**内容无损**：收口版＋新名 `253-…-done.md` 在 `88c7e752`，旧名 blob `02d5632…` 在历史里；⛔ 未 amend（禁）。教训具名：本仓 `git commit` 不带 pathspec 会提**整个 index**——下一腿请用 `git commit -- <显式 pathspec>`。

## §1 拆分表（九格 × ① 要写哪些 `file:line` ② 面 ③ 现读是否已被别的票满足）

| 格 | ① 要写的文件:行（本腿现读） | ② 面 | ③ 现读判定 ＋ 本波处置 |
|---|---|---|---|
| **AC#3** 更宽档走 L2 卡（300s=拒）＋审计三档全写 | `cmd/wisp/panel_inbound.go:254`／`:264` 两枚 `Confirm: nil` ⇒ 接 `cmd/wisp/run.go:677` 已握的真确认腿；`cmd/wisp/run.go:340` `rt.modeWrites`（字段）/`:677`（赋值）＝全仓**零读者**的死码（接它或具名删它＝编排者一语）；常驻腿卡面（`panel_resident_windows.go` 一族）；反向对照用例（真窗 `-tags winlive` 或 headless C18 卡） | **真机**（C18 卡）＋`cmd/wisp`（契约邻接：审批面注入形未定案） | **部分满足**：审计半支已在（`internal/perm/store.go:172-214` 三档全写；`internal/agent/approval/queue.go` 的 `DefaultApprovalTimeout` 300s）；卡片半支**生产不可达**——变宽请求在 `internal/panel/composer_handlers.go:133-138` 被 `ErrNoL2Confirm` 拒、走不到 `Set` 的卡（行号本腿复量；同 a2 普查 §2 AC#3）。⇒ **停手上报**（J1 形未定案＋真机；本波⛔不落、不翻框） |
| **AC#4** 门二扩展名集与 d22scan 文本类一致 | `internal/panel/composer_test.go`：`rendererSourceFile`（扩到并集）＋`TestPlantedRendererDoorShapesGoRed`（加 `.json` 种形）＋新 `TestDoorTwoReadsEveryTextTypeTheScannerWalks`（对齐尺） | **本队可写**（仪器） | **部分满足**：结构支已成立（92b 返工）；差 `.json`（`tools/d22scan/main.go:1312` `isTextFile` 有、门二无）＋两枚名册无对齐尺。⇒ **本波落地**（§2.1） |
| **AC#5** 只许显示＋发起请求；输入口不得变授权口；ban#6/ban#8 | 无产码要写。Go 半支＝既有钉：`internal/agent/approval/ui.go:167-171`（`PanelAPI` 无 `Allow`）＋`composer_test.go` 结构钉＋`internal/panel/config_handlers.go:108-136`（拒 `risk.`/`fs.`）；`frontend/` 半支＝d22scan／`frontend_hygiene_test.go` 读数 | **审计读数**（本队可写读数，无产码） | **Go 半支已成立**；frontend 半支＝跑仪器（§2.5 随行读数）。⚠ 票面射程 `U+2190–U+2BFF` 与 141 定案（六段、注释豁免）不是同一把——**按仪器读数报**，不按票面字面 |
| **AC#6** 真机差分截屏 | 前端侧（外部 agent）。本编队零枚 | **frontend（owner-delegated）** | **`skipped=frontend(owner-delegated)`**（票头 09-23 横幅＋台账 `A102`＋本波派单；保持未勾、⛔不替勾） |
| **AC#7** 门禁读数 | 无产码要写。读数＝`sh scripts/d22scan.sh`／`gofumpt -l`（`.exe` 全路径）／`bash -n`／`go build ./...`／长度门／POSIX 三合包 | **本队可写**（读数；⛔`scripts/**` 一字不动） | **"缺的那一发"已补**：本波改动后的全套原始读数在 §2.4（含 92b 记的 Git Bash `PATH` 坑：必须 `.exe` 全路径） |
| **AC#8** 两边同改变成门 | `internal/panel/bridge_test.go`（新 `TestEveryWhitelistedInboundNameIsAlsoSentByThePage`＋盘上种形正控；名字从 `bridge.go` 经 253 的 `parseFileInbound` 现取） | **本队可写**（仪器；不加名、不改 C17） | **253-r5 未覆盖此方向**（现读：其文件头逐字"reads nothing outside internal/panel's production sources"，它 gates 的是 declared/guard/router 三处 Go 侧登记脊柱，不读 `frontend/`）；旧唯一耦合＝N7 手抄四枚（`bridge_test.go:138-141`）⇒ 248 扩的 `config.get`/`config.set` 两枚名**此前无人守**。⇒ **本波落地**（§2.3） |
| **AC#9** 结构钉扫描根太窄 | `internal/panel/composer_test.go`：`TestTheRendererHoldsExactlyOneDoorToTheHost` 根 `frontend/src`→`frontend`；`TestPlantedRendererDoorShapesGoRed` 加 `frontend/public/*.js` 种形（R-92b-2 第二形） | **本队可写**（仪器）；"渲染器到底加载哪些文件"的**精确**答案仍归界面侧（a2 §6-J3） | **完全没做**（根写死在旧 `:554`；门二走整棵 `frontend`）⇒ **本波落地**（§2.2）：根＝整棵 `frontend`（任何加载集的超集，⛔不缩范围）；精确加载集留界面侧 |
| **AC#10** 两种宿主装法进同一枚钉 | 若落：(甲)`cmd/wisp/panel_host_windows.go`（宿主装配）或 (乙)`internal/panel/composer_test.go` 宿主调用点正则扩面／新钉（要读 `cmd/wisp` Go 源）。本腿现读补量：`AddHostObjectToScript` 在本仓 `internal/ cmd/ tools/` **零调用点**；真生产门是**第三种**装法 `w.Bind("wispDispatch")`（`panel_host_windows.go:82`/`:319`，253-r6 已把其静默丢弃改有声） | **frontend（owner-delegated）**（本波派单＋台账 `A925` 补具名 `AC#6/AC#10＝owner-delegated`） | **`skipped=frontend(owner-delegated)`**（不写不勾）。⚠ 具名给编排者：若改判要落"Go 侧可达性钉扩到 `wispDispatch` 绑定名"，落点在 `cmd/wisp` 宿主测试，另派 |
| **AC#11** `npm run render:composer && git diff --exit-code` 做成 CI 一步 | `.github/workflows/ci.yml`：现步骤 `:1026` `npm run render:composer`；全仓 `grep -c "exit-code" ci.yml`＝**0** ⇒ 防腐坏半支未做（a2 普查复认） | **`.github/**`＝契约邻接**（且要查票 111/85a 让渡） | ⇒ **停手上报**（本波⛔不动 `.github/**`；让渡与否＝编排者查） |

一句话：**九格里零枚被已收票整格做掉**；本波落三格仪器（AC#4/AC#8/AC#9），补一格读数（AC#7），随行一格审计（AC#5）；停手四格（AC#3 真机＋J1／AC#6＋AC#10 前端 owner-delegated／AC#11 `.github`）。

## §2 落地格读数（每格"改前 vs 改后"，命令原文见各 log 头）

### §2.1 AC#4 — 门二扩展名集对齐 d22scan 文本类（落地）

改动：`rendererSourceFile` 由 7 类扩为**并集**（+`.json .md .txt .yaml .yml .toml .go .svg .vue`，保留 `.mjs`；方向＝"scanner 的集必须被本过滤器包含"，**从不缩**，R16#4）；新增 `TestDoorTwoReadsEveryTextTypeTheScannerWalks`（从 `tools/d22scan/main.go` 的 `isTextFile` 现抽 15 类，逐一要求门二读；三枚地板名 `.go/.json/.html` 防"抽取器空读绿"；子测试用合成 switch 种 `.zzz114` 证尺有牙）。

- **改后绿**（logs/01）：`door 2 reads all 15 text types the scanner walks: .css .go .html .js .json .jsx .md .svg .toml .ts .tsx .txt .vue .yaml .yml`；真树其余全绿。
- **改前红**（logs/02，overlay 只回退扩展名集，diff 在件内）：同一发尺点名 **9 枚**未读类（`.go .json .md .svg .toml .txt .vue .yaml .yml`）；且种在 `.json` 里的 `"panel.mode.set"` 同时**逃过两枚仪器**——结构钉 `unknown` 读 2 枚（应 3）、门二只点 2 行（少 `ac92-plant-d.json`）。

### §2.2 AC#9 — 结构钉的扫描根（落地）

改动：`TestTheRendererHoldsExactlyOneDoorToTheHost` 的根由 `frontend/src` 扩到**整棵 `frontend`**（与门二同根；任何"渲染器加载集"都是它的子集，故不会漏 loadable 文件；精确加载集仍是界面侧的问题）；`TestPlantedRendererDoorShapesGoRed` 新增 `frontend/public/ac92-public-door.js` 种形（R-92b-2 第二形）并要求仪器**点名**它。

- **改后绿**（logs/01）：真树 `109 files scanned, 2 host call sites (all in src/lib/panel.ts)`（原 89 枚/2 点；扩根后两点不动＝真树干净）；种形树 `3 outside call sites, 3 unknown literals`，`public/ac92-public-door.js` 被点名。⚠ 窗口内 `frontend/` 正被界面侧写（settings-v2）：文件计数 15:39 读 **109**、POSIX 侧 ~16:0x 读 **111**，扩根后两发都绿（`111 files scanned, 2 host call sites`）。
- **改前红**（logs/03a，overlay 只回退两根）：同一发种形树读 `2 outside`、public 门**不被点名**、`unknown` 少 1 ⇒ 断言逐条红。
- **改前盲的现场**（logs/03b）：**不种形**时旧根在真树上照旧绿（`89 files scanned`）——这正是"绿与事实无关"、本格必须靠种形读的原因（复现 a2 普查：`21 files scanned, 2 host call sites` 一个数不动）。

### §2.3 AC#8 — Go 名册 ⇒ 页面字面量（落地）

改动：新增 `TestEveryWhitelistedInboundNameIsAlsoSentByThePage`（`bridge_test.go`）：名字用 253 的 `parseFileInbound` 从 `bridge.go` **现取**（前缀无关），逐一要求 `frontend/src/lib/panel.ts` 有 `"名字"` 字面量；自带**盘上种形正控**（把一枚未镜像名写进 `bridge.go` 的温度拷贝，同一 helper 必须点名它）。

- **改后绿**（logs/01/06）：真树 `page mirrors all 6 whitelisted inbound names`，正控报 `[config.unmirrored114]` 一枚不差；POSIX 侧同绿（`logs/06c`）。
- **改前绿（旧套件对同一漂移无感）**（logs/04a）：仓外模块拷贝（`D:/tmp/wisp114a5/repo8`，盘上把 `config.unmirrored114` 种进 `bridge.go`，放 HEAD 版 `bridge_test.go`）⇒ 旧 N7 **PASS**（手抄四枚之外看不见）。
- **改后红（同一棵种了漂移的树）**（logs/04b）：换上本腿新门 ⇒ **FAIL**，逐字点名 `frontend/src/lib/panel.ts never spells config.unmirrored114 as a literal`。
- ⚠ 诚实边界：overlay 换得到编译器、换不到 `os.ReadFile`（A713）；故"改前/改后红"用**盘上种**在模块拷贝里做（logs/04），⛔ 没用任何 overlay 绿抵账。

### §2.4 AC#7 — 门禁读数（本波改动后，原始输出入 logs/05、logs/06）

- `go build ./...`＝**rc=0**；`go vet ./internal/panel/`＝**rc=0**。
- `"$(go env GOPATH)\\bin\\gofumpt.exe" -l internal/panel/composer_test.go internal/panel/bridge_test.go`＝**零行输出**（rc=0；⚠ `.exe` 全路径写法照 92b 的 Git Bash 坑记）。
- `bash -n scripts/d22scan.sh`＝rc=0；`bash -n scripts/check-path-length-budget.sh`＝rc=0（⛔ 脚本一字未改）。
- `sh scripts/d22scan.sh`＝**rc=0 clean**：`scope bans #1-5 internal/=236 · cmd/=40 · ban #6 frontend/=111 · ban #7 internal/tools/=23 · ban #8 design/=39 · frontend/=111 · internal/=550 · cmd/=126`；自测 `runtests.sh: PASS=35 FAIL=0 SKIP=0`。
- 长度门 `bash scripts/check-path-length-budget.sh`＝**VERDICT GREEN（over-budget 60＝roster 60、not-in-roster 0）**。
- `internal/panel` 整包（windows，logs/06a）：`RUN=223 PASS=143 FAIL=5 SKIP=0`（四数＝原样 grep 口径：`=== RUN` 含子测试，`--- PASS/FAIL/SKIP` 只数顶层）；5 枚红**逐名与 HEAD 版文件对照件完全同集**（logs/06b 控制）⇒ **全部既有**，非本波：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree`（前端树耦合，35-r3 已记的既有红族）＋`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`（35-r8 既有）。
- **POSIX 三合包**（logs/06c，`golang:1.27` 容器，repo `:ro`＋宿主 modcache，`MSYS_NO_PATHCONV=1`）：`./internal/panel/ ./internal/tools/ ./internal/risk/` 四数 `RUN=621 PASS=441 FAIL=5 SKIP=5 rc=1`；**FAIL 5 名与 windows 逐名同**（同上，既有）；`internal/tools` ok／`internal/risk` ok；SKIP 5 名逐条：`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`／`TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3`／`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`／`TestD34WriteMatrix`／`TestWorkspaceSwitchRefusesAJunctionToOutside`（全为在册 POSIX-only 用例）。首跑（未挂宿主 modcache）在 `proxy.golang.org` connection refused 处搁浅、`RUN=0` —— **没被当读数**，换挂载后重跑。

### §2.5 AC#5 — 随行审计读数（无产码）

- Go 半支：`internal/agent/approval/ui.go:167-171` 的 `PanelAPI` 仍**无 `Allow`**（结构禁区在）；整包绿名单现读（logs/06a）含 `TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached`／`TestModeViewHasNoWriteSurface`／`TestThePanelSwitchCapabilityRidesOneWhitelistedRoute`。
- `frontend/` 半支（仪器读数，非本腿读码）：d22scan ban #6／ban #8 frontend/=**111 文件**、clean；包内孪生 `TestFrontendNeverNamesAnApprovalDecision`／`TestFrontendHasNoEmoji`（含 src+scripts+index.html）**PASS**（logs/06a）。

## §3 没做的（具名）

- **AC#3**：未接 `Confirm`（`panel_inbound.go:254/:264`）、未动 `rt.modeWrites` 死码、未开真窗/未造 headless 卡用例——J1（常驻腿拿什么举 C18 卡）在 spec/票面**无一处被具名**，按"未定义即停"上报。
- **AC#6／AC#10**：`skipped=frontend(owner-delegated)`，未写未勾（AC#10 的现读补量在 §1 那一行）。
- **AC#11**：`.github/workflows/ci.yml` 未动；票 111/85a 让渡未查（编排者面）。
- **AC#9 的精确加载集**（import 闭包＋入口 HTML＋public/）：未答，归界面侧；本波给的是不可漏的超集根。
- POSIX 首跑搁浅（网络）已如实记在 logs/06c 头部；四数为第二次真跑读数。

## §4 纪律自证

- 动过的仓内文件：`internal/panel/composer_test.go`、`internal/panel/bridge_test.go`（产码 `d8dd4b38`）＋本件＋`logs/01…06`＋票尾一条 bullet（交付 commit 另记）。⚠ `d8dd4b38` 带上 253 旧名删除一事见件头"纪律偏差自报"（具名、内容无损、未 amend）。
- ⛔ 未碰：三枚冻结件（`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）、`thresholds.go`、golden、`allowlist.txt`、`PLAN.md`、`docs/specs/**`、`tools/d22scan/**`、`scripts/**`、`.github/**`、`frontend/**`、`design/**`、`internal/risk/**`；⛔ 未改任何既有断言的判定（只追加/扩根/扩集，方向从宽）。
- 突变件全部只建不删、在仓外（`D:/tmp/wisp114a5/`）；`go test` 的 `-run` 一律小射程；⛔ 未跑整包 `cmd/wisp`（314 死区）。
