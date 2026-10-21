# 票 174 终裁腿 174-v1 — 它改的那句安全叙事对不对（两臂复量＋四枚承重＋写面纪律）

- 角色：**非实现者裁腿**（裁决者 非 实现者）。⛔ 未翻任何 AC 框、⛔ 未改票名/`Status:`/`-done`、⛔ 未碰台账 `docs/reports/**`、⛔ 未落任何 `.go` 字节、⛔ 未 push。
- 被裁对象：实现腿 `174-r1` 两笔 —— `a8eb0d51`（`internal/tools/task.go` 15 行 + 新判据件 262 行）／`56c07eeb`（交件 148 行）。起手锚 `b482968f`。本腿起手 `HEAD=56c07eeb`，收尾 `HEAD=d4f8962b`（`56c07eeb..d4f8962b` 只动 `docs/**` 两枚，`internal/ cmd/ scripts/ tools/` 零字节 ⇒ 本腿五门读数归属不变，已现核）。
- 尺册：`D:\tmp\wisp\20261021\174v1\sh\{10-gates,20-criteria-and-legs,30-mutations,40-premise-and-attribution,50-codepoints-and-agent,60-mut4-mut5,70-final-readings,80-anchors}.sh`；读数 `logs/`；**全部读数由本腿另起尺跑出，⛔ 未采信被告任何数**。变异与探针跑在 `git archive HEAD` 导出的私有副本 `174v1/copy/`（共享树零字节，`git status --porcelain -- internal/ cmd/` 起手与终态皆 **0 行**）。

---

## 1. 七项必判 —— 逐项判语与读数

### (1) 票面前提新鲜度 —— **成立（36/36 在 HEAD 上；行号大面积漂移但断言全部复量通过）**

尺：`40-premise-and-attribution.sh` 节 A。票面抽出的**36 枚**仓库相对路径 token，逐枚 `git cat-file -e HEAD:<path>`：**AT_HEAD=36，MISSING=0**。可判定断言对 **HEAD blob** 复量（`80-anchors.sh`）：

| 票面断言 | 票面行号 | 现量 | 结论 |
|---|---|---|---|
| 授权根只从 `cfg.FS.AllowedDirs` 拷 | `:327-331` | `run.go:510/511/514`；`grep -n allowed run.go` 全文件**只这 3 处**，零 `dataDir`→`allowed` 写入 | 断言真，行号漂 183 |
| spill 落盘＝`dataDir/artifacts` | `:609` | `run.go:1186`（另 `:1090` `ArtifactsDir:`） | 断言真，行号漂 577 |
| `paths.go` 空名单句 | `:50-52` | `paths.go:63` 逐字 `An EMPTY allowlist authorizes nothing: every fs call` | 断言真，行号漂 11 |
| `SPEC-03:35` 的 `allowed_dirs[]=[]` | `:35` | `:35` 逐字命中（🔒 同上栏） | 原地未动 |
| `task.go` 指针支只看 `ArtifactPath != ""` 零 stat | `:227` | `task.go:547`，条件里确无 stat | 断言真，行号漂 320 |

⇒ **前提无过期，不停手事由不成立。**

### (2) `AC#2b` 的"两条路" —— **兑现了（逐字命中），且两腿**运行时**逐字节同形**

- 票面（＋编排者 09-27 22:3x 第 3.2 条）要求回执照说出**两条**回来的路。改后文案（`task.go:858-860` 现读）逐字：
  `；注意：这条路径现在读不到，它不在你被授权的目录范围内：fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡`
  ⇒ **两条路各有一枚具名标记**（`[fs] allowed_dirs` / `那一张卡`），"唯一办法"字面 `才读得回来` 归零（`grep -c`＝**0**）。**兑现。**
- **同形**：不用眼看。私有副本里我自己驱动两腿（同一 C26、同一空名单、同一根外真实文件），切出 notice 区**逐字节比**：
  `SAME-BYTE-FOR-BYTE=true taskLen=246 spillLen=246`（`logs/` 探针 `TestProbe174v1TwoLegsSameWording`）。两侧文本逐字相同，含前缀 `；注意：`。
- **一枚必须记在它头上的失真（不影响本格判语）**：`task.go:852` 那句注释把判据件名写成 `task_output_pointer_notice_174rr1_test.go`——**该文件在 HEAD 与盘上都不存在**（`git ls-files | grep -c 174rr1`＝**0**；实际件名 `..._two_roads_174r1_test.go`）。这是 ban #9 形状的**幻影引用**，而 d22scan rc=0 只是因为它**没写目录前缀**、进不了 `repoPathRe`（`grep -c 'internal/tools/...174rr1'`＝**0**）——**不是仪器判它清白，是仪器扫不到**。

### (3) ★它那句安全断言的机制真实性 —— **成立，但写过头半句（详见 第 2 节的明确结论）**

三处逐字，另起尺：

1. **`fs.go` 的 `open` 不查授权根**（`fs.go:90-105` 现读逐字）：函数体＝`if d.Paths == nil` → `c, err := d.Paths.Canonicalize(raw)` → `strings.Contains(c,"??")` → `return c, nil`。**全文零 `InAllowlist`**。紧随其后的 `fsRead.Execute`（`:142` `t.d.open` / `:146` `os.Open(canon)`，`:167-171` 返回文本）**执行腿也不再问根**。
2. **allowlist 只作为判级输入进系统**：`run.go:510-514` 是唯一来源（`cfg.FS.AllowedDirs` → `NewPathCanonicalizer`）。生产码里 `InAllowlist` 的真调用点 **5 处**（其余命中是接口声明 `risk/assessor.go:143`／`agent/spill.go:70` 与定义 `tools/paths.go:201`）：`risk/rules_gateway.go:45`（判 R2）、`tools/bridge.go:1273`（审计行的 in_allowlist_scope）、`tools/refusal.go:93`（拒类分流）、`tools/task.go:847`（本票指针支）、`agent/spill.go:231`（另一条腿）——**五处全在"判级/说法"侧，没有一处在 `fs.go` 的执行腿上**。警示：本腿起手时编排者 09-27 那行记的是"三处／`bridge.go:928`／`task.go:315`"，**行号与枚数都已过期**（`bridge.go:928` 现为 `:1273`，且漏计 `refusal.go:93` 与 `spill.go:231`）；本腿按现量记 5。R2 那一行逐字（`risk/assessor.go:16`＋`:79`）：
   `//	R2  target path outside allowlist L2 (C26-canonicalized; resolver wired by 18)`
   `R2 RuleID = "R2" // path outside authorized allowlist`；`rules_gateway.go:45-50` 的贡献逐字 `level: L2, reason: "R2: 目标路径在授权目录之外: %s"`。⇒ **判级，不是硬边界。**
3. **★两臂对拼（我自己跑，正控＋负控齐）**——同一桥、同一路径、同一 C26，只换门的答复：

| 臂 | 门 | `isError` | level | bytes | window/approval | 判定 |
|---|---|---|---|---|---|---|
| **负控** A | `NoGate{}` | `true` | `L2` | 36（拒句） | 0 / **1** | 拒 |
| **正控** B | 答"允许" | **`false`** | `L2` | **20000（＝body 全长）** | 0 / **1** | **整读回** |
| 对照 C | 根加进 allowlist | `false` | **`L0`** | 20000 | 0 / **0** | 整读回、零张卡 |
| 反控 D | 答"超时" | `true` | `L2` | 39 | 0 / 1 | 拒（L2 超时＝拒，与 L1 极性相反） |

⇒ **"批准一张 L2 卡 ⇒ 根外路径真被读走"成立**（B 臂 20000 字节逐字节等于 body，`same-as-body=true`）。C 臂证明"能读回"确有对象；A 臂证明那句"没人批就是直接拒"为真。
- **但它写过头了一处（见 第 2 节）**：这句在**默认出厂形状**上不是"会要一张卡"，而可能是"连卡都不会出现、批也没用"。

### (4) 四枚新判据真的承重吗 —— **三枚承重、一枚半承重（有一发我自己的变异它抓不住）**

全名逐枚（`grep -oE 'func Test[A-Za-z0-9_]*174r1'`，4/4；`^func Test`＝4；`t.Skip`＝0；`//go:build`＝0；四名在 `scripts/portable-tests.sh` 命中各 **0**）；我这边的 `-v` 日志里四枚 **RUN=1 PASS=1 FAIL=0**。

| 我自造的变异 | 落地证明 | build | 读数（我自己跑） | 还原 |
|---|---|---|---|---|
| **MUT-1** 摘新句回旧句 | `grep -n '才读得回来' task.go`→`:859` | rc=0 | 恰 **1 枚**红 `…NamesEveryRoadBack174r1`（`:102`＋`:108` 两条断言） | md5 `432a52be…`＝HEAD blob，复跑 `ok` |
| **MUT-2** 让第二路变假（`bridge.go` `case risk.L2:` 后 `a = AnswerReject`） | `bridge.go:467` | rc=0 | 恰 **2 枚**红：`…IsRealOnThisBridge`／`…BookedAsACard` | md5 `0f754371…`＝HEAD，复跑 `ok` |
| **MUT-3** 把 allowlist 做成**执行硬边界**（`fs.go` open 加 `InAllowlist`） | `fs.go:105` | rc=0 | 恰 **2 枚**红（同上两枚） | md5 `cc22d191…`＝HEAD，复跑 `ok` |
| **MUT-4** 把 `canon` 拼进 notice 尾 | `task.go:860 '"+canon'` | rc=0 | 恰 **1 枚**红 `…LeaksNoPathOrRoots174r1` | md5 回 `432a52be…`，复跑 `ok` |
| **★MUT-5** 把**第二路整句删掉**（只剩"加进 allowed_dirs"，`L2 卡` 仍在） | `task.go:860`；`grep -c '要么批下那一张卡'`＝**0** | rc=0 | **rc=0，四枚全绿——一枚都不红** | md5 回 `432a52be…`，复跑 `ok` |

⇒ **判语**：#2/#3 承重（MUT-2、MUT-3 各咬住它们）；#4 承重（MUT-4）；**#1 只承"M2 卡这三个字在不在"，不承"第二条回来的路在不在"**——MUT-5 把"要么批下那一张卡"整句删掉，`L2 卡` 三字仍由前半句 `会要一张 L2 卡` 提供，于是四枚全绿。这正是它自称要挡的那件事（"句子准不准由事实兜底"）在**文本半**上的漏口：**#1 的标记集与它替换掉的旧句共享 `L2` 语义，所以"删掉第二路"这一发能穿过整个文件**。它自己的 M1（摘回整句旧句）抓得到，是因为旧句同时缺 `L2 卡`——**那是把两件事一起改才红**。

### (5) 写面与纪律 —— **合纪律；`.go` 2 枚、越界 0 枚；没有改别人件**

`git show --name-only` 逐笔：`a8eb0d51` = `internal/tools/task.go` ＋ `internal/tools/task_output_pointer_two_roads_174r1_test.go`（`.go` **2 枚**，越界射程 `internal/tools/**`＋`internal/agent/spill.go` 之外 **0 枚**）；`56c07eeb` = `docs/evidence/s1/174-r1-impl.md`（`.go` 0，**1 枚在 `.go` 射程外＝交件本身，派单允许**）。**正控**（证明尺会红）：同尺打 `5b14afe5` 读回 `IN RANGE=1 / OUT RANGE=1`（它改了 `internal/agent/loop.go`）⇒ 尺不是恒绿。
- **"顺手改别人件"专项**：`git diff b482968f..HEAD -- internal/agent/` ＝ **0 行**（`wc -l`＝0，`--stat` 0 行）⇒ 它自报的"两枚 `internal/agent` gofumpt 漂移＝存量"**在这条尺上成立**（那两枚与锚点比零差异）。`git status --porcelain -- internal/ cmd/` 起手/终态皆 **0 行**。
- 口径核正：`b482968f..HEAD` 里还有 `scripts/check-path-length-budget.sh`(+1/-1)、`docs/reports/**`、`.scratch/**` 的改动，**均不属被告**（`git log` 归到 `bd305c52`／`d4f8962b`／`5443fed6` 等编排者笔）——本腿未把区间当归属（⛔ 区间归因）。

### (6) 门禁复跑（我自己现跑；跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`）

| 门 | rc | 末行逐字 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | （无输出，0 行） |
| `go vet ./internal/tools/ ./internal/agent/` | **0** | （无输出，0 行） |
| `go test -count=1 -v ./internal/tools/ ./internal/agent/` | **0** | `ok  	github.com/CarlosShao/wisp/internal/tools	26.402s` ／ `ok  	github.com/CarlosShao/wisp/internal/agent	6.141s` |
| `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=244, bans #1-5 cmd/=42, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=608, ban #8 cmd/=143; …` |
| `sh scripts/check-path-length-budget.sh` | **0** | `check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` |

**四数分开（⛔ 相减）**：`^=== RUN`＝**443**；`^--- PASS`＝**336**；`^--- FAIL`＝**0**；`^--- SKIP`＝**0**。（另报：缩进层 `^    --- PASS`＝107，故 RUN 与 PASS 非同一分母。）
**另跑（不在五门内，供 AC#2b 参考）**：`go test -count=1 -run TestTaskOutputPointerNotice174r8 ./cmd/wisp/`＝**`ok 11.027s` rc=0**（需 `PATH` 带 `third_party\sherpa-onnx`；`bash` 形 PATH 会 `0xc0000135` 缺 DLL，那是本机环境事实）。

### (7) 红名册逐名 —— **新增红 0、既有红 0、skip 0（逐名清单为空，故无归属可判）**

我自己的 `-v` 全日志上 `^(---|    ---) FAIL` 命中 **0 行**、`SKIP` 命中 **0 行**（命令在 `logs/`，非只报枚数）。两包各 `ok`。⇒ 无一名需要归属。

---

## 2. ★明确结论：**"旧句那个『才』是半真话"——成立；但它换上的新句在默认形状上又造了一句过头话**

**(甲) "旧句那个『才』是半真话"：成立。** 成立在**这一条机制**上：`[fs] allowed_dirs` 是 **R2 的判级输入**（`risk/assessor.go:16`／`:79`、`rules_gateway.go:45-50`），**不是执行硬边界**——`FSDeps.open`（`fs.go:90-105`）执行腿只 `Canonicalize`、零 `InAllowlist`，我自己的 B 臂（批准门）真读回 20000 字节（C 臂对照 `L0`、零张卡）。旧句"fs.read 会被拒 … **才**读得回来"把"没人批"写成唯一情形 ⇒ **"半真"这个定性成立，且它给的那条机制理由逐字为真**。

**(乙) 新句"fs.read 会要一张 L2 卡"在默认形状上是过头话（越界）。** 三处我自己的读数：

1. **tier-A 目标根本不发卡**：`task.output` 的 allowlist 支只看 `InAllowlist`（`task.go:847`），**没有看 R3**。同桥同路径 `fs.read` 实测：`IsError=true level=L2 class="permission_denied"`、**approval 被问 0 次**、回执逐字 `该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）…`。此时 notice 仍逐字说 `会要一张 L2 卡 … 要么批下那一张卡` ⇒ **点了一张永不出示的卡，指了一条永不存在的路**（`risk/mode.go:203-208` 把 R3 列为不可静默红线；`assessor.go` 融合取最大，Deny 压 L2）。
2. **这不是构造出来的角落**：生产 artifacts 目录里完全可能出现 tier 命中——`risk.Classify` 现量：`…\wisp\artifacts\secrets.toml`→**B**、`…\artifacts\id_rsa`→**B**、`…\artifacts\credentials.json`→**B**（B 档＝仍需一张卡，不破句）；但**宿主登记一枚 A 档路径**（例如任务输出被回填成 `%APPDATA%\wisp\config.toml` 或 `~/.ssh/**` 下的文件）就落进 (1)。
3. **它自己那条腿的上游更早就知道这件事**：`internal/agent/spill.go` 的注释（`:218-222`）写着"a single approved card does fetch the bytes today"，**而生产上那条腿的 judge 恒为 nil**（`run.go:1087-1094` 的 `agent.Config` 字面量**没有 `PointerJudge`**，`grep -c PointerJudge run.go`＝**0`），我现跑：nil judge 下它说的是 `路径授权判定者未接线（fail-closed…）`——**它拿"另一条腿早就这么说"当依据，可那条腿在真机上根本走不到这句**。⇒ 论据的这一半是**引用不当**，不是读数错。

**(丙) 因此对实现腿那句话的定性（逐字指越界处）**：越界的是**这句**——`fs.read 会要一张 L2 卡，没人批就是直接拒`。它把"L2 卡"写成了**无条件机制**；而在 R3 A 档 / R8 不可逆等 `Deny` 形状下**没有卡**（我实测 approval=0），"没人批就是直接拒"里的"批"连发生的机会都没有。**该加而未加的条件词是"若这一发只被判到 L2"**（或等价限定）。**"两条路"本身没有越界**——在 R2-only 的形状上两条路都真（A/C 臂）。

---

## 3. 建议编排者翻哪几格（⛔ 本腿一枚未翻）

1. **`AC#2b` 的"两条路"半句：可翻。** 凭据＝第 1 节 (2) 逐字＋两腿运行时逐字节同形（246=246）＋第 1 节 (6) 的 `cmd/wisp` 端到端 `ok 11.027s`。
2. **但建议先开一枚新格（或并入 AC#2b 的条件）收 (乙)**：回执必须区分"只判到 L2（两条路都真）"与"被判 Deny（没有卡、批也无效）"。**这属措辞面，不属授权面**，可仍在 `internal/tools/**` 内做；另：它撞不撞 `task.go:299/:349` 两枚模板冻结钉＝`Q-63` 边界，仍须先有人裁一句。
3. **`AC#2d`／`AC#2` 不动**：本腿没有推翻旧判语，只**追加**了一条它没覆盖的形状。
4. **第 1 节 (4) 的 MUT-5 建议单列一枚小格**：#1 的标记集应加一枚只属第二路的标记（现句里 `要么批下那一张卡` 那一半无独立钉子）。
5. **第 1 节 (2) 结尾的幻影引用**（`task.go:852` 写 `…174rr1_test.go`）建议按 ban #9 同类处置：改注释，不留幻影。

---

## 4. 交付件与提交

- 本件：`docs/evidence/s1/174-v1-verdict.md`。
- 逐枚 `.go`：被告 `a8eb0d51`＝**2 枚**；越界（`internal/tools/**`＋`internal/agent/spill.go` 之外）＝**0**。本腿自己落的 `.go`＝**0 枚**（变异与探针只在 `D:\tmp\wisp\20261021\174v1\copy\`）。
- 正控：写面尺打 `5b14afe5` 读回 `OUT RANGE=1`；变异尺五发全部先证落地再读数、逐发单独还原并核 md5 回 HEAD blob（`432a52be…`／`0f754371…`／`cc22d191…`，`internal/agent/spill.go` `f31147b4…`）。
- 边界（具名未做）：`-race` 未跑；`cmd/wisp` 整包四数未跑（只跑 r8 那一发）；端到端真进程 `wisp run` 未跑；`gofumpt` 本机不在 PATH，未跑。
