# 票 251 r1 交件读数（写码腿 `251-r1`）

> 本节只交**读数**，判语归非实现者。所有行号／枚数＝本腿本机自跑，不是抄票面。

## 0. 起手锚

| 项 | 读数 |
|---|---|
| 取锚时刻 | `2026-10-02T09:29:26+08:00`（`date -Iseconds`） |
| 起手 HEAD | `f5f9cc3433385c13e24cb2a7d8054d8210f99cde`（`dev`，`git log -1`） |
| 起手 `git status --porcelain scripts/` | **空输出**（rc=0）⇒ 这一刻没有别的腿在 `scripts/` 这面上，未触发票面「交件要求」第 1 条的停手 |
| 起手 `wc -l -c` | `portable-tests.sh` 572/32543 · `portable-tests-selftest.sh` 313/13354 · `testdata/portable-tests/go` 100/4159（`wc -l` 少算末行无换行的情形，此处三件均以换行收尾，枚数一致） |
| 起手载具基线 | `bash scripts/portable-tests-selftest.sh` ⇒ `10 case(s) ran, 0 assertion(s) failed`、**rc=0**（`logs/selftest-baseline.txt`） |
| 起手真 census | `bash scripts/portable-tests.sh --scope=census` ⇒ rc=0，`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=9`，GOOS=windows（`logs/census-prefix.txt`） |

起手锚之后所有卫生尺都以 `git diff f5f9cc3..HEAD -- scripts/portable-tests.sh` 为准（§3）。

## 1. 改前必红（缺口凭据）

### 1.1 票面「现量」三把尺，本腿复跑结果

1. `bash scripts/portable-tests.sh --scope=winsec` ⇒ 逐字（stderr，`logs/prefix-scope-winsec.txt`）
   `portable-tests.sh: unknown --scope=winsec (known: core, windows, cli, census)`、**rc=2**。
   ⇒ 复认票面第 12-14 行：这一形**是响亮失败，不是静默绿**。
2. `--scope=census` 的真读数里 `internal/winsec` 那一行（`logs/census-prefix.txt:36`）逐字
   `portable-tests.sh: github.com/CarlosShao/wisp/internal/winsec     12/8        corewinsec`
   ⇒ census 把一枚 **`case $mode in` 里根本不存在的档名**（`winsec`）报成"有人认领"，且这一步 **rc=0**。
   这就是 AC#4 说的"两处各写一遍档名"今天的实际状态：**两处已经不等，而没有任何一响**。
3. census 全表里 `github.com/CarlosShao/wisp/internal/winsec` 只出现 **1 行**（`grep -c` 尺见 §2 AC#1 段），
   ⇒ 今日 `internal/winsec` 在模块里**没有子包**，故 `./internal/winsec/...` 与 `./internal/winsec/` 现在解析出同一枚单包
   —— 本腿把第五档写成 glob 形（§2 AC#1）正是踩在这条实测上，不是猜的。

### 1.2 两形种子打在**改前字节**上（拆包／删真包名，都走 winsec 那族今天真用的显式路径）

改前字节＝`git show f5f9cc3:scripts/portable-tests.sh` 落到 `/tmp/251-prefix-portable-tests.sh`，
用 `scripts/portable-tests-selftest.sh` 的 shadow-root 机制（`$2` 位）跑，`go` 一律是
`scripts/testdata/portable-tests/go`（假 `go`），⛔ 全程未起真 `go test`／真 `go build`。

| 种子形 | 命令 | 改前读数 | 改后读数 |
|---|---|---|---|
| 拆包：`internal/winsec` 多出第二包 | 假 `go` 解析出 2 行（pin 的那行 + `.../internal/winsec/acl`），调用形＝`bash portable-tests.sh ./internal/winsec/` | 续写在 §1.3 | 续写在 §1.3 |
| 删真包名：`internal/winsec` 一行都解析不出 | 假 `go` 对 `./internal/winsec/` 返回 0 行 | 续写在 §1.3 | 续写在 §1.3 |

§1.3 ＝ 把上面两发的逐字日志与退码摆出来（改前一发、改后一发，同一载具、同一种子），本表任何一行不许引用空节。

## 2. 逐格 AC

写法选择（AC#1 给了两条路，本腿**两条都走**，凭据在 §1.2 与 §2 各格）：
第五档 `winsec)` 进 `case $mode in`（scope 写成 `./internal/winsec/...` glob + `pinned=$winsec_pin`）；
**并且** GUARD C 在显式路径模式下也对账（CI 今天真走的那条路是 `winsec-tests.sh:94` 的显式路径，
只做前者等于再造一枚"钉了没人点"的档）。选型理由与残余限界写在 AC#1 段末。

| 格 | 读数出处 | 判语（归非实现者） |
|---|---|---|
| AC#1 钉咬得住 | 续写 | |
| AC#2 尺本机可判 | 续写 | |
| AC#3 不动别人的形状 | 续写 | |
| AC#4 census 自证 | 续写 | |

## 3. 门禁读数

- 载具：`bash scripts/portable-tests-selftest.sh` 的 case 数／失败数／rc（改前基线 10/0/rc=0，见 §0）。
- 卫生：`git diff f5f9cc3..HEAD -- scripts/portable-tests.sh` 的 **hunk 数与每枚 hunk 的起始行**逐枚摆出来，
  并说明每一枚落在哪一块（档位表／显式路径块／census 块／注释块），GUARD A／GUARD B／GUARD C 既有比对块的字节不许出现在改动行里。
- 占位符尺（字符类形，交件前跑）：0 命中。
- 未跑过的档：`--scope=core`／`--scope=windows`／`--scope=cli` 真档**一律没跑**（会起真 `go test`，
  `cmd/wisp`＋`internal/config` 正被腿 `198-r1` 占、`internal/ball` 被腿 `33-r8b` 占）；
  改后的 `--scope=winsec` 真档同理**没跑**——它现在会往下走进真测试，其读数一律由假 `go` 载具代取，逐字写明在哪一格。
- 空 scope 的 `exit 2`、unknown `--scope` 的 `rc=2`：两枚读数在 AC#3 格。

## 4. 本腿推翻编排者题面哪几句

| 题面位置 | 题面原话 | 实测 | 证据 |
|---|---|---|---|
| 票面 `:11` | `case $mode in`（`:170`）"只有 `core)`／`windows)`／`cli)`／`*)` **四支**" | **漏计 `census)`**：改前那份字节里是 `core)`(:171)／`windows)`(:184)／`cli)`(:191)／`census)`(:198)／`*)`(:202) ＝ **4 枚具名档 + 1 枚兜底，共 5 支** | `git show f5f9cc3:scripts/portable-tests.sh` 的 170-206 行；枚数以 `grep -c` 尺收尾 |
| 票面 `:15` | "winsec 那族测试真正进 CI 的路径是 `scripts/winsec-tests.sh:97` —— `bash "$portable" "${scope[@]}"`" | 该行在 **`:94`**；`:97` 是 `ran=$(count '^=== RUN')` | `grep -n 'bash "$portable"' scripts/winsec-tests.sh` ⇒ `94:` |
| 票面 `:10` | "`winsec_pin` 今天**唯一的读者**是 `:229-242` 那圈 census 循环" | **行号范围偏小**：那圈 while 循环是 `:226-242`（`:226` 是 `while IFS= read -r p; do`，`:229` 已在体内）；"唯一读者"这一支**复认成立**（`case` 里无 `winsec)`，显式路径 `pinned=''`） | `logs/census-prefix.txt` + 改前那份字节的 220-243 行 |
| 票面 `:9` | `:164-166` 的 `winsec_pin` ＝ **1 行** `github.com/CarlosShao/wisp/internal/winsec` | **复认成立**，逐字一致 | 改前那份字节的 164-166 行 |
| 票面 `:3` | 立票锚点 `db06a399` | 本腿起手锚是 `f5f9cc34`（票 250 结案之后又有 census(114-a2) 等提交落在 dev 上）；本票所有"改前"尺一律打在 `f5f9cc34` 那份字节上，⛔ 未拿 `db06a399` 当改前 | `git log -1` |

## 5. 判不动的地方（未定义即停）

- **甲（谁补得上）／乙（补不上）**逐条写在下面，本腿不自作主张。
