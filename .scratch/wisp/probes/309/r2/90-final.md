# 309-r2 · 90 收尾账（写腿，⛔ 裁决者）

起手 HEAD＝`202ed4a5`（派单给的锚，现跑一致）；交回 HEAD＝本腿最后一笔（`git rev-parse --short HEAD` 现跑，见 `logs/r359-14-handin.log`＝`d75db9ee`＋本证据笔）。
其间 `dev` 上落过 `308-r2` 的六笔（`6f389389`…`daa2f9be`），其中两笔改的是 `internal/risk/*_test.go` ⇒ 本腿所有整包名册/导出树读数与它们**同面⛔ 同树**，具名在 `logs/behaviour-vs-anchor.stdout` 的 V1/V1b。

## 1. 五件各自的一句话＋承重的尺

| 件 | 做了什么 | 证明它的尺（读数） |
|---|---|---|
| ① roots | `tools` 进、`scripts` 也进（具名答复） | `logs/roots-rulers.log`：三枚 nested module 前缀全在 `github.com/CarlosShao/wisp/` 之内；活例尺 1 命中＝`tools/signmodels/main.go:34` |
| ② 落点 | 排除支→"合计 1 枚＋具名 `internal/agent/loop.go`" | 钉自打印 `构造点=1（==1…）分母=293（cmd 40／internal 229／tools 13／scripts 11）`⇄第三把 `find` 尺 293 对得上；RIG-4b 往 agent 里加一条⇒`rc=1`、名册 2 枚、分母 294 |
| ③ 三形 | 补"声明＋赋值"与"同文件别名"；裸名按 `package` 子句放行 | `nail-rig-1/2/3.log` 各 `rc=1`（旧钉 2/3 形⛔ 绿）；六枚 control 双向在场（含"别的包也叫 ToolRequest"断 0）；`nail-rig-5-elided.log`／`nail-rig-6-crosspkg-alias.log` 各 `rc=0`＝具名漏形的读数 |
| ④ 注释 | 三处说实话＋具名 roots 判据＋漏形清单 | `logs/r359-10-final-rulers.log`＝71／69／2，剥行首 `//`＝1；`logs/r359-9-filterdemo.log`＝同枚 fixture 两把口径 2⇄1；R5＝`sed -n '211,213p'` 逐字 |
| ⑤ prov 肢 | ⛔ 造仪器，注释里另记一段（族／谁响／⛔ 具名仪器） | `logs/prov-rigs.stdout`：MUT-F（摘 `injected`）整包 `rc=0`；MUT-E（摘 `prov == nil`）整包 `rc=1`、顶层 16 枚红＝管道噪声；R359-4＝`Assessor:` 普查 2 行⇒生产写点 0 |

行为逐字节不变＝`git diff -U0 202ed4a5..d75db9ee -- internal/tools/bridge.go` 增 44／删 9 **全为 `//` 行**、非注释 0；两枚 blob 剥整行注释后 `cmp` **IDENTICAL**（sha `496ed6d6…`，854⇄854）；守卫那行 1⇄1 逐字同一；CR＝0、`i/lf w/lf attr/text eol=lf`。
门禁＝`rc_build=0`／`rc_vet=0`／`rc_d22scan=0`（两轮红⛔ 别腿：ban #9 phantom-citation，我自己写的 `cmd/internal/tools/scripts` 与一枚未落盘的件名）／`rc_pathlen=0`。
整包名册（两棵导出树，同面可比）＝顶层 210/1/0 ⇄ 210/1/0、含子测试 278/1/0 ⇄ 278/1/0、`comm` 作差 **added 0／removed 0**、`SKIP` 两向 0；两枚 236 空 id 载体两棵树照绿。

## 2. 新增欠读数 6 枚（⛔ 我盖章哪格该翻，判语归 `309-v2`／编排者）

1. **`Q-309-r2-word-ruler-still-self-referential`** — `ToolRequest{` 词面尺今天 71／69／2；注释⛔ 消灭这个咬尾，只是把它**写进尺的口径**（自带"按文件名字段滤＋剥行首 `//`"两支，并指解析器尺）。⇒ 残余＝"任何新注释都可能再挪一次分子"，⛔ 由注释负责，靠解析器尺兜。
2. **`Q-309-r2-nail-blind-shapes-unrigged`** — 钉头具名七条漏形，本腿只给两条装了实测读数（省类型嵌套／跨包别名）。其余五条（函数体内别名、改名 import 与 dot-import、`new()`／`make()`／反射／解码、定义型＋转换、别处一枚叫 `agent` 的包）＝**具名⛔ 实测**。⛔ 装成"已量过"。
3. **`Q-309-r2-pinned-total-is-constant-1`** — 新钉把"合计 1 枚"当常数钉住。将来若有**正当**的第二条生产构造点（带非空 id 的面板/球宿主），钉会红且红句给的是"接上非空 id 或进仪器入口"，⛔ 自动授权第二枚。归口＝下一次有人真要写它时由编排者裁（同 `AC#0` 形 ⓑ 的口径），⛔ 本票先放宽。
4. **`Q-309-r2-scripts-roots-admitted-by-prefix-only`** — `scripts` 进 roots 的根据只有"前缀放行"，它自己今天⛔ import `internal/`（尺＝1 命中全在 `tools/signmodels`）。⇒ 若将来有人把 `scripts/spike` 拆成 `wisp` 前缀**之外**的 module，`internal/agent` 就⛔ 可达，本钉仍会扫它（⛔ 错，只是射程变宽）；反之若把 `tools`/`scripts` 挪出前缀⛔ 更新注释，注释里那句"三枚都在前缀之内"就过期——它是**现读句**，⛔ 是保证。
5. **`Q-309-r2-tree-rigs-mixed-face`** — RIG-4b／6 的树＝`git archive 6f389389` ＋ `cp` 本腿⛔ commit 的两枚文件 ⇒ **混面**（钉是新码、其余是 `6f389389`）。commit 之后可逐字重跑（`git archive d75db9ee` 就⛔ 需要 cp），那一发本腿⛔ 补跑（时间轴在收尾），⛔ 读数算已固化。
6. **`Q-309-r2-prov-limb-colour-discrepancy`** — 派单 §5/票面 §2 那句"摘掉 `b.prov == nil` 或 `b.injected` 两发**各 rc=0**"是 `-run` 只打 309 仪器那把尺；**整包**那把尺下 `prov == nil` 摘掉＝`rc=1`。两把尺⛔ 矛盾，但"哪把尺支撑哪句判语"必须说清：支撑"`⛔ 具名该肢的仪器`"＝仪器那把；支撑"这肢有⛔ 任何反馈"＝整包那把，**结论⛔ 同**。⇒ 交 `309-v2` 定一句口径；本腿注释按"⛔ 具名＋整包会噪声地响"两半都写。

## 3. 越格自报（⛔ 瞒）

- **写面**：代码笔 `git show --name-only d75db9ee`＝`internal/tools/bridge.go` ＋ `internal/tools/toolrequest_nail_309_test.go`（2 枚，都在授予的 `internal/tools/**` 内）；其余笔只落 `.scratch/wisp/probes/309/r2/**` 与本票末尾**一条** bullet（1 加／0 删）。⛔ 翻框／⛔ 动 `Status:`／⛔ 动别家票／⛔ 动 `docs/**`。
- **多做的一件**（超出五件清单，⛔ 产品语义）：`assessorFor` 文档块里那句 `Test309... pins agent.ToolRequest{ construction points outside internal/agent to 0` 在 job 2 之后**会变成假话** ⇒ 一并改成"合计 1 枚、具名 `loop.go`、旧版整包跳过"的新口径（纯注释，落在 B1/B2 那两把尺里：44/9 全 `//`、剥注释 blob 同一）。
- **多跑的 go 命令**（都留了件）：`go test ./internal/tools/ -run 'Test309…'` 若干发；`-run 'Test309'` 一发；整包 `-count=1 -v` 三发（工作树 1 ⇄ 导出树 2）；突变台内整包 2 发（MUT-E/F，当场还原）；`go build ./...` 4 发、`go vet ./internal/tools/` 4 发、`sh scripts/d22scan.sh` 4 发（含两轮我自己的红）、`check-path-length-budget.sh` 3 发。**⛔ 跑过 `cmd/wisp` 那包的测试**（sherpa DLL 前置归编排者）。
- **多写的文件**：`logs/` 下 24 枚（脚本 `.txt`＋读数 `.log`／`.stdout`＋两枚 `.stripped`＋三枚 `.go.txt` 快照＋`msg-*.txt` commit 消息＋`tree-*.dir`）。⛔ 仓库根、⛔ `.out` 名（注：`logs/r359-*` 等一律 `.log`／`.txt`）；⛔ 在 `probes/**` 留 `.go`——中途落过一枚 `logs/bridge-head.go`，**同一次运行内**改名 `.go.txt` 并把脚本改成导到 `/tmp`，现 `find .scratch/wisp/probes/309/r2 -name '*.go'`＝**0**（尺见 `logs/r359-14-handin.log`）。
- **⛔ 删过任何临时件**：突变全部是"还原⇄⛔ 删除"（trap ＋ `cmp` ＋ `sha256sum`）；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；⛔ push。

## 4. 我认为编排者那份裁定里⛔ 成立的一处（顶回来请复跑再判）

§4 第 4 件与 §2 签收的 `Q-309-v1-prov-nil-limb-blind` 用的是"**两肢各 `rc=0`**"这个读数；整包那把尺下只有 `b.injected` 成立，`b.prov == nil` 是 `rc=1`（16 枚噪声红）。
⇒ **处置⛔ 变**（照样⛔ 造仪器、照样注释另记），变的是那句"同族同色"的说法：本腿注释写成"一肢全盲／一肢只有噪声式反馈、⛔ 具名仪器"。若编排者复跑 `logs/prov-rigs.txt` 那两发对得上，请把 §4 第 4 件的措辞按两把尺分口径记进台账；对不上就撤我这处。
