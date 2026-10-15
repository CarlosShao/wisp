# 111-r7 / g3 — 票 111 AC#7 逐格现读 ＋ 二选一落定（选 ②：当众改口径）

**本格判（第 0 步现读）**：**没做**（① 从未落地：`-run` 计数 0；② 的改口径文字只躺在票面 Progress log 里，
框未翻、且"当众写死"没有一处可引的硬读数）⇒ 本腿取 ② 并把它写死在本件＋票面尾 bullet；
⛔ 未改 `ci.yml`／`scripts/**` 一字节、⛔ 未动 `-skip` 与 ledger（那是刻意的机构，不是缺陷）。

## §1 现读（"① 落地"那一支今天仍 0 次）

- 尺 1：`git grep -c -- "-run TestSyncRegistryProbeLive" HEAD -- scripts tools .github` ⇒ **0 命中**（rc=1，无输出）。
  ⇒ 全 CI 面（含 `ci.yml`）没有任何一步用 `-run` 点名跑它。件 `logs/g3-rulers.txt`。
- 尺 2：run `38676514264`（head `39b0a246`，已推送）test-windows job `116231259514` **整发日志**里
  `grep -cE '(=== RUN|--- SKIP|--- PASS|--- FAIL).*TestSyncRegistryProbeLive'` = **0**
  ⇒ 它在 CI 里**从未被执行**（日志里出现的地方只有 ledger 行与 `-skip` 串）。同名日志 `logs/g3-rulers.txt`。
- 尺 3（押在盘上的机构，逐字）：`scripts/portable-tests.sh:677` 的 ledger 行——
  `"TestSyncRegistryProbeLive|./internal/risk/|windows|fixture|"`，理由串含
  `Measured on a windows host at HEAD 84e43af: --- SKIP at syncdirs_windows_test.go:133` 与
  `Remedy for a real run: a host with a sync record, or fold the shape check into the fixture-driven case.`
- 尺 4（被测件自己的行为，逐字）：`internal/risk/syncdirs_windows_test.go:132-134`
  `if len(roots) == 0 { t.Skip("no registry-grade sync record on this machine ...") }`
  ⇒ 没有同步记录的主机上它是 `t.Skip`，不是 FAIL 也不是 PASS。

## §2 ① 为什么被拒（两读数都写在这里，不假称"环境问题"）

- 把 `-run TestSyncRegistryProbeLive` 接进 **hosted** windows 腿只有两种可能读数，且都不是覆盖：
  1. 走 `runtests.sh`（本仓唯一合法运行器）⇒ `--- SKIP` 致命（`tools/d22scan/runtests.sh:24-26` 与 `:98`→`:102`
     "SKIP is NOT a pass"）⇒ **常红**；
  2. 走裸 `go test` ⇒ skip 印成 `ok` ⇒ **假绿**（正是票 93 建票要消灭的字节形状）。
- 这条论据**不是我新造的**：票 110 面 `110-no-ci-step-runs-internal-winsec-done.md:96` 逐字
  「**AC#4 `R-93-4` 的答复 = 它不该进 RUN 分母，但它该拿的步级证据我给了**」＋本机实测 `--- SKIP`；
  对抗验收表 `docs/evidence/s1/110-adversarial-acceptance.md:17` 已裁「"不该进 RUN 分母"的论证我认可」。
- 所以 ① 若照字面落地=把一台**没有主体（无同步记录）**的主机上的 SKIP 塞进分母，
  与 AC#7 想买的"真结论"正相反。**决定：不落地 ①**（并说明这不是"留在原地当已做"——见 §3）。

## §3 ② 当众改口径（落定文本，本件即公开面之一）

自本件/本 bullet 起，本票与票 110 相关的这条口径**改成**下面三句读（旧措辞＝票 93 `:201-202` 的
next= 第 4 条、票 111 AC#7 `:75-76` 的转述）：

1. **撤回**"`-run TestSyncRegistryProbeLive` 经 `runtests.sh` 接进 windows 腿"这一形：在任何没有 HKCU
   同步记录的主机上它只会是常红（SKIP 致命）或假绿（裸 `go test`），两读数都不是参评（§2）。
2. 票 110 AC#4 / `R-110-4` 真正欠的那发步级读数 = **windows 腿对 ledger 行的逐 run 复验＋class/reason 打印**。
   这发读数**今天存在**：run `38676514264` / job `116231259514` / 步名（ci.yml:763 逐字）
   `Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord)`
   （该步 03:40:31Z 起），步内 `03:40:45.3119022Z` 逐字打印该行全文（含 class `[fixture]` 与理由/R-93-4 那句），
   紧随 `-skip pattern built from the ledger: ^(...|TestSyncRegistryProbeLive)$` 点名它；
   件 `logs/g3-ci-ledger-print.txt`（去时间戳后的同文）。
3. **残余不伪造**："它真跑过并通过"至今 **0** 次（尺 2）。真要取那一读数，需一台 HKCU 有同步记录的机器
   （票 110 `:96` 具名的唯一候选＝self-hosted `wisp-slo`）；那是**另一个 runner 面/另一张票**的事，
   ⛔ 不是 hosted windows 腿，也⛔ 不许拿"接进 hosted 腿"充数。ledger 行里的
   `Remedy for a real run: a host with a sync record, or fold the shape check into the fixture-driven case`
   就是这条残余的下一手（且第二支要动 Go 源码＝本票禁面）。

## §4 本格判语

- AC#7 = **本腿按 ② 落定**：① 拒（逐条读数在 §2）、② 的文本写死在本件 §3＋票面尾 bullet（逐字可引）。
  这一格不再有"欠编排者 push 后读 step7"的挂账——② 引用的 run `38676514264` 已推送、读数已取；
  ① 若日后再被谁提出，它必须先回答 §2 的两读数（⛔ 不许默默落地）。
- ⛔ 本腿不翻框；`-skip`/ledger/被测件/`ci.yml` 均零字节改动。
