# 派单存档·2026-10-23 09:4x +08 编排者（机主批准三件后首批）

★**机主原话（逐字）**：「好的，这三件事我都允了，同意，你可以继续了」
⇒ ★**批准范围**：①票 348 加"删除那一跳"的产码接缝（**第二道人工批准已获**）②票 362 改 `.github/**` ③票 196 `D35` schema 约束。

★本文件是**三份派单的存档**（`348-r2`／`362-r1`／`363-r1`）。

---

## ① 票 348 续跑腿 `348-r2`（写腿，★**产码接缝已获第二道批准**）

**身份**：写腿（实现者）。⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ push。
**票面**：`.scratch/wisp/issues/348-*.md`（★**只剩 `AC#2`**）。

**★★机主已批准（逐字引用进交付件）**：「好的，这三件事我都允了，同意，你可以继续了」
⇒ ★**本腿获准在产码里加"删除那一跳"的接缝**——★**这是 `348-r1` 明确说"必须第二道人工批准"的那一件**。

**★`348-r1` 已量的两条成因（⛔ 重做，但★**你要自己复核**）**
- **(a) 祖先 `rename` 三态**：空手 `err=<nil>`／按名开着叶子 `Access is denied`／经 `os.Root` 开着叶子 `in use by another process`
  ⇒ ★**本包唯一声明的缝（`var applyDescriptor`，那时叶子句柄正开着）装不下"open 后祖先才翻"** ⇒ **要新缝**。
- **(b) 给外来对象装 `DENY(WRITE_DAC)` 后 seal 仍 `err=<nil>`**（属主隐式 `WRITE_DAC` 压过 deny ACE）⇒ **造不出"封不动的外来对象"**。

**★`348-r1` 已交的半边（⛔ 重做）**
- **POSIX 全满足**：①修前 >0（确定性红 ＋ 竞态 `deleted=76/1500`、`62/1500`）②修后 0（`0/1500`）③同发证窗口到达链接态（`provedLinkCycles=267/256`）。
- **Windows 只有 ②③**：②`deleted=0` ×4；③`provedReparse=75..81`、`formBLegRefusals=32..37`；★**①git 基线取不到**（`deleted=0` 但 `truncated=13/4/7` 证窗口到达过）。
- ★**①′摘跳树给了 ①>0**：`mut`＝把 `unlinkPinnedLeaf` 换回 `_ = os.Remove(...)` ⇒ **8 发 5 红**（`deleted=3,1,1,0,0,0,1,1`）。

**`AC#2` 要的（逐字见票面）**：成对三件——**①修前分子 > 0**、**②修后 0**、**③同一发证明窗口真到达过链接态**；★**两平台各一发**（⛔ 一发代两发）。

**做法（硬）**
1) ★**加接缝**（现在获准）：让"删除那一跳"**可被测试替换**（★最小面、⛔ 顺手重构）；★**先量后写**——加缝**之前**先给 Windows ① 的读数（应仍取不到）、加缝**之后**给 ①>0。
2) ★**三读数 × 两平台**（POSIX 已有基线可复用、Windows 要新做）＋ **反向判据**（摘掉那一跳 ⇒ 必须红并**逐名报红行**）。
3) ★**还原核**：★核 hash 前先问"谁算的、什么算法、几位"（`git hash-object`＝SHA-1／`sha256sum`＝SHA-256、⛔ 混比）。
4) ★**格式甲乙两尺都报**（甲＝blob 判据／乙＝工作树只诊断）。

**写面**：允许 `internal/winsec/**`（★**产码接缝 ＋ 测试**）＋ `docs/evidence/s1/348-r2-*.md`。
⛔ `internal/{tools,config,session,memory,risk,panel,agent}/**`（归别的腿、写面互斥）；⛔ `cmd/**`、`scripts/**`、`.github/**`、`tools/d22scan/**`、`go.mod`/`go.sum`、`docs/PLAN.md`／`docs/specs/**`、`docs/reports/**`、`frontend/**`／`design/**`、`thresholds.go`／golden／`allowlist.txt`／三枚冻结件。

**★★平台陷阱（本仓已确认）**：`cmd/wisp` 测试二进制在**仓外导出树**上会因缺 DLL 起不来（`0xc0000135`）⇒ 导出树跑 `cmd/wisp` 必须带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`。
★**跑门禁前先 `unset GOTMPDIR`**（`GOTMPDIR` 落非 `C:` 卷会让 `d22scan` 的 311 守卫 `rc=2` 拒跑、测试进程一次没起）。

**交件**：`docs/evidence/s1/348-r2-impl.md`（≤160 行）：★**机主批准逐字**／★**接缝加了什么（逐行）**／
★**Windows ① 加缝前 vs 加缝后两发读数**／★**POSIX 三读数（复用或重跑、要标清）**／★**两平台各一发（⛔ 一发代两发）**／反向判据＋还原核／门禁 rc。
**未定义即停**：①若**接缝必须动到 `internal/risk`／`internal/tools`**（越界）⇒ 停手上报；②若**加缝后 Windows ① 仍取不到** ⇒ 停手上报并给读数。

---

## ② 票 362 落地腿 `362-r1`（写腿，★**改 `.github/**` 已获批准**）

**身份**：写腿（实现者）。⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ push。
**票面**：`.scratch/wisp/issues/362-*.md`（★**五格，`AC#0` 已翻**）。

**★★机主已批准改 `.github/**`（逐字引用进交付件）**：「好的，这三件事我都允了，同意，你可以继续了」

**★编排者已拍的形＝〔甲′〕改调用面、⛔ 洗任何字节**（依据见票面 `AC#1` 的拍形块）：
- ★**把 `ci.yml:172` 那步的 `gofumpt -l . tools/d22scan tools/mockllm` 改成按目录集传参**：`gofumpt -l internal cmd tools scripts frontend`
  （★`tools/d22scan`／`tools/mockllm` **已含在 `tools` 里**）。
- ★**决定性读数（我现跑过）**：`gofumpt -l .` ⇒ **20 枚**、★**非 `.scratch` 的 ＝ 0**；`gofumpt -l internal cmd tools scripts frontend` ⇒ **0 枚**；**根目录无 `.go`**。
- ★★**必须同笔改那步注释**：它逐字写着 "**the denominator of the formatting gate is the TRACKED set, not the working tree**"
  ⇒ ★**改成目录集之后，"分母"从"全部 tracked `.go`"变成"产码目录集"** ⇒ ★**这是语义变化、必须在注释里写清**。
- ★**第 3 条也要写进注释**：**若将来有人的台件放在产码目录里（今天 0 枚）⇒ 会被扫**。

**⛔ 硬禁**
- ⛔ **洗任何 `.go` 字节**（★20 枚台件的价值**就是"故意坏"**；★其中 `probes/241/v1/posctl/badly_formatted.go` 与 `probes/259/r1/gofumpt-negctl/probe.go` 是**负控样本**、洗掉会**毁掉它们被造出来的目的**）。
- ⛔ 动 `.scratch/**` 以外的任何产码；⛔ 动 `internal/**`／`cmd/**`（本腿只改 CI 文本）。

**写面**：允许 `.github/workflows/ci.yml`（★**只那一处调用面 ＋ 那步注释**）＋ `docs/evidence/s1/362-r1-*.md`。
⛔ 其余一切。

**做法（硬）**
1) ★**先量**：改前在**仓外导出树**上跑 `ci.yml:172` 那步的现形（应列 20 行全 `.scratch`），★**留逐字**。
2) ★**改**：只改那一行 ＋ 注释（★**行号要现跑复认、⛔ 信票面**）。
3) ★**改后再跑**：同一树跑新形状 ⇒ ★**应 0 行**；★**再造一枚"产码目录里故意坏的 `.go`"（仓外）⇒ 新形状必须列它** ⇒ **证明新形状不是瞎的**。
4) ★**反向判据**：把新形状改回 `.` ⇒ **必须重新列出那 20 枚**。

**交件**：`docs/evidence/s1/362-r1-impl.md`（≤160 行）：★**机主批准逐字**／★**改了什么（逐行 diff）**／
★**改前 20 行逐字**／★**改后 0 行**／★**"产码目录里故意坏枚"被列出的证据**／★**反向判据**／★**注释改后逐字**（含"分母语义变了"那句）。

---

## ③ 票 363 落地腿 `363-r1`（写腿，★形已拍〔丙〕）

**身份**：写腿（实现者）。⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ push。
**票面**：`.scratch/wisp/issues/363-*.md`（★**`AC#0`／`AC#1` 已翻／已拍**，只剩 `AC#2`／`AC#3`／`AC#4`）。

**★形已拍＝〔丙〕两者都做**（依据见票面 `AC#1` 的拍形块）：
- **〔甲〕修日志**：让 `MODE-REDLINE … still asks` 在 **`L0`** 与 **`Deny`** 两个方向都**说真话**。
- **〔乙〕修 `Screen`**：让它**认得 `L0`**（★`L0` **是认得的级**：`assessor.go:41` `L0 < L1 < L2 < Deny`、`:45` "read-only, no side effects — execute directly"、`:46` `L0 Level = iota`；
  ★★**最硬的一条**：`Screen` **自己在 `mode.go:175`／`188` 返回 `Silenced{Level: L0, …}`** ⇒ **同函数把 `L0` 当输出级认得、却在入口判成"无法识别"**）。

**★`363-a1` 的三条承重实测（⛔ 重做，但★**你要自己复核**）**
- **L0 那发真读数**：真 `fs.read`、三个档各一发，逐字 `tools: MODE-REDLINE mode=<m> still asks tool=fs.read rules=[] because=无法识别的风险级，fail-closed 保留询问`；
  ★**同一发**：`CARDS TOTAL=0` ＋ `kind=success decision=allow outcome=success` ＋ `EXEC text="hello"`（**已执行完**）。
- **★`Deny` 上同样假**：`fs.read` 打 `%APPDATA%\wisp\config.toml` ⇒ `level=Deny`、`CARDS TOTAL=0`、`text="该目标被绝对禁止访问（R3 A 档…）"`、**而日志仍打 `still asks`**。
- **正控**：`L1/ask_every_step` ⇒ **1 张窗口卡且没有那行**；`L2` 两档 ⇒ **各 1 张审批卡且说的是真的**。

**★〔丙〕的硬边界（⛔ 越）**
- ★**⛔ 为了让日志"好看"而放宽任何判定**：`L0` **仍必须放行**（`bridge.go:432-435` 的 `case risk.L0: … return true, ""` **不动放行语义**）、`Deny` **仍必须拒**。
- ★**⛔ 把 `Screen` 改成"认得 L0 之后就不管了"**——`L0` 要**明确**走出「不询问」这一支。
- ★★**`Deny` 那条要先读清 `Kept` 契约**（`mode.go:143-146`：只在**红线拒绝静默**时非空）再动。

**写面**：允许 `internal/risk/mode.go`（`Screen` 一形）＋ `internal/tools/bridge.go`（日志一形）＋ 两包测试 ＋ `docs/evidence/s1/363-r1-*.md`。
★★**两包必须串行或同一腿**（★本腿一枚、两包都归它 ⇒ **⛔ 别的腿碰这两包**）。
⛔ `cmd/**`、`scripts/**`、`.github/**`、`tools/d22scan/**`、`go.mod`/`go.sum`、`docs/PLAN.md`／`docs/specs/**`、`docs/reports/**`、`frontend/**`／`design/**`、`thresholds.go`／golden／`allowlist.txt`／三枚冻结件；
⛔ `internal/{winsec,config,session,memory,panel,agent}/**`（归别的腿、写面互斥）。

**做法（硬）**
1) ★**先量**：两向各给"今天是什么读数"（⛔ 假设）。
2) ★**每处改动配一枚今天响、不修不响的用例**；★**反向判据**＝把新加的那一跳摘掉 ⇒ 必须**转红并逐名报红行**。
3) ★**行为无回归**：★**`L0` 仍放行、`L1` 仍窗口、`L2` 仍审批、`Deny` 仍拒** ⇒ 四档各一发读数。
4) ★**格式甲乙两尺都报**。

**交件**：`docs/evidence/s1/363-r1-impl.md`（≤160 行）：★**〔甲〕改了什么**（`L0` 与 `Deny` 两个方向的日志逐字）／★**〔乙〕`Screen` 改了什么**／
★**四档行为无回归读数**／★**反向判据三读数**／★**正控（`L1` 不动那行仍没有、`L2` 说的仍是真话）**／门禁 rc。
**未定义即停**：①若修 `Deny` 那条**必须动 `Kept` 契约语义** ⇒ 停手上报（那可能是改契约文字）；②若**行为有回归** ⇒ 停手上报。

---

## ★三腿共同的开工前必读

`.scratch/wisp/probes/orch/20261022/txt/ruler-scope-checklist.md`（**73 行**、**10 处「尺的射程」的错**）
**落笔前五问**：①这把尺量的是哪个对象？②输出形状？③我数的是「代码」还是「注释」？④这个数有没有和总数对上？⑤**谁算的、什么算法、几位？**
★判 `gofumpt`／`gofmt` 输出归属**必须用 `[\\/]`**；★判 CRLF **必须 `tr -cd "\r" | wc -c`**；★**名册取主检出**。
★**跑门禁前先 `unset GOTMPDIR`**。
★★**在工作树做突变 ⇒ 突变窗口记进件、⛔ 窗口期提交任何东西**；★**优先用 `-overlay` 或仓外导出树**。
