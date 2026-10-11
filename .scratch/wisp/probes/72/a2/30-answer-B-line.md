# 72-a2 — 必答 B："⛔ 把甲族的红用放宽断言消掉"那条线在哪

## B.0 一把可执行的尺（先给尺，再逐枚裁）

> **改完之后，这一枚"能失败的原因集合"变⛔ 变小？**
> 变小 ⇒ 那就是放宽判据（AGENTS §1.1"⛔ 为变绿放宽任何断言"＋票 72 `AC#5`"⛔ 动仪器"）。
> 一⛔ 变，只是把**比较的两端换成同一把折法产出的东西** ⇒ 那才是改比较写法。

这一把尺⛔ 是口头的：同一枚用例改完，把修复挪开（`formsOf` 退回 raw / `overrideApplies` 退回遍历全部形）**必须还是那一批红**，红句⛔ 换位置。⇒ 这就是票 307 `AC#2` 那"三缺一"的同一形，写腿照抄即可。

## B.1 逐枚：该动哪一侧，⛔ 该动哪一侧

### 形 I（#4 `junction:135`／#5 `:164`／#6 `:182`）——`want` 字面量侧

**该动的一侧**＝那一枚被叫作 `long` 的**期望串**（`junction:125`/`:150`/`:176` 三处 `long := filepath.Join(home, ".git-credentials")`）。
**⛔ 该动的一侧**＝`Classify(res.Canonical) != ClassA`（`:137`）与 `gateFor(t, res.Canonical).Allow`（`:140`/`:166`/`:184`）——这三行是这一族唯一的安全判级，**一字节⛔ 碰**（票 72 `AC#5` 逐字点名"`ClassA` 的期望值、A/B 表内容、豁免语义一律⛔ 改"）。

**最小诚实落点（候选甲，本机可闭，⛔ 依赖机器形状）**：期望串⛔ 再由 `t.TempDir()` 拼，改成**过一次 C26 再拼**——
`want := filepath.Join(mustResolveDir(t, home), ".git-credentials")`，其中 `mustResolveDir` ＝ `res, err := Resolve(dir, nil); …; return res.Canonical`（对**存在的目录**，`Resolve` 交回的正是句柄展开形：本腿件 `logs/aliased-7named.txt` 第 12/15/18 行的 `got` 三枚逐字都是那个形）。
⇒ 比较仍然是**一次逐字节 `EqualFold`**，判据集合⛔ 变小；只是把"人类的临时目录那一形"换成"OS 担保的那一形"。
⚠ 这一形**只让这 3 枚在别名台面上跑对**——在⛔ 别名的机器上它们照旧绿（本来也绿），所以它**⛔ 满足票 307 `AC#3` 那句"夹具⛔ 依赖机器形状"**。

**最小诚实落点（候选乙，多造 1 枚夹具，两台面同色）**：把"输入用短形"这件事自己种出来——同包已有**现成的两枚手**，⛔ 新造第三枚：
- `syncdirs_redteam_windows_test.go:29 tryShortPath(p) (string, bool)`（`windows.GetShortPathName`，注释逐字"Spellings are always read back from the OS — never typed by hand"）；
- `syncdirs_alias_307_windows_test.go:38 alias307Pair(t, base, name) (long, short string)`（票 307 刚落地的；取⛔ 到短形**即 `t.Fatal`**，⛔ `t.Skip`，票 115 定式）。
⇒ 落法＝`long, short := alias307Pair(t, t.TempDir(), "ThirtySevenAlias…")`，然后 **`Resolve(short)` 的期望＝`Resolve(long).Canonical`**（两侧都出自 C26）。
⚠ 两处**必带的具体坑**，⛔ 让写腿踩了才发现：
1. `alias307Pair` 交回的 `long` ＝ `filepath.Join(base, name)`——**它⛔ 把 `base` 折长**！在别名台面上那枚 `long` 自身还带着 `RISKAT~1\RUNNER~1`。⇒ 拿它的 `long` 当期望串**还会红同一行**；期望必须走 `Resolve(...).Canonical`。
2. 它住在带 `//go:build windows` 的文件里（`syncdirs_alias_307_windows_test.go:1`），与甲族那两枚文件同 tag ⇒ 同包可直接调；⛔ 需新文件。
⇒ **本腿推荐＝候选甲落地 ＋ 候选乙补一枚"自己种得出别名世界"的同族夹具**，因为票 307 `AC#3` 已经把那条判据写死了（"⛔ 依赖机器形状的夹具"），甲族吃同一条判据⛔ 讲得通。

**★哪一枚如果这样改就变成假绿**＝**#4**。若把 `:134` 那句比较改成"canonical ⛔⛔ 等于 short 也算过"（例如换成 `strings.Contains`、或"在同一棵子树下"、或干脆删掉那三行只留 `Classify`），那这一枚就不再钉"**8.3 短名被展开**"这件事——而 `short name did not expand to long path` 那一句正是票 18 红队 **Case 2 的全部内容**（票面 `:13-16` 强调"一条都⛔ 是占位"）。⛔ 删⛔ 降级＝假绿。#5/#6 同形（它们的 `:163`/`:181` 是各自唯一"折法对上了"的钉，之后只剩"⛔ 放行"一枚）。

### 形 II（#1 `anchor_spelling:57`）——前提侧 `t.Fatalf`

**该动的一侧**＝`if !strings.EqualFold(res.Canonical, p)`（`:56`）里那个 `p`（`:50` 由 `t.TempDir()` 一路 Join 而来）。
**⛔ 该动的一侧**＝`:59`（`Classify(res.Canonical) != ClassA`）与 `:63`（`Gate(...).Allow || .NeedL2`）——这一枚的**真判级**只在它那两行里。

**最小诚实落点**＝两条都行，判据集合都⛔ 变：
(a) 把 `profile` 先过一次 `Resolve` 再 Join 叶子（同候选甲那形）；
(b) ★更贴原意的一形：**判级的入参换成 `alias` 那一形**（这一枚的主题本来就是"USERPROFILE 是 junction 拼法"），即 `Classify(filepath.Join(alias, rel))` 也要 `ClassA`，同时保留 `Classify(res.Canonical)`。⇒ 它把票面 `:69-72` 那句"退 `formsOf` 就红"的变异靶**变宽而⛔ 变窄**。

**★哪一枚如果这样改就变成假绿**＝**#1**，而且形得很：把 `:57` 那枚 `t.Fatalf` 降级成 `t.Errorf` 或 `t.Logf` ⇒ 本枚立刻在别名台面绿（因为 `:59`/`:63` 的判级确实是 ClassA/确实⛔ 放行——本腿 `logs/aliased-all16-bothfiles.txt` 里同族的 `TestClassifySpellingInvariance` 绿就是侧证）。**但它会同时把票 72 `AC#3①` 那两发变异读数的靶子搬走**（票面 `:70-72`、log `:139-141`、`:146-150` 三处逐枚点名这一枚"退 `formsOf` 必红"）⇒ 那批证据当场过期。⛔ 是"⛔ 能改这一枚"，是**改完必须重跑那两形变异**（见 `40-nails` 第 4 条）。

### 形 III（#2 `anchor_spelling:140`）——`want 1` 计数面

**该动的一侧**＝输入名册里那一枚 `filepath.Join(home, "not-created-yet", "deep.txt")`（`:132`）——把它改成**在已展开的根之下重新拼**（`want := filepath.Join(mustResolveDir(t, home), "not-created-yet", "deep.txt")`），使输入本身只剩一种拼法。断言体 `len(f.spellings()) != 1`（`:139`）与 `!f.certain()`（`:142`）**一字节⛔ 动**。
**⛔ 该动的一侧**＝把 `want 1` 写成 `<= 2`、或把 `:132` 那枚从名册里**删掉**。

**★哪一枚如果这样改就变成假绿**＝**#2**，两条都是：
- 删掉那枚缺失叶子 ⇒ 票面 `:96-98` 那句"Resolve 交回的 canonical 只产**一种**比较形且 `certain()`"**悄悄缩成"只对存在的路径成立"**，而 AC#5 拿它当"爆炸半径⛔ 越界"的唯一仪器（票面 `:96-98` 原话"对没有第二种拼写的路径，本票一个字都不改变判定这句话现在是可跑的"）；
- 改阈值 ⇒ 这一枚彻底失去信号（它唯一的判据就是那个 `1`）。
⇒ 正解只有"把输入的拼法换成已展开那一形"一条。⚠ 还有第三条**⛔ 属票 72 射程**的岔路要具名堵住：**去改 `Resolve` 让它对缺失叶子也展开祖先**——那是产码改动，会波及 `Actable()`、票 307 刚落地的 `syncdirs.go:432-442`、`internal/tools/paths.go`，且正面撞上票 307 追加节 `:116` 那句裁语"**存值一字节⛔ 改**"的同族精神。谁提这条，谁就该被要求先报编排者判。

### 形 IV（#3 `anchor_spelling:238`／#7 `junction:292`）——把短拼当"OS 担保拼写"用

**该动的一侧**＝**测试递给 `Gate` 的那两样东西**：入参拼法 ＋ override key。两者都换成"OS 担保那一形"，即 `c := mustResolveDir(t, filepath.Dir(p)) + leaf`（或直接 `formsOf(p).real` 的产物），再 `Gate(c, map[string]bool{normPath(c): true})`。
**⛔ 该动的一侧**＝★**`blacklist.go:115-126 overrideApplies` 一字节⛔ 碰**。

★**这就是那条线本身，具名钉死**：
```go
if cand.raw != cand.real {
    return false // spelled differently from its own real path: confirm again
}
```
把这一支改成"遍历全部形式都算命中"，是**全票族最省事的一枚绿**（一发就同时治好 #3 与 #7），也是票 72 已经**走过一遍并且作废掉**的那一步。三处原话，逐字：
- 票面 `:95`＝`插单前我在 AC#5 写的那段"豁免查表改成全部拼写形"**作废**，那正是被指出来的放宽。`
- 票面 `:239`＝`判定侧漏一次 = A 静默降 B（放行风险由别人承担），批准侧多放一次 = 8.3 别名到另一个目录时放过同名文件；所以**拒绝可以宽，放行必须窄**，宁可让人多点一次。`
- 票面 `:171-173` 记那一发变异**实测过**：把 `overrideApplies` 改回遍历全部形式 ⇒ `TestOverrideOnlyAcceptsTheResolvedForm` 断言 (1)(2) 转红，报的正是 `{Allow:true … B-list single-file override … id_*}`。
⇒ **具名判语：#3/#7 如果靠动 `overrideApplies` 变绿，就是"为变绿放宽断言"，AGENTS §1.1 与票 72 `AC#5` 双重直接违反；且它现在⛔ 只是测试问题——生产 caller 已在（`internal/tools/mode.go:98` `risk.Gate(c, overrides)`，`c` 出自 `:93 b.paths.Canonicalize(p)`），动了它＝同时改生产分档（`BlacklistNote.AlreadyUnlocked` ⇄ `Unlockable`）。**

**恒真风险具名**：这一族的两枚**正控**（`anchor_spelling:237-239`、`junction:291-293`）**⛔ 能删也不能反转**——它们是"这条窄规则是规则、⛔ 是 blanket deny"的唯一反证（同族另一枚 `:192-197` 的注释逐字"control, so this is a rule and not a blanket deny"）。删掉它们，`overrideApplies` 就算恒 `return false` 也能全绿 ⇒ 那才是真正的假绿形。

## B.2 ★一条⛔ 在"比较写法"⛔ 在"放宽判据"上、但一定会被想到的第三条路——先把名字报出来

**把 `t.TempDir()` 换成手写的长形常量**（或反过来把期望"对齐"成短名）。
票 72 的 R17 段（票面 `:200-205`）把这一条列为**明令禁止的第 1 种解法**，逐字＝
> `把测试里的锚点也换成短名、或比较前对两侧做大小写/短名"对齐"的局部补丁——那是**在测试里复刻 bug**`

⇒ 同族禁止的第 2 条（`if isCI`／按环境变量放宽）与第 3 条（`t.Skip` 或 `//go:build` 挡出 `test-windows`）也一并写在票面 `:202-205`，本腿⛔ 新增任何一条。
⚠ 特别点名：**⛔ 用 `t.Skip` 或加/改 build tag 来消这 7 枚的红**——票面 `:205` 逐字"路径拼写⛔ 是平台 API 限制，是我们的 bug（A40③）"，而票 307 追加节 `:133` 已把"给那批加 tag 会把它们从 core 档分母里摘掉"写成先例（撞票 301 的包级 guard D）。

## B.3 验收形状（写给下一枚落地腿，⛔ 由它自勾）

1. **两台面各一发、名册同一集合**：别名台面＝本腿 `logs/aliased-all16-bothfiles.txt` 那把（16 枚／PASS=9／FAIL=7／SKIP=0 → 修后应＝**PASS=16／FAIL=0**）；默认台面＝`logs/default-7named.txt` 那把（**颜色⛔ 变**，仍 7/0/0）。
2. ⛔ 拿"默认台面全绿"当凭据（票 307 现量与本腿同形：本机 `%TEMP%` SHORT==LONG，逐字量过＝`C:\Users\swq\AppData\Local\Temp`）。
3. 每发**先断言名册行数**再谈颜色（本腿 `logs/t03-broken-surface-note.txt` 记了一发"PASS=0 却⛔ 全红"的 panic 假读数）。
4. 变异三形照抄票 307 `AC#2`：(a) 把 `formsOf` 退回 raw ⇒ #1/#2/#3 那族必红且**红位⛔ 换**；(b) `cmp` 逐字节还原 ⇒ 复绿；(c) 恒真形 ⇒ 指名用例仍响。⚠ 本腿按结构预判**⛔ 会沉默的那两枚**＝#1 的 `:57` 与 #2 的 `:140`（它们红的是前提／形状，⛔ 是判级），写腿要么给它们另配一枚能响的钉，要么具名交"本格⛔ 成立"。
