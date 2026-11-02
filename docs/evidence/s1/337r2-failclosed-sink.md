# 票 337 落地腿 `337r2` · 证据件 — 注入式宿主装配点补接 `risk.Logf`（残余那一接）

- 腿：`337r2`（写腿／实现者，`cmd/wisp` 面）。类型：接线（一行赋值＋注释），承 `337-r1` 之后票面唯一残余。
- 残余的具名出处＝`337-r1` 自己在 `internal/agent/approval/gate.go:175-179` 写下的那段：
  > "RESIDUAL, stated not hidden: a host that injects its own gate (run.go's s.gate != nil branch) never calls this constructor, so risk.Logf stays unwired there"
- 派单存档：`.scratch/wisp/probes/orch/20261102/txt/337r2-dispatch.txt`。主题：`票 337 落地腿 337r2`。
- 本腿**只**动 `cmd/wisp/**` ＋ 本件；`gate.go` 的 337-r1 逻辑**一字节未动**（尺见 §5）。

---

## 0. 起手三核

| 核 | 读数（现跑） | 判 |
|---|---|---|
| ① HEAD 在 `dev` ≥ `76fe90699` | `git log -1 --format='%H %d'`＝`76fe906997f7f824ad493ce1e3c5b3cef5c122ab  (HEAD -> dev, origin/dev, origin/HEAD)` | 通过 |
| ② 索引干净 | `git status --porcelain cmd/wisp internal/agent/approval internal/risk`＝**空输出** | 通过 |
| ③ 写面无他人脏 | 同上三面对 `cmd/wisp`＝空；编排者那一枚提交（`76fe90699`）正文逐字记"写腿只 1 条原因具名＝cmd/wisp 每波 1 条"⇒本腿是该面唯一在飞写腿 | 通过 |

起手三核文件（任务要求的"读三处"）：`cmd/wisp/run.go:691-783`（两 posture 分支：`s.gate == nil` 的谓词补全 691-693、注入分支 695-759、console 分支 760-783）／`cmd/wisp/resident_approval_windows.go:399-415`＋`:574-578`（宿主自己那枚 `Logf: ra.residentAuditf` 与它的 sink 形状）／`internal/agent/approval/gate.go:148-182`（337-r1 的接线与它自陈的残余）。

**未定义即停未触发**：① 残余**确实**没被 337-r1 覆盖（注入分支不调 `approval.New`，`gate.go:180-182` 那条接线在该分支根本不执行）；② 补这一接**不需要**改 `gate.go`；③ 判据是真实 `risk.Gate` 运行＋真实 `assembleRuntime`＋真实文件的读回，不是 mock 假报（见 §2）。

---

## 1. 接线落点（唯一一处产码改动）

`cmd/wisp/run.go`，注入式宿主分支（`if s.gate != nil { ... }`）末尾、`rt.liveCards = &nativeCards{h: s.cards}` 之后：

```go
risk.Logf = rt.auditf
```

- 位置读数：`git diff --numstat -- cmd/wisp/run.go`＝`31 0`（30 行注释＋**1 行代码**，删除列**0**）；该行现量在 `cmd/wisp/run.go:758`（尺＝`grep -n "risk.Logf = rt.auditf" cmd/wisp/run.go`）。
- 为什么是**赋值而不是 `if risk.Logf == nil` 守卫**（这一格是本腿的判语，写清落点）：注入式宿主自己那枚门**几乎必然已经**把 sink 接上了——`newResidentApprovalWithConfig` 走 `approval.New(approval.Options{... Logf: ra.residentAuditf ...})`（`resident_approval_windows.go:404`），337-r1 的 `gate.go:180-182` 因此已经把 `risk.Logf` 置成宿主那一行。守卫式写法在此**恰好**会让 AC#1 落空（实话继续只落宿主 stdout/slog 那一支，永远落不到本装配的 `rt.auditf`）；而 `approval.Gate` 把 `Options.Logf` 存成**只写字段**（`gate.go:180-182` 不回读），装配根**没有**任何 API 能把宿主那枚 sink 读回来比较。所以能兑现"注入式宿主下实话能落 `rt.auditf` sink"这句话的形状只有一种：装配根按与 console 腿**同一个值**（`Options.Logf: rt.auditf`，`run.go:767`）重新接上。运行现场由装配根拥有 ⇒ 安全 sink 由它接，两条分支一个 sink 家族。
- ⛔ 不是：授权、放宽、出向面。只喂**本机审计日志**（`agentRuntime.auditf`＝`[audit] ` 到 stderr ＋ `audit: ` 记录进本 run 的 sink，`run.go:1163-1167`＋`logsink.go`），⛔ 不进模型正文、⛔ 不进面板字节，也不能允许或拒绝任何东西。

---

## 2. AC#1 接线成立（新判据件，两枚用例）

新件＝`cmd/wisp/risk_sink_injected_337r2_windows_test.go`（`//go:build windows`）。台件＝真 `assembleRuntime` ＋ 真注入门（`newResidentApproval`）＋ 真 `risk.Gate` ＋ 真文件；注入面正是 AGENTS.md §1.3 点名的 CLI 接缝（`runSpec.gate/ui/cards/answer`），⛔ 无 mock 顶真件。

| 用例 | 测什么 | 读数 |
|---|---|---|
| `TestTicket337R2InjectedGateArmsRiskSinkOnRuntimeAudit` | **强形**：宿主**已**用另一枚 sink（测试内立的 foreign sentinel）接过 `risk.Logf` ⇒ 装配必须**改接**到 `rt.auditf`，且实话落本 run 的 stderr 审计面 | `--- PASS (1.59s)`；断言含"foreign sink 收行数必须为 0" |
| `TestTicket337R2InjectedGateArmsRiskSinkWhenHostNeverArmed` | **残余原形**：`risk.Logf` 起手为 nil（先把宿主 `approval.New` 那次接线清掉，确定性基线）⇒ 装配必须接上 | `--- PASS (1.33s)` |

落进 sink 的两条实话原文（用例里 `t.Logf` 打出的那一条，逐字）：```
[audit] risk: A-list DENY path=c:\users\...\001\profile\work\id_ed25519 rule=~/.aws/credentials [fail-closed, NOT a match on this rule]: the candidate ... is not proven handle-expanded (its own walk stopped at ...) and the anchor ... is not installed here (its expansion stopped at ...), so R17 protects the ancestry chain the two share instead of proving a miss (non-overridable)
```

同发另一条（保证型 A 拒，规则字面量不带 fail-closed 从句）：`[audit] risk: A-list DENY path=...profile\.ssh\id_installed rule=~/.ssh/** (non-overridable)`。⇒ 票 334 买到的那句 fail-closed 实话，现在在**注入式宿主**这条路径上也能被审计人读到。

尺（可逐字重跑，必须带 DLL PATH，见 §4 仪器坑）：

```
PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:/d/work/workspace/projects plans/Wisp/build:$PATH" \
  go test ./cmd/wisp -run "TestTicket337R2InjectedGate" -count=1 -timeout 30m
```

---

## 3. AC#2 反向判据（摘接线 ⇒ 红 ⇒ 还原核）

突变＝**只删** `run.go:758` 那一行（注释保留＝突变面恰 1 行，`758d757`），走**仓外 `-overlay`**，⛔ 不碰共享树：

- 件：`D:\tmp\wisp\20261102\337r2\mut\run.go`（＋`run.go.before` 快照、`overlay.json`）。
- 红读数（`go test ./cmd/wisp -overlay=... -run TestTicket337R2InjectedGate`，`rc=1`）：
  - `--- FAIL: TestTicket337R2InjectedGateArmsRiskSinkOnRuntimeAudit (1.45s)` ⇒ `risk_sink_injected_337r2_windows_test.go:184: an A deny on the injected-host assembly did not reach the runtime audit sink; got:`
  - `--- FAIL: TestTicket337R2InjectedGateArmsRiskSinkWhenHostNeverArmed (1.47s)` ⇒ `risk_sink_injected_337r2_windows_test.go:244: the injected assembly did not arm risk.Logf: the fail-closed reason is still dropped on this path (ticket 337 residual)`
  - ⇒ 两枚**各自红在自己那一格**（不落 sink／没接线），不是同一条红被数两遍。
- 还原核：`cmp -s cmd/wisp/run.go $M/run.go.before`⇒**`restore_cmp=IDENTICAL`**；`grep -n "risk.Logf = rt.auditf" cmd/wisp/run.go`⇒`758` 行仍在；还原后同一条命令⇒`ok github.com/CarlosShao/wisp/cmd/wisp 2.726s`，`rc=0`。

---

## 4. AC#3 零回归

| 尺 | 读数 |
|---|---|
| `go build ./cmd/wisp` | `rc=0`（改前基线同 `rc=0`，件 `D:\tmp\wisp\20261102\337r2\build-baseline.txt`／`build-after.txt`，两枚皆空输出） |
| `go test ./internal/agent/approval -count=1` | `ok github.com/CarlosShao/wisp/internal/agent/approval 0.469s`，**FAIL=0** |
| `go test ./internal/risk -count=1` | `ok github.com/CarlosShao/wisp/internal/risk 4.649s`，**FAIL=0**（未撞 `TestResolvePerCallBudget` 那枚计时族红） |
| 按名作差 | 定向 `-run TestTicket337R2InjectedGate` 本发**只**发出这 2 枚名、2 枚全绿；本腿未删/未改任何既有用例名（写面 §5 可核）。⛔ 追 `cmd/wisp` 整包绿＝按派单口径。 |
| **台面变更（具名，非本腿的伤）** | 11:10 起共享树多出**两枚别腿在飞脏件**＝`internal/config/migrate.go`（`"os" imported and not used`，票 345 面）与 `internal/ball/tokens.go`（`Palette` 缺 `GlassRing/GlassHi/GlassHiSoft` 字段，票 344 面），二者把 `go build ./cmd/wisp` 拖成 `rc=1`＋`[build failed]`。⛔ 不属本腿写面、⛔ 动它。⇒ §2/§3/§4 的**终态读数**一律改在"仓外 `-overlay` 把这两枚指回 HEAD"的基线上重取（manifest `D:\tmp\wisp\20261102\337r2\bench\overlay-head-two.json`）：`build rc=0`／`vet rc=0`／定向 `ok 4.410s`；摘线突变同基线仍 `rc=1` 两枚各自红。**先前那发（无 overlay、树还没脏）的读数一并留档**，两发同色＝不是台面形状。 |
| 仪器坑（在册，具名留档） | 裸 `go test ./cmd/wisp`（不带 sherpa/build DLL 的 PATH）＝`exit status 0xc0000135`＋`FAIL 0.056s`＋**零枚 `--- FAIL`／`=== RUN`**＝那包**根本没跑**，不是红。正解＝`PATH` 用 shell 自己的 POSIX `/d/…` 形（`scripts/wisp-cli-tests.sh:212-220` 逐字记了 `D:/` 形反而 `0xc0000135`）。另记一枚自踩：overlay manifest 里的替换路径必须是 `D:/tmp/...` 形，写成 `/d/tmp/...` 会被 Go 工具链解析成 `D:\d\tmp\...` 而报 `open ...: cannot find the path`（本腿第一发就是这么红的，`sed` 补 `g` 后修正＝台件形状问题，⛔ 判据）。本腿所有 `cmd/wisp` 读数都在带该前缀的台面上取；`GOTMPDIR` 一律 `unset`（只 unset，不另设）。 |

---

## 5. AC#4 纪律格

| 尺 | 读数 |
|---|---|
| 写面 | `git status --porcelain` 里本腿产出恰 2 枚＝`M cmd/wisp/run.go`＋`?? cmd/wisp/risk_sink_injected_337r2_windows_test.go`，加本件 `docs/evidence/s1/337r2-failclosed-sink.md`。余下 `??` 全是别腿/环境的 `.scratch/**`、`design/**`、`logs/**` 未跟踪件，本腿⛔ 不 add（两处 `add`／`commit` 均显式 pathspec）。 |
| ⛔ 改 337-r1 已落逻辑 | `git status --porcelain -- internal/agent/approval internal/risk`＝**空**；`git diff --numstat -- internal/agent/approval/gate.go`＝**空**（一字节未动）。 |
| ⛔ 禁区 | `thresholds`／`golden`／`frontend`／`design`／`speech`／`audio`／`testdata`＝零命中（写面只有上表 2＋1 枚）。 |
| `gofumpt` | `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/run.go cmd/wisp/risk_sink_injected_337r2_windows_test.go`＝**空输出**，`rc=0`。 |
| `d22scan` | `bash scripts/d22scan.sh`⇒`d22scan: clean - no D22 ban violations`，`rc=0`；分母读数 `ban #8 cmd/=186`（含本腿新件）、`bans #1-5 cmd/=52`、`ban #8 internal/=736`。 |
| ⛔ 翻框／改票面／台账／push | 本腿零操作；只 commit 不 push。 |

---

## 6. 残余与边界（照实写，不当已闭合）

1. **本腿闭合的是"注入式宿主这条装配路径上 `risk.Logf` 没人接"**这一格。`gate.go:175-179` 那段"RESIDUAL"文字因此**已过期**——本腿⛔ 改它（属 337-r1 面），登记为**归 `337-r1`/验收腿的一句话更正**：残余今天在 `cmd/wisp` 侧已补，注释里那句"stays unwired there"应改成"the console leg arms it here, the injected leg arms it at the assembly root"。
2. **多门进程的最后写者**：`risk.Logf` 是包级单点，谁最后构造谁定 sink。生产上装配根（`assembleRuntime`）跑在任何宿主门构造**之后**，所以本装配的 sink 定盘；一台件里若在同进程再造一次注入式装配，最后那一枚赢。这是既有形状（337-r1 在 `approval.New` 里同样是直接赋值），本腿未加宽也未收紧。
3. **⛔ 端到端**：本腿读数全在 `go test` 接缝上（真 `assembleRuntime`＋真 `risk.Gate`），**不是**真起 `wisp.exe` 常驻进程读它自己的 `<data>\logs\*.jsonl`。判语只能写到"注入式宿主装配路径上这句实话会进本 run 的审计 sink"，写不到"双击那发常驻进程今天日志里就有这行"。
4. `cmd/wisp` 整包终态名册⛔ 归本腿（派单口径＋在册死形 `MARK299/WebView2` 族），AC#3 那一格按名作差。
