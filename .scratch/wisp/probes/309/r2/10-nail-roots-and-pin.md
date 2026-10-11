# 309-r2 · 10 枚数钉的射程与落点（job 1＋job 2）

件：`internal/tools/toolrequest_nail_309_test.go`。判据⛔ 是"我记得的规矩"，下面每一句都带现跑的尺。

## 1. roots：`tools` 进，`scripts` 也进（两把尺的读数）

尺＝`logs/roots-rulers.log`（R359-1）＋`logs/package-agent-ruler.log`（R359-2），14:5x +08：

**① module 路径前缀尺**（`find . -name go.mod` 里三枚 nested 模块逐字）
```
go.mod                      :: module github.com/CarlosShao/wisp
scripts/spike/go.mod        :: module github.com/CarlosShao/wisp/scripts/spike
tools/d22scan/go.mod        :: module github.com/CarlosShao/wisp/tools/d22scan
tools/mockllm/go.mod        :: module github.com/CarlosShao/wisp/tools/mockllm
```
⇒ 三枚都落在 `github.com/CarlosShao/wisp/` 那一截**之内**。Go 的 `internal/` 可见性按 **import 路径前缀**判定（父目录是 `.../wisp` ⇒ 同族可 import），
所以"它是独立模块 ⇒ 进不来 `internal/agent`"⛔ 成立——编排者 14:1x 那句自纠（§2 里 `nail-roots-miss-tools-dir` 那枚）就是这一条。

**② 活例尺**（R359-3）＝`grep -rn "CarlosShao/wisp/internal" tools scripts --include=*.go`＝**1 命中**：
`tools/signmodels/main.go:34` 逐字 `"github.com/CarlosShao/wisp/internal/models"`。
⇒ nested 模块伸手进 `internal/` 在本仓**已是事实**（⛔ 假想）。`scripts` 自己今天 0 命中，但①已经放行它：
⛔ 用"谁此刻还没伸手"当排除理由＝把接线事实当结构保证，那正是本票在防的那族病。
⇒ **裁（本腿按两把尺读出来的结论）＝`tools` 与 `scripts` 两枚都进 roots。**

今天各 root 的非测试源枚数（尺＝`find <root> -name '*.go' -not -name '*_test.go' | wc -l`）：cmd 40／internal 229／tools 13／scripts 11。

## 2. 排除支换成"具名文件＋具名枚数"

- 旧：`strings.Contains(rel, "internal/agent/")` ⇒ `internal/agent` 的 **24 枚**非测试源一枚⛔ 解析；重试包装器最自然的落点恰在那包里。
- 新：⛔ 排任何目录，断言 **全 roots 生产构造点合计 == 1 枚**，且**具名**那枚落在 `internal/agent/loop.go`。
  ⇒ 有人往 `internal/agent` 里再加第二条构造点**也会响**（实测见 `20-nail-shapes-and-rigs.md` 的 RIG-4b：枚数 1→2、名册逐枚具名、rc=1）。

分母照旧打印且**逐 root 打印**：旧钉 245＝(cmd＋internal 269) − internal/agent 24；
新钉现量 **293**＝cmd 40＋internal 229＋tools 13＋scripts 11（钉自打印⇄第三把 `find` 尺逐 root 对得上，读数 `logs/r359-7-nail-roster.log`）。
新增两枚 root 各带一条"root 名册为空即 `t.Fatalf`"（先例逐字＝`task_state_188_test.go` 那句
*"a scan whose denominator is zero reports clean for a tree it never read"*），⛔ 一枚空名册当"扫过了"。
全局地板从 `<10` 提到 `<200`，且注释写明它是**地板⛔ 是常数**。

## 3. matcher 认的三形与⛔ 认的形（钉头逐字列了，这里给读数）

认：限定名字面量／同包裸名字面量（⛔ 只在被解析文件的 `package` 子句＝`agent` 时放行）／可解析到 `agent.ToolRequest` 且**无初始值**的 `var` 声明／同文件类型别名再写字面量。
尺＝R359-2：`grep -rl '^package agent$' cmd internal tools scripts` 去目录后＝**只有 `internal/agent`** ⇒ 裸名那一支今天⛔ 会从别处混进分母。
⚠ 已知偏差（写进注释）：那一支认的是 **package 子句**⛔ 是目录；别处真出现一枚叫 `agent` 的包，它自己的 `ToolRequest{}` 会进分子（多算，⛔ 少算）。

仍漏（⛔ 抄成"全覆盖"，钉头具名列了七条）：跨包别名／函数体内别名／改名 import 与 dot-import／`new()`·`make()`·反射·解码构造／定义型＋转换／省类型的嵌套元素／（外层是 ArrayType 时两向都计 0）。
其中两条各有一发实测绿读数当证据（RIG-5 省类型、RIG-6 跨包别名，各 rc=0），其余五条⛔ 造 Rig＝本腿欠读数，登记在 `90-final.md`。

## 4. 与票 154 门尺的关系（⛔ 揽、⛔ 改）

`gate-clauses.sh` 的 **G1** 那格仍登记 `want_n 0`，它问的是"internal/agent **之外**几枚"；本钉现在问"全 roots 共几枚、那枚在哪个文件"。
两问今天对得上（G1 侧＝0，本钉侧＝1 且那枚在 agent 里）。⛔ 动 `gate-clauses.sh` 一字（那 5 处绝对 `loop.go:N` 锚⛔ 在 `AC#3` 点名段里＝编排者 §2 已判"⛔ 本票账，归票 154/151 族"）。
