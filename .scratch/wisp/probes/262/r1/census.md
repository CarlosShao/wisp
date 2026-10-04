# 证据件 262-r1 — 跟踪路径长度门禁（丙＝本机 pre-push 尺 ＋ 丁＝CI 尺）落地读数

**腿**：实现腿 262-r1（写码腿，不是验收者）
**票**：`.scratch/wisp/issues/262-tracked-path-length-gate-for-ci-checkout.md`
**解冻凭据**：`docs/reports/pending-and-issues.md` 的 `A585`（只批 `.github/workflows/ci.yml` 加一枚 step）
**交付 commit**：`f5406c33`（本件骨架）· `519e1ade`（门禁脚本＋CI 那一步）· 第三枚＝本件终态＋脚本的"旧名可追"块
（枚号请现取：`git log --oneline -- .scratch/wisp/probes/262/r1/census.md`；本件不写自己那枚的 hash，写了就立刻过期）

**口径（先定死，全篇通用）**：长度＝**仓内相对字符数**＝`git ls-files` 那一串自身的字符数（含目录前缀、不含检出目录）。
**最坏前缀常量＝44**（＝self-hosted 检出目录 43 字符 ＋ 1 枚分隔符），出处＝票 §现量 1，本腿 §1.1 现跑复认。
**帽＝相对 >121 即红**（规则 9 名 ≤100 ＋ `.scratch/wisp/issues/` 那 21 字符），⛔ 不是墙（墙在相对开区间 (206,217]）。
本件只记读数与判语。每条读数带逐字命令；抄来的数一律标〔未复认〕。

---

## §1 起手读数复认（AC#0：票 §现量 1/2/5/6 四把尺本腿现跑）

### 1.1 §现量 1 — 检出目录两形 ＋ 两道配置

尺与读数（本机 `gh` 可用 ⇒ 本腿自己拉了真日志，不是抄）：

```
$ gh run view 37158259050 --log | grep -a -m4 -E "Working directory is"
slo-full   Run actions/checkout@v4   2026-10-04T00:46:09.3098770Z Working directory is 'E:\work\base\actions-runner\_work\wisp\wisp'
test-windows UNKNOWN STEP            2026-10-03T22:23:12.7300698Z Working directory is 'D:\a\wisp\wisp'
lint       UNKNOWN STEP              2026-10-03T22:23:11.3589386Z Working directory is '/home/runner/work/wisp/wisp'
slo-smoke  UNKNOWN STEP              2026-10-03T22:23:11.7095689Z Working directory is 'D:\a\wisp\wisp'
```
三串长度用文件计数（不让 shell 吞反斜杠）：
```
$ cat > /tmp/prefixes.txt <<'RAWEOF'
E:\work\base\actions-runner\_work\wisp\wisp
D:\a\wisp\wisp
/home/runner/work/wisp/wisp
RAWEOF
$ awk '{print length($0), $0}' /tmp/prefixes.txt
43 E:\work\base\actions-runner\_work\wisp\wisp
14 D:\a\wisp\wisp
27 /home/runner/work/wisp/wisp
```
⇒ **复认成立**：self-hosted＝43 字符（拼路径计 **44**）、托管 Windows＝**14**、托管 Linux＝**27**（票面没写 27，本腿补上，
因为门禁那一步就落在这台机器上，它自己的宽与本票无关——见 §5）。
```
$ git config --get core.longpaths            # rc=1，未设
$ powershell -NoProfile -Command "(Get-ItemProperty 'HKLM:\SYSTEM\CurrentControlSet\Control\FileSystem' -Name LongPathsEnabled).LongPathsEnabled"
0
```
⇒ `LongPathsEnabled=0x0` **复认成立**（262-a1 当时标"未复认"，本腿有注册表读面，现读到 0）。
⚠ **边界同形性仍未证**（本腿也没量）：本机 `git version 2.52.0.windows.1`，runner 侧票面记 `2.54.0` —— 两版对长路径处理是否同形**没量过**，
本件的墙读数是从**runner 自己的日志**取的，不是从本机复现的。

### 1.2 §现量 2 — 墙的位置（本腿独立复跑，不用票面数）

```
$ git ls-tree -r fd269de1^ --name-only | awk '/issues\/25[67]-/ {print length($0), $0}'
217 .scratch/wisp/issues/256-resident-leg-cannot-read-those-risk-config-keys-...-deferred-until-confirming-is-measured.md
219 .scratch/wisp/issues/257-on-a-clean-machine-all-seven-panel-settings-fields-are-unwritable-...-can-create-one.md
$ gh run view 37158259050 --log | grep -a "unable to create file" | sed 's/.*unable to create file //' | cut -c1-30 | sort -u | wc -l
2
```
⇒ 那一次 checkout 里**只有 256、257 两枚路径失败**，同一次检出、同一台机器上其余 5243 枚全部落盘成功；
其中含 260 的旧名（相对 **206**，本腿从 `fd269de1^` 的树里量到过）。换算全路径：
**206＋44＝250 通过／217＋44＝261 失败／219＋44＝263 失败** ⇒ **墙落在相对开区间 (206, 217]**。
⇒ **§现量 2 复认成立**，且和 §1.1 的 43/44 自洽（不需要"42"那个旧数）。
⛔ 本腿**没有**真跑一次 checkout 去把墙钉到点数；要钉死必须真撞一次，那属推送面（见 §7）。

### 1.3 §现量 5 — 改名后的分母

```
$ git -c core.quotepath=false ls-files | awk '{n=length($0); if(n<=121)a++; else if(n<=150)b++; else if(n<=180)c++; else d++; if(n>m){m=n;mf=$0}} END{printf "<=121:%d\n122-150:%d\n151-180:%d\n>180:%d\nmax:%d\ntotal:%d\n",a,b,c,d,m,NR}'
<=121:5142
122-150:38
151-180:19
>180:0
max:180
total:5199
$ git -c core.quotepath=false ls-files | awk 'length($0)>=180'
180 .scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-short-name-stays-short-long-name-folds-long-so-in-root-writes-read-as-out-of-bounds-r2-l2.md
```
⇒ **§现量 5 逐字复认成立**（`5142/38/19/0`、最长 180、全路径 180＋44＝**224**、距最低墙沿 207 还有 27 字符余量）。
（`total:5199` 是 09:1x 起手那一刻的分母；本腿后面每次跑门禁读数都在 5244～5253 之间——三枚在飞的腿此刻正在 commit，
这正说明门禁的量是**现跑**的，不是抄的。）

### 1.4 §现量 6 — 中间带

```
$ git -c core.quotepath=false ls-files | awk 'length($0)>121 && length($0)<=180' | wc -l
57
$ git -c core.quotepath=false ls-files | awk 'length($0)>121' | wc -l
57
```
⇒ **复认成立：中间带＝57 枚**（＝38＋19，与 §1.3 分档对得上）；且因 `>180` 已归零，**"相对 >121"这一档today 就是那 57 枚**
⇒ 名册的正确分母＝57，不是票面推导前的任何别的数。

### 1.5 口径的两面：字符 vs 字节 —— 实测，不是推理

```
$ python -c "（逐枚 len(p) 与 len(p.encode('utf-8')) 对比）"
tracked paths total: 5199
paths where bytes != chars: 28
chars=55 bytes=61 diff=6 .scratch/wisp/probes/171/r1/logs/roster-G5-负一负-post.txt   （最长的一枚非 ASCII 路径）
max chars overall: 180
max bytes overall: 180
over-121 by CHARS: 57
over-121 by BYTES : 57
$ git ls-files | grep -c '^"'      # 需要 git 加引号才能打印的路径枚数
28
```
⇒ **本仓中文名的确会让两口径不等（最多差 2 倍：3 字节/字），但不等的那 28 枚最长只到 61 字节**，
所以**今天的 57 枚超阈者在两口径下是同一批**——这句话现在是**两把尺各跑一遍得到的读数**（门禁脚本每次运行都重跑这两把尺并打印是否同集，见 §2.1 末行），
不是任何人的一次断言。
⇒ 另一枚帽的实测（脚本里第二道检查式的分母）：
```
$ git ls-files -z | awk 'BEGIN{RS="\0"} {n=split($0,c,"/"); if(length(c[n])>100) c++} END{print c}'
57
$ 末段名字最长的那枚
159 .scratch/wisp/issues/252-the-two-spellings-...-r2-l2.md
```
⇒ **末段名 >100 与 相对 >121 今天选出同一批 57 枚**（不存在"名 115 字符但相对 ≤121"的浅目录漏网者，
本腿专门量了：`length(末段)>100 && 相对<=121` 命中 **0 枚**），所以第二道帽不多花一枚名册位。

### 1.6 名册理由那句话是实测的

```
$ 逐枚取首次入库时刻（git log --diff-filter=A --format=%ct -1 -- <path>）中最晚的一枚
2026-10-03 16:38:28  .scratch/wisp/issues/251-winsec-pin-is-a-pin-...-done.md
$ git log -1 --format='%h %ci' fd269de1        # 规则 9 落地那一枚 commit
fd269de1 2026-10-04 08:55:30 +0800
```
⇒ **57 枚全部早于规则 9** ⇒ 名册用同一句类别理由（"filed before the hat"）是**有出处的读数**，不是懒。

---

## §2 门禁读数

### 2.1 第一发：HEAD 上跑＝绿（AC#2）

命令逐字：
```
$ sh scripts/check-path-length-budget.sh
check-path-length-budget.sh: repo=D:/work/workspace/projects plans/Wisp
check-path-length-budget.sh: unit=characters (probe: length of the 5-byte sample "a U+25FF b" = 3; LANG=C.UTF-8 LC_ALL=C.UTF-8)
check-path-length-budget.sh: hat: rule 9 name cap=100 + issues dir prefix=21 -> relative hat=121 ; worst checkout prefix=44 -> full-path budget=165
check-path-length-budget.sh: check: relative_length + 44 > 165, or a last name component > 100
check-path-length-budget.sh: denominator read: 5244 tracked paths
check-path-length-budget.sh: denominator: tracked paths=5244  over-budget=57  covered by roster=57  not in roster=0
check-path-length-budget.sh: longest=180 chars relative (.scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-short-name-stays-short-long-name-folds-long-so-in-root-writes-read-as-out-of-bounds-r2-l2.md)
check-path-length-budget.sh: worst full path on the self-hosted runner=224 chars (hat budget 165, wall open interval (206,217])
check-path-length-budget.sh: bands: over the hat=57  of which in the 122..180 middle=57  past the old debt line(180)=0  in the wall interval=0  roster entries=57
check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree
check-path-length-budget.sh: unit cross-check: the byte view and the character view select the SAME 57 paths on this tree, recomputed just now in both units
$ echo $?
0
```
AC#2 要的四枚数全在：**跟踪路径枚数 5244／最长那枚 180 与其全名／中间带 57／名册 57**；
外加两枚本腿自己要求打印的：**单位＝characters（探测所得）与两口径同集（现跑比对）**。

### 2.2 第二发：种一枚违反规则 9 的跟踪路径＝必须红（AC#3，真树那一侧）

⚠ 不在共享工作树里 `git add`（三枚腿在飞）。本腿用**本机 `git clone --local --no-hardlinks` 的一次性副本**跑，
副本的 HEAD＝`519e1ade`，名册与分母都是真的那一份。命令逐字：
```
$ git clone --local --no-hardlinks --quiet . /tmp/plb-clone-262r1
$ cd /tmp/plb-clone-262r1
$ sh scripts/check-path-length-budget.sh            # 基线，种之前
rc=0
check-path-length-budget.sh: denominator: tracked paths=5245  over-budget=57  covered by roster=57  not in roster=0
check-path-length-budget.sh: VERDICT GREEN - ...
$ PLANT=".scratch/wisp/issues/controlled-<42x y>-<42x y>-real-tree-planted.md"    # 相对 138 字符
$ printf '%s' "$PLANT" | awk '{print "planted relative length:", length($0)}'
planted relative length: 138
$ printf '%s\n' "seeded by the 262-r1 real-tree control" > "$PLANT"
$ git add -- "$PLANT"
$ sh scripts/check-path-length-budget.sh            # 种之后
rc=1
check-path-length-budget.sh: denominator: tracked paths=5246  over-budget=58  covered by roster=57  not in roster=1
check-path-length-budget.sh: RED - over budget and NOT in the roster: .scratch/wisp/issues/controlled-yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy-yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy-real-tree-planted.md  (relative 138 chars, name 117 chars, full path on the self-hosted runner 182 chars, budget 165)
check-path-length-budget.sh: RED - count guard: the roster holds 57 entries, the tree has 58 over-budget tracked paths
check-path-length-budget.sh: VERDICT RED
$ git rm --cached --quiet -- "$PLANT" ; rm -f "$PLANT"
$ sh scripts/check-path-length-budget.sh            # 拆掉之后
rc=0
check-path-length-budget.sh: denominator: tracked paths=5245  over-budget=57  covered by roster=57  not in roster=0
check-path-length-budget.sh: VERDICT GREEN - ...
```
⇒ **AC#3 两发齐**：种＝红且**逐字点名那枚路径与其长度**（138 相对／117 名／182 全路径／帽 165），拆＝绿。
⇒ 种的那枚是**相对 138、名 117**：违反规则 9 但离墙（206+）还远 ⇒ 正是要拦的那一档，不是撞墙才响的那一档。

### 2.3 脚本内建正控（`--self-test`，票 §要建什么 要的"自证"）

```
$ sh scripts/check-path-length-budget.sh --self-test
check-path-length-budget.sh: positive control bench=/tmp/tmp.Mfm5m4eJFt (kept on disk; this project never deletes temp artifacts)
check-path-length-budget.sh: control 1/3 ok - bench baseline green, so this gate is not stuck red
check-path-length-budget.sh: control 2/3 ok - the planted over-budget tracked path is rejected and named verbatim:
RED - over budget and NOT in the roster: .scratch/wisp/issues/seeded-<40x x>-<40x x>-over-the-hat.md  (relative 125 chars, name 104 chars, full path on the self-hosted runner 169 chars, budget 165)
check-path-length-budget.sh: control 3/3 ok - green again once the planted path is gone
check-path-length-budget.sh: positive control PASSED
$ echo $?
0
```
⇒ 一次性仓由 `mktemp -d` ＋ `git init` 造，跑的是**本文件删掉名册块之后的副本**（那枚仓里一枚真名册路径都没有，
带名册去跑会在基线那发就红成"57 枚陈旧"，控制就读不出是哪道守卫响的）。基线/种/拆三态各自断言 **rc 精确等于 0/1/0**
且种那发必须**逐字点名**——只要求"非 0"是假控制，rc=2（拒答）也能蒙过去。

### 2.4 反形控制："换个反形会不会也红"（证明它真在读那枚常量）

⚠ 定式：改前必红不够，还得问"把常量改成另一个值，行为跟着变吗"。两发：

**A 发＝只把阈值常量改成 999（名册照旧）**
```
$ sed "s/^HAT_NAME_LIMIT=100$/HAT_NAME_LIMIT=999/" scripts/check-path-length-budget.sh > /tmp/c999.sh
$ grep -c "^HAT_NAME_LIMIT=999$" /tmp/c999.sh
1
$ sh /tmp/c999.sh ; echo rc=$?
check-path-length-budget.sh: hat: rule 9 name cap=999 + issues dir prefix=21 -> relative hat=1020 ; worst checkout prefix=44 -> full-path budget=1064
check-path-length-budget.sh: denominator: tracked paths=5245  over-budget=0  covered by roster=0  not in roster=0
check-path-length-budget.sh: bands: over the hat=0  of which in the 122..180 middle=0  past the old debt line(180)=0  in the wall interval=0  roster entries=57
check-path-length-budget.sh: VERDICT RED
rc=1
$ grep -c "RED - over budget and NOT in the roster" /tmp/out999.txt
0
$ grep -c "RED - roster entry that is not an over-budget" /tmp/out999.txt
57
```
⇒ **没有任何一枚文件再被点名成"超预算"**（那一节 0 行），`over-budget` 从 57 掉到 **0** ⇒ 那枚常量**真的在被读**，
这把尺不是恒红也不是恒绿。整体仍 rc=1，但**红在别的牙齿上**：名册 57 枚此刻全都成了"树上已不存在的豁免"（守卫 B）＋枚数不等（守卫 C）。
这一发因此同时是 B/C 两道守卫的正控——本腿**没有**为了让它看起来干净去把名册也改小；那条读数在下面 B 发。

**B 发＝常量 999 ＋ 名册块整块删掉（把 B/C 两道守卫的输入也归零）**
```
$ sed "s/^HAT_NAME_LIMIT=100$/HAT_NAME_LIMIT=999/" scripts/check-path-length-budget.sh | sed "/^#ROSTER-AWK-BEGIN\$/,/^#ROSTER-AWK-END\$/d" > /tmp/c999nr.sh
$ sh /tmp/c999nr.sh ; echo rc=$?
check-path-length-budget.sh: denominator: tracked paths=5245  over-budget=0  covered by roster=0  not in roster=0
check-path-length-budget.sh: bands: over the hat=0 ... roster entries=0
check-path-length-budget.sh: VERDICT GREEN - ...
rc=0
```
⇒ 阈值常量单独一提一放，**超阈枚数跟着 57→0**，其余输入不变时整把尺回到绿 ⇒ 常量与判决之间的因果是本腿跑出来的，不是叙述的。

### 2.5 AC#4 两枚抽查：名册的理由必须是被解析的数据

```
C 发＝把第一枚名册的理由清空（`... = ""`），路径行照旧：
$ awk '...第一个匹配 /^    R\[".*"\] = "hist/ 的行，把值替换成 ""...' scripts/check-path-length-budget.sh > /tmp/cnoreason.sh
$ sh /tmp/cnoreason.sh ; echo rc=$?
check-path-length-budget.sh: RED - roster entry carries an EMPTY reason; the reason is data, not a comment: .scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-short-name-stays-short-long-name-folds-long-so-in-root-writes-read-as-out-of-bounds-r2-l2.md
check-path-length-budget.sh: VERDICT RED
rc=1        （红行枚数=1，即只有这一道响：枚数守卫没响，因为名册还是 57 枚）

D 发＝把那一整枚名册行删掉（路径还在树上超阈）：
$ sh /tmp/cnodrop.sh ; echo rc=$?
check-path-length-budget.sh: RED - over budget and NOT in the roster: .scratch/wisp/issues/252-the-two-spellings-...-r2-l2.md  (relative 180 chars, name 159 chars, full path on the self-hosted runner 224 chars, budget 165)
check-path-length-budget.sh: RED - count guard: the roster holds 56 entries, the tree has 57 over-budget tracked paths
check-path-length-budget.sh: VERDICT RED
rc=1
```
⇒ **删理由⇒门变；删名册行⇒门换个说法再变**，两发不同守卫、不同 rc 组成，都不是"注释里写写而已"。

---

## §3 名册、守卫与旧名可追（AC#7）

### 3.1 名册落点与构成

- 尺（现跑，本腿自己取的分母，未照抄票面的 57）：`git -c core.quotepath=false ls-files | awk 'length($0)>121 || 末段名>100'` ⇒ **57 枚**。
- 名册在 `scripts/check-path-length-budget.sh:314-371`（`#ROSTER-AWK-BEGIN` 与 `#ROSTER-AWK-END` 两枚标记之间，57 行赋值）。
- 构成（按目录与结案态）：**34 枚未结案工单 ＋ 20 枚已 `-done` 工单 ＋ 3 枚 `.scratch/wisp/dispatches/` 派单件 ＝ 57**。
  ```
  $ awk '{p=$2; if (p ~ /^\.scratch\/wisp\/dispatches\//) d++; else if (p ~ /-done\.md$/) c++; else o++} END{printf "open=%d closed=%d dispatch=%d total=%d\n",o,c,d,o+c+d}' /tmp/offenders.txt
  open tickets=34  closed(-done) tickets=20  dispatches=3  total=57
  ```
  三句类别理由与名册行数逐一对得上：`grep -c` 分别得 **34 / 20 / 3**。
- 为什么派单件也在册：规则 9 只管工单名（`README.md:57`），但 `.scratch/wisp/dispatches/` 前缀是 25 字符，
  那 3 枚的**末段名 102/106/111 字符本来就 >100** ⇒ 第二道帽抓到它们，名册必须点名，⛔ 不许"既不改名也不点名"。

### 3.2 四道守卫（本腿各跑过一发，见 §2）

| 守卫 | 语义 | 响法 | 本腿正控 |
|---|---|---|---|
| A | 超阈而不在名册 | 逐枚点名＋相对/名/全路径三个长度 | §2.2 种 138 ⇒ 点名；§2.5 D 发 ⇒ 点名 |
| B | 名册里那枚已不是树上超阈者（名册落后/路径被改名） | 逐枚点名 | §2.4 A 发 ⇒ 57 枚全部点名 |
| C | `枚(名册) != 枚(树上超阈)` | 打印两个数 | §2.2（57 vs 58）／§2.5 D（56 vs 57）／§2.4 A（57 vs 0） |
| D | 名册理由为空 | 点名那枚路径 | §2.5 C 发 |

⚠ C 是 A/B 的算术冗余，**故意两枚都留**：票面点名的就是"枚数不等⇒红"那枚一行算式，而 A/B 负责说清"哪一枚、往哪边差"。

### 3.3 旧名可追（AC#7）——11 枚逐枚在案

11 枚「曾名 …，现名 …」逐枚写在 **`scripts/check-path-length-budget.sh:178-213`**（票面 AC#7 指定的两处落点之一"脚本名册"），
每枚两行：`曾名` 全路径 ＋ `-> 现名` 全路径；2 枚肇事者出自 commit `fd269de1`、9 枚出自 commit `46079fcc`，
复核命令也写在那一段注释里（`git show --name-status -M --format= fd269de1`／`46079fcc`）。
本件不重复那 22 行长路径；两枚肇事者的相对长度实测（§1.2）＝**217／219**，改名后＝**60／73**（本腿从现名量的读数见下）。

### 3.4 过期指认的代价——三把尺，票面那个"11 行/10 份"本腿复认**不成**

```
尺一（逐字整串命中：把 11 枚旧名原样在 docs/evidence/s1/** 里找）
  exact-full-old-name citations: 3 occurrences across 3 files
    docs/evidence/s1/176-background-start-port-census-a1b.md
    docs/evidence/s1/177-c25-r4-path-exemption-c1.md
    docs/evidence/s1/177-mutation-shape-a-redlist-m1.md
尺二（前缀命中：允许旧名被缩写成 `...-trigger-gate.md` 这种半截）
  9 行 / 8 份（149 两份、154 两份含一行缩写、176、177 两份、183）
尺三（泛化：引用的 issues 路径不在当前 ls-files 里 ⇒ 判它过期）
  220 行 / 142 份 —— ★这把尺本腿判它**不适用**：它把 `issues/144-….md`、`issues/146-*.md` 这种**刻意的省略/通配写法**
  也算成过期指认（在 147-offset-naming-r1.md 里逐字看到三例），所以那 220 不是债的枚数，只说明"这形尺子不能用"。
```
⇒ 与票面 AC#7 那句"共有 **11 行 / 10 份**"对不上；与 262-a1 那句"仅 3 行"在**尺一**下同成立。
⇒ 本腿**不采认 11/10 这个数**，也不新造一枚权威数：差别在"缩写算不算过期指认"这一步（尺一＝3／尺二＝9），
请验收腿先定这把尺再定这格。⚠ 另：本腿**没有**去改 `docs/evidence/s1/**` 与 `.scratch/wisp/probes/**` 的任何历史引文（票面禁止，且改了读数就不是读数了）。

---

## §4 落点与接线（AC#1：调用点枚数 ≥1，具名 file:line）

| 物 | 落点（全路径:行） |
|---|---|
| 门禁脚本 | `scripts/check-path-length-budget.sh`（577 行／44,926 字节，文件自身名 35 字符 ⇒ 守规则 9） |
| **调用点①＝CI（丁）** | `.github/workflows/ci.yml:136`（step 声明 `- name: "Tracked path-length budget (ticket 262)"`）＋ **`.github/workflows/ci.yml:166`（命令行 `run: sh scripts/check-path-length-budget.sh --with-self-test`）** |
| **调用点②＝本机 pre-push（丙）** | `scripts/check-path-length-budget.sh:25-33` 具名"谁调它、什么时候调"：编排者（唯一持推送权的一方）每次 `git push` 前跑同一条命令 |
| 为何不是装 hook | `scripts/check-path-length-budget.sh:29-33`：本仓无 `.githooks/`、`git config --get core.hooksPath` rc=1（本腿现跑），装 hook 属 A585 未批的 git-config 变更 |
| 检查式 | `scripts/check-path-length-budget.sh:53`（`red iff relative_length + WORST_PREFIX > FULL_PATH_BUDGET`） |
| 阈值常量 | `scripts/check-path-length-budget.sh:149`（`HAT_NAME_LIMIT=100`，唯一旋钮）· `:156`（`WORST_PREFIX=44`，出处写死在紧邻注释里） |
| 扫描函数 | `scripts/check-path-length-budget.sh:285` |
| 内建正控 | `scripts/check-path-length-budget.sh:472` |
| 名册 | `scripts/check-path-length-budget.sh:314-371` |

⇒ AC#1 的"≥1 枚具名调用点"由 `.github/workflows/ci.yml:166` 满足；第二枚（本机）是文档化调用点，
本腿**没有**把它伪装成装好的钩子（盘上确实没有 hook，见 §6.5）。

---

## §5 `.github/workflows/ci.yml` 改动边界自证（A585 只批"加一步"）

```
$ git diff --numstat -- .github/workflows/ci.yml
32      0       .github/workflows/ci.yml
$ git diff -U0 -- .github/workflows/ci.yml | grep '^@@'
@@ -135,0 +136,32 @@ jobs:
```
⇒ **32 增／0 删，单一 hunk，插入点在 135 与 136 之间**＝`run: sh scripts/d22scan.sh`（原 `:134`）之后、
`- name: gofmt (gofumpt)`（原 `:136`，现 `:168`）之前——正是 A585 建议的位置。hunk 头 `-135,0` 即"这一侧没删过一行"。

未动的部分（本腿逐条核过）：
- 触发表 `:12`（`pull_request`）／`:36`（`schedule`）——在 135 之前，**行号与内容一字未变**（纯插入 ⇒ 前 135 行逐字节相同）。
- 并发组 `:50-51` 同上，未变。
- `slo-full` / `test-core` / `test-windows` / `slo-smoke` 任何现有 step：未动（hunk 只有一枚，在 `lint` job 内）。
- 任何 job 的 `runs-on`：未改。YAML 解析后的现值：
  ```
  jobs: ['lint','test-core','test-windows','slo-smoke','slo-full','lint-frontend']
  runs-on: lint=ubuntu-latest, test-core=ubuntu-latest, test-windows=windows-latest,
           slo-smoke=windows-latest, slo-full=['self-hosted','wisp-slo'], lint-frontend=ubuntu-latest
  lint steps: 12（改前 11 ⇒ ＋1）
  新 step 的键：['name','run']        # 无 if:、无 continue-on-error（D22 mode 6）
  ```
- `python -c "import yaml; yaml.safe_load(...)"` 通过 ⇒ 缩进/引号合法。

**这枚 step 照不到哪一侧的宽（具名，免得下一任读成"真机已复现"）**：它跑在 ubuntu-latest，
本机检出目录 27 字符、托管 Windows 14 字符，都 **不等于** self-hosted 那 43/44；
所以它的价值是**纸上最坏预算**（`相对 + 44 > 165`），不是那次真 checkout 的重放。
真正的 self-hosted 复现只在推送之后由 `slo-full` 自己给（见 §7）。

---

## §6 门禁读数与判不动的地方

### 6.1 本腿读到的（每条都有上面写过的尺）

- 两形检出宽 43/44 与 14（＋本腿补的 Linux 27）；`LongPathsEnabled=0`；`core.longpaths` 未设（rc=1）。
- 墙＝相对开区间 **(206, 217]**，换算全路径 250 通过／261 失败／263 失败；那次检出只有 2 枚路径失败。
- 改名后分母 5142／38／19／0，最长 180（全路径 224，距最低墙沿 27）。
- 中间带 **57 枚**＝名册分母；`>121` 与 `末段名>100` 两枚帽今天同集（浅目录漏网者实测 0 枚）。
- 字符/字节两口径在超阈带上同集（非 ASCII 那 28 枚最长 61 字节）；**这个同集是每跑必重测的读数**，不是断言。
- 门禁三发：HEAD 绿／种 138 红并逐字点名／常量 999 后**超阈枚数 0**；内建正控三态 rc 精确 0/1/0。
- 名册 34＋20＋3＝57，四道守卫 A/B/C/D 各有一发独立正控。
- 11 枚改名对逐枚在 `scripts/check-path-length-budget.sh:178-213` 可追。

### 6.2 判不动的地方（具名，不替自己圆）

1. **AC#5"CI 侧真出颜色"本腿做不了**：它要一次推送或 run 重跑后的**日志行**，而推送权在编排者手里（`issues/README` 规则 1/2），
   本腿只 commit。留给下一任的凭据见 §7.1。
2. **墙不能钉成点数**：要钉死必须真跑一次跨墙 checkout（落盘＋网络），本腿只从日志换算，
   所以 §1.2 给的是**区间**不是值。⛔ 本件没把 217 当已知墙高用（检查式里也没出现 217 参与判决）。
3. **本机 git 2.52.0 与 runner 2.54.0 是否同形**：没量。⇒ 本机的"绿"不等于 runner 的"绿"，本腿没有把前者说成后者。
4. **mawk 那一路本腿没跑过**：本机只有 gawk（`command -v mawk` → 不存在），字符探测在 CI 上会走 `unit=bytes` 那一支。
   该支的行为是**推的**（超阈者全为纯 ASCII ⇒ 两口径同集，见 §1.5），不是量过的。
   ⚠ 若 CI 上那一步报 `REFUSE - unit=bytes and N offender(s) carry non-ASCII bytes`，那是**设计中的拒答**，不是 bug，
   修法是把 gawk 装进 job 或让那枚名字缩短，⛔ 不是抬帽。**这是本件最大的未实测面，具名留给 262-v1。**
5. **丙那一枚调用点是"文档化"的**：`scripts/check-path-length-budget.sh:25-33` 写了谁在什么时候跑，但盘上没有强制机制
   （无 hook、无 `core.hooksPath`，且 A585 不批）。它响不响取决于编排者跑不跑——本腿没把这句写成"已接线"。
6. **门禁读的是索引（`git ls-files`），不是工作树**：此刻三枚腿在飞，同一枚 commit 上的分母在 5244→5253 之间漂，
   而 CI 检出的是**提交后的树**（§2.2 的 clone 读数 5245 就是那一款）。⇒ 本机读数与工作树脏度相关，这不是 bug，
   但任何"本腿绿了"的话都必须带上"绿在哪一枚树上"。
7. **相对帽对浅目录仍比规则 9 松**：第二道帽（末段名 >100）已把这一处补上；但规则 9 的词面只管工单名，
   所以"非工单的一枚 105 字符新文件名会不会被别的规矩要求短一点"——本腿没有那份规矩，判不动，没自己立。
8. **`docs/evidence/s1/**` 那 11/10 的代价读数复认不成**（§3.4）：三把尺各得 3／9／（不适用的）220，
   不知道票面那把是哪一枚 ⇒ 这格的数请验收腿定尺。
9. **卫生四门本腿一枚都没跑**（AC#6）：指令明令"⛔ 别跑任何 `go test`／`go build`"，而 `scripts/d22scan.sh` 第一步就是
   `runtests.sh`（＝go test）、`gofumpt` 也要装二进制。⇒ 见 §7.2 的替代自证与留件。

---

## §7 本腿没做成的（具名）与留给下一任的凭据

### 7.1 AC#5：今天只欠一次推送

- 没做成的：**"该步在 CI 里真跑过一次并留下日志行"**。原因＝本腿无推送权（规则 1/2），且不得 push。
- 留下的凭据（下一任可逐字核）：
  1. 该步的唯一命令文本：`.github/workflows/ci.yml:166` ＝ `sh scripts/check-path-length-budget.sh --with-self-test`；
  2. 本机同一命令的读数已在本件 §2.3＋§2.1（绿、rc=0，正控 3/3）；
  3. 推送后核验只需一条：`gh run view <new-run-id> --log | grep -a "Tracked path-length budget"`，
     期望看到 §2.1 那 10 行里的 `unit=`／`hat=`／`denominator:`／`VERDICT GREEN`，**且** `control 1/3..3/3 ok`；
  4. 期望的**红形**（同一命令在 CI 上的失败形状）＝§2.2 那两行 `RED - over budget and NOT in the roster: …` ＋ rc=1；
  5. ⚠ 若 CI 读到 `unit=bytes`（mawk 那一路），`denominator` 与 `VERDICT` 两行仍应逐字出现，
     只有含非 ASCII 的超阈者才会让它 `REFUSE`（§6.2 第 4 条）。
- 另欠：**`slo-full` 那一侧仍未被真复现过**。本尺是纸上预算；跨墙那枚路径今天依然会让 `slo-full` 死在 checkout，
  只是现在会**先**在 `lint` 死一次（更早、在任何机器上）。把"更早"说成"消除"就是假话，本件没说。

### 7.2 AC#6：卫生四门本腿一枚没跑

- 没做成的：`gofumpt`／`go vet`／`tools/d22scan`／票 212 的 ban #9 的**改后读数**。原因＝指令硬禁本测量期跑任何 go 命令
  （`tools/d22scan/**`、`internal/agent/**`、`internal/tools/**`、`cmd/wisp/**` 三枚腿在飞）。
- 本腿能给的替代自证（不是四门读数，别混用）：
  ```
  $ git show --name-only --format= 519e1ade
  .github/workflows/ci.yml
  scripts/check-path-length-budget.sh
  ```
  ⇒ 两枚文件**一枚 `.go` 都没有**；`tools/d22scan` 的 emoji 射程按 `ci.yml:56-62` 的自述是 `design/`＋`internal/`＋`cmd/` 的 `.go`，
  ban #9（票 212）扫的也是 Go 注释/形状 ⇒ 四门的**分母枚数不因本 commit 改变**这一句是**由文件清单推的**，
  ⛔ 不是"跑过且没扩大"的读数。**红名集合逐名比对＝本腿没做，具名交给 262-v1 现跑**：
  `sh scripts/d22scan.sh`、`gofumpt -l` 两形、`go vet ./...`，与本 commit 前的基线对表。
- 本件对 ban #9 的**形状自查**（本腿能做的部分）：脚本与注释里引仓内路径一律写全路径
  （`.scratch/wisp/issues/README.md`、`docs/reports/pending-and-issues.md`、`scripts/d22scan.sh`、
  `.scratch/wisp/probes/262/a1/census.md`、`.scratch/wisp/issues/262-...md`），
  且**没有指名任何一枚不存在的测试或脚本当凭据**——被引用的 `scripts/d22scan.sh`、
  `tools/d22scan/runtests.sh`、`.github/workflows/ci.yml`、`docs/reports/pending-and-issues.md` 全部盘上有物（本腿 ls/grep 过）。

### 7.3 其余没做成的（逐枚）

1. **AC#3 没在共享工作树里种路径**：只用了一次性 clone（§2.2）。理由＝三枚腿在飞，`git add` 一枚我不 commit 的路径
   会进别人的 pathspec 视野。代价＝那两发读数来自 `519e1ade` 的克隆，不是此刻的工作树（两枚分母都在 §2 标着）。
2. **没把 57 枚中间带改名**：票面给的是"要么一次改完，要么逐枚进名册各写一句理由"，本腿选了名册——
   改名不是本票的授权面（A584 那句"遗留 9 枚改短名"是编排者做的，且 57 枚里有 34 枚是**活票**，
   动活票名＝动防重领键的形状），⛔ 本腿没自己扩权。**后果具名**：这 57 枚今天合法、但**墙还在**——
   最长那枚距最低墙沿只剩 27 字符，`-done` 结案会再吃 5 字符；**门禁拦得住新的一员，拦不住已在册那 3 枚活雷**。
3. **`--self-test` 的一次性仓留在盘上没删**（`/tmp/tmp.Mfm5m4eJFt` 等多枚 ＋ `/tmp/plb-clone-262r1`）：
   规则 8"临时件只建不删"，本腿照做，路径全在正文里可回看。
4. **票面 AC 框一枚未碰**：`- [ ]` 七格原样，本腿只在 Progress log 追加了一行做了什么。
5. **`.gitignore`、`docs/PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、golden、
   `tools/d22scan/allowlist.txt`、`frontend/**`、`design/**`、三枚冻结件、`tools/d22scan/**`、`internal/agent/**`、
   `internal/tools/**`、`cmd/wisp/**`——一枚未动**（`git show --name-only 519e1ade` 即全集）。

### 7.4 交件自证（盘上读数）

```
$ git log --oneline -3
（本腿三枚：f5406c33＝本件骨架 · 519e1ade＝门禁脚本＋CI 那一步 · 第三枚＝本件终态＋脚本的"旧名可追"块）
逐枚可核：git log --oneline -- .scratch/wisp/probes/262/r1/census.md scripts/check-path-length-budget.sh .github/workflows/ci.yml
$ git show --name-only --format= 519e1ade
.github/workflows/ci.yml
scripts/check-path-length-budget.sh
$ wc -c scripts/check-path-length-budget.sh
44926
```
⇒ 本件**不写自己的字节数**（自引用一改就过期），那一枚读数在 §7.5 末与交件报告里给。
占位符自证见 §7.5（含"这把尺为什么会量到自己"那一处本腿真踩过的坑）。

### 7.5 占位符 grep（终态自证）

**尺的形状先说清，因为它决定这条自证能不能读**：本件**故意不把要扫的那六枚词形原样写出来**——写出来这把尺就会量到自己。
本腿在这一段上**踩了两次**同一处：第一版把六枚词形原样列在命令里 ⇒ 尺对本件得 **3**（全是命令自己）；
改成解释性散文后又得 **1＋1＋1**，因为散文里用了其中三枚中文词形。下面这版把词形拆开写，尺才只量目标。
⇒ 交付时的逐字命令本件不写（写了就污染读数）；验收腿要重跑，尺取六枚词形竖线相接：
英文那枚"尚未完成"形（T‑O‑D‑O 连写）／英文那枚"待定"形（T‑B‑D 连写）／中文"待‑填"那枚／中文"待‑补"那枚／
英文那枚"占位符"名（place‑holder 连写）／以及脚本里曾用以占位的"调用点锚点串"那枚。

本腿现跑（把词形用相邻引号拆开写，尺不被自己污染）：

| 词形（拆开写，故本件不含其原样） | 脚本命中 | 本件命中 | 命中是什么 |
|---|---|---|---|
| TO+DO | 0 | 0 | — |
| T+BD | 0 | 0 | — |
| 待+填 | 0 | 0 | — |
| 待+补 | 0 | 0 | — |
| place+holder | 0 | 0 | — |
| CALL+‑SITE | **2** | **2** | 两枚都**不是占位符**：见下一段那对历史工单名 |

⇒ 那 2 枚（连写的 CALL‑SITE）命中**都不是占位符**：`grep -c` 数的是**行**不是出现次数，
本件这 2 行＝§3.3"旧名可追"那对历史工单名在本节被引用的两行（本件 `:490-491`：
`…-production-call-sites-and-taskroster-…` 与它的现名 `…-no-spawner-zero-production-call-sites.md`），
脚本那 2 行＝同一对在 `scripts/check-path-length-budget.sh:200-201`。
逐字复核：`grep -n -i 'call.site' scripts/check-path-length-budget.sh` → 命中 5 行，
其中 3 行是散文里的 "call site(s)"（空格形，不计入上表），2 行是那对历史路径。
⇒ **判：脚本与本件都不含占位符，也不含未填的表行**（表里那些 0 是读数，不是缺项）。
⚠ `grep -c` 零命中返回 rc=1，所以这一格的"没命中"＝好消息那一侧，别当失败读。

文件大小（稳定读数，本件落定后不再改脚本）：
```
$ wc -c scripts/check-path-length-budget.sh
44926
$ wc -l scripts/check-path-length-budget.sh
577
```
⚠ 本件**不在自己身体里写自己的字节数**：那是一改就变的自引用读数。本件的 `wc -c`／`wc -l` 在回给编排者的交件报告里给，
文件正文只报脚本那一枚（§7.4）——它稳定，因为这枚 commit 之后 `scripts/**` 不再动。
