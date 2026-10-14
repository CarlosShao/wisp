# 174-r5 证据：票 174 AC#6（甲形）— ubuntu 腿那枚未记账 SKIP 登进名册

- 锚（起手）：HEAD `351214de` / 2026-10-14 10:57:55 +08 / `git status --porcelain -- scripts/ .github/` = **0 行**（两处干净）。
  写作时刻 HEAD 已被别腿推到 `65ff087d`（313-r3 两笔；与本腿零文件重叠）。
- 身份：实现者（174-r5）。⛔ 未翻任何框、⛔ 未动 `Status:`/票名、⛔ 未碰 `.github/**`、⛔ 未动任何测试文件/产码。
- 写面：`scripts/portable-tests.sh`（+1/-0）＋本证件目录。

## 1. 结论（先行）

- 现量（起手）：`grep -c 'TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3' scripts/portable-tests.sh` = **0**；`ledger=(`（`:659`）共 **11** 行。
- 落地：名册追加 **1 行**（`git diff --numstat` ＝ `1 0 scripts/portable-tests.sh`，逐字在 logs/01）。加的行逐字：

```text
    "TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3|./internal/tools/|linux|fixture|ticket174 AC#6 shape 甲 (the ledger row, not the POSIX-equivalent shape 乙). The case is AC#2d's junction arm: the filed path must be a REAL NTFS junction (cmd /c mklink /J) whose link spelling C26 refuses - per the fixture comment (task_output_pointer_notice_test.go:376-380) a non-exempt reparse traversal is the ONLY reachable input that gives this branch a path literal at all. On POSIX nothing can be refused, because the detection it exercises is a Windows implementation (risk/pathresolver_other.go reparseComponents returns nil elsewhere - the same wall the workspace-junction row above names), so this machine has no subject: task_output_pointer_notice_test.go:408 skips with 前置条件缺失：造不出真 NTFS junction（mklink /J %s %s: %v: %s）, and the %v the ubuntu leg filled in was exec: \"cmd\": executable file not found in \$PATH: (archived run-37166458550-failed.log, the file then at :428 - the line moves as other criteria land, so the row quotes the skip rather than the line). It RUNS AND PASSES on windows (real mklink /J, no skip), so the assertion is not lost, it is only un-evaluable here. Registered rather than deleted because runtests.sh makes every SKIP fatal (ticket 71 AC#3) and one unaccounted skip is exactly what turned test-core step7 red in the archived run (four numbers there: PASS=1010 FAIL=4 SKIP=1 - the FAILs have their own owners, the SKIP=1 is what this row accounts for): -skip filters an accounted name silently, so this is the repo's chosen mechanism for declaring the skip and re-checking it every run instead of letting it land as an unnamed red. OWNER of the POSIX half: whoever lands AC#6 shape 乙 on ticket 174, a POSIX-runnable equivalent of the same criterion (the notice may not relay the filed path, no junction needed); this row does not fake that half. Re-verified against the compiled LINUX test binary every run: delete the skip and the case starts running here; delete the case and the staleness check goes red. next= on ticket 174"
```

- 机制（为什么这行能消掉那枚红）：`-skip` 对**已入账**的名字是静默过滤（不产 `--- SKIP` 行）。这不是本腿推的，是归档日志自己的读数：`run-37166458550-failed.log` 里 core 腿 8 枚入账名（`-skip ^(8 名)$`）生效、`--- SKIP` 恰好只有这枚未入账用例那 1 枚（同档 test-windows 的 7 名同理，SKIP=1 是另一枚自己 `t.Skip` 的用例）。
- 平台影响：linux 腿 `-skip` 集合 8→9（唯一新增＝本行）；windows 腿集合逐名相同（7=7）。零测试文件改动、零机制改动。

## 2. 要登记的那条跳过（原文与现量行号）

- 用例：`TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3`，`internal/tools/task_output_pointer_notice_test.go:420`（AC#2d 的泄漏半格；`174-r3` 追加）。
- 跳过点：现量 **`:408`**（helper `mkJunctionForLeakProbe174r3`，`:404` 起，逐字）：
  `out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput()` →
  `if err != nil { t.Skipf("前置条件缺失：造不出真 NTFS junction（mklink /J %s %s: %v: %s）", link, target, err, out) }`
- ubuntu 腿渲染原文（归档 `run-37166458550-failed.log`，当时该文件在 `:428`；两处读数：用例自己那行、`unaccounted SKIP lines` 段的重印）：
  `task_output_pointer_notice_test.go:428: 前置条件缺失：造不出真 NTFS junction（mklink /J /tmp/TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r32408037487/001/ALLOWLIST-LEAK-PROBE-174r3 /tmp/TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r32408037487/003/RESOLVED-LEAK-PROBE-174r3: exec: "cmd": executable file not found in $PATH: ）`
  （逐字件：logs/05；三句步骤红文同在 logs/05。）
- 为什么这台机器没有可测对象：该支要求 C26（`risk.Resolve`）**拒绝**一枚非豁免 reparse 遍历，而文件自己注明这是该支"唯一可达的非空输入"（`:376-380`）。POSIX 上 `reparseComponents` 是 stub（`internal/risk/pathresolver_other.go:12`，`return nil`）⇒ 什么都不拒 ⇒ 该支不可达，没有对象可测。windows 上真建 junction、不 skip（a4 普查逐字：``:404`` 真 `mklink /J` 不 skip），所以断言不丢、只是在这边不可评。

## 3. 路由先例核（派单第 3 条；结论：现行裁法＝ledger 登记，无矛盾）

- 票 174 AC#6 本体（票面 `:90-96`，⛔ 本腿未动它的框）：甲＝"在同一张 `ledger=(` 里补一行，五字段齐全…并**写明 POSIX 半边归谁**（照票 111 AC#10 那一行的既有写法：`OWNER of the POSIX half: whoever lands ...`）"。
- 票 111 AC#10 现文（`:235-240` 附近）逐字："⇒ 走 AC#10 给的第一条路：**进 ledger，带平台=linux、类别=fixture、理由、归谁**"；其翻勾记录逐字："`scripts/portable-tests.sh:589` 的 `ledger=(` 名册里 `:598` 逐字登记了那枚 skip"（行号随改动漂移，名册行现在是 `:668`）。
- `scripts/portable-tests.sh` 自己的 GUARD 段（`:44-66` 注释、`:655-671` 名册、`:696-745` 校验）把"未记账 SKIP 即红 / 条目腐坏即红"写成机制本身；AC#6 的边界行逐字："票 111 AC#10 管的是**名册机制本身**…本格管的是**本票用例该不该进那本名册**。⛔ 两格不许合并、也不许由本腿去改机制。"
- ⇒ **两处口径一致**（ledger 登记就是现行裁法）⇒ 照 AC#6 甲形做；机制零改动、`.github/**` 零改动。

## 4. 字段为什么这么填

| 字段 | 值 | 依据 |
|---|---|---|
| name | `TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3` | 现量 `:420` 的 TEST 名；名册校验要求裸 Test 标识（`:703`） |
| package | `./internal/tools/` | 同包两行既有写法（`:668`/`:669`）；core 腿 scope 含 `./internal/tools/...`（`:306` 附近），归档日志的 scope 串逐字含之 |
| platform | `linux` | 它跳的平台＝ubuntu 腿（`goos=linux`）；windows 上不跳 ⇒ 若填 `any` 会在 windows 腿把这枚真跑的用例静默 skip 掉 |
| class | `fixture` | 名册自释义（`:656`）"the machine or OS cannot supply the subject at all"；非 opt-in（无 env 门）、非 reexec（非子进程入口） |
| reason | 见 §1 加的行 | 含现量 file:line（`:408`）＋跳过原文摘句（格式串逐字＋ubuntu 渲染）＋为什么本机无对象（POSIX 无 reparse 对象）＋`OWNER of the POSIX half`＋`next= on ticket 174`（照 `:668` 文风，一句到底） |

⚠ reason 里 `\"cmd\"` 与 `\$PATH` 是在双引号 shell 串里的转义写法，运行时得到的就是 `"cmd"` 与 `$PATH` 两个逐字字节；reason 内不含裸 `|`（本行按 `|` 分正好 5 字段，logs/02）。

## 5. 门禁与读数（logs/ 逐件；冗余原始件见 §6 表）

1. `bash -n scripts/portable-tests.sh`：**pre rc=0 / post rc=0**（logs/00、logs/01）。
2. 两把 grep 尺：`grep -c 'TestCanonicalize…'` **0 → 1**；名册枚数 **11 → 12**（logs/00、01、02）。
3. 脚本**自己的**校验代码重放（把 `ledger=(...)` 块原样 `eval`，用脚本逐字的 malformed／裸名／平台分支表达式跑；logs/02）：`entries=12`、本行按 `|` 分 **5** 字段、`malformed=0`；linux 活跃 **9** 枚（新名在模式里逐字可见）、windows 活跃 **7** 枚。
4. 平台集合 old-vs-new（`git show HEAD:` 对照，logs/02 尾）：windows `HEAD == worktree`（逐名相同）；linux 唯一新增本名（8→9）。
5. 腐坏尺代理（脚本 `:721` 用 `go test -list` 重验条目对应的测试名）：`GOTMPDIR=$LOCALAPPDATA/wisp-test-gotmp/174-r5-probe`（C: 盘）下 `go test -list '^TestCanonicalize…$' ./internal/tools/` **rc=0** 且逐字列出该名（logs/04）；`grep -c '^//go:build'`（该文件）= **0** ⇒ 同一枚也会编进 linux 测试二进制。
6. 附加自扫：新行在 d22scan 表情带（1F000-1FAFF/2200-22FF/2600-27BF/2B00-2BFF/FE0F/1F1E6-1F1FF）**0 命中**（logs/01 尾）。
7. ⛔ **没做，具名为什么**：
   - **linux 侧 `go test -list` 腐坏校验 ＋ `test-core` 那一步的真颜色**：win32 上无法 exec linux 测试二进制，仓外的 GitHub CI 是唯一能给出这枚读数的地方 ⇒ 只能等**下一发被推送的 CI**（ubuntu 腿、`scripts/portable-tests.sh` 默认 core scope）。预期读数：该用例不再出现在 `=== RUN` 里、由它产的那枚 `--- SKIP` 消失、`unaccounted SKIP lines` 段与它致的那份 `strict runner exited 1` 消失（同一 run 其余红各自有主，本行不认领）。
   - **全量 `portable-tests.sh` 干跑**：脚本无 dry-run/自检开关（Usage 只有 scope 选择，`--scope=census` 是 tier 审计、在名册段之前就退出）；派单门禁清单只要求 `bash -n`＋两尺。windows 腿与本行无关这点已由 §5.4 的集合层读数给出，不靠"跑一遍猜"。
   - `scripts/d22scan.sh` 与长度门：派单明示豁免（与 shell 名册无关），未跑。

## 6. 交件名册（逐件；行数＝入库跑现量）

| 文件 | 行数 | 提交 |
|---|---|---|
| `.scratch/wisp/probes/174/r5/00-anchor-and-gates.md`（本件） | 75 | 是 |
| `logs/00-pre.txt`（锚＋干净度＋两尺 pre＋`bash -n` pre） | 18 | 是 |
| `logs/01-post-anchor.txt`（两尺 post＋`bash -n` post＋`+1/-0`＋表情带自扫） | 23 | 是 |
| `logs/02-ledger-replay.txt`（脚本自身校验重放＋两平台模式＋old-vs-new 集合 diff） | 18 | 是 |
| `logs/04-staleness-proxy.txt`（`go test -list` 腐坏尺代理＋无 build tag 尺） | 12 | 是 |
| `logs/05-skip-verbatim.txt`（`:404-409` 源码＋归档日志渲染原文＋三句步骤红文） | 17 | 是 |
| `logs/03-platform-pattern-diff.txt`（同一比较的原始运行，内容已并入 02 尾） | 9 | **留盘未提交**（按 `issues/README` 规则 8"只建不删"，清单在此，交编排者统一处置） |

复现（任一枚）：`bash -n scripts/portable-tests.sh`；`grep -c 'TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3' scripts/portable-tests.sh`；`sed -n '/^ledger=(/,/^)/p' scripts/portable-tests.sh | grep -c '^    "'`；`git diff --numstat -- scripts/portable-tests.sh`。
