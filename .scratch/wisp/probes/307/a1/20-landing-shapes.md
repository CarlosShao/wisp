# 307-a1 · 20-landing-shapes —— 最小诚实落点候选（⛔ 选形，只给形状＋代价）

射程＝`internal/risk` 工作树；行号＝快照（内容锚在 `10-fold-roster.md`）。
每形写：动哪几行／会不会让别的判定变严或变松／会不会撞票面 `AC#1` 的四条硬约束。
本文件**没有任何一读来自跑测试**（本腿⛔ 跑编译面，见 `90-final.md` 自报），凡后果一律标〔读码推到〕。

`AC#1` 四条硬约束（票面逐字摘，作为尺）：
① 折法⛔ 新造一套——只许复用现成的 `Resolve`/`Actable` 或既有 `normPath`＋一次具名展开；⛔ 在 `risk.PathResolver` 之外用 `filepath.Clean|Abs` 做判定；
② ⛔ 让该 fail-closed 的那几支变松（`resolveTarget` 那两句 `errTargetUnverified`／`ErrReparseDenied` 一字不许动）；
③ ⛔ 顺手把 `allowed_dirs` 做成硬边界；
④ ⛔ 动三枚冻结件与 `thresholds.go`／`allowlist.txt`（票面逐字列的三枚＝`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`；本腿现跑跨面尺读到那两枚 panel 文件与本射程**零耦合**，见 `40-machine-alias.md` §4）。

现跑核对（工作树，尺＝`grep -rn --include=*.go "filepath.Clean\|filepath.Abs" internal cmd | grep -v "^internal/risk/pathresolver.go" | grep -v _test.go`）⇒ `internal/risk/*.go` 非测试面 **0 命中**（列出来的 13 行全在 `internal/agent`、`internal/memory`、`internal/models`、`internal/winsec`、`internal/tools`、`internal/projctx`、`cmd/wisp`，且除 `internal/memory/open.go:169` 与 `internal/models/manifest.go:293` 两枚真调用外全是注释里的自陈）。⇒ **`internal/risk` 里现在没有一枚绕过先例**，落点腿写下的任何 `filepath.Clean|Abs` 都会是这个包的第一枚，d22scan 那把尺（`tools/d22scan`）会当新增命中。

---

## 形 A —— 只折 `home` 侧：给 `home` 补一次具名展开

**形状**：新增一枚包内具名 helper（放 `syncdirs.go`），内容＝对 `orDefault(o.HomeDir, userHomeDir())` 走**同一枚** `Resolve(x, excs)` → `Actable()`，拿到形后仍过既有 `normPath`；`Resolve` 报错或形为空时**退回今天的 `normPath` 词汇形**。两枚赋值点（`syncdirs.go:88`、`provenance.go:345`）改成叫这枚 helper。判定行 `:428` 一字不动。

- 动几行：2 枚赋值点 ＋ 1 枚 helper（约 6-10 行新写）。**⛔ 碰 `resolveTarget`／`add()`／`finalize()`。**
- 变严还是变松：**只变严**——`home` 从短形变长形，`isUnder(cand, s.home)` 从"永不等"变"能等"，命中集合单调增。〔读码推到〕
- 残留的洞（诚实项）：`home` 自己打不开时（ACL、离线卷、`USERPROFILE` 指向不存在的目录）退回短形＝**今天这枚少报照旧**，且**没有任何声音**——形 A 不带 fail-closed 语义。〔读码推到〕
- `AC#1` 逐条：① 合规（复用 `Resolve`+`Actable`+`normPath`，⛔ 新折法、⛔ `filepath.*`）；② 合规（那两句不在射程）；③ 无关；④ 无关。
- 额外读数：`home` 变长形会**改变 `:431` 报出的 `Root.Path` 拼法**（现在报的是 `normPath` 比较形）。现跑读到 `Root.Path` 的用例只有 `syncdirs_ancestor_reparse_ticket105_windows_test.go:95,106`，两枚都用 `filepath.Clean` ＋ `strings.EqualFold` ⇒〔读码推到〕不受大小写/斜杠向影响，但它们**不在 `:431` 那支配 `s.home` 的分支上**（它们走根成员判定）。

## 形 B —— 两侧都过 `Resolve` 再比（把 `:160` 那枚短形出口一起堵上）

**形状**：形 A ＋ 把 `syncdirs.go:160`（`default:` 支）的 `normPath(expandInput(r.Path))` 换成带**空串守卫**的 `Resolve` 产出（`res.Canonical`）。

- 动几行：形 A 的全部 ＋ `:160` 一枚 ＋ 一枚守卫条件。
- ★**这一形有一枚会踩坏的陷阱（现跑代码锚，〔读码推到〕后果）**：`Resolve` 在 `ErrReparseDenied` 时是 `return res, ErrReparseDenied`（`pathresolver.go:123`），此刻 `res.Canonical` **还是空串**（`Canonical` 要到 `:131`/`:138` 才被赋值）。若把 `:160` 直接换成 `normPath(res.Canonical)` 而不守卫，就会造出 `e.canon == ""` 的根，而 `isUnder(path, "")`（`blacklist.go:169-171`）＝`path == "" || strings.HasPrefix(path, sepStr)` ⇒ **POSIX 上每一条绝对路径都算落进这枚根**（`cand` 恒以 `/` 开头），Windows 上则命中 UNC 拼法的候选。⇒ 后果＝**大面积假阳性（该判不可疑的全判可疑）＋ 把"一根未验证的根"升成"整棵树"**。这一枚是"顺手改 `:160`"最容易翻的地方，落点无论选哪形都建议看一眼。
- 变严还是变松：**两侧都变严**（`:422` 现在也能命中盘外/未建的根）。对 `AC#1`"本机默认台面⛔ 把别的台面弄红"这一条，形 B 是**风险最高的一形**：它改变的是根成员判定，而 `syncdirs_windows_test.go:54`、`:92`、`syncdirs_test.go:89`、`:375`、`:447`、`syncdirs_other_test.go:91`、`syncdirs_ancestor_actable_leg_116_test.go:209,212` 那一族**负向**断言（"plain/非根下写⛔ 是 sync"）全部建立在"⛔ 匹配＝⛔ 报警"上。〔读码推到：`AC#1` 的"⛔ 弄红别处"要现跑才能定，本腿跑不了。〕
- `AC#1` 逐条：① 合规**前提是**⛔ 用 `expandInput` 之外的新词汇步骤（`lexCanonical` 在 `pathresolver.go` 内，包外不可调）；② **有边缘**：若守卫写成"空串就 `return` 错误"，等于给 `add()` 新增一条错误面，虽不碰 `resolveTarget` 那两句，但会把"根不可验证"从 `canonical=false` 变成"报错"，与 `finalize()` 的既有语义（`canonical` 只是不`complete`，⛔ 报错）分家——⚠ 这一处是**语义扩张**，请编排者裁；③ 无关；④ 无关。

## 形 C —— 复用 `pathresolver.go` 里既有的 `formsOf`／`anchorForms` 那一族

**形状**：判定行 `:428` 从 `isUnder(cand, s.home)` 换成包内已经为票 72 建好的关系函数：`aUnder(formsOf(cand), anchorForms(s.home))`（`anchorForms`＝带 `anchorCache` 的 `formsOf`；`aUnder`→`aRelates`→逐 `spellings()` 并集比 ＋ `uncertainAnchorMiss` 兜 fail-closed）。〔现跑代码锚〕`blacklist.go:199` `homeS := normDir(userHomeDir())` ＋ `blacklist.go:221` `aUnder(cand, anchorForms(anchorPath(homeS, ".ssh")))` ＝**这条链早就在给 `home` 干同一件事**；`blacklist.go:189-193` 那句票 72 不变量逐字＝"the verdict must not depend on the spelling a path arrived in"。

- 动几行：**1 行判定**（`:428`）＋可能 `s.home` 的存取形式（若把字段留作 string 则 helper 内现折；`anchorForms` 已缓存 ⇒〔读码推到〕每条判定一次 `normPath` ＋ 至多一次 `resolveHandle` 走步）。可选把 `:422` 同形换成 `aUnder`，那部分＝形 B 的面＋形 C 的折法。
- 变严还是变松：**变严，且严的那一半是 R17 的 fail-closed 而不是新增匹配规则**——`spellings()` 是并集（`pathresolver.go:332-337`，注释逐字"only ever turn a miss into a deny"），所以**当前任何一枚命中都不会因换形而消失**〔读码推到〕；新增的是"两侧都是 handle 形但一侧 partial 且 root 同源"时判命中。
- 影响面（现跑尺，逐枚列在 `30-test-nails.md`）：负向断言那一族是风险位；`uncertainAnchorMiss`（`blacklist.go:314-322`）要求两侧 `root` 非空且互为祖先 ⇒〔读码推到〕`outOfProfile`/`profileevil` 那类"同一 `base` 下的兄弟目录"两侧 root 都能解析且互不含 ⇒ 仍判不命中 ⇒ 应不变红，**但这枚必须现跑两台面才算数**。
- 台面依赖：POSIX 上 `resolveHandle` 恒 `("", false)`（`pathresolver_other.go:10`）⇒ `formsOf` 只产 `.raw`、`.root`＝"" ⇒ `uncertainAnchorMiss` 早退 false ⇒〔读码推到〕POSIX 颜色与今天逐字一致（这一枚是 `AC#1`"⛔ 把别的台面弄红"的正证候选）。
- `AC#1` 逐条：① **最贴合**——这一族就是"现成的、已用在 `home` 锚点上的那一族"，⛔ 新造折法、⛔ `filepath.*`（`formsOf` 内部调的 `lexCanonical` 在 `pathresolver.go`，属 D22 许可面）；② 完全不碰 `resolveTarget`；③ 无关；④ 无关。
- 独有代价：`aUnder` 把 A 名单的 fail-closed 判据搬进同步判定 ⇒〔读码推到〕`Sync:true` 的集合变大，而 `Sync:true` 是"少一次 L2 询问还是多一次"的输入之一（`provenance.go:928` ⇒ `ChSyncWrite` ⇒ 门开着）⇒ **只会多问、不会少问**，与票 252 那族"多报"用户体验同向，与本票想消除的"少报"反向对齐。⚠ 但它同时给 `Root`/`Why` 报出更多"suspect-fallback"，面板那一侧的读数会变（本腿⛔ 量面板，具名欠读数）。

## 反形（⛔ 议，只登记代价）：把 `cand` 退回词汇形去迁就 `home`

`resolveTarget` 里去掉 `Resolve` 只留 `normPath(raw)` ⇒ 两边都短、5 枚当场复绿，**同时把票 252／对抗报告 B-1 修掉的 `\\?\` 与 8.3 拼法绕过重新放开** ⇒ 撞 `AC#1`②（放松 fail-closed 的方向）与票面"⛔ 为变绿放宽任何断言"。写在这儿只为让"颜色对上了"不被当凭据。

---

## 本文件⛔ 做到的

- ⛔ 选定落点（派单⛔ 由落地腿／本腿自裁）；⛔ 写"某形可行／已验证"——所有"不变红／仍绿"都是〔读码推到〕，⛔ 台面色（别名台面／本机台面）一发都没跑，那需要 `go test ./internal/risk/ -run '<名册>' -count=1 -v`（要编排者车道跑，名册见 `30-test-nails.md`）。
- ⛔ 给夹具设计（`AC#3` 的"自己造得出那个世界"）：只报到 `40-machine-alias.md` 的读数层。
