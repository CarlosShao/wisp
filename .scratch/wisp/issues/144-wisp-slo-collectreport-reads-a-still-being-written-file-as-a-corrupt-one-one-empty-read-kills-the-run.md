# 144 — `wisp slo` 收主体报告时把"还没写完的文件"读成"损坏的文件"：一发空读直接判死，整枚 `slo-full` 因此在并发下会红

- Status: ready-for-agent
- 来源：编排者 09-25 12:5x 归因 `ci` run `36094258734`（sha `cd87354`）里 `slo-full` 那枚**新增红**（基线 `36068900302`/`a1fd5bf` 这一道是 **success**）
- 关联：票 128（"拒绝启动"那批 fail-closed 裁定）、票 136 `AC#14`/`AC#15`（覆盖行能否决 pass 位）、`A235⑤`/`A237②`
- 地界：只改 `cmd/wisp/slo_windows.go` 里 `collectReport` 那一段＋`cmd/wisp/**` 新增一枚测试；**阈值 / golden / `thresholds.go` / `internal/observe/**` 一字节都不许动**

## 现量到的形状（我自己逐行核过源码，不是引 run 日志的转述）

`cmd/wisp/slo_windows.go:523-545` 的 `collectReport()`：

```go
if data, err := os.ReadFile(s.outPath); err == nil {
    var rep sloRun
    if err := json.Unmarshal(data, &rep); err != nil {
        return nil, fmt.Errorf("wisp slo: subject report: %w", err)   // :528-529  ← 一发即死，不重试
    }
    ...
}
```

⇒ 只要**文件已经存在但字节还没写完**（`os.WriteFile` 是先建后写，读者可以看见 0 字节或半截），
`json.Unmarshal` 拿到的就是 `unexpected end of JSON input`，这一支**立刻返回错误、不回预算里重试**。
CI 上那一发的原文读数就是这条链：

```
wisp slo: in-tree record unavailable (fail-closed, no silent downgrade to the tree basis): wisp slo: subject report: unexpected end of JSON input   （:383 那一层包住 :529）
slo-check.ps1: state Armed exit=2 pass=False
slo-check.ps1: report written to …/slo-report.json (all_pass=False)      ← 六态里只有 Armed 这一行红，其余 Warm/Conversation/PanelOpen/WorkPeak/settle 全 pass=True，leak 负控 flipped_to_fail=True 正常
```

⚠ **两处不对称是问题的核心**，本票就修这个不对称：
① 同一个文件里紧挨着的 `waitReady()`（`:503-518`）用 `fileHas(s.readyPath, "ready=1")` —— **它容忍"文件还没内容"**，在预算里轮询；
② `collectReport()` 不容忍，把"还没写完"与"内容真的坏了"**读成同一件事**。
而 `:520-522` 那句注释逐字写着 *"The subject only writes after subjectGrace past the end of its window, so this **never races** the observer's last sample"* ——
**这句话没成立**：它论证的是"不会撞上观察者的最后一发采样"，但**没有保证"读者不会在看文件时文件是空的"**（存在与写完是两回事）。
⇒ 这是一枚**"永远产不出它自己想量的读数"式**判据的姊妹形：注释把一次没验证的时序假设写成了不变式。

**方向要说清**：fail-closed 本身是**对的**（宁可红也不静默降级，票 128/136 那一族就是这规矩）；
本票**不放宽它**，只要求它**先分清"没写完"与"坏了"**——分不清的时候，红要红得有理由。

## AC

- [ ] **AC#1 先把这一形钉成能响的用例**（不许只加断言）。判据：构造"文件存在但内容为空/半截 JSON"的那一状态，
      白盒可（本仓同风格先例：直接问纯函数/小装配器），要求
      ⓐ 在预算内 ⇒ 它会**重试并最终读到完整报告**（或预算到点后给出**点名"预算与最后一次读到的字节数"**的错误，而不是 `unexpected end of JSON input`）；
      ⓑ 内容**确实是坏的**（非空、且不是"没写完"形状）⇒ **仍然当场红**（这一支不许被 AC#2 顺手放宽）。
      ⚠ 若你判断这一形在本机/CI 上都**造不出可重放的用例**（只能靠真实进程调度竞态），那就**诚实登记为"打不开"**并写明为什么，
      本票按"AC#1 自证作废"收——那比硬开一道永不响的 AC 好（`A234` 那批刚为此退回过程序）。
- [ ] **AC#2 修法只许一个方向：把"没写完"留在预算里重试，把"坏了"当场红**。
      两者必须**分得开且有凭据**（例如按"字节数为 0 / 尾部不完整"与"能解析出对象但内容不对"分开）。
      ⚠ **不许**改成"解析失败就 `continue` 到预算尽头"而不留最后一发的读数；**不许**加 `t.Skip`、**不许**动任何阈值/预算常量当修法。
- [ ] **AC#3 变异自证两向**：(i) 造一发"主体写的报告确实坏了"（非空、语义错）⇒ 必须红且红因点名；
      (ii) **摘掉 AC#2 新加的那一味**⇒ AC#1 那枚用例必须从绿变红（按本仓"承重"定义：摘掉它若什么都照旧，它就是装饰，本票作废并登记为什么）；
      并另问一句承重公式可能给误导答案的那条——**摘掉它，有没有任何外部可见读数变过**（`wisp slo` 的 exit 码、`slo-check.ps1` 那行 `state … pass=`）。
- [ ] **AC#4 门禁与名册**：`go test ./cmd/wisp/` 改前改后各一次（四数之外**点名册差集**；出现 `panic` 时那几十条记成"未取到"，不许改成 Skip）；
      `gofmt -l cmd/wisp/` 空；`go vet ./cmd/wisp/` 空；`sh scripts/d22scan.sh` rc=0
      （⚠ `tools/d22scan` 是独立 module，别在根目录 `go vet ./tools/d22scan/`；bans #1-5 不含 `_test.go`，只有 ban #8 含）。
      ⚠ **`GOOS=windows` 那半边**：`slo_windows.go` 只在 windows 编译，本机是 windows 所以能跑；
      若你为它补了 `!windows` 的对照腿，那条腿**只能在 linux 容器里量**（`golang:1.27` 在本地可用；
      Git Bash 下 `docker -v C:\…` 会静默挂空且 rc=0＝假绿，**必须当场证明挂载不为空**）。
- [ ] **AC#5 一句必须写进代码注释的话**：把 `:520-522` 那句 *"never races"* 改掉——它今天不成立。
      新措辞要说明**分得开的是哪两种状态**、以及**为什么"文件存在"不等于"报告完整"**。

## 规矩

只写 `cmd/wisp/slo_windows.go`、`cmd/wisp/**` 下新增的那一枚测试、以及 `docs/evidence/s1/144-slo-report-partial-read-r1.md`。
**禁碰**：`internal/panel/**`（门钉 r3 在写）、`tools/d22scan/**`（142 验收在地，含 `allowlist.txt`）、`internal/observe/**`、
`thresholds.go`、golden、`scripts/slo-check.ps1`（那是门禁脚本本体，改它＝改判据，须另批）、`docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`。
逐格 commit（每裁完一格提交一枚），显式 pathspec；commit 前现核 `git diff --cached --name-only`，出现别家路径停手报回；
禁 `add -A`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`push`；变异落仓外副本、临时件只建不删。
⚠ **不做任何计时类测量当判据**：此刻这台机器上可能有 3 程在跑测试、而 `slo-full` 就跑在本机的 self-hosted runner 上——
**任何"多少秒"的读数在这种编队下不作凭据**；要取整条 SLO 链的颜色可以，取秒数不行。
证据件先测后写再提交；"某节落在哪枚 commit"只认 `git log`/`show --name-only` 现量。

## 这一枚票**不**解决的事（登记，别顺手扩大）

- 它**不**让 `Armed` 那一发变绿之外还多保证什么；六态里其它五态今天本来就是绿的。
- 它**不**是"slo-full 每次推送必红"那件事（那是 `Q-36`/`Q-47` 那一族，owner 已明示"不改吧"）。
- 它**不**动 `subjectReportBudget`/`subjectGrace` 等任何预算常量——**用放宽预算当修法＝放水，直接退回**。
