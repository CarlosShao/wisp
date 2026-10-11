# 307-r1 · 20-ac2-shapes —— `AC#2` 三形并排＋反形敏感度（全部在**仓外导出树**里做）

⚠ **编号对照（具名报回一处会读歪的地方）**：票面 `AC#2` 原文的 (c) ＝"把判据换成反形它也不响"那一问；
追加节 §4 第 1 句把 ★中间形（只换 `:428`）也称为"第三形"，派单再把恒真形排为 (d)。
本文件按**派单的四个字母**编号，映射＝票面 (a)→§1、票面 (b)→§2、追加节 §4 的中间形→§3、票面 (c) 恒真形→§4。
**四格都交了**，⛔ 有一格被编号吃掉。

## 0. 导出树与台面（尺）

- 建树＝`git archive HEAD go.mod go.sum internal cmd | tar -x -C /c/Users/swq/tmp/wisp-307-r1-tree-main/`
  ⚠ **与派单不符处（具名报回）**：派单给的尺逐字是 `git archive HEAD internal cmd | tar -x …`，那样得到的树**没有 `go.mod`／`go.sum`**，
  `go test ./internal/risk/` 在其中直接 `go: cannot find main module`。本腿把 `go.mod go.sum` 一并 archive 才跑得起来 ⇒ 尺已按盘上事实改写。
- 母仓工作树在变异期间**一字节未动**（每发之后都 `Copy-Item -Force $pristine $target` 还原，且⛔ 有任何写操作落在 `D:\work\…\Wisp\internal\`）。
- pristine 副本＝`C:\Users\swq\tmp\wisp-307-r1-pristine-syncdirs.go`（从导出树里的 HEAD 那份 `syncdirs.go` 直接 `cp` 而来）。
- ★**变异工具自检（round-trip）**＝`-Mutation pristine` 那一发：用与所有变异同一对 `ReadAllText`/`WriteAllText` 写回再还原，
  `sha256_tree_after_restore == sha256_pristine == aa4779bb52704f4333a286c1a8396868579daa8d0c26634022ef9a8830ff2d36` ⇒ **写入路径本身⛔ 吃字节**。
- 台面＝`aliased`（`TEMP/TMP/TMPDIR = C:\Users\swq\tmp\RISKAT~1\RUNNER~1`）；尺＝`go test ./internal/risk/ -count=1 -v -run 'TestSync'`（`-Run` 参数见各节）。
- 变异脚本＝`logs/run-mutation.ps1`（字面量锚逐字写死；任何一枚锚没命中就 `throw`，⛔ 静默空变异）。
  件＝`logs/mut-<名>-<台面>.txt`，每份件自带 `---- mutation=… rc=… PASS=… FAIL=… SKIP=…` 行与两枚 sha256。

四枚变异用的**字面量锚**（＝HEAD `51a0a0f2` 里 `syncdirs.go` 的真句子）：

```
L_ROOT  if aUnder(candForms, anchorForms(e.canon)) {
L_HOME  if !s.complete && s.home != "" && aUnder(candForms, anchorForms(s.home)) {
L_FORMS <TAB>candForms := formsOf(cand)
```

---

## 1. (a) 把修复挪开：只把 `L_HOME`（＝票面 `:428` 那一支）换回 `isUnder`

`logs/mut-a-aliased.txt`：`rc=1 PASS=31 FAIL=6 SKIP=0`（`-run TestSync`，两台面里的 aliased）。

**那 5 枚指名用例当场红（5/5）**，红句逐字：

```
--- FAIL: TestSyncFixtureFallbackAndMatch (0.01s)
    syncdirs_test.go:175: under-profile write must stay sync-suspect while no confirmed root exists
--- FAIL: TestSyncFallbackNotDisarmableByWeakRoot (0.01s)
    syncdirs_test.go:213: source "default": under-profile fallback must stay armed
    syncdirs_test.go:213: source "fixture": under-profile fallback must stay armed
    syncdirs_test.go:213: source "options": under-profile fallback must stay armed
    syncdirs_test.go:213: source "": under-profile fallback must stay armed
--- FAIL: TestSyncUnverifiedRootKeepsFallback (0.01s)
    syncdirs_test.go:233: suspect fallback must stay armed
--- FAIL: TestSyncSuspectFallbackIsComponentBounded (0.00s)
    syncdirs_test.go:355: under-profile must be suspect
--- FAIL: TestSyncSuspectFallbackWhenUndetectable (0.01s)
    syncdirs_test.go:402: under-profile path must be sync-suspect, got {Sync:false Root:{Provider: Path: Source:} Why:write target is not under any sync root}
```

外加本票新夹具 `TestSync307AliasHomeUnderProfileStaysSuspect`（`syncdirs_alias_307_windows_test.go:71`）同形红＝
`ticket 307 fail-open: write under an alias-homed profile read as not-sync, home="C:\\Users\\swq\\tmp\\RISKAT~1\\RUNNER~1\\TESTSY~1\\001\\THIRTY~1" got {Sync:false …}`。

★**一枚与派单预期不符、必须具名报回的读数**：(a) 形里 `TestSyncFixtureFallbackAndMatch` 红的是 **`:175`**，⛔ 是 `:162`。
原因＝这一形⛔ 撤 `L_ROOT`，`L_ROOT` 已经修好，于是 `:162` 那枚"写在 fixture 根下"**改由 `:422` 根成员那一支拿到 `Sync:true`**，
断言 `:161` 过了；同一个函数体再往下走到 `:174`（`home\LocalOnly\a.txt`，只能来自 `:428`）才红。
⇒ **`AC#2(a)` 的"那 5 枚当场红"成立（5 枚全 FAIL），但 `:162` 这一枚的红在 (a) 形下换了行号**——
派单那句"`:428` 换回 `isUnder` ⇒ 那 5 枚当场红"⛔  specifying 哪一行红，本腿把行号交全，⛔ 顺着写成"和改前一样"。

## 2. (b) 逐字节还原 ⇒ 复绿

同一棵树、同一把尺，还原后（`-Mutation pristine`）：`logs/mut-pristine-aliased.txt` ⇒ `rc=0 PASS=37 FAIL=0 SKIP=1`。
还原前后 sha256（`Get-FileHash -Algorithm SHA256`）逐字：

```
sha256_pristine              = AA4779BB52704F4333A286C1A8396868579DAA8D0C26634022EF9A8830FF2D36
sha256_tree_after_restore    = AA4779BB52704F4333A286C1A8396868579DAA8D0C26634022EF9A8830FF2D36   (五发变异每一发之后都相同)
```

母仓侧的尺（bash，`cmp -s` 不带 `head` 吞退码）：

```
cmp -s /c/Users/swq/tmp/wisp-307-r1-pristine-syncdirs.go <tree>/internal/risk/syncdirs.go        rc_cmp_after_restore=0
cmp -s internal/risk/syncdirs.go <tree>/internal/risk/syncdirs.go                                 rc_cmp_worktree_vs_tree=0
```

⇒ 导出树里那枚文件在全部变异跑完之后**与 pristine、与母仓工作树、与 HEAD blob 三处逐字节相同**（`git status --porcelain -- internal cmd` 见 `40-ac4-gates.md`）。

## 3. (c) ★中间形＝只换 `:428`、⛔ 换 `:422`（`L_ROOT` 退回 `isUnder(cand, e.canon)`）

尺＝`-run 'TestSyncFixtureFallbackAndMatch|TestZZ307ProbeFixtureRootAttribution'`，件＝`logs/mut-c-aliased.txt`。
`TestZZ307ProbeFixtureRootAttribution` ＝**只活在导出树里**的一枚打印用例（`C:\Users\swq\tmp\wisp-307-r1-tree-main\internal\risk\zz307probe_test.go`），
逐字复制 `TestSyncFixtureFallbackAndMatch` 的夹具构造（同一个 `sandbox(t)`、同一个 `jsonEsc`、同一个未创建的 `home\Cloud\NutstoreSync`），
差别只在**把 `IsSyncPath` 的返回整个打出来而⛔ 断言**。⛔ 新建过仓内任何 `.go`。

读数（`rc=0 PASS=2 FAIL=0 SKIP=0` ⇒ 中间形下那枚指名用例**绿**）：

```
ZZ307PROBE home=C:\Users\swq\tmp\RISKAT~1\RUNNER~1\TestZZ307ProbeFixtureRootAttribution4263225922\001\profile
ZZ307PROBE syncRoot=…\001\profile\Cloud\NutstoreSync exists=false statErr=GetFileAttributesEx …: The system cannot find the path specified.
ZZ307PROBE target=…\profile\Cloud\NutstoreSync\out\secret.md
ZZ307PROBE roots=[{Provider:Nutstore Path:C:/Users/swq/tmp/RISKAT~1/RUNNER~1/…/001/profile/Cloud/NutstoreSync Source:fixture}]
ZZ307PROBE complete=false
ZZ307PROBE st={Sync:true Root:{Provider:sync-suspect Path:c:\users\swq\tmp\riskat~1\runner~1\testzz307probefixturerootattribution4263225922\001\profile Source:suspect-fallback} Why:no confirmed sync location (P12); path under the user profile treated as sync-suspect (safe default)}
ZZ307PROBE st.Root={Provider:sync-suspect Path:c:\users\swq\tmp\riskat~1\runner~1\…\001\profile Source:suspect-fallback}
ZZ307PROBE SPLIT Sync=true Provider="sync-suspect" Path="c:\\users\\…\\profile" Source="suspect-fallback"
```

★**具名回答那一问**：中间形下 `TestSyncFixtureFallbackAndMatch` 拿到的 `Sync:true`
**来自 `:428`（`s.home` 兜底那一支）**，证据＝`Root.Provider == "sync-suspect"`、`Root.Source == "suspect-fallback"`、
`Root.Path` ＝ 折过的 `s.home`（小写短形链），而**⛔ 是** `Root.Provider == "Nutstore"` / `Source == "fixture"`（那才是 `:422` 根成员会报的形）；
`Why` 逐字＝`no confirmed sync location (P12); path under the user profile treated as sync-suspect (safe default)` ＝ `:428` 那一支的句子。
⇒ **编排者那一枚读码推到成立**（那枚用例对"根成员判定"这件事**⛔ 敏感**），本腿⛔ 推翻它，⛔ 顺着写——这是跑出来的。

⇒ 按票面追加节 §4 第 1 句，`AC#2` 那一格**只有在另加一枚断言 `Root.Provider/Source` 的用例把它钉住时才闭**。
本票交的就是那枚钉＝`internal/risk/syncdirs_alias_307_windows_test.go` 的 `TestSync307AliasShortRootMembership`
（Windows 同形，判据形状照 `syncdirs_other_test.go:99`：`Root.Source != "suspect-fallback"` 且 `Root.Provider == "Nutstore"`，
夹具根**⛔ 在 profile 内**⇒ `:428` 结构上到不了，阳性**只能**来自 `:422`）。
它在 `-Mutation fullrevert`（两支都退回 `isUnder`）下发红，逐字（`logs/mut-fullrevert-aliased.txt`）：

```
syncdirs_alias_307_windows_test.go:120: ticket 307 fail-open at the root-membership branch: write under a short-spelled missing root read as not-sync, root="C:\\Users\\swq\\tmp\\RISKAT~1\\RUNNER~1\\TESTSY~1\\001\\THIRTY~2\\ThirtySevenAliasMissingRoot" got {Sync:false Root:{Provider: Path: Source:} Why:write target is not under any sync root}
```

**第四形对照（fullrevert，两支同时退回 `isUnder`＝票面缺陷原形）**：`rc=1 PASS=30 FAIL=7 SKIP=1`＝
那 5 枚指名（`:162`/`:213`×4/`:233`/`:355`/`:402`，红句逐字与 `10-ac1-colours.md` §2 那批发自**同一母仓工作树改前发**完全相同）
＋本票 2 枚新夹具全红。⇒ 新夹具量的**是这枚缺陷**，⛔ 是台面。

## 4. (d) 反形敏感度：把新判定换成**恒真**

尺＝`-Mutation d` ⇒ 两行都变成 `aUnder(…) || true`，件＝`logs/mut-d-aliased.txt`：`rc=1 PASS=26 FAIL=11 SKIP=1`。

**5 枚指名用例里响的＝3 枚**（逐字）：

```
--- FAIL: TestSyncFixtureFallbackAndMatch
    syncdirs_test.go:179: path outside the user profile must not be suspect
--- FAIL: TestSyncSuspectFallbackIsComponentBounded
    syncdirs_test.go:358: sibling directory with a name prefix of the profile must not be suspect
--- FAIL: TestSyncSuspectFallbackWhenUndetectable
    syncdirs_test.go:406: path outside the user profile must not be suspect
```

**⛔ 响的＝2 枚（具名，⛔ 掩）**：`TestSyncFallbackNotDisarmableByWeakRoot`、`TestSyncUnverifiedRootKeepsFallback`。
成因＝这两枚体内**只有正向断言**（`:209`/`:212`/`:216` 与 `:229`/`:232`），⛔ 任何一枚负向钉，
恒真形对它们**天生无信号** ⇒ 按票面 `AC#2` 第三形的判据，**那两枚用例对"换成反形它也不响"这件事⛔ 敏感**。
⇒ 本格**⛔ 由"5 枚全响"来交**，交的是：
- **包级**恒真形**⛔ 是静悄悄**的：同发另有 8 枚既有钉一起红（`TestSyncWriteNegative`、`TestSyncNormalNewFileWriteNotFlagged`(`:374` 前置闸)、
  `TestSyncDotDotTailFailsClosed`(`:440` 前置闸)、`TestSyncAncestorActableLegFailsClosedOnMovedAncestor116`、
  `TestSyncRedTeamGuardPlainNewFileWriteAllowed`、`TestSyncConfirmedGradesDisarmFallbackWindows`，加本票 2 枚新夹具
  `:82`（name-prefix sibling 被扫进兜底网）与 `:139`（plain 目录被当根））⇒ 换形⛔ 退化成"什么都算 sync"而不被仪器抓到；
- **枚级**具名欠账＝上面那 2 枚。本票⛔ 动它们的期望侧（票面硬禁区），只具名登记。

⇒ 结论交法：**`AC#2(d)` 在"5 枚指名用例"这个分母上⛔ 全格成立**（3/5 响、2/5 天生无信号），
在"包级仪器"这一层成立。⛔ 写"本格闭合"。

## 5. 本文件⛔ 做到的

- ⛔ 勾 `AC#2` 任何一格（裁决者≠实现者）。
- ⛔ 在共享工作树里做任何变异（全部在 `C:\Users\swq\tmp\wisp-307-r1-tree-main\`）。
- ⛔ 把 (a) 形的 `:162` 换成 `:175` 这件事写顺（第 1 节★那条就是与本腿读数相反的写法，本腿⛔ 那么写）。
