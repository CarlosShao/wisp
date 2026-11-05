# 179-a1p — 只读收口复核（fix 在 HEAD 上、判据全绿、AC#1 红前读数取到、AC#8 归 183）

- 复核腿：`179-a1p`（只读；**零产码字节**）
- 仓库：`D:\work\workspace\projects plans\Wisp`
- 分支：`dev`
- HEAD：`5cc9fefd4d5f7042e3bc4a90abea6dc3e492182b`（`5cc9fefd4`）
- 现量时刻：`2026-11-05T09:15:02+08:00`
- 票：`.scratch/wisp/issues/179-loop-deciderisk-treats-a-declared-l0-as-unclassified-so-every-readonly-builtin-is-refused-on-the-real-cli-and-no-criterion-holds-that-sentence.md`
- 本腿唯一提交件：本文件。临时件全部 create-only 于 `.scratch/wisp/probes/179/pre/`，未提交。

---

## 1. 修复在场 + HEAD 上判据全绿

### 1.1 `decideRisk` 的 L0 支

```
$ grep -n 'case RiskL0' internal/agent/loop.go
900:	case RiskL0:
```

上下文（`internal/agent/loop.go:891-912`，逐字）：

```go
func (l *Loop) decideRisk(ctx context.Context, j *taskJournal, rowID int64, info ToolInfo) (ok, unbooked bool, why string) {
	risk := info.RiskLevel
	switch risk {
	case RiskL1, RiskL2:
		if l.opt.AdmitTask == nil {
			j.decide(ctx, rowID, DecisionReject)
			return false, false, "该操作属于 " + risk + " 级，审批通道尚未接入（ticket 21），已拒绝执行"
		}
		return true, true, ""
	case RiskL0:
		// A declared tier, the lowest one: run it, interrupt nobody (D4).
		j.decide(ctx, rowID, DecisionAllow)
		return true, false, ""
	default:
		// Reached only by RiskUnclassified ("") and by anything the directory
		// declares no tier for; both fail closed while the switch is off.
		if !l.opt.Config.PassThroughUnclassifiedRisk {
			j.decide(ctx, rowID, DecisionReject)
			return false, false, "风险未分级且直通开关关闭，已拒绝执行"
		}
		j.decide(ctx, rowID, DecisionAllow)
		return true, false, ""
	}
}
```

判定：`"L0"` 现在走**显式的 `case RiskL0`**，不再落进 `default:`。修复在场。

### 1.2 编译

```
$ go build ./internal/agent ./internal/tools
build rc=0
```

（无输出，rc=0。）

### 1.3 判据改名册

```
$ go test ./internal/agent -run '179|Declared' -count=1 -v -timeout 200s
=== RUN   TestDeclaredL0RunsWhenPassThroughIsOff179
--- PASS: TestDeclaredL0RunsWhenPassThroughIsOff179 (0.15s)
=== RUN   TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179
--- PASS: TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179 (0.12s)
=== RUN   TestDeclaredL1L2RoutingUnchanged179
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L1/no-approval-layer
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L1/with-approval-layer
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L2/no-approval-layer
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L2/with-approval-layer
--- PASS: TestDeclaredL1L2RoutingUnchanged179 (0.13s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L1/no-approval-layer (0.01s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L1/with-approval-layer (0.00s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L2/no-approval-layer (0.00s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L2/with-approval-layer (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/agent	0.436s
test rc=0
```

- 顶层 3 枚全 PASS（C1 正向、C2 反向守卫、C3 L1/L2 守卫）；C3 下 4 枚子测全 PASS。
- 同一发读数另存：`.scratch/wisp/probes/179/pre/fixed-green-179.log`（`green rc=0`）。

### 1.4 格式

```
$ /d/work/base/gopath/bin/gofumpt.exe -l internal/agent/loop.go internal/agent/declared_l0_risk_179_test.go
gofumpt rc=0
```

（无文件名输出 = 两枚都干净。）

### 1.5 d22scan

```
$ sh scripts/d22scan.sh
...
runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0
d22scan: examined 408 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: scope bans #1-5 internal/      examined 353 production Go files
d22scan: scope bans #1-5 cmd/           examined  55 production Go files
d22scan: scope ban #6 frontend/         examined 116 text files
d22scan: scope ban #7 internal/tools/   examined  47 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined 116 text files
d22scan: scope ban #8 internal/         examined 865 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined 201 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; ...
d22scan rc=0
```

- 正反两向都跑：`tools/d22scan` 自测 `PASS=35 FAIL=0`；真实仓库扫描 `clean`、rc=0。

---

## 2. AC#1 红前读数（当前树、零产码）

### 2.1 第一步：票面指定的 `039efb47` 整枚 blob 不能直接 overlay

按票面`09-28 12:2x`补法取修法前那一版：

```
$ git cat-file blob 039efb47:internal/agent/loop.go > .scratch/wisp/probes/179/pre/loop.go
$ wc -c .scratch/wisp/probes/179/pre/loop.go
38376
$ grep -n 'case RiskL0' .scratch/wisp/probes/179/pre/loop.go
（无输出，rc=1）
```

该版确为修法前（无 `case RiskL0`；`风险未分级` 在 `:785`）。但把它整枚 overlay 到**当前树**编译失败，原因是它与 HEAD 之间已有 174 行插入/22 行删除的**无关漂移**（`PointerJudge`、`MaxOutputTokens`、`costMeter` 等）：

```
$ go test -overlay=.scratch/wisp/probes/179/pre/overlay-prefix.json ./internal/agent -run '179|Declared' -count=1 -v -timeout 200s
# github.com/CarlosShao/wisp/internal/agent [github.com/CarlosShao/wisp/internal/agent.test]
internal\agent\costbook_44.go:63:7: l.costMeter undefined (type *Loop has no field or method costMeter)
internal\agent\costbook_44.go:70:16: l.costMeter undefined ...
...
internal\agent\corr_roster_join_328_test.go:116:12: undefined: callCorr
internal\agent\corr_roster_join_328_test.go:116:12: too many errors
FAIL	github.com/CarlosShao/wisp/internal/agent [build failed]
rc=1
```

⇒ 这一形**不是判据红**，是 build failed；不能拿来当 AC#1 读数。

### 2.2 第二步：当前树只摘 L0 那一支（行为等价于修法前）→ 红

为把"红在**将被修改的那棵树上**"这件事做实，同时不让无关漂移进来，取 **HEAD 的 `loop.go`，只删 `900-903` 四行（`case RiskL0:` 及其三行体）**：

```
$ sed '900,903d' internal/agent/loop.go > .scratch/wisp/probes/179/pre/loop-head-minus-l0.go
$ grep -n 'case RiskL0' .scratch/wisp/probes/179/pre/loop-head-minus-l0.go
（无输出，rc=1）
$ grep -n '风险未分级' .scratch/wisp/probes/179/pre/loop-head-minus-l0.go
905:			return false, false, "风险未分级且直通开关关闭，已拒绝执行"
$ wc -c .scratch/wisp/probes/179/pre/loop-head-minus-l0.go
47358
```

overlay 只映射 `internal/agent/loop.go` 这一枚（见 `.scratch/wisp/probes/179/pre/overlay-head-minus-l0.json`），跑当前树的 179 判据：

```
$ go test -overlay=.scratch/wisp/probes/179/pre/overlay-head-minus-l0.json ./internal/agent -run '179|Declared' -count=1 -v -timeout 200s
time=2026-11-05T09:14:17.802+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestDeclaredL0RunsWhenPassThroughIsOff179
    declared_l0_risk_179_test.go:101: executions = 0, want 1: a declared L0 call must reach the provider
    declared_l0_risk_179_test.go:104: decision = "reject", want "allow"
    declared_l0_risk_179_test.go:107: outcome = "error", want success
    declared_l0_risk_179_test.go:110: error_class = "permission_denied", want none
    declared_l0_risk_179_test.go:120: recorded calls = 0, want 1
--- FAIL: TestDeclaredL0RunsWhenPassThroughIsOff179 (0.15s)
=== RUN   TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179
--- PASS: TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179 (0.15s)
=== RUN   TestDeclaredL1L2RoutingUnchanged179
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L1/no-approval-layer
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L1/with-approval-layer
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L2/no-approval-layer
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L2/with-approval-layer
=== NAME  TestDeclaredL1L2RoutingUnchanged179
    declared_l0_risk_179_test.go:252: declared L0 with a gate wired: ok/unbooked/why = false/false/"风险未分级且直通开关关闭，已拒绝执行", want true/false/empty
    declared_l0_risk_179_test.go:262: decision = [{ID:5 TaskID:task-179-l0-direct Seq:1 Tool:fs.read ArgsJSON:{} RiskLevel:L0 Decision:reject ...}], want one row with "allow" booked by the direct L0 arm
--- FAIL: TestDeclaredL1L2RoutingUnchanged179 (0.17s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L1/no-approval-layer (0.01s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L1/with-approval-layer (0.00s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L2/no-approval-layer (0.00s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L2/with-approval-layer (0.00s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/agent	0.515s
FAIL
rc=1
```

**逐字拒绝行**：

```
"风险未分级且直通开关关闭，已拒绝执行"
```

**尾态**（同一发）：`decision=reject`、`outcome=error`、`error_class=permission_denied`、`executions=0`、`recorded calls=0`。

**同发守卫**（证明 overlay 只摘了 L0 支、没有把门洗掉）：

- `TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179` → **PASS**（真未分级仍拒）；
- `TestDeclaredL1L2RoutingUnchanged179` 的 L1/L2 四枚子测 → **PASS**。

完整日志：`.scratch/wisp/probes/179/pre/overlay-red-179.log`（rc=1）。

### 2.3 AC#1 判读

- **实际取得：overlay-red（当前树）**。上句 `风险未分级且直通开关关闭，已拒绝执行` 就是 AC#1 要的那句逐字拒绝行；`decision=reject / outcome=error / permission_denied` 是尾态。
- 这一发是**当前树 + 只缺 L0 支**的 in-package overlay，不是票面原句里 `176-r1`/`probes/179/r2` 的整发 CLI e2e；本腿按派单允许的"smallest possible in-package check"取数。**没有**使用 cf527f32 的旧读数充当（旧读数只作佐证：`grep -ac '风险未分级' .scratch/wisp/probes/176/r1/logs/last-run.txt` ⇒ **7**）。
- 票面 AC#1 的框**仍未勾**：本腿是只读复核，不改票面/不翻任何框。事实已登记在此。

---

## 3. 票 183 状态（AC#8 的归宿）

- 盘上文件名：

```
$ ls .scratch/wisp/issues/183-*.md
.scratch/wisp/issues/183-host-pointer-exemption-fails-real-cli-done.md
```

- `Status` 行**逐字**（`183-*.md:3`）：

```
- Status: **ready-for-agent（先只读核，再落地）**。⚠ 本票是**缺陷票**，不是新功能票：判据在票 177 里已经写死并翻过勾（AC#2 反向／AC#3 正向〔成立·带条件〕），**本票要查的是"那枚豁免为什么在真机接缝上没生效"**。
```

⇒ 文件名已是 `-done`、十枚 AC 全 `[x]`，但 `Status` 行**没跟着改**（仍是 `ready-for-agent`）。这是票面自相矛盾的一处，登记为 caveat；按文件名 + 框 + Progress log 判，**183 已结线**。

- AC 框：`grep -nE '^- \[[ x]\] \*\*AC#'` ⇒ AC#1、#2、#3、#4、#5、#6、#7、#8、#2b、#9 **十枚全 `[x]`**。
- Progress log 末段逐字（`14:1x`）：

```
> **⇒ 票 183 十格全勾、结线改名 `-done`。** 残余两笔都归了口：**第二次续读＝票 185** ...
```

- 真机 e2e 读数（183 Progress log，`14:1x` 编排者自跑，带在册 PATH 前缀）：

```
第一发续读 kind=success risk=L0 decision=allow rules_hit=[]
第二发 kind=refused risk=L2 rules_hit=[R4] reason="R4: 包含来自 fs.read …"
```

⇒ **179 AC#8 要的那一发（照宿主指针续读、不被自家门拒、逐字节读回）在 183 已settled**：第一发 `fs.read` 成功、`risk=L0 decision=allow rules_hit=[]`、20000 字节逐字节回全。第二次有界续读仍拒，但那是**另一枚缺陷**，已归 **票 185**（盘上 `185-*.md` 存在，Status 逐字 `已立、按住不派`）。179 AC#8 按编排者 09-28 的转移裁定**不在本票单独跑 e2e、框保持不勾**。

---

## 4. 179 未勾框恰为 {AC#1, AC#8}

```
$ grep -nE '^- \[ \]' .scratch/wisp/issues/179-*.md
45:- [ ] **AC#1 红之前的读数不许"待跑"** ...
58:- [ ] **AC#8 端到端解封——本票存在的唯一目的** ...
```

- 未勾：**恰两枚，AC#1 + AC#8**。
- 已勾：`grep -cE '^- \[x\]'` ⇒ **8**（AC#2/#3/#4/#5/#6/#7/#9/#10）。
- 对应判据在 HEAD 全 PASS（见 §1.3）：AC#2 = `TestDeclaredL0RunsWhenPassThroughIsOff179`；AC#3 = `TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179`；AC#4 = `TestDeclaredL1L2RoutingUnchanged179`（含 L1/L2 四子测）；AC#9 = 同文件 `:262` 的 `decision` 列断言（本轮绿）；AC#10 = `loop.go:885-890` 的注释已改成实话（`The only caller (executeCalls) currently DISCARDS that value ...`）。

---

## 5. 零产码字节自证

```
$ git status --porcelain -- internal/ cmd/
?? cmd/wisp/zz_probe_temp_test.go.txt

$ git diff --name-only -- internal/ cmd/
（空）

$ sha256sum internal/agent/loop.go
b6600a0977370b15c34ae5cd3965bd539255d2caf8f5fffcd4630a7a12929f17 *internal/agent/loop.go
```

- `cmd/wisp/zz_probe_temp_test.go.txt` 是**起手就存在**的未跟踪件（不是本腿创建）；本腿未碰 `internal/`、`cmd/` 任何字节。
- 本腿只写 `.scratch/wisp/probes/179/pre/` 下临时件（overlay JSON、两枚 loop 副本、两枚日志）与本证据件。

---

## 6. 结论

- **closeable：YES**（编排者把 AC#1 按本读数翻勾、并按 09-28 转移裁定把 AC#8 记在 183 名下之后，本票可结）。
  - 严格说：AC#1 与 AC#8 的框**现在仍未勾**；本腿只读，未动票面。
- **fix-present + tests-green：YES**。
  - `case RiskL0` 在 `loop.go:900`；`go build` rc=0；179 判据 3 顶层 + 4 子测全 PASS；gofumpt 空；d22scan rc=0。
- **AC#1 读数（实际取得）：overlay-red（当前树）**。
  - `风险未分级且直通开关关闭，已拒绝执行` 逐字出现；`decision=reject`、`outcome=error`、`error_class=permission_denied`。
  - 票面指定的 `039efb47` 整枚 blob 在当前树**编译失败**（无关漂移：`costMeter`、`callCorr` 等），故改用"HEAD 只摘 `case RiskL0` 支"的行为等价副本；两向都如实记录。
  - 未使用 cf527f32 旧读数冒充；旧读数 `grep -ac` = 7 仅作佐证。
- **183 状态：`-done`**（文件名、十枚 AC 全勾、Progress log 明写结线），但 `Status` 行仍是 `ready-for-agent`（stale，caveat）。AC#8 的 e2e 已在 183 settled（第一发 `kind=success risk=L0 decision=allow rules_hit=[]`、20000 字节回全）；第二次有界续读归票 185（open）。
- **未勾框：恰 `{AC#1, AC#8}`**。

### Caveats

1. `183-*.md` 的 `Status` 行未随 `-done` 改名更新（仍 `ready-for-agent`），与文件名/框/Progress log 矛盾。
2. AC#1 的读数形状是**当前树上的 in-package overlay**，不是整发 `176-r1` CLI e2e；派单允许该较小形状，但若编排者要求票面字面的 CLI 尾态，仍需一枚带 PATH 前缀的 r2/176 台件复跑。
3. AC#8 的框在 179 上仍空，属**有意转移**；其证据在 183，第二次续读的残余缺陷在 185（`已立、按住不派`）。
4. 本腿未跑 `cmd/wisp` 包（票 98 的 `0xc0000135` 环境坑与本次收口无关），也未 push。
