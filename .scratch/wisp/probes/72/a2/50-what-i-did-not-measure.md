# 72-a2 — 本腿没量的那些事（⛔ 沉默；每条都给了下一程该用哪把尺）

## 1. 逐枚第 4 问（"会不会恒真／无信号"）里凡标〔结构判断〕的都⛔ 量过

本票⛔ 给变异授权（派单 §3 编译面边界逐字＝"⛔ 建树做变异"），且同包此刻有 `307-v1` 在飞。
⇒ 7 枚的"把判据换成反形它会不会响"全部**只按代码结构答**，一句都⛔ 当读数。
**下一程该用的尺**＝照票 307 `AC#2` 那三形，在**仓外导出树**上做（先例＝`307-r1` 的 `183032c1` 那笔"四形变异全在 `C:/Users/swq/tmp/wisp-307-r1-tree-main`"）：
- 形 (a)＝`pathresolver.go:345 formsOf` 退回"只交 raw 拼写"（票 72 `AC#3①` 原物）；
- 形 (b)＝`blacklist.go:115 overrideApplies` 改回遍历全部形（票 72 log `:171-173` 原物）；
- 形 (c)＝**本腿点名唯一一枚结构上就怀疑它没牙的**＝#2 `TestCanonicalInputGainsNoSecondForm`（见 `30-answer-B-line.md` 形 III 与 `40-nails` 第 4 条），它的判据是"枚数==1"，(a) 那一形很可能让它⛔ 再红。

## 2. ⛔ 跑整包，所以"`internal/risk` 只有这 7 枚红"⛔ 是本腿的读数

本腿四发全是具名 `-run`（名册 7 枚一把、两枚 windows 文件全 16 枚一把）。
⇒ 编排者 `A834` 那四发（`198/197/0/1`→`200/199/0/1`、`198/185/12/1`→`200/192/7/1`）本腿**只复用了"改后 aliased 剩 7 枚＝甲族"这一件事的形状**，⛔ 复现其分母。
**下一程该用的尺**＝`go test -C <目录> ./internal/risk/ -count=1 -timeout 300s -v` ＋ 对落盘件 `grep -cE '^[[:space:]]*(=== RUN|--- PASS|--- FAIL|--- SKIP)'`，**先断言名册行数**再谈颜色（`logs/t03-broken-surface-note.txt` 记了一发"PASS=0 却⛔ 全红"的 panic 假读数）。

## 3. ⛔ 复现乙族（⛔ 应属本腿射程）

`syncdirs_test.go` 那 5 枚与 `syncdirs.go` 的两支比较一⛔ 看、⛔ 读、⛔ 判。派单 §0 与本腿尺 2/3（`20-answer-A-cause.md` A.2 表的第 2、3 行）给的是**隔离性**，⛔ 是乙族的状态。

## 4. ⛔ 判"票 72 的修复是不是真有效"

本腿只证到"别名台面上那 7 枚红在写法侧、同族另外 9 枚绿"（`logs/aliased-all16-bothfiles.txt`）。
⇒ "票 72 的不变量已经成立"这句**是推论⛔ 是裁决**，归 `307-v1` 同级的非实现者终裁；且 `AC#4` 的 **runner 半边仍欠一次真 run**（票面 `:78` 逐字"runner 侧要 push 才有一次真 run，而本票明写'绝不 push'"），⛔ 一枚本地等价台面都⛔ 能替它（票面 `:75`）。

## 5. ⛔ 量 CI runner 的两把台面色

⛔ push（AGENTS §1.4／票面 `:112`），且 `A813` 那条"推送自启 `slo-full` 抢机主 CPU"仍⛔ 用（票 307 追加节 `:183`）。
**下一程**＝runner 色归编排者那一发，要与票 300／303／111 欠的 CI 色**同批作差**（`40-nails` 第 10 条）。

## 6. ⛔ 量"这台机器哪些卷开着 8.3"

`fsutil 8dot3name query C:` 需提权（先例＝票 307 追加节 `:127` 记的 `Error 5` 具名欠读数）。本腿只量到**行为**：`RISKAT~1\RUNNER~1` 与 `riskattrib1\RunnerTempAreaForAttribution1` 两形现在都在盘上（`ls -d` 各一发），且那 16 枚里 `--- SKIP`=0 ⇒ `getShortPath` 交回过非空。
⇒ 落地腿若要**自己种世界**，判据⛔ 是"这台机器开着 8.3"，⛔ 是"夹具自己造得出短形否则 `t.Fatal`"（票 307 追加节 `:132`）。

## 7. ⛔ 读 `307-v1` 的裁决件正文（它把一枚读数指回了本票族）

`ed1d2bad` 的标题逐字含"(d) **3/5 无信号裁归票 72**"。⇒ 那是**别的腿的裁决**，⛔ 是本腿的料；本腿只登记它存在并⛔ 去读 `probes/307/v1/**`（射程外＋它是裁决件）。**下一程必须先自己读它**再决定落地形状，否则会撞"⛔ 拿别人的判语当自己的读数"。

## 8. ⛔ 量"`Resolve` 对缺失叶子⛔ 展开祖先"这一条到⛔ 算缺陷

`20-answer-A-cause.md` A.3 第 1 条只登记**行为**与它的**方向是 fail-safe**（少放＝多问）。
要判它该⛔ 该修，需要两把本腿⛔ 打的尺：
- `git show <ref>:internal/risk/pathresolver.go` 与 `docs/specs/SPEC-06-security-gatekeeping.md §4` 那句"展开 8.3 短名"逐字对拉，看契约⛔ 管到"⛔ 存在的那一节"；
- `git log -S 'risk.Gate(' -- internal/tools/mode.go` 一枚 caller 长在哪个 commit（`A.4` 那句"票面 `:245-247` 已过期"目前只有**现状**证据，⛔ 有**何时变成现状**的证据）。

## 9. 卫生自检（本腿自己的，⛔ 交给别人代量）

本腿交件共 **6 笔**（逐笔尺＝`git show --name-only --format= <sha>`，先 `tr -d '\r'` 再数）＝
`7db4db4b`（锚）→ `493f2cb8`（七枚）→ `a6267fbb`（必答 A＋16 枚件）→ `25189588`（必答 B＋钉）→ 第 5 笔（本文件＋`90-handoff.md`）→ **第 6 笔＝`logs/t05-final-hygiene.txt`（前五笔的逐笔尺合账，落完第 5 笔之后才跑得出来）**。
⚠ 下面那张表量的是**前四笔**（第 5 笔⛔ 能在自己落笔前量自己，⛔ 让这件事变成一枚⛔ 存在的读数）。

| 尺 | 前四笔的读数（本腿现跑） |
|---|---|
| 写面外命中（逐笔 `grep -vc '^\.scratch/wisp/probes/72/a2/'`） | 四笔各＝**0** |
| 新建 `.go`（逐笔 `grep -c '\.go$'`） | 四笔各＝**0** |
| 碰过 `.scratch/wisp/issues/**`（逐笔 `grep -c 'issues/'`） | 四笔各＝**0** ⇒ **翻框 0 枚**这件事有逐笔尺，⛔ 只是自陈 |
| `.out` 件 | `find .scratch/wisp/probes/72/a2 -type f -name '*.out'`＝**0 行** |
| 0 字节件 | `find .scratch/wisp/probes/72/a2 -type f -size 0`＝**0 行**（带 `-type f`，先例＝票 307 追加节 `:92`） |
| 仓根临时件 | 0（中间件全落 `probes/72/a2/logs/` 或 Windows `%TEMP%`；⚠ 本腿有一发把票面 dump 写成 `"$OLDPWD/…"`，`$OLDPWD` 当时为空 ⇒ 落到了 Git Bash 根挂载点 `/tmp307.md` 之类，**⛔ 落进仓**，尺＝`git status --porcelain \| grep tmp307`＝**0 命中**；教训＝**重定向目标⛔ 用可能为空的变量拼**） |
| `git status --porcelain -- internal cmd` | 每发 Go 面之前各一把＝**0 行 ×4**；本腿⛔ 动过 `internal/`／`cmd/`／`frontend/`／`design/` 任何一字节 |
| 最长件名 | `50-what-i-did-not-measure.md`／`40-nails-that-would-break.md`＝**28 字符**（尺＝`find … -type f -printf '%f\n' \| awk '{print length, $0}' \| sort -rn \| head -3`）；⛔ 新立票 ⇒ 那条雷⛔ 触发 |
| push | **0 笔**（尺＝本腿⛔ 运行过任何 `git push`；`git log --oneline -1 origin/dev` 一类的比对也⛔ 做，免得失手） |

## 10. 一处本腿⛔ 干净的地方，具名自报（⛔ 等别人抓）

`logs/aliased-all16-bothfiles.txt` 这一枚件名本腿**用了两遍**：第 3 发（搭错台面那发）的正文被第 4 发覆盖。
⇒ 违反的是"临时件只建不删"的**精神**（那份读数事实上消失了）。
⇒ 处置＝第 3 发的原始读数已在本腿交件里逐字保了一枚 `logs/t03-broken-surface-note.txt`；定式＝**一发一个件名，别学本腿**。
