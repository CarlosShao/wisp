
## A1007（2026-10-16 13:0x +08，编排者）收 326-v1 并翻框六枚结案 ＋ 收 165-r3 交件转交裁决 ＋ 第十波放五枚腿（CPU 44/42/43）＋★台面配方第二改（门禁⛔ 打在部分导出树上）

**① 起手锚**＝`44f51fc8`（台账 A1006 那批，12:32）；放行前实测 CPU 连发三发＝**44／42／43（%）**、进程名册除机主 `wisp.exe` **PID 21892** 外 `go.exe`/`*.test.exe`/`gofumpt*` 命中 0 ⇒ 机主的尺"只要 cpu 不高于 70%，就可以继续补位"成立，`270-r2` 那枚按住 **解除**（按住件 `470ae1e5` 现场作废并追加更正，⛔ 当"已完成"——它此前连一次启动都没发生过）。

**② 收件（完成通知一律当未验证；三把尺＝提交真在／件真在／我自己复跑）**
- **`326-v1`（裁决腿·非实现者·只读）四笔 `c2e1c07a`→`7c481326` 全在**；逐枚 `git show --name-only`＝写面四枚（`internal/memory/dao_session_326v1_test.go` 新件／probes/326/v1 17 枚／`docs/evidence/s1/326-v1-verdict.md` 17177 B／票尾 1 条 bullet `numstat 1/0`），越界 0；`find probes/326/v1 -size 0`＝**0 枚**、logs 15 枚。判语＝`AC#0`…`AC#5` **六枚全"成立"**、⛔ 一枚没勾（翻框归我）。
- **`165-r3`（写·`internal/agent`）十笔 `4f2470703`→`6b18b1d4` 全在**；它自己申报的越界尺命中 0 我**未逐枚复跑**＝该复算的账转给裁决腿 `165-v2`（本波已放）；它 ⛔ 翻框、票 165 现量仍 **5 勾／1 空**。
- ★**两枚腿各自交回一枚"我没做"**：`326-v1` 具名 D1/D2/D3 三枚交付缺陷；`165-r3` 具名断口②③＋一处派单措辞更正。**入账、⛔ 折进判语成立**。

**③ 我自己重跑（这是翻框凭据，不是它的转述）**
- 两包四数：`go test ./internal/memory ./internal/session -count=1 -timeout 30m -v`＝**`RUN=117 PASS_TOP=63 FAIL_TOP=0 SKIP=2 rc=0`**（`ok 16.135s`／`ok 2.234s`）⇄ 它申报 `80/43/0/2`＋`37/20/0/0` 求和**逐位同形**（件 `logs/orch-326v1-two-packages.txt`）。
- ★**摘牙那一发我自己做**：`git archive HEAD` 全树建 `D:\tmp\orch-326v1-replay`，物理删 `internal/memory/dao_session.go:300` 那行 touch ⇒ 两枚常驻判据各红（`:63` "the refresh error is being swallowed"／`:110` "premise broken … proves nothing"），`rc=1`；还原 `rc=0`；起手⇄还原 sha256 两串同为 `ec6ee530b190…1d2b2`＝**与它申报的还原号逐字同串**（尺＝`txt/replay-326v1-detooth.sh`）。⇒ **`AC#1` 那格我认它有牙**。
- 放宽断言尺：`git diff --numstat 86298e24 HEAD -- internal/memory internal/session`＝只列本裁决腿那枚新件 `125/0`、**删除行 0 枚**。
- 长度门（工作树，本批提交前）＝`rc=0 VERDICT GREEN`，分母 14289 跟踪路径／超帽 60＝在册 60／名册外 0。

**④ 翻框与结案**：票 326 **翻 `AC#0`–`AC#5` 六枚全勾**（尺＝`txt/flip-326.sh`：6/0 → 0/6，diff 只见 `[ ]`→`[x]`，numstat 6/6，判据文字一字未动）＋**`git mv` 改名 `-done`**（改名后最后一名组件 95 字、相对 116 字 ≤ 帽 121 ⇒ 不需进名册）。⚠ 按 README `:15` 的机械规矩结案；**README `:201` 那句警告照旧生效**＝残余不是零，只是**归口在票 166**（会话列表今天⛔ 接进产品）。

**⑤ 我裁的事（裁决者≠实现者，四条都记在票尾 bullet 上）**
1. 裁决腿新加的常驻判据 `dao_session_326v1_test.go` **留**；⇒ `internal/memory` 现行分母 **顶层 41→43、`RUN` 78→80**，后续任何腿引这两包的分母必须现跑。
2. 三枚上交料（`ORDER BY` tie-break／缺母行 ⛔ `ErrNotFound`／迟到戳）**⛔ 判为本票欠账**——tie-break 是 `AC#2` 票面自己写的 ⛔ 顺手治；其余两枚**待现量查重**再决定立票（队列尾 `329-c1`，⛔ 现在不写票面判据）。
3. `326-v1` 自报三枚交付缺陷入账〔D1 射程标注缺／D2 红句节选未标／D3 反形逐发颜色漏报〕，**不改判语成立**。
4. 工作树 `go.mod` **仍脏且不是任何腿的活**：`git diff -- go.mod`＝`github.com/jchv/go-webview2` 由 indirect 提为 direct，文件 mtime `2026-10-10 22:28:50`。编排者⛔ 动 `go.mod`（禁区清单），**具名留给机主**；今天全部门禁绿＝零行为后果，⛔ 由任何腿顺手 `go mod tidy`。

**⑥ 我自抓（台面配方第二改，同族第二次）**
- ★**门禁三数⛔ 打在部分导出树上**：`tools/d22scan` 的 ban #8 射程**声明了 `design/`＋`frontend/`**（`tools/d22scan/main.go:581-582`），一枚声明了的射程扫到 0 枚是**硬退出 `rc=2`**，长得像红。本波两枚腿各撞过一次且都**具名上报没算进判据**：`165-r3` 报 `scope design/ examined 0 files`、`326-v1` 报三枚读运行库的自扫用例红（树里无 `.git`）＋长度门 `not inside a git work tree`。⇒ 定式改成：**八件套部分树只用于 `go build`/`go test`；`d22scan`／长度门／`gofumpt` 只在工作树或 `git archive`（不带 pathspec）全树上跑**。已把这条**只追加**进 `270-r2` 派单（`2e706255`）并写进本波四枚新派单。
- 上一枚同类＝A1006⑥①"七件套漏 `models/`"。两次的共同成因＝**我只列了"跑起来要读的面"，没列"仪器自己声明的射程"**⇒ 今后写配方前先 `grep -n 'dir: filepath.Join'` 那把尺读仪器自己的射程。

**⑦ 编队与队列（13:0x 现量）**：**在飞 5 枚＝写 1（`270-r2`，`internal/tools` 注释面）＋只读 3（`327-a1`／`328-a1`／`216-a2`）＋裁决 1（`165-v2`，独占 `internal/agent` 编译面）**；派单存档＝`2e706255`（`txt/327a1-dispatch.txt`／`txt/328a1-dispatch.txt`／`txt/165v2-dispatch.txt`／`txt/216a2-dispatch.txt`／`txt/append-270r2-release.md`）。⛔ `cmd/wisp/**`＋`internal/panel/**` 本波零写腿＝两枚面同时空出，排在队尾的落点包（票 313 那包、`panel_pump.go` 工作区那半、票 252 残余）按"每波只许一枚写腿"逐枚放。待人项照旧只有 **Q-85**（机主给 30 分钟退掉 `build\wisp.exe` PID 21892 的那一发整包终态名册，默认不催）。**⛔ push**：写腿在飞、本地领先 `origin/dev` 未复尺。
