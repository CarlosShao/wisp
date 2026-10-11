# 308-r1 · 90-final（落地腿终局；票 308 `AC#1`–`AC#5` 自报，判语归 `308-v1`）

锚笔 `1cd57ec1`；分支 `dev`（共享工作树）。起手 `git status --porcelain -- internal cmd`＝0，`wisp.exe`/`balldebug.exe` 在场尺＝0。

## ① commit 台账 ＋ 逐笔越界枚数（尺＝`git show --name-only --format= <sha> | grep -c .` ⇄ 写面外枚数应 0）

写面＝`internal/risk/*_test.go`(3) ＋ `probes/308/r1/**` ＋ 票 308 末尾 1 bullet。

| 短号 | 内容 | files | 越界枚数 |
|---|---|---|---|
| `510c95ca` | 00-anchor（commit-first） | 1 | 0 |
| `1ad48333` | 代码：两崩点件经 `driveShare` + 新 `pathresolver_volumeless_spelling_308_test.go`（helper＋直喂字面量用例）+ ②①注释缺陷 | 3 | 0 |
| `ee0a93e4` | 交付件 10/20/30/40 ＋ logs/（改前后名册、(a)(b) 反形证据、三把门禁 rc） | 35 | 0 |
| 本笔 | 90-final（本件，⛔ 票 308） | 1 | 0 |

越界尺逐笔跑通（`git show --name-only --format= <sha> | grep -vE '^\.scratch/wisp/probes/308/r1/|^internal/risk/pathresolver_(anchor_spelling|junction)_windows_test\.go$|^internal/risk/pathresolver_volumeless_spelling_308_test\.go$' | grep -c .`）＝各笔 **0**。

⚠ **票 308 末尾 bullet 本腿⛔ 提交（具名，见 §③ Q-308-r1-D）**：本腿已把那一枚 append-only bullet 写进工作树，但**票文件同时躺着编排者⛔ 提交的并发编辑**（起手 12:5x 读时票 40 行、AC#0＝`- [ ]`、⛔ `## 编排者裁定`；到 13:1x 工作树多出＝AC#0 翻 `- [x]` ＋ 一整节 `## 编排者裁定（A846）`——⛔ 本腿写、⛔ 本腿翻的框）。`git show HEAD:<票>`＝40 行（尺：`git show HEAD:<票> | grep -cE '^- \[x\] AC#0'`＝0、`grep -c 编排者裁定`＝0）。硬约束＝⛔ 翻任何框、⛔ 提交别人在飞中间态、本腿票增量应纯 +1/−0 ⇒ 本腿⛔ `git add` 票文件（stage 会把编排者的 tick＋A846 节一并纳入本腿 commit）。⇒ **票尾枚 bullet 落工作树、由编排者把自⛔ 的 A846 节连同本腿 bullet 一起 commit（或交 `308-v1`）**；本腿⛔ 代编排者 commit。

## ② 五格逐格答了什么

- `AC#1`（守卫落点，裁语 ⓐ–ⓒ）：共享 helper `driveShare(vol)->(share,ok)`（`len(vol)==0`⇒¬ok）＋`volumeInvariantSpellings(home,long)->(keep,dropped)`；两枚崩点件现量 0 处未守卫 `vol[:1]`。直喂无卷标字面量 `/tmp/...` 的单元用例断言三件事（①不 panic ②`unc-admin`/`unc-extended` 具名去掉＋`t.Logf` 印分母 ③长/大写/小写/扩展等 7 枚仍求值）。无整枚 `t.Skip`。详见 `10-guard-landing.md`。**闭合。**
- `AC#2`（反形敏感，三件）：(a) 拔守卫⇒单元用例回到 `panic: slice bounds out of range [:1] with length 0`，`rc=1`，逐字在 `logs/anti-a-mutation-broken.txt`。(b) 还原后 sha256＝`f2c01307…` ⇄ `logs/anti-b-restored-good.go.txt` 同 hash，复绿 `rc=0`。(c) 别名台面（`RISKAT~1\RUNNER~1`）改前⇄改后逐名作差＝只 ＋1 枚（新用例），旧名一枚不少一枚变色；两台面改后逐名相同。详见 `20-anti-form.md`。**闭合。**
- `AC#3`（分母与档位）：两崩点件 `//go:build windows`⇒test-core(ubuntu)不编译、test-windows(windows-latest)编译且求值；新件无 tag⇒两档都求值。详见 `40-tags-and-gates.md`。**部分闭合**：CI 实跑色未取（见 §③ 欠账 Q-308-r1-A/B），本腿只交静态档位定位。
- `AC#4`（门禁与越界）：`GOFLAGS= go build ./...`＝0、`go vet ./internal/risk/`＝0、`sh scripts/d22scan.sh`＝0，各 `logs/*.rc.txt` 自落 `rc=N`（非空）。格式名册两把并排：`gofmt -l internal/risk`＝1（`provenance.go` CRLF 假枚，blob 尺证）；`git ls-files '*.go'|xargs gofmt -l`＝27（23 `.scratch` 变异拷贝＋4 非-`.scratch` CRLF 假枚）；`gofmt -l` 我三枚＝0。既有债⛔ 我修⛔ 因它判红。逐笔越界＝0。**闭合。**
- `AC#5`（对账）：⛔ 翻票 72／307 任何格；"同文件不同件事"写清；⛔ 动产码、⛔ 改票面正文。详见 `40-tags-and-gates.md` §AC#5。**闭合。**

未闭合项：`AC#3` 的 CI 实跑色（本腿⛔ push，代答不了）。

## ③ 欠读数（编号＋该由谁用什么尺答；⛔ 本腿代答）

- **Q-308-r1-A**（`308-v1`／编排者，推送那批一起取）：test-windows `portable-tests.sh --scope=windows` 对 `TestClassifySpellingInvariance`／`TestPathResolverUNCAListDenied`／`TestVolumeInvariantSpellingsVolumeLess308` 三枚实跑 PASS/FAIL 色。尺＝CI 日志名册；⛔ 用"本机跑过"顶替。
- **Q-308-r1-B**（同上）：test-core(ubuntu) 对无 tag 新件 `TestVolumeInvariantSpellingsVolumeLess308` 的实跑色（helper 牙跨 OS 的正面证据）。
- **Q-308-r1-C**（编排者，票 §现量 第三条注释"never quietly stops measuring"）：那句在 `vol==""` 那形下为假——本腿⛔ 改（票面"登记、⛔ 顺手改对"，且 A843 只归 ①②）。问：要不要随守卫落点单独改写、归哪格。见 `30-comment-fixes-A843.md` §第三枚。
- **Q-308-r1-D（停手上报项）**：票 308 工作树⛔ 本腿一人编辑——本腿只写了末尾 1 枚 append-only bullet（`git diff` 里那枚 `+- 308-r1…` 行），但同一枚文件此刻还躺着编排者⛔ 提交的并发编辑（AC#0 由 `- [ ]` 翻成 `- [x]`、末尾新增一整节 `## 编排者裁定（A846）`）。三把尺坐实⛔ 本腿所为＝`git show HEAD:<票>` 40 行／`grep -cE '^- \[x\] AC#0'`＝0／`grep -c 编排者裁定`＝0（尺可复跑）。⇒ 本腿据"⛔ 翻任何框／⛔ 提交别人在飞中间态／票尾增量应纯 +1−0"⛔ `git add` 票文件；请编排者把自⛔ 的 A846 节连同本腿那枚 bullet 一起 commit，或由 `308-v1` 在核过框归属后处理。本腿⛔ 代裁那枚 `[x]` 属⛔。

## ④ 编译面命令自报（各几发，⛔ 瞒）

- `go test ./internal/risk/`：共 **7** 发（pre-default-full、pre-alias-full、post-default-full、post-alias-full、post-restored-default-full、anti-a 突变单用例、anti-b 还原单用例），原始输出全在 `logs/`。
- `GOFLAGS= go build ./...`：**1** 发（rc=0，⛔ 命中 `internal/tools/**`）。
- `go vet ./internal/risk/`：**1** 发（rc=0）。
- `sh scripts/d22scan.sh`：**1** 发（rc=0）。
- `gofmt -l`：三把射程各 1 发（我三枚／`internal/risk`／全仓 tracked）。

## ⑤ ⛔ 顺手动到的东西（具名）

- ⛔ 动产码（`pathresolver.go`／`pathresolver_windows.go`／任何非 `_test.go`）：`git diff --numstat` 对这两枚＝空。
- ⛔ 动 `internal/tools/**`（309-r1 在飞面）：⛔ stage⛔ 改；本发 build 全绿未撞其红。
- ⛔ 动 `.gitignore`／`go.mod`／`docs/reports/**`／`docs/evidence/s1/**`／`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／SLO／golden／`thresholds.go`／`allowlist.txt`／D43／C1–C32／D1–D47／三枚冻结件。
- ⛔ 放宽断言、⛔ 造 `--- SKIP` 换绿、⛔ mock 代真、⛔ 覆盖搬家式消失（去掉 2 枚拼法＝分母写清 7/9）。
- ⛔ `push`／⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；commit 全带显式 pathspec；⛔ `git add -A`／`.`。
- 临时件只建不删，全落 `probes/308/r1/logs/`，⛔ 仓根，⛔ `.out` 名。

## ⑥ 改前⇄改后逐名作差结论

- 顶层名册尺（`^--- (PASS|FAIL|SKIP):`）：改前 default／alias 各 **132** 行（131 PASS ＋ 1 既有 SKIP `TestSyncRegistryProbeLive`）；改后 **133** 行。**新增红枚数＝0**。
- 唯一名册差＝＋`PASS: TestVolumeInvariantSpellingsVolumeLess308`（本腿新用例），两台面同。⛔ 一枚旧名消失或变色。
- 子层级 `^    --- PASS`＝68⇄68（⛔ 变）；六份全名册 `panic|fatal|test executable failed` 命中均＝0。
- **别名台面（8.3 那把）颜色未因守卫而变**：`diff logs/pre-alias-names.txt logs/post-alias-names.txt` 只 ＋1 枚；`diff logs/post-default-names.txt logs/post-alias-names.txt` 空（两台面改后同集合）。⇒ 守卫⛔ 关掉用例。
