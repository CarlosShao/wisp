
## 编排者裁定（2026-10-11 13:1x，收 **`308-a1`** 只读普查腿，笔 `c023e8e0`→`dd638f98`→`24929768`，件 `.scratch/wisp/probes/308/a1/**`；台账 `A846`；⛔ 上面各节原句一字不改，本节追加）

★**本轮只翻 `AC#0` 一格**＝现尺 `grep -cE '^[[:space:]]*- \[x\]'`＝**1**／`- [ ]`＝**5**；`AC#1`/`AC#2`/`AC#3` 归非实现者验收（落地腿 `308-r1` 在飞）。裁语 ⓐ–ⓓ 之前只在我派单里说过 ⇒ **本节把它盘上化**（否则下一程只剩"派单里批过"那句⛔ 可核的转述）。

### 1. `AC#0` 我自己复跑的两发＋一处**我自己写错的口径**（具名回自己）

尺＝`go test ./internal/risk/ -count=1 -v -run "TestClassifySpellingInvariance|TestPathResolverUNCAListDenied"`，件＝`probes/orch/20261011/logs/orch-308a1-verify.txt`：

| 发 | env 写法 | 读数 |
|---|---|---|
| A | `TMPDIR=/tmp` | 两枚各 `--- PASS`、rc=0、`panic` 命中 **0**、终态行 `ok github.com/CarlosShao/wisp/internal/risk 0.081s` |
| B | `TMP=/tmp TEMP=/tmp` | 两枚各 `--- PASS`、rc=0、`panic` **0**、`ok … 0.070s` |

⚠★**我那句"两把台面"⛔ 成立**：Windows 档的 Go **⛔ 读 `TMPDIR`**（A 发走的是本机默认 temp），而 `TMP=/tmp` 会被补成 **`D:\tmp\…`**（腿的在场尺逐枚坐实，件 `probes/308/a1/logs/probe-T1-fwd-slash.txt`）⇒ 我那两发**都在"带盘符"那一侧**＝**两写法、一把台面**。⇒ 票面 `AC#0` 要的那发"无卷标台面"，**我和腿都⛔ 交出来**；本格闭合的是这一问的**答案**（哪一枚都不崩、为什么），⛔ 是"我复现了崩"。

产码⛔ 在射程（复认）：`internal/risk/pathresolver.go:262-263` 切 `drive[:1]` 之前逐字有 `if len(drive) == 2 && drive[1] == '$' && drive[0] >= 'a' && drive[0] <= 'z' || …`；`pathresolver_windows.go:61` 取了 `vol` 但⛔ 做 `[:1]`（那文件 `[:1]` 命中 0）。崩点取值源＝`tmpHome` 的 `home := t.TempDir()`。

### 2. ★守卫落点裁语 ⓐ–ⓓ（判语者＝我；写腿＝`308-r1`）

- ⓐ **⛔ 动产码**（射程＝那两枚 `//go:build windows` 测试件）。
- ⓑ **落点＝两枚件共用一枚 helper（把 `vol`→UNC 拼法那一步抽出、`vol == ""` 时具名报错）＋一枚"直接喂无卷标字面量"的单元用例**。★理由（这条是裁语的关键，⛔ 只写在派单里）：机器台面造不出无卷标 temp ⇒ 如果守卫只对 env 形状起作用，它就是一枚**永远不被求值的分支**；只有把 `""` 当**输入**直接喂进去，`AC#2` 的"拔掉守卫⇒回到 panic"才有一发可跑的载体。
- ⓒ **⛔ `t.Skip` 整枚用例**、**⛔ 把 `t.TempDir()` 换成手写常量**（票 72 `:201` 那条禁线同形）；只许**具名去掉 `unc-admin` 那一枚拼法**，且**分母要打印**（同片先例逐字＝`task_state_188_test.go` 那句"a scan whose denominator is zero reports clean for a tree it never read"）。
- ⓓ **两处注释缺陷并入本票写面**（腿复跑 `A843` §2 那把尺报回、我复认）：① `anchor_spelling:157` 指向 `probes/72/r1/20-canonical-form-choice.md`＝**盘上⛔ 此件**（真身＝同目录 `10-class-a-expected-side.md` 的 `## #2` 一节；⚠ 票面 `:77-78` 那两个行号锚今天也已漂到 `:103`，一律改**内容锚**）；② 同文件 `:12-16` 段称 `pathresolver_junction_windows_test.go` "byte-for-byte untouched"，而 `git show --numstat --format= f278c57f -- <那枚>`＝**`36 7`**＝**现时假话**。

### 3. ★历史 live panic 与"本机造不出"这枚矛盾具名登记（⛔ 选一边抹掉）

原始件在盘＝`probes/72/a2/logs/t03-broken-surface-note.txt`（24 行；`rc=1`／名册只 2 枚／`PASS=0`／`FAIL=2`；逐字 `panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`；崩位＝`TestClassifySpellingInvariance`）⇒ **崩发生过**；而同一枚工具台面今天⛔ 能重放（腿 T1/T2/T3/default/bare-root 五把在场尺全非空＋我上面两发）。⇒ 判＝**形状可达已被原始件钉住、⛔ 已"本机可复现"**；`AC#2` 凭据⛔ 用 env 台面去要，按 ⓑ 那枚输入面载体去要。⚠ 那发当时的成因归因（腿写"Go 从 Git Bash 继承无卷标 `TMPDIR`"）我⛔ 复认（Windows 档 Go⛔ 读 `TMPDIR`，见 §1）＝**归因⛔ 立、现象⛔ 否**，留给 `308-v1` 单独裁这一句。

### 4. 欠读数（逐枚具名归谁，⛔ 算在腿账上）

- **两枚崩点的 CI 色**＝只有 push 之后取得到（腿只交"哪一步会打到它们"：`ci.yml:797` 整包 `portable-tests.sh --scope=windows` 会；`ci.yml:813` 那步 `-run TestPathResolverJunctionWindows` **只打 Case 1＝⛔ 覆盖这两枚**）。⇒ 归编排者，与 `A813` 那批同批取；⛔ 用"本机跑过"顶替（`AC#3` 原句）。
- **真无卷标台面那一发**＝本机⛔ 可构造（要它只能换工具/换机器）。⇒ 记为本票**永久欠读数**，`AC#3` 那格的档位读数（`//go:build windows`＋`winlive` 命中 0）腿已交形、判语归 `308-v1`。

**Status（本程更新）：** `AC#0` **已由只读腿交付＋编排者现跑翻勾**（附 §1 那枚我自己口径错具名、§3 那枚矛盾具名）；`AC#1`–`AC#5` ⛔ 动框。⛔ 零 push。

- [2026-10-11 13:1x] agent=编排者 did=收 `308-a1`（三笔 `c023e8e0`/`dd638f98`/`24929768`）⇒ **`AC#0` 翻勾**＋守卫落点裁语 ⓐ–ⓓ **盘上化**（⛔ 动产码／helper＋"直接喂无卷标字面量"那枚用例＝唯一能让守卫有牙的载体／⛔ skip 整枚、⛔ 换手写常量、只许具名去掉 `unc-admin` 且分母必打印／两处注释缺陷并入写面）；我自己现跑两发 `go test` 复认两枚崩点各自 PASS＋rc=0＋panic 0，★并**抓到我自己那句"两把台面"是假的**（Windows 档 Go⛔ 读 `TMPDIR`、`TMP=/tmp` 被补成 `D:\tmp`＝两写法一把台面）；历史 live panic 与本机不可复现那枚矛盾**逐字两边都留**⛔ 抹一边；欠读数两枚具名（CI 色归我同批取／真无卷标台面＝永久欠）；⛔ 动 `AC#1`–`AC#5`、⛔ 改原句、⛔ push
