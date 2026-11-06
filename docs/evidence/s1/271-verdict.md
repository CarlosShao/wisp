# 票 271 — 非实现者终裁表（271-verdict）

- 角色：**非实现者终裁腿**（实现者＝`271-v1`／commit `dcd5bc825`；本腿与其不同人，`SPEC-12 §4.3` 双角色）
- 终裁读数 HEAD：`dffa9fd64`（派单锚 `446831ea9` 之后再推进一枚 `dffa9fd64`，后者**只动 `docs/reports/pending-and-issues.md`** 8 行——锚漂为纯台账漂，代码面零漂；271 两枚被碰文件自 `dcd5bc825` 后零触碰，`git log dcd5bc825..HEAD -- <两文件>` 空读数）
- 时间：2026-11-06 13:0x＋08；环境：`PATH` 前置 `build`＋`third_party/sherpa-onnx`；`GOTMPDIR/TMPDIR/TEMP/TMP` 全 unset；未设 `CGO_ENABLED`（全程未置 0）
- 写面：仅本件。台账、票名、产码零触碰（终态 `git status --porcelain -- internal cmd` 与起手逐位同，见 §3）。

---

## 0. 本腿独立做了什么（不是转述实现者）

1. 现读票面 49 行＋实现者凭据件 61 行＋新测试文件 264 行＋被改的 221 腿文件 549 行＋`task.go:600-605`＋236r3 钉 `tasklist_deferred_236r3_teeth_test.go:125-179`。
2. `git show dcd5bc825` 逐行核：commit 恰 3 文件（凭据件＋221 腿注释/文案＋新测试），**判据 `if len(names) != 2` 一字未动**（diff 里谓词行零出现，只有 Errorf 文案与注释）。
3. 五道门全部本腿现跑（fresh，非引用实现者读数），读数见 §3。
4. AC#0 选形独立重推（§1），先按票面三形各自推一遍代价，再对实现者的表核差——结论同形（甲），但理由链条是本腿自己的，且本腿补了一条实现者没写的 nit（§4 第 4 条）。

---

## 1. AC#0 独立选形裁定＝**甲**（升格成明写的警报线）

本腿先不看实现者的表，按票面三形自己推：

**这把尺量的是什么维**：`len(names) != 2` 量的是"task 这一家**变大了**"。全仓普查（271-a1 §1）确认：数 task 家族枚数的尺**只有这一把**；236r3 的名集钉只按名字认 `task.list` 一枚负向名，**其他任何名字的误注册都从它手下绿着过去**。⇒ 枚数维是**独一份的传感器**，名集维已有别人在看。

**乙为什么不行（本腿自己的推法，不只抄"重复覆盖"）**：乙把唯一的"家族变大"传感器拆掉，换来的覆盖 236r3 钉**已经有**——净效果是传感器集合变弱、零新增。乙唯一的真收益是"合法第三枚那天零摩擦"，但那一天**本来就要走流程**（新工具要有自己的票、D34 行、台账登记）；多看一枚带工序说明书的红，代价是分钟级，而传感器维度一旦拆掉是**永久**的。还有一层 236r3 覆盖不了的：**"批准维"没有代码尺**（271-a1 §4）——甲的哨兵红恰好是这一维的**流程强制停**：家族变大的那一刻，逼人回台账确认"这枚票是真的"。这个强制停乙给不了，丙也给不了。

**丙为什么不行**：丙的唯一理由（A648）是"几枚红全是转述、没真跑过"。那个重开条件已经兑现——`271-v1` 把 M-2c 补成了 **resident 实测**（本腿 §3 的 fresh 套件读数里它每次跑都重证一遍）。测量已经存在之后还选丙＝拒绝收一份已经付过钱、过了全部门的文档，是负价值。

**甲为什么最有价值＋最合主流 harness**：①**警报必须被验证会响**——从不验证的警报器是装饰（dead tripwire 反模式）；甲的正控 `Test271B3UnrelatedThirdToolTripsOnlyTheCountTripwire` 每次跑套件都证明"哨兵真的会红＋红句可归属＋名集钉同发不误红"，这是三形里唯一把"假红那发要能复现"（AC#1）变常驻判据的。②**失败文案即文档**——下一个接手的人在红的那一刻、在红的现场读到工序（分清红型→回票 271→不许就地放宽），零检索成本。③判据一字不动 → 不触碰 `-done` 票的判定语义，AC#4 的具名解冻程序不触发。④本腿独立 steelman 过乙（"永不假红"确实值钱），结论不变：假红没有消失、而是被**红现场的说明书**转化成了强制停——噪声变控制，比消失更有价值。

**裁定：甲。** 与实现者同形；差异只在论证路径（本腿从"传感器唯一性＋批准维无尺"推入，实现者从"A648 重开条件兑现"推入，两链互补不重复）。

---

## 2. 逐格判定

| AC | 判定 | 一句话依据（读数全在 §3） |
|---|---|---|
| AC#0 选形 | **甲** | §1，本腿独立重推 |
| AC#1（选甲形：B-3 台件仍红但红得有说明书＋还原自证） | **PASS** | resident 正控本腿 fresh 套件内 PASS（6.77s）：overlay 台件下哨兵 FAIL 且红句含"task 家族注册了 3 枚"＋`task.note`（可归属）；同台件名集钉 PASS＋Logf 报出现名册含 task.note；还原件 PASS（盘上 task.go 字节==编译期字节；真名集钉直调复跑；起手 2 枚只 Logf：`task 家族名册 2 枚（task.output / task.cancel）`）。台件走 `-overlay`＋temp 目录，共享树零编辑，锚点命中数≠1 即 Fatal（拒绝盲切） |
| AC#2 不许顺手放宽 | **PASS** | `git show dcd5bc825` 谓词零 diff；Errorf 文案与注释是票面甲形**明文要求**的落点，非放宽；本腿五道门 fresh 全绿，绿不是改断言买来的；新测试自身不造第二把枚数红（起手枚数只 Logf 不断言） |
| AC#3 门禁＋还原自证 | **PASS** | §3 全表：build rc=0／整包 ok 30.124s（PASS=398 FAIL=0 SKIP=0，逐名红名册起手=终态=空，作差=0）／gofumpt·gofumpt 对两枚被碰文件零输出／vet rc=0／d22scan rc=0 clean／`git status --porcelain -- internal cmd` 终态=起手（3 枚未跟踪在途件原样） |
| AC#4（仅选乙触发） | **Not Applicable** | 选甲，解冻程序不触发。附带核：221 腿文件被改的只有注释块＋Errorf 文案，判定语义（谓词/Fatalf/其余断言）零触碰，且"补注释＋写清文案"正是票面甲形定义本身 |

**closeable：yes。** 六格无一附条件；实现者 §5.4 建议的台账登记动作（甲形记录＋撤销口令入 `pending-and-issues.md`）归编排者，不在本票判定射程内、不阻收口。

---

## 3. 本腿 fresh 门禁读数（2026-11-06 13:0x，HEAD `dffa9fd64`，全部本腿现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| 编译 | `go build ./cmd/... ./internal/... ./tools/...` | rc=0 |
| 整包 | `go test ./internal/tools -count=1 -timeout 420s` | `ok github.com/CarlosShao/wisp/internal/tools 30.124s`；`PASS=398 FAIL=0 SKIP=0`；逐名红名册＝空 |
| 271 正控 | 同上日志内 | `--- PASS: Test271B3UnrelatedThirdToolTripsOnlyTheCountTripwire (6.77s)`，Logf 逐字：`B-3 台件 resident 读数：枚数哨兵红（假红，红句具名）＋名集钉绿——两把尺的分工在同一发突变下实测成立`；overlay JSON 日志在案（Replace 键=包内 task.go 绝对路径→temp 台件） |
| 271 还原 | 同上日志内 | `--- PASS: Test271OverlayOnlyDisciplineRosterBackAtBaseline (0.00s)`；Logf 逐字：`起手读数：task 家族名册 2 枚（task.output / task.cancel）` |
| 哨兵/名集钉 | 同上日志内 | `Test221TaskCancelIsRegisteredAtItsFrozenLevel` PASS；`Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment` PASS（Logf：`现名册：task.output / task.cancel`） |
| 格式 | `gofumpt -l` ＋ `gofmt -l`（两枚被碰文件） | 两枚均零输出，rc=0 |
| 静态 | `go vet ./internal/tools/` | rc=0 |
| D22 | `sh scripts/d22scan.sh` | rc=0，尾行逐字：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=363, bans #1-5 cmd/=58, … ban #8 cmd/ 216 Go files…`（内含正控 runtests.sh 自己先绿） |
| 还原自证 | `git status --porcelain -- internal cmd`（起手与终态各读一次） | 两次逐位同：仅 `?? cmd/wisp/panel_resident_323_windows_test.go`、`?? cmd/wisp/panel_resident_assembly_323_test.go`、`?? internal/tools/ticket23_image_channel_test.go`（三枚他腿在途件，非 271 写面；零 modified） |

**诚实记录一发噪声**：本腿第一次 `sh scripts/d22scan.sh`（与整包测试并发跑）以 rc=1 收场，日志只有 runtests.sh 的 GOTMPDIR 声明行、正控与正式扫都没有输出——与第一次整包测试（10 分钟作业上限被杀于半途）合看，判定为**两个重进程并发抢构建缓存/输出捕获的机器噪声**，不是门红：树零改动、立即单独重跑即 rc=0 clean（上面表内读数）。两次读数都留档在 `.scratch/271-verdict-d22scan.log` 与 `.scratch/271-verdict-tools-test.log`（未跟踪探针件，按规矩只建不删）。

**环境差异声明**：实现者读数跑在 HEAD `ac2c1d903`（其起手），本腿跑在 `dffa9fd64`；`dcd5bc825..dffa9fd64` 之间 271 两文件零提交（`git log` 空读数），`abf2a35a6` 那枚纯检出编译缺口已由 `446831ea9` 补交 helper 修复（本腿裸 `go build` rc=0 为证），故两代读数可同尺比较。整包规模已较实现者时点增长（其 23.8s/35.8s → 本腿 30.124s，PASS 数量级一致），无定性差异。

---

## 4. 本腿补记（实现者没写的一条 nit，不阻收口）

新注释块（221 腿文件 `:170-181`）与其 Errorf 文案里的行号字面量（`:184`／`:183`）指的是**加注释之前**的谓词位置；插入 12 行注释后真身漂到 `:195`（`if len(names) != 2`）。字面量锚自己把自己推高了 11 行。定性：纯化妆性漂移，红现场的读者按文案全文仍能落到哨兵（同函数体内），**不值得单独出一枚 commit**；下一次任何腿碰这枚文件时顺手刷新即可。

---

## 5. 诚实边界（本腿没做到/没验的）

1. **哨兵在 CI 历史上是否真响过**：仍然零证据（271-a1 §4 ①），任何形都取不到；本腿验证的是"它现在会响、响得可归属、响了有人知道怎么办"。
2. **批准维无代码尺**：三形共有盲区，甲不例外；"第三枚的票是不是真的"仍归台账流程。甲的哨兵红是这一维的**强制停**，不是**判定**。
3. **内层 overlay 台件**本腿未在套件外手工重放——resident 正控本身就是每次跑套件的自动重放，本腿 fresh 读数即其中一次；未另造第二份手工读数。
4. **台账动作未做**：派单写面=仅本件；`pending-and-issues.md` 的甲形登记（含撤销口令）留给编排者。

**撤销口令**：`git revert dcd5bc825`（三文件原样回退，哨兵回到无说明书状态）＋删除本件＝终裁记录一并消失；重推同形只需重放 271-v1 的三文件。

---

*终裁腿 271-verdict · 非实现者 · commit 只带本件 pathspec、不 push。*
