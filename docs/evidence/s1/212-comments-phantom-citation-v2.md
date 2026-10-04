# 212-v2 对抗验收裁决（票 212 第 2 轮修码 212-r2；非实现者腿，实现者＝212-r2）

- 本程：验收腿 **212-v2**。被验对象＝212-r2 的四枚 commit `7ac8965a`（骨架）→`5413f46d`（修码，五枚产码文件）
  →`1353ce66`（证据件满稿＋读数档＋两台突变脚本）→`d4b2e9b5`（§8 终态复量）。
- 起手锚：编排者给的 `d934e016`（212-r2 的未修码基线）；本程取数时 `dev` HEAD＝`095cfad8`（本程骨架枚）之上的共享树。
- 权威文本：票面 `.scratch/wisp/issues/212-comments-cite-evidence-files-that-do-not-exist.md`（只读，本程零写）
  ＋ 实现件 `.scratch/wisp/probes/212/r2/fix-and-readings.md` ＋ 前两任判语 `probes/212/v1/verdict.md`／`verdict-first-instance.md`。
- 本程全部读数落 `.scratch/wisp/probes/212/v2/`（`logs/v2-00..04` 与 `logs/v2-raw-*`），先落文件再判；
  派单里编排者给的数（37/37、行号、八枚 scope 值）一律按〔待验断言〕处理，下表凡"相符"都是本程现跑对上的。
- 方法学要点：**本程没有向共享工作树种任何假红**。所有"种样本"的探针落在仓外归档副本
  （`git archive HEAD` 与 `git archive 5e8748b3^` 解到系统临时区 `/tmp/212v2/{tree-head,preimage}`），
  用 `go build -o /tmp/212v2/d22scan-{fixed,mut}.exe` 两个二进制分别扫；临时件只建不删。

## §0 起手锚与本程自跑的三把尺（不复认编排者读数）

| 尺 | 本程逐字末行／rc | 与实现件 §3／§8 的关系 |
|---|---|---|
| `cd tools/d22scan && go run . -self-test` | `d22scan -self-test: clean - all 37 direction checks passed (20 expect-ring, 17 expect-silent)`／rc=0（档 `logs/v2-00-pristine-selftest.txt`，63 行） | **逐字相符**＝编排者那条〔待验断言〕转正；名册行逐字＝`roster read from main.go = 9 numbered ban(s) [1 bare-goroutine, 2 pathresolver-bypass, 3 plaintext-key, 4 wallclock-timeout, 5 mirror-hash, 6 panel-approval, 7 internal-artifact-tool, 8 emoji, 9 phantom-citation] + 1 finding type(s); 37 cases, 10 tag(s) covered, both directions required per tag` |
| `go vet ./...`（tools/d22scan 模块） | 无输出（`wc -c`＝0）／rc=0（档 `logs/v2-04-vet.txt`） | 相符 |
| `sh scripts/d22scan.sh`（全仓门禁，只读，未改码） | `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` ＋ `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, ... ban #8 cmd/=97`／rc=0（档 `logs/v2-03-gate-full.txt`，249 行） | 相符；八枚 scope 行与实现件 **`pre-scope-lines.txt` 逐字相等**（`diff` rc=0，见 §6） |
| `go test -count=1 ./...`（tools/d22scan 单包，由门禁第 1 步代跑） | `ok  	github.com/CarlosShao/wisp/tools/d22scan	24.711s` | 同一判语，秒数不同（负载），分母与 PASS 数才是要对的东西：PASS=34/FAIL=0/SKIP=0/RUN=76 相符 |

⛔ 本程**未跑**全仓 `go test ./...`（派单⛔＋`internal/ball` 有腿在飞）。

## §1 六格判语总表

| 格 | 问题 | 判 | 凭据（file:line 或逐字读数） |
|---|---|---|---|
| 1 | ⓐ 修法的牙（摘掉豁免，新 silent 样本必响） | **成立** | `logs/v2-01-mut-selftest-exemption-removed.txt:61` 那一枚 `silent FAIL`，三种拼法各响一次；末行 `1 direction(s) failed, 36/37 passed - the gate does not see what it claims` rc=1；还原后 37/37 rc=0，`md5sum tools/d22scan/main.go` 前后＝`e6d1745996e8e02c82ab691ac2cbadf8` 全等，`git diff -- tools/d22scan` 空（§2） |
| 2 | 排除的射程有没有过宽 | **带条件成立**：定向命中（真·全拼幻影仍响），但豁免按**区域起点**，尾部带 `…`/`*` 的**完整**幻影路径与同行粘连的缩写尾段会静默；**今日活体猎物＝0 枚**（delta 实测为空），且实现件自己在 `tools/d22scan/main.go:877-883` 与 §6 #4 具名登记了这一条 | P1 全拼幻影 `docs/evidence/s1/does-not-exist-full-name.md` ⇒ RING；P2 混形⇒静默（③）；P3/P4/P5 ⇒ 静默（过宽实证）；P6/P8/P9 ⇒ RING（控制组，证明豁免不外溢到别行／别 token）；clean HEAD 归档树上 fixed 与 mut 两扫的 token 差集＝空（`logs/v2-raw-tree-head-deleted.txt`／`logs/v2-raw-tree-head-deleted-mut.txt`）（§3） |
| 3 | 三处注释现在说的是不是真话 | **带条件成立**：三处逐枚对得上、路标全存在、ban #9 全仓不红；**一处指错具名**＝四枚字段的出处记错 | `internal/agent/approval/pending_read.go:41-43` 把 `Level`/`SessionOverrideBlocked` 也挂到 `internal/tools/gate.go:13` 名下，而那行只点名 `RulesHit` 与 `Reason`；四枚之说真身在 `cmd/wisp/panel_pump.go:44-46`（详见 §4 末两行） |
| 4 | 计数口径 §5 的 C／D 两式复算 | **成立**：C 排除在＝8／排除摘＝9／差集恰 `docs/evidence/s1/62-`（本程自跑，三数全等）；D 圈 pathspec＝9 hunk／8 文件（`queue.go` 2、其余各 1），整枚 commit＝35 hunk。A/B 两式不属可机算格（依赖 r1 死腿当日档），本程只核了口径 A 那份档的存在性（§5） | `logs/v2-raw-preimage-{fixed,mut}.txt` ＋ `git show 5e8748b3 \| grep -c "^@@"` → 35 |
| 5 | 票面禁区与 AC 框不许伤 | **成立**：四枚 commit 的文件名册里禁区一枚都不出现；票 212 那一枚在 `d934e016..HEAD` 区间**零 commit、零 diff**⇒AC 框零碰（§6） | `git show --name-only` 逐枚（§6 表）＋`git log d934e016..HEAD -- .scratch/wisp/issues/212-*.md`＝0 行＋`git diff` 同范围＝空 |
| 6 | 自报的"自咬"是否真消（ban #7 那侧） | **成立**：门禁 rc=0，`ban #7 internal/tools/=23` 与起手逐字相等，`internal/tools/` 侧今天零红；机制本程自证＝在同一枚归档副本的 `internal/tools/` 里种一行带引号的 `"internal-provider"` ⇒ `[internal-artifact-tool]` 当场红，删掉 ⇒ 该 tag 归 0（红句逐字见 §6） | `logs/v2-03-gate-full.txt`＋`logs/v2-raw-probequote.txt:14` |

## §2 格 1 取证：摘豁免突变（成对：改→跑→还原→md5 比）

突变＝把 `tools/d22scan/main.go:805` 的 `short := shorthandPathStarts(c.Text)` 换成
`short := map[int]bool{}`（＝摘回 2026-10-03 的形状），跑 `-self-test`，再逐字还原。

| 发 | 尺 | 逐字末行 | rc | md5(`tools/d22scan/main.go`) |
|---|---|---|---|---|
| 改前（起手） | `md5sum` | — | — | `e6d1745996e8e02c82ab691ac2cbadf8` |
| 红发 | `go run . -self-test` | `d22scan -self-test: 1 direction(s) failed, 36/37 passed - the gate does not see what it claims` | **1** | — |
| 还原后 | `go run . -self-test` | `d22scan -self-test: clean - all 37 direction checks passed (20 expect-ring, 17 expect-silent)` | **0** | `e6d1745996e8e02c82ab691ac2cbadf8`＝改前（`logs/md5-before.txt`／`md5-after-restore.txt`）；`git diff --stat -- tools/d22scan`＝空 |

红句逐字（档 `logs/v2-01-mut-selftest-exemption-removed.txt:61`，一枚样本响三次）：

```
d22scan -self-test: ban #9 phantom-citation    silent FAIL   CLASS 3 AS A PATH: "...", "*" and U+2026 inside a repo-relative token stay silent (ticket 212 裁 ⓐ) | rang on a sample that must stay silent:
  internal/probe/cites-shorthand.go:4: [phantom-citation] comment cites repo path "docs/evidence/s1/152-...-accept-r2.md" which does not exist on disk (ticket 212 ban #9): fix the citation or create the file; shorthand forms are a separate prescription, not a violation
  internal/probe/cites-shorthand.go:5: [phantom-citation] comment cites repo path "docs/evidence/s1/212-" ... （同上尾句）
  internal/probe/cites-shorthand.go:5: [phantom-citation] comment cites repo path "docs/evidence/s1/152" ... （同上尾句）
```

- 三枚 token＝三种缩写拼法（`...` 整吞、`*` 截成 `212-`、`…` 截成 `152`），与实现件 §4 的读数**逐字相同**；
  ⇒ 这枚 silent 样本不是只练了一种形状，ⓐ 的牙是**定向**的（只钉它说要钉的那三形）。
- 同一发里 ban #1-#8 与 ban #9 的 ring 向全绿（36/37 passed）⇒ 摘豁免不会连带把别的格洗红，读数是干净的。
- 反方向也验了（不是只会响）：还原后 37/37 rc=0，且 §3 的 P1 全拼形在两枚二进制下都响。

## §3 格 2 取证：本程自造的探针样本（仓外归档树，零污染共享工作树）

尺：`/tmp/212v2/d22scan-fixed.exe -root /tmp/212v2/tree-head`（`git archive HEAD` 的副本；探针文件种在
`internal/probe212v2*/`，**不在仓内**）；对照发＝同树跑 `d22scan-mut.exe`（豁免摘掉的二进制）。
基线读数（未种探针前）＝该副本只有 1 枚红：`cmd/wisp/panel_host_windows.go:343` 的
`.scratch/wisp/probes/33/r6/fullpack2.log`（**盘上存在但从未被 git 跟踪**，本程现跑：
`ls` 命中 197213 字节、`git ls-files --error-unmatch` → `Did you forget to 'git add'?`）⇒ 归档树复算必带这枚伪影，与实现件 §5 口径 C 的那条⚠具名一致。

| 探针 | 注释里写的 token（逐字） | 期望（按票面 ⓐ 口径） | 实测（fixed／mut） |
|---|---|---|---|
| P1 全拼、盘上没有 | `docs/evidence/s1/does-not-exist-full-name.md` | **必须仍响** | **RING**／RING ⇒ ⓐ 没把真猎物一起豁免掉（本格主判据） |
| P2 混形（中段 `...`） | `internal/xxx/.../foo.md` | ③（缩写＝只规定不问罪） | 静默／RING ⇒ 落 **③**，与 ⓐ 的字面口径相符 |
| P3 完整名＋尾 `…` | `docs/evidence/s1/phantom-tail-ellipsis.md…` | ⓐ 未说；按票面"…那一族"＝③ | 静默／RING ⇒ **过宽实证**：一枚完整的全拼幻影因尾随一个 `…` 免罪 |
| P4 markdown 星号裹名 | `*docs/evidence/s1/phantom-tail-star.md*` | 同上 | 静默／RING ⇒ 过宽实证（`*` 同形） |
| P5 同行粘连吞噬 | `docs/evidence/s1/phantom-swallow-a.md见docs/evidence/s1/...` | 同上 | 静默／RING（mut 报 `docs/evidence/s1/phantom-swallow-a.md见docs/evidence/s1`）⇒ 过宽实证：CJK 在 `repoPathRe` 字符类内，**一段无空格的粘连可把前面的完整幻影一起吃掉** |
| P6 控制（同形无标记） | `docs/evidence/s1/phantom-cjk-adj.md见docs/evidence/s1/subdir` | 必须响 | **RING**／RING ⇒ P5 的静默只归因于那个标记，不是粘连本身 |
| P7 真件＋缩写并存 | `internal/tools/gate.go`（在）／`docs/evidence/s1/152-...-accept.md` | 都静默 | 都静默（fixed）；mut 只红缩写那枚 ⇒ 真件路标不受影响 |
| P8 跨行 | 第 4 行 `docs/evidence/s1/phantom-two-line-a.md` ＋ 第 5 行缩写形 | 第 4 行必须响 | **RING**（`p8line.go:4`）⇒ 豁免按 `c.Text`＝**逐行**计算，不外溢到相邻注释行 |
| P9 空格分隔 | `… docs/evidence/s1/phantom-space-sep.md 见 docs/evidence/s1/... …` | 必须响 | **RING**（`p9space.go:4`）⇒ 区域被空格截断，吞噬只在"同一行且无非类字符隔断"时成立 |
| P10 NTFS 大小写洞 | `Internal/Tools/Gate.GO`（真件的大小写异体） | 本票不管辖（v1 §7 登记） | 静默 ⇒ **v1 那枚残洞仍在**，本程自跑实证；异体大小写命中＝stat 命中即不红 |
| P11 对照 | `docs/Evidence/s1/phantom-case-miss.md` | 必须响 | **RING** ⇒ P10 的静默确系大小写不敏感，不是射程漏了 |

**删样本必不响**（本程自跑的整树控制，不依赖夹具）：`rm -rf internal/probe212v2` 后同树重扫 ⇒
fixed 与 mut 的 phantom token 集合都回到那 1 枚归档伪影，`diff` 两者＝**空**
（`logs/v2-raw-tree-head-deleted.txt`／`-mut.txt`）。⇒ ①种必响、②删必不响在**真树形状**上成立，不只是夹具里。

**豁免今天的实际射程**：clean HEAD 归档树上 fixed 与 mut 的 token 差集＝**0 枚**
⇒ 现网产码注释里**没有一枚**"因 ⓐ 而免罪的全拼幻影"。活着的两处带标记形经本程逐枚核过都不构成免罪：
`cmd/wisp/logsink.go:131` 的 `internal/*` 与 `cmd/wisp/panel_host_windows.go:29` 的 `docs/evidence/s1/...`
——两枚在 `repoPathRe` 里根本截不出"末字符为字母数字"的 token／截出的 `docs/evidence/s1` **是存在的目录**，
所以它们在两发里都静默。唯一活体缩写形仍在 `_test.go`（`cmd/wisp/slo_report_144_windows_test.go:851`
`docs/evidence/s1/152-...-accept-r1.md §2.2`，本程现读该行仍在），而 `walkGo` 在 `main.go:673` 明确跳 `_test.go`
⇒ **ban #9 射程外**，与实现件 §6 #5 的叙述一致。

⇒ 格 2 判语：**带条件成立**。定向性过关（P1/P6/P8/P9/P11 全响），过宽是**真存在的三条形状**（P3/P4/P5），
但（a）实现件在仪器自己的注释 `main.go:877-883` 与 §6 #4 里把它写在了明处、没藏，（b）今日活体 0 枚，
（c）收窄＝射程变更＝人工批准，本程不裁它"跑歪"。

## §4 格 3 取证：三处注释逐枚现读复认（含 D37 名册与 d22scan 执法名册）

| 处（现读行界） | 声称 | 本程现读凭据 | 判 |
|---|---|---|---|
| `internal/agent/tools.go:88-94` | 落账走 `ErrorClassOfTurnError (internal/agent/guard.go:271)` | `internal/agent/guard.go:271` 逐字＝`func ErrorClassOfTurnError(err error) string {`，`:272-275` 为 `observe.ClassOf(err)`→`!ok` 即 `return string(observe.ClassInternal)` | **真**（行号、函数名、回退语义三者全对） |
| 同上 | "the loop books it" | `internal/agent/loop.go:692-694` 逐字＝`case execErr[i] != nil:` / `log.Outcome = OutcomeError` / `log.ErrorClass = ErrorClassOfTurnError(execErr[i])` | **真** |
| 同上 | 无类即回退 class `internal`（`observe.ClassInternal`） | `internal/observe/errors.go:34` `ClassInternal ErrorClass = "internal"` | **真** |
| 同上 | "D37's 17 classes have no `internal-tool` in them — see the roster in `internal/observe/errors.go`" | `errors.go:13` 契约句（"accept these 17 values and nothing else"）；常量块 `:18-36` 共 17 枚（config/auth/network/rate_limit/provider/model/audio_device/asr/tool/permission_denied/user_rejected/cancelled/budget/loop/injection/internal/resource）；`allClasses` 在 `:39`；`ValidateErrorClass` 在 `:61`。`git grep -nE "internal-(tool\|provider)"` 在 Go 产码面只有 2 处命中＝`tools.go:93` 与 `bridge.go:256`，都是**负断言本身** | **真**（名册 17 枚里确实没有这两枚连字符合成名） |
| 同上（附带形状） | 该行把 `"internal-tool"` **带引号**写进注释 | 本程在归档副本 `internal/tools/` 种同形引号串 ⇒ 当场 `[internal-artifact-tool]` 红（见 §6）；`tools.go` 属 `internal/agent/`，在 ban #7 的 `walkText(root/internal/tools)` 射程外（`main.go:260`） | **真而不安全**：它今天不红靠的是**射程边界**，不是措辞。若有人把这枚接口搬进 `internal/tools/`，仪器会红——那是仪器该红的时候，本程判为可接受，具名登记 |
| `internal/tools/bridge.go:251-259` | 非-nil 语义落 `internal/agent/loop.go:692-694`；`internal-provider` 不在 17 类；provider 故障＝class `provider`（`observe.ClassProvider`），"只有错误被按该类包装时才进得了那一行" | `loop.go:692-694` 同上；`errors.go:23` `ClassProvider ErrorClass = "provider"`；`ClassOf`（`errors.go:230-236`）只对携带合法类的 `*Error` 返回 `ok=true`⇒"only if wrapped" 语义成立。旧引 `loop.go:654-657` 本程在 `5e8748b3^` 上现读：`:652-656`＝`timeout := guard.PerToolTimeout()` 与 `req := ToolRequest{...}` 的字面量，`:657`＝`h := l.reg.Spawn("tool-exec-"...)`⇒**都不指向非-nil 分支**，旧锚确为错锚（实现件说"那五行"的行界差 2 枚，实质结论不变） | **真**（新锚逐字对得上；旧锚的"错锚"结论本程复认，行界描述差 2 行） |
| 同上 | "SPEC-07 §2 states the C3 rule as `未声明即拒绝调用（不是报错，是拒绝）`" | `docs/specs/SPEC-07-tools-and-plugins.md:11`＝`## 2. 工具接口与能力面`，`:25` 逐字＝`//                        未声明即拒绝调用（不是报错，是拒绝）` | **真**（章节号与引文逐字对上） |
| `internal/agent/approval/pending_read.go:43-46` | 冻结注真身＝`internal/tools/gate.go:13` | `gate.go:13` 逐字＝`// renders RulesHit and Reason VERBATIM (ticket 17's frozen-contract note), so`（`:12` 起为"confirmation card (tickets 21/37) renders"） | **真**（那一行就是那句冻结契约注，前缀也补齐了） |
| 同上 | "`tools/d22scan` does not enforce it（其 `pathresolver-bypass` ban 管的是 `risk.PathResolver` 之外的 `filepath.Clean/Abs`）" | 本程自跑名册：`roster read from main.go = 9 numbered ban(s) [1 bare-goroutine, 2 pathresolver-bypass, 3 plaintext-key, 4 wallclock-timeout, 5 mirror-hash, 6 panel-approval, 7 internal-artifact-tool, 8 emoji, 9 phantom-citation]`；`git grep -i "ruleshit\|verbatim" tools/d22scan` 在发射逻辑里 0 命中（只有两处无关注释 `main.go:305/:381/:529`）；ban #2 发射点 `main.go:726-727`＝`"filepath."+Sel+" outside the C26 PathResolver is banned (D22)"`，allowlist 里唯一豁免它的条目是 `pathresolver-bypass internal/risk/pathresolver.go` | **真**（不执法的断言成立，射程描述与发射点逐字一致） |
| 同上 | "copied per `LiveApprovals`"／"Position 是 FIFO 里的 1-based 位次" | `internal/agent/approval/pending_read.go:107` `func (q *Queue) LiveApprovals() []LiveApproval`，`:118-121` 逐枚 `Position: q.position(it)`；`internal/agent/approval/queue.go:268-275` 的 `position` 现读＝`return i + 1`（注释自称 1-based） | **真** |
| 同上（**唯一指错处**） | "RulesHit, Reason, **Level 和 SessionOverrideBlocked** are the fields the L2 card shows verbatim, **per ticket 17's frozen-contract note, which lives in internal/tools/gate.go:13**" | `gate.go:13` 只点名 **`RulesHit` 与 `Reason` 两枚**；四枚字段那句话说的是另一处——`cmd/wisp/panel_pump.go:44-46`（"Level/RulesHit/Reason/SessionOverrideBlocked are the fields internal/risk put on the queue item … panel renders them through the same CardViewFromDecision"） | **不成立（出处记错了一枚）**：四枚字段各自为真，但**冻结注只覆盖其中两枚**；把四枚都挂在 `gate.go:13` 上＝引用支撑不了句子。修法＝句中把 Level/SessionOverrideBlocked 的来源改指 `cmd/wisp/panel_pump.go:44`，或把 "per … note" 收窄到两枚。不构成 ban #9 违规（两枚路径都存在），属**注释精度缺陷**，具名上报，改不改归编排者 |

三枚 diff 本体（`git show 5413f46d -- internal/**` 全文读）：**每一条 `+`/`-` 行都以 `//` 开头**⇒
"只改注释、未动产码逻辑"这句话是实话。三枚文件今天都在 ban #9 射程内（`internal/` 产码面），门禁 rc=0
⇒ 它们引用的路径逐枚存在（本程另在 §4 逐行复认了 `guard.go:271`／`errors.go`／`loop.go:692-694`／`gate.go:13`／`panel_pump.go:44`）。

## §5 格 4 取证：计数口径复算（C＝归档树两发差集；D＝hunk 数；A/B 只核存在性）

| 式 | 实现件自报 | 本程读数 | 尺（本程自跑） | 判 |
|---|---|---|---|---|
| C 排除在 | 8 枚 distinct token | **8** | `git archive 5e8748b3^` 副本＋`d22scan-fixed.exe -root` ⇒ `grep -o 'repo path "[^"]*"' \| sort -u`＝8（档 `logs/v2-raw-preimage-fixed-tokens.txt`）：`.scratch/wisp/probes/33/r6/fullpack2.log`, `internal/engines`, `internal/provider`, `internal/tool`, `scripts/check-pathclean-ban.sh`, `tools/agent/cmd`, `tools/bridge.go`, `tools/gate.go` | **相符** |
| C 排除摘 | 9 枚；差集恰 `docs/evidence/s1/62-` | **9**；`comm -13` 差集＝**只有** `"docs/evidence/s1/62-"`；反向差集＝空 | 同树换 `d22scan-mut.exe`（档 `logs/v2-raw-preimage-mut-tokens.txt`） | **相符**（含"其余 8 枚不因①而变"这一句） |
| C 附带两句 | "9 枚＝8 真②＋1 伪影"；"3 枚 API 形两发都静默" | 8 真②名册与 §5 B 逐名相等（含 `docs/evidence/s1/62-`）；伪影＝`fullpack2.log`（本程现跑：盘上有、索引里没有）；`internal/tools.Result.AppliedSteps`／`internal/buildinfo.Name`／`internal/proc.WithRegistry` 在两发 token 集里都**不出现** | 同上 | **相符** |
| D 圈 pathspec | 9 hunk／8 文件（`queue.go`=2，其余各 1） | **9**；逐枚＝pathresolver 1／pending_read 1／**queue 2**／tools.go 1／models.go 1／statevisual.go 1／provenance.go 1／bridge.go 1 | `git show 5e8748b3 -- <8 枚文件> \| grep -c "^@@"` ＋逐枚同尺 | **相符** |
| D 整枚 commit | 35 hunk（含 258-r1 搭车料） | **35** | `git show 5e8748b3 \| grep -c "^@@"` | **相符**；搭车料本程点名复算：`cmd/wisp/resident_ball_windows.go` 7 hunk、`config_readers_255.go` 2、`resident_windows.go` 2、`resident_hotkey_258_test.go` 1、`tools/d22scan/{main.go 3,selftest.go 2,selftestsamples.go 1}` 等 ⇒ 不圈 pathspec 就会把 258 的料算进 212 的账，这一句成立 |
| A／B | A＝10 findings（r1 当日档）；B＝11＝8＋3 | **A 只核到尺的存在性与可重跑性**：`.scratch/wisp/probes/212/r1/post-scan-full.txt` 在库（`5e8748b3` 落入，268 行）；B 的 11 依赖当日工作树的 untracked 面，本程**无法**在归档树上复现（归档树给出的是 8＋1 伪影那一套） | `git show --stat 5e8748b3`／`ls` | **判不动（不是不符）**：B 是"当日工作树"的数，机器相依；实现件自己已写"本腿没直接复现 11"，本程同样不冒充复现它 |

## §6 格 5／格 6 取证：禁区、AC 框、八枚 ban scope 行、自咬

**四枚 commit 的文件名册（逐枚 `git show --name-only --format=`，本程现跑）**

| 枚 | 文件集 | 禁区内命中 |
|---|---|---|
| `7ac8965a` | `.scratch/wisp/probes/212/r2/fix-and-readings.md` | 0 |
| `5413f46d` | `internal/agent/approval/pending_read.go`·`internal/agent/tools.go`·`internal/tools/bridge.go`·`tools/d22scan/main.go`·`tools/d22scan/selftestsamples.go` | 0（`allowlist.txt` 不在册，`runtests.sh`／`scripts/d22scan.sh` 不在册） |
| `1353ce66` | `.scratch/wisp/probes/212/r2/**`（证据件＋25 份读数档＋两台突变脚本） | 0 |
| `d4b2e9b5` | `.scratch/wisp/probes/212/r2/**`（证据件＋6 份终值读数档） | 0 |

- 区间尺（**只作参考，不作单腿凭据**）：`git diff --name-only d934e016..d4b2e9b5` 共 68 行，
  其中命中 `^docs/PLAN\.md|^docs/specs/|^internal/observe/thresholds\.go|allowlist\.txt|^scripts/d22scan\.sh|probes/154/gate-clauses\.sh|^frontend/|^design/|golden` 的＝**0 行**。
  ⚠ 该区间含并发腿的料（`258-r2`/`262-a1`/`174-a4` 的 `.scratch/**` 与 `cmd/wisp/**`、编排者的 README 与台账），
  所以"某枚禁区文件在区间里动过"这种指控**不成立也不该提**；单腿账按上表逐枚看。
  `docs/reports/pending-and-issues.md` 在该区间被**编排者**改过（A588 补框那一枚 `9d089603`），不是 212-r2 的写面。
- **AC 框**：`git log --format="%h %s" d934e016..HEAD -- .scratch/wisp/issues/212-comments-cite-evidence-files-that-do-not-exist.md`＝**0 行**；
  `git diff d934e016..HEAD -- <同一枚>`＝**空**。⇒ 票面六枚框（含第 36 行**没有框**的 AC#4）本轮一字未动。
  本程亦未加框、未改框（那是编排者的账）。

**八枚 ban scope 行逐字（本程现跑，档 `logs/v2-03-scope-lines.txt`；与实现件 `logs/pre-scope-lines.txt` `diff` rc=0）**

```
d22scan: examined 266 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: scope bans #1-5 internal/      examined 228 production Go files
d22scan: scope bans #1-5 cmd/           examined  38 production Go files
d22scan: scope ban #6 frontend/         examined  85 text files
d22scan: scope ban #7 internal/tools/   examined  23 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined  85 text files
d22scan: scope ban #8 internal/         examined 498 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  97 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, bans #1-5 cmd/=38, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=498, ban #8 cmd/=97; ban #8 emoji coverage: design/ 39 text files; frontend/ 85 text files; internal/ 498 Go files, comments and _test.go included; cmd/ 97 Go files, comments and _test.go included
```

（上面十行＝`examined` 行＋八枚 scope 行＋`clean` 行，全部本程现跑逐字抄；与实现件起手档
`diff` rc=0 已证明这十行与**未修码基线**字节相等。）

- ⇒ 格 6 的 `internal/tools` 那侧今天不红＝**证成**（门禁 rc=0 且 `TestScannerSelfScanOfRealRepoIsGreen` 在 runtests 的 34 PASS 里）。
- 与实现件 §8 的起手名册 diff：无差异（它列的八枚＝本程现读的八枚；`cmd/#8=97` 这一枚两腿一致，A578 的 96 未再出现）。
- 自咬机制的本程自证（种在仓外副本 `internal/tools/zz212quote.go`，注释里写带引号的 `"internal-provider"`）：

```
internal/tools/zz212quote.go:5: [internal-artifact-tool] host-internal artifact write looks implemented as a gated tool name (D34 note 2/D22)
```

  删掉该文件后同 tag 计数＝0（档 `logs/v2-raw-probequote.txt`／`logs/v2-raw-probequote-after.txt`）⇒
  实现件 §1 ③.2 记的那口红是**真的会红**，且它的修法是**改措辞**：`allowlist.txt` 未被这轮动过（名册见上表），
  `artifactToolRe`（`main.go:158`）本体一字未改。

**争用检查**：本程在门禁与 `-self-test` 里**没有读到任何落在别人写面上的红**（rc 全 0），
所以本节没有"疑似争用"要报；`internal/ball`（260-r1 在飞）与 `scripts/**`（262-r1 在飞）在 `d22scan: clean` 下零命中。

## §7 AC#4 实质格（票面那一行没有框，判的是"够不够"，不是勾）

| 项 | 凭据 | 判 |
|---|---|---|
| "种一样本必响"＝夹具向 | `tools/d22scan/selftestsamples.go:376-385`（ring 样本引 `docs/evidence/s1/212-citation-ruler.md`，夹具不 seed）；本程现跑逐字＝`internal/probe/cites.go:3: [phantom-citation] comment cites repo path "docs/evidence/s1/212-citation-ruler.md" … ring OK` | **够** |
| "种一样本必响"＝真树向 | 本程在 `git archive HEAD` 副本里种 11 枚自造样本（§3），其中 5 枚**必响的都响了**（P1/P6/P8/P9/P11），10 枚红句逐字在 `logs/v2-raw-probe-fixed.txt`／`probe2-fixed.txt` | **够**（本程独立取数，不依赖实现件） |
| "删样本必不响" | 整树控制：删掉 `internal/probe212v2/` 后重扫 ⇒ phantom token 集回到未种探针前的那 1 枚归档伪影，fixed 与 mut 差集＝空（§3） | **够** |
| 反向准则（删掉样本会让 CI 红，不是让 CI 静） | 名册尺 `selftest.go:332`：`tag %q has only an expect-silent sample - a ban whose violating sample was deleted cannot be distinguished from a ban that was never implemented`；本程现跑名册行＝`9 numbered ban(s) … 37 cases, 10 tag(s) covered, both directions required per tag` | **够**（这一形由仪器自己执行，不是叙述） |
| ⓐ 新增的 silent 样本压得住几形 | 一枚样本压三种拼法（`...`／`*`／`…`），且本程 §2 证明三形同时依赖豁免（摘了就三形都响） | **够**（票面 212-r2 清单①要求的就是"补一枚缩写形静默样本"） |
| 欠的每一格 | §3 的 P3/P4/P5（完整名＋尾标记；同行粘连吞噬）与 P10（大小写洞）**今天都没有样本**；实现件 §6 #4 自述"尾部形没有样本" | **记为缺口，不扣 AC#4 的账**：票面 AC#4 只要求"引用一枚盘上不存在的证据件⇒必响；删⇒必不响"，那一形今天有夹具向＋真树向＋反向准则三重凭据。ⓐ 引入的新形状的样本属"射程变更的配套"，归编排者随「212 改 ⓑ／收窄」那道口令一起开（见 §8 #2） |

⇒ **AC#4 实质：够**（三条判据全有凭据；框本身缺失是编排者的账目缺陷，本程未碰）。

## §8 推翻清单／本程判不动／归后续程

| 枚 | 内容 | 依据 | 归谁 |
|---|---|---|---|
| 1 | ⓐ 的豁免**按区域起点**，故"完整全拼幻影＋尾随 `…`/`*`"与"同一行无非类字符隔断的粘连尾段"三形免罪（P3/P4/P5） | 本程 §3 实测；实现件 `main.go:877-883` 与 §6 #4 已自述；今日活体 0 枚（fixed/mut 差集为空） | **编排者裁**（收窄＝射程变更＝人工批准；与「212 改 ⓑ」同一道口令族）。**不判实现腿跑歪**——它把这条写在了仪器自己的注释里 |
| 2 | 上述新形状**没有样本**压着（`selftestsamples.go:394-404` 只压中段三形） | 本程 §3／§7 | 212-r3 或随 #1 那道裁决定案一起补（补样本＝加测试，不是射程变更，本程可判"该做"） |
| 3 | `shorthandRegionRe` 的前缀表与 `repoPathRe` 的前缀表是**手抄的两份**，没有任何尺钉它们同步（`git grep shorthandRegionRe` 只命中 `main.go`，测试面 0 命中）⇒ 将来给 `repoPathRe` 加前缀而忘了加区域表，AC#3 那一形会**悄悄复发** | 本程现跑 grep＋§2 突变证明豁免是唯一的牙 | 后续程（小料，可并入任何一枚 d22scan 票；本程不动手：改仪器射程＝人工批准面） |
| 4 | ban #9 的射程只覆盖 `internal/`＋`cmd/`（`main.go:250/253` 两次 `walkGo`），**`tools/**` 不在内**⇒仪器自己的源码里可以躺幻影路标：现例＝`tools/d22scan/selftestsamples.go:380` 引的 `docs/evidence/s1/212-citation-ruler.md` 在真仓不存在（本程 `ls docs/evidence/s1 \| grep 212` 只命中本程这枚新件），而两发扫描都没报它 | 本程 §3 基线读数（归档树里 `tools/` 路径 0 命中） | 编排者裁（扩射程＝射程变更）。**不是 212-r2 的账**：它是按 AC#2 的"只扫产码注释面"落地的 |
| 5 | NTFS 大小写不敏感的 stat 残洞仍在（P10 静默／P11 响） | 本程现跑（v1 §7 的〔待验断言〕今天由本程转正） | 登记即可；两任判语与本程同判"量得到的洞、今天无猎物" |
| 6 | `_test.go` 仍在 ban #9 射程外（`main.go:673`），唯一活体缩写形 `cmd/wisp/slo_report_144_windows_test.go:851` 因此不可见 | 本程现读该行＋walkGo 过滤 | 具名登记（实现件 §6 #5 同一句，本程复认）；要不要管＝编排者裁 |
| 7 | `pending_read.go:41-43` 的四枚字段**出处记错**（`gate.go:13` 只覆盖 RulesHit/Reason 两枚；四枚之说在 `cmd/wisp/panel_pump.go:44-46`） | 本程 §4 最后一行 | **212-r3**（一处措辞，属"改成实话"那一枚清单的余数；非射程变更，本程判"该改"但不动手，因为 `internal/agent/approval/**` 不在本程写面） |
| 8 | 行号路标没有仪器管：两例本程**现跑复认**——①`internal/tools/task_output_ac2_before_test.go:23` 引 `bridge.go:247-253` 指"Lookup 失败分支"，现读 `entry, ok := b.reg.Lookup(req.Name)` 在 **`internal/tools/bridge.go:272`**；②`internal/session/grants.go:256` 引 `internal/tools/bridge.go:874` 说"带上无法规范化标记"，现读那一行逐字＝`// happens are clauses of .scratch/wisp/probes/154/gate-clauses.sh - G1/G1b for`，与那句无关⇒同为错锚（顺带具名：`internal/tools/bridge.go:874` 自己引的 `.scratch/wisp/probes/154/gate-clauses.sh` **存在**，48783 字节，本程 `ls` 在案，所以 ban #9 不红它） | 本程现跑 grep/awk | 票 212 射程外（ban #9 只钉存在性）；值不值得立第十枚归编排者开票。实现件 §6 #3 这一格**说的是实话**，本程无异议 |
| 9 | 口径 B 的"11"本程复现不了（归档树给 8＋1 伪影那一套） | §5 表末行 | 不是缺陷：B 是"当日工作树 untracked 面"的数，机器相依；实现件自己已声明"本腿没直接复现 11"。**判不动，具名为判不动** |
| 10 | 票面 AC#4 那一行**没有框**（第 36 行既非 `[ ]` 也非 `[x]`） | 本程现读票面＋`git log` 零改动 | **编排者的账目缺陷**（不是 212-r2 的）；本程未加框。实质判语见 §7 |

## 终态

- 本程判语一句话：**212-r2 的四枚交付里，①②④三枚成立、③带一处具名缺陷（字段出处记错，非幻影、非射程问题）；
  ⓐ 修法有牙且定向，代价（区域起点式豁免的三条免罪形状）已由实现者自己写在仪器源码注释与 §6 #4 里，今日活体 0 枚。**
- 本程写面：`docs/evidence/s1/212-comments-phantom-citation-v2.md`＋`.scratch/wisp/probes/212/v2/logs/**`。
  产码面只在 §2 的成对突变窗口内瞬时改过 `tools/d22scan/main.go`，还原后 `md5` 前后全等、`git diff -- tools/d22scan` 为空；
  §3 的全部"种样本"都种在仓外归档副本 `/tmp/212v2/**`（只建不删），仓内 `internal/`／`cmd/` 零字节改动。
- 票面 AC 框：零碰。台账 `docs/reports/pending-and-issues.md`：零写。禁区面（`docs/PLAN.md`／`docs/specs/**`／
  `internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`scripts/d22scan.sh`／
  `probes/154/gate-clauses.sh`／`frontend/**`／`design/**`）：零触碰。未 push。
- 本程 commit 链：`095cfad8`（骨架：七节标题＋空表）→ 本枚（满稿＋读数档）。
