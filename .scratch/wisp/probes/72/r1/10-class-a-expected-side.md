# 72-r1 甲组 7 枚逐枚：坏形状 ／ 新期望侧逐字 ／ ⛔ 撞三种禁止解法 ／ 现在钉的是判定还是拼法

写面＝`internal/risk/pathresolver_anchor_spelling_windows_test.go`（#1/#2/#3）＋`internal/risk/pathresolver_junction_windows_test.go`（#4/#5/#6/#7）。产码⛔ 动（本腿所有 commit 的 `--numstat` 里非 `_test.go` 命中＝0 行，见 `40-two-surface-readings.md` 的逐笔 hunk 尺）。

票面不变式（`72-atble-classification-runner.md:215-222`）：**安全分类的结果不得依赖路径的拼写形式**。七枚红的**期望侧**全部写成"字符串等于我从环境里拿来的那一串"（`t.TempDir()`／`tmpHome(t)` 的返回拼写），而 `got` 那侧是解析器**干对了**的展开长形 ⇒ 修＝改期望侧，⛔ 改产码。

新引入的共用夹具（定义在 `anchor_spelling` 文件、两文件同包可用）：`namesSameFile(a, b string) bool` ＝拿包内既有的形族（`formsOf`，即句柄管线）问"**这两串是不是同一个文件**"，⛔ 问"是不是同一个字符串"。它是票面 §3 第 2 条要求的"**同一棵树两种拼法必须给同一个判定**"那一支的形状，⛔ 手写常量、⛔ 大小写／短名"对齐"。

---

## #1 `TestClassifyAnchorSpellingIsNotVerdict`（原红位 `anchor_spelling:57`）

ⓐ**坏形状＝前提句**（拿测试自己 `filepath.Join(t.TempDir(), …)` 的产物与 `Resolve` 的返回逐字符串比，挂在 `t.Fatalf` 前提位 ⇒ 别名台面上 `Classify`／`Gate` 两枚真断言一行都没跑）。

ⓑ**新期望侧逐字**：
```go
		if !namesSameFile(res.Canonical, p) {
			t.Fatalf("Resolve(%q) canonical %q: test premise broken (candidate must be the real path)", p, res.Canonical)
		}
```
`t.Fatalf` 那枚**没降级、没删**（票面 §3 第 3 条禁线），报错句逐字保留 `test premise broken (candidate must be the real path)` ⇒ 内容锚仍可一把 `grep` 取到。`:59`（`Classify(res.Canonical) != ClassA`）与 `:63`（`Gate(...).Allow || .NeedL2`）一字未动，且现在**两枚台面上都真的执行**（改前别名台面它们在 `t.Fatalf` 之后⛔ 到达＝a2 census 第 #1 枚第 3 问的"仪器哑"）。

ⓒ⛔ 撞禁止解法：⛔ 把锚点也换成短名（解法 1，在测试里复刻 bug）——比较两侧交给 `formsOf`，⛔ 出现"对齐"；⛔ `if isCI`／env 分档（解法 2）——本枚⛔ 读任何环境变量、⛔ 任何台面开关，同一份码两支台面同一句判据；⛔ `t.Skip`／`//go:build`（解法 3）——⛔ 新增任何一枚（尺见 `40` 件）。

ⓓ现在钉的是**判定**：前提问"是不是同一个文件"，真断言问"判不判 A／可不可放行"。⛔ 钉拼法。

## #2 `TestCanonicalInputGainsNoSecondForm`（原红位 `:140`，`want 1` 计数）

ⓐ**坏形状＝计数**（`len(f.spellings()) != 1`）。`want 1` 在别名台面上是**错的**：那枚名册里 `not-created-yet/deep.txt` 的叶子⛔ 存在，`Resolve` 对它交回**词汇形**（祖先仍写 `RISKAT~1\RUNNER~1`），而 `formsOf` 沿祖先把 `real` 折成长形 ⇒ 两种拼法＝**2 枚形是 `formsOf` 的正确行为**（`pathresolver.go:329-331`：形集是并集，⛔ 许⛔ 扩爆炸半径那侧＝`classifyForms` 只会多判⛔ 少判）。

ⓑ**新期望侧逐字**（形 Ⓐ；整枚函数体见文件，关键三句）：
```go
		if len(fs) == 0 {
			t.Fatalf("canonical %q produced no comparison form at all", res.Canonical)
		}
		if !f.certain() {
			t.Errorf("canonical %q is not certain: root=%q real=%q partial=%v", res.Canonical, f.root, f.real, f.partial)
		}
		if got := Classify(fs[0]); got != tc.want {
			t.Errorf("canonical %q: form %q classified %v, want %v (this ticket may not move a verdict)", res.Canonical, fs[0], got, tc.want)
		}
		for _, sp := range fs[1:] {
			if got := Classify(sp); got != tc.want { … }
			base := Gate(fs[0], nil)
			d := Gate(sp, nil)
			if d.Allow != base.Allow || d.NeedL2 != base.NeedL2 || d.Class != base.Class { … }
		}
```
名册五枚**一枚没删**（含那枚尚未创建的叶子），每枚带上**绝对期望判定**（`.ssh\id_testkey`→`ClassA`、`.git-credentials`→`ClassA`、`plain.txt`／`home`／缺失叶子→`ClassNone`，逐枚与 `blacklist.go:213-232` 的规则号对过），`certain()` 那枚原样保留。

ⓒ⛔ 撞禁止解法：⛔ 读 env／`if isCI` 分档（票面 `:201` 逐字禁），本枚⛔ 任何台面分支——两支台面跑同一句 `want`；⛔ 把 `want 1` 改成 `want<=2`（派单硬边界第 7 条），计数断言⛔ 存在，换成"每枚形各问一次判定、答案必须相同＋等于绝对期望"；⛔ 删那枚缺失叶子（同样⛔）；⛔ 用 `t.Skip`／build tag。

ⓓ现在钉的是**判定**（且是**绝对**判定，⛔ 是"两枚判定相等"那种可被恒真买下的相对式）。

**Ⓐ／Ⓑ 选型（这一枚单独交裁）**＝选 **Ⓐ**。理由（Ⓑ 为什么建不忠实）：`spellings()` 的枚数问的是"**我手里这串还需不需要展开**"，⛔ 问"这枚文件存不存在第二种拼法"。Ⓑ 的"独立测量"只能问后者——`GetShortPathNameW` 对**已经展开的长形**照样回一枚短形（`C:\Users\swq\AppData\Local\Temp` 那支台面上，`TESTCA~1` 这类短名同样存在，因为 8.3 生成是卷属性、⛔ 是别名属性），于是"有短名⇒恰 2 形"会要求四枚**存在的路径**交出 2 枚形，而产码正确行为是 1 枚 ⇒ Ⓑ 在两台面上都⛔ 绿。唯一忠实的答案是"再问一次句柄最终路径"，而那**正是 `formsOf` 内部已经跑的那一步**＝用被测物自己当尺（formsOf wearing a different hat）。⇒ Ⓑ⛔ 成立，本腿交 Ⓐ 并报代价（见下）。

**代价（具名，⛔ 用"应该没问题"含过去时）**：`AC#5` 那句"对没有第二种拼写的路径，本票一个字都不改变判定"改完**仍站得住**，但站的是**判定层**——现在逐枚路径钉"绝对判定⛔ 移＋各形同判定"，比旧的"恰 1 形"更贴那句承诺的语义。⛔ 站住的是**形层**那半句：**"对 `Resolve` 交回的形，`formsOf` 恒产恰 1 枚比较形"这句话现在零仪器**。能悄悄溜过去的形状＝有人把 `spellings()` 改成返回 3 枚同判定形（判定层全绿、形层爆炸半径承诺⛔ 再有人管）；⛔ 是本次改坏的吗——⛔，改前那一支在别名台面本来就红，但改前它**至少是枚真钉**（在无 8.3 的台面上钉得住）。要补回来的唯一路＝让 `Resolve` 对缺失叶子也展开祖先（**产码**，⛔ 本票射程：a2 census #2.7 具名它波及 `Actable()`、`syncdirs.go` 票 307 刚落地的两支比较、`internal/tools/paths.go:131`）⇒ 这一格请编排者与下一枚验收腿裁：是记一条永久欠读数（撤销口令自定），还是另立一票。

## #3 `TestAListWinsWhereBothTablesHit`（原红位 `:238`，B-only 对照组）

ⓐ**坏形状＝override 键**（`Gate(out, {strings.ToLower(out): true})`，`out` 是 `tmpHome(t)` 的短拼；票面追加节 §3 第 1 条的窄规则**故意**⛔ 认非已解析形 ⇒ 红在拼法上，与规则本身无关）。同族 a2 census 记作**形 IV**。

ⓑ**新期望侧逐字**：
```go
	out := filepath.Join(home, "proj", "id_bothkey")
	mustWrite(t, out, "secret")
	resOut, err := Resolve(out, nil)
	if err != nil {
		t.Fatalf("Resolve(%q): %v", out, err)
	}
	if got := Classify(out); got != ClassB {
		t.Fatalf("B-only control %q classified %v, want ClassB", out, got)
	}
	if d := Gate(resOut.Canonical, map[string]bool{strings.ToLower(resOut.Canonical): true}); !d.Allow {
		t.Fatalf("B-only control must stay overridable by design: %+v", d)
	}
```
报错句逐字保留（`B-only control must stay overridable by design`）；`Classify(out)` 那枚（拿短拼问判定，本来就绿）**留在原地**＝"分类⛔ 依赖拼法"在这枚对照组上仍然被问过。

ⓒ⛔ 撞禁止解法：**⛔ 动 `overrideApplies`**（本文件⛔ 一行产码；`git show --numstat` 见 `40` 件），修法正是票面追加节 §3 第 1 条给的窄规则配套——**把键换成解析器交回的那一枚已解析形**；⛔ env 分档：同一句在两台面都用 `resOut.Canonical`；⛔ Skip／⛔ 删正控（census `30-answer-B-line.md:72` 具名"正控⛔ 能删也⛔ 能反转，否则 `overrideApplies` 恒 false 也全绿"＝真正的假绿形，本腿没走）。

ⓓ钉的是**判定**（"这条窄规则是规则、⛔ 是 blanket deny"），并且是**生产形状**的判定：`internal/tools/mode.go:91→:93 Canonicalize→:98 risk.Gate(c, overrides)` 就是这个调用序列。

## #4 #5 #6 `TestPathResolverShortNameAListDenied` / `…UNCAListDenied` / `…ExtendedLengthPrefixAListDenied`（原红位 `junction:135`／`:164`／`:182`）

ⓐ**坏形状＝拼写字面量**（三枚都是 `if !strings.EqualFold(res.Canonical, long)`，`long := filepath.Join(tmpHome(t), ".git-credentials")`＝测试自己拼的那一串）。

ⓑ**新期望侧逐字**：
```go
	// #4
	if !namesSameFile(res.Canonical, long) {
		t.Fatalf("short name did not expand to the file the test wrote: got %q, written path %q", res.Canonical, long)
	}
	// #5（循环内，三枚 UNC 拼法逐枚问）
	if !namesSameFile(res.Canonical, long) {
		t.Fatalf("UNC spelling %q normalized to %q, which is not a spelling of the written file %q", spelling, res.Canonical, long)
	}
	if got := Classify(res.Canonical); got != ClassA {
		t.Fatalf("UNC spelling of A-list file classified %v, want ClassA", got)
	}
	// #6
	if !namesSameFile(res.Canonical, long) {
		t.Fatalf("\\\\?\\ spelling normalized to %q, not to a spelling of the written file %q", res.Canonical, long)
	}
	if got := Classify(res.Canonical); got != ClassA {
		t.Fatalf("\\\\?\\ spelling of A-list file classified %v, want ClassA", got)
	}
```
ⓣ #5/#6 各**新增**一枚 `Classify == ClassA`（改前两枚只问 `Gate(...).Allow`，⛔ 问分类）＝收紧、⛔ 放宽；`gateFor(t, res.Canonical).Allow` 那枚原样留在后面。#4 原有的 `Classify == ClassA`／`Gate` ⛔ 动。

ⓒ⛔ 撞禁止解法：⛔ 把 `long` 换成手写常量或换成短名（解法 1／派单硬边界"⛔ 把 `t.TempDir()` 换成手写常量＝在测试里复刻 bug"）——`long` 仍来自 `tmpHome(t)`，只是⛔ 再拿去逐字符串比；⛔ 台面分档——`junction:126` 那枚 `getShortPath` 返回空时的 `t.Skip` 是**票 18 既有合法件，原地留着、⛔ 顺手收掉**（本腿对该文件新增 `t.Skip`＝0，尺见 `40` 件），本枚改动与它无关；⛔ 新增 Skip／build tag。

ⓓ钉的是**判定**：#4/#5/#6 现在问"是不是同一枚文件（句柄管线说了算）＋判不判 A＋可不可放行"。⛔ 钉拼法。

## #7 `TestBListDefaultDenyAndOverride`（原红位 `junction:292`）

ⓐ**坏形状＝override 键**（同 #3，形 IV：`Gate(p, {strings.ToLower(p): true})`，`p` 是短拼，且 `p` 那枚文件**从没被创建**）。

ⓑ**新期望侧逐字**：
```go
	p := filepath.Join(home, "proj", ".env.local")
	mustWrite(t, p, "secret")
	resP, err := Resolve(p, nil)
	if err != nil {
		t.Fatalf("Resolve(%q): %v", p, err)
	}
	if d := Gate(resP.Canonical, map[string]bool{strings.ToLower(resP.Canonical): true}); !d.Allow {
		t.Fatalf("B-list single-file override must allow: %+v", d)
	}
```
ⓒ⛔ 撞禁止解法：⛔ 动 `overrideApplies`；⛔ env 分档；⛔ 删那枚 `:294` 的日志断言。★census #7.3 具名的**次生落点**这里正面接住：改前 `:294`（`if len(logs) == 0 { t.Fatal("B-list override must emit a log entry") }`）被 `:292` 的 `t.Fatalf` 挡着**从未到达**；改后它**第一次真的执行**——本腿⛔ 碰 `Logf` 那两行挂钩／还原（`:288` 接管、`:289` 还原成 `nil`），实跑结果＝绿（`Gate` 的 B+override 分支走 `logf`，`blacklist.go:52` 判空命中已接管的 `Logf`），读数在 `40` 件两发里。

ⓓ钉的是**判定**＋**日志副作用**：`Allow` 与"必须留一行安全日志"，⛔ 再问字符串。

---

## 本件用到的尺

- 红位内容锚＝`00-anchor.md` §2 那张表（七句 `sed -n` 逐字原文）。
- 名册尺＝`grep -n 'func Test' <两文件>`（前 6／10，两文件合计 16，⛔ 变）。
- 判定与颜色尺＝`40-two-surface-readings.md`（两台面 × 改前／改后 × 7 枚指名／16 枚全量 ＋ 红名册 `comm` 作差）。
- 牙尺＝`50-mutations.md`（ⓐ 甲组新期望侧退成恒真 ⇒ 指名用例必须红；ⓑ 乙组负向钉拔掉／恒真形 ⇒ 必须红；ⓒ 拔掉产码守卫分岔作对照）。
