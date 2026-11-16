# 313 R1 复审残留两条修正（写腿 agent-88，2026-11-16）

> 派单：`.scratch/wisp/probes/orch/20261116/txt/313-r1-residual-dispatch.txt`。
> 复审腿 `auditR1recheck`（`docs/evidence/s1/audit-r1-recheck-77abb659b.md`）判 R1 六跳全通、P1 有牙、P2 合规、P3 可达；
> 本单只修它抓到的两条残留：①absent-key 崩页（P0）＋②boot 三条阈值只建模一条（说假话）。
> ⛔ 未回退任何已修好的东西（R1 的判定式标签、回执分流、能力位三面全保留并向前扩展）。

## 0. 基线红（修前实测）

- **残留①崩点复现**：`npm --prefix frontend run render:model-settings`（node v24.21.0）在**未改任何产码**时
  **rc=1，`ProviderEdit` 渲染中抛 `TypeError: Cannot read properties of undefined`**——仪器夹具 ROW 无
  `modelCapabilities` 键（wire omitempty 形状），`model-settings.tsx` 旧判据 `row.modelCapabilities === null ?
  undefined : row.modelCapabilities[id]` 对 undefined 取下标即崩（复审腿 Node 实测同款）。
- 复审腿的探针读数（`audit-r1-recheck-77abb659b.md` §2.4）：`PROBE-WIRE-ROW name=emptyprov modelCapabilities
  present=false`；零模型 provider 行出厂可达（`internal/config/settings_313.go:154` 建行可不带模型 id ⇒
  `cmd/wisp/panel_config_store.go` 旧 `if len(p.Models) > 0` 不成立 ⇒ map 保持 nil ⇒ `omitempty` 抹键）。

## 1.【残留①】三态判据（有值／腿没答／该模型没答）——两处崩点＋同页同类全修

- `frontend/src/components/model-settings.tsx:674-683`：判据改为
  `const caps = row.modelCapabilities ?? null; const bits = caps === null ? undefined : (caps[id] ?? undefined);`
  —— `?? null` 把**键缺席（undefined）**与 null 归并为同一事实「腿没答」（Go 侧 wire 是
  `modelCapabilities,omitempty`，`internal/panel/config_handlers.go:776`，没答＝键不上线）；一条**已答但缺该
  模型 id** 的记录也归「等快照」，永不猜 false、永不抛。
- `frontend/src/components/voice-settings.tsx:122-138`：三态助手 `capabilitiesOf / capsOf / enabledOf /
  enabledOfModel`，`chainAssembly` 全部查表改走三态（`:189-202`）。
- **同页同类顺修**（同一 omitempty 指针门族，复审 §2.4-a 的缺席形态类）：`voice-settings.tsx:606-618`
  LocalAudioCard 的 `aecEnabled / inputDevice / sampleRate / aecBuildAvailable` 四键判据从 `!== null` /
  `!== false` 改为 `?? null` 归一——缺席不再被当成「已答＝关」。
- `models` 键无此问题（`providerPayload.Models` 无 omitempty，`config_handlers.go:759`，nil 走线为 null）。

## 2.【残留②】boot 三条阈值全建模——读面先补，前端不猜

boot 的拒绝顺序（`internal/speech/cloudcatalog.go` `resolveChainElement`）：目录存在 → `enabled` → 能力位 →
`api_key_ref`。页面 `chainAssembly`（`voice-settings.tsx:171-218`）现按同一把尺判三条＋同名先出：

| 阈值 | boot 拒绝句 | 读面（file:line） | 页面判定 |
|---|---|---|---|
| 能力位缺 | `:95-102` lacks capabilities.audio_in/audio_out | `modelCapabilities`（`config_handlers.go:776`；填于 `cmd/wisp/panel_config_store.go:176-190`） | **已有**，`lacks` ⇒ `装配不了·缺能力位` |
| 模型 `enabled=false` | `:87-90` names a model disabled in the catalog (enabled=false) | **新补读面**：`ProviderSettings.ModelsEnabled`（`internal/panel/config_handlers.go:274-283`）＋ wire `modelsEnabled,omitempty`（`:777-785`）＋ ReadSettings 填充（`cmd/wisp/panel_config_store.go:180-190`，从 `spec.Enabled` 读） | `disabled` ⇒ `装配不了·模型已停用`，点名元素（`stepfun/x 已停用（enabled=false）`） |
| `api_key_ref` 未配 | `:121-131` resolving api_key_ref ... failed (Unconfigured) | **已有读面，本腿接线**：`apiKeyRef`（`config_handlers.go:764`）＋ `credentialSet`（`:786`）＋ `credentialError`（`:787`）；填于 `panel_config_store.go:194-208` | `nokey` ⇒ `装配不了·缺 API Key`，点名 provider＋ref 状态（`x 缺 API Key（api_key_ref dpapi:... 未就绪）`；credentialError 形状则 `无法核对密钥（<原文>）`）。空 ref 放行——boot 对空 ref 同样跳过（`:121 if ref != ""`） |

- **四态标签最终文案**（`voice-settings.tsx:411-421`）：`已接入·重启生效`(ok，**唯一入口＝assembles**) /
  `装配不了·模型已停用`(warn) / `装配不了·缺 API Key`(warn) / `装配不了·缺能力位`(warn)；快照没答＝
  `已接入·装配未判（等快照）`(neutral)；链空＝`未设置`。多门槛同挂时标签取 boot 序第一个（disabled >
  nokey > lacks），橙色点名行列全（`:443-453`），无一藏。
- **保存回执分流**（`:469-494`）：disabled/nokey/lacks 各给 error＋点名＋可动作指引；unknown warn；全过才 ok。
- **credentialError 首次上语音页**（`voice-settings.tsx:290-302`）：凡 provider 行 `credentialError !== ""`
  即整行渲染 `无法核对密钥：<name>（<原文>）`——复审点名的「语音卡要渲染 credentialError 现在没渲染」。
- **第五态标签归位（复审 §2.5③）**：目录里没有这一行不再借用「缺能力位」标签——归 `lacks`
  的点名行原文 `目录里没有这一行`（标签仍 `装配不了·缺能力位`，低害项按派单范围未单开标签）。

## 3.【行为牙】

- **崩页牙（前端，`frontend/scripts/render-model-settings.tsx` tooth 6）**：
  - `:474-482` 夹具经 JSON round-trip＋delete 造出 **modelCapabilities 键真缺席** 的 wire 行；
  - `:505-514` 模型页喂该行 ⇒ 渲染 `等快照` 且不抛、且不挂任何能力位开关（不猜 false）；
  - `:558-564` 语音页喂存链 ⇒ 渲染 `已接入·装配未判（等快照）` 且不抛；
  - `:522-554` 正常值牙：`示例模型-1` audio in `aria-checked="true"`、audio out `"false"`，`示例模型-2` 双 false。
  - **基线红**：同仪器修前跑＝rc=1（ProviderEdit TypeError，§0）；修后 rc=0。
- **三枚假象牙（前端 tooth 6，各喂一次，标签都不许是已接入）**：
  - `:622-640` 假象 1：`apiKeyRef` 已声明＋`credentialSet: false` ⇒ 断言 `装配不了·缺 API Key` 且
    全文无 `已接入·重启生效`，点名 `stepfun 缺 API Key（api_key_ref …`；
  - `:642-656` 假象 2：`credentialError: "凭据存储未接入"` ⇒ 同上＋`无法核对密钥：stepfun（凭据存储未接入）`；
  - `:658-668` 假象 3：`modelsEnabled["示例模型-1"]=false` ⇒ `装配不了·模型已停用`＋点名
    `stepfun/示例模型-1 已停用（enabled=false）`；
  - `:671-677` 全过牙：三条全过 ⇒ `已接入·重启生效` 恰出现 1 次，且全文无 `装配不了·`。
- **Go 读面牙（`cmd/wisp/panel_modelcaps_313_windows_test.go`，windows 载具面，真 settings 腿）**：
  - `:271` `Test313ModelsEnabledReadFaceAnswersTheDisabledGate`：种一 enabled 一 disabled 模型、链指向
    disabled ⇒ 先断言 boot 原句 `names a model disabled in the catalog (enabled=false)`（PROBE-DISABLED 同句），
    再断言 ReadSettings 的 `row.ModelsEnabled` 答出 false/true；
  - `:343` `Test313CredentialReadFaceNamesTheUnresolvableKeyRef`：provider 声明 `api_key_ref` 且无凭据存储 ⇒
    boot 拒（api_key_ref 句族，PROBE-KEYREF 同族）＋ ReadSettings 答 `CredentialSet=false`＋
    `CredentialError` 非空——页面 key 门的两个字段就是这两个。
- **payload 键 pin 更新（非放宽，照上一腿先例）**：`internal/panel/config_313_test.go:387` wantRow 键集加入
  `modelsEnabled`（测试内注释写明这是点名的 re-adjudication，键集仍 exact-match；value-bearing 扫描照跑）；
  `:437-441` 加 `modelsEnabled` 字节 pin `{"m":true}`。
- **仪器陈旧断言重审（非放宽）**：tooth 5 的「wired=2」写于静态标签时代；R1 改判定式后未更新过、本次基线
  已不可达。重审为（`render-model-settings.tsx` tooth 5 注释写明）：VOICED（链已存且已答）断言
  `已接入·重启生效`=1 ＋ `装配不了·缺能力位`=1 ＋ 占位=1；未存链视图断言 `未设置`≥2 且 wired=0。每个计数
  仍是 exact-match，占位 pin 原样保留。

## 4. 门禁读数

- `npm --prefix frontend run typecheck` rc=**0**；`npm --prefix frontend run build` rc=**0**（node v24.21.0 @
  D:\work\server\node24）；构建产物已删回，`frontend/dist/` 仅 `.gitkeep`（`git status --porcelain frontend/dist` 空）。
- `npm --prefix frontend run render:model-settings`（行为牙仪器）rc=**0**（§3 全部断言过）。
- `gofmt -l internal/panel cmd/wisp` ＝ 仅存量 `cmd/wisp/approval_always.go`（上一腿复审 §5.3 已证与本提交无关）。
- `go vet ./internal/panel ./cmd/wisp` rc=**0**；`go test ./internal/panel ./internal/config` rc=**0**。
- `bash scripts/d22scan.sh`（D:\work\soft\Git\bin\bash.exe）＝ `d22scan: clean - no D22 ban violations`，
  `d22scan.sh: OK - control rc=0, scan rc=0`。
- `go test ./cmd/wisp` 只经载具 `scripts/wisp-cli-staged-tests.sh`（bash.exe），一次一发；跑前
  `Get-Process wisp.test` ＝ 0。整包读数见 §5。

## 5. 载具整包读数（A1882 口径）

- **载具**：`scripts/wisp-cli-staged-tests.sh`（经 `D:\work\soft\Git\bin\bash.exe`，一次一发；两发跑前
  `Get-Process wisp.test` 计数均 0）。日志留档：
  第一发 `.scratch/wisp/probes/orch/20261116/txt/313-r1-residual-vehicle-run.log`，
  第二发（终读数）`.../313-r1-residual-vehicle-run2.log`。
- **如实记一笔（第一发）**：整包 rc=1，含 **1 枚本腿自己的红**
  `Test313CredentialReadFaceNamesTheUnresolvableKeyRef`——红因＝keyref 种子的模型没声明
  `capabilities.audio_in`，boot 在能力位门（cloudcatalog.go:95）先拒、根本没走到密钥门（:121），
  探针自身种错了尺。修法＝种子补 `capabilities = { audio_in = true, audio_out = false }` 并写明原因
  （测试文件内注释）。按 A1893 单跑那一枚复量＝PASS（同载具，rc=0）。
- **第二发（种子修正后的整包，终读数）**：`FAIL github.com/CarlosShao/wisp/cmd/wisp 1138.00s`，rc=1；
  **全部红都是派单点名的环境红**（Wisp 球 PID 49348 在跑，别当本腿弄坏）：
  - `TestAC4FocusReturnToPriorWindowGap33r5`（红句：`the panel did not take the foreground on Show:
    foreground 0x5280c68, panel hwnd 0x23c0e88`＝常驻球占着前台）；
  - `TestAC246DevLegIgnoresTheTestTaskInjection`（stderr：`AllowSetForegroundWindow(ASFW_ANY) was
    refused ... Access is denied.`＝前台锁拒绝，同族）。
- **终点锚（A1882）**：`TestZZ315A1` 三枚全命中且全 PASS（RootLoop / SiblingChild /
  ProjectInstructionsSeed）；**最后一条 === RUN**＝
  `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`，日志以包结果行收尾
  （8041 行，无截断、无无声断掉）。
- **新测试位置**：在截断点（包尾）**之前**：`TestTicket255CloudChainAssemblyCallPointIsAlive`
  PASS（log:1284）、`Test313CapabilityBitFlipsAssemblyOutcome` PASS（log:3156）、
  `Test313CapabilityWritesRefuseByNameAndMoveNoByte` PASS（log:3158）、
  `Test313ModelsEnabledReadFaceAnswersTheDisabledGate` PASS（log:3160）、
  `Test313CredentialReadFaceNamesTheUnresolvableKeyRef` PASS（log:3162）。
- SKIP 均为载具名册的有名跳过（fixtures/活麦克风/活窗口族），FIXTURE-GAP＝**0**（无名册缺口）。
- 附一笔载具传参教训：经 pwsh 直传 `-test.run '^(...)$'` 会被 MSYS 拆成 `-test` ＋ `-run`（红句
  `flag provided but not defined: -test`，rc=2）；用 `bash -c "… -test.run '^(...)$' …"` 整串传参则正常。

## 6. 边界与声明

- ⛔ 未动 thresholds/golden/allowlist/PLAN/specs/D43/design/C17 白名单/dist 字节；⛔ 未 push；⛔ 未杀球 49348。
- `TestAC4Focus*`／`TestAC246*` 形环境红（球 PID 49348 在跑）按环境报，不算本腿弄坏。
- 「未知先于未知」的取序注释（voice-settings.tsx:160-163）：硬拒绝压过 unknown（boot 不受页面读不到影响），
  三硬门间按 boot 元素内检查序 disabled > nokey > lacks；点名行列全，故无门槛被藏成已接入。
- 顺修范围声明：LocalAudioCard 四键与 `models` 之外的同族缺席形态（§1 第三条）与残留①同型同页，
  属同一判据修正，非新开票。
