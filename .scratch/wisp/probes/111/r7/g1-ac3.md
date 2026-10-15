# 111-r7 / g1 — 票 111 AC#3（空分母守卫）逐格现读 + 当前 HEAD 两发人为变异

**腿**：`111-r7`（编排者 2026-10-15 派；派单存档 `.scratch/wisp/probes/orch/20261015/txt/111r7-dispatch.txt`）。
**起手锚**：`457e7513`（11:48 +0800 现读）。**漂移逐报**：11:50 现读 HEAD = `2ad0aa86`
（`319-r2` 落地一笔 + 编排者两笔）；`git diff --stat 457e7513..2ad0aa86 -- scripts/ .github/` = **空**
⇒ 本格靶件字节逐枚同锚点，两发读数对两枚 HEAD 同一。13:xx 收笔前复锚见 `verdict.md`。
**本格判（第 0 步现读）**：**半做**——守卫本体与"正控在 CI 真跑过"两半**已在盘上**（下面 §1 逐字钉），
缺的是"**人为抽掉一个包 ⇒ 红**"在**当前 HEAD** 的新鲜读数（旧读数都停在 10-05/09-21 世代）⇒ 本腿补 §2/§3 两发；
⛔ 未翻任何框、⛔ 未改 `scripts/**` 一字节。

## §1 现读（逐字钉住"已做的那半"）

- **守卫 A（AC#3 逐字要的那条）** 在盘上：`scripts/portable-tests.sh:651`（探针
  `denom=$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{else}}EMPTY {{.ImportPath}}`）→`:655` 报红 → `:664 exit 1`；
  逐字块 `logs/g1-runs.txt` 末段＋`ac3-guardA-block.txt` 同文；**守卫 C** `:629`（pin vs resolved 双向）、
  **守卫 D** `:474`（`0/0) ;;` 那一支）＋`:487`（totals 行）＋`:490`（报红）同在。
- **CI 侧读数（run `38676514264` / head `39b0a246`，2026-10-15 03:29Z 起）**，全在 `logs/g1-ci.txt`：
  - `test-windows` job `116231259514` **step 9** `Package coverage census (ticket 111 AC#1 + GUARD D)` = **success**；
    步内逐字 `census totals: packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`。
  - `lint` job `116231259522` **step 14** `Portable tests carrier self-test (ticket 111 AC#3 - GUARD D's positive control)` = **success**；
    步内逐字 `32 case(s) ran, 0 assertion(s) failed` + `GREEN - every seeded anomaly was refused, and the clean scope passed.`；
    正控逐条真跑过：case 8 `guard-a-empty-test-package: rc=1`、case 23 `census-unclaimed-package-goes-red: rc=1`、
    case 24（mutant 绿 / 真件红）。
  ⇒ 验收件（编排者 10-06 12:5x 那条）里缺的"**载体在 CI 上从未被跑过**"这一行读数，**今天已存在**。
- **空分母现状**（本格主题）：census 基线（快照内跑，rc=0，`logs/g1-census-base.txt`）逐字两行——
  `github.com/CarlosShao/wisp/internal/session 3/0 corewindows`（r1 腿放回，已有分母）、
  `github.com/CarlosShao/wisp/internal/watchdog 0/0 NO-SCOPE <-NO-TESTS`（仍然空、仍**不在**任何 scope）。
- **旧世代读数只具名、不重做**：`probes/111/r2/census-mutant-D1-plugin-pulled.txt`＋`D3-session-pulled.txt`＋
  `mutants/*.diff`（GUARD D 两发，已入库）；票面 Progress log 里旧腿的三发（GUARD C 抽 plugin、加 watchdog 两形）。
  本腿不引它们当本格凭据，只说明"这一发不是新发明"。

## §2 M1 —— 人为抽掉一个包：core scope 删 `./internal/plugin/`（pin 不动）⇒ 期望红

- **做法**：`git archive HEAD | tar -x -C D:/tmp/wisp111r7/snap-ac3`（仓外快照，⛔ 仓内零 worktree/checkout）；
  `cp -r snap-ac3 snap-ac3-m1` 后在 m1 里 `sed`；同一条 `&&` 链 `grep -n` 打印改后整行：
  `324:        ./cmd/llmrecord/ "$winsec_scope"`（原 `./internal/plugin/ ` 前缀没了）。`bash -n`（m1 内）**rc=0**。
- **mutant diff**：`logs/g1-mutants.diff` 第一段，恰好 **1 行**。
- **读数**（m1 内 `bash scripts/portable-tests.sh --scope=core`）：

```
portable-tests.sh: GUARD C - scope mode=core resolved to a DIFFERENT package
portable-tests.sh:   set than the one pinned next to it. Pinned: 27, resolved: 26.
portable-tests.sh:   19d18
portable-tests.sh:   < github.com/CarlosShao/wisp/internal/plugin
```

  **rc=1**，7 行全文在 `logs/g1-runs.txt` 第一段。
- **判**：抽掉一个包 ⇒ 这一步**不**会"少跑一个包照样绿"（点名被抽的那枚 + 分母自己动）。
  这就是票面 AC#3 第二半"并人为抽掉一个包证明它会红"要的那一发，当前 HEAD 成立。

## §3 M2 —— 空分母声明：core scope+pin 加回 `./internal/watchdog/` ⇒ 期望红（守卫 A 本体）

- **做法**：`snap-ac3-m2`（`cp -r` 自 snap-ac3）里 2 处 `sed`：pin 插 1 行、core scope 行尾追加 ` ./internal/watchdog/`。
  ⚠ 诚实边界：第一次切时那个 `sed` 把 **windows scope 的同形行**一起咬了（collateral）⇒ 从仓内
  `git show HEAD:scripts/portable-tests.sh` 还原重切；take-1 原文留痕 `logs/g1-runs.txt` 第三段；
  最终 diff = **恰好 2 行**（`logs/g1-mutants.diff` 第二段）。`bash -n` rc=0。
- **读数**（m2 内 `bash scripts/portable-tests.sh --scope=core`）：

```
portable-tests.sh: GUARD A - these packages are declared in scope mode=core but
portable-tests.sh:   compile NO test file at all for GOOS=windows:
portable-tests.sh:   EMPTY github.com/CarlosShao/wisp/internal/watchdog
portable-tests.sh: a scope entry with no denominator cannot fail, so keeping it there
portable-tests.sh: is a false claim of coverage (ticket 111 AC#3). ...
```

  **rc=1**（GUARD C 先放行：pinned 28 = resolved 28，红的只可能是 A）。
- **判**：AC#3 逐字那句"声明在范围内但该平台**没有任何测试文件** ⇒ 直接失败"有牙、报红句点名包与平台。
  M1 与 M2 合起来才把 AC#3 两个方向都钉住：**抽走（不许少跑）＋声明空壳（不许空跑）**。

## §4 本格判语（供非实现者复核）

- 守卫=已做（§1）；正控在 CI 真跑=已做（§1 的 step 14 读数，旧验收件缺的就是这行）；当前 HEAD 的
  "人为抽掉一个包⇒红"=本腿补齐（§2/§3，两发都 rc=1 且点名）。⇒ 建议可翻；⛔ 本腿不翻框、不碰票名/台账。
- 未做/具名边界：本格未在任何 CI 步里制造红（正控是**载体内部断言**的 rc=1，不是把 CI 弄红），
  这是载体设计的形状（`32 case(s) ran, 0 assertion(s) failed` 即"每一发的红/绿都按设计发生"）。
