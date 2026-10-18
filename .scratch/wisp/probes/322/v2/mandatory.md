# 票 322 终裁腿 322-v2 — 必答 (i)–(v)（判语件 `docs/evidence/s1/322-v2-verdict.md` 的具名分件）

> 全部读数＝本腿自己现跑（指名那一个原始 log／本腿 fresh 发），⛔ 抄腿、⛔ 抄编排者。
> 尺与落点＝verdict 第 1–5 节的读数件（`.scratch/wisp/probes/322/v2/logs/**`）。⛔ 翻任何 AC 框。

## (i) AC#4 总判语＝**成立**（逐条凭据：尺＋射程目录＋量/算）

1. **容器 `-count=2 -v` 四数（票面原句那一把）**＝`142/86/0/0` rc=0：本腿 fresh 两台各一发
   （`v2head`＝导出树 tree.v2B@`c11163c8` pristine、`v2mine`＝tree.v2A@HEAD 含 `322-r2` 那枚 `.go`；
   尺＝逐枚指名单文件 `grep -c '^=== RUN'`⇄顶格三把，`logs/recount.txt`；台面故障尺 0/0/0/0 逐发）
   ＋ r2 入库 `cw2-head/mine`（blob 对撞 MATCH）各复数一字不变＝**同一判语四台发在盘**。量的。
2. **`-run` 名册逐枚点名**＝本格红名册尺（第 2 节，roster-diff＋`roster/*.txt`，量的）＋ SKIP 名册三台逐字相同。
3. **gofumpt／d22scan／长度门／go vet（含 GOOS=darwin）**＝本腿共享树现跑、先 `unset GOTMPDIR TMPDIR`：
   `gofumpt -l internal/winsec/` rc=0 列 0；`-l internal cmd` rc=0 列 4＝票 338 存量假脏族（A1024 形，具名⛔ 归本票）；
   `d22scan.sh` rc=0 clean；长度门 rc=0 VERDICT GREEN；`go vet` winsec×(本机/darwin/linux)＋`./cmd/wisp/` 四发 rc=0。
   全在 `logs/gates-shared.txt`。量的。
4. **消费者 Linux 档（internal/risk 等五包）**＝本腿 fresh `v2cons`＝`1172/714/6/6` rc=1，六枚顶格红＝三枚唯一名×2 样本：
   两枚 `anchor_miss_reason_334_test.go`＝**已另立票 339**（编排者 10:4x 裁定③在册，⛔ 记本票）、
   一枚 `TestC26RewriteAccount…`＝台账 A380／本腿尺证历史 200/216 系件在册（`docs/reports/pending-and-issues.md:8415` 现读）。
   与 r2 `cc2-head/mine`（1160/702/6/6，红名册逐字相同＝`logs/roster/fail-r2.txt`⇄`fail-v2.txt` diff 空）对撞：
   RUN 增量 6 个用例名×2 样本＝227-r1 于 12:21 落定的 config 件（`logs/v2cons-delta.txt`＋`git log -S` 归因 `dd0e489d`）、
   消失名 0、红名册不变 ⇒ **本票收紧没把 Linux 档消费者推红**。量的。
5. **消费者 `cmd/wisp`（Windows 台面）**＝整包成对两台 D⇄C（`448/322/8/3` rc=1 ⇄ `448/323/7/3` rc=1，本腿从原始 log 重数，
   故障尺逐发 0；八枚红逐名归因＝verdict 第 2 节、本票新增 0 枚）＋estate 定向三档（resident `31/22/0/0` rc=0、
   deathpoints `2/2/0/0` rc=0、winlive `12/6/6/0` rc=1 见第 (iv) 条）＋live 首发 `0/0/0/0`＋ntr=2＝台面故障具名作废。量的。
   ★本腿未重跑 `cmd/wisp` 整包（写面腿 `338-r1` 12:22 后 cmd/wisp 树漂移 `54efbbbf→76a62271`＋车道同刻只许一枚）⇒
   整包色对 r2 锚 `1b9f602c` 有效、对现 HEAD 的复跑发＝**待车道**（verdict 第 6 节②）；归因不受影响＝因果尺（下条）＋装置/摆动枚逐名在册。
6. **因果尺（复认⛔ 靠推理）**＝本腿导出树现跑 `go list -e -f`：GOOS=windows 那枚 `.go` 进 `IgnoredGoFiles`(1)⇄`TestGoFiles`(0)，
   GOOS=linux 反之；`cmd/wisp` TestGoFiles 89(windows)/49(linux) 来自 winsec＝0；`cmd/wisp` Imports⇄TestImports 含
   `internal/winsec`＝0/0（前导空格假零已更正复跑）。⇒ `322-r2` 的机制级结论**复认成立**：本票字节机制上进不了 `cmd/wisp`。算的（工具链判定，非推理）。
7. **两台成对那一维够不够（派单③）**＝够：estate 定向 22 枚判定名全部出现在整包 D **与** C 两台顶格 PASS（comm 差 0/0，
   `logs/estate-dC-rulers.txt`）；deathpoints 两枚两台各 1 判定行；resident 族两台 RUN31（顶 22＋子 9）/判定 22/红 0 一字不变
   ⇒ estate 的"C 台一发"不改变任何判语，本格不因此降格（必答(v) 同文）。量的。

## (ii) "本票新增红＝零枚"说的是谁的零枚＋尺

＝**票 322 的字节（`322-r1`/`r1b`/`r2` 的 winsec 面）在 `cmd/wisp` 消费者色上的零枚**。尺（本腿独立跑，`logs/roster-diff.txt`＋`roster/{o,d,c,union-dc}.txt`）：
`UNION(D,C 顶格 FAIL 名) ⊖ 编排者 10-16 名册(orch-322v1-ac4-cmdwisp-full.txt 顶格 FAIL 6 枚)` ＝ `{TestCleanCheckoutBuilds_AC11, TestAC2SealNoticeLandsInTheRunLegLogFile}`；
再 ⊖ 装置枚（CleanCheckout：导出台无 `.git`，红句 `panel_host_gate_test.go:404: git rev-parse --show-toplevel: exit status 128`；
同一枚在 10-16 共享树台面＝`--- PASS`（该件 `:958` 逐字验））⊖ 摆动枚（SealNotice：两台差集唯一一枚 D 红 C 绿＋历史 24 枚件最早 10-13/10-14 在册、
`228/r3b` 两台同红＝`logs/sealnum-files.txt`＋`grep -n` 尺）＝**空集**。枚数相同而名字换一枚的"被盖住新红"形＝本腿查过：
in-O-not-D＝0 枚（旧六枚全在 D 名册里、没被换掉）。⇒ 这句话说得住，且⛔ 一句都⛔ 属"本票把门禁/消费者推红"。

## (iii) "resident_* 在整包从不执行"该收窄成什么形（追加更正，⛔ 抹历史原句）

**建议追加更正文字（归写面所有者落，本腿⛔ 动那些件一字）**：
"resident_*／estate 一族在 `cmd/wisp` 整包从不执行" ⇒ 收窄为
"**只在整进程于真窗区死亡的台面成立**（例：228-r3b surface c/d 两台分别死在 `TestPanelThreadIsSTAAndExitsCleanly`／
`TestAC13BringUpSurvivesAReusedThreadQuit`，本腿尺：两台顶格 Resident RUN＝0/判定＝0）；在跑到终点的台面（10-16 编排者台、
10-18 322-r2 D/C 两台）resident 族照常执行且全绿（顶格 RUN 22＋子 9、判定 22、红 0，三台一字不变）。
⇒ estate 色仍取定向 `-run`（那是防'死亡台面丢读数'的保守口径，⛔ 反过来当'整包不跑它'的证据）。"
**影响读法的已入库件（具名点名，⛔ 本腿改它们）**：
`docs/evidence/s1/322-v1-acceptance.md` 第 6 节（"票 314 死亡区"当整包通则引用）；
`docs/evidence/s1/322-openleaf.md` 第 7 节 AC#4 行（同句形）；
`.scratch/wisp/probes/228/r3b/report.md` 第 2 节（原句出处，读数件＝其腿遗产，更正归编排者）；
本票面 10-16/10-18 两条 bullet 的转述（Progress log 只追加、⛔ 回改）。
`docs/evidence/s1/145-*:207`／`245-*:157`／`75-*:890` 三处"从不执行"＝CI/tag 射程别事，**不受影响**（本腿逐处读过原句）。

## (iv) `-tags winlive` 那一档**不进** AC#4 的数

判语＝**环境存量、⛔ 计入本票**，不是"本格不可判"。凭据：① `12/6/6/0` 六枚红的红句自陈世界占用
（`resident_approval_live_246_windows_test.go:108: … Status:occupied by another app Err:Hot key is already registered`；
`resident_task_source_246…:434 AC#7 RED (desktop state, not our code)`；本腿现量 `wisp.exe` PID 21892 StartTime 2026-10-15 19:03:14，
起于本票所有腿之前，⛔ 杀）；② 同名六枚的历史红在 246/228/258 各家 own 台面早红过（本腿尺：`grep -rl "FAIL: TestLive…"` 命中 6/4/4/3 枚件，
＝别票在册形状）；③ CI 对 winlive 只有编译门（`ci.yml:623`/`:681`＝`go vet -tags winlive`，无 test 步）⇒ 这一档没有"因本票变色"的可比基线，
但**也没有任何一形能把六枚红接到 winsec 面**（①的红句里没有 privateFile/winsec；本腿尺：四枚 winlive 源文件 `grep -c winsec` 全＝0）；
④ AC#4 票面问的是"消费者不推红"，本格数＝默认档＋门禁；winlive 顶格 `=== RUN` 在整包两台＝0（本腿尺），它根本不在整包色里。
**缺的那一维要什么资源才量得到**（具名，⛔ 算腿的欠账）：一台"本机无驻留实例＋可注入鼠标点选托盘/前台焦点"的世界——
即临时停 PID 21892（归机主拍板，⛔ 任何腿动它）或独立非驻留桌面会话；外加 CI 若要把 winlive 变判据需要带窗口的 runner 档（228-v1 第 F-6 格在册）。

## (v) estate 缺 C 台那一发：**不必补**

结论＝整包两台已覆盖"随临时根台面摆动"这一维，estate 的 C 台一发不改变任何判语；台面敏感性另案（票 310/311 族）。
凭据（全部量过，⛔ 推理）：① 包含尺＝定向 resident 22 枚判定名 ⊂ full-D 顶格 PASS ⇄ ⊂ full-C 顶格 PASS（`comm` 差＝0/0）；
② 两台 resident 族 RUN/判定/红＝31/22/0 ⇄ 31/22/0 一字不变、deathpoints 两台各判定行（D=1/C=1×2 枚）；
③ 两台 SKIP 名册逐字相同、且与 10-16 台逐字相同（三台 diff 空）；
④ 会随台面换色的已知形（`TestGateJudgesThePrivateSetByResolvedSID` 一族）住在 `internal/winsec` 本机档、D 红 C 绿已被 r2 成对抓到
（wh/wd 两台），那一族归票 310/311，⛔ 拉进 estate 补发来"覆盖"。
⇒ 若编排者仍想要 estate C 台一发（比如为 228 票的口径统一），那是**另案加料**、不是 AC#4 成立性的必要条件；本腿⛔ 自派腿。

## 附：料级缺口两处（编排者已认，本腿具名点名，⛔ 回头算 322-r2 的账）

① estate 三档只在 D 台跑过、C 台只有整包——本腿裁决见 (v)：维度已被整包两台覆盖，缺的发⛔ 阻断本格。
② `internal/risk` 三枚 Linux 档红＝票 334 仪器在 POSIX 无世界，**已另立票 339**（`.scratch/wisp/issues/339-*.md` 在盘，本腿现量存在），⛔ 记在票 322 头上。
