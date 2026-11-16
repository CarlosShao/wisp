# 票 77（并入 34）解冻后第一批四格读数 — AC#1 / AC#4 / AC#10 / AC#12

- 实现腿：`agent-89`（本波唯一 frontend 写腿）· 2026-11-16 18:5x
- 判据数字全部现量；**AC 框全部未勾** —— 勾选与裁决表按规矩走非实现者验收（裁决者≠实现者）。
- 读数件：`77-scaffold-batch1-d22scan.log`（d22scan 双门全文）、`77-scaffold-batch1-token-diff.txt`（token 名册差集逐条全文，510 行）。
- 派单明确不做：AC#3（L2 卡真数据）、AC#6（CI run id）、AC#8/9/11。
- 硬约束遵守：未动 thresholds/golden/allowlist/PLAN/specs/D43/design；未 push；未放宽任何断言；未提交任何 dist 字节；UI/新文件零 emoji（新文件在 ban#8 扫描分母内实测通过）。

## AC#1 — 构建 → go:embed → Go 侧数 dist 条目 → dist 删回 .gitkeep

| 步 | 命令 | 读数 |
|---|---|---|
| 构建 | `npm --prefix frontend run build`（node v24.21.0 @ `D:\work\server\node24`） | **rc=0**（Vite 8.3.0，2457 modules） |
| 类型 | `npm --prefix frontend run typecheck` | **rc=0** |
| 编入 | `go build -o build/wisp.exe ./cmd/wisp` | **rc=0** |
| 嵌入测试 | `go test -count=1 ./internal/panel/ -run 'Built\|Assets'` | **rc=0**，4 个匹配测试全 PASS；`TestBuiltinBundleIsCompleteOrAbsentNeverHalf` 自报 `entry 1044 bytes, 2 asset refs, 5 embedded files` |
| Go 侧枚举 | `build\wisp.exe panel-assets` | **rc=0**：`panel assets embedded: 5 files, entry=index.html built=true` |
| 逐文件 | `build\wisp.exe panel-assets -manifest` | `.gitkeep 0` / `assets/index-CLdPDwCM.js 500879` / `assets/index-Cw4ODNoI.css 65368` / `assets/showcase-1wn9e3sF.js 155086` / `index.html 1044`（各带 sha256 前缀） |
| 自洽 | `build\wisp.exe panel-assets -check` | **rc=0**：`entry=index.html built=true 2 asset refs resolve [./assets/index-CLdPDwCM.js ./assets/index-Cw4ODNoI.css]` |

- **dist 条目明细（Go 侧读数）**：index.html 1 枚（1044 B）＋ assets/\*.js 2 枚（index-CLdPDwCM.js 500879 B、showcase-1wn9e3sF.js 155086 B）＋ assets/\*.css 1 枚（index-Cw4ODNoI.css 65368 B）＋ 锚点 .gitkeep 1 枚 ＝ **5 files**。
- **跑完删回**：`dist/` 仅剩 `.gitkeep`（`Get-ChildItem` 读数），`git status --porcelain -- frontend/dist` 空 ⇒ **零 dist 字节入库**。
- 框不勾口径不变：真 WebView2 窗口是票 33 的地界；本轮交付到命令级（无 node 也可取字节），AC#1 保持未勾。

## AC#4 — ban #6 / ban #8 对真实 frontend/ 的逐作用域文件数 ＋ 扫描台账

- 运行方式：`sh scripts/d22scan.sh`，shell 门用 **`D:\work\soft\Git\bin\bash.exe`**（A1884）；双门 rc=0（阳性对照 `runtests.sh`：`PASS=46 FAIL=0 SKIP=0, === RUN=103`；真扫 `d22scan: clean`）。全文 log ＝ `77-scaffold-batch1-d22scan.log`。
- **逐作用域自报行（>0，全部现量）**：

| 作用域 | 自报读数 |
|---|---|
| bans #1-5 internal/ | 407 production Go files |
| bans #1-5 cmd/ | 90 production Go files |
| **ban #6 frontend/** | **126 text files** |
| ban #7 internal/tools/ | 51 production Go files |
| ban #10 internal/risk/ | 21 production Go files（comment symbol citations） |
| **ban #8 design/** | **39 text files** |
| **ban #8 frontend/** | **126 text files**（everyFile，与 ban #6 同分母） |
| **ban #8 internal/** | **1027 Go files**（comments and _test.go included） |
| **ban #8 cmd/** | **316 Go files**（comments and _test.go included） |

- **扫描台账（被扫 / 未被扫）**：
  - 被扫：`design/`、`frontend/`（ban #6 与 ban #8 双覆盖，分母 126）、`internal/`、`cmd/`。
  - 声明作用域之外、**未被 ban #6/#8 扫**：`tools/`（d22scan 自身树，靠其自检测试另行覆盖）、根 `scripts/`（.sh/.ps1）、`docs/`、`.github/`、`.scratch/`、仓库根散件（AGENTS.md 等）。
  - 分母过滤：`node_modules/`、`frontend/dist` 构建产物、`build/` 由 gitignore 过滤排除（票 96 把 frontend/ 声明为 everyFile 的理由见 `tools/d22scan/main.go:597-602`；gitignore 过滤见同文件 walkEmoji 的 git-ignored 分母段）。
  - 作用域声明出处：ban #6 走 `walkText(frontend)`（`tools/d22scan/main.go:294`，scope 登记 :440-447）；ban #8 全部四树在 `emojiScopes()`（`tools/d22scan/main.go:612-619`）。
- 本轮**未重跑 planted-red**：票 77d 于 09-21 两次演示（`frontend/src/__ban6_probe.ts:1` 被点名、rc=1，探针即删，票面 Progress log 14:5x/15:1x 两款在案）；派单本轮 AC#4 只点名"逐作用域文件数＋台账"。

## AC#10 — token parity：design/assets/tokens.css vs 前端 token 名册（脚本化 diff）

- **具名事实：`design/assets/tokens.css` 不在仓** —— 盘上不存在、git 未跟踪；删除于 `72b65c614`（该 commit 对该文件只有删除）。**以哪份为准**：
  - 第一代（被删原件）的权威留档 ＝ `internal/ball/testdata/tokens_c21.css`：出处头 pin 原件 sha256 `8ee27dc5…`；本次用 `git show 72b65c614^:design/assets/tokens.css` 与留档副本 diff 实测 **header-only（+23 行出处头、正文无一处差异）** ⇒ 留档正文逐字节＝原件。
  - 第四代（现役）真相源 ＝ `frontend/scripts/gen-tokens.mjs` 的 `TOKENS` 表，`src/styles/tokens.generated.css` 由其逐值产出（`docs/evidence/s1/c21-native-tokens.md` 头部明文"权威真相源"背书）。
- **脚本化 diff 工具（本次新落盘，可复跑）**：`node scripts/token-diff-77.mjs`（`frontend/scripts/token-diff-77.mjs`）。输出（全文 510 行 ＝ `77-scaffold-batch1-token-diff.txt`）：
  - A（一代）129 名 / B（四代现役）43 名 / 双边同名 8 / **仅 A＝121 条** / **仅 B＝35 条** / **同名异值 8 条**（--accent、--danger、--shadow-card、--success、--success-soft、--warn、--warn-line、--warn-soft，逐条值对照在 diff 文件里）/ **完全一致 0 条**。
- 独立第二道：`npm run tokens:check` **rc=0**（`gen-tokens: src/styles/tokens.generated.css matches the generator (97 keys, 43 tokens, light+dark)`）——生成件与生成器逐值一致，且 src/ 下所有 var() 引用可解析。
- **"能补则补"的裁定：不补，一条都不补**。理由：164 条差集（121+35+8）全部是**世代更替**（owner 2026-09-25 pivot：第四代收编 beautiful-ui 主题，一代磨砂玻璃配色整体退役；--accent 从 `#86c2b9/#3e837a` 换成 `#3d9aff/#0285ff` 是裁定本身），不是同代漂移——往现役主题补一代 token ＝ 推翻 owner 裁定。同代一致性的机器守门已在位：`tokens:check`（生成器↔生成件）＋ `TestC21DesignTokensFourWayAgree`（C21 表↔native↔第四代，22 对 mapped pairs）。⛔ 没有凭空声称一致：本节数字就是差集本集。

## AC#12 — bundle 预算 ＋ motion/Shiki 懒载

- **首屏 JS（改后终态）**：`assets/index-CLdPDwCM.js` **500879 B raw**；gz ＝ 构建器自报 **147.70 kB** ／ node zlib level-9 现量 **146273 B**（＝ 142.8 KiB）。预算 ≤400KB gz：**达标，余量约 257 KiB**；超 0。
- **懒载（本轮真做了，不是只报告）**：
  - 改前事实链：`main.tsx` 静态 `import { Showcase }`（原 :6）⇒ `showcase.tsx:43-45` 静态 import BlurText/CountUp/Magnet ⇒ 四枚 vendored 组件静态 `import … from "motion/react"`（`blur-text.tsx:29`、`animated-list.tsx:34`、`count-up.tsx:25`、`shiny-text.tsx:33`）⇒ motion 在首屏 chunk（基线实测：`index-BlqFlbZ2.js` 654214 B raw / 197.21 kB gz）。
  - 改后：`main.tsx:4` import `ShowcaseRoute`（`components/showcase-lazy.tsx`，`lazy()+Suspense`，仅 `?harness=2` 挂载）⇒ motion 全部移入独立块 `showcase-1wn9e3sF.js`（155086 B raw / **50777 B gz**，按需取）。首屏 gz **197.21 → 147.70 kB，减 49.5 kB**。
  - `animated-list.tsx:34` / `shiny-text.tsx:33` 同样静态 import motion，但当前树零引用（react-bits 二期备用），本就不进 bundle。
  - **Shiki：frontend/src 零引用（grep `shiki|Shiki` ＝ 0 命中）＋ `package.json` dependencies 无 shiki ⇒ 不在 bundle，无懒载问题**；二期引入时必须走动态 import，别静态进首屏。
- ⛔ 口径声明：**没有删任何产品代码凑预算**——省下的字节全部来自把 demo 画廊（`?harness=2` 才挂载的 Showcase）挪出首屏，即 AC#12 字面要求的懒载本身；产品路径的模块一棵未动。

## 门禁复读（本次全部现跑）

| 门 | 读数 |
|---|---|
| `npm run typecheck` | rc=0 |
| `npm run build` | rc=0 |
| `npm run lint`（oxlint） | rc=0（36 warnings / **0 errors**；新文件 `showcase-lazy.tsx`、`token-diff-77.mjs` 各 0 条，`main.tsx` 终态 0 条） |
| `npm run tokens:check` | rc=0 |
| `go build ./cmd/wisp` | rc=0 |
| `go test -count=1 ./internal/panel/`（全包） | ok rc=0（21.9s） |
| `sh scripts/d22scan.sh`（Git bash 双门） | rc=0 clean |
| `node scripts/token-diff-77.mjs` | rc=0（diff 报告完整产出） |

## 遗留与去向

- AC#1 真窗口半格 → 票 33；AC#3 → 票 145 乙段口径；AC#6 → 编排者 push 后取 `lint-frontend` 步级 run id；AC#8/9/11 未派。
- 裁决表与 AC 勾框 → 非实现者按本文件三件读数（含两个 log/txt）复核后落。
