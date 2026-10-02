# e2e-panel-1 — 「从双击到看见设置页并存下去」跳表（只读取证腿）

> 本件唯一问题：机主那句「我要能看到主面板，我要点击设置，自己配置模型这些参数」到今天之间**差几跳**。
> 本腿**不裁任何 AC、不判任何票好还坏、不提任何修法**。逐跳四问：这一跳是什么／今天通没通／承重处 `file:line`／不通缺哪一环、归哪张票。
> 写点只有本文件。产码／测试／票面／`docs/**` 零字节改动。`frontend/**` 与 `design/**` 既不读也不引 ⇒ 本件**不出现**那两棵树的任何路径与行号；页面内容不是本腿射程。

## §0 起手锚（进场第一发，逐字）

```
$ date -Iseconds
2026-10-02T11:34:18+08:00

$ git log -1 --format='%H %ad %s'
21a15d2704d72e7469423b72aaf2741c899614c4 Fri Oct 2 11:33:27 2026 +0800 ledger(A529)＋票 254 结案改名 -done：…（完整标题见 git log）

$ git rev-parse --abbrev-ref HEAD
dev

$ git status --porcelain -- cmd internal | wc -l
0
```

- 代号复认：`.scratch/wisp/probes/` 起手 listing 共 98 枚，无 `e2e-panel-1`；本目录由本腿新建（⛔ 只建不删）。
- ⚠ 共享工作树整树非空是常态，本腿只在 `cmd internal` 两枚根上取分母；`0` 是**起手时刻**读数，不是恒量。
- 中途复量（写 §4/§5 之前）：`date -Iseconds` = `2026-10-02T12:11:34+08:00`，`git status --porcelain -- cmd internal | head -20` **仍为空**。
- ⚠ `248-v1c` 正在往跟踪文件种临时变异 ⇒ 本件的 `file:line` 一律**双口径**给出：工作树读数 + `git cat-file blob HEAD:<path>` 读数。两把尺已对 `panel_config_store.go`／`panel_host_windows.go`／`manager.go` 三枚承重件逐条复认（见 §6）。

## §1 跳表（本腿正身）

## §2 今天已经通的跳（调用点尺＋真实读数）

## §3 今天不通的跳（缺哪一环＋归哪张票，⛔ 不含修法）

## §4 我可能判错的条目

1. **「通」的尺可能偏宽**。本件判「通」＝**非测试调用者 ≥1 枚**。但面板那条链上多枚承重件带 `windows` 构建标签（`cmd/wisp/panel_host_windows.go:1`、`cmd/wisp/panel_resident_windows.go:1`、`cmd/wisp/resident_windows.go`），非 Windows 构建图里它们**不存在**。我按 D7（平台范围＝Windows）把这族判成「通」，如果编排者要的尺是「跨平台可编译可跑」，我这几格一律要降级。
2. **行号可能在我读数之后漂移**。起手 `cmd internal` 分母为 0，12:11:34 复量仍为 0；但 `248-v1c` 的种变异窗口不落 `git status` 之外的时间。我只对三枚文件做了 HEAD-blob 双口径复认，其余（`bridge.go`／`config_handlers.go`／`settings.go`／`resident_ball_windows.go`／`main.go`／`loader.go`／`manager.go` 的段表）**只有工作树单口径**。
3. **`res.Tier` 那枚行号我读成 `:228`，派单与票 255 的 AC#5 都写 `:227`**。两把尺（工作树 grep 与 HEAD blob grep）**一致给 228**（且 `:227` 是 `res.Written = written`）。我倾向认为不是漂移而是票面/派单少一行，但也可能是我对「`:222-229` 前面没有任何按键/按段判断」这句话的射程理解偏窄——那一段里确实有一枚 `switch w.Field`（按键分派），只是它不分派**档位**。这条留给编排者定文字，我不改票面一字。
4. **「机主能存下去」的判据我按落盘取**。`writeOneKey` 同时改内存 `m.cur` 与文件（`internal/config/settings.go:237-246`），我据此判「落盘＝通」。若编排者的尺是「下一次任务真的用上这个值」，那一跳我今天只在 §3 里列成跳，没有判它通。
5. **6249 这个数字我怀疑但不敢用**。我用仓库级 census 量到 embed 目标目录下 **跟踪件 1 枚**（复认派单锚 3），同一目录口径下 **ignored 件 6249 枚**。我只数了路径尾带 `/dist/` 的条目，**没有**（也不能）点名那枚目录，因此**无法排除**其中混有第三方的 `dist/`。这条直接决定 §1 里「页面产物」那一跳是「本机已备、干净检出仍缺」还是「本机也缺」——见 §5 第 2 条。
6. **`config.get` 的回话形状我判成「一句人话，不是字段集」**（`internal/panel/config_handlers.go:445-462` 产 `renderSettingsView` 字符串，宿主把 `Handle` 的返回值原样交给绑定回话）。如果机主那句「看到设置页、自己配参数」**不要求**结构化字段（页侧把这句话渲染出来就算满足），我把它列成跳是**过度解读**；反过来如果我列轻了，那是把 C17 契约面漏了。这一格我不裁，只标尺不确定。
7. **常驻进程里有几枚 `config.Manager` 我判不动**。尺上 `config.NewManager` 非测试调用点是 `cmd/wisp/panel_inbound.go:230`（常驻面板链经 `newComposerDispatchChain` 走它）与 `cmd/wisp/run.go:409`；但 `run.go:409` 所在的 `assembleRuntime` 只在 `interactiveStdin()` 非 nil 时才被常驻腿调到（`cmd/wisp/resident_task_source_windows.go:218-231` 那一支 `return nil`）。⇒ **纯双击（无控制台）时常驻进程可能只有 1 枚 Manager**，这时「两枚各说各话」不成立，成立的是「任务管线根本没装配」。我没跑，判不了，列进 §5 第 4 条。
8. **票面 `-done` 状态我只看了文件名**。本腿只 `ls` 出六枚票（33/145/198/248/253/255）都**不带** `-done` 后缀，池内另有 84 枚带后缀。`-done` 是防重领键、不等于验收结论，我不拿它当「好/坏」判据，也**没有**逐枚读六张票的 AC 勾选态（33 的票面我只读了 AC#13/#14 与裁定段，见 §5 第 5 条）。
9. **锚 1 里「窗口尺寸写死在 `:304-306`」**：我复认 `:304 Width: 420,`／`:305 Height: 260,`，`:306` 是 `},`。射程一致，只是分界行差一枚；我在跳表里按 `:304-305` 写。

## §5 判不动／没测到的地方

1. ⛔ **本腿零跑**：`go build`／`go vet`／`go test`／`./...`／`go list -deps` 一律未执行（`248-v1c` 正在整包跑 `cmd/wisp`）。⇒ 本件所有「通/不通」都是**调用点枚数＋码内分支阅读**的读数，**不含任何执行证据**；凡是「只有真跑才能定」的格子，我在 §1 的「通」列里写的是**接线态**，不是**运行时态**。
2. **页面产物在运行时到底是哪一档，我测不到**。可读的码侧事实：`internal/panel/assets.go:54-60` 用 `fs.Stat(tree, "index.html")` 决定 `built`；`:75-78` 未 built 时 `Resolve` 直接 `errNotBuilt`；`cmd/wisp/panel_host_windows.go:342-344` 供页失败就走 `serveNotBuiltNoticeLocked`（`:356-365` 那枚自造告示页）。至于**本机今天**embed 里到底有没有入口字节，需要跑一次 `panel.BuiltinAssets()` 或读那枚目录的文件名——前者被禁、后者属禁区。⚠ 现存两枚**互相冲突的台件读数**（都在 `cmd/wisp` 与 `docs/evidence/s1/` 里，不是我推的）：一枚记录工作树形如「有页面包」，另一枚的具名 skip 文案写着「embed 解析不出入口，本节在本树无对象」（`cmd/wisp/panel_resident_windows_test.go:317`）。我**不裁**哪枚对。
3. **「回执真到不到页面」测不到**。码侧接线在：`cmd/wisp/panel_host_windows.go:319-322` 绑定回话＝`dispatchRaw` 的返回串；`cmd/wisp/panel_resident_windows.go:250` 把线程泵交进库的 `Run()`（`Run()` 是派发队列唯一读者——这句是 `:21-27` 注释自陈，我没读第三方模块源码复认）。判「到达」要一发真页面 `await`，本腿无此射程。
4. **常驻进程运行态的对象数量判不了**（§4 第 7 条）：面板写入用的是链上那枚 Manager，热加载 tick 只在 `cmd/wisp/run.go:813 rt.startConfigReload()` 装（属 `assembleRuntime`），而 tick 报的那句「这些段已立即生效」只写 `rt.stdout`（`cmd/wisp/config_reload.go:171`）。⇒「谁在什么进程里读到面板刚写的值」我没跑，判不了。
5. **票 33 的 AC 勾选态与裁定 2/3/4 的落地程度我没读全**（票面 342 行，我只读了 AC#13/#14 与随后裁定段）。凡我把某一跳「归票 33」的地方，归的是**那一格的名字与射程**，不是「这格还没做完」的判语。
6. **界面侧我一律不量**：设置页存在与否、页面上有没有输入框、`wispDispatch` 有没有被页面调、页面里那批方法名字面量——`frontend/**`／`design/**` 禁令覆盖，具名口径按票 253「需要界面侧配合的那一件由 owner 带话」执行，本腿不转达、不推断。
7. **凭据面我只写字段名/常量名**：`credentialValue`（`cmd/wisp/panel_config_store.go:53`）、`FieldProviderCredential`、`provider_api_key_ref`、引用形状 `dpapi:`/`env:` 前缀名。⛔ 零凭据值，零本机 blob 名。
8. **`d22scan` 会不会因这棵树今天的形状判红，我判不了**：`tools/d22scan/runtests.sh:98-102` 的尺我读到了（`skipped != 0` ⇒ `exit 1`，「SKIP 不是过」），但**它今天到底跳过哪几枚**只有跑才知道。派单锚 3 那句「AC#13 因此必跳」我**只能确认尺的形状**，不能确认它今天真跳。
9. **票 145（快照字段）在本跳表里的射程**：我量到快照泵的非测试构造点只有 `cmd/wisp/run.go:699`（`panel.NewSnapshotPump`），常驻面板链**没有**泵；`Credential` 那枚读者挂在 `cmd/wisp/run.go:710 rt.settings.credentialStatus`。⇒「面板窗口能收到快照」这一跳我列进 §3，但**票 145 的题面是载体字段**，泵装配的归口我判不动，标在这里。
10. **未定义即停（本腿没有自填的格）**：我没有碰到需要新造规矩的情况；遇到「两套说法相反」「两枚读数冲突」一律**上报不裁**（§4 第 3 条、§5 第 2 条、§6）。

## §6 我推翻编排者哪一句

## §7 一句话结论（按依赖顺序列跳，不列修法）
