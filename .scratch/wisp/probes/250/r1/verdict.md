# 票 250 · 落地腿 `250-r1` 裁决件 —— `portable-tests.sh` GUARD C 的分母被 `go list` 的 stderr 污染

## ① 起手锚与写面

- 起手锚（`git log -1 --format='%H %ad' --date=iso`，取数时刻 2026-10-02 08:49:08 +0800）：
  `50d7a88172f397befcac49e08778bbc6f6355ead 2026-10-02 08:47:53 +0800`，分支 `dev`。
- 起手写面：`git status --porcelain scripts/` ＝**空**（08:49 现跑，与票面 08:5x 的断言一致）。
- 起手全树：`git status --porcelain | wc -l` ＝ **304**（取数时刻 08:51；共享树，此数会漂，任何"此刻多少枚"只在此刻有效）。
  其中可见在飞腿的痕迹（`M internal/panel/bridge.go`、`M internal/panel/l2_grant_boundary_test.go`、`M cmd/wisp/run.go` 等），
  本腿一律不碰、不提。
- 票面：`.scratch/wisp/issues/250-portable-tests-guard-c-merges-go-list-stderr-into-its-own-denominator-so-a-cold-module-cache-kills-the-whole-core-scope-reading.md`
  （46 行 / 7,679 字节，`wc -l -c` 现跑）。AC 勾选框一枚未动。

### ①.1 编排者题面那几行尺的复认（全部**成立**，逐字核对 08:51 现跑 `grep -n "" scripts/portable-tests.sh | sed -n '247,296p'`）

| 题面断言 | 实测 | 结论 |
|---|---|---|
| `:261` 逐字 `if ! go list "${scope[@]}" >"$resolved" 2>&1; then` | 逐字一致 | 成立 |
| `:268` `grep -v '^$' "$resolved" \| sort -u >"$resolved.sorted"` | 实测行尾还有 ` \|\| true`（题面省略了这一段，语义无差） | 成立（措辞不全） |
| `:270` `pkgcount=$(wc -l <"$resolved" \| tr -d '[:space:]')` | 逐字一致 | 成立 |
| GUARD C 本体 `:278-292`，`want` `:279` / `got` `:280` / 比较 `:281` / 红句 `:283-284` / `exit 1` `:291` | 逐一对位 | 成立 |
| `:254-256` 空 scope 硬退出 `exit 2` | `:254` if / `:255` echo / `:256` `exit 2` | 成立 |
| `:262-266` `go list` 真失败支（`cat` 在 `:262`、`exit 1` 在 `:266`） | 成立 | 成立 |
| 本机 `GOMODCACHE` 温热 ⇒ stderr 为空 ⇒ 缺陷本机静默 | 现跑真实 `go list <core 的 20 枚 glob>`（温热缓存）：`rc=0`、stdout **25 行**、stderr **0 字节** | 成立 |

一处**推翻/修正**：题面（与 `ci-delta-1`）把污染说成"那 10 行"。本机真实冷缓存复现拿到的是 **11 行**（多 `github.com/mattn/go-isatty`、`github.com/ncruces/go-strftime`，少 `github.com/google/uuid`），
于是红句是 `Pinned: 25, resolved: 36.` 而非 CI 的 `35`。机制同形，枚数随主机的模块集/GOOS 变——见 §②。

## ② 改前那发红句逐字（真实冷缓存 + 真实脚本，本机可重跑）

载具（最省形＝票面 AC#3 点名的那一形：`GOMODCACHE` 指到空目录再跑 `--scope=core`）：

```bash
cd "D:/work/workspace/projects plans/Wisp"
cold=$(mktemp -d)
GOMODCACHE="$cold" GOPATH="$cold/gopath" bash scripts/portable-tests.sh --scope=core
# rc=1，12.2 秒内死在 GUARD C；$cold 落到 /tmp/tmp.2ztPH1RoMJ（303M，按 issues/README 规则 8 只建不删）
```

读数逐字（存档件 `.scratch/wisp/probes/250/r1/logs/prefix-real-coldcache.txt`）：

```
portable-tests.sh: GUARD C - scope mode=core resolved to a DIFFERENT package
portable-tests.sh:   set than the one pinned next to it. Pinned: 25, resolved: 36.
portable-tests.sh:   25a26,36
portable-tests.sh:   > go: downloading github.com/dustin/go-humanize v1.0.1
portable-tests.sh:   > go: downloading github.com/mattn/go-isatty v0.0.24
portable-tests.sh:   > go: downloading github.com/ncruces/go-strftime v1.0.0
portable-tests.sh:   > go: downloading github.com/pelletier/go-toml/v2 v2.2.4
portable-tests.sh:   > go: downloading github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec
portable-tests.sh:   > go: downloading golang.org/x/crypto v0.57.0
portable-tests.sh:   > go: downloading golang.org/x/sys v0.48.0
portable-tests.sh:   > go: downloading modernc.org/libc v1.75.7
portable-tests.sh:   > go: downloading modernc.org/mathutil v1.7.1
portable-tests.sh:   > go: downloading modernc.org/memory v1.12.1
portable-tests.sh:   > go: downloading modernc.org/sqlite v1.59.0
portable-tests.sh: a deleted or newly-uncovered package must fail the step, not
portable-tests.sh: shorten it. Either restore the scope entry or update the pin
portable-tests.sh: in the SAME commit and say why in the CI log.
```

要点：`25a26,36` 即 diff 说"resolved 比 want 多出 11 行"，那 11 行逐字是 `go: downloading …` 进度；
`go test` 一行没跑（脚本在 `:278-292` 就 `exit 1`，GUARD A／B、台账、strict runner 全在其后）⇒
**整段 core scope 的 `=== RUN` 读数在这一步永久采不到**，票面 §代价 与本机复现同形。

## ③ 改法与逐枚 `file:line`

射程两枚文件：`scripts/portable-tests.sh`（被修物）＋ 常驻载具（AC#3 要求的"本机就能判"的那把尺）。
改法三块，逐枚行号在改后复量并记于此节末尾：

1. **AC#1 分流**：`go list` 的 stdout 与 stderr 各落一枚临时文件；成功时 stderr 若有内容 ⇒ **原样回声进步骤日志**
   （"不许进分母"不等于"不许看"），失败时**两路都打**并仍 `exit 1`。
2. **AC#2 能力形**：分母（已 `sort -u` 的那份文件，与 `pkgcount` 数的是同一份字节）逐行过 `[[ $line =~ $import_path_re ]]`，
   判据＝整行无空白 + 每段只含模块路径文法允许的字符；任一行不过 ⇒ 指名打印并 `exit 1`，
   ⛔ 不 grep 掉继续。
3. **常驻载具**：`scripts/portable-tests-selftest.sh` + `scripts/testdata/portable-tests/go`（一层假 `go`，
   只替换 `go` 这一个进程，跑的是**真脚本**），八形模式见 §④。假 `go` 的核心 import path 集合是从
   `scripts/portable-tests.sh` 的 `core_pin` 现场 `sed` 抽的，不另存第二份真相。

## ④ 改后读数 + 正控读数（四格各自对应）

四格与对应载具模式（读数在 §④.1 起逐条展开；每条都是本机实跑，非推演）：

- AC#1 分流 ⇒ `coldcache-stderr` 模式：stderr 有进度、stdout 完全正确 ⇒ `pkgcount` 必须仍是钉住的那枚数、守卫不响；
  反面 ⇒ `golist-fails` 模式：`go list` 真失败时仍非 0 退出且**两路都打**。
- AC#2 能力形 ⇒ `seed-stdout` 模式：往 stdout 种一行带空格的伪包名 ⇒ 必须响亮拒绝并非 0 退出。
- AC#3 常驻 ⇒ 上面两形全在本机跑，改前/改后各一遍，读数落 `.scratch/wisp/probes/250/r1/logs/`。
- AC#4 不许改软 ⇒ `guard-a`／`guard-b`／`empty-scope` 三形各一枚"种 X 必响"，另加三段守卫的**字节不变**证明。

## ⑤ 越界检查

- 本腿提交只带显式 pathspec；`git add -A`/`.`/`commit -a` 一次未用。
- 提交前 `git status --porcelain` 与起手 304 枚的差集＝只本腿那几枚路径（差集读数记在 §⑤.1）。
- `frontend/**`／`design/**`：零读取、零引用、零写入（起手 `git status` 顺带列出的 `design/...` 行是别的腿的痕迹，本腿未打开任何一枚）。
- 未跑整树 `go test ./...`；未跑 `cmd/wisp`／`internal/panel`／`internal/ball` 任何一包（载具里的 `go` 是假进程，真 `go test` 一行未执行）。
- 台账 `docs/reports/pending-and-issues.md` 未写；`thresholds.go`／golden／`allowlist.txt`／三枚冻结测试件／`PLAN.md`／`docs/specs/**` 未碰。

## ⑥ 门禁读数（这一步在 CI 那台上会是什么颜色）

本机载具读数只证明**形状**：改前红、改后绿都在 §②／§④ 落地，且用的就是 CI 那一步的同一段代码。
**CI 的颜色本腿不判**，原因有三，逐条具名：

1. 本腿不推送、不 `gh` 触发任何东西（硬边界 8）；CI 上的 `test-core` 第 7 步跑的是 `origin/dev` 的 SHA，本机 HEAD 不是它。
2. CI 那台 runner 的冷缓存枚数随实例而异（本机 11 行、CI 那发 10 行）——修好后**无论几行都不该进分母**，
   但"这一发真的采到了 1365 枚 `=== RUN`"只有编排者推上去看日志才能钉。
3. 台账 `A512` 归编排者写。

本腿能给的承诺读数：`scripts/portable-tests-selftest.sh` 全模式在改后应 `rc=0` 且逐模式打 `VERDICT=ok`，
改前应至少两形红（`coldcache-stderr`、`seed-stdout`）。实跑读数见 §④。

## ⑦ 判不动的地方

1. **CI 上 `test-core` 改后的真实颜色与 `=== RUN` 枚数**——判不动的原因：本腿禁推送禁 `gh`，且冷缓存只在 runner 实例上随机出现；
   能判的人＝编排者（推上去看一次 run，或对比同 SHA 的 schedule run `36940536372`）。
2. **那 4 枚"消失"的前端契约红（`TestC21DesignTokensFourwayAgree` 等）到底是变绿还是被跳过**——判不动的原因：
   判据要在真 `go test` 跑起来之后才存在，本票修好之前 CI 上根本没有那发读数；
   而且其中三枚落在 `internal/panel`（在飞腿 `248-r1` 的写面 + 冻结测试件）。能判的人＝编排者 + `248-r1`。
3. **假 `go` 载具能证"形状"，证不了"真 `go list` 的 stdout 永远不含进度"**——载具的 stdout 是我写的；
   真实冷缓存那发（§②）只证明 stderr 有内容、stdout 是 25 行干净路径。要把"stdout 永不掺进度"钉成契约，
   能判的人＝编排者（或票 250 的验收腿在 runner 上取一发）。
4. **GUARD C 改后在 windows legs（`--scope=windows`）上的同形性**——本腿载具只造了 core 一发钉；
   windows/cli 两枚钉共用同一段代码，但票面 AC 未点名，扩到三枚钉各一发要不要做＝编排者定。
5. **`empty scope` 那支（`:254-256`）从 CLI 面上今天不可达**——判不动的原因：`core/windows/cli` 三枚钉的 scope 都非空，
   `census` 在 `:213` 就 `exit 0` 走掉了，显式传参支 `:247-253` 只在 `$# -gt 0` 时接管，
   所以没有任何入口能把 `scope` 摆成空数组。AC#4 要的那枚"种 X 必响"只能用**逐字切片载具**量
   （`sed -n '247,257p'` 抽出真字节再喂 `scope=()`），读数见 §④.6；
   要不要给脚本补一个可达入口＝编排者的活。
