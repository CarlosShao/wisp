# 300-a1r — `90` 卫生与门禁（三把＋越界自查＋名册）

锚＝`b2933dc9`。本程**零产码**：⛔ 改任何跟踪文件（唯一被追加的既有件＝票面 `Progress log` 那一行，见 §4）。
每一把都带**命令逐字**＋**`rc` 取法**＋**读数**；⛔ 从 `tee`／`head` 的退码冒充被测命令的退码。

## 1. 三把门禁（⛔ 顺序、⛔ 口径）

### 1.1 ① `sh scripts/d22scan.sh` ⇒ **`rc-d22scan=0`**

命令逐字：`sh scripts/d22scan.sh > /tmp/wisp300-a1r/logs/hygiene-d22scan.txt 2>&1`（`rc` 取法＝`echo $?` 紧跟该命令，⛔ 管道）

末三行逐字（全文＝`logs/hygiene-d22scan.txt`，23084 字节，尺＝`wc -c`）：

```
d22scan: scope ban #8 internal/         examined 524 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined 119 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=229, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=524, ban #8 cmd/=119; ban #8 emoji coverage: design/ 39 text files; frontend/ 85 text files; internal/ 524 Go files, comments and _test.go included; cmd/ 119 Go files, comments and _test.go included
```

★**报告口径＝当前树**（工作树，⛔ blob）：那 **约 805 枚**既有脏项（`design/**`／`frontend/**` 的删除项）**⛔ 我造成**，
本把尺对它们的态度是"照扫"——它回的是 `clean`，⇒ 我这三枚 `.md`／八枚 `.txt` 落在 `.scratch/**`，
**根本不在这把尺的射程里**（射程＝上面那行点名的 `internal/` `cmd/` `frontend/` `design/` `internal/tools/`）。
⇒ 本格对"我这枚 `.md` 里那些 `→`／`⇒`／`⛔`"⛔ 免疫凭 d22scan，⛔ 免疫凭"它没扫"这件事当作合规——
**它没扫＝它不在分母里**，与本仓"注释豁免、字符串不豁免"那套仪器口径是两件事（见 AGENTS.md §1.2 的"规格 vs 仪器"那段，⛔ 我据它做任何判断）。

### 1.2 ② `GOFLAGS= go build ./...` ⇒ **`rc-build=0`**

命令逐字：`GOFLAGS= go build ./... > /tmp/wisp300-a1r/logs/hygiene-build.txt 2>&1`（`rc` 取法＝`echo $?`）
⇒ stdout＋stderr **0 字节**（尺＝`wc -c`，原始件 `logs/hygiene-build.txt`；空**是读数本身**＝`go build` 成功时不打印，
我在该件顶上追加了 3 行 provenance 头，⛔ 让一枚 0 字节的件被读成"那格没交"）。
⛔ 意外：产码零改动，这把本来就该绿；它绿**⛔ 证不了**任何行为，只证我没把谁的树写坏。

### 1.3 ③ 越界自查 ⇒ **`git status --porcelain -- internal cmd` ＝ 0 行**

命令逐字：`git status --porcelain -- internal cmd | tee /tmp/wisp300-a1r/logs/hygiene-scope.txt`
- 尺＝`wc -l < logs/hygiene-scope.txt` ⇒ **0**（⚠ 这条管道里 `$?` 是 `tee` 的，⛔ 拿它当 `git status` 的退码；
  本把尺要的是**行数**，0 行＝`internal/**` 与 `cmd/**` 两棵子树我一枚没动）
- ★`cmd/wisp` 的仪器面（票 255 名册 `cmd/wisp/config_readers_255.go` 用行号引产码行）：本程**没碰 `cmd/**`**，
  也没碰 `internal/**` ⇒ 那把尺⛔ 被我打红（先例 `A799`/`A801` 的形我这一程**结构性不可能犯**）。
- ⚠ **口径**：整棵工作树**本来就脏约 805 枚**（⛔ 我造成、⛔ 我复位）⇒ 本件⛔ 任何一处写"工作树干净"；
  越界回报只认上面那把"逐路径 0 行"的尺。

## 2. 我自己跑过什么（自报，⛔ 藏着当没跑）

| 命令 | 射程 | 结果 |
|---|---|---|
| `go version` / `go env GOOS GOARCH GOFLAGS` / `go env GOMODCACHE` | 环境 | 见 `00-anchor.md`；`GOMODCACHE=D:\work\base\gopath\pkg\mod` |
| `sh scripts/d22scan.sh` | 全仓当前树 | rc=0（§1.1） |
| `GOFLAGS= go build ./...` | 全仓 | rc=0（§1.2） |
| `go test ./internal/audio/ -count=1 -run TestParseWaveFormatSubFormatOffset300 -v` | **仅仓外导出树**，跑 **2 发**（改前／突变后） | rc=1／rc=0，逐字见 `20-ac1-rig.md` §2.1／§2.2 |
| `go vet ./internal/audio/` | 仅导出树 | rc=0，空输出 |
| `go test -c -gcflags='-m' -o ../probe300-test.exe ./internal/audio/` | 仅导出树（堆证据） | `logs/ac1-escape.txt`；产物 `probe300-test.exe` 落在 `/tmp/wisp300-a1r/`（⛔ 仓内），留着⛔ 删 |
| `git archive b2933dc9 \| tar -x -C /tmp/wisp300-a1r/tree` | 建导出树 | 存在性读数见 `20-ac1-rig.md` §0 |
| `find`/`grep`/`sed`/`ls`/`wc`/`cygpath` 只读尺 | Windows SDK 头文件＋仓内若干行 | 全部逐字落 `logs/*.txt` |
| ⛔ **没跑** | 整包 `cmd/wisp` 测试（约 8 分钟且占测试面，归编排者）、真麦克风／真设备那一发（`GetMixFormat`，`AC#3`）、`gofumpt`／格式名册（`AC#4` 落点）、CI 三跳（`AC#5`） | 本格⛔ 交这些读数 |

## 3. `git show --stat` 逐笔名册（期望只含 `.scratch/wisp/probes/300/**`）

尺＝`git show --stat --oneline <号>`，口径＝**blob**（⛔ 工作树）；三笔已存在的号：

- `70a6b9d9` → `.scratch/wisp/probes/300/a1r/00-anchor.md`（1 file changed, 68 insertions）
- `cc878604` → `.scratch/wisp/probes/300/a1r/10-ac0-authority.md`（193＋）
  ＋ `.scratch/wisp/probes/300/a1r/logs/ac0-mmreg-extract.txt`（78＋）（2 files changed, 271 insertions）
- `1c007bf6` → 8 files changed, 609 insertions：`20-ac1-rig.md` ＋ `logs/` 下 **7** 枚 `.txt`（尺＝`git show --name-only` 的行数，⛔ 数成 8 枚 `.txt`——8 是**文件总数**，其中 `.md` 一枚），
  即 `ac1-escape` `ac1-mutated-24` `ac1-pristine` `ac1-vet` `hygiene-build` `hygiene-d22scan` `hygiene-scope`
  ⚠ 上面 §2 那张"逐字"读数取自 `tail -8`，**它把 `20-ac1-rig.md` 那一行截掉了**（⛔ 拿被截的读数当名册全集）
  ⇒ 完整名册在**最后一笔 commit 之前**又用 `git show --name-only --format=%h 1c007bf6` 现跑了一遍，逐字贴在本件末节。
- **最后一笔**（本件＋票面那一行的载体）＝⛔ 我写本件时还不存在，号在交回正文里给编排者；
  它的名册＝`90-hygiene.md` ＋ `logs/hygiene-build.txt`／`logs/hygiene-scope.txt`（provenance 头追加）＋票面 `Progress log` 一行。

⛔ push（机主从未授权）；每笔都带**显式 pathspec**；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；
⛔ `--no-verify`——本仓 `.git/hooks` 只有 `post-checkout`／`post-commit`（先例 `A800`：钩子缺位⛔ 是"可以绕"的许可证）。

### 3.1 CRLF 口径（本仓实测的那把假枚数尺）

`git add` 时 Git 现报（逐字，节选）：

```
warning: in the working copy of '.scratch/wisp/probes/300/a1r/logs/ac1-escape.txt', LF will be replaced by CRLF the next time Git touches it
```

⇒ 这些 `.txt`／`.md` **blob 里是 LF**，被 checkout／重新 touch 时工作树会变 **CRLF**。
影响两把尺：`wc -l`（枚数）**两把一致**，`wc -c`（字节数）**会差**（每行差 1 字节）
⇒ 本件凡报枚数都⛔ 报字节数当枚数；要复跑我的读数请在**工作树**里跑，或先 `git show <号>:<路径>` 取 blob 再数。

## 4. 票面那一行（唯一被追加的既有件）

- 落点＝`.scratch/wisp/issues/300-parsewaveformat-subformat-offset-reads-two-bytes-past-the-guid.md` 的
  `## Progress log` 节**末行追加一行**，钟点由 `date '+%Y-%m-%d %H:%M:%S'` 的 stdout 插值（⛔ 手打）。
- ⛔ 勾任何 `- [ ]` 框（`AC#0`／`AC#1` 都**⛔ 由我翻**，翻框＝编排者）、⛔ 改票面其它一字。
- ⇒ 票面那两格**盘上仍显示未开工状态**，这是**故意的**，⛔ 读成"腿没做"。

## 5. 我不同意票面／派单的地方（★具名，⛔ 替谁圆场）

1. **票面 `:13` 那句关于 union 的话是错的**（原文逐字：`标准结构里没有独立的 validBits 字段，nValidSamples 是与 dwChannelMask 同一个 union 的成员`）。
   前半对，**后半与两把权威都矛盾**：`Samples`（三枚 WORD 的 union，2 字节）与 `dwChannelMask`（DWORD）是**两个独立成员**，
   没有任何一枚 union 同时装着它们；标准结构里⛔ 存在叫 `wValidSamples`/`nValidSamples` 的字段
   （尺＝`grep -n "ValidSamples"` 射程＝`shared/mmreg.h` `shared/ksmedia.h` `um/mmdeviceapi.h` `um/audioclient.h` ⇒ **0 行**）。
   ⛔ 影响 `AC#0` 结论（偏移 24 靠字段序＋`cbSize=22`＋总长 40 三枚闭合，⛔ 依赖这句），但它在票面**现量节**里 ⇒ 归编排者处置（我⛔ 改票面那一节）。
2. **派单示例路径 `Include/*/um/mmreg.h` 在盘上不成立**：那枚头文件在 `shared/`（`ls .../um/mmreg.h` ⇒ rc=2、0 行）。
   ⛔ 影响结论；写下来只为下一枚腿⛔ 照示例再撞一次。
3. **票面 `AC#1` 的"预期今天两枚都红（都读出 0）"低估了一形**：我按权威布局造了 4 形，
   其中 **S4（`Data1` 故意挪到 26 的坏面）今天读出的是 `3`、⛔ 是 0**（逐字见 `20-ac1-rig.md` §2.1）。
   这不是分歧而是加强：**"恒读 0"只对良构面成立**；把"都读出 0"写成通则会把 S4 这一形藏起来，
   而 S4 正是"这枚尺两种形状都放行⛔"的那枚正控。⇒ 建议 `AC#2` 入库那枚用例**保留 S4**（连同 `want tag == 0` 的断言），
   ⛔ 只留 S1/S2（那是一把单形尺）。
4. **派单里 `cbSize=22`／总长 40 那一串是编排者推的**（票面 `:13` 自己写着"⛔ 是读数"）——
   本程把它从"推的"升成"读到的"：`mmreg.h` 第 **2540**／**2550** 行两处逐字 `/* Format.cbSize = 22 */`，
   Learn 那句逐字 `The **cbSize** member must be at least 22.` ⇒ **⛔ 再当假设引用**。
   ⚠ 一处仍开着：Learn 写的是 "**at least** 22"，头文件写的是"＝22"，两把在"能不能大于 22"上⛔ 同形
   （大于 22 的附加数据是另一枚形状，与本票的偏移问题无关，我⛔ 判它）。

## 6. 本件的完整名册复核（最后一笔 commit 之前现跑）

`git show --name-only --format=%h 1c007bf6` 逐字输出：

```
1c007bf6

.scratch/wisp/probes/300/a1r/20-ac1-rig.md
.scratch/wisp/probes/300/a1r/logs/ac1-escape.txt
.scratch/wisp/probes/300/a1r/logs/ac1-mutated-24.txt
.scratch/wisp/probes/300/a1r/logs/ac1-pristine.txt
.scratch/wisp/probes/300/a1r/logs/ac1-vet.txt
.scratch/wisp/probes/300/a1r/logs/hygiene-build.txt
.scratch/wisp/probes/300/a1r/logs/hygiene-d22scan.txt
.scratch/wisp/probes/300/a1r/logs/hygiene-scope.txt
```

- 同一发再补一把越界尺（⛔ 靠肉眼数上面那 8 行）：
  `git show --name-only --format= 1c007bf6 | grep -vc '^\.scratch/wisp/probes/300/a1r/'` ⇒ **0**
  （尺＝**不匹配**我射程目录的路径行数 ⇒ 0＝该笔一枚越界件都没带）
- ⚠ 两枚 `hygiene-*.txt` 在**这次名册之后**被 `>>` 追加了 provenance 头（§1.2／§1.3 那两句"空是读数本身"），
  ⇒ 名册里它们显示 `| 0`；追加的那几行随**最后一笔** commit 入库，⛔ 改写历史、⛔ `--amend`。
