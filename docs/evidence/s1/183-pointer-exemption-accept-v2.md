# 183-v2 — 非实现者对抗验收（收窄版）：宿主指针豁免的三发变异恒真性 / 安全退让面进攻 / 三条边界 / AC#2 措辞判读

- 验收代理：`183-v2`（角色＝非实现者对抗验收；**只裁不改**：产码零改动，AC 框一枚不勾）
- 派单：`.scratch/wisp/dispatches/2026-09-28-132x-accept-183-v2-narrowed-incremental-commit-after-the-previous-leg-died.md`
- 被验收对象：`183-r1` 交件 `docs/evidence/s1/183-pointer-exemption-r1.md` ＋ 产码 commit `0662a35a`
- 本表是**增量交付**：第 1 节＝step-0 与起手现场；§2＝恒真性矩阵；§3＝安全退让面进攻；§4＝三条边界；§5＝AC 判读；§6＝本程没裁什么；§7＝现场与护栏自证；§8＝next。
- ⚠ 本节以下标 `待填` 的格，是"写表时还没跑"，**不是**"判定通过"。

## 1. step-0 五件（本程现量）

```
$ date "+%Y-%m-%d %H:%M:%S %z"
2026-09-28 13:23:30 +0800

$ git rev-parse --abbrev-ref HEAD ; git rev-parse --short HEAD
dev
7ec24d04

$ git merge-base --is-ancestor d61aee1b HEAD
（真＝派单起手锚 d61aee1b 是当前 HEAD 7ec24d04 的祖先；HEAD 多出的那一枚正是派单 183-v2 自己）

$ git status --porcelain -- internal/ cmd/
（空）

# 三向 md5（逐行管道到 md5sum，与本程起手一致）
$ sed -n '468,473p' internal/risk/provenance.go | md5sum
858e45116383caa3e7c1dd4b0924fad1 *-      ← 派单预期值，一致
$ sed -n '11,15p'   internal/risk/taintmatch.go | md5sum
5680ddd18e2d2ec2a85e485b54f4c12e *-      ← 一致
$ sed -n '482,506p' internal/risk/provenance.go | md5sum
f89e891e5eee3f3ea2b4f89d921072c4 *-      ← 一致

$ git show --numstat --format="%h" 0662a35a
```

`0662a35a` 的 numstat（本程现量，删除列**逐枚全 0**）：名下产码件 = `internal/risk/provenance.go` 48/0、`internal/risk/taintmatch.go` 59/0、`internal/risk/pointer_183_test.go` 331/0；其余为 `docs/evidence/s1/183-pointer-exemption-r1.md` 178/0、`probes/183/r1/**`（mut-f1/f2/freq/head 副本＋logs＋overlay＋zz183r1_e2e_test.go）、票 183 面 1/0。没有一枚落在 `frontend/**`／`design/**`。

判据件顶层 `func Test*` 枚数（本程现量，尺＝`grep -n "^func Test" internal/risk/pointer_183_test.go`）：**7 枚**，派单与被验收表写"顶层 6 枚" ⇒ 数目差一枚，§2 逐名列出这 7 枚。

## 2. 恒真性矩阵（六发自取变异 × 逐名判据）

尺（本程现跑，可复制；逐份读数留在 `.scratch/wisp/probes/183/v2/logs/mut-<tag>.txt`）：

```
bash .scratch/wisp/probes/183/v2/mut-matrix.sh HEAD M1 M2 M3 M4 M5 M0
# 台件＝同一脚本：python 现改现跑，每发之后 git cat-file blob HEAD:<path> > <path> 还原
# 判据跑法＝go test ./internal/risk/ -run 'TestPointer183' -v（单包，A359）
```

七枚常驻判据（`grep -n "^func Test" internal/risk/pointer_183_test.go` 现量＝**7 枚**；派单写"顶层 6 枚"＝派单侧数目漂了，被验收表 `183-r1` §11/§13 自报 7 枚，与盘上一致）：

| 号 | 判据全名（`pointer_183_test.go`） | M0 整支豁免摘掉 | M1 `contains()` 值跳过摘掉 | M2 `attachDeclaredPath` 那一支摘掉（留跨度） | M3 不要求 `runeIndexOf` 就挂 `declaredNorm` | M4 `spellsDeclaredPath` 容一枚 rune 差 | M5 声明过的 mark 整枚不算证据 |
|---|---|---|---|---|---|---|---|
| R1 | `TestPointer183RefereeTwinFragmentOutsideSpanStaysExempted` | **红** | **红** | **红** | 绿 | 绿 | 绿 |
| R2 | `TestPointer183RefereeControlBodyWithoutTwinIsClean` | **红** | 绿 | 绿 | 绿 | 绿 | 绿 |
| R3 | `TestPointer183PathAppearingTwiceInBodyStaysExempted` | **红** | **红** | **红** | 绿 | 绿 | 绿 |
| R4 | `TestPointer183BodySecretOutsideSpanStillHits` | 绿 | 绿 | 绿 | 绿 | 绿 | **红** |
| R5 | `TestPointer183ForeignMarkCarryingSamePathStillHits` | **红**（归因错人） | **红**（归因 `fs.read`） | **红**（同） | 绿 | 绿 | 绿 |
| R6 | `TestPointer183AlmostPathInSameHostMarkStillHits` | 绿 | 绿 | 绿 | 绿 | **红** | **红** |
| R7 | `TestPointer183WorstCaseOfTheLandedExemptionIsPinned` | **红** | **红** | **红** | 绿 | 绿 | **红** |
| — | 本发红掉的枚数 | 5 | 4 | 4 | **0** | 1 | 3 |

`HEAD` 基线（`logs/mut-HEAD.txt`）：七枚全绿、`ok …/internal/risk 0.120s`。变异终态：`git status --porcelain -- internal/ cmd/` 空，两支冻结文字与文档块三向 md5 在六发之后复量仍等于起手值（见 §7）。

**逐条判读**

1. **没有一枚是装饰**：七枚里每一枚都能被我六发中的至少一发打红（R2 只在"整支豁免被摘"时红＝它的职务就是位置性那一半的总开关探测器，与 r1 §2"归因控制腿"自述一致）。⇒ 这一格**判通过**。
2. **M1 与 M2 的红名册逐字相同**（R1/R3/R5/R7）⇒ 值规则是一枚行为，摘"复核侧那一跳"与摘"挂载那一跳"等价；这也说明 r1 §4 主张的"索引侧 map／复核侧逐字＝语义等价"**没有被判据区分开来**（等价性仍是未裁推理，见 §6 第 5 格）。
3. **M3 零枚红＝本程最要紧的一笔**：代码注释承诺的那支 fail-closed（`provenance.go:568-572`"声明的路径不在这段正文里 ⇒ 什么也不排除"）**没有任何常驻判据守着**——把"必须 `runeIndexOf` 命中才挂 `declaredNorm`"整条去掉，七枚照样全绿。⇒ **判不通过（判的是判据集，不是产码）**：那一支今天只靠一行注释活着。要补的腿形状很窄：`body` 不含 `hostPath` ⇒ `Inspect(fs.read,{path:hostPath})` **必须仍命中**（红在 M3、绿在 HEAD）。列入 §8 `next=`；我没写、也不该由我写。
4. **R5 在 M1/M2 下的红是"归因"红**（`hit must be attributed to the foreign mark, got "fs.read"`）而不是"没命中"红：宿主那枚 mark 也匹配同一条路径时，命中被记到宿主头上。⇒ 落地那族**顺带**修好了未修码上就存在的归因错位；A364 说"归因要当判据写进票 177 家族"——实测**已有一枚常驻腿在当它**（R5），不必新立。
5. M4/M5 证明防放宽的两枚边界腿各就各位：R6 拦"逐字比较放宽成近似／前缀同级"，R4 拦"值规则退化整枚豁免"；R6 在 M5 下也红 ⇒ 它不是只测一个方向。M4 只红一枚＝防"前缀放宽"那一族今天只有一枚腿在响，可接受但脆。

## 3. 安全退让面的进攻（自造样本）＋ 频率读数口径

台件（本程自写夹具，非复用 r1 夹具；跑法＝overlay 一枚虚拟件，跟踪件零改动）：

```
go test -count=1 -overlay=.scratch/wisp/probes/183/v2/overlay-attack.json -run 'TestV183Attack' -v ./internal/risk/
# rc=0，四枚全 PASS；逐字读数＝logs/attack.txt（我在落地码 HEAD 上那一发）
```

| 腿 | 造的形状 | 读数 | 判 |
|---|---|---|---|
| L1 | 同一枚宿主 mark，正文里除声明路径外再写**另一条资源**路径 `…\wisp\config.toml`（与声明串共享 `c:\users\swq\appdata\roaming\wisp\` ＝ >8 rune 同款前缀），参数去读那条 | `taint-hit … source=fs.read task.output fragment_len=8` ⇒ **命中** | 共享前缀**不够** laundering：候选参数里只要有一枚窗口拼不进声明串，证据还在 |
| L2 | 参数＝声明串的**连续子串**（`…\wisp\artifacts\` 目录段） | `hit=false (clean)` ⇒ **逃过 R4** | 逃得出的集合＝**整条候选参数都是声明串的逐字子串**；这是退让面的准确边界，比 r1 §7"少判一枚卡"可操作 |
| L3 | 宿主 mark 在 `task-1` 声明路径，外来 mark（`web.fetch`）在 `task-2` 携带同款前缀 | `task-1` 对自己声明的路径仍 clean；`task-2` 命中且归因 `web.fetch` | **跨 scope 借不到** ⇒ 没有名册行为（与 §4(b) 同证） |
| L4 | 票 185 那一形的包级版本：第一枚声明 mark ⇒ 首读 clean；`fs.read` 读回来的正文再被 `Mark`（**没声明**）盖戳 ⇒ 同一条路径第二发 | 第二发**仍命中**，`source=fs.read <那条产物路径>` | 影响面**只在同一枚 mark 内**；第二发被拒不是豁免死了，是另一枚没声明的 mark 在挡＝票 185 的事实，AC#3 禁止把它压安静 |

**进攻结论：我这发没结出可用的外带。** 要把内容从 R4 摘掉，攻击者必须让"想外带的那串"整体等于宿主自己写下的那条路径的子串；路径模板（Temp＋`artifacts`＋uuid＋`tool-output-…txt`）由宿主铸，可逃集合的实际价值＝"宿主已经公开给模型的那段字"。**真退让**在别处：同一枚 mark 正文里凡逐字复述了这条路径的 ≥8 rune（现实＝盘符＋`AppData\Roaming`＋`artifacts`）不再是续读证据＝少一枚卡（R7 钉语义，我 L2 钉集合）。

**频率读数复量与口径**（尺＝本程自己跑 `go test -count=1 -overlay=.scratch/wisp/probes/183/r1/overlay-freq-head.json -run TestFreq183 -v ./internal/risk/`，rc=0，读数 `logs/freq-recount.txt`）：

- `3 of 10 bodies make the host's own pointer unreadable` **逐字复现**；命中三枚＝`AGENTS.md`、`internal/agent/loop.go`、`synthetic:dir-listing`；响的窗口是同一组 profile 前缀（`\appdata\`／`\roaming`／`\swq\app`），不是夹具那枚 `-output-`。
- 分母 10 的组成（同一发读数自己打出来的）：**8 枚真仓内文件 + 2 枚作者自造合成正文**；另有 2 枚原拟语料（`docs/specs/SPEC-06-…`／`SPEC-08-…`）**因文件名取错读不到、没计入**。
- ⇒ **口径判定：数没错、话不够准**。拆开应是"真文件 2/8 响、自造样本 1/2 响"；`3 of 10` 把自造夹具算进了分母，而 §8 结论行是拿它当"真内容也会响"的凭据——那半句实际只由 2 枚真文件支撑（支撑得住，但要用 2/8 说）。
- r1 表 §8 末的 ⚠ 段**已具名披露**组成、丢样与"真实 LLM 输出语料发生率本程零读数"⇒ 不算瞒；要改的是**结论行与 ⚠ 段口径不一致**这一处文字（按 `A363` 规矩只许追加更正）。
- 再记性质一句：那 8 枚真语料是**仓内文件**，不是用户文档；"自用第一天就会撞上"是从"真实输出里会出现用户 AppData 路径"推出来的**推断**，不是读数，措辞上该标〔推断〕。


## 4. 三条边界

待填：(a) 参数侧按值放行有没有换张脸回来；(b) `declaredPath` 跨 mark 可见性；(c) `taintmatch.go:11-15` "逐 token 追踪已被否决"算不算被复活。

## 5. AC 判读（只给判读，不勾框）

待填：票 183 AC#2 档位；票 177 AC#3 条件格／票 175 AC#5／票 176 AC#3-5 可否翻；票 185 独立还是并进 183。

## 6. 本程没裁什么（逐枚具名＋为什么）

待填。

## 7. 现场与护栏自证

待填（含被拒/没成功的调用、有没有跑过删除命令、工具调用枚数 vs 硬顶 40、伪授权两栏）。

## 8. next=

待填。
