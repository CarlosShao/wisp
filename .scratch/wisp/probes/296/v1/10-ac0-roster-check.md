# 296-v1 / AC#0 —— 射程普查名册复核判语：**成立**（两处具名缺陷，不推翻结论）

腿：`296-v1`（非实现者对抗验收）。HEAD＝`28b484ec`。被裁件＝`296-a1` 的
`.scratch/wisp/probes/296/a1/10-ac0-roster.md`＋`00-anchor-and-rulers.md`（写于 `d432d728` 时刻，改产码之前）。

## 1. 本腿自己现跑的尺（⛔ 不采信名册的枚数）

| 尺名 | 命令逐字 | 读数 |
|---|---|---|
| V-R1 | `git grep -n 'ball\.HotkeyConfig{' HEAD -- cmd internal \| grep -v '_test\.go'` | **6 行**：`cmd/balldebug/main.go:239`（已过合并）、`cmd/wisp/resident_windows.go:49/:53`（`hotkeyReloadSource296` 两支，均已过合并）、`:223/:230`（`hotCfg258` 两支，裸、由消费者补）、`internal/ball/hotkey_reload.go:18`（注释示例） |
| V-R1pre | 同尺跑 `ff5c193e^`（落地前一刻） | **6 行**：`balldebug:239`／`resident_windows.go:188/:195/:213/:215`／`hotkey_reload.go:18` ⇒ **与名册 R1 的 6 行逐字相同**（含行号） |
| V-R2 | `git grep -n 'ApplyHotkeyDefaults' HEAD -- cmd internal \| grep -v '_test\.go'` | HEAD 真调用 4 处：`balldebug:239`、`resident_ball_windows.go:235`、`resident_windows.go:49`、`:53`；余为注释/日志文案 ⇒ 名册 R2 那句"改前 `cmd/wisp` 真调用只有 `:235`"对 `ff5c193e^` 成立（本腿另尺复核） |
| V-R3 | `git grep -n '\.Hotkey\b' HEAD -- cmd internal \| grep -v '_test\.go'` | 5 行：`balldebug:238`、`resident_windows.go:53`、`:225`、`hotkey_reload.go:17`（注释）、`hotkey_windows.go:81`（注释）、`internal/config/manager.go:281`（config 侧重载赋值，**不喂 ball**，不属本族） ⇒ **第二把尺（字段读点）与第一把尺（复合字面量）交叉覆盖，名册 6 枚真落点没有第三种形状** |
| V-R4 | `git grep -n 'NewHotkeyReloader' HEAD -- cmd internal \| grep -v '_test\.go'` | 生产装配点 **2 枚**：`cmd/balldebug/main.go:237`、`cmd/wisp/resident_ball_windows.go:353`；定义/注释在 `internal/ball/hotkey_reload.go:16/:70/:73` ⇒ 与名册 A3/A5/B1 一一对应，**无第三枚宿主** |

⚠ 本腿所有尺的射程＝`cmd`＋`internal` 两棵树、**排除 `_test.go`**（与名册 R1 同射程）；注释行**包含**（V-R1/V-R2 都命中注释行并单列档）。
含 `_test.go` 的同尺读数（只为"如果射程换一把该看到什么"而跑，不计入名册枚数）：`git grep -c 'ball\.HotkeyConfig{' HEAD -- cmd internal` ⇒ 4 枚测试文件共 12 行
（`resident_hotkey_258_windows_test.go` 6／`resident_hotkey_296_windows_test.go` 2／`resident_hotkey_v1probe_test.go` 4／其余 0）。

## 2. 三问逐答

**① 名册是不是整族？——是。** 两条独立尺（V-R1 字面量／V-R3 字段读点／V-R4 装配点）在"改前"与"HEAD"两个时刻都落在名册已列的同一批位置：
`resident_windows.go` 两枚闭包（A1 构造源、A2 热加载源，**各含两枚 `return`**，名册逐枚点名 `:188/:195` 与 `:213/:215`）、
`cmd/balldebug/main.go:237-242`（A3）、`resident_ball_windows.go:316`（A4）、`:353`（A5）、`internal/ball/ball_windows.go:813`（A6）。
名册还多给了"消费者补默认"的落点 A4 的**第二重保险**（`ball_windows.go:153-154` 只在四格全空时生效），这一枚我复读到、且它支持"机主那形救不回来"的判语。

**② 每枚枚数有没有写清"哪把尺＋射程目录＋含不含 `_test.go` 与注释行"？——写了，但有一枚口径撞名（具名缺陷 1）。**
名册配套件 `00-anchor-and-rulers.md` 把 R1／R1b／R2／R3／R4 五条命令**逐字**落盘，射程（`cmd internal`／`.scratch`）、
`grep -v _test`（排除测试文件）、变异拷贝排除枚数（`.scratch/wisp/probes/258/v1/mut0|mut3` **8 枚**，逐条路径行号）都齐 ⇒ ② 合格。
撞名处：R1 的原始读数是 **6 行字面量**，名册的"真落点"也是 **6 枚**（A1..A6），**两把不同的尺都叫"6 枚"**（前者按 `return` 行、后者按站点/消费者）。
本仓第 4 次抓的正是这一形（票 294 收口节 `startResidentBall` 12 行＝真调用 10＋注释 2）。名册靠逐枚 `file:line` 仍可重建，**不推翻结论**，但今后枚数须自带尺名。

**③ 有没有漏掉票面点名的落点？——没有。** 票面 AC#0 点名的四组全部在册：`:181/:205` 两枚闭包（A1/A2）、
`cmd/balldebug/main.go:237-241`（A3）、`internal/ball/hotkey_reload.go:13-20` 示例注释（B1）、`cmd/wisp` 别处字面量（R1 全量）。
**balldebug 反例是具名独立复核的**（名册第四节表格该行写"本腿自读，未采信转述"），本腿 V-R1/V-R2/V-R4 三把尺在 HEAD 复现同一读数：
`cmd/balldebug/main.go:239` 逐字 `			return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{` ⇒ **"同一个坑 balldebug 没有、只有 wisp 有"成立**。

**具名缺陷 2（不推翻结论，处置归编排者）**：名册第四节复跑票面"测试里 0 枚空串夹具"那把尺时用的是 `grep -rn ... --include=*_test.go cmd/wisp`，
它今天仍成立，但**同一目录下另有 `cmd/wisp/resident_hotkey_v1probe_test.go`（票 258 的 v1 探针件，含 4 枚 `ball.HotkeyConfig{` 字面量、
`Test258V1Probe*` 两枚用例在整包名册里跑）**；名册因射程声明排除了 `_test.go` 所以不算漏，但**"哪些测试件常驻在 `cmd/wisp` 里替这个形状说话"这一档没有单列**，
下一枚改 `[hotkey]` 链路的腿容易以为只有 258/296 两件。

## 3. 判语

**`AC#0` 成立**（整族、尺具名、票面点名无遗漏、balldebug 反例三方复核成立）。
两处具名缺陷＝①"6 枚"撞名不带尺名；②`cmd/wisp` 测试族那一档未单列。**均不改变任何判据。**
