# 票 77 AC#9 —— CI job 全绿覆盖 ＋ 种子违规必被抓（种子扫描器正/负控）

- 日期：2026-11-17（写腿，派单 `.scratch/wisp/probes/orch/20261117/txt/77-ac9-ci-seeds-dispatch.txt`）
- 范围：**只做 AC#9**，其余格未动；勾框留非实现者。
- node 实测：**v24.21.0 @ D:\work\server\node24**（CI 侧 `setup-node node-version: 24`，已钉）。

## 1. 交付物（写面逐条对齐派单）

| 件 | 内容 |
|---|---|
| `frontend/scripts/lint-seeds.mjs` | 种子扫描器（正/负控一体，**零新依赖**，仅 node:fs/path/url） |
| `frontend/package.json` | 仅加两条 script：`seeds` / `seeds:self-test` |
| `.github/workflows/ci.yml` | `lint-frontend` job 末尾 **ADDED STEP ONLY** 追加 2 步（无既有步骤移动/改写/加 `if:`/continue-on-error，D22 mode-6） |
| `docs/evidence/s1/77-ac9-ci-seeds.md` | 本件 |

## 2. CI job 与步骤清单

Job 名：**`lint-frontend`**（`.github/workflows/ci.yml`，ubuntu-latest，`working-directory: frontend`，`timeout-minutes: 10`，无 `if:false`/continue-on-error/path 过滤，node 钉 24）。步骤（既有 8 步 + 本次新增 2 步，加粗为新增）：

1. checkout
2. setup-node（`node-version: 24`，npm cache，lockfile 键）
3. `npm ci`
4. `npm run typecheck`（tsc -b）
5. `npm run lint`（oxlint）
6. `npm run tokens:check`
7. `npm run build`（vite build → dist）
8. `npm run render:l2 -- fixtures/l2-card-fs-delete.json`
9. `npm run render:composer` ＋ fixture 逐字节 `git diff --exit-code`
10. `npm run render:stream`
11. `npm run render:nav`
12. **`npm run seeds:self-test`**（新增：正控绿＋三族种子负控各红，不过即红）
13. **`npm run seeds`**（新增：真扫 `frontend/src`，hex/emoji/禁用导入）

AC#9 第 1 条（CI 覆盖 tsc/lint/build 全绿 + node 钉 24）由既有步骤 2–7 承载；本腿不改它，只把种子扫描器接成步骤 12–13。run id 归 AC#6（编排者 push 后取），本单不做、不 push。

## 3. 三族规则的口径（含具名排除项，非隐藏豁免）

- **hex**：src 下任何 `.ts/.tsx/.css` 出现 `#rgb/#rgba/#rrggbb/#rrggbbaa` 字面量即红。唯一具名排除＝`src/styles/tokens.generated.css`（C21 token 生成单一来源，另有 `tokens:check` 与 `TestC21DesignTokensFourWayAgree` 独立守）。与 Go 侧 `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` 同口径，但**不限可达性**（全 src 皆扫）。
- **emoji**：字符类**逐字照抄** `tools/d22scan/main.go` 的 `emojiRe`（`1F000–1FAFF / 2200–22FF / 2600–27BF / 2B00–2BFF / FE0F / 1F1E6–1F1FF`），**去注释豁免**——原始行直扫，注释与字符串同罪。与 Go 规格文字冲突时以「UI 零 emoji」为准（派单明令）。⚠ 扫描器自身在 `TestFrontendHasNoEmoji` 的分母内（它扫 src+scripts），故样本字符用 `String.fromCodePoint(0x2713)` 运行时构造，源文件零 emoji 区间字节。
- **禁用导入**：两种形状——(a) specifier 解析逃出 src 落进 `fixtures/`/`demo/` 树（渲染 fixture JSON 不许进 bundle）；(b) `src/fixtures/**`（harness 演示数据）被 harness 脸之外的文件 import。harness 脸＝具名作用域：`src/main.tsx`、`src/components/harness/**`、`src/components/showcase.tsx`——即规则落地时恰有的引用面（现树 5 处引用全部在脸内，生而绿）。
- 退出码：0=净；1=有 finding 或自检失灵；2=分母为空（扫不到东西不许读成绿）。

## 4. 现跑读数（本机 2026-11-17，node v24.21.0）

### 4.1 自检（`npm run seeds:self-test`，rc=0）

```
PASS (green) positive control (clean file must stay green) -> 0 findings
RED  (caught) negative control 1 - hardcoded hex colour must ring -> src/components/seed-bad-hex.tsx:2 [hex] hardcoded colour literal "#0285ff" - reference a C21 token (var(--...)) instead
RED  (caught) negative control 2 - UI emoji must ring, even inside a COMMENT (the d22scan comment exemption is deliberately absent) -> src/components/seed-bad-emoji.tsx:1 [emoji] UI emoji character present (d22scan emojiRe bands, no comment exemption): U+2713
RED  (caught) negative control 3 - forbidden import must ring (product-face file importing src/fixtures) -> src/components/seed-bad-import.tsx:1 [forbidden-import] forbidden-import: "src/components/seed-bad-import.tsx" imports the harness fixture "@/fixtures/harness" (src/fixtures/harness) but is not the harness face (src/main.tsx, src/components/harness/**, src/components/showcase.tsx)
self-test verdict: positive control green + 3 seeded violations red - the scanner has teeth
selftest rc=0
```

正控＝干净 tsx 样本（token 引色、`@/lib` 导入）→ 0 findings；负控 2 的种子字符放在**注释行**内被点红，即证明「注释豁免已去」。四枚样本与真树走**同一个 `scanFile` 函数**，无第二套检测逻辑。

### 4.2 真树扫描（`npm run seeds`，rc=0）

```
lint-seeds: scan frontend/src (ticket 77 AC#9)
denominators: emoji files=104 | hex files=103 (excluded by name: src/styles/tokens.generated.css) | import files=100
findings: 0
scan rc=0
```

生而绿已实证：现树 104 文件 0 命中（hex 唯一命中地在具名排除的生成 token 文件内）。

### 4.3 串行门禁

| 命令 | rc | 读数 |
|---|---|---|
| `npm --prefix frontend run typecheck` | **0** | tsc -b 干净 |
| `npm --prefix frontend run lint` | **0** | oxlint：Found **38 warnings and 0 errors**（新文件 **0 新增**——初版 2 条 `no-misleading-character-class`/`no-useless-escape` 已改为 `new RegExp` 构造与去转义消除；修后复跑自检/扫描读数不变） |
| `npm --prefix frontend run build` | **0** | vite 8.3.0，2719 modules，dist：index.html 1044 B / index-CrOLLgyp.js 703528 B / index-BfOthHF3.css 65931 B / showcase-BWhv63et.js 152393 B |

dist 产物**零字节入库**：build 后已删回锚点，`git status --porcelain -- frontend/dist`＝空，目录仅剩 `dist\.gitkeep`（0 B）。

## 5. 硬约束遵守清单

- ⛔ 未动 `docs/PLAN.md`／`docs/specs/**`／thresholds／golden／allowlist／D43／design 契约文字（本枚 diff 仅 §1 表列四件）。
- ⛔ 未放宽任何断言（Go 侧 hygiene/colour/emoji 测试零触碰；新扫描器比 Go 版覆盖面更宽——hex 不限可达性）。
- ⛔ 未 push（AC#6 的 run id 归编排者 push 后取，本单不做）。
- ⛔ 未提交 dist 字节（§4.3 末段读数）。
- ⛔ UI emoji 未进产品文案（新文件源码零 emoji 区间字符，样本字符运行时构造；复跑 `TestFrontendHasNoEmoji` 口径的分母内通过）。
- git：显式 pathspec 提交，无 add -A／amend／reset／rebase／stash。

## 6. 边界与诚实口径

- 种子负控样本是**内嵌字符串**（不落盘、不产生临时文件），由 `--self-test` 经同一检测函数驱动；真树扫描与自检共用规则函数，没有"仅测试路径"。
- 禁用导入的 (b) 形状以「harness 脸」具名作用域实现，因为现树 `main.tsx`/`components/harness/**`/`showcase.tsx` 本就合法引用 `src/fixtures/**`（harness 架构是票 77 在先交付）；把脸外的任何新引用点红。脸的成员清单写在扫描器与 ci.yml 注释里，改脸＝改扫描器＝可审。
- hex 家族对 9 位以上十六进制串另有兜底分支（`#[0-9a-fA-F]{9,}`）；3–8 位段用 `\b` 边界，`#feedback` 类 id 选择器不误报。
- CI 新增两步**尚未有 run id**（本单 ⛔ push）；AC#9 勾框仍留待「编排者 push 后 lint-frontend 步级结论」与其他格一起非实现者验收。
