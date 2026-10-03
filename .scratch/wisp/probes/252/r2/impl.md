# 票 252 · R2 腿（`252-r2`）—— 给"两段包含不许并成一段"这枚硬禁造仪器

## §0 起手锚（本腿自己现取，⛔ 不引用别人给的号）

- `date -Iseconds` = `2026-10-03T10:33:43+08:00`
- `git log -1 --format=%H` = `e16e303773f59be81d2f02eee3cc638870aa1630`（分支 `dev`）
- `git status --porcelain -- internal/tools internal/risk` = **0 行**（起手即为空 ⇒ 允许开工）
  - ⚠ 一次工具侧误读要记下：第一次跑这枚尺时 shell 的工作目录是 `.scratch/wisp/issues`，
    于是相对 pathspec 打空并吐出一行
    `warning: could not open directory '.scratch/wisp/issues/internal/': No such file or directory`（rc 仍为 0）。
    在仓库根 `D:/work/workspace/projects plans/Wisp` 重跑才是那枚 0 行的读数。**结论：pathspec 必须在仓库根求值。**
- 起手名册（`go test -count=1 -v ./internal/tools/`，落
  `.scratch/wisp/probes/252/r2/baseline-v-test.log`）：
  - 全仓规矩的那枚正则 `^(--- )?(PASS|FAIL|SKIP)` 命中 **187 行**＝186 枚顶层 + 1 行裸 `PASS`；
  - 含子测试的 `^[[:space:]]*--- (PASS|FAIL|SKIP)` 命中 **257 行＝257 PASS／0 FAIL／0 SKIP**；
  - 包体行 `ok github.com/CarlosShao/wisp/internal/tools 14.166s`；`go test -count=1 ./internal/risk/` = `ok 5.221s`。
  - ⇒ 编排者派单里"257"这一枚**复量为真**，但**它不是他那条正则能数出来的数**（那条数出 187）。见 §7 第 R2-1 条。

## §1 问题：那枚硬禁今天到底有没有尺

票面 AC#2 的硬禁原文（`252-...-r2-l2.md:30`）：⛔ **绝不许**把两次包含改成一次（那是把票 107 的洞重新打开）。
非实现者验收表 `docs/evidence/s1/252-allowlist-sameform-v1.md` 量出这条零仪器（M-merge 257 全绿），
本腿的第一件事是**自己复认**，第二件事是**造出能让它当场红的尺**。

现跑的码体（`internal/tools/paths.go`，锚点 `e16e3037`，本腿一个字都不改）：

| 段 | 行 | 逐字形状 | 谁喂它键 |
|---|---|---|---|
| 词法段（LEG 1） | `:196` | `if !rootsContain(p.roots, f) { return false }`，`f := foldPath(canonical)`（`:190`） | **只有** `canonical` 这一枚入参 |
| `resolvedForm` 段（LEG 2） | `:208-209` | `rf, ok := resolvedForm(canonical)` + `if !ok \|\| !rootsContain(p.roots, foldPath(rf)) { return false }` | `resolvedForm` 的返回值 |
| 工作区段（LEG 3，票 92） | `:215` | `if p.workspace != "" && !rootsContain([]string{p.workspace}, f)` | 也是 `f`，⛔ 但它的书不是 `p.roots` |

为什么 M-merge 量不出来（本腿的机制解释，不是抄来的）：修法 `2b1a3071` 把同形那一步挪进了
`Canonicalize`（`:145-148` → `sameFormOfUnresolved`），于是**生产能递给 `InAllowlist` 的每一枚串都已经两侧同形**
⇒ 对任何这类输入 `LEG1 == LEG2` ⇒ 包内那把 `measure252`（`paths_shortname_252_probe_test.go:114-126`）
算的恒等式 `final == leg1 && leg2` 里的 `leg1` 是**测试自己重算的**，不是从生产函数里读出来的，
生产函数少一条腿它照样绿。⇒ **凡是对"同形输入"取值的尺都天生看不见 M-merge**，这不是偶然失手，是射程不覆盖。

## §2 修法选择：三把尺＋为什么不是别的形

本腿只加仪器，写面 `internal/tools/**_test.go`（＋必要时 testdata）。生产码一字未动（终态 porcelain 见 §5）。

1. **W1 词法段的分歧见证（ behavioural，`//go:build windows`）**：
   取真 8.3 短名 + 一枚还不存在的叶子，**绕过 `Canonicalize` 直接**递给 `InAllowlist`。
   这一枚输入的形状是 `LEG1=false / LEG2=true` —— 全仓**只有**这种"两段不同答"的输入能让删掉词法段变红。
   断言 `InAllowlist=false`，并附一枚"同一条 raw ask 走完整生产回合必为 true"的对照，
   说明这一枚红不是"多拦了一次真授权"，而是"少了第二段之外还少的第一段"。
2. **W2 `resolvedForm` 段的分歧见证（portable，复用票 107b 的 `makeDirLink107b`／`dirLinkEvidence107b`）**：
   根真树、根内一枚 junction（POSIX 上是 symlink）通到根外。`LEG1=true / LEG2=false`，
   断言 `InAllowlist=false` —— 拆掉第二段当场红。
3. **S1 调用图尺（portable，`go/parser`＋`go/ast` 读 `paths.go`）**：
   只对 `InAllowlist` 的**函数体**取断言，两格：
   - 计数：体内对**授权本（`p.roots`）取用的包含调用**必须 **恰为 2 枚**，各带一条独立的拒绝支；
   - 来源：**两枚的键不许同源自同一枚表达式**——一枚必须只由入参折叠得来、另一枚必须经过一次再解析调用。
     这一格是行为尺天生量不到的那一种突变：把词法段的键悄悄换成 `foldPath(rf)`，
     对**所有**输入（含 W1、W2）取值都不变，两段却已经在实现上并成一段。
   - ⛔ 不给生产码加任何导出钩子：读文件用包内相对名，`go test` 的工作目录就是包目录。

选边理由（⛔ 不是"结构尺更好"的一句话）：**W1/W2 有牙但不完全，S1 完全但没有牙**。
删除任一段 → W1 或 W2 必红（这是票面要的"当场变红"）；把两段并成同源的"伪装合并" → 只有 S1 红。
两把互为补集，缺一把就留一个洞。编排者建议的"按调用图而不是按返回值"本腿**接受其一半**：
调用图尺保留（S1），但**不许只交它**，因为纯结构尺对"改成两段同源的两次包含"以外的情形
（例如把 `rootsContain` 内联成循环）会误响，而它自己宣称的"两次独立判定都在"最好有一枚取值证据兜着。

脆性自陈（命中面／下次合法重构会不会误响）写在 §3 末与 §6，读数一律现跑。

## §3 正控读数（骨架；读数与红句逐字在本腿第二次 commit 填入）

## §4 恒真自查（骨架；判据换反形的实测在本腿第二次 commit 填入）

## §5 门禁读数（骨架；四数与名册 comm 对比在本腿第二次 commit 填入）

## §6 判不动的地方（骨架；射程边界与脆性清单在本腿第二次 commit 填入）

## §7 推翻清单（骨架；对编排者派单逐句复跑的结果在本腿第二次 commit 填入）
