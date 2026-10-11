# 309-a1 · 90 本腿终局（只读普查腿，格＝`AC#0`，⛔ 裁形）

## 一、交了几笔 commit

| 笔 | 内容 | `git show --name-only --format=` 枚数 | 落在 `probes/309/a1/` 外 |
|---|---|---|---|
| `2a37d612` | `00-anchor.md` | 1 | **0** |
| `27059960` | `10-three-shapes.md`·`20-shape-a-cost.md`·`21-shape-b-cost.md` | 3 | **0** |
| 本笔（第 3 笔） | `30-structure-guarantee.md`·`90-final.md`·`logs/**`（12 枚原始读数件） | 14 | **0** |

- 尺（可逐字重跑）＝`git show --name-only --format= <sha> | grep -c .` ⇄ `… | grep -v "probes/309/a1" | grep -c .`；前两枚已现跑＝1/0 与 3/0（本笔的 14/0 由 `git add` 的显式 pathspec 保证：`git add -- .scratch/wisp/probes/309/a1` ＋ commit 再带同一 pathspec，⛔ `git add -A`/`.`）。
- ⛔ push（子代理只 commit）；⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`--no-verify`；⛔ 在仓内建 worktree/checkout；⛔ 动 `internal/**`·`cmd/**`·票面·台账·`docs/**`（写面只有本目录）。
- `logs/` 一并进仓的理由＝车队卫生那族（台账 `A841`：台账引用的证据件"只在盘上从未提交"）；件名⛔ 用 `.out`（根 `.gitignore` 第 8 行是全仓 `*.out`），三态均为 **HEAD blob 落盘**，总量 222 KB。

## 二、本腿答了什么（一句一条）

1. 票面三处（①②③）**逐字全部成立**，三个行号锚（`:925`/`:164`/`:263`）**现量未漂**；仪器锚 `failclosed_236_teeth_test.go:244`/`:288` 未漂。
2. **名册该补第 4 枚**：`internal/tools/task.go:708-711`（`task.cancel` 对空宿主 id 的 fail-closed 拒绝）＋它的仪器 `failclosed_236_teeth_test.go:107/:118/:122` 与参数化载体 `task_cancel_221_legs_test.go:116-119` ⇒ "同形三种处置"其实是**同形四种、两拒一自命名一静默**，静默那枚⛔ 仪器。
3. **`b.assess` 本来就没挂探测器**（`bridge.go:208-214`）⇒ "摘掉 R4 探测器"是结构事实；机械短路点在 `internal/risk/rules_gateway.go:102` `if ctx.taint == nil {` ⇒ 走 `:926` 之后**门照走、`IsSyncPath` 那一跳整条不发生**。
4. 到门距离：① **0 跳**（守卫长在 `:297` 的调用表达式上）、② 到桥门 **0 跳**（自造 assessor 走 `panel.NewApprovalCardView`）、③④ 在门**下游**（改不动任何判级）。
5. **ⓐ形代价**：必动 2 枚文件（`internal/tools/bridge.go` 加 `b.reject` 分支＋`failclosed_236_teeth_test.go` 两枚用例改载体）＋条件必动 1 枚（新仪器保 ③④ 覆盖）；**静态必红 2 枚断言**（`:257`、`:122`，同一枚文件）；条件红 2 处（若误用 Go error ⇒ `:292`/`:123` 的 `t.Fatalf`）；用户侧今天**没人看到**，但代价含"③④ 的教学文案被上游遮蔽"。
6. **ⓑ形代价**：必动 1 枚（`bridge.go:920-923` 文档块＋守卫上方注释）；**静态 0 枚必红**（逐枚查过谁会把桥的注释读进断言：桥里被文字钉住的是 `panel_inbound_guards_35r3_test.go:65-67` 那三支别的分支）；用户侧零变化；实际防御增量为零，"响"必须由 `AC#1`/形 B/形 C 另给。
7. **"⛔ 删那一肢"这条边界本腿⛔ 越**：`provenance.go:694-710` 现量证明地板真实存在，但**比票面措辞窄一档**——`unbound` 只在"引擎里另有 scope 带污点"时置真（`:702` `if reg != nil { return nil, false }`、`:704-709` 的 `for … if len(r.marks) > 0`）⇒ 未读敏感源的任务撞⛔ 报。方向不变，措辞该修。
8. **`AC#2` 有料可交**：仓里已有**四族零接线尺形**（形 A 枚数钉 `gate-clauses.sh` 的 `want_n`、形 B `go/parser` 源码头文件如 `panel_assets_143_test.go:267`／`firstrun_257_test.go:430`、形 C 入口拒绝（今天就被 2 枚测试载体可达）、形 D reflect API 钉 `grant_test.go:509`／`ticket90_test.go:454`），三类未来候选者的半成品现量＝候选者 1 有（②本身；生产 `.Execute(` 非桥命中 1 枚＝`subagent_197.go:581` 的纯转发）、候选者 2 有形状无 id 语义（同那枚装饰器）、候选者 3 只有 `scheduler/doc.go` 零实现。

## 三、欠读数具名清单（只有实测／现跑能答，⛔ 归本腿，一律归编排者或落地腿）

| 号 | 内容 | 该由谁用什么尺答 |
|---|---|---|
| U-1 | `internal/tools` 整包名册里，除已点名的 2 枚载体之外还有没有别的用例依赖空 id 走进 `assessorFor` | `go test ./internal/tools/ -count=1`（⛔ 本腿禁跑；`72-r1` 交完之后由 `309-r1` 跑） |
| U-2 | ⓐ 新增的 reject 分支会不会打红"审计行数／`OutcomeKind` 取值"族钉（含 `cmd/wisp` 面板侧） | `go test ./internal/tools/ ./cmd/wisp/` |
| U-3 | `internal/tools/*_test.go` 内有没有"运行时读 `.go` 源码文本"的钉（我只在 `cmd/wisp` 找到 5 枚范本，`internal/tools` 只静态看了 `bridge.go` 提及） | `git grep -n "ReadFile\|ParseFile" HEAD -- 'internal/tools/*_test.go'` ＋跑色 |
| U-4 | 今日 `gate-clauses.sh` 的 G1/G1b/G3 命中行数 vs 登记 `want_n`（静态等价尺给 G1b＝1 vs 登记 0） | `bash .scratch/wisp/probes/154/gate-clauses.sh` |
| U-5 | `bridge.go:861` 那句 `internal/agent/loop.go:366` 是否也漂（我只核了 `:332/:321` 两枚） | `grep -n "defer.*CloseTask\|AdmitTask" logs/loop.go` 现读 |
| U-6 | 票 143 那枚测试对 ② 的覆盖度（钉的是 assessor 装配形状，还是也钉 `taintSourceScopeID` 的命名故意） | 读 `cmd/wisp/panel_assets_143_test.go` 全文＋跑色 |
| U-7 | 参数化载体 `r183Call`/`r185Call`/`sub197Harness.spawn` 的**全部调用点**里有没有别的传空 id（我只现跑 `cancel(t, ""`＝1 命中、字面量 `TaskID: ""`＝1 命中） | 逐枚展开调用点或 `-run` 实测 |

⇒ **共 7 条欠读数**，全部具名登记于此，⛔ 本腿代答。

## 四、我跑过的编译面命令（硬约束 1 自报，终局口径）

- `go build`／`go vet`／`go test`／`gofmt`／`gofumpt`：**0 条**。
- `go env`／`go list`：**0 条**（连允许的两条也没用上，本腿所有问题都能由 `git show`/`git grep` 静态答）。
- 起手之后补跑的一次干净度复核（见下）与全部读数：`git show HEAD:<path>`·`git grep <pat> HEAD`·`sed -n`·`git ls-tree`·`git log`·`git status --porcelain -- internal cmd`·`git show --name-only --format=`。

## 五、我没做、但射程边上有想法的事（⛔ 动手，一句话一条）

1. 全仓普查"注释里用现在时写没有消费者的页面"这一族（`git grep -n "Today no production\|no production code" HEAD -- internal cmd`）——票 40/`A838` §2 那族病可能不止一处。
2. `bridge.go:132` 的 `injected` 与 `Options.Assessor` 生产 0 赋值 ⇒ 值得一枚"这字段该不该存在"的票（票面 `:27` 已说归装配族另裁，我只是确认现量仍是 0）。
3. `assessorFor` 每次调用 `risk.NewRiskAssessor()` 重建一枚判级器 ⇒ SLO/延迟预算面上有没有人读过这个开销，没人管过。
4. `subagent_197.go:581` 那枚纯转发是今日唯一的非桥 `.Execute(` ⇒ 154 的 G1b 若不加排除说明，下次别人现跑会以为是漂移（属门自己的一票）。
5. 形 B（`go/parser` 枚数钉）如果在 `internal/tools` 落地，`cmd/wisp` 那 5 枚范本的尺形可以抽成公共 helper，避免每票现造一遍遍历。
6. `:920-923` 文档块那句 "the binding is per scope, hence per task" 与 ② 的"scope id 可以⛔ 是 task id"（`panel-assets-l2`）之间有措辞张力 ⇒ 无论裁ⓐ还是ⓑ，那段文字大概率都得改一句，⛔ 别只当行号腐烂处理（`AC#3` 射程）。
7. `task.go:708-711` 与 `subagent_197.go:262-266` 两枚字面量共享前缀 `这条调用没有宿主给的任务 id（fail-closed：` ⇒ 若裁ⓐ，新拒绝的字面量选词最好与该前缀族一致，否则 `Contains` 型尺会误判（本腿⛔ 选词）。

## 六、本格交付状态

- 只交料、⛔ 裁形、⛔ 翻任何一格框、⛔ push；票面 `.scratch/wisp/issues/309-…md` 逐字未动（`git status --porcelain -- internal cmd` 起手＝0 行，交付前再核一次并记在下面）。
