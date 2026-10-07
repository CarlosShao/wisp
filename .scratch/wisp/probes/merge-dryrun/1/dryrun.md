# 合并预演（只读，⛔ 不落工作树／不动 index／不碰任何分支）

时刻 `2026-10-07 10:0x +08`；`git version 2.52.0`；dev 起手 `ca4ae2af`。
机主拍的口径＝**乙：先预演**。方法＝`git merge-tree --write-tree`（Git 的纯对象层合并，不 checkout、不写 index），随后所有数都打在那枚**结果树**上。

- 合并双方：`dev` × `dsh/feat/frontend-p0-v2`（工作树 `D:/wt/fe`，⛔ 我没进那棵树，一个字节都没读）
- 分岔点：`c1b10089`（2026-09-27 11:15）；那之后 fe 分支 **21 枚提交**、dev **1565 枚提交**
- 结果树：`21223d573af443cd2b5e7e6aa724a29b10a0e5b8`

## §1 真冲突枚数＝0（且不是"没检查"）
`git merge-tree --write-tree --name-only dev dsh/feat/frontend-p0-v2` → **rc=0，输出只有树号、没有任何冲突文件行**。
独立对拉一把（不依赖 merge 算法）：两边相对分岔点各自动过的文件清单取交集＝**0 枚**。
- fe 侧动过 **43** 枚文件，⛔ **全部在 `frontend/` 底下**（`grep -vc '^frontend/'`＝0）
- dev 侧动过 **4771** 枚文件
- 只在 fe 上存在（dev 没有的）＝**22** 枚，全是页面源件（`frontend/src/components/**.tsx`、`frontend/scripts/vendor-beautifului.mjs` 那一族；⛔ 我只列文件名，未开内容）

## §2 门分母不会被顶歪（两条各自有尺）
| 问 | 尺 | 读数 |
|---|---|---|
| frontend 作用域跟踪枚数会变多少 | `git ls-tree -r --name-only <树> \| grep -c '^frontend/'` | dev **85** → 合并树 **107**（＋22；其中 `.ts/.tsx` 66 → 86） |
| 有没有哪把门把 frontend 的枚数**写死成常数** | `grep -rn 'frontend' tools/d22scan/*.go \| grep -E '== *[0-9]+\|want'` | 那些常数断言（`want 5`、`+1`、`3 file(s)`）全在**合成台架**里（自己造 `frontend/.gitignore`、`frontend/dist/tracked-or-not.tsx` 那类假件），对真实仓的断言形状是"两把尺互相一致"（`TestRealRepoBan8CoversFrontendTreeAtBan6sCount`）⇒ ＋22 枚不会顶歪它 |
| 表情符号那道门（ban #8）会不会在这 22 枚上响 | `git grep -c -P '[1F300-1FAFF\|2200-22FF\|2600-27BF\|2B00-2BFF\|FE0F]' <树> -- 'frontend/'` | **0 行**（dev 现树也 0 行）。⚠ 这是**上界**（我这把尺不管"注释豁免、字符串不豁免"那条区分，注释也算进去了），上界是 0 ⇒ 门不会响 |
| 上面那把尺本身命中得了真名吗（正控） | 同一模式打 `HEAD:docs/reports/pending-and-issues.md`／`HEAD:AGENTS.md` | 台账 **2719 行**、`AGENTS.md` **6 行** ⇒ 尺有效，那个 0 是真读数（⛔ 不是"尺坏了所以 0"） |

## §3 ★并进来**不解决**"出货 exe 带不带页面字节"
- 合并树里 `frontend/dist/` 底下**仍只有 `.gitkeep`**（`git ls-tree -r --name-only <树> \| grep '^frontend/dist'`）⇒ **fe 分支那边也没把产物入库**（`.gitignore` 里 `dist/*` 那条规则是票 77 AC#1 定的，⛔ 不是谁的疏漏）。
- 而我自己在 `scripts/build.ps1:73-74` 现读到那句：
  `build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)`
  ⇒ 这句话的**前半已经过期**：`frontend/embed.go:19 //go:embed all:dist` 早就在树里了（`bundle-1` 报的，我按名核过作用域），**"embed 到 S5 才落"没落**；而后半"没有前端"更是反的——前端有，只是**构建这一步谁也不生成它**。
  ⇒ 后果（这条与合并无关，dev 侧就能修）：`go build` 时 `all:dist` 吃的是**编译机器上碰巧存在的那份 dist**（本机是 09-27 那份），换台干净机器就嵌进一只空目录。
- 我已实测过一枚反证，防止自己把上面说重：`build/wisp.exe panel-assets` → **`panel assets embedded: 4 files, entry=index.html built=true`** ⇒ 本机那枚 exe **确实带了 4 个文件**，所以"空目录"是**换机器时的风险**、不是今天已经发生的故障。

### §3 自纠（写完后 20 分钟内被我自己推翻的一句，⛔ 原句保留不抹）
本节第一段我当时写的是**"没有门看得见空 bundle"**——**这句是重的**。立票 274 前我按名字去核那把尺，现读到：
- `internal/panel/assets.go:34 errNotBuilt` ＋ `:54 newAssets` 里 `:57` 只有 `index.html` 在场才置 `built=true` ＋ `:75 Resolve` 在 `:76` 未构建时 `:77` **返回错误而不是假页面**；
- `cmd/wisp/panel_host_gate_test.go:77 TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12` 每次 CI 的 `cmd/wisp` 那一发都在跑（`ci.yml:573`），它 `:115-:127` 管"带页面"那一支、`:131-:137` 管"只有锚文件"那一支。
⇒ **真实缺口不是"没人管"，而是"那把尺两种形状都放行"**：它钉的是"降级必须 fail-closed"，不是"出货不许降级"。所以 CI 恒走 `built=false` 支、**整包照绿**。
⇒ 新增一把尺把因果钉死（本件起手没有这枚数）：`npm run build` 全 `ci.yml` **只有 1 处**＝`:903`（`lint-frontend`，ubuntu），而跑 `scripts/build.ps1` 造 exe 的是 `:571/:757/:820` 三处，**没有一处先跑 npm**；job 之间不共享工作目录 ⇒ **CI 造的每一枚 exe 都是未构建那一形**。
⇒ 另有一枚我上一段没量到的语义细节，已用仓外隔离台件量掉（⛔ 没动本机 dist）：`$HOME/tmp/embedprobe-274/` 只放 `dist/.gitkeep` ＋ 同样的 `//go:embed all:dist` ⇒ `go list -f '{{.EmbedFiles}}'` 回 **`[dist/.gitkeep]`**、`go build` **rc=0** ⇒ "能不能编译"与"里面有没有页面"今天**毫无关系**。
⇒ 票 274 的射程已按这节改写（原句"出货判据零票那一格"仍成立，但凭据从"没人看得见"换成"有把尺两形都放行"）。

## §4 我判不了、也不该由我判的那一格
**页面消费的字段跟 dev 这侧 Go 交出去的字段对不对得上**——量它必须打开 `frontend/src/**`，⛔ 这在我的边界里（机主 09-28 定的口径：前端那棵树不派、不写、不读、不转达）。两条合法出路，等他点：① 他给一句**一次性例外授权**，我只读"字段名清单"这一维、不看实现；② 由页面那边自己出一份对照表（他们本来就有那张表）。
⇒ ⛔ 我没有因为他选了"预演"就自己把这条读过去。

## §5 我的建议（等他拍，不是替他拍）
**现在不要并。** 理由不是技术（技术侧 0 冲突、门不歪、表情符号上界 0），而是**并了也不解决他真正关心的那件事**（干净机器双击能不能开出内容），而那件事的钥匙在 `build.ps1` 那一步与 CI，⛔ 不在分支合并里。⇒ 建议顺序：**先立并派"出货 exe 真带页面字节"那条票**（本预演 §3 的现量就是它的 AC#1 凭据），并完之后**再**决定分支——那时"并进来值不值"才有可测的判据（并完那份 exe 开出来是不是跟现在不一样）。
本件同时立成新票：`.scratch/wisp/issues/274-no-nail-requires-the-shipped-exe-to-carry-the-page-build-ps-step-2-comment-is-false-in-both-halves.md`（已入库，同一批 commit）。
⚠ 该票 §4 已具名写清**与票 34 的分工**（34 声称的 `docker/frontend.Dockerfile` 与 `assets/web` 在树里跟踪枚数都是 **0**，尺＝`git ls-files docker | grep -ci front`、`git ls-files | grep -c '^assets/web'`；34 的 AC#1 要真开一扇窗才能判、今天 `winlive` 未批 ⇒ 本票只补"不需要开窗就能自动判"的那半句，⛔ 不动 34、不并 34）——机主 10-07 那句"不要重复劳动"是按这一节落的。
