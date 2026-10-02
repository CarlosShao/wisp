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
