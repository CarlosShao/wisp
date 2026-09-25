**ADAPTED 2026-09-25**: props-driven ({text,streaming,sources?,followUps?}), 外链降级 span、光标归 reveal-text，P9；vendor.mjs 标 adapted | yes - Showcase 01 |**ADAPTED 2026-09-25**: props-driven four-state chips ({tools[]}), demo rows removed, P9；vendor.mjs 标 adapted | yes - Showcase 01（生产等票 35 字段） |**ADAPTED 2026-09-25**: props-driven ({seconds,steps,reasoning,sources}), sparkles head replaced by pixel grid (D23 禁图标), P9；vendor.mjs 标 adapted 重跑跳过 | yes - Showcase 01 |# Vendored code ledger (ticket 77 AC#7)

Every third-party file that enters `frontend/` is listed here, per file, with its
source repository, the component it is, its license, the exact upstream revision
or registry URL, and what this repo changed. Each such file also carries the same
facts in its own header so `git grep` on any file answers the provenance question
without opening this one.

Delivery model: all three sources below are copy-paste libraries (Q-20). **None of
them is an npm dependency.** `package.json` names only real runtime/build packages
(react, react-dom, tailwindcss, lucide-react, clsx, tailwind-merge,
class-variance-authority, radix-ui, vite, typescript, oxlint).

## How a file here got here

```
git clone --depth 1 https://github.com/TurboKach/ai-native-react-components.git /tmp/beautifului
cd frontend && node scripts/vendor.mjs /tmp/beautifului

curl -o /tmp/sh-button.json https://ui.shadcn.com/r/styles/new-york-v4/button.json   # card, badge likewise
cd frontend && node scripts/vendor-shadcn.mjs /tmp
```

Both scripts are in the repo, both are idempotent, and both refuse to run if the
transformation they claim to perform does not apply (a surviving `use client`
directive, or a dingbat with no ASCII mapping, exits 2). Re-running them is the
way to check this ledger against the tree; hand-editing a vendored file is what
this file exists to make visible.

## Beautiful UI - `TurboKach/ai-native-react-components` (MIT)

Upstream commit vendored: `05dab2d2b5f1f3e40029776e339a486d70491079`
License text: upstream `LICENSE`, "MIT License, Copyright (c) 2026 Turbo".
Owner ruling **R19**: this library is the base of the first version.

| File here | Upstream source | Component | Changes | Mounted? |
|---|---|---|---|---|
| `src/components/ai-native/stream-text.tsx` | `components/atoms/StreamText.tsx` | StreamText | header, dropped `use client`, LF endings | no - replaced on 2026-09-25, not rejected: it builds reveal units with `text.split(" ")`, and a Chinese answer has no spaces, so `PLAN.md:3476`'s 逐段 was unsatisfiable in the product's main language. Our own `src/components/reveal-text.tsx` now carries that job; the atom stays as the visual reference it was cut from, body bytes unchanged |
| `src/components/ai-native/shimmer.tsx` | `components/atoms/Shimmer.tsx` | Shimmer | header, dropped `use client`, LF endings | **yes** - `src/components/panel-skeleton.tsx` |
| `src/components/ai-native/approval-card.tsx` | `components/approval-card.tsx` | ApprovalCard | header, dropped `use client`, LF endings | no - upstream demo questionnaire; it is the visual reference `src/components/l2-approval-card.tsx` was cut from |
| `src/components/ai-native/loading-state.tsx` | `components/loading-state.tsx` | LoadingState | **ADAPTED 2026-09-25**: props-driven ({label}), demo timer/state removed, P9；vendor.mjs 标 adapted，重跑跳过（--force 覆盖） | yes - mounted by PanelSkeleton/Showcase |
| `src/components/ai-native/thinking.tsx` | `components/thinking.tsx` | Thinking | **ADAPTED 2026-09-25**: props-driven ({seconds, steps, reasoning, sources}), sparkles 头换 pixel 点阵（D23 禁图标），P9；vendor.mjs 标 adapted 重跑跳过 | yes - Showcase 01 |
| `src/components/ai-native/tool-chips.tsx` | `components/tool-chips.tsx` | ToolChips | **ADAPTED 2026-09-25**: props-driven 四态 chips ({tools[]})，演示行删除，P9；vendor.mjs 标 adapted | yes - Showcase 01（生产等票 35 字段） |
| `src/components/ai-native/task-rows.tsx` | `components/task-rows.tsx` | TaskRows | header, dropped `use client`, LF endings | no - history rows are ticket 35+ |
| `src/components/ai-native/streaming-text.tsx` | `components/streaming-text.tsx` | StreamingText | **ADAPTED 2026-09-25**: props-driven ({text, streaming, sources?, followUps?})，外链降级 span、光标归 reveal-text，P9；vendor.mjs 标 adapted | yes - Showcase 01 |

What was **not** taken from Beautiful UI: its `app/globals.css` `:root` block. Those
40-odd declarations are that library's own palette (`--surface: #fff`, `--ink:
#1f2124`, ...), and copying them would have put a second style truth source into the
tree, which C21 forbids. Instead `src/styles/theme.css` maps their *names* onto Wisp
tokens with `var(...)` only - `--surface: var(--bg-raised-color)` - so the vendored
components render in C21's palette. That mapping is asserted literal-free by
`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`.

The 7 unmounted files are kept in the tree on purpose (R19 makes this library the
base, and ticket 35+ mounts most of them), and `TestVendoredDemoComponentsAreNotMounted`
proves which ones are unreachable, so "unmounted" is a fact the gate maintains
rather than a claim in this table. Vite tree-shakes an unimported module out of the
bundle entirely, so no demo string reaches the binary.

## shadcn/ui (MIT) - base primitives, registry style `new-york-v4`

License: "MIT, Copyright (c) 2023 shadcn" (upstream `LICENSE` in
`github.com/shadcn-ui/ui`). Retrieved from the public registry on 2026-09-21; the
registry items carry no commit id, so the URL and date are recorded per file header.

| File here | Registry item | Changes |
|---|---|---|
| `src/components/ui/button.tsx` | `r/styles/new-york-v4/button.json` | header; `from "cn"` -> `from "@/lib/cn"`; `text-white` -> `text-on-solid` |
| `src/components/ui/card.tsx` | `r/styles/new-york-v4/card.json` | header; `from "cn"` -> `from "@/lib/cn"` |
| `src/components/ui/badge.tsx` | `r/styles/new-york-v4/badge.json` | header; `from "cn"` -> `from "@/lib/cn"`; `text-white` -> `text-on-solid` |
| `src/lib/cn.ts` | `r/styles/new-york-v4/utils.json` (the `cn` helper) | header only |

shadcn's colour vocabulary (`primary`, `muted`, `destructive`, `card`, `ring`, ...)
is aliased onto C21 in `src/styles/theme.css`, and its `rounded-md`/`rounded-xl`
utilities are re-pointed at C21's `--r-*` scale, so even a radius written by
upstream resolves to a Wisp token.

## React Bits - **phase two, zero code in this tree**

| Item | Value |
|---|---|
| Repository | `DavidHDev/react-bits` (reactbits.dev; 48,058 stars, created 2024-08-06, not a fork) - measured 2026-09-25 with `gh api repos/DavidHDev/react-bits`. **This row used to name `dillionverma/react-bits`, which is not a repository**: that exact API call returns 404 and the user has no react-bits project in its repo list. Nothing in `docs/` carried the wrong name - PLAN.md never names an upstream owner for React Bits - so the fix stays here. |
| License | **MIT + Commons Clause** - not a plain MIT |
| Licence text as measured | upstream `LICENSE.md`, blob sha `6425315416e94469f28d0223a09f7285b2f785ab`, 1303 bytes; header line "MIT + Commons Clause License Condition v1.0"; copyright line "Copyright (c) 2026 David Haz". Three sentences bind Wisp - the grant, the restriction and the notice condition - and they are quoted verbatim in `docs/reports/frontend-session-log.md` §53.1 rather than paraphrased here. |
| Availability of the 11, measured 2026-09-25 | Only `AnimatedList` has source in the upstream tree (`src/content/Components/AnimatedList/`). The other ten names are absent from `src/content/**` **and** from the jsrepo registry `public/r/registry.json` (832 items, matched on name, title and description); each appears in the repo only as a preview image under `public/assets/pro/components/`. So for those ten, phase two's first question is "is there an item to take, and is it the Pro tier", not "does the Commons Clause allow it". Per-component rows: log §53.2. |
| Status in this repo | **not vendored, one line of it is in the tree.** R19: the owner picked 12 animated components for a later phase and ordered phase one to ship on Beautiful UI alone |
| Gate before it may enter | **owner re-review of the Commons Clause terms is required before any React Bits code is vendored** (the clause restricts selling/Redistribution of the components themselves, which is a product-shape question, not a code question) |
| Phase-two constraints already agreed | at most one full-screen WebGL ambience layer on screen at a time (`Glass Flow` / `Aura Blob` / `Neural Float` / `Fog Sphere` are mutually exclusive), and it must be destroyed when the panel hides, not merely paused - WebView2 CPU counts against the D32 budget. `Agentic Ball` may only ever be a secondary indicator inside the panel; the desktop ball stays Win32/Direct2D (ticket 77's first hard constraint) |

Verification that the exclusion holds, runnable at any time:

```
git grep -in "react-bits\|reactbits" -- frontend/ | grep -v VENDORED.md
```

Expected: no hits. This row exists so React Bits cannot disappear from the books
merely because phase one does not use it.

## Not vendored, and why it is named here anyway

- `design/` (11 prototype screens + `design/assets/tokens.css`) is this repo's own
  work and, per ruling **R18**, a *reference* rather than a blueprint. Only
  `tokens.css` is load-bearing: it is the C21 source the frontend theme generates
  from (`scripts/gen-tokens.mjs`).
- `lucide-react` (ISC) is a genuine npm dependency, used for the D23 icon rule
  (SVG icons, never emoji). It is not copy-paste code, so it is not in this ledger.

## Beautiful UI - `slev12397/beautiful-ui` (MIT)

Upstream commit vendored: `44a274e598395ab61e7c96c26fda2758780253b7`
License text: upstream `LICENSE`, "MIT License, Copyright (c) 2026 Shane Levine".
A second Beautiful-UI checkout, distinct from the `TurboKach/ai-native-react-components`
table above; the two atoms below were cut from this one because owner's 2026-09-25
ruling ("组件都按照真组件库的来") asked for the library pieces by name.

| File here | Upstream source | Component | Changes | Mounted? |
|---|---|---|---|---|
| `src/components/ai-native/chip.tsx` | `components/atoms/Chip.tsx` | Chip | provenance header; upstream doc comment kept verbatim; nothing else (row appended 2026-09-25 - the file shipped earlier the same day with its header but without this ledger row) | **yes** - `src/components/config-screen.tsx` |
| `src/components/ai-native/status-pill.tsx` | `components/atoms/StatusPill.tsx` | StatusPill | provenance header; `from "@/lib/utils"` re-pointed to `from "@/lib/cn"` | **yes** - `src/components/panel-skeleton.tsx` (title-bar pending-approval pill) |
| `src/components/ai-native/progress-ring.tsx` | `components/atoms/ProgressRing.tsx` | ProgressRing | provenance header; nothing else - body verbatim | **yes** - `src/components/config-screen.tsx` (the alpha knob's live ring) + `src/components/showcase.tsx` (section 05 demos) |
| `src/components/ai-native/segmented-control.tsx` | `components/atoms/SegmentedControl.tsx` | SegmentedControl | provenance header; dropped `use client`; nothing else | **yes** - `src/components/showcase.tsx` (permission-tier + retention demos, sections 06/07) |
| `src/components/ai-native/text-row.tsx` | `components/atoms/TextRow.tsx` | TextRow | provenance header; nothing else - body verbatim | **yes** - `src/components/showcase.tsx` (privacy cards, section 07) |
| `src/components/ai-native/value-pill.tsx` | `components/atoms/ValuePill.tsx` | ValuePill | provenance header; nothing else - body verbatim (its ring values are color-mix() formulas over hue vars, not literals) | **yes** - `src/components/showcase.tsx` (sections 06/07/08) |
| `src/components/ai-native/button.tsx` | `components/atoms/Button.tsx` | Button (+ buttonVariants) | provenance header; dropped `use client`; `from "@/lib/utils"` re-pointed to `from "@/lib/cn"` and made type-only (verbatimModuleSyntax); **1 literal swap**: the `filledShadow` inset white-highlight shadow (upstream wrote it as a colour literal) replaced with the `shadow-btn` token utility - the only permitted rewrite, the rest is upstream verbatim | **yes** - `src/components/showcase.tsx` header (light/dark toggle); the panel's own buttons stay on the shadcn primitive |
| `src/components/ai-native/entity-chip.tsx` | `components/atoms/EntityChip.tsx` | Monogram + EntityChip | provenance header; **the sanctioned literal fix changes a signature**: upstream defaulted Monogram's `color` to a hex literal (and EntityChip's `color` was optional), so the default is deleted and BOTH `color` props are now required - callers must name a colour, which in this tree is always a `var()` token | **yes** - `src/components/showcase.tsx` (authorization-record actor marks, section 06) |
| `src/components/ai-native/switch.tsx` | `components/atoms/Switch.tsx` | Switch | provenance header; dropped `use client`; **1 literal swap**: the thumb's drop-shadow (upstream wrote it as a colour literal) replaced with the `shadow-btn` token utility; nothing else | **yes** - `src/components/showcase.tsx` (local-retention toggle, section 07) |

Every utility class StatusPill names (the `green`/`orange`/`red`/`accent`/`neutral`
tones and their dot colours) resolves through `src/styles/theme.css`'s
library-vocabulary alias block; nothing was missing at mount time and no alias had
to be added for it.

## Gate notes for this tree (AC#4)

- **ban #6 (the panel-side approval-decision identifier, `frontend/`)** - the very
  first hit this gate produced was **this file**: writing the forbidden
  identifier into the ledger to explain the rule tripped both the scanner and
  `TestFrontendNeverNamesAnApprovalDecision`. It is written around here instead,
  because the ban is on the identifier appearing in this tree, not on
  understanding it - `tools/d22scan/main.go` and PLAN.md carry the wording.
  `tools/d22scan` walks
  `frontend/` as text files and the planted-violation run named the file
  (`frontend/src/__ban6_probe.ts:1: [panel-approval]`). The panel tree never
  contains that identifier; the callback that reports what a user asked for is
  `requestApprovalResolution` and the bridge method is `panel.approval.request`,
  because an allow decision is made natively and only natively (D33/F2). The
  scanner's `frontend/` entry is still registered `absentOK`, so the drift guard
  exits 2 while `frontend/` exists - arming it (`live:true`) is a
  `tools/d22scan` change, which is not this ticket's to make.
- **ban #8 (zero emoji)** - its declared scope is `design/`, `internal/`, `cmd/`;
  `frontend/` is not in it. Choice recorded per the ticket: **enforce it ourselves**
  with `TestFrontendHasNoEmoji`, using the scanner's own character ranges, rather
  than register `frontend/` as a permitted-emoji tree. 24 files scanned, 0 hits;
  the two glyphs that upstream shipped in `tool-chips.tsx` were ASCII-ized during
  vendoring, not after.
  > **2026-09-25 追加更正（前端会话，随本轮那 6 处字符的改动一起落）：上面这段的三句已不成立，原句保留只为旧引用可核。**
  > ① **`frontend/` 现在在扫描射程里** - 票 96 把它补进了 `emojiScopes()`，且声明为 `everyFile: true`
  > （`tools/d22scan/main.go:491-498`、`:867-871`），所以"门没照到我们、我们自行执行"这个前提已经作废；
  > ② **"0 hits" 当时是假的** - 同一把尺在 `composer.tsx:170`、`fixtures/composer-states.html:2/5/8` 的
  > U+2264 与 `thinking.tsx:213`、`tool-chips.tsx:186` 的 U+2212 上报 6 行，`internal/panel` 那份
  > `TestFrontendHasNoEmoji` 报 0，因为它的字符类（`frontend_hygiene_test.go:66`）**自称逐字抄自扫描器却没有
  > `\x{2200}-\x{22FF}`** 那一段（票 141 只加宽了扫描器那一份）。**这一处副本归编排者修（台账 `U1`），不是本文件能自结的。**
  > ③ **"during vendoring, not after" 现在只对 2 枚** - 本轮这 6 处是**手工在 vendoring 之后**改的，
  > 且 `scripts/vendor.mjs` 的 `EMOJI_RE` 同样缺数学段 ⇒ **它既看不见 U+2212、也就无法在重新拉取时把它挡回来**：
  > 一次 `npm run vendor:beautifului` 会把 `thinking.tsx` / `tool-chips.tsx` 的 U+2212 恢复原状，
  > 而那句"unmapped glyph is fatal rather than silently shipped"对它不生效。
  > ⇒ 本轮之后**盘上确实 0 命中**，但**这一格的持久性依赖 `U1` 那批把三份副本一起加宽**，不依赖本文件的自述。
  > **同轮稍后更正（前端会话，选了 A 案）：上面第 ③ 条描述的洞已在本文件同批关掉。**
  > `scripts/vendor.mjs` 的 `EMOJI_RE` 已补 `\u{2200}-\u{22FF}`，`GLYPH_MAP` 已加 `U+2212 → "-"`。
  > 取数依据：改前那把旧尺在 `3b59512~1` 的 8 枚 vendored 文件上报 **0** 命中、补宽后报 **2**
  > （`thinking.tsx` 与 `tool-chips.tsx` 各 1 枚 U+2212）——**旧尺结构性看不见那两枚**，这正是"手工改完、
  > 下次 re-vendor 静默回退"的机制；在今天这棵树上两把尺都报 0，所以加宽**不产生任何新命中**，只买一个前置拦截。
  > 从此"将来 vendor 进来一个带数学符号的文件谁负责拦"有了答案：**`vendor.mjs` 自己**——无映射即 `exit 2`
  > 并指名文件与码位，而不是等 CI 事后红。
  > 一枚本轮才看清的事实：**`ai-native/*.tsx` 头部那段 provenance 注释是 `vendor.mjs` 生成的**
  > （`header(job, commit, cleaned.replaced)`），所以**手改那几行是徒劳的**——一次 re-vendor 会整段重写。
  > ⇒ 本轮把 `thinking.tsx:11` 改回生成器会产出的样子（`1 dingbat glyph(s) ASCII-ized`），
  > 说明性文字只写在这里，因为只有本文件是手工维护、不会被机器覆盖。
  > 已证：把补宽后的 `asciiize` 跑在改前的文件上，**正文的 `-` 与当前盘上一字不差**（可复现）。
  > 未证：`tool-chips.tsx` 头部那个 **3** 是**推断**值（上游 2 枚 U+2713 ＋ 1 枚 U+2212），因为本轮
  > **没有上游 checkout**（取它要联网、且会重写文件），模拟只能拿"已被 ASCII 化过的本地拷贝"当输入、
  > 只会少数不会多算。生成器下次真跑时若给出别的数，**以它为准**，那一行不归本文件主张。
