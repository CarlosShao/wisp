# 票 174 · `174-v2` 验收腿裁决（AC#6 形 乙；非实现者）

- 时刻 2026-10-14 14:30–14:5x +08；起手 HEAD `4c7f0726`，交卷 HEAD `76b9e214`（他腿 253-r6 收件＋编排者修复 `runtests.sh`；`internal/tools/` 与 `portable-tests.sh` 零重叠，被测 blob 仍 `15e0ef42…`＝ `a65bb386` 的 blob，工作树零差）。
- 被裁＝落地腿 `174-r6` 两笔：`a65bb386`（测试 `+135/−0`，新用例 `…Portable174r6`）＋`11cfbc93`（件＋票尾）。
- 三把尺复核（编排者给的）：两 commit 在（13:34:06／13:45:52）｜件齐（`probes/174/r6` 八节＋logs 五枚）｜亲跑 `-run 'TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath'` = **`ok 1.959s` rc=0**。
- 原始读数：`.scratch/wisp/probes/174/v2/`（`00-anchor-and-gates.md`＋logs 五枚；raw 18 枚留盘）。本件每句可回指。

## 0. 判语

**形 乙＝成立。** 同一判据（canonicalize-error notice 不许转发 filed path）在 POSIX 上有真可跑的形：空白字面量在 `paths.go` 的 `TrimSpace` 闸（`:131-134`，先于 `risk.Resolve`）被拒 ⇒ `pointerNotice` 的 `case err != nil`（`task.go:839`）平台无关可达；承重未降层（真 Bridge＋真 C26＋真名册＋真毒值根目录）。

- **接受**r6 自报的强度差（"POSIX 上路径半边弱一档"）：如实、且补强（4×Unicode 空格锐序列＋`%q` 拼法）已落地。**一处更正它没声称的强**：新枚在"转发行为"上另抓两档 r3 臂抓不到的形——**M2**（转述 `Roots()`，命中 fold 后小写形）与**我的 M4**（`strings.ToLower(raw)`，全小写路径穿 r3 扫而过）。新枚的双拼法禁表让它 ≥ r3，不是"弱化替代"。
- **无保留性反对**；两条"如实登记"照单接受：① 被拒字面量是空白串而非毒值路径字面量（POSIX 唯一可达，判据本体不受损）；② 判据的路径半边在 POSIX 上永远大小写无关（空白），故 M4 一类 fold 规避对它是**不可能形**而不是漏网——反而更严。
- ⛔ 本腿未翻任何框、未动 `Status:`/票名，判语与收口建议交回编排者。

## 1. 独立四证（逐条现跑，不复引 r6）

| # | 证 | 读数 |
|---|---|---|
| 1 | windows 双跑（dispatch 原命令） | rc=0；r3 **真建 mklink /J 并 PASS**（notice 逐字 `risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions…`）＋新枚 **PASS**；logs/01 |
| 2 | linux 静态 | `GOOS=linux go vet` rc=0；`go test -c` rc=0；linux 二进制含**新名 ×3**（r3 名 ×3、运行时才 SKIP）；`//go:build`=0 行；新枚体内 `t.Skip`=0。logs/02 |
| 3 | linux 实跑（Docker linux/amd64，**非 CI ubuntu**，如实） | r3 名 **SKIP**（逐字 `exec: "cmd": executable file not found in $PATH`）＋新枚**真跑 PASS**（notice 与 windows 同串 `（tools: empty path）`）；rc=0。logs/03 |
| 4 | r3 行为零回归 | 同发 r3 PASS；r3 侧字节 1-522 行与原版逐字节同、六个锚（`:376/:408/:420/:428/:439/:486`）同号同文；ledger 行 sha `b1ffdb59…`＝r5 落地时。logs/05 |

## 2. 反证四发（mutation 全 `-overlay`，tracked 零字节；还原 sha 证在 logs/04）

| 突变 | 形 | windows 红（16 名全家谱） | linux 红（docker） |
|---|---|---|---|
| M1（r6 的） | `（"+err.Error()+"："+raw+"）` | 恰 {r3, 新}（14 绿） | 新（r3 SKIP） |
| M2（r6 的） | `；授权根 "+strings.Join(Roots(), "、")` | 恰 {新}；r3 绿（**fold 小写漏**，verbatim `…\\rootlist-leak-probe-174r6-a` 命中） | 未跑（r6 亦未跑；禁表平台无关） |
| M3（我造的**新形**） | `，"+fmt.Sprintf("%q", raw)` | {r3, 新}（r3 靠引号内 token 原样命中；新枚命 `%q` 专禁） | 新：`must not relay "\"\\u00a0\u2003\u2009\u205f\""` |
| M4（我造的**新形**） | `，"+strings.ToLower(raw)` | 恰 {**新**}；r3 臂绿**且其 notice 日志证明漏**（`…，c:\\…\\allowlist-leak-probe-174r3）`） | 新：`must not relay "\u00a0\u2003\u2009\u205f"` |

- ≥2 发自造、两发（M3/M4）都是 r6 没种过的形；windows 红名逐枚点名如上；M4 linux 红＋M1 linux 红＝"POSIX 本体上判据也咬得住"。
- 一处**自摆乌龙如实登记**：M3 首次 linux 跑传了 host 键 overlay ⇒ 容器内 no-op、假绿 rc=0；我识别（notice 未变异）后换 `/wisp` 键 overlay 重跑得红，假读数不记入案。
- 还原证：`task.go` sha `e7419719…`（突变前后同值）＋ `git status --porcelain -- internal/tools/` 空＋`git diff --numstat` 空。

## 3. 共存（r3 半边一字节未动）

- `a65bb386` = `135 0`（纯追加，单 hunk `@@ -523,6 +523,141 @@`）；文件 1-522 行同旧版；r3 fixture/helper/两枚用例原名原位。
- ledger：`grep -c` r3 名 = **1**（`:682`，五字段名|pkg=./internal/tools/|linux|fixture|reason 含 `OWNER of the POSIX half: whoever lands AC#6 shape 乙…`）、新名 = **0**（不入名册）；行 sha 与 r5 的 `32f9823f` 逐字节同。
- `probes/174/{r1,r2,r3,r5}` porcelain 空；`2b31d728..HEAD -- probes/174/` 只增 r6 六件。
- **下一发 CI linux 腿判法（写清）**：名册机制对 linux-active 名做 `-skip '^(…)$'`（r3 名在内）⇒ r3 名**不出现在输出**（无 `--- SKIP`）；新名不在名册、不可能 SKIP（0 tag／0 t.Skip）⇒ **真跑、应在 PASS 侧；本对贡献 SKIP=0**。docker 已用同形 `-skip` 预演（只新名 RUN/PASS，rc=0）。

## 4. 门禁现量（各现跑）

`go build ./...` rc=0 ｜ `GOOS=linux go vet ./internal/tools/` rc=0 ｜ gofumpt v0.12.0 `-l` 空 ｜ `scripts/d22scan.sh` **clean** rc=0（`ban #8 internal/`=550、`cmd/`=126）｜ 长度门 **VERDICT GREEN** rc=0（tracked=12755、over-budget=60、roster=60、not-in-roster=0 ⇒ 与派单"现 60=60"一致）。

## 5. 风落一枚（非 r6；14:38 撞见，14:49 已被编排者修复，本腿复验）

`bash scripts/d22scan.sh` 在 **GOTMPDIR 未设**时 rc=1：`tools/d22scan/runtests.sh:92: LOCALAPPDATAwisp: unbound variable`。该行 `gotmp_run="${WISP_TEST_GOTMP_ROOT:-"$LOCALAPPDATA\wisp-test-gotmp"}/…"` 是 `581dd252`「微修 runtests.sh 第四处折法」引入的**单反斜杠**形；`set -eu` 下反斜杠吞进变量名。三处同类脚本都是**双反斜杠**且实测正确（`portable-tests.sh:155`／`winsec-tests.sh:69`／`wisp-cli-tests.sh:82`）⇒ 第四处镜像少打一枚。仅波及**本机 windows**（ubuntu CI 不进该分支）。
**收尾（不欠账）**：同窗口 `253-r6` 也撞见同源；编排者 `76b9e214`（2026-10-14 14:49:32）单→双修复并登记 `A929`。本腿**修复后复验**：GOTMPDIR 未设 ⇒ `d22scan clean rc=0`（修复前 rc=1，两读数都在 logs/05）。六形对照 repro（仓外 `C:\tmp\wisp174v2gotmp\parse-test.sh`）与修复 diff 逐字对上 ⇒ 记为**独立复现**，不是新账。

## 6. 票 174 整票框态与收口建议

现读（HEAD `4c7f0726`，`grep '^- \[[ x]\]'`）：**已勾 4**＝AC#1、AC#2、AC#2d、**AC#6**；**未勾 5**＝AC#3、AC#2b、AC#2c、AC#4、AC#5。⇒ **整票不可收口**；本腿对象可收。

- **AC#6（已 [x]；13:18 甲形 CI 实证，遗留句"POSIX 半（形乙）由 r6 落地中"）**：本腿判语＝形乙成立 ⇒ **遗留可标结**。建议动作：编排者在票面追加一行（或并入 A 账）写"形乙落地已由 `174-v2` 收口（确立；docker linux 实跑＋四发反证＋下一发 CI 为最终确认）"；**框不用再动**（已 [x]）。留档句建议＝「AC#6 两形齐（甲＝ledger 行 CI 实证 run 38508311376；乙＝`a65bb386` POSIX 判据，`174-v2` 判成立 2026-10-14）」；撤销口令＝「撤 174 AC#6 乙」。
- **AC#3（未勾）**：三条禁区。现量：`internal/risk/**`、`cmd/**`、`PLAN.md`／`SPEC-03` 由各腿反复证零字节；(ii) 的 C26 唯一决定者钉有 `174-r4` 两腿四枚（有牙）。建议：派一枚非实现者腿按 r4 表出"三禁区自证"判语，或编排者据 r4＋本腿（产码零字节链）裁"可否结"。
- **AC#2b（未勾）**：装配行已落（145-r3 补 `Paths`），缺端到端三臂读数＋"或批一张 L2 卡"第二半句（撞 `:299/:349` 模板冻结钉＝Q-63 边界）。建议：动文字前先裁"算不算动冻结文字"；端到端要一台不缺 DLL 的宿主。
- **AC#2c（未勾）**：`spill.go` 桩，票面自己写"排程等票 177 裁完再动"。建议：维持未勾；收口说明里逐条写明（AC#5 同款要求）。
- **AC#4（未勾）**：契约轴，须 owner 单独批准（未批前零字节）。建议：owner 不批则收口时**明写**"文字不改/未批"，不许静默。
- **AC#5（未勾）**：收口动作本身——逐条说明 `Q-59` 与本票各自状态。建议：编排者收口时写两行（本票＝授权根指针诚实性，已由 AC#2/2d＋形乙钉住；`Q-59`＝偏移参数，另账）。
- 结论：**本轮可收的只有 AC#6 遗留句**；整票收口要等 2b/2c/4 的裁决与 3/5 的结题陈述。

## 7. 没做的（具名）

CI `test-core` 真颜色（推不了；docker≠ubuntu，如实）｜整包 `-race` 与全量 `portable-tests.sh`（他腿同包＋脚本无 dry-run）｜端到端 `wisp run`（本机 DLL 存量事实）｜M2 的 linux 重跑（与 r6 同口径未跑；其禁表平台无关，M1/M3/M4 已给 linux 红）｜`internal/risk/**`、`cmd/**` 新内容未评（零字节）。
