# 票 235 裁决表 —— 235-v1（非实现者验收腿）

- **验收人**：编排者之外的非实现代理（腿号 `235-v1`）。**本表一行实现码都没写**：零产码、零测试码改动，
  所有突变只经 `go test -overlay` 施加，每发之后 `git status --porcelain -- internal cmd` 的读数逐条记在
  `.scratch/wisp/probes/235/v1/readings.md` §5。
- **被验对象**：票
  `.scratch/wisp/issues/235-the-pool-nail-comment-reason-was-falsified-by-ticket-222-and-only-one-of-three-cases-sees-the-real-bridge.md`
  的四格判据，交件腿＝`235-r1`。
- **起手锚点（本腿现跑 `git rev-parse --short HEAD`）**：`98df640a`。
  被验实现腿的五枚 commit＝`795ed767`／`1d2ad737`／`a71f0be9`／`feb88b91`／`9c96f103`；
  改动面（`git diff --name-only 25556879..HEAD -- internal cmd`）＝
  `internal/tools/subagent_197_test.go`＋`internal/tools/subagent_222_test.go`，**零产码**（本腿现跑对上）。
- **行号口径**：票面／台件里的旧行号一律只当线索。下表"票面行号"＝**本腿在锚点 `98df640a` 现读的该 AC 判据所在行**。
- **读数标记**：〔我现跑〕＝本腿自己执行的命令；〔读台件，未复跑〕＝只看 `235-r1` 的现场件。两者不混。

## 0. 逐格判据（本表正文见 §1-§4；骨架先落，逐格裁完一格追加一节）

| 格 | 票面行号（现读） | 判语 | 凭据（本腿读数所在节） |
|---|---|---|---|
| AC#1 改注释、结论与断言零字符变更；M6 抬池仍是有效检测力证明 | `:19` | 待裁 | §1 |
| AC#2 三枚 222 用例的前置读数；M3 下三枚全红且红不来自 30s 护栏 | `:20` | 待裁 | §2 |
| AC#3 两枚自证（① 摘掉前置读数退回 1/3；② 注释改回旧句是否有仪器保护） | `:21` | 待裁 | §3 |
| AC#4 门禁与终态读数（四数只从 `-v` 量／gofumpt／vet／d22scan／逐名双向差集） | `:22` | 待裁 | §4 |

## 1. AC#1（票面 `:19`）——**成立**

**判语**：注释段的活是真的、约束也真守住了；末句"抬池⇒这枚钉子仍红"是一枚**有效的检测力证明**，
"4 枚红"这个报法**不算夸大**（但里面有一枚是仪器自造的阻塞，见下面第 5 条）。

凭据（全部本腿现跑，台件 §0／§2／§4）：

1. **删除列只落在注释行**〔我现跑〕：`git diff -U0 25556879..HEAD -- internal/tools/subagent_197_test.go`
   ＝ `-` 五行全是 `// ` 开头的旧注释、`+` 十三行全是 `// `，函数名行与其下两个 `t.Errorf` 一枚字符都没进 diff
   （`--stat` 同向：13 insertion／5 deletion）。比较方向 `MaxConcurrentSubagents > MaxToolConcurrency` 原样（现读 `:412`）。
2. **假理由真的从盘上消失**〔我现跑〕：`grep -rniE "bridge slot for its child|whole life|bigger pool is a lie" --include=*.go internal cmd`
   ⇒ `subagent_197_test.go` **零命中**（只剩 `internal/agent`／`internal/audio`／`bridge.go` 三处无关件）。
3. **历史句逐字在位**〔我现跑〕：`grep -n` ⇒ 命中 **`:407`**（`which is exactly the 7/8-red shape ticket 197 leg A reported.`），
   同段还把它降级成出处（下一行就写 "Raise MaxConcurrentSubagents and this leg still goes red"）。
4. **新理由的载荷句有行为尺托着，不是故事**〔我现读〕：注释里 "a waiting parent holds none of the four permits"
   ＝`Test222WaitingParentHoldsNoBridgeSlot` 的 `len(h.bridge.sem)` 判据（终态 0.00s PASS，本腿 `-count=3` 那发 18 名全绿）。
   ⚠ 一处**措辞级不准**：注释写 "giveBackWhileWaiting in `subagent_197.go`"，本腿现读＝**调用点**在 `subagent_197.go:365`、
   **定义**在 `bridge.go:694`。指调用点不算假，但严格按盘读应当是 "`subagent_197.go` 里的那次调用"。不影响判据。
5. **M6（抬池⇒红）本腿自己重造、自己读数**〔我现跑，⛔ 工作树零突变字节〕：overlay 只把产码常量
   `MaxConcurrentSubagents = 4` 改成 `8`（一行），终态其余一字未动 ⇒ `=== RUN` 218／全名 218／`--- FAIL` **4**，
   与交件报的红名册**逐枚同名同秒**：三枚 197 钉各 **0.00s**（正文 `:413` "池 8 大于桥的 D38d 天花板 4…名册却说它们在跑"、
   `:452` 同族）＋ `Test222SpawnConclusionArrivesThroughRealBridgeChildren` **3.00s**。
   - **有效的那三枚**：`Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)` 正是注释那句 "Raise MaxConcurrentSubagents and
     **this leg** still goes red" 的正面证明 ⇒ 牙还在，这条支撑成立。
   - **3.00s 那枚的性质**〔我现读正文，不是信注②转述〕：它红在 `subagent_222_test.go:488`（"结论…出现在 0 枚父任务回复里"），
     四枚父任务收到的是 `工具 task.spawn 超时（3000ms），已协作式中止`，而 `3000ms` ＝该台件自己的
     `h222PreFixBudget = 3 * time.Second`（`:81`，喂给 `DefaultTimeout :272`），gate 又要求 8 枚孩子同时进桥、桥只有 4 枚许可
     ⇒ **确证是测试构造的阻塞**，与 `222-v1` 注② 同判。
   - **"报成 4 枚红"是否夸大＝不夸大**，三条理由：① 票面现量 #3 自己就是这个数（"4 枚红＝三枚 197 钉＋一枚 222 用例"）；
     ② 交件同一段里**具名限定**了它不当"池 8 有害"的证据（不是把四枚并列当四条理由用）；
     ③ 代码注释正文里被写进去的只有 "this leg still goes red" 那一枚 0.00s 钉，四枚的数只出现在"顺带指出处"的句子里。
     ⚠ 唯一可挑的一处**读者风险**（不影响本格判语）：注释末三行把 "four legs red" 写进了源码头，
     下一位读者如果不去点 `docs/evidence/s1/222-…v1.md` 的注②，可能把这四枚读成四条"池 8 有害"的证据。


## 2. AC#2 —— 待裁

## 3. AC#3 —— 待裁

## 4. AC#4 —— 待裁

## 5. 我攻不动的地方

（逐格裁完后回填，不许为空）

## 6. 没做完／留给编排者

（逐格裁完后回填，不许为空）
