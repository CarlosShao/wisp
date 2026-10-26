# close2-r1（收口腿 B·写腿）停手上报 —— 三张票的最后一格各有具名前置，其中两格的剩余写面落在 `cmd/wisp/**`

- 腿：`close2-r1`（写腿）｜起手锚 `cbb43c948`｜时点 2026-10-26 19:1x +08｜工作树＝共享（别的腿在飞的 `.scratch/**` 有改动，本腿一字未碰）。
- **结论**：三张票确实**各只剩 1 格**（现跑 `unchecked=1`），但派单那句「`Blocked by` 全 done」**不成立** ⇒ 按派单「未定义即停」②③④ 停手：**零产码字节、零测试字节；票面／`Status:`／`-done`／台账一字未动；⛔ push**。
- **为什么⛔ 交 `close2-r1-impl.md`**：没有 implementation 就⛔ 造一枚叫 `impl` 的件（"件名说 impl、盘上零字节"正是本仓反复咬人的那一族）。本件是停手上报件，写面内合法（`docs/evidence/s1/close2-r1-*.md`）。

## 0. 开工前三条（逐条现跑）

| 尺 | 读数 | rc |
|---|---|---|
| `git diff --cached --name-status` | 空 | 0 |
| `git status --porcelain -- cmd/wisp internal/speech internal/session internal/config internal/secret internal/memory` | 空 | 0 |
| 三票未勾／已勾（`^- \[ \] \*\*AC#` ⇄ `^- \[x\]`） | 216 **1／5**｜339 **1／6**｜223 **1／6** | 0 |

⚠ 工作树里有别腿在飞的 `.scratch/**`（`M .scratch/wisp/issues/115-…-done.md`、`M .scratch/wisp/probes/161/r6/logs/*`、`?? .scratch/…`）＝**不是我的字节，未碰**。

## 1. 三张票最后一格：原文 ＋ 判定 ＋ 读数

### 1.1 票 216 AC#5（票面 `:28`，逐字）

> `- [ ] **AC#5**：错误提示**按错误码分支**，**不许按英文错误串里有没有某个数字**（openchamber 那形 `ChatInput.tsx:2205` 靠 `413` 子串决定弹哪句＝反面教材，我 22:3x 只读到腿的转述、**未逐字复核**，落地前自己再看原文）。`

**判定：⛔ 本腿不可落——不是"没人做"，是编排者已把它判出本票射程。**
- 票面 `:68` 逐字「ⓑ 归口那行 DEFERRED 的 S6 子项、**⛔ 在本票翻它**、**⛔ 为它另立接线票**」；`:69` 逐字「此后任何人⛔ 把 AC#5ⓐ 读成"已钉"」；`:78` 逐字「`AC#5` 仍不勾＝该腿判"不适用"（那半已归口 S6 的 DEFERRED，账 `A1010`）」。
- 归口行现读：`docs/specs/SPEC-12-roadmap-governance.md:77` 逐字 `| DEFERRED | 完整错误文案体系 | D23 | … | S6 |`。ⓐ（禁令本体）＝`A1010` 登记为"仪器工程，本波不放"。

| 事实 | 我的读数 | 尺 |
|---|---|---|
| `MessageKey` 生产消费者 | **0**（只剩 `internal/observe/errors.go:160` 注释＋`:167` 声明；测试消费者 `errors_test.go:191`、`llm/matrix_14_2_test.go:175/249`） | `git grep -n MessageKey -- '*.go' \| grep -v _test` |
| `ErrorClass` 枚举 | **17** 枚（`errors.go:19-35`） | 现读 |
| 按错误码选提示的处数 | **0**；`runExitCode`（`cmd/wisp/run.go:77-95`）按码只分**退出码** 0/1/2/3/4 | 现读 |
| CLI 人话句 | **20+ 处**全是"调用点写死前缀＋`%v`"（`run.go:229/247/265/431/454/475/481/489/512/521/557/583/754/1300/1457/1463`…） | 现读 |
| `TestMessageKeys`（唯一消费者） | `--- PASS` `ok 0.034s` | `go test ./internal/observe -run TestMessageKeys -v` |

**"错误码 → 提示"映射表今天＝0 行**（空表就是读数，⛔ 把票面要求当已实现）。
**反例（现读）**：同一串 `Detail="payload too large"` 分别装进 `ClassProvider` 与 `ClassBudget` ⇒ **同一句提示**——没有任何一句提示去读 `Class`；全仓唯一读 `Class` 给人看的地方是诊断行 `run.go:1463` 的 `退出码 %d（error_class=%s）`，以及 `runExitCode` 的退出码。⇒ AC#5 要的"同一串不同码 ⇒ 不同提示"**今天为假**；做它得先落 S6 那层。

**另一条硬边界**：人话面里 CLI 那一半只能落 `cmd/wisp/run.go` ⇒ 派单②（`wire-r1` 独占）⇒ **停手上报 ②**。

### 1.2 票 339 AC#2（票面 `:27`，逐字）

> `- [ ] **AC#2 ★不许用放宽断言换绿**：⛔ 把 `want A` 改成"红也算对"、⛔ 删用例、⛔ `t.Skip`（无具名登记面的裸 Skip＝本仓禁手）、⛔ 摘掉那三句承重断言里的任何一句；★判据自带**去牙反控**＝把回落支整支摘掉（同 `334-v1` 的 M-V1 那一形）⇒ 修好后的仪器**必须**在**两平台**都点名红（这才是"补世界"而不是"藏红"）。`

**判定：⛔ 本腿不可落；且这⛔ 是"放宽断言"问题——是"缺世界"。**
- 票面 `:43` 逐字「那一格唯一的归处＝**必答 (iv) 那张"真乙"票**（我已据此立票 **347**…）」；`:39` 逐字「★**"去牙反控在两平台都点名红"这一维在本票写面内不可达**…＝**⛔ 任何腿为凑它去动产码**」；`:41` 逐字「`AC#2`＝**部分成立**…POSIX 半区＝结构不可达⇒摘牙零变化…这一维＝缺世界已登记」。

**结构三跳（我自己现跑）**：
1. `internal/risk/pathresolver_other.go:10` 逐字 `func resolveHandle(p string) (string, bool) { return "", false }`（首行 `//go:build !windows`）。
2. `git grep -n '\.root = ' -- internal/risk | grep -v _test` ⇒ **唯一一枚** `pathresolver.go:357`，在 `if final, ok := resolveHandle(cut); ok` 成功支内 ⇒ POSIX 上 `root` **恒空**。
3. `internal/risk/blacklist.go:362-365` 逐字 `if cand.root == "" || anchor.root == "" { return false }` ⇒ A 支（R17 fail-closed）在 POSIX **结构不可达**。

**改前 Linux 档为什么红**（339-r1 已落，我读码复核）：Linux 上假主目录造得出来（`:86` 同发设了 `HOME`，票面第 9 行那句"`USERPROFILE` 独有"已由 `:39` 作废），但候选落不进锚的祖先链 ⇒ 判级 A→B ⇒ `:131/:156` fixture guard 红、`:193/:197/:201` 三句承重断言红。
**改后为什么是对的**：339-r1 选**甲形平台拆分**——`anchor_miss_reason_334_test.go` 首行现读 `//go:build windows`；新增 `anchor_miss_reason_339_other_test.go` 首行现读 `//go:build !windows`（8004 B），把"这一维在本平台无世界"写成**具名绊线**。
**是修实现还是修断言？** 都不是：判据文字一字未动、`want A`/`ClassA` 承重句逐字在场、两件 `t.Skip` 命中 **0**（`339-v1` 复跑）。它做的是"把仪器收进它真有世界的平台 ＋ 让另一平台说出'我没有这个世界'"。
**"这不是放宽"的理由**：放宽的定义是"同一断言在本来会红的世界上变绿"；这里那半格**从来没绿过、且在那条结构下永远绿不了**——摘牙（去掉 `aRelates` 尾行回落支）在 POSIX 上是 `1/1/0/0 rc=0` **零变化**，即它连"红"都没有过。同时绊线**不是装饰**：把 `resolveHandle` 换成可解析后，`:140/:143` 当场点名红（`P0=37c42e94… ⇄ MUT=48f26a7c… ⇄ R=37c42e94… match=YES`，`339-v1` 重打）。**⛔ 我把"结构不可达"说成"放宽"**。
**本腿自己那一发（原样绿）**：`go test ./internal/risk -run 'AnchorMiss|AnchorFallback|FallbackBranch' -count=1 -v` ⇒ `TestAnchorFallbackTeethPerAnchor336`／`TestAnchorMissReasonNamesCandidateAndMissingAnchor`（3 子用例）／`TestAnchorMissKeepsNonOverridableClassAndNeedL2` 全 `--- PASS`，`ok 0.377s`，rc=0。
**唯一通路**：让 POSIX 有世界＝落 `resolveHandle` 解析器＝**票 347 的产码面**；347 首行现读 `Status: 已立，未派`、`AC#0` 6 格全空、逐字「开工前须机主拍一次」⇒ **停手上报 ④（要机主批，⛔ 腿自己扩面）**。

### 1.3 票 223 AC#2（票面 `:35`，逐字）

> `- [ ] **AC#2 三档生效级别各有读数**：立即生效的那一档要能观测到"改了→不重启→行为变了"；重启后生效的那一档要能观测到"改了→不重启→**明确告知未生效**"。⛔ **不许用"静默不生效"充当第二档**。`

**判定：⛔ 本腿不可落。** 票面 `:83` 逐字「**⛔ 不许拿"已归口 233"去翻 AC#2 的勾——233 落地之后再谈**」；`:52-54` 逐字给三条缺项：①reload 生产消费者（`OnReload` 在 `cmd/wisp` 今天 **0**）＋操作员读数；②hot 的"行为"今天只是 `configStore.ReadSettings` 读面答新值、而 provider chain 仍是 **boot snapshot**；③若要 config 层射程须**显式撤销 09-29 那句**。

**三档读数（我自己现跑）**：

| 档 | 生产读数 | 判定 |
|---|---|---|
| 热（立即） | 配置层真应用（`CheckAndReload` 生产调用者 `cmd/wisp/config_reload.go:157`），**但**运行时 `cfg := mgr.Config()`（`run.go:460`，注释自陈 `snapshot copy … frozen at what this boot read`）⇒ 模型/provider 链不换 | ★**行为未变 ⇒ 这档不成立** |
| reload | `Tier*` 三常量非测试命中**只在** `internal/config/schema.go:34-43` ⇒ **零赋值点**；`OnReload` 在 `wisp` **0**（唯一赋值点 `cmd/balldebug/main.go:244`） | ★**生产者 0／消费者 0** |
| 重启后 | `rt.mgr.OnRestartPending = rt.reportRestartPending`（`config_reload.go:119`）＋操作员句（`:325-330`） | ★**成立（AC#7 已翻）** |

`233` 现读：`Status: 待派`、6 格全空（`unchecked=6`）。⇒ 缺的两档都压在 `cmd/wisp/**`（装配根＋stdout 面）＝派单②；reload 那档另需 `233` ⇒ **停手上报 ②③**。

## 2. 反向判据（本腿零改动 ⇒ 无本腿变异；只给"原样绿"＋引用，⛔ 冒领）

- **本腿改动字节＝0** ⇒ "把改动摘掉 ⇒ 必须红"**没有可摘对象**；⛔ 我把别腿的牙证写成我的。
- 本腿**自己跑的**只有"原样绿"：339 §1.2 那发 `ok 0.377s`；216 §1.1 的 `TestMessageKeys` `PASS`。
- 引用（**是引用，⛔ 本腿复跑**）：339 Windows 半区 `339-v1` M-V1＝未突变 `22/3/0/0 rc=0` ⇄ 摘 `aRelates` 尾行回落支 `22/0/3/0 rc=1` 三枚点名红、还原哈希 `c1799cd2…` `match=YES`；**POSIX 同发 `1/1/0/0 rc=0` 零变化**——这行零变化**就是**本条停手的凭据。216 去牙（`return s`）⇒ 三枚红；223 摘 `:119` ⇒ `FAIL 42.44s`（各自裁决腿既有读数）。
- 还原核：本腿无变异、无还原动作 ⇒ N/A（⛔ 假报一发"还原绿"）。

## 3. 零回归（本腿未改任何包；改前＝改后）

| 包 | 改前 | 改后 |
|---|---|---|
| `./internal/risk` | `FAIL`，恰 **1 枚**红 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（既存，点名 `internal/panel/display_scrub_216.go:26/37/218`＋`composer.go:127`）＋SKIP 1 | **未动，同读数** |
| `./internal/observe` | `ok 4.088s` | 同 |
| `./internal/config` | `ok 1.569s` | 同 |
| `go vet ./internal/{risk,panel,observe,config}` | `rc=0` | 同 |

## 4. 四数（零字节，逐条）

| 面 | 字节 |
|---|---|
| `cmd/wisp/**` | **0** |
| `internal/speech/**` | **0** |
| `internal/session/**` | **0** |
| `internal/{config,secret,memory}/**` | **0** |
| `internal/{risk,panel,observe}/**` | **0** |
| 新增 | 本件 1 枚；仓外临时件 `D:\tmp\wisp\20261026\close2r1\**`（只建不删） |

## 5. 门禁（共享工作树，逐条 rc）

| 门 | 读数 | rc |
|---|---|---|
| gofumpt 甲（仓根 `.`） | ⚠ 已知假红：射程含 `.scratch/**` 的探针／突变件（不可解析 Go） | 2 |
| gofumpt 乙（`./internal ./cmd ./tools`，绝对路径 exe） | **恰 3 枚既存**：`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`cmd/wisp/models.go`；三枚工作树均净（`git status` 空），最后提交 `d03d166fe`／`5413f46d4`／`ca8834d72` ⇒ ⛔ 我的责任 | 0 |
| `bash scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations` | **0** |
| `go vet` 四包 | 静 | 0 |

## 6. 三条解锁路径（⛔ 我自己选，请编排者拍）

1. **216**：要么撤销／收窄 `A1010` 的「⛔ 在本票翻它」并放开 `cmd/wisp` 或 `internal/panel` 写面（＝提前做 S6 那层）；要么把 AC#5 按"不适用"交非实现者裁→你翻框；要么维持挂着（现状）。
2. **339**：把 `347 AC#0` 的一／二／三摆给机主。选"一"⇒347 落解析器后那半格才第一次有世界；选"二/三"⇒AC#2 后半按"缺世界"收窄（同 AC#3 的收窄笔法）。
3. **223**：先派 `233`（reload 生产者＋消费者）＋腾出 `cmd/wisp` 写面；若走 config 层射程，须显式撤销 09-29 那句。

— 本腿未翻任何框、未改票名／`Status:`／未加 `-done`、未改台账、⛔ push。
