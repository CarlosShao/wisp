# 票 251 — `winsec_pin` 是一枚"钉了没人对账"的钉：winsec 那档走显式路径，GUARD C 今天看不到它

**立票时刻**：2026-10-02 09:2x +08，锚点 HEAD `db06a399`（`dev`）
**来源**：`250-r1` 交件时具名上报（它推翻我题面第 4 条），我现量复认后立票；不夹进票 250（那票四格是 GUARD C 的分母口径，本票是"另一档根本没进那道门"）
**性质**：仪器／CI 口径缺口，⛔ 不是产品行为，owner 用软件撞不到

## 现量（三把尺都是编排者本机自跑，2026-10-02 09:2x）

1. `scripts/portable-tests.sh:164-166` 里钉了一枚 `winsec_pin`，内容＝**1 行** `github.com/CarlosShao/wisp/internal/winsec`。
2. 它今天**唯一的读者**是 `:229-242` 那圈 census 循环（`--scope=census` 用来回答"这包归哪个档"，`:233` 处 `winsec) pin=$winsec_pin ;;`）——**只是把钉当名册查，从不拿它去和实际解析出的包集作差**。
3. `case $mode in`（`:170`）只有 `core)`／`windows)`／`cli)`／`*)` 四支，⛔ 没有 `winsec)`；
   我现跑 `bash scripts/portable-tests.sh --scope=winsec` ⇒
   逐字 `portable-tests.sh: unknown --scope=winsec (known: core, windows, cli, census)`、**`rc=2`**。
   ⇒ 这一形**是响亮失败、不是静默绿**（这条很重要：上报那句"GUARD C 对 winsec 根本不跑"字面对，但容易被读成"那档悄悄少测了"；实际是"有人想点它，脚本当场拒绝"）。
4. winsec 那族测试真正进 CI 的路径是 `scripts/winsec-tests.sh:97` —— `bash "$portable" "${scope[@]}"`，**传的是显式包路径**，不走任何具名档 ⇒ `pinned` 为空 ⇒ GUARD C 的钉-vs-解析对账**在这一档一次都没执行过**。

## 不做的话，哪天会怎样（人话）

`internal/winsec` 这包要是哪天拆出第二包（或某枚腿把它的名字改了），CI 里那一步**照跑照绿**，只是少测了一块——这正是本仓反复钉的"绿不等于测过"那一族；而 census 那张表会继续把旧名字报成"有人认领"。

## 要建什么

- [ ] **AC#1 让那枚钉真的咬得住**：`internal/winsec` 那一档要么进 `case $mode in` 成为第五档（则 `winsec_pin` 与解析集每次都对账），要么 GUARD C 在"显式路径模式"下**也**做一次钉-vs-解析对账（谁能证明选了后者、给出凭据）。判据＝种一枚"winsec 包名被改／少一行"的异常 ⇒ 指名那一步必须红；改之前它必须绿。
- [ ] **AC#2 尺必须本机可判**：⛔ 不许写成〔仅 CI 可量〕。载具形沿用票 250 的 `scripts/portable-tests-selftest.sh`（同一台件、同一批 case 名风格），往 stdout 种假包名与删真包名两形各一枚正控。
- [ ] **AC#3 不动别人的形状**：空 scope 硬退出 `exit 2`（含 `--scope=winsec` 现在那句 unknown ⇒ `rc=2`）**一字不许软**；⛔ 不许为了"让 winsec 能点"把 unknown scope 改成默认档或静默跳过。GUARD A／B／C 的既有比对块用 `git diff <锚>..HEAD -- scripts/portable-tests.sh` 证明只有本票射程内的 hunk。
- [ ] **AC#4 顺手把 census 的自证补上**：`--scope=census` 的"这包归哪个档"读数今天没有任何钉证明它和 `case $mode in` 的档位表**同源**（两处在 `:170` 与 `:233` 各写了一遍档名）。判据＝档名集合在两处不等时，census 那一步必须红（或明确具名登记为"故意分两处、由 X 钉住"）。

## 禁区

- ⛔ 不许顺手把 `internal/winsec` 的**测试内容**改了来"凑对上"——本票射程＝脚本口径。
- ⛔ 不许新增第六档而不交代 `winsec_pin` 之外那三档的钉是否也缺读者（AC#4 就是这一格）。
- ⛔ 不许用"把 GUARD C 对显式路径也跳过"来实现 AC#1 ——那是把门拆了不是把门修好。
- ⛔ 不碰 `frontend/**`／`design/**`（**既不读也不引**）；⛔ 不碰 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden。
- 只在接缝注入；⛔ 不许 `t.Skip` 让哪一步变绿；⛔ 不许把 SKIP 读成绿。

## 交件要求（腿只交读数，判语归非实现者）

- 起手先 `date -Iseconds`＋`git log -1` 取锚，并 `git status --porcelain scripts/` 确认这一刻没别的腿在这面上；有则**停手上报**，不许并行写。
- 交件必带：**改前必红**的那一发（种异常 → 指名步红）与改后同一载具的绿；`wc -l -c` 实测尺寸；占位符 0（尺用字符类，别把尺文本写进被扫文件）；⛔ 表里任何一行不许引用一节空的凭据。
- 写明**它推翻编排者上面哪一句**（上面那些行号与"1 行"的枚数全是待验断言）。

## 归口与排程

- 排在票 250 之后、与票 250 同一射程（`scripts/`），**同面串行**；本票不与任何在飞写腿并发。
- 与 `winlive-census-1`／`pushpreflight-1` 那两笔 CI 账无关，别混。

---

## 落地回执（腿 `251-r1`，2026-10-02 10:1x +08）

> ⛔ 四枚 AC 勾选框一枚未碰（`grep -c '^- \[ \]'` 交件前后都是 **4**、`^- \[x\]` 都是 **0**）。
> 全节只交读数，逐字日志在 `.scratch/wisp/probes/251/r1/`（`verdict.md` + `seed-winsec.sh` + `logs/` 17 份）。

**起手锚** `f5f9cc34`（dev，`date -Iseconds` 2026-10-02T09:29:26+08:00）；起手 `git status --porcelain scripts/` 空 ⇒ 无并发写腿。
本腿四笔（逐枚 hash＋pathspec）：`2ca152d7`（`.scratch/wisp/probes/251`）· `709d2589`（`scripts/portable-tests.sh` + `scripts/portable-tests-selftest.sh`）· `4cca2c0e`（`scripts/portable-tests.sh`）· `aa0dbb69`（`scripts/portable-tests-selftest.sh` + `scripts/testdata/portable-tests/go`）。⛔ 未 push。

**改了什么**：`case $mode in` 加第五档 `winsec)`（scope 用 glob `./internal/winsec/...`）；显式路径块认出"这正是 winsec 档自己的那一条"（`scripts/winsec-tests.sh:94` 递过来的 `./internal/winsec/`）后按 `mode=winsec` + `winsec_pin` 走**既有** GUARD C 比对块；档名只写一处（`tiers=`），census 在 `go list` 之前先把名册与本文件自己的 `case` 分支集合双向作差、不等即拒答（rc=1），循环里没钉的档名同样拒答。载具 10→18 枚 case；假 `go` 新增 `FAKEGO_SCOPE_FILTER`（dir 形只认那一包、glob 形认子树），否则"拆包会响"是仪器自己造的假象。

**改前必红那一发**（同一载具、同一种子、CI 真走的显式路径）
- 改前（锚点字节，假 `go`）：`rc=0`，全日志一字 `GUARD` 都没有，末行逐字
  `portable-tests.sh: four numbers (all from -v output): === RUN=2  --- PASS=2  --- FAIL=0  --- SKIP=0`
  （`logs/prefix-split-full.txt`；删真包名那一发同样 rc=0，且解析集为空也照绿，`logs/prefix-gone-full.txt`）
- 改后同一发逐字红句：
  `portable-tests.sh: GUARD C - scope mode=winsec resolved to a DIFFERENT package`
  `portable-tests.sh:   set than the one pinned next to it. Pinned: 1, resolved: 2.`
  `portable-tests.sh:   > github.com/CarlosShao/wisp/internal/winsec/acl`（`logs/postfix-split-full.txt`）
  删真包名那一发 ⇒ `Pinned: 1, resolved: 0` + `< github.com/CarlosShao/wisp/internal/winsec`（rc=1）。
- 整台件打在改前字节：`18 case(s) ran, 20 assertion(s) failed`、rc=1；其中两枚显式路径种子在改前都是 rc=0（＝缺口本体），`census-tier-roster-*` 与 case 18 报"种子在这份字节上没有可下的地方"而**不是**报绿。

**四格读数**
- AC#1：两半都交（第五档＋显式路径对账）。选型凭据不是话术，是跑出来的反证：只按题面 (b) 的字面形（配钉、scope 仍用调用方的目录形）在同一颗拆包种子上 **rc=0 绿**，本腿交的形 **rc=1 红** ⇒ "glob 承重"已固化成载具 case 18 `winsec-explicit-glob-is-what-bites`。
- AC#2：本机可判，载具 `bash scripts/portable-tests-selftest.sh` ⇒ **18 case / 0 assertion failed / rc=0**；两形正控齐（种假包名 case 12/13/18、删真包名 case 14），case 名沿用票 250 风格。
- AC#3：`bash scripts/portable-tests.sh --scope=winsecfoo` ⇒ **rc=2**、逐字 `unknown --scope=winsecfoo (known: core, windows, cli, winsec, census)`；`--scope=` 同 rc=2；空 scope 那 4 行字节未动，载具 case 10 仍逐字切片跑出 rc=2。卫生尺：`git diff f5f9cc3..HEAD -- scripts/portable-tests.sh` ⇒ **7 枚 hunk**，全部落在锚点 69-255 之间；8 枚删除行的锚点号＝73/203/230/248/249/250/251/252，**最大 252 < 空 scope 块的 254**，GUARD C 比对块（锚点 353-368）／GUARD A（370-388）／GUARD B（513-568）／shape check（312-345）／ledger（390-406）零行进入删除集。
- AC#4：选了"钉住"这支。正向：真 `--scope=census` 仍 rc=0，与起手那份逐字对差只剩别人的一行（`internal/tools` 测试枚数，出处 `3465dcee`）；反向：名册删 `winsec`、分支留着 ⇒ **rc=1** 逐字 `census - the tier roster (tiers=) and the 'case $mode in'` + 两张集合都打出来（`logs/postfix-census-drift-full.txt`）。

**本腿推翻／修正题面哪几句**（详见 `verdict.md` §4，七条）
1. `:11` "只有 `core)`／`windows)`／`cli)`／`*)` 四支" ⇒ 漏计 `census)`：锚点是 **4 枚具名 + 1 枚兜底＝5 支**（awk 抽取尺⇒具名 4）。
2. `:15` "`scripts/winsec-tests.sh:97`" ⇒ 那行在 **`:94`**；`:97` 是 `ran=$(count '^=== RUN')`。同一处错号还在 `docs/reports/pending-and-issues.md:10505`（本腿不改台账，只报）。
3. `:10` "唯一读者是 `:229-242` 那圈 census 循环" ⇒ "唯一读者"成立，**行号是 `:226-242`**。
4. `:23` AC#1 判据"种'包名被改／少一行'的异常" ⇒ **覆盖不到票面自己 `:19` 那句"拆出第二包"**：目录形解析集看不见旁侧新包，题面 (b) 的字面形在同一颗种子上绿（作用面本腿跑过，见上）。
5. 隐含前提"假 `go` 的输出＝解析读数" ⇒ 不成立，原 shim 对任何参数打印整份名册（本腿第一次跑反证就中了这个雷，那两发已标作废：`logs/variant-b-*.txt`）。
6. `:9` `winsec_pin`＝1 行、`:3` 立票锚点 `db06a399` ⇒ 前者复认成立；后者本腿未当改前用，所有"改前"尺打在 `f5f9cc34`。

**没跑过／判不动**（`verdict.md` §3 末行、§5 全表）
- ⛔ 真 `go test`／`go build` 一发未跑（`cmd/wisp`＋`internal/config`＝腿 `198-r1`、`internal/ball`＝腿 `33-r8b` 在飞），故第五档与 `scripts/winsec-tests.sh` 的**真机 rc 读数缺失**，属〔待验〕；本腿只跑过 `--scope=census`（只做 `go list`）与取 rc 的 unknown 档。
- `ci.yml` 未动：第五档目前**没有 CI 调用者**（CI 仍只 `run: bash scripts/winsec-tests.sh`，`ci.yml:439`）。要不要把 workflow 指向新档＝契约级，本腿**停手未动**（§5 第 2 条）。
- `scripts/winsec-tests.sh` 不在本腿写面 ⇒ 改它去点 `--scope=winsec` 这件事留给编排者（§5 第 1 条）。
- core 档里那行 `./internal/winsec/`（`:211`）未改成 glob：改了会把新子包算进 core 的解析集、进而要求 `core_pin` 同步加行＝射程外（§5 第 3 条）。
- 占位符尺（字符类）交件前 **0 命中**；尺寸 `wc -l -c`：`portable-tests.sh` 668/38629 · `portable-tests-selftest.sh` 529/24878 · `testdata/portable-tests/go` 143/6280。
