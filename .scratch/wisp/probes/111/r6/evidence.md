# 票 111 AC#11 - 写腿 `111-r6` 证据件

范围（只有这一格）：给 `.github/workflows/ci.yml` 的 `test-windows` 腿加一步
`go vet -tags winlive ./cmd/wisp/ ./internal/ball/`，让带 `windows && winlive` 的
12 枚文件在 CI 里至少被编译面检查一次。零开窗、只编译、不执行。ⓐ/ⓑ/ⓒ 三半判据
逐一取数。所有读数件在同目录 `logs/`。

---

## §0 起手锚 + 在飞冲突检查

- 起手锚（`git rev-parse HEAD`）：`8e98b8f099fcdab75eaec42c55bca8ff7696f490`（短号 `8e98b8f0`）。
- `date`：`Thu Oct 8 08:55:26 CST 2026`。
- 分支：`dev`。
- `ci.yml` 起手行数：见 `logs/anchor.txt`。

在飞冲突（起手 `git status --porcelain` 摘，全部属别人、我一字未动、未 stage）：
`.gitignore` 为 ` M`；`design/**` 有 ` M`/` D` 一批（base.css / icons.js / theme.js /
tokens.css / index.html / screens/*.html 等）；`.scratch/wisp/probes/**` 下有他腿日志在改；
`.scratch/` 里一堆 `??` commit-msg 临时件。另：我这几笔 commit 之间被别的腿的 commit
`cfc9a5ff`（"35-v4 交付"）插进来一次——共享工作树、正常，不是我造的、我没碰它。
⇒ 我只 stage 我自己造的件（`ci.yml` + `.scratch/wisp/probes/111/r6/**`），每笔 `git add`
与 `git commit` 都带显式 pathspec，无 `-A`/无 `.`/无 `--amend`/`reset`/`rebase`/`stash`/
`checkout .`/`clean`。

编排者 §2 三条现量前提，落笔前各自复跑（输出见 logs）：
1. 作业/runner 名册 `logs/ruler-jobs.txt`：`lint`(65,ubuntu)/`test-core`(395,ubuntu)/
   `test-windows`(505,**windows-latest**)/`slo-smoke`(739,windows)/`slo-full`(797,
   `[self-hosted, wisp-slo]`)/`lint-frontend`(871,ubuntu)。与转述一致 ⇒ 落点 = `test-windows`。
2. 新步必须带 `if: ${{ !cancelled() }}`（编排者自踩的坑）。已带，逐字见 §1。
3. `if:` 两把尺 `logs/ruler-if.txt`：作业级 `if:` = **0**、步级 `if:` = **12**（其中
   11 枚 `!cancelled()`、1 枚 `always()`）。与转述一致 ⇒ 我只加步级 `!cancelled()`，
   未给任何 job 加 `if:`、无 `continue-on-error`、无 skip。
   ⚠ 注意（自抓、非结论性）：严格 8 空格尺 `'^        if:'` 只数到 11（漏了 :462 的
   `if: always()`，那行虽显示 8 空格但被该尺漏计），故本腿用 `'^\s+if:'`（任意前导空白）
   这一把得 12，才与 §2 转述的 12/11 对得上；两把尺都在件里留了原文。

winlive 名册 `logs/roster-winline.txt`：尺 `grep -rlE '//go:build.*winlive' --include=*.go`
现量 = **12 枚**（`cmd/wisp/` 7 + `internal/ball/` 5），与 AC#11 与转述一致；12 枚全部
带 `//go:build windows && winlive`（逐枚标签行见本次会话记录）。ci.yml 里改前
`-tags` 0 命中、`winlive` 0 命中（`logs/ruler-tags.txt`），唯一的 `go vet ./...`(:238)
在 ubuntu（`windows` 维即排除这 12 枚）⇒ 这 12 枚此前在 CI 里连编译都没人查，AC 描述成立。
逐枚 `git hash-object` 基线在 `logs/roster-sha.txt`（供 §3 还原自证）。

---

## §1 加了哪一步（逐字）

只新增一步，放在 `test-windows` 腿内、紧跟既有 `run: bash scripts/wisp-cli-tests.sh`
(:593) 之后、既有的 `Package coverage census`（下移到 :634）之前。其余步骤、其它 job、
排序：一字未动。完整 commit diff 存 `logs/ciyml-step.diff`（3909 bytes / 39 增行）。

关键三行（文件当前行号，逐字）：
```
630:        shell: bash
631:        if: ${{ !cancelled() }}
632:        run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/
```
⇒ §2 前提 2（"必须带 `if: ${{ !cancelled() }}`"）自证：第 631 行原文即上。

步名与注释块（`name: "winlive compile gate (go vet -tags winlive, ticket 111 AC#11)"` +
一段英文注释）逐字见 `logs/ciyml-step.diff`。注释里具名写死了射程：
compile coverage only；明确 `go vet` 不产测试二进制、不开窗、不需要 DLL；
并明确这步**不是**票 35 `:52`/`:75(c)` 真窗读数的 CI 载体（那按 `A690` 走"仅本机可量"另一支）。

为何没加 `scripts/` 包壳：这一步就是一条 `go vet` 命令，与本文件既有的内联
`go vet ./...`(:238/:241/:332) 同形，无需新件；引入脚本反而是多一枚要维护的跟踪件，
不满足"只有确实必要才造包壳"的边界，故未造。

---

## §2 ⓐ 当前读数（rc 与红名册）

命令：`go vet -tags winlive ./cmd/wisp/ ./internal/ball/`（本机 windows 宿主，唯一能取此数处）。
输出件：`logs/vet-base.txt`（含出处：go1.27.1 windows/amd64）与 `logs/vet-winline-base.txt`（裸 rc）。

读数：**rc = 0（绿）**。go vet 成功时零输出，故件里 `rc=0` 行是自落读数。
⇒ 今天的 HEAD 上这一步是绿的，**没有红名册要入台账**；且我**未**为让它绿动过那 12 枚
文件里的任何断言（⛔ 边界遵守）。

⚠ 出处限制（具名）：这是**本机**读数，不等于 CI 读数；CI 里这一步今天**跑没跑**见 §6 欠账。

---

## §3 ⓑ 反形两发 + 还原自证

全程 `go vet -overlay`，**未改任何跟踪文件**：被塞错的副本在仓外 `D:/tmp/wisp111r6/mut/`，
`overlay.json` 的 Replace 键用正斜杠绝对路径（`cygpath -m` 换算，路径里的空格原样保留在 JSON
字符串内）。被 mutate 的枚：`internal/ball/live_windows_test.go`，追加一行
`this is not valid Go syntax @@@@ zzz`（真实文件未被写入）。

- 发 (i) - 带门命令 + overlay（塞了语法错）⇒ **必须红**：
  `go vet -tags winlive -overlay <json> ./cmd/wisp/ ./internal/ball/` →
  `rc = 1`，报 `internal\ball\live_windows_test.go:691:1: expected declaration, found this`。
  件：`logs/negform-with-tag-red.txt`。⇒ 这步确实有牙。
- 发 (ii) - 同一 overlay、去掉 `-tags winlive`（= 今天 CI 编译这两包所用的形状）⇒
  **同样的错不红**：`go vet -overlay <json> ./cmd/wisp/ ./internal/ball/` → `rc = 0`，
  上面那行错误**完全不可见**（该枚被 `windows && winlive` 排除，vet 根本不解析它）。
  件：`logs/negform-no-tag-green.txt`。⇒ 证今天确实没人查这 12 枚。
- 还原自证（`logs/negform-restore-proof.txt`）：overlay 用完即弃、未动跟踪文件 ⇒
  对 12 枚重跑 `git hash-object` 与 §0 基线 `logs/roster-sha.txt` 逐枚同（diff 空，
  baseline 12 / now 12）；`git status --porcelain -- cmd/wisp internal/ball` = 空。
  ⇒ 名册与跑前逐枚同，且零跟踪改动。

---

## §4 门禁各件 rc

- `d22scan`（独立 module，`cd tools/d22scan && go run . -root ../../`）→ **rc = 0**，
  `clean - no D22 ban violations`。件 `logs/gate-d22scan.txt`。
  ⚠ 注：其 emoji 门 #8 射程 = design/(39)/frontend/(85)/internal/(514)/cmd/(108)，**不**含
  `.github/` 与 `.scratch/`；但新步的 ci.yml 注释本就零 emoji（本文件/界面纪律同样守）。
- `ci.yml` YAML 合法性 → **rc = 0**（`python -c "import yaml; yaml.safe_load(...)"` 解析通过，
  `test-windows` 步数 9 → 10，含我新步）。件 `logs/gate-yaml.txt`。
  无 yamllint（本机无），改用 python yaml；两者选一按 §4.7 说明。
- 未跑（按令）：`go build ./...`（§4.7 明说"不需要，别跑"）；`gofmt` 不适用（未新增 .go）；
  ⛔ 未跑任何 `go test`（尤其 `./cmd/wisp/`，35-v4 独占其窗，且 AC#11 三半都不要求 go test）。

---

## §5 ⓒ 牙的口径（买到什么 / 买不到什么）

买到：**这 12 枚带 `windows && winlive` 的文件，若被改坏到编译/vet 不过，`test-windows`
腿会红。** 仅此一面＝编译覆盖（compile coverage）。

买不到（具名、防误读）：
- ⛔ 不执行任何用例：`go vet` 不产测试二进制、不 `--- RUN`/`--- PASS`，一个用例都没跑。
  故**不许**把本步读成"票 35 `:52`/`:75(c)` 的真窗读数有了 CI 载体"——那两格按 `A690`
  走**另一支**（明写"仅本机可量"），本步与它无关、不抵它的账。
- ⛔ 不冒充 AC#7（某枚 live 用例跑没跑）、不冒充 AC#9（winsec 执行覆盖）：AC#11 盘上原文
  与本条 ⓐ/ⓑ/ⓒ 均在同一 AC#11 之下，只主张编译面。
- ⛔ 不 disturb `cmd/wisp` 的 go 测试窗：本步只 vet、不建测试二进制、不碰原生 DLL 运行，
  与另一腿的运行窗口无争（这也是 AC#11 里"执行面 vs 编译面"两个方向的纪律）。
- ⛔ 不改那 12 枚断言/加 build 条件/换包名/放宽判据（红也照实报，本例今绿）。

---

## §6 具名欠账（没取的数 / 没跑的尺 / 与转述不一致处）

1. **CI-success 读数：欠（要等推送之后）**。⛔ 零 push ⇒ 本腿没有任何一次真实 GitHub run
   可查。"这一步在 CI 里出现过 success/failure"这一格**取不到**，只能记欠；
   ⛔ 不把"yaml 里有这一步"读成"它跑过"，因此本腿**不自称 AC#11 完成**（完成与否由编排者
   在非实现者验收后判定）。§2/§3 的读数都是**本机**读数。
2. ⓐ 的出处是 windows 宿主机，非 ubuntu、非托管 runner；runner 上 `-tags winlive` 的 cgo
   编译能否 rc=0（third_party/build.ps1 之后）属"待 CI 首跑验"的范畴，同样落在 1 里。
3. 与编排者转述的差异（**以盘上原文为准**，均已按原文做，此处只点名不改动）：
   - AC#11 盘上原文（`sed -n '91,94p'`）比转述多出两条 ⚠ 交叉核对子弹（"与 AC#7 不同轴"
     "与 AC#9 的告诫相反要小心"），并把真窗载体那支具名到票 35 `:52`/`:75(c)` + `A690`。
     实质不冲突，本腿 ⓒ/§5 已按原文覆盖这两轴不互相抵账。
   - §2 前提 1 的 `if:` 尺：转述"步级 12/作业级 0"成立；但我发现严格 8 空格尺会漏数 :462
     的 `always()`（得 11），用任意空白尺才得 12。已两把尺都落件，取 `^\s+if:` 为准。
   - 其余 §2 前提（作业名册/落点/必须带 `!cancelled()`）逐条复跑一致。
4. 未造 `scripts/` 包壳（见 §1 理由）；未跑 `go build ./...`/`go test`（见 §4）。
5. 预算：本腿未触及 40/55 次告警线，全部 6 格（锚、ⓐ、yaml、ⓑ、门、票面追加+本件）均落 commit，
   无因额度跳步或自称全绿。

---

## commit 链（本腿造，均在 `8e98b8f0` 之后）

- `02505442` §0 起手锚读数件
- `779714f8` §2 ⓐ 现读数
- `351e5a5e` §1 ci.yml 新步
- `04b01c6c` §3 ⓑ 反形+还原自证
- `e6981776` §4 门禁
- （本件与票面追加随后一笔）
