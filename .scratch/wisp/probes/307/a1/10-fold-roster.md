# 307-a1 · 10-fold-roster —— `AC#0` 的折法名册（只读数）

锚点：进场 HEAD `1eaa91f`（见 `00-anchor.md`）。射程＝工作树（`internal/cmd` 对 HEAD 干净，0 行）。
每一件**行号都现跑过**（尺附在每节末），行号＝**快照**，判据＝抄出来的原文（内容锚）。
标记：〔现跑〕＝这条命令这条输出；〔读码推到〕＝从码面读出的可达性/后果，**没跑过码**。

---

## 0. 票面/派单的行号先对一遍（现跑）

尺＝`grep -n '<原文>' internal/risk/syncdirs.go internal/risk/blacklist.go internal/risk/provenance.go`（工作树）。

| 票面写的号 | 盘上现量 | 内容锚（逐字） | 判定 |
|---|---|---|---|
| `syncdirs.go:425` | **`syncdirs.go:422`** | `		if isUnder(cand, e.canon) {` | **过期 −3 行** |
| `syncdirs.go:429` | **`syncdirs.go:428`** | `	if !s.complete && s.home != "" && isUnder(cand, s.home) {` | **过期 −1 行** |
| `syncdirs.go:88` | `syncdirs.go:88` | `	s := &syncSet{home: normPath(orDefault(o.HomeDir, userHomeDir())), excs: o.ReparseExceptions}` | 对 |
| `syncdirs.go:197` | `syncdirs.go:197` | `func (s *syncSet) resolveTarget(raw string) (string, error) {` | 对 |
| `blacklist.go:133` | `blacklist.go:133` | `func normPath(p string) string {` | 对 |
| `blacklist.go:169` | `blacklist.go:169` | `func isUnder(path, dir string) bool {` | 对 |
| `provenance.go:928` | `provenance.go:928` | `		if p.IsSyncPath(t).Sync {` | 对 |
| 测试侧 `:162/:213/:233/:355/:402` | 同号 | 见 `30-test-nails.md` | 对 |

⇒ 具名报回：**只有产码那两枚 `:425`/`:429` 腐烂**（同一文件里 `resolveTarget` 与 `home` 赋值点没动，所以是 `match()` 内部被删过 3 行）。落点腿一律以内容锚定位。

---

## 1. `home` 侧（`:428` 的右操作数 `s.home`）——折法名册

尺＝`grep -rn --include=*.go "syncSet{" internal`（工作树，含注释行；`internal/risk`＋`internal/cmd` 全 `*.go`）⇒ **带 `home:` 的构造点＝2 枚**，其余 4 枚是 `&syncSet{}` 空构造（测试件里 `pathresolver_expansion_test.go:85,100`、`pathresolver_rewrite_account_test.go:134,145`）。

| 赋值点 `file:行` | 原文（内容锚） | 过了哪几层归一 | 过不过 `Resolve` |
|---|---|---|---|
| `internal/risk/syncdirs.go:88` | `s := &syncSet{home: normPath(orDefault(o.HomeDir, userHomeDir())), …}` | 只有 `normPath` | ⛔ **不过** |
| `internal/risk/provenance.go:345` | 与上一枚**逐字同形**的复制件（`if o.NoProbe {` 分支内） | 只有 `normPath` | ⛔ **不过** |

- `normPath`（`blacklist.go:133`）逐字两层：`u := strings.ToLower(unifySeparators(p))` ＋ 剥尾分隔符 ⇒ **折小写、折斜杠，⛔ 折 8.3、⛔ 走解析器**〔现跑 `sed -n '133,141p'`〕。
- `userHomeDir()`（`blacklist.go:354`）逐字＝`os.UserHomeDir()` ⇒ Windows 上就是 `%USERPROFILE%` 的**原始拼法**（短名机器上即 `C:\Users\RUNNER~1` 那一族）。
- 消费点只有 2 处：`syncdirs.go:428`（判定）与 `syncdirs.go:431`（`Root: SyncRoot{Provider: "sync-suspect", Path: s.home, …}`＝把**比较形**当路径报出去；`normPath` 自己的注释写着"the shape it produces must never be handed to the OS or printed as a path (ticket 75)"）。〔现跑〕`grep -rn "\.home\b" internal/risk/*.go` 非测试面＝`:428`、`:431` 两行。

### 1.1 两枚构造点各自谁能走到（现跑）

尺＝`grep -rn --include=*.go "NoProbe" internal cmd | grep -v _test.go` 与 `grep -rn --include=*.go "NewProvenance(" internal cmd | grep -v _test.go`。

- `internal/risk/provenance.go:344` `if o.NoProbe {` ⇒ **`NoProbe:true` 走 `:345`**；否则走 `detectSyncSet(o)`＝`syncdirs.go:88`。
- **产码面只有 2 个 `NewProvenance` 调用者，且两个都是 `NoProbe: true`**：
  - `cmd/wisp/run.go:743` `c25 := risk.NewProvenance(risk.ProvOptions{NoProbe: true})`
  - `cmd/wisp/panel_assets.go:239` `prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})`
  - ⇒ 〔现跑〕尺＝上面那条 `NewProvenance(` 非测试 grep，命中 3 行（1 枚是 `provenance.go:293` 的定义）。
- 非 `NoProbe`（＝能走 `syncdirs.go:88`）的调用者＝**只有 2 枚测试**：`syncdirs_test.go:488` `p := NewProvenance(ProvOptions{}) // live env` 与 `syncdirs_redteam_windows_test.go:210` 同形。〔现跑〕
- 票面点名的 5 枚用例（`:162/:213/:233/:355/:402`）**全部**是 `ProvOptions{NoProbe: true, HomeDir: …}` ⇒ 全走 `provenance.go:345`。〔现跑：见 `30-test-nails.md` 的 grep〕

⇒ **读数（不含结论）**：`home` 这枚折法有 2 个赋值点；**生产面走的是 `provenance.go:345` 那一枚**，`syncdirs.go:88` 那一枚在生产面当前不可达（因为两个产码调用者都 `NoProbe:true`），但它是 `syncdirs_test.go:501` 那枚用例的入口。派单问"同形复制件要不要一起算进落点"——**两枚逐字同形、且其中一枚正是生产入口**，要不要一起动请编排者裁。

---

## 2. `e.canon` 侧（`:422` 的右操作数）——折法名册

尺＝`cat -n internal/risk/syncdirs.go | sed -n '139,168p'`（工作树）⇒ `(*syncSet).add()` 内 `c.canon` 赋值点＝**3 枚**，别无其它（`grep -rn "c\.canon" internal/risk/*.go`＝3 命中）。

| 赋值点 `file:行` | 分支条件（原文） | 右操作数原文 | 过了哪几层归一 | 过不过 `Resolve` |
|---|---|---|---|---|
| `syncdirs.go:152` | `case err == nil && res.Rewritten:` | `c.canon, c.canonical = normPath(res.Canonical), false` | `Resolve`(step1-5) → `normPath`；`res.Canonical` 的形**取决于根在不在盘上**（下条） | ✅ **过**（`:140` 那次 `Resolve`） |
| `syncdirs.go:156` | `case err == nil && res.Resolved && res.Canonical != "":` | `c.canon, c.canonical = normPath(res.Canonical), true` | 同上，且此支 `Resolved==true` ⇒ `Canonical` 必是 `resolveHandle` 的**长形**（`pathresolver.go:127-134`） | ✅ **过（长形）** |
| `syncdirs.go:160` | `default:` | `c.canon, c.canonical = normPath(expandInput(r.Path)), false` | **只有 `expandInput`＝step 1（env/`~`）＋ `normPath`**；⛔ `lexCanonical`、⛔ `resolveHandle` | ⛔ **不过**（`res` 的答案被丢弃，重新从 `r.Path` 折） |

⚠ `syncdirs.go:82` 那句注释 `canon string // normPath(Resolve(Path))` 只对 `:156` 成立，对 `:160` **不成立**——派单要求"⛔ 按注释读"，这一枚就是会读歪的地方。〔现跑 `sed -n '80,84p'`〕

### 2.1 `Resolve` 的两种产出（决定 `:152` 是长还是短）

`pathresolver.go:127-139` 原文（现跑）：

```
	if final, ok := resolveHandle(p); ok {
		final = stripExtendedPrefix(final)
		final = normalizeLocalUNC(final)
		if final != "" {
			res.Canonical = final
			res.Resolved = true
			return res, nil
		}
	}
	// Path does not exist (or handle query unsupported): lexical fallback,
	res.Canonical = p
```

⇒ `Canonical`＝**盘上有 ⇒ 长形（8.3 已展开）**；**盘上没有 ⇒ `lexCanonical` 的词汇形（Abs+Clean，⛔ 展 8.3）**。〔现跑〕
Windows 的 `resolveHandle` 在 `pathresolver_windows.go:25`；POSIX 半边是 `pathresolver_other.go:10` `func resolveHandle(p string) (string, bool) { return "", false }`（`//go:build !windows` ＋ `DEFERRED(macOS/Linux)`）⇒ **POSIX 上每一枚根都落 `:160`**。〔现跑〕

### 2.2 `roots` 的上游赋值链（逐枚跟到底，现跑）

`grep -rn --include=*.go "\.add(\|detectSyncSet(" internal/risk/*.go | grep -v _test.go` ⇒ **3 枚 `add()` 调用点**：

1. `syncdirs.go:114` `s.add(r)` ← `found`（`syncdirs.go:101-108`）＝`envConfiguredRoots` ＋ `registryProbe` ＋ `clientConfigProbe` ＋ `defaultLocationProbe` ＋ `fixtureRoots(o.SyncFixturePath)` ＋ `o.SyncRoots`（仅 `!NoProbe` 时）
2. `provenance.go:347` `s.add(r)` ← `for _, r := range o.SyncRoots`（`NoProbe` 分支）
3. `provenance.go:351` `s.add(r)` ← `fixtureRoots(o.SyncFixturePath)`（`NoProbe` 分支）

- `defaultLocationProbe`（`syncdirs.go:544-561`）只收 `os.Stat` 得到且 `IsDir()` 的目录 ⇒〔读码推到〕这一族的根必在盘上 ⇒ 走 `:156`（Windows）或 `:160`（POSIX）。
- `envConfiguredRoots` / `registryProbeFor` / `clientConfigProbe` 收回来的拼法**不保证在盘上**（配置里写了但目录被搬/未建）⇒〔读码推到〕落 `:160`。
- **生产面 `roots`＝空**：`cmd/wisp/run.go:743` 与 `cmd/wisp/panel_assets.go:239` 传的 `ProvOptions` 里**没有** `SyncRoots`／`SyncFixturePath`（〔现跑〕尺＝`grep -rn --include=*.go "SyncRoots:" internal cmd | grep -v _test.go` ⇒ **0 行**，rc=1）⇒ `finalize()`（`syncdirs.go:123-136`）在 0 枚根上循环 ⇒ `complete` 恒 `false` ⇒〔读码推到〕**生产面 `:422` 那个 for 循环零次迭代，`:428` 是唯一的阳性出口**。

### 2.3 左操作数 `cand`（两支共用，现跑）

`resolveTarget`（`syncdirs.go:197`）＝ `Resolve(raw, s.excs)` → `res.Actable()`（⛔ 直接读 `Canonical`）→ `res.Resolved` 时 `normPath(canon)`（`:216`）；否则 `deepestExistingAncestor(canon)` → 对**既存祖先**再 `Resolve` → `ares.Actable()` → `normPath(anceCanon + sepStr + strings.Join(rest, sepStr))`（`:255`）。
⇒ 两支的左边**总是**"既存前缀的 handle 长形 ＋ 词汇尾巴"；错误面 `:408` 一律 `Sync:true`（`errTargetUnverified`／`ErrReparseDenied` 两支 fail-closed 就住在这里）。〔现跑 `sed -n '197,256p'`〕

---

## 3. `AC#0` 那一问的读数回答（⛔ 选落点、⛔ 勾框）

问：`roots` 侧的 `e.canon` 与 `home` 侧的 `s.home` 是不是**同一把折法**？

- **不是同一把**：`e.canon` 的 3 枚赋值点里 **2 枚过 `Resolve`**（`:152`/`:156`，`:156` 必长形），`s.home` 的 2 枚赋值点 **0 枚过 `Resolve`**。
- 但 `e.canon` 里**有 1 枚（`:160`）也不过 `Resolve`**，且它那枚连 `lexCanonical` 都没过（`expandInput` 只有 step 1），而 `cand` 总是过了 `Resolve`/`lexCanonical` 的形 ⇒〔读码推到〕**同一枚"一边折长、一边留短"的机制在 `:422` 也成立**，触发条件＝那枚根**盘上不存在**（POSIX 上则是**恒成立**）。票面点名的 `:162` 那一枚（fixture 根 `home\Cloud\NutstoreSync` 从未被 `sandbox()` 创建，`syncdirs_test.go:120-130/152-163` 现跑可证）走的正是这条。
- ⇒ 出口枚数的读数＝**2 处**（`:422` 的 `:160` 支 ＋ `:428`），其中**生产可达＝1 处**（`:428`，因为生产 `roots` 为空）；`home` 这枚折法自身的赋值点＝**2 处**（`syncdirs.go:88`／`provenance.go:345`），生产走后者。
- 另有一枚**顺带读数**：`:431` 把 `normPath` 的比较形当 `Root.Path` 报出去，与 `blacklist.go:128-132` 那段注释（ticket 75）相抵；它不是本票的判定，但**任何改 `home` 折法的落点都会改变这枚报出的拼法**。现跑的读它的用例＝`syncdirs_ancestor_reparse_ticket105_windows_test.go:95,106`（两枚都用 `filepath.Clean` ＋ `strings.EqualFold`，不锁大小写之外的形）。

⛔ 本格不写"闭合"，⛔ 由本腿裁落点形状——形状候选见 `20-landing-shapes.md`。
