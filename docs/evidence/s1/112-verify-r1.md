# 112 裁决表 — winsec 步在 CI runner 首跑红：非实现者复判 r1

**裁决者：** agent-ticket112-verify-r1（**非实现者**；实现腿为 `agent-ticket112` / `agent-ticket112b`）
**裁决时刻：** 2026-11-09 19:0x +08
**判据锚点：** HEAD `18ce9f5bc8d8a96ad0465d64916c1ce8b21c2c11`（`dev`）
**裁决对象：** 票 112 的 AC#1–AC#5 五格
**结论：** **五格全部成立，票可结案。**

本表只做判，不改任何 AC 原文。所有红/绿读数均**现跑**，取自导出树
`.scratch/wisp/probes/112/r1/clone`（`git clone --shared`，`-c core.autocrlf=false -c core.eol=lf`；
`*.go text eol=lf` 由 `.gitattributes` 强制，实测 `tree_ownership_112_windows_test.go` CRLF 对 **0**）与
**CI 原始日志**（`gh api .../jobs/<id>/logs`，认字节数与正向形状，不认"没出现我背的词"）。

---

## 0. 起手三核（先核派单前提，再判）

| 核 | 判据 | 读数 |
|---|---|---|
| ① 派单点名的串/行还在不在 | 现跑 `grep`/`sed` | 票面五枚 AC 框**全部未勾**（`- [ ]` x5），与派单"五枚未勾"一致；三条红的用例名在 HEAD 树里都在 |
| ② 引文出处文件/行 | 逐条比对源文件 | 票面引的三条断言原文，**逐字命中** `job 106319703680` 日志（见 §1） |
| ③ 红/绿是不是别人在飞件造成的 | 导出树跑门禁 + `git status` | 共享树 `git status --porcelain -- internal/winsec/ internal/risk/ scripts/ .github/` **空**（零脏）；导出树 `HEAD` 与共享树同 SHA；门禁在导出树全绿 |

**派单前提核验结论：** 派单说"票 112 五枚未勾、待非实现者复判"，**与盘上现状一致，前提未过期**。
但派单**未提到**的一件事必须补记：**AC#5 的 CI 读数已经在盘外存在**（见 §5），
实现腿当时写"欠编排者 push 后复跑"的那一格，在 `dev` 上已经由后续推送自然补齐。

---

## 1. AC#1 — 三条各自：CI 失败断言原文 + 本机读数 + 定因 + 环境差异/判定缺陷

**取数判据（正向形状，四项全中才算读到）：**
`http=200`、日志 **214419 字节**、首行是真时间戳、且本步自己的 `##[group]Run bash scripts/winsec-tests.sh` 在文件里。
落盘：`.scratch/wisp/probes/112/r1/logs/ac1-old-job-106319703680.log`（214419 字节，与票面记的字节数**一致**）。

票面自带的四数（`winsec-tests.sh: four numbers`，日志第 1359 行）逐字：
`=== RUN=71  --- PASS=32  --- FAIL=3  --- SKIP=0`；同 job 第 1353 行 `FAIL	github.com/CarlosShao/wisp/internal/winsec	13.707s`。

### ① `TestSealNarrowsAndNamesThePrincipalItRemovedBySID`

**CI 原文（日志第 261 行，`job 106319703680`）逐字：**

```
2026-09-21T11:46:30.4011140Z     private_set_sid_windows_test.go:268: the notice named the cleared principal by spelling only, not by the resolved SID it holds: cleared="S-1-1-0(A;OICI;FA;;;WD)", want S-1-5-21-3699639565-2515463329-295617607-500 in it (root DACL now [0/0x0=S-1-5-18 0/0xb=S-1-5-18 0/0x0=S-1-5-32-544 0/0xb=S-1-5-32-544 0/0x0=S-1-5-21-3699639565-2515463329-295617607-500 0/0xb=S-1-5-21-3699639565-2515463329-295617607-500])
```

**本机对应读数：** 同一用例在 HEAD 导出树 `-count=1 -v` = `--- PASS`（0.15s）。
**定因：** 夹具把"`LA` 一定是外人"当成机器无关事实；Windows job 以**内置 Administrator（RID 500）**登录，
票 106 的私有集按 SID 判定后**含 LA**，于是密封保留它、通知里当然没有它。
而报错文本自带的 `root DACL now [...]` 证明**密封做对了**（LA 那两个 ACE 还在，Everyone 被清掉）。
**判定：环境差异暴露的判定缺陷——缺陷在测试夹具，不在产品判定。**

### ② `TestC26PipelineIsWiredIntoWinsec`（0.00s 就红）

**CI 原文（日志第 666 行）逐字：**

```
2026-09-21T11:46:30.4598381Z     resolve_windows_test.go:30: no C26 pipeline installed in winsec: internal/risk/winsec_c26.go's init did not run
```

同 job 第 685 / 687 行是它下游的两条同名 ERROR 读数（`resolve_windows_test.go:90` / `:97`，
逐字 `C26 is not installed, so this leg measured the fallback instead`）。

**上游同因（日志第 204 行，守卫自己写的 ERROR）逐字：**

```
2026-09-21T11:46:30.3942522Z 2026/09/21 11:46:16 ERROR winsec: refusing to install a path resolver into the sealing seam resolver=risk.c26Pipeline reason="it answered \"C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\wisp-108-tree-ownership-probe\" with \"C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\wisp-108-tree-ownership-probe\", which is not inside the tree \"C:\\Users\\runneradmin\\AppData\\Local\\Temp\" it answered for that path's own parent \"C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\": the seam may not be used to move a seal into another tree" consequence="the incumbent resolver, or the built-in floor, stays in place"
```

**本机对应读数：** HEAD 导出树同一用例 `--- PASS`（0.15s）。
**定因：** 守卫问"已存在的父目录"（解析器给**真实长名**，8.3 展开）与"不存在的子路径"（解析器按契约给**词法原样**，
仍是 `RUNNER~1`），再**逐段比拼写** ⇒ 把同一棵树两种拼写读成"搬树"，拒绝装 C26，整包退到 built-in floor。
只有 profile 目录长于 8 字符的机器才看得见（runner 是 `runneradmin`，本机是 `swq`）。
**判定：判定缺陷——守卫拿拼写当树。"只有 runner 那种机器看得见"是暴露条件，不是豁免理由。**

### ③ `TestAC3JunctionInputIsRefusedNotSealed`（两个子形状）

**CI 原文（日志第 713–715 行）逐字：**

```
2026-09-21T11:46:30.4661415Z --- FAIL: TestAC3JunctionInputIsRefusedNotSealed (0.09s)
2026-09-21T11:46:30.4662413Z     --- FAIL: TestAC3JunctionInputIsRefusedNotSealed/existing_directory_behind_the_link (0.00s)
2026-09-21T11:46:30.4663805Z     --- FAIL: TestAC3JunctionInputIsRefusedNotSealed/missing_directory_under_the_link (0.00s)
```

第三腿 `with_only_the_built-in_verifier` 同一块里是 `--- PASS (0.03s)`（日志第 716 行）。
**本机对应读数：** HEAD 导出树 `--- PASS`（2.26s），三腿全 PASS。
**定因：** 与 ② **同一上游**（C26 没装上 ⇒ 走 floor ⇒ 两条腿测的是 fallback 而非真接线）。
**判定：判定缺陷，同 ② 一根。**

**AC#1 判：成立。** 三条**分开**给了原文（未打包），且"环境差异 vs 判定缺陷"逐条点名：
三条全部是**判定缺陷**，环境只是暴露条件；①的缺陷在夹具，②③的缺陷在守卫。

---

## 2. AC#2 — 修的方向只能"更严或更响亮"：逐条核

判据不是"这一步绿了"，而是**每一处改动是否让判据变松**。逐条现读 diff（`a50560793`）核：

| 改动 | 方向 | 判据 |
|---|---|---|
| `resolve.go` 守卫 | **加第二见证，不减第一见证** | 原包含判据 `answerInsideTree(childAns, parentAns)` 保留为**第一件事**，通过即 `return ""`；只有它不过时才追问 `anchor = filepath.Dir(childAns)`。**没有删掉任何一条旧判据**，也没有把比较改成大小写无关的"看起来像同一棵树"（`sameTree` 逐段比，折叠规则与 `answerInsideTree` 同一把尺，未新起正规化器 ⇒ 不撞 D22 ban #2）。 |
| 守卫的松紧 | **净收紧** | 诚实解析器（两种拼写同一对象）过；**搬树伪件仍被拒**，且拒绝理由**多了一条**（"nor inside the tree that answer names once the candidate is asked about it directly"）。拒绝/如实记账搬家**依旧算窄**（`return ""`），没被改成"必须报错"。 |
| 夹具候选名单 | **净加强** | 从硬编码 `[everyone, la]` 改成**本机现量**：`foreign = [everyone, guests]`，`member` 按 `inPrivateSet(set, la)` 现测决定。**新增两条断言**：`member` 不得出现在 cleared 里（`:348`）、`member` 必须仍站在 root 上（`:352`）。旧形状一条都没少。 |
| 种子自检 | **新增（更响亮）** | `for _, sid := range plant { if !standsOn(t, root, sid) { t.Fatalf("the fixture planted nothing to judge for %s") } }` —— 种不进去就红，不许静默测空。 |
| 新增仪器 | **新增（更响亮）** | `tree_ownership_112_windows_test.go` 先读仪器自己：`strings.EqualFold(parentAns, short)` 或子答案被展开 ⇒ `t.Fatalf("the instrument measured nothing")`。测不到形状就红，不许"因为没测到所以过"。 |

**逐条核"四禁"（AC#2 点名禁的四样）：**

- **放宽断言：** 无。反证 = §4 的 MUT1/MUT3/MUT5 三发变异，**每一发都把新用例打红**；放宽过的断言打不红。
- **改阈值：** 无。`git diff --stat 2ab6519..HEAD -- internal/winsec/ internal/risk/ scripts/winsec-tests.sh .github/workflows/ci.yml` **为空**（本票改的四处全部在更早的 `a50560793` 里，且不含阈值文件）。
- **加 `t.Skip`：** 本票**没有新增** `t.Skip`。树里唯一的 skip 是 `tree_ownership_112_windows_test.go:122` 的**子测试能力门**（卷不生成 8.3 名 ⇒ 形状造不出来），
  这正是 AC#2 允许的形状（"若某条确实依赖本机才有的条件，正解是把它做成平台/能力门并说明为什么"），且
  **放在子测试内**（顶层 `--- SKIP` 才是致命的，票 71 AC#3），注释写明了理由。**runner 上它没 skip，真跑了**（§5 读数）。
- **加 `//go:build`：** 新文件带 `//go:build windows`，但 8.3 短名是 **Windows 文件系统属性、POSIX 没有对应物**可种；
  同一守卫的 POSIX 腿在 `ancestor_separator_108_other_test.go`（票 113 的文件），**没被删**。

**AC#2 判：成立。** 方向确实是"更严或更响亮"，且"更响亮"有三处新增（种子自检、仪器自检、守卫的第二条拒绝理由）。

---

## 3. AC#3 — 复现仪器：没有那些条件的机器上也能验证同一性质

**仪器：** `internal/winsec/tree_ownership_112_windows_test.go`（`package winsec_test`），
用**真** C26 适配器 `c26Pipeline112`（三腿照 `internal/risk/winsec_c26.go` 的生产接线）+ 测试专用口
`winsec.TreeOwnershipProbeForTest`。

**造法（票 106 的形状）：** 在 `t.TempDir()` 下建一个长名目录 `wisp11283parentdir` ⇒ 向 OS 要它的 **8.3 短名**
（`windows.GetShortPathName`，不手拼）⇒ 在**短拼写**下挂一条**不存在**的子路径。
这正是 runner 的不对称：已存在对象被展开成真实长名，不存在的叶子保持词法短名。

**先读仪器自己（关键，防"没测到所以过"）：**
`:141` 若父答案**没被展开**（`parentAns == short`）⇒ `t.Fatalf("the instrument measured nothing ...")`；
`:144` 若子答案**被展开**了 ⇒ 同样 Fatalf。即：这台机器造不出不对称，它必须**说"我什么都没测到"**，不许静默通过。

**本机现跑（HEAD 导出树，`-count=1 -v`）：**

```
=== RUN   TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent
=== RUN   TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent/on_a_volume_that_generates_8.3_names
--- PASS: TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent (0.01s)
```

子测试**真跑了**（不是 SKIP）⇒ 本机卷生成 8.3 名，不对称种成功。

**修前必须红（AC#3 的硬要求）——本裁决者独立复现，不引实现腿的读数：** 见 §4 MUT1。
MUT1 把守卫退回旧单见证后，同一仪器给出：

```
tree_ownership_112_windows_test.go:149: AC#1 RED: the guard still refuses the honest pipeline on the runner's own shape: it answered "C:\\Users\\swq\\TESTTR~1\\001\\WISP11~1\\wisp-112-missing-child" with "C:\\Users\\swq\\TESTTR~1\\001\\WISP11~1\\wisp-112-missing-child", which is not inside the tree "C:\\Users\\swq\\TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParenton_1881802828\\001\\wisp11283parentdir" it answered for that path's own parent "C:\\Users\\swq\\TESTTR~1\\001\\WISP11~1": the seam may not be used to move a seal into another tree
--- FAIL: TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent (0.01s)
```

**AC#3 判：成立。** 仪器在本机（`swq`，非 runner）上把 runner 的形状**显式种**出来并验证同一性质；
修前红、修后绿，两头都现跑过。

---

## 4. AC#4 — 变异：被修的判定退回旧实现 ⇒ 新用例必须红

三发变异，全部在**导出树的另一份克隆** `.scratch/wisp/probes/112/r1/mut`（同 SHA `18ce9f5bc`）里做；
每发**同链打印被改后整行**，先 `go build` rc=0（编译失败不算变异）再 `go test`。

### MUT1 — 守卫退回 ticket 108 的单一拼写见证

改动（`internal/winsec/resolve.go:394-397`，被改后整行同链打印）：

```
394:	// MUT1 (ticket 112 AC#4): the guard is reverted to ticket 108's single
395-	// spelling witness. Everything above this point is untouched.
396-	return fmt.Sprintf("it answered %q with %q, which is not inside the tree %q it answered for that path's own parent %q: the seam may not be used to move a seal into another tree",
397-		child, childAns, parentAns, parent)
```

`go build ./internal/winsec/` rc=0。`go test -run 'TestTreeOwnershipProbe|TestC26Pipeline|TestAC3Junction|TestAC1SeamRejects|TestAC1SeamCannot' -count=1 -v`：

```
--- PASS: TestC26PipelineIsWiredIntoWinsec (0.08s)
--- PASS: TestAC3JunctionInputIsRefusedNotSealed (3.18s)
--- PASS: TestAC1SeamCannotBeFreedThenGivenATreeMovingFake (1.17s)
--- PASS: TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone (1.48s)
--- FAIL: TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent (0.01s)
FAIL	github.com/CarlosShao/wisp/internal/winsec	6.020s
```

**红在守卫，不是红在噪声**：同一次运行里另四条既有 seam 用例**仍绿**，唯一红的是新仪器的 `:149 AC#1 RED`。

### MUT3 — 通知只写 OS 渲染（`token := ace.text`）

改动（`internal/winsec/winsec_windows.go:181-183`）：

```
181:		// MUT3 (ticket 112 AC#4): the notice is named by the OS rendering only,
182:		// which is the shape ticket 106 removed.
183:		token := ace.text
```

`go build` rc=0。`go test -run TestSealNarrowsAndNamesThePrincipalItRemovedBySID -count=1 -v`：

```
private_set_sid_windows_test.go:340: the notice named the cleared principal by spelling only, not by the resolved SID it holds: cleared="A;OICI;FA;;;LA A;OICI;FA;;;BG A;OICI;FA;;;WD", want S-1-1-0 in it (...)
--- FAIL: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (0.28s)
```

**红在"按 SID 命名"这条腿**，且本机**真复现了"只剩拼写"的形状**（`cleared` 里三个都是名字不是 SID）。

### MUT5 — 去掉按 SID 的私有集成员豁免（票面自认欠着的那一发）

**为什么补这一发：** 实现腿在进度日志里**诚实交代**过——"集合内主体不得被报告"这条新腿
**没有**被它那四发变异单独打红过，理由是"要 runner 那个条件（job 以 RID 500 登录），本地无法造"。
本裁决者判：**这个理由不成立**。该断言判的是"通知里不得含 member 的 SID"，与"LA 是不是当前用户"无关——
本机 `member = BUILTIN\Administrators`（S-1-5-32-544），把它从豁免里拿掉就必然被报告。**能造，就造。**

改动（`internal/winsec/winsec_windows.go:174-177`）：

```
174:		// MUT5 (ticket 112 AC#4): the private-set membership exemption is removed,
175:		// so a principal inside the set is reported as cleared. This is the leg the
176:		// ticket's own progress log said no mutation had ever reddened.
177:		if false && ace.grant && set[ace.trustee] {
```

`go build` rc=0。`go test -run TestSealNarrowsAndNamesThePrincipalItRemovedBySID -count=1 -v`：

```
private_set_sid_windows_test.go:348: the notice named a private-set member as cleared: member=S-1-5-32-544 cleared="S-1-5-18(A;;FA;;;SY) S-1-5-32-544(A;OICI;FA;;;BA) ..." (set [S-1-5-18 S-1-5-32-544 S-1-5-21-...-1001], trustees [...])
--- FAIL: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (0.26s)
```

**红在 `:348`，正是那条"至今没有结论"的腿** ⇒ 票面那一格的空白由本裁决补齐：该腿**有牙**。

**AC#4 判：成立（且比实现腿自认的更完整）。** 三发变异分别打红守卫腿、SID 命名腿、成员豁免腿；
每发都 `go build` rc=0 在前，红名逐条点名，红在守卫而非噪声（MUT1 的同链对照四条绿）。

---

## 5. AC#5 — 结案条件：这一步在 CI 上给出过 success

**主读数（HEAD 当枚，最硬的一枚）：**

| 项 | 值 |
|---|---|
| run id | **`42605029431`** |
| head | `18ce9f5bc8d8a96ad0465d64916c1ce8b21c2c11`（= 本裁决锚点 = 当前 `dev` HEAD） |
| job id | **`129118693717`**（`test-windows`） |
| **step 号** | **4** — `Windows ACL sealing gate (internal/winsec's own tests, ticket 110)` |
| step 结论 | **`success`** |
| 日志 | 1593038 字节，落盘 `.scratch/wisp/probes/112/r1/logs/step4-success-129118693717.log` |

**这一步自己的四数（日志第 1888 行）逐字：**
`winsec-tests.sh: four numbers (all from -v output): === RUN=129  --- PASS=78  --- FAIL=0  --- SKIP=0`
**结果行（第 1887 行）逐字：** `winsec-tests.sh: winsec result line: ok  	github.com/CarlosShao/wisp/internal/winsec	7.887s`

**三条红的用例在这枚 run 上的逐字读数：**

```
--- PASS: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (0.03s)
--- PASS: TestC26PipelineIsWiredIntoWinsec (0.01s)
--- PASS: TestAC3JunctionInputIsRefusedNotSealed (0.08s)
    --- PASS: TestAC3JunctionInputIsRefusedNotSealed/existing_directory_behind_the_link (0.02s)
    --- PASS: TestAC3JunctionInputIsRefusedNotSealed/missing_directory_under_the_link (0.00s)
    --- PASS: TestAC3JunctionInputIsRefusedNotSealed/with_only_the_built-in_verifier (0.02s)
--- PASS: TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent (0.00s)
    --- PASS: TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent/on_a_volume_that_generates_8.3_names (0.00s)
```

**②那枚"集合内主体必须被保留且不得被报告"的新腿，在 runner 上真跑到并过了**（日志第 554 行）：
`private_set_sid_windows_test.go:314: planted as strangers: [S-1-1-0 S-1-5-32-546]; planted as a private-set member: S-1-5-21-3289802278-1561746829-1498425070-500 (LA=..., in set: true, set=[S-1-5-18 S-1-5-32-544 S-1-5-21-...-500])`
⇒ 选桶那步打印了 runner 形状（LA 在集合内），随后 `:348`/`:352` 两条断言**没有被 Fatalf 打断**，
用例整体 `--- PASS`。实现腿当时写"这条腿的结案证据仍欠着"，**现在补齐了**。

**旁证（连续五枚 run 的 step 4 全绿，排除"恰好那一枚"）：**

| run | job | step 4 |
|---|---|---|
| `42565983659` | `128992481751` | success |
| `42572132704` | `129011661041` | success |
| `42572328515` | `129012287773` | success |
| `42578363048` | `129031082782` | success |
| `42578644540` | `129031947465` | success |

**一个必须说清的读数：** run `42605029431` 的**整体 conclusion 是 `failure`**，
但红的是 **step 7（cmd/wisp CLI tests）与 step 10（Portable windows tests）**，
**不是 step 4**；step 5/6/8/9/11 均 `success`。AC#5 判的是**这一步**的结论（票 111 AC#6 已把判据改成 per-step），
故"job 红"不抵消"step 4 绿"。step 7/10 的红属别的票，**不在本票射程**。

**AC#5 判：成立。** 三元组齐全（run / job / step 号），且是本裁决锚点那一枚 HEAD。

---

## 6. 门禁（本裁决零产码：只写 `docs/evidence/s1/` 与票面）

全部在导出树 `.scratch/wisp/probes/112/r1/clone` 跑（`unset GOTMPDIR TMPDIR TEMP TMP`，
`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`，**未设 `CGO_ENABLED=0`**）：

| 门禁 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/` | 输出**空**，rc=0 |
| build | `go build ./...` | rc=0 |
| vet | `go vet ./internal/winsec/` | rc=0 |
| d22scan | `sh scripts/d22scan.sh` | **clean**，rc=0；`ban #8 internal/ = 972 Go files, comments and _test.go included` |
| 路径长度 | `bash scripts/check-path-length-budget.sh` | **VERDICT GREEN**；`tracked paths=19225  over-budget=63  covered by roster=63  not in roster=0`；rc=0 |

---

## 7. 逐格判决

| AC | 判决 | 一句话判据 |
|---|---|---|
| AC#1 | **成立** | 三条**分开**给了 CI 逐字原文（含上游 ERROR），本机对照 + 定因 + "环境差异 vs 判定缺陷"逐条点名；三条全是判定缺陷 |
| AC#2 | **成立** | 四禁逐条核过：零放宽、零阈值、零新增 skip、build tag 是 8.3 的正当平台门；且三处"更响亮"新增 |
| AC#3 | **成立** | 仪器在本机显式种 runner 形状、先自检"测到了没有"；修前红（MUT1 独立复现）修后绿 |
| AC#4 | **成立** | 三发变异分别打红守卫/SID 命名/成员豁免三腿，含实现腿自认从未打红过的那一发 |
| AC#5 | **成立** | run `42605029431` / job `129118693717` / **step 4** = `success`，四数 `RUN=129 PASS=78 FAIL=0 SKIP=0`，锚点即当前 HEAD |

**总判：五格全成立，票 112 可结案（改名 `-done`）。**

---

## 8. 与实现腿自述的两处出入（如实登记，不改 AC）

1. **"AC#5 欠编排者 push 后复跑"已不再欠。** 实现腿写这句时确实没推（代理不 push 是对的）；
   但后续推送已把它带上 `dev`，本裁决取到 step 4 = `success` 五枚连绿。**这一格从"欠"变成"已达成"**，
   不需要再等一次 push。
2. **"成员豁免腿没有变异打红"这个自认已经可以销掉。** 实现腿说它"要 runner 那个条件，本地无法造"；
   本裁决用 MUT5 证明**能造**（member 取本机的 `S-1-5-32-544`，与 LA 无关），并打红了 `:348`。
   这是对实现腿**过谦**的一处更正，不是对它的指控。

**未改一字：** 票面 AC 原文、阈值、golden、`internal/winsec/**`、`internal/risk/**`、`scripts/winsec-tests.sh`、`ci.yml`。
本裁决**零产码**，只新增本文件 + 票面翻勾与 append。
