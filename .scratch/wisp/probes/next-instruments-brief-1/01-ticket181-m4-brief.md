# 简案一：票 181 的 M4 盲区——`internal/panel/workspace.go:141` 的成功路径填充点

> 备料腿 `next-instruments-brief-1`（只读）｜取锚 `2026-10-08 19:01` HEAD `5f52f310`（A736）｜词面/行号一律本程现取（`sed -n`／Read 带行号），包先定为 `internal/panel` 再量行数。
> 本件不装仪器、不种突变；只给将来写腿的派单级形状。出处：台账 `A731`（`docs/reports/pending-and-issues.md:14344-14352`，含编排者亲手复现）、`probes/181/v3/evidence.md` §5。
> 射程声明：未读未引 `frontend/**`／`design/**`；三枚冻结件未碰（射程判断，非内容引用）。

## 1. 可达性判语（先给结论，再给两层证据）

**判语：成功路径本身可达；该填充点的「值域」是单点 `{false}`——`res.Rewritten == true` 今天结构上到不了 `:139-147`，`:144-146` 是死枝。**

第一层（成功路径可达，有现役用例在跑）：
- `:116` `res, err := scope.ResolveWorkspace(input)`；`:117` `if err == nil {`；`:120-124` `if actable, aErr := res.Actable(); aErr != nil { err = aErr } else if sErr := scope.SetWorkspaceRoot(actable, res); sErr != nil { err = sErr }`；`:126` 之后任何 `err != nil` 走拒绝路径。
- 成功路径（`:126 err == nil` 起）由 `workspace_test.go:50` `TestWorkspaceSwitchAuditsAndAppliesACleanPath` 现跑：serialized `risk.Result{... Resolved: true}`（`Rewritten` 零值 false）⇒ `err=nil` ⇒ 返回 `view`。

第二层（携带真值的变体不可达，`Actable` 闸逐字挡死）：
- `internal/risk/pathresolver.go:93-98` 逐字：`func (r Result) Actable() (string, error) { if r.Rewritten { return "", fmt.Errorf("%w: ...", ErrRewrittenPath, ...) } return r.Canonical, nil }` —— `Rewritten=true` **必**返非 nil error。
- `:120-121` 收到 error ⇒ `err = aErr` ⇒ `:126` 拒绝（已有用例 `workspace_test.go:99` `TestWorkspaceSwitchRefusesARewrittenAccountEvenIfTheResolverSaysOK` 守此）。
- ⇒ `:141` 处 `Rewritten: res.Rewritten,` 的成功路径真值恒为 `false`；`:144 if res.Rewritten` 的文案分支（"展开改写了这条路径……已拒绝授权改写后的树"）永不可达——**死枝**。
- 旁证（本程现取）：`:91` 与 `:141` 词面**逐字相同**（`Reparse: res.Reparse, Rewritten: res.Rewritten,`），两处都靠行号定位；`Rewritten:` 非测试全仓填充点仅 `workspace.go:91`／`:141`（`risk/pathresolver.go:116` 是账的来源，非视图）。

## 2. 「该不该现在就装」：代价摆全，⛔ 不替编排者选

- **A｜现在装**：代价＝给单点值域装仪器（"给死枝装仪器"的准确含义：它守的是"成功 view 不得谎称 rewritten"，防的坏法只有一种——有人把 `:141` 写成常量 `true`；它**不能**证明"改写真值沿此路传递"，因为"成功 + rewritten=true"今天不存在）。好处＝一旦 186 或人工批准放开该路，仪器就位即咬。
- **B｜等票 186（"the-composer-area-should-detect-git-and-let-the-user-switch-between-local-worktree-and-branch"）接上面板 handler 后再装**：代价＝期间 `:141` 继续零仪器（M4 现状）；186 落地后必须先重测三读数再决定形状：①`RequestWorkspaceSwitch` 非测试调用者从 0 枚变为几枚；②`:120 Actable()` 闸是否仍结构恒挡（若不动，仪器仍只能一向）；③`res.Rewritten=true` 能否走到 `:139`（要能动得改 `internal/risk` 判级语义＝须人工批准，禁区）。⇒ B 的"两向可造"**不是 186 自动带来的**，是"186 + 判级语义变更"才带来。
- 两个选项都不提红线：⛔ 不许为装仪器去放宽 `Actable`／改 `risk` 判级语义（契约变更，人工批准）。

## 3. 若选 A：最小形状（两向区分要求下的诚实版）

**测试文件（新建）**：`internal/panel/workspace_switch_success_181_test.go`（新测试函数；同包可直接复用 `workspace_test.go:17-40` 的 `scriptedScope` 与 `:42-48` 的 `auditLog`，⛔ 不改既有断言行）。

**断言（一句一牙）**：
1. 用 `scriptedScope{res: risk.Result{Spelling: plain, Canonical: plain, Resolved: true}}` 调 `RequestWorkspaceSwitch(scope, plain, log.f)` ⇒ 先 `err==nil`、`view.Set`、`view.Canonical==plain`；
2. `view.Rewritten == false`（抓 `:141` 填 `true`）；
3. `view.Reparse == false` 且 `view.Reason` 不含 `"改写"`；
4. **对账式**：审计行 `:135-137` 写的是 `rewritten=%v`（直接读 `res`），断言 `view.Rewritten == 该审计行的字面值`（今天两侧应都为 false）。这枚把"视图 vs 它自己写的账"钉在一起（`:135` 审计行与 `:139-141` 视图共用 `res`，只有视图被填错时两头打架）。

**两向区分：一向可造、一向结构上造不出（具名）**：
- 可造的一向＝填 `true` 红。正控种法（写腿照做，**先复量再种**）：`sed -n '141p'` 应逐字 `\t\tReparse: res.Reparse, Rewritten: res.Rewritten,`；`sed -i '141s/Rewritten: res.Rewritten,/Rewritten: true,/' internal/panel/workspace.go`；改后 `sed -n '141p'` 复量＝`... Rewritten: true,`；`git diff --stat` 应只见这一处。预期：新测试 rc=1，断言 2 与断言 4 各自红。
- **造不出的一向＝填 `false` 红**：填 `false` 与真值在此路**外延等价**（§1 第二层），今天不存在"成功 + rewritten=true"的可达输入 ⇒ 行为测试在原理上不可区分。若编排者坚持要这一向，唯一可买的形状＝**静态文本尺（A2）**：读 `internal/panel/workspace.go` 源码，在 `func RequestWorkspaceSwitch` 到下一个函数边界之间断言"恰一处 `Rewritten:` 填充且右值逐字为 `res.Rewritten`，不得是 `true`／`false` 常量"。先例＝risk 的 AC#5 static criterion（`internal/risk/pathresolver_rewrite_account_test.go:163` `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`，其性质见 `pathresolver.go:85-86` 注释"the static criterion ... fails the build-level gate if one stops"）。⚠ A2 的红是**锁字面**不是行为（填 `false` 时行为上与真值无差，"红"只因字面不等于 `res.Rewritten`）⇒ 不许把 A2 的红当成"行为两向"来记账。

**降级复认实验（可选读数，不是失败）**：种 `:141`→`false` ⇒ 预期**全绿**；它复认的是"此路只可达 false"（M4 判语的独立复现），用途与"突变见红"相反，写读数时须写清。

**会不会撞既有在册红**：`internal/panel` 今天在册 6 枚（`A731`／`181-v3` 在册：本程 2 枚 flood/dropped＋既有 4 枚）。新仪器用**焦点 `-run`** 跑（如 `-run 'TestWorkspaceSwitchSuccess'`），⛔ 不进整包跑、⛔ 不许把在册 6 枚算进判据、⛔ 不许为它们放宽任何断言。

## 4. 与 `:91` 已有仪器的关系（一句话）

`:91`（`WorkspaceViewFromRoot`）是**真生产链**的填充点（`cmd/wisp/panel_pump.go:87` 调用），两向有牙：clean⇒`false` 由 `workspace_account_181r3_test.go:65`、rewritten⇒`true` 由 `:82`；`:141` 是**零生产调用者**＋`Actable` 恒挡路径上的填充点，今天只能一向——`:91` 的牙覆盖不到它（M4 实测全绿），两处是同字段在两个面上的仪器，各自记账。
