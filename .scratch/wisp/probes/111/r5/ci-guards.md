# 111 / r5 — CI 步守卫普查与补守卫（写码腿 `111-r5`，2026-10-07 09:06 +08 起手）

> 本件是这一枚写腿的台件。写面只有 `.github/workflows/ci.yml` 与本目录；
> **零 `.go` 读写**；只 commit 不 push。凡编排者转述与本腿读到的原文冲突，一律以原文为准并具名报回（见 §7）。

## §0 起手锚（现量，不是抄来的）

| 尺 | 读数 |
|---|---|
| `date` | `2026-10-07 09:06 +0800` |
| `git log --oneline -1` | `4dab3fbc probes(A651 落账)：两枚隔夜死腿都零提交⇒"第45次调用先commit"` |
| `git status --porcelain \| wc -l` | **753** 行（共树在飞的未跟踪件，本腿一枚都没 add） |
| 票面 `- [ ]` / `- [x]` 枚数 | **4 / 6**（编排者说的"十枚框、四枚未勾"复现；未勾＝`:27` AC#3、`:31` AC#5、`:75` AC#7、`:77` AC#8） |
| `.github/workflows/ci.yml` 起手状态 | 工作树与 HEAD **逐字节相同**：`git diff --numstat HEAD -- .github/workflows/ci.yml` 空；blob 与 `git hash-object` 同为 `d080b5b51b3fa47cce97b795629551ab27371266` |
| 起手期间 HEAD 是否移动 | **移动了**：本腿第一次读的是 `82a10da4`，第二次读到 `4dab3fbc`；`git log 82a10da4..HEAD -- .github/workflows/ci.yml` = **0 枚** ⇒ 那几枚提交没碰 `ci.yml`，下面 §1 的名册对本腿实际编辑的那份文件成立 |

⚠ 上面那行"起手行数"这一尺本腿换了写法：`git status --porcelain` 共树里有 753 行未跟踪/在飞件，
本腿只按**显式 pathspec** 取自己那几枚（`-- .github/workflows/ci.yml`、`-- scripts`），不引全仓计数当自己的足迹。

## §1 名册：`ci.yml` 每一步 ×（有没有 `if:` / 跑什么尺）

两份产出，都是**从文件逐行抽出来的**，不是人记的：

- `logs/step-roster-before.txt` — **39 行**（＝改前**有名字的**步枚数），抽法 `logs/roster.awk`（awk，纯文本，零 go 命令）。
  列义：`job|步序|- name: 所在行|守卫|类型|步名 + run/uses 摘要`。
- `logs/yaml-census-before.txt` — 同一份文件经**真 YAML 解析器**（`python` + PyYAML，尺＝`logs/yaml-guard-census.py`）的读数：
  **含未命名的 `uses:` 步**，所以是 **51 枚步**（39 有名字 + 12 只有 `uses:`）。逐作业：

| 作业 | 步数 | 带 `if:` | 带 `if:` 的是哪几枚 |
|---|---|---|---|
| `lint` | 13 | 3 | `staticcheck`、`mockllm module vet`、`Portable tests carrier self-test`（票 111 AC#3 正控载体，r5b 那笔 `1309757b` 加的） |
| `test-core` | 7 | 1 | `Stop compose services`（`if: always()`，清容器用的） |
| `test-windows` | 9 | 7 | 该作业**全部**具名步都带 `if: ${{ !cancelled() }}`（票 111 AC#6 那一笔） |
| `slo-smoke` | 6 | 0 | — |
| `slo-full` | 5 | 0 | — |
| `lint-frontend` | 11 | 0 | — |
| 合计 | **51** | **11** | 另：`continue-on-error` 在整份文件里出现 **0** 次（解析器数的，不是 grep 数的） |

改前那枚**本腿的靶子**（`lint` 作业，名册里第 6 枚具名步）：

```
lint|6|184|NOGUARD|sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only
     name: "gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)"
     前一步 = lint|5|168|NOGUARD|gofmt (gofumpt)   ← 票面 §305-306 记的正是这一步今天在红
```

同族对照（编排者给的 `:650` 那一枚；本腿逐行读到的形状＝**同一枚步的 `if:` 落在 `:650`、`run:` 落在 `:651`、
`- name:` 落在 `:564`**，编排者只给了守卫那一行，行号对得上）：
`test-windows` 的 `Package coverage census (ticket 111 AC#1 + GUARD D)` **带** `if: ${{ !cancelled() }}`
⇒ 同一把普查尺，windows 腿那枚有守卫、`lint` 腿这枚没有。

## §2 票面四格原文（逐字抽自 `.scratch/wisp/issues/111-ci-tests-20-of-33-packages.md`，全件在 `logs/ac-verbatim.txt`，8 行 / 1,161 字节）

`:27-28` **AC#3**

> - [ ] **AC#3** `session`/`watchdog` 这种"空分母"要**响亮**：scope 校验加上
>       "**声明在范围内但该平台没有任何测试文件 ⇒ 直接失败**"的守卫（与票 93 的"条目腐坏即红"同族），并人为抽掉一个包证明它会红。

`:31-32` **AC#5**

> - [ ] **AC#5** 门禁：`bash -n` 改动脚本 rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 不降；
>       ⚠ 新步若排在"会失败的步骤"之后 ⇒ 必须放前面或 `if: always()`（本仓实测过这道门因此从未执行）。

`:75-76` **AC#7**

> - [ ] **AC#7** `R-110-4`：票 110 承诺过"`-run TestSyncRegistryProbeLive` 经 `runtests.sh` 接进 windows 腿（step7）"，
>       验收复算发现**至今 0 次** ⇒ 要么落地并给 step7 读数，要么在票 110/111 面把它**当众改口径**（不许留在原地当已做）。

`:77-78` **AC#8**

> - [ ] **AC#8** `R-110-3`：包匹配式缺前缀锚定（`"winsec"` 宽松串曾让我把 0 命中说成 18 次）
>       ⇒ 匹配式要能区分"**被测包**"与"日志里出现过这个词"，并用一次阳性自证（种一个只在字符串里出现的包名 ⇒ 不许计入分母）。

⚠ 原文里 AC#5 那一行写的是 **`if: always()`**，而本仓**实际落地的形状**是 `if: ${{ !cancelled() }}`
（票 111 自己在 `ci.yml:438-465` 写明了为什么不取 `always()`；AC#6 的凭据读数也按 `!cancelled()` 记）。
本腿**照仓里已落地的形状**补守卫，不照派单字面把 `always()` 引进来，理由与出处见 §3 末。
