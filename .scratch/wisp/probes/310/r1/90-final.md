# `90-final.md` — 腿 `310-r1` 结论（实现者；判语⛔ 由本腿下）

**做了什么**：票 310 `AC#1` 落地形＝甲-G（编排者 10-12 裁形）。在三枚会起 `go test` 的 `scripts/**` 入口件里，跑测期把 **`GOTMPDIR`** 指到一枚 `%TEMP%` 的**兄弟**专用根（`%LOCALAPPDATA%\wisp-test-gotmp\<入口>-<时间>-<pid>`），起手建、收尾只把路径打进读数件、⛔ 删任何东西；第四枚名册件 `scripts/slo-check.ps1` 经现跑证明**零枚 `go test` 调用点**⇒ 本腿⛔ 在那一枚里造步子（具名顶回，见 `00-anchor.md` §5）。

**买到了什么**（读数⛔ 在本件，全在 `20-paired-readings.md`／`30-gates.md`）：
1. 撤掉那一步 ⇒ 死亡形当场在 `%TEMP%` 顶层多一枚（超时 panic 一发 `TestRoundTripLoadMarshalLoad2461500944`；外力杀两发 `…2236627062`／`…826293033`），每一发同行带 `panic: test timed out` 或 `target_pid_still_present=0`。
2. 那一步在位 ⇒ 顶层 0 新增、残渣逐字落在专用根里（`…\portable-20261012-115401-2013\TestRoundTripLoadMarshalLoad1761995879` 等三枚），同一批发同时给出"那一形真发生了"的行。两发作废的发（a1、b2）本腿自己标了、留了非零 rc，⛔ 拿它们的 0 当凭据。
3. 门禁：`sh scripts/d22scan.sh` rc=0；`sh scripts/check-path-length-budget.sh --with-self-test` rc=0；`--scope=census` totals 行跑前（HEAD 导出树）⇄ 跑后（工作树）`diff` rc=0 逐字未变；winsec 档改前⇄改后同色（58 PASS／0 FAIL）；cli 档改前⇄改后同红同三名。
4. ★**本腿撞出一条派单⛔ 写的代价**：专用根**落哪一枚卷不是自由变量**。第一稿落在仓所在卷（`D:/tmp`）当场把 winsec 那档打成 `56 PASS／2 FAIL`（`TestGateJudgesThePrivateSetByResolvedSID` 两枚子例＋`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`，判的是 ACL 继承自哪儿）；换到与 `%TEMP%` 同卷同 ACL 族的 `%LOCALAPPDATA%` 兄弟根后同数复绿。⇒ 派单/票面那句"`os.TempDir()`⛔ 读 `GOTMPDIR` ⇒ 那 17 枚的答案不变"在**词面层面对、在颜色层面对⛔ 成立**；丙那一行我⛔ 动（`gate-roster-304.md` 一字未碰）。

**三张表指回**：表①（改了哪几枚＋调用点行号＋命名规则）＝`10-impl.md` §1/§2；表②（成对两发＋作废两发＋本腿造的目录名册）＝`20-paired-readings.md` §1/§2/§4；表③（门禁＋格式名册两把＋totals 作差＋入口件跑法/颜色）＝`30-gates.md` §1–§4。

## 1. ★本腿欠的读数（逐枚编号；⛔ 这些⛔ 成立，只写"没跑"）

- `Q-310-r1-1`＝**派单点名的可复用台件⛔ 适用本票**：`probes/orch/20261011/shape-c-orch.ps1`（38 行）判"世界造出来了"靠盯一枚带本仓 golden 路径的 `mockllm.exe` 孩子（票 304 进程面形状），`internal/config` 那枚用例⛔ 起子进程 ⇒ 原样跑读不回任何形。本腿改用普查腿 `20-shape-c-watch.ps1` 的派生件（`logs/10c-shapec-watch.ps1`，两处改动写在件头）。⛔ 下一程再拿 304 那枚台件量临时目录这一族。
- `Q-310-r1-2`＝**第五枚入口⛔ 在名册上、⛔ 在本腿落点**：`scripts/d22scan.sh` 经 `tools/d22scan/runtests.sh:75` 起真 `go test`，本票四枚名册⛔ 含它 ⇒ 走那一条路死掉的测试二进制**照旧**在 `%TEMP%` 顶层留目录。⛔ 本腿自行扩面（派单写"唯一落点＝那四枚"）；要不要把 `scripts/d22scan.sh`（或 `tools/**` 那枚 runner）纳入射程＝编排者一句话的事。
- `Q-310-r1-3`＝**托管 runner／CI 那一侧⛔ 本腿能交**（归推送批次，且本腿⛔ push）：① ubuntu 腿——本腿用 `go env GOOS=windows` 把这一步平台门住了，非 Windows 只打一行 `left as-is`，所以 CI 的 ubuntu 档行为一字未改，但那一族残渣在 `/tmp` 继续长、无人管；② windows runner——同一段代码会在 `C:\Users\<runner>\AppData\Local\wisp-test-gotmp` 建根，本腿⛔ 能测那里可不可写，所以降级路径（`DEGRADED` 那两行、⛔ 翻红）是唯一⛔ 被实测过的分支。裁决者若要 fail-closed，改点在 `10-impl.md` §2 第 4 条。
- `Q-310-r1-4`＝**格式名册那一处 4⇄0 的差⛔ 由本腿归因**：四枚 tracked `.go` 在工作树里 `gofmt -l` 命中、在 HEAD blob 抽取树里零命中，而 `git diff --name-only HEAD -- '*.go'`＝0 枚 ⇒ 差在检出面形状（行尾族），本腿没深挖；另本机⛔ `gofumpt`、本腿用 `gofmt` 代尺（替换记录在 `30-gates.md` §2）。CI 那把 gofumpt 的颜色⛔ 本腿能交。
- `Q-310-r1-5`＝**第四枚名册件⛔ 跑**：`scripts/slo-check.ps1` 整枚没跑（它要出厂 `wisp.exe`＋采样，真窗／音频那一族要机主在场），依据只是"它零枚 `go test` 调用点、本腿⛔ 动它一根指头"。
- `Q-310-r1-6`＝**⛔ 跑过的死亡形还有两枚**：① shell 超时／Ctrl+C 那一族（`SIGINT` 给 `go test` 而非测试二进制本身）；② 机器重启／进程被回收之后 `t.TempDir()` 根里那批残渣的**再归属**问题（本票甲-G⛔ 治"没人收尾"，票面 §4 已写）。
- `Q-310-r1-7`＝**存量 55 枚的处置⛔ 本票射程**：本腿起手 55、交件 58（差三枚全是本腿的，`20-paired-readings.md` §4 逐名），⛔ 删过任何一枚；票面那条"先逐枚名册、⛔ 按名册逐枚、先经非实现者复核"仍⛔ 满足，且名册尺要先按 `A857` §5#2 那把改正再重跑。
- `Q-310-r1-8`＝**枚数会随任何一腿再长**：丙那一行买的是"看见它在长"，⛔ 是"有人收尾"；本票的甲-G 只改"落在哪儿"。乙（起手清自己上一发留下的）按票面 §4 是**甲-G 落地后可再评的增量**，本腿⛔ 自行采用。

## 2. 纪律自证（本腿做了什么、没做什么）

只 commit、⛔ push；每笔 `add`／`commit` 双显式 pathspec（⛔ `-A`／`.`／目录一把扫）；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内⛔ worktree/checkout；`probes/**` 下⛔ 新建 `.go`（新增件全部 `.md`／`.txt`／`.ps1`）；临时件只建不删、只建在 `probes/310/r1/**` 与仓外 `/d/tmp/310r1-*`；证据件⛔ 叫 `.out`；⛔ 整包 `go test ./...`（只跑窄用例与三枚入口档）；机主/别家那批未提交改动（`.gitignore`／`go.mod`／`design/**`）⛔ 提交、⛔ revert、⛔ 修改；起手暂存区 0 枚、全程未动过别家路径。票尾 Progress bullet＝一条 append-only（1 加／0 删），⛔ 动框、⛔ 动 `Status:`、⛔ 改名 `-done`。
