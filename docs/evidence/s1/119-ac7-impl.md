# 119 AC#7 —— 实现侧交件：那两枚未解析根用例接了记名断言，并同时换到已解析基根

**格子**＝票 `.scratch/wisp/issues/119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md`
的 **AC#7 一格**（票面 `:61-80`）。
**本件不是裁决表**：AC#7 的勾**没翻**，本票文件名**没改成 `-done`**，终裁归非实现者
（`docs/evidence/s1/119-ac7-*.md` 由另一枚 agent 出）。本件里不写"成立／通过"这类判语。
**程名**＝`worker-ticket119-ac7`（实现者）。

## §0 锚点自量（全部现跑，不抄派单里的号）

本机 `dev`，2026-09-24 17:12 +0800（`date` 原文见下）：

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-24 17:12 +0800
$ git rev-parse --abbrev-ref HEAD
dev
```

| 号 | 是什么 | 怎么量的 |
|---|---|---|
| `a9c8b6e6dce4ad63a6fc04b15385651a1443ea65` | **改前那棵树**＝派单开工时 `git rev-parse HEAD` 的现量值 | 派单第一枚命令现量；随后 `git log --oneline a9c8b6e..HEAD -- internal/winsec/` = **空**（到本件落盘为止 `internal/winsec/` 无邻居改动）⇒ 基线树与交件树只差本文件那一枚 |
| `ddc1583` | 码第一枚 commit（两味药） | `git log --oneline -1` 当场回显 |
| `c94927d` | 码第二枚 commit（两句红句点名本用例种的链接）；**本件所有改后读数的树** | `git log --oneline -1` 当场回显 |
| `28c997c` | 派单开工后落在我两枚 commit 之上的邻居 commit（编排者的台账） | `git log --oneline -3` |

⚠ 派单简报里给的锚点写作 `4e66817`（AGENTS.md 的生成锚）与"`b1ea719`"（我第一条命令量到的 HEAD）；
**盘上现量到的是 `a9c8b6e`／`866c849`**（同一天邻居在落盘），本件一律用我自己现量的那组号，按简报要求写明差异。

工作树**不干净且与我的活无关**：`design/**` 16 枚 tracked 文件被 owner 那侧挪成未跟踪的 `design/old/`
（另有 `design/doubao/`）。本程**一枚未还原、一枚未提交、一枚未删**；`git add` 只用显式路径，
每枚 commit 前跑过 `git diff --cached --name-only`（原文在 §5）。

一处**推送事实登记**（不是本程做的，本程只 commit 未 push）：`git log` 现量到编排者的
`28c997c docs(A173): 如实记——我推自己那枚 A172 时把 #50 的两枚在飞中间态 commit（ddc1583/c94927d）
一起发布到两远程了` ⇒ 我那两枚中间态 commit 已被那次推送带走。本程未发过任何 `git push`。

## §1 判据①：两味药都下了，落在哪几行

改后文件 `internal/winsec/dataroot_symlink_119_other_test.go`（**只此一枚文件**）：

```
$ grep -n "errors.Is(err, winsec.ErrUnresolvedPath)\|refusalCreditsLink137(err, shape.link)\|base := cleanSpelling119" internal/winsec/dataroot_symlink_119_other_test.go
225:	base := cleanSpelling119(t, t.TempDir())            （AC#1 母项：换已解析基根）
241:	} else if !errors.Is(err, winsec.ErrUnresolvedPath)  （原 `:203` 那一处，条件未动）
243:	} else if !refusalCreditsLink137(err, shape.link)     （新：记名断言）
291:	} else if !errors.Is(err, winsec.ErrUnresolvedPath)  （`:251` 那一处：行号漂到 291，内容一字未动，见 §4）
337:	base := cleanSpelling119(t, t.TempDir())            （AC#2 leg 1：换已解析基根）
361:	} else if !errors.Is(err, winsec.ErrUnresolvedPath)  （原 `:308` 那一处，条件未动）
363:	} else if !refusalCreditsLink137(err, shape.link)     （新：记名断言）
```

行号漂移是**我自己现量的**（派单给的 `:203/:251/:308` 是改前号；开工前 `grep -n` 复量过一次，
改完后又量过一次，就是上面这组）。

- 记名那味用的 helper＝同包既有的 `refusalCreditsLink137`，签名
  `func refusalCreditsLink137(err error, planted string) bool`，定义在
  `internal/winsec/placement_symlink_113_other_test.go:160`（**改前就有**，本程零新增依赖、零新函数）。
  可达性核过：两枚文件同为 `//go:build !windows` 的 `package winsec_test`，同目录同测试包。
- 换根那味走的是本文件自己的先例 `cleanSpelling119`（`filepath.EvalSymlinks`，`:112` 那枚函数），
  **不是** `proc.SealableRoot`——理由与票 119 Rules `:94`（`R-119-9`"不许拿被测函数算 fixture"）
  以及 `119-ac2-r2b-acceptance.md` §3-A 末尾那句"真要落地该跟同包 `:251` 那枚的先例走 `cleanSpelling119`"
  一致。`winsec.SealableTempDirForTest124`（137 那批用的那枚）也可达，本程没用它，理由：
  本文件已有同义的 `cleanSpelling119`，多引一枚跨文件 helper 只会让"哪味药换的根"更难对上。
- AC#1 母项另加**一句前提门**（`PrivateDirAll(<解析后基根>/premise)` 必须为 nil，否则 `t.Fatalf`）：
  它钉的是"这棵树上除了我种的链接不许再有别的链接"，用的判定方向与本文件 `:227` 那枚 AC#3 的同形先例一致；
  AC#2 那枚没另加，因为它的 leg 2（同一棵树下拼写干净的根必须封成功）已经是同一句话的门。
- **未动的东西**：两枚用例的判定条件方向（`err == nil` / `!errors.Is(err, ErrUnresolvedPath)`）一字未动，
  只在后面**追加**了 `else if !refusalCreditsLink137(...)` 一支；`c94927d` 那枚 commit 只改两句红句的
  **措辞**（把本用例种的链接名打印进红句），条件一字未动。全文 7 枚删除行逐行如下（`git diff -U0` 原文，已截去路径前缀）：

```
-// Two instrument notes, because both have bitten this repository:
-	base := t.TempDir()
-		t.Errorf("AC#1 RED: an unresolved data root through a symlink was accepted, so the leg is gone")
-		t.Errorf("AC#1: refusal did not name ErrUnresolvedPath: %v", err)
-	base := t.TempDir()
-	// creating anything in the tree the link names.
-		t.Errorf("AC#2: refusal did not name ErrUnresolvedPath: %v", err)
```

即：**两枚换根、三枚换措辞、两枚注释**。没有任何断言被删或被放宽。

`git diff --numstat a9c8b6e c94927d -- internal/winsec/dataroot_symlink_119_other_test.go` ⇒
**62 增 / 7 删**（同一枚文件，本程在 `internal/winsec/` 只动这一枚）。

## §2 仪器（读数全部出自 Linux 容器，Windows 上那两枚 `!windows` 用例没有分母）

- `docker version --format '{{.Server.Version}}'` ⇒ `29.6.2`；镜像 `golang:1.27` **本机已有**（`docker images`
  现量 `golang:1.27 3680233e3204`），**未 `docker pull`**；容器内 `go version go1.27.1 linux/amd64`，uid 0。
- 挂载**一律 Windows 风格**：`-v "D:\tmp\wisp119-ac7\snap-xxx:/wisp"`（外加
  `-v "D:\tmp\wisp119-ac7\scripts:/scripts:ro" -v "D:\tmp\wisp119-ac7\logs:/logs"` 与两枚卷
  `wisp119ac7-gomod`／`wisp119ac7-gobuild`）。
- **挂载非空自证**：每发容器第一句 `ls -l /wisp/go.mod` 进该发的 `head.txt`，原文
  `-rwxrwxrwx 1 root root 883 … /wisp/go.mod` ＋ `883 /wisp/go.mod` ＋
  `f6ef661732b1851e5c3db348113cb605  /wisp/go.mod`；容器内另有硬门（`go.mod` 为空 ⇒ `GATE97` 直接 exit、**不取颜色**）。
- ⚠ **一次仪器自拒（登记，不抹）**：改后四发第一次跑时我把快照路径写成**相对路径**
  （`sh scripts/drive.sh b1-... snap-after plain full`），`cygpath` 给出相对串 ⇒ 容器里 `/wisp` 挂成空目录，
  `GATE95` 在取颜色之前 `exit 95`，**那一遍零枚颜色入账**；改用绝对路径 `-v "D:\tmp\...\snap-after:/wisp"` 重跑才是本件的数。
  简报警告的"`-v /d/tmp/...:/wisp` 静默挂空"这一形我也顺手加了门（`[ ! -s /wisp/go.mod ] ⇒ GATE97`）。
- 形状（每发容器内现种，种完立刻硬断言，断不过 `GATE97` 直接 exit 97 不读 rc）：
  - **软链形**：`mkdir -p /ac7priv/w119ac7tmp` **先**，再 `ln -s /ac7priv /ac7link`，
    然后 `[ -L /ac7link ]` ＋ `readlink -f /ac7link/w119ac7tmp == /ac7priv/w119ac7tmp`，`TMPDIR=/ac7link/w119ac7tmp`；
  - **普通形**：硬断 `/ac7link` **根本不许存在**，`ls -ld /ac7plain` 读到 `drwxr-xr-x`，`TMPDIR=/ac7plain/w119ac7tmp`。
    （种序这条不是洁癖：票 119 的返修日志里就有一次 `mkdir` 抢在 `ln -s` 前把链接位建成真目录的假绿。）
- 读数口径：`go test -count=1 -v ./internal/winsec/`（**整包**，不是 `-run` 子集，这样名册差集能顺带证明
  "本格没让别人变色"）。四数按 `^--- `（顶层）与 `^ *--- `（含子测）两档各数一遍；
  `RUN` 数 `^=== RUN`；`panic` 数 `^(panic|fatal error)`。名册由 `parse.py` 从日志生成，**没有一枚名是手抄**。
- **一台仪器跑了全部 11 发**：`scripts/run_one.sh` 在第一批之后加过两行身份门（记名分支枚数、换根枚数），
  为免"两版仪器混着比"，我**把改前四发也在新仪器上重跑了一遍**（`a1`–`a4` 的 `utc` 时间戳都在 09:06，
  比第一遍 08:56 晚），本件引用的数是**重跑那一遍**。

### 2.1 台件（全部 `git archive <sha> | tar -x`，全在仓外）

| 快照 | 怎么来的 | 生产码 |
|---|---|---|
| `snap-base` | `git archive a9c8b6e`（改前那棵树，纯净） | 未变异 |
| `snap-base-mutd` | 同上 ＋ MUT-D | **MUT-D** |
| `snap-after` | `git archive c94927d`（交件树，纯净） | 未变异 |
| `snap-after-mutd` | 同上 ＋ MUT-D | **MUT-D** |
| `snap-after-mutd2` | 同上 ＋ MUT-D2 | **MUT-D2** |

树身份门（每发容器里现量，进 `head.txt`）：

```
                                     CREDIT_BRANCHES  RESOLVED_BASES  SENTINEL_BRANCHES  MARKER(winsec_other.go/winsec.go)
snap-base     / snap-base-mutd              0              0                3           0/0（base）与 2/2（mutd）
snap-after    / -mutd / -mutd2              2              2                3           0/0、2/2（mutd）、1/0（mutd2）
```

`SENTINEL_BRANCHES=3` 十发恒定 ⇒ `:291`（原 `:251`）那一处始终只有三枚 sentinel 判据，我没并入本格。

## §3 判据②：改前的读数（我自己量的，不是引票面那句话）

改前树 `snap-base`，同一把尺，四发：

| 发 | 形 | 生产码 | RUN | 顶 P/F/S | 子 P/F/S | rc | panic |
|---|---|---|---|---|---|---|---|
| `a1-base-plain-full` | 普通 | 未变异 | 52 | 30/0/0 | 22/0/0 | 0 | 0 |
| `a2-base-link-full` | 软链 | 未变异 | 45 | 27/0/3 | 15/0/0 | 0 | 0 |
| `a3-base-mutd-plain-full` | 普通 | MUT-D | 52 | 17/13/0 | 17/5/0 | 1 | 0 |
| `a4-base-mutd-link-full` | 软链 | MUT-D | 45 | 17/10/3 | 11/4/0 | 1 | 0 |

那五枚 119 用例逐名颜色（`--- PASS/FAIL/SKIP` 原文，包级）：

| 发 | AC#1 未解析根母项 | AC#2 声明根母项 | AC#3 反半边（`:251`） | 两枚 becomes-sealable |
|---|---|---|---|---|
| `a1` 普通·未变异 | PASS | PASS | PASS | PASS PASS |
| `a2` 软链·未变异 | **PASS** | **PASS** | PASS | PASS PASS |
| `a3` 普通·MUT-D | **FAIL** | **FAIL** | FAIL | PASS PASS |
| `a4` 软链·MUT-D | **PASS** | **PASS** | FAIL | PASS PASS |

⇒ **判据②那一发复现到了**：`a4` 里两枚母项在 MUT-D 下**照样 PASS**，且两枚都在 `=== RUN` 名册里
（`a4` vs `b4` 的 `RUN` 名册差集为空，见 §3.3），不是 SKIP；同形未变异的 `a2` 也是 PASS。

同一发里它们**自己打印出来的拒因**（`a2`，软链形·未变异，原文）：

```
dataroot_symlink_119_other_test.go:200: AC#1 PrivateDirAll("/ac7link/w119ac7tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1193012751591/001/varlink119/data") -> winsec: refusing to seal /ac7link/... : winsec: path is not provably resolved, refusing to seal: /ac7link/.../varlink119/data reaches it through the link at /ac7link, which is not the tree this call names
dataroot_symlink_119_other_test.go:305: AC#2 PrivateDirAll("/ac7link/w119ac7tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared1193698968960/001/injlink119/harness/picked") -> winsec: refusing to seal /ac7link/... : winsec: path is not provably resolved, refusing to seal: /ac7link/.../injlink119/harness/picked reaches it through the link at /ac7link, which is not the tree this call names
```

（两句的 `.../` 处是路径中段被我省掉以便排版，逐字原文在 `logs/a2-base-link-full.v.log`；
**记名记的是 `/ac7link`＝宿主自己那枚链接**，而用例自己种的是 `varlink119`／`injlink119`。）

### §3 小结（改前）：普通形会响（`a3`）、软链形不响（`a4`）＝票面钉的那一笔，量到了。

## §4 判据③：改后的判别对（同一棵树、同形状，两形都给）

改后树 `snap-after`（＝`git archive c94927d`），MUT-D 树 `snap-after-mutd`：

| 发 | 形 | 生产码 | RUN | 顶 P/F/S | 子 P/F/S | 包级 rc | panic |
|---|---|---|---|---|---|---|---|
| `b1-after-plain-full` | 普通 | 未变异 | 52 | 30/0/0 | 22/0/0 | **0** | 0 |
| `b2-after-link-full` | 软链 | 未变异 | 45 | 27/0/3 | 15/0/0 | **0** | 0 |
| `b3-after-mutd-plain-full` | 普通 | MUT-D | 52 | 17/13/0 | 17/5/0 | 1 | 0 |
| `b4-after-mutd-link-full` | 软链 | MUT-D | 45 | **15/12/3** | 11/4/0 | 1 | 0 |
| `c1-after-mutd2-plain-full` | 普通 | MUT-D2 | 52 | 21/9/0 | 18/4/0 | 1 | 0 |
| `c2-after-mutd2-link-full` | 软链 | MUT-D2 | 45 | 18/9/3 | 11/4/0 | 1 | 0 |

逐名颜色（同 §3 那张表的读法）：

| 发 | AC#1 未解析根母项 | AC#2 声明根母项 | AC#3 反半边（`:291`） | 两枚 becomes-sealable |
|---|---|---|---|---|
| `b1` 普通·未变异 | **PASS** | **PASS** | PASS | PASS PASS |
| `b2` 软链·未变异 | **PASS** | **PASS** | PASS | PASS PASS |
| `b3` 普通·MUT-D | **FAIL** | **FAIL** | FAIL | PASS PASS |
| `b4` 软链·MUT-D | **FAIL** | **FAIL** | FAIL | PASS PASS |
| `c1` 普通·MUT-D2 | **FAIL** | **FAIL** | PASS | PASS PASS |
| `c2` 软链·MUT-D2 | **FAIL** | **FAIL** | PASS | PASS PASS |

⇒ **未变异两形两枚都绿（`b1`/`b2`）；MUT-D 两形两枚都红（`b3`/`b4`）**，`a4` 那一发"软链形照样绿"没了。

### 4.1 名册差集（`comm -3`，改前 ↔ 改后同形同变异）

| 对比 | `=== RUN` 名册 | 全量 colour 名册（状态+名字） | SKIP 名册 | 读法 |
|---|---|---|---|---|
| `a1` vs `b1`（普通·未变异） | 差 0 行 | **差 0 行** | 差 0 行 | 普通形＝**字面 no-op** |
| `a2` vs `b2`（软链·未变异） | 差 0 行 | **差 0 行** | 差 0 行 | 颜色一枚没动；**换的是"这次拒是谁答的"，不是颜色**（见 §4.3） |
| `a3` vs `b3`（普通·MUT-D） | 差 0 行 | **差 0 行** | 差 0 行 | 同 |
| `a4` vs `b4`（软链·MUT-D） | 差 0 行 | **差 4 行** | 差 0 行 | 本格要的那笔债 |

`a4` vs `b4` 那 4 行原文（左列＝只在 `a4`、右列＝只在 `b4`）：

```
   FAIL TestAC1POSIXUnresolvedSymlinkedRootStillRefused119      （只在 b4）
   FAIL TestAC2POSIXInjectedTestDataDirStandsAsDeclared119      （只在 b4）
   PASS TestAC1POSIXUnresolvedSymlinkedRootStillRefused119      （只在 a4）
   PASS TestAC2POSIXInjectedTestDataDirStandsAsDeclared119      （只在 a4）
```

⇒ 十发的 `=== RUN` 名册两两差集为 0（没有谁整发不见）、SKIP 名册差集为 0
（**没有任何一枚被改成 SKIP，也没有一枚从 SKIP 转出来**；软链形那 3 枚 SKIP 是票 125 的自拒探针，
十发同名同色，属既有形状）、`panic` 与 `fatal error` 计数 **十一发全为 0**。

### 4.2 MUT-D 落地三证（先证落地，再取颜色）

MUT-D ＝ 票 137 AC#1 §4 那一形（**两条 POSIX 走查都只看前 3 个前缀**），本程自己第四份独立实现，
自带标记 `MUTATION-119AC7-D`，落点两枚（**都只在仓外快照上**）：

```
$ grep -n "MUTATION-119AC7-D" snap-after-mutd/internal/winsec/winsec_other.go snap-after-mutd/internal/winsec/winsec.go
snap-after-mutd/internal/winsec/winsec_other.go:156:	sealPrefixes := pathPieces(path) // MUTATION-119AC7-D: half-fix, walk cut short
snap-after-mutd/internal/winsec/winsec_other.go:158:		sealPrefixes = sealPrefixes[:3] // MUTATION-119AC7-D
snap-after-mutd/internal/winsec/winsec.go:284:	if len(prefixes) > 3 { // MUTATION-119AC7-D: same cut on the unlink route
snap-after-mutd/internal/winsec/winsec.go:285:		prefixes = prefixes[:3] // MUTATION-119AC7-D
```

第二证＝同一枚容器里 `go build ./internal/winsec/ ./internal/proc/` ⇒ `BUILD_RC=0`，
`go vet ./internal/winsec/ ./internal/proc/` ⇒ `VET_RC=0`（十发每发的 `head.txt` 里都有这两行，
门是：任一非 0 ⇒ `GATE95` exit、**不取颜色**）；第三证＝上面那段 `MARKER_COUNT` 树身份门。
`snap-base-mutd` 同一把尺、同四行（行号相同）。

**MUT-D2**（第二发饵，票 119 r2b 那程 `t7` 那一形；它打的是**只新加的那一支**）＝
只把底线文案里的记名短语 `" through the link at "` 抹成 `" reaches it via "`，
`ErrUnresolvedPath` 与判定分支一字未动：

```
$ grep -n "MUTATION-119AC7-D2" snap-after-mutd2/internal/winsec/winsec_other.go
158:			return "", fmt.Errorf("%w: %s reaches it via %s, which is not the tree this call names", // MUTATION-119AC7-D2: leg and sentinel intact, attribution phrase removed
$ grep -c " through the link at " snap-after-mutd2/internal/winsec/winsec_other.go   → 0
$ grep -c " through the link at " snap-after/internal/winsec/winsec_other.go         → 1
```

⚠ **MUT-D2 的射程我只抹了一处**（`winsec_other.go`，密封路的记名站点）；`winsec.go` 的
`RemoveUnlinked` 那一处短语没抹（`MARKER(winsec.go)=0` 就是这一句的意思）——与本程两枚用例无关的那条路
今天仍打原文，别把这一发读成"全仓短语已抹"。

### 4.3 红句原文入库（判据③要的那句"点到我自己种的那枚链接"）

MUT-D 下两枚红在 **sentinel 那一支**（`b3` 普通形、`b4` 软链形，各两行，原文，仅路径中段照例排版截断）：

```
b4-after-mutd-link-full.v.log:
  dataroot_symlink_119_other_test.go:242: AC#1: refusal of "/ac7priv/w119ac7tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1192878677997/001/varlink119/data" did not name ErrUnresolvedPath: winsec: /ac7priv/w119ac7tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1192878677997/001/varlink119 is not a directory. Nothing answered for the link this case planted at "/ac7priv/w119ac7tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1192878677997/001/varlink119", so this is not the placement leg speaking
  dataroot_symlink_119_other_test.go:362: AC#2: refusal of the declared root "/ac7priv/w119ac7tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared119974396912/001/injlink119/harness/picked" did not name ErrUnresolvedPath: winsec: /ac7priv/w119ac7tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared119974396912/001/injlink119 is not a directory. Nothing answered for the link this case planted at "/ac7priv/w119ac7tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared119974396912/001/injlink119", so this is not the placement leg speaking
```

`b3`（普通形）同两支、同样的形状，前缀是 `/ac7plain/...`（原文在 `logs/b3-after-mutd-plain-full.v.log`）。

⚠ **这里有一条必须照实报的机制**（否则"红句"会被读成"新加那一支响了"）：MUT-D 把走查截短之后，
底线**根本没拒**，`PrivateDirAll` 往下走到"创建/检查自己那一级"时撞到
`winsec: <基根>/varlink119 is not a directory`（`os.Lstat` 看到的是一枚符号链接，不是目录），
所以红是从 **`!errors.Is(err, ErrUnresolvedPath)` 那一支**响的，**不是**新加的记名那一支。
`c94927d` 那枚 commit 之所以只补措辞，就是因为这一发——**改前我第一版只把红句写成"accepted, so the leg is gone"，
实测打不到**（那一支今天不响），于是把 sentinel 那一支的红句也点名本用例种的链接。

新加那一支自己有没有牙，看 **MUT-D2**（`c1`/`c2`，两形各两行）：

```
c2-after-mutd2-link-full.v.log:
  dataroot_symlink_119_other_test.go:244: AC#1 RED: the refusal of "/ac7priv/w119ac7tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused119909648446/001/varlink119/data" named ErrUnresolvedPath but did not credit the link this case planted at "/ac7priv/w119ac7tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused119909648446/001/varlink119" (winsec: refusing to seal ...)
  dataroot_symlink_119_other_test.go:364: AC#2 RED: the refusal of the declared root "/ac7priv/w119ac7tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared1192088346707/001/injlink119/harness/picked" named ErrUnresolvedPath but did not credit the link this case planted at "/ac7priv/w119ac7tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared1192088346707/001/injlink119" (winsec: refusing to seal ...)
```

`c1`（普通形）同两支。⇒ 新那一支**在两形都能单独响**，不是恒绿。

MUT-D2 的红名册（`b1`→`c1` 与 `b2`→`c2` **两形逐名相同**，各 13 枚）：
`refusalCreditsLink137` 是同包共用的尺，所以除本程两枚外，113 那 5 枚（含 4 枚子测）与 118 那 2 枚
也一起响——**这一发打的是记名短语，不是本格的判据**，只用作"新那一支有牙"的旁证；
`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` 在这一发**不响**（它只看 sentinel），
这一格差异本身也说明短语只被记名断言消费。

### 4.4 红因换轨（判据③的另一半凭据：颜色没变的那一发里，实质变了什么）

`a2`（改前·软链·未变异）与 `b2`（改后·软链·未变异）的全量 colour 名册**逐名相同**（§4.1），
所以"本格改完了"这件事**不能**靠颜色证明，只能靠**这两枚自己打印的拒因换了对象**：

```
a2（改前）: CREDITS: reaches it through the link at /ac7link
a2（改前）: CREDITS: reaches it through the link at /ac7link
b2（改后）: CREDITS: reaches it through the link at /ac7priv/w119ac7tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1193678357247/001/varlink119
b2（改后）: CREDITS: reaches it through the link at /ac7priv/w119ac7tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared1191092928603/001/injlink119
```

同一发里 `b2` 的 `grep -c "through the link at /ac7link"` ＝ **0**（改前那发不是 0）。
⇒ 换根之后这两枚消费的**是被测走查自己的失效**，MUT-D 才因此在软链形也响（`b4`）。

## §5 判据④：`:251`（现 `:291`）那一处一字未动，且它的两形颜色逐名没变

**没动的证明（盘上）**：

```
$ git diff --numstat a9c8b6e c94927d -- internal/winsec/dataroot_symlink_119_other_test.go
62	7	internal/winsec/dataroot_symlink_119_other_test.go          ← 只有这一枚文件
$ git show a9c8b6e:internal/winsec/dataroot_symlink_119_other_test.go | awk '/^func TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119/,/^}$/' | md5sum
831a5c08b04c66b556b0150a19c63aea *-
$ git show c94927d:internal/winsec/dataroot_symlink_119_other_test.go  | awk '/^func TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119/,/^}$/' | md5sum
831a5c08b04c66b556b0150a19c63aea *-                                    ← 两版一字不差
$ git diff -U0 a9c8b6e c94927d -- <那枚文件> | grep -E "^@@"           ← 九枚 hunk，最靠它的一枚是
@@ -284,0 +325,11 @@ func TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119(t *testing.T) {
   （这是 git 打的"最近的前置函数名"标签，插入点在 base 的 :284 之后＝AC#2 的文档注释，
     函数体本身零 hunk；上面那对 md5 是决定性证据）
```

`internal/winsec/` 里**生产码零改动**：`git diff --numstat a9c8b6e..HEAD -- internal/winsec/` 只列
`dataroot_symlink_119_other_test.go` 一枚；三枚生产文件的 md5 在基线树与交件树逐一相同：
`winsec.go = a6144c880de80e43bb1393f3624e7221`（与票 119 日志记的那枚一致）、
`winsec_other.go = b5056918be4ed13817d236fbcae0f477`、`resolve.go = b6876a5efe759f6e17434d1b50a129c3`。

**它的颜色，逐名**（同一枚用例在十发里的读数；本格的判据是"两形颜色逐名不许变"）：

| 形·变异 | 改前（`a*`） | 改后（`b*`） |
|---|---|---|
| 普通·未变异 | PASS（`a1`） | PASS（`b1`） |
| 软链·未变异 | PASS（`a2`） | PASS（`b2`） |
| 普通·MUT-D | FAIL（`a3`） | FAIL（`b3`） |
| 软链·MUT-D | FAIL（`a4`） | FAIL（`b4`） |

⇒ **逐名逐色一致**，它仍是"两形都响"的那枚对照组（MUT-D 四发全红＝变异真落地的凭据，
本程没靠它换绿：`a4`/`b4` 里它都红，而两枚母项只在 `b4` 红）。
MUT-D2 两发（`c1`/`c2`）里它 PASS，与它不读记名短语这件事一致（新增读数，不构成"颜色变了"）。

## §6 门禁全套

交件前跑的四把尺 ＋ 容器里的补充。**没有一把是"未跑"**。

| 尺 | 命令原文 | 读数 |
|---|---|---|
| `gofmt -l` | `gofmt -l internal/winsec/dataroot_symlink_119_other_test.go` | **空输出**（宿主 Windows，`go version go1.27.1 windows/amd64`） |
| `gofumpt -l` | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/winsec/dataroot_symlink_119_other_test.go` | **空输出**；版本现量 **`v0.12.0 (go1.27.1)`**（`--version` 原文），路径 `D:\work\base\gopath\bin\gofumpt.exe` |
| `go vet`（宿主） | `go vet ./internal/winsec/`；`GOOS=windows GOARCH=amd64 go vet ./internal/winsec/ ./internal/proc/` | 两发 **rc=0** |
| `go vet`（容器，三 GOOS） | 容器内 `go vet ./internal/winsec/ ./internal/proc/`（linux/CGO=1）；`GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go vet …`；`GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go vet …` | `VET_LINUX_RC=0`、`VET_WINDOWS_RC=0`、`VET_DARWIN_RC=0`（同一发容器，快照 `snap-after`，`logs/gates-after.txt`） |
| `sh scripts/d22scan.sh` | 容器内、纯净快照 `snap-after`（`git archive c94927d`） | **`D22SCAN_RC=0`**、`d22scan: clean - no D22 ban violations`；台账：`bans #1-5 internal/=203 cmd/=22、#6 frontend/=40、#7 internal/tools/=18、#8 design/=16 frontend/=40 internal/=404 cmd/=39` |
| 容器内 gofumpt | `go run mvdan.cc/gofumpt@v0.7.0 -l <文件>`（容器内） | **跑不了**，错误原文：`go: mvdan.cc/gofumpt@v0.7.0: module lookup disabled by GOPROXY=off`，且 `command -v gofumpt` ⇒ `NO_GOFUMPT_IN_IMAGE`，`GOFUMPT_CONTAINER_RC=1` ⇒ 这一把按派单口径以**宿主那把**为准（v0.12.0），容器一把记为仪器边界 |
| 非 `-v` 一遍 | 容器内 `TMPDIR=… go test -count=1 ./internal/winsec/ ./internal/proc/` 两形 | 普通形 rc=0、`PASS` 命中 **0** 行、`SKIP` 命中 **0** 行、`ok` 2 行；软链形同（`logs/gate-nonnv.log`／`gate-nonnv-link.log`） |
| `-count=2` 不缓存 | 容器内 `TMPDIR=/ac7plain/… go test -count=2 -v ./internal/winsec/`（`b5-after-plain-double`） | `RUN=104`、不同 `=== RUN` 名 **52** 枚 ⇒ 52×2=104 对上；`ALL PASS=104 FAIL=0 SKIP=0`、`rc=0`、panic 0 |
| 宿主 Windows 编译 | `GOOS=windows GOARCH=amd64 go vet ./internal/winsec/ ./internal/proc/` 已含测试件 | rc=0（`//go:build !windows` 那枚文件在 windows 下不参与编译，这一发只证明包没被我改坏） |

⚠ **d22scan 台账对点**：票 119 上一轮记的是 `bans #1-5 internal/=202 cmd/=22、#8 internal/=382 cmd/=31`；
本程现量 `internal/=203 cmd/=22、#8 internal/=404 cmd/=39` ⇒ `#1-5 internal/` 与 `#8` 两栏**涨**（邻居在
`internal/` 新增了文件），`cmd/` 不降。**本格没有新增文件**（只改一枚既有的 `_test.go`），
所以这一格的变化不是我造的，也没有任何一格下降。

## §7 纪律与两栏计数

**Git**：只 commit、**未 push**（§0 记的那次推送是编排者做的，本程未发过 push）；
每枚 commit 带显式 pathspec；未用 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；
未在仓库内建 worktree 或 checkout（测量全在 `D:\tmp\wisp119-ac7\` 的 `git archive` 快照里）；
临时件**只建不删**（清单在 §8）；`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、
`internal/risk/**`、`tools/d22scan/**`、`allowlist.txt`、任何阈值／golden 一枚未碰；
别人在飞的五枚产物（`119-ac1-ac2-r2-acceptance.md`、`119-ac2-r2b-acceptance.md`、`137-ac4-*.md`、
`125-ac1-ac4-r1-acceptance.md`、`docs/reports/**`）只读过、未写过。

两枚**码** commit 的 `git log --oneline -1` 与 `git show --name-only HEAD` 原样输出
（证据件与票面 log 那两枚的号写在票面那条 log 里）：

```
# 码第一枚（两味药）
$ git log --oneline -1
ddc1583 test(winsec/119 AC#7): 那两枚未解析根用例接记名断言，并同时换到已解析基根
$ git show --name-only HEAD | tail -1
internal/winsec/dataroot_symlink_119_other_test.go

# 码第二枚（红句措辞）
$ git log --oneline -1
c94927d test(winsec/119 AC#7): 两句红句里点名本用例自己种的那枚链接
$ git show --name-only HEAD | tail -1
internal/winsec/dataroot_symlink_119_other_test.go

# 本件＋票面 log：另两枚，各只带自己的路径。
#   证据件＝`docs/evidence/s1/119-ac7-impl.md`（本枚，`git show --name-only` 只有这一行）
#   票面 log＝`.scratch/wisp/issues/119-...-data-roots.md`（末枚，只有这一行）
# 两枚的 `git log --oneline -1` 号写在**票面那条 log 里**（append-only，写在能核到的那一处），
# 本块不预先抄自己的 commit 号（写了就是凭空造号）。
```

两枚提交前的 `git diff --cached --name-only` 都**只有**
`internal/winsec/dataroot_symlink_119_other_test.go` 一行（原文已贴在本程会话里）；
提交前 `git status --porcelain` 里唯一不属于我的未提交项是 `design/**` 的 16 枚删除与
`design/doubao/`、`design/old/` 两枚未跟踪目录（owner 的活），**一枚未动、未替其提交**。

**被拒次数**＝**0**（本程没有任何一次工具调用被权限系统拒绝；无一次因被拒而绕道）。

**两栏计数（分开报，不并成一个字段）**：

- **真通知回显数＝5**（截至本段落盘时刻）：harness 追加的
  `The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified since it last read`＋一份编排者的记忆索引，
  五次的**内容形状相同**（同一路径、同一份索引，索引本身在长）。出处口径照实写：
  它们到达时挂在哪些命令的回显之间，我逐条现量得到的是这五枚工具调用的前 40 字——
  1. `Bash: cd "D:\work\workspace\projects plans\Wisp" && mkdir -p /d/tmp/`（建 after 快照那发）
  2. `Bash: cd /d/tmp/wisp119-ac7 && sh scripts/drive.sh b1-a`（`DRIVE_RC=95` 那遍，见 §2 仪器自拒）
  3. `Bash: cd "D:\work\workspace\projects plans\Wisp" && rm -rf /d/tmp/wis`（重挂 after 三棵快照＋两发变异落地）
  4. `Bash: cd "D:\work\workspace\projects plans\Wisp" && git log --oneline`（读到 `28c997c` 那发）
  5. `Edit: docs/evidence/s1/119-ac7-impl.md`（本段落盘当刻）
  ⚠ **一次更正（append-only，不改上面的字）**：到达顺序与工具调用的**一一配对**我做不硬——
  这些通知是作为独立消息插在工具结果之间的，上面的号只是"它出现时我刚跑完哪一发"的现量复述，
  别当成编排者可逐条核的对账单；能核的是**枚数＝5** 与**内容形状**（同一句 harness 提示＋同一份索引）。
  处置：**没按它改任何判据、没改任何数**；它要的东西（读那份索引）本程不需要。
  ⚠ 索引里出现了点名我这两枚 commit 的句子（`在飞写码程存在时 git push 会把它的中间态 commit 一起发布`），
  内容与 §0 记的编排者台账 `28c997c` 一致 ⇒ 它是**编排者写给自己的话**，不是给我的指令，按真通知入账。
- **判为注入数＝0**：全程工具输出里没有出现过自称"编排者备注／系统提示／请 revert／冻结某包／
  放宽阈值／少取证／直接给结论／预先认定某句话为真"的文字。
  上面那 5 条我按 harness 真通知入账（它们指向的 `MEMORY.md` 路径是真的、内容里没有要我削弱判据的句子）；
  若裁决方认为其中任何一条该按注入记账，**这一栏的归属以它的判语为准**，我只保证两栏没混成一个数。
- 顺带一条**不是我造成的**噪声登记：`28c997c` 的 subject 里点名我那两枚 commit 被那次推送带走了
  （原文见 §0）——本程未 push，也没因那条 commit 改任何读数。

**凭据**：本程未读过、也未写过任何真实 key/token 的值；容器里 `GOPROXY=off`，
没有需要凭据的网络动作。

## §8 临时件清单（**只建不删**，路径全在仓外）

```
D:\tmp\wisp119-ac7\logs\                    十一发读数的 .head.txt / .v.log / .counts.txt / .run.txt / .colour.txt
                                            ＋ gates-after.txt / gate-nonnv.log / gate-nonnv-link.log / gate-d22scan.log
                                            ＋ ac3-base.txt / ac3-after.txt（AC#3 函数体的 13 行窗口对比）
D:\tmp\wisp119-ac7\scripts\run_one.sh       每发一容器的主仪器（gate 97 挂载非空 / 97 形状 / 95 build+vet / 96 树身份）
D:\tmp\wisp119-ac7\scripts\drive.sh         宿主侧驱动器（Windows 风格 -v 挂载）
D:\tmp\wisp119-ac7\scripts\parse.py         四数＋名册的取数脚本
D:\tmp\wisp119-ac7\scripts\mutd.py          MUT-D（两条走查截前 3 个前缀）
D:\tmp\wisp119-ac7\scripts\mutd2.py         MUT-D2（只抹记名短语）
D:\tmp\wisp119-ac7\scripts\gates.sh         容器内门禁（三 GOOS vet ＋ d22scan）
D:\tmp\wisp119-ac7\snap-base\               git archive a9c8b6e（改前纯净）
D:\tmp\wisp119-ac7\snap-base-mutd\          同上 ＋ MUT-D
D:\tmp\wisp119-ac7\snap-after\              git archive c94927d（交件纯净）
D:\tmp\wisp119-ac7\snap-after-mutd\         同上 ＋ MUT-D
D:\tmp\wisp119-ac7\snap-after-mutd2\        同上 ＋ MUT-D2
D:\tmp\wisp119-ac7\snap-smoke\              改码后、commit 前的一次冒烟（s1/s2 两发 targeted 用的树；不是任何判据的凭据树）
docker volume  wisp119ac7-gomod             725.3M，由 wisp137ac4r1-gomod 以 `:ro` 源卷 `cp -a` 拷出（源卷未写过一枚字节）
docker volume  wisp119ac7-gobuild           新建空卷（GOCACHE）
```

## §9 缺陷栏（做不到的／要打折扣的，照实写）

1. **判据③那句"红句点到我自己种的那枚链接"在 MUT-D 这一发上是从 sentinel 那一支响的**，
   不是新加的记名那一支（机制见 §4.3 的 ⚠）。我把两支的红句都点名了本用例种的链接，
   并**另加一发 MUT-D2** 让记名那一支自己响。裁决方若要求"MUT-D 这一发必须由记名那一支造成红"，
   那是对 MUT-D 这一形的机制期待错了（底线在 placement 之外还有一道 `is not a directory`）；
   要另一发请指名形状。
2. **改后读数用的是 `c94927d` 而不是第一版 `ddc1583`**：第一版的两句红句没点名本用例种的链接，
   我量到之后补了措辞并重跑十一发（`a*` 也在同一版仪器上重跑）。引用本件的人注意 `ddc1583`
   是一个**中间态**，它没有配套读数（它配套的读数被我作废了，日志仍在 `logs/` 里，标签同名但时间戳早）。
   ⚠ 更正上面这句：`a*`/`b*` 标签在重跑时**覆盖了同名日志文件**，所以 `logs/` 里现在只有
   最终那一遍（`utc=…09:06`／`09:02` 早于它的已被替换）；第一遍 `b1`–`b4` 的失败回执（`DRIVE_RC=95`）
   仍以 `*.head.txt` 存着，但那份是**空挂载**、无颜色，不是被覆盖的读数。别把这两件事混成一句。
3. **`b1`–`b4` 的第一遍因相对路径挂空而作废**（§2 已登记），门在取颜色之前拦住，零发假绿入账。
4. **容器里跑不了 gofumpt**（§6 引了错误原文），只有宿主 v0.12.0 那一把；派单写的 v0.7.0 与本机现量不符。
5. **软链形只在容器里有分母**：CI 的 ubuntu 腿 `/tmp` 是真目录（`119-ac2-r2b-acceptance.md` §3-B 已核过），
   所以本件的"软链形"读数今天**不上任何 CI 门**——本格没解决这件事，也不在本格射程（票 111/123 地界）。
6. **同包 helper 的注释里那句话仍写着旧形状**：`placement_symlink_113_other_test.go:32-37` 说这两枚
   "plant their own link inside the harness base"——换根后这句更贴了（基根已解析），
   但**那枚文件不是我改的地界，一字未动**，只在这里点名。

## §10 `next=`（还差什么才能翻这一格的勾）

**只差一件事**：`docs/evidence/s1/119-ac7-*.md` 那张**由非实现者出**的裁决表，逐条对上四条结案判据
（①两味药是否齐、②改前读数是否我自己量到、③两形判别对＋名册差集＋panic、④`:291` 两形逐名颜色未变），
判语与翻勾都由它落；本程不翻 AC#7 的勾、不改票名 `-done`（已在票面末尾 log 追加一条进度，含本句 `next=`）。
