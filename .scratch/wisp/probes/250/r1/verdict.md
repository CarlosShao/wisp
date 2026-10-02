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

### ①.1 编排者题面那几行尺的复认（**全部成立**，逐字核对 08:51 现跑 `grep -n "" scripts/portable-tests.sh | sed -n '247,296p'`）

| 题面断言 | 实测 | 结论 |
|---|---|---|
| `:261` 逐字 `if ! go list "${scope[@]}" >"$resolved" 2>&1; then` | 逐字一致 | 成立 |
| `:268` `grep -v '^$' "$resolved" \| sort -u >"$resolved.sorted"` | 实测行尾还有 ` \|\| true`（题面省略，语义无差） | 成立（措辞不全） |
| `:270` `pkgcount=$(wc -l <"$resolved" \| tr -d '[:space:]')` | 逐字一致 | 成立 |
| GUARD C 本体 `:278-292`，`want` `:279` / `got` `:280` / 比较 `:281` / 红句 `:283-284` / `exit 1` `:291` | 逐一对位 | 成立 |
| `:254-256` 空 scope 硬退出 `exit 2` | `:254` if / `:255` echo / `:256` `exit 2` | 成立 |
| `:262-266` `go list` 真失败支（`cat` 在 `:262`、`exit 1` 在 `:266`） | 成立 | 成立 |
| 本机 `GOMODCACHE` 温热 ⇒ stderr 为空 ⇒ 缺陷本机静默 | 现跑真实 `go list <core 的 20 枚 glob>`（温热缓存）：`rc=0`、stdout **25 行**、stderr **0 字节** | 成立 |

**本腿推翻/修正编排者题面与票面的三处**（其余行号全对）：

1. **"那 10 行"不是常数。** 本机真实冷缓存复现拿到 **11 行**（多 `github.com/mattn/go-isatty v0.0.24`、`github.com/ncruces/go-strftime v1.0.0`，
   少 `github.com/google/uuid v1.6.0`），于是红句逐字是 `Pinned: 25, resolved: 36.` 而非 CI 的 `resolved: 35`。
   机制同形，枚数随主机 GOOS 与冷集而定——**所以任何把"35"写进断言的尺都会在第一台不同主机上假红**；
   本腿的判据是"stderr 行数 ≡ 0 进分母"，不是"分母等于某个常数"。
2. **毒不只落在 GUARD C。** 同一份 `$resolved` 还是 `:477`（旧号）`done <"$resolved"` 那枚 GUARD B 逐包循环的分母
   （现号 `:552`）。票面只算了 GUARD C 一枪。**实测量**（`.scratch/wisp/probes/250/r1/logs/second-victim-guardb.txt`）：
   显式路径形态（`pinned=''` ⇒ GUARD C 不响）下把 CI 那 10 行进度喂进 stderr——改前 `rc=1`，
   GUARD B 打 `GUARD B - 10 of 35 packages in scope mode=core printed` 并把十行 `go: downloading …` 逐行点名成
   "测了零次的包"（10 枚 `<NO TOP-LEVEL RESULT LINE> missing`）；改后同一发 `rc=0`、missing 0、真包 25 行。
   ⇒ 分流这一刀同时救 GUARD B 的账；票面那句"3 秒内死守卫"之下还有第二枚受害者，只是被 GUARD C 挡在了前面。
3. **`--scope=winsec` 不存在。** 票面 §现量 未提，但 `scripts/portable-tests.sh` 有 `winsec_pin`（`:164-166`）却不在 `:170` 的 `case $mode` 里
   （只有 core/windows/cli/census），所以那枚钉今天从 CLI 面上**取不到**；`scripts/winsec-tests.sh:94` 是以显式路径委托的，
   显式路径支 `:247-253` 把 `pinned` 置空 ⇒ **GUARD C 对 winsec 那枚钉今天根本不跑**。不归本票修（本票不许顺带扩面），
   但它是"GUARD C 覆盖面"这一格的既有空洞，具名报给编排者。

## ② 改前那发红句逐字（真实冷缓存 + 真实脚本，本机可重跑）

载具（最省形＝票面 AC#3 点名的那一形：`GOMODCACHE` 指到空目录再跑 `--scope=core`）：

```bash
cd "D:/work/workspace/projects plans/Wisp"
cold=$(mktemp -d)
GOMODCACHE="$cold" GOPATH="$cold/gopath" bash scripts/portable-tests.sh --scope=core
# 实测：rc=1，12.2 秒；冷缓存目录 /tmp/tmp.2ztPH1RoMJ 落了 303M（按 issues/README 规则 8 只建不删）
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

要点：`25a26,36` 即 diff 报"resolved 比 want 多出 11 行"，那 11 行逐字是 `go: downloading …` 进度；
`go test` 一行没跑（脚本在 GUARD C 就 `exit 1`，GUARD A／B、台账、strict runner 全在其后）⇒
**整段 core scope 的 `=== RUN` 读数在这一步永久采不到**，与票面 §代价 同形。

真实 `go list` 的流形状（同一台机器，1.3 秒，非载具）：

```bash
cold=$(mktemp -d); out=$(mktemp); err=$(mktemp)
GOMODCACHE="$cold" GOPATH="$cold/gopath" go list ./internal/config/... >"$out" 2>"$err"
# rc=0；stdout 1 行 = github.com/CarlosShao/wisp/internal/config；stderr 2 行 = 两条 go: downloading
```

⇒ 前提坐实：**成功退出的 `go list` 可以 stdout 干净、stderr 有货**。

## ③ 改法与逐枚 `file:line`（全部为改后实测行号）

三枚文件，本腿独占写面：

| 文件 | 行 | 是什么 |
|---|---|---|
| `scripts/portable-tests.sh` | `:259-287` | 为什么（把上面那发读数写进注释，含 `Pinned: 25, resolved: 36.` 逐字） |
| 同上 | `:288-289` | 两枚临时件：`resolved`（stdout）与 `resolved_err`（stderr），各自 `mktemp` |
| 同上 | **`:290`** | `if ! go list "${scope[@]}" >"$resolved" 2>"$resolved_err"; then` —— **AC#1 分流的那一行**（旧 `:261` 的 `2>&1` 没了） |
| 同上 | `:291-301` | 失败支：`out\|`／`err\|` 两路各自加前缀打全，末句仍是 `the scope cannot be audited`，`rm -f` 两枚临时件，**`exit 1`** |
| 同上 | `:302-307` | 成功支：stderr 若非空 ⇒ 打一行 `go list wrote N line(s) to stderr (progress, not packages, not in the denominator):` 再逐行 `err\|` 回声；随后 `rm -f "$resolved_err"`。⛔ 不是丢弃，是不计数 |
| 同上 | `:308-310` | **一字未动**（`grep -v '^$' \| sort -u` / `mv` / `pkgcount=$(wc -l <"$resolved")`），此刻 `$resolved` 里只有 stdout |
| 同上 | `:312-326` | AC#2 的"为什么是能力形不是词表" |
| 同上 | **`:327`** | `import_path_re='^[A-Za-z0-9][A-Za-z0-9._+~!-]*(/[A-Za-z0-9][A-Za-z0-9._+~!-]*)*$'` |
| 同上 | `:328-334` | `while IFS= read -r line \|\| [ -n "$line" ]` 逐行过形状（`\|\| [ -n ]` 那半句是防"末行无换行 ⇒ 一行逃过检查"）；读的是**归一化之后、与 `pkgcount` 同一份字节** |
| 同上 | `:335-346` | 任一行不过 ⇒ 打印枚数与**逐行原句**、`refused rather than filtered` 两句、`rm -f`、**`exit 1`** |
| 同上 | `:347-368` | GUARD C 钉比较本体（旧 `:272-293`）——**字节未动**，见 §④.5 的 sha1 |
| 同上 | `:370-388` | GUARD A 本体（旧 `:295-313`）——字节未动 |
| 同上 | `:513-568` | GUARD B 本体（旧 `:438-493`）——字节未动 |
| `scripts/portable-tests-selftest.sh`（新增，313 行） | `:1-31` 用途与边界；`:70-76` 从被审脚本现场 `awk` 抽 `core_pin`（不留第二份真相）；`:110-127` `run/want_rc/has/hasnt/denominator_is` 四把尺；`:143-273` 十枚用例；`:274-283` 汇总与 rc | AC#3 的常驻载具 |
| `scripts/testdata/portable-tests/go`（新增，100 行） | 假 `go`：`env GOOS`／`list`（含 GUARD A 的 `-f` 形）／`test`（`-list` 与真跑两形），全由 `FAKEGO_*` 驱动；未预期的调用 **`exit 99` 而非编一个像样的空成功** | 载具的注入面 |

`scripts/portable-tests.sh`：497 行 → **572 行**；`git diff --stat` ＝ **80 insertions / 5 deletions**，
`git diff -U0` 的 hunk 头只有四枚，全落在旧 `:259-271` 这一整段 resolve 里：
`@@ -259,0 +260,28 @@`／`@@ -261,5 +289,11 @@`／`@@ -267,0 +302,6 @@`／`@@ -271,0 +312,35 @@`。⇒ 改动**没有越出 resolve 那一步**。

形状判据的选择理由（AC#2 要"能力形"）：判的是 `go/src/internal/module` 对路径元素的文法（字母数字起头，元素内允许 `- . _ + ~` 与转义用的 `!`，`/` 分段），
**"整行不含空白"才是杀进度行那一刀**——`go: downloading x v1.2.3` 无论 Go 把它换成什么措辞，只要它是一句人话就带空格；
而单段路径（`fmt`、`errors`）仍通过，所以显式路径形态（`bash scripts/portable-tests.sh errors`）不会被这把新尺砍倒。

## ④ 改后读数 + 正控读数（四格各自对应）

一把总闸：`bash scripts/portable-tests-selftest.sh`（十枚用例，全绿＝rc 0）。
改后实跑：`rc=0`，`10 case(s) ran, 0 assertion(s) failed`，逐枚 `ok`（全文 `.scratch/wisp/probes/250/r1/logs/postfix-selftest-all.txt`，
逐用例日志 `.scratch/wisp/probes/250/r1/logs/cases/postfix-*.log`；`logs/` 顶层另有同名十枚是我 `cp` 的两份之一份，未删，规则 8）。

**"分母"这一格怎么读**：脚本成功时不印 `pkgcount`，但 GUARD B 对 `$resolved` 的**每一行**打一行 `ok (own line)  <path>`，
所以 `grep -c '^portable-tests\.sh:   \(ok\|FAIL\) (own line)'` ＝ 分母枚数。载具里这把尺叫 `denominator_is`。

### ④.1 AC#1 正控（票面点名的"stderr 有内容、stdout 完全正确"）——载具用例 `coldcache-stderr`

stderr 喂 CI 那 10 行逐字进度，stdout 喂 `core_pin` 抽出的 25 枚：
`rc=0` ｜ 分母 **25** ｜ `hasnt 'GUARD [ABC] -'` 命中 0 次 ｜ 回声句出现 1 次 ｜ `err| go: downloading ` 出现 **10 次**。

同形的**真 go** 版（一次性混合载具，不入库；`go list` 走真 `go`＋真空 `GOMODCACHE`，只有 `go test` 是假的）：

```
portable-tests.sh: go list wrote 11 line(s) to stderr (progress, not packages, not in the denominator):
portable-tests.sh:   err| go: downloading github.com/pelletier/go-toml/v2 v2.2.4
... 共 11 行 ...
portable-tests.sh: four numbers (all from -v output): === RUN=25  --- PASS=25  --- FAIL=0  --- SKIP=0
rc=0（26 秒）；GUARD 命中 0；分母（GUARD B 行）25
```
（`.scratch/wisp/probes/250/r1/logs/postfix-real-coldcache-hybrid.txt`；那 25 枚 PASS 是假 `go` 打的，
**这条读数只关于分母与守卫，不关于任何 Go 测试**。）

### ④.2 AC#2 正控（往 stdout 种一行带空格的伪包名）——用例 `seed-stdout-space`

`rc=1`，红句逐字（`.scratch/wisp/probes/250/r1/logs/cases/postfix-seed-stdout-space.log`）：

```
portable-tests.sh: GUARD C - go list's stdout carried 1 line(s) that are not
portable-tests.sh:   import paths, so this 26-line denominator cannot be trusted:
portable-tests.sh:   github.com/CarlosShao/wisp/internal/fake pkg
portable-tests.sh: refused rather than filtered: a guard that drops the lines it does not
portable-tests.sh: recognise measures less than it claims to (ticket 250 AC#2).
```
注意 `this 26-line denominator`：那 26 行**全留在分母里被点名**，一行没被 grep 掉——这正是票面 AC#2 要防的形状。

第二形 `seed-stdout-progress`（把**原缺陷那 10 行**种进 stdout，模拟"将来有人把 `2>&1` 写回来"）：`rc=1`，
`carried 10 line(s) that are not` ⇒ 分母的形状检查是**独立于分流**的第二道闸。

### ④.3 AC#1 反面（`go list` 真失败仍 `exit 1`，两路都打）——用例 `golist-fails`

`FAKEGO_LIST_RC=1`、stdout 有 26 行（25 真 + 1 伪）、stderr 有 10 行 ⇒ `rc=1`，且
`exited non-zero. Its stdout was:` ×1、`out| github.com/CarlosShao/wisp/internal/partiallyresolved` ×1、
`err| go: downloading modernc.org/sqlite v1.59.0` ×1、`the scope cannot be audited` ×1。

### ④.4 AC#4 三形各一枚"种 X 必响"（证明没被改软）

| 形 | 注入 | 读数 |
|---|---|---|
| GUARD A | `FAKEGO_GUARD_A_SEED=.../internal/session`（假 `go list -f` 打 `EMPTY `） | `rc=1` ＋ `GUARD A - these packages are declared in scope mode=core` ＋ 指名那枚 |
| GUARD B | `FAKEGO_DROP_PKG=.../internal/perm`（分母里有它、日志里没有它的结果行） | `rc=1` ＋ `GUARD B - 1 of .* packages in scope mode=core printed` ＋ 指名那枚 |
| 空 scope `exit 2` | 用例 `empty-scope-slice`：把真字节 `if [ ${#scope[@]} -eq 0 ]; then` / echo / `exit 2` / `fi` **四行原样切出来**（切片已打印在输出里）喂 `scope=()` | `rc=2` ＋ `refusing to be a green no-op` |
| GUARD C 两向 | `pin-drift-one-package-missing`（少一枚）／`seed-stdout-valid-but-unpinned`（形状合法但不在钉里） | 两枚均 `rc=1`，diff 分别打 `< …winsec` 与 `> …notpinned` |
| **字节未动** | `sha1sum` 对比 `git show 50d7a88:scripts/portable-tests.sh` 与工作版，按锚点抽块 | GUARD C 钉块 `5eb0993c4724c82c171c15e331788a24fe152b86` 两版相同；GUARD A 块 `d9aae29e6d9518de9566918d50007d4af582611d` 相同；GUARD B 块 `7b67fc33eeb4a82ee0c4caaa011e1b88207093cf` 相同；空 scope 四行 `70c95a865e22af511ac93467fb36a0caaec90b55` 相同；台账 `ledger=(…)` 块 `14fde2c1f98dc4d1ea6640da79a6e09cf326ff29` 相同 |
| 反证：同一把载具跑**改前**那份字节 | `bash scripts/portable-tests-selftest.sh all /tmp/portable-tests-prefix-250.sh` | `rc=1`，`10 case(s) ran, 10 assertion(s) failed`，红全落在 `coldcache-stderr`（GUARD C 假红 + 分母读不到）、`seed-stdout-space`／`seed-stdout-progress`（无形状检查）、`golist-fails`（两路未标注）；而 `clean`／`pin-drift`／`guard-a`／`guard-b`／`empty-scope-slice`／`seed-stdout-valid-but-unpinned` **改前改后同形**（⇒ 这几枚的行为我一行业没动）。逐字 `.scratch/wisp/probes/250/r1/logs/prefix-selftest-all.txt`，其中 `coldcache-stderr` 用例在改前逐字打出 **`Pinned: 25, resolved: 35.`**（`logs/cases/prefix-coldcache-stderr.log`）＝CI run `36889094435` 那一发的**本机确定性复现**，一把尺子内同时钉住 AC#3 的"改前必红/改后必绿" |

### ④.4b GUARD B 的第二受害者读数（票面没算那一枪）

见 §①.1 第 2 条与 `logs/second-victim-guardb.txt`：同一条 stderr 注入、GUARD C 不响的形态下，
改前 `rc=1`（`GUARD B - 10 of 35 packages ... printed` ＋ 10 行逐字点名），改后 `rc=0`（missing 0）。
这格是"我没把守卫改成少测东西"的另一面：**分母小了会被 GUARD B 咬，分母脏了也会被咬，两向都在。**

### ④.5 真 `go` 真测试的一发冒烟（证明改后脚本在非载具下仍能干活）

`bash scripts/portable-tests.sh ./internal/buildinfo/` ⇒ `rc=0`，16.9 秒，真 `go list`（温热，stderr 0 字节）、真形状检查（1 行合法路径）、
真 GUARD A、真 `runtests.sh`、真 `go test`：`four numbers ... === RUN=2  --- PASS=2  --- FAIL=0  --- SKIP=0`
（`.scratch/wisp/probes/250/r1/logs/postfix-real-warm-explicit-singlepkg.txt`）。
单包、非在飞包，未跑 `cmd/wisp`／`internal/panel`／`internal/ball`，未跑整树。

### ④.6 载具自身的两处坏尺（我先怀疑尺，改尺之后才下结论）

1. `run()` 里 `"$@"` 展开的 `VAR=val` 被引号保护后**不再被 bash 认作赋值前缀**，九枚用例全 `rc=127`。
   ⇒ 加 `env "$@"`。这是编排者说的"三把坏尺"在本腿的对应物，读数留在第一版 `postfix-selftest-all.txt` 的覆盖之前（已重跑）。
2. 空 scope 切片我按"三行"断言，实际锚点到 `fi` 是**四行**（`if`/echo/`exit 2`/`fi`），载具自判 `FAIL`。
   ⇒ 断言改 4，并把切出的四行**原样打印**在用例输出里，避免下一次再靠行数猜。

## ⑤ 越界检查

- 本腿提交两枚 commit，各带**显式 pathspec**（`git add -- <逐枚路径>`；`git add -A`/`.`/`commit -a` 一次未用）：
  - 代码：`scripts/portable-tests.sh`、`scripts/portable-tests-selftest.sh`、`scripts/testdata/portable-tests/go`
  - 证据：`.scratch/wisp/probes/250/r1/verdict.md` 与 `logs/**`
- 骨架 commit：`c9c8de95be8ce6ef0f49175215c6392cb8f53f55`（只含 `verdict.md` + `logs/prefix-real-coldcache.txt` 两枚路径）。
- 写面前后差集（09:11 现跑 `git status --porcelain scripts/ .scratch/wisp/probes/250/`）：
  `M scripts/portable-tests.sh` ＋ `?? scripts/portable-tests-selftest.sh` ＋ `?? scripts/testdata/` ＋ `?? .scratch/wisp/probes/250/r1/...`
  ——**没有第 5 枚路径是本腿写的**。起手 08:49 `scripts/` 为空 ⇒ 无他人写面被卷入。
- 全树 `git status --porcelain | wc -l`：起手 08:51 ＝ 304，09:11 ＝ 306（多出的两枚是 `.scratch/wisp/probes/250/r1/logs/` 与其 `cases/` 子目录，本腿所建；
  共享树里此数随别的腿一起漂，只在各自取数时刻有效）。
- `frontend/**`／`design/**`：零读取、零引用、零写入（起手 `git status` 顺带列出的 `design/...` 行是别的腿的痕迹，本腿未打开任何一枚）。
- 未跑整树 `go test ./...`；未跑 `cmd/wisp`／`internal/panel`／`internal/ball`。载具与混合载具里真 `go test` 一行未执行；
  唯一一次真实测试是 §④.5 的单包 `internal/buildinfo`。
- ⚠ **一处必须自曝的灰区**：为了拿"我的新文件会不会点红仓门"这一格读数，本腿跑了 `scripts/d22scan.sh` 的第 2 步
  （`cd tools/d22scan && go run . -root <repo>`）。那把门**自己**会走 `frontend/`（85 枚文本）与 `design/`（39 枚文本）。
  本腿没有打开那两棵树里的任何一枚文件、输出里也只有枚数与 `clean` 判语；但"零读取"这一格我按字面纪律**上报而不是藏**，
  要不要把"跑仓门即算越界"钉成规矩＝编排者裁。**下次同类腿默认不跑 d22scan 的整树步。**
- 一字未动：`docs/PLAN.md`、`docs/specs/**`、`docs/BUILD.md`、`docs/SLO.md`、`internal/observe/thresholds.go`、任何 golden、
  `tools/d22scan/allowlist.txt`、`tools/d22scan/runtests.sh`（含 `:98-102` 那句「SKIP is not a pass」）、三枚冻结测试件、`-skip` 名单（台账 11 行 sha1 见 §④.4）。
- 未写 `docs/reports/pending-and-issues.md`；未改票面任何勾选框；未 push、未 `gh`。
- 新增 `.go` 台件：**零枚**（两枚新文件是 bash；`scripts/testdata/portable-tests/go` 无 `.go` 后缀、且在 `testdata` 下，
  `go list ./...`／`go vet ./...`／`gofumpt -l` 的枚举都取不到它——`git ls-files -s scripts/` 里既有 `.sh` 多为 `100644`，
  载具在运行时 `cp`＋`chmod +x`，不依赖入库 exec 位）。故 `gofumpt` 一格无读数可贴。
  两枚新 `.sh` 与改后的 `portable-tests.sh` 均 `bash -n` 通过；两枚新文件纯 ASCII（`LC_ALL=C grep -n '[^ -~\t]'` 零命中），
  `portable-tests.sh` 里剩下的非 ASCII 只有既有注释 `:219`／`:345`（未动）。

## ⑥ 门禁读数（这一步在 CI 那台上会是什么颜色）

本机读数（能给的都给了）：

- 改前，真实冷缓存真脚本：`--scope=core` **红**，`Pinned: 25, resolved: 36.`，12 秒，`go test` 未起（§②）。
- 改后，同一把载具：十枚用例全绿、`rc=0`；同形的真 `go list` 混合载具：`rc=0`、分母 25、GUARD 零命中（§④.1）。
- 改后，真 `go` 真 `go test` 的单包显式形态：`rc=0`、四数 `=== RUN=2 --- PASS=2 --- FAIL=0 --- SKIP=0`（§④.5）。
- 本机其他门禁面：`bash -n` 三枚 `.sh` 全过；`cd tools/d22scan && go run . -root <repo>` ⇒ `rc=0`、
  `clean - no D22 ban violations`、`grep -nE "scripts/(portable-tests|testdata)"` 在输出里**零命中**
  （那把门的射程只有 internal/ cmd/ frontend/ design/，`scripts/*.sh` 今天不在其内——记一笔，别以为它看过我的文件）；
  `gofumpt` 一格**无读数**：本腿新增 `.go` 零枚（`scripts/testdata/portable-tests/go` 无 `.go` 后缀且在 `testdata/` 下，
  `git ls-files '*.go'`／`go list ./...` 都取不到它）。
- 载具跑改前字节：`rc=1`、10 枚断言红，且红只落在本票四格所指的用例（§④.4 末行）＝**这把尺不是永绿的**。

**CI 的颜色本腿不判**，三条具名理由：

1. 本腿不推送、不 `gh` 触发（硬边界 8）；`test-core` 第 7 步跑的是远端 SHA，本机 HEAD 不是它。
2. CI runner 的冷缓存枚数随实例而异（本机 11 行、CI 那发 10 行）：修好后**几行都不该进分母**，
   但"这一发真采到了 1365 枚 `=== RUN`"只能推上去看日志才能钉。
3. `ci-delta-1` 那条"同码二次采样"（schedule run `36940536372`）说明这枚缺陷在 CI 上是**概率性**的：
   一发绿不能证明修好，得连看两发以上冷拉起——那是编排者的核过窗口，不是本腿的断言。

本腿能钉的承诺：`bash scripts/portable-tests-selftest.sh` 在 CI 那台 ubuntu runner 上同样应 `rc=0`
（它只依赖 bash、awk、sed、grep、mktemp，不依赖网络、不依赖 `GOMODCACHE` 冷热、不建任何 Go 测试件）。

## ⑦ 判不动的地方

1. **CI 上 `test-core` 改后的真实颜色与 `=== RUN` 枚数**——判不动：禁推送禁 `gh`，冷缓存只在 runner 实例上随机出现。
   能判＝编排者（推上去连看两发冷拉起，或对比 schedule run `36940536372` 那一族）。
2. **那 4 枚"消失"的前端契约红到底是变绿还是被跳过**——判不动：判据要在真 `go test` 跑起来之后才存在，
   且其中三枚在 `internal/panel`（在飞腿 `248-r1` 写面 + 两枚冻结测试件射程内）。能判＝编排者 + `248-r1`。
3. **假 `go` 载具证的是"形状"，不是"真 `go list` 的 stdout 永不掺进度"**——载具的 stdout 由我写；
   §②/§④.1 的真实冷缓存两发只证明"stderr 有货时 stdout 仍干净"。要把 stdout 纯度钉成契约，能判＝编排者
   （或在 runner 上取一发的验收腿）。
4. **`--scope=windows`／`cli` 两枚钉未逐钉配载具**——共用同一段 resolve 代码，但载具只造了 core 一发钉。
   扩到三枚钉各一发要不要做＝编排者定（票面 AC 未点名，本腿不擅自扩面）。
5. **`winsec_pin`（`:164-166`）从 CLI 取不到、GUARD C 对 winsec 那枚钉今天不跑**（见 §①.1 第 3 条）——
   改它＝往 `case $mode` 里加一支，超出本票"一块石头"的射程。能判＝编排者（另立票或并进门控族）。
6. **`GOMODCACHE` 冷缓存这发读数依赖网络与 goproxy.cn 可达**——离线主机上 §② 那形拿不到（载具不依赖，仍是绿的）。
   把"离线也要能判"钉进 CI＝编排者定；本腿选择让常驻载具**不依赖网络**（假 `go`），代价是它只判形状。
7. **`/tmp` 里三枚冷缓存目录共约 0.9 GB 未删**（`/tmp/tmp.2ztPH1RoMJ`、`/tmp/tmp.GaMxazmt6T`、`/tmp/tmp.rDi8OXsiVX`）——
   按"临时件只建不删"留着，取数时刻 09:0x；能判＝编排者（一句 `rm -rf` 的事，但那不是本腿的授权）。
