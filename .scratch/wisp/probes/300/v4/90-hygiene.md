# 300-v4 · 门禁件、越界自查、边界自陈、具名欠账

锚＝`34e1962b`（`git rev-parse HEAD` 起手现量，件 `00-anchor.md`）。
台面＝**仓外导出树** `/tmp/wisp300v4/tree`（`git archive 34e1962bb166b4494670b77296263c79bd8671a3 | tar -x`，件 `logs/00-archive-rc.txt`）。
⚠ 同台另有一枚腿在飞：`303-v1` 于 **20:19** 落了 `8fc42230`（在我两笔之间）。本腿所有 Go 读数取自 **34e1962b 的导出树** ⇒ 世代⛔ 被它污染；`git show --name-only 8fc42230` 名册⛔ 含 `internal/audio/**`（我这发核过），故工作树侧也无碰撞。

---

## 1. 本腿现跑的门禁尺，逐把一行 `rc=N`（⛔ 0 字节件）

| 尺 | rc | 读数 | 件 |
|---|---|---|---|
| `go vet ./internal/audio/` | **0** | stdout 0 字节而 rc=0（两件事⛔ 混写，照票面 `AC#6` 那节顶回②的口径） | `logs/30-vet.txt` |
| `gofmt -l internal/audio`（导出树＝blob 内容） | **0** | `count=0`；工作树侧 `git status --porcelain internal/audio`=**0 行**⇒ 工作树≡blob ⇒ 两把尺在本射程**collapse 成同一个 0**；全仓 5 枚既有残留⛔ 我射程、⛔ 动 | `logs/31-gofmt-scope.txt` |
| `sh scripts/d22scan.sh`（导出树） | **0** | 含 ban#8 实扫 | `logs/32-d22scan-tail.txt` |
| `GOFLAGS= go build ./...`（导出树） | **0** | 独立第三人复跑那一枚派单漏列的尺 | `logs/33-build-all.txt` |
| `go test ./internal/audio/ -count=1 -coverprofile=…` | **0** | `parseWav 70.0%`、`Drain 0.0%`、`convertPacket 0.0%` | `logs/34-cov-run.txt`、`39-…` |
| 定向用例 pristine（尺＝`^--- PASS` 顶层 1 枚／`^    --- PASS` 子测试 6 枚） | **0** | 三次 pristine（`10`/`13`/`25`）逐字同形 | `logs/10-…`、`logs/13-…`、`logs/25-…` |
| 整包名册 ×4 发（`-v`，改后 2＋改前 2） | **0** | 四发 `rc=0`、`^--- FAIL`＋`^    --- FAIL`=**0**、`--- SKIP`=**1**（同名 `TestLiveWasapiSmoke`）、`^--- PASS` 顶层 **42（含新用例）⇔ 41（拆掉）** | `logs/70-*.txt`、`logs/44-*.txt`、`logs/45-*.txt`（0 字节＝名册空＝新增红 0，读数见 `logs/75-…`） |
| `find … -name mmreg.h`（仓外只读） | **0** | 全机恰好 **1** 枚命中 | `logs/80-authority-mmreg-v4.txt` |
| `sed -n '2110p/2376p/2418p'` 逐字 | **0** | 三行与件内引文逐字同形（含 `#define` 后两枚空格与 `0x0003` 写法） | 同上 |
| `grep -n KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` × 两枚头文件 | **0／0** | `mmreg.h:2480/2481/2483/2484` 与 `ksmedia.h:851/852/854/855` **都有** | 同上 |
| `tasklist` 四把（起手闸门） | **0** | wisp=0／balldebug=0／mockllm=0／webview2=24（⛔ 残留，先例 `A821`） | `00-anchor.md` |

突变台（全部仓外，`cp` 备份→`sed`→跑→写回→`cmp`）：**六枚 `restore_cmp=IDENTICAL`**（`logs/60-restore-cmp-all.txt`）＋ injector 那两枚各自 `cp` 还原 `IDENTICAL`＋三枚文件与 HEAD blob 逐字节 `IDENTICAL`（`prod_vs_blob`/`test_vs_blob`/`inj_vs_blob`）。
导出树还原终尺：`cmp <导出树> <(git show HEAD:<同路径>)` 三枚全 `IDENTICAL` ⇒ 本腿⛔ 在共享工作树里留下任何字节改动，也⛔ 做过 checkout/reset/stash/clean。

---

## 2. 越界自查（逐笔，尺＝`git show --name-only --format= <sha>`）

| 笔 | 号 | 文件数 | 越出 `.scratch/wisp/probes/300/v4/**` | 三枚冻结件／`thresholds.go`／`allowlist.txt` | `scripts`／`.github`／`frontend`／`design`／`tools`／`internal/audio`／`docs` |
|---|---|---|---|---|---|
| 1 起手锚 | `2f2597ed` | 2 | **0** | **0** | **0** |
| 2 判语＋件 | `438b0145` | 39 | **0** | **0** | **0** |

- `commit` 一律显式 pathspec（`git commit -F <件> -- .scratch/wisp/probes/300/v4`），⛔ `git add -A`／`git add .`；commit 消息一律 `Write` 成件再 `-F`（本轮⛔ 一次反引号内联执行事故）。
- 件总数尺：`v4/` 下 **41** 枚（`*.md` 3 ＋ `*.txt` 38），⛔ 一枚命名 `.out`（根 `.gitignore:8` 全仓 `*.out` ⇒ 会被静默跳过而 commit 仍回显成功），`find -name '*.out'` = **0**。
- 0 字节件 **4** 枚（`44`/`45` comm 双向空、`71`/`72` 红名册空）＝**名册为空的合法读数**，逐枚解释写死在 `logs/75-roster-and-comm-reading.txt`，防"0 字节＝那格没交"那把尺误读。

---

## 3. 三数口径自守

- ① 所有枚数**带尺名**：顶层尺＝`^--- PASS`／`^--- FAIL`，子测试尺＝`^    --- PASS`／`^    --- FAIL`，包级行＝`^ok|^FAIL`。本格两处并存：定向件＝**顶层 1／子测试 6**；整包＝**顶层 PASS 42（含新用例）⇔ 41（拆掉那枚）**，与 `300-a2` 的 `F5`（源码 `^func Test` 尺＝42）⛔ 冲突、⛔ 混用。
- ② 词面尺结论**逐字抄命中行**：`logs/20-expected-side-hits.txt`（9 行）、`logs/80-authority-mmreg-v4.txt`（三行引文＋两枚 GUID 四处命中）、`logs/50/51-*-after.txt`（injector 那两行原文）。
- ③ 格式名册写明**射程＋blob/工作树**：`internal/audio` ＋ **导出树（blob 内容）** count=0 ＋ 工作树≡blob（脏行数 0）；全仓 5 枚既有残留⛔ 本射程。

---

## 4. 具名欠账（取不到的读数＋为什么）

1. **`cmd/wisp` 整包的红名交集（`AC#4` 那半枚）**——⛔ 取。为什么＝本腿派单逐字 `⛔ 跑 cmd/wisp 整包（那是 303-v1 的独占面）`，且那一半需要 sherpa `PATH` 铺设＋真窗台面（`TestAC4FocusReturnToPriorWindowGap33r5` 本机 13 发 5 绿／8 红）。归属裁语见 `20-residuals-and-ac4.md` ⓕ2＝**记编排者／与 `AC#5` 那发 CI 同批销**，⛔ 记落地腿。
2. **CI 那一发的颜色（`AC#5`）**——⛔ 取、也⛔ 该本腿取：只有推送能给（票面逐字欠着，编排者已写"⛔ 为它单推"）。
3. **`scripts/portable-tests.sh` 的 `win_pin`/`core_pin` 求值顺序全貌**——我只读了 `:18`/`:587`/`:593`/`:594`/`:704` 那五行（⛔ 跑 `scripts/**` 任何一支，⛔ 动它）。够坐实"未登记 SKIP＝红"这一条前提，⛔ 够裁 CI 分母全貌——那一裁本腿⛔ 需要。
4. **`R1b` 的两枚 injector 突变我只跑了一发一把尺**（`rc=0`／名册 0 行，各一发）。裁它"零牙"用的是**两把独立尺交叉**（覆盖 0 命中 ＋ 突变整包照绿），⛔ 成对两发；要写成票凭据请落地腿按票面 `AC#1` 那形补成对。
5. **`Drain` 那一行在真机上的实际执行**——⛔ 取（⛔ 真设备归腿，票面 `:31` 逐字）。本腿只取到覆盖尺的 `0.0%`（**测试面**零执行），⛔ 等价于"生产从不执行它"。

---

## 5. 三句自陈（本腿逐字）

- **⛔ 翻框**：本腿⛔ 翻任何 `- [ ]`、⛔ 改票面原句、⛔ 改 `-done`；判语只写在本目录的件里，翻勾由编排者做（票面 `AC#6` 原文那句"⛔ 由实现腿自勾"与本腿的"⛔ 由编排者翻"⛔ 冲突：前者禁实现者自勾，后者禁任何人非判而翻）。
- **⛔ 产码**：本腿⛔ 改任何产码或测试文件（`internal/audio/**` 工作树脏行数 **0**；三枚文件与 HEAD blob 逐字节 `IDENTICAL`；所有突变只在仓外导出树、跑完即还原并 `cmp` 证明）；⛔ 放宽任何断言、⛔ 造 `--- SKIP`、⛔ 动 SLO 阈值／golden／`thresholds.go`／`allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件／`scripts/**`／`.github/**`／`frontend/**`／`design/**`／`tools/**`；⛔ 在仓内建 worktree。
- **⛔ push**：`git status -sb` = `## dev...cnb/dev [ahead 977]` ⇒ 本腿**零推送**，推送归编排者（`A813` fast-forward）。已推送历史⛔ 改写：要更正就追加新 commit（本腿三笔皆新笔，⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`，⛔ 动 git config，⛔ 杀任何进程）。
