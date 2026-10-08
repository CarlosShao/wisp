# 35-r7 — impl.md（票 35 `:52`「No secret leakage」残差①：入向 raw 不泄密仪器）

腿：`35-r7`。写面交付＝`internal/panel/inbound_raw_leak_35r7_test.go`（1 枚，278 行，纯测试、**零产码改动**）＋本目录证据件＋票 35 文末一节。
起手锚＝`00-anchor.md`（commit `62891ef0`，锚 HEAD `d1dddb77`，先于任何 `go` 命令）；仪器件 commit＝`6038531c`。

## 1. 做了什么

`ComposerDispatch.Handle`＝`internal/panel` 的入向整跳（网页递给 Go 的那一份原始 envelope 从这里进）。票 248 那族尺从没把**入向 raw 本身**当被检输入（普查腿 `35-a6` 的残差①结论，本腿现量复核同色）。本仪器：

- 用带凭据形状的 raw 打**五路由**：凭据写成功（`config.set`/`provider_credential`，canary 藏在专用键 `credentialValue`）、普通字段值即通用 `sk-` 形状（`provider_base_url`）、未在册名拒、伪造来源拒、处理器未接拒。
- 逐路由扫**本跳自己产出的三类落盘面**：audit 行（`[audit]` Sink 捕获＝真装配里落盘的那面）、返回宿主的 receipt/refusal 串、返回的 error（宿主的 `INBOUND-DISPATCH ... err=%v` 会把它 echo 进审计行）。
- 断言**写通道真的发生了**（`credCalls==1`、`lastRawSeen==raw`、receipt 里引用在）——干净读数若是静音就不算覆盖（248 原尺同款纪律）。
- 判定尺＝对 `cmd/wisp/panel_config_248_test.go` 的 `canary248`（`:39`）/`secretShape248`（`:43-46`）/`hits248`（`:52`）**逐字镜像**，零新增规则；另立 `TestDecisionRulesMirrorTicket248Ruler35r7` 读原件对拉（含「原件 MustCompile 枚数＝2」的形制尺），任一单独漂移即红，原件不可读＝ `t.Fatal`（孤儿镜像不算绿）。一把尺、两端钉死。
- 失败路径不回显命中原文（沿 248 头注纪律：测试日志本身就是被检面），红句只报**面名＋规则**：`leak surface: audit[0] matches sk-[A-Za-z0-9]{12,}`、命中哨兵处打 `<canary-redacted>`。

## 2. 残差②为什么是具名死格（⛔ 没为它造实现）

「桥载荷不泄密」里 **Go→页推出去的那半句**今天没有被测物：出向整跳零产码——`pageTransport` 只有 `Bind`/`Init`（`cmd/wisp/panel_host_windows.go:652-655`，普查现量、本腿复跑同色：`.Eval(`/`PostWebMessage`/`EvaluateScript`/`CreateWebMessageAsJson` 非 test 全 0，`lastPanelSnapshot()` 非 test 只有定义）。没有推载荷还立「推载荷不泄密」的判据＝造一个不存在的实现来假绿。先例＝票 248 `AC#4` 按同样理由从队列摘。等出向落地后，把推载荷补进本仪器/248 原尺的面清单即可（两者同形，扩面＝清单加一行）。

## 3. 具名顶回（派单转述 vs 盘上原文）

1. ⛔ **「复用既有载具 `hits248` 那一族（⛔ 不要新造第二套判定器）」物理不可达**：`hits248`/`secretShape248` 是 `cmd/wisp` 包的**包内测试符号**——Go 不能 import 测试文件；`internal/panel` 反向 import `cmd/wisp` 成环。普查腿建议的落点（`cmd/wisp/` 同包复用）又与本程硬约束（⛔ 写面禁 `cmd/wisp/**`＝`255-r1` 编译面）正面冲突。⇒ 处置＝逐字镜像＋防漂钉（§1），冲突具名上报，不静默照抄错话。
2. ⚠ 派单沿普查写「票面 `:51`」——盘上现量（本腿 `grep -n`）：未勾框行号＝`:42/:44/:45/:47/:49/:52`，「No secret leakage」在 **`:52`**（`:51` 是 Backpressure 的续行）。普查腿锚更早，非错话，名册以现量为准。
3. ⚠ 行数超预估：派单「预计 1 枚文件、60–150 行」，交付 278 行——超出全在文件头纪律注释与控制枚数，判定/断言本体约 120 行。具名报备。
4. 其余零枚：派单里「`panel_config_248_test.go:154`/`:218`」「六面清单」「`tools/d22scan` ban #3」逐枚现量复核同色（`:154`/`:218` 分毫不差）。

## 4. 突变自证（定向突变＝摘掉真不抄的那一支，逐发红句）

三形全部：临时改产码 → 跑 `-run TestInboundRawCredentialShapeReachesNoArtifact35r7 -v` → 红 → **原样摘回**（终态 `git diff internal/panel/composer_dispatch.go internal/panel/bridge.go internal/panel/config_handlers.go` ＝ 0 行）。

| 形 | 突变 | 指名红（rc） | 红句（面＋形状，节选） |
|---|---|---|---|
| M1 审计行抄内容字段 | `composer_dispatch.go` record 加 `text=%q`（req.Text） | rc=1；refused-unlisted / forged-source / unattached 三子例红 | `leak surface: audit[0] matches sk-[A-Za-z0-9]{12,}` ＋ `<canary-redacted>`（logs/mutation-m1） |
| M2 拒绝句回写原文（来源判定改宽） | `bridge.go` 来源拒绝 err 带 `原文=%q` | rc=1；forged-source 子例红，receipt＋error＋audit 三面 | `leak surface: receipt/error/audit[0] matches ...`（logs/mutation-m2） |
| M3 凭据 accept 审计抄 raw | `config_handlers.go` credential-accepted 追加 raw 行 | rc=1；credential-write-accepted 子例红 | `leak surface: audit[1] matches sk-...`（logs/mutation-m3） |

恒真自查（两枚正控，缺一不作凭据）：
- `TestInboundLeakRulerFiresOnPlantedShapes35r7`：哨兵与**通用 `sk-` 形状（不带哨兵）**各种进 audit/receipt 面具名发火⇒尺不是瞎的、也不只是名字尺。
- `TestLegalPanelTrafficIsStillWritten35r7`：合法流量（`role_chat_model=deepseek-chat`）照样到腿、审计照落、receipt 照写，尺不响⇒门没被刷黑。

## 5. 门禁 rc 名册（逐件自带 `rc=` 行；测 rc 句前无管道）

| 门禁 | 改前 | 改后 | 件 |
|---|---|---|---|
| `go vet ./internal/panel/` | rc=0 | rc=0 | logs/vet-after.txt, logs/vet-final.txt |
| `go test ./internal/panel/ -count=1`（整包名册化） | rc=1，红 4 枚 | rc=1，红**同 4 枚** | logs/test-before-fullpackage.txt, logs/test-after-fullpackage.txt |
| `go test -run '...35r7' -v -count=1`（枚枚指名） | —（件不存在） | rc=0，4 枚 PASS（含 5 子例） | logs/named-green-initial.txt, logs/named-green-final.txt |
| `gofmt -l internal/panel` | rc=0、空 | rc=0、空（我的件经 `gofmt -w` 格式化后） | logs/gofmt-before.txt, logs/gofmt-final.txt |
| `sh scripts/d22scan.sh`（既有支持起法） | — | rc=0，clean，`internal/=515` 含 `_test.go` 与注释 | logs/d22scan-final.txt |

改前红 4 枚（既有、非本程账、⛔ 一枚没放宽）：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`。作差＝**新增红 0、新增绿 4（枚级）**。⚠ 踩坑具名：仓根 `go run ./tools/d22scan` 报「main module does not contain package」（`tools/d22scan` 是独立 module）——起法以 `scripts/d22scan.sh` 为准（`.github/workflows/ci.yml:135` 原文「Wrapper is the supported entry point」）。

## 6. 纪律自报

- 写面：仅 `internal/panel/inbound_raw_leak_35r7_test.go` ＋ `.scratch/wisp/probes/35/r7/**` ＋票 35 文末一节。`cmd/wisp/**` 零笔——起手锚时该面 4 行脏（`255-r1` 在飞），本腿交件前复量 `git status --porcelain -- cmd internal docs .github scripts`＝**0 行**（他们已自行入库，与本腿无关），本腿全程一枚没 stage。
- ⛔ 未动 `frontend/**`；⛔ 未新增 C17 名/白名单条目；⛔ 无 L2「允许」面改动；⛔ 无 `filepath.Clean|Abs` 文件系统决策（新件里 `filepath.Join` 仅拼测试夹具只读路径，沿 `frontend_hygiene_test.go`/`panelRepoRoot` 既有形制）；⛔ 无 `go func(`；⛔ 无墙钟超时；⛔ SLO/golden/thresholds 零字节。
- Git：每笔显式 pathspec；⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／push；⛔ worktree；临时件只建不删；证据件无 `.out`、无 0 字节（`wc -c` 见目录清单）。
- 票面 AC 框：追加前后 `grep -cE '^[[:space:]]*- \[ \]'` ＝ **6 → 6**，`git diff --numstat` ＝ 纯追加（0 删），节头 `grep -c` ＝ 1（无重复块）。

## 7. 具名限制（交后续程，⛔ 不在本程做）

1. `record()` 把处理器/封套错误 `err` **逐字 echo** 进审计行：若 `cmd/wisp` 的腿把 raw 塞进 error 字符串，本包不 scrub——腿契约「只回状态」写在 `config_handlers.go` 头注，验腿面在 `cmd/wisp`（本程写面禁）。⛔ 不据此降档：判据没有「来源」这一栏就不按来源豁免。
2. 普通字段的 `SettingWrite.Value` 直通写腿＝票 248 已裁的设计通道；值落 `config.toml` 前的凭据形状校验在腿侧，同样不在本程。
3. 残差②见 §2。
4. 页面声明面（`frontend/src/lib/panel.ts`）本仪器不需要、亦未动。
